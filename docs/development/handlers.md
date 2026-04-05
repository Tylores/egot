# HTTP Handlers

## Naming Convention

Handler methods follow a deterministic naming pattern:

```
<HTTP_METHOD><ResourceName>
```

Examples:
- `GETCustomerAccountList` — GET /bill
- `POSTCustomerAccountList` — POST /bill
- `GETCustomerAccount` — GET /bill/{id1}
- `DELETECustomerAccount` — DELETE /bill/{id1}
- `HEADCustomerAccountList` — HEAD /bill

All public handler methods follow this pattern; private helpers use lowercase.

## Registering Handlers

### Option 1: scaffold-gen (recommended for new services)

`scaffold-gen` reads a WADL specification and generates a fully wired `cmd/<Service>/main.go` with correct route registrations. See [docs/tools/scaffold-gen.md](../tools/scaffold-gen.md).

### Option 2: Regenerate from existing handlers

For services that already have handler files:

```bash
./scripts/regenerate_cmd_handlers.sh              # all services
./scripts/regenerate_cmd_handlers.sh Bill         # specific service
```

The script inspects `internal/<Service>/handler/handler.go`, extracts public handler methods, and rewrites `cmd/<Service>/main.go`.

### Manual registration

Register handlers with HTTP method and path pattern:

```go
h := handler.NewHandler(repo)

// List resources
http.HandleFunc("GET /tp", h.GETTariffProfileList)
http.HandleFunc("HEAD /tp", h.HEADTariffProfileList)
http.HandleFunc("POST /tp", h.POSTTariffProfileList)
http.HandleFunc("PUT /tp", h.PUTTariffProfileList)
http.HandleFunc("DELETE /tp", h.DELETETariffProfileList)

// Single resources
http.HandleFunc("GET /tp/{id1}", h.GETTariffProfile)
http.HandleFunc("HEAD /tp/{id1}", h.HEADTariffProfile)
http.HandleFunc("DELETE /tp/{id1}", h.DELETETariffProfile)
http.HandleFunc("PUT /tp/{id1}", h.PUTTariffProfile)
http.HandleFunc("POST /tp/{id1}", h.POSTTariffProfile)
```

## WADL-Compliant Response Patterns

The SEP 2 API defines specific response patterns per HTTP method and resource type. All handlers must follow these rules.

### Status Code Decision Table

| Method | Resource Type | Status | Body | Location Header |
|--------|--------------|--------|------|-----------------|
| GET | List | 200 | Yes (struct) | No |
| GET | Single | 200 | Yes (struct) | No |
| HEAD | Any | 200 | No | No |
| POST | List | 201 | No | **Required** |
| POST | Single | 405 | No | No |
| PUT | List | 405 | No | No |
| PUT | Single | 405 | No | No |
| DELETE | List | 405 | No | No |
| DELETE | Single | 200 | Yes (struct) | No |

### HTTP Header Ordering (Critical)

Go's `http.ResponseWriter` requires strict ordering:

```go
// 1. Set all headers first
w.Header().Set("location", ...)      // if needed
w.Header().Set("Content-Type", ...)  // if body present

// 2. Write status code once
w.WriteHeader(http.StatusXXX)

// 3. Write body after
xml.NewEncoder(w).Encode(&sep.MyStruct{})
```

**Important:** `WriteHeader` is called implicitly on the first body write. Any `Header().Set()` call after that is silently ignored.

### Template 1: GET List Resource

```go
func (h *Handler) GETTariffProfileList(w http.ResponseWriter, req *http.Request) {
    _, err := h.getLFDI(req)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)
    if err := xml.NewEncoder(w).Encode(&sep.TariffProfileList{}); err != nil {
        log.Printf("Response encode error: %v\n", err)
    }
}
```

### Template 2: POST to List Resource

```go
func (h *Handler) POSTTariffProfileList(w http.ResponseWriter, req *http.Request) {
    _, err := h.getLFDI(req)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("location", "/tp/1") // TODO: generate actual resource ID
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusCreated)
}
```

### Template 3: HEAD (Any Resource)

```go
func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
    _, err := h.getLFDI(req)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

### Template 4: DELETE Single Resource

```go
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
    if err := xml.NewEncoder(w).Encode(&sep.TariffProfile{}); err != nil {
        log.Printf("Response encode error: %v\n", err)
    }
}
```

### Template 5: PUT/DELETE to List or POST/PUT to Single (Not Allowed)

```go
func (h *Handler) PUTTariffProfileList(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
}
```

## Reference Implementation

`internal/rsps/handler/handler.go` is the reference implementation. It demonstrates the correct helper pattern:

```go
func (h *Handler) getLFDI(req *http.Request) (string, error) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    if _, err := h.repo.GetEntity(lfdi); err != nil {
        return "", err
    }
    return lfdi, nil
}
```

Always use `getLFDI` as the first step in any handler that accesses repository data.

## Handler File Structure

```go
package handler

import (
    "crypto/sha256"
    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/Tylores/egot/internal/<ServiceName>/repository/memory"
    "github.com/Tylores/egot/sep"
)

type Handler struct {
    repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
    return &Handler{repo}
}

func (h *Handler) getLFDI(req *http.Request) (string, error) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    if _, err := h.repo.GetEntity(lfdi); err != nil {
        return "", err
    }
    return lfdi, nil
}

// ... handler methods follow
```

## SEP Element Types

Common types from the `sep` package used in response encoding:

| Resource | List Type | Single Type |
|----------|-----------|-------------|
| TariffProfile | `sep.TariffProfileList` | `sep.TariffProfile` |
| Bill | `sep.CustomerAccountList` | `sep.CustomerAccount` |
| EndDevice | `sep.EndDeviceList` | `sep.EndDevice` |
| ResponseSet | `sep.ResponseSetList` | `sep.ResponseSet` |
| BillingPeriod | `sep.BillingPeriodList` | `sep.BillingPeriod` |

Use `grep -r "type " sep/` to find all available types.

## Code Quality Checklist

- [ ] Handler methods follow `<HTTP_METHOD><ResourceName>` naming
- [ ] LFDI validation happens before any resource access
- [ ] Headers set before `WriteHeader`
- [ ] `WriteHeader` called once and before body encoding
- [ ] POST to list sets `location` header
- [ ] HEAD methods return 200 with no body
- [ ] PUT/DELETE to list and POST/PUT to single return 405 with no body
- [ ] `Encode(nil)` replaced with `Encode(&sep.ActualType{})`
- [ ] No duplicate `WriteHeader` calls
- [ ] `go build ./...` succeeds

## WADL Files Reference

| Service | WADL File | Resources |
|---------|-----------|-----------|
| TariffProfile | wadl/tp.wadl | 9 |
| Bill | wadl/bill.wadl | 14 |
| BRS | wadl/brs.wadl | 4 |
| DCAP | wadl/dcap.wadl | 1 |
| DERP | wadl/derp.wadl | 8 |
| DR | wadl/dr.wadl | 5 |
| EDevice | wadl/edev.wadl | 53 |
| File | wadl/file.wadl | 2 |
| Messaging | wadl/msg.wadl | 5 |
| MUP | wadl/mup.wadl | 2 |
| Notify | wadl/ntfy.wadl | 2 |
| PPY | wadl/ppy.wadl | 9 |
| RSPS | wadl/rsps.wadl | 10 |
| SDevice | wadl/sdev.wadl | 1 |
| TimeOfUse | wadl/tm.wadl | 1 |
| UPT | wadl/upt.wadl | 9 |
