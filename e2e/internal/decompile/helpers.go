package decompile

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/actions"
	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/database"
	"github.com/cookiengineer/godecompose/database/syscall"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/matcher"
	"github.com/cookiengineer/godecompose/patterns/golang"

	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
)

// CompileAndOpen compiles a testdata/src/<name> Go program for linux/amd64 and opens the binary.
func CompileAndOpen(t *testing.T, name string) binary.Binary {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srcDir := filepath.Join(baseDir, "..", "..", "..", "testdata", "src", name)

	dir := t.TempDir()
	outPath := filepath.Join(dir, name)

	cmd := exec.Command("go", "build", "-gcflags=all=-l", "-o", outPath, ".")
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

	db := loadTestDb(t)
	output, err := actions.DecompileBinary(b, db)
	if err != nil {
		t.Fatalf("DecompileBinary: %v", err)
	}

	t.Logf("functions: %d total (user: %d)",
		len(output.FuncResult.Functions), len(output.FuncResult.UserFunctions))
	t.Logf("user instructions: %d / %d total", len(output.UserInstructions), len(output.Instructions))
	t.Logf("found %d matches", len(output.Matches))

	t.Logf("Go module: %s", output.GoModule)
	for pkg, funcs := range output.FuncResult.Packages {
		t.Logf("  package %s: %d functions", pkg, len(funcs))
	}

	return Result{
		Matches:      output.Matches,
		Output:       output.GeneratedSource,
		Instructions: output.UserInstructions,
		FuncResult:   output.FuncResult,
	}
}

func loadTestDb(t *testing.T) *database.Database {
	t.Helper()

	db := database.New()
	if err := golang.LoadStdlib(db); err != nil {
		t.Logf("loading stdlib patterns: %v", err)
	}
	if err := golang.LoadRuntime(db); err != nil {
		t.Logf("loading runtime patterns: %v", err)
	}
	if err := golang.LoadFallback(db); err != nil {
		t.Logf("loading fallback patterns: %v", err)
	}
	if err := golang.LoadControlFlow(db); err != nil {
		t.Logf("loading controlflow patterns: %v", err)
	}
	if err := db.LoadSyscallsFromFS(syscall.TablesFS); err != nil {
		t.Logf("loading syscall tables: %v", err)
	}
	return db
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

// WriteToDir runs the full decompile pipeline and writes project output to a temp dir.
func WriteToDir(t *testing.T, b binary.Binary, db *database.Database) (string, *function.RecoverResult) {
	t.Helper()

	output, err := actions.DecompileBinary(b, db)
	if err != nil {
		t.Fatalf("DecompileBinary: %v", err)
	}

	dir := t.TempDir()
	if err := actions.WriteProject(output, dir); err != nil {
		t.Fatalf("WriteProject: %v", err)
	}

	t.Logf("written to: %s", dir)
	t.Logf("functions: %d total (user: %d)", len(output.FuncResult.Functions), len(output.FuncResult.UserFunctions))
	t.Logf("packages: %d", len(output.FuncResult.Packages))
	t.Logf("structs: %d", len(output.FuncResult.Structs))
	t.Logf("match count: %d", len(output.Matches))

	return dir, output.FuncResult
}
