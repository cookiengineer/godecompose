package decompile_structfields

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestStructFields(t *testing.T) {
	b := decompile.CompileAndOpen(t, "structfields")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "structfields")

	if len(r.Matches) == 0 {
		t.Fatal("structfields: no patterns matched at all")
	}

	matchNames := make(map[string]int)
	for _, m := range r.Matches {
		matchNames[m.Pattern.Name]++
	}
	for name, count := range matchNames {
		t.Logf("  matched: %s (%dx)", name, count)
	}

	totalMatches := len(r.Matches)
	t.Logf("structfields: %d total matches across %d unique patterns", totalMatches, len(matchNames))

	fieldMatches := 0
	for _, m := range r.Matches {
		if strings.Contains(m.Pattern.Name, "field") {
			fieldMatches++
			t.Logf("  field access: %s @ 0x%x", m.Pattern.Name, m.StartAddr)
		}
	}
	t.Logf("struct field access matches: %d", fieldMatches)

	output := r.Output
	if strings.Contains(output, "field_0x") {
		t.Logf("  output contains struct field references")
	}

	if totalMatches < 5 {
		t.Errorf("structfields: expected at least 5 matches, got %d", totalMatches)
	}
}
