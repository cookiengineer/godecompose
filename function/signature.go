package function

import (
	"fmt"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
)

var argRegs = []string{"RAX", "RBX", "RCX", "RDI", "RSI", "R8", "R9", "R10", "R11"}

type Signature struct {
	Name      string
	Args      []Param
	Returns   []Param
	IsPointer bool
	Receiver  string
}

type Param struct {
	Name string
	Type string
}

func ReconstructSignature(f *Function) *Signature {
	sig := &Signature{
		Name:      f.ShortName,
		Receiver:  f.ReceiverType,
		IsPointer: f.IsPointerReceiver,
	}

	blocks := f.Blocks
	if len(blocks) == 0 {
		return sig
	}

	entryBlock := blocks[0]
	reads := make(map[string]bool)
	writes := make(map[string]bool)
	spilled := make(map[string]bool)
	bodyTypes := inferBodyTypes(f)

	for _, inst := range entryBlock.Instructions {
		srcRegs, dstRegs := extractRegs(inst.IntelSyntax)
		for _, r := range srcRegs {
			n := normalizeReg(r)
			if !writes[n] {
				reads[n] = true
			}
		}
		for _, r := range dstRegs {
			n := normalizeReg(r)
			writes[n] = true
			if strings.Contains(inst.IntelSyntax, "[") && n != "RSP" && n != "RBP" && n != "R14" {
				spilled[n] = true
			}
		}
	}

	for _, argReg := range argRegs {
		if reads[argReg] || spilled[argReg] {
			p := Param{Name: fmt.Sprintf("arg%d", len(sig.Args))}
			if bodyTypes.argTypeHint != "" {
				p.Type = bodyTypes.argTypeHint
			}
			sig.Args = append(sig.Args, p)
		}
	}

	returnWrites := make(map[string]bool)
	for _, block := range blocks {
		lastInst := block.Instructions[len(block.Instructions)-1]
		if !lastInst.IsReturn && lastInst.Opcode != "RET" && lastInst.Opcode != "RETQ" {
			continue
		}
		for _, inst := range block.Instructions {
			_, dstRegs := extractRegs(inst.IntelSyntax)
			for _, r := range dstRegs {
				n := normalizeReg(r)
				returnWrites[n] = true
			}
		}
	}

	hasErrorReturn := bodyTypes.hasError || nameSuggestsError(f.ShortName) ||
		(len(blocks) > 0 && hasErrorPattern(blocks))

	for _, retReg := range argRegs {
		if returnWrites[retReg] {
			p := Param{}
			if hasErrorReturn && len(sig.Returns) == 0 {
				continue
			} else if bodyTypes.retType != "" {
				p.Type = bodyTypes.retType
			} else if retType := inferReturnFromName(f.ShortName); retType != "" {
				p.Type = retType
			}
			sig.Returns = append(sig.Returns, p)
		}
	}

	if hasErrorReturn {
		sig.Returns = append(sig.Returns, Param{Type: "error"})
	}

	return sig
}

type bodyTypeInfo struct {
	argTypeHint string
	retType     string
	hasError    bool
}

func inferBodyTypes(f *Function) bodyTypeInfo {
	info := bodyTypeInfo{}

	for _, block := range f.Blocks {
		for _, inst := range block.Instructions {
			if !inst.IsCall {
				continue
			}
			intel := strings.ToLower(inst.IntelSyntax)

			switch {
			case containsAny(intel, "runtime.convtstring", "convTstring"):
				info.argTypeHint = "string"
			case containsAny(intel, "runtime.convt64", "convT64"):
				info.argTypeHint = "int64"
			case containsAny(intel, "runtime.convt32", "convT32"):
				info.argTypeHint = "int32"
			case containsAny(intel, "runtime.convt2e", "convT2E"):
				info.argTypeHint = "interface{}"
			case containsAny(intel, "runtime.newobject"):
				info.retType = "interface{}"
			case containsAny(intel, "runtime.makeslice"):
				info.retType = "[]interface{}"
			case containsAny(intel, "runtime.makemap"):
				info.retType = "map[K]V"
			case containsAny(intel, "runtime.makechan"):
				info.retType = "chan T"
			case containsAny(intel, "runtime.growslice"):
				_ = 1
			case containsAny(intel, "fmt.errorf", "errors.new"):
				info.hasError = true
			case containsAny(intel, "runtime.panic", "runtime.gopanic"):
				_ = 1
			}
		}
	}
	return info
}

func nameSuggestsError(name string) bool {
	suffixes := []string{"Error", "error"}
	for _, s := range suffixes {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func inferReturnFromName(name string) string {
	if strings.HasPrefix(name, "is") || strings.HasPrefix(name, "has") ||
		strings.HasPrefix(name, "can") || strings.HasPrefix(name, "should") {
		return "bool"
	}
	return ""
}

func hasErrorPattern(blocks []*disasm.BasicBlock) bool {
	for _, block := range blocks {
		for _, inst := range block.Instructions {
			if inst.IsReturn || inst.Opcode == "RET" || inst.Opcode == "RETQ" {
				continue
			}
			if inst.IsCall {
				intel := strings.ToLower(inst.IntelSyntax)
				if containsAny(intel, "runtime.convTstring", "runtime.convT64") {
					return true
				}
			}
		}
	}
	return false
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func (s *Signature) String() string {
	var b strings.Builder
	name := sanitizeFuncName(s.Name)

	if s.IsPointer && s.Receiver != "" {
		recv := strings.ToLower(s.Receiver[:1])
		b.WriteString(fmt.Sprintf("func (%s *%s) %s(", recv, s.Receiver, name))
	} else if s.Receiver != "" {
		recv := strings.ToLower(s.Receiver[:1])
		b.WriteString(fmt.Sprintf("func (%s %s) %s(", recv, s.Receiver, name))
	} else {
		b.WriteString(fmt.Sprintf("func %s(", name))
	}

	for i, a := range s.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		if a.Type != "" {
			b.WriteString(fmt.Sprintf("%s %s", a.Name, a.Type))
		} else {
			b.WriteString(fmt.Sprintf("%s interface{}", a.Name))
		}
	}
	b.WriteString(")")

	if len(s.Returns) > 0 {
		if len(s.Returns) == 1 && s.Returns[0].Type == "error" {
			b.WriteString(" error")
		} else {
			b.WriteString(" (")
			for i, r := range s.Returns {
				if i > 0 {
					b.WriteString(", ")
				}
				if r.Type != "" {
					b.WriteString(r.Type)
				} else {
					b.WriteString(fmt.Sprintf("ret%d", i))
				}
			}
			b.WriteString(")")
		}
	}

	return b.String()
}

func extractRegs(intel string) (src []string, dst []string) {
	intel = strings.TrimSpace(intel)
	spaceIdx := strings.Index(intel, " ")
	if spaceIdx < 0 {
		return nil, nil
	}
	operands := intel[spaceIdx+1:]
	parts := strings.Split(operands, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		reg := extractRegister(part)
		if reg != "" {
			if i == 0 {
				dst = append(dst, reg)
			} else {
				src = append(src, reg)
			}
		}
	}
	return src, dst
}

func extractRegister(part string) string {
	known := []string{
		"rax", "rbx", "rcx", "rdx", "rsi", "rdi", "rsp", "rbp",
		"r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
		"eax", "ebx", "ecx", "edx", "esi", "edi", "esp", "ebp",
		"ax", "bx", "cx", "dx", "si", "di",
	}
	part = strings.TrimPrefix(part, "qword ptr [")
	part = strings.TrimPrefix(part, "dword ptr [")
	part = strings.TrimPrefix(part, "word ptr [")
	part = strings.TrimPrefix(part, "byte ptr [")
	part = strings.TrimSuffix(part, "]")
	part = strings.TrimSpace(part)
	if idx := strings.Index(part, "+"); idx >= 0 {
		part = part[:idx]
	}
	if idx := strings.Index(part, "-"); idx >= 0 {
		part = part[:idx]
	}
	part = strings.TrimSpace(part)

	for _, k := range known {
		if strings.EqualFold(part, k) {
			return k
		}
	}
	return ""
}

func normalizeReg(reg string) string {
	reg = strings.ToUpper(reg)
	switch reg {
	case "EAX", "AX", "AL":
		return "RAX"
	case "EBX", "BX", "BL":
		return "RBX"
	case "ECX", "CX", "CL":
		return "RCX"
	case "EDX", "DX", "DL":
		return "RDX"
	case "ESI", "SI", "SIL":
		return "RSI"
	case "EDI", "DI", "DIL":
		return "RDI"
	case "ESP", "SP", "SPL":
		return "RSP"
	case "EBP", "BP", "BPL":
		return "RBP"
	}
	if strings.HasPrefix(reg, "R") {
		return reg
	}
	return reg
}

// FieldInfo describes a detected struct field from offset analysis.
type FieldInfo struct {
	Offset string
	Name   string
	Type   string
	Count  int
}

// InferStructFields analyzes memory accesses across methods of a struct type
// and groups them by offset to infer field names.
func InferStructFields(st *StructType) []FieldInfo {
	if st == nil || len(st.Methods) == 0 {
		return nil
	}

	offsets := make(map[string]*FieldInfo)
	offsetHints := make(map[string]map[string]int)

	for _, m := range st.Methods {
		if len(m.Blocks) == 0 {
			continue
		}
		hints := extractMethodNameHints(m.ShortName)
		for _, block := range m.Blocks {
			for _, inst := range block.Instructions {
				intel := inst.IntelSyntax
				if strings.Contains(intel, "rsp") || strings.Contains(intel, "rbp") || strings.Contains(intel, "RSP") || strings.Contains(intel, "RBP") {
					continue
				}
				offset := extractMemOffset(intel)
				if offset == "" {
					continue
				}
				key := fmt.Sprintf("0x%s", offset)
				if fi, ok := offsets[key]; ok {
					fi.Count++
				} else {
					fi = &FieldInfo{Offset: key, Count: 1}
					fi.Name = fmt.Sprintf("field_%s", offset)
					fi.Type = inferFieldTypeFromInst(inst.IntelSyntax)
					offsets[key] = fi
				}

				if offsetHints[key] == nil {
					offsetHints[key] = make(map[string]int)
				}
				for _, hint := range hints {
					offsetHints[key][hint]++
				}
			}
		}
	}

	for key, fi := range offsets {
		hints := offsetHints[key]
		if len(hints) > 0 {
			best := ""
			bestCount := 0
			for hint, count := range hints {
				if count > bestCount {
					best = hint
					bestCount = count
				}
			}
			if best != "" {
				fi.Name = best
			}
		}

		if fi.Type == "" {
			fi.Type = inferTypeFromOffset(fi.Offset)
		}
	}

	result := make([]FieldInfo, 0, len(offsets))
	for _, fi := range offsets {
		result = append(result, *fi)
	}
	return result
}

func inferTypeFromOffset(offsetHex string) string {
	var offset int64
	fmt.Sscanf(offsetHex, "0x%x", &offset)
	switch {
	case offset < 0x08:
		return ""
	case offset < 0x20:
		return "int"
	case offset < 0x40:
		return "int"
	default:
		return "int"
	}
}

var knownRegsSet = map[string]bool{
	"rax": true, "rbx": true, "rcx": true, "rdx": true,
	"rsi": true, "rdi": true, "rsp": true, "rbp": true,
	"r8": true, "r9": true, "r10": true, "r11": true,
	"r12": true, "r13": true, "r14": true, "r15": true,
	"eax": true, "ebx": true, "ecx": true, "edx": true,
	"esi": true, "edi": true, "esp": true, "ebp": true,
	"ax": true, "bx": true, "cx": true, "dx": true,
}

func isKnownRegister(s string) bool {
	return knownRegsSet[strings.ToLower(s)]
}

func extractMethodNameHints(methodName string) []string {
	var hints []string

	type prefixRule struct {
		prefix string
		suffix string
	}

	rules := []prefixRule{
		{"Get", ""},
		{"Set", ""},
		{"Is", ""},
		{"Has", ""},
		{"Calc", ""},
		{"Calc", "ulate"},
	}

	for _, r := range rules {
		if strings.HasPrefix(methodName, r.prefix) {
			rest := methodName[len(r.prefix):]
			if r.suffix != "" {
				if !strings.HasPrefix(rest, r.suffix) {
					continue
				}
				rest = rest[len(r.suffix):]
			} else {
				if len(rest) == 0 || rest[0] < 'A' || rest[0] > 'Z' {
					continue
				}
			}
			if len(rest) > 0 {
				hints = append(hints, strings.ToLower(rest))
			}
		}
	}

	return hints
}

func extractMemOffset(intel string) string {
	bStart := strings.Index(intel, "[")
	bEnd := strings.Index(intel, "]")
	if bStart < 0 || bEnd < 0 {
		return ""
	}
	inner := intel[bStart+1 : bEnd]
	inner = strings.TrimSpace(inner)

	if strings.Contains(inner, "rip") || strings.Contains(inner, "eip") {
		return ""
	}

	plus := strings.Index(inner, "+")
	minus := strings.Index(inner, "-")
	if plus >= 0 {
		inner = inner[plus+1:]
	} else if minus >= 0 {
		inner = inner[minus+1:]
	} else {
		if isKnownRegister(inner) {
			return "0"
		}
		return ""
	}
	inner = strings.TrimSpace(inner)
	if strings.Contains(inner, "*") {
		return ""
	}

	if strings.HasPrefix(inner, "0x") {
		inner = inner[2:]
	}

	if inner == "" {
		return ""
	}

	// Must be a valid hex number (not a register name)
	isHex := true
	for _, ch := range inner {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			isHex = false
			break
		}
	}
	if !isHex {
		return ""
	}

	val := int64(0)
	fmt.Sscanf(inner, "%x", &val)
	if val < 0 || val > 0x200 {
		return ""
	}

	return inner
}

func inferFieldTypeFromInst(intel string) string {
	lower := strings.ToLower(intel)
	if strings.Contains(lower, "movb") || strings.Contains(lower, "mov byte") || strings.Contains(lower, "cmpb") || strings.Contains(lower, "cmp byte") {
		return "bool"
	}
	if strings.Contains(lower, "convtstring") || strings.Contains(lower, "conv.tstring") || strings.Contains(lower, "convTstring") {
		return "string"
	}
	if strings.Contains(lower, "convt64") || strings.Contains(lower, "conv.t64") || strings.Contains(lower, "convT64") {
		return "int64"
	}
	if strings.Contains(lower, "qword") {
		return "int"
	}
	return ""
}

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

func sanitizeFuncName(name string) string {
	name = strings.ReplaceAll(name, ".", "_")

	if name == "init" || name == "main_init" {
		return "_init"
	}

	parts := strings.Split(name, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		if goKeywords[p] {
			parts[i] = "_" + p
		}
	}
	return strings.Join(parts, "_")
}
