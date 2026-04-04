# Microserver Scaffold Generator - Specification

## Overview

The Microserver Scaffold Generator creates complete microservices from WADL (Web Application Description Language) XML specifications. Each generated service follows the established pattern: isolated `internal/{service}/` packages with `handler`, `repository`, and `server` components, mirroring the structure of existing services like `core` and `flowreservation`.

## Architecture Pattern

Each microservice follows this structure:

```
cmd/{service}/
  main.go                    # Entry point: TLS setup, repository init, route registration

internal/{service}/
  handler/
    handler.go              # HTTP handlers for each resource method
  repository/
    memory/
      memory.go            # In-memory storage implementation
    error.go               # Repository errors
  server/
    server.go              # Server setup utilities

internal/routes/
  routes.go                # Service port constants (auto-updated)
```

## WADL Specification Format

Services are defined using WADL (Web Application Description Language) XML files. WADL is a machine-readable XML format for describing HTTP-based web applications.

### WADL XML Schema

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="service-name" port="8002" max_entities="50">
  <resources>
    <resource path="/resource-path">
      <methods>
        <method name="MethodName" http_method="GET|POST|PUT|DELETE" 
                request_type="TypeName" response_type="TypeName">
          <path_params>
            <param name="id" type="int|string"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

### Attributes & Elements

**Application Element:**
- `name` (required): Service name in kebab-case (e.g., "device-manager")
- `port` (required): Service listening port (e.g., 8002)
- `max_entities` (optional): Maximum concurrent entities in memory pool (default: 100)

**Resource Element:**
- `path` (required): URI path for this resource (supports `{id}` style path parameters)

**Method Element:**
- `name` (required): Method function name (e.g., "GetDevice", "CreateDevice")
- `http_method` (required): HTTP verb (GET, POST, PUT, DELETE)
- `request_type` (optional): Type name for request body (POST/PUT methods)
- `response_type` (required): SEP type name or custom type for response

**Param Element** (child of path_params):
- `name` (required): Parameter name (must match path placeholder, e.g., {id})
- `type` (required): Parameter type (int, string)

### WADL Example

See `wadl/device-manager.wadl` for a complete example:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="device-manager" port="8002" max_entities="50">
  <resources>
    <resource path="/capability">
      <methods>
        <method name="GetCapability" http_method="GET" response_type="DeviceCapability"/>
      </methods>
    </resource>
    
    <resource path="/devices">
      <methods>
        <method name="GetDevices" http_method="GET" response_type="EndDeviceList"/>
        <method name="CreateDevice" http_method="POST" 
                request_type="EndDevice" response_type="EndDevice"/>
      </methods>
    </resource>
    
    <resource path="/devices/{id}">
      <methods>
        <method name="GetDevice" http_method="GET" response_type="EndDevice">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
        <method name="UpdateDevice" http_method="PUT" 
                request_type="EndDevice" response_type="EndDevice">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

## Using the Scaffold Generator

### 1. Create WADL Specification

Create an XML file describing your microservice resources:

```bash
cat > wadl/my-service.wadl << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<application name="my-service" port="8003" max_entities="100">
  <resources>
    <resource path="/status">
      <methods>
        <method name="GetStatus" http_method="GET" response_type="Status"/>
      </methods>
    </resource>
  </resources>
</application>
EOF
```

### 2. Run Scaffold Generator

```bash
go run ./cmd/scaffold-gen -wadl wadl/my-service.wadl -output .
```

The generator will create:
- `cmd/my-service/main.go` - Service entry point
- `internal/my-service/handler/handler.go` - HTTP handlers
- `internal/my-service/repository/memory/memory.go` - Data storage
- `internal/my-service/repository/error.go` - Error definitions
- `internal/my-service/server/server.go` - Server utilities
- Updated `internal/routes/routes.go` with new service port

### 3. Implement Business Logic

Edit the generated files to add your service logic:

**handler/handler.go:** Implement the handler methods
```go
func (h *Handler) GetStatus(w http.ResponseWriter, req *http.Request) {
    // Your implementation here
    status := h.repo.GetStatus(entityID)
    w.Header().Set("Content-Type", sep.ContentType)
    xml.NewEncoder(w).Encode(status)
}
```

**repository/memory/memory.go:** Add storage for your resources
```go
type Pool struct {
    status []sep.Status  // Add slice for each resource type
}
```

### 4. Build & Run

```bash
go build ./cmd/my-service
./my-service
```

## Generated Handler Pattern

For each method in the WADL, the generator creates a handler function following this pattern:

```go
func (h *Handler) GetDevice(w http.ResponseWriter, req *http.Request) {
    // 1. Extract and verify client certificate
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    // 2. Validate entity exists in repository
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        log.Printf("Repository get error: %v\n", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    // 3. Extract and validate path parameters (if any)
    id, err := strconv.Atoi(req.PathValue("id"))
    if err != nil {
        log.Printf("path id value error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    // 4. TODO: Retrieve data from repository
    // 5. TODO: Encode response and return
    w.Header().Set("Content-Type", sep.ContentType)
    err = xml.NewEncoder(w).Encode(nil)  // Replace nil with actual data
    if err != nil {
        log.Printf("Response encode error: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}
```

## Generated Repository Pattern

For each resource type in the WADL, storage methods are generated:

```go
type Repository struct {
    sync.RWMutex
    tag_lookup      map[string]Entity        // LFDI -> Entity ID
    active_entities []bool                   // Track allocated entities
    pool            Pool                     // Data storage
}

type Pool struct {
    devices []sep.EndDevice    // One slice per resource type
    status  []sep.Status       // Add more as needed
}

// Generated getters and setters (implement in memory.go)
func (r *Repository) GetDevice(id Entity) sep.EndDevice {
    r.RLock()
    defer r.RUnlock()
    return r.pool.devices[id]
}

func (r *Repository) PutDevice(id Entity, data sep.EndDevice) {
    r.Lock()
    defer r.Unlock()
    r.pool.devices[id] = data
}
```

## Generated cmd/main.go Pattern

The entry point generated from WADL will look like:

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
    
    "github.com/Tylores/egot/internal/device-manager/handler"
    "github.com/Tylores/egot/internal/device-manager/repository/memory"
    "github.com/Tylores/egot/internal/routes"
)

const MAX_ENTITIES memory.Entity = 50

func main() {
    cfg := &tls.Config{
        MinVersion: tls.VersionTLS12,
        ClientAuth: tls.RequireAndVerifyClientCert,
    }
    server := http.Server{
        Addr:      routes.DeviceManager,  // Auto-populated from WADL
        TLSConfig: cfg,
    }
    
    repo := memory.NewRepository(MAX_ENTITIES)
    repo.InitRepository("./ssl")
    
    h := handler.NewHandler(repo)
    
    // Routes auto-generated from WADL resources
    http.Handle("GET /devices", http.HandlerFunc(h.GetDevices))
    http.Handle("POST /devices", http.HandlerFunc(h.CreateDevice))
    http.Handle("GET /devices/{id}", http.HandlerFunc(h.GetDevice))
    
    err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Integration with Core Services

Each generated service:
- Uses mutual TLS (client certificate verification) for security
- Implements the Entity/Repository pattern from core service
- Uses XML encoding matching `sep.ContentType`
- Registers its port in `internal/routes/routes.go`
- Follows the same directory structure for consistency

## Notes

- All services listen on mutual TLS with required client certificates
- Request/response handling uses XML encoding (SEP format)
- Entity lifecycle: Client cert LFDI → Entity ID lookup → Data retrieval
- Memory repository uses RWMutex for concurrent read/serialized write safety
- Service ports are constants in `internal/routes/routes.go`
- Path parameters must match placeholder names in resource paths
