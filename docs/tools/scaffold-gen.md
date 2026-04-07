# scaffold-gen

Generates complete microservices from SEP 2 WADL specifications, following the same patterns as all existing services.

## Workflow

```bash
# 1. Extract a service WADL from the SEP 2 spec
go run ./cmd/wadl-extract -wadl wadl/sep_wadl.xml -path /dcap -output wadl/dcap.wadl

# 2. Generate (run from project root; filename = service name)
go run ./cmd/scaffold-gen -wadl wadl/dcap.wadl

# 3. Build
go build ./cmd/dcap
```

## What Gets Generated

```
cmd/dcap/
  main.go                          # Service entry point with TLS setup + route registration

internal/dcap/
  handler/handler.go               # HTTP handler stubs (matching WADL methods)
  repository/memory/memory.go      # In-memory storage
  repository/error.go              # Error types
  server/server.go                 # Server utilities
```

`internal/routes/routes.go` is also updated with the new service port constant.

## WADL Format

scaffold-gen reads the SEP 2 WADL format produced by `wadl-extract`. The relevant attributes are:

**`<application>`**
- `port` (required) — added by `wadl-extract`; must be set before running scaffold-gen
- `max_entities` (optional) — entity pool size, default 100

**`<resource>`**
- `wx:samplePath` — resource URI (e.g., `/brs/{id1}/br/{id2}`); path params are extracted automatically

**`<method>`**
- `id` — becomes the handler function name (e.g., `GETDeviceCapability`)
- `name` — HTTP verb (`GET`, `POST`, `PUT`, `DELETE`, `HEAD`)
- `wx:mode` — drives generated behaviour (see table below)

**`<representation>`** (inside `<response>`)
- `element` — SEP type (e.g., `sep:DeviceCapability`); stripped to `sep.DeviceCapability{}`

## Generated handler behaviour

| `wx:mode` | HTTP verb | Generated stub |
|-----------|-----------|----------------|
| `E` | any | `StatusMethodNotAllowed` |
| `M` / `D` | `GET` | `StatusOK` + `xml.NewEncoder(w).Encode(&sep.Type{})` |
| `M` / `D` | `HEAD` | `StatusOK` |
| `M` / `D` | `POST` | `StatusCreated` + `location` header |
| `M` / `D` | `DELETE` | `StatusOK` + encode |
| `M` / `D` | `PUT` | `StatusOK` (stub) |

## Generated Patterns

### Handler stub

```go
// getLFDI helper is generated once
func (h *Handler) getLFDI(req *http.Request) (string, error) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
    if _, err := h.repo.GetEntity(lfdi); err != nil {
        return "", err
    }
    return lfdi, nil
}

// DeviceCapability resource handlers

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

func (h *Handler) PUTDeviceCapability(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed) // mode E
}
```

### Repository stub

```go
type Repository struct {
    sync.RWMutex
    tag_lookup      map[string]Entity
    active_entities []bool
    pool            Pool
}

type Pool struct {
    // TODO: add slices for each resource type
    // items []sep.Item
}
```

## Architecture Principles

- **Isolated services** — each service has its own packages with no cross-service imports
- **Mutual TLS** — all services require valid client certificates (`RequireAndVerifyClientCert`)
- **Entity-based access** — client cert LFDI maps to an entity in the repository pool
- **Thread-safe storage** — `sync.RWMutex` for concurrent read / serialised write
- **XML serialisation** — all responses use `sep.ContentType` (`application/sep+xml`)

## Troubleshooting

**`WADL must specify application port`:**
Re-run `wadl-extract` — the updated tool now injects `port` and `max_entities` into the output. If editing the WADL manually, add `port="8XXX"` to the `<application>` tag.

**Server fails to start:**
- Verify SSL certificates exist in `./ssl`
- Ensure the port is not already in use
- Confirm the routes constant was added to `internal/routes/routes.go`

**TLS certificate errors:**
- Client certificates must be in `./ssl` with `client` in the filename
- Certificates must be valid and readable by the service


## What Gets Generated

```
cmd/my-service/
  main.go                          # Service entry point with TLS setup + route registration

internal/my-service/
  handler/handler.go               # HTTP handler stubs
  repository/memory/memory.go      # In-memory storage
  repository/error.go              # Error types
  server/server.go                 # Server utilities
```

`internal/routes/routes.go` is also updated with the new service port constant.

## WADL Specification

### Schema

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application port="8002" max_entities="50">
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

### Attributes

**`<application>`**
- `port` (required) — service port (e.g., `8002`)
- `max_entities` (optional) — concurrent entity pool size (default: 100)

**`<resource>`**
- `path` (required) — URI path, supports `{param}` placeholders

**`<method>`**
- `name` (required) — handler function name (CamelCase)
- `http_method` (required) — `GET`, `POST`, `PUT`, or `DELETE`
- `request_type` (optional) — request body type (POST/PUT)
- `response_type` (required) — response type name

**`<param>`** (child of `<path_params>`)
- `name` (required) — must match placeholder in `resource.path`
- `type` (required) — `int` or `string`

### Examples

**Minimal (single endpoint):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<application port="8004" max_entities="10">
  <resources>
    <resource path="/status">
      <methods>
        <method name="GetStatus" http_method="GET" response_type="Status"/>
      </methods>
    </resource>
  </resources>
</application>
```

**Full CRUD:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<application port="8005" max_entities="100">
  <resources>
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
          <path_params><param name="id" type="int"/></path_params>
        </method>
        <method name="UpdateDevice" http_method="PUT"
                request_type="EndDevice" response_type="EndDevice">
          <path_params><param name="id" type="int"/></path_params>
        </method>
        <method name="DeleteDevice" http_method="DELETE" response_type="void">
          <path_params><param name="id" type="int"/></path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

## Generated Patterns

### Handler stub

```go
func (h *Handler) GetItem(w http.ResponseWriter, req *http.Request) {
    cert := req.TLS.PeerCertificates[0]
    lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

    _, err := h.repo.GetEntity(lfdi)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    id, err := strconv.Atoi(req.PathValue("id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    _ = id

    // TODO: retrieve from repository
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)
    xml.NewEncoder(w).Encode(nil) // replace with actual struct
}
```

After generation, replace `Encode(nil)` with `Encode(&sep.ActualType{})` and implement the repository logic. See [docs/development/handlers.md](../development/handlers.md) for WADL-compliant response patterns.

### Repository stub

```go
type Repository struct {
    sync.RWMutex
    tag_lookup      map[string]Entity
    active_entities []bool
    pool            Pool
}

type Pool struct {
    // TODO: add slices for each resource type
    // items []sep.Item
}
```

### Generated main.go

```go
func main() {
    cfg, err := tlsutil.NewServerConfig("./ssl")
    if err != nil {
        log.Fatal(err)
    }
    server := http.Server{
        Addr:      routes.MyService,
        TLSConfig: cfg,
    }

    repo := memory.NewRepository(MAX_ENTITIES)
    repo.InitRepository("./ssl")

    h := handler.NewHandler(repo)
    http.Handle("GET /items", http.HandlerFunc(h.GetItems))
    http.Handle("POST /items", http.HandlerFunc(h.CreateItem))
    http.Handle("GET /items/{id}", http.HandlerFunc(h.GetItem))
    http.Handle("DELETE /items/{id}", http.HandlerFunc(h.DeleteItem))

    err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Architecture Principles

- **Isolated services** — each service has its own packages with no cross-service imports
- **Mutual TLS** — all services require valid client certificates (`RequireAndVerifyClientCert`)
- **Entity-based access** — client cert LFDI maps to an entity in the repository pool
- **Thread-safe storage** — `sync.RWMutex` for concurrent read / serialised write
- **XML serialisation** — all responses use `sep.ContentType` (`application/sep+xml`)

## Troubleshooting

**"declared and not used" errors:**
The generated stubs have `TODO` comments where you add implementation. Replace `nil` responses with actual repository calls.

**Server fails to start:**
- Verify SSL certificates exist in `./ssl`
- Ensure the port is not already in use
- Confirm the routes constant was added to `internal/routes/routes.go`

**TLS certificate errors:**
- Client certificates must be in `./ssl` with `client` in the filename
- Certificates must be valid and readable by the service
