package lexer

import (
	"testing"

	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

func TestLexEmpty(t *testing.T) {
	l := New("")
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	if len(tokens) != 1 || tokens[0].Type != token.EOF {
		t.Errorf("expected single EOF token, got %d tokens: %v", len(tokens), tokens)
	}
}

func TestLexKeywords(t *testing.T) {
	keywords := []string{
		"struct", "union", "enum", "bitfield", "using",
		"u8", "u16", "u32", "u64", "u128",
		"s8", "s16", "s32", "s64",
		"char", "bool", "float", "double", "str", "auto", "padding",
		"if", "else", "while", "for", "match", "return", "break", "continue", "try", "catch",
		"fn", "namespace", "import", "const", "in", "out", "reference",
		"true", "false", "null", "parent", "this", "as", "is", "from",
		"little_endian", "big_endian", "signed", "unsigned",
		"sizeof", "addressof", "typenameof",
		"instr", "gen", "bind", "pattern", "arch", "platform",
	}

	for _, kw := range keywords {
		l := New(kw)
		tokens, err := l.Lex()
		if err != nil {
			t.Errorf("Lex(%q): %v", kw, err)
			continue
		}
		if len(tokens) < 2 {
			t.Errorf("Lex(%q): expected 2 tokens, got %d", kw, len(tokens))
			continue
		}
		if tokens[0].Type != token.Keyword || tokens[0].Literal != kw {
			t.Errorf("Lex(%q): got %v %q, want KEYWORD %q", kw, tokens[0].Type, tokens[0].Literal, kw)
		}
	}
}

func TestLexIdentifiers(t *testing.T) {
	tests := []string{"myVar", "_private", "test123", "msg_and_data", "RAX", "MOVQ", "x86_64"}
	for _, id := range tests {
		l := New(id)
		tokens, err := l.Lex()
		if err != nil {
			t.Errorf("Lex(%q): %v", id, err)
			continue
		}
		if tokens[0].Type != token.Identifier || tokens[0].Literal != id {
			t.Errorf("Lex(%q): got %v %q, want IDENT %q", id, tokens[0].Type, tokens[0].Literal, id)
		}
	}
}

func TestLexIntegers(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0", "0"},
		{"42", "42"},
		{"0xFF", "0xFF"},
		{"0xDEADBEEF", "0xDEADBEEF"},
		{"0o777", "0o777"},
		{"0b1010", "0b1010"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens, err := l.Lex()
		if err != nil {
			t.Errorf("Lex(%q): %v", tt.input, err)
			continue
		}
		if tokens[0].Type != token.Integer || tokens[0].Literal != tt.want {
			t.Errorf("Lex(%q): got %v %q, want INTEGER %q", tt.input, tokens[0].Type, tokens[0].Literal, tt.want)
		}
	}
}

func TestLexFloats(t *testing.T) {
	tests := []string{"3.14", "1.0", "0.5", "1e10", "2.5e-3", "6.022e23"}
	for _, input := range tests {
		l := New(input)
		tokens, err := l.Lex()
		if err != nil {
			t.Errorf("Lex(%q): %v", input, err)
			continue
		}
		if tokens[0].Type != token.Float {
			t.Errorf("Lex(%q): got %v, want FLOAT", input, tokens[0].Type)
		}
	}
}

func TestLexStrings(t *testing.T) {
	l := New(`"hello world" "escaped\nstring\there"`)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	if tokens[0].Literal != "hello world" {
		t.Errorf("string 1 = %q, want %q", tokens[0].Literal, "hello world")
	}
	if tokens[1].Literal != "escaped\nstring\there" {
		t.Errorf("string 2 = %q, want %q", tokens[1].Literal, "escaped\nstring\there")
	}
}

func TestLexChar(t *testing.T) {
	l := New(`'a' '\n' '\t' '\\'`)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	if len(tokens) < 5 { // 4 chars + EOF
		t.Fatalf("expected 5 tokens, got %d", len(tokens))
	}
	expected := []string{"a", "\n", "\t", "\\"}
	for i, exp := range expected {
		if tokens[i].Literal != exp {
			t.Errorf("char[%d] = %q, want %q", i, tokens[i].Literal, exp)
		}
	}
}

func TestLexOperators(t *testing.T) {
	tests := []struct {
		input string
		typ   token.Type
	}{
		{"+", token.Plus},
		{"-", token.Minus},
		{"*", token.Asterisk},
		{"/", token.Slash},
		{"%", token.Percent},
		{"=", token.Assign},
		{"+=", token.PlusAssign},
		{"-=", token.MinusAssign},
		{"*=", token.StarAssign},
		{"/=", token.SlashAssign},
		{"%=", token.PctAssign},
		{"&=", token.AmpAssign},
		{"|=", token.PipeAssign},
		{"^=", token.CaretAssign},
		{"<<", token.LShift},
		{">>", token.RShift},
		{"<<=", token.LShiftAssign},
		{">>=", token.RShiftAssign},
		{"<", token.Less},
		{">", token.Greater},
		{"<=", token.LEqual},
		{">=", token.GEqual},
		{"==", token.Equal},
		{"!=", token.NotEqual},
		{"&&", token.And},
		{"||", token.Or},
		{"^^", token.Xor},
		{"::", token.Scope},
		{"->", token.Arrow},
		{"...", token.Range},
		{"!", token.Exclamation},
		{"?", token.Question},
		{"$", token.Dollar},
		{"&", token.Ampersand},
		{"|", token.Pipe},
		{"^", token.Caret},
		{"~", token.Tilde},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens, err := l.Lex()
		if err != nil {
			t.Errorf("Lex(%q): %v", tt.input, err)
			continue
		}
		if tokens[0].Type != tt.typ {
			t.Errorf("Lex(%q): got %v, want %v", tt.input, tokens[0].Type, tt.typ)
		}
	}
}

func TestLexSeparators(t *testing.T) {
	l := New("(){}[],.;:@")
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	expectedTypes := []token.Type{
		token.LParen, token.RParen,
		token.LBrace, token.RBrace,
		token.LBracket, token.RBracket,
		token.Comma, token.Dot, token.Semicolon, token.Colon, token.At,
		token.EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, typ := range expectedTypes {
		if tokens[i].Type != typ {
			t.Errorf("token[%d] = %v, want %v", i, tokens[i].Type, typ)
		}
	}
}

func TestLexComments(t *testing.T) {
	input := `// line comment
42 /* block comment */ 3.14 /// doc comment
/** doc block */ "string"`
	l := New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	if len(tokens) != 4 { // 42, 3.14, "string", EOF
		t.Fatalf("expected 4 tokens, got %d: %v", len(tokens), tokens)
	}

	if tokens[0].Type != token.Integer || tokens[0].Literal != "42" {
		t.Error("integer token missing/mismatched")
	}
	if tokens[1].Type != token.Float || tokens[1].Literal != "3.14" {
		t.Error("float token missing/mismatched")
	}
	if tokens[2].Type != token.String || tokens[2].Literal != "string" {
		t.Error("string token missing/mismatched")
	}
}

func TestLexDirectives(t *testing.T) {
	input := `#include "file.hexpat"
#define FOO 42
#ifdef BAR
#error "not supported"
#endif
42`
	l := New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	directiveCount := 0
	for _, tok := range tokens {
		if tok.Type == token.Directive {
			directiveCount++
		}
	}
	if directiveCount != 5 {
		t.Errorf("expected 5 directives, got %d", directiveCount)
	}

	lastToken := tokens[len(tokens)-2]
	if lastToken.Type != token.Integer || lastToken.Literal != "42" {
		t.Errorf("last substantive token = %v %q, want INTEGER 42", lastToken.Type, lastToken.Literal)
	}
}

func TestLexNestedBlockComments(t *testing.T) {
	input := "1 /* outer /* inner */ */ 2"
	l := New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	if len(tokens) != 3 { // 1, 2, EOF
		t.Fatalf("expected 3 tokens, got %d: %v", len(tokens), tokens)
	}
	if tokens[0].Literal != "1" || tokens[1].Literal != "2" {
		t.Errorf("got %q %q, want '1' '2'", tokens[0].Literal, tokens[1].Literal)
	}
}

func TestLexLineTracking(t *testing.T) {
	input := "a\nb\nc"
	l := New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	if tokens[0].Line != 1 {
		t.Errorf("a line = %d, want 1", tokens[0].Line)
	}
	if tokens[1].Line != 2 {
		t.Errorf("b line = %d, want 2", tokens[1].Line)
	}
	if tokens[2].Line != 3 {
		t.Errorf("c line = %d, want 3", tokens[2].Line)
	}
}

func TestLexPatternExample(t *testing.T) {
	input := `arch x86_64;
platform linux, darwin;

pattern go_memmove {
    name: "runtime.memmove";
    
    instr match_copy {
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

	l := New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("no tokens produced from valid pattern")
	}

	foundInstr := false
	foundGen := false
	foundBind := false
	for _, tok := range tokens {
		if tok.Type == token.Keyword {
			switch tok.Literal {
			case "instr":
				foundInstr = true
			case "gen":
				foundGen = true
			case "bind":
				foundBind = true
			}
		}
	}

	if !foundInstr {
		t.Error("instr keyword not found in pattern")
	}
	if !foundGen {
		t.Error("gen keyword not found in pattern")
	}
	if !foundBind {
		t.Error("bind keyword not found in pattern")
	}
}

func TestLexErrorUnterminatedString(t *testing.T) {
	l := New(`"unterminated`)
	_, err := l.Lex()
	if err == nil {
		t.Error("expected error for unterminated string")
	}
}

func TestLexErrorUnterminatedChar(t *testing.T) {
	l := New(`'a`)
	_, err := l.Lex()
	if err == nil {
		t.Error("expected error for unterminated char")
	}
}

func TestLexErrorUnexpectedChar(t *testing.T) {
	l := New("`")
	_, err := l.Lex()
	if err == nil {
		t.Error("expected error for unexpected character")
	}
}
