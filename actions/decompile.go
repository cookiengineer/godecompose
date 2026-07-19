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

	symLookup := buildSymLookup(b)
	instructions, err := disasm.DecodeStreamWithSymbols(textSection.Data, textSection.Address, symLookup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "disassembly error: %v\n", err)
	}

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "function recovery: %v\n", err)
	}

	var userInstructions []disasm.Instruction
	for _, f := range result.UserFunctions {
		for _, inst := range instructions {
			if inst.Address >= f.EntryPoint && inst.Address < f.EndAddr {
				userInstructions = append(userInstructions, inst)
			}
		}
	}

	patterns := db.FindPatterns(b.Architecture(), platformFromBinary(b))

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

	return &DecompileOutput{
		Matches:          matches,
		GeneratedSource:  generatedSource,
		Instructions:     instructions,
		UserInstructions: userInstructions,
		FuncResult:       result,
		GoModule:         goModule,
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
