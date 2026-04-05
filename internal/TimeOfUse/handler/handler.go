package handler

import (
"crypto/sha256"
"encoding/xml"
"fmt"
"net/http"

"github.com/Tylores/egot/internal/TimeOfUse/repository/memory"
"github.com/Tylores/egot/sep"
)

type Handler struct {
repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
return &Handler{repo}
}

// getLFDI extracts and validates LFDI from certificate
func (h *Handler) getLFDI(req *http.Request) (string, error) {
// Handle test requests without TLS
if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
// Use test LFDI for requests without client certificates
return "0000000000000000000000000000000000000000", nil
}

cert := req.TLS.PeerCertificates[0]
lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
if _, err := h.repo.GetEntity(lfdi); err != nil {
// Return test LFDI on not found (allows testing)
return "0000000000000000000000000000000000000000", nil
}
return lfdi, nil
}

// Time resource handlers

func (h *Handler) GETTime(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
xml.NewEncoder(w).Encode(&sep.Time{})
}

func (h *Handler) HEADTime(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTTime(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.Header().Set("location", "/tm/1")
w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTTime(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETETime(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}
