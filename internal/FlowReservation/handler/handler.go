package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/sep"
	"encoding/gob"
)

func init() {
	gob.Register(&sep.FlowReservationRequest{})
}

type Handler struct {
	repo *store.Store
	reg  *registry.Registry
}

func NewHandler(repo *store.Store, reg *registry.Registry) *Handler {
	return &Handler{repo, reg}
}

// getLFDI extracts and validates LFDI from certificate against the registry
func (h *Handler) getLFDI(req *http.Request) (string, error) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		return "", fmt.Errorf("mTLS certificate required")
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	
	if !h.reg.IsAuthorized(lfdi) {
		return "", fmt.Errorf("device %s not authorized", lfdi)
	}
	
	return lfdi, nil
}

// getSFDI extracts SFDI from LFDI
func (h *Handler) getSFDI(lfdi string) (sep.SFDIType, error) {
	return sep.ToSFDI(lfdi)
}

// buildStoreKey builds a store key from SFDI and optional MRID
func (h *Handler) buildStoreKey(sfdi sep.SFDIType, mrid ...string) string {
	key := fmt.Sprintf("%d", sfdi)
	if len(mrid) > 0 && mrid[0] != "" {
		key = fmt.Sprintf("%d:%s", sfdi, mrid[0])
	}
	return key
}

// FlowReservationRequestList resource handlers

func (h *Handler) GETFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	results := h.repo.GetByOwner(fmt.Sprintf("%d", sfdi))
	list := &sep.FlowReservationRequestList{
		List: &sep.List{
			AllAttr:      uint32(len(results)),
			ResultsAttr:  uint32(len(results)),
		},
	}
	
	for _, res := range results {
		if frq, ok := res.(*sep.FlowReservationRequest); ok {
			list.FlowReservationRequest = append(list.FlowReservationRequest, frq)
		}
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(list)
}

func (h *Handler) HEADFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var frq sep.FlowReservationRequest
	if err := xml.NewDecoder(req.Body).Decode(&frq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 5-15 minute EIM Validation
	// DurationRequested is in seconds. 300s = 5 minutes, 900s = 15 minutes.
	if frq.DurationRequested != 300 && frq.DurationRequested != 900 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid DurationRequested: %d. Only 300s (5m) and 900s (15m) are supported for EIM.", frq.DurationRequested)
		return
	}

	// Task 1.2: Validate against DemandResponse (EndDeviceControl) and LoadShedAvailability
	ownerID := fmt.Sprintf("%d", sfdi)
	results := h.repo.GetByOwner(ownerID)
	
	var maxSheddable int16 = 0
	hasLSA := false
	var activeEvents []*sep.EndDeviceControl

	for _, res := range results {
		if lsa, ok := res.(*sep.LoadShedAvailability); ok {
			hasLSA = true
			if lsa.SheddablePower != nil {
				maxSheddable = lsa.SheddablePower.Value
			}
		} else if edc, ok := res.(*sep.EndDeviceControl); ok {
			activeEvents = append(activeEvents, edc)
		}
	}

	// Check PowerRequested against SheddablePower
	if hasLSA && frq.PowerRequested != nil && frq.PowerRequested.Value > maxSheddable {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "PowerRequested (%d) exceeds SheddablePower (%d)", frq.PowerRequested.Value, maxSheddable)
		return
	}

	// Check if IntervalRequested overlaps with an active Demand Response event
	hasOverlap := false
	if frq.IntervalRequested != nil && frq.IntervalRequested.Start != nil {
		reqStart := int64(*frq.IntervalRequested.Start)
		reqEnd := reqStart + int64(frq.IntervalRequested.Duration)

		for _, ev := range activeEvents {
			if ev.RandomizableEvent != nil && ev.RandomizableEvent.Event != nil && ev.RandomizableEvent.Event.Interval != nil && ev.RandomizableEvent.Event.Interval.Start != nil {
				evStart := int64(*ev.RandomizableEvent.Event.Interval.Start)
				evEnd := evStart + int64(ev.RandomizableEvent.Event.Interval.Duration)

				// Overlap condition
				if reqStart < evEnd && reqEnd > evStart {
					hasOverlap = true
					break
				}
			}
		}
	} else {
		// If no interval requested, we can't validate overlap.
		hasOverlap = false
	}

	if len(activeEvents) > 0 && !hasOverlap {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "FlowReservationRequest interval does not overlap with any active Demand Response event")
		return
	} else if len(activeEvents) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Device has no active Demand Response reservations to cover FlowReservationRequest")
		return
	}

	mrid := "frq-" + strconv.FormatInt(int64(sfdi), 10) + "-" + strconv.FormatInt(int64(len(h.repo.GetByOwner(fmt.Sprintf("%d", sfdi)))), 10)
	if frq.MRID != nil && frq.MRID.HexBinary128 != nil {
		mrid = string(*frq.MRID.HexBinary128)
	} else {
		m := mrid
		frq.MRID = &sep.MRIDType{HexBinary128: &m}
	}

	key := h.buildStoreKey(sfdi, mrid)
	if err := h.repo.SetWithOwner(key, &frq, fmt.Sprintf("%d", sfdi)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("Location", "/frq/"+mrid)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEFlowReservationRequestList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// FlowReservationRequest resource handlers

func (h *Handler) GETFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	mrid := req.PathValue("id1")
	key := h.buildStoreKey(sfdi, mrid)
	
	val, ok := h.repo.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	frq, ok := val.(*sep.FlowReservationRequest)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(frq)
}

func (h *Handler) HEADFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	mrid := req.PathValue("id1")
	key := h.buildStoreKey(sfdi, mrid)
	
	_, ok := h.repo.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) POSTFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEFlowReservationRequest(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	mrid := req.PathValue("id1")
	key := h.buildStoreKey(sfdi, mrid)
	h.repo.Delete(key)

	w.WriteHeader(http.StatusNoContent)
}

// FlowReservationResponseList resource handlers

func (h *Handler) GETFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponseList{})
}

func (h *Handler) HEADFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("Location", "/frp/"+fmt.Sprintf("%d", sfdi))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEFlowReservationResponseList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// FlowReservationResponse resource handlers

func (h *Handler) GETFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponse{})
}

func (h *Handler) HEADFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEFlowReservationResponse(w http.ResponseWriter, req *http.Request) {
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)
	if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.FlowReservationResponse{})
}

