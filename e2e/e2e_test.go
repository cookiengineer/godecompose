package e2e

import (
	"strings"
	"testing"

	"github.com/cookiengineer/godecompose/binary"
	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/goutil"
	_ "github.com/cookiengineer/godecompose/binary/elf"
	_ "github.com/cookiengineer/godecompose/binary/macho"
	_ "github.com/cookiengineer/godecompose/binary/pe"
)

func TestEndToEndSimpleELF(t *testing.T) {
	binaries := goutil.CompileSimple(t)
	path := goutil.GetBinary(t, binaries, "linux")
	if path == "" {
		t.Fatal("no linux binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	assertFormat(t, b, "ELF")
	assertArchitecture(t, b, "x86_64")
	assertSections(t, b)
	assertEntryPoint(t, b)
	assertGoBuildInfo(t, b)
	assertPclntab(t, b)
	assertSymbolContains(t, b, "main.main")
	assertSymbolContains(t, b, "main.factorial")
}

func TestEndToEndSimplePE(t *testing.T) {
	binaries := goutil.CompileSimple(t)
	path := goutil.GetBinary(t, binaries, "windows")
	if path == "" {
		t.Fatal("no windows binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	assertFormat(t, b, "PE")
	assertArchitecture(t, b, "x86_64")
	assertSections(t, b)
	assertEntryPoint(t, b)
	assertGoBuildInfo(t, b)
	assertSymbolContains(t, b, "main.main")
	assertSymbolContains(t, b, "main.factorial")
}

func TestEndToEndSimpleMachO(t *testing.T) {
	binaries := goutil.CompileSimple(t)
	path := goutil.GetBinary(t, binaries, "darwin")
	if path == "" {
		t.Fatal("no darwin binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	assertFormat(t, b, "Mach-O")
	assertArchitecture(t, b, "x86_64")
	assertSections(t, b)
	assertEntryPoint(t, b)
	assertGoBuildInfo(t, b)
	assertSymbolContains(t, b, "main.main")
	assertSymbolContains(t, b, "main.factorial")
}

func TestEndToEndDisassembly(t *testing.T) {
	binaries := goutil.CompileSimple(t)
	path := goutil.GetBinary(t, binaries, "linux")
	if path == "" {
		t.Fatal("no linux binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	textSection, ok := b.Section(".text")
	if !ok {
		t.Fatal("no .text section")
	}

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	if len(instructions) == 0 {
		t.Fatal("no instructions decoded from .text section")
	}

	t.Logf("decoded %d instructions from .text section (%d bytes)", len(instructions), textSection.Size)

	assertContainsInstruction(t, instructions, "CALL")
	assertContainsInstruction(t, instructions, "RET")

	blocks := disasm.BuildControlFlowGraph(instructions, nil)
	if len(blocks) == 0 {
		t.Error("no basic blocks built from .text")
	}
	t.Logf("  built %d basic blocks", len(blocks))
}

func TestEndToEndFunctionRecovery(t *testing.T) {
	binaries := goutil.CompileSimple(t)
	path := goutil.GetBinary(t, binaries, "linux")
	if path == "" {
		t.Fatal("no linux binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	textSection, ok := b.Section(".text")
	if !ok {
		t.Fatal("no .text section")
	}

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		t.Fatalf("RecoverFromBinary: %v", err)
	}

	if len(result.Functions) == 0 {
		t.Fatal("no functions recovered")
	}

	t.Logf("recovered %d functions:", len(result.Functions))
	for _, f := range result.Functions {
		if len(f.Blocks) > 0 {
			t.Logf("  %s @ 0x%x (blocks: %d)", f.Name, f.EntryPoint, len(f.Blocks))
		}
	}

	foundMain := false
	foundFactorial := false
	for _, f := range result.Functions {
		name := strings.ToLower(f.Name)
		if strings.Contains(name, "main.main") {
			foundMain = true
		}
		if strings.Contains(name, "factorial") {
			foundFactorial = true
		}
	}

	if !foundMain {
		t.Error("'main.main' function not found in recovered functions")
	}
	if !foundFactorial {
		t.Error("'factorial' function not found in recovered functions")
	}

	assertCFGForFunctions(t, result.Functions)
}

func TestEndToEndComplexBinary(t *testing.T) {
	binaries := goutil.CompileComplex(t)
	path := goutil.GetBinary(t, binaries, "linux")
	if path == "" {
		t.Fatal("no linux binary compiled")
	}

	b, err := binary.Open(path)
	if err != nil {
		t.Fatalf("binary.Open: %v", err)
	}
	defer b.Close()

	assertFormat(t, b, "ELF")
	assertArchitecture(t, b, "x86_64")
	assertGoBuildInfo(t, b)
	assertPclntab(t, b)

	textSection, ok := b.Section(".text")
	if !ok {
		t.Fatal("no .text section")
	}

	instructions, err := disasm.DecodeStream(textSection.Data, textSection.Address)
	if err != nil {
		t.Fatalf("DecodeStream: %v", err)
	}

	result, err := function.RecoverFromBinary(b, instructions)
	if err != nil {
		t.Fatalf("RecoverFromBinary: %v", err)
	}

	if len(result.Functions) == 0 {
		t.Fatal("no functions recovered from complex binary")
	}

	t.Logf("complex binary: %d functions recovered, %d instructions decoded",
		len(result.Functions), len(instructions))

	expectedFunctions := []string{
		"goroutineFanIn",
		"main.main",
	}

	found := make(map[string]bool)
	for _, f := range result.Functions {
		lower := strings.ToLower(f.Name)
		for _, expected := range expectedFunctions {
			if strings.Contains(lower, strings.ToLower(expected)) {
				found[expected] = true
			}
		}
	}

	for _, expected := range expectedFunctions {
		if !found[expected] {
			t.Errorf("function containing %q not found in recovered functions (%d total)",
				expected, len(result.Functions))
		}
	}

	assertCFGForFunctions(t, result.Functions)
}

func TestEndToEndCrossPlatformFormats(t *testing.T) {
	binaries := goutil.CompileSimple(t)

	expectedFormats := map[string]string{
		"linux":   "ELF",
		"windows": "PE",
		"darwin":  "Mach-O",
	}

	for osName, format := range expectedFormats {
		path := goutil.GetBinary(t, binaries, osName)
		if path == "" {
			t.Errorf("no binary for %s", osName)
			continue
		}

		b, err := binary.Open(path)
		if err != nil {
			t.Errorf("%s: binary.Open: %v", osName, err)
			continue
		}

		if b.Format().String() != format {
			t.Errorf("%s: format = %q, want %q", osName, b.Format().String(), format)
		}

		if b.Architecture().String() != "x86_64" {
			t.Errorf("%s: arch = %q, want x86_64", osName, b.Architecture().String())
		}

		sections := b.Sections()
		if len(sections) == 0 {
			t.Errorf("%s: no sections", osName)
		}

		b.Close()
	}
}

func assertFormat(t *testing.T, b binary.Binary, expected string) {
	t.Helper()
	got := b.Format().String()
	if got != expected {
		t.Errorf("Format() = %q, want %q", got, expected)
	}
}

func assertArchitecture(t *testing.T, b binary.Binary, expected string) {
	t.Helper()
	got := b.Architecture().String()
	if got != expected {
		t.Errorf("Architecture() = %q, want %q", got, expected)
	}
}

func assertSections(t *testing.T, b binary.Binary) {
	t.Helper()
	sections := b.Sections()
	if len(sections) == 0 {
		t.Error("Sections() returned empty")
		return
	}

	hasText := false
	for _, s := range sections {
		name := strings.ToLower(s.Name)
		if strings.Contains(name, ".text") || strings.Contains(name, "__text") {
			hasText = true
			if s.Flags&binary.SectionExecutable == 0 {
				t.Errorf("text section %q not marked executable", s.Name)
			}
			break
		}
	}
	if !hasText {
		t.Error("no executable text section found")
	}

	t.Logf("  %d sections (text found: %v)", len(sections), hasText)
}

func assertEntryPoint(t *testing.T, b binary.Binary) {
	t.Helper()
	entry := b.EntryPoint()
	if entry == 0 {
		t.Error("EntryPoint() returned 0")
	}
	t.Logf("  entry point: 0x%x", entry)
}

func assertGoBuildInfo(t *testing.T, b binary.Binary) {
	t.Helper()
	info, ok := b.GoBuildInfo()
	if !ok {
		t.Error("GoBuildInfo() returned false for a Go binary")
		return
	}
	if info.Version == "" {
		t.Error("GoBuildInfo.Version is empty")
	}
	t.Logf("  Go version: %s, path: %s", info.Version, info.Path)
}

func assertPclntab(t *testing.T, b binary.Binary) {
	t.Helper()
	data, addr, ok := b.Pclntab()
	if !ok {
		t.Error("Pclntab() returned false for a Go binary")
		return
	}
	if len(data) == 0 {
		t.Error("Pclntab() returned empty data")
	}
	t.Logf("  pclntab: %d bytes at 0x%x", len(data), addr)
}

func assertSymbolContains(t *testing.T, b binary.Binary, substr string) {
	t.Helper()
	syms, err := b.Symbols()
	if err != nil {
		t.Logf("  Symbols() error: %v (may be stripped)", err)
		return
	}

	for _, s := range syms {
		if strings.Contains(s.Name, substr) {
			t.Logf("  found symbol: %s", s.Name)
			return
		}
	}
	t.Errorf("symbol containing %q not found among %d symbols", substr, len(syms))
}

func assertContainsInstruction(t *testing.T, instructions []disasm.Instruction, opcode string) {
	t.Helper()
	for _, inst := range instructions {
		if inst.Opcode == opcode {
			return
		}
	}
	t.Errorf("no instruction with opcode %q found among %d instructions", opcode, len(instructions))
}

func assertCFGForFunctions(t *testing.T, functions []*function.Function) {
	t.Helper()
	for _, f := range functions {
		if len(f.Blocks) == 0 {
			continue
		}
		for _, block := range f.Blocks {
			if block.StartAddr < f.EntryPoint {
				t.Errorf("function %s: block starts at 0x%x before function entry 0x%x",
					f.Name, block.StartAddr, f.EntryPoint)
			}
		}
	}
}
