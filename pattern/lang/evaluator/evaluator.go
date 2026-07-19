// Package evaluator implements a tree-walking interpreter for the pattern
// language AST. It supports binary data mode (creating patterns at offsets),
// instruction matching mode (compiling assembly patterns), and source
// generation mode (template expansion).
package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cookiengineer/godecompose/pattern/lang/ast"
)

// Value represents a runtime value in the pattern language.
type Value struct {
	Type  ValueType
	Int   int64
	Float float64
	Str   string
	Bool  bool
	Null  bool
}

type ValueType int

const (
	ValNull ValueType = iota
	ValInt
	ValFloat
	ValString
	ValBool
)

func intValue(v int64) Value     { return Value{Type: ValInt, Int: v} }
func floatValue(v float64) Value { return Value{Type: ValFloat, Float: v} }
func stringValue(v string) Value { return Value{Type: ValString, Str: v} }
func boolValue(v bool) Value     { return Value{Type: ValBool, Bool: v} }
func nullValue() Value           { return Value{Type: ValNull, Null: true} }

func (v Value) String() string {
	switch v.Type {
	case ValInt:
		return fmt.Sprintf("%d", v.Int)
	case ValFloat:
		return fmt.Sprintf("%g", v.Float)
	case ValString:
		return v.Str
	case ValBool:
		if v.Bool {
			return "true"
		}
		return "false"
	default:
		return "null"
	}
}

// CompiledPattern is the output of evaluating an instr block: a sequence of
// instruction matchers with capture variable assignments.
type CompiledPattern struct {
	Name        string
	Arch        string
	Platforms   []string
	Library     string
	Version     string
	Alternatives [][]CompiledInstruction
	Bindings    []CompiledBinding
	GenTemplate string
}

type CompiledInstruction struct {
	Opcode     string
	IsLabel    bool
	LabelName  string
	Operands   []CompiledOperand
}

type CompiledOperand struct {
	IsWildcard   bool
	IsImmediate  bool
	Literal      string
	Register     string
	CaptureVar   string
	BaseReg      string
	IndexReg     string
	Scale        string
	Offset       string
}

type CompiledBinding struct {
	CaptureVar string
	Alias      string
}

// Evaluator walks the AST and produces runtime results.
type Evaluator struct {
	scopes    []map[string]Value
	builtins  map[string]func([]Value) (Value, error)
	patterns  []*CompiledPattern
}

// New creates a new evaluator with registered built-in functions.
func New() *Evaluator {
	e := &Evaluator{
		builtins: make(map[string]func([]Value) (Value, error)),
	}
	e.pushScope()
	e.registerBuiltins()
	return e
}

func (e *Evaluator) pushScope() {
	e.scopes = append(e.scopes, make(map[string]Value))
}

func (e *Evaluator) popScope() {
	if len(e.scopes) > 1 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *Evaluator) setVar(name string, val Value) {
	e.scopes[len(e.scopes)-1][name] = val
}

func (e *Evaluator) getVar(name string) (Value, bool) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if v, ok := e.scopes[i][name]; ok {
			return v, true
		}
	}
	return Value{}, false
}

// Evaluate walks a program AST and produces compiled patterns.
func (e *Evaluator) Evaluate(prog *ast.Program) ([]*CompiledPattern, error) {
	e.scopes = []map[string]Value{make(map[string]Value)}
	e.patterns = nil

	for _, node := range prog.Nodes {
		if err := e.evalNode(node); err != nil {
			return nil, err
		}
	}

	return e.patterns, nil
}

func (e *Evaluator) evalNode(node ast.Node) error {
	switch n := node.(type) {
	case *ast.VariableDeclaration:
		val := nullValue()
		if n.Value != nil {
			var err error
			val, err = e.evalExpr(n.Value)
			if err != nil {
				return err
			}
		}
		e.setVar(n.Name, val)

	case *ast.FunctionDefinition:
		e.pushScope()
		for _, p := range n.Params {
			e.setVar(p.Name, nullValue())
		}
		if n.Body != nil {
			for _, s := range n.Body.Statements {
				if err := e.evalNode(s); err != nil {
					e.popScope()
					return err
				}
			}
		}
		e.popScope()

	case *ast.ExpressionStatement:
		_, err := e.evalExpr(n.Expr)
		return err

	case *ast.ReturnStatement:
		return nil

	case *ast.CompoundStatement:
		e.pushScope()
		for _, s := range n.Statements {
			if err := e.evalNode(s); err != nil {
				e.popScope()
				return err
			}
		}
		e.popScope()

	case *ast.ConditionalStatement:
		cond, err := e.evalExpr(n.Condition)
		if err != nil {
			return err
		}
		if cond.Bool {
			return e.evalNode(n.Body)
		} else if n.ElseBody != nil {
			return e.evalNode(n.ElseBody)
		}

	case *ast.WhileStatement:
		for {
			cond, err := e.evalExpr(n.Condition)
			if err != nil {
				return err
			}
			if !cond.Bool {
				break
			}
			if err := e.evalNode(n.Body); err != nil {
				return err
			}
		}

	case *ast.ArrayVariableDeclaration:
		e.setVar(n.Name, nullValue())

	case *ast.PatternDefinition:
		return e.evalPattern(n)

	case *ast.ArchDirective, *ast.PlatformDirective, *ast.EndianDirective,
		*ast.ImportDeclaration, *ast.UsingDeclaration, *ast.NamespaceDeclaration:
		return nil

	case *ast.StructDefinition, *ast.UnionDefinition, *ast.EnumDefinition, *ast.BitfieldDefinition:
		return nil
	}
	return nil
}

func (e *Evaluator) evalPattern(pat *ast.PatternDefinition) error {
	cp := &CompiledPattern{
		Name:      pat.Name,
		Library:   pat.Library,
		Version:   pat.Version,
		Arch:      pat.Arch,
		Platforms: pat.Platforms,
	}

	for _, block := range pat.InstrBlocks {
		var altInstructions [][]CompiledInstruction
		for _, alt := range block.Alternatives {
			var instructions []CompiledInstruction
			for _, ip := range alt {
				ci := CompiledInstruction{
					Opcode:    ip.Opcode,
					IsLabel:   ip.Label != "",
					LabelName: ip.Label,
				}
				for _, op := range ip.Operands {
					ci.Operands = append(ci.Operands, CompiledOperand{
						IsWildcard:  op.IsWildcard,
						IsImmediate: op.IsImmediate,
						Literal:     op.LiteralValue,
						Register:    op.RegisterName,
						CaptureVar:  op.CaptureVar,
					})
					if op.MemoryRef != nil {
						ci.Operands[len(ci.Operands)-1].BaseReg = op.MemoryRef.Base
						ci.Operands[len(ci.Operands)-1].IndexReg = op.MemoryRef.Index
						ci.Operands[len(ci.Operands)-1].Scale = op.MemoryRef.Scale
						ci.Operands[len(ci.Operands)-1].Offset = op.MemoryRef.Offset
					}
				}
				instructions = append(instructions, ci)
			}
			altInstructions = append(altInstructions, instructions)
		}
		cp.Alternatives = altInstructions
	}

	if pat.BindBlock != nil {
		for _, b := range pat.BindBlock.Bindings {
			cp.Bindings = append(cp.Bindings, CompiledBinding{
				CaptureVar: b.CaptureVar,
				Alias:      b.Alias,
			})
		}
	}

	if pat.GenBlock != nil {
		cp.GenTemplate = e.evalGenBlock(pat.GenBlock, cp.Bindings)
	}

	e.patterns = append(e.patterns, cp)
	return nil
}

func (e *Evaluator) evalGenBlock(block *ast.GenBlock, bindings []CompiledBinding) string {
	var buf strings.Builder
	var prevWasExpr bool
	for _, stmt := range block.Statements {
		isExpr := false
		switch stmt.(type) {
		case *ast.GenExpr:
			isExpr = true
		}
		s := e.evalGenStmt(stmt, bindings)
		if len(s) == 0 {
			continue
		}
		if prevWasExpr && !isExpr && !isPunctOrSpace(s[0]) {
			buf.WriteByte(' ')
		}
		if !prevWasExpr && isExpr && buf.Len() > 0 {
			last := buf.String()[buf.Len()-1]
			if !isPunctOrSpace(last) {
				buf.WriteByte(' ')
			}
		}
		buf.WriteString(s)
		prevWasExpr = isExpr
	}
	return buf.String()
}

func isPunctOrSpace(b byte) bool {
	return b == '(' || b == ')' || b == ',' || b == ';' || b == '.' || b == '{' || b == '}' || b == ' ' || b == '\n' || b == '\t'
}

func (e *Evaluator) evalGenStmt(stmt ast.GenStatement, bindings []CompiledBinding) string {
	switch s := stmt.(type) {
	case *ast.GenText:
		return e.expandVariables(s.Text, bindings)
	case *ast.GenExpr:
		val, err := e.evalGenExpr(s.Expr, bindings)
		if err != nil {
			return "<error>"
		}
		return val.String()
	case *ast.GenConditional:
		cond, err := e.evalExpr(s.Condition)
		body := s.Body
		if err != nil || !cond.Bool {
			body = s.ElseBody
		}
		var buf strings.Builder
		for _, st := range body {
			buf.WriteString(e.evalGenStmt(st, bindings))
		}
		return buf.String()
	case *ast.GenLoop:
		count, err := e.evalExpr(s.Cond)
		if err != nil {
			return ""
		}
		n := int(count.Int)
		if n < 0 {
			n = 0
		}
		if n > 1000 {
			n = 1000
		}
		var buf strings.Builder
		for i := 0; i < n; i++ {
			for _, st := range s.Body {
				buf.WriteString(e.evalGenStmt(st, bindings))
			}
		}
		return buf.String()
	}
	return ""
}

func (e *Evaluator) evalGenExpr(expr ast.Expression, bindings []CompiledBinding) (Value, error) {
	if ident, ok := expr.(*ast.Identifier); ok {
		for _, b := range bindings {
			if b.CaptureVar == ident.Name {
				return stringValue(b.Alias), nil
			}
		}
		return stringValue(ident.Name), nil
	}
	if unary, ok := expr.(*ast.UnaryExpression); ok && unary.Operator == "$" {
		if ident, ok := unary.Right.(*ast.Identifier); ok {
			for _, b := range bindings {
				if b.CaptureVar == ident.Name {
					return stringValue(b.Alias), nil
				}
			}
			return stringValue(ident.Name), nil
		}
		return e.evalExpr(unary.Right)
	}
	return e.evalExpr(expr)
}

func (e *Evaluator) expandVariables(text string, bindings []CompiledBinding) string {
	result := text
	for _, b := range bindings {
		placeholder := "$" + b.CaptureVar
		result = strings.ReplaceAll(result, placeholder, b.Alias)
	}
	return result
}

func (e *Evaluator) evalExpr(expr ast.Expression) (Value, error) {
	if expr == nil {
		return nullValue(), nil
	}

	switch n := expr.(type) {
	case *ast.IntegerLiteral:
		i, err := strconv.ParseInt(n.Tok.Literal, 0, 64)
		if err != nil {
			return intValue(0), nil
		}
		return intValue(i), nil

	case *ast.FloatLiteral:
		f, err := strconv.ParseFloat(n.Tok.Literal, 64)
		if err != nil {
			return floatValue(0), nil
		}
		return floatValue(f), nil

	case *ast.StringLiteral:
		return stringValue(n.Value), nil

	case *ast.CharLiteral:
		return stringValue(n.Tok.Literal), nil

	case *ast.BoolLiteral:
		return boolValue(n.Value), nil

	case *ast.NullLiteral:
		return nullValue(), nil

	case *ast.Identifier:
		if val, ok := e.getVar(n.Name); ok {
			return val, nil
		}
		return stringValue(n.Name), nil

	case *ast.BinaryExpression:
		return e.evalBinary(n)

	case *ast.UnaryExpression:
		right, err := e.evalExpr(n.Right)
		if err != nil {
			return nullValue(), err
		}
		switch n.Operator {
		case "-":
			right.Int = -right.Int
			right.Float = -right.Float
			return right, nil
		case "!":
			return boolValue(!right.Bool), nil
		default:
			return right, nil
		}

	case *ast.CallExpression:
		return e.evalCall(n)

	case *ast.ScopeExpression:
		return stringValue(n.Member), nil

	case *ast.CastExpression:
		return e.evalExpr(n.Expr)

	default:
		return nullValue(), nil
	}
}

func (e *Evaluator) evalBinary(expr *ast.BinaryExpression) (Value, error) {
	left, err := e.evalExpr(expr.Left)
	if err != nil {
		return nullValue(), err
	}
	right, err := e.evalExpr(expr.Right)
	if err != nil {
		return nullValue(), err
	}

	switch expr.Operator {
	case "+":
		if left.Type == ValString || right.Type == ValString {
			return stringValue(left.String() + right.String()), nil
		}
		return intValue(left.Int + right.Int), nil
	case "-":
		return intValue(left.Int - right.Int), nil
	case "*":
		return intValue(left.Int * right.Int), nil
	case "/":
		if right.Int == 0 {
			return intValue(0), nil
		}
		return intValue(left.Int / right.Int), nil
	case "%":
		if right.Int == 0 {
			return intValue(0), nil
		}
		return intValue(left.Int % right.Int), nil
	case "==":
		return boolValue(left.String() == right.String()), nil
	case "!=":
		return boolValue(left.String() != right.String()), nil
	case "<":
		return boolValue(left.Int < right.Int), nil
	case ">":
		return boolValue(left.Int > right.Int), nil
	case "<=":
		return boolValue(left.Int <= right.Int), nil
	case ">=":
		return boolValue(left.Int >= right.Int), nil
	case "&&":
		return boolValue(left.Bool && right.Bool), nil
	case "||":
		return boolValue(left.Bool || right.Bool), nil
	default:
		return left, nil
	}
}

func (e *Evaluator) evalCall(expr *ast.CallExpression) (Value, error) {
	funcName := ""
	if ident, ok := expr.Function.(*ast.Identifier); ok {
		funcName = ident.Name
	}

	var args []Value
	for _, arg := range expr.Args {
		val, err := e.evalExpr(arg)
		if err != nil {
			return nullValue(), err
		}
		args = append(args, val)
	}

	if fn, ok := e.builtins[funcName]; ok {
		return fn(args)
	}

	return nullValue(), nil
}

func (e *Evaluator) registerBuiltins() {
	e.builtins["print"] = func(args []Value) (Value, error) {
		for _, a := range args {
			fmt.Print(a.String())
		}
		return nullValue(), nil
	}
}
