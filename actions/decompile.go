package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/generate"
	"github.com/cookiengineer/godecompose/pattern/matcher"
	"github.com/cookiengineer/godecompose/types"
)

func DecompileBinary(b binary.Binary, db *database.Database) (*DecompileOutput, error) {
	textSection, ok := b.Section(".text")
	if !ok {
		return nil, fmt.Errorf("no .text section found")
	}

	fmt.Fprintf(os.Stderr, "[1/5] disassembling...\n")
	symLookup := buildSymLookup(b)
	instructions, err := disasm.DecodeStreamWithSymbols(textSection.Data, textSection.Address, symLookup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "disassembly error: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "[2/5] recovering functions...\n")
	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "function recovery: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "[3/5] filtering user instructions (%d functions)...\n", len(result.UserFunctions))

	var userInstructions []disasm.Instruction

	addrToIdx := make(map[uint64]int, len(instructions))
	for i, inst := range instructions {
		addrToIdx[inst.Address] = i
	}

	visited := make(map[uint64]bool)
	for _, f := range result.UserFunctions {
		for addr := f.EntryPoint; addr < f.EndAddr; {
			if idx, ok := addrToIdx[addr]; ok {
				inst := instructions[idx]
				if !visited[inst.Address] {
					userInstructions = append(userInstructions, inst)
					visited[inst.Address] = true
				}
				addr += uint64(inst.Size)
			} else {
				addr++
			}
		}
	}

	patterns := db.FindPatterns(b.Architecture(), platformFromBinary(b))

	fmt.Fprintf(os.Stderr, "[4/5] pattern matching (%d instructions, %d patterns)...\n", len(userInstructions), len(db.AllPatterns()))

	m := matcher.New(patterns)
	matches := m.Match(userInstructions)

	goModule := "decompiled"
	if info, hasInfo := b.GoBuildInfo(); hasInfo && info != nil {
		if info.Path != "" && info.Path != goModule {
			goModule = info.Path
		}
	}
	if goModule == "decompiled" {
		for _, f := range result.UserFunctions {
			if f.PackagePath != "" && f.PackagePath != "main" {
				parts := strings.Split(f.PackagePath, "/")
				if len(parts) > 0 && parts[0] != "" {
					goModule = parts[0]
					break
				}
			}
		}
	}

	g := generate.New(matches, userInstructions)
	generatedSource := g.Generate()

	metrics := computeMetrics(result, userInstructions, matches, generatedSource)

	return &DecompileOutput{
		Matches:          matches,
		GeneratedSource:  generatedSource,
		Instructions:     instructions,
		UserInstructions: userInstructions,
		FuncResult:       result,
		GoModule:         goModule,
		Metrics:          metrics,
	}, nil
}

func WriteProject(output *DecompileOutput, dir string) error {
	g := generate.NewForProject(output.Matches, output.UserInstructions, output.FuncResult.UserFunctions, output.FuncResult.Packages, output.FuncResult.Structs)
	return g.WriteProject(dir, output.GoModule)
}

func buildSymLookup(b binary.Binary) disasm.SymLookup {
	syms, err := b.Symbols()
	if err != nil {
		return nil
	}
	entries := make([]disasm.SymbolEntry, 0, len(syms))
	for _, s := range syms {
		if s.Name != "" && s.Size > 0 {
			entries = append(entries, disasm.SymbolEntry{
				Name:    s.Name,
				Address: s.Address,
				Size:    s.Size,
			})
		}
	}
	return disasm.BuildSymLookup(entries)
}

func platformFromBinary(b binary.Binary) types.Platform {
	switch b.Format() {
	case binary.FormatELF:
		return types.PlatformLinux
	case binary.FormatPE:
		return types.PlatformWindows
	case binary.FormatMachO:
		return types.PlatformDarwin
	}
	return types.PlatformUnknown
}

func computeMetrics(result *function.RecoverResult, userInstructions []disasm.Instruction, matches []matcher.Match, generatedSource string) Metrics {
	m := Metrics{
		TotalInstructions: len(userInstructions),
		TotalUserFuncs:    len(result.UserFunctions),
		TotalStructs:      len(result.Structs),
	}

	covered := make(map[uint64]bool)
	for _, match := range matches {
		for addr := match.StartAddr; addr < match.EndAddr; addr++ {
			covered[addr] = true
		}
	}
	for _, inst := range userInstructions {
		if covered[inst.Address] {
			m.MatchedInstructions++
		}
	}
	if m.TotalInstructions > 0 {
		m.RecoveryPct = float64(m.MatchedInstructions) / float64(m.TotalInstructions) * 100
	}

	noiseOps := map[string]bool{"NOP": true, "NOPL": true, "NOPW": true, "INT": true, "INT3": true, "DATA16": true}
	for _, inst := range userInstructions {
		if noiseOps[inst.Opcode] {
			continue
		}
		m.DFAExpressions++
	}
	if m.TotalInstructions > 0 {
		m.DFAExpressionPct = float64(m.DFAExpressions) / float64(m.TotalInstructions) * 100
	}

	for _, f := range result.UserFunctions {
		if f.Blocks != nil && len(f.Blocks) > 0 {
			sig := function.ReconstructSignature(f)
			if len(sig.Args) > 0 || len(sig.Returns) > 0 {
				m.FuncsWithSignatures++
			}
		}
	}

	for _, st := range result.Structs {
		fields := function.InferStructFields(st)
		if len(fields) > 0 {
			m.StructsWithFields++
		}
	}

	callCount := 0
	identifiedCalls := 0
	callCovered := make(map[uint64]bool)
	for _, match := range matches {
		for _, inst := range userInstructions {
			if inst.IsCall && match.StartAddr <= inst.Address && inst.Address < match.EndAddr {
				callCovered[inst.Address] = true
			}
		}
	}
	for _, inst := range userInstructions {
		if inst.IsCall {
			callCount++
		}
	}
	identifiedCalls = len(callCovered)
	m.TotalCallSites = callCount
	m.IdentifiedCallSites = identifiedCalls
	if callCount > 0 {
		m.CallSiteRecoveryPct = float64(identifiedCalls) / float64(callCount) * 100
	}

	return m
}
