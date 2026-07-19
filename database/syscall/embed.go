// Package syscall manages kernel syscall tables for pattern-based
// decompilation, mapping platform-specific syscall numbers to names.
package syscall

import "embed"

//go:embed tables/*/*.json
var TablesFS embed.FS
