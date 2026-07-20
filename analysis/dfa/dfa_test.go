package dfa

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
)

func TestParseGoSyntax_RegisterOps(t *testing.T) {
	tests := []struct {
		goSyntax string
		opcode   string
		numOps   int
	}{
		{"MOVQ AX, BX", "MOVQ", 2},
		{"ADDQ BX, AX", "ADDQ", 2},
		{"CMPQ AX, BX", "CMPQ", 2},
		{"RET", "RET", 0},
		{"CALL fmt.Println(SB)", "CALL", 1},
		{"JEQ 0x12345", "JEQ", 1},
	}

	for _, tt := range tests {
		op, ops := ParseGoSyntaxOperands(tt.goSyntax)
		if op != tt.opcode {
			t.Errorf("%q: opcode = %q, want %q", tt.goSyntax, op, tt.opcode)
		}
		if len(ops) != tt.numOps {
			t.Errorf("%q: got %d operands, want %d: %v", tt.goSyntax, len(ops), tt.numOps, ops)
		}
	}
}

func TestParseOperand_Immediate(t *testing.T) {
	_, val := ParseOperand("$42")
	if val == nil || val.Kind != ValConst || val.Const != 42 {
		t.Errorf("$42 → %v", val)
	}

	_, val = ParseOperand("$0x1A")
	if val == nil || val.Kind != ValConst || val.Const != 0x1A {
		t.Errorf("$0x1A → %v", val)
	}
}

func TestParseOperand_Register(t *testing.T) {
	for _, reg := range []string{"AX", "BX", "CX", "DX", "SI", "DI", "R8", "R9", "R10", "R14", "SP", "BP"} {
		_, val := ParseOperand(reg)
		if val == nil || val.Kind != ValReg || val.Reg != reg {
			t.Errorf("%q → %v", reg, val)
		}
	}
}

func TestParseOperand_Memory_Simple(t *testing.T) {
	mem, val := ParseOperand("8(SP)")
	if val != nil {
		t.Errorf("8(SP) should not be a value: %v", val)
	}
	if mem == nil {
		t.Fatal("8(SP) should be a memory ref")
	}
	if mem.Base != "SP" {
		t.Errorf("base = %q, want SP", mem.Base)
	}
	if mem.Offset != 8 {
		t.Errorf("offset = %d, want 8", mem.Offset)
	}
}

func TestParseOperand_Memory_Bare(t *testing.T) {
	mem, val := ParseOperand("(AX)")
	if val != nil {
		t.Errorf("(AX) should not be a value")
	}
	if mem == nil || mem.Base != "AX" {
		t.Errorf("(AX) base = %v", mem)
	}
}

func TestParseOperand_Memory_Indexed(t *testing.T) {
	mem, _ := ParseOperand("8(SP)(BX*2)")
	if mem == nil {
		t.Fatal("should be a memory ref")
	}
	if mem.Base != "SP" {
		t.Errorf("base = %q, want SP", mem.Base)
	}
	if mem.Index != "BX" {
		t.Errorf("index = %q, want BX", mem.Index)
	}
	if mem.Scale != 2 {
		t.Errorf("scale = %d, want 2", mem.Scale)
	}
}

func TestParseOperand_Symbol(t *testing.T) {
	mem, _ := ParseOperand("fmt.Println(SB)")
	if mem == nil {
		t.Fatal("should be a memory ref")
	}
	if mem.Symbol != "fmt.Println" {
		t.Errorf("symbol = %q, want fmt.Println", mem.Symbol)
	}
	if mem.Base != "SB" {
		t.Errorf("base = %q, want SB", mem.Base)
	}
}

func TestBlockState_RegTracking(t *testing.T) {
	s := NewBlockState()

	s.SetReg("AX", ConstValue(42))
	v := s.GetReg("AX")
	if v.Kind != ValConst || v.Const != 42 {
		t.Errorf("AX = %v", v)
	}

	v = s.GetReg("BX")
	if v.Kind != ValReg || v.Reg != "BX" {
		t.Errorf("default BX = %v", v)
	}
}

func TestBlockState_StackSlot(t *testing.T) {
	s := NewBlockState()

	name := s.GetStackSlot(8)
	if name != "v0" {
		t.Errorf("first slot = %q", name)
	}

	name = s.GetStackSlot(8)
	if name != "v0" {
		t.Errorf("same slot should be v0, got %q", name)
	}

	name = s.GetStackSlot(16)
	if name != "v1" {
		t.Errorf("second slot = %q", name)
	}
}

func TestTranslate_MOV_Immediate(t *testing.T) {
	a := NewBlockAnalyzer(nil)
	inst := disasm.Instruction{GoSyntax: "MOVQ $42, AX"}
	a.translateInstruction(inst)

	v := a.State().GetReg("AX")
	if v.Kind != ValConst || v.Const != 42 {
		t.Errorf("AX should be const(42), got %v", v)
	}
}

func TestTranslate_MOV_RegToReg(t *testing.T) {
	a := NewBlockAnalyzer(nil)
	a.State().SetReg("BX", ConstValue(100))

	inst := disasm.Instruction{GoSyntax: "MOVQ BX, AX"}
	a.translateInstruction(inst)

	v := a.State().GetReg("AX")
	if v.Kind != ValConst || v.Const != 100 {
		t.Errorf("AX should be const(100) from BX, got %v", v)
	}
}

func TestTranslate_ADD(t *testing.T) {
	a := NewBlockAnalyzer(nil)
	a.State().SetReg("AX", ConstValue(10))
	a.State().SetReg("BX", ConstValue(5))

	inst := disasm.Instruction{GoSyntax: "ADDQ BX, AX"}
	a.translateInstruction(inst)

	v := a.State().GetReg("AX")
	if v.Kind != ValOp || v.Op != "+" {
		t.Errorf("AX should be binop +, got %v", v)
	}
}

func TestTranslate_LEA(t *testing.T) {
	a := NewBlockAnalyzer(nil)

	inst := disasm.Instruction{GoSyntax: "LEAQ 8(SP), AX"}
	a.translateInstruction(inst)

	v := a.State().GetReg("AX")
	if v.Kind != ValAddrOf {
		t.Errorf("AX should be addrof, got %v", v)
	}
}

func TestTranslate_CMP_and_Condition(t *testing.T) {
	a := NewBlockAnalyzer(nil)
	a.State().SetReg("AX", ConstValue(5))
	a.State().SetReg("BX", ConstValue(10))

	cmpInst := disasm.Instruction{GoSyntax: "CMPQ AX, BX"}
	a.translateInstruction(cmpInst)

	cond := a.buildCondition("JLT")
	if !strings.Contains(cond, "<") {
		t.Errorf("JLT condition should contain <: %q", cond)
	}
}

func TestTranslate_CALL(t *testing.T) {
	a := NewBlockAnalyzer(nil)

	inst := disasm.Instruction{
		GoSyntax: "CALL fmt.Println(SB)",
		IsCall:   true,
	}
	a.translateInstruction(inst)

	v := a.State().GetReg("AX")
	if v.Kind != ValCall || v.Func != "fmt.Println" {
		t.Errorf("AX should be call to fmt.Println, got %v", v)
	}
}

func TestOptimize_ConstantFolding(t *testing.T) {
	left := ConstValue(10)
	right := ConstValue(5)
	bin := BinOpValue("+", left, right)

	result := simplifyValue(bin)
	if result.Kind != ValConst || result.Const != 15 {
		t.Errorf("10+5 should be const(15), got %v", result)
	}
}

func TestOptimize_AddZero(t *testing.T) {
	left := ConstValue(0)
	right := &Value{Kind: ValReg, Reg: "ax"}
	bin := BinOpValue("+", left, right)

	result := simplifyValue(bin)
	if result.Kind != ValReg || result.Reg != "ax" {
		t.Errorf("0+ax should be ax, got %v", result)
	}
}

func TestEmit_Const(t *testing.T) {
	v := ConstValue(42)
	out := emitValueGo(v)
	if out != "42" {
		t.Errorf("const(42) → %q", out)
	}
}

func TestEmit_Register(t *testing.T) {
	v := RegValue("AX")
	out := emitValueGo(v)
	if out != "ax" {
		t.Errorf("reg(AX) → %q", out)
	}
}

func TestEmit_BinaryOp(t *testing.T) {
	left := &Value{Kind: ValReg, Reg: "ax"}
	right := ConstValue(10)
	v := BinOpValue("+", left, right)

	out := emitValueGo(v)
	if out != "(ax + 10)" {
		t.Errorf("ax+10 → %q", out)
	}
}

func TestExtractCallTarget(t *testing.T) {
	target := ExtractCallTarget("CALL fmt.Println(SB)")
	if target != "fmt.Println" {
		t.Errorf("CALL fmt.Println(SB) → %q", target)
	}

	target = ExtractCallTarget("CALL runtime.newobject(SB)")
	if target != "runtime.newobject" {
		t.Errorf("CALL runtime.newobject(SB) → %q", target)
	}
}

func TestParseImmediate(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"$42", 42},
		{"$0x1A", 26},
		{"$0x100", 256},
	}

	for _, tt := range tests {
		got := ParseImmediate(tt.input)
		if got != tt.expected {
			t.Errorf("ParseImmediate(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestShortFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fmt.Println", "Println"},
		{"runtime.newobject", "newobject"},
		{"main.greet", "greet"},
		{"encoding/json.Marshal", "Marshal"},
	}

	for _, tt := range tests {
		got := shortFunc(tt.input)
		if got != tt.expected {
			t.Errorf("shortFunc(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
