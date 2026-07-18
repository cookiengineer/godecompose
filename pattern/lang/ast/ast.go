// Package ast defines the Abstract Syntax Tree node types for the
// ImHex-compatible Pattern Language with godecompose extensions.
package ast

import "github.com/cookiengineer/godecompose/pattern/lang/token"

// Node is the base interface for all AST nodes.
type Node interface {
	Token() token.Token
	String() string
}

// Position returns the token position from a node's associated token.
func position(n Node) token.Token {
	if n == nil {
		return token.Token{}
	}
	return n.Token()
}

// --- Literals ---

type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (n *IntegerLiteral) Token() token.Token { return n.Tok }
func (n *IntegerLiteral) String() string      { return n.Tok.Literal }

type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (n *FloatLiteral) Token() token.Token { return n.Tok }
func (n *FloatLiteral) String() string      { return n.Tok.Literal }

type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (n *StringLiteral) Token() token.Token { return n.Tok }
func (n *StringLiteral) String() string      { return n.Tok.Literal }

type CharLiteral struct {
	Tok   token.Token
	Value rune
}

func (n *CharLiteral) Token() token.Token { return n.Tok }
func (n *CharLiteral) String() string      { return n.Tok.Literal }

type BoolLiteral struct {
	Tok   token.Token
	Value bool
}

func (n *BoolLiteral) Token() token.Token { return n.Tok }
func (n *BoolLiteral) String() string      { return n.Tok.Literal }

type NullLiteral struct {
	Tok token.Token
}

func (n *NullLiteral) Token() token.Token { return n.Tok }
func (n *NullLiteral) String() string      { return "null" }

// --- Identifiers and Names ---

type Identifier struct {
	Tok   token.Token
	Name  string
}

func (n *Identifier) Token() token.Token { return n.Tok }
func (n *Identifier) String() string      { return n.Name }

// --- Types ---

type BuiltinType struct {
	Tok  token.Token
	Name string // e.g., "u8", "u32", "float", "str"
}

func (n *BuiltinType) Token() token.Token { return n.Tok }
func (n *BuiltinType) String() string      { return n.Name }

type CustomType struct {
	Tok  token.Token
	Name string
	Args []Node // template arguments
}

func (n *CustomType) Token() token.Token { return n.Tok }
func (n *CustomType) String() string      { return n.Name }

type PointerType struct {
	Tok  token.Token
	Base Node // the type being pointed to
}

func (n *PointerType) Token() token.Token { return n.Tok }
func (n *PointerType) String() string      { return "*" }

// --- Expressions ---

type Expression interface {
	Node
	expressionNode()
}

// BinaryExpression represents binary operations: +, -, *, /, &&, ||, ==, etc.
type BinaryExpression struct {
	Tok      token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (n *BinaryExpression) Token() token.Token { return n.Tok }
func (n *BinaryExpression) String() string      { return n.Operator }
func (n *BinaryExpression) expressionNode()     {}

// UnaryExpression represents unary operations: -, !, ~
type UnaryExpression struct {
	Tok      token.Token
	Operator string
	Right    Expression
}

func (n *UnaryExpression) Token() token.Token { return n.Tok }
func (n *UnaryExpression) String() string      { return n.Operator }
func (n *UnaryExpression) expressionNode()     {}

// TernaryExpression represents condition ? trueExpr : falseExpr
type TernaryExpression struct {
	Tok       token.Token
	Condition Expression
	TrueExpr  Expression
	FalseExpr Expression
}

func (n *TernaryExpression) Token() token.Token { return n.Tok }
func (n *TernaryExpression) String() string      { return "?:" }
func (n *TernaryExpression) expressionNode()     {}

// CastExpression represents type(expr) or expr as type
type CastExpression struct {
	Tok  token.Token
	Type Node
	Expr Expression
}

func (n *CastExpression) Token() token.Token { return n.Tok }
func (n *CastExpression) String() string      { return "cast" }
func (n *CastExpression) expressionNode()     {}

// CallExpression represents a function call: name(args)
type CallExpression struct {
	Tok      token.Token
	Function Expression
	Args     []Expression
}

func (n *CallExpression) Token() token.Token { return n.Tok }
func (n *CallExpression) String() string      { return "call" }
func (n *CallExpression) expressionNode()     {}

// IndexExpression represents array/member access: expr[index] or expr.name
type IndexExpression struct {
	Tok   token.Token
	Left  Expression
	Index Expression
}

func (n *IndexExpression) Token() token.Token { return n.Tok }
func (n *IndexExpression) String() string      { return "[]" }
func (n *IndexExpression) expressionNode()     {}

// ScopeExpression represents namespace::member access
type ScopeExpression struct {
	Tok    token.Token
	Left   Expression
	Member string
}

func (n *ScopeExpression) Token() token.Token { return n.Tok }
func (n *ScopeExpression) String() string      { return "::" }
func (n *ScopeExpression) expressionNode()     {}

// TypeOperator represents sizeof/addressof/typenameof(expr)
type TypeOperator struct {
	Tok      token.Token
	Operator string
	Expr     Expression
}

func (n *TypeOperator) Token() token.Token { return n.Tok }
func (n *TypeOperator) String() string      { return n.Operator }
func (n *TypeOperator) expressionNode()     {}

// MatchExpression represents match(value) { cases }
type MatchExpression struct {
	Tok   token.Token
	Value Expression
	Cases []MatchCase
}

func (n *MatchExpression) Token() token.Token { return n.Tok }
func (n *MatchExpression) String() string      { return "match" }
func (n *MatchExpression) expressionNode()     {}

type MatchCase struct {
	Patterns []MatchPattern
	Body     []Node
}

type MatchPattern struct {
	Tok   token.Token
	Start Expression // single value, or lower bound for range
	End   Expression // nil for non-range patterns
}

// Make literals implement Expression interface
func (n *IntegerLiteral) expressionNode() {}
func (n *FloatLiteral) expressionNode()   {}
func (n *StringLiteral) expressionNode()  {}
func (n *CharLiteral) expressionNode()    {}
func (n *BoolLiteral) expressionNode()    {}
func (n *NullLiteral) expressionNode()    {}
func (n *Identifier) expressionNode()     {}

// --- Statements ---

type Statement interface {
	Node
	statementNode()
}

// CompoundStatement represents { ... } block
type CompoundStatement struct {
	Tok      token.Token
	Statements []Node
}

func (n *CompoundStatement) Token() token.Token { return n.Tok }
func (n *CompoundStatement) String() string      { return "{...}" }
func (n *CompoundStatement) statementNode()      {}

// ConditionalStatement represents if/else
type ConditionalStatement struct {
	Tok       token.Token
	Condition Expression
	Body      Node
	ElseBody  Node // nil if no else
}

func (n *ConditionalStatement) Token() token.Token { return n.Tok }
func (n *ConditionalStatement) String() string      { return "if" }
func (n *ConditionalStatement) statementNode()      {}

// WhileStatement represents while loop
type WhileStatement struct {
	Tok       token.Token
	Condition Expression
	Body      Node
}

func (n *WhileStatement) Token() token.Token { return n.Tok }
func (n *WhileStatement) String() string      { return "while" }
func (n *WhileStatement) statementNode()      {}

// ForStatement represents for(init; cond; post) loop
type ForStatement struct {
	Tok       token.Token
	Init      Node
	Condition Expression
	Post      Node
	Body      Node
}

func (n *ForStatement) Token() token.Token { return n.Tok }
func (n *ForStatement) String() string      { return "for" }
func (n *ForStatement) statementNode()      {}

// ReturnStatement represents return [value]
type ReturnStatement struct {
	Tok   token.Token
	Value Expression
}

func (n *ReturnStatement) Token() token.Token { return n.Tok }
func (n *ReturnStatement) String() string      { return "return" }
func (n *ReturnStatement) statementNode()      {}

// BreakStatement represents break
type BreakStatement struct {
	Tok token.Token
}

func (n *BreakStatement) Token() token.Token { return n.Tok }
func (n *BreakStatement) String() string      { return "break" }
func (n *BreakStatement) statementNode()      {}

// ContinueStatement represents continue
type ContinueStatement struct {
	Tok token.Token
}

func (n *ContinueStatement) Token() token.Token { return n.Tok }
func (n *ContinueStatement) String() string      { return "continue" }
func (n *ContinueStatement) statementNode()      {}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Tok  token.Token
	Expr Expression
}

func (n *ExpressionStatement) Token() token.Token { return n.Tok }
func (n *ExpressionStatement) String() string      { return "expr" }
func (n *ExpressionStatement) statementNode()      {}

// TryCatchStatement represents try/catch
type TryCatchStatement struct {
	Tok     token.Token
	TryBody Node
	CatchID string
	CatchBody Node
}

func (n *TryCatchStatement) Token() token.Token { return n.Tok }
func (n *TryCatchStatement) String() string      { return "try" }
func (n *TryCatchStatement) statementNode()      {}

// --- Declarations ---

// VariableDeclaration represents a variable declaration: type name [@ offset] [= value]
type VariableDeclaration struct {
	Tok      token.Token
	Type     Node
	Name     string
	Offset   Expression // nil if no @ placement
	Value    Expression // nil if no initializer
	IsConst  bool
	IsIn     bool
	IsOut    bool
	IsRef    bool
}

func (n *VariableDeclaration) Token() token.Token { return n.Tok }
func (n *VariableDeclaration) String() string      { return n.Name }
func (n *VariableDeclaration) statementNode()      {}

// ArrayVariableDeclaration represents a fixed-size array: type name[size] [@ offset]
type ArrayVariableDeclaration struct {
	Tok      token.Token
	Type     Node
	Name     string
	Size     Expression
	Offset   Expression
}

func (n *ArrayVariableDeclaration) Token() token.Token { return n.Tok }
func (n *ArrayVariableDeclaration) String() string      { return n.Name }
func (n *ArrayVariableDeclaration) statementNode()      {}

// PointerVariableDeclaration represents type *name @ offset
type PointerVariableDeclaration struct {
	Tok    token.Token
	Type   Node
	Name   string
	Offset Expression
}

func (n *PointerVariableDeclaration) Token() token.Token { return n.Tok }
func (n *PointerVariableDeclaration) String() string      { return n.Name }
func (n *PointerVariableDeclaration) statementNode()      {}

// Assignment represents lhs = rhs
type Assignment struct {
	Tok      token.Token
	Left     Expression
	Operator string // =, +=, -=, etc.
	Right    Expression
}

func (n *Assignment) Token() token.Token { return n.Tok }
func (n *Assignment) String() string      { return n.Operator }
func (n *Assignment) statementNode()      {}

// --- Type Definitions ---

// StructDefinition represents struct Name [: Base] { members };
type StructDefinition struct {
	Tok      token.Token
	Name     string
	Base     string // empty if no inheritance
	Members  []Node
}

func (n *StructDefinition) Token() token.Token { return n.Tok }
func (n *StructDefinition) String() string      { return n.Name }

// UnionDefinition represents union Name { members };
type UnionDefinition struct {
	Tok     token.Token
	Name    string
	Members []Node
}

func (n *UnionDefinition) Token() token.Token { return n.Tok }
func (n *UnionDefinition) String() string      { return n.Name }

// EnumDefinition represents enum Name : type { values };
type EnumDefinition struct {
	Tok    token.Token
	Name   string
	Base   string // underlying type, e.g., "u16"
	Values []EnumValue
}

func (n *EnumDefinition) Token() token.Token { return n.Tok }
func (n *EnumDefinition) String() string      { return n.Name }

type EnumValue struct {
	Name  string
	Value Expression // nil if auto-increment
}

// BitfieldDefinition represents bitfield Name { fields };
type BitfieldDefinition struct {
	Tok    token.Token
	Name   string
	Fields []BitfieldField
}

func (n *BitfieldDefinition) Token() token.Token { return n.Tok }
func (n *BitfieldDefinition) String() string      { return n.Name }

type BitfieldField struct {
	Name string
	Size Expression
}

// UsingDeclaration represents using Name = Type;
type UsingDeclaration struct {
	Tok  token.Token
	Name string
	Type Node
}

func (n *UsingDeclaration) Token() token.Token { return n.Tok }
func (n *UsingDeclaration) String() string      { return n.Name }

// --- Namespace ---

// NamespaceDeclaration represents namespace Name { ... }
type NamespaceDeclaration struct {
	Tok  token.Token
	Name string
	Body []Node
}

func (n *NamespaceDeclaration) Token() token.Token { return n.Tok }
func (n *NamespaceDeclaration) String() string      { return n.Name }

// --- Import ---

// ImportDeclaration represents import std::name; or import "file";
type ImportDeclaration struct {
	Tok  token.Token
	Path string
}

func (n *ImportDeclaration) Token() token.Token { return n.Tok }
func (n *ImportDeclaration) String() string      { return n.Path }

// --- Function ---

// FunctionDefinition represents fn name(params) { body }
type FunctionDefinition struct {
	Tok      token.Token
	Name     string
	Params   []FunctionParam
	Body     *CompoundStatement
}

func (n *FunctionDefinition) Token() token.Token { return n.Tok }
func (n *FunctionDefinition) String() string      { return n.Name }

type FunctionParam struct {
	Name         string
	Type         Node
	DefaultValue Expression // nil if no default
}

// --- Attributes ---

// Attribute represents [[name(args)]]
type Attribute struct {
	Tok  token.Token
	Name string
	Args []Expression
}

func (n *Attribute) Token() token.Token { return n.Tok }
func (n *Attribute) String() string      { return n.Name }

// --- Endianness Directive ---

// EndianDirective represents little_endian; or big_endian;
type EndianDirective struct {
	Tok      token.Token
	IsBigEndian bool
}

func (n *EndianDirective) Token() token.Token { return n.Tok }
func (n *EndianDirective) String() string {
	if n.IsBigEndian {
		return "big_endian"
	}
	return "little_endian"
}

// --- Godecompose Extensions ---

// ArchDirective represents arch x86_64; or arch arm64;
type ArchDirective struct {
	Tok  token.Token
	Arch string
}

func (n *ArchDirective) Token() token.Token { return n.Tok }
func (n *ArchDirective) String() string      { return n.Arch }

// PlatformDirective represents platform linux, darwin;
type PlatformDirective struct {
	Tok      token.Token
	Platforms []string
}

func (n *PlatformDirective) Token() token.Token { return n.Tok }
func (n *PlatformDirective) String() string {
	if len(n.Platforms) > 0 {
		return n.Platforms[0]
	}
	return "platform"
}

// PatternDefinition represents a named pattern with instr/gen/bind blocks.
type PatternDefinition struct {
	Tok         token.Token
	Name        string
	Library     string
	Version     string
	Description string
	Arch        string
	Platforms   []string
	InstrBlocks []*InstrBlock
	GenBlock    *GenBlock
	BindBlock   *BindBlock
}

func (n *PatternDefinition) Token() token.Token { return n.Tok }
func (n *PatternDefinition) String() string      { return n.Name }

// InstrBlock represents an instr block describing an assembly pattern.
type InstrBlock struct {
	Tok        token.Token
	Name       string
	Labels     []InstrLabel
	Alternatives [][]InstructionPattern // each alt is a sequence of instructions
}

func (n *InstrBlock) Token() token.Token { return n.Tok }
func (n *InstrBlock) String() string      { return n.Name }

type InstrLabel struct {
	Name    string
	Address uint64 // 0 if not known
}

// InstructionPattern represents a single instruction line in an instr block.
type InstructionPattern struct {
	Tok      token.Token
	Opcode   string            // e.g., "MOVQ", "CALL", "SYSCALL"
	Operands []OperandPattern
	Label    string            // @labelname if this instruction is a label target
}

type OperandPattern struct {
	Tok token.Token
	// Exactly one of the following is non-nil/true:
	IsWildcard     bool
	IsImmediate    bool   // $ prefix
	LiteralValue   string // literal value (number, string)
	RegisterName   string // e.g., "RAX", "X0"
	CaptureVar     string // identifier to bind matched value to
	MemoryRef      *MemoryRefPattern
}

// MemoryRefPattern represents an addressing mode like offset(base)(index*scale).
type MemoryRefPattern struct {
	Offset    string // capture variable or literal
	Base      string // register or capture var
	Index     string // register or capture var
	Scale     string // 1, 2, 4, 8 or capture var
	Segment   string // segment override: FS, GS, etc.
}

// GenBlock represents a source code generation template block.
type GenBlock struct {
	Tok       token.Token
	Statements []GenStatement
}

func (n *GenBlock) Token() token.Token { return n.Tok }
func (n *GenBlock) String() string      { return "gen" }

type GenStatement interface {
	Node
	genStatementNode()
}

// GenText represents raw template text with $variable substitution.
type GenText struct {
	Tok  token.Token
	Text string
}

func (n *GenText) Token() token.Token { return n.Tok }
func (n *GenText) String() string      { return n.Text }
func (n *GenText) genStatementNode()   {}

// GenExpr represents ${expression} interpolation.
type GenExpr struct {
	Tok  token.Token
	Expr Expression
}

func (n *GenExpr) Token() token.Token { return n.Tok }
func (n *GenExpr) String() string      { return "${...}" }
func (n *GenExpr) genStatementNode()   {}

// GenConditional represents an if/else in a gen block.
type GenConditional struct {
	Tok       token.Token
	Condition Expression
	Body      []GenStatement
	ElseBody  []GenStatement
}

func (n *GenConditional) Token() token.Token { return n.Tok }
func (n *GenConditional) String() string      { return "gen:if" }
func (n *GenConditional) genStatementNode()   {}

// GenLoop represents a for loop in a gen block.
type GenLoop struct {
	Tok  token.Token
	Init string
	Cond Expression
	Post string
	Body []GenStatement
}

func (n *GenLoop) Token() token.Token { return n.Tok }
func (n *GenLoop) String() string      { return "gen:for" }
func (n *GenLoop) genStatementNode()   {}

// BindBlock represents variable renaming bindings.
type BindBlock struct {
	Tok      token.Token
	Bindings []BindEntry
}

func (n *BindBlock) Token() token.Token { return n.Tok }
func (n *BindBlock) String() string      { return "bind" }

type BindEntry struct {
	CaptureVar string
	Alias      string
}

// --- Helper types ---

// Program represents the top-level AST root containing all nodes.
type Program struct {
	Nodes      []Node
	Imports    []*ImportDeclaration
	Structs    []*StructDefinition
	Unions     []*UnionDefinition
	Enums      []*EnumDefinition
	Bitfields  []*BitfieldDefinition
	Functions  []*FunctionDefinition
	Variables  []*VariableDeclaration
	Patterns   []*PatternDefinition
	Namespaces []*NamespaceDeclaration
}
