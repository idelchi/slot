// Package store provides persistent storage operations for slot data with atomic file operations.
package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/slot/internal/slot"
)

type slotsFile struct {
	Include []string `json:"include,omitempty"`
	Slots   slot.Slots
}

type includeStack struct {
	stores []Store
}

// Store handles persistent storage operations for slot data.
type Store string

// Path returns the file path of the slot store.
func (store Store) Path() string {
	return string(store)
}

// New creates a new Store instance from the given file path.
func New(slotsFile string) (Store, error) {
	store := Store(slotsFile)

	dataDir := filepath.Dir(slotsFile)

	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return store, fmt.Errorf("creating data directory: %w", err)
	}

	return store, nil
}

// Load reads slots from disk and recursively included files.
func (store Store) Load() (slot.Slots, error) {
	store, err := store.clean()
	if err != nil {
		return nil, err
	}

	slots, err := store.load(true, includeStack{}, map[Store]bool{})
	if err != nil {
		return nil, err
	}

	return slots.Unique(), nil
}

// LoadLocal reads only the slots directly defined in this store.
func (store Store) LoadLocal() (slot.Slots, error) {
	store, err := store.clean()
	if err != nil {
		return nil, err
	}

	file, err := store.read(true)
	if err != nil {
		return nil, err
	}

	return file.Slots, nil
}

// Delete removes the visible slot with the given name from the file that defines it.
func (store Store) Delete(name string) (bool, error) {
	store, err := store.clean()
	if err != nil {
		return false, err
	}

	foundStore, found, err := store.find(name, true, includeStack{}, map[Store]bool{})
	if err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	file, err := foundStore.read(false)
	if err != nil {
		return false, err
	}

	file.Slots.Delete(name)

	if err := foundStore.write(file); err != nil {
		return false, err
	}

	return true, nil
}

// Save writes the slots to disk.
func (store Store) Save(slots slot.Slots) error {
	store, err := store.clean()
	if err != nil {
		return err
	}

	file, err := store.read(true)
	if err != nil {
		return err
	}

	file.Slots = slots

	return store.write(file)
}

// load reads the store's slots and recursively includes its dependencies.
func (store Store) load(allowMissing bool, stack includeStack, visited map[Store]bool) (slot.Slots, error) {
	if slices.Contains(stack.stores, store) {
		return nil, fmt.Errorf("recursive include: %s", stack.formatCycle(store))
	}

	if visited[store] {
		return nil, nil
	}

	visited[store] = true

	file, err := store.read(allowMissing)
	if err != nil {
		return nil, err
	}

	stack.stores = append(stack.stores, store)

	slots := append(slot.Slots{}, file.Slots...)

	for _, include := range file.Include {
		includeStore, err := store.resolveInclude(include)
		if err != nil {
			return nil, err
		}

		includedSlots, err := includeStore.load(false, stack, visited)
		if err != nil {
			return nil, err
		}

		slots = append(slots, includedSlots...)
	}

	return slots, nil
}

// find returns the store that defines the visible slot with name.
func (store Store) find(
	name string,
	allowMissing bool,
	stack includeStack,
	visited map[Store]bool,
) (Store, bool, error) {
	if slices.Contains(stack.stores, store) {
		return "", false, fmt.Errorf("recursive include: %s", stack.formatCycle(store))
	}

	if visited[store] {
		return "", false, nil
	}

	visited[store] = true

	file, err := store.read(allowMissing)
	if err != nil {
		return "", false, err
	}

	if file.Slots.Exists(name) {
		return store, true, nil
	}

	stack.stores = append(stack.stores, store)

	for _, include := range file.Include {
		includeStore, err := store.resolveInclude(include)
		if err != nil {
			return "", false, err
		}

		foundStore, found, err := includeStore.find(name, false, stack, visited)
		if err != nil || found {
			return foundStore, found, err
		}
	}

	return "", false, nil
}

// read reads one slots file from disk.
func (store Store) read(allowMissing bool) (slotsFile, error) {
	file := slotsFile{}

	data, err := os.ReadFile(store.Path())
	if err != nil {
		if allowMissing && os.IsNotExist(err) {
			return file, nil
		}

		return file, fmt.Errorf("reading slots file %q: %w", filepath.ToSlash(store.Path()), err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return file, nil
	}

	if err := yaml.Unmarshal(data, &file); err != nil {
		return file, fmt.Errorf("unmarshalling slots file %q: %w", filepath.ToSlash(store.Path()), err)
	}

	return file, nil
}

// write writes one slots file to disk.
func (store Store) write(file slotsFile) error {
	data, err := yaml.MarshalWithOptions(
		file,
		yaml.IndentSequence(true),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	if err != nil {
		return fmt.Errorf("marshalling slots: %w", err)
	}

	data = bytes.TrimRight(data, "\n")

	if err := os.WriteFile(store.Path(), data, 0o600); err != nil {
		return fmt.Errorf("writing file %q: %w", filepath.ToSlash(store.Path()), err)
	}

	return nil
}

// resolveInclude returns the store for an include declared by this store.
func (store Store) resolveInclude(includePath string) (Store, error) {
	if includePath == "" {
		return "", fmt.Errorf("empty include in %q", filepath.ToSlash(store.Path()))
	}

	if !filepath.IsAbs(includePath) {
		includePath = filepath.Join(filepath.Dir(store.Path()), includePath)
	}

	return Store(includePath).clean()
}

// clean returns the absolute, cleaned store path.
func (store Store) clean() (Store, error) {
	absolute, err := filepath.Abs(store.Path())
	if err != nil {
		return "", fmt.Errorf("resolving slots file path %q: %w", store.Path(), err)
	}

	return Store(filepath.Clean(absolute)), nil
}

// formatCycle formats the recursive include path for error messages.
func (stack includeStack) formatCycle(repeated Store) string {
	start := slices.Index(stack.stores, repeated)
	if start == -1 {
		start = 0
	}

	cycle := append(append([]Store{}, stack.stores[start:]...), repeated)
	paths := make([]string, len(cycle))

	for i, store := range cycle {
		paths[i] = filepath.ToSlash(store.Path())
	}

	return strings.Join(paths, " -> ")
}

// DefaultSlotsFile returns the full path to the default slots file location.
func DefaultSlotsFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "slots.yaml", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.ToSlash(filepath.Join(home, ".config", "slot", "slots.yaml")), nil
}
