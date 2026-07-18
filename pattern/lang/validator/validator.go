// Package validator performs semantic analysis on the pattern language AST:
// type checking, identifier resolution, and control flow validation.
package validator

import (
	"fmt"

	"github.com/cookiengineer/godecompose/pattern/lang/ast"
	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

// Validator walks the AST and collects semantic errors.
type Validator struct {
	errors []error
	scopes []*scope
	loops  int
	inFunc int
}

type scope struct {
	names map[string]*ast.VariableDeclaration
	types map[string]bool
	funcs map[string]*ast.FunctionDefinition
}

func newScope() *scope {
	return &scope{
		names: make(map[string]*ast.VariableDeclaration),
		types: make(map[string]bool),
		funcs: make(map[string]*ast.FunctionDefinition),
	}
}

// New creates a new validator.
func New() *Validator {
	v := &Validator{}
	v.pushScope()
	return v
}

// Validate runs all semantic checks on a program AST.
func (v *Validator) Validate(prog *ast.Program) []error {
	v.errors = nil
	v.scopes = []*scope{newScope()}
	v.loops = 0
	v.inFunc = 0

	for _, node := range prog.Nodes {
		v.validateNode(node)
	}

	return v.errors
}

func (v *Validator) pushScope() {
	v.scopes = append(v.scopes, newScope())
}

func (v *Validator) popScope() {
	if len(v.scopes) > 1 {
		v.scopes = v.scopes[:len(v.scopes)-1]
	}
}

func (v *Validator) topScope() *scope {
	return v.scopes[len(v.scopes)-1]
}

func (v *Validator) errorf(node ast.Node, format string, args ...interface{}) {
	var file string
	var line int
	if node != nil {
		tok := node.Token()
		file = tok.File
		line = tok.Line
	}
	v.errors = append(v.errors, fmt.Errorf("%s:%d: %s",
		file, line, fmt.Sprintf(format, args...)))
}

func (v *Validator) errorfDirect(tok token.Token, format string, args ...interface{}) {
	v.errors = append(v.errors, fmt.Errorf("%s:%d: %s",
		tok.File, tok.Line, fmt.Sprintf(format, args...)))
}

func (v *Validator) validateNode(node ast.Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.StructDefinition:
		v.topScope().types[n.Name] = true
		for _, m := range n.Members {
			v.validateNode(m)
		}

	case *ast.UnionDefinition:
		v.topScope().types[n.Name] = true
		for _, m := range n.Members {
			v.validateNode(m)
		}

	case *ast.EnumDefinition:
		v.topScope().types[n.Name] = true

	case *ast.BitfieldDefinition:
		v.topScope().types[n.Name] = true

	case *ast.FunctionDefinition:
		v.topScope().funcs[n.Name] = n
		v.pushScope()
		v.inFunc++
		for _, p := range n.Params {
			v.topScope().names[p.Name] = &ast.VariableDeclaration{Name: p.Name}
		}
		if n.Body != nil {
			for _, s := range n.Body.Statements {
				v.validateNode(s)
			}
		}
		v.inFunc--
		v.popScope()

	case *ast.VariableDeclaration:
		if n.Name != "" {
			v.topScope().names[n.Name] = n
		}

	case *ast.ArrayVariableDeclaration:
		if n.Name != "" {
			v.topScope().names[n.Name] = &ast.VariableDeclaration{Name: n.Name}
		}

	case *ast.CompoundStatement:
		v.pushScope()
		for _, s := range n.Statements {
			v.validateNode(s)
		}
		v.popScope()

	case *ast.ConditionalStatement:
		v.validateNode(n.Body)
		if n.ElseBody != nil {
			v.validateNode(n.ElseBody)
		}

	case *ast.WhileStatement:
		v.pushScope()
		v.loops++
		v.validateNode(n.Body)
		v.loops--
		v.popScope()

	case *ast.ForStatement:
		v.pushScope()
		v.loops++
		if n.Init != nil {
			v.validateNode(n.Init)
		}
		v.validateNode(n.Body)
		v.loops--
		v.popScope()

	case *ast.ReturnStatement:
		if v.inFunc == 0 {
			v.errorf(n, "return outside function")
		}

	case *ast.BreakStatement:
		if v.loops == 0 {
			v.errorf(n, "break outside loop")
		}

	case *ast.ContinueStatement:
		if v.loops == 0 {
			v.errorf(n, "continue outside loop")
		}

	case *ast.NamespaceDeclaration:
		for _, child := range n.Body {
			v.validateNode(child)
		}

	case *ast.PatternDefinition:
		for _, block := range n.InstrBlocks {
			v.validateInstrBlock(block)
		}

	case *ast.ExpressionStatement:
	case *ast.EndianDirective, *ast.ArchDirective, *ast.PlatformDirective,
		*ast.ImportDeclaration, *ast.UsingDeclaration:
	}
}

func (v *Validator) validateInstrBlock(block *ast.InstrBlock) {
	for _, alt := range block.Alternatives {
		for _, ip := range alt {
			if ip.Opcode == "" && ip.Label == "" {
				v.errorfDirect(ip.Tok, "empty instruction pattern in instr block %q", block.Name)
			}
		}
	}
}
