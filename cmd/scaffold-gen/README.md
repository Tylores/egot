# Microservice Scaffold Generator

The scaffold generator (`cmd/scaffold-gen`) creates complete microservices from SEP 2 WADL specifications extracted by `wadl-extract`. Each generated service follows the established pattern from the existing services.

## Quick Start

### 1. Extract a service WADL from the SEP 2 spec

Use `wadl-extract` to pull a service's resources out of `wadl/sep_wadl.xml` and write a standalone WADL file:

```bash
go run ./cmd/wadl-extract -wadl wadl/sep_wadl.xml -path /dcap -output wadl/dcap.wadl
```

The `-path` flag is the resource path prefix for the service (e.g., `/dcap`, `/brs`, `/edev`).

### 2. Generate the scaffold

Run from the project root — the filename becomes the service name:

```bash
go run ./cmd/scaffold-gen -wadl wadl/dcap.wadl
```

This creates:
- `cmd/dcap/main.go` — service entry point with TLS and route registration
- `internal/dcap/handler/handler.go` — HTTP handler stubs matching the WADL methods
- `internal/dcap/repository/memory/memory.go` — in-memory storage
- `internal/dcap/repository/error.go` — error definitions
- `internal/dcap/server/server.go` — server utilities
- Updates `internal/routes/routes.go` with the new service constant and next available `egot.internal.com` port

### 3. Implement business logic

The generated handler stubs follow the same pattern as all other services:

```go
// getLFDI helper is already generated
func (h *Handler) GETDeviceCapability(w http.ResponseWriter, req *http.Request) {
    _, err := h.getLFDI(req)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)
    xml.NewEncoder(w).Encode(&sep.DeviceCapability{}) // replace with repo call
}
```

Add repository methods in `repository/memory/memory.go` and wire them in handlers.

### 4. Build and run

```bash
go build ./cmd/dcap
./dcap
```

## Generated handler behaviour

The generator maps each WADL `wx:mode` to a default HTTP response:

| Mode | Meaning | Generated behaviour |
|------|---------|---------------------|
| `M`  | Mandatory | implement: GET→200+encode, HEAD→200, POST→201+location, DELETE→200+encode, PUT→200 |
| `D`  | Discoverable | same as M |
| `E`  | Error / Not Allowed | `StatusMethodNotAllowed` |

## Architecture

Generated services follow the same pattern as all existing services:

1. **Handler** — `getLFDI()` helper validates client cert; each method has the correct stub body
2. **Repository** — in-memory storage with `sync.RWMutex`
3. **Main** — mutual TLS via `tlsutil.NewServerConfig`, repository init, route registration

## Tips

- Service name is the WADL filename without extension (`dcap.wadl` → `dcap`)
- `scaffold-gen` assigns the next available `egot.internal.com` port when it updates `internal/routes/routes.go`
- `max_entities` defaults to 100 if not set in the WADL
