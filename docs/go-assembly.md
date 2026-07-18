# Go Assembly Dialect (Plan 9 Assembler for amd64)

## Overview

Go uses the **Plan 9 assembler** syntax, which differs significantly from both AT&T and Intel syntax. The key insight: Go's assembler is a **semi-abstract instruction set**, not a direct 1:1 mapping to machine code. The toolchain performs instruction selection _after_ code generation, meaning some pseudo-instructions may be expanded or replaced by the linker.

For godecompose, we disassemble machine code using `golang.org/x/arch/x86/x86asm`, which provides `GoSyntax()` for Plan 9 output and `IntelSyntax()` for traditional output.

## Operand Order

**Left-to-right (source, destination)**:

```
MOVQ $1, AX    ; Move immediate 1 into AX (source first, then destination)
ADDQ BX, CX    ; CX = CX + BX
```

This is the opposite of Intel syntax (`mov rax, 1`) and different from AT&T (`movq $1, %rax`).

## Operand Size Suffixes

Suffixes on opcodes indicate operand width:

| Suffix | Width   | Example  |
|--------|---------|----------|
| `B`    | 8-bit   | `MOVB`   |
| `W`    | 16-bit  | `MOVW`   |
| `L`    | 32-bit  | `MOVL`   |
| `Q`    | 64-bit  | `MOVQ`   |
| `O`    | 128-bit | `MOVO` (SSE/AVX) |

Some opcodes are always 64-bit on amd64 and don't need a suffix:

```
SYSCALL, CPUID, SWAPGS, LFENCE, MFENCE, SFENCE, PAUSE
```

## Register Naming

### Hardware Registers (Go Plan 9 Names)

| Go Name | Standard Name | Notes                  |
|---------|---------------|------------------------|
| `AX`    | RAX           | General purpose        |
| `BX`    | RBX           | General purpose        |
| `CX`    | RCX           | General purpose        |
| `DX`    | RDX           | General purpose        |
| `SI`    | RSI           | Source index           |
| `DI`    | RDI           | Destination index      |
| `BP`    | RBP           | Frame pointer          |
| `SP`    | RSP           | Stack pointer          |
| `R8`    | R8            | Extended GP            |
| `R9`    | R9            | Extended GP            |
| `R10`   | R10           | Extended GP            |
| `R11`   | R11           | Extended GP            |
| `R12`   | R12           | Scratch (ABIInternal)  |
| `R13`   | R13           | Scratch (ABIInternal)  |
| `R14`   | R14           | Goroutine pointer (g)  |
| `R15`   | R15           | Scratch / GOT pointer  |

8-bit versions: `AL`, `CL`, `DL`, `BL`, `SIB`, `DIB`, `R8B`-`R15B`

16-bit versions: `AX`, `CX`, `DX`, `BX`, `SP`, `BP`, `SI`, `DI`, `R8W`-`R15W`

SSE/AVX: `X0`-`X15` (and `X16`-`X31` with AVX-512), `Y0`-`Y31`, `Z0`-`Z31`

### Pseudo-Registers (Virtual, Not Hardware)

These are maintained by the toolchain and do not correspond to any hardware register:

| Pseudo | Meaning                          | Usage                                      |
|--------|----------------------------------|--------------------------------------------|
| `FP`   | Virtual frame pointer            | Access function arguments: `first_arg+0(FP)` |
| `SP`   | Virtual stack pointer            | Local variables: `x-8(SP)` (name REQUIRED) |
| `SB`   | Static base                      | Global symbols: `foo(SB)`, `foo+4(SB)`     |
| `PC`   | Program counter                  | Branch targets / labels                    |

**Critical distinction**: On amd64, `x-8(SP)` (with name prefix) references the _virtual_ stack pointer for local variables. `$-8(SP)` (with `$` prefix) or bare `-8(SP)` references the _hardware_ SP register. This is a common source of confusion.

### Special Register Roles (Go-specific)

```go
RARG  = -1   // First integer argument (used as BP in ABI0)
REGRET  = AX  // Integer return value register
FREGRET = X0  // Floating-point return value register
REGSP   = SP  // Stack pointer
REGCTXT = DX  // Closure context pointer
REGG    = R14 // Current goroutine pointer (g)
REGEXT  = R15 // External register (allocated from R15 down)
R12, R13      // Permanent scratch (ABIInternal)
X15           // Zero register (always zero in ABIInternal)
```

## Addressing Modes

```
(R1)              ; Register indirect: dereference R1
8(R1)             ; Offset: [R1 + 8]
(DI)(BX*2)        ; Scaled index: [DI + BX*2]  (scale: 1, 2, 4, 8)
64(DI)(BX*2)      ; Full SIB: [DI + BX*2 + 64]
$foo(SB)          ; Absolute address of global symbol foo
foo+8(SB)         ; Offset from global symbol foo
first_arg+0(FP)   ; Function argument (name REQUIRED on FP references)
x-8(SP)           ; Local variable (name REQUIRED on virtual SP)
$-8(SP)           ; Hardware SP (no name prefix, $ prevents virtual interpretation)
$42               ; Immediate constant
```

## Go Calling Conventions (x86_64)

Go has two coexisting ABIs:

### ABI0 (Stack-based, Legacy)

- All arguments and results passed on stack
- `BP` serves as `RARG` (first argument register alias)
- Return value in `AX` (integer) or `X0` (floating-point)
- Used by hand-written assembly with Go prototypes

```
TEXT ·myFunc(SB), NOSPLIT, $0-16
    MOVQ arg1+0(FP), AX     ; Load first arg (8 bytes at FP+0)
    MOVQ arg2+8(FP), BX     ; Load second arg (8 bytes at FP+8)
    ADDQ BX, AX              ; AX = AX + BX
    MOVQ AX, ret+16(FP)     ; Store return value (8 bytes at FP+16)
    RET
```

### ABIInternal (Register-based, Go 1.17+)

Uses registers for arguments and results:

**Integer argument registers** (in order): `RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11`

**Floating-point argument registers** (in order): `X0`-`X14`

Result registers use the same sequences.

**Fixed-purpose registers**:
- `RDX` — Closure context pointer (`REGCTXT`)
- `R14` — Current goroutine pointer `g` (`REGG`)
- `X15` — Zero register (always contains 0)
- `R12, R13` — Permanent scratch
- `R15` — Scratch (GOT pointer in dynamically linked binaries)
- `RBP` — Frame pointer (callee-save)

**Stack layout**: The caller reserves "spill space" on the stack for all register-based arguments (so the callee can spill them if needed). Stack grows downward, 8-byte aligned.

### Detecting the ABI

For a given function, we determine the ABI by analyzing the prologue:

- **ABIInternal**: Function accesses arguments from registers at entry (`MOVQ AX, ...` or `MOVQ BX, ...`)
- **ABI0**: Function accesses arguments via FP-relative addressing (`MOVQ arg+0(FP), ...`)

In practice, most modern Go binaries (Go 1.17+) use ABIInternal for compiler-generated code. ABI0 is primarily used by hand-written assembly in the runtime.

## Go Compiler `-S` Output Format

When using `go tool compile -S`, output lines have this structure:

```
0x0000 00000 (main.go:3)  MOVQ $42, AX
0x000a 00005 (main.go:4)  CALL runtime.printlock(SB)
```

Columns:
1. **Byte offset** from function start (hex)
2. **PC offset** (instruction counter, for branch targets)
3. **Source position** `(file:line)`
4. **Instruction** with operands

After linking (`go tool objdump`):

```
main.go:3  0x10501c0  48c7c02a000000  MOVQ $0x2a, AX
```

Format: `<source:line> <absolute_address> <hex_bytes> <GoSyntax>`

## Assembler Directives

| Directive  | Purpose                              | Example                                    |
|------------|--------------------------------------|--------------------------------------------|
| `TEXT`     | Define a function entry point        | `TEXT ·Add(SB), NOSPLIT, $0-24`           |
| `DATA`     | Initialize static data               | `DATA ·message+0(SB)/8, $"Hello, "`       |
| `GLOBL`    | Make symbol globally visible         | `GLOBL ·counter(SB), NOPTR, $8`           |
| `FUNCDATA` | GC metadata (compiler-generated)     | `FUNCDATA $0, gclocals·xxx(SB)`           |
| `PCDATA`   | PC-relative GC metadata              | `PCDATA $0, $1`                           |
| `NOP`      | No-op (may be removed by linker)     | `NOP`                                      |
| `PCALIGN`  | Align next instruction               | `PCALIGN $32`                              |
| `BYTE`     | Raw byte in instruction stream       | `BYTE $0x0f`                              |

### TEXT Flags

```
NOPROF   = 1     ; Don't profile this function
DUPOK    = 2     ; Multiple definitions allowed
NOSPLIT  = 4     ; Don't need stack split preamble
RODATA   = 8     ; Read-only data
NOPTR    = 16    ; No pointers in data
WRAPPER  = 32    ; Wrapper function
NEEDCTXT = 64    ; Closure uses context register
LOCAL    = 128   ; Local to DSO
TLSBSS   = 256   ; Thread-local storage
NOFRAME  = 512   ; No stack frame
TOPFRAME = 2048  ; Outermost frame (traceback termination)
```

## Condition Codes (Plan 9 Convention)

Go uses the Plan 9 convention for branch conditions:

| Go Asm | Intel      | Meaning              |
|--------|------------|----------------------|
| `JOS`  | `JO`       | Overflow set         |
| `JOC`  | `JNO`      | Overflow clear       |
| `JCS`  | `JB`/`JC`  | Carry set (Below)    |
| `JCC`  | `JNB`/`JNC`| Carry clear          |
| `JEQ`  | `JZ`/`JE`  | Equal (Zero)         |
| `JNE`  | `JNZ`/`JNE`| Not equal            |
| `JLS`  | `JBE`      | Lower or same        |
| `JHI`  | `JA`       | Higher (Above)       |
| `JMI`  | `JS`       | Minus (Sign)         |
| `JPL`  | `JNS`      | Plus (Not sign)      |
| `JPS`  | `JP`       | Parity set           |
| `JPC`  | `JNP`      | Parity clear         |
| `JLT`  | `JL`       | Less than            |
| `JGE`  | `JGE`      | Greater or equal     |
| `JLE`  | `JLE`      | Less or equal        |
| `JGT`  | `JG`       | Greater than         |

Same naming applies to `SETcc` and `CMOVcc`.

## Symbol Naming

- In assembly source files: `fmt·Printf` (middle dot U+00B7, not period)
- In compiler `-S` output: `fmt.Printf` (regular period)
- Within same package: `·FunctionName` (linker prepends package path)
- File-static symbols: `foo<>(SB)` (visible only in current file)

Example:

```asm
TEXT "".main(SB), NOSPLIT, $0
    CALL fmt·Println(SB)
    RET
```

The `"".main` means the symbol is in the unnamed (main) package.

## Implications for Disassembly

### What We Get from `x86asm.GoSyntax()`

`x86asm.GoSyntax(inst, pc, symname)` produces valid Go Plan 9 assembly syntax. Given a `SymLookup` function, it resolves PC-relative addresses to symbol names:

```go
symname := func(addr uint64) (string, uint64) {
    // Look up addr in symbol table
    // Return (symbol_name, symbol_base_address)
}

syntax := x86asm.GoSyntax(inst, pc, symname)
// e.g.: "MOVQ $42, AX" or "CALL runtime.printlock(SB)"
```

### What We Need to Handle Ourselves

1. **Pseudo-register classification**: `GoSyntax()` uses bare register names. We need to detect when `SP` refers to the virtual stack pointer (local variable access) vs the hardware register based on context (preceding TEXT directive, symbol names).

2. **FUNCDATA/PCDATA stripping**: These are compiler-generated metadata and should be filtered from the decoded instruction stream. They appear as NOPs in the decoded output.

3. **Linker transformations**: The Go linker can fold branches, duplicate code, and elide NOPs. Our disassembly must handle:
   - Synthetic branches (jumps inserted by linker)
   - NOP padding (for alignment, may be elided)
   - Function prologue duplication (linker copies prologues in some cases)

4. **TEXT directive reconstruction**: `x86asm.Decode()` only decodes machine instructions, not assembler directives. We reconstruct `TEXT` directives from:
   - Symbol table entries (function names + addresses)
   - pclntab entries (function boundaries)
   - Frame size analysis (for `$framesize-argsize`)

5. **ABI detection from register usage**: We analyze which registers a function reads at entry to determine ABI0 vs ABIInternal. This affects how we interpret argument access patterns.
