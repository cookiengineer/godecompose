package decompile_encoding

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestEncoding(t *testing.T) {
	b := decompile.CompileAndOpen(t, "encoding")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "encoding")

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "base64.") || decompile.HasMatch([]matcher.Match{m}, "hex.") ||
			decompile.HasMatch([]matcher.Match{m}, "xml.") || decompile.HasMatch([]matcher.Match{m}, "binary.") ||
			decompile.HasMatch([]matcher.Match{m}, "regexp.") || decompile.HasMatch([]matcher.Match{m}, "filepath.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
