package decompile_controlflow

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestControlFlow(t *testing.T) {
	b := decompile.CompileAndOpen(t, "controlflow")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "controlflow")

	if len(r.Matches) == 0 {
		t.Fatal("controlflow: no patterns matched at all")
	}

	matchNames := make(map[string]int)
	for _, m := range r.Matches {
		matchNames[m.Pattern.Name]++
	}
	for name, count := range matchNames {
		t.Logf("  matched: %s (%dx)", name, count)
	}

	totalMatches := len(r.Matches)
	t.Logf("controlflow: %d total matches across %d unique patterns", totalMatches, len(matchNames))

	if totalMatches < 10 {
		t.Errorf("controlflow: expected at least 10 matches, got %d", totalMatches)
	}

	output := r.Output
	expectedGoKeywords := []string{
		"if",
		"return",
		"make(chan",
		"m[key]",
		"os.CreateTemp",
		"os.Stat",
		"defer",
		"go fn",
		"new(T)",
	}
	for _, kw := range expectedGoKeywords {
		if strings.Contains(output, kw) {
			t.Logf("  output contains: %q", kw)
		} else {
			t.Logf("  output missing: %q", kw)
		}
	}
}
