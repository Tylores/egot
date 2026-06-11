package store_test

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/Tylores/egot/internal/store"
)

func init() {
	// Register a concrete type used in persistence tests.
	gob.Register(item{})
}

type item struct {
	Name  string
	Value int
}

func TestGetSetDelete(t *testing.T) {
	s := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if _, ok := s.Get("missing"); ok {
		t.Fatal("expected missing key to not be found")
	}

	s.Set("a", 42)
	v, ok := s.Get("a")
	if !ok {
		t.Fatal("expected key 'a' to exist")
	}
	if v.(int) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}

	s.Delete("a")
	if _, ok := s.Get("a"); ok {
		t.Fatal("expected deleted key to be absent")
	}

	// Delete of non-existent key should be a no-op.
	s.Delete("missing")
}

func TestKeys(t *testing.T) {
	s := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	s.Set("b", 1)
	s.Set("a", 2)
	s.Set("c", 3)

	keys := s.Keys()
	sort.Strings(keys)

	want := []string{"a", "b", "c"}
	for i, k := range keys {
		if k != want[i] {
			t.Fatalf("keys mismatch: got %v, want %v", keys, want)
		}
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	s1 := store.New(path)
	if err := s1.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	s1.Set("x", item{Name: "hello", Value: 7})
	s1.Set("y", item{Name: "world", Value: 99})

	if err := s1.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2 := store.New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	v, ok := s2.Get("x")
	if !ok {
		t.Fatal("expected key 'x' after load")
	}
	got := v.(item)
	if got.Name != "hello" || got.Value != 7 {
		t.Fatalf("unexpected value: %+v", got)
	}

	v, ok = s2.Get("y")
	if !ok {
		t.Fatal("expected key 'y' after load")
	}
	got = v.(item)
	if got.Name != "world" || got.Value != 99 {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestLoadMissingFile(t *testing.T) {
	s := store.New(filepath.Join(t.TempDir(), "nonexistent.db"))
	if err := s.Load(); err != nil {
		t.Fatalf("Load on missing file should return nil, got: %v", err)
	}
}

func TestSaveCreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "test.db")
	s := store.New(path)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	s.Set("k", item{Name: "v"})
	if err := s.Save(); err != nil {
		t.Fatalf("Save with missing parent dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestAtomicSave(t *testing.T) {
	// After Save, no .tmp file should remain.
	path := filepath.Join(t.TempDir(), "test.db")
	s := store.New(path)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	s.Set("k", 1)
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp file should not exist after successful Save")
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := "k"
			s.Set(key, n)
			s.Get(key)
			s.Keys()
		}(i)
	}
	wg.Wait()
}
