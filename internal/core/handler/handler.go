package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/repository/memory"
	"github.com/Tylores/egot/internal/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) GetDeviceCapability(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = xml.NewEncoder(w).Encode(h.repo.GetDeviceCapability())
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
}

func (h *Handler) HeadDeviceCapability(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	xml_bytes, err := xml.Marshal(h.repo.GetDeviceCapability())
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("ContentLength", string(len(xml_bytes)))
}
