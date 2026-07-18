# Pattern Language Specification

Godecompose implements an **ImHex-compatible** pattern language, extended with constructs for assembly pattern matching and source code generation.

## Relationship to ImHex PatternLanguage

The lexer, parser, and AST are designed to be fully compatible with the [ImHex PatternLanguage](https://github.com/WerWolv/PatternLanguage) specification for the base language features. Pattern files written for ImHex's binary data parsing _should_ parse correctly in godecompose.

Extensions (described below) use additional keywords and block types that would not be valid in upstream ImHex.

## Base Language (ImHex-Compatible)

### Types

All standard ImHex integer types:

| Type     | Width    | Signed |
|----------|----------|--------|
| `u8`     | 8-bit    | No     |
| `u16`    | 16-bit   | No     |
| `u24`    | 24-bit   | No     |
| `u32`    | 32-bit   | No     |
| `u48`    | 48-bit   | No     |
| `u64`    | 64-bit   | No     |
| `u96`    | 96-bit   | No     |
| `u128`   | 128-bit  | No     |
| `s8`     | 8-bit    | Yes    |
| `s16`    | 16-bit   | Yes    |
| `s24`    | 24-bit   | Yes    |
| `s32`    | 32-bit   | Yes    |
| `s48`    | 48-bit   | Yes    |
| `s64`    | 64-bit   | Yes    |
| `s96`    | 96-bit   | Yes    |
| `s128`   | 128-bit  | Yes    |
| `char`   | 8-bit    | —      |
| `char16` | 16-bit   | —      |
| `bool`   | 1-bit    | —      |
| `float`  | 32-bit   | —      |
| `double` | 64-bit   | —      |
| `str`    | variable | —      |
| `auto`   | inferred | —      |
| `padding`| N bytes  | —      |

Compound types: `struct`, `union`, `enum`, `bitfield`.

### Structs

```rust
struct MyStruct {
    u32 magic;
    u32 version;
    padding[16];
    double timestamp;
};

MyStruct header @ 0x00;
```

Structs support inheritance and attributes:

```rust
struct Named [[name("Display Name"), color("FF0000")]] {
    u32 id;
};

struct Extended : Named {
    u64 extra;
};
```

### Enums

```rust
enum FileType : u16 {
    ELF = 0x457F,
    PE  = 0x5A4D,
    MachO_32 = 0xFEEDFACE,
    MachO_64 = 0xFEEDFACF
};
```

### Unions

```rust
union FloatOrInt {
    float as_float;
    u32 as_int;
};
```

### Bitfields

```rust
bitfield Flags {
    read    : 1;
    write   : 1;
    execute : 1;
    reserved : 5;
};
```

### Variables

Placed variables (at a specific offset in the binary):

```rust
u32 magic @ 0x00;
u8 data[256] @ 0x100;
u32 *count_ptr @ 0x04;
```

Local variables (in heap memory, not tied to binary offset):

```rust
u32 counter;
str name = "unknown";
```

Discard variable: `_` for values you don't need to bind.

### Functions

```rust
fn read_string(u32 offset, u32 max_len) {
    u8 bytes[max_len] @ offset;
    str result;
    // ... build string from bytes
    return result;
};

fn main() {
    std::print("Parsing: {}", read_string(0x10, 256));
};
```

### Control Flow

`if`/`else`:

```rust
if (magic == 0x50) {
    std::print("Type A");
} else if (magic == 0x60) {
    std::print("Type B");
} else {
    std::print("Unknown");
}
```

`while` loops:

```rust
u32 i = 0;
while (i < count) {
    u32 entry @ $;
    i = i + 1;
}
```

`for` loops (C-style):

```rust
for (u32 i = 0, i < count, i = i + 1) {
    u32 entry @ $;
}
```

`match` statement:

```rust
match (value) {
    (0x50): std::print("Type A");
    (0x60 ... 0x6F): std::print("Range");
    (0x80 | 0x81): std::print("Alternatives");
    (_): std::print("Default");
};
```

### Operators

Arithmetic: `+`, `-`, `*`, `/`, `%`
Bitwise: `&`, `|`, `^`, `~`, `<<`, `>>`
Logical: `&&`, `||`, `^^`, `!`
Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
Assignment: `=`, `+=`, `-=`, `*=`, `/=`, `%=`, `&=`, `|=`, `^=`, `<<=`, `>>=`
Ternary: `? :`
Scope resolution: `::`
Placement: `@` (place at offset)
Current offset: `$`
Type operators: `sizeof`, `addressof`, `typenameof`
Cast: `type(value)` or `value as type`

### Endianness

```rust
little_endian;  // default
big_endian;
```

### Namespaces

```rust
namespace my_format {
    struct Header { u32 magic; };
};

my_format::Header h @ 0x00;
```

### Imports

```rust
import std::string;
import "other_file.hexpat";
```

### Attributes

```rust
struct MyType [[name("Display Name"), color("FF0000"), hidden]] {
    u32 value [[format_read("hex")]];
};
```

---

## Godecompose Extensions

These constructs extend the ImHex PatternLanguage for decompilation use cases.

### Architecture and Platform Declarations

At the top of a pattern file:

```rust
arch x86_64;
platform linux, darwin, freebsd;
```

Or with version ranges:

```rust
arch x86_64;
platform windows(10.0, 11.0);
```

These act as filters: the pattern matcher only considers patterns matching the target binary's arch/platform.

### `instr` Block — Assembly Pattern Matching

An `instr` block describes a sequence of assembly instructions to match against the disassembled instruction stream.

```rust
instr match_write_syscall {
    MOVQ $1, RAX               // Linux x86_64: syscall number for write
    MOVQ fd_reg, RDI           // $fd_reg captures any register used for fd
    MOVQ buf_reg, RSI
    MOVQ count_reg, RDX
    SYSCALL
}
```

**Syntax rules for `instr` blocks:**

1. Each line is `<OPCODE> [operands...]`
2. Opcodes use Go Plan 9 syntax: `MOVQ`, `ADDQ`, `CALL`, `LEAQ`, `JMP`, `JE`, `JNE`, etc.
3. Operands can be:
   - **Literal values**: `$1`, `$0x3F`, `$42` (must match exactly)
   - **Named registers**: `RAX`, `RDI`, `RSI`, `R8-R15`, `X0-X15`
   - **Capture variables**: identifiers not matching known register names → bind to matched value
   - **Wildcard**: `_` matches any operand (including none)
   - **Addressing modes**: `offset(base)`, `(base)(index*scale)`, `symbol(SB)`, `name+offset(FP)`
4. Labels use `@labelname:` (Plan 9 style)
5. `;` separates instructions on the same line (for REP prefix etc.)
6. `//` for line comments
7. Multiple alternative matches: use `|` between instructions:

```rust
instr match_any_64bit_move {
    MOVQ src, dst       // Matches register-to-register
    | MOVQ $imm, dst    // Matches immediate-to-register
    | MOVQ (src), dst   // Matches memory-to-register
}
```

### Capture Variables in `instr` Blocks

Identifiers in operand position that are not known register names become capture variables. On a successful match, the captured value is bound to the variable.

```rust
instr match_add {
    ADDQ $amount, target_reg
    // $amount captures the immediate value
    // target_reg captures which register was modified
}
```

Capture variables can also be used in addressing modes:

```rust
instr match_stack_access {
    MOVQ value_reg, offset(SP)
    // value_reg: captures source register
    // offset: captures stack offset value
}
```

### `gen` Block — Source Code Generation

A `gen` block defines what source code to emit when the matching `instr` block succeeds.

```rust
gen write_source {
    write($fd_reg, $buf_reg, $count_reg);
}
```

**Syntax rules for `gen` blocks:**

1. Template text with `$variable` placeholders that are substituted with bound values
2. Multiline templates
3. Expression substitution: `${expression}` evaluates an expression and inserts its string form
4. Conditional generation:

```rust
gen {
    if ($platform == "linux") {
        syscall(SYS_write, $fd_reg, $buf_reg, $count_reg);
    } else {
        write($fd_reg, $buf_reg, $count_reg);
    }
}
```

5. Loop generation (for repeated structures):

```rust
gen {
    for (i := 0; i < $count; i++) {
        array[$i] = ${read_u32($base + i * 4)};
    }
}
```

### `bind` Block — Variable Naming

```rust
bind {
    fd_reg as "fd"
    buf_reg as "buf"
    count_reg as "count"
    value_reg as "result"
    target_reg as "counter"
}
```

This maps captured register names to human-readable variable names used in generated source code.

### Complete Pattern Example

A full decompilation pattern for a function:

```rust
// patterns/libs/golang/runtime/memmove.hexpat

arch x86_64;
platform linux, darwin, windows;

import std::mem;

pattern go_runtime_memmove {
    name: "runtime.memmove";
    library: "go-runtime";
    version: ">=1.0";
    description: "Go runtime memmove - copies n bytes from src to dst";

    instr match_forward_copy {
        // Detect direction by comparing src and dst
        CMPQ src_ptr, dst_ptr
        JAE @forward

        // Backward copy: STD; REP; MOVSB
        STD
        MOVQ src_ptr, SI
        ADDQ len_minus_1, SI
        MOVQ dst_ptr, DI
        ADDQ len_minus_1, DI
        MOVQ n_reg, CX
        REP; MOVSB
        CLD
        JMP @done

    @forward:
        // Forward copy: REP; MOVSQ + MOVSB for remainder
        MOVQ src_ptr, SI
        MOVQ dst_ptr, DI
        MOVQ n_reg, CX
        SHRQ $3, CX
        REP; MOVSQ
        MOVQ n_reg, CX
        ANDQ $7, CX
        REP; MOVSB

    @done:
        RET
    }

    gen {
        memmove($dst_ptr, $src_ptr, $n_reg);
    }

    bind {
        src_ptr as "src"
        dst_ptr as "dst"
        n_reg as "n"
        len_minus_1 as "n"
    }
}
```

## Compilation and Execution Model

The pattern language follows ImHex's pipeline:

1. **Preprocessor**: Resolves `#include`, `#define`, `#pragma`, `import` statements.
2. **Lexer**: Tokenizes the source into a `[]Token` stream.
3. **Parser**: Recursive descent parser produces an AST from tokens. Uses backtracking for ambiguous constructs.
4. **Validator**: Type-checks the AST, resolves identifiers, reports semantic errors.
5. **Evaluator**: Tree-walking interpreter. For godecompose, evaluator mode depends on usage:
   - **Binary data parsing mode**: Like ImHex — evaluates patterns against raw bytes, creates Pattern objects at offsets.
   - **Assembly matching mode**: Evaluates `instr` blocks against an instruction stream, produces variable bindings.
   - **Code generation mode**: Evaluates `gen` blocks with bound variables, produces source code strings.

## Pattern File Organization

Pattern files use the `.hexpat` extension and are organized by target:

```
patterns/
├── kernels/                  # OS-level patterns
│   ├── linux/
│   │   ├── syscall_table.json     # Syscall number → name mapping
│   │   ├── syscall_write.hexpat   # Pattern for write() syscall
│   │   ├── syscall_read.hexpat
│   │   └── ...
│   ├── windows/
│   │   ├── syscall_nt10.json      # Windows 10 syscall table
│   │   ├── syscall_nt11.json      # Windows 11 syscall table
│   │   └── ntdll_*.hexpat
│   ├── darwin/
│   │   └── syscall_table.json
│   └── freebsd/
│       └── syscall_table.json
│
└── libs/                     # Library-specific patterns
    ├── golang/
    │   ├── runtime/               # Go runtime function patterns
    │   │   ├── memmove.hexpat
    │   │   ├── newobject.hexpat
    │   │   ├── typedmemmove.hexpat
    │   │   ├── chan_send.hexpat
    │   │   ├── chan_recv.hexpat
    │   │   ├── mapaccess.hexpat
    │   │   └── deferreturn.hexpat
    │   └── stdlib/                # Go standard library patterns
    │       ├── fmt_sprintf.hexpat
    │       ├── io_read.hexpat
    │       ├── sync_mutex.hexpat
    │       └── net_dial.hexpat
    └── openssl/
        ├── aes_encrypt.hexpat
        ├── aes_decrypt.hexpat
        ├── sha256.hexpat
        └── rand_bytes.hexpat
```

Syscall tables use JSON for easy editing and contribution:

```json
{
  "platform": "linux",
  "arch": "x86_64",
  "version": "5.0",
  "syscalls": {
    "0":  {"name": "read",    "args": ["unsigned int fd", "char *buf", "size_t count"]},
    "1":  {"name": "write",   "args": ["unsigned int fd", "const char *buf", "size_t count"]},
    "2":  {"name": "open",    "args": ["const char *pathname", "int flags", "mode_t mode"]},
    "3":  {"name": "close",   "args": ["unsigned int fd"]},
    "59": {"name": "execve",  "args": ["const char *pathname", "char *const argv[]", "char *const envp[]"]},
    "60": {"name": "exit",    "args": ["int error_code"]},
    "231": {"name": "exit_group", "args": ["int error_code"]}
  }
}
```
