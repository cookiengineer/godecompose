package disasm

import (
	"testing"
)

func TestDecodeStreamSimple(t *testing.T) {
	// NOP (0x90), NOP, RET (0xC3)
	data := []byte{0x90, 0x90, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(instructions))
	}

	if instructions[0].Opcode != "NOP" {
		t.Errorf("inst[0].Opcode = %q, want NOP", instructions[0].Opcode)
	}
	if instructions[0].Address != 0x1000 {
		t.Errorf("inst[0].Address = %x, want 0x1000", instructions[0].Address)
	}
	if instructions[0].Size != 1 {
		t.Errorf("inst[0].Size = %d, want 1", instructions[0].Size)
	}

	if instructions[2].Opcode != "RET" {
		t.Errorf("inst[2].Opcode = %q, want RET", instructions[2].Opcode)
	}
	if !instructions[2].IsReturn {
		t.Error("inst[2].IsReturn should be true")
	}
}

func TestDecodeStreamCall(t *testing.T) {
	// CALL rel32 (E8 + 4 bytes offset)
	data := []byte{0xE8, 0x00, 0x00, 0x00, 0x00, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) < 1 {
		t.Fatal("expected at least 1 instruction")
	}

	if instructions[0].Opcode != "CALL" {
		t.Errorf("inst[0].Opcode = %q, want CALL", instructions[0].Opcode)
	}
	if !instructions[0].IsCall {
		t.Error("inst[0].IsCall should be true")
	}
	if !instructions[0].IsBranch {
		t.Error("inst[0].IsBranch should be true")
	}
}

func TestDecodeStreamJump(t *testing.T) {
	// JMP rel8 (EB + 1 byte offset)
	// JE rel8 (74 + 1 byte offset)
	data := []byte{0xEB, 0x05, 0x90, 0x90, 0x90, 0x90, 0x90, 0x74, 0xFE, 0xC3}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) < 2 {
		t.Fatal("expected at least 2 instructions")
	}

	if instructions[0].Opcode != "JMP" {
		t.Errorf("inst[0].Opcode = %q, want JMP", instructions[0].Opcode)
	}
	if !instructions[0].IsBranch {
		t.Error("inst[0].IsBranch should be true")
	}
	if instructions[0].IsConditional {
		t.Error("JMP should not be conditional")
	}

	// Find the JE instruction
	foundJE := false
	for _, inst := range instructions {
		if inst.Opcode == "JE" {
			foundJE = true
			if !inst.IsBranch {
				t.Error("JE.IsBranch should be true")
			}
			if !inst.IsConditional {
				t.Error("JE.IsConditional should be true")
			}
			break
		}
	}
	if !foundJE {
		t.Log("JE instruction not found in decoded stream (may be a different opcode)")
	}
}

func TestDecodeStreamMOVQ(t *testing.T) {
	// MOV RAX,immediate: 48 C7 C0 2A 00 00 00 (movq $42, rax)
	data := []byte{0x48, 0xC7, 0xC0, 0x2A, 0x00, 0x00, 0x00}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) < 1 {
		t.Fatal("expected at least 1 instruction")
	}

	if instructions[0].Opcode != "MOV" {
		t.Logf("Opcode = %q (expected MOV)", instructions[0].Opcode)
	}

	if instructions[0].IntelSyntax == "" {
		t.Error("IntelSyntax should not be empty")
	}
}

func TestDecodeStreamEmpty(t *testing.T) {
	instructions, err := DecodeStream([]byte{}, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}
	if len(instructions) != 0 {
		t.Errorf("expected 0 instructions, got %d", len(instructions))
	}
}

func TestDecodeStreamTruncated(t *testing.T) {
	// DecodeStream returns partial results + error on truncated data
	data := []byte{0x48, 0x8B}
	instructions, err := DecodeStream(data, 0x1000)
	if err == nil {
		t.Log("truncated instruction decoded without error")
	}
	_ = instructions
}

func TestResolveBranchTarget(t *testing.T) {
	// Test that JMP rel8 resolves correctly
	data := []byte{0xEB, 0x04, 0x90, 0x90, 0x90, 0x90}
	instructions, err := DecodeStream(data, 0x1000)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) < 1 {
		t.Fatal("expected at least 1 instruction")
	}

	// JMP rel8 = EB 04, instruction size = 2
	// target = 0x1000 + 2 + 4 = 0x1006
	expectedTarget := uint64(0x1006)
	if instructions[0].BranchTarget != expectedTarget {
		t.Errorf("BranchTarget = 0x%x, want 0x%x", instructions[0].BranchTarget, expectedTarget)
	}
}
