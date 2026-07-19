# Phase 10: User-Code Control Flow & Data Type Reconstruction

**Status: In Progress** — ~40% of call sites recovered on ysco (2,209 matches). Full user-code reconstruction not yet achieved.

## Current Recovery Rate

On the ysco benchmark (a real-world Go service manager):

| Metric | Original | After Phase 10 | Target |
|--------|----------|----------------|--------|
| Patterns loaded | 550 | 694 | — |
| Pattern matches | 7 | 2,209 | — |
| User function calls identified | 0% | ~40% | 100% |
| Remaining unresolved asm lines | 5,700 | 20,682 | 0 |
| Decompilable Go source | ~0% | ~15% (high-level trace) | 100% |
| Decompilation time | ∞ (timeout) | ~30s | — |

The output is a high-level trace: function boundaries are correct, call sites are labeled with inferred Go operations, and control flow patterns mark if/else branches with offset-based goto targets. However, the output is not valid compilable Go — it contains unresolved assembly blocks and approximate Go statements stitched together.

## What Phase 10 Achieved

### 1. Code Infrastructure Fixes (9 files)

| Change | Why needed |
|--------|-----------|
| Go Plan 9 opcode extraction in disassembler | Patterns use Go assembler names (JEQ/JGT/JLT), disasm used Intel (JE/JG/JL) |
| TEST/CMP size normalization | 8-bit/32-bit/64-bit comparisons normalized to base opcode |
| GenConditional evaluator fix | Only first statement of each branch was processed (silent data loss) |
| GenLoop evaluation | Parse-and-drop → functional for template iteration |
| Gen block balanced braces | Nested `{`/`}` inside gen blocks closed the block prematurely |
| Gen block `@if`/`@for` prefix | Compile-time gen conditionals must not conflict with Go `if`/`for` output |
| Fuzzy matcher underscore splitting | Runtime names use `_` (e.g. `mapassign_faststr`) but matcher only split on dots |
| O(n²) → O(n) performance | 3 nested loops caused ∞ timeout on 1.5M-instruction binary |
| Memory operand matching | Struct field access requires parsing `[reg+offset]` memory refs from Intel syntax |

### 2. New Patterns (10 files, ~127 patterns)

| File | Category | Count |
|------|----------|-------|
| `if_else.hexpat` | Control flow: nil checks, CMP comparisons | 9 |
| `go_idioms.hexpat` | Go idioms: defer, goroutines, channels, maps | 12 |
| `data_types.hexpat` | Data types: new/make, type conversions | 12 |
| `calls_patterns.hexpat` | slog, time, os, syscall, mutex, GC barriers | 18 |
| `runtime_extras.hexpat` | morestack, panics, concatstrings, assertions | 22 |
| `more_calls.hexpat` | flag, net/http, fmt, sconf, prometheus, template | 24 |
| `remaining_calls.hexpat` | os, net, syscall, strings, strconv, time | 18 |
| `more_data.hexpat` | memclr, slicecopy, newarray, gostring, convT2E/I | 18 |
| `struct_fields.hexpat` | Memory operand read/write with offset capture | 5 |
| `switch_case.hexpat` | CMP+JEQ chain for sparse switch statements | 2 |

### 3. E2E Tests (8 test files)

| Test | Matches | Unique Patterns | Status |
|------|---------|-----------------|--------|
| `controlflow` | 50 | 18 | PASS |
| `dataops` | 35 | 11 | PASS |
| `structfields` | 16 | 6 (11 field access matches) | PASS |
| `switchcase` | 43 | 7 (0 switch — uses jump tables) | PASS |
| Existing stdlib (37) | varies | varies | PASS |

---

## Remaining Tasks for Full User-Code Reconstruction

The remaining 20,682 unresolved assembly lines break down as:

| Opcode | Lines | % | What it represents |
|--------|-------|---|--------------------|
| mov | 7,337 | 35% | Register/spill moves, field access, value assignment |
| movups | 2,965 | 14% | SSE struct/array copies (runtime internal) |
| lea | 2,499 | 12% | Address computation, local var references |
| call | 1,218 | 6% | **Pattern matching bug**: many calls SHOULD match but don't |
| int3 | 819 | 4% | Alignment padding |
| nop | 746 | 4% | No-ops (alignment/linker artifacts) |
| xor | 648 | 3% | Variable zeroing |
| cmp | 622 | 3% | Comparisons (should be matched by control flow) |
| jnz/jz/jmp | 1,538 | 7% | Conditional/unconditional branches |
| add/sub/dec/inc | 752 | 4% | Arithmetic, counter increments |
| other | 1,538 | 7% | ret, pop, push, movzx, and, etc. |

### Category A: Pattern Matching Bugs (must fix)

| Task | Description | Impact |
|------|-------------|--------|
| **A1 — CALL pattern matching regression** | ~1,218 CALL lines remain unresolved despite having patterns. The fuzzy matcher and/or conflict resolver is dropping valid matches. Needs root-cause investigation. | Recovers ~400 additional matches |
| **A2 — CMP/TEST + Jcc not matching** | ~622 CMP + 576 TEST lines + 1,538 Jcc lines should match control flow patterns but don't. TEST/CMP operands in Intel syntax (`dword ptr [addr]`, `al`, `byte ptr [reg]`) may not match pattern operands (which expect simple register names). | Recovers ~500 control flow statements |
| **A3 — struct field false positives** | Struct field patterns match ALL MOV-with-memory, including stack accesses (`mov [sp+0x8], rax`). Need to distinguish struct field access (receiver-register-based offsets) from stack locals. | Makes struct field output meaningful |
| **A4 — gen block `/* comment */` output** | Patterns produce comment output (`/* mem read: ... */`) instead of Go code because captured register names aren't Go variables. | Makes output compilable Go |

### Category B: New Pattern Development (more patterns needed)

| Task | Description | Impact |
|------|-------------|--------|
| **B1 — For loop patterns** | Counted loops (`XORL i,i; JMP @check; ...; INCQ i; CMPQ i,n; JLT @body`) and range loops need multi-instruction patterns with captured loop variables. | Recovers iteration constructs |
| **B2 — Switch/case (jump table)** | Dense integer switches compile to jump tables. Need to recognize the jump table pattern (`LEAQ jumptable, base; MOVSXD [base+idx*8], target; ADDQ base, target; JMP target`) and reconstruct case labels. | Recovers switch statements |
| **B3 — Defer/finally cleanup** | Deferred close/cleanup patterns (resource cleanup after main logic). | Recovers defer blocks |
| **B4 — Error return idioms** | `val, err := fn(); if err != nil { return zero, err }` multi-instruction patterns. | Recovers common Go error handling |
| **B5 — Type assertion & nil checks** | `v, ok := x.(T)` and `v != nil` combined patterns. | Recovers interface assertions |
| **B6 — String switch (hash dispatch)** | String switches compile to hash+table dispatch. Complex multi-instruction pattern. | Recovers string switch |

### Category C: Architectural Limitations (requires new capabilities)

| Task | Description | Impact |
|------|-------------|--------|
| **C1 — Variable tracking across instructions** | We can match individual MOV instructions but can't track that `MOV $42, AX` and `MOV AX, local_offset(SP)` are the same variable. Need data flow analysis across basic blocks. | Fundamental for variable declarations |
| **C2 — Type inference from runtime calls** | `CALL runtime.convTstring` means a string was boxed into `interface{}`. Need to propagate type information backward through the instruction stream. | Recovers Go types |
| **C3 — Function signature reconstruction** | From calling convention (ABIInternal: AX,BX,CX,DI,SI = first 5 args) and return value placement, reconstruct `func name(args) returns`. | Recovers function prototypes |
| **C4 — Struct field naming from offsets** | Multiple methods access the same offset → that offset is a field. Cross-function analysis to group offsets into struct definitions with inferred field names. | Recovers struct type definitions |
| **C5 — Control flow structuring** | Current if/else patterns produce flat `goto $addr` output. Need to structure basic blocks into proper nested `if/else { ... }` with balanced braces. Requires dominator tree analysis or interval structuring. | Makes output compilable Go |
| **C6 — Block-level code generation** | Patterns match individual instructions but functions need coherent block-level generation. Current `generate.go` concatenates pattern matches with unresolved asm comments. Need a generator that understands basic blocks and can emit structured code. | Output reads like Go, not annotated asm |
| **C7 — Compiler-optimized code** | Inlined functions, loop unrolling, dead code elimination all obscure the original Go source. Need deoptimization passes (outline inlined code, recognize loop patterns despite unrolling). | Handles optimized builds |

### Category D: Infrastructure & Tooling

| Task | Description |
|------|-------------|
| **D1 — Pattern discovery tool** | `godecompose patterns discover <source.go>` that compiles a Go file with `-S`, extracts instruction sequences, and generates candidate `.hexpat` files. |
| **D2 — Match debugging tool** | `godecompose patterns explain <binary>` that shows WHY a pattern did or didn't match at a given address (fuzzy match trace, operand comparison, conflict resolution audit). |
| **D3 — Recovery rate metrics** | Automated measurement of % statements recovered, % functions reconstructed, % types identified for benchmark binaries. |

---

## Path to 100% Recovery

Achieving full user-code reconstruction requires a fundamental architectural evolution beyond pattern matching:

### Phase A: Fix Current Bugs (near-term)
- Fix CALL pattern matching regression (Category A1)
- Fix CMP/TEST + Jcc operand matching (Category A2)
- Add for loop and switch/jump-table patterns (Category B1, B2)
- Estimated recovery improvement: 40% → 60% of call sites

### Phase B: Basic Block Analysis (medium-term)
- Implement control flow structuring (Category C5)
- Replace flat pattern matching with block-level code generation (Category C6)
- Add variable liveness tracking across basic blocks (Category C1)
- Estimated recovery improvement: 60% → 75% of statements

### Phase C: Data Flow & Type Recovery (long-term)
- Implement type inference from runtime call context (Category C2)
- Build struct field offset consensus across methods (Category C4)
- Reconstruct function signatures from ABI analysis (Category C3)
- Estimated recovery improvement: 75% → 90% of statements

### Phase D: Deoptimization (advanced)
- Outline inlined functions back to original call sites (Category C7)
- Recognize and reconstruct loop structures
- Handle compiler-specific optimizations (bounds check elimination, nil check removal)
- Estimated recovery improvement: 90% → 95%+ of statements

---

## Known Limitations of Current Architecture

1. **Pattern matching is fundamentally local**: Each pattern matches a window of instructions independently. There's no cross-pattern communication (e.g., variable X in one pattern is the same variable X in another).

2. **Register-to-variable mapping is absent**: The disassembler works with hardware registers (AX, BX, R8...). Go variables live in registers or on the stack. Without register allocation analysis, we can't name variables.

3. **Control flow is flat**: The matcher doesn't understand basic blocks, dominator trees, or loop structures. Branch targets are hex addresses, not structured labels.

4. **Types are opaque**: Runtime type information exists in the binary but we don't extract it. Calls like `runtime.convTstring` reveal the type but we don't propagate that information.

5. **Inlined code is lost**: The Go compiler aggressively inlines small functions. The original function call is gone and can't be recovered without deoptimization.

6. **Gen blocks can't reference match data**: The `gen` template is compiled at pattern-load time, before any matches exist. `$variable` substitution is simple string replacement. Gen conditionals (`@if`) can only use compile-time metadata, not match-time data.
