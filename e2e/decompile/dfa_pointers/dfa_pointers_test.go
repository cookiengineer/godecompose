package dfa_pointers

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestDFAPointers(t *testing.T) {
	b := decompile.CompileAndOpen(t, "dfa_pointers")
	r := decompile.Decompile(t, b)

	decompile.AssertPipelineOk(t, r, "dfa_pointers")
	decompile.LogMatches(t, r.Matches)

	l := strings.ToLower(r.Output)
	noiseTerms := []string{"nop", "int3", "data16"}
	for _, noise := range noiseTerms {
		if strings.Contains(l, noise) {
			t.Errorf("output contains noise: %q", noise)
		}
	}

	t.Logf("output (%d bytes):\n%s", len(r.Output), r.Output)

	hasUpdate := false
	hasPrint := false
	for _, f := range r.FuncResult.UserFunctions {
		if strings.Contains(f.Name, "updateItem") {
			hasUpdate = true
		}
		if strings.Contains(f.Name, "printItem") {
			hasPrint = true
		}
	}
	if !hasUpdate {
		t.Error("updateItem function not found")
	}
	if !hasPrint {
		t.Error("printItem function not found")
	}

	if !decompile.HasMatch(r.Matches, "fmt") && !decompile.HasMatch(r.Matches, "Println") {
		t.Log("no fmt patterns matched — printItem call might not be resolved")
	}
}
