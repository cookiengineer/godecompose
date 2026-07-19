// Package disasm provides control flow structuring on top of the CFG.
// It identifies if/else branches, loops, and sequential blocks, enabling
// structured Go code generation instead of flat goto-based output.
package disasm

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
	Order       []*StructuredBlock // structured order for code gen
	DomTree     []int              // immediate dominator for each block
	Loops       [][2]int           // (header, latch) loop edges
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

	// Build dominator tree
	computeDominators(sb, addrToIdx)

	// Identify loops (back edges where target dominates source)
	identifyLoops(sb, addrToIdx)

	// Classify blocks
	classifyBlocks(sb, addrToIdx)

	return &StructuredFunc{
		Name:       name,
		EntryBlock: blocks[0],
		Blocks:     sb,
		DomTree:    getDomTree(sb),
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
	// a dominates b if a is on the path from root to b
	for b >= 0 && b != a {
		b = sb[b].idom
	}
	return b == a
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
