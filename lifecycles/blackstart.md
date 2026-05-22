# ESI Lifecycle Sequence Diagrams: Blackstart & DR Dispatch

This document covers the client/server interactions and sequence diagrams for the **Blackstart & DR Dispatch** grid service (cold-start recovery / feeder load-shed dispatcher) across all 5 ESI lifecycle phases, aligned with IEEE 2030.5 CSIP test procedures.

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
    participant DR as "DR Service (:8014)"

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

    Client->>GW: POST /edev/123/lsl (Register sheddable capacity baseline)
    GW->>EDev: Forward POST /edev/123/lsl
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: GET /dr (Discover available DR Programs)
    GW->>DR: Forward GET /dr
    DR-->>GW: 200 OK (DemandResponseProgramList with pollRate=900)
    GW-->>Client: 200 OK (DemandResponseProgramList)
```

---

## 2. Scheduling

In the Scheduling phase, the client monitors and parses future planned load-shed events using either Polling or Subscriptions.

### Option A: Polling Interaction
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DR as "DR Service (:8014)"

    Note over Client: [CORE-003] Client polls at interval defined by pollRate (e.g. 900s)
    Client->>GW: GET /dr/1 (Fetch target DemandResponseProgram details)
    GW->>DR: Forward GET /dr/1
    DR-->>GW: 200 OK (DemandResponseProgram detailing event windows)
    GW-->>Client: 200 OK (DemandResponseProgram)

    Client->>GW: GET /dr/1/edc (Discover scheduled EndDeviceControl list)
    GW->>DR: Forward GET /dr/1/edc
    DR-->>GW: 200 OK (EndDeviceControlList showing future events)
    GW-->>Client: 200 OK (EndDeviceControlList)
```

### Option B: Subscription/Notification
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DR as "DR Service (:8014)"

    Client->>GW: POST /dr/1/edc/sub (Subscribe to EndDeviceControlList changes)
    GW->>DR: Forward POST /dr/1/edc/sub (Includes Client Notification URI)
    DR-->>GW: 201 Created (Location: /dr/1/edc/sub/5)
    GW-->>Client: 201 Created (Location: /dr/1/edc/sub/5)

    Note over DR: [CORE-018] / [CORE-019] Event is scheduled by Grid Operator
    DR->>GW: POST /client/notification (Notify client of resource change)
    GW->>Client: Forward POST /client/notification (Contains updated EndDeviceControlList)
    Client-->>GW: 204 No Content
    GW-->>DR: 204 No Content
```

---

## 3. Operation

During the Operation phase, the client tracks active controls, applies event randomization, manages priorities, and triggers shedding.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DR as "DR Service (:8014)"

    Client->>GW: GET /dr/1/actedc (Fetch Active End Device Controls)
    GW->>DR: Forward GET /dr/1/actedc
    DR-->>GW: 200 OK (ActiveEndDeviceControl detailing immediate load shed event)
    GW-->>Client: 200 OK (ActiveEndDeviceControl)

    Note over Client: [CORE-021] Client applies randomizeStart / randomizeDuration values
    Note over Client: [BASIC-021] Priority Check: If overlapping events occur, SP (Service Point) level supersedes SY (System) level
    Note over Client: Client executes emergency load-shed trigger at randomized start time
```

---

## 4. Verify

In the Verification phase, the client posts status confirmations, reports execution states, updates statuses, and logs alarms.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant RSPS as "Rsps Service (:8041)"
    participant EDev as "EDevice Service (:8015)"

    Note over Client: [CORE-022] / [BASIC-017] Response State Reporting
    Client->>GW: POST /rsps/123/rsp (Post Status = 1 [Received])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: Event starts executing
    Client->>GW: POST /rsps/123/rsp (Post Status = 2 [Started])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: [BASIC-028] Inverter Status Update
    Client->>GW: PUT /edev/123/dstat (Periodically report device status)
    GW->>EDev: Forward PUT /edev/123/dstat
    EDev-->>GW: 204 No Content
    GW-->>Client: 204 No Content

    Note over Client: [BASIC-027] Alarm Logging (If fault occurs during event)
    Client->>GW: POST /edev/123/lel (POST LogEvent for fault alarm)
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

In the Settlement phase, telemetry data is evaluated by the Billing service to reconcile consumption against sheddable capacity agreements.

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
    
    Bill->>Bill: Reconcile baseline sheddable energy vs actual consumption during event window
    Bill->>Bill: Fetch customer agreement active billing periods (/bill/123/ca/1/actbp)
    Bill->>Bill: Reconcile & credit CustomerAccount
```
