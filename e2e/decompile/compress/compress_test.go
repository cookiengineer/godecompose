package decompile_compress

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestCompress(t *testing.T) {
	b := decompile.CompileAndOpen(t, "compress")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "compress")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "gzip.") || decompile.HasMatch([]matcher.Match{m}, "zlib.") || decompile.HasMatch([]matcher.Match{m}, "flate.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
