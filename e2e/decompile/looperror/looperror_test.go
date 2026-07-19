package decompile_looperror

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestLoopError(t *testing.T) {
	b := decompile.CompileAndOpen(t, "looperror")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "looperror")

	if len(r.Matches) == 0 {
		t.Fatal("looperror: no patterns matched at all")
	}

	matchNames := make(map[string]int)
	for _, m := range r.Matches {
		matchNames[m.Pattern.Name]++
	}
	for name, count := range matchNames {
		t.Logf("  matched: %s (%dx)", name, count)
	}

	totalMatches := len(r.Matches)
	t.Logf("looperror: %d total matches across %d unique patterns", totalMatches, len(matchNames))

	loopMatches := 0
	errorMatches := 0
	assertMatches := 0
	for _, m := range r.Matches {
		n := m.Pattern.Name
		if strings.Contains(n, "for loop") || strings.Contains(n, "range") {
			loopMatches++
			t.Logf("  loop: %s @ 0x%x", n, m.StartAddr)
		}
		if strings.Contains(n, "error") || strings.Contains(n, "return err") || strings.Contains(n, "return nil") {
			errorMatches++
			t.Logf("  error: %s @ 0x%x", n, m.StartAddr)
		}
		if strings.Contains(n, "type assert") || strings.Contains(n, "assertI") || strings.Contains(n, "assertE") {
			assertMatches++
			t.Logf("  assert: %s @ 0x%x", n, m.StartAddr)
		}
	}

	t.Logf("loop: %d, error: %d, assert: %d", loopMatches, errorMatches, assertMatches)

	output := r.Output
	expected := []string{"for i := 0", "if err != nil", "v := x.(T)", "os.CreateTemp", "os.Stat", "errors.Is"}
	for _, kw := range expected {
		if strings.Contains(output, kw) {
			t.Logf("  output contains: %q", kw)
		} else {
			t.Logf("  output missing: %q (may be inlined)", kw)
		}
	}

	if totalMatches < 5 {
		t.Errorf("looperror: expected at least 5 matches, got %d", totalMatches)
	}
}
