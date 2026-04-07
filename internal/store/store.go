// Package store provides a simple persistent key/value store backed by
// map[string]any with a sync.RWMutex for concurrent access and gob encoding
// for persistence.
//
// Concrete types stored as any must be registered by the caller with
// gob.Register before calling Load or Save.
package store

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store is a thread-safe key/value store that can be saved to and loaded from
// a file using gob encoding.
type Store struct {
	mu   sync.RWMutex
	data map[string]any
	path string
}

// New creates a new Store whose persistence file is at path. The file is not
// loaded automatically; call Load to populate the store from an existing file.
func New(path string) *Store {
	return &Store{
		data: make(map[string]any),
		path: path,
	}
}

// Get returns the value for key and whether it was found.
func (s *Store) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Set stores val under key, replacing any existing value.
func (s *Store) Set(key string, val any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
}

// Delete removes key from the store. It is a no-op if the key does not exist.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Keys returns a snapshot of all keys currently in the store.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Save gob-encodes the store to its file atomically (write to a temp file then
// rename). Parent directories are created if they do not exist.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("store: mkdir: %w", err)
	}

	tmp := s.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("store: create temp: %w", err)
	}

	if err := gob.NewEncoder(f).Encode(s.data); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("store: encode: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("store: close temp: %w", err)
	}

	if err := os.Rename(tmp, s.path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("store: rename: %w", err)
	}
	return nil
}

// Load gob-decodes the store from its file. If the file does not exist, Load
// returns nil without modifying the store.
func (s *Store) Load() error {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("store: open: %w", err)
	}
	defer f.Close()

	var data map[string]any
	if err := gob.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("store: decode: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = data
	return nil
}
