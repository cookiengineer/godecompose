# Phase 10: User-Code Control Flow & Data Type Reconstruction

**Status: Completed** (668 patterns, 2,209 matches on ysco -- 316x improvement from original 7)

## Why This Phase Is Necessary

The original 550 patterns (stdlib + runtime + fallback) matched **only CALL-based linear
sequences**. In the ysco decompilation test, this produced a ~0.03% source code
recovery rate -- just 7 trivial one-liners across 44 decompiled functions. The
remaining ~5700 lines of assembly are `// unresolved`.

The core gap: **there are no patterns for Go constructs that map to non-CALL assembly
sequences**. We need patterns that understand:

1. **Control flow** -- `CMP`/`TEST` + `Jcc` sequences map to `if`/`else`/`for`
2. **Variable assignment** -- `MOVQ`/`LEAQ` between registers or from stack
3. **Data type creation** -- composite literals, string construction, slice/map creation
4. **Go idioms** -- `if err != nil { return err }`, `v, ok := m[key]`, defer patterns
5. **Struct field access** -- offset-based field reads/writes
6. **Method dispatch** -- virtual calls through itab

The pattern language fully supports multi-instruction matching, `gen` block
`if`/`for` conditionals, and pipe `|` alternatives -- but **none of this is used**
in existing patterns.

---

## Task 1: Fix GenConditional Evaluator Bug [BLOCKING]

**File**: `pattern/lang/evaluator/evaluator.go:312-323`

**Bug**: `evalGenStmt` for `GenConditional` uses `return` inside the `for` loop body,
so only the **first** gen statement in each conditional branch is processed. Multi-statement
`if`/`else` gen bodies silently drop all but the first statement.

**Fix**: Change the `for`/`return` loop to accumulate into a `strings.Builder`.

---

## Task 2: Implement GenLoop Evaluation

**File**: `pattern/lang/evaluator/evaluator.go:324`

**Status**: GenLoop AST nodes are parsed but silently dropped (falls through to
`return ""`).

**Required**: Add `GenLoop` handling to `evalGenStmt`. The simplest useful form
evaluates a loop condition from bindings (a captured `$count` variable) and repeats
the body `n` times.

---

## Task 3: Control Flow Patterns

### 3a. if/else -- Error Check

**Go source**:
```go
if err != nil {
    log.Fatal("error")
}
```

**Go assembly (x86_64, ABIInternal)**:
```
MOVQ err_reg, AX
TESTQ AX, AX          ; check nil
JEQ  @after
LEAQ str_ptr(SB), AX
MOVQ $len, BX
XORL CX, CX           ; args = nil
CALL log.Fatal(SB)
@after:
```

**Pattern** (`patterns/golang/controlflow/if_err_nil_fatal.hexpat`):
```rust
arch x86_64;
platform linux, darwin, windows, freebsd;
pattern go_if_err_nil_fatal {
    name: "go: if err != nil { log.Fatal }";
    library: "go-controlflow";
    description: "if err != nil { log.Fatal } pattern";
    instr match {
        MOVQ err, tmp
        TESTQ tmp, tmp
        JEQ @after
        LEAQ format, AX
        MOVQ len, BX
        CALL log.Fatal
    @after:
    }
    gen {
        if err != nil {
            log.Fatal($err);
        }
    }
    bind {
        tmp as "err"
        err as "err"
        format as "format"
        len as "len"
    }
}
```

### 3b. if err != nil { return err }

```
TESTQ err_reg, err_reg       ; check nil
JEQ @noerr
MOVQ err_reg, AX             ; return err as first return
MOVQ err_type_reg, BX        ; return error type
RET
@noerr:
```

Pattern: `if $err != nil { return $err }`.

### 3c. if v == const / if v != const

```
CMPQ reg, $imm
JEQ @equal    ; or JNE @notequal
```

### 3d. if/else with CMP comparison

```
CMPQ reg1, reg2
JLT @else_block    ; or JGT / JLE / JGE
; if body
JMP @after
@else_block:
; else body
@after:
```

### 3e. for loop (counted)

```
XORL i, i             ; i := 0
JMP @check
@body:
; loop body
INCQ i
@check:
CMPQ i, n
JLT @body
```

---

## Task 4: Variable Assignment & Data Movement Patterns

| Pattern | Assembly Signature | Gen Output |
|---------|-------------------|------------|
| `int literal` | `MOVQ $imm, reg` | `x := $imm` |
| `string literal` | `LEAQ str_data(SB), AX; MOVQ $len, BX` | `x := "$str"` |
| `zero value` | `XORL AX, AX` | `x := 0` |
| `bool true` | `MOVB $1, AL` | `x := true` |
| `bool false` | `XORL AL, AL` | `x := false` |
| `copy variable` | `MOVQ reg_src, reg_dst` | `dest = src` |
| `address of` | `LEAQ (SP) offset, AX` | `x := &local` |
| `struct field write` | `MOVQ val, offset(receiver)` | `receiver.field = val` |
| `struct field read` | `MOVQ offset(receiver), val` | `val := receiver.field` |
| `slice element read` | `MOVQ base(reg)(idx*8), val` | `val := slice[idx]` |

---

## Task 5: Data Type Creation Patterns

| Pattern | Assembly Signature | Gen Output |
|---------|-------------------|------------|
| `string from data + len` | `LEAQ data(SB), AX; MOVQ $len, BX` | `x := "$data"` |
| `make slice` | `MOVQ len, AX; MOVQ cap, BX; CALL runtime.makeslice` | `x := make([]T, $len, $cap)` |
| `slice literal` | `runtime.newobject + MOVQ sequence + CALL runtime.growslice` | `x := []T{...}` |
| `make map` | `MOVQ hint, BX; CALL runtime.makemap` | `x := make(map[K]V)` |
| `map literal` | `CALL runtime.makemap_small + CALL runtime.mapassign_faststr` | `x := map[K]V{...}` |
| `struct literal (stack)` | `MOVQ field1, offset(SP); MOVQ field2, offset+8(SP)` | `x := Type{...}` |
| `new(struct)` | `LEAQ type(SB), AX; CALL runtime.newobject` | `x := new(Type)` |
| `interface{} from value` | `CALL runtime.convTstring` / `convT64` / `convT32` | `interface{}(x)` |

---

## Task 6: Go Idiom Patterns

### 6a. `val, err := fn(); if err != nil { return err }`

```
CALL someFunc(SB)
TESTQ BX, BX            ; BX = second return (err)
JEQ @noerr
MOVQ AX, RET1           ; return err
MOVQ BX, RET2
RET
@noerr:
```

### 6b. `v, ok := m[key]` (map access with comma-ok)

```
CALL runtime.mapaccess2_faststr(SB)  ; returns (val_ptr, ok) in AX/BX
TESTB BL, BL                          ; check ok
```

### 6c. `defer fn(args)`

```
LEAQ deferred_fn(SB), AX
CALL runtime.deferproc(SB)
```

### 6d. `go fn(args)` (goroutine)

```
LEAQ goroutine_fn(SB), AX
CALL runtime.newproc(SB)
```

---

## Task 7: Struct Field Access Patterns

Recognize field access by consistent offsets from a struct pointer receiver.

**Strategy**: Since field offsets depend on struct layout, create dynamic offset-based
patterns that capture the offset value and emit symbolic names:

```
MOVQ val, 0x69(receiver_reg)    →  receiver.unknown_field_0x69 = val
MOVQ 0x28(receiver_reg), BX     →  v := receiver.unknown_field_0x28
```

Multi-function analysis can infer field count and sizes by comparing offsets used
across methods of the same struct.

---

## Task 8: Enhanced CALL Pattern Matching

### 8a. Interface method calls (itab dispatch)

```
LEAQ go:itab.*Type,Interface(SB), AX
MOVQ qword ptr [AX+0x18], AX     ; load method from itab
CALL AX                           ; indirect call
```

### 8b. Closure calls (via DX / REGCTXT)

```
CALL DX    ; indirect call through closure context register
```

### 8c. Method calls on structs

Recognize `CALL pkg.(*Type).Method(SB)` with argument setup as method call.

---

## Task 9: E2E Verification

For each pattern category, create test programs:

```
patterns/golang/controlflow/if_err_nil_fatal.hexpat
testdata/src/controlflow/main.go
e2e/decompile/controlflow/controlflow_test.go
```

---

## Task 10: Pattern Discovery Automation

Tool that:
1. Takes a Go source file as input
2. Compiles it with `go tool compile -S` to get assembly
3. Extracts assembly sequences between known function boundaries
4. Generates candidate `.hexpat` pattern files from assembly
5. Flags sequences that don't fit existing patterns as "need pattern"

Command: `godecompose patterns discover <source.go>`

---

## Summary of All Pattern Categories

| Priority | Category | Estimated Patterns | Impact |
|----------|----------|-------------------|--------|
| **P0** | Fix GenConditional bug | 0 (code fix) | Unblocks multi-statement gen |
| **P0** | if/else control flow | ~12 | Recovers 80% of unresolved blocks |
| **P0** | Error handling idioms | ~8 | Recovers `if err != nil` everywhere |
| **P1** | Variable assignment (MOV/LEA) | ~15 | Recovers local variable declarations |
| **P1** | String creation | ~5 | Recovers string literals |
| **P1** | Slice/map creation | ~8 | Recovers composite literals |
| **P1** | Struct field access | Dynamic | Recovers struct member reads/writes |
| **P1** | For loop patterns | ~6 | Recovers iteration constructs |
| **P2** | Switch/case patterns | ~4 | Recovers switch statements |
| **P2** | Interface method dispatch | ~3 | Recovers virtual method calls |
| **P2** | Closure/goroutine patterns | ~4 | Recovers `go fn()` and closures |
| **P2** | Defer patterns | ~2 | Recovers defer statements |
| **P3** | Pattern discovery automation | Tool | Accelerates future pattern creation |

---

## Implementation Order

1. Fix `evalGenStmt` for `GenConditional`
2. Implement `GenLoop` evaluation
3. Write if/else control flow patterns -- highest impact
4. Write error handling idiom patterns -- every Go function has `if err != nil`
5. Write variable assignment patterns
6. Write data type creation patterns
7. Write method dispatch/closure patterns
8. Build pattern discovery tool

---

## Completed (as of Phase 10 iteration 1)

### Code Fixes (8 files)

| File | Change |
|------|--------|
| `disasm/disasm.go` | Extract opcodes from GoSyntax (Plan 9); normalize TEST/CMP across size variants; condition code mapping (JE→JEQ, JG→JGT, etc.) |
| `pattern/lang/evaluator/evaluator.go` | Fixed GenConditional accumulator bug; added GenLoop evaluation; added whitespace insertion between gen statements |
| `pattern/lang/parser/parser.go` | Balanced brace tracking in gen blocks; `@if`/`@for` for compile-time gen conditionals; plain `if`/`for` as gen text |
| `pattern/matcher/matcher.go` | O(n×10k)→O(n×m) conflict resolution; fixed fuzzy matcher to split on underscores (`_`); reduced max raw matches to 10000 |
| `actions/decompile.go` | O(n²)→O(n) instruction collection via address map; progress logging |
| `function/pclntab.go` | `extractFunctionBlocks` binary search instead of full scan |
| `patterns/golang/embed.go` | Added `controlflow/**` embed + `LoadControlFlow()` |
| `cmd/godecompose/main.go` | Load all 4 pattern modules (stdlib/runtime/fallback/controlflow) |

### New Pattern Files (6 files, 93 patterns)

| File | Patterns | Type |
|------|----------|------|
| `controlflow/if_else.hexpat` | 9 | Control flow: if/else, nil checks, CMP comparisons |
| `controlflow/go_idioms.hexpat` | 12 | Go idioms: defer, goroutines, channels, maps |
| `controlflow/data_types.hexpat` | 12 | Data types: new(T), make(), type conversions |
| `controlflow/calls_patterns.hexpat` | 20 | Specific calls: slog, time, os, syscall, mutex |
| `controlflow/runtime_extras.hexpat` | 22 | Runtime: morestack, panics, concatstrings, assertions |
| `controlflow/more_calls.hexpat` | 24 | Stdlib: flag, net/http, fmt, sconf, prometheus |
| `controlflow/remaining_calls.hexpat` | 18 | Remaining stdlib: os, net, syscall, strings, strconv |

### Key Bugs Fixed

1. **Plan 9 opcode mismatch**: Disassembler used Intel opcodes (JE) while patterns used Go Plan 9 (JEQ). Fixed by extracting opcodes from GoSyntax with condition code normalization.
2. **Fuzzy matcher underscore bug**: Runtime function names use underscores (`mapassign_faststr`) but fuzzy matcher only split on dots/parens. Fixed by adding `_` to the normalizer.
3. **GenConditional accumulator**: Only first statement of each conditional branch was returned. Fixed to accumulate all.
4. **Gen block brace nesting**: `{`/`}` inside gen blocks closed the gen block prematurely. Fixed with balanced depth tracking.
5. **Inline `//` comments in gen blocks**: Lexer strips `//` to end-of-line, eating the closing `}`. Fixed by using `/* */` style comments or removing them.
6. **O(n²) performance**: Three separate O(n²) loops (instruction collection, function block extraction, conflict resolution). All fixed to O(n) or O(n log n).

### Results (ysco decompilation)

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Patterns loaded | 550 | 668 | +118 |
| Matches | 7 | 2,209 | **316x** |
| Unresolved CALLs | 769 | 0 (all remaining are user functions) | All stdlib/runtime matched |
| Recovery quality | ~0.03% | ~40% of call sites identified | High-level trace recovered |

### Remaining Work

| Priority | Category | Status |
|----------|----------|--------|
| **P1** | Struct field access via offset inference | Not started |
| **P1** | For loop patterns (counted iteration) | Not started |
| **P2** | Switch/case patterns | Not started |
| **P2** | Interface method dispatch (itab) | Not started |
| **P2** | Closure call detection (via DX) | Not started |
| **P3** | Pattern discovery automation (`godecompose patterns discover`) | Not started |
| **P3** | E2E tests for controlflow patterns | Not started |
