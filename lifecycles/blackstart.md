# ESI Lifecycle Sequence Diagrams: Blackstart & DR Dispatch

This document covers the client/server interactions and sequence diagrams for the **Blackstart & DR Dispatch** grid service (cold-start recovery / feeder load-shed dispatcher) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client registers as an End Device, registers its load-shed capacity telemetry (SheddablePower), and discovers available Demand Response (DR) programs.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant EDev as "EDevice Service (:8015)"
    participant DR as "DR Service (:8014)"

    Note over Client, GW: mTLS Handshake initiated by Client
    Client->>GW: POST /edev (SFDI / registration payload)
    GW->>EDev: Forward POST /edev
    EDev-->>GW: 201 Created (Location: /edev/123)
    GW-->>Client: 201 Created (Location: /edev/123)

    Client->>GW: POST /edev/123/rg (Pin/SFDI verification)
    GW->>EDev: Forward POST /edev/123/rg
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: POST /edev/123/lsl (Register LoadShedAvailability)
    GW->>EDev: Forward POST /edev/123/lsl (Sends SheddablePower baseline)
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: GET /dr (Discover available DR Programs)
    GW->>DR: Forward GET /dr
    DR-->>GW: 200 OK (DemandResponseProgramList)
    GW-->>Client: 200 OK (DemandResponseProgramList)
```

---

## 2. Scheduling

In the Scheduling phase, the client queries scheduled End Device Controls (EDC) associated with the Demand Response Program to identify future planned events.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DR as "DR Service (:8014)"

    Client->>GW: GET /dr/1 (Fetch target DemandResponseProgram details)
    GW->>DR: Forward GET /dr/1
    DR-->>GW: 200 OK (DemandResponseProgram detailing event windows)
    GW-->>Client: 200 OK (DemandResponseProgram)

    Client->>GW: GET /dr/1/edc (Discover scheduled EndDeviceControl list)
    GW->>DR: Forward GET /dr/1/edc
    DR-->>GW: 200 OK (EndDeviceControlList showing future events)
    GW-->>Client: 200 OK (EndDeviceControlList)
```

---

## 3. Operation

During the Operation phase (e.g. during a cold-start recovery or emergency feeder overload), the operator triggers a critical load-shed event. The client fetches active controls, parses the emergency signal, and immediately disconnects or limits its load.

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

    Note over Client: Client executes emergency load-shed trigger
    Note over Client: Client sheds requested sheddable load (kW)
```

---

## 4. Verify

In the Verification phase, the client posts confirmation responses to the Rsps service and updates its Device Status to prove command compliance.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant RSPS as "Rsps Service (:8041)"
    participant EDev as "EDevice Service (:8015)"

    Client->>GW: POST /rsps/123/rsp (Post control execution status)
    GW->>RSPS: Forward POST /rsps/123/rsp (Sends status = event started/completed)
    RSPS-->>GW: 201 Created (Stores Response payload)
    GW-->>Client: 201 Created

    Client->>GW: GET /edev/123/dstat (Fetch DeviceStatus for confirmation)
    GW->>EDev: Forward GET /edev/123/dstat
    EDev-->>GW: 200 OK (DeviceStatus reflecting reduced load level)
    GW-->>Client: 200 OK (DeviceStatus)
```

---

## 5. Settlement

In the Settlement phase, telemetry data from the Mirror Usage Point is evaluated by the Billing service against the customer agreement and the baselines to credit the account for emergency load shedding.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Client->>GW: POST /mup/123 (Submit telemetry meter readings)
    GW->>MUP: Forward POST /mup/123
    MUP-->>GW: 201 Created (Stores usage telemetry logs)
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch delivered energy readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list)
    Bill->>Bill: Reconcile shed energy baseline vs actual consumption during event
    Bill->>Bill: Apply blackstart/DR tariff & credit CustomerAccount
```
