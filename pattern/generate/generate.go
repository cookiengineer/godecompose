// Package generate produces decompiled source code from matched patterns
// by expanding gen templates with captured variable bindings.
package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

// Generator produces source code output from a list of pattern matches.
type Generator struct {
	matches      []matcher.Match
	instructions []disasm.Instruction
}

// New creates a generator from matched patterns and the original instruction stream.
func New(matches []matcher.Match, instructions []disasm.Instruction) *Generator {
	return &Generator{
		matches:      matches,
		instructions: instructions,
	}
}

// Generate produces the full source code output. Unmatched instruction ranges
// are emitted as raw assembly comments.
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

func (g *Generator) expandTemplate(match matcher.Match) string {
	template := match.Pattern.GenTemplate
	if template == "" {
		return fmt.Sprintf("// matched %s @ 0x%x", match.Pattern.Name, match.StartAddr)
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
		if _, alreadySet := match.Bindings[b.CaptureVar]; !alreadySet {
			result = strings.ReplaceAll(result, placeholder, b.Alias)
		}
	}

	return result
}

func (g *Generator) emitRawRange(buf *strings.Builder, start, end uint64) {
	buf.WriteString("// --- unresolved code ---\n")
	instCount := 0
	for _, inst := range g.instructions {
		if inst.Address >= start && inst.Address < end {
			buf.WriteString(fmt.Sprintf("// %016x: %s\n", inst.Address, inst.IntelSyntax))
			instCount++
		}
	}
	if instCount == 0 {
		buf.WriteString(fmt.Sprintf("// no instructions in range 0x%x-0x%x\n", start, end))
	}
}
