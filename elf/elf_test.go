package elf

import (
	"debug/elf"
	"os"
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/binary"
)

func TestOpenSelf(t *testing.T) {
	path := "/proc/self/exe"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("%s not available: %v", path, err)
	}

	b, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%s) = %v", path, err)
	}
	defer b.Close()

	if b.Format() != binary.FormatELF {
		t.Errorf("Format = %s, want ELF", b.Format())
	}

	if b.Architecture().String() == "unknown" {
		t.Error("architecture should not be unknown")
	}

	if b.EntryPoint() == 0 {
		t.Error("entry point is 0")
	}

	sections := b.Sections()
	if len(sections) == 0 {
		t.Fatal("no sections")
	}

	foundText := false
	for _, s := range sections {
		if s.Name == ".text" {
			foundText = true
			if len(s.Data) == 0 {
				t.Error(".text section data is empty")
			}
			break
		}
	}
	if !foundText {
		t.Log("no .text section found (stripped binary?)")
	}

	syms, err := b.Symbols()
	if err != nil {
		t.Logf("Symbols() error: %v (possibly stripped)", err)
	} else if len(syms) > 0 {
		hasMain := false
		for _, s := range syms {
			if strings.Contains(s.Name, "main") {
				hasMain = true
				break
			}
		}
		if !hasMain && len(syms) > 10 {
			t.Logf("no 'main' symbol found among %d symbols", len(syms))
		}
	}

	binarySection, ok := b.Section(".text")
	if ok {
		if binarySection.Name != ".text" {
			t.Errorf("Section(\".text\").Name = %q", binarySection.Name)
		}
	}
}

func TestFormatGoSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fmt·Println", "fmt.Println"},
		{"runtime·memmove", "runtime.memmove"},
		{"normal.name", "normal.name"},
		{"", ""},
	}

	for _, tt := range tests {
		got := FormatGoSymbol(tt.input)
		if got != tt.expected {
			t.Errorf("FormatGoSymbol(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestElfSectionFlagConversion(t *testing.T) {
	tests := []struct {
		name  string
		flags elf.SectionFlag
		check func(binary.SectionFlag) bool
	}{
		{
			"executable",
			elf.SHF_ALLOC | elf.SHF_EXECINSTR,
			func(f binary.SectionFlag) bool {
				return f&binary.SectionReadable != 0 && f&binary.SectionExecutable != 0 && f&binary.SectionWritable == 0
			},
		},
		{
			"writable",
			elf.SHF_ALLOC | elf.SHF_WRITE,
			func(f binary.SectionFlag) bool {
				return f&binary.SectionReadable != 0 && f&binary.SectionWritable != 0 && f&binary.SectionExecutable == 0
			},
		},
		{
			"readonly",
			elf.SHF_ALLOC,
			func(f binary.SectionFlag) bool {
				return f&binary.SectionReadable != 0 && f&binary.SectionWritable == 0 && f&binary.SectionExecutable == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := elfSectionFlags(tt.flags)
			if !tt.check(got) {
				t.Errorf("elfSectionFlags(%x) = %v", tt.flags, got)
			}
		})
	}
}
