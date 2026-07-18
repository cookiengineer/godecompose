// Package function provides function boundary recovery from disassembled
// binaries. For Go binaries, it parses the pclntab (PC-line table) for
// exact function information. For generic C/C++ binaries, it uses
// prologue/epilogue heuristics.
package function

import "github.com/cookiengineer/godecompose/disasm"

// Variable represents a function parameter, return value, or local variable.
type Variable struct {
	Name     string
	Offset   int64
	TypeName string
	IsArg    bool
	IsReturn bool
}

// Function represents a recovered function with its basic blocks and metadata.
type Function struct {
	Name       string
	EntryPoint uint64
	EndAddr    uint64
	FrameSize  int
	ArgSize    int
	Blocks     []*disasm.BasicBlock
	Args       []Variable
	Returns    []Variable
	Locals     []Variable
	IsGoFunc   bool
	SourceFile string
	SourceLine int
}

// RecoverResult holds all recovered functions and any warnings.
type RecoverResult struct {
	Functions []*Function
	Warnings  []string
}
