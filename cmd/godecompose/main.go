// Command godecompose is a pattern-based decompiler for x86_64 binaries.
// It supports ELF, PE/COFF, and Mach-O formats.
package main

import (
	"fmt"
	"os"

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

	case "disassemble":
		if len(args) < 1 {
			fmt.Println("usage: godecompose disassemble <binary>")
			os.Exit(1)
		}
		b := openBinary(args[0])
		defer b.Close()
		actions.Disassemble(b)

	case "decompile":
		if len(args) < 2 {
			printUsage()
			os.Exit(1)
		}
		namespace := args[0]
		binaryPath := args[1]
		outputDir := namespace

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

		if err := actions.WriteProject(output, outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "write project: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\nProject written to: %s\n", outputDir)

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
  godecompose info <binary>                   Show binary metadata
  godecompose disassemble <binary>            Disassemble binary
  godecompose decompile <namespace> <binary>  Decompile binary into namespace folder
  godecompose patterns list                   List available patterns
  godecompose patterns validate <file>        Validate a pattern file

Examples:
  godecompose info ./myapp
  godecompose disassemble ./myapp
  godecompose decompile myproject ./myapp
  godecompose patterns list`)
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
	if err := golang.LoadFallback(db); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading fallback patterns: %v\n", err)
	}
	if err := golang.LoadControlFlow(db); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading controlflow patterns: %v\n", err)
	}
	if err := db.LoadSyscallsFromFS(syscall.TablesFS); err != nil {
		fmt.Fprintf(os.Stderr, "warning: loading syscall tables: %v\n", err)
	}

	return db
}
