package actions

import (
	"fmt"

	"github.com/cookiengineer/godecompose/binary"
)

func Info(b binary.Binary) error {
	fmt.Printf("File:        %s\n", "embedded")
	fmt.Printf("Format:      %s\n", b.Format())
	fmt.Printf("Arch:        %s\n", b.Architecture())
	fmt.Printf("Entry:       0x%x\n", b.EntryPoint())
	fmt.Printf("PIE:         %v\n", b.IsPIE())
	fmt.Printf("Stripped:    %v\n", b.IsStripped())

	sections := b.Sections()
	fmt.Printf("Sections:    %d\n", len(sections))
	for _, s := range sections {
		fmt.Printf("  %-20s addr=0x%x size=0x%x flags=%c%c%c\n",
			s.Name, s.Address, s.Size,
			flagChar(s.Flags, binary.SectionExecutable, 'X'),
			flagChar(s.Flags, binary.SectionWritable, 'W'),
			flagChar(s.Flags, binary.SectionReadable, 'R'),
		)
	}

	syms, err := b.Symbols()
	if err == nil {
		fmt.Printf("Symbols:     %d\n", len(syms))
	}

	if info, ok := b.GoBuildInfo(); ok {
		fmt.Printf("Go version:  %s\n", info.Version)
		fmt.Printf("Go path:     %s\n", info.Path)
	}

	return nil
}

func flagChar(flags binary.SectionFlag, flag binary.SectionFlag, ch byte) byte {
	if flags&flag != 0 {
		return ch
	}
	return '-'
}
