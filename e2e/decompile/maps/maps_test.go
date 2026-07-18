package decompile_maps

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestMaps(t *testing.T) {
	b := decompile.CompileAndOpen(t, "maps")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "maps")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "map") || decompile.HasMatch([]matcher.Match{m}, "fmt.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
