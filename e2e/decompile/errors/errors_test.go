package decompile_errors

import (
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

func TestErrors(t *testing.T) {
	b := decompile.CompileAndOpen(t, "errors")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "errors")

	expectedPatterns := []string{"errors.New", "errors.Is", "errors.As", "errors.Unwrap", "errors.Join"}
	for _, name := range expectedPatterns {
		if !decompile.HasMatch(r.Matches, name) {
			t.Logf("  pattern %s not matched", name)
		}
	}

	for _, m := range r.Matches {
		if decompile.HasMatch([]matcher.Match{m}, "errors.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}
