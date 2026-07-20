package generate

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/pattern/matcher"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
)

func makeInst(opcode string, intel string, addr uint64, size int) disasm.Instruction {
	return disasm.Instruction{
		Opcode:      opcode,
		IntelSyntax: intel,
		Address:     addr,
		Size:        size,
	}
}

func TestGenerateSingleMatch(t *testing.T) {
	pat := &evaluator.CompiledPattern{
		Name:     "syscall_write",
		GenTemplate: "write($fd, $buf, $count);",
		Bindings: []evaluator.CompiledBinding{
			{CaptureVar: "fd", Alias: "fd"},
			{CaptureVar: "buf", Alias: "buf"},
			{CaptureVar: "count", Alias: "count"},
		},
	}

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov eax, 1", 0x1000, 5),
		makeInst("MOV", "mov edi, 0", 0x1005, 5),
		makeInst("SYSCALL", "syscall", 0x100a, 2),
	}

	matches := []matcher.Match{
		{
			Pattern:    pat,
			StartAddr:  0x1000,
			EndAddr:    0x100c,
			Bindings: map[string]matcher.Binding{
				"fd":    {CaptureVar: "fd", Value: "edi"},
				"buf":   {CaptureVar: "buf", Value: "rsi"},
				"count": {CaptureVar: "count", Value: "edx"},
			},
		},
	}

	g := New(matches, instructions)
	output := g.Generate()

	if !strings.Contains(output, "write(") {
		t.Errorf("output doesn't contain write(): %q", output)
	}
	if !strings.Contains(output, "edi") {
		t.Errorf("output doesn't contain captured fd value: %q", output)
	}
	t.Logf("Generated output:\n%s", output)
}

func TestGenerateUnmatchedRange(t *testing.T) {
	pat := &evaluator.CompiledPattern{
		Name:     "test",
		GenTemplate: "matched();",
	}

	instructions := []disasm.Instruction{
		makeInst("NOP", "nop", 0x1000, 1),
		makeInst("NOP", "nop", 0x1001, 1),
		makeInst("MOV", "mov eax, 1", 0x1002, 5),
		makeInst("RET", "ret", 0x1007, 1),
	}

	matches := []matcher.Match{
		{
			Pattern:   pat,
			StartAddr: 0x1002,
			EndAddr:   0x1007,
		},
	}

	g := New(matches, instructions)
	output := g.Generate()

	if !strings.Contains(output, "unresolved") {
		t.Errorf("output doesn't contain unresolved code marker: %q", output)
	}
	if !strings.Contains(output, "matched()") {
		t.Errorf("output doesn't contain matched code: %q", output)
	}
	t.Logf("Generated output:\n%s", output)
}

func TestGenerateEmpty(t *testing.T) {
	g := New(nil, nil)
	output := g.Generate()
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestGenerateWithAliasBindings(t *testing.T) {
	pat := &evaluator.CompiledPattern{
		Name:     "go_memmove",
		GenTemplate: "memmove($dst, $src, $len);",
		Bindings: []evaluator.CompiledBinding{
			{CaptureVar: "src", Alias: "source"},
			{CaptureVar: "dst", Alias: "dest"},
			{CaptureVar: "len", Alias: "count"},
		},
	}

	instructions := []disasm.Instruction{
		makeInst("MOV", "mov rsi, rdi", 0x1000, 3),
	}

	matches := []matcher.Match{
		{
			Pattern:   pat,
			StartAddr: 0x1000,
			EndAddr:   0x1003,
			Bindings: map[string]matcher.Binding{
				"src": {CaptureVar: "src", Value: "rsi"},
				"dst": {CaptureVar: "dst", Value: "rdi"},
			},
		},
	}

	g := New(matches, instructions)
	output := g.Generate()

	if !strings.Contains(output, "memmove") {
		t.Errorf("output doesn't contain memmove: %q", output)
	}
	if strings.Contains(output, "$src") || strings.Contains(output, "$dst") {
		t.Errorf("output contains unexpanded placeholders: %q", output)
	}
	t.Logf("Generated output:\n%s", output)
}
