package dfa

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func valueToExpr(v *Value) ast.Expr {
	if v == nil {
		return ast.NewIdent("_")
	}
	switch v.Kind {
	case ValConst:
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprintf("%d", v.Const),
		}
	case ValReg:
		return ast.NewIdent(strings.ToLower(v.Reg))
	case ValOp:
		left := valueToExpr(v.Left)
		right := valueToExpr(v.Right)
		op := binOpToToken(v.Op)
		return &ast.BinaryExpr{X: left, Op: op, Y: right}
	case ValUnary:
		left := valueToExpr(v.Left)
		op := unaryOpToToken(v.Op)
		return &ast.UnaryExpr{Op: op, X: left}
	case ValLoad:
		return &ast.StarExpr{X: memToExpr(v.Mem)}
	case ValAddrOf:
		if v.Mem != nil {
			return &ast.UnaryExpr{Op: token.AND, X: memToExpr(v.Mem)}
		}
		if v.Left != nil {
			return &ast.UnaryExpr{Op: token.AND, X: valueToExpr(v.Left)}
		}
		return ast.NewIdent("_")
	case ValCall:
		args := make([]ast.Expr, len(v.Args))
		for i, a := range v.Args {
			args[i] = valueToExpr(a)
		}
		fn := cleanCallIdent(v.Func)
		return &ast.CallExpr{Fun: ast.NewIdent(fn), Args: args}
	case ValSym:
		if v.Left != nil {
			return &ast.BinaryExpr{
				X:  ast.NewIdent(v.Func),
				Op: token.ADD,
				Y:  valueToExpr(v.Left),
			}
		}
		return ast.NewIdent(v.Func)
	}
	return ast.NewIdent("_")
}

func binOpToToken(op string) token.Token {
	switch op {
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	case "&":
		return token.AND
	case "|":
		return token.OR
	case "^":
		return token.XOR
	case "<<":
		return token.SHL
	case ">>":
		return token.SHR
	case "==":
		return token.EQL
	case "!=":
		return token.NEQ
	case "<":
		return token.LSS
	case "<=":
		return token.LEQ
	case ">":
		return token.GTR
	case ">=":
		return token.GEQ
	}
	return token.ILLEGAL
}

func unaryOpToToken(op string) token.Token {
	switch op {
	case "-":
		return token.SUB
	case "^":
		return token.XOR
	case "!":
		return token.NOT
	}
	return token.ILLEGAL
}

func memToExpr(m *MemRef) ast.Expr {
	if m == nil {
		return ast.NewIdent("_")
	}
	if m.Base == "SB" {
		if m.Symbol != "" {
			name := CleanAssemblyIdent(m.Symbol)
			if m.Offset != 0 {
				return &ast.BinaryExpr{
					X:  ast.NewIdent(name),
					Op: token.ADD,
					Y:  &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", m.Offset)},
				}
			}
			return ast.NewIdent(name)
		}
		return ast.NewIdent("sb_unk")
	}
	if m.Base == "SP" {
		return ast.NewIdent(fmt.Sprintf("stack_0x%x", m.Offset))
	}
	return ast.NewIdent(strings.ToLower(m.Base))
}

func cleanCallIdent(fn string) string {
	fn = CleanAssemblyIdent(fn)
	if idx := strings.LastIndexByte(fn, '.'); idx >= 0 {
		return fn[idx+1:]
	}
	if idx := strings.LastIndexByte(fn, '/'); idx >= 0 {
		return fn[idx+1:]
	}
	return fn
}

func StatementToAST(stmt Statement) ast.Stmt {
	switch stmt.Kind {
	case StmtAssign:
		return assignToAST(stmt)
	case StmtStore:
		return storeToAST(stmt)
	case StmtCall:
		return callToAST(stmt)
	case StmtReturn:
		return returnToAST(stmt)
	case StmtIf:
		return ifToAST(stmt)
	case StmtGoto:
		return gotoToAST(stmt)
	case StmtLabel:
		return labelToAST(stmt)
	case StmtExpr:
		return exprToAST(stmt)
	case StmtComment:
		return &ast.ExprStmt{X: &ast.Ident{
			Name: "// " + stmt.Dst,
		}}
	}
	return &ast.EmptyStmt{}
}

func assignToAST(stmt Statement) ast.Stmt {
	dstIdent := ast.NewIdent(stmt.Dst)
	if stmt.Value == nil {
		return &ast.AssignStmt{
			Lhs: []ast.Expr{dstIdent},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent("nil")},
		}
	}
	if stmt.Value.Kind == ValCall {
		return &ast.AssignStmt{
			Lhs: []ast.Expr{dstIdent},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{valueToExpr(stmt.Value)},
		}
	}
	return &ast.AssignStmt{
		Lhs: []ast.Expr{dstIdent},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{valueToExpr(stmt.Value)},
	}
}

func storeToAST(stmt Statement) ast.Stmt {
	if stmt.Value == nil {
		return &ast.EmptyStmt{}
	}
	dst := CleanAssemblyIdent(stmt.Dst)
	return &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.StarExpr{X: ast.NewIdent(dst)}},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{valueToExpr(stmt.Value)},
	}
}

func callToAST(stmt Statement) ast.Stmt {
	fn := stmt.Func
	if fn == "" && stmt.Value != nil {
		fn = stmt.Value.Func
	}
	if fn == "" {
		return &ast.EmptyStmt{}
	}
	return &ast.ExprStmt{X: &ast.CallExpr{
		Fun: ast.NewIdent(cleanCallIdent(fn)),
	}}
}

func returnToAST(stmt Statement) ast.Stmt {
	if stmt.Value != nil && stmt.Value.Kind != ValUnknown {
		return &ast.ReturnStmt{Results: []ast.Expr{valueToExpr(stmt.Value)}}
	}
	return &ast.ReturnStmt{}
}

func ifToAST(stmt Statement) ast.Stmt {
	cond := stmt.Cond
	if cond == "" || cond == "_" {
		cond = "true"
	}
	label := cleanLabel(stmt.Label)
	condExpr := parseCondExpr(cond)
	body := &ast.BlockStmt{List: []ast.Stmt{
		&ast.BranchStmt{Tok: token.GOTO, Label: ast.NewIdent(label)},
	}}
	return &ast.IfStmt{
		Cond: condExpr,
		Body: body,
	}
}

func gotoToAST(stmt Statement) ast.Stmt {
	return &ast.BranchStmt{
		Tok:   token.GOTO,
		Label: ast.NewIdent(cleanLabel(stmt.Label)),
	}
}

func labelToAST(stmt Statement) ast.Stmt {
	return &ast.LabeledStmt{
		Label: ast.NewIdent(cleanLabel(stmt.Label)),
		Stmt:  &ast.EmptyStmt{},
	}
}

func exprToAST(stmt Statement) ast.Stmt {
	if stmt.Value != nil {
		return &ast.ExprStmt{X: valueToExpr(stmt.Value)}
	}
	return &ast.EmptyStmt{}
}

func parseCondExpr(cond string) ast.Expr {
	cond = strings.TrimSpace(cond)
	if cond == "" || cond == "_" || cond == "true" {
		return ast.NewIdent("true")
	}

	if expr, err := parseGoExpr(cond); err == nil {
		return expr
	}

	if idx := strings.Index(cond, "=="); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+2:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.EQL,
			Y:  parseIfIntLiteral(right),
		}
	}
	if idx := strings.Index(cond, "!="); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+2:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.NEQ,
			Y:  parseIfIntLiteral(right),
		}
	}
	if idx := strings.Index(cond, "<="); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+2:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.LEQ,
			Y:  parseIfIntLiteral(right),
		}
	}
	if idx := strings.Index(cond, ">="); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+2:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.GEQ,
			Y:  parseIfIntLiteral(right),
		}
	}
	if idx := strings.Index(cond, "<"); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+1:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.LSS,
			Y:  parseIfIntLiteral(right),
		}
	}
	if idx := strings.Index(cond, ">"); idx >= 0 {
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+1:])
		return &ast.BinaryExpr{
			X:  ast.NewIdent(left),
			Op: token.GTR,
			Y:  parseIfIntLiteral(right),
		}
	}

	return ast.NewIdent(cond)
}

func parseIfIntLiteral(s string) ast.Expr {
	s = strings.TrimSpace(s)
	if s == "nil" {
		return ast.NewIdent("nil")
	}
	if val, err := strconv.ParseInt(s, 0, 64); err == nil {
		return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", val)}
	}
	return ast.NewIdent(s)
}

func parseGoExpr(s string) (ast.Expr, error) {
	expr, err := tokenizeGoExpr(s)
	return expr, err
}

func tokenizeGoExpr(s string) (ast.Expr, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "nil" || s == "true" || s == "false" {
		return ast.NewIdent(s), nil
	}
	if n, err := strconv.ParseInt(s, 0, 64); err == nil {
		return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", n)}, nil
	}
	if s[0] == '"' || s[0] == '`' {
		return &ast.BasicLit{Kind: token.STRING, Value: s}, nil
	}
	return ast.NewIdent(s), nil
}

func cleanLabel(s string) string {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimSpace(s)
	return "L_" + s
}

func EmitStatementsAST(stmts []Statement) []ast.Stmt {
	var result []ast.Stmt
	for _, stmt := range stmts {
		result = append(result, StatementToAST(stmt))
	}
	return result
}

func EmitSafeStatementsAST(stmts []Statement) []ast.Stmt {
	var result []ast.Stmt
	for _, stmt := range stmts {
		switch stmt.Kind {
		case StmtGoto, StmtLabel, StmtIf:
			continue
		}
		result = append(result, StatementToAST(stmt))
	}
	return result
}
