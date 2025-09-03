// Package model defines the core data structures for slot command storage.
package model

// Slot represents a saved command with metadata.
type Slot struct {
	// Cmd is the command template with placeholders.
	Cmd string
	// Tags are optional labels for organizing slots.
	Tags []string `json:"tags,omitempty"`
}

// DB represents the slot database structure.
type DB struct {
	// Slots maps slot names to their definitions.
	Slots map[string]Slot
}
