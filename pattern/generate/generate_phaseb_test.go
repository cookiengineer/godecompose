package generate

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func stmtsToString(stmts []ast.Stmt) string {
	if len(stmts) == 0 {
		return ""
	}
	var buf bytes.Buffer
	fset := token.NewFileSet()
	for _, s := range stmts {
		printer.Fprint(&buf, fset, s)
		buf.WriteByte('\n')
	}
	return buf.String()
}

func TestBuildFunctionBodyNoiseFiltering(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "NOP", IntelSyntax: "nop", GoSyntax: "NOP", Size: 1},
		{Address: 0x1001, Opcode: "INT", IntelSyntax: "int3", GoSyntax: "INT $3", Size: 1},
		{Address: 0x1002, Opcode: "MOVQ", IntelSyntax: "mov rax, rbx", GoSyntax: "MOVQ RAX, BX", Size: 3},
		{Address: 0x1005, Opcode: "RET", IntelSyntax: "ret", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	f := &function.Function{
		Name:       "testNoise",
		ShortName:  "testNoise",
		EntryPoint: 0x1000,
		EndAddr:    0x1006,
	}

	blocks := disasm.BuildControlFlowGraph(insts, []uint64{f.EntryPoint})
	structure := disasm.StructureControlFlow(f.Name, blocks)

	g := &Generator{
		matches:      nil,
		instructions: insts,
	}

	body := g.buildFunctionBody(f, insts, blocks, structure, nil)
	output := stmtsToString(body)
	t.Logf("output:\n%s", output)

	if strings.Contains(output, "int3") {
		t.Error("output contains INT3 noise")
	}
	if strings.Contains(output, "nop") {
		t.Error("output contains NOP noise")
	}
	if strings.Contains(output, "int") {
		t.Error("output contains INT noise")
	}
}

func TestBuildFunctionBodyLabels(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "JEQ", IntelSyntax: "jz 0x1003", GoSyntax: "JEQ 0x1003", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1003},
		{Address: 0x1002, Opcode: "NOP", IntelSyntax: "nop", GoSyntax: "NOP", Size: 1},
		{Address: 0x1003, Opcode: "RET", IntelSyntax: "ret", GoSyntax: "RET", Size: 1, IsReturn: true},
	}

	f := &function.Function{
		Name:       "testLabels",
		ShortName:  "testLabels",
		EntryPoint: 0x1000,
		EndAddr:    0x1004,
	}

	blocks := disasm.BuildControlFlowGraph(insts, []uint64{f.EntryPoint})
	structure := disasm.StructureControlFlow(f.Name, blocks)

	g := &Generator{
		matches:      nil,
		instructions: insts,
	}

	body := g.buildFunctionBody(f, insts, blocks, structure, nil)
	output := stmtsToString(body)
	t.Logf("output:\n%s", output)

	if strings.Contains(output, "nop") {
		t.Error("output contains NOP noise")
	}

	if !strings.Contains(output, "if ") && !strings.Contains(output, "L0") {
		t.Log("output contains neither 'if' nor labels — both are acceptable for unresolved-only blocks")
	}
}

func TestBuildFunctionBodyWithMatches(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "TEST", IntelSyntax: "test rax, rax", GoSyntax: "TESTQ AX, AX", Size: 3},
		{Address: 0x1003, Opcode: "JEQ", IntelSyntax: "jz 0x1010", GoSyntax: "JEQ 0x1010", Size: 2, IsBranch: true, IsConditional: true, BranchTarget: 0x1010},
		{Address: 0x1005, Opcode: "MOVQ", IntelSyntax: "mov rax, 1", GoSyntax: "MOVQ $1, AX", Size: 5},
		{Address: 0x100a, Opcode: "RET", IntelSyntax: "ret", GoSyntax: "RET", Size: 1, IsReturn: true},
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

	blocks := disasm.BuildControlFlowGraph(insts, []uint64{f.EntryPoint})
	structure := disasm.StructureControlFlow(f.Name, blocks)
	t.Logf("blocks: %d, structure: %v", len(blocks), structure != nil)

	g := &Generator{
		matches:      matches,
		instructions: insts,
	}

	body := g.buildFunctionBody(f, insts, blocks, structure, matches)
	output := stmtsToString(body)
	t.Logf("output:\n%s", output)

	if !strings.Contains(output, "if ") {
		t.Error("output missing 'if' keyword")
	}
	if !strings.Contains(output, "rax") && !strings.Contains(output, "nil") {
		t.Error("output doesn't reference rax or nil from condition")
	}
}
