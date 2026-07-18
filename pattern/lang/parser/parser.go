// Package parser implements a recursive descent parser for the ImHex-compatible
// Pattern Language with godecompose extensions. It uses a backtracking mechanism
// via mark/reset to handle ambiguous syntax.
package parser

import (
	"fmt"

	"github.com/cookiengineer/godecompose/pattern/lang/ast"
	"github.com/cookiengineer/godecompose/pattern/lang/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}

	for !p.isEOF() {
		if p.peekIs(token.Directive) {
			p.advance()
			continue
		}

		node, err := p.parseTopLevel()
		if err != nil {
			return program, fmt.Errorf("line %d: %w", p.current().Line, err)
		}
		if node == nil {
			if !p.isEOF() {
				return program, fmt.Errorf("line %d: unexpected token %q", p.current().Line, p.current().Literal)
			}
			break
		}

		switch n := node.(type) {
		case *ast.ImportDeclaration:
			program.Imports = append(program.Imports, n)
		case *ast.StructDefinition:
			program.Structs = append(program.Structs, n)
		case *ast.UnionDefinition:
			program.Unions = append(program.Unions, n)
		case *ast.EnumDefinition:
			program.Enums = append(program.Enums, n)
		case *ast.BitfieldDefinition:
			program.Bitfields = append(program.Bitfields, n)
		case *ast.FunctionDefinition:
			program.Functions = append(program.Functions, n)
		case *ast.VariableDeclaration:
			program.Variables = append(program.Variables, n)
		case *ast.ArrayVariableDeclaration:
			program.Variables = append(program.Variables, &ast.VariableDeclaration{
				Tok:  n.Tok,
				Name: n.Name,
				Type: n.Type,
			})
		case *ast.PointerVariableDeclaration:
			program.Variables = append(program.Variables, &ast.VariableDeclaration{
				Tok:  n.Tok,
				Name: n.Name,
				Type: n.Type,
			})
		case *ast.PatternDefinition:
			program.Patterns = append(program.Patterns, n)
		case *ast.NamespaceDeclaration:
			program.Namespaces = append(program.Namespaces, n)
		}
		program.Nodes = append(program.Nodes, node)
	}

	return program, nil
}

func (p *Parser) parseTopLevel() (ast.Node, error) {
	tok := p.current()

	if tok.Type != token.Keyword {
		return p.parseVariableDeclaration()
	}

	switch tok.Literal {
	case "struct":
		return p.parseStruct()
	case "union":
		return p.parseUnion()
	case "enum":
		return p.parseEnum()
	case "bitfield":
		return p.parseBitfield()
	case "using":
		return p.parseUsing()
	case "fn":
		return p.parseFunction()
	case "namespace":
		return p.parseNamespace()
	case "import":
		return p.parseImport()
	case "little_endian", "big_endian":
		return p.parseEndian()
	case "pattern":
		return p.parsePattern()
	case "arch":
		return p.parseArch()
	case "platform":
		return p.parsePlatform()
	default:
		return p.parseVariableDeclaration()
	}
}

func (p *Parser) parseStruct() (*ast.StructDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	st := &ast.StructDefinition{Tok: tok, Name: name.Literal}

	if p.peekIs(token.Colon) {
		p.advance() // :
		baseTok := p.expect(token.Identifier)
		st.Base = baseTok.Literal
	}

	p.expect(token.LBrace)
	st.Members = p.parseStructMembers()
	p.expect(token.RBrace)
	p.expect(token.Semicolon)

	return st, nil
}

func (p *Parser) parseStructMembers() []ast.Node {
	var members []ast.Node
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		member := p.parseStructMember()
		if member != nil {
			members = append(members, member)
		}
	}
	return members
}

func (p *Parser) parseStructMember() ast.Node {
	if p.peekIs(token.Keyword, "padding") {
		p.advance() // padding
		p.expect(token.LBracket)
		_ = p.parseExpression()
		p.expect(token.RBracket)
		p.expect(token.Semicolon)
		return &ast.VariableDeclaration{
			Tok:  token.Token{},
			Name: "padding",
			Type: &ast.BuiltinType{Name: "padding"},
		}
	}

	node, _ := p.parseVariableDeclaration()
	return node
}

func (p *Parser) parseUnion() (*ast.UnionDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	p.expect(token.LBrace)
	var members []ast.Node
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		member, _ := p.parseVariableDeclaration()
		if member != nil {
			members = append(members, member)
		}
	}
	p.expect(token.RBrace)
	p.expect(token.Semicolon)

	return &ast.UnionDefinition{Tok: tok, Name: name.Literal, Members: members}, nil
}

func (p *Parser) parseEnum() (*ast.EnumDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	p.expect(token.Colon)
	baseType := p.parseType()

	p.expect(token.LBrace)
	var values []ast.EnumValue
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		if p.peekIs(token.Comma) {
			p.advance()
			continue
		}
		valName := p.expect(token.Identifier)
		if valName.Type != token.Identifier {
			p.advance()
			continue
		}
		entry := ast.EnumValue{Name: valName.Literal}
		if p.peekIs(token.Assign) {
			p.advance() // =
			entry.Value = p.parseExpression()
		}
		values = append(values, entry)
		if p.peekIs(token.Comma) {
			p.advance()
		}
		if p.peekIs(token.Semicolon) {
			p.advance()
		}
	}
	p.expect(token.RBrace)
	p.expect(token.Semicolon)

	return &ast.EnumDefinition{Tok: tok, Name: name.Literal, Base: baseName(baseType), Values: values}, nil
}

func (p *Parser) parseBitfield() (*ast.BitfieldDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	p.expect(token.LBrace)
	var fields []ast.BitfieldField
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		fieldName := p.expect(token.Identifier)
		p.expect(token.Colon)
		size := p.parseExpression()
		fields = append(fields, ast.BitfieldField{Name: fieldName.Literal, Size: size})
		p.expect(token.Semicolon)
	}
	p.expect(token.RBrace)
	p.expect(token.Semicolon)

	return &ast.BitfieldDefinition{Tok: tok, Name: name.Literal, Fields: fields}, nil
}

func (p *Parser) parseUsing() (*ast.UsingDeclaration, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)
	p.expect(token.Assign)
	typ := p.parseType()
	p.expect(token.Semicolon)

	return &ast.UsingDeclaration{Tok: tok, Name: name.Literal, Type: typ}, nil
}

func (p *Parser) parseFunction() (*ast.FunctionDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	p.expect(token.LParen)
	var params []ast.FunctionParam
	for !p.peekIs(token.RParen) && !p.isEOF() {
		param := p.parseFunctionParam()
		params = append(params, param)
		if p.peekIs(token.Comma) {
			p.advance()
		}
	}
	p.expect(token.RParen)

	body := p.parseBlock()

	return &ast.FunctionDefinition{Tok: tok, Name: name.Literal, Params: params, Body: body}, nil
}

func (p *Parser) parseFunctionParam() ast.FunctionParam {
	param := ast.FunctionParam{
		Name: "",
		Type: p.parseType(),
	}

	if p.peekIs(token.Identifier) {
		param.Name = p.advance().Literal
	}

	if p.peekIs(token.Assign) {
		p.advance() // =
		param.DefaultValue = p.parseExpression()
	}

	return param
}

func (p *Parser) parseBlock() *ast.CompoundStatement {
	tok := p.expect(token.LBrace)
	var stmts []ast.Node
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	p.expect(token.RBrace)

	return &ast.CompoundStatement{Tok: tok, Statements: stmts}
}

func (p *Parser) parseNamespace() (*ast.NamespaceDeclaration, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	p.expect(token.LBrace)
	var body []ast.Node
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		node, err := p.parseTopLevel()
		if err != nil {
			break
		}
		if node != nil {
			body = append(body, node)
		}
	}
	p.expect(token.RBrace)

	return &ast.NamespaceDeclaration{Tok: tok, Name: name.Literal, Body: body}, nil
}

func (p *Parser) parseImport() (*ast.ImportDeclaration, error) {
	tok := p.advance()
	var path string

	if p.peekIs(token.String) {
		path = p.advance().Literal
	} else {
		for !p.peekIs(token.Semicolon) && !p.isEOF() {
			path += p.advance().Literal
		}
	}
	p.expect(token.Semicolon)

	return &ast.ImportDeclaration{Tok: tok, Path: path}, nil
}

func (p *Parser) parseEndian() (*ast.EndianDirective, error) {
	tok := p.advance()
	p.expect(token.Semicolon)
	return &ast.EndianDirective{Tok: tok, IsBigEndian: tok.Literal == "big_endian"}, nil
}

func (p *Parser) parsePattern() (*ast.PatternDefinition, error) {
	tok := p.advance()
	name := p.expect(token.Identifier)

	pat := &ast.PatternDefinition{Tok: tok, Name: name.Literal}

	p.expect(token.LBrace)
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		switch {
		case p.peekIs(token.Semicolon):
			p.advance()
			continue
		case p.peekIs(token.Keyword, "instr"):
			block, err := p.parseInstrBlock()
			if err != nil {
				break
			}
			pat.InstrBlocks = append(pat.InstrBlocks, block)
		case p.peekIs(token.Keyword, "gen"):
			block, err := p.parseGenBlock()
			if err != nil {
				break
			}
			pat.GenBlock = block
		case p.peekIs(token.Keyword, "bind"):
			block, err := p.parseBindBlock()
			if err != nil {
				break
			}
			pat.BindBlock = block
		case p.peekIs(token.Identifier):
			field := p.advance().Literal
			p.expect(token.Colon)
			switch field {
			case "name":
				pat.Name = p.expect(token.String).Literal
			case "library":
				pat.Library = p.expect(token.String).Literal
			case "version":
				pat.Version = p.expect(token.String).Literal
			case "description":
				pat.Description = p.expect(token.String).Literal
			default:
				p.advance()
			}
			p.expect(token.Semicolon)
		default:
			p.advance()
		}
	}
	p.expect(token.RBrace)

	return pat, nil
}

func (p *Parser) parseArch() (*ast.ArchDirective, error) {
	tok := p.advance()
	arch := p.expect(token.Identifier)
	p.expect(token.Semicolon)
	return &ast.ArchDirective{Tok: tok, Arch: arch.Literal}, nil
}

func (p *Parser) parsePlatform() (*ast.PlatformDirective, error) {
	tok := p.advance()
	var platforms []string
	for !p.peekIs(token.Semicolon) && !p.isEOF() {
		platforms = append(platforms, p.advance().Literal)
		if p.peekIs(token.Comma) {
			p.advance()
		}
	}
	p.expect(token.Semicolon)
	return &ast.PlatformDirective{Tok: tok, Platforms: platforms}, nil
}

func (p *Parser) parseInstrBlock() (*ast.InstrBlock, error) {
	tok := p.advance() // instr
	name := p.advance() // block name (Identifier or Keyword)
	block := &ast.InstrBlock{Tok: tok, Name: name.Literal}

	p.expect(token.LBrace)
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		line, err := p.parseInstructionPattern()
		if err != nil {
			break
		}
		if line != nil {
			if len(block.Alternatives) == 0 {
				block.Alternatives = append(block.Alternatives, nil)
			}
			block.Alternatives[len(block.Alternatives)-1] = append(
				block.Alternatives[len(block.Alternatives)-1], *line)
		}
		if p.peekIs(token.Pipe) {
			p.advance()
			block.Alternatives = append(block.Alternatives, nil)
		}
	}
	p.expect(token.RBrace)

	return block, nil
}

func (p *Parser) parseInstructionPattern() (*ast.InstructionPattern, error) {
	if p.peekIs(token.RBrace) || p.peekIs(token.Pipe) || p.isEOF() {
		return nil, nil
	}

	ip := &ast.InstructionPattern{}

	if p.peekIs(token.At) {
		p.advance() // @
		ip.Label = p.advance().Literal
		if p.peekIs(token.Colon) {
			p.advance() // :
		}
		return ip, nil
	}

	if p.peekIs(token.Identifier) || p.peekIs(token.Keyword) {
		ip.Opcode = p.advance().Literal
	} else {
		return nil, fmt.Errorf("expected opcode")
	}

	for !p.peekIs(token.RBrace) && !p.peekIs(token.Pipe) && !p.isEOF() {
		if p.peekIs(token.Semicolon) {
			p.advance()
			continue
		}
		if p.peekIs(token.Identifier) && isOpcode(p.current().Literal) {
			break
		}
		if p.peekIs(token.At) {
			break
		}
		if p.peekIs(token.Comma) {
			p.advance()
			continue
		}
		op := p.parseOperandPattern()
		if op != nil && (op.IsWildcard || op.IsImmediate || op.RegisterName != "" || op.CaptureVar != "" || op.LiteralValue != "" || op.MemoryRef != nil) {
			ip.Operands = append(ip.Operands, *op)
		}
	}

	return ip, nil
}

func isOpcode(name string) bool {
	if len(name) < 2 {
		return false
	}
	if isRegister(name) {
		return false
	}
	for _, ch := range name {
		if ch >= 'a' && ch <= 'z' {
			return false
		}
	}
	return true
}

func (p *Parser) parseOperandPattern() *ast.OperandPattern {
	op := &ast.OperandPattern{}

	if p.peekIs(token.Asterisk) {
		p.advance()
		op.IsWildcard = true
		return op
	}

	if p.peekIs(token.Dollar) {
		p.advance() // $
		op.IsImmediate = true
		if p.peekIs(token.Identifier) {
			op.CaptureVar = p.advance().Literal
		} else if p.peekIs(token.Integer) {
			op.LiteralValue = p.advance().Literal
		} else {
			op.CaptureVar = p.advance().Literal
		}
		return op
	}

	if p.peekIs(token.Integer) || p.peekIs(token.Float) {
		op.LiteralValue = p.advance().Literal
		return op
	}

	if p.peekIs(token.Identifier) || p.peekIs(token.Keyword) {
		ident := p.advance().Literal
		if isRegister(ident) {
			op.RegisterName = ident
		} else {
			op.CaptureVar = ident
		}

		if p.peekIs(token.LParen) {
			// Addressing mode: offset(reg) or (reg)(index*scale)
			ref := &ast.MemoryRefPattern{}
			if op.RegisterName != "" && isRegister(ident) {
				ref.Base = ident
			} else {
				ref.Offset = ident
			}
			p.advance() // (
			if p.peekIs(token.Identifier) {
				ref.Base = p.advance().Literal
				if p.peekIs(token.RParen) {
					p.advance() // )
					if p.peekIs(token.LParen) {
						p.advance() // (
						ref.Index = p.advance().Literal
						if p.peekIs(token.Asterisk) {
							p.advance() // *
							ref.Scale = p.advance().Literal
						}
						p.expect(token.RParen)
					}
				}
			} else if p.peekIs(token.Integer) {
				ref.Offset = p.advance().Literal
				p.expect(token.RParen)
			} else {
				p.expect(token.RParen)
			}
			op.MemoryRef = ref
			return op
		}

		return op
	}

	if p.peekIs(token.LParen) {
		ref := &ast.MemoryRefPattern{}
		p.advance() // (
		if p.peekIs(token.Identifier) {
			ref.Base = p.advance().Literal
		} else if p.peekIs(token.Integer) {
			ref.Offset = p.advance().Literal
		}
		p.expect(token.RParen)

		if p.peekIs(token.LParen) {
			p.advance() // (
			if p.peekIs(token.Identifier) {
				ref.Index = p.advance().Literal
			}
			if p.peekIs(token.Asterisk) {
				p.advance() // *
				if p.peekIs(token.Identifier) || p.peekIs(token.Integer) {
					ref.Scale = p.advance().Literal
				}
			}
			p.expect(token.RParen)
		}

		op.MemoryRef = ref
		return op
	}

	return op
}

func (p *Parser) parseGenBlock() (*ast.GenBlock, error) {
	tok := p.advance() // gen
	block := &ast.GenBlock{Tok: tok}
	p.expect(token.LBrace)
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		stmt := p.parseGenStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}
	p.expect(token.RBrace)
	return block, nil
}

func (p *Parser) parseGenStatement() ast.GenStatement {
	switch {
	case p.peekIs(token.Dollar):
		return p.parseGenExpr()
	case p.peekIs(token.Keyword, "if"):
		return p.parseGenConditional()
	case p.peekIs(token.Keyword, "for"):
		return p.parseGenLoop()
	default:
		return p.parseGenText()
	}
}

func (p *Parser) parseGenText() *ast.GenText {
	tok := p.current()
	text := ""
	for !p.peekIs(token.RBrace) && !p.peekIs(token.Dollar) &&
		!p.peekIs(token.Keyword, "if") && !p.peekIs(token.Keyword, "for") && !p.isEOF() {
		text += p.advance().Literal
		if p.peekIs(token.Semicolon) {
			text += p.advance().Literal
		}
	}
	return &ast.GenText{Tok: tok, Text: text}
}

func (p *Parser) parseGenExpr() *ast.GenExpr {
	tok := p.current()
	expr := p.parseExpression()
	return &ast.GenExpr{Tok: tok, Expr: expr}
}

func (p *Parser) parseGenConditional() *ast.GenConditional {
	tok := p.advance() // if
	cond := p.parseExpression()
	p.expect(token.LBrace)
	var body []ast.GenStatement
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		stmt := p.parseGenStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}
	p.expect(token.RBrace)

	var elseBody []ast.GenStatement
	if p.peekIs(token.Keyword, "else") {
		p.advance()
		p.expect(token.LBrace)
		for !p.peekIs(token.RBrace) && !p.isEOF() {
			stmt := p.parseGenStatement()
			if stmt != nil {
				elseBody = append(elseBody, stmt)
			}
		}
		p.expect(token.RBrace)
	}

	return &ast.GenConditional{Tok: tok, Condition: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseGenLoop() *ast.GenLoop {
	tok := p.advance() // for
	p.expect(token.LParen)

	// for (init; cond; post)
	init := ""
	for !p.peekIs(token.Semicolon) && !p.isEOF() {
		init += p.advance().Literal
	}
	p.expect(token.Semicolon)

	cond := p.parseExpression()
	p.expect(token.Semicolon)

	post := ""
	for !p.peekIs(token.RParen) && !p.isEOF() {
		post += p.advance().Literal
	}
	p.expect(token.RParen)

	p.expect(token.LBrace)
	var body []ast.GenStatement
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		stmt := p.parseGenStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}
	p.expect(token.RBrace)

	return &ast.GenLoop{Tok: tok, Init: init, Cond: cond, Post: post, Body: body}
}

func (p *Parser) parseBindBlock() (*ast.BindBlock, error) {
	tok := p.advance() // bind
	block := &ast.BindBlock{Tok: tok}
	p.expect(token.LBrace)
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		captureVar := p.expect(token.Identifier).Literal
		p.expect(token.Keyword) // as
		alias := p.expect(token.String).Literal
		block.Bindings = append(block.Bindings, ast.BindEntry{
			CaptureVar: captureVar,
			Alias:      alias,
		})
		p.expect(token.Semicolon)
	}
	p.expect(token.RBrace)
	return block, nil
}

func (p *Parser) parseVariableDeclaration() (ast.Node, error) {
	if !p.peekIs(token.Identifier) && !p.peekIs(token.Keyword) {
		return nil, nil
	}

	if p.peekIs(token.Keyword) && !isTypeKeyword(p.current().Literal) {
		return nil, nil
	}

	typ := p.parseType()
	if typ == nil {
		return nil, nil
	}

	decl := &ast.VariableDeclaration{Type: typ}

	if p.peekIs(token.Asterisk) {
		p.advance()
		decl.Type = &ast.PointerType{Base: typ}
	}

	if p.peekIs(token.Identifier) {
		decl.Name = p.advance().Literal
	}

	if p.peekIs(token.LBracket) {
		p.advance() // [
		if !p.peekIs(token.RBracket) {
			size := p.parseExpression()
			_ = size
			p.expect(token.RBracket)
		} else {
			p.advance() // ]
		}
		arrDecl := &ast.ArrayVariableDeclaration{
			Type: decl.Type,
			Name: decl.Name,
		}
		if p.peekIs(token.At) {
			p.advance() // @
			arrDecl.Offset = p.parseExpression()
		}
		p.expect(token.Semicolon)
		return arrDecl, nil
	}

	if p.peekIs(token.At) {
		p.advance() // @
		decl.Offset = p.parseExpression()
	}

	if p.peekIs(token.Assign) {
		p.advance() // =
		decl.Value = p.parseExpression()
	}

	p.expect(token.Semicolon)

	return decl, nil
}

func (p *Parser) parseType() ast.Node {
	if p.peekIs(token.Keyword, "unsigned") || p.peekIs(token.Keyword, "signed") {
		p.advance()
	}
	if p.peekIs(token.Keyword, "const") {
		p.advance()
	}

	tok := p.current()
	if !p.peekIs(token.Keyword) && !p.peekIs(token.Identifier) {
		return nil
	}

	name := p.advance()
	typ := &ast.BuiltinType{Tok: name, Name: name.Literal}

	if p.peekIs(token.Scope) {
		p.advance() // ::
		member := p.advance() // member name
		return &ast.CustomType{Tok: tok, Name: name.Literal + "::" + member.Literal}
	}

	return typ
}

func (p *Parser) parseStatement() ast.Node {
	switch {
	case p.peekIs(token.LBrace):
		return p.parseBlock()
	case p.peekIs(token.Keyword, "if"):
		return p.parseIf()
	case p.peekIs(token.Keyword, "while"):
		return p.parseWhile()
	case p.peekIs(token.Keyword, "for"):
		return p.parseFor()
	case p.peekIs(token.Keyword, "return"):
		return p.parseReturn()
	case p.peekIs(token.Keyword, "break"):
		return p.parseBreak()
	case p.peekIs(token.Keyword, "continue"):
		return p.parseContinue()
	case p.peekIs(token.Keyword, "match"):
		return p.parseMatch()
	case p.peekIs(token.Keyword, "try"):
		return p.parseTryCatch()
	case p.peekIs(token.Keyword, "fn"):
		fn, err := p.parseFunction()
		if err != nil {
			return nil
		}
		return fn
	case p.peekIs(token.Keyword, "struct"), p.peekIs(token.Keyword, "union"),
		p.peekIs(token.Keyword, "enum"), p.peekIs(token.Keyword, "bitfield"):
		return nil
	case p.peekIs(token.Semicolon):
		p.advance()
		return nil
	default:
		if p.peekIs(token.Identifier) || p.peekIs(token.Keyword) {
			if p.peekIs(token.Keyword) && isTypeKeyword(p.current().Literal) {
				decl, _ := p.parseVariableDeclaration()
				return decl
			}
		}
		expr := p.parseExpression()
		if expr == nil {
			p.advance()
			return nil
		}
		p.expect(token.Semicolon)
		return &ast.ExpressionStatement{Tok: token.Token{}, Expr: expr}
	}
}

func (p *Parser) parseIf() *ast.ConditionalStatement {
	tok := p.advance() // if
	p.expect(token.LParen)
	cond := p.parseExpression()
	p.expect(token.RParen)
	body := p.parseStatement()

	var elseBody ast.Node
	if p.peekIs(token.Keyword, "else") {
		p.advance()
		elseBody = p.parseStatement()
	}

	return &ast.ConditionalStatement{Tok: tok, Condition: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseWhile() *ast.WhileStatement {
	tok := p.advance() // while
	p.expect(token.LParen)
	cond := p.parseExpression()
	p.expect(token.RParen)
	body := p.parseStatement()

	return &ast.WhileStatement{Tok: tok, Condition: cond, Body: body}
}

func (p *Parser) parseFor() *ast.ForStatement {
	tok := p.advance() // for
	p.expect(token.LParen)

	var init ast.Node
	if !p.peekIs(token.Semicolon) {
		initVar, _ := p.parseVariableDeclaration()
		if initVar == nil {
			initExpr := p.parseExpression()
			if initExpr != nil {
				init = &ast.ExpressionStatement{Expr: initExpr}
			}
		} else {
			init = initVar
		}
	}
	p.expect(token.Semicolon)

	var cond ast.Expression
	if !p.peekIs(token.Semicolon) {
		cond = p.parseExpression()
	}
	p.expect(token.Semicolon)

	var post ast.Node
	if !p.peekIs(token.RParen) {
		postExpr := p.parseExpression()
		if postExpr != nil {
			post = &ast.ExpressionStatement{Expr: postExpr}
		}
	}
	p.expect(token.RParen)

	body := p.parseStatement()

	return &ast.ForStatement{Tok: tok, Init: init, Condition: cond, Post: post, Body: body}
}

func (p *Parser) parseReturn() *ast.ReturnStatement {
	tok := p.advance()
	var val ast.Expression
	if !p.peekIs(token.Semicolon) {
		val = p.parseExpression()
	}
	p.expect(token.Semicolon)
	return &ast.ReturnStatement{Tok: tok, Value: val}
}

func (p *Parser) parseBreak() *ast.BreakStatement {
	tok := p.advance()
	p.expect(token.Semicolon)
	return &ast.BreakStatement{Tok: tok}
}

func (p *Parser) parseContinue() *ast.ContinueStatement {
	tok := p.advance()
	p.expect(token.Semicolon)
	return &ast.ContinueStatement{Tok: tok}
}

func (p *Parser) parseMatch() *ast.MatchExpression {
	p.advance() // match
	tok := p.current()
	p.expect(token.LParen)
	val := p.parseExpression()
	p.expect(token.RParen)

	p.expect(token.LBrace)
	var cases []ast.MatchCase
	for !p.peekIs(token.RBrace) && !p.isEOF() {
		cases = append(cases, p.parseMatchCase())
	}
	p.expect(token.RBrace)

	return &ast.MatchExpression{Tok: tok, Value: val, Cases: cases}
}

func (p *Parser) parseMatchCase() ast.MatchCase {
	p.expect(token.LParen)
	var patterns []ast.MatchPattern
	for !p.peekIs(token.RParen) && !p.isEOF() {
		pattern := ast.MatchPattern{Tok: p.current()}
		pattern.Start = p.parseExpression()
		if p.peekIs(token.Range) {
			p.advance() // ...
			pattern.End = p.parseExpression()
		}
		patterns = append(patterns, pattern)
		if p.peekIs(token.Pipe) {
			p.advance()
		}
	}
	p.expect(token.RParen)
	p.expect(token.Colon)

	var body []ast.Node
	if p.peekIs(token.LBrace) {
		block := p.parseBlock()
		body = block.Statements
	} else {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	return ast.MatchCase{Patterns: patterns, Body: body}
}

func (p *Parser) parseTryCatch() *ast.TryCatchStatement {
	tok := p.advance() // try
	tryBody := p.parseStatement()

	var catchID string
	var catchBody ast.Node
	if p.peekIs(token.Keyword, "catch") {
		p.advance()
		p.expect(token.LParen)
		catchID = p.expect(token.Identifier).Literal
		p.expect(token.RParen)
		catchBody = p.parseStatement()
	}

	return &ast.TryCatchStatement{Tok: tok, TryBody: tryBody, CatchID: catchID, CatchBody: catchBody}
}

// Expression parsing with precedence climbing
func (p *Parser) parseExpression() ast.Expression {
	return p.parseTernary()
}

func (p *Parser) parseTernary() ast.Expression {
	left := p.parseOr()
	if p.peekIs(token.Question) {
		tok := p.advance()
		trueExpr := p.parseExpression()
		p.expect(token.Colon)
		falseExpr := p.parseExpression()
		return &ast.TernaryExpression{Tok: tok, Condition: left, TrueExpr: trueExpr, FalseExpr: falseExpr}
	}
	return left
}

func (p *Parser) parseOr() ast.Expression {
	left := p.parseAnd()
	for p.peekIs(token.Or) {
		op := p.advance()
		right := p.parseAnd()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseAnd() ast.Expression {
	left := p.parseXor()
	for p.peekIs(token.And) {
		op := p.advance()
		right := p.parseXor()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseXor() ast.Expression {
	left := p.parseEquality()
	for p.peekIs(token.Xor) {
		op := p.advance()
		right := p.parseEquality()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseEquality() ast.Expression {
	left := p.parseRelational()
	for p.peekIs(token.Equal) || p.peekIs(token.NotEqual) {
		op := p.advance()
		right := p.parseRelational()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseRelational() ast.Expression {
	left := p.parseShift()
	for p.peekIs(token.Less) || p.peekIs(token.Greater) || p.peekIs(token.LEqual) || p.peekIs(token.GEqual) {
		op := p.advance()
		right := p.parseShift()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseShift() ast.Expression {
	left := p.parseAdditive()
	for p.peekIs(token.LShift) || p.peekIs(token.RShift) {
		op := p.advance()
		right := p.parseAdditive()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseAdditive() ast.Expression {
	left := p.parseMultiplicative()
	for p.peekIs(token.Plus) || p.peekIs(token.Minus) {
		op := p.advance()
		right := p.parseMultiplicative()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseMultiplicative() ast.Expression {
	left := p.parseUnary()
	for p.peekIs(token.Asterisk) || p.peekIs(token.Slash) || p.peekIs(token.Percent) {
		op := p.advance()
		right := p.parseUnary()
		left = &ast.BinaryExpression{Tok: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Expression {
	if p.peekIs(token.Minus) || p.peekIs(token.Exclamation) || p.peekIs(token.Tilde) {
		op := p.advance()
		right := p.parseUnary()
		return &ast.UnaryExpression{Tok: op, Operator: op.Literal, Right: right}
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() ast.Expression {
	left := p.parsePrimary()

	for {
		if p.peekIs(token.LParen) {
			p.advance() // (
			var args []ast.Expression
			for !p.peekIs(token.RParen) && !p.isEOF() {
				args = append(args, p.parseExpression())
				if p.peekIs(token.Comma) {
					p.advance()
				}
			}
			p.expect(token.RParen)
			left = &ast.CallExpression{Function: left, Args: args}
		} else if p.peekIs(token.LBracket) {
			p.advance() // [
			index := p.parseExpression()
			p.expect(token.RBracket)
			left = &ast.IndexExpression{Left: left, Index: index}
		} else if p.peekIs(token.Dot) {
			p.advance() // .
			member := p.expect(token.Identifier)
			left = &ast.IndexExpression{Left: left, Index: &ast.Identifier{Name: member.Literal}}
		} else if p.peekIs(token.Scope) {
			p.advance() // ::
			member := p.expect(token.Identifier)
			left = &ast.ScopeExpression{Left: left, Member: member.Literal}
		} else if p.peekIs(token.Keyword, "as") {
			p.advance() // as
			typ := p.parseType()
			left = &ast.CastExpression{Type: typ, Expr: left}
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parsePrimary() ast.Expression {
	switch {
	case p.peekIs(token.Integer):
		tok := p.advance()
		return &ast.IntegerLiteral{Tok: tok, Value: 0}
	case p.peekIs(token.Float):
		tok := p.advance()
		return &ast.FloatLiteral{Tok: tok, Value: 0}
	case p.peekIs(token.String):
		tok := p.advance()
		return &ast.StringLiteral{Tok: tok, Value: tok.Literal}
	case p.peekIs(token.Char):
		tok := p.advance()
		return &ast.CharLiteral{Tok: tok}
	case p.peekIs(token.Keyword, "true"):
		return &ast.BoolLiteral{Tok: p.advance(), Value: true}
	case p.peekIs(token.Keyword, "false"):
		return &ast.BoolLiteral{Tok: p.advance(), Value: false}
	case p.peekIs(token.Keyword, "null"):
		return &ast.NullLiteral{Tok: p.advance()}
	case p.peekIs(token.Identifier) || p.peekIs(token.Keyword):
		tok := p.advance()
		return &ast.Identifier{Tok: tok, Name: tok.Literal}
	case p.peekIs(token.Keyword, "sizeof") || p.peekIs(token.Keyword, "addressof") || p.peekIs(token.Keyword, "typenameof"):
		op := p.advance()
		p.expect(token.LParen)
		expr := p.parseExpression()
		p.expect(token.RParen)
		return &ast.TypeOperator{Tok: op, Operator: op.Literal, Expr: expr}
	case p.peekIs(token.Keyword, "match"):
		return p.parseMatch()
	case p.peekIs(token.LParen):
		p.advance() // (
		// Could be a cast: type(expr) or parenthesized expression
		if p.peekIs(token.Keyword) || p.peekIs(token.Identifier) {
			mark := p.mark()
			typ := p.parseType()
			if typ != nil && p.peekIs(token.RParen) {
				p.advance() // )
				if p.peekIs(token.LParen) || p.peekIs(token.Integer) || p.peekIs(token.Identifier) {
					expr := p.parsePrimary()
					return &ast.CastExpression{Type: typ, Expr: expr}
				}
			}
			p.reset(mark)
		}
		expr := p.parseExpression()
		p.expect(token.RParen)
		return expr
	case p.peekIs(token.Dollar):
		p.advance() // $
		expr := p.parsePrimary()
		return &ast.UnaryExpression{Operator: "$", Right: expr}
	default:
		return nil
	}
}

// --- Helpers ---

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() token.Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *Parser) isEOF() bool {
	return p.pos >= len(p.tokens) || p.tokens[p.pos].Type == token.EOF
}

func (p *Parser) peekIs(typ token.Type, lit ...string) bool {
	tok := p.current()
	if tok.Type != typ {
		return false
	}
	if len(lit) > 0 && tok.Literal != lit[0] {
		return false
	}
	return true
}

func (p *Parser) expect(typ token.Type) token.Token {
	tok := p.current()
	if tok.Type != typ {
		return tok
	}
	return p.advance()
}

func (p *Parser) skipSemicolon() {
	if p.peekIs(token.Semicolon) {
		p.advance()
	}
}

type mark struct {
	pos int
}

func (p *Parser) mark() mark {
	return mark{pos: p.pos}
}

func (p *Parser) reset(m mark) {
	p.pos = m.pos
}

func baseName(typ ast.Node) string {
	if bt, ok := typ.(*ast.BuiltinType); ok {
		return bt.Name
	}
	return ""
}

func isRegister(name string) bool {
	regs := map[string]bool{
		"AL": true, "CL": true, "DL": true, "BL": true,
		"AH": true, "CH": true, "DH": true, "BH": true,
		"AX": true, "CX": true, "DX": true, "BX": true,
		"SI": true, "DI": true, "BP": true, "SP": true,
		"R8": true, "R9": true, "R10": true, "R11": true,
		"R12": true, "R13": true, "R14": true, "R15": true,
		"R8B": true, "R9B": true, "R10B": true, "R11B": true,
		"R12B": true, "R13B": true, "R14B": true, "R15B": true,
		"R8W": true, "R9W": true, "R10W": true, "R11W": true,
		"R12W": true, "R13W": true, "R14W": true, "R15W": true,
		"RAX": true, "RBX": true, "RCX": true, "RDX": true,
		"RSI": true, "RDI": true, "RBP": true, "RSP": true,
		"X0": true, "X1": true, "X2": true, "X3": true, "X4": true,
		"X5": true, "X6": true, "X7": true, "X8": true, "X9": true,
		"X10": true, "X11": true, "X12": true, "X13": true, "X14": true, "X15": true,
		"Y0": true, "Y1": true, "Y2": true, "Y3": true,
		"FS": true, "GS": true, "CS": true, "DS": true, "ES": true, "SS": true,
		"FP": true, "SP_VIRTUAL": true, "SB": true, "PC": true,
	}
	return regs[name]
}

func isTypeKeyword(kw string) bool {
	switch kw {
	case "u8", "u16", "u24", "u32", "u48", "u64", "u96", "u128",
		"s8", "s16", "s24", "s32", "s48", "s64", "s96", "s128",
		"char", "char16", "bool", "float", "double", "str", "auto", "padding",
		"signed", "unsigned", "const":
		return true
	}
	return false
}
