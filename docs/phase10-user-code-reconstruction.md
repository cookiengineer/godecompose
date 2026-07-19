# Phase 10: User-Code Control Flow & Data Type Reconstruction

**Status: Phase D Complete** — 709 patterns, 3,141 matches on ysco (80.9% call-site recovery, 98.6% function signatures, 100% struct fields). 4 phases (A: bugs, B: blocks, C: signatures/types, D: tooling). 44 E2E suites + 16 unit tests.

## Current Recovery Rate

On the ysco benchmark (a real-world Go service manager):

| Metric | Original | After Phase 10 | Target |
|--------|----------|----------------|--------|
| Patterns loaded | 550 | 701 | — |
| Pattern matches | 7 | 3,141 | — |
| User function calls identified | 0% | ~40% | 100% |
| Remaining unresolved asm lines | 5,700 | 20,682 | 0 |
| Decompilable Go source | ~0% | ~15% (high-level trace) | 100% |
| Decompilation time | ∞ (timeout) | ~30s | — |

The output is a high-level trace: function boundaries are correct, call sites are labeled with inferred Go operations, and control flow patterns mark if/else branches with offset-based goto targets. However, the output is not valid compilable Go — it contains unresolved assembly blocks and approximate Go statements stitched together.

## What Phase 10 Achieved

### 1. Code Infrastructure Fixes (15 files modified, 2 new)

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
| ParsePackageName rewrite | Main-package closures (`main.cmdInit.Printf.func3`) went to separate folders; pointer/value receivers embedded in package path were missed |
| `operandParts` opcode stripping | Opcode was included as first operand part (`["test", "rax", "rax"]`), silently breaking all multi-instruction patterns |
| CFG-based block labeling (C5/C6) | Flat pattern concatenation → per-block emission with labels, `if`/`goto`, noise filtering |
| Function signature reconstruction (C3) | ABI register analysis in entry + return blocks; normalizeReg maps EAX→RAX, AL→RAX |
| Type inference (C2) | Runtime calls (convTstring→string, convT64→int64); name heuristics (is→bool, Error→error) |
| Struct field naming (C4) | Cross-method offset consensus; hex-valid filter <0x200; byte→bool, qword→int |
| Pattern discovery tool (D1) | `godecompose patterns discover <source.go>` — compiles with `-S`, extracts CALLs, generates `.hexpat` |
| Recovery rate metrics (D3) | Instructions (13.1%), Functions (98.6%), Structs (100%), Call sites (80.9%) on ysco |

### 2. New Patterns (12 files, 159 patterns)

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
| `final_calls.hexpat` | filepath.Base, time.Time.*, os.MkdirAll, json.Marshal | 10 |
| `loops_errors_assertions.hexpat` | For loops, error returns, type assertions | 8 |

### 3. E2E Tests (44 test files — 37 stdlib + 7 Phase 10)

| Test | Matches | Unique Patterns | Status |
|------|---------|-----------------|--------|
| `controlflow` | 50 | 18 | PASS |
| `dataops` | 35 | 11 | PASS |
| `structfields` | 7 | 6 | PASS (7 Point methods, struct recovered, fields inferred) |
| `switchcase` | 43 | 7 | PASS |
| `funcgroup` | 16 | — | PASS (methods with receivers) |
| `looperror` | 30 | 11 | PASS (7 user funcs, for/error/assert patterns) |
| `structout` | 5 | — | PASS (3 functions, packages verified, signatures) |
| Existing stdlib (37) | varies | varies | PASS |

### 4. Unit Tests (16 new tests across 3 files)

| File | Count | What it verifies |
|------|-------|-----------------|
| `disasm/structure_test.go` | 5 | Nil, single-block, if-else, loop, return classification |
| `pattern/generate/generate_phaseb_test.go` | 3 | Noise filtering, block labels, condition extraction |
| `function/signature_test.go` | 8 | normalizeReg, extractRegister, extractMemOffset, ReconstructSignature (4 variants), InferStructFields, NameHeuristics |

---

## Remaining Tasks for Full User-Code Reconstruction

The remaining unresolved assembly lines are primarily:

| Category | What remains | Why hard |
|----------|-------------|----------|
| MOV/LEA/XOR | ~10,000 lines | Hardware registers vs Go variable names; need data flow analysis (Category C1) |
| NOP/INT3/data16 | ~1,700 lines | Alignment/padding — can be filtered as noise |
| CMP/TEST + Jcc | ~1,200 lines | Many now matched; remaining ones have complex memory operands |
| Virtual dispatch CALLs | ~800 lines | CALL through register (rax/rcx/etc.) — need itab analysis (C2) |
| Stack prologue/epilogue | ~400 lines | `runtime.morestack` + push/pop — function infrastructure
| add/sub/dec/inc | 752 | 4% | Arithmetic, counter increments |
| other | 1,538 | 7% | ret, pop, push, movzx, and, etc. |

### Category A: Pattern Matching Bugs

| Task | Status |
|------|--------|
| A1 — CALL pattern matching regression | ✓ Fixed — was function grouping bug and operandParts opcode issue |
| A2 — CMP/TEST + Jcc operand matching | ✓ Fixed — operandParts now strips opcode; `if v != nil` patterns match |
| A3 — Struct field false positives | Fixed — GC barrier pattern added; mem patterns use comments |
| A4 — Gen block produces comments, not Go code | Open — requires cross-pattern variable tracking (Category C1) |

### Category B: New Pattern Development

| Task | Description | Status |
|------|-------------|--------|
| B1 — For loop patterns | Counted loop `(i=0; i<n; i++)` and range loop patterns | ✓ Created (`loops_errors_assertions.hexpat`) — fragile due to compiler variance |
| B2 — Switch/case (jump table) | Jump table dispatch pattern | ✓ Created (`switch_case.hexpat`) — limited to CMP+JEQ chains |
| B3 — Defer/finally cleanup | Defer patterns | ✓ Existing patterns (`go_idioms.hexpat`) |
| B4 — Error return idioms | `val, err := fn(); if err != nil { return err }` | ✓ Created (`loops_errors_assertions.hexpat`) — multi-inst, fragile |
| B5 — Type assertion patterns | `v := x.(T)` and `v, ok := x.(T)` | ✓ Created (`loops_errors_assertions.hexpat`) |
| B6 — String switch (hash dispatch) | String switch via hash tables | Open — complex multi-instruction pattern needed |

### Category C: Architectural Limitations (requires new capabilities)

| Task | Description | Impact |
|------|-------------|--------|
| **C5 — Control flow structuring** | ✓ Phase B: CFG-based block labeling with `goto` resolution. Conditional blocks emit `if condition { ... }` with conditions extracted from matched patterns. Branch targets replaced with labeled blocks. |
| **C6 — Block-level code generation** | ✓ Phase B: `writeFunctionBody` uses per-block emission. Instruction noise filtered (NOP/INT3/DATA16). Blocks grouped with labels. Pattern matches per-block with indented body. |
| **C2 — Type inference from runtime calls** | ✓ Phase C: `inferBodyTypes()` scans function body for runtime calls (convTstring→string, convT64→int64, newobject→interface{}, makeslice→[]T). `inferReturnFromName()` detects bool returns from `is`/`has`/`can` prefixes. `nameSuggestsError()` detects error returns from `Error` suffix. |
| **C4 — Struct field naming from offsets** | ✓ Phase C: `InferStructFields()` analyzes memory operands across all methods of a struct, groups by offset, filters to plausible struct offsets (< 0x200, hex valid), infers types (byte→bool, qword→int, convT calls→string/int64). Output: `type State struct { field_50 int; field_69 bool; ... }`. |
| **C7 — Compiler-optimized code** | Inlined functions, loop unrolling, dead code elimination obscure original Go source. E2E tests use `-gcflags=all=-l` to disable inlining for accurate function recovery. |

### Category D: Infrastructure & Tooling

| Task | Description |
|------|-------------|
| **D1 — Pattern discovery tool** | ✓ Phase D: `godecompose patterns discover <source.go>` — compiles with `-S`, extracts CALL targets, generates candidate `.hexpat` files |
| **D2 — Match debugging tool** | Not started |
| **D3 — Recovery rate metrics** | ✓ Phase D: Instructions (13.1%), Functions (98.6%), Structs (100%), Call sites (80.9%) measured on ysco |

---

## Path to 100% Recovery

Achieving full user-code reconstruction requires a fundamental architectural evolution beyond pattern matching:

### Phase A: Fix Current Bugs (near-term) ✓ COMPLETED
- ✓ Fix CALL pattern matching (A1 + A2)
- ✓ Add for loop and switch patterns (B1, B2)
- ✓ Add error return and type assertion patterns (B4, B5)
- ✓ 7 → 3,141 matches (449x improvement)

### Phase B: Basic Block Analysis (medium-term) — IN PROGRESS
- ✓ Implement control flow block labeling (C5): basic blocks labeled L0, L1, ...; branch targets resolved; conditional blocks emit `if condition { ... }` with conditions extracted from matched patterns
- ✓ Replace flat generation with block-level generation (C6): `writeFunctionBody` uses per-block emission with buildControlFlowGraph; NOP/INT3 noise filtered
- □ Add variable liveness tracking across basic blocks (C1) — requires data flow analysis

### Phase C: Data Flow & Type Recovery — COMPLETE
- ✓ Function signature reconstruction from ABI register analysis (C3)
- ✓ Type inference from runtime calls and function name heuristics (C2)
- ✓ Struct field naming via cross-method offset consensus with type inference (C4)
- □ Full type inference from runtime call context (deeper analysis of data flow)
- □ Deoptimization — outline inlined functions, reconstruct loop structures (C7)

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
