// Package lexer implements a hand-written lexer for the ImHex-compatible
// Pattern Language, extended with godecompose constructs.
package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

type Lexer struct {
	source  []byte
	pos     int
	line    int
	column  int
	file    string
	tokens  []token.Token
}

func New(source string) *Lexer {
	return NewWithFile(source, "<input>")
}

func NewWithFile(source, file string) *Lexer {
	return &Lexer{
		source: []byte(source),
		pos:    0,
		line:   1,
		column: 1,
		file:   file,
	}
}

func (l *Lexer) Lex() ([]token.Token, error) {
	l.tokens = nil
	l.pos = 0
	l.line = 1
	l.column = 1

	for l.pos < len(l.source) {
		ch := l.current()

		if unicode.IsSpace(rune(ch)) {
			l.skipWhitespace()
			continue
		}

		if ch == '#' && l.column == 1 {
			tok := l.scanDirective()
			l.tokens = append(l.tokens, tok)
			continue
		}

		if ch == '/' {
			if l.peek() == '/' || l.peek() == '*' {
				l.scanComment()
				continue
			}
		}

		if isDigit(ch) {
			tok := l.scanNumber()
			l.tokens = append(l.tokens, tok)
			continue
		}

		if isLetter(ch) || ch == '_' {
			tok := l.scanIdentifier()
			l.tokens = append(l.tokens, tok)
			continue
		}

		if ch == '"' {
			tok, err := l.scanString('"')
			if err != nil {
				return nil, err
			}
			l.tokens = append(l.tokens, tok)
			continue
		}

		if ch == '\'' {
			tok, err := l.scanChar()
			if err != nil {
				return nil, err
			}
			l.tokens = append(l.tokens, tok)
			continue
		}

		tok := l.scanOperator()
		if tok.Type != token.Illegal {
			l.tokens = append(l.tokens, tok)
			continue
		}

		tok = l.scanSeparator()
		if tok.Type != token.Illegal {
			l.tokens = append(l.tokens, tok)
			continue
		}

		return nil, fmt.Errorf("%s:%d:%d: unexpected character %q", l.file, l.line, l.column, ch)
	}

	l.tokens = append(l.tokens, token.Token{
		Type: token.EOF, Line: l.line, Column: l.column, File: l.file,
	})

	return l.tokens, nil
}

func (l *Lexer) current() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) peek() byte {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.pos+1]
}

func (l *Lexer) advance() byte {
	ch := l.current()
	l.pos++
	l.column++
	if ch == '\n' {
		l.line++
		l.column = 1
	}
	return ch
}

func (l *Lexer) backup() {
	if l.pos > 0 {
		l.pos--
		l.column--
		if l.pos > 0 && l.source[l.pos-1] == '\n' {
			l.line--
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.source) {
		ch := l.current()
		if ch == '\n' {
			l.line++
			l.column = 1
			l.pos++
		} else if unicode.IsSpace(rune(ch)) {
			l.pos++
			l.column++
		} else {
			break
		}
	}
}

func (l *Lexer) readWhile(pred func(byte) bool) string {
	start := l.pos
	for l.pos < len(l.source) && pred(l.current()) {
		l.advance()
	}
	return string(l.source[start:l.pos])
}

func (l *Lexer) scanIdentifier() token.Token {
	startLine := l.line
	startCol := l.column

	ident := l.readWhile(func(ch byte) bool {
		return isLetter(ch) || isDigit(ch) || ch == '_'
	})

	if len(ident) == 0 {
		ident = string(l.advance())
	}

	// Check if it's a keyword
	if _, isKw := token.LookupKeyword(ident); isKw {
		return token.Token{
			Type:    token.Keyword,
			Literal: ident,
			Line:    startLine,
			Column:  startCol,
			File:    l.file,
		}
	}

	return token.Token{
		Type:    token.Identifier,
		Literal: ident,
		Line:    startLine,
		Column:  startCol,
		File:    l.file,
	}
}

func (l *Lexer) scanNumber() token.Token {
	startLine := l.line
	startCol := l.column
	isFloat := false

	if l.current() == '0' {
		if l.peek() == 'x' || l.peek() == 'X' {
			l.advance() // 0
			l.advance() // x
			digits := l.readWhile(isHexDigit)
			return token.Token{
				Type:    token.Integer,
				Literal: "0x" + digits,
				Line:    startLine,
				Column:  startCol,
				File:    l.file,
			}
		}
		if l.peek() == 'o' || l.peek() == 'O' {
			l.advance() // 0
			l.advance() // o
			digits := l.readWhile(isOctalDigit)
			return token.Token{
				Type:    token.Integer,
				Literal: "0o" + digits,
				Line:    startLine,
				Column:  startCol,
				File:    l.file,
			}
		}
		if l.peek() == 'b' || l.peek() == 'B' {
			l.advance() // 0
			l.advance() // b
			digits := l.readWhile(isBinaryDigit)
			return token.Token{
				Type:    token.Integer,
				Literal: "0b" + digits,
				Line:    startLine,
				Column:  startCol,
				File:    l.file,
			}
		}
	}

	digits := l.readWhile(func(ch byte) bool {
		if ch == '\'' {
			return true
		}
		if ch == '.' && (isDigit(l.peek()) || l.peek() == 'e' || l.peek() == 'E') {
			isFloat = true
			return true
		}
		if (ch == 'e' || ch == 'E') && (isDigit(l.peek()) || l.peek() == '+' || l.peek() == '-') {
			isFloat = true
			return true
		}
		if ch == '+' || ch == '-' {
			prevIdx := l.pos - 1
			if prevIdx >= 0 && (l.source[prevIdx] == 'e' || l.source[prevIdx] == 'E') {
				return true
			}
		}
		return isDigit(ch)
	})

	if isFloat {
		return token.Token{
			Type:    token.Float,
			Literal: digits,
			Line:    startLine,
			Column:  startCol,
			File:    l.file,
		}
	}

	return token.Token{
		Type:    token.Integer,
		Literal: digits,
		Line:    startLine,
		Column:  startCol,
		File:    l.file,
	}
}

func (l *Lexer) scanString(quote byte) (token.Token, error) {
	startLine := l.line
	startCol := l.column
	l.advance() // opening quote

	var buf strings.Builder
	for l.pos < len(l.source) {
		ch := l.current()
		if ch == quote {
			l.advance()
			return token.Token{
				Type:    token.String,
				Literal: buf.String(),
				Line:    startLine,
				Column:  startCol,
				File:    l.file,
			}, nil
		}
		if ch == '\\' {
			l.advance()
			if l.pos >= len(l.source) {
				return token.Token{}, fmt.Errorf("%s:%d:%d: unexpected end of string", l.file, l.line, l.column)
			}
			esc := l.advance()
			switch esc {
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			case '\\':
				buf.WriteByte('\\')
			case '"':
				buf.WriteByte('"')
			case '\'':
				buf.WriteByte('\'')
			case '0':
				buf.WriteByte(0)
			default:
				buf.WriteByte('\\')
				buf.WriteByte(esc)
			}
			continue
		}
		if ch == '\n' {
			return token.Token{}, fmt.Errorf("%s:%d:%d: unexpected newline in string", l.file, l.line, l.column)
		}
		buf.WriteByte(ch)
		l.advance()
	}

	return token.Token{}, fmt.Errorf("%s:%d:%d: unterminated string", l.file, startLine, startCol)
}

func (l *Lexer) scanChar() (token.Token, error) {
	startLine := l.line
	startCol := l.column
	l.advance() // opening '

	if l.pos >= len(l.source) {
		return token.Token{}, fmt.Errorf("%s:%d:%d: unexpected end of char literal", l.file, l.line, l.column)
	}

	var ch rune
	if l.current() == '\\' {
		l.advance()
		if l.pos >= len(l.source) {
			return token.Token{}, fmt.Errorf("%s:%d:%d: unexpected end of char literal", l.file, l.line, l.column)
		}
		switch l.current() {
		case 'n':
			ch = '\n'
		case 't':
			ch = '\t'
		case 'r':
			ch = '\r'
		case '\\':
			ch = '\\'
		case '\'':
			ch = '\''
		case '0':
			ch = 0
		default:
			ch = rune(l.current())
		}
		l.advance()
	} else {
		r, size := utf8.DecodeRune(l.source[l.pos:])
		ch = r
		l.pos += size
		l.column += size
	}

	if l.pos >= len(l.source) || l.current() != '\'' {
		return token.Token{}, fmt.Errorf("%s:%d:%d: expected closing ' in char literal", l.file, l.line, l.column)
	}
	l.advance() // closing '

	return token.Token{
		Type:    token.Char,
		Literal: string(ch),
		Line:    startLine,
		Column:  startCol,
		File:    l.file,
	}, nil
}

func (l *Lexer) scanComment() {
	startLine := l.line
	startCol := l.column

	l.advance() // first /
	ch := l.current()

	if ch == '/' {
		// Line comment
		l.readWhile(func(b byte) bool { return b != '\n' })
		return
	}

	if ch == '*' {
		l.advance() // *
		depth := 1
		for l.pos < len(l.source) && depth > 0 {
			c := l.advance()
			if c == '/' && l.current() == '*' {
				l.advance()
				depth++
			} else if c == '*' && l.current() == '/' {
				l.advance()
				depth--
			}
		}
		return
	}

	// Not a comment — put back the slash
	l.pos--
	l.column--
	_ = startCol
	_ = startLine
}

func (l *Lexer) scanOperator() token.Token {
	startLine := l.line
	startCol := l.column

	ch1 := l.advance()
	ch2 := l.current()

	// Three-character operators
	if ch1 == '<' && ch2 == '<' && l.peekNext() == '=' {
		l.advance()
		l.advance()
		return token.Token{Type: token.LShiftAssign, Literal: "<<=", Line: startLine, Column: startCol, File: l.file}
	}
	if ch1 == '>' && ch2 == '>' && l.peekNext() == '=' {
		l.advance()
		l.advance()
		return token.Token{Type: token.RShiftAssign, Literal: ">>=", Line: startLine, Column: startCol, File: l.file}
	}
	if ch1 == '.' && ch2 == '.' && l.peekNext() == '.' {
		l.advance()
		l.advance()
		return token.Token{Type: token.Range, Literal: "...", Line: startLine, Column: startCol, File: l.file}
	}

	// Two-character operators
	if ch2 != 0 {
		combo := string([]byte{ch1, ch2})
		switch combo {
		case "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>", "<=", ">=", "==", "!=", "&&", "||", "^^", "::", "->":
			l.advance()
			return token.Token{Type: opType(combo), Literal: combo, Line: startLine, Column: startCol, File: l.file}
		}
	}

	// Single-character operator with required second char
	if ch1 == '<' && ch2 == '<' {
		l.advance()
		return token.Token{Type: token.LShift, Literal: "<<", Line: startLine, Column: startCol, File: l.file}
	}
	if ch1 == '>' && ch2 == '>' {
		l.advance()
		return token.Token{Type: token.RShift, Literal: ">>", Line: startLine, Column: startCol, File: l.file}
	}

	// Single-character operators
	if typ, ok := singleOpType(ch1); ok {
		return token.Token{Type: typ, Literal: string(ch1), Line: startLine, Column: startCol, File: l.file}
	}

	l.backup()
	return token.Token{Type: token.Illegal}
}

func opType(s string) token.Type {
	switch s {
	case "+=":
		return token.PlusAssign
	case "-=":
		return token.MinusAssign
	case "*=":
		return token.StarAssign
	case "/=":
		return token.SlashAssign
	case "%=":
		return token.PctAssign
	case "&=":
		return token.AmpAssign
	case "|=":
		return token.PipeAssign
	case "^=":
		return token.CaretAssign
	case "<<":
		return token.LShift
	case ">>":
		return token.RShift
	case "<=":
		return token.LEqual
	case ">=":
		return token.GEqual
	case "==":
		return token.Equal
	case "!=":
		return token.NotEqual
	case "&&":
		return token.And
	case "||":
		return token.Or
	case "^^":
		return token.Xor
	case "::":
		return token.Scope
	case "->":
		return token.Arrow
	default:
		return token.Illegal
	}
}

func singleOpType(ch byte) (token.Type, bool) {
	switch ch {
	case '+':
		return token.Plus, true
	case '-':
		return token.Minus, true
	case '*':
		return token.Asterisk, true
	case '/':
		return token.Slash, true
	case '%':
		return token.Percent, true
	case '=':
		return token.Assign, true
	case '<':
		return token.Less, true
	case '>':
		return token.Greater, true
	case '!':
		return token.Exclamation, true
	case '?':
		return token.Question, true
	case '$':
		return token.Dollar, true
	case '&':
		return token.Ampersand, true
	case '|':
		return token.Pipe, true
	case '^':
		return token.Caret, true
	case '~':
		return token.Tilde, true
	default:
		return token.Illegal, false
	}
}

func (l *Lexer) scanSeparator() token.Token {
	startLine := l.line
	startCol := l.column
	ch := l.current()
	if !isSeparator(ch) {
		return token.Token{Type: token.Illegal}
	}
	l.advance()

	switch ch {
	case '(':
		return token.Token{Type: token.LParen, Literal: "(", Line: startLine, Column: startCol, File: l.file}
	case ')':
		return token.Token{Type: token.RParen, Literal: ")", Line: startLine, Column: startCol, File: l.file}
	case '{':
		return token.Token{Type: token.LBrace, Literal: "{", Line: startLine, Column: startCol, File: l.file}
	case '}':
		return token.Token{Type: token.RBrace, Literal: "}", Line: startLine, Column: startCol, File: l.file}
	case '[':
		return token.Token{Type: token.LBracket, Literal: "[", Line: startLine, Column: startCol, File: l.file}
	case ']':
		return token.Token{Type: token.RBracket, Literal: "]", Line: startLine, Column: startCol, File: l.file}
	case ',':
		return token.Token{Type: token.Comma, Literal: ",", Line: startLine, Column: startCol, File: l.file}
	case '.':
		return token.Token{Type: token.Dot, Literal: ".", Line: startLine, Column: startCol, File: l.file}
	case ';':
		return token.Token{Type: token.Semicolon, Literal: ";", Line: startLine, Column: startCol, File: l.file}
	case ':':
		return token.Token{Type: token.Colon, Literal: ":", Line: startLine, Column: startCol, File: l.file}
	case '@':
		return token.Token{Type: token.At, Literal: "@", Line: startLine, Column: startCol, File: l.file}
	}

	return token.Token{Type: token.Illegal}
}

func isSeparator(ch byte) bool {
	switch ch {
	case '(', ')', '{', '}', '[', ']', ',', '.', ';', ':', '@':
		return true
	}
	return false
}

func (l *Lexer) scanDirective() token.Token {
	startLine := l.line
	startCol := l.column

	content := l.readWhile(func(ch byte) bool { return ch != '\n' })
	return token.Token{
		Type:    token.Directive,
		Literal: content,
		Line:    startLine,
		Column:  startCol,
		File:    l.file,
	}
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.pos+1]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func isOctalDigit(ch byte) bool {
	return '0' <= ch && ch <= '7'
}

func isBinaryDigit(ch byte) bool {
	return ch == '0' || ch == '1'
}
