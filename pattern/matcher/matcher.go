// Package matcher matches compiled decompilation patterns against
// disassembled instruction streams, producing variable bindings for
// capture variables.
package matcher

import (
	"sort"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
)

// Match represents a successful pattern match against a sequence of instructions.
type Match struct {
	Pattern    *evaluator.CompiledPattern
	Alternative int
	StartAddr  uint64
	EndAddr    uint64
	Bindings   map[string]Binding
	Confidence float64
}

// Binding captures a variable name and its matched value.
type Binding struct {
	CaptureVar string
	Value      string
	Alias      string
}

// Matcher indexes compiled patterns by opcode and runs the matching algorithm.
type Matcher struct {
	patterns  []*evaluator.CompiledPattern
	byOpcode  map[string][]*evaluator.CompiledPattern
}

// New creates a matcher with the given compiled patterns.
func New(patterns []*evaluator.CompiledPattern) *Matcher {
	m := &Matcher{
		patterns: patterns,
		byOpcode: make(map[string][]*evaluator.CompiledPattern),
	}
	for _, p := range patterns {
		for _, alt := range p.Alternatives {
			if len(alt) == 0 {
				continue
			}
			firstOp := strings.ToUpper(alt[0].Opcode)
			m.byOpcode[firstOp] = append(m.byOpcode[firstOp], p)
		}
	}
	return m
}

// Match scans the instruction stream for pattern matches and returns the
// best matches with minimal overlap. Limits to 10000 raw matches.
func (m *Matcher) Match(instructions []disasm.Instruction) []Match {
	const maxRawMatches = 10000
	var allMatches []Match

	for i, inst := range instructions {
		if len(inst.Opcode) == 0 {
			continue
		}
		opcode := strings.ToUpper(inst.Opcode)
		candidates := m.byOpcode[opcode]
		if len(candidates) == 0 {
			continue
		}

		for _, pat := range candidates {
			for altIdx, alt := range pat.Alternatives {
				if len(alt) == 0 {
					continue
				}
				if strings.ToUpper(alt[0].Opcode) != opcode {
					continue
				}

				bindings, endIdx, confidence := m.tryMatch(instructions, i, alt)
				if bindings != nil {
					endAddr := instructions[endIdx].Address + uint64(instructions[endIdx].Size)
					match := Match{
						Pattern:     pat,
						Alternative: altIdx,
						StartAddr:   inst.Address,
						EndAddr:     endAddr,
						Bindings:    bindings,
						Confidence:  confidence,
					}
					allMatches = append(allMatches, match)

					if len(allMatches) >= maxRawMatches {
						return resolveConflicts(allMatches)
					}
				}
			}
		}
	}

	return resolveConflicts(allMatches)
}

func (m *Matcher) tryMatch(instructions []disasm.Instruction, startIdx int, alt []evaluator.CompiledInstruction) (map[string]Binding, int, float64) {
	bindings := make(map[string]Binding)
	matchCount := 0
	wildcardCount := 0
	operandCount := 0

	si := startIdx
	pi := 0

	for si < len(instructions) && pi < len(alt) {
		patternInst := alt[pi]
		streamInst := instructions[si]

		if patternInst.IsLabel {
			pi++
			continue
		}

		matched := instructionMatches(streamInst, patternInst)
		if !matched {
			return nil, 0, 0
		}

		if !matchOperands(streamInst, patternInst.Operands, bindings) {
			return nil, 0, 0
		}

		for _, op := range patternInst.Operands {
			operandCount++
			if op.IsWildcard {
				wildcardCount++
			}
		}

		matchCount++
		si++
		pi++
	}

	if pi < len(alt) {
		return nil, 0, 0
	}

	confidence := 1.0
	if operandCount > 0 {
		confidence = 1.0 - float64(wildcardCount)/float64(operandCount)
	}
	confidence += float64(matchCount) * 0.1
	if confidence > 1.0 {
		confidence = 1.0
	}

	di := si - 1
	if di < 0 {
		di = 0
	}

	return bindings, di, confidence
}

func instructionMatches(inst disasm.Instruction, pat evaluator.CompiledInstruction) bool {
	patOp := strings.ToUpper(pat.Opcode)
	instOp := strings.ToUpper(inst.Opcode)

	if patOp == "CALL" && instOp == "CALL" {
		if len(pat.Operands) > 0 && pat.Operands[0].CaptureVar != "" {
			target := pat.Operands[0].CaptureVar
			target = strings.ReplaceAll(target, "_", ".")
			if fuzzyMatchCall(inst.GoSyntax, target) || fuzzyMatchCall(inst.IntelSyntax, target) {
				return true
			}
			return false
		}
	}

	if !strings.EqualFold(instOp, patOp) {
		return false
	}

	return true
}

// fuzzyMatchCall checks if a CALL instruction's syntax matches a target
// function pattern like "sync.Mutex.Lock" against "sync.(*Mutex).lockSlow(SB)".
func fuzzyMatchCall(goSyntax, target string) bool {
	if strings.Contains(goSyntax, target) {
		return true
	}
	lower := strings.ToLower(goSyntax)
	lowerTarget := strings.ToLower(target)

	if strings.Contains(lower, lowerTarget) {
		return true
	}

	norm := strings.NewReplacer(".", " ", "(", " ", ")", " ", "/", " ", "*", " ", "_", " ").Replace(lower)
	normTarget := strings.NewReplacer(".", " ", "_", " ").Replace(lowerTarget)
	normParts := strings.Fields(norm)
	normTargetParts := strings.Fields(normTarget)

	if len(normTargetParts) == 0 {
		return false
	}

	for _, tp := range normTargetParts {
		found := false
		for _, p := range normParts {
			if strings.HasPrefix(p, tp) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func matchOperands(inst disasm.Instruction, patternOps []evaluator.CompiledOperand, bindings map[string]Binding) bool {
	if len(patternOps) == 0 {
		return true
	}

	intel := inst.IntelSyntax
	parts := operandParts(intel)
	if len(parts) > 0 {
		parts = parts[1:]
	}

	pi := 0
	for _, part := range parts {
		if pi >= len(patternOps) {
			break
		}
		pat := patternOps[pi]

		if pat.IsWildcard {
			pi++
			continue
		}

		if isTypeQualifier(part) {
			continue
		}

		if pat.BaseReg != "" && isMemoryRef(part) {
			if bindMemoryOperand(part, pat, bindings) {
				pi++
			}
			continue
		}

		matched := matchSingleOperand(part, pat, bindings)
		if matched {
			pi++
		}
	}

	return pi >= len(patternOps)
}

func isTypeQualifier(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "qword", "dword", "word", "byte", "ptr":
		return true
	}
	return false
}

func isMemoryRef(s string) bool {
	return strings.Contains(s, "[") && strings.Contains(s, "]")
}

func matchSingleOperand(opStr string, pat evaluator.CompiledOperand, bindings map[string]Binding) bool {
	opStr = strings.TrimSpace(opStr)

	if pat.IsImmediate && strings.HasPrefix(opStr, "0x") || strings.HasPrefix(opStr, "$") {
		if pat.CaptureVar != "" {
			bindings[pat.CaptureVar] = Binding{
				CaptureVar: pat.CaptureVar,
				Value:      opStr,
			}
		}
		return true
	}

	if pat.IsImmediate {
		if pat.CaptureVar != "" {
			bindings[pat.CaptureVar] = Binding{
				CaptureVar: pat.CaptureVar,
				Value:      opStr,
			}
		}
		return true
	}

	if pat.Register != "" && strings.EqualFold(opStr, pat.Register) {
		return true
	}

	if pat.CaptureVar != "" && !pat.IsWildcard {
		if existing, ok := bindings[pat.CaptureVar]; ok {
			return strings.EqualFold(opStr, existing.Value)
		}
		bindings[pat.CaptureVar] = Binding{
			CaptureVar: pat.CaptureVar,
			Value:      opStr,
		}
		return true
	}

	return false
}

func operandParts(intelSyntax string) []string {
	var parts []string
	var current strings.Builder
	inParen := 0

	for _, ch := range intelSyntax {
		switch ch {
		case '(':
			inParen++
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
			current.WriteRune(ch)
		case ')':
			inParen--
			current.WriteRune(ch)
			if inParen == 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		case ',':
			if inParen > 0 {
				current.WriteRune(ch)
			} else {
				if current.Len() > 0 {
					parts = append(parts, strings.TrimSpace(current.String()))
					current.Reset()
				}
			}
		case ' ':
			if inParen > 0 {
				current.WriteRune(ch)
			} else if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

func parseMemOp(opStr string) (base string, offset string) {
	bStart := strings.Index(opStr, "[")
	bEnd := strings.Index(opStr, "]")
	if bStart < 0 || bEnd < 0 || bEnd <= bStart {
		return "", ""
	}
	inner := opStr[bStart+1 : bEnd]

	plusIdx := strings.Index(inner, "+")
	minusIdx := strings.Index(inner, "-")
	if minusIdx >= 0 {
		plusIdx = minusIdx
	}
	if plusIdx < 0 {
		if strings.Contains(inner, "*") {
			parts := strings.Split(inner, "+")
			if len(parts) >= 2 {
				base = strings.TrimSpace(parts[len(parts)-1])
				offset = "0"
				return base, offset
			}
		}
		base = strings.TrimSpace(inner)
		offset = "0"
		return base, offset
	}

	basePart := strings.TrimSpace(inner[:plusIdx])
	offsetPart := strings.TrimSpace(inner[plusIdx+1:])

	if strings.Contains(basePart, "*") {
		base = strings.TrimSpace(inner[plusIdx+1:])
		offset = "0"
		return base, offset
	}

	base = basePart
	if minusIdx >= 0 {
		offset = "-" + offsetPart
	} else {
		offset = offsetPart
	}
	return base, offset
}

func bindMemoryOperand(opStr string, pat evaluator.CompiledOperand, bindings map[string]Binding) bool {
	base, offset := parseMemOp(opStr)
	if base == "" {
		return false
	}
	if pat.BaseReg != "" {
		if isKnownReg(pat.BaseReg) {
			if !strings.EqualFold(base, pat.BaseReg) {
				return false
			}
		} else {
			bindings[pat.BaseReg] = Binding{
				CaptureVar: pat.BaseReg,
				Value:      base,
			}
		}
	}
	if pat.Offset != "" && offset != "" {
		bindings[pat.Offset] = Binding{
			CaptureVar: pat.Offset,
			Value:      offset,
		}
	}
	return true
}

func isKnownReg(name string) bool {
	n := strings.ToUpper(name)
	switch n {
	case "AX", "BX", "CX", "DX", "SI", "DI", "BP", "SP",
		"R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15",
		"AL", "BL", "CL", "DL", "AH", "BH", "CH", "DH",
		"X0", "X1", "X2", "X3", "X4", "X5", "X6", "X7",
		"X8", "X9", "X10", "X11", "X12", "X13", "X14", "X15":
		return true
	}
	return false
}

func resolveConflicts(matches []Match) []Match {
	if len(matches) <= 1 {
		return matches
	}

	sorted := make([]Match, len(matches))
	copy(sorted, matches)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Confidence != sorted[j].Confidence {
			return sorted[i].Confidence > sorted[j].Confidence
		}
		return (sorted[i].EndAddr - sorted[i].StartAddr) > (sorted[j].EndAddr - sorted[j].StartAddr)
	})

	occupied := make([]Match, 0, len(sorted))
	for _, m := range sorted {
		overlap := false
		for _, o := range occupied {
			if m.StartAddr < o.EndAddr && m.EndAddr > o.StartAddr {
				overlap = true
				break
			}
		}
		if !overlap {
			occupied = append(occupied, m)
		}
	}

	return occupied
}
