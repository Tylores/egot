# Microservice Scaffold Generator

A Go-based tool for generating complete microservices from WADL (Web Application Description Language) XML specifications.

## Overview

The scaffold generator creates new microservices that follow the established patterns of existing services (`core`, `flowreservation`). Each generated service includes:

- **Handler layer** - HTTP request processing with TLS client certificate validation
- **Repository layer** - In-memory data storage with thread-safe access
- **Server setup** - Mutual TLS configuration and route registration
- **Main entry point** - Complete service startup code

## Quick Start

### Generate a new service:

```bash
go run ./cmd/scaffold-gen -wadl wadl/my-service.wadl -output .
```

### What you get:

```
cmd/my-service/
  main.go                          # Service entry point

internal/my-service/
  handler/handler.go              # HTTP handlers
  repository/memory/memory.go     # Data storage
  repository/error.go             # Error types
  server/server.go                # Server utilities
```

### Next steps:

1. Implement business logic in `internal/{service}/handler/handler.go`
2. Add storage methods in `internal/{service}/repository/memory/memory.go`
3. Build: `go build ./cmd/{service}`
4. Run: `./{service}`

## WADL Specification

Services are defined in XML using the WADL format. Here's a minimal example:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="my-service" port="8003" max_entities="50">
  <resources>
    <resource path="/items">
      <methods>
        <method name="GetItems" http_method="GET" response_type="ItemList"/>
      </methods>
    </resource>
    <resource path="/items/{id}">
      <methods>
        <method name="GetItem" http_method="GET" response_type="Item">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

**WADL Elements:**

| Element | Attribute | Type | Description |
|---------|-----------|------|-------------|
| application | name | string | Service name (kebab-case) |
| application | port | int | Service port |
| application | max_entities | int | Max entities in pool (default: 100) |
| resource | path | string | Resource URI path with `{param}` support |
| method | name | string | Handler method name (CamelCase) |
| method | http_method | string | GET, POST, PUT, DELETE |
| method | response_type | string | Response type name |
| method | request_type | string | Request body type (POST/PUT) |
| param | name | string | Parameter name |
| param | type | string | int or string |

## Generated Code Patterns

### Handler Pattern

Each method gets a handler stub:

```go
func (h *Handler) GetItem(w http.ResponseWriter, req *http.Request) {
    // 1. Extract and verify client certificate (LFDI)
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    
    // 2. Validate entity exists
    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    // 3. Extract path parameters if present
    id, _ := strconv.Atoi(req.PathValue("id"))
    
    // 4. TODO: Get data from repository
    // 5. TODO: Encode response
    w.Header().Set("Content-Type", sep.ContentType)
    xml.NewEncoder(w).Encode(nil)
}
```

### Repository Pattern

Generated repository with entity management:

```go
type Repository struct {
    sync.RWMutex
    tag_lookup      map[string]Entity    // LFDI -> Entity ID
    active_entities []bool               // Entity allocation tracking
    pool            Pool                 // Resource storage
}

type Pool struct {
    items []sep.Item  // TODO: Add resource types
}
```

### Main Entry Point

```go
func main() {
    // TLS configuration
    cfg := &tls.Config{
        MinVersion: tls.VersionTLS12,
        ClientAuth: tls.RequireAndVerifyClientCert,
    }
    
    // Initialize repository from client certificates
    repo := memory.NewRepository(MAX_ENTITIES)
    repo.InitRepository("./ssl")
    
    // Register routes
    h := handler.NewHandler(repo)
    http.Handle("GET /items", http.HandlerFunc(h.GetItems))
    
    // Start server
    server := &http.Server{
        Addr:      routes.MyService,
        TLSConfig: cfg,
    }
    server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
}
```

## Usage Examples

### Example 1: Simple Service

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="status-service" port="8004" max_entities="10">
  <resources>
    <resource path="/status">
      <methods>
        <method name="GetStatus" http_method="GET" response_type="Status"/>
      </methods>
    </resource>
  </resources>
</application>
```

### Example 2: Full CRUD Service

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="device-manager" port="8005" max_entities="100">
  <resources>
    <resource path="/devices">
      <methods>
        <method name="ListDevices" http_method="GET" response_type="DeviceList"/>
        <method name="CreateDevice" http_method="POST" request_type="Device" response_type="Device"/>
      </methods>
    </resource>
    <resource path="/devices/{id}">
      <methods>
        <method name="GetDevice" http_method="GET" response_type="Device">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
        <method name="UpdateDevice" http_method="PUT" request_type="Device" response_type="Device">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
        <method name="DeleteDevice" http_method="DELETE" response_type="void">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

## Project Structure

```
.
├── cmd/
│   ├── core/              # Existing core service
│   ├── flowreservation/   # Existing flow service
│   └── scaffold-gen/      # The scaffold generator
│       ├── main.go        # CLI entry point
│       └── README.md      # Scaffold generator docs
├── internal/
│   ├── core/              # Core service implementation
│   ├── flow/              # Flow service implementation
│   ├── routes/            # Service port registry
│   └── scaffold/          # Generator implementation
│       ├── generator.go   # Main generator logic
│       └── generator_templates.go  # Code templates
├── wadl/                  # WADL specification files
│   └── device-manager.wadl  # Example WADL
├── SCAFFOLD_SPEC.md       # Complete specification
└── cmd/scaffold-gen/README.md  # Quick start guide
```

## Architecture Principles

1. **Isolated Services** - Each service has its own packages; no shared code
2. **Mutual TLS** - All services require valid client certificates
3. **Entity-Based Access** - Client certificate LFDI maps to Entity ID
4. **Thread-Safe Storage** - sync.RWMutex for concurrent access
5. **XML Serialization** - All requests/responses use SEP XML format

## Documentation

- **[SCAFFOLD_SPEC.md](./SCAFFOLD_SPEC.md)** - Complete specification with WADL schema
- **[cmd/scaffold-gen/README.md](./cmd/scaffold-gen/README.md)** - Scaffold generator quick start
- **[wadl/device-manager.wadl](./wadl/device-manager.wadl)** - Example WADL file

## Implementation Tips

- Keep service names in kebab-case (device-manager, flow-service)
- Use descriptive method names (GetDevice, UpdateDevice, DeleteDevice)
- Plan entity pool size based on expected concurrent clients
- Implement repository methods for each resource type in Pool
- Replace `nil` placeholder responses with actual data from repository
- Test with mutual TLS certificates in ./ssl directory

## Troubleshooting

**Build errors with "declared and not used":**
- These are placeholder TODO comments where you'll add implementation
- Replace `nil` responses with actual repository calls

**Server fails to start:**
- Verify SSL certificates exist in ./ssl directory
- Ensure port is not already in use
- Check that routes constant was added to internal/routes/routes.go

**TLS certificate errors:**
- Ensure client certificates are in ./ssl with "client" in filename
- Certificates must be valid and current
- Check certificate permissions (should be readable by service)
