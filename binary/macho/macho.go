// Package macho provides a Mach-O binary parser implementing the
// binary.Binary interface. Supports thin (single-arch) and fat (universal)
// binaries for x86_64.
package macho

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"fmt"
	"strings"

	gbinary "github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/types"
)

func init() {
	gbinary.RegisterFormat(gbinary.FormatMachO, func(path string) (gbinary.Binary, error) {
		return Open(path)
	})
}

// Open parses a Mach-O file (thin or fat). For fat binaries, prefers x86_64.
func Open(path string) (gbinary.Binary, error) {
	f, err := macho.Open(path)
	if err != nil {
		return nil, fmt.Errorf("macho.Open: %w", err)
	}

	b := &Binary{file: f}
	b.parseSections()
	b.parseGoBuildInfo()
	b.findPclntab()

	return b, nil
}

type Binary struct {
	file        *macho.File
	sections    []gbinary.Section
	symbols     []gbinary.Symbol
	goInfo      *gbinary.GoBuildInfo
	goInfoOK    bool
	pclntabData []byte
	pclntabAddr uint64
	pclntabOK   bool
}

func (b *Binary) Format() gbinary.Format {
	return gbinary.FormatMachO
}

func (b *Binary) Architecture() types.Arch {
	switch b.file.Cpu {
	case macho.CpuAmd64:
		return types.ArchX86_64
	case macho.CpuArm64:
		return types.ArchAArch64
	default:
		return types.ArchUnknown
	}
}

func (b *Binary) EntryPoint() uint64 {
	for _, sect := range b.sections {
		if sect.Name == "__TEXT,__text" && sect.Flags&gbinary.SectionExecutable != 0 {
			return sect.Address
		}
	}
	for _, sect := range b.sections {
		if sect.Name == "__text" && sect.Flags&gbinary.SectionExecutable != 0 {
			return sect.Address
		}
	}
	for _, sect := range b.sections {
		if stringsContains(sect.Name, "__text") && sect.Flags&gbinary.SectionExecutable != 0 {
			return sect.Address
		}
	}
	return 0
}

func (b *Binary) Sections() []gbinary.Section {
	return b.sections
}

func (b *Binary) Section(name string) (gbinary.Section, bool) {
	for _, s := range b.sections {
		if s.Name == name {
			return s, true
		}
	}
	return gbinary.Section{}, false
}

func (b *Binary) Symbols() ([]gbinary.Symbol, error) {
	if b.symbols != nil {
		return b.symbols, nil
	}
	syms := b.file.Symtab.Syms
	b.symbols = make([]gbinary.Symbol, 0, len(syms))
	for _, s := range syms {
		b.symbols = append(b.symbols, gbinary.Symbol{
			Name:    s.Name,
			Type:    gbinary.SymbolFunction,
			Address: s.Value,
		})
	}
	return b.symbols, nil
}

func (b *Binary) DynamicSymbols() ([]gbinary.Symbol, error) {
	// Mach-O dynamic symbols are a subset of the total symbol table.
	// Return all symbols with external visibility.
	syms, err := b.Symbols()
	if err != nil {
		return nil, err
	}
	dsyms := make([]gbinary.Symbol, 0)
	for _, s := range syms {
		if s.Name != "" && !contains(s.Name, '.') {
			dsyms = append(dsyms, s)
		}
	}
	return dsyms, nil
}

func (b *Binary) IsPIE() bool {
	const flagPIE = 0x200000
	return b.file.Flags&flagPIE != 0
}

func (b *Binary) IsStripped() bool {
	return len(b.file.Symtab.Syms) == 0
}

func (b *Binary) ByteOrder() binary.ByteOrder {
	return b.file.ByteOrder
}

func (b *Binary) GoBuildInfo() (*gbinary.GoBuildInfo, bool) {
	return b.goInfo, b.goInfoOK
}

func (b *Binary) Pclntab() ([]byte, uint64, bool) {
	return b.pclntabData, b.pclntabAddr, b.pclntabOK
}

func (b *Binary) Close() error {
	return b.file.Close()
}

func (b *Binary) parseSections() {
	b.sections = make([]gbinary.Section, 0, len(b.file.Sections))
	for _, s := range b.file.Sections {
		data := make([]byte, s.Size)
		n, _ := s.ReadAt(data, 0)
		data = data[:n]

		flag := gbinary.SectionReadable
		if stringsContains(s.Seg, "__TEXT") {
			flag |= gbinary.SectionExecutable
		}
		if stringsContains(s.Seg, "__DATA") || stringsContains(s.Seg, "__BSS") {
			flag |= gbinary.SectionWritable
		}

		sec := gbinary.Section{
			Name:    fmt.Sprintf("%s,%s", s.Seg, s.Name),
			Address: s.Addr,
			Size:    s.Size,
			Offset:  uint64(s.Offset),
			Data:    data,
			Flags:   flag,
			Align:   uint64(s.Align),
		}
		b.sections = append(b.sections, sec)
	}
}

func (b *Binary) parseGoBuildInfo() {
	for _, s := range b.sections {
		if info, ok := parseGoBuildInfoFromData(s.Data); ok {
			b.goInfo = info
			b.goInfoOK = true
			return
		}
	}
}

func (b *Binary) findPclntab() {
	for _, s := range b.sections {
		if s.Name == "__TEXT,__gopclntab" || s.Name == "__DATA_CONST,__gopclntab" {
			b.pclntabData = s.Data
			b.pclntabAddr = s.Address
			b.pclntabOK = true
			return
		}
	}
}

func parseGoBuildInfoFromData(data []byte) (*gbinary.GoBuildInfo, bool) {
	if len(data) < 16 {
		return nil, false
	}

	// Try V1 header: "\x10Go build info V1\n"
	headerV1 := []byte("\x10Go build info V1\n")
	if idx := bytes.Index(data, headerV1); idx >= 0 {
		return decodeGoBuildInfoV1(data[idx+len(headerV1):]), true
	}

	// Try V2 header: "\xff Go buildinf:" + 2 bytes + 2 pointer words padding
	headerV2 := []byte("\xff Go buildinf:")
	if idx := bytes.Index(data, headerV2); idx >= 0 {
		meta := idx + len(headerV2)
		if meta+2 > len(data) {
			return nil, false
		}
		ptrSize := int(data[meta+0])
		padding := ptrSize * 2
		start := meta + 2 + padding
		if start < len(data) {
			return decodeGoBuildInfoV1(data[start:]), true
		}
	}

	return nil, false
}

func decodeGoBuildInfoV1(data []byte) *gbinary.GoBuildInfo {
	info := &gbinary.GoBuildInfo{
		Settings: make(map[string]string),
	}

	pos := 0

	readUVarint := func() (uint64, bool) {
		var x uint64
		var s uint
		for i := 0; i < binary.MaxVarintLen64; i++ {
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

	readString := func() (string, bool) {
		length, ok := readUVarint()
		if !ok {
			return "", false
		}
		if pos+int(length) > len(data) {
			return "", false
		}
		result := string(data[pos : pos+int(length)])
		pos += int(length)
		return result, true
	}

	// V1 format: Version, Path, Main, then deps and settings
	info.Version, _ = readString()
	modData, _ := readString()

	if strings.HasPrefix(modData, "0w") && len(modData) > 16 {
		parseGoBuildInfoText(modData[16:], info)
		return info
	}

	// Old V1 format: second string is Path, third is Main
	info.Path = modData
	info.Main, _ = readString()

	numDeps, ok := readUVarint()
	if !ok {
		return info
	}
	info.Deps = make([]gbinary.GoModuleDep, 0, int(numDeps))
	for i := uint64(0); i < numDeps; i++ {
		dep := gbinary.GoModuleDep{}
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

func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func contains(s string, ch byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return true
		}
	}
	return false
}

func parseGoBuildInfoText(data string, info *gbinary.GoBuildInfo) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "path":
			if len(parts) > 1 && info.Path == "" {
				info.Path = parts[1]
			}
		case "mod":
			if len(parts) >= 3 {
				dep := gbinary.GoModuleDep{
					Path:    parts[1],
					Version: parts[2],
				}
				if len(parts) > 3 && parts[3] != "" {
					dep.Sum = parts[3]
				}
				info.Deps = append(info.Deps, dep)
			}
		case "build":
			if len(parts) > 1 {
				kv := strings.SplitN(parts[1], "=", 2)
				if len(kv) == 2 {
					info.Settings[kv[0]] = kv[1]
				}
			}
		}
	}
}
