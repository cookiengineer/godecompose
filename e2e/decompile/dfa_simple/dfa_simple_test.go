package dfa_simple

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestDFASimple(t *testing.T) {
	b := decompile.CompileAndOpen(t, "dfa_simple")
	r := decompile.Decompile(t, b)

	decompile.AssertPipelineOk(t, r, "dfa_simple")
	decompile.LogMatches(t, r.Matches)

	hasAdd42 := false
	hasMain := false
	for _, f := range r.FuncResult.UserFunctions {
		if strings.Contains(f.Name, "add42") {
			hasAdd42 = true
		}
		if strings.Contains(f.Name, "main") {
			hasMain = true
		}
	}
	if !hasAdd42 {
		t.Error("add42 function not found")
	}
	if !hasMain {
		t.Error("main function not found")
	}

	t.Logf("decompiled output (%d bytes):\n%s", len(r.Output), r.Output)

	oldAsmMarkers := []string{
		"00000000", // hex addresses in asm comments
	}
	for _, marker := range oldAsmMarkers {
		if strings.Contains(r.Output, marker) {
			t.Logf("note: output contains address markers (expected for some ranges)")
		}
	}
}
