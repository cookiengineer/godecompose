package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func PatternsDiscover(sourcePath string) error {
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	dir := filepath.Dir(absPath)
	file := filepath.Base(absPath)

	tmpDir, err := os.MkdirTemp("", "godecompose-discover")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("go", "build", "-gcflags=-S", "-o", filepath.Join(tmpDir, "out"), file)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOARCH=amd64", "GOOS=linux")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compiling: %w\n%s", err, out)
	}

	asm := string(out)
	callRe := regexp.MustCompile(`\s+CALL\s+(\S+)`)
	seen := make(map[string]bool)

	var discovered []string
	for _, line := range strings.Split(asm, "\n") {
		matches := callRe.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}
		target := matches[1]
		target = strings.TrimSuffix(target, "(SB)")
		target = strings.TrimSuffix(target, "(SB")

		if strings.Contains(target, "0x") || strings.Contains(target, "(") {
			continue
		}

		if target == "runtime.morestack" || target == "runtime.morestack_noctxt" ||
			strings.Contains(target, "gcWriteBarrier") || target == "runtime.duffzero" ||
			target == "runtime.duffcopy" {
			continue
		}

		patternName := strings.ReplaceAll(target, ".", "_")
		patternName = strings.ReplaceAll(patternName, "/", "_")
		patternName = strings.ReplaceAll(patternName, "(", "_")
		patternName = strings.ReplaceAll(patternName, ")", "_")
		patternName = strings.ReplaceAll(patternName, "*", "_")

		if seen[patternName] {
			continue
		}
		seen[patternName] = true

		genName := strings.TrimPrefix(target, "runtime.")
		genName = strings.TrimPrefix(genName, "main.")

		pattern := fmt.Sprintf(`pattern go_%s {
    name: "go: %s";
    library: "go-discovered";
    description: "Auto-discovered from %s";
    instr match { CALL %s }
    gen { %s(...) }
}
`, patternName, target, file, patternName, genName)

		discovered = append(discovered, pattern)
	}

	if len(discovered) == 0 {
		fmt.Println("No CALL instructions found in compiler output.")
		return nil
	}

	fmt.Printf("// Auto-generated patterns from: %s\n", file)
	fmt.Printf("// Found %d unique CALL targets\n\n", len(discovered))
	fmt.Printf("arch x86_64;\n")
	fmt.Printf("platform linux, darwin, windows, freebsd;\n\n")
	for _, p := range discovered {
		fmt.Println(p)
	}
	fmt.Fprintf(os.Stderr, "Discovered %d patterns\n", len(discovered))

	return nil
}
