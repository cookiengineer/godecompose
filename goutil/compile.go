// Package goutil provides test helpers for compiling Go source to
// cross-compiled binaries (ELF, PE, Mach-O) for integration testing.
package goutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Target specifies a cross-compilation target: OS, architecture, output format, and file extension.
type Target struct {
	GOOS   string
	GOARCH string
	Format string
	Ext    string
}

// DefaultTargets returns the three standard amd64 targets.
func DefaultTargets() []Target {
	return []Target{
		{GOOS: "linux", GOARCH: "amd64", Format: "ELF", Ext: ""},
		{GOOS: "windows", GOARCH: "amd64", Format: "PE", Ext: ".exe"},
		{GOOS: "darwin", GOARCH: "amd64", Format: "MachO", Ext: ""},
	}
}

// GetBinary returns the compiled binary path for a given GOOS from a
// CompileProgram result map.
func GetBinary(t *testing.T, binaries map[Target]string, goos string) string {
	t.Helper()
	for _, target := range DefaultTargets() {
		if target.GOOS == goos {
			if path, ok := binaries[target]; ok {
				return path
			}
		}
	}
	return ""
}

// CompileProgram cross-compiles the Go program at srcDir for every target and
// returns a map from target to compiled binary path. Binaries are placed in a
// temp directory cleaned up when the test finishes.
func CompileProgram(t *testing.T, srcDir string, targets []Target) map[Target]string {
	t.Helper()

	if len(targets) == 0 {
		targets = DefaultTargets()
	}

	dir, err := os.MkdirTemp("", "godecompose-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	binaries := make(map[Target]string)

	for _, target := range targets {
		outPath := filepath.Join(dir, fmt.Sprintf("%s_%s_%s%s",
			filepath.Base(srcDir), target.GOOS, target.GOARCH, target.Ext))

		cmd := exec.Command("go", "build",
			"-o", outPath,
			".",
		)
		cmd.Dir = srcDir
		cmd.Env = append(os.Environ(),
			"GOOS="+target.GOOS,
			"GOARCH="+target.GOARCH,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("compile GOOS=%s GOARCH=%s: %v\n%s", target.GOOS, target.GOARCH, err, output)
		}

		if _, err := os.Stat(outPath); err != nil {
			t.Fatalf("compiled binary not found at %s: %v", outPath, err)
		}

		binaries[target] = outPath
	}

	return binaries
}

// CompileSimple cross-compiles the simple test program for all default targets.
func CompileSimple(t *testing.T) map[Target]string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srcDir := filepath.Join(baseDir, "..", "testdata", "src", "simple")
	return CompileProgram(t, srcDir, DefaultTargets())
}

// CompileComplex cross-compiles the complex test program for all default targets.
func CompileComplex(t *testing.T) map[Target]string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srcDir := filepath.Join(baseDir, "..", "testdata", "src", "complex")
	return CompileProgram(t, srcDir, DefaultTargets())
}
