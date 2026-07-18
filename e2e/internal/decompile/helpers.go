package decompile

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/generate"
	"github.com/cookiengineer/godecompose/pattern/matcher"
	"github.com/cookiengineer/godecompose/types"

	_ "github.com/cookiengineer/godecompose/elf"
	_ "github.com/cookiengineer/godecompose/macho"
	_ "github.com/cookiengineer/godecompose/pe"
)

// CompileAndOpen compiles a testdata/src/<name> Go program for linux/amd64 and opens the binary.
func CompileAndOpen(t *testing.T, name string) binary.Binary {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srcDir := filepath.Join(baseDir, "..", "..", "..", "testdata", "src", name)

	dir := t.TempDir()
	outPath := filepath.Join(dir, name)

	cmd := exec.Command("go", "build", "-o", outPath, ".")
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("compile %s: %v\n%s", name, err, out)
	}

	b, err := binary.Open(outPath)
	if err != nil {
		t.Fatalf("open %s: %v", outPath, err)
	}
	t.Cleanup(func() { b.Close() })
	return b
}

// Result holds the output of a decompilation run.
type Result struct {
	Matches      []matcher.Match
	Output       string
	Instructions []disasm.Instruction
	FuncResult   *function.RecoverResult
}

// Decompile runs the full decompile pipeline on a binary.
func Decompile(t *testing.T, b binary.Binary) Result {
	t.Helper()

	textSection, ok := b.Section(".text")
	if !ok {
		t.Fatal("no .text section")
	}

	symLookup := buildSymLookup(b)
	instructions, err := disasm.DecodeStreamWithSymbols(textSection.Data, textSection.Address, symLookup)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		t.Logf("function recovery: %v", err)
	}

	runtimeCount := 0
	stdlibCount := 0
	for _, f := range result.Functions {
		switch f.Classification {
		case function.ClassRuntime:
			runtimeCount++
		case function.ClassStdlib:
			stdlibCount++
		}
	}
	t.Logf("functions: %d total (runtime: %d, stdlib: %d, user: %d)",
		len(result.Functions), runtimeCount, stdlibCount, len(result.UserFunctions))

	var userInstructions []disasm.Instruction
	for _, f := range result.UserFunctions {
		for _, inst := range instructions {
			if inst.Address >= f.EntryPoint && inst.Address < f.EndAddr {
				userInstructions = append(userInstructions, inst)
			}
		}
	}
	t.Logf("user instructions: %d / %d total", len(userInstructions), len(instructions))

	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	patternsDir := filepath.Join(baseDir, "..", "..", "..", "patterns")

	db := database.New()
	_ = db.LoadSyscallsFromDir(filepath.Join(patternsDir, "kernels"))
	if err := db.LoadPatternsFromDir(filepath.Join(patternsDir, "libs")); err != nil {
		t.Logf("loading patterns: %v", err)
	}

	patterns := db.FindPatterns(b.Architecture(), platformGuess(b))
	t.Logf("loaded %d patterns for matching", len(patterns))

	m := matcher.New(patterns)
	matches := m.Match(userInstructions)
	t.Logf("found %d matches", len(matches))

	g := generate.New(matches, userInstructions)
	output := g.Generate()

	return Result{
		Matches:      matches,
		Output:       output,
		Instructions: userInstructions,
		FuncResult:   result,
	}
}

// AssertPipelineOk checks that the decompile pipeline produced non-empty results.
func AssertPipelineOk(t *testing.T, r Result, category string) {
	t.Helper()
	if len(r.Instructions) == 0 {
		t.Errorf("%s: no instructions decoded", category)
	}
	if r.FuncResult == nil || len(r.FuncResult.Functions) == 0 {
		t.Errorf("%s: no functions recovered", category)
	}
	if len(r.Output) == 0 {
		t.Errorf("%s: no decompiler output", category)
	}
}

// HasMatch returns true if any match name contains the given substring.
func HasMatch(matches []matcher.Match, nameSubstr string) bool {
	for _, m := range matches {
		if strings.Contains(m.Pattern.Name, nameSubstr) {
			return true
		}
	}
	return false
}

// LogMatches logs all matches at debug level.
func LogMatches(t *testing.T, matches []matcher.Match) {
	t.Helper()
	for _, m := range matches {
		t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
	}
}

// FilterMatchesByPackage returns matches whose name contains the package prefix.
func FilterMatchesByPackage(matches []matcher.Match, pkgPrefix string) []matcher.Match {
	var out []matcher.Match
	for _, m := range matches {
		if strings.Contains(m.Pattern.Name, pkgPrefix) {
			out = append(out, m)
		}
	}
	return out
}

func buildSymLookup(b binary.Binary) disasm.SymLookup {
	syms, err := b.Symbols()
	if err != nil {
		return nil
	}
	entries := make([]disasm.SymbolEntry, 0, len(syms))
	for _, s := range syms {
		if s.Name != "" && s.Size > 0 {
			entries = append(entries, disasm.SymbolEntry{
				Name:    s.Name,
				Address: s.Address,
				Size:    s.Size,
			})
		}
	}
	return disasm.BuildSymLookup(entries)
}

func platformGuess(b binary.Binary) types.Platform {
	switch b.Format() {
	case binary.FormatELF:
		return types.PlatformLinux
	case binary.FormatPE:
		return types.PlatformWindows
	case binary.FormatMachO:
		return types.PlatformDarwin
	}
	return types.PlatformUnknown
}
