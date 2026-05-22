# ESI Lifecycle Sequence Diagrams: Regulation

This document covers the client/server interactions and sequence diagrams for the **Regulation** grid service (fast active power tracking / 5-minute Energy Imbalance Market regulation) across all 5 ESI lifecycle phases, aligned with IEEE 2030.5 CSIP test procedures.

---

## 1. Registration

During Registration, the client performs discovery, synchronizes time, and registers its identity.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DCAP as "DCAP Service (:8012)"
    participant TimeOfUse as "TimeOfUse Service (:8023)"
    participant EDev as "EDevice Service (:8015)"
    participant DER as "DER Service (:8026)"

    Note over Client, GW: [CORE-001] mTLS Handshake initiated by Client
    
    Client->>GW: GET /dcap (Discover root capabilities)
    GW->>DCAP: Forward GET /dcap
    DCAP-->>GW: 200 OK (DeviceCapability List Links)
    GW-->>Client: 200 OK
    
    Client->>GW: GET /tm (Retrieve Server Time for sync)
    GW->>TimeOfUse: Forward GET /tm
    Note over TimeOfUse: [CORE-005] Confirm quality metric = 7 (uncoordinated)
    TimeOfUse-->>GW: 200 OK (Time resource with quality = 7)
    GW-->>Client: 200 OK (Client synchronizes local clock)

    Client->>GW: POST /edev (SFDI / registration payload)
    GW->>EDev: Forward POST /edev (with peer certificate metadata)
    EDev-->>GW: 201 Created (Location: /edev/123)
    GW-->>Client: 201 Created (Location: /edev/123)

    Client->>GW: POST /edev/123/rg (Registration PIN verification)
    GW->>EDev: Forward POST /edev/123/rg (PIN validation)
    Note over EDev: [CORE-003] / [BASIC-001] Verify LFDI matches and PIN is 111115
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: GET /der (Discover associated DER instances)
    GW->>DER: Forward GET /der
    Note over DER: [CORE-009] Retrieve DER capability registration links
    DER-->>GW: 200 OK (DERList containing DER MRID)
    GW-->>Client: 200 OK (DERList)
```

---

## 2. Scheduling

In the Scheduling phase, the client requests participation in fast regulation by submitting a short-horizon (e.g. 5-minute / 300s) FlowReservationRequest, and fetches or subscribes to the approved FlowReservationResponse.

### Option A: Polling Interaction
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant FR as "FlowReservation Service (:8027)"

    Client->>GW: POST /frq (Create FlowReservationRequest)
    GW->>FR: Forward POST /frq (Contains 5-minute duration request)
    Note over FR: [CORE-009] FR checks real-time EIM window limits
    FR-->>GW: 201 Created (Location: /frp/15)
    GW-->>Client: 201 Created (Location: /frp/15)

    Note over Client: [CORE-003] Client polls at interval defined by pollRate (e.g. 900s)
    Client->>GW: GET /frp/15 (Fetch approved FlowReservationResponse)
    GW->>FR: Forward GET /frp/15
    FR-->>GW: 200 OK (FlowReservationResponse with status code [Accepted])
    GW-->>Client: 200 OK (FlowReservationResponse)
```

### Option B: Subscription/Notification
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant FR as "FlowReservation Service (:8027)"

    Client->>GW: POST /frp/sub (Subscribe to FlowReservationResponse list changes)
    GW->>FR: Forward POST /frp/sub (Includes Client Notification URI)
    FR-->>GW: 201 Created (Location: /frp/sub/8)
    GW-->>Client: 201 Created (Location: /frp/sub/8)

    Note over FR: [CORE-018] / [CORE-019] Schedule changes or dynamic dispatch updates applied by Grid Operator
    FR->>GW: POST /client/notification (Notify client of response updates)
    GW->>Client: Forward POST /client/notification (Contains updated FlowReservationResponse)
    Client-->>GW: 204 No Content
    GW-->>FR: 204 No Content
```

---

## 3. Operation

During the Operation phase, the client monitors real-time active controls, applies event randomization, manages priorities, and modulates its active power generation or load consumption rapidly in response to dynamic dispatch signals.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /derp/1/actderc (Fetch Active DER Controls)
    GW->>DER: Forward GET /derp/1/actderc
    DER-->>GW: 200 OK (ActiveDERControl with active regulation signals)
    GW-->>Client: 200 OK (ActiveDERControl)

    Note over Client: [CORE-021] Client applies randomizeStart / randomizeDuration values
    Note over Client: [BASIC-021] Priority Check: SP (Service Point) level supersedes SY (System) level
    Note over Client: Client parses dynamic active power setpoint
    Note over Client: Client adjusts power injection/consumption (kW) dynamically
```

---

## 4. Verify

In the Verification phase, the client posts status confirmations, reports execution states, updates statuses, and logs alarms to verify exact active power tracking performance.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant RSPS as "Rsps Service (:8041)"
    participant DER as "DER Service (:8026)"
    participant EDev as "EDevice Service (:8015)"

    Note over Client: [CORE-022] / [BASIC-017] Response State Reporting
    Client->>GW: POST /rsps/123/rsp (Post Status = 1 [Received])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: Event starts executing (Dynamic tracking begins)
    Client->>GW: POST /rsps/123/rsp (Post Status = 2 [Started])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: [BASIC-028] Inverter Status Update
    Client->>GW: PUT /der/1/ders (Periodically report device/inverter status)
    GW->>DER: Forward PUT /der/1/ders
    DER-->>GW: 204 No Content
    GW-->>Client: 204 No Content

    Note over Client: [BASIC-027] Alarm Logging (If active power output tracking error exceeds limit)
    Client->>GW: POST /edev/123/lel (POST LogEvent for tracking deviation alert)
    GW->>EDev: Forward POST /edev/123/lel
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: Event completes execution
    Client->>GW: POST /rsps/123/rsp (Post Status = 3 [Completed])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, telemetry data is evaluated by the Billing service to credit the customer for active power delivery and performance-based tracking.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Note over Client: [BASIC-029] Client periodically posts cumulative Wh readings
    Client->>GW: POST /mup/123 (Submit telemetry meter readings payload)
    GW->>MUP: Forward POST /mup/123 (XML payload validation)
    MUP-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch delivered energy readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list)
    
    Bill->>Bill: Reconcile actual delivered active energy against scheduled regulation tracking signals
    Bill->>Bill: Fetch customer agreement active billing periods (/bill/123/ca/1/actbp)
    Bill->>Bill: Reconcile & credit CustomerAccount
```
