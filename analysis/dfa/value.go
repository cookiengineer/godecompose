package dfa

import "fmt"

type ValueKind int

const (
	ValUnknown ValueKind = iota
	ValConst
	ValReg
	ValOp
	ValUnary
	ValLoad
	ValStore
	ValAddrOf
	ValSym
	ValCall
)

func (k ValueKind) String() string {
	switch k {
	case ValConst:
		return "const"
	case ValReg:
		return "reg"
	case ValOp:
		return "op"
	case ValUnary:
		return "unary"
	case ValLoad:
		return "load"
	case ValStore:
		return "store"
	case ValAddrOf:
		return "addrof"
	case ValSym:
		return "sym"
	case ValCall:
		return "call"
	}
	return "unknown"
}

type Value struct {
	ID       int
	Kind     ValueKind
	Const    uint64
	Reg      string
	Op       string
	Left     *Value
	Right    *Value
	Func     string
	Args     []*Value
	Mem      *MemRef
	Size     int
	TypeHint string
}

func ConstValue(c uint64) *Value {
	return &Value{Kind: ValConst, Const: c, Size: 8}
}

func RegValue(reg string) *Value {
	return &Value{Kind: ValReg, Reg: reg, Size: 8}
}

func BinOpValue(op string, left, right *Value) *Value {
	return &Value{Kind: ValOp, Op: op, Left: left, Right: right, Size: 8}
}

func UnaryValue(op string, left *Value) *Value {
	return &Value{Kind: ValUnary, Op: op, Left: left, Size: 8}
}

func LoadValue(mem *MemRef) *Value {
	return &Value{Kind: ValLoad, Mem: mem, Size: 8}
}

func CallValue(fn string, args []*Value) *Value {
	return &Value{Kind: ValCall, Func: fn, Args: args, Size: 8}
}

func (v *Value) GoString() string {
	switch v.Kind {
	case ValConst:
		return fmt.Sprintf("0x%x", v.Const)
	case ValReg:
		return v.Reg
	case ValOp:
		return fmt.Sprintf("(%s %s %s)", v.Left.GoString(), v.Op, v.Right.GoString())
	case ValUnary:
		return fmt.Sprintf("%s(%s)", v.Op, v.Left.GoString())
	case ValLoad:
		return fmt.Sprintf("*(%s)", v.Mem.String())
	case ValAddrOf:
		return fmt.Sprintf("&%s", v.Left.GoString())
	case ValSym:
		if v.Left != nil {
			return fmt.Sprintf("%s+%s", v.Func, v.Left.GoString())
		}
		return v.Func
	case ValCall:
		args := make([]string, len(v.Args))
		for i, a := range v.Args {
			args[i] = a.GoString()
		}
		if v.Func != "" {
			return fmt.Sprintf("%s(%s)", v.Func, join(args, ", "))
		}
		return fmt.Sprintf("call(%s)", join(args, ", "))
	}
	return fmt.Sprintf("?(%d)", v.Kind)
}

type MemRef struct {
	Base   string
	Index  string
	Scale  int
	Offset int64
	Symbol string
}

func (m *MemRef) String() string {
	if m == nil {
		return "?"
	}
	if m.Symbol != "" {
		if m.Offset != 0 {
			return fmt.Sprintf("%s+%d(SB)", m.Symbol, m.Offset)
		}
		return fmt.Sprintf("%s(SB)", m.Symbol)
	}
	s := ""
	if m.Offset != 0 {
		s += fmt.Sprintf("%d", m.Offset)
	}
	s += "(" + m.Base
	if m.Index != "" && m.Scale != 0 {
		s += fmt.Sprintf(")(%s*%d", m.Index, m.Scale)
	}
	s += ")"
	return s
}

type BlockState struct {
	regs       map[string]*Value
	stackSlots map[int64]string
	params     []Param
	varCount   int
	lastCmp    *CmpInfo
	stmts      []Statement
}

type Param struct {
	Name string
	Type string
	Reg  string
}

type CmpInfo struct {
	Left  *Value
	Right *Value
}

type Statement struct {
	Address  uint64
	Kind     StatementKind
	Dst      string
	Value    *Value
	Cond     string
	Func     string
	Args     []*Value
	Label    string
}

type StatementKind int

const (
	StmtAssign   StatementKind = iota
	StmtStore
	StmtCall
	StmtReturn
	StmtIf
	StmtElse
	StmtFor
	StmtGoto
	StmtLabel
	StmtExpr
	StmtComment
)

func NewBlockState() *BlockState {
	return &BlockState{
		regs:       make(map[string]*Value),
		stackSlots: make(map[int64]string),
	}
}

func (s *BlockState) Clone() *BlockState {
	clone := NewBlockState()
	for k, v := range s.regs {
		clone.regs[k] = v
	}
	for k, v := range s.stackSlots {
		clone.stackSlots[k] = v
	}
	clone.params = s.params
	clone.varCount = s.varCount
	return clone
}

func (s *BlockState) AddParam(p Param) {
	s.params = append(s.params, p)
	if p.Reg != "" {
		s.regs[p.Reg] = RegValue(p.Reg)
	}
}

func (s *BlockState) GetReg(name string) *Value {
	if v, ok := s.regs[name]; ok {
		return v
	}
	return RegValue(name)
}

func (s *BlockState) SetReg(name string, v *Value) {
	s.regs[name] = v
}

func (s *BlockState) GetStackSlot(offset int64) string {
	if name, ok := s.stackSlots[offset]; ok {
		return name
	}
	name := fmt.Sprintf("v%d", s.varCount)
	s.varCount++
	s.stackSlots[offset] = name
	return name
}

func (s *BlockState) StackVarName(offset int64) string {
	if name, ok := s.stackSlots[offset]; ok {
		return name
	}
	return fmt.Sprintf("stack_0x%x", offset)
}

func (s *BlockState) AddStmt(stmt Statement) {
	s.stmts = append(s.stmts, stmt)
}

func (s *BlockState) Statements() []Statement {
	return s.stmts
}

func join(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
