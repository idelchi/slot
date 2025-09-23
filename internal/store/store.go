// Package store provides persistent storage operations for slot data with atomic file operations.
package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/slot/internal/slot"
)

// Store handles persistent storage operations for slot data.
type Store string

// Path returns the file path of the slot store.
func (store Store) Path() string {
	return string(store)
}

// ExpandHome resolves home directory references in paths.
// Replaces leading ~ with the user's home directory path.
// Returns the original path if expansion fails or isn't needed.
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == "~" {
		return home
	}

	return filepath.Join(home, path[1:])
}

// New creates a new Store instance from the given file path.
func New(slotsFile string) (Store, error) {
	slotsFile = os.ExpandEnv(ExpandHome(slotsFile))

	store := Store(slotsFile)

	dataDir := filepath.Dir(slotsFile)

	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return store, fmt.Errorf("creating data directory: %w", err)
	}

	return store, nil
}

// Load reads the slots from disk, returning an empty slots object if the file doesn't exist.
func (store *Store) Load() (slot.Slots, error) {
	slots := slot.Slots{}

	data, err := os.ReadFile(store.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return slots, nil
		}

		return slots, fmt.Errorf("reading slots file: %w", err)
	}

	if err := yaml.Unmarshal(data, &slots); err != nil {
		return slots, fmt.Errorf("unmarshalling slots: %w", err)
	}

	return slots, nil
}

// Save writes the slots to disk.
func (store Store) Save(slots slot.Slots) error {
	data := []byte{}

	for i := range slots {
		slot, err := yaml.MarshalWithOptions(slots.Slice(i, i+1), yaml.UseLiteralStyleIfMultiline(true))
		if err != nil {
			return fmt.Errorf("marshalling slots: %w", err)
		}

		data = append(data, slot...)

		data = append(data, '\n')
	}

	// Trim the last newline
	if len(data) > 0 {
		data = data[:len(data)-1]
	}

	if err := os.WriteFile(store.Path(), data, 0o600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// DefaultSlotsFile returns the full path to the default slots file location.
func DefaultSlotsFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "slots.yaml", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.ToSlash(filepath.Join(home, ".config", "slot", "slots.yaml")), nil
}
