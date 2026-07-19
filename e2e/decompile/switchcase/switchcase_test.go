package decompile_switchcase

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestSwitchCase(t *testing.T) {
	b := decompile.CompileAndOpen(t, "switchcase")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "switchcase")

	if len(r.Matches) == 0 {
		t.Fatal("switchcase: no patterns matched at all")
	}

	matchNames := make(map[string]int)
	for _, m := range r.Matches {
		matchNames[m.Pattern.Name]++
	}
	for name, count := range matchNames {
		t.Logf("  matched: %s (%dx)", name, count)
	}

	totalMatches := len(r.Matches)
	t.Logf("switchcase: %d total matches across %d unique patterns", totalMatches, len(matchNames))

	switchMatches := 0
	for _, m := range r.Matches {
		if strings.Contains(m.Pattern.Name, "switch") || strings.Contains(m.Pattern.Name, "case") {
			switchMatches++
			t.Logf("  switch match: %s @ 0x%x", m.Pattern.Name, m.StartAddr)
		}
	}
	t.Logf("switch/case pattern matches: %d", switchMatches)

	if totalMatches < 5 {
		t.Errorf("switchcase: expected at least 5 matches, got %d", totalMatches)
	}
}
