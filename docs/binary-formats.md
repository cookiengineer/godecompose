# Binary Format Parsing

## Overview

The `binary` package provides a unified interface over ELF, PE (COFF), and Mach-O binary formats. Each format sub-package wraps the Go standard library's parsers where available and supplements them with additional metadata extraction needed by the disassembler and function recovery stages.

## Common Interface

### Binary

```go
type Format int

const (
    ELF   Format = iota
    PE
    MachO
)

type Binary interface {
    Format() Format
    Architecture() Arch          // x86_64 (future: arm64, riscv64)
    EntryPoint() uint64
    Sections() []Section
    Segment(name string) (Section, bool)
    Symbols() ([]Symbol, error)
    DynamicSymbols() ([]Symbol, error)
    IsPIE() bool                // Position-Independent Executable
    IsStripped() bool
    ByteOrder() binary.ByteOrder
    GoBuildInfo() (*GoBuildInfo, error)  // Go-specific: from .go.buildinfo note
    Pclntab() ([]byte, uint64, error)    // Go-specific: pclntab section data + base
    Close() error
}
```

### Section

```go
type SectionFlag int

const (
    SectionExecutable SectionFlag = 1 << iota
    SectionWritable
    SectionReadable
)

type Section struct {
    Name    string
    Address uint64         // Virtual address (where it loads in memory)
    Size    uint64          // Size in memory
    Offset  uint64          // File offset
    Data    []byte          // Section contents (may be empty for BSS)
    Flags   SectionFlag
    EntSize uint64          // Entry size (for symbol tables, etc.)
    Link    uint32           // Section link index
    Info    uint32           // Section info
    Align   uint64          // Alignment
}
```

### Symbol

```go
type SymbolType int

const (
    SymbolFunction SymbolType = iota
    SymbolObject
    SymbolSection
    SymbolFile
    SymbolUnknown
)

type Symbol struct {
    Name    string
    Address uint64
    Size    uint64
    Type    SymbolType
    Section string         // Section this symbol is defined in
    Version string         // Symbol version (@GLIBC_2.2.5 etc.)
    Binding SymbolBinding  // Local, Global, Weak
}
```

### GoBuildInfo

```go
type GoBuildInfo struct {
    Version    string      // Go compiler version (e.g., "go1.22.0")
    Path       string      // Module path
    Main       string      // Main package path
    Deps       []GoModuleDep
    Settings   map[string]string  // -gcflags, -ldflags, etc.
}

type GoModuleDep struct {
    Path    string
    Version string
    Sum     string
}
```

---

## ELF Format (`binary/elf`)

### Go Standard Library Usage

The `debug/elf` package handles:
- Header parsing (class, endianness, machine, entry point)
- Program headers (LOAD segments for runtime layout)
- Section headers (all sections)
- Symbol tables (.symtab, .dynsym)
- String tables (.strtab, .dynstr)
- Dynamic section (.dynamic)

### What We Add

**Go-specific ELF notes** (`.note.go.buildinfo`, `.note.gobuildid`):

Go binaries embed build metadata in ELF notes. The `.go.buildinfo` note contains:
- Go compiler version
- Module path and version
- Build settings (-gcflags, -ldflags, -tags, etc.)
- Go module dependencies (path, version, hash)

This data is encoded as a Go-specific format within the ELF note. Parsing logic:

1. Find the ELF note section with type `GoBuildID` (4) or name `.go.buildinfo`
2. Parse the build info structure (versioned, currently V1)
3. Extract the module graph and build settings

**pclntab detection**:

The pclntab (PC-line table) is stored in its own section (`.gopclntab` or `.data.rel.ro.gopclntab`) in modern Go versions, or at the start of `.text` in older versions (Go <1.2).

Detection strategy:
1. Look for `.gopclntab` section by name
2. If not found, search for pclntab magic bytes in `.text` or read-only data sections
3. Verify against pclntab header format

### Section Name Conventions for Go ELF

```
.text              — Executable code
.rodata            — Read-only data (strings, constants)
.typelink          — Go type descriptors
.itablink          — Go interface table links
.gosymtab          — Go symbol table (deprecated, now merged into pclntab)
.gopclntab         — Go PC-line table (function metadata)
.noptrdata         — Non-pointer data
.data              — Mutable data (with pointers)
.bss               — Zero-initialized data
.noptrbss          — Zero-initialized non-pointer data
.go.buildinfo      — Go build info ELF note
.note.go.buildid   — Go build ID ELF note
.debug_*           — DWARF debug info (present in most Go binaries)
```

### Implementation Notes

- Use `debug/elf.Open()` to open files
- `debug/elf.File.Sections` provides all sections
- `debug/elf.File.Symbols()` provides static symbols
- `debug/elf.File.DynamicSymbols()` provides dynamic symbols
- `debug/elf.File.ImportedSymbols()` provides symbols from shared libraries
- For PIE binaries: sections have zero `Addr` in the file — the runtime base must be determined separately
- Go's internal linking produces ELF files with `.go.buildinfo` note; external linking may use C linker conventions

---

## PE Format (`binary/pe`)

### Go Standard Library Usage

The `debug/pe` package handles:
- DOS header + PE signature
- COFF file header (machine, number of sections, timestamp)
- Optional header (entry point, image base, section alignment)
- Section headers
- COFF symbol table
- Import directory (IAT / INT)
- Export directory (EAT)

### What We Add

**Export directory parsing** (if not fully covered by `debug/pe`):

The export directory provides function names and ordinals for DLL exports. For Go binaries, this contains the set of exported functions from the main binary or any Go DLLs.

**Go build info in PE**:

Go binaries on Windows embed build information differently from ELF. Go 1.13+ stores build info in a PE section named `.go.buildinfo` or in a resource. We need to:

1. Search for `.go.buildinfo` section
2. Parse the build info V1 format
3. Extract module dependencies and build settings

**pclntab detection in PE**:

Similar to ELF: look for `.gopclntab` section, or search for magic bytes.

**PE-specific details**:

- Sections use `VirtualAddress` (RVA, relative to image base) — need to track image base
- Import/export tables use RVAs
- Relocations stored in `.reloc` section (base relocations, needed for ASLR)
- Go PE binaries typically use `IMAGE_FILE_MACHINE_AMD64` (0x8664)

### Section Name Conventions for Go PE

```
.text    — Executable code
.rdata   — Read-only data (includes Go type data, pclntab in some versions)
.data    — Mutable data
.bss     — Zero-initialized
.idata   — Import data
```

### Implementation Notes

- Use `debug/pe.Open()` to open files
- `debug/pe.File.Sections` provides all sections
- `debug/pe.File.Symbols` provides COFF symbols (may be absent in Go-linked PEs)
- `debug/pe.File.ImportedSymbols()` provides import table entries
- Go PE binaries may have unusual section alignment (minimal, as Go linker optimizes for size)

---

## Mach-O Format (`binary/macho`)

### Go Standard Library Usage

Go's standard library has limited Mach-O support via `debug/macho`. However, this package is in the `x/` repositories and has been removed from newer Go versions. For this project, we implement our own Mach-O parser based on the Apple specification.

### Mach-O Binary Structure

```
┌─────────────────────┐
│   Mach-O Header      │  magic, cputype, cpusubtype, filetype, ncmds, sizeofcmds
├─────────────────────┤
│   Load Commands      │  Array of load commands (segment, symtab, dysymtab, ...)
│   - LC_SEGMENT_64    │    Segment definitions (name, vmaddr, vmsize, offset)
│   - LC_SYMTAB        │    Symbol table location (symoff, nsyms, stroff, strsize)
│   - LC_DYSYMTAB      │    Dynamic symbol table indices
│   - LC_LOAD_DYLIB    │    Dynamic library dependencies
│   - LC_MAIN          │    Entry point (replaces LC_UNIXTHREAD on modern macOS)
│   - LC_UUID          │    Binary UUID
│   - LC_VERSION_MIN_*  │    Minimum OS version
│   - ...              │
├─────────────────────┤
│   Segment Data       │  Raw bytes for each segment
├─────────────────────┤
│   Symbol Table       │  nlist_64 entries
├─────────────────────┤
│   String Table       │  Symbol name strings
└─────────────────────┘
```

### Mach-O Header

```go
type MachOHeader64 struct {
    Magic      uint32  // MH_MAGIC_64 (0xFEEDFACF) or MH_CIGAM_64 (0xCFFAEDFE)
    CPUType    uint32  // CPU_TYPE_X86_64 (0x01000007)
    CPUSubtype uint32  // CPU_SUBTYPE_X86_64_ALL (0x00000003)
    FileType   uint32  // MH_EXECUTE (2), MH_DYLIB (6), MH_BUNDLE (8)
    NumCmds    uint32  // Number of load commands
    SizeOfCmds uint32  // Total size of all load commands
    Flags      uint32  // MH_NOUNDEFS, MH_DYLDLINK, MH_PIE, MH_TWOLEVEL, etc.
}
```

### Fat (Universal) Binaries

macOS frequently ships "fat" binaries containing code for multiple architectures (x86_64 + arm64). The fat binary format wraps multiple Mach-O images:

```
┌─────────────────────┐
│   Fat Header         │  magic (0xCAFEBABE or 0xBEBAFECA), nfat_arch
├─────────────────────┤
│   Fat Arch 1         │  cputype, cpusubtype, offset, size, align
│   Fat Arch 2         │
│   ...                │
├─────────────────────┤
│   Mach-O Image 1     │  Full Mach-O binary at specified offset
├─────────────────────┤
│   Mach-O Image 2     │
└─────────────────────┘
```

Our parser will auto-detect fat binaries and allow the caller to select the desired architecture.

### Load Commands

Key load commands we need to parse:

| Command            | Value    | Purpose                                  |
|--------------------|----------|------------------------------------------|
| `LC_SEGMENT_64`    | 0x19     | Defines a segment in memory              |
| `LC_SYMTAB`        | 0x02     | Symbol table location                    |
| `LC_DYSYMTAB`      | 0x0B     | Dynamic symbol table indices             |
| `LC_LOAD_DYLIB`    | 0x0C     | Dylib dependency path                    |
| `LC_MAIN`          | 0x80000028 | Entry point offset from __TEXT segment |
| `LC_UUID`          | 0x1B     | UUID for symbolication                   |
| `LC_BUILD_VERSION` | 0x32     | Platform, minos, tools (modern)          |
| `LC_VERSION_MIN_MACOSX` | 0x24 | Minimum macOS version (legacy)       |

### Section Naming

Mach-O sections are nested within segments:

```
__TEXT segment
    __text          — Executable code
    __stubs         — Indirect symbol stubs
    __stub_helper   — Stub helper code
    __cstring       — C string constants
    __const         — Constants
    __rodata        — Read-only data (Go)
    __typelink      — Go type links
    __itablink      — Go interface table links
    __gosymtab      — Go symbol table
    __gopclntab     — Go pclntab

__DATA segment
    __data          — Mutable data
    __bss           — Zero-initialized
    __noptrdata     — Non-pointer data (Go)
    __noptrbss      — Non-pointer BSS (Go)
    __go_buildinfo  — Go build info

__LINKEDIT segment
    — Symbol table, string table, indirect symbol table
```

### Go-specific Mach-O Notes

- Go uses the same build info format as ELF, stored in `__go_buildinfo` section
- pclntab is stored in `__gopclntab` section
- Go Mach-O binaries always include DWARF debug info
- Entry point is specified via `LC_MAIN` (not `LC_UNIXTHREAD` which is used by C compilers)

---

## Error Handling Pattern

All binary parsers follow the same error pattern:

```go
func Open(path string) (*Binary, error) {
    // Detect format by magic bytes
    // Open with format-specific parser
    // Validate required sections exist
    // Extract metadata
    return binary, nil
}

// Specific error types
type ErrUnsupportedArch struct { Arch string }
type ErrSectionNotFound struct { Name string }
type ErrNotGoBinary struct{}  // No Go build info found
```

## Byte Order Handling

- ELF: Big-endian or little-endian, specified in header
- PE: Always little-endian
- Mach-O: Big-endian or little-endian, encoded by magic number (0xFEEDFACE = native, 0xCEFAEDFE = swapped)

The `Binary` interface exposes `ByteOrder()` so consumers can correctly interpret multi-byte values in the binary.
