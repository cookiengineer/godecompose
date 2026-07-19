package decompile_structout

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestStructuredOutput(t *testing.T) {
	b := decompile.CompileAndOpen(t, "structout")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "structout")

	if len(r.Matches) == 0 {
		t.Fatal("structout: no patterns matched at all")
	}
	if len(r.FuncResult.UserFunctions) == 0 {
		t.Fatal("structout: no user functions recovered")
	}

	hasBranch := false
	hasLoop := false
	for _, f := range r.FuncResult.UserFunctions {
		t.Logf("  func: %s (short:%s) recv=%q method=%v pkg=%s",
			f.Name, f.ShortName, f.ReceiverType, f.IsMethod, f.PackagePath)
		if strings.Contains(f.ShortName, "branch") {
			hasBranch = true
		}
		if strings.Contains(f.ShortName, "loop") || strings.Contains(f.ShortName, "Count") {
			hasLoop = true
		}
	}
	if !hasBranch {
		t.Error("structout: branch function not found (may be inlined)")
	}
	if !hasLoop {
		t.Log("structout: loop function not found (may be inlined)")
	}

	output := r.Output
	checks := []struct {
		name string
		fn   func(string) bool
	}{
		{"has fmt.Println match", func(s string) bool {
			return strings.Contains(s, "fmt.Println") || strings.Contains(s, "Println")
		}},
		{"has function signatures", func(s string) bool {
			for _, f := range r.FuncResult.UserFunctions {
				if strings.Contains(s, f.ShortName) {
					return true
				}
			}
			return false
		}},
	}
	for _, check := range checks {
		if check.fn(output) {
			t.Logf("  ✓ %s", check.name)
		} else {
			t.Logf("  ✗ %s", check.name)
		}
	}

	if len(r.FuncResult.Packages) < 1 {
		t.Error("structout: no packages found")
	}
	t.Logf("packages: %d", len(r.FuncResult.Packages))
	for pkg := range r.FuncResult.Packages {
		t.Logf("  pkg: %s (%d funcs)", pkg, len(r.FuncResult.Packages[pkg]))
	}
}
