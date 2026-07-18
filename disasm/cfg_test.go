package disasm

import "testing"

func TestBuildCFGSimple(t *testing.T) {
	// Straight-line code: NOP; NOP; RET
	data := []byte{0x90, 0x90, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	blocks := BuildControlFlowGraph(instructions, nil)
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if len(block.Instructions) != 3 {
		t.Errorf("expected 3 instructions in block, got %d", len(block.Instructions))
	}
	if block.StartAddr != 0x1000 {
		t.Errorf("StartAddr = 0x%x, want 0x1000", block.StartAddr)
	}
	if len(block.Successors) != 0 {
		t.Errorf("expected 0 successors for straight-line RET, got %d", len(block.Successors))
	}
}

func TestBuildCFGJump(t *testing.T) {
	// JMP target; NOP; target: RET
	data := []byte{0xEB, 0x01, 0x90, 0xC3}
	// 0x1000: JMP 0x1003 (+2+1)
	// 0x1002: NOP
	// 0x1003: RET

	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	blocks := BuildControlFlowGraph(instructions, nil)
	if len(blocks) < 2 {
		t.Fatalf("expected at least 2 blocks, got %d", len(blocks))
	}

	// Find the JMP block and verify it has a successor targeting RET
	for _, block := range blocks {
		lastInst := block.Instructions[len(block.Instructions)-1]
		if lastInst.Opcode == "JMP" {
			if len(block.Successors) != 1 {
				t.Errorf("JMP block has %d successors, want 1", len(block.Successors))
			}
			if block.Successors[0].StartAddr != lastInst.BranchTarget {
				t.Errorf("JMP successor at 0x%x, want 0x%x",
					block.Successors[0].StartAddr, lastInst.BranchTarget)
			}
		}
	}
}

func TestBuildCFGConditional(t *testing.T) {
	// JE target; NOP; target: RET
	data := []byte{0x74, 0x01, 0x90, 0xC3}
	// 0x1000: JE 0x1003 (+2+1)
	// 0x1002: NOP (fallthrough)
	// 0x1003: RET

	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	blocks := BuildControlFlowGraph(instructions, nil)

	blockMap := make(map[uint64]*BasicBlock)
	for _, b := range blocks {
		blockMap[b.StartAddr] = b
	}

	jeBlock := blockMap[0x1000]
	if jeBlock == nil {
		t.Fatal("no block at 0x1000")
	}
	if len(jeBlock.Successors) < 1 {
		t.Errorf("JE block has %d successors, want at least 1 (taken)", len(jeBlock.Successors))
	}
}

func TestBuildCFGEmpty(t *testing.T) {
	blocks := BuildControlFlowGraph(nil, nil)
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks for nil input, got %d", len(blocks))
	}

	blocks = BuildControlFlowGraph([]Instruction{}, nil)
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks for empty input, got %d", len(blocks))
	}
}

func TestBuildCFGEntryPoints(t *testing.T) {
	data := []byte{0x90, 0x90, 0xC3, 0x90, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	// Treat 0x1003 as a function entry point
	blocks := BuildControlFlowGraph(instructions, []uint64{0x1003})
	if len(blocks) < 2 {
		t.Fatalf("expected at least 2 blocks with entry point, got %d", len(blocks))
	}

	found := false
	for _, b := range blocks {
		if b.StartAddr == 0x1003 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected block at entry point 0x1003")
	}
}
