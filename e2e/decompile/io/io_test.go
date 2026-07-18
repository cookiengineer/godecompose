package decompile_io

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestIo(t *testing.T) {
	b := decompile.CompileAndOpen(t, "io")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "io")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "io.") || decompile.HasMatch([]matcher.Match{m}, "time.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
