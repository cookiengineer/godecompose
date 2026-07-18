// Package binary defines types and the interface for parsed executable files.
// Format-specific packages (elf, pe, macho) implement this interface and
// register themselves via RegisterFormat.
package binary

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/cookiengineer/godecompose/types"
)

type Format int

const (
	FormatUnknown Format = iota
	FormatELF
	FormatPE
	FormatMachO
)

func (f Format) String() string {
	switch f {
	case FormatELF:
		return "ELF"
	case FormatPE:
		return "PE"
	case FormatMachO:
		return "Mach-O"
	default:
		return "unknown"
	}
}

type SectionFlag int

const (
	SectionExecutable SectionFlag = 1 << iota
	SectionWritable
	SectionReadable
)

type Section struct {
	Name    string
	Address uint64
	Size    uint64
	Offset  uint64
	Data    []byte
	Flags   SectionFlag
	EntSize uint64
	Link    uint32
	Info    uint32
	Align   uint64
}

type SymbolType int

const (
	SymbolUnknown SymbolType = iota
	SymbolFunction
	SymbolObject
	SymbolSection
	SymbolFile
)

type SymbolBinding int

const (
	SymbolLocal SymbolBinding = iota
	SymbolGlobal
	SymbolWeak
)

type Symbol struct {
	Name    string
	Address uint64
	Size    uint64
	Type    SymbolType
	Section string
	Version string
	Binding SymbolBinding
}

type GoModuleDep struct {
	Path    string
	Version string
	Sum     string
}

type GoBuildInfo struct {
	Version  string
	Path     string
	Main     string
	Deps     []GoModuleDep
	Settings map[string]string
}

type Binary interface {
	Format() Format
	Architecture() types.Arch
	EntryPoint() uint64
	Sections() []Section
	Section(name string) (Section, bool)
	Symbols() ([]Symbol, error)
	DynamicSymbols() ([]Symbol, error)
	IsPIE() bool
	IsStripped() bool
	ByteOrder() binary.ByteOrder
	GoBuildInfo() (*GoBuildInfo, bool)
	Pclntab() ([]byte, uint64, bool)
	Close() error
}

type Opener func(path string) (Binary, error)

var (
	mu      sync.RWMutex
	openers = make(map[Format]Opener)
)

func RegisterFormat(format Format, opener Opener) {
	mu.Lock()
	openers[format] = opener
	mu.Unlock()
}

// Open detects the binary format from magic bytes and returns a parsed Binary.
func Open(path string) (Binary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("binary.Open: %w", err)
	}

	magic := make([]byte, 4)
	n, err := f.Read(magic)
	f.Close()
	if err != nil {
		return nil, fmt.Errorf("binary.Open: read magic: %w", err)
	}
	if n < 4 {
		return nil, fmt.Errorf("binary.Open: file too short (%d bytes)", n)
	}

	format := detectFormat(magic)

	mu.RLock()
	opener, ok := openers[format]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("binary.Open: no opener registered for format %s (magic: %x)", format, magic)
	}

	return opener(path)
}

func detectFormat(magic []byte) Format {
	if len(magic) >= 4 && magic[0] == 0x7F && magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F' {
		return FormatELF
	}
	if len(magic) >= 2 && magic[0] == 'M' && magic[1] == 'Z' {
		return FormatPE
	}
	if len(magic) >= 4 {
		m := uint32(magic[0])<<24 | uint32(magic[1])<<16 | uint32(magic[2])<<8 | uint32(magic[3])
		switch m {
		case 0xFEEDFACE, 0xFEEDFACF, 0xCEFAEDFE, 0xCFFAEDFE:
			return FormatMachO
		case 0xCAFEBABE, 0xBEBAFECA:
			return FormatMachO
		}
	}
	return FormatUnknown
}
