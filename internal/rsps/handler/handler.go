package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/rsps/repository/memory"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}


func (h *Handler) GETResponseSetList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADResponseSetList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTResponseSetList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTResponseSetList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEResponseSetList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETResponseSet(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADResponseSet(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTResponseSet(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTResponseSet(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEResponseSet(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETResponseList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADResponseList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTResponseList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTResponseList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEResponseList(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETPriceResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADPriceResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTPriceResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTPriceResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEPriceResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETTextResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADTextResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTTextResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTTextResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETETextResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEDERControlResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GETDrResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HEADDrResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PUTDrResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) POSTDrResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DELETEDrResponse(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ , err = strconv.Atoi(req.PathValue("id1"))
	if err != nil {
		log.Printf("path id1 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ , err = strconv.Atoi(req.PathValue("id2"))
	if err != nil {
		log.Printf("path id2 value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

