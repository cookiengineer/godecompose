// Package generate produces decompiled source code from matched patterns
// by expanding gen templates with captured variable bindings, and can
// write a complete Go project directory structure.
package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

// Generator produces source code output from a list of pattern matches.
type Generator struct {
	matches      []matcher.Match
	instructions []disasm.Instruction
	functions    []*function.Function
	packages     map[string][]*function.Function
}

// New creates a generator from matched patterns and instructions.
func New(matches []matcher.Match, instructions []disasm.Instruction) *Generator {
	return &Generator{
		matches:      matches,
		instructions: instructions,
	}
}

// NewForProject creates a generator for project-directory output.
func NewForProject(
	matches []matcher.Match,
	instructions []disasm.Instruction,
	funcs []*function.Function,
	pkgs map[string][]*function.Function,
) *Generator {
	return &Generator{
		matches:      matches,
		instructions: instructions,
		functions:    funcs,
		packages:     pkgs,
	}
}

// Generate produces the full source code output as a single flat string.
func (g *Generator) Generate() string {
	sort.Slice(g.matches, func(i, j int) bool {
		return g.matches[i].StartAddr < g.matches[j].StartAddr
	})

	var buf strings.Builder
	lastAddr := uint64(0)

	if len(g.instructions) > 0 {
		lastAddr = g.instructions[0].Address
	}

	for _, match := range g.matches {
		if match.StartAddr > lastAddr {
			g.emitRawRange(&buf, lastAddr, match.StartAddr)
		}
		buf.WriteString(g.expandTemplate(match))
		buf.WriteString("\n")
		lastAddr = match.EndAddr
	}

	endAddr := lastAddr
	if len(g.instructions) > 0 {
		last := g.instructions[len(g.instructions)-1]
		endAddr = last.Address + uint64(last.Size)
	}
	if lastAddr < endAddr {
		g.emitRawRange(&buf, lastAddr, endAddr)
	}

	return buf.String()
}

// WriteProject creates a Go project directory with main.go and sub-packages.
func (g *Generator) WriteProject(dir string, goModule string) error {
	if goModule == "" {
		goModule = "decompiled"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Collect instructions per function for targeted generation
	funcInsts := make(map[string][]disasm.Instruction)
	for _, f := range g.functions {
		for _, inst := range g.instructions {
			if inst.Address >= f.EntryPoint && inst.Address < f.EndAddr {
				funcInsts[f.Name] = append(funcInsts[f.Name], inst)
			}
		}
	}

	// Generate main package (entry point)
	if mainPkg, ok := g.packages["main"]; ok {
		if err := g.writePackage(dir, "", mainPkg, funcInsts); err != nil {
			return err
		}
	}

	// Generate sub-packages
	for pkgPath, funcs := range g.packages {
		if pkgPath == "main" {
			continue
		}
		pkgDir := filepath.Join(dir, pkgPath)
		if err := g.writePackage(pkgDir, pkgPath, funcs, funcInsts); err != nil {
			return err
		}
	}

	// Write go.mod
	modContent := fmt.Sprintf("module %s\n\ngo 1.21\n", goModule)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	return nil
}

func (g *Generator) writePackage(dir string, pkgPath string, funcs []*function.Function, funcInsts map[string][]disasm.Instruction) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create package dir %s: %w", dir, err)
	}

	pkgName := "main"
	if pkgPath != "" {
		parts := strings.Split(pkgPath, "/")
		pkgName = parts[len(parts)-1]
	}

	fileName := pkgName + ".go"
	if pkgName == "main" {
		fileName = "main.go"
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))

	if pkgName == "main" {
		buf.WriteString("func main() {\n")
	} else {
		buf.WriteString("import \"fmt\"\n\n")
	}

	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].EntryPoint < funcs[j].EntryPoint
	})

	for _, f := range funcs {
		insts := funcInsts[f.Name]
		if len(insts) == 0 {
			continue
		}

		if pkgName == "main" && f.ShortName == "main" {
			g.writeFunctionBody(&buf, f, insts, "\t")
		} else if pkgName != "main" || f.ShortName != "main" {
			buf.WriteString(fmt.Sprintf("\nfunc %s() {\n", f.ShortName))
			g.writeFunctionBody(&buf, f, insts, "\t")
		}
	}

	if pkgName == "main" {
		buf.WriteString("}\n")
	}

	return os.WriteFile(filepath.Join(dir, fileName), []byte(buf.String()), 0644)
}

func (g *Generator) writeFunctionBody(buf *strings.Builder, f *function.Function, insts []disasm.Instruction, indent string) {
	if len(insts) == 0 {
		return
	}

	// Group matches by function
	var funcMatches []matcher.Match
	for _, m := range g.matches {
		if m.StartAddr >= f.EntryPoint && m.EndAddr <= f.EndAddr {
			funcMatches = append(funcMatches, m)
		}
	}

	sort.Slice(funcMatches, func(i, j int) bool {
		return funcMatches[i].StartAddr < funcMatches[j].StartAddr
	})

	lastAddr := f.EntryPoint
	for _, match := range funcMatches {
		if match.StartAddr > lastAddr {
			g.emitFunctionRange(buf, lastAddr, match.StartAddr, insts, indent)
		}
		template := g.expandTemplate(match)
		buf.WriteString(indent + strings.TrimSpace(template) + "\n")
		lastAddr = match.EndAddr
	}

	endAddr := f.EndAddr
	if lastAddr < endAddr {
		g.emitFunctionRange(buf, lastAddr, endAddr, insts, indent)
	}

	if f.ShortName != "main" {
		buf.WriteString("}\n")
	}
}

func (g *Generator) emitFunctionRange(buf *strings.Builder, start, end uint64, insts []disasm.Instruction, indent string) {
	count := 0
	for _, inst := range insts {
		if inst.Address >= start && inst.Address < end {
			if count == 0 {
				buf.WriteString(indent + "// unresolved\n")
			}
			buf.WriteString(fmt.Sprintf("%s// %016x: %s\n", indent, inst.Address, inst.IntelSyntax))
			count++
		}
	}
}

func (g *Generator) expandTemplate(match matcher.Match) string {
	template := match.Pattern.GenTemplate
	if template == "" {
		return fmt.Sprintf("// matched %s", match.Pattern.Name)
	}

	result := template
	for _, b := range match.Bindings {
		placeholder := "$" + b.CaptureVar
		value := b.Value
		if b.Alias != "" {
			value = b.Alias
		}
		result = strings.ReplaceAll(result, placeholder, value)
	}

	for _, b := range match.Pattern.Bindings {
		placeholder := "$" + b.CaptureVar
		if _, ok := match.Bindings[b.CaptureVar]; !ok {
			result = strings.ReplaceAll(result, placeholder, b.Alias)
		}
	}

	return result
}

func (g *Generator) emitRawRange(buf *strings.Builder, start, end uint64) {
	buf.WriteString("// --- unresolved code ---\n")
	for _, inst := range g.instructions {
		if inst.Address >= start && inst.Address < end {
			buf.WriteString(fmt.Sprintf("// %016x: %s\n", inst.Address, inst.IntelSyntax))
		}
	}
}
