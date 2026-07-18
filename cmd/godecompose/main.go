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

	fmt.Printf("Disassembling .text section: %d bytes at 0x%x\n", textSection.Size, textSection.Address)

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "disassembly error: %v\n", err)
	}
	if len(instructions) > 0 {
		fmt.Printf("Decoded %d instructions\n", len(instructions))
		blocks := disasm.BuildControlFlowGraph(instructions, nil)
		fmt.Printf("Built %d basic blocks\n", len(blocks))
	}

	result, err := function.RecoverFromBinary(b, instructions)
	if err == nil && len(result.Functions) > 0 {
		fmt.Printf("Recovered %d functions\n", len(result.Functions))
		for _, f := range result.Functions {
			if len(f.Blocks) > 0 && f.Name != "" {
				fmt.Printf("  %s @ 0x%x (blocks: %d)\n", f.Name, f.EntryPoint, len(f.Blocks))
			}
		}
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
}

func cmdDecompile(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: godecompose decompile <binary>")
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

	db := loadDatabase()
	fmt.Fprintf(os.Stderr, "Loading database...\n%s\n", db.Stats())

	patterns := db.FindPatterns(b.Architecture(), platformFromBinary(b))
	fmt.Fprintf(os.Stderr, "Matching %d patterns...\n", len(patterns))

	m := matcher.New(patterns)
	matches := m.Match(instructions)

	fmt.Fprintf(os.Stderr, "Found %d matches\n", len(matches))

	g := generate.New(matches, instructions)
	output := g.Generate()

	fmt.Print(output)
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
