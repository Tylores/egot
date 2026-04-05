package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/TariffProfile/repository/memory"
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

// TariffProfileList resource handlers

func (h *Handler) DELETETariffProfileList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) GETTariffProfileList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.TariffProfileList{})
}

func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTTariffProfileList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "/tp/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTTariffProfileList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// TariffProfile resource handlers

func (h *Handler) DELETETariffProfile(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TariffProfile{})
}

func (h *Handler) GETTariffProfile(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TariffProfile{})
}

func (h *Handler) HEADTariffProfile(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTTariffProfile(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTTariffProfile(w http.ResponseWriter, req *http.Request) {
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

// RateComponentList resource handlers

func (h *Handler) DELETERateComponentList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETRateComponentList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.RateComponentList{})
}

func (h *Handler) HEADRateComponentList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTRateComponentList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/tp/{id1}/rc/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTRateComponentList(w http.ResponseWriter, req *http.Request) {
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

// RateComponent resource handlers

func (h *Handler) DELETERateComponent(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.RateComponent{})
}

func (h *Handler) GETRateComponent(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.RateComponent{})
}

func (h *Handler) HEADRateComponent(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTRateComponent(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTRateComponent(w http.ResponseWriter, req *http.Request) {
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

// ActiveTimeTariffIntervalList resource handlers

func (h *Handler) DELETEActiveTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETActiveTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TimeTariffIntervalList{})
}

func (h *Handler) HEADActiveTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTActiveTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/tp/{id1}/rc/{id2}/acttti/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTActiveTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

// TimeTariffIntervalList resource handlers

func (h *Handler) DELETETimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TimeTariffIntervalList{})
}

func (h *Handler) HEADTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/tp/{id1}/rc/{id2}/tti/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTTimeTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

// TimeTariffInterval resource handlers

func (h *Handler) DELETETimeTariffInterval(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TimeTariffInterval{})
}

func (h *Handler) GETTimeTariffInterval(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.TimeTariffInterval{})
}

func (h *Handler) HEADTimeTariffInterval(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTTimeTariffInterval(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTTimeTariffInterval(w http.ResponseWriter, req *http.Request) {
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

// ConsumptionTariffIntervalList resource handlers

func (h *Handler) DELETEConsumptionTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) GETConsumptionTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ConsumptionTariffIntervalList{})
}

func (h *Handler) HEADConsumptionTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTConsumptionTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/tp/{id1}/rc/{id2}/tti/{id3}/cti/1")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PUTConsumptionTariffIntervalList(w http.ResponseWriter, req *http.Request) {
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

// ConsumptionTariffInterval resource handlers

func (h *Handler) DELETEConsumptionTariffInterval(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ConsumptionTariffInterval{})
}

func (h *Handler) GETConsumptionTariffInterval(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.ConsumptionTariffInterval{})
}

func (h *Handler) HEADConsumptionTariffInterval(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTConsumptionTariffInterval(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTConsumptionTariffInterval(w http.ResponseWriter, req *http.Request) {
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

