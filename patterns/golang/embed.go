// Package golang provides embedded Go patterns for the decompiler database.
// Patterns are organized into three modules loaded independently:
//
//	stdlib   — multi-instruction patterns (highest confidence)
//	runtime  — Go runtime function patterns
//	fallback — single-CALL idiomatic patterns for stripped binaries
//
// Each module is loaded via a dedicated Load* function that compiles
// patterns into a database.Database instance.
package golang

import (
	"embed"

	"github.com/cookiengineer/godecompose/database"
)

//go:embed stdlib/** runtime/** fallback/** controlflow/**
var PatternsFS embed.FS

// LoadStdlib loads all stdlib multi-instruction patterns into the database.
func LoadStdlib(db *database.Database) error {
	return db.LoadFromFS(PatternsFS, "stdlib")
}

// LoadRuntime loads all runtime patterns into the database.
func LoadRuntime(db *database.Database) error {
	return db.LoadFromFS(PatternsFS, "runtime")
}

// LoadFallback loads single-CALL fallback patterns for stripped binaries.
func LoadFallback(db *database.Database) error {
	return db.LoadFromFS(PatternsFS, "fallback")
}

// LoadControlFlow loads control flow reconstruction patterns.
func LoadControlFlow(db *database.Database) error {
	return db.LoadFromFS(PatternsFS, "controlflow")
}
