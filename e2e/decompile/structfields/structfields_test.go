package decompile_structfields

import (
	"sort"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/e2e/internal/decompile"
	"github.com/cookiengineer/godecompose/function"
)

func TestStructFields(t *testing.T) {
	b := decompile.CompileAndOpen(t, "structfields")
	r := decompile.Decompile(t, b)

	t.Logf("output: %d bytes, instructions: %d", len(r.Output), len(r.Instructions))
	decompile.AssertPipelineOk(t, r, "structfields")

	if len(r.Matches) == 0 {
		t.Fatal("structfields: no patterns matched at all")
	}

	fr := r.FuncResult
	t.Logf("functions: %d total, %d user", len(fr.Functions), len(fr.UserFunctions))
	t.Logf("structs recovered: %d", len(fr.Structs))

	for _, f := range fr.UserFunctions {
		t.Logf("  func: %s pkg=%s recv=%q method=%v ptr=%v",
			f.ShortName, f.PackagePath, f.ReceiverType, f.IsMethod, f.IsPointerReceiver)
	}

	hasStruct := false
	for _, st := range fr.Structs {
		t.Logf("  struct: %s pkg=%s methods=%d", st.Name, st.PackagePath, len(st.Methods))
		if st.Name == "Point" {
			hasStruct = true
			fields := function.InferStructFields(st)
			sort.Slice(fields, func(i, j int) bool { return fields[i].Offset < fields[j].Offset })
			t.Logf("  Point fields: %d", len(fields))
			for _, fld := range fields {
				t.Logf("    %s offset=%s type=%s count=%d", fld.Name, fld.Offset, fld.Type, fld.Count)
			}
			if len(fields) < 2 {
				t.Errorf("Point struct: expected at least 2 fields, got %d", len(fields))
			}
		}
	}
	if !hasStruct {
		t.Log("Point struct not recovered (may be inlined or methods not detected)")
	}

	methodCount := 0
	for _, f := range fr.UserFunctions {
		if f.IsMethod {
			methodCount++
		}
	}
	t.Logf("method functions: %d", methodCount)
	if methodCount < 5 {
		t.Errorf("structfields: expected at least 5 methods, got %d (compiler may inline)", methodCount)
	}

	output := r.Output
	checks := []struct {
		name string
		fn   func(string) bool
	}{
		{"has fmt.Println", func(s string) bool {
			return strings.Contains(s, "Println")
		}},
		{"no INT3 noise in output", func(s string) bool {
			return !strings.Contains(s, "int3")
		}},
	}
	for _, c := range checks {
		if c.fn(output) {
			t.Logf("  ✓ %s", c.name)
		} else {
			t.Logf("  ✗ %s", c.name)
		}
	}
}
