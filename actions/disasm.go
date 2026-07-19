package actions

import (
	"fmt"
	"os"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
)

func Disassemble(b binary.Binary) error {
	textSection, ok := b.Section(".text")
	if !ok {
		return fmt.Errorf("no .text section found")
	}

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "disassembly error: %v\n", err)
	}
	fmt.Printf("Decoded %d instructions\n", len(instructions))

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "function recovery: %v\n", err)
		return nil
	}

	runtimeCount := 0
	stdlibCount := 0
	otherCount := 0
	for _, f := range result.Functions {
		switch f.Classification {
		case function.ClassRuntime:
			runtimeCount++
		case function.ClassStdlib:
			stdlibCount++
		default:
			otherCount++
		}
	}

	fmt.Printf("Functions: %d total\n", len(result.Functions))
	fmt.Printf("  runtime:  %d (skipped)\n", runtimeCount)
	fmt.Printf("  stdlib:   %d (skipped)\n", stdlibCount)
	fmt.Printf("  user:     %d\n", len(result.UserFunctions))
	fmt.Printf("  other:    %d\n", otherCount)

	if len(result.UserFunctions) > 0 {
		fmt.Printf("\nUser functions:\n")
		for _, f := range result.UserFunctions {
			fmt.Printf("  %-50s @ 0x%x (blocks: %d)\n", f.Name, f.EntryPoint, len(f.Blocks))
		}
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	return nil
}
