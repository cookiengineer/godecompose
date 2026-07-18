package decompile_text_template

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestTextTemplate(t *testing.T) {
	b := decompile.CompileAndOpen(t, "text_template")
	r := decompile.Decompile(t, b)
	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "text_template")
	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "template.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
