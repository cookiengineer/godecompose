package decompile_archive

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestArchive(t *testing.T) {
	b := decompile.CompileAndOpen(t, "archive")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "archive")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "tar.") || decompile.HasMatch([]matcher.Match{m}, "zip.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
