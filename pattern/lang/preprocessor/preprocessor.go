// Package preprocessor handles preprocessor directives in the pattern
// language: #include, #define, #undef, #ifdef/#ifndef/#endif, #pragma, #error.
// It operates on token streams from the lexer before parsing.
package preprocessor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

// IncludeResolver resolves #include paths to source text.
type IncludeResolver interface {
	Resolve(path string) (string, error)
}

// FileResolver resolves includes from the filesystem relative to a base directory.
type FileResolver struct {
	BaseDir string
}

func (r *FileResolver) Resolve(path string) (string, error) {
	fullPath := filepath.Join(r.BaseDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("include %s: %w", path, err)
	}
	return string(data), nil
}

// Preprocessor processes directives and macros in a token stream.
type Preprocessor struct {
	resolver IncludeResolver
	macros   map[string][]token.Token
	ifStack  []bool
	source   []token.Token
	pos      int
	output   []token.Token
}

// New creates a preprocessor with an optional include resolver.
func New(resolver IncludeResolver) *Preprocessor {
	return &Preprocessor{
		resolver: resolver,
		macros:   make(map[string][]token.Token),
	}
}

// Process runs the preprocessor over a token stream and returns the
// transformed token stream ready for parsing.
func (p *Preprocessor) Process(tokens []token.Token) ([]token.Token, error) {
	p.source = tokens
	p.pos = 0
	p.output = nil
	p.ifStack = []bool{true}
	p.macros = make(map[string][]token.Token)

	return p.processTokens()
}

func (p *Preprocessor) processTokens() ([]token.Token, error) {
	for p.pos < len(p.source) {
		tok := p.source[p.pos]

		if tok.Type == token.EOF {
			break
		}

		if tok.Type == token.Directive {
			if err := p.processDirective(tok); err != nil {
				return p.output, err
			}
			continue
		}

		if !p.shouldEmit() {
			p.pos++
			continue
		}

		if tok.Type == token.Identifier {
			if replacement, ok := p.macros[tok.Literal]; ok {
				p.pos++
				p.output = append(p.output, replacement...)
				continue
			}
		}

		p.output = append(p.output, tok)
		p.pos++
	}

	p.output = append(p.output, token.Token{
		Type:    token.EOF,
		Line:    p.lastLine(),
		Literal: "",
	})

	return p.output, nil
}

func (p *Preprocessor) processDirective(tok token.Token) error {
	p.pos++
	directive := strings.TrimSpace(tok.Literal)
	parts := strings.Fields(directive)

	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]
	cmd = strings.TrimPrefix(cmd, "#")
	args := strings.TrimSpace(strings.TrimPrefix(directive, parts[0]))

	switch cmd {
	case "include":
		return p.handleInclude(args)
	case "define":
		return p.handleDefine(args)
	case "undef":
		return p.handleUndef(args)
	case "ifdef":
		return p.handleIfdef(args)
	case "ifndef":
		return p.handleIfndef(args)
	case "endif":
		return p.handleEndif()
	case "pragma":
		return nil
	case "error":
		msg := strings.TrimSpace(args)
		msg = strings.Trim(msg, "\"")
		return fmt.Errorf("%s:%d: #error %s", tok.File, tok.Line, msg)
	default:
		return nil
	}
}

func (p *Preprocessor) handleInclude(args string) error {
	path := strings.TrimSpace(args)
	path = strings.Trim(path, "\"")

	if p.resolver == nil {
		return nil
	}

	source, err := p.resolver.Resolve(path)
	if err != nil {
		return fmt.Errorf("#include %s: %w", path, err)
	}

	l := lexer.NewWithFile(source, path)
	includedTokens, err := l.Lex()
	if err != nil {
		return fmt.Errorf("#include %s: %w", path, err)
	}

	// Remove the EOF token from included tokens
	if len(includedTokens) > 0 && includedTokens[len(includedTokens)-1].Type == token.EOF {
		includedTokens = includedTokens[:len(includedTokens)-1]
	}

	// Insert included tokens at current position
	tail := make([]token.Token, len(p.source)-p.pos)
	copy(tail, p.source[p.pos:])
	p.source = append(p.source[:p.pos], includedTokens...)
	p.source = append(p.source, tail...)

	return nil
}

func (p *Preprocessor) handleDefine(args string) error {
	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	name := parts[0]
	if name == "" {
		return nil
	}

	var replacement string
	if len(parts) > 1 {
		replacement = strings.TrimSpace(parts[1])
	}

	if replacement == "" {
		p.macros[name] = nil
		return nil
	}

	l := lexer.New(replacement)
	tokens, err := l.Lex()
	if err != nil {
		return nil
	}
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == token.EOF {
		tokens = tokens[:len(tokens)-1]
	}

	p.macros[name] = tokens
	return nil
}

func (p *Preprocessor) handleUndef(args string) error {
	name := strings.TrimSpace(args)
	delete(p.macros, name)
	return nil
}

func (p *Preprocessor) handleIfdef(args string) error {
	name := strings.TrimSpace(args)
	_, defined := p.macros[name]
	p.pushConditional(defined)
	return nil
}

func (p *Preprocessor) handleIfndef(args string) error {
	name := strings.TrimSpace(args)
	_, defined := p.macros[name]
	p.pushConditional(!defined)
	return nil
}

func (p *Preprocessor) handleEndif() error {
	if len(p.ifStack) <= 1 {
		return fmt.Errorf("#endif without matching #ifdef/#ifndef")
	}
	p.ifStack = p.ifStack[:len(p.ifStack)-1]
	return nil
}

func (p *Preprocessor) pushConditional(emit bool) {
	currentEmit := p.ifStack[len(p.ifStack)-1]
	p.ifStack = append(p.ifStack, currentEmit && emit)
}

func (p *Preprocessor) shouldEmit() bool {
	if len(p.ifStack) == 0 {
		return true
	}
	return p.ifStack[len(p.ifStack)-1]
}

func (p *Preprocessor) lastLine() int {
	if len(p.output) > 0 {
		return p.output[len(p.output)-1].Line
	}
	return 1
}
