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
	Name             string
	EntryPoint       uint64
	EndAddr          uint64
	FrameSize        int
	ArgSize          int
	Blocks           []*disasm.BasicBlock
	Args             []Variable
	Returns          []Variable
	Locals           []Variable
	IsGoFunc         bool
	SourceFile       string
	SourceLine       int
	Classification   Classification
	PackagePath      string
	ShortName        string
	ReceiverType     string
	IsMethod         bool
	IsPointerReceiver bool
}

// StructType represents a recovered struct type with its methods.
type StructType struct {
	Name        string
	PackagePath string
	Methods     []*Function
}

// ParsePackageName extracts the package path, function name, and optional
// method receiver from a Go symbol name. Symbols can be:
//
//	main.main              -> pkg=main, func=main
//	main.greet             -> pkg=main, func=greet
//	main.cmdInit.func3     -> pkg=main, func=cmdInit.func3 (closure)
//	main.Type.Method       -> pkg=main, receiver=Type, method=Method
//	main.(*Type).Method    -> pkg=main, receiver=Type, pointer, method=Method
//	main.(Type).Method     -> pkg=main, receiver=Type, method=Method
//	example.com/pkg.Func   -> pkg=example.com/pkg, func=Func
//	example.com/pkg.T.M    -> pkg=example.com/pkg, receiver=T, method=M
//	example.com/pkg.(*T).M -> pkg=example.com/pkg, receiver=T, pointer, method=M
func ParsePackageName(fullName string) (pkgPath string, shortName string, receiverType string, isPointerReceiver bool, isMethod bool) {
	if fullName == "" {
		return "", "", "", false, false
	}

	rest := fullName

	if strings.HasPrefix(rest, "main.") {
		pkgPath = "main"
		rest = rest[5:]
	} else {
		receiverDot := strings.Index(rest, ".(")
		if receiverDot >= 0 {
			pkgPath = rest[:receiverDot]
			rest = rest[receiverDot+1:]
		} else {
			lastDot := strings.LastIndex(rest, ".")
			if lastDot < 0 {
				return "", rest, "", false, false
			}
			pkgPath = rest[:lastDot]
			rest = rest[lastDot+1:]
		}
	}

	if strings.HasSuffix(rest, ".abi0") {
		rest = rest[:len(rest)-5]
	}

	if strings.HasPrefix(rest, "(") {
		closing := strings.IndexByte(rest, ')')
		if closing > 1 && closing+1 < len(rest) && rest[closing+1] == '.' {
			receiverStr := rest[1:closing]
			isPointerReceiver = strings.HasPrefix(receiverStr, "*")
			receiverType = strings.TrimPrefix(receiverStr, "*")
			shortName = rest[closing+2:]
			isMethod = true
			return
		}
	}

	if pkgPath != "main" {
		if funcDot := strings.IndexByte(rest, '.'); funcDot > 0 {
			receiverType = rest[:funcDot]
			shortName = rest[funcDot+1:]
			isMethod = true
			return
		}
	}

	if pkgPath != "main" {
		if funcDot := strings.IndexByte(rest, '.'); funcDot > 0 {
			receiverType = rest[:funcDot]
			shortName = rest[funcDot+1:]
			isMethod = true
			return
		}
	} else {
		dots := 0
		for _, ch := range rest {
			if ch == '.' {
				dots++
			}
		}
		if dots == 1 {
			funcDot := strings.IndexByte(rest, '.')
			candidate := rest[funcDot+1:]
			if !isClosureName(candidate) {
				receiverType = rest[:funcDot]
				shortName = candidate
				isMethod = true
				return
			}
		}
	}

	shortName = rest
	return
}

// isClosureName returns true if the name looks like a Go closure
// (e.g., "func1", "func2", "func3", "func1.1", "func2.1").
func isClosureName(name string) bool {
	return strings.HasPrefix(name, "func") && len(name) > 4 &&
		(name[4] >= '0' && name[4] <= '9')
}

// SetPackageInfo parses the function name and sets PackagePath, ShortName,
// and method receiver fields.
func (f *Function) SetPackageInfo() {
	f.PackagePath, f.ShortName, f.ReceiverType, f.IsPointerReceiver, f.IsMethod = ParsePackageName(f.Name)
}

// BuildStructs groups user functions by their receiver type into struct
// definitions. The initial PackagePath is taken from the first method.
func (r *RecoverResult) BuildStructs() {
	groups := make(map[string]*StructType)

	for _, f := range r.UserFunctions {
		if !f.IsMethod || f.ReceiverType == "" {
			continue
		}
		st, ok := groups[f.ReceiverType]
		if !ok {
			st = &StructType{
				Name:        f.ReceiverType,
				PackagePath: f.PackagePath,
			}
			groups[f.ReceiverType] = st
		}
		st.Methods = append(st.Methods, f)
	}

	r.Structs = make([]*StructType, 0, len(groups))
	for _, st := range groups {
		r.Structs = append(r.Structs, st)
	}
}

// RefineStructPackages updates each struct's PackagePath to the consensus
// of its methods' packages. If methods disagree (e.g. after callgraph
// refinement), the majority package wins.
func (r *RecoverResult) RefineStructPackages() {
	for _, st := range r.Structs {
		pkgCounts := make(map[string]int)
		for _, m := range st.Methods {
			pkgCounts[m.PackagePath]++
		}
		best := ""
		bestCount := 0
		for pkg, count := range pkgCounts {
			if count > bestCount {
				best = pkg
				bestCount = count
			}
		}
		if best != "" {
			st.PackagePath = best
		}
	}
}

// RecoverResult holds all recovered functions and any warnings.
type RecoverResult struct {
	Functions     []*Function
	UserFunctions []*Function
	Warnings      []string
	GoMainPackage string
	Packages      map[string][]*Function
	Structs       []*StructType
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
