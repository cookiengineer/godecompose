package decompile_math

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestMath(t *testing.T) {
	b := decompile.CompileAndOpen(t, "math")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "math")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "math.") || decompile.HasMatch([]matcher.Match{m}, "sort.") || decompile.HasMatch([]matcher.Match{m}, "rand.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
