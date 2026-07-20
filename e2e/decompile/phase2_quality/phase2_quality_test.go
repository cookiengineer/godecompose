package phase2_quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/patterns/golang"

	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
)

func TestQualityOutput(t *testing.T) {
	b := decompile.CompileAndOpen(t, "phase2_quality")

	db := database.New()
	golang.LoadStdlib(db)
	golang.LoadRuntime(db)
	golang.LoadFallback(db)
	golang.LoadControlFlow(db)

	dir, result := decompile.WriteToDir(t, b, db)

	mainGo := filepath.Join(dir, "main.go")
	content, err := os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	text := string(content)
	t.Logf("main.go (%d bytes):\n%s", len(content), text)

	checks := []string{
		"package main",
		"func main()",
		"type Counter struct",
		"func NewCounter",
		"func (c *Counter) Increment",
		"func (c *Counter) GetValue",
		"func processData",
		"func classify",
		"func switchExample",
	}
	for _, c := range checks {
		if !strings.Contains(text, c) {
			t.Errorf("missing: %q", c)
		}
	}

	qualityChecks := []string{
		"return 0",
		"*(AX) = (*AX + 1)",
		"return *AX",
	}
	for _, c := range qualityChecks {
		if !strings.Contains(text, c) {
			t.Logf("note: missing quality marker %q — may need DFA improvement", c)
		}
	}

	noiseTerms := []string{"data16", "int3"}
	for _, n := range noiseTerms {
		if strings.Contains(strings.ToLower(text), n) {
			t.Errorf("noise in output: %q", n)
		}
	}

	if len(result.Structs) >= 1 {
		st := result.Structs[0]
		t.Logf("struct: %s.%s (%d methods)", st.PackagePath, st.Name, len(st.Methods))
		if st.Name != "Counter" {
			t.Errorf("expected Counter struct, got %q", st.Name)
		}
		if len(st.Methods) < 2 {
			t.Errorf("expected >=2 methods on Counter, got %d", len(st.Methods))
		}
	}

	if len(result.UserFunctions) < 7 {
		t.Errorf("expected >=7 user functions, got %d", len(result.UserFunctions))
	}
}
