package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/core/repository/memory"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) GetDeviceCapability(w http.ResponseWriter, req *http.Request) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(h.repo.GetDeviceCapability())
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}

}

func (h *Handler) GetTime(w http.ResponseWriter, req *http.Request) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(h.repo.GetTime())
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}
}

func (h *Handler) GetEndDevices(w http.ResponseWriter, req *http.Request) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	e, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	list := sep.NewEndDeviceList()
	list.AllAttr = 1
	list.ResultsAttr = 1
	edev := h.repo.GetEndDevice(e)
	list.EndDevice = append(list.EndDevice, &edev)

	// query := req.URL.Query()
	// s := query.Get("s")
	// a := query.Get("a")
	// l := query.Get("l")
	// currently only aggregators would have an actual list of values
	// since we arn't working through aggregators yet we ignore queries
	// becuase there is only a single value per client anyways

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}
}

func (h *Handler) GetEndDevice(w http.ResponseWriter, req *http.Request) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	e, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		log.Printf("path id value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if memory.Entity(id) != e {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(h.repo.GetEndDevice(e))
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}
}

func (h *Handler) GetRegistration(w http.ResponseWriter, req *http.Request) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	e, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		log.Printf("path id value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if memory.Entity(id) != e {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	err = xml.NewEncoder(w).Encode(h.repo.GetRegistration(e))
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		return
	}
}
