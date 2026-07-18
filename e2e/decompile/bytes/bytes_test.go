package decompile_bytes

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestBytes(t *testing.T) {
	b := decompile.CompileAndOpen(t, "bytes")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "bytes")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "bytes.") || decompile.HasMatch([]matcher.Match{m}, "bufio.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
