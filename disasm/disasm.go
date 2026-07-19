// Package disasm provides x86_64 instruction decoding and control flow
// analysis using the Go standard library's x86asm package.
package disasm

import (
	"strings"

	"golang.org/x/arch/x86/x86asm"
)

const defaultMode = 64

// Instruction represents a single decoded x86_64 instruction with its
// address, raw bytes, and syntax representations.
type Instruction struct {
	Address       uint64
	Bytes         []byte
	Opcode        string
	IntelSyntax   string
	GoSyntax      string
	IsCall        bool
	IsReturn      bool
	IsBranch      bool
	IsConditional bool
	BranchTarget  uint64
	Size          int
}

// SymLookup resolves an address to a symbol name and its base address.
// Used by syntax formatters to produce symbolic names for call/jump targets.
type SymLookup func(addr uint64) (name string, base uint64)

// SymbolEntry is a minimal symbol representation for building lookup tables.
type SymbolEntry struct {
	Name    string
	Address uint64
	Size    uint64
}

// BuildSymLookup creates a SymLookup from a list of symbols. It returns
// the best-matching symbol name and base address for any given address.
func BuildSymLookup(symbols []SymbolEntry) SymLookup {
	if len(symbols) == 0 {
		return nil
	}
	return func(addr uint64) (string, uint64) {
		best := ""
		var bestAddr uint64
		for _, s := range symbols {
			if addr >= s.Address && addr < s.Address+s.Size && s.Size > 0 {
				if s.Name > best {
					best = s.Name
					bestAddr = s.Address
				}
			}
		}
		return best, bestAddr
	}
}

// DecodeStream decodes a linear byte stream into an instruction sequence.
// GoSyntax uses nil symbol lookup (raw hex addresses).
func DecodeStream(data []byte, baseAddr uint64) ([]Instruction, error) {
	return DecodeStreamWithSymbols(data, baseAddr, nil)
}

// DecodeStreamWithSymbols decodes instructions and resolves PC-relative
// addresses to symbol names using the provided lookup function.
// Pass nil for raw hex addresses.
func DecodeStreamWithSymbols(data []byte, baseAddr uint64, lookup SymLookup) ([]Instruction, error) {
	var instructions []Instruction
	offset := 0

	for offset < len(data) {
		inst, err := x86asm.Decode(data[offset:], defaultMode)
		if err != nil {
			offset++
			continue
		}

		pc := baseAddr + uint64(offset)
		var symLookup x86asm.SymLookup
		if lookup != nil {
			symLookup = x86asm.SymLookup(lookup)
		}

		goSyntax := x86asm.GoSyntax(inst, pc, symLookup)

		instr := Instruction{
			Address:      pc,
			Bytes:        data[offset : offset+inst.Len],
			Opcode:       extractOpcode(goSyntax),
			IntelSyntax:  x86asm.IntelSyntax(inst, pc, symLookup),
			GoSyntax:     goSyntax,
			Size:         inst.Len,
			IsCall:       inst.Op == x86asm.CALL,
			IsReturn:     inst.Op == x86asm.RET,
			IsBranch:     isBranchOp(inst.Op),
			IsConditional: isConditionalBranch(inst.Op),
			BranchTarget: resolveBranchTarget(inst, pc),
		}

		instructions = append(instructions, instr)
		offset += inst.Len
	}

	return instructions, nil
}

func isBranchOp(op x86asm.Op) bool {
	switch op {
	case x86asm.JMP, x86asm.JA, x86asm.JAE, x86asm.JB, x86asm.JBE,
		x86asm.JCXZ, x86asm.JE, x86asm.JECXZ, x86asm.JG, x86asm.JGE,
		x86asm.JL, x86asm.JLE, x86asm.JNE, x86asm.JNO, x86asm.JNP,
		x86asm.JNS, x86asm.JO, x86asm.JP, x86asm.JRCXZ, x86asm.JS,
		x86asm.CALL, x86asm.RET, x86asm.LOOP, x86asm.LOOPE, x86asm.LOOPNE,
		x86asm.IRET, x86asm.IRETD, x86asm.IRETQ:
		return true
	default:
		return false
	}
}

func isConditionalBranch(op x86asm.Op) bool {
	switch op {
	case x86asm.JA, x86asm.JAE, x86asm.JB, x86asm.JBE,
		x86asm.JCXZ, x86asm.JE, x86asm.JECXZ, x86asm.JG, x86asm.JGE,
		x86asm.JL, x86asm.JLE, x86asm.JNE, x86asm.JNO, x86asm.JNP,
		x86asm.JNS, x86asm.JO, x86asm.JP, x86asm.JRCXZ, x86asm.JS,
		x86asm.LOOP, x86asm.LOOPE, x86asm.LOOPNE:
		return true
	default:
		return false
	}
}

func resolveBranchTarget(inst x86asm.Inst, pc uint64) uint64 {
	if inst.Op == x86asm.RET || inst.Op == x86asm.IRET || inst.Op == x86asm.IRETD || inst.Op == x86asm.IRETQ {
		return 0
	}

	for _, arg := range inst.Args {
		if arg == nil {
			continue
		}
		switch a := arg.(type) {
		case x86asm.Rel:
			return pc + uint64(inst.Len) + uint64(a)
		case x86asm.Imm:
			if inst.Op == x86asm.CALL || inst.Op == x86asm.JMP || isConditionalBranch(inst.Op) {
				return uint64(a)
			}
		}
	}

	return 0
}

func extractOpcode(goSyntax string) string {
	for i, ch := range goSyntax {
		if ch == ' ' || ch == '\t' || ch == ';' {
			opcode := goSyntax[:i]
			if opcode == "REP" && ch == ';' {
				rest := strings.TrimLeft(goSyntax[i+1:], " \t")
				for j, c := range rest {
					if c == ' ' || c == '\t' {
						return normalizeOpcode(rest[:j])
					}
				}
				return normalizeOpcode(rest)
			}
			return normalizeOpcode(opcode)
		}
	}
	return normalizeOpcode(goSyntax)
}

func normalizeOpcode(opcode string) string {
	base := opcode
	if len(opcode) > 1 {
		switch opcode[len(opcode)-1] {
		case 'Q', 'L', 'W', 'B':
			base = opcode[:len(opcode)-1]
		}
	}
	switch opcode {
	case "JE":
		return "JEQ"
	case "JG":
		return "JGT"
	case "JL":
		return "JLT"
	case "JB":
		return "JLO"
	case "JA":
		return "JHI"
	case "JAE":
		return "JCC"
	case "JBE":
		return "JLS"
	case "JS":
		return "JMI"
	case "JNS":
		return "JPL"
	case "JO":
		return "JOS"
	case "JNO":
		return "JOC"
	case "JP":
		return "JPS"
	case "JNP":
		return "JPC"
	case "NOPL", "NOPW":
		return "NOP"
	}
	switch base {
	case "TEST", "CMP":
		return base
	}
	return opcode
}

func opcodeToX86asm(opcode string) x86asm.Op {
	opMap := map[string]x86asm.Op{
		"CALL": x86asm.CALL,
		"RET":  x86asm.RET,
		"JMP":  x86asm.JMP,
		"JE":   x86asm.JE,
		"JNE":  x86asm.JNE,
		"JG":   x86asm.JG,
		"JGE":  x86asm.JGE,
		"JL":   x86asm.JL,
		"JLE":  x86asm.JLE,
		"JA":   x86asm.JA,
		"JAE":  x86asm.JAE,
		"JB":   x86asm.JB,
		"JBE":  x86asm.JBE,
		"JO":   x86asm.JO,
		"JNO":  x86asm.JNO,
		"JS":   x86asm.JS,
		"JNS":  x86asm.JNS,
		"JP":   x86asm.JP,
		"JNP":  x86asm.JNP,
	}
	if op, ok := opMap[opcode]; ok {
		return op
	}
	return 0
}
