package decompile_slice_string

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestSliceString(t *testing.T) {
	b := decompile.CompileAndOpen(t, "slice_string")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "slice_string")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "slice") || decompile.HasMatch([]matcher.Match{m}, "string") || decompile.HasMatch([]matcher.Match{m}, "fmt.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
