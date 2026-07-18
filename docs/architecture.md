# Godecompose Architecture

## Project Purpose

Godecompose is a **decompiler framework** focused on decompilation of amd64 (x86_64) binaries via pattern matching against a pattern database, rather than classical iterative decompilation.

The core idea: instead of trying to reverse engineer what was lost between source → compilation → assembly, we identify _known patterns_ in the assembly output and map them back to source code. This approach is particularly effective for:
- Statically linked Go binaries (large, but with recognizable runtime patterns)
- Libraries with known compilation patterns (OpenSSL, zlib, etc.)
- Syscall wrappers that follow predictable calling conventions

## Design Principles

1. **Framework-first**: All packages are public and importable. The `cmd/godecompose` CLI is one consumer of the framework.
2. **Pure Go**: No CGo dependencies. Disassembly uses `golang.org/x/arch/x86/x86asm`.
3. **Pattern-driven**: Decompilation quality scales with the pattern database, not with algorithm sophistication.
4. **Positive-style branches**: Early returns, guard clauses, minimal nesting.
5. **Human-readable code**: Descriptive variable names, clear function boundaries, no clever tricks.

## Pipeline

```
Binary File (ELF / PE / Mach-O)
    │
    ▼
┌─────────────────────────────────┐
│ binary package                  │  Parse binary format headers
│ └─ elf / pe / macho subpackages │  Extract sections, symbols, entry point
└───────────────┬─────────────────┘
                │ .text section bytes
                ▼
┌─────────────────────────────────┐
│ disasm package                  │  Decode x86_64 instructions
│ └─ DecodeStream() → []Inst      │  Produce linear instruction stream
│ └─ BuildCFG() → []BasicBlock    │  Construct control flow graph
│ └─ goasm subpackage             │  Go Plan 9 dialect handling, ABI detection
└───────────────┬─────────────────┘
                │ []Inst + CFG
                ▼
┌─────────────────────────────────┐
│ function package                │  Recover function boundaries
│ └─ ParsePclntab() (Go binaries) │  Go: use pclntab for precise boundaries
│ └─ HeuristicRecovery()          │  Generic: prologue/epilogue heuristics
└───────────────┬─────────────────┘
                │ []Function with basic blocks
                ▼
┌─────────────────────────────────┐
│ pattern package                 │  ImHex-compatible pattern language
│ └─ lang/ (lexer, parser, AST)   │  Lex and parse .hexpat pattern files
│ └─ matcher/ (instruction match) │  Match compiled patterns against instream
│ └─ generate/ (source output)    │  Generate source code from matched patterns
└───────────────┬─────────────────┘
                │ Matched patterns + variable bindings
                ▼
┌─────────────────────────────────┐
│ database package                │  Pattern database management
│ └─ syscall subpackage           │  Syscall tables (Linux, Windows, Darwin, BSD)
│ └─ Pattern index + lookup       │  Fast opcode-indexed pattern retrieval
└───────────────┬─────────────────┘
                │ Compiled patterns, syscall metadata
                ▼
┌─────────────────────────────────┐
│ cmd/godecompose (CLI consumer)  │  User-facing commands
│ └─ disasm, decompile, patterns  │  Orchestrates the pipeline
└─────────────────────────────────┘
```

## Package Layout

```
godecompose/
├── cmd/godecompose/          # CLI tool (consumes the framework)
│   └── main.go
│
├── types/                    # Shared types: Arch, Platform, enums
│
├── binary/                   # Common interface + format detection (public API)
│   ├── binary.go             # Binary interface, Section, Symbol, registry, Open()
│   ├── detect_test.go        # Format detection unit tests
│   └── binary_test.go        # Integration tests (external package)
│
├── elf/                      # ELF parser (wraps debug/elf, imports binary)
├── pe/                       # PE/COFF parser (wraps debug/pe, imports binary)
├── macho/                    # Mach-O parser (wraps debug/macho, imports binary)
│
├── disasm/                   # Disassembly layer (public API)
│   ├── disasm.go             # Core: Instruction, DecodeStream, branch detection
│   ├── cfg.go                # Control flow graph construction
│   ├── goasm/                # Go Plan 9 assembly dialect support
│   └── syntax/               # Assembly syntax formatters
│
├── function/                 # Function boundary recovery (public API)
│   ├── function.go           # Function, Variable, RecoverResult types
│   └── pclntab.go            # Go pclntab parser + symbol/heuristic recovery
│
├── goutil/                   # Test helpers for cross-compilation
│
├── e2e/                      # End-to-end integration tests
│
├── pattern/                  # Pattern language engine (public API)
│   ├── lang/                 # ImHex-compatible pattern language
│   │   ├── token/            # Token types and keywords
│   │   ├── lexer/            # Hand-written lexer
│   │   ├── parser/           # Recursive descent parser with backtracking
│   │   ├── ast/              # AST node hierarchy (30+ node types)
│   │   ├── preprocessor/     # #include, #define, #pragma support
│   │   ├── validator/        # Type checking and semantic analysis
│   │   └── evaluator/        # Tree-walking interpreter
│   ├── matcher/              # Assembly instruction pattern matching engine
│   └── generate/             # Source code generation from matched patterns
│
├── database/                 # Pattern database (public API)
│   ├── database.go           # Loader, indexer, query interface
│   └── syscall/              # Syscall table types + per-platform data
│
├── patterns/                 # Shipped pattern files
│   ├── kernels/              # Syscall tables per kernel
│   │   ├── linux/
│   │   ├── windows/
│   │   ├── darwin/
│   │   └── freebsd/
│   └── libs/                 # Library-specific decompilation patterns
│       ├── golang/
│       └── openssl/
│
├── testdata/                 # Test binaries for integration tests
├── docs/                     # Technical documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Key Design Decisions

### Why Pattern Matching Instead of Classical Decompilation?

Classical decompilers (Hex-Rays, Ghidra) use iterative algorithms:
1. Lift assembly to an intermediate representation (IR)
2. Apply simplification passes (SSA, copy propagation, dead code elimination)
3. Perform type recovery via constraint solving
4. Structure control flow (gotos → if/while/for)
5. Emit C-like pseudo-code

This works well for C/C++ but struggles with:
- Go binaries (unusual calling convention, runtime-coupled code)
- Heavily optimized code (inlined functions, vectorized loops)
- Statically linked binaries (massive code volume)

Pattern matching inverts the problem: instead of trying to _reconstruct_ what was lost, we _recognize_ what was produced. This is faster and more reliable for common patterns, at the cost of being incomplete for novel code.

### Why ImHex Compatibility?

[ImHex PatternLanguage](https://github.com/WerWolv/PatternLanguage) is:
- Well-documented and actively maintained
- Already designed for binary data pattern description
- Has a familiar C-like/Rust-like syntax
- Can potentially share patterns with the ImHex ecosystem (binary format descriptions)

We extend it with assembly-matching and source-generation constructs (`instr`, `gen`, `bind` blocks).

### Why Pure Go?

- Single binary distribution (no C libraries to link)
- Cross-compilation trivial (`GOOS=linux GOARCH=amd64 go build`)
- Go's toolchain already includes production-quality x86 decoding (`golang.org/x/arch`)
- No build system complexity (no CMake, no CGo)

### Why Go Binaries First?

Go binaries have distinctive characteristics that make them ideal for pattern-matching decompilation:
- The `pclntab` section provides exact function boundaries and names
- The runtime is large but highly predictable (same patterns across all Go programs)
- Go's SSA-based compiler produces consistent code patterns for language constructs
- Statically linked → no dynamic linking ambiguities
