package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/parser"
	"github.com/cookiengineer/godecompose/pattern/lang/validator"
)

func PatternsList(db *database.Database) error {
	fmt.Print(db.Stats())

	if len(db.AllPatterns()) == 0 {
		fmt.Println("\nNo patterns loaded.")
		return nil
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

	return nil
}

func PatternsValidate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	l := lexer.NewWithFile(string(data), filePath)
	tokens, err := l.Lex()
	if err != nil {
		return fmt.Errorf("lex error: %w", err)
	}

	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	v := validator.New()
	errs := v.Validate(prog)
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "validation errors (%d):\n", len(errs))
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  %v\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errs))
	}

	e := evaluator.New()
	patterns, err := e.Evaluate(prog)
	if err != nil {
		return fmt.Errorf("evaluation error: %w", err)
	}

	fmt.Printf("Pattern file %s is valid.\n", filePath)
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

	return nil
}
