// Package token defines the lexical token types for the ImHex-compatible
// Pattern Language, extended with godecompose decompilation constructs.
package token

// Type categorizes a lexical token.
type Type int

const (
	// Special
	Illegal Type = iota
	EOF
	Comment

	// Literals
	Identifier
	Integer // decimal, hex 0x, octal 0o, binary 0b
	Float
	String // "string"
	Char   // 'c'

	// Separators
	LParen    // (
	RParen    // )
	LBrace    // {
	RBrace    // }
	LBracket  // [
	RBracket  // ]
	Comma     // ,
	Dot       // .
	Semicolon // ;
	Colon     // :
	At        // @

	// Operators
	Assign        // =
	Plus          // +
	Minus         // -
	Asterisk      // *
	Slash         // /
	Percent       // %
	Ampersand     // &
	Pipe          // |
	Caret         // ^
	Tilde         // ~
	Less          // <
	Greater       // >
	Exclamation   // !
	Question      // ?
	Dollar        // $

	// Compound operators
	PlusAssign   // +=
	MinusAssign  // -=
	StarAssign   // *=
	SlashAssign  // /=
	PctAssign    // %=
	AmpAssign    // &=
	PipeAssign   // |=
	CaretAssign  // ^=
	LShift       // <<
	RShift       // >>
	LShiftAssign // <<=
	RShiftAssign // >>=
	LEqual       // <=
	GEqual       // >=
	Equal        // ==
	NotEqual     // !=
	And          // &&
	Or           // ||
	Xor          // ^^
	Scope        // ::
	Range        // ...
	Arrow        // ->

	// Keywords
	Keyword // base for keyword tokens

	// Godecompose extension keywords
	GodecomposeExt // base for godecompose extensions

	// Preprocessor directives
	Directive
)

// IsLiteral reports whether t is a literal token type.
func (t Type) IsLiteral() bool {
	switch t {
	case Integer, Float, String, Char:
		return true
	}
	return false
}

// IsSeparator reports whether t is a separator token type.
func (t Type) IsSeparator() bool {
	return t >= LParen && t <= At
}

// IsOperator reports whether t is an operator token type.
func (t Type) IsOperator() bool {
	return t >= Assign && t <= Arrow
}

func (t Type) String() string {
	if s, ok := typeNames[t]; ok {
		return s
	}
	return "<unknown>"
}

var typeNames = map[Type]string{
	Illegal:    "ILLEGAL",
	EOF:        "EOF",
	Comment:    "COMMENT",
	Identifier: "IDENT",
	Integer:    "INTEGER",
	Float:      "FLOAT",
	String:     "STRING",
	Char:       "CHAR",
	LParen:     "(",
	RParen:     ")",
	LBrace:     "{",
	RBrace:     "}",
	LBracket:   "[",
	RBracket:   "]",
	Comma:      ",",
	Dot:        ".",
	Semicolon:  ";",
	Colon:      ":",
	At:         "@",
	Assign:     "=",
	Plus:       "+",
	Minus:      "-",
	Asterisk:   "*",
	Slash:      "/",
	Percent:    "%",
	Ampersand:  "&",
	Pipe:       "|",
	Caret:      "^",
	Tilde:      "~",
	Less:       "<",
	Greater:    ">",
	Exclamation: "!",
	Question:   "?",
	Dollar:     "$",
	PlusAssign:   "+=",
	MinusAssign:  "-=",
	StarAssign:   "*=",
	SlashAssign:  "/=",
	PctAssign:    "%=",
	AmpAssign:    "&=",
	PipeAssign:   "|=",
	CaretAssign:  "^=",
	LShift:       "<<",
	RShift:       ">>",
	LShiftAssign: "<<=",
	RShiftAssign: ">>=",
	LEqual:       "<=",
	GEqual:       ">=",
	Equal:        "==",
	NotEqual:     "!=",
	And:          "&&",
	Or:           "||",
	Xor:          "^^",
	Scope:        "::",
	Range:        "...",
	Arrow:        "->",
	Keyword:      "KEYWORD",
	Directive:    "DIRECTIVE",
}

// keywordToToken maps keyword strings to their token Type.
var keywordToToken = map[string]Type{
	// Type keywords
	"struct":  Keyword,
	"union":   Keyword,
	"using":   Keyword,
	"enum":    Keyword,
	"bitfield": Keyword,

	// Value types
	"u8":     Keyword,
	"u16":    Keyword,
	"u24":    Keyword,
	"u32":    Keyword,
	"u48":    Keyword,
	"u64":    Keyword,
	"u96":    Keyword,
	"u128":   Keyword,
	"s8":     Keyword,
	"s16":    Keyword,
	"s24":    Keyword,
	"s32":    Keyword,
	"s48":    Keyword,
	"s64":    Keyword,
	"s96":    Keyword,
	"s128":   Keyword,
	"char":   Keyword,
	"char16": Keyword,
	"bool":   Keyword,
	"float":  Keyword,
	"double": Keyword,
	"str":    Keyword,
	"auto":   Keyword,
	"padding": Keyword,

	// Control flow
	"if":       Keyword,
	"else":     Keyword,
	"while":    Keyword,
	"for":      Keyword,
	"match":    Keyword,
	"return":   Keyword,
	"break":    Keyword,
	"continue": Keyword,
	"try":      Keyword,
	"catch":    Keyword,

	// Declarations
	"fn":        Keyword,
	"namespace": Keyword,
	"import":    Keyword,
	"const":     Keyword,
	"in":        Keyword,
	"out":       Keyword,
	"reference": Keyword,

	// Values
	"true":  Keyword,
	"false": Keyword,
	"null":  Keyword,

	// Other
	"parent": Keyword,
	"this":   Keyword,
	"as":     Keyword,
	"is":     Keyword,
	"from":   Keyword,

	// Endianness
	"little_endian": Keyword,
	"big_endian":    Keyword,
	"signed":        Keyword,
	"unsigned":      Keyword,

	// Type operators
	"sizeof":     Keyword,
	"addressof":  Keyword,
	"typenameof": Keyword,

	// Godecompose extensions
	"instr":    Keyword,
	"gen":      Keyword,
	"bind":     Keyword,
	"pattern":  Keyword,
	"arch":     Keyword,
	"platform": Keyword,
}

// LookupKeyword checks if an identifier string is a keyword.
// Returns (Type, true) if it is a keyword, (Illegal, false) otherwise.
func LookupKeyword(ident string) (Type, bool) {
	if t, ok := keywordToToken[ident]; ok {
		return t, true
	}
	return Illegal, false
}

// IsKeyword reports whether ident is a keyword in the pattern language.
func IsKeyword(ident string) bool {
	_, ok := keywordToToken[ident]
	return ok
}

// Token represents a single lexical token with its position in the source.
type Token struct {
	Type     Type
	Literal  string
	Line     int
	Column   int
	File     string
}

func (t Token) String() string {
	if t.Type == Identifier || t.Type == Keyword || t.Type == Directive {
		return t.Literal
	}
	if t.Type.IsLiteral() {
		return t.Literal
	}
	if s, ok := typeNames[t.Type]; ok {
		return s
	}
	return t.Literal
}

// Position returns the source location as a "file:line:col" string.
func (t Token) Position() string {
	if t.File == "" {
		return ""
	}
	return t.File
}
