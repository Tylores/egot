package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/flow/repository/memory"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) GetFlowReservationRequests(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	frql, err := h.repo.GetFlowReservationRequests(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(frql)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}

}

func (h *Handler) GetFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		log.Printf("path id value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	frq, err := h.repo.GetFlowReservationRequest(lfdi, memory.Entity(id))
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(frq)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}

}
