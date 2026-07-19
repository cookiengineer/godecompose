package pe

import (
	"bytes"
	"debug/pe"
	encodingbinary "encoding/binary"
	"fmt"
	"strings"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/types"
)

func init() {
	binary.RegisterFormat(binary.FormatPE, func(path string) (binary.Binary, error) {
		return Open(path)
	})
}

type Binary struct {
	file        *pe.File
	sections    []binary.Section
	symbols     []binary.Symbol
	goInfo      *binary.GoBuildInfo
	goInfoOK    bool
	pclntabData []byte
	pclntabAddr uint64
	pclntabOK   bool
}

type sectionData struct {
	data []byte
}

func Open(path string) (*Binary, error) {
	f, err := pe.Open(path)
	if err != nil {
		return nil, fmt.Errorf("pe.Open: %w", err)
	}

	b := &Binary{file: f}
	b.parseSections()
	b.parseGoBuildInfo()
	b.findPclntab()

	return b, nil
}

func (b *Binary) Format() binary.Format {
	return binary.FormatPE
}

func (b *Binary) Architecture() types.Arch {
	switch b.file.Machine {
	case pe.IMAGE_FILE_MACHINE_AMD64:
		return types.ArchX86_64
	case pe.IMAGE_FILE_MACHINE_ARM64:
		return types.ArchAArch64
	default:
		return types.ArchUnknown
	}
}

func (b *Binary) EntryPoint() uint64 {
	switch oh := b.file.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		return oh.ImageBase + uint64(oh.AddressOfEntryPoint)
	case *pe.OptionalHeader32:
		return uint64(oh.ImageBase) + uint64(oh.AddressOfEntryPoint)
	}
	return 0
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
	syms := b.file.Symbols
	b.symbols = make([]binary.Symbol, 0, len(syms))
	for _, s := range syms {
		b.symbols = append(b.symbols, b.convertSymbol(s))
	}
	return b.symbols, nil
}

func (b *Binary) DynamicSymbols() ([]binary.Symbol, error) {
	imps, err := b.file.ImportedSymbols()
	if err != nil {
		return nil, fmt.Errorf("pe.DynamicSymbols: %w", err)
	}
	syms := make([]binary.Symbol, 0, len(imps))
	for _, s := range imps {
		syms = append(syms, binary.Symbol{
			Name: s,
			Type: binary.SymbolFunction,
		})
	}
	return syms, nil
}

func (b *Binary) IsPIE() bool {
	// PE uses the DLL characteristics flag for ASLR-compatible images
	if oh, ok := b.file.OptionalHeader.(*pe.OptionalHeader64); ok {
		return oh.DllCharacteristics&pe.IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE != 0
	}
	return false
}

func (b *Binary) IsStripped() bool {
	return len(b.file.Symbols) == 0
}

func (b *Binary) ByteOrder() encodingbinary.ByteOrder {
	return encodingbinary.LittleEndian
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
		data := make([]byte, s.Size)
		read, _ := s.ReadAt(data, 0)
		data = data[:read]

		sec := binary.Section{
			Name:    s.Name,
			Address: uint64(s.VirtualAddress),
			Size:    uint64(s.VirtualSize),
			Offset:  uint64(s.Offset),
			Data:    data,
			Flags:   peSectionFlags(s.Characteristics),
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
		if s.Name == ".gopclntab" {
			b.pclntabData = s.Data
			b.pclntabAddr = s.Address
			b.pclntabOK = true
			return
		}
	}
}

func (b *Binary) convertSymbol(s *pe.Symbol) binary.Symbol {
	sectionName := ""
	if s.SectionNumber > 0 && int(s.SectionNumber) <= len(b.file.Sections) {
		sectionName = b.file.Sections[s.SectionNumber-1].Name
	}

	symType := binary.SymbolUnknown
	if s.SectionNumber > 0 && s.StorageClass == 2 {
		symType = binary.SymbolFunction
	}

	return binary.Symbol{
		Name:    s.Name,
		Address: uint64(s.Value),
		Type:    symType,
		Section: sectionName,
	}
}

func peSectionFlags(ch uint32) binary.SectionFlag {
	var f binary.SectionFlag
	if ch&pe.IMAGE_SCN_MEM_EXECUTE != 0 {
		f |= binary.SectionExecutable
	}
	if ch&pe.IMAGE_SCN_MEM_WRITE != 0 {
		f |= binary.SectionWritable
	}
	if ch&pe.IMAGE_SCN_MEM_READ != 0 {
		f |= binary.SectionReadable
	}
	return f
}

func parseGoBuildInfoFromData(data []byte) (*binary.GoBuildInfo, bool) {
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

func parseGoBuildInfoText(data string, info *binary.GoBuildInfo) {
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
				dep := binary.GoModuleDep{
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
