// Package syntax provides assembly instruction formatting for different
// assembler dialects.
package syntax

import (
	"fmt"

	"golang.org/x/arch/x86/x86asm"

	"github.com/cookiengineer/godecompose/disasm/goasm"
)

// Format returns a Go Plan 9 syntax string for the instruction.
// If symLookup is nil, no symbol resolution is performed.
func Format(inst x86asm.Inst, pc uint64, symLookup x86asm.SymLookup) string {
	return x86asm.GoSyntax(inst, pc, symLookup)
}

// FormatIntel returns Intel syntax for the instruction.
func FormatIntel(inst x86asm.Inst, pc uint64, symLookup x86asm.SymLookup) string {
	return x86asm.IntelSyntax(inst, pc, symLookup)
}

// FormatGNU returns GNU/AT&T syntax for the instruction.
func FormatGNU(inst x86asm.Inst, pc uint64, symLookup x86asm.SymLookup) string {
	return x86asm.GNUSyntax(inst, pc, symLookup)
}

// RegisterName returns the descriptive (Go Plan 9) name for an x86_64 register.
func RegisterName(reg x86asm.Reg) string {
	return goasm.GoRegisterName(reg)
}

// ConditionCode returns the Plan 9 condition code suffix for a Jcc/SETcc/CMOVcc
// opcode, or an empty string if not applicable.
func ConditionCode(op x86asm.Op) string {
	opStr := op.String()
	switch {
	case len(opStr) >= 2 && opStr[0] == 'J' && opStr[1] >= 'A' && opStr[1] <= 'Z':
		return opStr[1:]
	case len(opStr) >= 4 && opStr[:4] == "CMOV":
		return opStr[4:]
	case len(opStr) >= 3 && opStr[:3] == "SET":
		return opStr[3:]
	default:
		return ""
	}
}

// PseudoInstructionName maps x86asm opcodes to Go assembler pseudo-op names.
var PseudoInstructionName = map[x86asm.Op]string{
	x86asm.RET: "RET",
}

// ConditionCodeGo maps Plan 9 condition codes to their verbose names.
func ConditionCodeGo(cc string) string {
	names := map[string]string{
		"OS": "overflow set",
		"OC": "overflow clear",
		"CS": "carry set (unsigned below)",
		"CC": "carry clear (unsigned above or equal)",
		"EQ": "equal (zero)",
		"NE": "not equal (not zero)",
		"LS": "lower or same (unsigned below or equal)",
		"HI": "higher (unsigned above)",
		"MI": "minus (negative)",
		"PL": "plus (non-negative)",
		"PS": "parity set (even)",
		"PC": "parity clear (odd)",
		"LT": "less than (signed)",
		"GE": "greater or equal (signed)",
		"LE": "less or equal (signed)",
		"GT": "greater than (signed)",
	}
	if name, ok := names[cc]; ok {
		return name
	}
	return cc
}

// DisassembleRange decodes a range of bytes and returns formatted output
// with addresses. This is a convenience helper for debugging output.
func DisassembleRange(data []byte, baseAddr uint64, mode int) ([]string, error) {
	var lines []string
	offset := 0

	for offset < len(data) {
		inst, err := x86asm.Decode(data[offset:], mode)
		if err != nil {
			lines = append(lines, fmt.Sprintf("%016x: ??? (%v)", baseAddr+uint64(offset), err))
			offset++
			continue
		}

		goSyntax := x86asm.GoSyntax(inst, baseAddr+uint64(offset), nil)
		intelSyntax := x86asm.IntelSyntax(inst, baseAddr+uint64(offset), nil)

		line := fmt.Sprintf("%016x: %-40s ; %s",
			baseAddr+uint64(offset),
			goSyntax,
			intelSyntax,
		)
		lines = append(lines, line)
		offset += inst.Len
	}

	return lines, nil
}
