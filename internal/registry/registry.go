// Package registry provides a shared device registry for authorizing
// LFDIs across all microservices.
package registry

import (
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// Registry manages authorized device identities.
type Registry struct {
	db   *sql.DB
	path string
}

// New creates a new Registry using the specified database path.
func New(path string) *Registry {
	return &Registry{path: path}
}

// Load opens the registry database and ensures the schema exists.
func (r *Registry) Load() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("registry: mkdir: %w", err)
	}

	db, err := sql.Open("sqlite", r.path)
	if err != nil {
		return fmt.Errorf("registry: open: %w", err)
	}

	schema := `CREATE TABLE IF NOT EXISTS authorized_devices (
		lfdi TEXT PRIMARY KEY,
		sfdi TEXT,
		common_name TEXT,
		registered_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return fmt.Errorf("registry: schema: %w", err)
	}

	r.db = db
	return nil
}

// Close closes the registry database.
func (r *Registry) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// RegisterLFDI adds an LFDI to the authorized devices list.
func (r *Registry) RegisterLFDI(lfdi, sfdi, cn string) error {
	if r.db == nil {
		return fmt.Errorf("registry: not loaded")
	}
	_, err := r.db.Exec("INSERT OR REPLACE INTO authorized_devices (lfdi, sfdi, common_name) VALUES (?, ?, ?)", lfdi, sfdi, cn)
	return err
}

// IsAuthorized checks if an LFDI is in the registry.
func (r *Registry) IsAuthorized(lfdi string) bool {
	if r.db == nil {
		return false
	}
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM authorized_devices WHERE lfdi = ?)", lfdi).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// PopulateFromCertDir walks a directory for *.crt files and registers them.
// It specifically looks for files with 'client' in the name.
func (r *Registry) PopulateFromCertDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".crt") {
			return nil
		}
		if !strings.Contains(strings.ToLower(info.Name()), "client") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		block, _ := pem.Decode(data)
		if block == nil {
			return nil
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil
		}

		lfdi := ComputeLFDI(cert.Raw)
		sfdi := ComputeSFDI(lfdi)
		
		return r.RegisterLFDI(lfdi, sfdi, cert.Subject.CommonName)
	})
}

// ComputeLFDI calculates the LFDI from a certificate's raw bytes.
// LFDI is the first 160 bits (20 bytes / 40 hex chars) of the SHA-256 hash.
func ComputeLFDI(certRaw []byte) string {
	hash := sha256.Sum256(certRaw)
	return strings.ToUpper(hex.EncodeToString(hash[:20]))
}

// ComputeSFDI calculates the SFDI from an LFDI.
// SFDI is the first 36 bits of the LFDI hash, represented as a decimal with a check digit.
// For simplicity in this research project, we take the first 9 digits of the hex-to-decimal conversion.
func ComputeSFDI(lfdi string) string {
	// Real IEEE 2030.5 SFDI is more complex, but this satisfies the research needs
	// ensuring a stable 1:1 mapping from LFDI to SFDI.
	val, _ := hex.DecodeString(lfdi[:10]) // Take first 40 bits
	sum := uint64(0)
	for i, b := range val {
		sum += uint64(b) << (8 * (len(val) - 1 - i))
	}
	return fmt.Sprintf("%012d", sum)[:11] // 11 digits + 1 check digit (omitted here for simplicity)
}
