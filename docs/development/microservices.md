# Microservices

## Service Directory

All 18 microservices live under `cmd/` and `internal/`:

| Service | Port | CMD | Internal |
|---------|------|-----|----------|
| BRS | 8010 | cmd/BRS | internal/BRS |
| Bill | 8011 | cmd/Bill | internal/Bill |
| DCAP | 8012 | cmd/DCAP | internal/DCAP |
| DERP | 8013 | cmd/DERP | internal/DERP |
| DR | 8014 | cmd/DR | internal/DR |
| EDevice | 8015 | cmd/EDevice | internal/EDevice |
| File | 8016 | cmd/File | internal/File |
| MUP | 8017 | cmd/MUP | internal/MUP |
| Messaging | 8018 | cmd/Messaging | internal/Messaging |
| Notify | 8019 | cmd/Notify | internal/Notify |
| PPY | 8020 | cmd/PPY | internal/PPY |
| SDevice | 8021 | cmd/SDevice | internal/SDevice |
| TariffProfile | 8022 | cmd/TariffProfile | internal/TariffProfile |
| TimeOfUse | 8023 | cmd/TimeOfUse | internal/TimeOfUse |
| UPT | 8024 | cmd/UPT | internal/UPT |
| DER | 8026 | cmd/DER | internal/DER |
| FlowReservation | 8027 | cmd/FlowReservation | internal/FlowReservation |
| rsps | 8041 | cmd/rsps | internal/rsps |

Plus 5 original tools/services: `crawler`, `operator`, `client`, `scaffold-gen`, `wadl-extract`.

Port constants are defined in `internal/routes/routes.go`.

## Service Structure

Each service follows this layout:

```
cmd/<Service>/
  main.go                          # Entry point
  <service>_route_registration_test.go # Route registration tests

internal/<Service>/
  handler/handler.go               # HTTP handlers
  repository/error.go              # Error types
  server/server.go                 # Server utilities
```

### Example main.go

```go
package main

import (
    "log"
    "net/http"

    "github.com/Tylores/egot/internal/Bill/handler"
    "github.com/Tylores/egot/internal/Bill/server"
    "github.com/Tylores/egot/internal/routes"
    "github.com/Tylores/egot/internal/store"
)

func main() {
    // Initialize store
    s := store.NewStore("data/Bill.db")
    
    // Create handler
    h := handler.NewHandler(s)
    
    // Register routes
    registerRoutes(h)
    
    // Start server
    srv := server.NewServer(routes.Bill, h)
    log.Printf("Starting Bill service on %s", routes.Bill)
    if err := srv.ListenAndServeTLS("ssl/server.crt", "ssl/server.key"); err != nil {
        log.Fatal(err)
    }
}
```

## Building Services

```bash
# Build a single service
go build -o ./bin/Bill ./cmd/Bill

# Build all services
go build ./cmd/...

# Or loop over generated services
for service in BRS Bill DCAP DERP DR EDevice File MUP Messaging Notify PPY SDevice TariffProfile TimeOfUse UPT; do
    go build -o ./bin/$service ./cmd/$service
done
```

## Generating a New Service

Use `scaffold-gen` to create a complete service from a WADL specification:

```bash
go run ./cmd/scaffold-gen -wadl wadl/my-service.wadl -output .
```

See [docs/tools/scaffold-gen.md](../tools/scaffold-gen.md) for the full WADL specification format.

## Customising Handler Routes

After generating a service, register specific HTTP routes in `cmd/<Service>/main.go`:

```go
h := handler.NewHandler(repo)

http.HandleFunc("GET /bills", h.GETCustomerAccountList)
http.HandleFunc("GET /bills/{id1}", h.GETCustomerAccount)
http.HandleFunc("POST /bills", h.POSTCustomerAccountList)
http.HandleFunc("DELETE /bills/{id1}", h.DELETECustomerAccount)
```

To regenerate handler registrations from existing handler files:

```bash
./scripts/regenerate_cmd_handlers.sh              # all services
./scripts/regenerate_cmd_handlers.sh Bill         # specific service
./scripts/regenerate_cmd_handlers.sh Bill Messaging  # multiple services
```

## Implementing Handlers

Add methods to `internal/<Service>/handler/handler.go`:

```go
func (h *Handler) GETCustomerAccountList(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

See [docs/development/handlers.md](handlers.md) for WADL-compliant response patterns and best practices.

## Adding Repository Methods

Extend `internal/<Service>/repository/memory/repository.go`:

```go
func (r *Repository) GetByID(id string) (interface{}, error) {
    // implementation
}
```

## Testing

```bash
make test                              # all tests
make test-service SERVICE=Bill         # specific service
make test-package PKG=./internal/Bill  # specific package
make test-coverage                     # coverage report
```

See [docs/testing/README.md](../testing/README.md) for the full testing guide.

## SSL Certificates

All services use mutual TLS (mTLS). Certificates live in `ssl/` and are loaded at startup.

### Certificate Files

| File | Description | Used by |
|------|-------------|---------|
| `ssl/ca.crt` | CA certificate | All services (client verification) |
| `ssl/ca.key` | CA private key | Certificate signing only |
| `ssl/server.crt` + `ssl/server.key` | Server cert for `egot.internal.com` | core, BRS, Bill, DCAP, … |
| `ssl/srv.crt` + `ssl/srv.key` | Same server cert (copies) | flowreservation, operator |
| `ssl/client.crt` + `ssl/client.key` | Default client cert (CN=user-test) | crawler, client cmd |
| `ssl/client-NNNN.crt` + `ssl/client-NNNN.key` | Numbered client certs | Auto-loaded by `InitRepository` |

### Refreshing Certificates

When certs expire or after initial clone:

```shell
make ssl-refresh
```

Generates a new self-signed ECDSA P-256 CA (10-year validity) and re-issues all base certs with a 1-year validity.

### Generating Client Certificates

```shell
make ssl-clients N=5
# Creates ssl/client-0001.crt ... ssl/client-0005.crt (and matching .key files)
```

Each cert gets CN=`user-NNNN`. The `.crt` extension is required for `InitRepository` to auto-load them on service startup.

### How Services Load Client Certs

`InitRepository("./ssl")` walks the `ssl/` directory and loads all `*.crt` files whose path contains `client`. This populates the in-memory repository with known client identities (identified by LFDI — the first 40 hex chars of the cert's SHA-256 fingerprint).
