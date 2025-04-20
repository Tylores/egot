package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/repository/memory"
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

	e, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = xml.NewEncoder(w).Encode(h.repo.GetDeviceCapability(e))
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
