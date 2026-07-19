
// Package actions provides reusable decompilation pipeline steps.
// Each exported function accepts descriptive parameters so that the
// CLI layer (cmd/godecompose) only handles argument parsing.
package actions

import (
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

// DecompileOutput holds the complete result of a decompilation run.
type DecompileOutput struct {
	Matches          []matcher.Match
	GeneratedSource  string
	Instructions     []disasm.Instruction
	UserInstructions []disasm.Instruction
	FuncResult       *function.RecoverResult
	GoModule         string
	Metrics          Metrics
}

// Metrics holds decompilation recovery rate statistics.
type Metrics struct {
	TotalInstructions    int
	MatchedInstructions  int
	RecoveryPct          float64
	TotalUserFuncs       int
	FuncsWithSignatures  int
	TotalStructs         int
	StructsWithFields    int
	TotalCallSites       int
	IdentifiedCallSites  int
	CallSiteRecoveryPct  float64
}
