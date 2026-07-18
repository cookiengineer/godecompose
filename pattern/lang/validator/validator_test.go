package validator

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/parser"
)

func validate(t *testing.T, input string) []error {
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

	v := New()
	return v.Validate(prog)
}

func assertNoErrors(t *testing.T, errs []error) {
	t.Helper()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func assertErrorContains(t *testing.T, errs []error, substr string) {
	t.Helper()
	for _, err := range errs {
		if strings.Contains(err.Error(), substr) {
			return
		}
	}
	t.Errorf("expected error containing %q, got %v", substr, errs)
}

func TestValidateStruct(t *testing.T) {
	input := `struct Header {
    u32 magic;
    u32 version;
};`
	errs := validate(t, input)
	assertNoErrors(t, errs)
}

func TestValidateFunction(t *testing.T) {
	input := `fn add(u32 a, u32 b) {
    return a + b;
}`
	errs := validate(t, input)
	assertNoErrors(t, errs)
}

func TestValidateBreakOutsideLoop(t *testing.T) {
	input := `fn test() {
    break;
}`
	errs := validate(t, input)
	assertErrorContains(t, errs, "break outside loop")
}

func TestValidateContinueOutsideLoop(t *testing.T) {
	input := `fn test() {
    continue;
}`
	errs := validate(t, input)
	assertErrorContains(t, errs, "continue outside loop")
}

func TestValidateBreakInLoop(t *testing.T) {
	input := `fn test() {
    while (true) {
        break;
    }
}`
	errs := validate(t, input)
	assertNoErrors(t, errs)
}

func TestValidateNestedLoop(t *testing.T) {
	input := `fn test() {
    while (true) {
        while (true) {
            break;
        }
        break;
    }
}`
	errs := validate(t, input)
	assertNoErrors(t, errs)
}

func TestValidatePattern(t *testing.T) {
	input := `arch x86_64;
pattern test {
    instr match {
        MOVQ src, dst
    }
}`
	errs := validate(t, input)
	assertNoErrors(t, errs)
}
