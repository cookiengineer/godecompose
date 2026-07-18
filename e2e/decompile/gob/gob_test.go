package decompile_gob

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestGob(t *testing.T) {
	b := decompile.CompileAndOpen(t, "gob")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "gob")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "gob.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
