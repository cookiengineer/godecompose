# godecompose

**Pattern-based decompiler for Go binaries.** Recovers original Go source code by matching known compiler output patterns against disassembled machine code.

Unlike classical decompilers that try to reverse-engineer what was lost during compilation, godecompose identifies _known patterns_ in the assembly output and maps them back to source code. It works especially well for Go because Go binaries are statically linked, include symbol tables by default, and the `pclntab` section provides exact function boundaries and names.

## Example

Take this Go program that uses `fmt`, `strings`, and `os`:

```go
// main.go
package main

import (
    "fmt"
    "os"
    "strings"
)

func greet(name string) string {
    s := fmt.Sprintf("hello, %s!", name)
    return strings.ToUpper(s)
}

func main() {
    for _, name := range os.Args[1:] {
        fmt.Println(greet(name))
    }
}
```

### Build it

```bash
$ GOOS=linux GOARCH=amd64 go build -o app .
```

### Decompile it

```bash
$ godecompose disasm app

Decoded 180725 instructions
Functions: 1966 total
  runtime:  1414 (skipped)
  stdlib:   219 (skipped)
  user:     2          ← only your code matters

User functions:
  main.greet @ 0x4a11a0 (blocks: 8)
  main.main  @ 0x4a1220 (blocks: 9)
```

```bash
$ godecompose decompile app

// main.greet:
s := fmt.Sprintf($format, $args...);   ← recovered from CALL fmt.Sprintf(SB)
strings.ToUpper($s);                   ← recovered from CALL strings.ToUpper(SB)

// main.main:
fmt.Println($args...);                 ← recovered from CALL fmt.Fprintln(SB)
```

Each recovered line maps directly to a pattern in the database. The `$format`, `$args`, and `$s` are captured register variables from the pattern's `bind` block — in a full project generation pass (`--output`), these get resolved to concrete variable names.

### Generate a project

```bash
$ godecompose decompile app --output=./recovered/

./recovered/
├── go.mod
├── main.go              ← package main with reconstructed main() and greet()
```

Since `greet` is lowercase (unexported), the callgraph analysis determines it must belong to the `main` package — the only package where it can be called. Struct types and their methods are similarly placed in the correct package based on method callers.

For programs with multiple packages:

```bash
$ godecompose decompile myapp --output=./recovered/

./recovered/
├── go.mod
├── main.go              ← package main
├── utils/
│   └── utils.go         ← recovered utils package with struct definitions
└── models/
    └── models.go        ← recovered models package with struct definitions
```

## Quick Start

```bash
# Install
git clone https://github.com/cookiengineer/godecompose
cd godecompose
go build ./cmd/godecompose/

# Show binary metadata
./godecompose info ./myprogram

# Disassemble and show user functions (skipping runtime/stdlib)
./godecompose disasm ./myprogram

# Decompile to stdout
./godecompose decompile ./myprogram

# Decompile to a Go project directory
./godecompose decompile ./myprogram --output=./recovered/
```

## How It Works

### 1. Parse the binary

Supports ELF (Linux), PE (Windows), and Mach-O (macOS) formats for x86_64. Extracts sections, symbol tables, Go build info, and the pclntab (PC-line table) that gives exact function boundaries.

### 2. Disassemble with symbol resolution

Uses `golang.org/x/arch/x86/x86asm` to decode x86_64 machine code. Builds a `SymLookup` from the binary's symbol table so that CALL instructions resolve to human-readable names like `CALL fmt.Fprintln(SB)` instead of `CALL 0x49fc25`.

### 3. Recover and classify functions

The pclntab provides exact function entry points. Symbol names are classified:
- **User code**: `main.*`, or functions matching the detected module prefix
- **Runtime**: `runtime.*`, `type:.*`, and internal Go symbols
- **Stdlib**: `fmt.*`, `sync.*`, `os.*`, `net.*`, `encoding/*`, etc.

Only user code is decompiled. Runtime and stdlib are skipped.

Method receivers (e.g. `(*IntHeap).Push` or `User.Method`) are parsed from symbol names. **Callgraph analysis** then refines package placement: lowercase (unexported) functions are reassigned to their caller's package using Go visibility rules, and struct types are placed based on the consensus of their methods' packages.

### 4. Match patterns

The pattern database (350+ patterns across all Go stdlib packages) describes known compiler output sequences:

```rust
// patterns/libs/golang/highlevel/highlevel.hexpat
pattern fmt.Printf {
    name: "fmt.Printf";
    instr match {
        CALL fmt_Fprintf    // matches CALL fmt.Fprintf(SB) in disassembly
    }
    gen {
        fmt.Printf($format, $args...);
    }
    bind { format as "fmt"; args as "args"; }
}
```

Patterns match by:
- **CALL target name**: fuzzy-matches `CALL fmt_Fprintf` against `CALL fmt.Fprintf(SB)` or `CALL fmt.(*pp).doPrintf(SB)`
- **Instruction sequence**: matches multi-instruction patterns like the `runtime.memmove` forward/backward copy sequence
- **Operand capture**: extracts register names and immediate values as variables for source reconstruction

### 5. Generate source code

Matched patterns expand their `gen` templates with captured variable bindings. Unmatched instruction ranges are emitted as assembly comments.

### 6. Write project (optional)

Functions are grouped by their Go package path (extracted from symbol names like `myprogram/utils.Greet`). A `main.go` is generated for the entry point, and sub-package directories are created for each recovered package. Recovered struct types emit stub definitions with method receiver syntax (`func (r *Type) Method() { ... }`).

## Requirements

| Requirement | Details |
|---|---|
| Architecture | x86_64 (amd64) |
| Binary format | ELF, PE, or Mach-O |
| Symbol table | **Required** for high-level patterns. Go binaries include this by default (`go build`). Stripping with `-ldflags="-s"` removes symbols. |
| Debug info | Not needed. DWARF sections are unused. |
| Go version | Works with Go 1.16+ binaries (pclntab format evolved but is handled) |

## Pattern Language

Godecompose uses an ImHex-compatible pattern description language with extensions for decompilation. Pattern files use the `.hexpat` extension.

### Basic structure

```rust
arch x86_64;
platform linux, darwin;

pattern my_pattern {
    name: "display name";
    library: "identifier";
    description: "what this matches";

    // Assembly view: what instructions to match
    instr block_name {
        MOVQ src, dst       // capture variables: src, dst
        CALL function_name  // match by symbol name
    }

    // Source view: what to generate
    gen {
        myFunction($dst, $src);
    }

    // Variable renaming
    bind {
        src as "source";
        dst as "dest";
    }
}
```

### Pattern types

- **CALL patterns** (high-level): `CALL fmt_Fprintln` — matches any CALL to a function whose symbol contains `fmt.Fprintln`. Requires symbol table.
- **Instruction sequence patterns** (low-level): Multi-instruction sequences with register capture. Works even without symbol table.
- **Alternative patterns**: Use `|` for alternative instruction sequences within the same pattern.

## Project Structure

```
godecompose/
├── cmd/godecompose/          # CLI tool
├── actions/                  # Reusable decompilation pipeline steps
├── types/                    # Arch, Platform enums
├── binary/                   # Binary format interface + Open() dispatcher
│   ├── elf/                  # ELF parser
│   ├── pe/                   # PE/COFF parser
│   └── macho/                # Mach-O parser
├── disasm/                   # x86_64 disassembler + CFG builder
│   └── goasm/                # Go Plan 9 assembly dialect support
├── function/                 # Function recovery (pclntab, classification, callgraph, structs)
├── pattern/
│   ├── lang/                 # Pattern language engine (lexer/parser/AST/evaluator)
│   ├── matcher/              # Instruction pattern matcher
│   └── generate/             # Source code + project generator
├── database/                 # Pattern database loader + syscall tables
│   └── syscall/              # Syscall table types
├── patterns/                 # Pattern files (.hexpat) + syscall tables (.json)
│   ├── kernels/              # Linux, Windows, Darwin, FreeBSD syscall tables
│   └── libs/golang/          # Go runtime + stdlib decompilation patterns
├── docs/                     # Technical documentation
├── e2e/                      # End-to-end integration tests
├── testdata/                 # Test Go programs for cross-compilation
└── goutil/                   # Test helpers (cross-compilation)
```

## Contributing

### Adding new patterns

1. Create a `.hexpat` file in the appropriate directory under `patterns/libs/golang/`
2. Write patterns using the pattern language (see existing files for examples)
3. Validate: `./godecompose patterns validate path/to/file.hexpat`
4. Write an E2E test that compiles a Go program using the functions and verifies pattern matches

### Adding syscall tables

Edit the JSON files in `patterns/kernels/<os>/syscall_table.json` and add entries:

```json
{"number": 123, "name": "syscall_name", "args": "arg descriptions", "returns": "return type"}
```

## License

AGPL-3.0
