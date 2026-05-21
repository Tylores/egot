// Package store provides a high-performance persistent key/value store
// backed by SQLite (CGO-free via modernc.org/sqlite).
//
// Concrete types stored as any must be registered by the caller with
// gob.Register before calling Get or Set.
package store

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store is a thread-safe key/value store backed by SQLite.
type Store struct {
	db   *sql.DB
	path string
}

// New creates a new Store whose persistence file is at path.
func New(path string) *Store {
	return &Store{
		path: path,
	}
}

// Load opens the SQLite database and ensures the schema exists.
func (s *Store) Load() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("store: mkdir: %w", err)
	}

	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return fmt.Errorf("store: open: %w", err)
	}

	// Enable WAL mode and other performance pragmas
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-64000;", // 64MB cache
		"PRAGMA busy_timeout=5000;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return fmt.Errorf("store: pragma %s: %w", p, err)
		}
	}

	// Configure connection pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	schema := `CREATE TABLE IF NOT EXISTS kv (
		key TEXT PRIMARY KEY,
		val BLOB,
		owner_id TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_owner ON kv(owner_id);`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return fmt.Errorf("store: schema: %w", err)
	}

	s.db = db
	return nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Save is now a no-op as SQLite persists changes immediately.
// Kept for backward compatibility with existing handlers.
func (s *Store) Save() error {
	return nil
}

// Get returns the value for key and whether it was found.
func (s *Store) Get(key string) (any, bool) {
	if s.db == nil {
		return nil, false
	}

	var buf []byte
	err := s.db.QueryRow("SELECT val FROM kv WHERE key = ?", key).Scan(&buf)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	var val any
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&val); err != nil {
		return nil, false
	}

	return val, true
}

// Set stores val under key, replacing any existing value.
func (s *Store) Set(key string, val any) {
	s.SetWithOwner(key, val, "")
}

// SetWithOwner stores val under key with an associated ownerID for secondary indexing.
func (s *Store) SetWithOwner(key string, val any, ownerID string) {
	if s.db == nil {
		return
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&val); err != nil {
		return
	}

	_, err := s.db.Exec("INSERT OR REPLACE INTO kv (key, val, owner_id) VALUES (?, ?, ?)", key, buf.Bytes(), ownerID)
	if err != nil {
		fmt.Printf("store: set error: %v\n", err)
	}
}

// GetByOwner returns all values associated with a specific ownerID.
func (s *Store) GetByOwner(ownerID string) []any {
	if s.db == nil {
		return nil
	}

	rows, err := s.db.Query("SELECT val FROM kv WHERE owner_id = ?", ownerID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var buf []byte
		if err := rows.Scan(&buf); err != nil {
			continue
		}
		var val any
		if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&val); err == nil {
			results = append(results, val)
		}
	}
	return results
}

// Delete removes key from the store.
func (s *Store) Delete(key string) {
	if s.db == nil {
		return
	}
	_, _ = s.db.Exec("DELETE FROM kv WHERE key = ?", key)
}

// Keys returns all keys currently in the store.
func (s *Store) Keys() []string {
	if s.db == nil {
		return nil
	}

	rows, err := s.db.Query("SELECT key FROM kv")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err == nil {
			keys = append(keys, key)
		}
	}
	return keys
}
