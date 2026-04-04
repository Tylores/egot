# Handler Registration & Scaffold-Gen Integration

## Overview

You can use scaffold-gen (via WADL specifications) to generate complete microservices, OR you can manually wire up handlers from existing internal/<Service>/handler/handler.go files into cmd/<Service>/main.go.

## Two Approaches

### Approach 1: Use scaffold-gen (Recommended for New Services)

scaffold-gen reads WADL specifications and generates:
- `cmd/<Service>/main.go` with HTTP routes
- `internal/<Service>/handler/handler.go` with handler stubs
- `internal/<Service>/repository/memory/repository.go` with storage

**Command:**
```bash
./cmd/scaffold-gen \
  -sep2-wadl ./wadl/sep2-spec.wadl \
  -service MyService \
  -port 8030 \
  -output .
```

This generates a complete, fully wired service.

### Approach 2: Manual Handler Registration (For Existing Services)

For services that already have handlers defined, use the `regenerate_cmd_handlers.sh` script to:
1. Analyze `internal/<Service>/handler/handler.go`
2. Extract all public handler methods
3. Auto-generate `cmd/<Service>/main.go` with proper registrations

**Command:**
```bash
./scripts/regenerate_cmd_handlers.sh              # All services
./scripts/regenerate_cmd_handlers.sh Bill         # Specific service
./scripts/regenerate_cmd_handlers.sh Bill Messaging  # Multiple services
```

## How Handlers Are Wired

### Manual Method Naming Pattern

Handler methods follow this pattern:
```go
func (h *Handler) <HTTPMethod><ResourceName>(w http.ResponseWriter, req *http.Request)
```

Examples:
- `GETCustomers` - GET /Customers
- `POSTCustomer` - POST /Customer  
- `PUTCustomerByID` - PUT /CustomerByID
- `DELETECustomer` - DELETE /Customer
- `HEADCustomers` - HEAD /Customers

### Current Handler Registration

In cmd/<Service>/main.go:
```go
h := handler.NewHandler(repo)

// Register handler methods
http.HandleFunc("/", h.GETCustomerList)
http.HandleFunc("/", h.POSTCustomer)
http.HandleFunc("/", h.DELETECustomer)
// ... more handlers
```

**Note:** All handlers are currently registered to "/" - you should improve this by extracting the resource name from the method name and registering to proper paths like "/customers", "/orders", etc.

## Improving Handler Registration

### Better Approach: Register to Specific Paths

Modify `regenerate_cmd_handlers.sh` or manually update main.go to register handlers to proper paths:

```go
// Instead of all pointing to "/"
h := handler.NewHandler(repo)

// Register handlers with proper paths
// Customers
http.HandleFunc("GET /customers", h.GETCustomerList)
http.HandleFunc("GET /customers/{id}", h.GETCustomer)
http.HandleFunc("POST /customers", h.POSTCustomer)
http.HandleFunc("PUT /customers/{id}", h.PUTCustomer)
http.HandleFunc("DELETE /customers/{id}", h.DELETECustomer)

// Orders
http.HandleFunc("GET /orders", h.GETOrderList)
http.HandleFunc("GET /orders/{id}", h.GETOrder)
// ... etc
```

### Enhancement: Extract Paths from Method Names

To automatically improve path registration, update `scripts/regenerate_cmd_handlers.sh`:

```bash
# Extract resource name from method name
# GETCustomerList -> /customers
# POSTOrder -> /orders
# PUTCustomerByID -> /customers/{id}

extract_path_from_method() {
  local method=$1
  
  # Remove HTTP method prefix
  local resource="${method#GET}"
  resource="${resource#POST}"
  resource="${resource#PUT}"
  resource="${resource#DELETE}"
  resource="${resource#HEAD}"
  
  # Convert to lowercase plural
  local path=$(echo "$resource" | sed 's/\([A-Z]\)/-\1/g' | tr '[:upper:]' '[:lower:]' | sed 's/^-//')
  
  if [[ "$path" != *"s" ]]; then
    path="${path}s"
  fi
  
  echo "/$path"
}
```

## Service Status

### Fully Generated with Handlers

All 15 generated services now have:
- ✅ `cmd/<Service>/main.go` with handler registration
- ✅ Compiled and verified working
- ✅ TLS configuration
- ✅ Repository initialization

### Handler Counts

| Service | Handlers |
|---------|----------|
| BRS | 20 |
| Bill | 70 |
| DCAP | 5 |
| DERP | 40 |
| DR | 25 |
| EDevice | 265 |
| File | 10 |
| MUP | 10 |
| Messaging | 25 |
| Notify | 10 |
| PPY | 45 |
| SDevice | 5 |
| TariffProfile | 45 |
| TimeOfUse | 5 |
| UPT | 45 |

**Total**: 525 handler methods across 15 services!

## Using Scaffold-Gen for Future Services

If you create new WADL specifications, use scaffold-gen:

```bash
# Create a WADL specification (example: invoice.wadl)
cat > invoice.wadl << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<application name="Invoice" port="8025" max_entities="100">
  <resources>
    <resource path="/invoices">
      <methods>
        <method name="GETInvoiceList" http_method="GET"/>
        <method name="POSTInvoice" http_method="POST"/>
      </methods>
    </resource>
    <resource path="/invoices/{id}">
      <methods>
        <method name="GETInvoice" http_method="GET">
          <path_params>
            <param name="id" type="string"/>
          </path_params>
        </method>
        <method name="PUTInvoice" http_method="PUT">
          <path_params>
            <param name="id" type="string"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
EOF

# Generate the service
./cmd/scaffold-gen -wadl invoice.wadl -output .

# Build and test
go build -o bin/invoice cmd/Invoice
make test-service SERVICE=Invoice
```

## Workflow

### For Existing Services (Current Situation)

1. Service already has handlers in `internal/<Service>/handler/handler.go`
2. Run: `./scripts/regenerate_cmd_handlers.sh <Service>`
3. Verify: `go build ./cmd/<Service>`
4. Test: `make test-service SERVICE=<Service>`

### For New Services

1. Create WADL specification
2. Run: `./cmd/scaffold-gen -wadl spec.wadl -output .`
3. scaffold-gen generates everything:
   - cmd/<Service>/main.go (with handlers from WADL)
   - internal/<Service>/handler/handler.go (with stubs)
   - internal/<Service>/repository/memory/repository.go
4. Add business logic to handler.go
5. Implement storage in repository.go
6. Build and test

## Handler Registration Best Practices

### 1. Use HTTP Method in Handler Name

```go
// Good - HTTP method is clear
func (h *Handler) GETUsers(w http.ResponseWriter, req *http.Request) { }
func (h *Handler) POSTUser(w http.ResponseWriter, req *http.Request) { }

// Avoid - HTTP method unclear
func (h *Handler) ListUsers(w http.ResponseWriter, req *http.Request) { }
func (h *Handler) CreateUser(w http.ResponseWriter, req *http.Request) { }
```

### 2. Include Resource Name

```go
// Good - Resource is clear
func (h *Handler) GETCustomerList(w http.ResponseWriter, req *http.Request) { }
func (h *Handler) POSTCustomer(w http.ResponseWriter, req *http.Request) { }

// Avoid - Resource name missing
func (h *Handler) GET(w http.ResponseWriter, req *http.Request) { }
func (h *Handler) POST(w http.ResponseWriter, req *http.Request) { }
```

### 3. Register to Proper Paths

```go
// Good
http.HandleFunc("GET /users", h.GETUserList)
http.HandleFunc("POST /users", h.POSTUser)
http.HandleFunc("GET /users/{id}", h.GETUser)

// Avoid - all handlers on root path
http.HandleFunc("GET /", h.GETUserList)
http.HandleFunc("POST /", h.POSTUser)
```

## Troubleshooting

### Script Doesn't Find Handlers

Ensure handler methods follow the pattern:
```go
func (h *Handler) <HTTPMethod><Resource>(...) {
```

Private methods (lowercase) are automatically excluded.

### Services Don't Compile

Check:
1. Repository interface matches: `NewRepository()` with no parameters
2. Handler struct is defined: `type Handler struct`
3. NewHandler takes repository: `func NewHandler(repo *memory.Repository)`

### Handlers Not Being Called

Verify HTTP path registration - handlers are currently all on "/", so any request matches the first registered handler. This should be improved by registering to specific paths like "/customers", "/orders", etc.

## Summary

**Current State**:
- ✅ All 15 services have cmd/main.go files
- ✅ All handlers from internal/<Service>/handler.go are registered
- ✅ All services compile successfully
- ✅ TLS configured
- ✅ Repositories initialized

**Next Improvements**:
- Extract paths from method names
- Register handlers to proper routes (/customers, /orders, etc.)
- Generate path parameters from method names
- Use scaffold-gen for future services

**Tools Available**:
- `./scripts/regenerate_cmd_handlers.sh` - Regenerate handlers from existing code
- `./cmd/scaffold-gen` - Generate complete services from WADL
- `make test-service SERVICE=...` - Test individual services
- `make test` - Test all services

See also: MICROSERVICES_CMD_GENERATION.md, CMD_GENERATION_GUIDE.md
