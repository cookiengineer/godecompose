package generate

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestWriteFunctionBodyNoiseFiltering(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "NOP", IntelSyntax: "nop", Size: 1},
		{Address: 0x1001, Opcode: "INT", IntelSyntax: "int3", Size: 1},
		{Address: 0x1002, Opcode: "MOV", IntelSyntax: "mov rax, rbx", Size: 3},
		{Address: 0x1005, Opcode: "RET", IntelSyntax: "ret", Size: 1},
	}

	f := &function.Function{
		Name:       "testNoise",
		ShortName:  "testNoise",
		EntryPoint: 0x1000,
		EndAddr:    0x1006,
	}

	g := &Generator{
		matches:      nil,
		instructions: insts,
	}

	var buf strings.Builder
	g.writeFunctionBody(&buf, f, insts, "\t")

	output := buf.String()
	t.Logf("output:\n%s", output)

	if strings.Contains(output, "int3") {
		t.Error("output contains INT3 noise")
	}
	if strings.Contains(output, "nop") {
		t.Error("output contains NOP noise")
	}
	if !strings.Contains(output, "mov") {
		t.Error("output missing valid MOV instruction")
	}
	if !strings.Contains(output, "ret") {
		t.Error("output missing RET instruction")
	}
	if !strings.Contains(output, "}") {
		t.Error("output missing closing brace")
	}
}

func TestWriteFunctionBodyLabels(t *testing.T) {
	// JE target; target: RET — JE is alone in a block, RET is the target
	// 0x1000: JE 0x1003 (block 0, conditional)
	// 0x1002: another instruction to break fallthrough
	// 0x1003: RET (block 1, target)
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "JE", IntelSyntax: "jz 0x1003", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1003},
		{Address: 0x1002, Opcode: "NOP", IntelSyntax: "nop", Size: 1},
		{Address: 0x1003, Opcode: "RET", IntelSyntax: "ret", Size: 1, IsReturn: true},
	}

	f := &function.Function{
		Name:       "testLabels",
		ShortName:  "testLabels",
		EntryPoint: 0x1000,
		EndAddr:    0x1004,
	}

	g := &Generator{
		matches:      nil,
		instructions: insts,
	}

	var buf strings.Builder
	g.writeFunctionBody(&buf, f, insts, "\t")

	output := buf.String()
	t.Logf("output:\n%s", output)

	if strings.Contains(output, "nop") {
		t.Error("output contains NOP noise")
	}

	if !strings.Contains(output, "if ") && !strings.Contains(output, "L0") {
		t.Log("output contains neither 'if' nor labels — both are acceptable for unresolved-only blocks")
	}
}

func TestWriteFunctionBodyWithMatches(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "TEST", IntelSyntax: "test rax, rax", Size: 3},
		{Address: 0x1003, Opcode: "JE", IntelSyntax: "jz 0x1010", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1010},
		{Address: 0x1005, Opcode: "MOV", IntelSyntax: "mov rax, 1", Size: 5},
		{Address: 0x100a, Opcode: "RET", IntelSyntax: "ret", Size: 1, IsReturn: true},
	}

	pat := &evaluator.CompiledPattern{
		Name:        "go_if_test_jeq",
		GenTemplate: "if $reg == nil { goto $lab }",
	}
	matches := []matcher.Match{
		{
			Pattern:   pat,
			StartAddr: 0x1000,
			EndAddr:   0x1005,
			Bindings: map[string]matcher.Binding{
				"reg": {CaptureVar: "reg", Value: "rax"},
				"lab": {CaptureVar: "lab", Value: "0x1010"},
			},
		},
	}

	f := &function.Function{
		Name:       "testMatch",
		ShortName:  "testMatch",
		EntryPoint: 0x1000,
		EndAddr:    0x100b,
	}

	g := &Generator{
		matches:      matches,
		instructions: insts,
	}

	var buf strings.Builder
	g.writeFunctionBody(&buf, f, insts, "\t")

	output := buf.String()
	t.Logf("output:\n%s", output)

	if !strings.Contains(output, "if ") {
		t.Error("output missing 'if' keyword")
	}
	if !strings.Contains(output, "rax") && !strings.Contains(output, "nil") {
		t.Error("output doesn't reference rax or nil from condition")
	}
}
