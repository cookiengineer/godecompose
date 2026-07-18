package preprocessor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

func lexTokens(t *testing.T, input string) []token.Token {
	t.Helper()
	l := lexer.New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	return tokens
}

func TestDefine(t *testing.T) {
	input := `#define FOO 42
u32 x = FOO;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found42 := false
	for _, tok := range result {
		if tok.Literal == "42" {
			found42 = true
			break
		}
	}
	if !found42 {
		t.Error("FOO was not expanded to 42")
	}
}

func TestDefineNoValue(t *testing.T) {
	input := `#define DEBUG
#ifdef DEBUG
u32 x = 1;
#endif
u32 y = 2;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found1 := false
	found2 := false
	for _, tok := range result {
		if tok.Literal == "1" {
			found1 = true
		}
		if tok.Literal == "2" {
			found2 = true
		}
	}
	if !found1 {
		t.Error("DEBUG block was not emitted")
	}
	if !found2 {
		t.Error("outside block not emitted")
	}
}

func TestIfdefFalse(t *testing.T) {
	input := `#ifdef MISSING
u32 x = 1;
#endif
u32 y = 2;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	for _, tok := range result {
		if tok.Literal == "1" {
			t.Error("MISSING block was emitted but should not be")
		}
	}

	found2 := false
	for _, tok := range result {
		if tok.Literal == "2" {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Error("y=2 not found")
	}
}

func TestIfndefTrue(t *testing.T) {
	input := `#ifndef MISSING
u32 x = 1;
#endif`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found1 := false
	for _, tok := range result {
		if tok.Literal == "1" {
			found1 = true
			break
		}
	}
	if !found1 {
		t.Error("block not emitted for undefined macro")
	}
}

func TestUndef(t *testing.T) {
	input := `#define FOO 42
#undef FOO
#ifdef FOO
u32 x = 1;
#endif
u32 x = 2;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	for _, tok := range result {
		if tok.Literal == "1" {
			t.Error("FOO block emitted after undef")
		}
	}

	found2 := false
	for _, tok := range result {
		if tok.Literal == "2" {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Error("x=2 not found")
	}
}

func TestNestedIfdef(t *testing.T) {
	input := `#define OUTER
#ifdef OUTER
#define INNER
#ifdef INNER
u32 x = 1;
#endif
#endif
u32 y = 2;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found1 := false
	found2 := false
	for _, tok := range result {
		if tok.Literal == "1" {
			found1 = true
		}
		if tok.Literal == "2" {
			found2 = true
		}
	}
	if !found1 {
		t.Error("nested block not emitted")
	}
	if !found2 {
		t.Error("outer block not emitted")
	}
}

func TestPragmaIgnored(t *testing.T) {
	input := `#pragma once
u32 x = 1;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found := false
	for _, tok := range result {
		if tok.Literal == "1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("pragma prevented emission")
	}
}

func TestErrorDirective(t *testing.T) {
	input := `#error "not supported"`

	tokens := lexTokens(t, input)
	p := New(nil)
	_, err := p.Process(tokens)
	if err == nil {
		t.Error("expected error from #error directive")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("error message = %q, want containing 'not supported'", err.Error())
	}
}

func TestFileResolver(t *testing.T) {
	dir := t.TempDir()
	resolver := &FileResolver{BaseDir: dir}

	included := `u32 included_var @ 0x00;`
	path := filepath.Join(dir, "common.hexpat")
	if err := os.WriteFile(path, []byte(included), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	input := `#include "common.hexpat"
u32 main_var @ 0x04;`

	tokens := lexTokens(t, input)
	p := New(resolver)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	foundIncluded := false
	foundMain := false
	for _, tok := range result {
		if tok.Literal == "included_var" {
			foundIncluded = true
		}
		if tok.Literal == "main_var" {
			foundMain = true
		}
	}
	if !foundIncluded {
		t.Error("included_var not found - include not processed")
	}
	if !foundMain {
		t.Error("main_var not found")
	}
}

func TestMacroWithIdentifierExpansion(t *testing.T) {
	input := `#define SIZE 256
u8 data[SIZE] @ 0x00;`

	tokens := lexTokens(t, input)
	p := New(nil)
	result, err := p.Process(tokens)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	found256 := false
	for _, tok := range result {
		if tok.Type == token.Integer && tok.Literal == "256" {
			found256 = true
			break
		}
	}
	if !found256 {
		t.Error("SIZE macro not expanded")
	}
}
