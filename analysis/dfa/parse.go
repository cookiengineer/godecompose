package dfa

import (
	"fmt"
	"strconv"
	"strings"
)

// Known registers in Go Plan9 x86_64 assembler.
var knownRegs = map[string]bool{
	"AX": true, "BX": true, "CX": true, "DX": true,
	"SI": true, "DI": true, "BP": true, "SP": true,
	"R8": true, "R9": true, "R10": true, "R11": true,
	"R12": true, "R13": true, "R14": true, "R15": true,
	"X0": true, "X1": true, "X2": true, "X3": true,
	"X4": true, "X5": true, "X6": true, "X7": true,
	"X8": true, "X9": true, "X10": true, "X11": true,
	"X12": true, "X13": true, "X14": true, "X15": true,
}

// ParseGoSyntaxOperands splits a GoSyntax string into opcode and operands.
// GoSyntax format: "OPCODE operand1, operand2" or "OPCODE"
func ParseGoSyntaxOperands(goSyntax string) (opcode string, operands []string) {
	goSyntax = strings.TrimSpace(goSyntax)

	spaceIdx := strings.IndexByte(goSyntax, ' ')
	if spaceIdx < 0 {
		return goSyntax, nil
	}

	opcode = goSyntax[:spaceIdx]
	rest := strings.TrimSpace(goSyntax[spaceIdx+1:])

	parts := splitPlan9Operands(rest)
	return opcode, parts
}

func splitPlan9Operands(s string) []string {
	var parts []string
	current := ""
	parenDepth := 0

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '(':
			parenDepth++
			current += string(ch)
		case ')':
			parenDepth--
			current += string(ch)
		case ',':
			if parenDepth == 0 && !inScaleContext(current) {
				parts = append(parts, strings.TrimSpace(current))
				current = ""
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}

func inScaleContext(s string) bool {
	starIdx := strings.LastIndexByte(s, '*')
	if starIdx < 0 {
		return false
	}
	after := s[starIdx+1:]
	_, err := strconv.Atoi(after)
	return err == nil
}

// ParseOperand parses a single Plan9 operand into its components.
func ParseOperand(operand string) (*MemRef, *Value) {
	if operand == "" {
		return nil, nil
	}

	if strings.HasPrefix(operand, "$") {
		imm := operand[1:]
		val, err := strconv.ParseUint(imm, 0, 64)
		if err != nil {
			return nil, nil
		}
		return nil, ConstValue(val)
	}

	if isRegister(operand) {
		return nil, RegValue(operand)
	}

	if strings.Contains(operand, "(") {
		mem := parseMemOperand(operand)
		return mem, nil
	}

	if strings.Contains(operand, ".") && strings.HasSuffix(operand, "(SB)") {
		mem := parseMemOperand(operand)
		return mem, nil
	}

	val, err := strconv.ParseUint(operand, 0, 64)
	if err == nil {
		return nil, ConstValue(val)
	}

	return nil, nil
}

func isRegister(s string) bool {
	return knownRegs[s]
}

func parseMemOperand(operand string) *MemRef {
	mem := &MemRef{}

	sbIdx := strings.Index(operand, "(SB)")
	if sbIdx >= 0 {
		prefix := operand[:sbIdx]
		if idx := strings.LastIndexByte(prefix, '+'); idx >= 0 {
			mem.Symbol = strings.TrimSpace(prefix[:idx])
			offStr := strings.TrimSpace(prefix[idx+1:])
			if off, err := strconv.ParseInt(offStr, 0, 64); err == nil {
				mem.Offset = off
			}
		} else {
			mem.Symbol = strings.TrimSpace(prefix)
		}
		mem.Base = "SB"
		return mem
	}

	openParen := strings.IndexByte(operand, '(')
	if openParen < 0 {
		return mem
	}

	offsetStr := strings.TrimSpace(operand[:openParen])
	if offsetStr != "" {
		if off, err := strconv.ParseInt(offsetStr, 0, 64); err == nil {
			mem.Offset = off
		}
	}

	firstClose := findMatchingParen(operand, openParen)
	if firstClose < 0 {
		return mem
	}

	baseInner := operand[openParen+1 : firstClose]
	mem.Base = baseInner

	after := operand[firstClose+1:]
	if strings.HasPrefix(after, "(") {
		idxClose := findMatchingParen(after, 0)
		if idxClose > 0 {
			idxInner := after[1:idxClose]
			parts := strings.SplitN(idxInner, "*", 2)
			if len(parts) == 2 {
				mem.Index = parts[0]
				if s, err := strconv.Atoi(parts[1]); err == nil {
					mem.Scale = s
				}
			}
		}
	}

	return mem
}

func findMatchingParen(s string, openIdx int) int {
	depth := 1
	for i := openIdx + 1; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// ParseInstructions splits GoSyntax into a list of (opcode, operands) pairs.
// Handles REP prefix instructions like "REP; MOVSQ".
func ParseInstructions(goSyntax string) []struct{ Opcode string; Operands []string } {
	goSyntax = strings.TrimSpace(goSyntax)
	if strings.Contains(goSyntax, ";") {
		var result []struct{ Opcode string; Operands []string }
		parts := strings.Split(goSyntax, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			op, ops := ParseGoSyntaxOperands(part)
			result = append(result, struct{ Opcode string; Operands []string }{op, ops})
		}
		return result
	}

	op, ops := ParseGoSyntaxOperands(goSyntax)
	return []struct{ Opcode string; Operands []string }{{op, ops}}
}

// ParseImmediate parses a $imm operand and returns the uint64 value.
func ParseImmediate(imm string) uint64 {
	if strings.HasPrefix(imm, "$") {
		imm = imm[1:]
	}
	val, err := strconv.ParseUint(imm, 0, 64)
	if err != nil {
		return 0
	}
	return val
}

// ExtractCallTarget returns the function name from a CALL GoSyntax string.
// "CALL fmt.Println(SB)" → "fmt.Println"
func ExtractCallTarget(goSyntax string) string {
	op, ops := ParseGoSyntaxOperands(goSyntax)
	if op != "CALL" || len(ops) == 0 {
		return ""
	}
	target := ops[0]
	if sb := strings.Index(target, "(SB)"); sb >= 0 {
		return strings.TrimSpace(target[:sb])
	}
	return target
}

// ExtractBranchTarget returns the branch target address from a Jcc or JMP GoSyntax.
func ExtractBranchTarget(goSyntax string) uint64 {
	_, ops := ParseGoSyntaxOperands(goSyntax)
	if len(ops) == 0 {
		return 0
	}
	val, _ := strconv.ParseUint(ops[0], 0, 64)
	return val
}

func fmtHex(n uint64) string {
	if n >= 0x10000 {
		return fmt.Sprintf("0x%x", n)
	}
	return fmt.Sprintf("%d", n)
}
