// Package disasm provides control flow structuring on top of the CFG.
// It identifies if/else branches, loops, and sequential blocks, enabling
// structured Go code generation instead of flat goto-based output.
package disasm

import "strings"

// BlockKind classifies a basic block by its role in control flow.
type BlockKind int

const (
	BlockPlain    BlockKind = iota // sequential block
	BlockIfThen                    // if-then body
	BlockIfElse                    // if-else body
	BlockLoopHead                  // loop header (back-edge target)
	BlockLoopBody                  // loop body
	BlockExit                      // function exit (return)
	BlockMerge                     // merge point after if/else
)

func (k BlockKind) String() string {
	switch k {
	case BlockPlain:
		return "plain"
	case BlockIfThen:
		return "if-then"
	case BlockIfElse:
		return "if-else"
	case BlockLoopHead:
		return "loop-head"
	case BlockLoopBody:
		return "loop-body"
	case BlockExit:
		return "exit"
	case BlockMerge:
		return "merge"
	}
	return "unknown"
}

// StructuredBlock wraps a basic block with structural information.
type StructuredBlock struct {
	Block    *BasicBlock
	Kind     BlockKind
	Children []*StructuredBlock
	ID       int
	dom      []int // dominators
	idom     int   // immediate dominator
	dfnum    int   // DFS number
	semi     int   // semi-dominator
	parent   int   // DFS parent
	bucket   []int
	ancestor int
	label    int
}

// StructuredFunc represents a function with structured control flow.
type StructuredFunc struct {
	Name        string
	EntryBlock  *BasicBlock
	Blocks      []*StructuredBlock
	Order       []*StructuredBlock
	DomTree     []int
	PdomTree    []int
	Loops       [][2]int
	Regions     []*CFGRegion
}

// RegionKind classifies a control flow region.
type RegionKind int

const (
	RegionSeq    RegionKind = iota
	RegionIfThen
	RegionIfElse
	RegionLoop
	RegionSwitch
	RegionJumpTable
)

func (k RegionKind) String() string {
	switch k {
	case RegionSeq:
		return "seq"
	case RegionIfThen:
		return "if-then"
	case RegionIfElse:
		return "if-else"
	case RegionLoop:
		return "loop"
	case RegionSwitch:
		return "switch"
	}
	return "unknown"
}

// CFGRegion represents a structured subgraph of the CFG.
type CFGRegion struct {
	Kind     RegionKind
	Header   int
	Merge    int
	Cond     *StructuredBlock
	Then     *CFGRegion
	Else     *CFGRegion
	Body     *CFGRegion
	Exit     int
	Blocks   []int
	Children []*CFGRegion
}

// StructureControlFlow analyzes a function's CFG and returns a structured
// representation suitable for code generation.
func StructureControlFlow(name string, blocks []*BasicBlock) *StructuredFunc {
	if len(blocks) == 0 {
		return nil
	}

	n := len(blocks)
	sb := make([]*StructuredBlock, n)
	addrToIdx := make(map[uint64]int)

	for i, b := range blocks {
		sb[i] = &StructuredBlock{
			Block:    b,
			Kind:     BlockPlain,
			ID:       i,
			idom:     -1,
			dfnum:    -1,
			semi:     i,
			ancestor: -1,
			label:    i,
		}
		addrToIdx[b.StartAddr] = i
	}

	computeDominators(sb, addrToIdx)
	pdom := computePostDominators(sb, addrToIdx)
	identifyLoops(sb, addrToIdx)
	classifyBlocks(sb, addrToIdx)

	var exitBlocks []int
	for _, b := range sb {
		if b.Kind == BlockExit {
			exitBlocks = append(exitBlocks, b.ID)
		}
	}

	return &StructuredFunc{
		Name:       name,
		EntryBlock: blocks[0],
		Blocks:     sb,
		DomTree:    getDomTree(sb),
		PdomTree:   pdom,
	}
}

func getDomTree(sb []*StructuredBlock) []int {
	tree := make([]int, len(sb))
	for i, b := range sb {
		tree[i] = b.idom
	}
	return tree
}

func computeDominators(sb []*StructuredBlock, addrToIdx map[uint64]int) {
	n := len(sb)
	dfnum := 0
	vertex := make([]int, n)

	var dfs func(v int)
	dfs = func(v int) {
		sb[v].dfnum = dfnum
		vertex[dfnum] = v
		sb[v].semi = v
		sb[v].label = v
		dfnum++

		for _, succ := range sb[v].Block.Successors {
			w, ok := addrToIdx[succ.StartAddr]
			if !ok {
				continue
			}
			if sb[w].dfnum < 0 {
				sb[w].parent = v
				dfs(w)
			}
		}
	}

	dfs(0)

	// Lengauer-Tarjan: compute semi-dominators
	for i := dfnum - 1; i > 0; i-- {
		w := vertex[i]
		for _, pred := range sb[w].Block.Predecessors {
			v, ok := addrToIdx[pred.StartAddr]
			if !ok {
				continue
			}
			if sb[v].dfnum >= 0 {
				u := evalLT(v, sb)
				if sb[u].semi >= 0 && sb[sb[u].semi].dfnum < sb[sb[w].semi].dfnum {
					sb[w].semi = sb[u].semi
				}
			}
		}
		sb[sb[w].semi].bucket = append(sb[sb[w].semi].bucket, w)
		linkLT(sb[w].parent, w, sb)
		for _, v := range sb[sb[w].parent].bucket {
			u := evalLT(v, sb)
			if sb[u].semi >= 0 && sb[sb[u].semi].dfnum < sb[sb[v].semi].dfnum {
				sb[v].idom = sb[u].semi
			} else {
				sb[v].idom = sb[w].parent
			}
		}
		sb[sb[w].parent].bucket = nil
	}

	for i := 1; i < dfnum; i++ {
		w := vertex[i]
		if sb[w].idom != sb[sb[w].semi].semi {
			sb[w].idom = sb[sb[w].idom].idom
		}
	}

	sb[0].idom = -1
}

func linkLT(v, w int, sb []*StructuredBlock) {
	sb[w].ancestor = v
}

func evalLT(v int, sb []*StructuredBlock) int {
	if sb[v].ancestor < 0 {
		return sb[v].label
	}
	compressLT(v, sb)
	if sb[sb[sb[sb[v].ancestor].label].semi].dfnum < sb[sb[sb[v].label].semi].dfnum {
		return sb[sb[v].ancestor].label
	}
	return sb[v].label
}

func compressLT(v int, sb []*StructuredBlock) {
	if sb[sb[v].ancestor].ancestor >= 0 {
		compressLT(sb[v].ancestor, sb)
		if sb[sb[sb[sb[v].ancestor].label].semi].dfnum < sb[sb[sb[v].label].semi].dfnum {
			sb[v].label = sb[sb[v].ancestor].label
		}
		sb[v].ancestor = sb[sb[v].ancestor].ancestor
	}
}

func identifyLoops(sb []*StructuredBlock, addrToIdx map[uint64]int) {
	if len(sb) < 1 {
		return
	}
	for _, b := range sb {
		for _, succ := range b.Block.Successors {
			si, ok := addrToIdx[succ.StartAddr]
			if !ok {
				continue
			}
			// Skip entry block as successor from other blocks (morestack fake loops)
			if si == 0 && b.ID != 0 && b.ID != si {
				continue
			}

			if dominates(si, b.ID, sb) {
				b.Kind = BlockLoopBody
				if sb[si].Kind != BlockLoopHead {
					sb[si].Kind = BlockLoopHead
				}
			}
		}
	}
}

func dominates(a, b int, sb []*StructuredBlock) bool {
	for b >= 0 && b != a {
		b = sb[b].idom
	}
	return b == a
}

func postDominates(a, b int, pdom []int) bool {
	for b >= 0 && b != a {
		b = pdom[b]
	}
	return b == a
}

func DetectJumpTable(block *BasicBlock) bool {
	if block == nil || len(block.Instructions) == 0 {
		return false
	}
	for _, inst := range block.Instructions {
		if inst.Opcode == "JMP" && strings.Contains(inst.GoSyntax, "(") && !strings.Contains(inst.GoSyntax, "(SB)") {
			return true
		}
		if inst.Opcode == "JMP" && strings.Contains(inst.Opcode, "JMP") && !inst.IsConditional && !inst.IsCall {
			syntax := inst.GoSyntax
			if strings.Contains(syntax, "(") && !strings.Contains(syntax, "(SB)") {
				return true
			}
		}
	}
	return false
}

func extractSingleCondition(cmpOp string, jccOp string) string {
	switch {
	case jccOp == "JEQ":
		return "=="
	case jccOp == "JNE":
		return "!="
	case jccOp == "JGT":
		return ">"
	case jccOp == "JLT":
		return "<"
	case jccOp == "JGE":
		return ">="
	case jccOp == "JLE":
		return "<="
	case jccOp == "JHI":
		return ">"
	case jccOp == "JLS":
		return "<="
	}
	return "_"
}

func classifyBlocks(sb []*StructuredBlock, addrToIdx map[uint64]int) {
	for _, b := range sb {
		hasCond := false
		for _, inst := range b.Block.Instructions {
			if inst.IsConditional {
				hasCond = true
				break
			}
		}
		last := b.Block.Instructions[len(b.Block.Instructions)-1]
		hasRet := last.IsReturn || last.Opcode == "RET" || last.Opcode == "RETQ"

		if hasRet && !hasCond && b.Kind == BlockPlain {
			b.Kind = BlockExit
			continue
		}

		if hasCond && len(b.Block.Successors) >= 1 {
			if b.Kind == BlockPlain {
				b.Kind = BlockIfThen
			}
		}

		if (last.Opcode == "JMP" || last.Opcode == "JMPL" || last.Opcode == "JMPQ") && b.ID > 0 {
			if b.Kind == BlockPlain && len(b.Block.Predecessors) >= 1 && len(b.Block.Successors) == 1 {
				b.Kind = BlockIfElse
			}
		}
	}
}

func computePostDominators(sb []*StructuredBlock, addrToIdx map[uint64]int) []int {
	n := len(sb)
	pdom := make([]int, n)
	for i := range pdom {
		pdom[i] = -1
	}

	exitIDs := make([]int, 0)
	for _, b := range sb {
		if b.Kind == BlockExit {
			exitIDs = append(exitIDs, b.ID)
		}
	}
	if len(exitIDs) == 0 {
		return pdom
	}

	revSuccs := make([][]int, n)
	revPreds := make([][]int, n)
	for i, b := range sb {
		for _, succ := range b.Block.Successors {
			s, ok := addrToIdx[succ.StartAddr]
			if !ok {
				continue
			}
			revSuccs[s] = append(revSuccs[s], i)
			revPreds[i] = append(revPreds[i], s)
		}
	}

	virtualExit := n
	revSuccsExt := append(revSuccs, []int{})
	revPredsExt := append(revPreds, []int{})
	for _, exitID := range exitIDs {
		revSuccsExt[virtualExit] = append(revSuccsExt[virtualExit], exitID)
		revPredsExt[exitID] = append(revPredsExt[exitID], virtualExit)
	}

	pdomBlocks := make([]*StructuredBlock, n+1)
	fakeAddr := uint64(0)
	for i := 0; i <= n; i++ {
		fakeAddr = uint64(i * 8)
		pdomBlocks[i] = &StructuredBlock{
			Block: &BasicBlock{StartAddr: fakeAddr},
			Kind:  BlockPlain,
			ID:    i,
			idom:  -1,
			dfnum: -1,
			semi:  i,
			ancestor: -1,
			label: i,
		}
	}

	fakeAddrToIdx := make(map[uint64]int)
	for i := range pdomBlocks {
		fakeAddrToIdx[pdomBlocks[i].Block.StartAddr] = i
	}

	for i := range pdomBlocks {
		for _, pred := range revSuccsExt[i] {
			pdomBlocks[i].Block.Predecessors = append(pdomBlocks[i].Block.Predecessors, pdomBlocks[pred].Block)
		}
		for _, succ := range revPredsExt[i] {
			pdomBlocks[i].Block.Successors = append(pdomBlocks[i].Block.Successors, pdomBlocks[succ].Block)
		}
	}

	computeDominators(pdomBlocks, fakeAddrToIdx)

	for i := 0; i < n; i++ {
		idom := pdomBlocks[i].idom
		pdom[i] = idom
	}

	return pdom
}
