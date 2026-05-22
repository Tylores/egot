package handler

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/sep"
)

func TestGETTimeCSIP(t *testing.T) {
	repo := store.New(":memory:")
	repo.Load()
	defer repo.Close()
	reg := registry.New(":memory:")
	reg.Load()

	h := NewHandler(repo, reg)

	// Mock TLS and LFDI
	cert := &x509.Certificate{Raw: []byte("test-cert")}
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	reg.RegisterLFDI(lfdi, "123", "Test Device")

	req := httptest.NewRequest("GET", "/tm", nil)
	req.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}

	w := httptest.NewRecorder()
	h.GETTime(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp sep.Time
	if err := xml.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Quality != 7 {
		t.Errorf("Expected Quality = 7, got %v", resp.Quality)
	}

	now := time.Now().Unix()
	if int64(resp.CurrentTime) < now-5 || int64(resp.CurrentTime) > now+5 {
		t.Errorf("Expected CurrentTime close to %d, got %d", now, resp.CurrentTime)
	}
}
