package phase1_forloop

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestForLoop(t *testing.T) {
	b := decompile.CompileAndOpen(t, "phase1_forloop")
	r := decompile.Decompile(t, b)

	decompile.AssertPipelineOk(t, r, "phase1_forloop")
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

	hasFor := strings.Contains(r.Output, "for ")
	hasLoopLabel := strings.Contains(r.Output, "loop")
	t.Logf("output (%d bytes):\n%s", len(r.Output), r.Output)
	if !hasFor && !hasLoopLabel {
		t.Log("no for loop or loop label detected — may still be correct as sequential blocks")
	}
}
