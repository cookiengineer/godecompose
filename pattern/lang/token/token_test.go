package token

import "testing"

func TestIsKeyword(t *testing.T) {
	validKeywords := []string{
		"struct", "union", "enum", "bitfield", "using",
		"u8", "u16", "u32", "u64", "s8", "s16", "s32", "s64",
		"char", "bool", "float", "double", "str", "auto", "padding",
		"if", "else", "while", "for", "match", "return", "break", "continue", "try", "catch",
		"fn", "namespace", "import", "const", "in", "out", "reference",
		"true", "false", "null", "parent", "this", "as", "is", "from",
		"little_endian", "big_endian", "signed", "unsigned",
		"sizeof", "addressof", "typenameof",
		"instr", "gen", "bind", "pattern", "arch", "platform",
	}

	for _, kw := range validKeywords {
		if !IsKeyword(kw) {
			t.Errorf("IsKeyword(%q) = false, want true", kw)
		}
	}

	notKeywords := []string{"hello", "myFunc", "foo_bar", "x86_64", "RAX", "MOVQ"}
	for _, kw := range notKeywords {
		if IsKeyword(kw) {
			t.Errorf("IsKeyword(%q) = true, want false", kw)
		}
	}
}

func TestLookupKeyword(t *testing.T) {
	tests := []struct {
		ident string
		want  bool
	}{
		{"struct", true},
		{"u8", true},
		{"match", true},
		{"instr", true},
		{"gen", true},
		{"pattern", true},
		{"arch", true},
		{"platform", true},
		{"myVariable", false},
		{"", false},
	}

	for _, tt := range tests {
		_, ok := LookupKeyword(tt.ident)
		if ok != tt.want {
			t.Errorf("LookupKeyword(%q) ok=%v, want %v", tt.ident, ok, tt.want)
		}
	}
}

func TestTokenTypeIsLiteral(t *testing.T) {
	literals := []Type{Integer, Float, String, Char}
	for _, typ := range literals {
		if !typ.IsLiteral() {
			t.Errorf("%s.IsLiteral() = false, want true", typ)
		}
	}

	nonLiterals := []Type{Identifier, LParen, Assign, Keyword, EOF, Illegal}
	for _, typ := range nonLiterals {
		if typ.IsLiteral() {
			t.Errorf("%s.IsLiteral() = true, want false", typ)
		}
	}
}

func TestTokenTypeIsSeparator(t *testing.T) {
	separators := []Type{LParen, RParen, LBrace, RBrace, LBracket, RBracket, Comma, Dot, Semicolon, Colon, At}
	for _, typ := range separators {
		if !typ.IsSeparator() {
			t.Errorf("%s.IsSeparator() = false, want true", typ)
		}
	}
	if Identifier.IsSeparator() {
		t.Error("Identifier.IsSeparator() = true, want false")
	}
}

func TestTokenTypeIsOperator(t *testing.T) {
	operators := []Type{Assign, Plus, Minus, Asterisk, Slash, Percent, LShift, RShift, Equal, NotEqual, Scope}
	for _, typ := range operators {
		if !typ.IsOperator() {
			t.Errorf("%s.IsOperator() = false, want true", typ)
		}
	}
	if Identifier.IsOperator() {
		t.Error("Identifier.IsOperator() = true, want false")
	}
}

func TestTokenString(t *testing.T) {
	tok := Token{Type: Identifier, Literal: "myVar", Line: 5, Column: 10, File: "test.hexpat"}
	if tok.String() != "myVar" {
		t.Errorf("Token.String() = %q, want %q", tok.String(), "myVar")
	}

	if tok.Position() != "test.hexpat" {
		t.Errorf("Token.Position() = %q, want %q", tok.Position(), "test.hexpat")
	}

	intTok := Token{Type: Integer, Literal: "42", Line: 1, Column: 1}
	if intTok.String() != "42" {
		t.Errorf("Integer Token.String() = %q, want %q", intTok.String(), "42")
	}

	opTok := Token{Type: Equal, Literal: "==", Line: 1, Column: 1}
	if opTok.String() != "==" {
		t.Errorf("Operator Token.String() = %q, want %q", opTok.String(), "==")
	}
}
