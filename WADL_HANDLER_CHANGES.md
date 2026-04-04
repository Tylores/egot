# WADL to Handler Code Update - Comprehensive Analysis

## Executive Summary

**Status:** Analysis Complete - Ready for Implementation

**Scope:** 16 WADL files, 16 handler packages, ~700 API methods across all services

**Key Finding:** All handlers have **5 systematic issues** that need correction:

1. ❌ Encoding `nil` instead of proper SEP structs
2. ❌ Missing HTTP status codes per WADL specification  
3. ❌ Missing required headers (especially `location` for POST)
4. ❌ Incorrect HTTP header ordering (headers after WriteHeader)
5. ❌ HEAD methods incorrectly encoding response bodies

---

## Technical Background

### SEP 2 API Pattern

The SEP 2 API is a RESTful XML-based API with specific response patterns:

- **Content-Type:** `application/sep+xml` for all responses with bodies
- **Key Resource Types:** TariffProfile, Bill, Device, etc.
- **List Resources:** Collections that support GET, POST, HEAD, and reject PUT/DELETE
- **Single Resources:** Individual items that support GET, HEAD, DELETE, and reject POST/PUT

### HTTP Compliance Requirements

Go's `http.ResponseWriter` requires strict header order:

```
1. w.Header().Set(...) - Set all headers FIRST
2. w.WriteHeader(http.StatusXXX) - Call status ONCE
3. w.Body.Write(...) - Write body AFTER
```

**Critical:** Once `WriteHeader()` is called, subsequent `Header().Set()` calls are ignored!

---

## Current Issues (Root Cause Analysis)

### Issue 1: Encoding nil Instead of Structs

**Current (WRONG):**
```go
err = xml.NewEncoder(w).Encode(nil)  // produces empty XML element
```

**Should be:**
```go
err = xml.NewEncoder(w).Encode(&sep.TariffProfileList{})  // proper struct
```

**Impact:** Clients receive `<nil/>` or similar invalid responses, breaking API contracts

---

### Issue 2: Missing/Incorrect Status Codes

**Current (WRONG):**
```go
// ALL methods use default 200, errors use 500/404
w.WriteHeader(http.StatusInternalServerError)  // wrong for all cases
```

**Should be (per WADL):**
- GET: 200 (OK)
- HEAD: 200 (OK)  
- POST to list: 201 (Created)
- PUT/DELETE to list: 405 (MethodNotAllowed)
- DELETE single: 200 (OK)

**Impact:** Clients can't properly interpret responses; breaks HTTP semantics

---

### Issue 3: Missing Required Headers

**Current (WRONG):**
```go
// POST operations never set location header
err = xml.NewEncoder(w).Encode(nil)
```

**Should be:**
```go
w.Header().Set("location", fmt.Sprintf("/tp/%d", resourceID))
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusCreated)
err = xml.NewEncoder(w).Encode(&sep.TariffProfileList{})
```

**Impact:** Clients can't locate newly created resources

**WADL Spec:** POST responses include required `location` header parameter

---

### Issue 4: Header Ordering (HTTP Violation)

**Current (WRONG):**
```go
w.Header().Set("Content-Type", sep.ContentType)  // Set header
err = xml.NewEncoder(w).Encode(nil)              // This internally calls WriteHeader()
```

**Problem:** XML encoder calls `WriteHeader()` implicitly on first write. Any prior `Header().Set()` calls are lost.

**Should be:**
```go
w.Header().Set("Content-Type", sep.ContentType)
w.WriteHeader(http.StatusOK)  // Explicit call
err = xml.NewEncoder(w).Encode(&sep.TariffProfileList{})
```

---

### Issue 5: HEAD Methods with Response Bodies

**Current (WRONG):**
```go
func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", sep.ContentType)
    err = xml.NewEncoder(w).Encode(nil)  // WRONG: HEAD must not have body
}
```

**Should be:**
```go
func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusOK)  // Just return status, no headers or body
}
```

**RFC 7231:** HEAD responses must not include a message body

---

## WADL Specification Patterns

### Pattern 1: List Resources (e.g., `/tp`)

```xml
<resource id="TariffProfileList" path="/tp">
  <!-- GET: retrieve all items -->
  <method id="GETTariffProfileList" name="GET">
    <response>
      <representation mediaType="application/sep+xml" element="sep:TariffProfileList"/>
    </response>
  </method>
  
  <!-- HEAD: check if exists -->
  <method id="HEADTariffProfileList" name="HEAD"/>
  
  <!-- POST: create new item -->
  <method id="POSTTariffProfileList" name="POST">
    <response status="200">
      <param name="location" style="header" required="true"/>
    </response>
    <response status="201">
      <param name="location" style="header" required="true"/>
    </response>
  </method>
  
  <!-- PUT: not allowed -->
  <method id="PUTTariffProfileList" name="PUT">
    <response status="400"/>
    <response status="405"/>
  </method>
  
  <!-- DELETE: not allowed -->
  <method id="DELETETariffProfileList" name="DELETE">
    <response status="400"/>
    <response status="405"/>
  </method>
</resource>
```

**Handler Rules:**
- GET → 200, encode struct
- HEAD → 200, no body
- POST → 201, set location header, encode struct
- PUT → 405, no body
- DELETE → 405, no body

---

### Pattern 2: Single Resources (e.g., `/tp/{id1}`)

```xml
<resource id="TariffProfile" path="/tp/{id1}">
  <!-- GET: retrieve item -->
  <method id="GETTariffProfile" name="GET">
    <response>
      <representation mediaType="application/sep+xml" element="sep:TariffProfile"/>
    </response>
  </method>
  
  <!-- DELETE: remove item -->
  <method id="DELETETariffProfile" name="DELETE">
    <response>
      <representation mediaType="application/xml" element="sep:TariffProfile"/>
    </response>
  </method>
  
  <!-- POST: not allowed -->
  <method id="POSTTariffProfile" name="POST">
    <response status="400"/>
    <response status="405"/>
  </method>
  
  <!-- PUT: not allowed -->
  <method id="PUTTariffProfile" name="PUT">
    <response status="400"/>
    <response status="405"/>
  </method>
</resource>
```

**Handler Rules:**
- GET → 200, encode struct
- HEAD → 200, no body
- DELETE → 200, encode struct
- POST → 405, no body
- PUT → 405, no body

---

## Decision Table: HTTP Status Codes

| Method | Resource | WADL Status | Decision | Handler Code |
|--------|----------|-------------|----------|--------------|
| GET | List | 200 | Always 200 | `http.StatusOK` |
| GET | Single | 200 | Always 200 | `http.StatusOK` |
| HEAD | Any | 200 | Always 200 | `http.StatusOK` |
| POST | List | 200/201 | Use 201 (higher) | `http.StatusCreated` |
| POST | Single | 400/405 | Use 405 | `http.StatusMethodNotAllowed` |
| PUT | List | 400/405 | Use 405 | `http.StatusMethodNotAllowed` |
| PUT | Single | 400/405 | Use 405 | `http.StatusMethodNotAllowed` |
| DELETE | List | 400/405 | Use 405 | `http.StatusMethodNotAllowed` |
| DELETE | Single | 200 | Always 200 | `http.StatusOK` |

---

## Decision Table: Response Bodies

| Method | Resource | Response Body | Content-Type | Example |
|--------|----------|---------------|--------------|---------|
| GET | List | Yes (struct) | `application/sep+xml` | `&sep.TariffProfileList{}` |
| GET | Single | Yes (struct) | `application/sep+xml` | `&sep.TariffProfile{}` |
| HEAD | Any | No | (none) | (no body) |
| POST | List | Yes (struct) | `application/sep+xml` | `&sep.TariffProfileList{}` |
| POST | Single | No | (none) | (no body) |
| PUT | List | No | (none) | (no body) |
| PUT | Single | No | (none) | (no body) |
| DELETE | List | No | (none) | (no body) |
| DELETE | Single | Yes (struct) | `application/sep+xml` | `&sep.TariffProfile{}` |

---

## Code Transformation Templates

### Template 1: GET List Resource

```go
// BEFORE (WRONG)
func (h *Handler) GETTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)
    err = xml.NewEncoder(w).Encode(nil)  // ← WRONG
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)  // ← TOO LATE
        return
    }
}

// AFTER (CORRECT)
func (h *Handler) GETTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)  // ← ADD THIS
    err = xml.NewEncoder(w).Encode(&sep.TariffProfileList{})  // ← CHANGE nil TO STRUCT
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        return  // ← REMOVE WriteHeader call
    }
}
```

---

### Template 2: POST to List Resource

```go
// BEFORE (WRONG)
func (h *Handler) POSTTariffProfileList(w http.ResponseWriter, req *http.Request) {
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

// AFTER (CORRECT)
func (h *Handler) POSTTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    // ADD LOCATION HEADER - derive from request/ID
    w.Header().Set("location", fmt.Sprintf("/tp/%d", rand.Intn(1000)))  // ← ADD THIS
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusCreated)  // ← CHANGE TO 201
    err = xml.NewEncoder(w).Encode(&sep.TariffProfileList{})  // ← CHANGE nil TO STRUCT
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        return
    }
}
```

---

### Template 3: HEAD (Any Resource)

```go
// BEFORE (WRONG)
func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)  // ← REMOVE THIS
    err = xml.NewEncoder(w).Encode(nil)  // ← REMOVE THIS
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)  // ← REMOVE THIS
        return  // ← REMOVE THIS
    }
}

// AFTER (CORRECT)
func (h *Handler) HEADTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusOK)  // ← ONLY THIS
}
```

---

### Template 4: DELETE to Single Resource

```go
// BEFORE (WRONG)
func (h *Handler) DELETETariffProfile(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)
    err = xml.NewEncoder(w).Encode(nil)  // ← WRONG
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}

// AFTER (CORRECT)
func (h *Handler) DELETETariffProfile(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)  // ← ADD THIS
    err = xml.NewEncoder(w).Encode(&sep.TariffProfile{})  // ← CHANGE nil TO STRUCT
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        return
    }
}
```

---

### Template 5: PUT/DELETE to List (Read-Only)

```go
// BEFORE (WRONG)
func (h *Handler) PUTTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", sep.ContentType)  // ← REMOVE
    err = xml.NewEncoder(w).Encode(nil)  // ← REMOVE
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)  // ← REMOVE
        return  // ← REMOVE
    }
}

// AFTER (CORRECT)
func (h *Handler) PUTTariffProfileList(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusMethodNotAllowed)  // ← ONLY THIS - 405
}
```

---

## Mapping of SEP Element Types

The following structs are defined in `github.com/Tylores/egot/sep`:

| Element Type | Go Type |
|--------------|---------|
| TariffProfileList | `sep.TariffProfileList` |
| TariffProfile | `sep.TariffProfile` |
| RateComponentList | `sep.RateComponentList` |
| RateComponent | `sep.RateComponent` |
| CustomerAccountList | `sep.CustomerAccountList` |
| CustomerAccount | `sep.CustomerAccount` |
| BillingPeriodList | `sep.BillingPeriodList` |
| BillingPeriod | `sep.BillingPeriod` |
| TimeIntervalList | `sep.TimeIntervalList` |
| TimeInterval | `sep.TimeInterval` |
| EndDeviceList | `sep.EndDeviceList` |
| EndDevice | `sep.EndDevice` |
| And ~50+ more...| |

---

## Implementation Checklist

For each handler file in `internal/*/handler/handler.go`:

- [ ] Identify all methods and their WADL patterns
- [ ] For each GET/DELETE method: add proper status code before encoding
- [ ] For each POST to list: add location header + 201 status
- [ ] For each POST/PUT to single: change to 405 with no body
- [ ] For each HEAD method: remove body encoding
- [ ] For each PUT/DELETE to list: change to 405 with no body
- [ ] Verify header order: Headers → WriteHeader → Body
- [ ] Test compilation: `go build ./...`
- [ ] Verify no duplicate WriteHeader calls
- [ ] Verify all struct types exist in `sep` package

---

## Validation Steps

1. **Compile Check:**
   ```bash
   cd /home/tylor/dev/egot
   go build ./...
   ```

2. **HTTP Header Order Verification:**
   - grep all handler files for patterns of headers after WriteHeader
   - Ensure all Header().Set() calls come before WriteHeader()

3. **Status Code Verification:**
   - Verify all GET methods return 200
   - Verify all POST to list methods return 201
   - Verify all HEAD methods return 200
   - Verify all PUT/DELETE to list return 405

4. **Struct Encoding Verification:**
   - No `Encode(nil)` calls remain
   - All body-returning methods encode proper types
   - HEAD methods encode nothing

5. **Header Presence Verification:**
   - All POST to list methods set location header
   - All responses with bodies set Content-Type header

---

## Files Affected

### Handler Packages (16 total):
- `internal/Bill/handler/handler.go`
- `internal/BRS/handler/handler.go`
- `internal/DCAP/handler/handler.go`
- `internal/DERP/handler/handler.go`
- `internal/DR/handler/handler.go`
- `internal/EDevice/handler/handler.go`
- `internal/File/handler/handler.go`
- `internal/Messaging/handler/handler.go`
- `internal/MUP/handler/handler.go`
- `internal/Notify/handler/handler.go`
- `internal/Pricing/handler/handler.go`
- `internal/SDevice/handler/handler.go`
- `internal/TariffProfile/handler/handler.go`
- `internal/TimeOfUse/handler/handler.go`
- `internal/UPT/handler/handler.go`
- `internal/rsps/handler/handler.go`

### WADL Files (16 total):
- `wadl/tp.wadl`
- `wadl/brs.wadl`
- `wadl/bill.wadl`
- `wadl/dcap.wadl`
- `wadl/derp.wadl`
- `wadl/dr.wadl`
- `wadl/edev.wadl`
- `wadl/file.wadl`
- `wadl/msg.wadl`
- `wadl/mup.wadl`
- `wadl/ntfy.wadl`
- `wadl/ppy.wadl`
- `wadl/rsps.wadl`
- `wadl/sdev.wadl`
- `wadl/tm.wadl`
- `wadl/upt.wadl`

---

## Next Steps

1. ✅ **Analysis Complete** (this document)
2. ⏳ Generate automated code fixes for each handler
3. ⏳ Validate fixes compile and match WADL specs
4. ⏳ Update git with migration commit
5. ⏳ Test API responses match WADL contracts

---

*Analysis generated by WADL Parser v1.0*
*Document: WADL_HANDLER_CHANGES.md*
