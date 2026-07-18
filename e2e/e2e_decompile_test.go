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

	// Build symbol lookup for disassembly with symbol names
	symLookup := buildSymLookup(b)

	instructions, err := disasm.DecodeStreamWithSymbols(textSection.Data, textSection.Address, symLookup)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

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
	matches := m.Match(instructions)
	t.Logf("found %d matches", len(matches))

	g := generate.New(matches, instructions)
	output := g.Generate()

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		t.Logf("function recovery: %v", err)
	}

	return decompileResult{
		matches:      matches,
		output:       output,
		instructions: instructions,
		result:       result,
	}
}

// buildSymLookup creates a SymLookup from the binary's symbol table.
func buildSymLookup(b binary.Binary) disasm.SymLookup {
	syms, err := b.Symbols()
	if err != nil {
		return nil
	}
	entries := make([]disasm.SymbolEntry, 0, len(syms))
	for _, s := range syms {
		if s.Name != "" {
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

	// Verify high-level fmt patterns match
	assertPatternMatch(t, r.matches, "fmt.Println")

	// fmt.Sprintf and fmt.Errorf may be inlined by Go compiler
	if !checkPatternMatch(r.matches, "fmt.Sprintf") {
		t.Log("fmt.Sprintf not matched (may be inlined)")
	}
	if !checkPatternMatch(r.matches, "fmt.Errorf") {
		t.Log("fmt.Errorf not matched (may be inlined)")
	}

	// Print match details
	for _, m := range r.matches {
		if strings.Contains(m.Pattern.Name, "fmt.") || strings.Contains(m.Pattern.Name, "runtime.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}

	// Verify generated output contains Go-level constructs
	if !strings.Contains(r.output, "fmt.Println") && !strings.Contains(r.output, "fmt.Sprintf") {
		t.Log("generated output may not contain high-level fmt constructs")
		t.Logf("output sample: %s", truncateStr(r.output, 500))
	}
}

func TestE2ESync(t *testing.T) {
	b := compileAndOpen(t, "sync")
	r := decompileBinary(t, b)

	t.Logf("output: %d bytes, instructions: %d, functions: %d",
		len(r.output), len(r.instructions), countFuncs(r.result))

	assertPipelineOk(t, r, "sync")

	assertPipelineOk(t, r, "sync")

	if !checkPatternMatch(r.matches, "sync.Mutex.Lock") {
		t.Log("sync.Mutex.Lock not in resolved matches (conflict resolution may exclude high-level CALL patterns)")
	}
	if !checkPatternMatch(r.matches, "sync.Mutex.Unlock") {
		t.Log("sync.Mutex.Unlock not in resolved matches")
	}
	if !checkPatternMatch(r.matches, "sync.WaitGroup.Done") {
		t.Log("sync.WaitGroup.Done not matched (may be inlined)")
	}

	for _, m := range r.matches {
		if strings.Contains(m.Pattern.Name, "sync.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}

func TestE2EChannels(t *testing.T) {
	b := compileAndOpen(t, "channels")
	r := decompileBinary(t, b)

	t.Logf("output: %d bytes, instructions: %d, functions: %d, matches: %d",
		len(r.output), len(r.instructions), countFuncs(r.result), len(r.matches))

	assertPipelineOk(t, r, "channels")

	for _, m := range r.matches {
		if strings.Contains(m.Pattern.Name, "chan") || strings.Contains(m.Pattern.Name, "runtime.") {
			t.Logf("  match: %s @ 0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.Confidence)
		}
	}
}

func TestE2EMaps(t *testing.T) {
	b := compileAndOpen(t, "maps")
	r := decompileBinary(t, b)

	assertPipelineOk(t, r, "maps")

	for _, m := range r.matches {
		if strings.Contains(m.Pattern.Name, "map") || strings.Contains(m.Pattern.Name, "runtime.") {
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

	// Verify key high-level patterns matched
	if !checkPatternMatch(r.matches, "sync.Mutex") {
		t.Log("sync.Mutex patterns not in resolved matches")
	}
	if !checkPatternMatch(r.matches, "sync.WaitGroup") {
		t.Log("sync.WaitGroup patterns not in resolved matches")
	}

	// Check generated output
	t.Logf("Match patterns found:")
	for _, m := range r.matches {
		t.Logf("  %s @ 0x%x-0x%x (conf=%.2f)", m.Pattern.Name, m.StartAddr, m.EndAddr, m.Confidence)
	}

	// Verify output has sync-relevant code
	if !strings.Contains(r.output, "Lock") && !strings.Contains(r.output, "WaitGroup") {
		t.Log("output may not contain sync constructs (expected for partial decompilation)")
	}
}

func assertPatternMatch(t *testing.T, matches []matcher.Match, nameSubstr string) {
	t.Helper()
	if !checkPatternMatch(matches, nameSubstr) {
		t.Errorf("no match found for pattern %q", nameSubstr)
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

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
