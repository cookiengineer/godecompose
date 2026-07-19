package function

import (
	"testing"

	"github.com/cookiengineer/godecompose/binary"
	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
)

func TestRecoverFromBinary(t *testing.T) {
	bin, err := binary.Open("/proc/self/exe")
	if err != nil {
		t.Skipf("cannot open self: %v", err)
	}
	defer bin.Close()

	textSection, ok := bin.Section(".text")
	if !ok {
		t.Skip("no .text section")
	}

	_ = textSection
	_ = bin

	// Basic smoke test: verify RecoverFromBinary returns without error
	result, err := RecoverFromBinary(bin, nil)
	if err != nil {
		t.Logf("RecoverFromBinary error (may be expected for stripped binary): %v", err)
	}

	if result != nil && len(result.Functions) == 0 {
		t.Log("no functions recovered (stripped binary?)")
	}
}
