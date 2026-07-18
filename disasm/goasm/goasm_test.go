package goasm

import (
	"testing"

	"golang.org/x/arch/x86/x86asm"
)

func TestGoRegisterName(t *testing.T) {
	tests := []struct {
		reg  x86asm.Reg
		want string
	}{
		{x86asm.RAX, "AX"},
		{x86asm.RBX, "BX"},
		{x86asm.RCX, "CX"},
		{x86asm.RDX, "DX"},
		{x86asm.RSI, "SI"},
		{x86asm.RDI, "DI"},
		{x86asm.RBP, "BP"},
		{x86asm.RSP, "SP"},
		{x86asm.R8, "R8"},
		{x86asm.R9, "R9"},
		{x86asm.R10, "R10"},
		{x86asm.R11, "R11"},
		{x86asm.R12, "R12"},
		{x86asm.R13, "R13"},
		{x86asm.R14, "R14"},
		{x86asm.R15, "R15"},
		{x86asm.X0, "X0"},
		{x86asm.X15, "X15"},
	}

	for _, tt := range tests {
		got := GoRegisterName(tt.reg)
		if got != tt.want {
			t.Errorf("GoRegisterName(%s) = %q, want %q", tt.reg, got, tt.want)
		}
	}
}

func TestClassifyRegister(t *testing.T) {
	tests := []struct {
		reg  x86asm.Reg
		role SpecialRegister
	}{
		{x86asm.R14, RegGoroutine},
		{x86asm.RDX, RegClosure},
		{x86asm.RBP, RegFramePointer},
		{x86asm.RSP, RegStackPointer},
		{x86asm.X15, RegZeroValue},
		{x86asm.R12, RegScratch},
		{x86asm.R13, RegScratch},
		{x86asm.R15, RegGOT},
		{x86asm.RAX, RegNormal},
		{x86asm.RBX, RegNormal},
		{x86asm.RCX, RegNormal},
	}

	for _, tt := range tests {
		got := ClassifyRegister(tt.reg)
		if got != tt.role {
			t.Errorf("ClassifyRegister(%s) = %v, want %v", tt.reg, got, tt.role)
		}
	}
}

func TestIsPseudoRegister(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"FP", true},
		{"SP", true},
		{"SB", true},
		{"PC", true},
		{"fp", true},
		{"RAX", false},
		{"AX", false},
		{"", false},
	}

	for _, tt := range tests {
		got := IsPseudoRegister(tt.name)
		if got != tt.want {
			t.Errorf("IsPseudoRegister(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestABIArgumentRegisters(t *testing.T) {
	regs := ABIArgumentRegisters()
	if len(regs) != 9 {
		t.Errorf("expected 9 integer argument registers, got %d", len(regs))
	}

	expected := map[x86asm.Reg]int{
		x86asm.RAX: 0, x86asm.RBX: 1, x86asm.RCX: 2,
		x86asm.RDI: 3, x86asm.RSI: 4,
		x86asm.R8: 5, x86asm.R9: 6, x86asm.R10: 7, x86asm.R11: 8,
	}

	for reg, idx := range expected {
		if regs[idx] != reg {
			t.Errorf("ABI argument register at index %d: %s, want %s", idx, regs[idx], reg)
		}
	}
}

func TestABIFloatArgumentRegisters(t *testing.T) {
	regs := ABIFloatArgumentRegisters()
	if len(regs) != 15 {
		t.Errorf("expected 15 float argument registers, got %d", len(regs))
	}
}

func TestDetectABI(t *testing.T) {
	tests := []struct {
		name   string
		code   []string
		expect ABI
	}{
		{
			"abi_internal_goroutine_check",
			[]string{"CMPQ R14, 16(SP)", "JBE morestack"},
			ABIInternal,
		},
		{
			"abi0_args",
			[]string{"MOVQ arg+0(FP), AX", "MOVQ arg2+8(FP), BX"},
			ABI0,
		},
		{
			"systemv_regs",
			[]string{"MOVQ RDI, AX", "CALL foo"},
			ABISystemV,
		},
		{
			"unknown",
			[]string{},
			ABIUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectABI(tt.code)
			if got != tt.expect {
				t.Errorf("DetectABI() = %v, want %v", got, tt.expect)
			}
		})
	}
}
