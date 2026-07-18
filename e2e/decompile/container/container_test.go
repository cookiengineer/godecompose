package decompile_container

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestContainer(t *testing.T) {
	b := decompile.CompileAndOpen(t, "container")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "container")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "list.") || decompile.HasMatch([]matcher.Match{m}, "heap.") || decompile.HasMatch([]matcher.Match{m}, "ring.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
