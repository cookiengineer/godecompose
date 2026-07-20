package disasm

import (
	"testing"
)

func TestPostDominators_Simple(t *testing.T) {
	insts := []Instruction{
		{Address: 0x1000, Opcode: "MOVQ", GoSyntax: "MOVQ $1, AX", Size: 5},
		{Address: 0x1005, Opcode: "TESTQ", GoSyntax: "TESTQ AX, AX", Size: 3},
		{Address: 0x1008, Opcode: "JEQ", GoSyntax: "JEQ 0x1020", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1020},
		{Address: 0x100a, Opcode: "MOVQ", GoSyntax: "MOVQ $2, AX", Size: 5},
		{Address: 0x100f, Opcode: "JMP", GoSyntax: "JMP 0x1030", Size: 3, IsBranch: true, BranchTarget: 0x1030},
		{Address: 0x1020, Opcode: "MOVQ", GoSyntax: "MOVQ $3, AX", Size: 5},
		{Address: 0x1025, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
		{Address: 0x1030, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	blocks := BuildControlFlowGraph(insts, []uint64{0x1000})
	if len(blocks) == 0 {
		t.Fatal("no blocks built")
	}

	sf := StructureControlFlow("test", blocks)
	if sf == nil {
		t.Fatal("StructureControlFlow returned nil")
	}

	t.Logf("blocks: %d, domTree: %v", len(sf.Blocks), sf.DomTree)
	t.Logf("pdomTree: %v", sf.PdomTree)

	if len(sf.PdomTree) != len(blocks) {
		t.Errorf("pdomTree length %d != blocks %d", len(sf.PdomTree), len(blocks))
	}

	exitCount := 0
	for _, b := range sf.Blocks {
		if b.Kind == BlockExit {
			exitCount++
		}
	}
	if exitCount == 0 {
		t.Error("no exit blocks found")
	}
}

func TestPostDominators_IfElse(t *testing.T) {
	insts := []Instruction{
		{Address: 0x1000, Opcode: "TESTQ", GoSyntax: "TESTQ AX, AX", Size: 3},
		{Address: 0x1003, Opcode: "JEQ", GoSyntax: "JEQ 0x1010", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1010},
		{Address: 0x1005, Opcode: "MOVQ", GoSyntax: "MOVQ $1, AX", Size: 5},
		{Address: 0x100a, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
		{Address: 0x1010, Opcode: "MOVQ", GoSyntax: "MOVQ $2, AX", Size: 5},
		{Address: 0x1015, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	blocks := BuildControlFlowGraph(insts, []uint64{0x1000})
	sf := StructureControlFlow("ifelse", blocks)

	if sf == nil {
		t.Fatal("StructureControlFlow returned nil")
	}

	t.Logf("blocks: %d", len(sf.Blocks))
	for i, b := range sf.Blocks {
		t.Logf("  block %d: kind=%s start=0x%x succs=%d preds=%d idom=%d pdom=%d",
			i, b.Kind, b.Block.StartAddr, len(b.Block.Successors), len(b.Block.Predecessors), b.idom, sf.PdomTree[i])
	}

	if len(sf.Blocks) >= 2 {
		condBlock := sf.Blocks[0]
		if condBlock.Kind != BlockIfThen {
			t.Errorf("first block should be if-then, got %s", condBlock.Kind)
		}
	}
}

func TestDomTree_Root(t *testing.T) {
	insts := []Instruction{
		{Address: 0x1000, Opcode: "MOVQ", GoSyntax: "MOVQ $1, AX", Size: 5},
		{Address: 0x1005, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	blocks := BuildControlFlowGraph(insts, []uint64{0x1000})
	sf := StructureControlFlow("root", blocks)

	if sf == nil || len(sf.Blocks) == 0 {
		t.Fatal("no blocks")
	}

	if sf.DomTree[0] != -1 {
		t.Errorf("entry block should have idom=-1, got %d", sf.DomTree[0])
	}

	if sf.PdomTree[0] != 1 {
		t.Logf("entry pdom=%d (may be -1 if single block)", sf.PdomTree[0])
	}
}

func TestLoopDetection(t *testing.T) {
	insts := []Instruction{
		{Address: 0x1000, Opcode: "MOVQ", GoSyntax: "MOVQ $0, CX", Size: 5},
		{Address: 0x1005, Opcode: "CMPQ", GoSyntax: "CMPQ CX, AX", Size: 3},
		{Address: 0x1008, Opcode: "JGE", GoSyntax: "JGE 0x1020", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1020},
		{Address: 0x100a, Opcode: "ADDQ", GoSyntax: "ADDQ $1, CX", Size: 5},
		{Address: 0x100f, Opcode: "JMP", GoSyntax: "JMP 0x1005", Size: 3, IsBranch: true, BranchTarget: 0x1005},
		{Address: 0x1020, Opcode: "MOVQ", GoSyntax: "MOVQ CX, AX", Size: 3},
		{Address: 0x1023, Opcode: "RET", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	blocks := BuildControlFlowGraph(insts, []uint64{0x1000})
	sf := StructureControlFlow("loop", blocks)

	t.Logf("blocks: %d, loops: %v", len(sf.Blocks), sf.Loops)

	loopFound := false
	for _, b := range sf.Blocks {
		if b.Kind == BlockLoopHead || b.Kind == BlockLoopBody {
			loopFound = true
			t.Logf("loop block: kind=%s start=0x%x", b.Kind, b.Block.StartAddr)
		}
	}
	if !loopFound {
		t.Error("no loop blocks detected — back-edge (JMP 0x1005 → 0x1005 block) should be detected")
	}
}
