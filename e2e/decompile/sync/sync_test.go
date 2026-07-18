package decompile_sync

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestSync(t *testing.T) {
	b := decompile.CompileAndOpen(t, "sync")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "sync")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "sync.") || decompile.HasMatch([]matcher.Match{m}, "fmt.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
