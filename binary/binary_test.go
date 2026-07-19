package binary_test

import (
	"os"
	"testing"

	"github.com/cookiengineer/godecompose/binary"
	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
)

func TestOpenELF(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skipf("cannot find executable: %v", err)
	}

	b, err := binary.Open(exe)
	if err != nil {
		t.Fatalf("Open(%s) = %v", exe, err)
	}
	defer b.Close()

	if b.Format() != binary.FormatELF {
		t.Errorf("Format() = %v, want ELF", b.Format())
	}

	if b.EntryPoint() == 0 {
		t.Error("EntryPoint() returned 0")
	}

	sections := b.Sections()
	if len(sections) == 0 {
		t.Error("no sections found")
	}

	foundText := false
	for _, s := range sections {
		if s.Name == ".text" {
			foundText = true
			if len(s.Data) == 0 {
				t.Error(".text section has no data")
			}
			break
		}
	}
	if !foundText {
		t.Log("no .text section found (stripped binary?)")
	}
}

func TestOpenNonExistent(t *testing.T) {
	_, err := binary.Open("/nonexistent/binary/path")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
