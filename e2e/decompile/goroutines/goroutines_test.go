package decompile_goroutines

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestGoroutines(t *testing.T) {
	b := decompile.CompileAndOpen(t, "goroutines")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "goroutines")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "goroutine.") || decompile.HasMatch([]matcher.Match{m}, "defer") || decompile.HasMatch([]matcher.Match{m}, "fmt.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
