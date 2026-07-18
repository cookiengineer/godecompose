# Godecompose Architecture

## Overview

Godecompose is a **pattern-based decompiler** built as a Go library with a CLI frontend. It recovers Go source code from compiled binaries by matching known compiler output patterns against disassembled machine code.

The pipeline flows: **binary → parse → disassemble → recover functions → match patterns → generate source**.

```
                         ┌─────────────────────┐
                         │   Binary (ELF/PE/   │
                         │   Mach-O)           │
                         └─────────┬───────────┘
                                   │
                         ┌─────────▼───────────┐
                         │  binary.Open()      │
                         │  Format detection   │
                         │  Section extraction │
                         │  Symbol table       │
                         │  Go build info      │
                         │  pclntab section    │
                         └─────────┬───────────┘
                                   │
                    ┌──────────────┼──────────────┐
                    │              │              │
            ┌───────▼──────┐ ┌────▼─────┐ ┌──────▼──────┐
            │ disasm       │ │ function │ │ pattern/lang │
            │ DecodeStream │ │ pclntab  │ │ Lexer/Parser │
            │ GoSyntax     │ │ classify │ │ AST/Evaluate │
            │ SymLookup    │ │ filter   │ │ Compile      │
            └───────┬──────┘ └────┬─────┘ └──────┬──────┘
                    │              │              │
                    │    ┌─────────▼─────────┐    │
                    │    │ ClassifyFunction  │    │
                    │    │ User/Runtime/     │    │
                    │    │ Stdlib/Vendor     │    │
                    │    └─────────┬─────────┘    │
                    │              │              │
                    │    ┌─────────▼─────────┐    │
                    │    │ Filter user fns   │    │
                    │    │ Only user insts   │    │
                    │    └─────────┬─────────┘    │
                    │              │              │
                    │    ┌─────────▼─────────┐    │
                    └────┤ pattern/matcher   ├────┘
                         │ Opcode-indexed    │
                         │ fuzzy CALL match  │
                         │ operand capture   │
                         └─────────┬─────────┘
                                   │
                         ┌─────────▼───────────┐
                         │ pattern/generate    │
                         │ Template expansion  │
                         │ Project generation  │
                         └─────────────────────┘
```

## Component Design

### 1. Binary Parsers (`binary/`, `elf/`, `pe/`, `macho/`)

**Interface**: `binary.Binary` defines the common API for all binary formats.

```go
type Binary interface {
    Format() Format
    Architecture() types.Arch
    EntryPoint() uint64
    Sections() []Section
    Section(name string) (Section, bool)
    Symbols() ([]Symbol, error)
    GoBuildInfo() (*GoBuildInfo, bool)
    Pclntab() ([]byte, uint64, bool)
    IsPIE() bool
    IsStripped() bool
    ByteOrder() binary.ByteOrder
    Close() error
}
```

**Format detection**: `binary.Open(path)` reads the first 4 bytes and dispatches:
- `0x7F ELF` → `elf.Open()`
- `MZ` → `pe.Open()`
- `0xFEEDFACE` / `0xCFFAEDFE` / `0xCAFEBABE` → `macho.Open()`

**Go-specific extraction**:
- ELF: `.go.buildinfo` section (V1/V2 format), `.gopclntab` section, `.note.go.buildid`
- PE: `.go.buildinfo` section, `.gopclntab` section
- Mach-O: `__go_buildinfo` section, `__gopclntab` section

**Registry pattern**: Each format package registers itself via `init()`:
```go
func init() {
    binary.RegisterFormat(binary.FormatELF, func(path string) (binary.Binary, error) {
        return Open(path)
    })
}
```

### 2. Disassembler (`disasm/`)

Uses `golang.org/x/arch/x86/x86asm` for pure-Go x86_64 instruction decoding.

```go
type Instruction struct {
    Address       uint64
    Bytes         []byte
    Opcode        string       // e.g., "MOV", "CALL", "JMP"
    IntelSyntax   string       // Intel syntax: "mov rax, rbx"
    GoSyntax      string       // Plan 9 syntax: "MOVQ BX, AX"
    IsCall        bool
    IsReturn      bool
    IsBranch      bool
    IsConditional bool
    BranchTarget  uint64
    Size          int
}
```

**Symbol resolution**: `DecodeStreamWithSymbols(data, baseAddr, lookup)` passes a `SymLookup` function to `x86asm.GoSyntax()`. This resolves PC-relative CALL targets to symbol names:

```
nil lookup:    CALL 0x49fc25
with lookup:   CALL fmt.Fprintln(SB)
```

**SymLookup construction**: `BuildSymLookup(symbols)` builds an address→name map from the binary's symbol table:

```go
type SymLookup func(addr uint64) (name string, base uint64)
```

**CFG building**: `BuildControlFlowGraph(instructions, entryPoints)` identifies basic block leaders at jump targets, call targets, and RET successors. Builds predecessor/successor edges for control flow analysis.

**Go-specific** (`disasm/goasm/`): Maps x86asm register names to Plan 9 assembler names (RAX→AX, R14→R14), detects ABI (ABI0 vs ABIInternal), classifies special registers (goroutine pointer, closure context, zero register).

### 3. Function Recovery (`function/`)

**pclntab parser**: Supports Go 1.2, 1.16, and 1.18+ pclntab formats. Extracts function entry points from the PC-line table:
- Go 1.2: Magic `0xFFFFFFFA`, fixed-size entries
- Go 1.16: Magic `0xFFFFFFFB`, compact format
- Go 1.18+: Magic `0xFFFFFFF0`/`0xFFFFFFF1`, generics-aware

**Symbol merging**: Matches pclntab entry points against symbol table addresses to assign function names.

**Function classification**: Each recovered function is classified:

| Class | Criteria | Examples |
|---|---|---|
| `ClassUser` | `main.*` or module-prefixed | `main.main`, `mymod/pkg.Func` |
| `ClassRuntime` | `runtime.*`, `type:.*`, `_rt0_*` | `runtime.memmove` |
| `ClassStdlib` | `fmt.*`, `sync.*`, `encoding/*`, etc. | `fmt.Println` |
| `ClassInternal` | `internal/*` | `internal/poll.FD.Init` |
| `ClassVendor` | Other dotted names | Third-party deps |

**User function filtering**: Only `ClassUser` functions pass through to the pattern matcher. Runtime and stdlib functions are skipped, reducing instruction count from ~180K to ~80 for a typical program.

**Module name extraction**: The Go module path is detected from:
1. `GoBuildInfo.Main` (if build info is parsed correctly)
2. `GoBuildInfo.Path` (if valid, not a Go version string)
3. Longest common prefix from non-stdlib symbol names (fallback)

### 4. Pattern Language Engine (`pattern/lang/`)

Implements an ImHex-compatible pattern language with decompilation extensions.

**Pipeline**: `Source → Lexer → Preprocessor → Parser → Validator → Evaluator`

**Lexer**: Hand-written single-pass scanner. Supports:
- All ImHex token types (keywords, value types, operators, separators)
- Numeric literals: decimal, hex `0x`, octal `0o`, binary `0b`
- String/char literals with escape sequences
- 35+ compound operators with greedy max-length matching
- Nested block comments `/* /* */ */`
- Preprocessor directive detection (`#include`, `#define`, etc.)

**Parser**: Recursive descent with backtracking. Full operator precedence table (13 levels). Supports:
- All ImHex constructs: struct/union/enum/bitfield, variables, functions, control flow
- **Godecompose extensions**: `pattern`, `instr`, `gen`, `bind`, `arch`, `platform`
- Assembly-specific: opcode detection heuristic, register recognition, addressing modes `(base)(index*scale)`
- Backtracking via `mark()`/`reset()` for cast vs. parenthesized expression disambiguation

**AST**: 40+ node types. All standard ImHex nodes plus:
- `PatternDefinition`, `InstrBlock`, `InstructionPattern`, `OperandPattern`, `MemoryRefPattern`
- `GenBlock`, `GenText`, `GenExpr`, `GenConditional`, `GenLoop`
- `BindBlock`, `BindEntry`, `ArchDirective`, `PlatformDirective`

**Evaluator**: Tree-walking interpreter. Produces `CompiledPattern` structures:
- Instruction sequences compiled from `instr` blocks
- Gen templates with variable substitution markers
- Binding tables mapping capture variables to aliases

### 5. Pattern Matcher (`pattern/matcher/`)

**Pre-filtering**: Patterns are indexed by their first instruction's opcode (`byOpcode` map). For each CALL instruction, only CALL-based patterns are considered.

**CALL matching**: `instructionMatches()` uses a multi-strategy fuzzy matcher:

1. **Exact substring**: Check if the GoSyntax contains the pattern's expected function name
2. **Case-insensitive**: Lowercase both the GoSyntax and the expected name
3. **Prefix matching**: Normalize separators (`.`, `(`, `)`, `/`, `*` → space), split into words, match each target word as a prefix against GoSyntax words

This handles Go's symbol name variations:
```
Pattern:  sync_Mutex_Lock  → target "sync.Mutex.Lock"
GoSyntax: CALL sync.(*Mutex).lockSlow(SB)  → normalizes to "call sync mutex lockslow sb"
Target:   "sync mutex lock"  → each word found as prefix → MATCH
```

**Operand matching**: `matchOperands()` compares pattern operand constraints against disassembled instruction operands. Supports:
- Wildcard (`_`): matches anything
- Immediate (`$imm`): matches immediate values
- Register (`RAX`, `X0`): exact register match
- Capture variable (`src`): captures the matched operand value

**Conflict resolution**: Matches are sorted by confidence (longer, more specific patterns score higher). Overlapping matches are resolved by preferring the highest-confidence match.

### 6. Code Generator (`pattern/generate/`)

**Template expansion**: `expandTemplate()` substitutes `$variable` placeholders with captured values or binding aliases.

**Flat output** (`Generate()`): Produces a single text stream with matched gen templates and unresolved code comments.

**Project output** (`WriteProject()`): Groups functions by Go package path and writes a directory structure:
- `go.mod` with the detected module name
- `main.go` for the `main` package with `func main()` entry point
- Sub-package directories with `.go` files for each recovered package

### 7. Pattern Database (`database/`)

**Loading**: `LoadPatternsFromDir()` recursively walks a directory, lexes/parses/evaluates each `.hexpat` file, and adds compiled patterns to the database.

**Indexing**: Patterns are indexed by first opcode for fast matcher pre-filtering.

**Filtering**: `FindPatterns(arch, platform)` returns patterns matching the target binary's architecture and platform.

**Syscall tables**: JSON files in `patterns/kernels/` provide per-platform syscall number→name mappings. Four tables included:
- Linux x86_64: 137 syscalls
- Windows NT 10.0: 121 syscalls
- Darwin/macOS: 70 syscalls
- FreeBSD: 57 syscalls

### 8. CLI (`cmd/godecompose/`)

| Command | Description |
|---|---|
| `info <binary>` | Format, arch, sections, symbols, Go build info |
| `disasm <binary>` | Disassemble, recover functions, show user code summary |
| `decompile <binary> [--output=<dir>]` | Full pipeline with pattern matching and source generation |
| `patterns list` | List loaded patterns and syscall tables with statistics |
| `patterns validate <file>` | Validate a `.hexpat` pattern file through lex→parse→validate→evaluate |

## Key Design Decisions

### Why pattern matching instead of classical decompilation?

Classical decompilers (Hex-Rays, Ghidra) lift assembly to IR, apply simplification passes, recover types, structure control flow, and emit pseudo-code. This works for C/C++ but struggles with:
- **Go binaries**: Unusual calling convention, runtime coupling, static linking
- **Heavily optimized code**: Inlined functions, vectorized loops
- **Scale**: Go binaries are large (runtime included) — classical decompilation is slow

Pattern matching inverts the problem: instead of reconstructing what was lost, we recognize what was produced. Go's SSA-based compiler produces consistent code patterns for language constructs, making this approach viable.

### Why symbol-based matching for high-level patterns?

Go binaries include symbol tables by default (`go build` without `-ldflags="-s"`). This means CALL instructions resolve to human-readable names like `CALL fmt.Fprintln(SB)`. By matching against these symbol names, we can identify function calls with high confidence without needing to understand the surrounding instruction sequence.

### Why focus on user code?

A typical Go binary contains ~180K instructions, of which ~95% are runtime and stdlib. Decompiling all of it produces unreadable output. By classifying functions using pclntab names and filtering to user code only, we reduce the decompilation target to ~80 instructions — making the output clean and focused.

### Why pure Go?

- Single binary distribution (no C libraries to link)
- Cross-compilation trivial (`GOOS=linux GOARCH=amd64 go build`)
- Go's `x/arch/x86/x86asm` provides production-quality x86 decoding
- No build system complexity (no CMake, no CGo)

## Dependencies

| Package | Purpose | Pure Go |
|---|---|---|
| `golang.org/x/arch/x86/x86asm` | x86_64 instruction decoder | Yes |
| `debug/elf` (stdlib) | ELF parsing | Yes |
| `debug/pe` (stdlib) | PE/COFF parsing | Yes |
| `debug/macho` (stdlib) | Mach-O parsing | Yes |
| `encoding/json` (stdlib) | Syscall table loading | Yes |
