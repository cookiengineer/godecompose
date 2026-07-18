// Package syscall provides syscall table types and per-platform data
// for mapping syscall numbers back to function signatures.
package syscall

import (
	"github.com/cookiengineer/godecompose/types"
)

// Table represents a platform's complete syscall table.
type Table struct {
	Platform types.Platform `json:"-"`
	PlatformStr string      `json:"platform"`
	Arch       string       `json:"arch"`
	Version    string       `json:"version"`
	Entries    []Entry      `json:"syscalls"`
}

// Entry describes a single syscall with its number, name, and arguments.
type Entry struct {
	Number  uint64 `json:"number"`
	Name    string `json:"name"`
	Args    string `json:"args"`
	Returns string `json:"returns"`
}

// Lookup finds a syscall entry by number.
func (t *Table) Lookup(num uint64) (Entry, bool) {
	for _, e := range t.Entries {
		if e.Number == num {
			return e, true
		}
	}
	return Entry{}, false
}

// LookupByName finds a syscall entry by name.
func (t *Table) LookupByName(name string) (Entry, bool) {
	for _, e := range t.Entries {
		if e.Name == name {
			return e, true
		}
	}
	return Entry{}, false
}
