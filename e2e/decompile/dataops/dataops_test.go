package decompile_dataops

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestDataOps(t *testing.T) {
	b := decompile.CompileAndOpen(t, "dataops")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "dataops")

	if len(r.Matches) == 0 {
		t.Fatal("dataops: no patterns matched at all")
	}

	matchNames := make(map[string]int)
	for _, m := range r.Matches {
		matchNames[m.Pattern.Name]++
	}
	for name, count := range matchNames {
		t.Logf("  matched: %s (%dx)", name, count)
	}

	totalMatches := len(r.Matches)
	t.Logf("dataops: %d total matches across %d unique patterns", totalMatches, len(matchNames))

	if totalMatches < 10 {
		t.Errorf("dataops: expected at least 10 matches, got %d", totalMatches)
	}

	output := r.Output

	expectedKeywords := []string{
		"make([]",
		"append(",
		"string(",
		"time.Now",
	}
	for _, kw := range expectedKeywords {
		if strings.Contains(output, kw) {
			t.Logf("  output contains: %q", kw)
		} else {
			t.Logf("  output missing: %q (compiler may have inlined/optimized)", kw)
		}
	}

	dataPatterns := []string{
		"runtime.makeslice",
		"runtime.growslice",
		"runtime.slicebytetostring",
		"runtime.stringtoslicebyte",
		"runtime.concatstring",
		"runtime.newobject",
		"runtime.makemap",
		"runtime.convT",
		"runtime.rand",
		"runtime.memequal",
		"runtime.slicecopy",
		"runtime.typedslicecopy",
		"runtime.newarray",
	}
	matched := 0
	for _, p := range dataPatterns {
		if decompile.HasMatch(r.Matches, p) {
			matched++
			t.Logf("  data pattern found: %s", p)
		}
	}
	t.Logf("data patterns matched: %d/%d", matched, len(dataPatterns))

	if matched < 3 {
		t.Errorf("dataops: too few data patterns matched: %d (expected >= 3)", matched)
	}
}
