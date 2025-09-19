// Package slots defines the data structures for slot command storage.
package slots

import (
	"fmt"
	"io"
	"slices"
)

// Slot represents a saved command with metadata.
type Slot struct {
	// Name is the unique identifier for the slot.
	Name string
	// Description provides a brief explanation of the slot's purpose.
	Description string `json:"description,omitempty"`
	// Cmd is the command template with placeholders.
	Cmd string
	// Tags are optional labels for organizing slots.
	Tags []string `json:"tags,omitempty"`
}

// Slots is a slice of Slot structs.
type Slots []Slot

// Add adds a new slot..
func (s *Slots) Add(slot Slot) {
	*s = append(*s, slot)
}

// Delete removes the slot with the specified name, returning true if found and deleted.
func (s *Slots) Delete(name string) bool {
	i := s.index(name)
	if i == -1 {
		return false
	}

	*s = slices.Delete(*s, i, i+1)

	return true
}

// Exists checks if a slot with the given name exists.
func (s Slots) Exists(name string) bool {
	return s.index(name) != -1
}

// Get retrieves a pointer to the slot with the specified name, or nil if not found.
func (s Slots) Get(name string) *Slot {
	i := s.index(name)
	if i == -1 {
		return nil
	}

	return &s[i]
}

// Render outputs the slots in the specified format ("table" or "tsv") to the given writer.
func (s Slots) Render(format string, writer io.Writer) error {
	switch format {
	case "table":
		return asTable(s, writer)
	case "tsv":
		return asTSV(s, writer)
	default:
		return fmt.Errorf("unknown format: %q", format)
	}
}

// index returns the index of the slot with the given name, or -1 if not found.
func (s Slots) index(name string) int {
	return slices.IndexFunc(s, func(slot Slot) bool {
		return slot.Name == name
	})
}
