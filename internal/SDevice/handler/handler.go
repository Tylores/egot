package handler

import (
"crypto/sha256"
"encoding/xml"
"fmt"
"net/http"

"github.com/Tylores/egot/internal/SDevice/repository/memory"
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
cert := req.TLS.PeerCertificates[0]
lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
if _, err := h.repo.GetEntity(lfdi); err != nil {
return "", err
}
return lfdi, nil
}

// SelfDevice resource handlers

func (h *Handler) GETSelfDevice(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
xml.NewEncoder(w).Encode(&sep.SelfDevice{})
}

func (h *Handler) HEADSelfDevice(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTSelfDevice(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.Header().Set("location", "/sdev/1")
w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTSelfDevice(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETESelfDevice(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}
