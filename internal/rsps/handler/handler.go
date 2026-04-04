package handler

import (
"crypto/sha256"
"encoding/xml"
"fmt"
// "log"
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

// Helper to get and validate LFDI from certificate
func (h *Handler) getLFDI(req *http.Request) (string, error) {
cert := req.TLS.PeerCertificates[0]
lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
if _, err := h.repo.GetEntity(lfdi); err != nil {
return "", err
}
return lfdi, nil
}

// ResponseSetList resource handlers

func (h *Handler) GETResponseSetList(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
xml.NewEncoder(w).Encode(&sep.ResponseSetList{})
}

func (h *Handler) HEADResponseSetList(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
}

func (h *Handler) PUTResponseSetList(w http.ResponseWriter, req *http.Request) {
// WADL spec: PUT on ResponseSetList returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTResponseSetList(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
// WADL spec: POST returns 200 or 201 with location header
w.Header().Set("Content-Type", sep.ContentType)
w.Header().Set("location", "/rsps/1")
w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEResponseSetList(w http.ResponseWriter, req *http.Request) {
// WADL spec: DELETE on ResponseSetList returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

// ResponseSet resource handlers

func (h *Handler) GETResponseSet(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.ResponseSet{})
}

func (h *Handler) HEADResponseSet(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTResponseSet(w http.ResponseWriter, req *http.Request) {
// WADL spec: PUT on ResponseSet returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTResponseSet(w http.ResponseWriter, req *http.Request) {
// WADL spec: POST on ResponseSet returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEResponseSet(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
w.WriteHeader(http.StatusBadRequest)
return
}
// WADL spec: DELETE on ResponseSet returns 200 with representation
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
xml.NewEncoder(w).Encode(&sep.ResponseSet{})
}

// ResponseList resource handlers

func (h *Handler) GETResponseList(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.ResponseSetList{})
}

func (h *Handler) HEADResponseList(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTResponseList(w http.ResponseWriter, req *http.Request) {
// WADL spec: PUT on ResponseList returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTResponseList(w http.ResponseWriter, req *http.Request) {
_, err := h.getLFDI(req)
if err != nil {
w.WriteHeader(http.StatusNotFound)
return
}
if _, err := strconv.Atoi(req.PathValue("id1")); err != nil {
w.WriteHeader(http.StatusBadRequest)
return
}
// WADL spec: POST returns 200, 201, or 204
w.Header().Set("Content-Type", sep.ContentType)
w.Header().Set("location", "/rsps/1/rsp/1")
w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DELETEResponseList(w http.ResponseWriter, req *http.Request) {
// WADL spec: DELETE on ResponseList returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

// Response resource handlers (base Response type)

func (h *Handler) GETResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.Response{})
}

func (h *Handler) HEADResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTResponse(w http.ResponseWriter, req *http.Request) {
// WADL spec: PUT on Response returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTResponse(w http.ResponseWriter, req *http.Request) {
// WADL spec: POST on Response returns 400 or 405
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEResponse(w http.ResponseWriter, req *http.Request) {
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
// WADL spec: DELETE returns 200 with representation
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)
xml.NewEncoder(w).Encode(&sep.Response{})
}

// PriceResponse resource handlers

func (h *Handler) GETPriceResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.PriceResponse{})
}

func (h *Handler) HEADPriceResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTPriceResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTPriceResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEPriceResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.PriceResponse{})
}

// TextResponse resource handlers

func (h *Handler) GETTextResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.TextResponse{})
}

func (h *Handler) HEADTextResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTTextResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTTextResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETETextResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.TextResponse{})
}

// DefaultDERControlResponse resource handlers

func (h *Handler) GETDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DefaultDERControlResponse{})
}

func (h *Handler) HEADDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDefaultDERControlResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DefaultDERControlResponse{})
}

// DERControlResponse resource handlers

func (h *Handler) GETDERControlResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DERControlResponse{})
}

func (h *Handler) HEADDERControlResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDERControlResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDERControlResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDERControlResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DERControlResponse{})
}

// FlowReservationResponseResponse resource handlers

func (h *Handler) GETFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.FlowReservationResponseResponse{})
}

func (h *Handler) HEADFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEFlowReservationResponseResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.FlowReservationResponseResponse{})
}

// DrResponse resource handlers

func (h *Handler) GETDrResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DrResponse{})
}

func (h *Handler) HEADDrResponse(w http.ResponseWriter, req *http.Request) {
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

func (h *Handler) PUTDrResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) POSTDrResponse(w http.ResponseWriter, req *http.Request) {
w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) DELETEDrResponse(w http.ResponseWriter, req *http.Request) {
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
xml.NewEncoder(w).Encode(&sep.DrResponse{})
}
