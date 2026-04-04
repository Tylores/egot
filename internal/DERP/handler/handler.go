package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/DERP/repository/memory"
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

// DERProgramList resource handlers

func (h *Handler) DELETEDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERProgramList{})
}

func (h *Handler) HEADDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/derp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERProgram resource handlers

func (h *Handler) DELETEDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERProgram{})
}

func (h *Handler) GETDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERProgram{})
}

func (h *Handler) HEADDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ActiveDERControlList resource handlers

func (h *Handler) DELETEActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERControlList{})
}

func (h *Handler) HEADActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/derp/{id1}/actderc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERControlList resource handlers

func (h *Handler) DELETEDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERControlList{})
}

func (h *Handler) HEADDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/derp/{id1}/derc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERControl resource handlers

func (h *Handler) DELETEDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERControl{})
}

func (h *Handler) GETDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERControl{})
}

func (h *Handler) HEADDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DefaultDERControl resource handlers

func (h *Handler) DELETEDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DefaultDERControl{})
}

func (h *Handler) HEADDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/derp/{id1}/dderc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERCurveList resource handlers

func (h *Handler) DELETEDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERCurveList{})
}

func (h *Handler) HEADDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/derp/{id1}/dc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERCurve resource handlers

func (h *Handler) DELETEDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERCurve{})
}

func (h *Handler) GETDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERCurve{})
}

func (h *Handler) HEADDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

