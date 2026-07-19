package decompile_funcgroup

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
)

func TestFuncGroup(t *testing.T) {
	b := decompile.CompileAndOpen(t, "funcgroup")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "funcgroup")

	funcsPerPkg := make(map[string]int)
	for pkg, funcs := range r.FuncResult.Packages {
		funcsPerPkg[pkg] = len(funcs)
		t.Logf("  pkg %s: %d functions", pkg, len(funcs))
	}
	t.Logf("funcgroup: %d total matches across %d packages", len(r.Matches), len(r.FuncResult.Packages))

	mainPkgCount, hasMain := funcsPerPkg["main"]
	if !hasMain {
		t.Error("funcgroup: no main package found")
	} else {
		t.Logf("main package has %d functions", mainPkgCount)
	}

	totalPkgs := len(r.FuncResult.Packages)
	if totalPkgs > 2 {
		t.Errorf("funcgroup: too many packages (%d), expected main + optional subpackage", totalPkgs)
	}

	var methodFuncs []string
	for _, f := range r.FuncResult.UserFunctions {
		t.Logf("  func: pkg=%s name=%s receiver=%s method=%v",
			f.PackagePath, f.ShortName, f.ReceiverType, f.IsMethod)
		if f.IsMethod && f.ReceiverType == "Counter" && f.PackagePath == "main" {
			methodFuncs = append(methodFuncs, f.ShortName)
			t.Logf("    -> main.Counter method: %s", f.ShortName)
		}
	}

	for _, f := range r.FuncResult.UserFunctions {
		if f.PackagePath != "main" && f.PackagePath != "" &&
			!strings.Contains(f.PackagePath, "internal") &&
			!strings.HasPrefix(f.PackagePath, "funcgroup") {
			t.Errorf("funcgroup: unexpected package: %s", f.PackagePath)
		}
	}

	t.Logf("Counter methods found: %v", methodFuncs)

	if len(r.Matches) < 3 {
		t.Errorf("funcgroup: expected at least 3 matches, got %d", len(r.Matches))
	}
}
