package phase1_structs

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestStructFields(t *testing.T) {
	b := decompile.CompileAndOpen(t, "phase1_structs")
	r := decompile.Decompile(t, b)

	decompile.AssertPipelineOk(t, r, "phase1_structs")
	decompile.LogMatches(t, r.Matches)

	l := strings.ToLower(r.Output)
	noiseTerms := []string{"nop", "int3", "data16"}
	for _, noise := range noiseTerms {
		if strings.Contains(l, noise) {
			t.Errorf("output contains noise: %q", noise)
		}
	}

	t.Logf("output (%d bytes):\n%s", len(r.Output), r.Output)

	if !strings.Contains(r.Output, "Point") {
		t.Log("no Point struct found")
	}
}
