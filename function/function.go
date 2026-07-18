// Package function provides function boundary recovery from disassembled
// binaries. For Go binaries, it parses the pclntab (PC-line table) for
// exact function information. For generic C/C++ binaries, it uses
// prologue/epilogue heuristics.
package function

import (
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
)

// Variable represents a function parameter, return value, or local variable.
type Variable struct {
	Name     string
	Offset   int64
	TypeName string
	IsArg    bool
	IsReturn bool
}

// Classification identifies whether a function is user code or runtime/stdlib.
type Classification int

const (
	ClassUnknown  Classification = iota
	ClassUser                    // Application code (main package)
	ClassRuntime                 // Go runtime (runtime.*)
	ClassStdlib                  // Go standard library (fmt.*, sync.*, os.*, etc.)
	ClassInternal                // Go internal packages (internal/*)
	ClassVendor                  // Third-party dependencies
)

func (c Classification) String() string {
	switch c {
	case ClassUser:
		return "user"
	case ClassRuntime:
		return "runtime"
	case ClassStdlib:
		return "stdlib"
	case ClassInternal:
		return "internal"
	case ClassVendor:
		return "vendor"
	default:
		return "unknown"
	}
}

// Function represents a recovered function with its basic blocks and metadata.
type Function struct {
	Name           string
	EntryPoint     uint64
	EndAddr        uint64
	FrameSize      int
	ArgSize        int
	Blocks         []*disasm.BasicBlock
	Args           []Variable
	Returns        []Variable
	Locals         []Variable
	IsGoFunc       bool
	SourceFile     string
	SourceLine     int
	Classification Classification
	PackagePath    string
	ShortName      string
}

// ParsePackageName extracts the package path and short function name from
// a Go symbol name like "main.main" or "testproject/utils.Greet".
// Returns (packagePath, shortName).
func ParsePackageName(fullName string) (string, string) {
	// Handle main package specially
	if strings.HasPrefix(fullName, "main.") {
		return "main", fullName[5:]
	}

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot < 0 {
		return "", fullName
	}

	pkgPath := fullName[:lastDot]
	shortName := fullName[lastDot+1:]

	// Clean up: remove .abi0 suffix, .func1 suffix, etc.
	if strings.HasSuffix(shortName, ".abi0") {
		shortName = shortName[:len(shortName)-5]
	}

	return pkgPath, shortName
}

// SetPackageInfo parses the function name and sets PackagePath and ShortName.
func (f *Function) SetPackageInfo() {
	f.PackagePath, f.ShortName = ParsePackageName(f.Name)
}

// RecoverResult holds all recovered functions and any warnings.
type RecoverResult struct {
	Functions      []*Function
	UserFunctions  []*Function
	Warnings       []string
	GoMainPackage  string
	Packages       map[string][]*Function
}

// FilterUserFunctions separates user code from runtime/stdlib functions
// and groups them by package.
func (r *RecoverResult) FilterUserFunctions() {
	r.UserFunctions = nil
	r.Packages = make(map[string][]*Function)
	for _, f := range r.Functions {
		f.SetPackageInfo()
		f.Classification = ClassifyFunction(f.Name, r.GoMainPackage)
		if f.Classification == ClassUser {
			r.UserFunctions = append(r.UserFunctions, f)
			r.Packages[f.PackagePath] = append(r.Packages[f.PackagePath], f)
		}
	}
}

// ClassifyFunction determines the category of a Go function by its name.
// mainPackage is the module path from GoBuildInfo (e.g., "testproject").
func ClassifyFunction(name, mainPackage string) Classification {
	if name == "" {
		return ClassUnknown
	}

	// Main package functions — always user code
	if strings.HasPrefix(name, "main.") {
		return ClassUser
	}

	// Functions from the user's main module package
	// e.g., "testproject/utils.Greet" when mainPackage is "testproject"
	if mainPackage != "" {
		if strings.HasPrefix(name, mainPackage+".") || strings.HasPrefix(name, mainPackage+"/") {
			return ClassUser
		}
	}

	// Go runtime
	if strings.HasPrefix(name, "runtime.") ||
		strings.HasPrefix(name, "runtime/") ||
		strings.Contains(name, ".abi") ||
		strings.HasPrefix(name, "_rt0_") ||
		name == "runtime" {
		return ClassRuntime
	}

	// Internal packages
	if strings.HasPrefix(name, "internal/") ||
		strings.HasPrefix(name, "internal.") {
		return ClassInternal
	}

	// Standard library packages
	stdlibPrefixes := []string{
		"fmt.", "sync.", "os.", "io.", "time.", "net.", "http.",
		"strings.", "bytes.", "strconv.", "errors.", "log.", "flag.",
		"math.", "sort.", "path.", "context.", "reflect.", "regexp.",
		"encoding/", "crypto/", "bufio.", "archive/", "compress/",
		"database/", "debug/", "go/", "hash/", "image/", "index/",
		"mime/", "text/", "unicode/", "syscall.", "container/",
		"sync/", "os/", "io/", "time/", "net/",
	}
	for _, prefix := range stdlibPrefixes {
		if strings.HasPrefix(name, prefix) {
			return ClassStdlib
		}
	}

	// Generated/internal symbols
	if strings.HasPrefix(name, "type:") ||
		strings.HasPrefix(name, "go:") ||
		strings.HasPrefix(name, "gc:") ||
		strings.HasPrefix(name, "gogo") ||
		strings.HasPrefix(name, "callRet") ||
		strings.HasPrefix(name, "gosave") ||
		strings.HasPrefix(name, "setg_") ||
		strings.HasPrefix(name, "debugCall") ||
		strings.HasPrefix(name, "aeshash") ||
		strings.HasPrefix(name, "gcWriteBarrier") {
		return ClassRuntime
	}

	// If it has a dot, it's likely a package.function (stdlib or vendor)
	if strings.Contains(name, ".") {
		// Could be user code from a module — if mainPackage is set and
		// the name starts with a non-stdlib prefix, treat as vendor
		return ClassVendor
	}

	return ClassUnknown
}
