// Package store provides persistent storage operations for slot data with atomic file operations.
package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/slot/internal/model"
)

// Store handles persistent storage operations for slot data.
type Store struct {
	// Path is the location of the slots database file.
	Path string
}

// New creates a new Store instance with resolved data directory paths.
func New() (*Store, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	return &Store{
		Path: filepath.Join(dataDir, "slots.yaml"),
	}, nil
}

// Load reads the slot database from disk, returning an empty database if the file doesn't exist.
func (store *Store) Load() (model.DB, error) {
	database := model.DB{Slots: make(map[string]model.Slot)}

	data, err := os.ReadFile(store.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return database, nil
		}

		return database, fmt.Errorf("reading slots file: %w", err)
	}

	if err := yaml.Unmarshal(data, &database); err != nil {
		return database, fmt.Errorf("unmarshalling slots: %w", err)
	}

	return database, nil
}

// Save writes the slot database to disk.
func (store *Store) Save(database model.DB) error {
	data, err := yaml.MarshalWithOptions(database)
	if err != nil {
		return fmt.Errorf("marshalling slots: %w", err)
	}

	if err := os.WriteFile(store.Path, data, 0o600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GetDataDir resolves the data directory path.
func GetDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.Join(home, ".config", "slot"), nil
}
