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
	// LogPath is the location of the audit log file.
	LogPath string
}

// New creates a new Store instance with resolved data directory paths.
func New() (*Store, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &Store{
		Path:    filepath.Join(dataDir, "slots.yaml"),
		LogPath: filepath.Join(dataDir, "slots.log"),
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

		return database, fmt.Errorf("failed to read slots file: %w", err)
	}

	if err := yaml.Unmarshal(data, &database); err != nil {
		return database, fmt.Errorf("failed to unmarshal slots: %w", err)
	}

	return database, nil
}

// Save writes the slot database to disk atomically using a temporary file.
func (store *Store) Save(database model.DB) error {
	data, err := yaml.MarshalWithOptions(
		database,
	) // yaml.UseLiteralStyleIfMultiline(true), // Use literal style for multiline strings
	if err != nil {
		return fmt.Errorf("failed to marshal slots: %w", err)
	}

	tempPath := store.Path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, store.Path); err != nil {
		_ = os.Remove(tempPath) // Ignore cleanup errors

		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// getDataDir resolves the data directory path.
func getDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".config", "slot"), nil
}

// GetLogPath returns the audit log file path.
func (store *Store) GetLogPath() string {
	return store.LogPath
}
