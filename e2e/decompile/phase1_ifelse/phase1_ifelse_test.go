package phase1_ifelse

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestIfElse(t *testing.T) {
	b := decompile.CompileAndOpen(t, "phase1_ifelse")
	r := decompile.Decompile(t, b)

	decompile.AssertPipelineOk(t, r, "phase1_ifelse")
	decompile.LogMatches(t, r.Matches)

	if !strings.Contains(r.Output, "return") {
		t.Error("output does not contain return statement")
	}

	l := strings.ToLower(r.Output)
	noiseTerms := []string{"nop", "int3", "data16"}
	for _, noise := range noiseTerms {
		if strings.Contains(l, noise) {
			t.Errorf("output contains noise: %q", noise)
		}
	}

	t.Logf("output (%d bytes):\n%s", len(r.Output), r.Output)
}
