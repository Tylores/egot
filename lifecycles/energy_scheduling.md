# ESI Lifecycle Sequence Diagrams: Energy Scheduling

This document covers the client/server interactions and sequence diagrams for the **Energy Scheduling** grid service (longer-term / Day-Ahead hourly power injection schedules) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client registers as an End Device, registers its DER capabilities, and initializes registration settings.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant EDev as "EDevice Service (:8015)"
    participant DER as "DER Service (:8026)"

    Note over Client, GW: mTLS Handshake initiated by Client
    Client->>GW: POST /edev (SFDI / registration payload)
    GW->>EDev: Forward POST /edev
    EDev-->>GW: 201 Created (Location: /edev/123)
    GW-->>Client: 201 Created (Location: /edev/123)

    Client->>GW: POST /edev/123/rg (Pin/SFDI verification)
    GW->>EDev: Forward POST /edev/123/rg
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: POST /der (Register DER capabilities)
    GW->>DER: Forward POST /der
    DER-->>GW: 201 Created (Location: /der/1)
    GW-->>Client: 201 Created (Location: /der/1)
```

---

## 2. Scheduling

In the Scheduling phase, the client requests a multi-hour or Day-Ahead energy schedule by posting a FlowReservationRequest, and fetches the approved FlowReservationResponse indicating the hourly power limits allocated by the server.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant FR as "FlowReservation Service (:8027)"

    Client->>GW: POST /frq (Create FlowReservationRequest)
    GW->>FR: Forward POST /frq (Contains hourly profile / durations)
    Note over FR: FR checks day-ahead window limits & feeder capacities
    FR-->>GW: 201 Created (Location: /frp/10)
    GW-->>Client: 201 Created (Location: /frp/10)

    Client->>GW: GET /frp/10 (Fetch approved FlowReservationResponse)
    GW->>FR: Forward GET /frp/10
    FR-->>GW: 200 OK (FlowReservationResponse with status code [Accepted])
    GW-->>Client: 200 OK (FlowReservationResponse)
```

---

## 3. Operation

During the Operation phase, the client tracks the current time and executes active power injections or draws according to the scheduled active hour blocks.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"

    Note over Client: Client monitors time-of-day clock
    Note over Client: Client starts scheduled active power generation (kW)
    Note over Client: Client adjusts output matching scheduled profile
```

---

## 4. Verify

In the Verification phase, the client posts energy meter readings to its Usage Point (UPT) service to record cumulative Wh energy delivery.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant UPT as "UPT Service (:8024)"

    Client->>GW: POST /upt/123/mr/1/rs/1/r (Post MeterReading telemetry)
    GW->>UPT: Forward POST /upt/123/mr/1/rs/1/r (Sends Reading payload)
    UPT-->>GW: 201 Created (Stores reading)
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, the Billing service processes the registered customer agreement, retrieves actual metered energy data, and reconciles output vs schedule to credit the customer.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant UPT as "UPT Service (:8024)"
    participant Bill as "Bill Service (:8011)"

    Note over Bill: Periodic billing process runs
    Bill->>UPT: GET /upt/123/mr/1/rs/1/r (Fetch energy delivery readings)
    UPT-->>Bill: 200 OK (Readings list)
    Bill->>Bill: Reconcile actual delivered energy against scheduled profile
    Bill->>Bill: Apply tariff settings & credit CustomerAccount
```
