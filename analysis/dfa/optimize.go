package dfa

func OptimizeBlock(stmts []Statement) []Statement {
	if len(stmts) == 0 {
		return stmts
	}

	result := make([]Statement, 0, len(stmts))
	for _, stmt := range stmts {
		switch stmt.Kind {
		case StmtAssign, StmtReturn, StmtStore:
			if stmt.Value != nil {
				stmt.Value = simplifyValue(stmt.Value)
			}
		case StmtCall, StmtExpr:
			if stmt.Value != nil {
				stmt.Value = simplifyValue(stmt.Value)
			}
			for i, arg := range stmt.Args {
				stmt.Args[i] = simplifyValue(arg)
			}
		}
		result = append(result, stmt)
	}

	result = eliminateDeadStores(result)
	return result
}

func simplifyValue(v *Value) *Value {
	if v == nil {
		return nil
	}

	switch v.Kind {
	case ValOp:
		v.Left = simplifyValue(v.Left)
		v.Right = simplifyValue(v.Right)

		if v.Left != nil && v.Right != nil {
			if v.Left.Kind == ValConst && v.Right.Kind == ValConst {
				result := evalConstBinOp(v.Op, v.Left.Const, v.Right.Const)
				if result != nil {
					return result
				}
			}

			if v.Op == "^" && v.Left.Kind == ValReg && v.Right.Kind == ValReg && v.Left.Reg == v.Right.Reg {
				return ConstValue(0)
			}

			if v.Op == "-" && v.Left.Kind == ValReg && v.Right.Kind == ValReg && v.Left.Reg == v.Right.Reg {
				return ConstValue(0)
			}

			if v.Op == "+" {
				if v.Left.Kind == ValConst {
					v.Left, v.Right = v.Right, v.Left
				}
				if v.Left != nil && v.Left.Kind == ValConst && v.Left.Const == 0 {
					return v.Right
				}
				if v.Right != nil && v.Right.Kind == ValConst && v.Right.Const == 0 {
					return v.Left
				}
			}
			if v.Op == "*" {
				if v.Right.Kind == ValConst && v.Right.Const == 1 {
					return v.Left
				}
				if v.Left.Kind == ValConst && v.Left.Const == 1 {
					return v.Right
				}
				if v.Right.Kind == ValConst && v.Right.Const == 0 {
					return ConstValue(0)
				}
			}
		}

	case ValUnary:
		v.Left = simplifyValue(v.Left)
		if v.Left != nil {
			if v.Op == "-" && v.Left.Kind == ValConst {
				return ConstValue(^v.Left.Const + 1)
			}
			if v.Op == "^" && v.Left.Kind == ValConst {
				return ConstValue(^v.Left.Const)
			}
		}
	}

	return v
}

func evalConstBinOp(op string, a, b uint64) *Value {
	switch op {
	case "+":
		return ConstValue(a + b)
	case "-":
		return ConstValue(a - b)
	case "*":
		return ConstValue(a * b)
	case "/":
		if b != 0 {
			return ConstValue(a / b)
		}
	case "&":
		return ConstValue(a & b)
	case "|":
		return ConstValue(a | b)
	case "^":
		return ConstValue(a ^ b)
	case "<<":
		return ConstValue(a << b)
	case ">>":
		return ConstValue(a >> b)
	}
	return nil
}

func eliminateDeadStores(stmts []Statement) []Statement {
	usedVars := make(map[string]bool)

	for _, stmt := range stmts {
		switch stmt.Kind {
		case StmtCall, StmtStore, StmtReturn, StmtIf:
			if stmt.Value != nil {
				collectUses(stmt.Value, usedVars)
			}
		case StmtAssign:
			usedVars[stmt.Dst] = false
			if stmt.Value != nil {
				collectUses(stmt.Value, usedVars)
			}
		}
	}

	for i := len(stmts) - 1; i >= 0; i-- {
		if stmts[i].Kind == StmtAssign {
			name := stmts[i].Dst
			for j := i + 1; j < len(stmts); j++ {
				if stmts[j].Kind == StmtAssign && stmts[j].Dst == name {
					usedVars[name] = false
					break
				}
				if usesVar(stmts[j], name) {
					usedVars[name] = true
					break
				}
			}
		}
	}

	var result []Statement
	for _, stmt := range stmts {
		if stmt.Kind == StmtAssign && !usedVars[stmt.Dst] {
			continue
		}
		result = append(result, stmt)
	}
	return result
}

func collectUses(v *Value, used map[string]bool) {
	if v == nil {
		return
	}
	switch v.Kind {
	case ValReg:
		used[v.Reg] = true
	case ValOp:
		collectUses(v.Left, used)
		collectUses(v.Right, used)
	case ValCall:
		for _, a := range v.Args {
			collectUses(a, used)
		}
	}
}

func usesVar(stmt Statement, name string) bool {
	return valueUsesVar(stmt.Value, name)
}

func valueUsesVar(v *Value, name string) bool {
	if v == nil {
		return false
	}
	switch v.Kind {
	case ValReg:
		return v.Reg == name
	case ValOp:
		return valueUsesVar(v.Left, name) || valueUsesVar(v.Right, name)
	case ValCall:
		for _, a := range v.Args {
			if valueUsesVar(a, name) {
				return true
			}
		}
	}
	return false
}
