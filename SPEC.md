# egot — System Specification

> IEEE 2030.5 / SEP 2 Server & Device Emulation Platform

---

## Table of Contents

1. [Overview](#1-overview)
2. [Architecture](#2-architecture)
3. [Service Inventory](#3-service-inventory)
4. [mTLS & Identity](#4-mtls--identity)
5. [Cross-Service Discovery & Linking](#5-cross-service-discovery--linking)
6. [Repository Pattern](#6-repository-pattern)
7. [Business Logic Implementation Plan](#7-business-logic-implementation-plan)
8. [Persistence](#8-persistence)
9. [Device Emulator Clients](#9-device-emulator-clients)
10. [Testing Strategy](#10-testing-strategy)
11. [Development Workflow](#11-development-workflow)
12. [Open Questions & Tradeoffs](#12-open-questions--tradeoffs)

---

## 1. Overview

**egot** is a Go implementation of the IEEE 2030.5 Smart Energy Profile 2 (SEP 2) server
and device client platform. It hosts a fleet of microservices — each responsible for a
distinct SEP 2 function set — and provides device emulators that interact with the
platform as real-world energy devices (smart meters, distributed energy resources, etc.).

All communication is over HTTPS with mutual TLS (mTLS). Device identity is derived
from client certificate fingerprints (LFDI/SFDI).

**Current state:** All services are scaffolded and build successfully. Core is the only
fully implemented service. All others have correct HTTP routing and status-code semantics
but return empty SEP 2 objects — business logic (repository data) is not yet wired up.

---

## 2. Architecture

### 2.1 Microservice Layout

```
egot.internal.com
├── :8000  core            — Device discovery (DCAP), End Device registration, Time
├── :8001  flowreservation — Flow Reservation requests/responses
├── :8003  rsps            — Response sets (client → server submissions)
├── :8010  BRS             — Billing Reading Sets
├── :8011  Bill            — Customer Accounts, Billing Periods, Agreements
├── :8012  DCAP            — (reserved; core handles /dcap)
├── :8013  DERP            — DER Programs, Controls, Curves
├── :8014  DR              — Demand Response Programs & Controls
├── :8015  EDevice         — End Device management (full edev.wadl)
├── :8016  File            — Firmware & config file delivery
├── :8017  MUP             — Mirror Usage Points (device-push metering)
├── :8018  Messaging       — Text messaging programs
├── :8019  Notify          — Server-push notifications
├── :8020  PPY             — Prepayment
├── :8021  SDevice         — Self Device info
├── :8022  TariffProfile   — Tariff Profiles, Rate Components
├── :8023  TimeOfUse       — Time-of-Use intervals
├── :8024  UPT             — Usage Points, Meter Readings
└── :8025  DER             — Per-device DER resources
```

Each service follows the same three-layer structure:

```
cmd/<Service>/main.go              — HTTP server, route registration
internal/<Service>/handler/        — HTTP handlers (SEP 2 request/response logic)
internal/<Service>/repository/memory/ — In-memory data store
```

### 2.2 Request Lifecycle

```
Device (mTLS client cert)
   │
   ▼
HTTP/2 over TLS (:port)
   │
   ├─ TLS layer: RequireAndVerifyClientCert (ca.crt)
   │
   ▼
Handler.getLFDI(req)
   ├─ Extract cert from req.TLS.PeerCertificates[0]
   ├─ LFDI = SHA256(cert.Raw)[0:40]    (40 hex chars = 160 bits)
   └─ repo.GetEntity(lfdi)             → 404 if not registered
   │
   ▼
Business logic (read/write repository)
   │
   ▼
XML encode response (sep.ContentType = "application/sep+xml")
```

### 2.3 Inter-Service Communication

Services are independent HTTP servers. They do **not** call each other at runtime.
Cross-resource links (e.g., an EndDevice linking to its BillingReadingSets) are
expressed as `href` strings embedded in SEP 2 response objects. Clients are
responsible for following links to the appropriate service.

`DeviceCapability` (served by Core on `/dcap`) acts as the root service registry —
it contains links to all enabled function-set top-level resources. Clients start
here and traverse the link tree.

### 2.4 Technology Stack

| Concern | Choice |
|---------|--------|
| Language | Go 1.23 |
| HTTP server | `net/http` stdlib (HTTP/1.1 + HTTPS) |
| TLS | `crypto/tls` — mTLS, TLS 1.2+ |
| Serialization | `encoding/xml` (SEP 2 schema-generated types) |
| Persistence | `internal/store` (gob-encoded key/value, atomic file writes) |
| SEP 2 types | Auto-generated from `sep.xsd` (`sep/sep.go`) |
| Testing | `testing` + `github.com/stretchr/testify` |
| XSD validation | `github.com/terminalstatic/go-xsd-validate` |
| Certificates | OpenSSL EC P-256 (via `scripts/gen-ssl.sh`) |

---

## 3. Service Inventory

### 3.1 Core Services

#### `core` (`:8000`) — Device Capability & Registration
The entry point for all SEP 2 clients. Fully implemented.

| Route | Handler | Status |
|-------|---------|--------|
| `GET /dcap` | GetDeviceCapability | ✅ Implemented — returns links to all function sets |
| `GET /tm` | GetTime | ✅ Implemented — returns synchronized time |
| `GET /edev` | GetEndDevices | ✅ Implemented — returns list for this client |
| `GET /edev/{id}` | GetEndDevice | ✅ Implemented |
| `GET /rg/{id}` | GetRegistration | ✅ Implemented |

Repository: Pre-allocated pool of `EndDevice` + `Registration` structs, initialized
from client certificates at startup via `InitRepository("./ssl")`.

#### `rsps` (`:8003`) — Response Sets
Collects device responses to server-sent events (DR activations, DER controls).

#### `flowreservation` (`:8001`) — Flow Reservation
Handles flow reservation requests/responses for network bandwidth management.

### 3.2 Metering & Billing Services

#### `UPT` (`:8024`) — Usage Points & Meter Readings
Manages `UsagePoint`, `MeterReading`, `ReadingType`, `ReadingSet`, `Reading`.
Primary data collection endpoint for smart meters.

#### `BRS` (`:8010`) — Billing Reading Sets
Presents formatted billing data: `BillingReadingSet`, `BillingReadingList`, `BillingReading`.

#### `Bill` (`:8011`) — Customer Accounts
Customer billing hierarchy: `CustomerAccount` → `CustomerAgreement` →
`BillingPeriod` → `ProjectionReading`.

#### `MUP` (`:8017`) — Mirror Usage Points
Device-push metering: devices POST readings to `MirrorUsagePoint` endpoints
rather than being polled. Primary metering path for resource-constrained devices.

### 3.3 Energy Management Services

#### `DERP` (`:8013`) — DER Programs
Manages DER control programs: `DERProgram`, `DERControl`, `DERCurve`,
`DefaultDERControl`, `ActiveDERControlList`. The utility operator creates
programs here; devices pull them.

#### `DER` (`:8025`) — Per-Device DER Resources
Per-EndDevice DER state: `DERList`, `DERSettings`, `DERStatus`, `DERAvailability`,
`DERCapability`, `DERComponent`. Devices report their DER capabilities and status here.

#### `DR` (`:8014`) — Demand Response
Demand response programs: `DemandResponseProgram`, `EndDeviceControl`,
`ActiveEndDeviceControlList`, `LoadShedAvailability`.

### 3.4 Pricing & Tariff Services

#### `TariffProfile` (`:8022`) — Tariff Profiles
`TariffProfile` → `RateComponent` → `ConsumptionTariffInterval` hierarchies.
Utility operators publish rate schedules here.

#### `TimeOfUse` (`:8023`) — Time-of-Use
`TimeTariffInterval` with time-bounded rate specifications.

#### `PPY` (`:8020`) — Prepayment
Prepayment accounts, credit registers, supply interruption overrides.

### 3.5 Device Management Services

#### `EDevice` (`:8015`) — End Device Management
Full end-device lifecycle from the edev.wadl (the largest WADL at ~1935 lines):
device configuration, status, capabilities, FSA, DER links, network interfaces.

#### `SDevice` (`:8021`) — Self Device
Server-side self-device resource (`/sdev`). Describes the utility server itself.

#### `File` (`:8016`) — File Delivery
Firmware updates and configuration file delivery to devices.

### 3.6 Communication Services

#### `Messaging` (`:8018`) — Text Messages
Utility-to-device text messaging: `MessagingProgram` → `TextMessage`.

#### `Notify` (`:8019`) — Notifications
Server-push notification delivery. Devices register `Subscription` resources
pointing to their callback URIs; this service delivers `Notification` objects.

---

## 4. mTLS & Identity

### 4.1 Certificate Hierarchy

```
ssl/ca.crt       — Self-signed root CA (10-year validity, EC P-256)
ssl/server.crt   — Server certificate (CN=egot.internal.com, 1-year)
ssl/client.crt   — Default development client certificate
ssl/client-NNNN.crt — Numbered device client certificates
```

All services share the same CA, server cert, and client certs from `./ssl/`.
Services must be started from the repository root directory.

### 4.2 Device Identity

```
LFDI = hex(SHA256(client_cert.Raw))[0:40]   // 40 hex chars, 160 bits
SFDI = uint64(parseInt(LFDI[0:9], 16))      // 36-bit truncated identifier
```

LFDI is the primary device key. Every handler validates the requesting device's
LFDI against the repository before serving data. Unregistered devices get HTTP 404.

### 4.3 Registration Flow

1. Client certificate is provisioned (via `make ssl-clients N=<count>`)
2. Server starts → `InitRepository("./ssl")` walks `ssl/`, finds all `client*.crt`
3. For each cert: extracts LFDI, allocates an `Entity` slot, populates initial data
4. Device connects with mTLS → LFDI extracted from cert → mapped to Entity slot

New devices require a server restart to be registered in the current implementation.
See §12 for a discussion of dynamic registration.

---

## 5. Cross-Service Discovery & Linking

### 5.1 Link Structure

SEP 2 resources carry typed link structs that embed `href` strings:

```go
type EndDevice struct {
    *ExternalDevice
    // href="/brs" — points to BRS service on :8010
    BillingReadingSetListLink *BillingReadingSetListLink
    // href="/derp" — points to DERP service on :8013
    DERProgramListLink *DERProgramListLink
    ...
}
```

Links use **relative paths only** (e.g., `/brs`, `/brs/{id}`). This means a client
must already know the host:port of each service before following a link.

### 5.2 Service Registry via DeviceCapability

The `DeviceCapability` response from Core is the root discovery document.
Business logic implementation must populate all active service links:

```go
// Core repository — GetDeviceCapability()
func (r *Repository) GetDeviceCapability() sep.DeviceCapability {
    dcap := sep.NewDeviceCapability(...)
    dcap.EndDeviceListLink      = link(uri.EndDeviceList)       // core :8000
    dcap.TimeLink               = link(uri.Time)                // core :8000
    // Must also include links to other services:
    dcap.BillingReadingSetListLink    = link(uri.BillingReadingSetList) // :8010
    dcap.TariffProfileListLink        = link(uri.TariffProfileList)     // :8022
    dcap.DemandResponseProgramListLink = link(uri.DemandResponseProgramList) // :8014
    dcap.MessagingProgramListLink      = link(uri.MessagingProgramList)  // :8018
    dcap.UsagePointListLink            = link(uri.UsagePointList)        // :8024
    dcap.DERProgramListLink            = link(uri.DERProgramList)        // :8013
    dcap.PrepaymentListLink            = link(uri.PrepaymentList)        // :8020
    dcap.FileListLink                  = link(uri.FileList)              // :8016
    dcap.ResponseSetListLink           = link(uri.ResponseSetList)       // :8003
    return *dcap
}
```

Clients connect to Core first, parse the `DeviceCapability` hrefs, and connect
directly to the relevant service ports for subsequent requests.

### 5.3 Client Port Discovery

Device clients must resolve service ports. Two approaches:

**Option A — Static config (current):** Clients import `internal/routes` constants.
Simple, zero overhead, works for local development and testing.

**Option B — DNS-SD:** Each service registers itself with a local mDNS/DNS-SD
responder. Clients discover services dynamically. Appropriate for production
deployments where service ports may vary.

**Decision:** Implement static config (Option A) now. The `routes` package is the
single source of truth for port assignments. Add DNS-SD in a future iteration.

---

## 6. Repository Pattern

### 6.1 Two Patterns in Use

**Pattern A — Entity Pool** (core, BRS): Pre-allocated typed slices indexed by
`Entity` (uint32). Fast O(1) lookups. Used when the service manages complex
per-device state.

```go
type Entity uint32

type Pool struct {
    edev []sep.EndDevice          // pool[entity] = device data
    reg  []sep.Registration
}

type Repository struct {
    sync.RWMutex
    tag_lookup      map[string]Entity   // LFDI → Entity index
    active_entities []bool              // allocation bitmap
    pool            Pool
}
```

**Pattern B — Simple Map** (Bill, DCAP, DERP, DR, and most others): Generic
`map[string]interface{}` keyed on LFDI. Currently returned by scaffolded services.
Authentication works; business data is not yet stored.

```go
type Repository struct {
    entities map[string]interface{}   // LFDI → whatever
}
```

### 6.2 Target Pattern for Business Logic

All services implementing actual data should graduate to the Entity Pool pattern
with typed resource collections. Each service manages its own resource hierarchy:

```go
// Example: UPT service repository
type Pool struct {
    usagePoints     []sep.UsagePoint
    meterReadings   [][]sep.MeterReading      // [entity][readingIdx]
    readingSets     [][][]sep.ReadingSet       // [entity][mrIdx][rsIdx]
    readings        [][][][]sep.Reading        // [entity][mrIdx][rsIdx][rIdx]
}
```

Resources that are per-device are indexed by Entity. Resources that are global
(e.g., TariffProfile, DERProgram) are indexed by their resource ID, not by LFDI.

### 6.3 Global vs Per-Device Resources

| Service | Resource Scope | Key |
|---------|---------------|-----|
| Core — EndDevice | per-device | Entity (LFDI-derived) |
| UPT — UsagePoint | per-device | Entity |
| UPT — MeterReading | per-device | Entity + reading index |
| BRS — BillingReadingSet | per-device | Entity + set index |
| DERP — DERProgram | global (all devices) | program ID |
| DERP — DERControl | per-program | program ID + control index |
| DR — DemandResponseProgram | global | program ID |
| TariffProfile — TariffProfile | global | tariff ID |
| Bill — CustomerAccount | per-device | Entity |
| Messaging — MessagingProgram | global | program ID |

---

## 7. Business Logic Implementation Plan

### 7.1 Priority Order

Implementation should proceed in dependency order:

**Phase 1 — Foundation** (Core already done)
- [x] Core: DeviceCapability, EndDevice, Registration, Time
- [ ] Core: populate all function-set links in DeviceCapability response
- [ ] Core: implement FunctionSetAssignments (`/fsa`) — tells each device which services are available to it

**Phase 2 — Metering** (primary data collection)
- [ ] **MUP**: Accept device-pushed readings; store per-device `MirrorUsagePoint` collections
- [ ] **UPT**: Serve `UsagePoint` → `MeterReading` → `ReadingSet` → `Reading` hierarchy

**Phase 3 — Energy Programs** (utility → device commands)
- [ ] **TariffProfile**: CRUD for `TariffProfile` and `RateComponent`
- [ ] **TimeOfUse**: `TimeTariffInterval` schedule management
- [ ] **DERP**: `DERProgram` creation, `DERControl` scheduling, `DERCurve` library
- [ ] **DR**: `DemandResponseProgram` and `EndDeviceControl` event management

**Phase 4 — Device State** (device → server telemetry)
- [ ] **DER**: Per-device `DERSettings`, `DERStatus`, `DERAvailability`, `DERCapability`
- [ ] **EDevice**: Extended end-device attributes, configuration, FSA links

**Phase 5 — Billing & Prepay**
- [ ] **BRS**: `BillingReadingSet` hierarchies linked from UsagePoint
- [ ] **Bill**: `CustomerAccount` → `CustomerAgreement` → `BillingPeriod`
- [ ] **PPY**: Prepayment account balance, credit registers

**Phase 6 — Communications**
- [ ] **Messaging**: `MessagingProgram` → `TextMessage` delivery
- [ ] **Notify**: Subscription management + notification delivery
- [ ] **rsps**: Collect and validate device responses to events

### 7.2 Handler Implementation Rules

Every non-stub handler must follow this pattern:

```go
func (h *Handler) GETResourceList(w http.ResponseWriter, req *http.Request) {
    lfdi, err := h.getLFDI(req)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // 1. Parse path params (if any)
    id1, err := strconv.Atoi(req.PathValue("id1"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // 2. Parse query params (s=start, l=limit, a=after)
    query := req.URL.Query()
    start, _ := strconv.Atoi(query.Get("s"))
    limit, _ := strconv.Atoi(query.Get("l"))

    // 3. Fetch from repository
    list, err := h.repo.GetResourceList(lfdi, id1, start, limit)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // 4. Respond
    w.Header().Set("Content-Type", sep.ContentType)
    w.WriteHeader(http.StatusOK)
    xml.NewEncoder(w).Encode(list)
}
```

**POST handlers** must:
- Decode the request body: `xml.NewDecoder(req.Body).Decode(&resource)`
- Validate required fields
- Allocate a resource ID
- Set `Location` header to the new resource URI
- Return 201

**PUT handlers** must:
- Decode the request body
- Validate the path ID matches an existing resource
- Update repository
- Return 204 (no content) or 200

### 7.3 Query Parameter Support (Pagination)

All List resources support:
- `s` — start index (default 0)
- `l` — limit / results per page (default `MAX_ENTITIES`)
- `a` — after timestamp filter (for event lists)

Repository `GetXxxList(lfdi, start, limit int)` methods must respect these.
The response `List.AllAttr` = total count, `List.ResultsAttr` = count in this page.

### 7.4 PollRate

`DeviceCapability.PollRateAttr` and per-resource `PollRateAttr` fields tell clients
how often to poll. Default is `sep.PollRate = 900` seconds (15 min). Set lower
values (e.g., 60s) on resources that change frequently (DERControl active list,
DR active controls).

---

## 8. Persistence

### 8.1 Current State

All repository data is in-memory and lost on restart. The only persistent state
is derived from the `./ssl/*.crt` files, which are re-read on each startup.

### 8.2 `internal/store` Package

The `internal/store` package provides a thread-safe gob-encoded key/value store
with atomic file writes:

```go
s := store.New("./data/derp.gob")
s.Load()                            // load persisted data on startup
s.Set("program:1", myDERProgram)    // write to in-memory map
s.Save()                            // persist atomically to disk
```

Callers must register concrete types before Load/Save:
```go
gob.Register(sep.DERProgram{})
gob.Register(sep.DERControl{})
```

### 8.3 Persistence Plan

| Service | What to persist | File |
|---------|-----------------|------|
| Core | EndDevice, Registration | `./data/core.gob` |
| DERP | DERProgram, DERControl, DERCurve | `./data/derp.gob` |
| DR | DemandResponseProgram, EndDeviceControl | `./data/dr.gob` |
| TariffProfile | TariffProfile, RateComponent | `./data/tp.gob` |
| UPT | UsagePoint, MeterReading, ReadingSet | `./data/upt.gob` |
| Bill | CustomerAccount, CustomerAgreement | `./data/bill.gob` |
| PPY | Prepayment, AccountBalance | `./data/ppy.gob` |

**Decision:** Ephemeral services (Notify, Messaging, File) do not need persistence —
content is operator-supplied at startup.

### 8.4 Write-Through Pattern

Repositories that use persistence should write-through on every mutation:

```go
func (r *Repository) PutDERProgram(id int, prog sep.DERProgram) error {
    r.Lock()
    defer r.Unlock()
    r.pool.programs[id] = prog
    return r.store.Save()   // atomic write after every mutation
}
```

On startup: `store.Load()` → populate pool from stored data.

---

## 9. Device Emulator Clients

### 9.1 Overview

Device emulators simulate real-world energy devices interacting with the egot
server stack. They implement the SEP 2 client protocol: discovery, registration,
polling, and event response.

Three device profiles are planned:

| Profile | Binary | Simulates |
|---------|--------|-----------|
| Smart Meter | `cmd/emulator-meter` | AMI meter pushing interval readings |
| DER Inverter | `cmd/emulator-der` | Solar inverter reporting DER status + executing DER controls |
| EV Charger | `cmd/emulator-evse` | EVSE responding to DR load-shed events |

### 9.2 Common Client Infrastructure

All emulators share a common client library in `internal/emulator/client/`:

```go
type Client struct {
    http    *http.Client          // mTLS-configured
    baseURL string                // https://egot.internal.com:<port>
    lfdi    string                // derived from cert at startup
}

func NewClient(sslDir, certName string, port int) (*Client, error)
func (c *Client) Get(path string, out any) error
func (c *Client) Post(path string, body any) (location string, err error)
func (c *Client) Put(path string, body any) error
func (c *Client) Delete(path string) error
func (c *Client) Head(path string) (int, error)
```

### 9.3 SEP 2 Client Protocol

#### Discovery Sequence

```
1. GET https://core:8000/dcap
   → DeviceCapability{EndDeviceListLink, BillingReadingSetListLink, ...}

2. GET /edev  → find own EndDevice (match by LFDI)

3. For each link in DeviceCapability → record service address
   (Links map to specific service ports via routes package)

4. GET /fsa/{id} → FunctionSetAssignments
   → Know which function sets the server expects this device to use
```

#### Polling Lifecycle

```go
// After discovery, poll each subscribed resource at its PollRate
func (d *Device) pollLoop(ctx context.Context, resource string, port int, interval time.Duration) {
    ticker := time.NewTicker(interval)
    for {
        select {
        case <-ticker.C:
            d.poll(resource, port)
        case <-ctx.Done():
            return
        }
    }
}
```

#### Event Response Flow (DER / DR)

```
1. Device polls /actderc (Active DER Control List) at reduced interval
2. Server returns new DERControl event
3. Device reads DERControl.DERControlBase (power setpoints, ramp rates)
4. Device applies control locally (simulated)
5. Device POSTs response to /rsps or /rsp
```

### 9.4 Smart Meter Emulator (`cmd/emulator-meter`)

**Behavior:**
- On startup: discover DCAP → find `/mup` on MUP service → POST `MirrorUsagePoint`
- Every 15 minutes: POST interval reading to `/mup/{id}` (MirrorUsagePoint push)
- Configurable via flags: `-ssl`, `-interval`, `-watt-hour-base`, `-variance`

**Data generation:**
```go
// Simulate realistic energy consumption
type MeterSimulator struct {
    baseLoadW    float64     // average watt draw
    variancePct  float64     // random variance ±%
    intervalSec  int         // reading interval (default 900)
    accumulated  float64     // Wh accumulator
}

func (m *MeterSimulator) NextReading() sep.Reading {
    // generate reading with realistic variance
    // increment accumulated kWh counter
}
```

**Configuration flags:**
```
-ssl          path to ssl directory (default: ./ssl)
-cert         client cert name (default: client)
-interval     reading interval seconds (default: 900)
-load         base load in watts (default: 1500)
-variance     variance percentage (default: 0.15)
```

### 9.5 DER Inverter Emulator (`cmd/emulator-der`)

**Behavior:**
- Startup: discover DCAP → find `/der` on DER service → PUT `DERSettings`, `DERCapability`
- Every 60s: PUT `DERStatus`, `DERAvailability` (generation/state telemetry)
- Every 300s: Poll `/actderc` on DERP service for active DER control events
- On new DERControl: apply setpoints (simulated), POST response to `/rsps`

**Data generation:**
```go
type InverterSimulator struct {
    ratedCapacityW  float64    // rated AC output capacity
    currentOutputW  float64    // current generation (varies with simulated irradiance)
    stateOfCharge   float64    // 0.0–1.0 if battery-attached
    irradiance      float64    // W/m² (time-of-day simulation)
}
```

**Simulated irradiance model:** A sine curve peaking at solar noon to simulate
realistic PV generation profiles. Includes cloud cover random variation.

**Configuration flags:**
```
-ssl          path to ssl directory (default: ./ssl)
-cert         client cert name (default: client)
-capacity     rated capacity in watts (default: 5000)
-battery      enable battery storage simulation (default: false)
-latitude     latitude for solar angle calculation (default: 37.0)
```

### 9.6 EV Charger Emulator (`cmd/emulator-evse`)

**Behavior:**
- Startup: discover DCAP → register capabilities
- Poll `/dr/{id}/actedc` (Active End Device Control) for load-shed events
- On DR event: reduce charging rate to specified `durationSec` and `randomizeStart`
- Report `LoadShedAvailability` to indicate how much load can be shed
- POST `DrResponse` to `/rsps` after event ends

**Configuration flags:**
```
-ssl          path to ssl directory (default: ./ssl)
-cert         client cert name (default: client)
-max-kw       maximum charging rate kW (default: 7.2)
-min-kw       minimum charging rate kW during DR (default: 0)
```

### 9.7 Emulator Launcher (`cmd/emulator-fleet`)

Launches multiple emulators to simulate a realistic device population:

```
emulator-fleet -ssl ./ssl -meters 10 -inverters 5 -evse 3
```

Generates one numbered client cert per device (or reuses existing), spawns each
emulator as a goroutine with its own `http.Client` and cert.

### 9.8 Crawler (existing — `cmd/crawler`)

The existing crawler validates that all reachable resources respond correctly and
their XML bodies validate against the SEP 2 XSD schema. It starts from `/dcap`
and follows all `href` links recursively, checking HEAD/GET/PUT/POST/DELETE.

---

## 10. Testing Strategy

### 10.1 Unit Tests

Each service handler should have a unit test that:
1. Creates a handler with a test repository (no TLS required — `getLFDI` returns test LFDI when `req.TLS == nil`)
2. Calls the handler directly with `httptest.NewRecorder()` and `httptest.NewRequest()`
3. Asserts status code, Content-Type header, and decoded XML body

```go
func TestGETDERProgramList(t *testing.T) {
    repo := memory.NewRepository(10)
    h := handler.NewHandler(repo)
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/derp", nil)
    h.GETDERProgramList(w, r)
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, sep.ContentType, w.Header().Get("Content-Type"))
}
```

### 10.2 Integration Tests (`test/integration/`)

Pattern: start the service binary → run a client against it → verify responses.
Existing: `test/integration/core/` and `test/integration/access/`.

Each service should add an integration test that:
1. Starts the service binary with `exec.Command`
2. Waits for it to be ready (TCP dial loop)
3. Runs a minimal `cmd/client`-style verification
4. Kills the process

### 10.3 XSD Validation

The `cmd/crawler` performs XSD validation of all GET responses. After implementing
business logic in any service, run the crawler to confirm XML conformance:

```
make build-all && make start && ./bin/crawler
```

### 10.4 End-to-End Simulation

Full pipeline test: start all services + an emulator fleet, let them run for one
poll cycle, verify no errors in service logs.

```
make start
./bin/emulator-fleet -ssl ./ssl -meters 2 -inverters 1 -evse 1
make stop
```

### 10.5 Route Coverage

`make test-routes` verifies that every handler in every service is registered on
the expected path and returns the correct status code for each HTTP method.

---

## 11. Development Workflow

### 11.1 Initial Setup

```bash
# 1. Generate all certificates
make ssl-refresh

# 2. (Optional) Generate additional client certs for multi-device testing
make ssl-clients N=10

# 3. Build everything
make build-all

# 4. Start all services
make start

# 5. Verify with crawler
./bin/crawler
```

### 11.2 Adding a New Service

```bash
# 1. Extract the WADL for this service, set port and max_entities
./bin/wadl-extract -wadl wadl/example.wadl -output wadl/example-extracted.wadl
# Edit the port in the output to match the next available port in routes.go

# 2. Generate the scaffold
./bin/scaffold-gen -wadl wadl/example-extracted.wadl

# 3. Implement business logic in internal/<service>/handler/handler.go
# 4. Implement repository in internal/<service>/repository/memory/memory.go
# 5. Add make build-<service> and SERVICES list entry to Makefile
# 6. Build and test
```

### 11.3 Certificate Management

```bash
make ssl-refresh          # Regenerate all certs (auto-restarts running services)
make ssl-clients N=5      # Add 5 more numbered client certs
```

When `ssl-refresh` runs while services are running, it automatically restarts them.

### 11.4 Logs

Service logs are written to `./bin/logs/<ServiceName>.log` when started via `make start`.

---

## 12. Open Questions & Tradeoffs

### 12.1 Dynamic Device Registration

**Current behavior:** Devices are registered from `./ssl/*.crt` files on startup.
Adding a new device requires generating a cert and restarting all services.

**Alternative:** Implement a registration endpoint where a pre-authorized device
can POST its certificate to `Core`, which adds it to the repository live.
SEP 2 supports this via the `Registration` resource and a PIN-based flow.

**Tradeoff:** Dynamic registration increases complexity and attack surface.
For a controlled lab/test environment, static registration is sufficient and simpler.
**Decision:** Implement static registration first; spec dynamic registration separately.

### 12.2 Service-to-Service Communication

Some SEP 2 behaviors require cross-service awareness. For example:
- When a DER program is activated, the DR service might want to notify the DERP service.
- When a device's `EndDeviceControl` fires, the `rsps` service collects its response.

**Option A — Event bus:** An internal pub/sub channel (Go channels) or a lightweight
message broker (NATS, Redis PubSub) that services subscribe to.

**Option B — Direct HTTP calls:** Services call each other's internal APIs over
localhost using the routes constants.

**Option C — Polling / eventual consistency:** Services don't communicate; the
client drives all state transitions by polling each service independently.

**Decision for now:** Option C (client-driven). It matches the SEP 2 specification's
design intent. Revisit if server-side coordination becomes necessary.

### 12.3 Subscription / Notification Delivery

SEP 2 defines a `Subscription` mechanism where devices register a callback URI
and the server POSTs `Notification` objects to it when resources change.

This requires the server to initiate HTTPS connections **to** devices — reversing
the usual direction. This is architecturally significant:

- Server needs to maintain a notification delivery queue
- Delivery must be retried on failure
- The device's callback URI must be resolvable from the server
- Device must present a server cert the egot CA can verify (or mutual trust)

**Decision:** Notifications are deferred to Phase 6. The `Notify` service scaffold
handles the Subscription CRUD. Actual delivery (POSTing to device callbacks) is
a follow-on implementation.

### 12.4 Concurrency Model for Emulators

The emulator fleet spawns many goroutines, each making periodic HTTPS requests.
With 100+ simulated devices each polling multiple services, connection pool
exhaustion becomes a concern.

**Decision:** Each emulator profile type shares a single `*http.Client` with
connection pooling via `http.Transport`. Individual device goroutines serialize
their own requests but share the transport. `MaxIdleConns` and dial timeouts
should be tuned based on fleet size.

### 12.5 Data Volume & Pool Sizing

Each service allocates `MAX_ENTITIES` slots at startup. At 100 entities × 15
services, the footprint is small. But nested resources (readings per meter, controls
per DER program) grow multiplicatively.

**Decision:** Use flat slices for per-device resources (`readings[entity][idx]`).
Cap list lengths at `MAX_ENTITIES` per level. For services with unbounded growth
(UPT readings over time), implement a ring buffer or time-windowed retention
(e.g., keep last 96 15-minute intervals = 24 hours).

### 12.6 DCAP vs Core Consolidation

The codebase has both `cmd/core` (the working DCAP + EndDevice implementation)
and `cmd/DCAP` (a scaffolded stub from the WADL). These overlap.

**Decision:** `cmd/core` is the authoritative DCAP implementation. `cmd/DCAP` is
vestigial and should either be removed or repurposed. The `routes.DCAP` constant
at `:8012` should not conflict with `routes.Core` at `:8000`.
