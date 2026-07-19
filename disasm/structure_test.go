package disasm

import (
	"testing"
)

func TestStructureControlFlowNil(t *testing.T) {
	sf := StructureControlFlow("test", nil)
	if sf != nil {
		t.Error("expected nil for empty blocks")
	}

	sf = StructureControlFlow("test", []*BasicBlock{})
	if sf != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestStructureSingleBlock(t *testing.T) {
	instructions := []Instruction{
		{Address: 0x1000, Opcode: "NOP", IntelSyntax: "nop", Size: 1},
		{Address: 0x1001, Opcode: "NOP", IntelSyntax: "nop", Size: 1},
		{Address: 0x1002, Opcode: "NOP", IntelSyntax: "nop", Size: 1},
	}
	blocks := BuildControlFlowGraph(instructions, []uint64{0x1000})
	sf := StructureControlFlow("test", blocks)

	if sf == nil {
		t.Fatal("expected non-nil for single block")
	}
	if len(sf.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(sf.Blocks))
	}
	if sf.Blocks[0].Kind != BlockPlain {
		t.Errorf("expected BlockPlain, got %s", sf.Blocks[0].Kind)
	}
}

func TestStructureIfElse(t *testing.T) {
	// JE target; NOP; target: RET
	// Block 0: JE -> taken or fallthrough
	// Block 1: NOP (fallthrough from JE)
	// Block 2: RET (target of JE)
	data := []byte{0x74, 0x01, 0x90, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}
	blocks := BuildControlFlowGraph(instructions, []uint64{0x1000})
	sf := StructureControlFlow("test", blocks)

	if sf == nil {
		t.Fatal("expected non-nil")
	}
	if len(sf.Blocks) < 2 {
		t.Errorf("expected at least 2 blocks, got %d", len(sf.Blocks))
	}

	// Block 0 (JE) should be if-then
	jeBlock := sf.Blocks[0]
	if jeBlock.Kind != BlockIfThen {
		t.Errorf("expected BlockIfThen for JE block, got %s", jeBlock.Kind)
	}
	if jeBlock.idom != -1 {
		t.Errorf("entry block should have idom=-1, got %d", jeBlock.idom)
	}
	t.Logf("JE block: kind=%s idom=%d successors=%d", jeBlock.Kind, jeBlock.idom, len(jeBlock.Block.Successors))
}

func TestStructureLoop(t *testing.T) {
	// JMP -2 (infinite loop)
	data := []byte{0xEB, 0xFE}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}
	blocks := BuildControlFlowGraph(instructions, []uint64{0x1000})
	sf := StructureControlFlow("test", blocks)

	if sf == nil {
		t.Fatal("expected non-nil")
	}

	// Block 0 (JMP to self) should be loop-head AND loop-body
	block := sf.Blocks[0]
	t.Logf("loop block: kind=%s idom=%d", block.Kind, block.idom)
	if block.Kind != BlockLoopHead && block.Kind != BlockLoopBody {
		t.Errorf("expected loop classification, got %s", block.Kind)
	}
}

func TestStructureReturns(t *testing.T) {
	data := []byte{0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}
	blocks := BuildControlFlowGraph(instructions, []uint64{0x1000})
	sf := StructureControlFlow("test", blocks)

	if sf == nil {
		t.Fatal("expected non-nil")
	}

	block := sf.Blocks[0]
	if block.Kind != BlockExit {
		t.Errorf("expected BlockExit for RET, got %s", block.Kind)
	}
}
