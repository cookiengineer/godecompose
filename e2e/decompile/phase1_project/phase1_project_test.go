package phase1_project

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

func TestProjectOutput(t *testing.T) {
	b := decompile.CompileAndOpen(t, "phase1_structs")

	db := database.New()
	golang.LoadStdlib(db)
	golang.LoadRuntime(db)
	golang.LoadFallback(db)
	golang.LoadControlFlow(db)

	dir, result := decompile.WriteToDir(t, b, db)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	t.Logf("project dir contents: %d entries", len(entries))
	for _, e := range entries {
		t.Logf("  %s (isDir=%v)", e.Name(), e.IsDir())
	}

	mainGo := filepath.Join(dir, "main.go")
	content, err := os.ReadFile(mainGo)
	if err != nil {
		t.Errorf("read main.go: %v", err)
	} else {
		text := string(content)
		t.Logf("main.go (%d bytes):\n%s", len(content), text)

		checks := []string{
			"package main",
			"func main()",
			"type Point struct",
		}
		for _, c := range checks {
			if !strings.Contains(text, c) {
				t.Errorf("main.go missing: %q", c)
			}
		}

		noiseTerms := []string{"data16", "int3"}
		for _, n := range noiseTerms {
			if strings.Contains(strings.ToLower(text), n) {
				t.Errorf("main.go contains noise: %q", n)
			}
		}
	}

	if len(result.Structs) > 0 {
		for _, st := range result.Structs {
			t.Logf("struct: %s.%s (%d methods)", st.PackagePath, st.Name, len(st.Methods))
			for _, m := range st.Methods {
				t.Logf("  method: %s", m.ShortName)
			}
		}
	}

	if result.GoMainPackage == "" {
		t.Log("no main package detected")
	}
}
