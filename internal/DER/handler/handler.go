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
	"time"
)

func init() {
	gob.Register(&sep.DERCurve{})
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

// DERProgramList resource handlers

func (h *Handler) GETDERProgramList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERProgramList{})
}

func (h *Handler) HEADDERProgramList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERProgramList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/derp/"+fmt.Sprintf("%d", sfdi))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEDERProgramList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERProgram resource handlers

func (h *Handler) GETDERProgram(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERProgram{})
}

func (h *Handler) HEADDERProgram(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERProgram(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDERProgram(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERProgram{})
}

// ActiveDERControlList resource handlers

func (h *Handler) GETActiveDERControlList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERControlList{})
}

func (h *Handler) HEADActiveDERControlList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEActiveDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERControlList resource handlers

func (h *Handler) GETDERControlList(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERControlList{})
}

func (h *Handler) HEADDERControlList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERControlList(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("location", "/derp/{id1}/derc/"+fmt.Sprintf("%d", sfdi))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEDERControlList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERControl resource handlers

func (h *Handler) GETDERControl(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERControl{})
}

func (h *Handler) HEADDERControl(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDERControl(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERControl{})
}

// DefaultDERControl resource handlers

func (h *Handler) GETDefaultDERControl(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DefaultDERControl{})
}

func (h *Handler) HEADDefaultDERControl(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDefaultDERControl(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) POSTDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDefaultDERControl(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERCurveList resource handlers

func (h *Handler) GETDERCurveList(w http.ResponseWriter, req *http.Request) {
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
	list := &sep.DERCurveList{
		List: &sep.List{
			AllAttr:     uint32(len(results)),
			ResultsAttr: uint32(len(results)),
		},
	}

	for _, res := range results {
		if curve, ok := res.(*sep.DERCurve); ok {
			list.DERCurve = append(list.DERCurve, curve)
		}
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(list)
}

func (h *Handler) HEADDERCurveList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERCurveList(w http.ResponseWriter, req *http.Request) {
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

	var curve sep.DERCurve
	if err := xml.NewDecoder(req.Body).Decode(&curve); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if curve.CurveType == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "CurveType is required")
		return
	}

	cType := *curve.CurveType
	// 0: Frequency-Watt, 11: Volt-Var
	if cType == 0 || cType == 11 {
		if len(curve.CurveData) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "CurveData must contain at least 2 points for interpolation")
			return
		}
	}

	if curve.IdentifiedObject == nil {
		curve.IdentifiedObject = &sep.IdentifiedObject{}
	}

	mrid := "curve-" + strconv.FormatInt(int64(sfdi), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	if curve.MRID != nil && curve.MRID.HexBinary128 != nil {
		mrid = string(*curve.MRID.HexBinary128)
	} else {
		m := mrid
		curve.MRID = &sep.MRIDType{HexBinary128: &m}
	}

	key := h.buildStoreKey(sfdi, mrid)
	h.repo.SetWithOwner(key, &curve, fmt.Sprintf("%d", sfdi))

	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("Location", "/derp/"+req.PathValue("id1")+"/dc/"+mrid)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEDERCurveList(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// DERCurve resource handlers

func (h *Handler) GETDERCurve(w http.ResponseWriter, req *http.Request) {
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

	mrid := req.PathValue("id2")
	key := h.buildStoreKey(sfdi, mrid)

	val, ok := h.repo.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	curve, ok := val.(*sep.DERCurve)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(curve)
}

func (h *Handler) HEADDERCurve(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERCurve(w http.ResponseWriter, req *http.Request) {
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDERCurve(w http.ResponseWriter, req *http.Request) {
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
	xml.NewEncoder(w).Encode(&sep.DERCurve{})
}
