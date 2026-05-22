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

	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/sep"
)

func TestGETRegistrationCSIP(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/edev/123/rg", nil)
	req.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}
	req.SetPathValue("id1", "123")

	w := httptest.NewRecorder()
	h.GETRegistration(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp sep.Registration
	if err := xml.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.PIN == nil {
		t.Fatalf("Expected PIN to be not nil")
	}

	if *resp.PIN != 111115 {
		t.Errorf("Expected PIN = 111115, got %v", *resp.PIN)
	}
}
