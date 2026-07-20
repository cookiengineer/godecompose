# Godecompose Architecture

## Overview

Godecompose is a **pattern-based decompiler** built as a Go library with a CLI frontend. It recovers Go source code from compiled binaries by matching known compiler output patterns against disassembled machine code.

The pipeline flows: **binary → parse → disassemble → recover functions → classify → callgraph refine → [DFA analysis] → [structural analysis] → match patterns → generate source**.

The two new analysis stages (added in Phase 1) run between callgraph refinement and pattern generation:

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
                     │    │ Callgraph refine  │    │
                     │    │ Lowercase fns →   │    │
                     │    │ caller's pkg      │    │
                     │    │ Build structs     │    │
                     │    │ Refine struct pkgs│    │
                     │    └─────────┬─────────┘    │
                     │              │              │
                     │    ┌─────────▼─────────┐    │
                     │    │ analysis/dfa      │ NEW│
                     │    │ Plan9 op parser   │    │
                     │    │ Inst→Value graph  │    │
                     │    │ Copy prop/DCE     │    │
                     │    └─────────┬─────────┘    │
                     │              │              │
                     │    ┌─────────▼─────────┐    │
                     │    │ disasm/structure  │ NEW│
                     │    │ Region analysis   │    │
                     │    │ if/else diamond   │    │
                     │    │ for loop detect   │    │
                     │    │ switch/case       │    │
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
                          │ DFA Go expressions  │
                          │ Structured CFG emit │
                          │ Struct stubs        │
                          │ Method receivers    │
                          │ Project generation  │
                          └─────────────────────┘
```

## Component Design

### 1. Binary Parsers (`binary/`, `binary/elf/`, `binary/pe/`, `binary/macho/`)

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
- `0x7F ELF` → `binary/elf.Open()`
- `MZ` → `binary/pe.Open()`
- `0xFEEDFACE` / `0xCFFAEDFE` / `0xCAFEBABE` → `binary/macho.Open()`

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
    Opcode        string       // Plan 9 syntax: "MOVQ", "CALL", "JEQ", "JGT", "JLT"
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

**Opcode extraction**: `extractOpcode()` derives Plan 9 opcodes from the `GoSyntax` string (first token). Condition codes are mapped to Go Plan 9 convention (JE→JEQ, JG→JGT, JL→JLT, JB→JLO, JA→JHI, JBE→JLS, etc.). TEST and CMP opcodes normalize across operand-size variants (TESTQ/TESTL/TESTB → TEST, CMPQ/CMPL/CMPB → CMP). REP prefix opcodes extract the underlying instruction (REP; MOVSQ → MOVSQ).

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
- Go 1.18+: Magic `0xFFFFFFF0`/`0xFFFFFFF1`, generics-aware

**Symbol merging**: Matches pclntab entry points against symbol table addresses to assign function names.

**Method receiver parsing**: `ParsePackageName()` extracts package path, function name, and optional method receiver from Go symbol names:

| Symbol | Package | Function/Receiver |
|---|---|---|
| `main.greet` | `main` | func `greet` |
| `main.Type.Method` | `main` | `(Type).Method` |
| `main.(*Type).Method` | `main` | `(*Type).Method` |
| `pkg/path.Func` | `pkg/path` | func `Func` |
| `pkg/path.Type.Method` | `pkg/path` | `(Type).Method` |

Additional fields on `Function`: `ReceiverType`, `IsMethod`, `IsPointerReceiver`.

**Function classification**: Each recovered function is classified:

| Class | Criteria | Examples |
|---|---|---|
| `ClassUser` | `main.*` or module-prefixed | `main.main`, `mymod/pkg.Func` |
| `ClassRuntime` | `runtime.*`, `type:.*`, `_rt0_*` | `runtime.memmove` |
| `ClassStdlib` | `fmt.*`, `sync.*`, `encoding/*`, etc. | `fmt.Println` |
| `ClassInternal` | `internal/*` | `internal/poll.FD.Init` |
| `ClassVendor` | Other dotted names | Third-party deps |

**User function filtering**: Only `ClassUser` functions pass through to the pattern matcher. Runtime and stdlib functions are skipped, reducing instruction count from ~180K to ~80 for a typical program.

**Callgraph refinement** (`RefinePackagesByCallGraph`): Analyzes all CALL instructions across the binary to correctly place unexported (lowercase) functions. In Go, a lowercase function can only be called from within its defining package, so if all callers of a lowercase function belong to the same package, the function is reassigned there. Stdlib and runtime callers are excluded from analysis.

**Struct tracking** (`BuildStructs`, `RefineStructPackages`): Groups user functions by their receiver type into `StructType` definitions. Each struct's package path is determined by consensus of its methods' packages. The callgraph refinement ensures lowercase methods (and their containing structs) are placed in the calling package.

**Module name extraction**: The Go module path is detected from:
1. `GoBuildInfo.Main` (if build info is parsed correctly)
2. `GoBuildInfo.Path` (if valid, not a Go version string)
3. Longest common prefix from non-stdlib symbol names (fallback)

### 3b. Data Flow Analysis (`analysis/dfa/`)

Symbolically executes Plan9 Go assembler instructions per basic block to build a value graph. Replaces raw assembly comments with Go expressions.

**Plan9 operand parser** (`parse.go`): Parses Go Plan9 assembler syntax operands: registers (`AX`, `BX`, `R8`), immediates (`$42`, `$0x3F`), memory references (`8(SP)`, `(AX)(BX*8)`, `main.x(SB)`).

**Instruction translation** (`translate.go`): Dispatches by normalized opcode to build `Value` nodes in a block-local state machine. Covers MOV, ADD, SUB, AND, OR, XOR, SHL, SHR, INC, DEC, NEG, NOT, LEA, CMP, TEST, CALL, Jcc, SETcc, RET.

**Optimization** (`optimize.go`): Copy propagation, dead store elimination, constant folding, identity simplification (e.g., `x+0`→`x`, `x*1`→`x`).

**Code emission** (`emit.go`): Renders `Value` trees to Go source text. Unrecognized instructions fall back to assembly comments.

### 3c. Structural Analysis (`disasm/structure.go`)

Post-dominator-based control flow structuring for proper Go emission.

**Post-dominator computation**: Lengauer-Tarjan on reversed CFG with virtual exit node. Enables merge-point detection for if/else diamonds.

**Region detection** (in `pattern/generate/generate.go`):
- **If/else diamonds**: Conditional block with two successors merging at IPD — emits `} else {` with inline else body
- **For loops**: Loop headers detected via back-edges from dominator analysis — emits `for { ... }` or `for cond { ... }`
- **If-else-if chains**: Consecutive conditional blocks — emits `else if { ... }`
- **Condition extraction**: Uses DFA-analyzed CMP/TEST to reconstruct Go comparison expressions

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
- Gen template expansion with `$variable` → alias substitution via bindings
- Gen block brace balancing: nested `{`/`}` tracked by depth for Go-style gen output
- `@if`/`@for` for compile-time template conditionals; plain `if`/`for` emitted as literal text
- GenLoop evaluation: repeats gen body based on loop condition
- Whitespace insertion between adjacent gen statements for readable output

### 5. Pattern Matcher (`pattern/matcher/`)

**Pre-filtering**: Patterns are indexed by their first instruction's opcode (`byOpcode` map). For each CALL instruction, only CALL-based patterns are considered.

**CALL matching**: `instructionMatches()` uses a multi-strategy fuzzy matcher:

1. **Exact substring**: Check if the GoSyntax contains the pattern's expected function name
2. **Case-insensitive**: Lowercase both the GoSyntax and the expected name
3. **Prefix matching**: Normalize separators (`.`, `(`, `)`, `/`, `*`, `_` → space), split into words, match each target word as a prefix against GoSyntax words

This handles Go's symbol name variations:
```
Pattern:  sync_Mutex_Lock  → target "sync.Mutex.Lock"
GoSyntax: CALL sync.(*Mutex).lockSlow(SB)  → normalizes to "call sync mutex lockslow sb"
Target:   "sync mutex lock"  → each word found as prefix → MATCH

Pattern:  runtime_chansend1  → target "runtime.chansend1"
GoSyntax: CALL runtime.chansend1(SB)  → normalizes to "call runtime chansend1 sb"
Target:   "runtime chansend1"  → MATCH (underscores split into separate words)
```

**Operand matching**: `matchOperands()` compares pattern operand constraints against disassembled instruction operands (Intel syntax). Supports:
- Wildcard (`_`): matches anything
- Immediate (`$imm`): matches immediate values
- Register (`RAX`, `X0`): case-insensitive exact register match
- Capture variable (`src`): captures the matched operand value

**Conflict resolution**: Matches are sorted by confidence (longer, more specific patterns score higher). Overlapping matches are resolved by preferring the highest-confidence match.

### 6. Code Generator (`pattern/generate/`)

**Template expansion**: `expandTemplate()` substitutes `$variable` placeholders with captured values or binding aliases.

**Flat output** (`Generate()`): Produces a single text stream with matched gen templates and unresolved code comments.

**Project output** (`WriteProject()`): Groups functions by Go package path and writes a directory structure:
- `go.mod` with the detected module name
- `main.go` for the `main` package with `func main()` entry point
- Sub-package directories with `.go` files for each recovered package
- **Struct stubs**: For each recovered struct type, an empty struct definition with a `// fields unknown` comment
- **Method receivers**: Functions that are methods include their receiver syntax (`func (r *Type) Method() { ... }`) in generated output

### 7. Pattern Database (`database/`)

**Loading**: `LoadPatternsFromDir()` recursively walks a directory, lexes/parses/evaluates each `.hexpat` file, and adds compiled patterns to the database.

**Indexing**: Patterns are indexed by first opcode for fast matcher pre-filtering.

**Filtering**: `FindPatterns(arch, platform)` returns patterns matching the target binary's architecture and platform.

**Syscall tables**: JSON files in `patterns/kernels/` provide per-platform syscall number→name mappings. Four tables included:
- Linux x86_64: 137 syscalls
- Windows NT 10.0: 121 syscalls
- Darwin/macOS: 70 syscalls
- FreeBSD: 57 syscalls

### 8. CLI (`cmd/godecompose/`) and Actions (`actions/`)

The `actions/` package provides reusable decompilation pipeline steps with descriptive parameters. The CLI layer (`cmd/godecompose/main.go`) handles argument parsing, binary opening, and database loading, then delegates work to the actions package.

| Command | Description |
|---|---|
| `info <binary>` | Format, arch, sections, symbols, Go build info |
| `disasm <binary>` | Disassemble, recover functions, show user code summary |
| `decompile <binary> [--output=<dir>]` | Full pipeline with pattern matching and source generation |
| `patterns list` | List loaded patterns and syscall tables with statistics |
| `patterns validate <file>` | Validate a `.hexpat` pattern file through lex→parse→validate→evaluate |

**Actions API**:
- `actions.Info(b binary.Binary) error`
- `actions.Disassemble(b binary.Binary) error`
- `actions.DecompileBinary(b binary.Binary, db *database.Database) (*DecompileOutput, error)`
- `actions.PatternsList(db *database.Database) error`
- `actions.PatternsValidate(filePath string) error`
- `actions.WriteProject(output *DecompileOutput, dir string) error`

## Key Design Decisions

### Why data flow analysis in addition to pattern matching?

Pattern matching identifies high-level operations (function calls, if/else structure) but cannot reconstruct the Go expressions between them. A function like `result := x + y; fmt.Println(result)` compiles to `MOVQ RAX, DI; ADDQ RBX, DI; MOVQ DI, 8(SP); LEAQ 8(SP), AX; CALL fmt.Println`. Only the CALL matches a pattern. The MOV/ADD/LEA instructions become assembly comments.

The DFA (`analysis/dfa/`) fills this gap by symbolically executing each instruction to build a value graph, then emitting Go expressions from that graph. This turns the 10,000+ unresolved assembly lines into readable Go statements.

### Why post-dominator-based structuring?

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

## Project Structure

```
godecompose/
├── actions/              # Reusable decompilation pipeline steps
├── analysis/             # Data flow analysis for Go expression reconstruction
│   └── dfa/              # Intra-block symbolic execution + optimization + emission
├── binary/               # Common binary interface + format parsers
│   ├── elf/              # ELF binary parser
│   ├── pe/               # PE/COFF binary parser
│   └── macho/            # Mach-O binary parser
├── cmd/godecompose/      # CLI entry point
├── database/             # Pattern database + syscall tables (JSON)
│   └── syscall/          # syscalls
│       └── tables/       # syscall tables (JSON)
├── disasm/               # x86_64 disassembler + CFG + structural analysis
├── docs/                 # Documentation
├── e2e/
│   ├── e2e_test.go               # Binary parsing/disassembly E2E tests
│   ├── decompile/               # Per-package decompilation E2E tests
│   │   ├── fmt/fmt_test.go
│   │   ├── sync/sync_test.go
│   │   ├── dfa_simple/
│   │   ├── phase1_forloop/
│   │   ├── phase1_ifelse/
│   │   ├── phase1_switch/
│   │   ├── phase1_structs/
│   │   └── ...
│   └── internal/decompile/     # Shared test helpers
├── function/             # Function recovery (pclntab, classification, callgraph, structs, signatures)
├── goutil/               # Go compilation test utilities
├── pattern/              # Pattern language engine (lang/, matcher/, generate/)
├── patterns/             # Pattern files (.hexpat)
│   ├── golang/
│   │   ├── stdlib/       # Go stdlib patterns (one subdir per package)
│   │   ├── runtime/      # Go runtime patterns
│   │   ├── fallback/     # Single-CALL high-level patterns
│   │   ├── controlflow/  # Control flow, idioms, data types, stdlib calls
│   │   └── embed.go      # //go:embed all four modules
├── testdata/src/         # Test Go source programs (one subdir per package)
├── types/                # Arch/Platform enums
├── go.mod
└── go.sum
```

### Adding a New Stdlib Pattern

To add a new stdlib pattern, create three files in parallel directories:

1. **Pattern**: `patterns/libs/golang/stdlib/<pkg>/<pkg>.hexpat`
   - Add a `pattern` block matching the CALL instruction
2. **Test source**: `testdata/src/<pkg>/main.go`
   - Go program exercising the package functions
3. **E2E test**: `e2e/decompile/<pkg>/<pkg>_test.go`
   - Import `e2e/internal/decompile`, call `CompileAndOpen` + `Decompile`, assert pipeline OK

The 1:1:1 mapping ensures every pattern has a corresponding end-to-end test.
