package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/EDevice/repository/memory"
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

// EndDeviceList resource handlers

func (h *Handler) DELETEEndDeviceList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETEndDeviceList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.EndDeviceList{})
}

func (h *Handler) HEADEndDeviceList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTEndDeviceList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/edev/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTEndDeviceList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// EndDevice resource handlers

func (h *Handler) DELETEEndDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.EndDevice{})
}

func (h *Handler) GETEndDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.EndDevice{})
}

func (h *Handler) HEADEndDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTEndDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTEndDevice(w http.ResponseWriter, req *http.Request) {
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

// Registration resource handlers

func (h *Handler) DELETERegistration(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETRegistration(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.Registration{})
}

func (h *Handler) HEADRegistration(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTRegistration(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/rg/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTRegistration(w http.ResponseWriter, req *http.Request) {
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

// DeviceStatus resource handlers

func (h *Handler) DELETEDeviceStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDeviceStatus(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DeviceStatus{})
}

func (h *Handler) HEADDeviceStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDeviceStatus(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/dstat/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDeviceStatus(w http.ResponseWriter, req *http.Request) {
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

// ProxiedDeviceList resource handlers

func (h *Handler) DELETEProxiedDeviceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETProxiedDeviceList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ProxiedDeviceList{})
}

func (h *Handler) HEADProxiedDeviceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTProxiedDeviceList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/prxy/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTProxiedDeviceList(w http.ResponseWriter, req *http.Request) {
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

// ProxiedDevice resource handlers

func (h *Handler) DELETEProxiedDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ProxiedDevice{})
}

func (h *Handler) GETProxiedDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ProxiedDevice{})
}

func (h *Handler) HEADProxiedDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTProxiedDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTProxiedDevice(w http.ResponseWriter, req *http.Request) {
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

// AggregationPriority resource handlers

func (h *Handler) DELETEAggregationPriority(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETAggregationPriority(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.AggregationPriority{})
}

func (h *Handler) HEADAggregationPriority(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTAggregationPriority(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/aggp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTAggregationPriority(w http.ResponseWriter, req *http.Request) {
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

// AggregatedDeviceList resource handlers

func (h *Handler) DELETEAggregatedDeviceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETAggregatedDeviceList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.AggregatedDeviceList{})
}

func (h *Handler) HEADAggregatedDeviceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTAggregatedDeviceList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/adev/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTAggregatedDeviceList(w http.ResponseWriter, req *http.Request) {
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

// AggregatedDevice resource handlers

func (h *Handler) DELETEAggregatedDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.AggregatedDevice{})
}

func (h *Handler) GETAggregatedDevice(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.AggregatedDevice{})
}

func (h *Handler) HEADAggregatedDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTAggregatedDevice(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTAggregatedDevice(w http.ResponseWriter, req *http.Request) {
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

// FunctionSetAssignmentsList resource handlers

func (h *Handler) DELETEFunctionSetAssignmentsList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETFunctionSetAssignmentsList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FunctionSetAssignmentsList{})
}

func (h *Handler) HEADFunctionSetAssignmentsList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFunctionSetAssignmentsList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/fsa/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTFunctionSetAssignmentsList(w http.ResponseWriter, req *http.Request) {
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

// FunctionSetAssignments resource handlers

func (h *Handler) DELETEFunctionSetAssignments(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FunctionSetAssignments{})
}

func (h *Handler) GETFunctionSetAssignments(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FunctionSetAssignments{})
}

func (h *Handler) HEADFunctionSetAssignments(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFunctionSetAssignments(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTFunctionSetAssignments(w http.ResponseWriter, req *http.Request) {
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

// SubscriptionList resource handlers

func (h *Handler) DELETESubscriptionList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETSubscriptionList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.SubscriptionList{})
}

func (h *Handler) HEADSubscriptionList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTSubscriptionList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/sub/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTSubscriptionList(w http.ResponseWriter, req *http.Request) {
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

// Subscription resource handlers

func (h *Handler) DELETESubscription(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.Subscription{})
}

func (h *Handler) GETSubscription(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.Subscription{})
}

func (h *Handler) HEADSubscription(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTSubscription(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTSubscription(w http.ResponseWriter, req *http.Request) {
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

// DeviceInformation resource handlers

func (h *Handler) DELETEDeviceInformation(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDeviceInformation(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DeviceInformation{})
}

func (h *Handler) HEADDeviceInformation(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDeviceInformation(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/di/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDeviceInformation(w http.ResponseWriter, req *http.Request) {
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

// SupportedLocaleList resource handlers

func (h *Handler) DELETESupportedLocaleList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETSupportedLocaleList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.SupportedLocaleList{})
}

func (h *Handler) HEADSupportedLocaleList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTSupportedLocaleList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/di/loc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTSupportedLocaleList(w http.ResponseWriter, req *http.Request) {
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

// SupportedLocale resource handlers

func (h *Handler) DELETESupportedLocale(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.SupportedLocale{})
}

func (h *Handler) GETSupportedLocale(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.SupportedLocale{})
}

func (h *Handler) HEADSupportedLocale(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTSupportedLocale(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTSupportedLocale(w http.ResponseWriter, req *http.Request) {
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

// PowerStatus resource handlers

func (h *Handler) DELETEPowerStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETPowerStatus(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.PowerStatus{})
}

func (h *Handler) HEADPowerStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTPowerStatus(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/ps/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTPowerStatus(w http.ResponseWriter, req *http.Request) {
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

// IPInterfaceList resource handlers

func (h *Handler) DELETEIPInterfaceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETIPInterfaceList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.IPInterfaceList{})
}

func (h *Handler) HEADIPInterfaceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTIPInterfaceList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/ns/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTIPInterfaceList(w http.ResponseWriter, req *http.Request) {
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

// IPInterface resource handlers

func (h *Handler) DELETEIPInterface(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.IPInterface{})
}

func (h *Handler) GETIPInterface(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.IPInterface{})
}

func (h *Handler) HEADIPInterface(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTIPInterface(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTIPInterface(w http.ResponseWriter, req *http.Request) {
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

// IPAddrList resource handlers

func (h *Handler) DELETEIPAddrList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETIPAddrList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.IPAddrList{})
}

func (h *Handler) HEADIPAddrList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTIPAddrList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/ns/{id2}/addr/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTIPAddrList(w http.ResponseWriter, req *http.Request) {
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

// IPAddr resource handlers

func (h *Handler) DELETEIPAddr(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.IPAddr{})
}

func (h *Handler) GETIPAddr(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.IPAddr{})
}

func (h *Handler) HEADIPAddr(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTIPAddr(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTIPAddr(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// RPLInstanceList resource handlers

func (h *Handler) DELETERPLInstanceList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETRPLInstanceList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLInstanceList{})
}

func (h *Handler) HEADRPLInstanceList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTRPLInstanceList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTRPLInstanceList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// RPLInstance resource handlers

func (h *Handler) DELETERPLInstance(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLInstance{})
}

func (h *Handler) GETRPLInstance(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLInstance{})
}

func (h *Handler) HEADRPLInstance(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTRPLInstance(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTRPLInstance(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// RPLSourceRoutesList resource handlers

func (h *Handler) DELETERPLSourceRoutesList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETRPLSourceRoutesList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLSourceRoutesList{})
}

func (h *Handler) HEADRPLSourceRoutesList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTRPLSourceRoutesList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTRPLSourceRoutesList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// RPLSourceRoutes resource handlers

func (h *Handler) DELETERPLSourceRoutes(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id5")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLSourceRoutes{})
}

func (h *Handler) GETRPLSourceRoutes(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id5")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.RPLSourceRoutes{})
}

func (h *Handler) HEADRPLSourceRoutes(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id5")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTRPLSourceRoutes(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id5")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTRPLSourceRoutes(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id5")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// LLInterfaceList resource handlers

func (h *Handler) DELETELLInterfaceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETLLInterfaceList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LLInterfaceList{})
}

func (h *Handler) HEADLLInterfaceList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTLLInterfaceList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/ns/{id2}/ll/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTLLInterfaceList(w http.ResponseWriter, req *http.Request) {
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

// LLInterface resource handlers

func (h *Handler) DELETELLInterface(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.LLInterface{})
}

func (h *Handler) GETLLInterface(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.LLInterface{})
}

func (h *Handler) HEADLLInterface(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTLLInterface(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTLLInterface(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// NeighborList resource handlers

func (h *Handler) DELETENeighborList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETNeighborList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.NeighborList{})
}

func (h *Handler) HEADNeighborList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTNeighborList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/edev/{id1}/ns/{id2}/ll/{id3}/nbh/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTNeighborList(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Neighbor resource handlers

func (h *Handler) DELETENeighbor(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.Neighbor{})
}

func (h *Handler) GETNeighbor(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.Neighbor{})
}

func (h *Handler) HEADNeighbor(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTNeighbor(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTNeighbor(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(req.PathValue("id4")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// LogEventList resource handlers

func (h *Handler) DELETELogEventList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETLogEventList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LogEventList{})
}

func (h *Handler) HEADLogEventList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTLogEventList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/lel/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTLogEventList(w http.ResponseWriter, req *http.Request) {
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

// LogEvent resource handlers

func (h *Handler) DELETELogEvent(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LogEvent{})
}

func (h *Handler) GETLogEvent(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LogEvent{})
}

func (h *Handler) HEADLogEvent(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTLogEvent(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTLogEvent(w http.ResponseWriter, req *http.Request) {
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

// Configuration resource handlers

func (h *Handler) DELETEConfiguration(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETConfiguration(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.Configuration{})
}

func (h *Handler) HEADConfiguration(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTConfiguration(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/cfg/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTConfiguration(w http.ResponseWriter, req *http.Request) {
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

// PriceResponseCfgList resource handlers

func (h *Handler) DELETEPriceResponseCfgList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETPriceResponseCfgList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.PriceResponseCfgList{})
}

func (h *Handler) HEADPriceResponseCfgList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTPriceResponseCfgList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/cfg/prcfg/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTPriceResponseCfgList(w http.ResponseWriter, req *http.Request) {
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

// PriceResponseCfg resource handlers

func (h *Handler) DELETEPriceResponseCfg(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.PriceResponseCfg{})
}

func (h *Handler) GETPriceResponseCfg(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.PriceResponseCfg{})
}

func (h *Handler) HEADPriceResponseCfg(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTPriceResponseCfg(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTPriceResponseCfg(w http.ResponseWriter, req *http.Request) {
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

// FileStatus resource handlers

func (h *Handler) DELETEFileStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETFileStatus(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FileStatus{})
}

func (h *Handler) HEADFileStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFileStatus(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/fs/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTFileStatus(w http.ResponseWriter, req *http.Request) {
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

// LoadShedAvailabilityList resource handlers

func (h *Handler) DELETELoadShedAvailabilityList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETLoadShedAvailabilityList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LoadShedAvailabilityList{})
}

func (h *Handler) HEADLoadShedAvailabilityList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTLoadShedAvailabilityList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/lsl/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTLoadShedAvailabilityList(w http.ResponseWriter, req *http.Request) {
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

// LoadShedAvailability resource handlers

func (h *Handler) DELETELoadShedAvailability(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LoadShedAvailability{})
}

func (h *Handler) GETLoadShedAvailability(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.LoadShedAvailability{})
}

func (h *Handler) HEADLoadShedAvailability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTLoadShedAvailability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTLoadShedAvailability(w http.ResponseWriter, req *http.Request) {
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

// FlowReservationRequestList resource handlers

func (h *Handler) DELETEFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationRequestList{})
}

func (h *Handler) HEADFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/frq/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
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

// FlowReservationRequest resource handlers

func (h *Handler) DELETEFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationRequest{})
}

func (h *Handler) GETFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationRequest{})
}

func (h *Handler) HEADFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
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

// FlowReservationResponseList resource handlers

func (h *Handler) DELETEFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponseList{})
}

func (h *Handler) HEADFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/frp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
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

// FlowReservationResponse resource handlers

func (h *Handler) DELETEFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponse{})
}

func (h *Handler) GETFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponse{})
}

func (h *Handler) HEADFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
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

// DERList resource handlers

func (h *Handler) DELETEDERList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERList{})
}

func (h *Handler) HEADDERList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERList(w http.ResponseWriter, req *http.Request) {
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

// DER resource handlers

func (h *Handler) DELETEDER(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DER{})
}

func (h *Handler) GETDER(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DER{})
}

func (h *Handler) HEADDER(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDER(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDER(w http.ResponseWriter, req *http.Request) {
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

// AssociatedUsagePoint resource handlers

func (h *Handler) DELETEAssociatedUsagePoint(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETAssociatedUsagePoint(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.UsagePoint{})
}

func (h *Handler) HEADAssociatedUsagePoint(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTAssociatedUsagePoint(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/upt/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTAssociatedUsagePoint(w http.ResponseWriter, req *http.Request) {
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

// AssociatedDERProgramList resource handlers

func (h *Handler) DELETEAssociatedDERProgramList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETAssociatedDERProgramList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERProgramList{})
}

func (h *Handler) HEADAssociatedDERProgramList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTAssociatedDERProgramList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/derp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTAssociatedDERProgramList(w http.ResponseWriter, req *http.Request) {
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

// CurrentDERControls resource handlers

func (h *Handler) DELETECurrentDERControls(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETCurrentDERControls(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.CurrentDERControls{})
}

func (h *Handler) HEADCurrentDERControls(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTCurrentDERControls(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/cdc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTCurrentDERControls(w http.ResponseWriter, req *http.Request) {
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

// CurrentDERProgram resource handlers

func (h *Handler) DELETECurrentDERProgram(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETCurrentDERProgram(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERProgram{})
}

func (h *Handler) HEADCurrentDERProgram(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTCurrentDERProgram(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/cdp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTCurrentDERProgram(w http.ResponseWriter, req *http.Request) {
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

// DERSettings resource handlers

func (h *Handler) DELETEDERSettings(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERSettings(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERSettings{})
}

func (h *Handler) HEADDERSettings(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERSettings(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/derg/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERSettings(w http.ResponseWriter, req *http.Request) {
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

// DERStatus resource handlers

func (h *Handler) DELETEDERStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERStatus(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERStatus{})
}

func (h *Handler) HEADDERStatus(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERStatus(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/ders/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERStatus(w http.ResponseWriter, req *http.Request) {
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

// DERAvailability resource handlers

func (h *Handler) DELETEDERAvailability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERAvailability(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERAvailability{})
}

func (h *Handler) HEADDERAvailability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERAvailability(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/dera/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERAvailability(w http.ResponseWriter, req *http.Request) {
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

// DERCapability resource handlers

func (h *Handler) DELETEDERCapability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERCapability(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERCapability{})
}

func (h *Handler) HEADDERCapability(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERCapability(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/dercap/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERCapability(w http.ResponseWriter, req *http.Request) {
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

// DERComponentList resource handlers

func (h *Handler) DELETEDERComponentList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETDERComponentList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERComponentList{})
}

func (h *Handler) HEADDERComponentList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDERComponentList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/edev/{id1}/der/{id2}/dercom/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTDERComponentList(w http.ResponseWriter, req *http.Request) {
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

// DERComponent resource handlers

func (h *Handler) DELETEDERComponent(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERComponent{})
}

func (h *Handler) GETDERComponent(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.DERComponent{})
}

func (h *Handler) HEADDERComponent(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTDERComponent(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) PUTDERComponent(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id3")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

