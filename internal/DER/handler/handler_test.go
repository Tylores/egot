package handler

import (
	"bytes"
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

func TestDERCurveStorage(t *testing.T) {
	repo := store.New("test_der.db")
	repo.Load()
	defer repo.Close()
	reg := registry.New("test_reg.db")

	h := NewHandler(repo, reg)

	reg.Load()
	// Mock TLS and LFDI
	cert := &x509.Certificate{Raw: []byte("test-cert")}
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	reg.RegisterLFDI(lfdi, "123", "Test Device")

	mrid := "test-curve"
	curve := &sep.DERCurve{
		IdentifiedObject: &sep.IdentifiedObject{
			MRID: &sep.MRIDType{HexBinary128: &mrid},
		},
		CurveType: func() *sep.DERCurveType { v := sep.DERCurveType(11); return &v }(), // Volt-Var
		CurveData: []*sep.CurveData{
			{Xvalue: 95, Yvalue: 100},
			{Xvalue: 100, Yvalue: 0},
			{Xvalue: 105, Yvalue: -100},
		},
	}

	body, _ := xml.Marshal(curve)
	req := httptest.NewRequest("POST", "/derp/1/dc", bytes.NewBuffer(body))
	req.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}
	req.SetPathValue("id1", "1")

	w := httptest.NewRecorder()
	h.POSTDERCurveList(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", w.Code)
	}

	// Verify retrieval
	reqGet := httptest.NewRequest("GET", "/derp/1/dc/test-curve", nil)
	reqGet.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}
	reqGet.SetPathValue("id1", "1")
	reqGet.SetPathValue("id2", "test-curve")

	wGet := httptest.NewRecorder()
	h.GETDERCurve(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", wGet.Code)
	}

	var retrieved sep.DERCurve
	xml.NewDecoder(wGet.Body).Decode(&retrieved)

	if retrieved.CurveType == nil || *retrieved.CurveType != 11 {
		t.Errorf("Expected CurveType 11, got %v", retrieved.CurveType)
	}

	if len(retrieved.CurveData) != 3 {
		t.Errorf("Expected 3 CurveData points, got %d", len(retrieved.CurveData))
	}
}
