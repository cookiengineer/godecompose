// Command godecompose is a pattern-based decompiler for x86_64 binaries.
// It supports ELF, PE/COFF, and Mach-O formats.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cookiengineer/godecompose/actions"
	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/database/syscall"
	"github.com/cookiengineer/godecompose/patterns/golang"

	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
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
		if len(args) < 1 {
			fmt.Println("usage: godecompose info <binary>")
			os.Exit(1)
		}
		b := openBinary(args[0])
		defer b.Close()
		actions.Info(b)

	case "disasm":
		if len(args) < 1 {
			fmt.Println("usage: godecompose disasm <binary>")
			os.Exit(1)
		}
		b := openBinary(args[0])
		defer b.Close()
		actions.Disassemble(b)

	case "decompile":
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

		db := loadDatabase()
		output, err := actions.DecompileBinary(b, db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decompile error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "=== Decompilation Summary ===\n")
		fmt.Fprintf(os.Stderr, "Go module:  %s\n", output.GoModule)
		fmt.Fprintf(os.Stderr, "Patterns:   %d loaded, %d matches\n",
			len(db.AllPatterns()), len(output.Matches))
		fmt.Fprintf(os.Stderr, "Packages:   %d\n", len(output.FuncResult.Packages))
		for pkg := range output.FuncResult.Packages {
			fmt.Fprintf(os.Stderr, "  %s (%d functions)\n", pkg, len(output.FuncResult.Packages[pkg]))
		}

		if outputDir != "" {
			if err := actions.WriteProject(output, outputDir); err != nil {
				fmt.Fprintf(os.Stderr, "write project: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "\nProject written to: %s\n", outputDir)
		} else {
			fmt.Print(output.GeneratedSource)
		}

	case "patterns":
		if len(args) < 1 {
			fmt.Println("usage: godecompose patterns <list|validate> [args]")
			os.Exit(1)
		}
		switch args[0] {
		case "list":
			db := loadDatabase()
			actions.PatternsList(db)
		case "validate":
			if len(args) < 2 {
				fmt.Println("usage: godecompose patterns validate <file.hexpat>")
				os.Exit(1)
			}
			if err := actions.PatternsValidate(args[1]); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
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
	db := database.New()

	if err := golang.LoadStdlib(db); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading stdlib patterns: %v\n", err)
	}
	if err := golang.LoadRuntime(db); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading runtime patterns: %v\n", err)
	}
	if err := db.LoadSyscallsFromFS(syscall.TablesFS); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading syscall tables: %v\n", err)
	}

	return db
}
