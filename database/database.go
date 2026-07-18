// Package database manages the pattern database: loading pattern files,
// indexing them by opcode, and querying by target platform/arch.
package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cookiengineer/godecompose/database/syscall"
	"github.com/cookiengineer/godecompose/pattern/lang/ast"
	"github.com/cookiengineer/godecompose/pattern/lang/evaluator"
	"github.com/cookiengineer/godecompose/pattern/lang/lexer"
	"github.com/cookiengineer/godecompose/pattern/lang/parser"
	"github.com/cookiengineer/godecompose/pattern/lang/preprocessor"
	"github.com/cookiengineer/godecompose/types"
)

// Database holds compiled patterns and syscall tables indexed for fast lookup.
type Database struct {
	Patterns []*evaluator.CompiledPattern
	Syscalls map[types.Platform]*syscall.Table
	byOpcode map[string][]*evaluator.CompiledPattern
}

// New creates an empty database.
func New() *Database {
	return &Database{
		Syscalls: make(map[types.Platform]*syscall.Table),
		byOpcode: make(map[string][]*evaluator.CompiledPattern),
	}
}

// LoadPatternsFromDir loads all .hexpat files from a directory recursively.
func (db *Database) LoadPatternsFromDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".hexpat") {
			return nil
		}

		patterns, err := db.loadPatternFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		for _, p := range patterns {
			db.addPattern(p)
		}

		return nil
	})
}

// LoadSyscallsFromDir loads all JSON syscall tables from a directory recursively.
func (db *Database) LoadSyscallsFromDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		table, err := db.loadSyscallFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		db.Syscalls[table.Platform] = table
		return nil
	})
}

// FindPatterns returns patterns matching the given arch and platform filters.
func (db *Database) FindPatterns(arch types.Arch, platform types.Platform) []*evaluator.CompiledPattern {
	var result []*evaluator.CompiledPattern
	for _, p := range db.Patterns {
		if p.Arch != "" && p.Arch != arch.String() {
			continue
		}
		if len(p.Platforms) > 0 {
			platformMatch := false
			platformStr := platform.String()
			for _, plat := range p.Platforms {
				if strings.EqualFold(plat, platformStr) {
					platformMatch = true
					break
				}
			}
			if !platformMatch {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}

// SyscallTable returns the syscall table for a given platform.
func (db *Database) SyscallTable(platform types.Platform) (*syscall.Table, bool) {
	t, ok := db.Syscalls[platform]
	return t, ok
}

// Stats returns string statistics about the database.
func (db *Database) Stats() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Patterns: %d", len(db.Patterns)))
	lines = append(lines, fmt.Sprintf("Opcode index: %d entries", len(db.byOpcode)))
	lines = append(lines, fmt.Sprintf("Syscall tables: %d", len(db.Syscalls)))
	for plat, t := range db.Syscalls {
		lines = append(lines, fmt.Sprintf("  %s: %d syscalls", plat.String(), len(t.Entries)))
	}

	platforms := make(map[string]int)
	archs := make(map[string]int)
	for _, p := range db.Patterns {
		if p.Arch != "" {
			archs[p.Arch]++
		}
		for _, plat := range p.Platforms {
			platforms[plat]++
		}
	}

	if len(archs) > 0 {
		var a []string
		for arch, count := range archs {
			a = append(a, fmt.Sprintf("%s=%d", arch, count))
		}
		sort.Strings(a)
		lines = append(lines, fmt.Sprintf("Arch distribution: %s", strings.Join(a, ", ")))
	}
	if len(platforms) > 0 {
		var p []string
		for plat, count := range platforms {
			p = append(p, fmt.Sprintf("%s=%d", plat, count))
		}
		sort.Strings(p)
		lines = append(lines, fmt.Sprintf("Platform distribution: %s", strings.Join(p, ", ")))
	}

	return strings.Join(lines, "\n")
}

func (db *Database) loadPatternFile(path string) ([]*evaluator.CompiledPattern, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source := string(data)
	l := lexer.NewWithFile(source, path)
	tokens, err := l.Lex()
	if err != nil {
		return nil, err
	}

	pp := preprocessor.New(&preprocessor.FileResolver{BaseDir: filepath.Dir(path)})
	tokens, err = pp.Process(tokens)
	if err != nil {
		return nil, err
	}

	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	e := evaluator.New()
	return e.Evaluate(prog)
}

func (db *Database) loadSyscallFile(path string) (*syscall.Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var table syscall.Table
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, err
	}

	table.Platform = types.PlatformFromString(table.PlatformStr)

	return &table, nil
}

func (db *Database) addPattern(p *evaluator.CompiledPattern) {
	db.Patterns = append(db.Patterns, p)
	for _, alt := range p.Alternatives {
		if len(alt) == 0 {
			continue
		}
		key := strings.ToUpper(alt[0].Opcode)
		db.byOpcode[key] = append(db.byOpcode[key], p)
	}
}

// ByOpcode returns patterns indexed by their first opcode, useful for matcher pre-filtering.
func (db *Database) ByOpcode() map[string][]*evaluator.CompiledPattern {
	return db.byOpcode
}

// AllPatterns returns all compiled patterns.
func (db *Database) AllPatterns() []*evaluator.CompiledPattern {
	return db.Patterns
}

// unused import guard
var _ = ast.Program{}
