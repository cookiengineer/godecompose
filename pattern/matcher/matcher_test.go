package matcher

import (
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
)

func makePattern(opcodes []string, operands [][]evaluator.CompiledOperand) *evaluator.CompiledPattern {
	var alt []evaluator.CompiledInstruction
	for i, op := range opcodes {
		var compOps []evaluator.CompiledOperand
		if i < len(operands) {
			compOps = operands[i]
		}
		alt = append(alt, evaluator.CompiledInstruction{
			Opcode:   op,
			Operands: compOps,
		})
	}
	return &evaluator.CompiledPattern{
		Name: "test",
		Alternatives: [][]evaluator.CompiledInstruction{alt},
	}
}

func makeInst(opcode string, intel string, addr uint64) disasm.Instruction {
	return disasm.Instruction{
		Opcode:      opcode,
		IntelSyntax: intel,
		Address:     addr,
		Size:        3,
	}
}

func TestMatchExactOpcode(t *testing.T) {
	pat := makePattern(
		[]string{"MOV", "ADD", "RET"},
		nil,
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("ADD", "add rax, 1", 0x1003),
		makeInst("RET", "ret", 0x1006),
	}

	m := New([]*evaluator.CompiledPattern{pat})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].StartAddr != 0x1000 {
		t.Errorf("StartAddr = 0x%x, want 0x1000", matches[0].StartAddr)
	}
	if matches[0].EndAddr != 0x1009 {
		t.Errorf("EndAddr = 0x%x, want 0x1009", matches[0].EndAddr)
	}
}

func TestMatchConsecutive(t *testing.T) {
	pat := makePattern(
		[]string{"ADD", "RET"},
		nil,
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("ADD", "add rax, 1", 0x1003),
		makeInst("RET", "ret", 0x1006),
	}

	m := New([]*evaluator.CompiledPattern{pat})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match (ADD+RET), got %d", len(matches))
	}
	if matches[0].StartAddr != 0x1003 {
		t.Errorf("StartAddr = 0x%x, want 0x1003", matches[0].StartAddr)
	}
}

func TestMatchNoMatch(t *testing.T) {
	pat := makePattern(
		[]string{"CALL", "RET"},
		nil,
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("RET", "ret", 0x1003),
	}

	m := New([]*evaluator.CompiledPattern{pat})
	matches := m.Match(instructions)

	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestMatchOperandCapture(t *testing.T) {
	pat := makePattern(
		[]string{"MOV", "RET"},
		[][]evaluator.CompiledOperand{
			{{Register: "RAX"}, {CaptureVar: "src"}},
			nil,
		},
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("RET", "ret", 0x1003),
	}

	m := New([]*evaluator.CompiledPattern{pat})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Bindings["src"].Value != "rbx" {
		t.Errorf("captured src = %q, want rbx", matches[0].Bindings["src"].Value)
	}
}

func TestMatchWildcard(t *testing.T) {
	pat := makePattern(
		[]string{"MOV", "RET"},
		[][]evaluator.CompiledOperand{
			{{IsWildcard: true}, {IsWildcard: true}},
			nil,
		},
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("RET", "ret", 0x1003),
	}

	m := New([]*evaluator.CompiledPattern{pat})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	t.Logf("wildcard confidence: %f", matches[0].Confidence)
}

func TestMatchConfidence(t *testing.T) {
	exactPat := makePattern(
		[]string{"MOV", "ADD", "RET"},
		[][]evaluator.CompiledOperand{
			{{Register: "RAX"}, {Register: "RBX"}},
			{{Register: "RAX"}, {IsImmediate: true}},
			nil,
		},
	)

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("ADD", "add rax, 1", 0x1003),
		makeInst("RET", "ret", 0x1006),
	}

	m := New([]*evaluator.CompiledPattern{exactPat})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Confidence < 0.5 {
		t.Errorf("confidence too low for exact match: %f", matches[0].Confidence)
	}
}

func TestConflictResolution(t *testing.T) {
	patA := makePattern(
		[]string{"MOV", "ADD", "RET"},
		nil,
	)
	patA.Name = "A_longer"

	patB := makePattern(
		[]string{"MOV", "RET"},
		nil,
	)
	patB.Name = "B_shorter"

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rax, rbx", 0x1000),
		makeInst("ADD", "add rax, 1", 0x1003),
		makeInst("RET", "ret", 0x1006),
	}

	m := New([]*evaluator.CompiledPattern{patA, patB})
	matches := m.Match(instructions)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match after conflict resolution, got %d", len(matches))
	}
	t.Logf("selected: %s (confidence: %f)", matches[0].Pattern.Name, matches[0].Confidence)
}
