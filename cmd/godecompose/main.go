// Command godecompose is a pattern-based decompiler for x86_64 binaries.
// It supports ELF, PE/COFF, and Mach-O formats.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/generate"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/parser"
	"github.com/cookiengineer/godecompose/pattern/lang/validator"
	"github.com/cookiengineer/godecompose/pattern/matcher"
	"github.com/cookiengineer/godecompose/types"

	_ "github.com/cookiengineer/godecompose/elf"
	_ "github.com/cookiengineer/godecompose/macho"
	_ "github.com/cookiengineer/godecompose/pe"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "info":
		cmdInfo(args)
	case "disasm":
		cmdDisasm(args)
	case "decompile":
		cmdDecompile(args)
	case "patterns":
		if len(args) < 1 {
			fmt.Println("usage: godecompose patterns <list|validate> [args]")
			os.Exit(1)
		}
		switch args[0] {
		case "list":
			cmdPatternsList(args[1:])
		case "validate":
			cmdPatternsValidate(args[1:])
		default:
			fmt.Printf("unknown patterns subcommand: %s\n", args[0])
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`godecompose — pattern-based decompiler for x86_64 binaries

Usage:
  godecompose info <binary>               Show binary metadata
  godecompose disasm <binary>             Disassemble binary
  godecompose decompile <binary>          Full decompilation pipeline
  godecompose patterns list               List available patterns
  godecompose patterns validate <file>    Validate a pattern file`)
}

func openBinary(path string) binary.Binary {
	b, err := binary.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v\n", path, err)
		os.Exit(1)
	}
	return b
}

func loadDatabase() *database.Database {
	exe, _ := os.Executable()
	baseDir := filepath.Dir(exe)
	patternsDir := filepath.Join(baseDir, "..", "..", "patterns")

	if _, err := os.Stat(filepath.Join(patternsDir, "kernels")); err != nil {
		patternsDir = filepath.Join(baseDir, "patterns")
	}
	if _, err := os.Stat(filepath.Join(patternsDir, "kernels")); err != nil {
		patternsDir = "patterns"
	}

	db := database.New()
	if err := db.LoadSyscallsFromDir(filepath.Join(patternsDir, "kernels")); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading syscall tables: %v\n", err)
	}
	if err := db.LoadPatternsFromDir(filepath.Join(patternsDir, "libs")); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading patterns: %v\n", err)
	}

	return db
}

func cmdInfo(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: godecompose info <binary>")
		os.Exit(1)
	}

	b := openBinary(args[0])
	defer b.Close()

	fmt.Printf("File:        %s\n", args[0])
	fmt.Printf("Format:      %s\n", b.Format())
	fmt.Printf("Arch:        %s\n", b.Architecture())
	fmt.Printf("Entry:       0x%x\n", b.EntryPoint())
	fmt.Printf("PIE:         %v\n", b.IsPIE())
	fmt.Printf("Stripped:    %v\n", b.IsStripped())

	sections := b.Sections()
	fmt.Printf("Sections:    %d\n", len(sections))
	for _, s := range sections {
		fmt.Printf("  %-20s addr=0x%x size=0x%x flags=%c%c%c\n",
			s.Name, s.Address, s.Size,
			flagChar(s.Flags, binary.SectionExecutable, 'X'),
			flagChar(s.Flags, binary.SectionWritable, 'W'),
			flagChar(s.Flags, binary.SectionReadable, 'R'),
		)
	}

	syms, err := b.Symbols()
	if err == nil {
		fmt.Printf("Symbols:     %d\n", len(syms))
	}

	if info, ok := b.GoBuildInfo(); ok {
		fmt.Printf("Go version:  %s\n", info.Version)
		fmt.Printf("Go path:     %s\n", info.Path)
	}
}

func cmdDisasm(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: godecompose disasm <binary>")
		os.Exit(1)
	}

	b := openBinary(args[0])
	defer b.Close()

	textSection, ok := b.Section(".text")
	if !ok {
		fmt.Fprintf(os.Stderr, "no .text section found\n")
		os.Exit(1)
	}

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "disassembly error: %v\n", err)
	}
	fmt.Printf("Decoded %d instructions\n", len(instructions))

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "function recovery: %v\n", err)
		return
	}

	runtimeCount := 0
	stdlibCount := 0
	otherCount := 0
	for _, f := range result.Functions {
		switch f.Classification {
		case function.ClassRuntime:
			runtimeCount++
		case function.ClassStdlib:
			stdlibCount++
		default:
			otherCount++
		}
	}

	fmt.Printf("Functions: %d total\n", len(result.Functions))
	fmt.Printf("  runtime:  %d (skipped)\n", runtimeCount)
	fmt.Printf("  stdlib:   %d (skipped)\n", stdlibCount)
	fmt.Printf("  user:     %d\n", len(result.UserFunctions))
	fmt.Printf("  other:    %d\n", otherCount)

	if len(result.UserFunctions) > 0 {
		fmt.Printf("\nUser functions:\n")
		for _, f := range result.UserFunctions {
			fmt.Printf("  %-50s @ 0x%x (blocks: %d)\n", f.Name, f.EntryPoint, len(f.Blocks))
		}
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
}

func cmdDecompile(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: godecompose decompile <binary> [--output=<dir>]")
		os.Exit(1)
	}

	binaryPath := args[0]
	outputDir := ""
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--output=") {
			outputDir = arg[9:]
		}
	}

	b := openBinary(binaryPath)
	defer b.Close()

	textSection, ok := b.Section(".text")
	if !ok {
		fmt.Fprintf(os.Stderr, "no .text section found\n")
		os.Exit(1)
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

	runtimeCount := 0
	stdlibCount := 0
	for _, f := range result.Functions {
		switch f.Classification {
		case function.ClassRuntime:
			runtimeCount++
		case function.ClassStdlib:
			stdlibCount++
		}
	}

	db := loadDatabase()
	patterns := db.FindPatterns(b.Architecture(), platformFromBinary(b))

	var userInstructions []disasm.Instruction
	for _, f := range result.UserFunctions {
		for _, inst := range instructions {
			if inst.Address >= f.EntryPoint && inst.Address < f.EndAddr {
				userInstructions = append(userInstructions, inst)
			}
		}
	}

	m := matcher.New(patterns)
	matches := m.Match(userInstructions)

	// Determine Go module name — try build info first, then infer from symbols
	goModule := "decompiled"
	if info, hasInfo := b.GoBuildInfo(); hasInfo && info != nil {
		if info.Path != "" && info.Path != goModule {
			goModule = info.Path
		}
	}
	// Fallback: extract module from non-main user function names
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

	fmt.Fprintf(os.Stderr, "=== Decompilation Summary ===\n")
	fmt.Fprintf(os.Stderr, "Go module:  %s\n", goModule)
	fmt.Fprintf(os.Stderr, "Functions:  %d total (runtime: %d, stdlib: %d, user: %d)\n",
		len(result.Functions), runtimeCount, stdlibCount, len(result.UserFunctions))
	fmt.Fprintf(os.Stderr, "Patterns:   %d loaded, %d matches\n", len(patterns), len(matches))
	fmt.Fprintf(os.Stderr, "Packages:   %d\n", len(result.Packages))

	for pkg := range result.Packages {
		fmt.Fprintf(os.Stderr, "  %s (%d functions)\n", pkg, len(result.Packages[pkg]))
	}

	if outputDir != "" {
		g := generate.NewForProject(matches, userInstructions, result.UserFunctions, result.Packages)
		if err := g.WriteProject(outputDir, goModule); err != nil {
			fmt.Fprintf(os.Stderr, "write project: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\nProject written to: %s\n", outputDir)
	} else {
		g := generate.New(matches, userInstructions)
		fmt.Print(g.Generate())
	}
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

func cmdPatternsList(args []string) {
	db := loadDatabase()
	fmt.Print(db.Stats())

	if len(db.AllPatterns()) == 0 {
		fmt.Println("\nNo patterns loaded.")
		return
	}

	fmt.Println("\nLoaded patterns:")
	for _, p := range db.AllPatterns() {
		fmt.Printf("  %s", p.Name)
		if p.Library != "" {
			fmt.Printf(" [%s]", p.Library)
		}
		if p.Version != "" {
			fmt.Printf(" %s", p.Version)
		}
		if p.Arch != "" {
			fmt.Printf(" arch=%s", p.Arch)
		}
		if len(p.Platforms) > 0 {
			fmt.Printf(" platform=%s", strings.Join(p.Platforms, ","))
		}
		fmt.Println()
	}
}

func cmdPatternsValidate(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: godecompose patterns validate <file.hexpat>")
		os.Exit(1)
	}

	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	l := lexer.NewWithFile(string(data), path)
	tokens, err := l.Lex()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lex error: %v\n", err)
		os.Exit(1)
	}

	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	v := validator.New()
	errs := v.Validate(prog)
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "validation errors (%d):\n", len(errs))
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  %v\n", e)
		}
		os.Exit(1)
	}

	e := evaluator.New()
	patterns, err := e.Evaluate(prog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "evaluation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Pattern file %s is valid.\n", path)
	fmt.Printf("  Parsed %d patterns\n", len(patterns))
	if len(patterns) > 0 {
		for _, pat := range patterns {
			altCount := len(pat.Alternatives)
			opCount := 0
			for _, alt := range pat.Alternatives {
				opCount += len(alt)
			}
			fmt.Printf("  - %s (alternatives: %d, operations: %d)\n", pat.Name, altCount, opCount)
			if pat.GenTemplate != "" {
				fmt.Printf("    gen: %s\n", pat.GenTemplate)
			}
		}
	}
}

func flagChar(flags binary.SectionFlag, flag binary.SectionFlag, ch byte) byte {
	if flags&flag != 0 {
		return ch
	}
	return '-'
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
