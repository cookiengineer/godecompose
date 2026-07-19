package function

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
)

func TestNormalizeReg(t *testing.T) {
	tests := []struct{ in, want string }{
		{"rax", "RAX"}, {"eax", "RAX"}, {"ax", "RAX"}, {"al", "RAX"},
		{"rbx", "RBX"}, {"ebx", "RBX"}, {"bl", "RBX"},
		{"rcx", "RCX"}, {"ecx", "RCX"}, {"cl", "RCX"},
		{"rdx", "RDX"}, {"edx", "RDX"}, {"dl", "RDX"},
		{"rsi", "RSI"}, {"esi", "RSI"}, {"si", "RSI"},
		{"rdi", "RDI"}, {"edi", "RDI"}, {"di", "RDI"},
		{"r8", "R8"}, {"r10", "R10"}, {"r14", "R14"},
		{"unknown", "UNKNOWN"},
	}
	for _, tc := range tests {
		got := normalizeReg(tc.in)
		if got != tc.want {
			t.Errorf("normalizeReg(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestExtractRegister(t *testing.T) {
	tests := []struct{ part, want string }{
		{"rax", "rax"},
		{"rbx", "rbx"},
		{"r12", "r12"},
		{"qword ptr [rax+0x8]", "rax"},   // memory operand -> extract base
		{"dword ptr [rbx]", "rbx"},         // memory operand no offset
		{"word ptr [rsi+0x10]", "rsi"},    // word memory
		{"byte ptr [rdi]", "rdi"},          // byte memory
		{"qword ptr [rsp+0x50]", "rsp"},   // stack
		{"ptr [rip+0x12345]", ""},         // RIP-relative (not in known list)
		{"0x2a", ""},                       // immediate
	}
	for _, tc := range tests {
		got := extractRegister(tc.part)
		if got != tc.want {
			t.Errorf("extractRegister(%q) = %q, want %q", tc.part, got, tc.want)
		}
	}
}

func TestExtractMemOffset(t *testing.T) {
	tests := []struct{ intel, want string }{
		{"mov qword ptr [rax+0x8], rbx", "8"},
		{"mov qword ptr [rbx+0x10], rcx", "10"},
		{"mov qword ptr [rdi+0x28], rsi", "28"},
		{"mov qword ptr [rsp+0x50], rax", "50"},
		{"mov qword ptr [rax], rbx", ""},                 // no offset
		{"lea rax, ptr [rip+0x12345]", ""},               // RIP-relative
		{"mov qword ptr [rax+rbx], rcx", ""},             // register index
		{"mov qword ptr [rax+0xrbx], rcx", ""},           // non-hex offset
		{"mov qword ptr [rcx-0x8], rax", "8"},            // negative (but takes absolute)
		{"mov qword ptr [rax+0x300], rbx", ""},           // too large
		{"mov qword ptr [rax+0x0], rbx", ""},             // zero offset
	}
	for _, tc := range tests {
		got := extractMemOffset(tc.intel)
		if got != tc.want {
			t.Errorf("extractMemOffset(%q) = %q, want %q", tc.intel, got, tc.want)
		}
	}
}

func TestReconstructSignatureArgs(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "MOV", IntelSyntax: "mov qword ptr [rsp+0x50], rax", Size: 5},
		{Address: 0x1005, Opcode: "MOV", IntelSyntax: "mov qword ptr [rsp+0x48], rbx", Size: 5},
		{Address: 0x100a, Opcode: "MOV", IntelSyntax: "mov qword ptr [rsp+0x40], rcx", Size: 5},
		{Address: 0x100f, Opcode: "CALL", IntelSyntax: "call runtime.convT64", IsCall: true, Size: 5},
		{Address: 0x1014, Opcode: "MOV", IntelSyntax: "mov rax, rbx", Size: 3},
		{Address: 0x1017, Opcode: "RET", IntelSyntax: "ret", IsReturn: true, Size: 1},
	}
	blocks := disasm.BuildControlFlowGraph(insts, []uint64{0x1000})
	f := &Function{
		Name:       "test.process",
		ShortName:  "process",
		EntryPoint: 0x1000,
		EndAddr:    0x1018,
		Blocks:     blocks,
	}

	sig := ReconstructSignature(f)
	if len(sig.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(sig.Args))
	}
	if len(sig.Returns) == 0 {
		t.Errorf("expected at least 1 return, got %d", len(sig.Returns))
	}
	s := sig.String()
	if !strings.Contains(s, "arg0") {
		t.Errorf("expected arg0 in signature, got %q", s)
	}
	if !strings.Contains(s, "int64") {
		t.Errorf("expected int64 type hint from convT64, got %q", s)
	}
	t.Logf("signature: %s", s)
}

func TestReconstructSignatureMethod(t *testing.T) {
	insts := []disasm.Instruction{
		{Address: 0x1000, Opcode: "MOV", IntelSyntax: "mov qword ptr [rax+0x50], rbx", Size: 5},
		{Address: 0x1005, Opcode: "RET", IntelSyntax: "ret", Size: 1, IsReturn: true},
	}
	blocks := disasm.BuildControlFlowGraph(insts, []uint64{0x1000})
	f := &Function{
		Name:             "main.(*Point).setX",
		ShortName:        "setX",
		ReceiverType:     "Point",
		IsPointerReceiver: true,
		IsMethod:         true,
		PackagePath:      "main",
		EntryPoint:       0x1000,
		EndAddr:          0x1006,
		Blocks:           blocks,
	}

	sig := ReconstructSignature(f)
	if !sig.IsPointer {
		t.Error("expected pointer receiver")
	}
	if sig.Receiver != "Point" {
		t.Errorf("expected receiver Point, got %q", sig.Receiver)
	}
	s := sig.String()
	if !strings.Contains(s, "(p *Point)") {
		t.Errorf("expected pointer receiver, got %q", s)
	}
	t.Logf("signature: %s", s)
}

func TestInferStructFields(t *testing.T) {
	f1 := &Function{
		Name: "main.(*Point).setX", ShortName: "setX",
		ReceiverType: "Point", IsPointerReceiver: true, IsMethod: true,
		PackagePath: "main",
		Blocks: []*disasm.BasicBlock{{
			Instructions: []disasm.Instruction{
				{IntelSyntax: "mov qword ptr [rax+0x8], rbx"},
				{IntelSyntax: "mov qword ptr [rax+0x10], rcx"},
				{IntelSyntax: "ret", IsReturn: true},
			},
		}},
	}
	f2 := &Function{
		Name: "main.(*Point).getX", ShortName: "getX",
		ReceiverType: "Point", IsPointerReceiver: true, IsMethod: true,
		PackagePath: "main",
		Blocks: []*disasm.BasicBlock{{
			Instructions: []disasm.Instruction{
				{IntelSyntax: "mov rax, qword ptr [rbx+0x8]"},
				{IntelSyntax: "mov qword ptr [rbx+0x20], 0x1"},
				{IntelSyntax: "ret", IsReturn: true},
			},
		}},
	}
	f1.SetPackageInfo()
	f2.SetPackageInfo()

	st := &StructType{Name: "Point", PackagePath: "main", Methods: []*Function{f1, f2}}
	fields := InferStructFields(st)

	if len(fields) < 3 {
		t.Errorf("expected at least 3 fields, got %d", len(fields))
	}

	found := map[string]bool{}
	for _, fld := range fields {
		found[fld.Offset] = true
		t.Logf("field: %s offset=%s type=%s count=%d", fld.Name, fld.Offset, fld.Type, fld.Count)
	}
	if !found["0x8"] {
		t.Error("expected field at offset 0x8")
	}
	if !found["0x10"] {
		t.Error("expected field at offset 0x10")
	}
	if !found["0x20"] {
		t.Error("expected field at offset 0x20")
	}
}

func TestReconstructSignatureEmpty(t *testing.T) {
	f := &Function{Name: "test.empty", ShortName: "empty"}
	sig := ReconstructSignature(f)
	if sig == nil {
		t.Fatal("expected non-nil signature")
	}
	s := sig.String()
	if !strings.Contains(s, "func empty()") {
		t.Errorf("expected empty signature, got %q", s)
	}
}

func TestNameHeuristics(t *testing.T) {
	tests := []struct{ name, wantRet string }{
		{"isOk", "bool"},
		{"hasValue", "bool"},
		{"canFail", "bool"},
		{"shouldRetry", "bool"},
		{"greet", ""},
		{"compute", ""},
	}
	for _, tc := range tests {
		got := inferReturnFromName(tc.name)
		if got != tc.wantRet {
			t.Errorf("inferReturnFromName(%q) = %q, want %q", tc.name, got, tc.wantRet)
		}
	}

	errTests := []string{"myError", "fatalError", "readError"}
	for _, n := range errTests {
		if !nameSuggestsError(n) {
			t.Errorf("nameSuggestsError(%q) = false, want true", n)
		}
	}
	if nameSuggestsError("greet") {
		t.Error("greet should not suggest error")
	}
}
