package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/sep"
	"encoding/gob"
)

func init() {
	gob.Register(&sep.EndDeviceControl{})
	gob.Register(&sep.DemandResponseProgram{})
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

// DemandResponseProgramList resource handlers

func (h *Handler) GETDemandResponseProgramList(w http.ResponseWriter, req *http.Request) {
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
	list := &sep.DemandResponseProgramList{
		SubscribableList: &sep.SubscribableList{},
	}
	for _, res := range results {
		if drp, ok := res.(*sep.DemandResponseProgram); ok {
			list.DemandResponseProgram = append(list.DemandResponseProgram, drp)
		}
	}
	list.SubscribableList.AllAttr = uint32(len(list.DemandResponseProgram))
	list.SubscribableList.ResultsAttr = uint32(len(list.DemandResponseProgram))

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(list)
}

func (h *Handler) HEADDemandResponseProgramList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDemandResponseProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDemandResponseProgramList(w http.ResponseWriter, req *http.Request) {
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

	var drp sep.DemandResponseProgram
	if err := xml.NewDecoder(req.Body).Decode(&drp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mrid := "drp-" + strconv.FormatInt(int64(sfdi), 10)
	if drp.IdentifiedObject != nil && drp.MRID != nil && drp.MRID.HexBinary128 != nil {
		mrid = string(*drp.MRID.HexBinary128)
	} else {
		if drp.IdentifiedObject == nil {
			drp.IdentifiedObject = &sep.IdentifiedObject{}
		}
		m := mrid
		drp.MRID = &sep.MRIDType{HexBinary128: &m}
	}

	key := h.buildStoreKey(sfdi, mrid)
	if err := h.repo.SetWithOwner(key, &drp, fmt.Sprintf("%d", sfdi)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("Location", "/dr/"+mrid)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEDemandResponseProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DemandResponseProgram resource handlers

func (h *Handler) GETDemandResponseProgram(w http.ResponseWriter, req *http.Request) {
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
	drp, ok := val.(*sep.DemandResponseProgram)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(drp)
}

func (h *Handler) HEADDemandResponseProgram(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDemandResponseProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDemandResponseProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDemandResponseProgram(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DemandResponseProgram{})
}

// ActiveEndDeviceControlList resource handlers

func (h *Handler) GETActiveEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.EndDeviceControlList{})
}

func (h *Handler) HEADActiveEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTActiveEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTActiveEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEActiveEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// EndDeviceControlList resource handlers

func (h *Handler) GETEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
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
	list := &sep.EndDeviceControlList{
		SubscribableList: &sep.SubscribableList{},
	}
	for _, res := range results {
		if edc, ok := res.(*sep.EndDeviceControl); ok {
			list.EndDeviceControl = append(list.EndDeviceControl, edc)
		}
	}
	list.SubscribableList.AllAttr = uint32(len(list.EndDeviceControl))
	list.SubscribableList.ResultsAttr = uint32(len(list.EndDeviceControl))

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(list)
}

func (h *Handler) HEADEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
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

	var edc sep.EndDeviceControl
	if err := xml.NewDecoder(req.Body).Decode(&edc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mrid := "edc-" + strconv.FormatInt(int64(sfdi), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	if edc.RespondableSubscribableIdentifiedObject == nil {
		edc.RespondableSubscribableIdentifiedObject = &sep.RespondableSubscribableIdentifiedObject{}
	}
	if edc.MRID != nil && edc.MRID.HexBinary128 != nil {
		mrid = string(*edc.MRID.HexBinary128)
	} else {
		m := mrid
		edc.MRID = &sep.MRIDType{HexBinary128: &m}
	}

	key := h.buildStoreKey(sfdi, mrid)
	if err := h.repo.SetWithOwner(key, &edc, fmt.Sprintf("%d", sfdi)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("Location", "/dr/"+req.PathValue("id1")+"/edc/"+mrid)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEEndDeviceControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// EndDeviceControl resource handlers

func (h *Handler) GETEndDeviceControl(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.EndDeviceControl{})
}

func (h *Handler) HEADEndDeviceControl(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTEndDeviceControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTEndDeviceControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEEndDeviceControl(w http.ResponseWriter, req *http.Request) {
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
	if _, err := strconv.Atoi(req.PathValue("id2")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(&sep.EndDeviceControl{})
}
