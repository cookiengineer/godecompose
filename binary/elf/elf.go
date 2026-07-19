package elf

import (
	"bytes"
	"debug/elf"
	encodingbinary "encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/types"
)

func init() {
	binary.RegisterFormat(binary.FormatELF, func(path string) (binary.Binary, error) {
		return Open(path)
	})
}

type Binary struct {
	file        *elf.File
	sections    []binary.Section
	symbols     []binary.Symbol
	dynsyms     []binary.Symbol
	goInfo      *binary.GoBuildInfo
	goInfoOK    bool
	pclntabData []byte
	pclntabAddr uint64
	pclntabOK   bool
}

func Open(path string) (*Binary, error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("elf.Open: %w", err)
	}

	b := &Binary{file: f}
	b.parseSections()
	b.parseGoBuildInfo()
	b.findPclntab()

	return b, nil
}

func (b *Binary) Format() binary.Format {
	return binary.FormatELF
}

func (b *Binary) Architecture() types.Arch {
	switch b.file.Machine {
	case elf.EM_X86_64:
		return types.ArchX86_64
	case elf.EM_AARCH64:
		return types.ArchAArch64
	default:
		return types.ArchUnknown
	}
}

func (b *Binary) EntryPoint() uint64 {
	return b.file.Entry
}

func (b *Binary) Sections() []binary.Section {
	return b.sections
}

func (b *Binary) Section(name string) (binary.Section, bool) {
	for _, s := range b.sections {
		if s.Name == name {
			return s, true
		}
	}
	return binary.Section{}, false
}

func (b *Binary) Symbols() ([]binary.Symbol, error) {
	if b.symbols != nil {
		return b.symbols, nil
	}
	syms, err := b.file.Symbols()
	if err != nil {
		return nil, fmt.Errorf("elf.Symbols: %w", err)
	}
	b.symbols = make([]binary.Symbol, 0, len(syms))
	for _, s := range syms {
		b.symbols = append(b.symbols, b.convertSymbol(s))
	}
	return b.symbols, nil
}

func (b *Binary) DynamicSymbols() ([]binary.Symbol, error) {
	if b.dynsyms != nil {
		return b.dynsyms, nil
	}
	syms, err := b.file.DynamicSymbols()
	if err != nil {
		return nil, fmt.Errorf("elf.DynamicSymbols: %w", err)
	}
	b.dynsyms = make([]binary.Symbol, 0, len(syms))
	for _, s := range syms {
		b.dynsyms = append(b.dynsyms, b.convertSymbol(s))
	}
	return b.dynsyms, nil
}

func (b *Binary) IsPIE() bool {
	return b.file.Type == elf.ET_DYN
}

func (b *Binary) IsStripped() bool {
	syms, err := b.Symbols()
	if err != nil {
		return true
	}
	return len(syms) == 0
}

func (b *Binary) ByteOrder() encodingbinary.ByteOrder {
	return b.file.ByteOrder
}

func (b *Binary) GoBuildInfo() (*binary.GoBuildInfo, bool) {
	return b.goInfo, b.goInfoOK
}

func (b *Binary) Pclntab() ([]byte, uint64, bool) {
	return b.pclntabData, b.pclntabAddr, b.pclntabOK
}

func (b *Binary) Close() error {
	return b.file.Close()
}

func (b *Binary) parseSections() {
	b.sections = make([]binary.Section, 0, len(b.file.Sections))
	for _, s := range b.file.Sections {
		data, _ := s.Data()
		sec := binary.Section{
			Name:    s.SectionHeader.Name,
			Address: s.SectionHeader.Addr,
			Size:    s.SectionHeader.Size,
			Offset:  uint64(s.SectionHeader.Offset),
			Data:    data,
			Flags:   elfSectionFlags(s.SectionHeader.Flags),
			EntSize: s.SectionHeader.Entsize,
			Link:    s.SectionHeader.Link,
			Info:    s.SectionHeader.Info,
			Align:   s.SectionHeader.Addralign,
		}
		b.sections = append(b.sections, sec)
	}
}

func (b *Binary) parseGoBuildInfo() {
	for _, s := range b.file.Sections {
		if s.Type != elf.SHT_NOTE && s.Type != elf.SHT_PROGBITS {
			continue
		}
		if s.Type == elf.SHT_PROGBITS && s.Name != ".go.buildinfo" {
			continue
		}
		data, err := s.Data()
		if err != nil {
			continue
		}
		info := parseGoBuildInfo(bytes.NewReader(data))
		if info != nil {
			b.goInfo = info
			b.goInfoOK = true
			return
		}
	}
}

func (b *Binary) findPclntab() {
	for _, s := range b.sections {
		if s.Name == ".gopclntab" || s.Name == ".data.rel.ro.gopclntab" {
			b.pclntabData = s.Data
			b.pclntabAddr = s.Address
			b.pclntabOK = true
			return
		}
	}
}

func (b *Binary) convertSymbol(s elf.Symbol) binary.Symbol {
	return binary.Symbol{
		Name:    s.Name,
		Address: s.Value,
		Size:    s.Size,
		Type:    convertSymbolType(s),
		Section: b.sectionNameByIndex(int(s.Section)),
		Version: s.Version,
		Binding: convertSymbolBinding(s),
	}
}

func (b *Binary) sectionNameByIndex(idx int) string {
	if idx < 0 || idx >= len(b.sections) {
		return ""
	}
	return b.sections[idx].Name
}

func convertSymbolType(s elf.Symbol) binary.SymbolType {
	switch elf.ST_TYPE(s.Info) {
	case elf.STT_FUNC:
		return binary.SymbolFunction
	case elf.STT_OBJECT:
		return binary.SymbolObject
	case elf.STT_SECTION:
		return binary.SymbolSection
	case elf.STT_FILE:
		return binary.SymbolFile
	default:
		return binary.SymbolUnknown
	}
}

func convertSymbolBinding(s elf.Symbol) binary.SymbolBinding {
	switch elf.ST_BIND(s.Info) {
	case elf.STB_LOCAL:
		return binary.SymbolLocal
	case elf.STB_GLOBAL:
		return binary.SymbolGlobal
	case elf.STB_WEAK:
		return binary.SymbolWeak
	default:
		return binary.SymbolLocal
	}
}

func elfSectionFlags(flags elf.SectionFlag) binary.SectionFlag {
	var f binary.SectionFlag
	if flags&elf.SHF_ALLOC != 0 {
		f |= binary.SectionReadable
	}
	if flags&elf.SHF_WRITE != 0 {
		f |= binary.SectionWritable
	}
	if flags&elf.SHF_EXECINSTR != 0 {
		f |= binary.SectionExecutable
	}
	return f
}

// parseGoBuildInfo searches for and parses the Go build info from ELF note data.
func parseGoBuildInfo(r io.Reader) *binary.GoBuildInfo {
	data, err := io.ReadAll(r)
	if err != nil || len(data) < 16 {
		return nil
	}

	// Try the V1 header format: "\x10Go build info V1\n"
	headerV1 := []byte("\x10Go build info V1\n")
	if idx := bytes.Index(data, headerV1); idx >= 0 {
		return decodeGoBuildInfoV1(data[idx+len(headerV1):])
	}

	// Try the V2 / current format: "\xff Go buildinf:" + 2 bytes metadata + 2 pointer words padding
	headerV2 := []byte("\xff Go buildinf:")
	if idx := bytes.Index(data, headerV2); idx >= 0 {
		meta := idx + len(headerV2)
		if meta+2 > len(data) {
			return nil
		}
		ptrSize := int(data[meta+0])
		padding := ptrSize * 2
		start := meta + 2 + padding
		if start < len(data) {
			return decodeGoBuildInfoV1(data[start:])
		}
	}

	return nil
}

func decodeGoBuildInfoV1(data []byte) *binary.GoBuildInfo {
	info := &binary.GoBuildInfo{
		Settings: make(map[string]string),
	}

	pos := 0

	readUVarint := func() (uint64, bool) {
		var x uint64
		var s uint
		for i := 0; i < encodingbinary.MaxVarintLen64; i++ {
			if pos >= len(data) {
				return 0, false
			}
			b := data[pos]
			pos++
			x |= uint64(b&0x7f) << s
			if b < 0x80 {
				return x, true
			}
			s += 7
		}
		return 0, false
	}

	readBytes := func(n int) ([]byte, bool) {
		if pos+n > len(data) {
			return nil, false
		}
		result := data[pos : pos+n]
		pos += n
		return result, true
	}

	readString := func() (string, bool) {
		length, ok := readUVarint()
		if !ok {
			return "", false
		}
		b, ok := readBytes(int(length))
		if !ok {
			return "", false
		}
		return string(b), true
	}

	// V1 format: Version, Path, Main, then deps and settings
	info.Version, _ = readString()
	info.Path, _ = readString()
	info.Main, _ = readString()

	numDeps, ok := readUVarint()
	if !ok {
		return info
	}
	info.Deps = make([]binary.GoModuleDep, 0, int(numDeps))
	for i := uint64(0); i < numDeps; i++ {
		dep := binary.GoModuleDep{}
		dep.Path, _ = readString()
		dep.Version, _ = readString()
		dep.Sum, _ = readString()
		info.Deps = append(info.Deps, dep)
	}

	numSettings, ok := readUVarint()
	if !ok {
		return info
	}
	for i := uint64(0); i < numSettings; i++ {
		key, _ := readString()
		value, _ := readString()
		info.Settings[key] = value
	}

	return info
}

// FormatGoSymbol replaces Go assembler middle-dot separators (U+00B7) with periods.
func FormatGoSymbol(name string) string {
	return strings.ReplaceAll(name, "\u00b7", ".")
}
