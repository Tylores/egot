# Microservice Scaffold Generator

The scaffold generator (`cmd/scaffold-gen`) creates complete microservices from WADL (Web Application Description Language) XML specifications. Each generated service follows the established pattern from `core` and `flowreservation` services.

## Quick Start

### 1. Create a WADL specification file

```bash
cat > wadl/my-service.wadl << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<application name="my-service" port="8003" max_entities="50">
  <resources>
    <resource path="/items">
      <methods>
        <method name="GetItems" http_method="GET" response_type="ItemList"/>
        <method name="CreateItem" http_method="POST" request_type="Item" response_type="Item"/>
      </methods>
    </resource>
    <resource path="/items/{id}">
      <methods>
        <method name="GetItem" http_method="GET" response_type="Item">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
        <method name="UpdateItem" http_method="PUT" request_type="Item" response_type="Item">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
EOF
```

### 2. Generate the scaffold

```bash
go run ./cmd/scaffold-gen -wadl wadl/my-service.wadl -output .
```

This creates:
- `cmd/my-service/main.go` - Service entry point with TLS and route setup
- `internal/my-service/handler/handler.go` - HTTP handler stubs for each method
- `internal/my-service/repository/memory/memory.go` - In-memory storage
- `internal/my-service/repository/error.go` - Error definitions
- `internal/my-service/server/server.go` - Server utilities
- Updates to `internal/routes/routes.go` with the new service port

### 3. Implement business logic

Edit the generated handler and repository files to add your service logic.

**handler/handler.go:**
```go
func (h *Handler) GetItems(w http.ResponseWriter, req *http.Request) {
    // Replace nil with actual repository call
    items := h.repo.GetItems()
    w.Header().Set("Content-Type", sep.ContentType)
    xml.NewEncoder(w).Encode(items)
}
```

**repository/memory/memory.go:**
```go
type Pool struct {
    items []sep.Item  // Add your resource types
}

func (r *Repository) GetItems() []sep.Item {
    r.RLock()
    defer r.RUnlock()
    return r.pool.items
}
```

### 4. Build and run

```bash
go build ./cmd/my-service
./my-service
```

## WADL Reference

See [SCAFFOLD_SPEC.md](../SCAFFOLD_SPEC.md) for complete WADL specification and examples.

### Key Elements

- **application**: Root element with `name`, `port`, and optional `max_entities`
- **resource**: REST resource with `path` attribute
- **method**: HTTP method with `name`, `http_method`, `response_type`, and optional `request_type`
- **path_params**: Path parameters with `name` and `type` (int, string)

## Architecture

Generated services follow this pattern:

1. **Handler** - Extracts client certificate (LFDI), validates entity, processes request
2. **Repository** - In-memory storage with sync.RWMutex for thread safety
3. **Main** - Sets up mutual TLS, initializes repository from certificates, registers routes

All services use:
- Mutual TLS with client certificate verification
- XML encoding for requests/responses (sep.ContentType)
- Entity-based access control (LFDI → Entity ID)

## Examples

See `wadl/device-manager.wadl` for a complete example WADL file.

## Tips

- Service names should be kebab-case (e.g., "device-manager")
- Port numbers should not conflict with existing services (8000-8001 are taken)
- Method names should be descriptive (GetItem, CreateItem, UpdateItem, etc.)
- Path parameters must match placeholders: `{id}` → `<param name="id"/>`
- Keep max_entities reasonable for memory usage (default: 100)
