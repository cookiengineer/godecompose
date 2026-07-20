# Roadmap

## Phase 1: Decompiler Rewrite — Correct Go Source Reconstruction ✓

**Completed.** Goes beyond pattern matching to produce legitimate Go syntax: structured control flow with inline else bodies, Go expressions via data flow analysis (instead of assembly comments), for-loops, if-else-if chains, and semantically-named struct fields.

See `docs/phase1-decompiler-rewrite.md` for the full implementation plan.

### Deliverables

| Sub-phase | Description |
|---|---|
| 1a — Data Flow Analysis | New `analysis/dfa/` package (5 files, 23 unit tests): Plan9 operand parser, instruction→Value translation, copy propagation, dead store elimination, constant folding, Go text emission |
| 1b — DFA Integration | Replaced `emitFunctionRange` and `emitRawRange` in `pattern/generate/generate.go` with DFA-based Go expression emission; asm fallback only for unrecognized operations |
| 1c — Control Flow Structuring | Post-dominator computation, if/else diamond with inline else bodies, for-loop detection, if-else-if chains, compound condition extraction, jump table detection |
| 1d — Struct Field Naming | Method name analysis (GetX→x, SetY→y, IsZ→bool hint), cross-method type consensus (MOVB/CMPB→bool, convT→string/int64), offset-based Go convention heuristics |
| 1e — Integration & Tests | 8 new E2E suites: dfa_simple, dfa_pointers, phase1_forloop, phase1_ifelse, phase1_switch, phase1_structs, phase1_project, phase2_quality. Unit tests: structure_test.go (4), signature_test.go (15 new). Recovery metrics with DFA expression count |

### New Packages and Files

```
analysis/dfa/value.go          — Value, ValueKind, MemRef, BlockState, Statement types
analysis/dfa/parse.go          — Plan9 operand parser (registers, immediates, memory refs, symbols)
analysis/dfa/translate.go      — Instruction→Value translation (MOV, ADD, SUB, AND, OR, XOR, SHL, SHR, LEA, CMP, CALL, RET)
analysis/dfa/optimize.go       — Copy propagation, dead store elimination, constant folding
analysis/dfa/emit.go           — Value graph→Go text emission
analysis/dfa/dfa_test.go       — 23 unit tests
docs/phase1-decompiler-rewrite.md  — Full implementation plan and design
```

### Key API Additions

| Package | Addition |
|---|---|
| `disasm/structure.go` | `CFGRegion`, `RegionKind`, `PdomTree` field, `computePostDominators()`, `postDominates()` |
| `function/signature.go` | `extractMethodNameHints()` — GetX/SetX/IsX→field name |
| `pattern/generate/generate.go` | `emitLoop()`, `emitIfElseIfChain()`, `emitCondBlock()` (rewritten), `emitPlainBlock()`, `findMergeBlock()`, `isSimpleBlockContent()`, `addrLookup()` |

### Results

| Metric | Before (Phase 0) | After (Phase 1) |
|---|---|---|
| Assembly comment lines | ~10,000 `// 0x...: mov rax, rbx` | DFA Go expressions; asm fallback for unrecognized ops only |
| Control flow | Flat `L0:`, `goto L3` | Nested `if`/`else` with inline bodies; `for` loops; `else if` chains |
| Else bodies | `} else { goto L4 }` | Inline else code between `} else {` and `}` |
| Struct fields | `field_0x50` | Method-name-derived: `GetX`→`x`, `SetY`→`y` |
| Noise filtering | NOP/INT3 in output | NOP/INT3/DATA16 filtered everywhere |
| E2E suites | 44 | 52 |
| Test packages passing | All | All (68 packages, 0 failures, `go vet` clean) |
| Unit tests added | 16 (Phase 10) | 42 new (23 DFA + 4 structure + 15 signature) |

---

## Historical Phases

### Phase 10 — Pattern Development & Infrastructure ✓

709 patterns, 3,141 matches on ysco benchmark (80.9% call-site recovery). 44 E2E suites (37 stdlib + 7 Phase 10) + 16 unit tests. Pattern discovery tool (`godecompose patterns discover`), recovery rate metrics. CFG-based structured output, function signatures, struct fields.

### Phase 8-9 — CLI, Actions, User-Code Focus ✓

CLI commands (`info`, `disasm`, `decompile`, `patterns`). Reusable `actions/` package. Function classification (User/Runtime/Stdlib). Callgraph refinement. Struct tracking. Project directory generation via `WriteProject`.

### Phase 6-7 — Pattern Matcher, Code Generator, Database ✓

Opcode-indexed pattern matching, gen template expansion. Four kernel syscall tables (Linux 137, Windows 121, Darwin 70, FreeBSD 57). Pattern database with recursive `.hexpat` loading.

### Phase 4-5 — Pattern Language Engine ✓

ImHex-compatible pattern language with godecompose extensions (`instr`, `gen`, `bind`). Lexer (50+ token types), recursive-descent parser (13-level precedence), validator, evaluator. Preprocessor (`#include`, `#define`, `#ifdef`).

### Phase 1-3 — Foundation ✓

Binary format parsers (ELF, PE, Mach-O) via `binary.Open()`. Go build info + pclntab extraction. x86_64 disassembler (`golang.org/x/arch/x86/x86asm`) with Plan9 syntax. CFG construction. Function recovery via pclntab + symbol tables.