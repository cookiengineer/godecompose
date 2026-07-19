# Roadmap

Phase 10 completed. 668 patterns (stdlib + runtime + fallback + controlflow), 4 syscall tables, full decompilation pipeline with project generation. Every stdlib package has end-to-end tests.

## Phase 1: Foundation — Binary Format Parsers ✓
[...]

## Phase 9: User-Code Focus and Project Generation ✓

**Completed.** Filters runtime noise, recovers package structure, and generates Go project directories.

### Completed Tasks

- [x] Function classification: `ClassUser`, `ClassRuntime`, `ClassStdlib`, `ClassInternal`, `ClassVendor`
- [x] User-code-only decompilation: skip runtime/stdlib, match patterns only against user function instructions
- [x] Package path extraction from symbol names: `testproject/utils.Greet` → pkg `testproject/utils`, name `Greet`
- [x] Module name detection: GoBuildInfo + symbol-based fallback
- [x] Project directory generation: `go.mod` + `main.go` + sub-package `.go` files
- [x] CLI `--output=<dir>` flag for project output
- [x] E2E tests verify user-function filtering and project generation
- [x] Method receiver parsing: `(*Type).Method` and `Type.Method` patterns detected from symbols
- [x] Callgraph-based package refinement: lowercase functions placed in caller's package
- [x] Struct type tracking: methods grouped by receiver, consensus package assignment
- [x] Struct stubs generated in project output with method receiver syntax
- [x] `actions/` package: reusable decomposition pipeline steps (`Info`, `Disassemble`, `DecompileBinary`, `PatternsList`, `PatternsValidate`)
- [x] CLI refactored to thin parameter-parsing layer delegating to actions

### Realized API

- `function.ParsePackageName(fullName) (pkgPath, shortName, receiverType, isPointerReceiver, isMethod)`
- `function.Function` fields: `ReceiverType`, `IsMethod`, `IsPointerReceiver`
- `function.StructType` — groups methods by receiver, tracks package consensus
- `function.RecoverResult.Structs` — recovered struct types
- `function.RecoverResult.BuildStructs()` — groups user methods by receiver
- `function.RecoverResult.RefineStructPackages()` — consensus package for structs
- `function.RecoverResult.RefinePackagesByCallGraph(instructions)` — callgraph-based placement
- `function.ClassifyFunction(name, mainPackage) Classification`
- `RecoverResult.Packages` — functions grouped by package path
- `generate.WriteProject(dir, goModule)` — writes complete Go project directory with struct stubs and method receivers
- `actions.DecompileBinary(b, db) (*DecompileOutput, error)` — reusable pipeline
- `actions.WriteProject(output, dir) error` — project generation

---

**Completed.** Parses ELF, PE, and Mach-O binaries through a unified API.

### Completed Tasks

- [x] Initialize Go module (`github.com/cookiengineer/godecompose`)
- [x] Implement `types/` package (Arch, Platform enums)
- [x] Implement `binary/binary.go` — common interface (Binary, Section, Symbol), registry-based `Open()`
- [x] Implement `binary/elf/` — ELF parser wrapping `debug/elf`, with Go build info (V1/V2) and pclntab extraction
- [x] Implement `binary/pe/` — PE/COFF parser wrapping `debug/pe`
- [x] Implement `binary/macho/` — Mach-O parser wrapping `debug/macho`
- [x] Unit tests: format detection, ELF integration test
- [x] E2E tests: cross-compile Go → parse ELF/PE/Mach-O → verify sections, symbols, build info

### Realized API

- `binary.Open(path) (Binary, error)` — auto-detects format and returns parsed binary
- Each Section exposes Name, Address, Size, Data, Flags
- Symbol table access, Go build info (version, deps, settings), pclntab extraction
- Registry pattern: format packages self-register via `init()`

---

## Phase 2: Disassembler ✓

**Completed.** Decodes x86_64 machine code, builds control flow graphs, supports Go Plan 9 dialect.

### Completed Tasks

- [x] Add `golang.org/x/arch` dependency
- [x] Implement `disasm/disasm.go` — core types: Instruction, DecodeStream (resilient, skips bad bytes)
- [x] Implement `disasm/BuildControlFlowGraph()` — CFG construction with leader detection
- [x] Implement `disasm/goasm/` — Go Plan 9 specific handling:
  - Register name mapping (x86asm.Reg → Go asm names)
  - Pseudo-register detection (FP, SP, SB, PC)
  - ABI detection (ABI0 stack-based vs ABIInternal register-based vs SystemV)
  - Special register classification (goroutine pointer, closure context, zero register, scratch)
- [x] Implement `disasm/syntax/` — format helpers (Intel, Go Plan 9, GNU/AT&T output, condition codes)
- [x] Unit tests: NOP, CALL, JMP, JE decoding, branch target resolution, CFG construction
- [x] E2E tests: decode 180K+ instructions from .text, build 245+ basic blocks

---

## Phase 3: Function Recovery ✓

**Completed.** Identifies function boundaries using Go pclntab metadata and heuristics.

### Completed Tasks

- [x] Implement `function/function.go` — Function type (Name, EntryPoint, EndAddr, Blocks, Args, Returns)
- [x] Implement `function/pclntab.go` — Go pclntab parser:
  - Detect pclntab format version (Go 1.2, 1.16, 1.18+)
  - Parse function table entries (entry points)
  - Merge symbol table names with pclntab entry points
- [x] Implement `function/heuristic.go` — Generic function boundary detection:
  - CALL target tracking
  - RET-based boundary detection
  - Entry point identification
- [x] E2E tests: recover 1966 functions from simple/complex Go binaries, verify `main.factorial` and `main.main`

---

## Phase 4: Pattern Language — Lexer, Parser, AST ✓

**Completed.** ImHex-compatible pattern language with full godecompose extensions.

### Completed Tasks

- [x] Implement `pattern/lang/token/token.go` — 50+ token types, keyword registry, source positions
- [x] Implement `pattern/lang/lexer/lexer.go` — Hand-written single-pass lexer:
  - Number literals (decimal, hex `0x`, octal `0o`, binary `0b`, digit separators)
  - String/char literals with escape sequences
  - 35+ compound operators (greedy max-length matching)
  - Nested block comment support (`/* outer /* inner */ */`)
  - Directives (#include, #define, etc.)
- [x] Implement `pattern/lang/ast/` — 40+ AST node types:
  - Literals, identifiers, expressions (binary, unary, ternary, cast, call, index, scope)
  - Statements (compound, conditional, while, for, return, break, match, try/catch)
  - Declarations (struct, union, enum, bitfield, using, fn, namespace, import)
  - Variables (plain, array, pointer, assignment)
  - **Godecompose extensions**: PatternDefinition, InstrBlock, InstructionPattern, OperandPattern, MemoryRefPattern, GenBlock, GenText, GenExpr, GenConditional, GenLoop, BindBlock, ArchDirective, PlatformDirective
- [x] Implement `pattern/lang/parser/parser.go` — Recursive descent with backtracking:
  - Full operator precedence table (13 levels)
  - All statement and declaration types
  - Pattern definitions with metadata, instr/gen/bind blocks
  - Assembly instruction line detection via opcode heuristic
  - Register detection, capture variables, addressing modes `(base)(index*scale)`
  - Backtracking for cast vs. parenthesized expression disambiguation
- [x] 40 tests: 6 token + 17 lexer + 17 parser

### Deliverables

- `token.LookupKeyword(ident)` — keyword detection
- `lexer.New(source).Lex() ([]Token, error)` — tokenization
- `parser.New(tokens).Parse() (*ast.Program, error)` — parsing
- Full ImHex base language support + all godecompose extensions

### Phase 4e: Preprocessor ✓

- [x] Implement `pattern/lang/preprocessor/`:
  - `#include "file"` with file resolver interface
  - `#define NAME [value]` with lexed replacement tokens
  - `#undef NAME`
  - `#ifdef`/`#ifndef`/`#endif` with nested conditional support
  - `#pragma` (pass-through)
  - `#error "message"` with error propagation
- [x] 10 tests: define, undef, ifdef/ifndef, nested, #include with file resolver, #error

---

## Phase 5: Pattern Language — Validator and Evaluator ✓

**Completed.** Type-checks patterns and evaluates them through a tree-walking interpreter.

### Completed Tasks

- [x] Implement `pattern/lang/validator/` — semantic analysis:
  - Scope tracking with push/pop (function bodies, compound statements, loops)
  - Control flow validation (break/continue in loops, return in functions)
  - Empty instruction pattern detection in instr blocks
- [x] Implement `pattern/lang/evaluator/` — tree-walking interpreter:
  - Lexical scope stack with push/pop
  - Variable creation, lookup, assignment
  - Expression evaluation: arithmetic, comparison, logical, string ops
  - Pattern compilation from `instr` blocks → `CompiledPattern` structures
  - Gen template expansion with `$variable` → alias substitution via bindings
  - Builtin function dispatch (`print`)
- [x] 7 validator tests + 4 evaluator tests

### Deliverables

- `validator.Validate(prog) []error` — semantic validation
- `evaluator.Evaluate(prog) ([]*CompiledPattern, error)` — pattern compilation with gen expansion

---

## Phase 6: Pattern Matcher and Code Generator ✓

**Completed.** Matches compiled patterns against instruction streams and generates decompiled source.

### Completed Tasks

- [x] Implement `pattern/matcher/`:
  - Opcode-indexed pre-filtering for fast candidate lookup
  - Multi-instruction sequence matching with operand-level detail
  - Capture variable extraction from matched operands
  - Wildcard support with confidence scoring
  - Overlapping match conflict resolution (highest confidence first)
- [x] Implement `pattern/generate/`:
  - Gen template expansion with captured variable values
  - Alias binding substitution from pattern bind blocks
  - Unmatched instruction ranges emitted as raw assembly comments
  - Sorted match output by address
- [x] 7 matcher tests + 4 generator tests

### Deliverables

- `matcher.New(patterns).Match(instructions) []Match` — instruction stream matching
- `generate.New(matches, instructions).Generate() string` — source code generation

---

## Phase 7: Pattern Database and Syscall Tables ✓

**Completed.** Four kernel syscall tables and initial Go runtime decompilation patterns.

### Completed Tasks

- [x] `database/database.go` — Pattern DB loader, opcode indexer, recursive directory walk, platform/arch filtering, stats
- [x] `database/syscall/syscall.go` — Table and Entry types with lookup methods
- [x] Linux x86_64 syscall table: 137 syscalls (read, write, open, mmap, socket, clone, execve, futex, ...)
- [x] Windows NT 10.0 x86_64 syscall table: 121 syscalls (NtClose, NtCreateFile, NtReadFile, NtWriteFile, NtWaitFor*, ...)
- [x] Darwin/macOS x86_64 syscall table: 70 syscalls (exit, fork, read, write, mach syscalls, ...)
- [x] FreeBSD x86_64 syscall table: 57 syscalls (kqueue, kevent, shm_open2, close_range, ...)
- [x] Initial Go runtime patterns: `runtime.memmove`, `runtime.newobject`

### Deliverables

- `database.LoadPatternsFromDir(dir)` — recursive .hexpat loading
- `database.LoadSyscallsFromDir(dir)` — recursive JSON syscall table loading
- `database.FindPatterns(arch, platform)` — filtered pattern lookup
- `database.SyscallTable(platform)` — per-platform syscall lookup

---

## Phase 8: CLI and Integration ✓

**Completed.** Working command-line decompiler with full pipeline integration. CLI logic extracted into reusable `actions/` package.

### Completed Tasks

- [x] `godecompose info <binary>` — Format, arch, entry point, sections (with flags), symbols, Go build info
- [x] `godecompose disasm <binary>` — Disassembles .text, builds CFG, recovers functions with names
- [x] `godecompose decompile <binary>` — Full pipeline: parse → disasm → load DB → match patterns → generate source
- [x] `godecompose patterns list` — List loaded patterns and syscall tables with stats
- [x] `godecompose patterns validate <file>` — Lex → parse → validate → evaluate a .hexpat file
- [x] `actions/` package with `Info`, `Disassemble`, `DecompileBinary`, `PatternsList`, `PatternsValidate`
- [x] CLI refactored to thin parameter-parsing layer delegating to `actions/*`

### Verified CLI Output

- `info`: Shows 16 sections, Go build info, symbols for ELF binary
- `disasm`: Decodes 378K instructions into 57K basic blocks, recovers 8,261 functions
- `decompile`: Loads 4 syscall tables (137+121+70+57 entries) + 2 patterns, runs full pipeline
- `patterns validate`: Validates .hexpat files through lexer → parser → validator → evaluator

---

## Phase 10: User-Code Control Flow & Data Type Reconstruction ✓

**Completed.** Added control flow recognition (if/else, nil checks, CMP-based comparisons), Go idiom patterns, data type creation patterns, and comprehensive stdlib/runtime call coverage. Improved decompilation match rate from 7 to 2,209 (316x) on the ysco benchmark.

### Completed Tasks

- [x] Fix GenConditional evaluator bug (only first statement of conditional branches was processed)
- [x] Implement GenLoop evaluation (was silently dropped)
- [x] Fix disassembler opcode extraction: Go Plan 9 opcodes from GoSyntax with condition code mapping (JE→JEQ, JG→JGT, etc.)
- [x] Normalize TEST/CMP opcodes across size variants (TESTQ/TESTL/TESTB → TEST, CMPQ/CMPL/CMPB → CMP)
- [x] Gen block parser: balanced brace depth tracking for nested `{`/`}` support
- [x] Gen block parser: `@if`/`@for` for compile-time conditionals/loops; plain `if`/`for` treated as gen text
- [x] Gen block evaluator: whitespace insertion between gen text and gen expressions
- [x] Fuzzy matcher: add underscore (`_`) to normalizer for runtime function names (e.g., `mapassign_faststr`)
- [x] Conflict resolution: optimized from O(n×10000) to O(n×m) interval overlap check
- [x] Performance: 3 O(n²) bottlenecks fixed (instruction collection, function block extraction, conflict resolution)
- [x] Load all 4 pattern modules: stdlib, runtime, fallback, controlflow
- [x] 93 new patterns in 6 files: control flow, Go idioms, data types, runtime extras, stdlib calls

### Realized API

- `disasm.extractOpcode(goSyntax)` — extracts Go Plan 9 opcode from GoSyntax string with size and CC normalization
- `disasm.normalizeOpcode(op)` — Plan 9 condition code mapping (JE→JEQ, JG→JGT, JB→JLO, etc.)
- `evaluator.evalGenStmt` — handles GenConditional (accumulates all statements), GenLoop (repeats body N times)
- `parser.parseGenBlock` — balanced brace depth tracking for nested gen blocks
- `parser.parseGenStatement` — `@if`/`@for` for compile-time, plain `if`/`for` as gen text
- `matcher.fuzzyMatchCall` — splits on `_` in addition to `.()/*` for runtime function matching
- `actions.DecompileBinary` — O(n) address-indexed instruction collection with progress logging

### Benchmarks (ysco decompilation)

| Metric | Before | After |
|--------|--------|-------|
| Patterns | 550 | 668 |
| Matches | 7 | 2,209 |
| Unresolved stdlib/runtime CALLs | 769 | 0 |
| Decompilation time | ∞ (timeout) | ~30s |
