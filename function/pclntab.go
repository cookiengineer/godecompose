package function

import (
	encodingbinary "encoding/binary"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/disasm"
)

// RecoverFromBinary recovers function boundaries from a parsed binary.
func RecoverFromBinary(bin binary.Binary, instructions []disasm.Instruction) (*RecoverResult, error) {
	result := &RecoverResult{}

	pclntab, pclntabAddr, hasPclntab := bin.Pclntab()
	goInfo, hasGoInfo := bin.GoBuildInfo()

	// Set the main package path for classification
	if hasGoInfo && goInfo != nil {
		result.GoMainPackage = goInfo.Main
	}
	if result.GoMainPackage == "" && hasGoInfo && goInfo != nil && goInfo.Path != "" {
		if isValidModulePath(goInfo.Path) {
			result.GoMainPackage = goInfo.Path
		}
	}

	syms, symsErr := bin.Symbols()
	if symsErr != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("symbols: %v", symsErr))
	}

	if result.GoMainPackage == "" {
		// Fallback: extract from function symbols
		result.GoMainPackage = extractModuleFromSymbols(syms)
	}

	if hasPclntab && len(pclntab) > 0 {
		funcs, err := parsePclntab(pclntab, pclntabAddr)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("pclntab parse: %v", err))
		} else {
			isGo := hasGoInfo || len(funcs) > 10
			mergeSymbolNames(funcs, syms)
			for _, f := range funcs {
				f.IsGoFunc = isGo
				f.Blocks = extractFunctionBlocks(instructions, f.EntryPoint, f.EndAddr)
			}
			result.Functions = funcs
			result.FilterUserFunctions()
			result.RefinePackagesByCallGraph(instructions)
			result.BuildStructs()
			result.RefineStructPackages()
			return result, nil
		}
	}

	funcs := recoverBySymbols(instructions, syms)
	for _, f := range funcs {
		f.Blocks = extractFunctionBlocks(instructions, f.EntryPoint, f.EndAddr)
	}
	result.Functions = funcs
	result.FilterUserFunctions()
	result.RefinePackagesByCallGraph(instructions)
	result.BuildStructs()
	result.RefineStructPackages()

	if len(funcs) == 0 {
		funcs = recoverByHeuristics(instructions, bin.EntryPoint())
		for _, f := range funcs {
			f.Blocks = extractFunctionBlocks(instructions, f.EntryPoint, f.EndAddr)
		}
		result.Functions = funcs
		result.FilterUserFunctions()
		result.RefinePackagesByCallGraph(instructions)
		result.BuildStructs()
		result.RefineStructPackages()
	}

	return result, nil
}

// mergeSymbolNames matches function entries from pclntab with symbol names
// by address. Functions without a matching symbol get a "func_<addr>" name.
func mergeSymbolNames(funcs []*Function, syms []binary.Symbol) {
	if len(syms) == 0 {
		for _, f := range funcs {
			if f.Name == "" {
				f.Name = fmt.Sprintf("func_%x", f.EntryPoint)
			}
		}
		return
	}

	symByName := make(map[uint64]string)
	for _, s := range syms {
		if s.Address != 0 && s.Name != "" {
			symByName[s.Address] = s.Name
		}
	}

	for _, f := range funcs {
		if name, ok := symByName[f.EntryPoint]; ok {
			f.Name = name
		} else {
			f.Name = fmt.Sprintf("func_%x", f.EntryPoint)
		}
	}
}

func parsePclntab(data []byte, baseAddr uint64) ([]*Function, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("pclntab too short: %d bytes", len(data))
	}

	magic := encodingbinary.LittleEndian.Uint32(data[0:4])

	var funcs []*Function

	switch magic {
	case 0xFFFFFFF0, 0xFFFFFFF1:
		funcs = parsePclntabV118(data, baseAddr)
	case 0xFFFFFFFB:
		funcs = parsePclntabV116(data, baseAddr)
	case 0xFFFFFFFA:
		funcs = parsePclntabV12(data, baseAddr)
	default:
		return nil, fmt.Errorf("unknown pclntab magic: 0x%08X", magic)
	}

	if len(funcs) == 0 {
		return nil, fmt.Errorf("no functions found in pclntab")
	}

	return funcs, nil
}

func parsePclntabV118(data []byte, baseAddr uint64) []*Function {
	if len(data) < 16 {
		return nil
	}

	quantum := data[4]
	pointerSize := data[5]
	_ = quantum

	if pointerSize != 8 {
		return nil
	}

	nfunc := int(encodingbinary.LittleEndian.Uint16(data[6:8]))
	tabOffset := uint32(8)

	var funcs []*Function
	for i := 0; i < nfunc; i++ {
		fentry := tabOffset + uint32(i)*funcTabSize(pointerSize)
		if fentry+8 > uint32(len(data)) {
			break
		}

		entryOff := readPtr(data[fentry:], pointerSize)
		_ = entryOff

		funcs = append(funcs, &Function{
			EntryPoint: baseAddr + uint64(entryOff),
		})
	}

	return funcs
}

func parsePclntabV116(data []byte, baseAddr uint64) []*Function {
	if len(data) < 16 {
		return nil
	}

	pointerSize := data[5]
	if pointerSize != 8 {
		return nil
	}

	nfunc := int(encodingbinary.LittleEndian.Uint32(data[8:12]))
	funcOffset := encodingbinary.LittleEndian.Uint32(data[12:16])

	var funcs []*Function
	for i := 0; i < nfunc; i++ {
		entry := funcOffset + uint32(i)*uint32(funcTabSizeV116())
		if entry+16 > uint32(len(data)) {
			break
		}

		entryOff := readPtrV116(data[entry:])
		_ = entryOff

		funcs = append(funcs, &Function{
			EntryPoint: baseAddr + uint64(entryOff),
		})
	}

	return funcs
}

func parsePclntabV12(data []byte, baseAddr uint64) []*Function {
	if len(data) < 16 {
		return nil
	}

	pointerSize := data[5]
	if pointerSize != 8 {
		return nil
	}

	nfunc := int(encodingbinary.LittleEndian.Uint32(data[8:12]))
	funcOffset := int(encodingbinary.LittleEndian.Uint32(data[12:16]))

	var funcs []*Function
	for i := 0; i < nfunc; i++ {
		entry := funcOffset + i*int(funcTabSizeV12())
		if entry+8 > len(data) {
			break
		}

		entryOff := encodingbinary.LittleEndian.Uint32(data[entry : entry+4])
		funcs = append(funcs, &Function{
			EntryPoint: baseAddr + uint64(entryOff),
		})
	}

	return funcs
}

func extractFunctionBlocks(instructions []disasm.Instruction, entry, end uint64) []*disasm.BasicBlock {
	funcInsts := make([]disasm.Instruction, 0)
	for _, inst := range instructions {
		if inst.Address >= entry && (end == 0 || inst.Address < end) {
			funcInsts = append(funcInsts, inst)
		}
	}

	if len(funcInsts) == 0 {
		return nil
	}

	return disasm.BuildControlFlowGraph(funcInsts, []uint64{entry})
}

func recoverBySymbols(instructions []disasm.Instruction, syms []binary.Symbol) []*Function {
	if len(syms) == 0 {
		return nil
	}

	addrMap := make(map[uint64][]binary.Symbol)
	for _, s := range syms {
		if s.Type == binary.SymbolFunction && s.Name != "" {
			addrMap[s.Address] = append(addrMap[s.Address], s)
		}
	}

	sortedAddrs := make([]uint64, 0, len(addrMap))
	for addr := range addrMap {
		sortedAddrs = append(sortedAddrs, addr)
	}
	sort.Slice(sortedAddrs, func(i, j int) bool { return sortedAddrs[i] < sortedAddrs[j] })

	var funcs []*Function
	for i, addr := range sortedAddrs {
		nextAddr := uint64(0)
		if i+1 < len(sortedAddrs) {
			nextAddr = sortedAddrs[i+1]
		}

		for _, s := range addrMap[addr] {
			f := &Function{
				Name:       s.Name,
				EntryPoint: addr,
				EndAddr:    nextAddr,
				IsGoFunc:   false,
			}
			funcs = append(funcs, f)
		}
	}

	return funcs
}

func recoverByHeuristics(instructions []disasm.Instruction, entryPoint uint64) []*Function {
	if len(instructions) == 0 {
		return nil
	}

	candidates := make(map[uint64]bool)

	for _, inst := range instructions {
		if inst.IsCall && inst.BranchTarget != 0 {
			candidates[inst.BranchTarget] = true
		}
	}

	candidates[entryPoint] = true

	for i, inst := range instructions {
		if inst.Opcode == "RET" && i+1 < len(instructions) {
			candidates[instructions[i+1].Address] = true
		}
	}

	sortedAddrs := make([]uint64, 0, len(candidates))
	for addr := range candidates {
		sortedAddrs = append(sortedAddrs, addr)
	}
	sort.Slice(sortedAddrs, func(i, j int) bool { return sortedAddrs[i] < sortedAddrs[j] })

	var funcs []*Function
	for i, addr := range sortedAddrs {
		nextAddr := uint64(0)
		if i+1 < len(sortedAddrs) {
			nextAddr = sortedAddrs[i+1]
		}
		funcs = append(funcs, &Function{
			Name:       fmt.Sprintf("func_%x", addr),
			EntryPoint: addr,
			EndAddr:    nextAddr,
			IsGoFunc:   false,
		})
	}

	return funcs
}

// RefinePackagesByCallGraph analyzes the call graph from the disassembly
// to correctly place unexported (lowercase) functions. In Go, a lowercase
// function can only be called from within its defining package, so if all
// callers of a lowercase function belong to the same package, the function
// belongs to that package. Stdlib and runtime callers are excluded.
func (r *RecoverResult) RefinePackagesByCallGraph(instructions []disasm.Instruction) {
	if len(r.Functions) == 0 || len(instructions) == 0 {
		return
	}

	funcStarts := make([]uint64, len(r.Functions))
	for i, f := range r.Functions {
		funcStarts[i] = f.EntryPoint
	}

	entryToFunc := make(map[uint64]*Function, len(r.Functions))
	for _, f := range r.Functions {
		entryToFunc[f.EntryPoint] = f
	}

	findCaller := func(addr uint64) *Function {
		idx := sort.Search(len(funcStarts), func(i int) bool {
			return funcStarts[i] > addr
		}) - 1
		if idx < 0 {
			return nil
		}
		if idx+1 < len(r.Functions) {
			if addr >= funcStarts[idx+1] {
				return nil
			}
		}
		return r.Functions[idx]
	}

	isLowercase := func(s string) bool {
		if len(s) == 0 {
			return false
		}
		for _, r := range s {
			return unicode.IsLower(r)
		}
		return false
	}

	callerPkgs := make(map[uint64]map[string]struct{})

	for _, inst := range instructions {
		if !inst.IsCall || inst.BranchTarget == 0 {
			continue
		}
		callee, ok := entryToFunc[inst.BranchTarget]
		if !ok || callee.Classification != ClassUser {
			continue
		}
		if !isLowercase(callee.ShortName) {
			continue
		}

		caller := findCaller(inst.Address)
		if caller == nil {
			continue
		}
		if caller.Classification == ClassStdlib || caller.Classification == ClassRuntime {
			continue
		}

		if callerPkgs[callee.EntryPoint] == nil {
			callerPkgs[callee.EntryPoint] = make(map[string]struct{})
		}
		callerPkgs[callee.EntryPoint][caller.PackagePath] = struct{}{}
	}

	changed := false
	for _, f := range r.UserFunctions {
		pkgs := callerPkgs[f.EntryPoint]
		if len(pkgs) == 1 {
			for pkg := range pkgs {
				if pkg != f.PackagePath && pkg != "" {
					f.PackagePath = pkg
					changed = true
				}
				break
			}
		}
	}

	if changed {
		r.Packages = make(map[string][]*Function)
		for _, f := range r.UserFunctions {
			r.Packages[f.PackagePath] = append(r.Packages[f.PackagePath], f)
		}
	}
}

func readPtr(data []byte, ptrSize uint8) uint32 {
	if ptrSize == 8 {
		return uint32(encodingbinary.LittleEndian.Uint64(data[:8]))
	}
	return encodingbinary.LittleEndian.Uint32(data[:4])
}

func readPtrV116(data []byte) uint32 {
	return encodingbinary.LittleEndian.Uint32(data[:4])
}

func funcTabSize(ptrSize uint8) uint32 {
	return uint32(ptrSize) * 2
}

func funcTabSizeV116() int64 {
	return 16
}

func funcTabSizeV12() uint32 {
	return 8
}

// extractModuleFromSymbols inspects function symbols to determine the
// Go module path. It looks for names like "example.com/mod/pkg.Func"
// and returns the longest common prefix before the first slash.
func extractModuleFromSymbols(syms []binary.Symbol) string {
	moduleCounts := make(map[string]int)
	for _, s := range syms {
		if s.Type != binary.SymbolFunction || s.Name == "" {
			continue
		}
		if strings.HasPrefix(s.Name, "main.") || strings.HasPrefix(s.Name, "runtime.") {
			continue
		}
		lastDot := strings.LastIndex(s.Name, ".")
		if lastDot < 0 {
			continue
		}
		pkgPath := s.Name[:lastDot]
		if pkgPath == "" || strings.ContainsAny(pkgPath, "()") {
			continue
		}
		// Find the module root (first component before slash or dot)
		slashIdx := strings.Index(pkgPath, "/")
		dotIdx := strings.Index(pkgPath, ".")
		moduleEnd := len(pkgPath)
		if slashIdx > 0 {
			moduleEnd = slashIdx
		} else if dotIdx > 0 {
			moduleEnd = dotIdx
		}
		module := pkgPath[:moduleEnd]
		if module != "" && !isStdlibModule(module) {
			moduleCounts[module]++
		}
	}

	best := ""
	bestCount := 0
	for m, c := range moduleCounts {
		if c > bestCount {
			best = m
			bestCount = c
		}
	}
	return best
}

func isStdlibModule(name string) bool {
	stdlib := []string{
		"fmt", "sync", "os", "io", "time", "net", "http",
		"strings", "bytes", "strconv", "errors", "log", "flag",
		"math", "sort", "path", "context", "reflect", "regexp",
		"encoding", "crypto", "bufio", "archive", "compress",
		"database", "debug", "go", "hash", "image", "index",
		"mime", "text", "unicode", "syscall", "container",
		"internal", "runtime",
	}
	for _, s := range stdlib {
		if name == s {
			return true
		}
	}
	return false
}

// isValidModulePath checks if a string looks like a valid Go module path
// (no control characters, starts with a letter).
func isValidModulePath(path string) bool {
	if len(path) == 0 || len(path) > 256 {
		return false
	}
	if path[0] < 'a' || path[0] > 'z' {
		return false
	}
	for _, ch := range path {
		if ch < 0x20 || ch > 0x7e {
			return false
		}
	}
	return true
}
