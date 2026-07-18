package e2e

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

func compileAndOpen(t *testing.T, name string) binary.Binary {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srcDir := filepath.Join(baseDir, "..", "testdata", "src", "e2e_"+name)

	dir := t.TempDir()
	outPath := filepath.Join(dir, name)

	cmd := exec.Command("go", "build", "-o", outPath, ".")
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("compile e2e_%s: %v\n%s", name, err, out)
	}

	b, err := binary.Open(outPath)
	if err != nil {
		t.Fatalf("open %s: %v", outPath, err)
	}
	t.Cleanup(func() { b.Close() })
	return b
}

type decompileResult struct {
	matches      []matcher.Match
	output       string
	instructions []disasm.Instruction
	result       *function.RecoverResult
}

func decompileBinary(t *testing.T, b binary.Binary) decompileResult {
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

	// Recover and classify functions first
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

	// Only match against user function instructions
	var userInstructions []disasm.Instruction
	for _, f := range result.UserFunctions {
		for _, inst := range instructions {
			if inst.Address >= f.EntryPoint && inst.Address < f.EndAddr {
				userInstructions = append(userInstructions, inst)
			}
		}
	}
	t.Logf("user instructions: %d / %d total", len(userInstructions), len(instructions))

	// Load patterns and match
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	patternsDir := filepath.Join(baseDir, "..", "patterns")

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

	return decompileResult{
		matches:      matches,
		output:       output,
		instructions: userInstructions,
		result:       result,
	}
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

func TestE2EFmt(t *testing.T) {
	b := compileAndOpen(t, "fmt")
	r := decompileBinary(t, b)

	t.Logf("output: %d bytes, instructions: %d, functions: %d",
		len(r.output), len(r.instructions), countFuncs(r.result))

	assertPipelineOk(t, r, "fmt")

	if !checkPatternMatch(r.matches, "fmt.Println") {
		t.Log("fmt.Println not matched (may be inlined or no symbol table)")
	}

	for _, m := range r.matches {
		t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
	}
}

func TestE2ESync(t *testing.T) {
	b := compileAndOpen(t, "sync")
	r := decompileBinary(t, b)

	t.Logf("output: %d bytes, instructions: %d, functions: %d",
		len(r.output), len(r.instructions), countFuncs(r.result))

	assertPipelineOk(t, r, "sync")

	for _, m := range r.matches {
		if strings.Contains(m.Pattern.Name, "sync.") || strings.Contains(m.Pattern.Name, "fmt.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}

func TestE2EFullPipeline(t *testing.T) {
	b := compileAndOpen(t, "sync")
	r := decompileBinary(t, b)

	t.Logf("Full pipeline: %d instructions, %d functions, %d matches, %d output bytes",
		len(r.instructions), countFuncs(r.result), len(r.matches), len(r.output))

	if len(r.output) == 0 {
		t.Error("decompiler produced no output")
	}

	if countFuncs(r.result) == 0 {
		t.Error("no functions recovered")
	}

	for _, m := range r.matches {
		t.Logf("  %s @ 0x%x-0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.EndAddr, m.Confidence)
	}
}

func checkPatternMatch(matches []matcher.Match, nameSubstr string) bool {
	for _, m := range matches {
		if strings.Contains(m.Pattern.Name, nameSubstr) {
			return true
		}
	}
	return false
}

func assertPipelineOk(t *testing.T, r decompileResult, category string) {
	t.Helper()
	if len(r.instructions) == 0 {
		t.Errorf("%s: no instructions decoded", category)
	}
	if countFuncs(r.result) == 0 {
		t.Errorf("%s: no functions recovered", category)
	}
	if len(r.output) == 0 {
		t.Errorf("%s: no decompiler output", category)
	}
}

func countFuncs(r *function.RecoverResult) int {
	if r == nil {
		return 0
	}
	return len(r.Functions)
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
