package parser

import (
	"testing"

	"github.com/cookiengineer/godecompose/pattern/lang/ast"
	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
)

func parse(t *testing.T, input string) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	p := New(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return program
}

func TestParseStruct(t *testing.T) {
	input := `struct Header {
    u32 magic;
    u32 version;
    padding[16];
    double timestamp;
};`

	prog := parse(t, input)
	if len(prog.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(prog.Structs))
	}
	s := prog.Structs[0]
	if s.Name != "Header" {
		t.Errorf("struct name = %q, want Header", s.Name)
	}
	if len(s.Members) != 4 {
		t.Errorf("expected 4 members, got %d", len(s.Members))
	}
}

func TestParseEnum(t *testing.T) {
	input := `enum FileType : u16 {
    ELF = 0x457F,
    PE  = 0x5A4D,
    MachO
};`

	prog := parse(t, input)
	if len(prog.Enums) != 1 {
		t.Fatalf("expected 1 enum, got %d", len(prog.Enums))
	}
	e := prog.Enums[0]
	if e.Name != "FileType" {
		t.Errorf("enum name = %q, want FileType", e.Name)
	}
	if e.Base != "u16" {
		t.Errorf("enum base = %q, want u16", e.Base)
	}
	if len(e.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(e.Values))
	}
	if e.Values[0].Name != "ELF" {
		t.Errorf("val[0].Name = %q, want ELF", e.Values[0].Name)
	}
}

func TestParseVariable(t *testing.T) {
	input := `u32 magic @ 0x00;
str name = "test";
u8 data[256] @ 0x100;`

	prog := parse(t, input)
	if len(prog.Variables) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(prog.Variables))
	}
	v0 := prog.Variables[0]
	if v0.Name != "magic" {
		t.Errorf("var[0].Name = %q, want magic", v0.Name)
	}
}

func TestParseFunction(t *testing.T) {
	input := `fn add(u32 a, u32 b) {
    return a + b;
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(prog.Functions))
	}
	f := prog.Functions[0]
	if f.Name != "add" {
		t.Errorf("fn name = %q, want add", f.Name)
	}
	if len(f.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(f.Params))
	}
	if f.Body == nil {
		t.Fatal("body is nil")
	}
}

func TestParseEndian(t *testing.T) {
	input := `big_endian;
u32 test @ 0x00;`

	prog := parse(t, input)
	found := false
	for _, n := range prog.Nodes {
		if _, ok := n.(*ast.EndianDirective); ok {
			found = true
			break
		}
	}
	if !found {
		t.Error("endian directive not found in AST")
	}
}

func TestParseIfElse(t *testing.T) {
	input := `fn test(u32 x) {
    if (x == 0) {
        return 1;
    } else {
        return x;
    }
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParseMatch(t *testing.T) {
	input := `fn classify(u32 val) {
    match (val) {
        (0x50): return 1;
        (0x60 ... 0x6F): return 2;
        (_): return 0;
    }
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParseWhile(t *testing.T) {
	input := `fn loop(u32 n) {
    u32 i = 0;
    while (i < n) {
        i = i + 1;
    }
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParsePatternInstrGenBind(t *testing.T) {
	input := `arch x86_64;
platform linux, darwin;

pattern go_memmove {
    name: "runtime.memmove";
    library: "go-runtime";
    version: ">=1.0";

    instr match_move {
        MOVQ src, dst
        MOVQ len, CX
        REP; MOVSQ
    @done:
        RET
    }

    gen {
        memmove($dst, $src, $len);
    }

    bind {
        src as "source";
        dst as "dest";
        len as "count";
    }
}`

	prog := parse(t, input)
	if len(prog.Patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(prog.Patterns))
	}
	pat := prog.Patterns[0]
	if pat.Name != "runtime.memmove" {
		t.Errorf("pattern name = %q", pat.Name)
	}
	if len(pat.InstrBlocks) != 1 {
		t.Errorf("expected 1 instr block, got %d", len(pat.InstrBlocks))
	}
	if pat.GenBlock == nil {
		t.Error("gen block is nil")
	}
	if pat.BindBlock == nil {
		t.Error("bind block is nil")
	}
	if len(pat.BindBlock.Bindings) != 3 {
		t.Errorf("expected 3 bindings, got %d", len(pat.BindBlock.Bindings))
	}
}

func TestParseArchPlatform(t *testing.T) {
	input := `arch x86_64;
platform linux, darwin;`

	prog := parse(t, input)
	foundArch := false
	foundPlatform := false
	for _, n := range prog.Nodes {
		if _, ok := n.(*ast.ArchDirective); ok {
			foundArch = true
		}
		if _, ok := n.(*ast.PlatformDirective); ok {
			foundPlatform = true
		}
	}
	if !foundArch {
		t.Error("arch directive not found")
	}
	if !foundPlatform {
		t.Error("platform directive not found")
	}
}

func TestParseNamespace(t *testing.T) {
	input := `namespace myns {
    u32 counter;
    fn init() {
        counter = 0;
    }
}`

	prog := parse(t, input)
	if len(prog.Namespaces) != 1 {
		t.Fatalf("expected 1 namespace, got %d", len(prog.Namespaces))
	}
}

func TestParseImport(t *testing.T) {
	input := `import std::string;
import "other.hexpat";`

	prog := parse(t, input)
	if len(prog.Imports) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(prog.Imports))
	}
}

func TestParseBinaryExpressions(t *testing.T) {
	input := `fn calc() {
    return (a + b) * (c - d) / 2;
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParseTernary(t *testing.T) {
	input := `fn pick(u32 x) {
    return x > 0 ? x : 0;
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParseCast(t *testing.T) {
	input := `fn conv(u32 x) {
    return x as float;
}`

	prog := parse(t, input)
	if len(prog.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
}

func TestParseInstrBlockMultipleAlternatives(t *testing.T) {
	input := `pattern test {
    instr match_add {
        ADDQ $imm, reg
        | ADDQ src, dst
        | LEAQ (base)(index*scale), reg
    }
}`

	prog := parse(t, input)
	if len(prog.Patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(prog.Patterns))
	}
	block := prog.Patterns[0].InstrBlocks[0]
	if len(block.Alternatives) != 3 {
		t.Errorf("expected 3 alternatives, got %d", len(block.Alternatives))
	}
}
