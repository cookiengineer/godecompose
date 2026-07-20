package dfa

import (
	"fmt"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
)

type BlockAnalyzer struct {
	state  *BlockState
	params []Param
}

func NewBlockAnalyzer(params []Param) *BlockAnalyzer {
	s := NewBlockState()
	for _, p := range params {
		s.AddParam(p)
	}
	return &BlockAnalyzer{state: s, params: params}
}

func (a *BlockAnalyzer) State() *BlockState {
	return a.state
}

func (a *BlockAnalyzer) Analyze(insts []disasm.Instruction) {
	for _, inst := range insts {
		a.translateInstruction(inst)
	}
}

func (a *BlockAnalyzer) translateInstruction(inst disasm.Instruction) {
	subs := ParseInstructions(inst.GoSyntax)
	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		a.translate(sub.Opcode, sub.Operands, inst)
	}
}

func (a *BlockAnalyzer) translate(opcode string, operands []string, inst disasm.Instruction) {
	switch {
	case opcode == "MOVQ" || opcode == "MOVL" || opcode == "MOVW" || opcode == "MOVB",
		opcode == "MOVUPS", opcode == "MOVSD", opcode == "MOVSS":
		a.translateMOV(opcode, operands)

	case opcode == "LEAQ" || opcode == "LEAL":
		a.translateLEA(opcode, operands)

	case opcode == "ADDQ" || opcode == "ADDL":
		a.translateBinOp("+", operands)
	case opcode == "SUBQ" || opcode == "SUBL":
		a.translateBinOp("-", operands)
	case opcode == "IMULQ" || opcode == "IMULL" || opcode == "MULQ":
		a.translateBinOp("*", operands)
	case opcode == "ANDQ" || opcode == "ANDL" || opcode == "ANDB":
		a.translateBinOp("&", operands)
	case opcode == "ORQ" || opcode == "ORL" || opcode == "ORB":
		a.translateBinOp("|", operands)
	case opcode == "XORQ" || opcode == "XORL" || opcode == "XORB":
		a.translateBinOp("^", operands)
	case opcode == "SHLQ" || opcode == "SHLL", opcode == "SALQ" || opcode == "SALL":
		a.translateBinOp("<<", operands)
	case opcode == "SHRQ" || opcode == "SHRL":
		a.translateBinOp(">>", operands)
	case opcode == "INCQ" || opcode == "INCL":
		a.translateIncDec("+", operands)
	case opcode == "DECQ" || opcode == "DECL":
		a.translateIncDec("-", operands)
	case opcode == "NEGQ" || opcode == "NEGL":
		a.translateUnary("-", operands)
	case opcode == "NOTQ" || opcode == "NOTL":
		a.translateUnary("^", operands)

	case opcode == "CMPQ" || opcode == "CMPL" || opcode == "CMP" || opcode == "CMPB":
		a.translateCMP(operands)
	case opcode == "TESTQ" || opcode == "TESTL" || opcode == "TEST":
		a.translateTEST(operands)

	case opcode == "CALL":
		a.translateCALL(inst)

	case opcode == "RET":
		retVal := a.state.GetReg("AX")
		if retVal != nil && retVal.Kind != ValReg {
			a.state.AddStmt(Statement{Kind: StmtReturn, Address: inst.Address, Value: retVal})
		} else {
			a.state.AddStmt(Statement{Kind: StmtReturn, Address: inst.Address})
		}

	case opcode == "JMP":
		a.state.AddStmt(Statement{
			Kind:    StmtGoto,
			Address: inst.Address,
			Label:   fmtHex(inst.BranchTarget),
		})

	case opcode == "JEQ" || opcode == "JNE" || opcode == "JGT" || opcode == "JLT" ||
		opcode == "JGE" || opcode == "JLE" || opcode == "JHI" || opcode == "JLS":
		a.translateCondJump(opcode, inst)

	case opcode == "PUSHQ" || opcode == "POPQ":
		// stack manipulation — track basic push/pop
	case opcode == "SETEQ" || opcode == "SETNE" || opcode == "SETGT":
		a.translateSET(opcode, operands)
	}
}

func (a *BlockAnalyzer) translateMOV(_ string, operands []string) {
	if len(operands) < 2 {
		return
	}
	srcMem, srcVal := ParseOperand(operands[0])
	dstMem, _ := ParseOperand(operands[1])

	src := a.resolveSrc(srcMem, srcVal)
	dstReg := dstRegisterName(operands[1])

	if dstReg != "" {
		a.state.SetReg(dstReg, src)
		return
	}

	if dstMem != nil && dstMem.Base != "" {
		if dstMem.Base == "SP" {
			offset := dstMem.Offset
			name := a.state.GetStackSlot(offset)
			a.state.AddStmt(Statement{
				Kind: StmtAssign, Dst: name, Value: src,
			})
			return
		}
		a.state.AddStmt(Statement{
			Kind: StmtStore, Dst: dstMem.String(), Value: src,
		})
	}
}

func (a *BlockAnalyzer) translateLEA(_ string, operands []string) {
	if len(operands) < 2 {
		return
	}
	mem, _ := ParseOperand(operands[0])
	dstReg := dstRegisterName(operands[1])

	if dstReg != "" && mem != nil {
		v := &Value{Kind: ValAddrOf, Mem: mem, Size: 8}
		a.state.SetReg(dstReg, v)
	}
}

func (a *BlockAnalyzer) translateBinOp(op string, operands []string) {
	if len(operands) < 2 {
		return
	}
	srcMem, srcVal := ParseOperand(operands[0])
	dstReg := dstRegisterName(operands[1])

	if dstReg == "" {
		return
	}

	src := a.resolveSrc(srcMem, srcVal)
	dst := a.state.GetReg(dstReg)
	result := BinOpValue(op, dst, src)
	a.state.SetReg(dstReg, result)
}

func (a *BlockAnalyzer) translateIncDec(op string, operands []string) {
	if len(operands) < 1 {
		return
	}
	dstReg := dstRegisterName(operands[0])
	if dstReg == "" {
		return
	}
	dst := a.state.GetReg(dstReg)
	one := ConstValue(1)
	result := BinOpValue(op, dst, one)
	a.state.SetReg(dstReg, result)
	a.state.AddStmt(Statement{
		Kind: StmtAssign, Dst: dstReg, Value: result,
	})
}

func (a *BlockAnalyzer) translateUnary(op string, operands []string) {
	if len(operands) < 1 {
		return
	}
	dstReg := dstRegisterName(operands[0])
	if dstReg == "" {
		return
	}
	dst := a.state.GetReg(dstReg)
	result := UnaryValue(op, dst)
	a.state.SetReg(dstReg, result)
}

func (a *BlockAnalyzer) translateCMP(operands []string) {
	if len(operands) < 2 {
		return
	}
	_, srcVal := ParseOperand(operands[0])
	dstMem, dstVal := ParseOperand(operands[1])

	var left, right *Value

	if dstMem != nil && dstMem.Base != "" {
		left = &Value{Kind: ValLoad, Mem: dstMem}
	} else if dstVal != nil {
		left = dstVal
	} else {
		left = a.resolveSrc(dstMem, dstVal)
	}

	if srcVal != nil {
		right = srcVal
	} else {
		right = a.resolveSrc(nil, srcVal)
	}

	a.state.lastCmp = &CmpInfo{Left: left, Right: right}
}

func (a *BlockAnalyzer) translateTEST(operands []string) {
	if len(operands) < 2 {
		return
	}
	_, leftVal := ParseOperand(operands[0])
	_, rightVal := ParseOperand(operands[1])

	left := a.resolveSrc(nil, leftVal)
	right := a.resolveSrc(nil, rightVal)

	a.state.lastCmp = &CmpInfo{Left: left, Right: right}
}

func (a *BlockAnalyzer) translateCALL(inst disasm.Instruction) {
	target := ExtractCallTarget(inst.GoSyntax)
	val := &Value{Kind: ValCall, Func: target, Size: 8}

	a.state.SetReg("AX", val)
	a.state.AddStmt(Statement{
		Kind: StmtCall, Address: inst.Address, Func: target, Value: val,
	})
}

func (a *BlockAnalyzer) translateCondJump(opcode string, inst disasm.Instruction) {
	cond := a.buildCondition(opcode)
	a.state.AddStmt(Statement{
		Kind:    StmtIf,
		Address: inst.Address,
		Cond:    cond,
		Label:   fmtHex(inst.BranchTarget),
	})
}

func (a *BlockAnalyzer) translateSET(opcode string, operands []string) {
	if len(operands) < 1 {
		return
	}
	dstReg := dstRegisterName(operands[0])
	if dstReg == "" {
		return
	}
	result := ConstValue(1)
	a.state.SetReg(dstReg, result)
}

func (a *BlockAnalyzer) buildCondition(jccOpcode string) string {
	cmp := a.state.lastCmp
	if cmp == nil || cmp.Left == nil || cmp.Right == nil {
		return "_"
	}

	left := simplifyValue(cmp.Left)
	right := simplifyValue(cmp.Right)

	leftStr := emitValueGo(left)
	rightStr := emitValueGo(right)

	if leftStr == rightStr || (left.Kind == ValReg && right.Kind == ValReg && left.Reg == right.Reg) {
		switch {
		case jccOpcode == "JEQ":
			return leftStr + " == 0"
		case jccOpcode == "JNE":
			return leftStr + " != 0"
		}
	}

	switch {
	case jccOpcode == "JEQ":
		return leftStr + " == " + rightStr
	case jccOpcode == "JNE":
		return leftStr + " != " + rightStr
	case jccOpcode == "JGT":
		return leftStr + " > " + rightStr
	case jccOpcode == "JLT":
		return leftStr + " < " + rightStr
	case jccOpcode == "JGE":
		return leftStr + " >= " + rightStr
	case jccOpcode == "JLE":
		return leftStr + " <= " + rightStr
	case jccOpcode == "JHI":
		return leftStr + " > " + rightStr
	case jccOpcode == "JLS":
		return leftStr + " <= " + rightStr
	}
	return leftStr + " ? " + rightStr
}

func (a *BlockAnalyzer) resolveSrc(mem *MemRef, val *Value) *Value {
	if mem != nil && mem.Base != "" {
		if mem.Base == "SP" {
			name := a.state.StackVarName(mem.Offset)
			return &Value{Kind: ValReg, Reg: name, Size: 8}
		}
		return &Value{Kind: ValLoad, Mem: mem, Size: 8}
	}
	if val != nil {
		if val.Kind == ValReg {
			return a.state.GetReg(val.Reg)
		}
		return val
	}
	return nil
}

func dstRegisterName(operand string) string {
	if isRegister(operand) {
		return operand
	}
	return ""
}

func emitValueGo(v *Value) string {
	if v == nil {
		return "_"
	}
	switch v.Kind {
	case ValConst:
		return fmtHex(v.Const)
	case ValReg:
		return strings.ToLower(v.Reg)
	case ValOp:
		return "(" + emitValueGo(v.Left) + " " + v.Op + " " + emitValueGo(v.Right) + ")"
	case ValUnary:
		return v.Op + emitValueGo(v.Left)
	case ValLoad:
		return "*" + emitLoadName(v.Mem)
	case ValAddrOf:
		if v.Mem != nil {
			return "&" + emitLoadName(v.Mem)
		}
		if v.Left != nil {
			return "&" + emitValueGo(v.Left)
		}
		return "&_"
	case ValCall:
		args := make([]string, len(v.Args))
		for i, a := range v.Args {
			args[i] = emitValueGo(a)
		}
		return shortFunc(v.Func) + "(" + join(args, ", ") + ")"
	}
	return "_"
}

func emitLoadName(m *MemRef) string {
	if m == nil {
		return "?"
	}
	if m.Base == "SB" {
		if m.Symbol != "" {
			name := CleanAssemblyIdent(m.Symbol)
			if m.Offset != 0 {
				return fmt.Sprintf("%s_%d", name, m.Offset)
			}
			return name
		}
		return "sb_unk"
	}
	if m.Base == "SP" {
		return fmt.Sprintf("stack_0x%x", m.Offset)
	}
	return m.Base
}

func shortFunc(fn string) string {
	if idx := strings.LastIndexByte(fn, '.'); idx >= 0 {
		return fn[idx+1:]
	}
	if idx := strings.LastIndexByte(fn, '/'); idx >= 0 {
		return fn[idx+1:]
	}
	return fn
}
