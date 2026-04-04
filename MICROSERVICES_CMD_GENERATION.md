# Microservice CMD Directories Generated

## Overview

All 15 microservices in the `internal/` directory now have corresponding `cmd/` entry point directories.

### Services Generated

| Service | Port | CMD Directory | Internal Package | Status |
|---------|------|---------------|------------------|--------|
| BRS | 8010 | cmd/BRS | internal/BRS | ✅ Created |
| Bill | 8011 | cmd/Bill | internal/Bill | ✅ Created |
| DCAP | 8012 | cmd/DCAP | internal/DCAP | ✅ Created |
| DERP | 8013 | cmd/DERP | internal/DERP | ✅ Created |
| DR | 8014 | cmd/DR | internal/DR | ✅ Created |
| EDevice | 8015 | cmd/EDevice | internal/EDevice | ✅ Created |
| File | 8016 | cmd/File | internal/File | ✅ Created |
| MUP | 8017 | cmd/MUP | internal/MUP | ✅ Created |
| Messaging | 8018 | cmd/Messaging | internal/Messaging | ✅ Created |
| Notify | 8019 | cmd/Notify | internal/Notify | ✅ Created |
| PPY | 8020 | cmd/PPY | internal/PPY | ✅ Created |
| SDevice | 8021 | cmd/SDevice | internal/SDevice | ✅ Created |
| TariffProfile | 8022 | cmd/TariffProfile | internal/TariffProfile | ✅ Created |
| TimeOfUse | 8023 | cmd/TimeOfUse | internal/TimeOfUse | ✅ Created |
| UPT | 8024 | cmd/UPT | internal/UPT | ✅ Created |

Plus 8 original services (core, crawler, flowreservation, operator, rsps, client, scaffold-gen, wadl-extract).

## Structure

Each generated cmd directory contains:

```
cmd/<Service>/
└── main.go                 # Entry point with TLS setup and route registration
```

### Example Structure (cmd/Bill/main.go)

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
    
    "github.com/Tylores/egot/internal/Bill/handler"
    "github.com/Tylores/egot/internal/Bill/repository/memory"
)

const MAX_ENTITIES memory.Entity = 100

func main() {
    cfg := &tls.Config{
        MinVersion: tls.VersionTLS12,
        ClientAuth: tls.RequireAndVerifyClientCert,
    }
    server := http.Server{
        Addr:      "egot.internal.com:8011",
        TLSConfig: cfg,
    }
    
    repo := memory.NewRepository(MAX_ENTITIES)
    repo.InitRepository("./ssl")
    
    h := handler.NewHandler(repo)
    http.HandleFunc("/", h.ServeHTTP)
    
    err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Routes Configuration

All service routes are now defined in `internal/routes/routes.go`:

```go
const (
    BRS             = "egot.internal.com:8010"
    Bill            = "egot.internal.com:8011"
    DCAP            = "egot.internal.com:8012"
    // ... all 15 services
    UPT             = "egot.internal.com:8024"
)
```

## Building Services

### Build Individual Service

```bash
go build -o ./bin/bill ./cmd/Bill
go build -o ./bin/messaging ./cmd/Messaging
# etc.
```

### Build All Services

```bash
for service in BRS Bill DCAP DERP DR EDevice File MUP Messaging Notify PPY SDevice TariffProfile TimeOfUse UPT; do
    go build -o ./bin/$service ./cmd/$service
done
```

### Or use Make

```bash
make test-service SERVICE=bill
make test-service SERVICE=messaging
```

## Testing Services

Each service now has tests set up with the testing framework:

```bash
# Test all services
make test

# Test specific services
make test-service SERVICE=bill
make test-service SERVICE=messaging
make test-service SERVICE=tariffprofile
# etc.

# Test specific package
make test-package PKG=./internal/Bill
```

## Customization

### Update Handler Routes

Edit each service's `cmd/<Service>/main.go` to register appropriate HTTP handlers:

```go
h := handler.NewHandler(repo)

// Add service-specific route handlers
http.HandleFunc("/api/bills", h.GetBills)
http.HandleFunc("/api/bills/{id}", h.GetBill)
http.HandleFunc("/api/bills", h.CreateBill)
```

### Update Ports

If you need to change port assignments, edit:
1. `internal/routes/routes.go` - Update constants
2. `cmd/<Service>/main.go` - Update server.Addr
3. DNS/network configuration for the service endpoints

### Add More Handler Methods

Edit the handler files in `internal/<Service>/handler/handler.go` and wire them up in the cmd/main.go.

## Verification

All generated services compile successfully:

```bash
✓ BRS - compiles OK
✓ Bill - compiles OK
✓ DCAP - compiles OK
✓ DERP - compiles OK
✓ DR - compiles OK
✓ EDevice - compiles OK
✓ File - compiles OK
✓ MUP - compiles OK
✓ Messaging - compiles OK
✓ Notify - compiles OK
✓ PPY - compiles OK
✓ SDevice - compiles OK
✓ TariffProfile - compiles OK
✓ TimeOfUse - compiles OK
✓ UPT - compiles OK
```

## Next Steps

1. **Test individual services**
   ```bash
   make test-service SERVICE=bill
   make test-service SERVICE=messaging
   ```

2. **Customize handler routes** in each `cmd/<Service>/main.go`

3. **Add business logic** to `internal/<Service>/handler/handler.go`

4. **Implement storage** in `internal/<Service>/repository/memory/repository.go`

5. **Build services**
   ```bash
   go build -o ./bin/bill ./cmd/Bill
   ```

## Summary

✅ All 15 microservices have cmd entry points  
✅ All services compile successfully  
✅ Routes configured in internal/routes/routes.go  
✅ TLS setup ready in each main.go  
✅ Repository initialization configured  
✅ Testing framework ready for all services  

Your microservices are now ready for development and deployment!
