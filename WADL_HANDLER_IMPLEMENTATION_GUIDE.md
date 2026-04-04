# WADL Handler Implementation Guide

## Overview

This guide explains how to implement HTTP handlers that follow the response specifications defined in WADL (Web Application Description Language) files for each microservice.

## Pattern Reference: RSPS Microservice

The RSPS (Response Sets) handler at `internal/rsps/handler/handler.go` serves as the reference implementation. It demonstrates the correct pattern for all microservices.

### Key Features of Reference Implementation

1. **Helper Method for Common Logic**
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

2. **HTTP Method Patterns**

   **GET (List Resources)**
   ```go
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
   ```

   **POST on List (Create Resource)**
   ```go
   func (h *Handler) POSTResponseSetList(w http.ResponseWriter, req *http.Request) {
       _, err := h.getLFDI(req)
       if err != nil {
           w.WriteHeader(http.StatusNotFound)
           return
       }
       w.Header().Set("Content-Type", sep.ContentType)
       w.Header().Set("location", "/rsps/1")  // TODO: generate actual ID
       w.WriteHeader(http.StatusCreated)
   }
   ```

   **PUT on List (Not Allowed)**
   ```go
   func (h *Handler) PUTResponseSetList(w http.ResponseWriter, req *http.Request) {
       w.WriteHeader(http.StatusMethodNotAllowed)
   }
   ```

   **DELETE on Individual Resource**
   ```go
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
       w.Header().Set("Content-Type", sep.ContentType)
       w.WriteHeader(http.StatusOK)
       xml.NewEncoder(w).Encode(&sep.ResponseSet{})
   }
   ```

## Response Status Code Mapping

### List Resources (e.g., `/rsps`, `/rsps/{id}/items`)
| Method | Status | Headers | Body |
|--------|--------|---------|------|
| GET | 200 OK | Content-Type | Resource list |
| HEAD | 200 OK | Content-Type | (empty) |
| POST | 201 Created | location | (empty) |
| PUT | 405 Method Not Allowed | — | — |
| DELETE | 405 Method Not Allowed | — | — |

### Individual Resources (e.g., `/rsps/{id}`, `/rsps/{id}/items/{id2}`)
| Method | Status | Headers | Body |
|--------|--------|---------|------|
| GET | 200 OK | Content-Type | Resource |
| HEAD | 200 OK | Content-Type | (empty) |
| POST | 405 Method Not Allowed | — | — |
| PUT | 405 Method Not Allowed | — | — |
| DELETE | 200 OK | Content-Type | Resource |

## Implementing a New Handler

### Step 1: Identify Resource Structure from WADL

For example, for `tm.wadl` (TimeOfUse):
```xml
<resource id="Time" wx:samplePath="/tm">
  <method id="GETTime" name="GET">
    <response>
      <representation mediaType="application/sep+xml" element="sep:Time"/>
    </response>
  </method>
  <!-- ... other methods ... -->
</resource>
```

### Step 2: Create Handler File

Path: `internal/{ServiceName}/handler/handler.go`

Start with the basic structure:
```go
package handler

import (
    "crypto/sha256"
    "encoding/xml"
    "fmt"
    "net/http"
    "strconv"

    "github.com/Tylores/egot/internal/{ServiceName}/repository/memory"
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
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    if _, err := h.repo.GetEntity(lfdi); err != nil {
        return "", err
    }
    return lfdi, nil
}
```

### Step 3: Implement Handler Methods

For each resource in the WADL, implement 5 methods (GET, HEAD, POST, PUT, DELETE).

Follow these rules:
- Always validate LFDI first
- Validate path parameters (id1, id2, etc.)
- Set Content-Type header
- Return appropriate status code
- Encode empty resource for GET/DELETE on resources

### Step 4: Test Compilation

```bash
cd /home/tylor/dev/egot
go build ./internal/{ServiceName}/handler/...
```

### Step 5: Commit

```bash
git add internal/{ServiceName}/handler/handler.go
git commit -m "feat: implement WADL-compliant handlers for {ServiceName} microservice

- Implement all {N} HTTP handlers following WADL specification
- Return proper HTTP status codes (200, 201, 405, etc.)
- Add location header for POST responses (201 Created)
- Validate LFDI from certificate
- All handlers compile successfully

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

## WADL Files by Service

| Service | File | Resources | Methods | Complexity |
|---------|------|-----------|---------|------------|
| RSPS | rsps.wadl | 10 | 50 | ✅ Implemented |
| TimeOfUse | tm.wadl | 1 | 5 | Easy |
| DCAP | dcap.wadl | 1 | 5 | Easy |
| SDevice | sdev.wadl | 1 | 5 | Easy |
| File | file.wadl | 2 | 10 | Easy |
| Notify | ntfy.wadl | 2 | 10 | Easy |
| MUP | mup.wadl | 2 | 10 | Easy |
| BRS | brs.wadl | 4 | 20 | Medium |
| Messaging | msg.wadl | 5 | 25 | Medium |
| DR | dr.wadl | 5 | 25 | Medium |
| TariffProfile | tp.wadl | 9 | 45 | Medium |
| DERP | derp.wadl | 8 | 40 | Medium |
| Bill | bill.wadl | 14 | 70 | Hard |
| PPY | ppy.wadl | 9 | 45 | Hard |
| UPT | upt.wadl | 9 | 45 | Hard |
| EDevice | edev.wadl | 53 | 265 | Very Hard |

## Common Issues and Solutions

### Issue: Import not found
**Solution**: Ensure the service repository type exists at `internal/{ServiceName}/repository/memory/repository.go`

### Issue: sep.{ResourceType} not defined  
**Solution**: Check that the type exists in `sep/sep.go`. Use `grep -r "type {ResourceType}" sep/`

### Issue: Handlers compile but handler methods not found in main
**Solution**: Ensure handler method names match what main.go expects. Run `go build ./cmd/{Service}` to identify missing methods.

### Issue: Path values always nil
**Solution**: Ensure route handlers use named groups like `{id1}`, `{id2}`. Check routing setup in `cmd/{Service}/main.go`

## Code Quality Checklist

- [ ] All methods follow WADL specification status codes
- [ ] LFDI validation happens before any resource access
- [ ] Path parameters validated with strconv.Atoi()
- [ ] Content-Type header set for all responses with bodies
- [ ] location header set for POST responses (201 Created)
- [ ] Empty response structures returned (no mock data)
- [ ] All methods compile without errors
- [ ] Build succeeds: `go build ./...`
- [ ] Commit message follows pattern
- [ ] Git history is clean

## References

- WADL Specification: `wadl/wadl.xsd`
- SEP 2.0 Master WADL: `wadl/sep_wadl.xml`
- Reference Implementation: `internal/rsps/handler/handler.go`
- Type Definitions: `sep/sep.go`

