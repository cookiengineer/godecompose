package decompile_math_big

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestMathBig(t *testing.T) {
	b := decompile.CompileAndOpen(t, "math_big")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "math_big")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "big.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
