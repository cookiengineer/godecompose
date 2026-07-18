// Package goasm provides Go-specific assembly analysis helpers for
// Plan 9 assembler output, including pseudo-register detection and
// ABI detection.
package goasm

import (
	"strings"

	"golang.org/x/arch/x86/x86asm"
)

// GoRegister maps x86asm register names to Plan 9 assembler register names.
var goRegisterNames = map[x86asm.Reg]string{
	x86asm.AL:   "AL",
	x86asm.CL:   "CL",
	x86asm.DL:   "DL",
	x86asm.BL:   "BL",
	x86asm.AH:   "AH",
	x86asm.CH:   "CH",
	x86asm.DH:   "DH",
	x86asm.BH:   "BH",
	x86asm.SPB:  "SP",
	x86asm.BPB:  "BP",
	x86asm.SIB:  "SI",
	x86asm.DIB:  "DI",
	x86asm.R8B:  "R8",
	x86asm.R9B:  "R9",
	x86asm.R10B: "R10",
	x86asm.R11B: "R11",
	x86asm.R12B: "R12",
	x86asm.R13B: "R13",
	x86asm.R14B: "R14",
	x86asm.R15B: "R15",

	x86asm.AX:   "AX",
	x86asm.CX:   "CX",
	x86asm.DX:   "DX",
	x86asm.BX:   "BX",
	x86asm.SP:   "SP",
	x86asm.BP:   "BP",
	x86asm.SI:   "SI",
	x86asm.DI:   "DI",
	x86asm.R8W:  "R8",
	x86asm.R9W:  "R9",
	x86asm.R10W: "R10",
	x86asm.R11W: "R11",
	x86asm.R12W: "R12",
	x86asm.R13W: "R13",
	x86asm.R14W: "R14",
	x86asm.R15W: "R15",

	x86asm.RAX: "AX",
	x86asm.RCX: "CX",
	x86asm.RDX: "DX",
	x86asm.RBX: "BX",
	x86asm.RSP: "SP",
	x86asm.RBP: "BP",
	x86asm.RSI: "SI",
	x86asm.RDI: "DI",
	x86asm.R8:  "R8",
	x86asm.R9:  "R9",
	x86asm.R10: "R10",
	x86asm.R11: "R11",
	x86asm.R12: "R12",
	x86asm.R13: "R13",
	x86asm.R14: "R14",
	x86asm.R15: "R15",

	x86asm.X0:  "X0",
	x86asm.X1:  "X1",
	x86asm.X2:  "X2",
	x86asm.X3:  "X3",
	x86asm.X4:  "X4",
	x86asm.X5:  "X5",
	x86asm.X6:  "X6",
	x86asm.X7:  "X7",
	x86asm.X8:  "X8",
	x86asm.X9:  "X9",
	x86asm.X10: "X10",
	x86asm.X11: "X11",
	x86asm.X12: "X12",
	x86asm.X13: "X13",
	x86asm.X14: "X14",
	x86asm.X15: "X15",
}

// GoRegisterName returns the Plan 9 assembler name for an x86asm register.
func GoRegisterName(reg x86asm.Reg) string {
	if name, ok := goRegisterNames[reg]; ok {
		return name
	}
	return reg.String()
}

// ABI represents the calling convention used by a function.
type ABI int

const (
	ABIUnknown    ABI = iota
	ABI0               // Stack-based (legacy Go ABI)
	ABIInternal        // Register-based (Go 1.17+)
	ABISystemV         // System V AMD64 ABI (C/C++)
)

func (a ABI) String() string {
	switch a {
	case ABI0:
		return "ABI0"
	case ABIInternal:
		return "ABIInternal"
	case ABISystemV:
		return "SystemV"
	default:
		return "unknown"
	}
}

// SpecialRegister indicates the role of a register in Go's ABI.
type SpecialRegister int

const (
	RegNormal     SpecialRegister = iota
	RegGoroutine                  // R14: current goroutine pointer (g)
	RegClosure                    // RDX: closure context pointer
	RegFramePointer               // BP: frame pointer
	RegStackPointer               // SP: stack pointer
	RegZeroValue                  // X15: always-zero register (ABIInternal)
	RegScratch                    // R12/R13: permanent scratch (ABIInternal)
	RegGOT                        // R15: GOT pointer (dynamic linking)
)

// ClassifyRegister identifies the special role of a register in Go's ABI.
func ClassifyRegister(reg x86asm.Reg) SpecialRegister {
	switch reg {
	case x86asm.R14, x86asm.R14B, x86asm.R14W:
		return RegGoroutine
	case x86asm.RDX, x86asm.DL, x86asm.DH, x86asm.DX:
		return RegClosure
	case x86asm.RBP, x86asm.BPB, x86asm.BP:
		return RegFramePointer
	case x86asm.RSP, x86asm.SPB, x86asm.SP:
		return RegStackPointer
	case x86asm.X15:
		return RegZeroValue
	case x86asm.R12, x86asm.R12W, x86asm.R12B,
		x86asm.R13, x86asm.R13W, x86asm.R13B:
		return RegScratch
	case x86asm.R15, x86asm.R15W, x86asm.R15B:
		return RegGOT
	default:
		return RegNormal
	}
}

// ABIArgumentRegisters returns the registers used for arguments in Go's
// register-based ABI (ABIInternal) in order.
func ABIArgumentRegisters() []x86asm.Reg {
	return []x86asm.Reg{
		x86asm.RAX, x86asm.RBX, x86asm.RCX,
		x86asm.RDI, x86asm.RSI,
		x86asm.R8, x86asm.R9, x86asm.R10, x86asm.R11,
	}
}

// ABIFloatArgumentRegisters returns the floating-point argument registers.
func ABIFloatArgumentRegisters() []x86asm.Reg {
	regs := make([]x86asm.Reg, 15)
	for i := 0; i < 15; i++ {
		regs[i] = x86asm.Reg(int(x86asm.X0) + i)
	}
	return regs
}

// IsPseudoRegister checks if a symbol name in Go assembly syntax refers to
// a pseudo-register (FP, SP, SB, PC).
func IsPseudoRegister(name string) bool {
	switch strings.ToUpper(name) {
	case "FP", "SP", "SB", "PC":
		return true
	}
	return false
}

// GoSpecialRegisters maps Go pseudo-register names to descriptions.
var GoSpecialRegisters = map[string]string{
	"FP": "Virtual frame pointer (function arguments: name+offset(FP))",
	"SP": "Virtual stack pointer (local variables: name-offset(SP))",
	"SB": "Static base pointer (global symbols: symbol<>(SB))",
	"PC": "Program counter (branch targets and labels)",
}

// DetectABI tries to determine which ABI a function uses based on its first
// few instructions. Returns ABIUnknown if it cannot determine.
func DetectABI(goSyntax []string) ABI {
	if len(goSyntax) == 0 {
		return ABIUnknown
	}

	for _, line := range goSyntax {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "morestack") || strings.Contains(trimmed, "CMP") && strings.Contains(trimmed, "R14") {
			return ABIInternal
		}

		if strings.Contains(trimmed, "arg") && strings.Contains(trimmed, "FP") {
			return ABI0
		}

		if strings.Contains(trimmed, "RDI") || strings.Contains(trimmed, "RSI") || strings.Contains(trimmed, "R8") || strings.Contains(trimmed, "R9") {
			return ABISystemV
		}
	}

	return ABIUnknown
}
