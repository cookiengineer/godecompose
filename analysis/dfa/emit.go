package dfa

import (
	"fmt"
	"strings"
)

func EmitStatements(stmts []Statement, indent string) string {
	var buf strings.Builder

	for _, stmt := range stmts {
		line := emitStatement(stmt)
		if line == "" {
			continue
		}
		buf.WriteString(indent)
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return buf.String()
}

func EmitSafeStatements(stmts []Statement, indent string) string {
	var buf strings.Builder

	for _, stmt := range stmts {
		switch stmt.Kind {
		case StmtGoto, StmtLabel, StmtIf:
			continue
		}
		line := emitStatement(stmt)
		if line == "" {
			continue
		}
		buf.WriteString(indent)
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return buf.String()
}

func emitStatement(stmt Statement) string {
	switch stmt.Kind {
	case StmtAssign:
		return emitAssign(stmt)
	case StmtStore:
		return emitStore(stmt)
	case StmtCall:
		return emitCall(stmt)
	case StmtReturn:
		return emitReturn(stmt)
	case StmtIf:
		return emitIf(stmt)
	case StmtGoto:
		return emitGoto(stmt)
	case StmtLabel:
		return emitLabel(stmt)
	case StmtExpr:
		return emitExpr(stmt)
	case StmtComment:
		return "// " + stmt.Dst
	}
	return ""
}

func CleanAssemblyIdent(s string) string {
	s = strings.TrimSuffix(s, "(SB)")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "+", "_")
	s = strings.ReplaceAll(s, "-", "n")
	if idx := strings.LastIndexByte(s, '.'); idx >= 0 {
		s = s[idx+1:]
	}
	return s
}

func cleanAssemblyIdent(s string) string {
	return CleanAssemblyIdent(s)
}

func goLabel(s string) string {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimSpace(s)
	return "L_" + s
}

func emitAssign(stmt Statement) string {
	if stmt.Value == nil {
		return stmt.Dst + " = _"
	}
	if stmt.Value.Kind == ValCall {
		return stmt.Dst + " := " + emitValueGo(stmt.Value)
	}
	return stmt.Dst + " = " + emitValueGo(stmt.Value)
}

func emitStore(stmt Statement) string {
	if stmt.Value == nil {
		return ""
	}
	dst := cleanAssemblyIdent(stmt.Dst)
	return fmt.Sprintf("*%s = %s", dst, emitValueGo(stmt.Value))
}

func emitCall(stmt Statement) string {
	fn := stmt.Func
	if fn == "" && stmt.Value != nil {
		fn = stmt.Value.Func
	}
	if fn == "" {
		return "// call ?"
	}
	fn = cleanAssemblyIdent(fn)
	return shortFunc(fn) + "()"
}

func emitReturn(stmt Statement) string {
	if stmt.Value != nil && stmt.Value.Kind != ValUnknown {
		return "return " + emitValueGo(stmt.Value)
	}
	return "return"
}

func emitIf(stmt Statement) string {
	cond := stmt.Cond
	if cond == "" || cond == "_" {
		cond = "true"
	}
	return "// if " + cond + " goto " + goLabel(stmt.Label)
}

func emitGoto(stmt Statement) string {
	return "goto " + goLabel(stmt.Label)
}

func emitLabel(stmt Statement) string {
	return goLabel(stmt.Label) + ":"
}

func emitExpr(stmt Statement) string {
	if stmt.Value != nil {
		return emitValueGo(stmt.Value)
	}
	return "_"
}
