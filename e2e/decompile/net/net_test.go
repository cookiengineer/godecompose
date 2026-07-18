package decompile_net

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestNet(t *testing.T) {
	b := decompile.CompileAndOpen(t, "net")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "net")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "net.") || decompile.HasMatch([]matcher.Match{m}, "http.") || decompile.HasMatch([]matcher.Match{m}, "url.") || decompile.HasMatch([]matcher.Match{m}, "tls.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
