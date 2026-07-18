package evaluator

import (
	"testing"

	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/parser"
)

func eval(t *testing.T, input string) (*Evaluator, []*CompiledPattern) {
	t.Helper()
	l := lexer.New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	e := New()
	patterns, err := e.Evaluate(prog)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	return e, patterns
}

func TestEvalSimplePattern(t *testing.T) {
	input := `arch x86_64;
platform linux;

pattern go_memmove {
    name: "runtime.memmove";
    
    instr match_move {
        MOVQ src, dst
        MOVQ len, CX
        REP; MOVSQ
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

	_, patterns := eval(t, input)
	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(patterns))
	}

	p := patterns[0]
	if p.Name != "runtime.memmove" {
		t.Errorf("name = %q, want runtime.memmove", p.Name)
	}
	if p.Library != "" {
		t.Logf("library = %q (expected empty for simple test)", p.Library)
	}
	if len(p.Alternatives) != 1 {
		t.Fatalf("expected 1 alternative, got %d", len(p.Alternatives))
	}
	if len(p.Alternatives[0]) != 4 {
		t.Errorf("expected 4 instructions (MOVQ, MOVQ, REP, RET), got %d", len(p.Alternatives[0]))
	}
	if len(p.Bindings) != 3 {
		t.Errorf("expected 3 bindings, got %d", len(p.Bindings))
	}
	if p.GenTemplate == "" {
		t.Error("gen template is empty")
	}

	t.Logf("GenTemplate: %q", p.GenTemplate)
	for i, b := range p.Bindings {
		t.Logf("  binding[%d]: %s -> %s", i, b.CaptureVar, b.Alias)
	}
}

func TestEvalPatternMultipleAlternatives(t *testing.T) {
	input := `pattern test {
    instr match {
        ADDQ $imm, reg
        | MOVQ src, dst
        | LEAQ (base)(index*scale), reg
    }
}`

	_, patterns := eval(t, input)
	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(patterns))
	}

	p := patterns[0]
	if len(p.Alternatives) != 3 {
		t.Fatalf("expected 3 alternatives, got %d", len(p.Alternatives))
	}

	alt0 := p.Alternatives[0][0]
	if alt0.Opcode != "ADDQ" {
		t.Errorf("alt0 opcode = %q, want ADDQ", alt0.Opcode)
	}
	if len(alt0.Operands) < 2 {
		t.Fatalf("expected 2 operands, got %d", len(alt0.Operands))
	}
	if !alt0.Operands[0].IsImmediate {
		t.Error("first operand should be immediate")
	}
	if alt0.Operands[0].CaptureVar != "imm" {
		t.Errorf("imm capture var = %q", alt0.Operands[0].CaptureVar)
	}

	alt2 := p.Alternatives[2][0]
	if alt2.Opcode != "LEAQ" {
		t.Errorf("alt2 opcode = %q, want LEAQ", alt2.Opcode)
	}
	if len(alt2.Operands) < 1 {
		t.Fatal("expected at least 1 operand")
	}
	if alt2.Operands[0].BaseReg != "base" {
		t.Errorf("base reg = %q, want base", alt2.Operands[0].BaseReg)
	}
	if alt2.Operands[0].IndexReg != "index" {
		t.Errorf("index reg = %q, want index", alt2.Operands[0].IndexReg)
	}
	if alt2.Operands[0].Scale != "scale" {
		t.Errorf("scale = %q, want scale", alt2.Operands[0].Scale)
	}
}

func TestEvalGenTemplateExpansion(t *testing.T) {
	input := `pattern test {
    instr match {
        CALL func
    }
    gen {
        call($func);
    }
    bind {
        func as "printf";
    }
}`

	_, patterns := eval(t, input)
	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern")
	}

	if patterns[0].GenTemplate != "call(printf);" {
		t.Errorf("GenTemplate = %q, want \"call(printf);\"", patterns[0].GenTemplate)
	}
}

func TestEvalGenTemplateWithDollar(t *testing.T) {
	input := `pattern test {
    instr match {
        MOVQ src, dst
    }
    gen {
        memmove($dst, $src);
    }
    bind {
        src as "src_ptr";
        dst as "dst_ptr";
    }
}`

	_, patterns := eval(t, input)
	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern")
	}

	expected := "memmove(dst_ptr,src_ptr);"
	if patterns[0].GenTemplate != expected {
		t.Errorf("GenTemplate = %q, want %q", patterns[0].GenTemplate, expected)
	}
}
