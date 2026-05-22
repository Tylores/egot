# ESI Lifecycle Sequence Diagrams: Regulation

This document covers the client/server interactions and sequence diagrams for the **Regulation** grid service (fast active power tracking / 5-minute Energy Imbalance Market regulation) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client registers as an End Device, registers its dynamic active power/regulation capabilities, and completes mTLS authentication.

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

    Client->>GW: POST /edev/123/der (Register DER capabilities)
    GW->>DER: Forward POST /edev/123/der
    DER-->>GW: 201 Created (Location: /edev/123/der/1)
    GW-->>Client: 201 Created (Location: /edev/123/der/1)
```

---

## 2. Scheduling

In the Scheduling phase, the client requests participation in fast regulation by submitting a short-horizon (e.g. 5-minute / 300s) FlowReservationRequest, and retrieves the approved FlowReservationResponse.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant FR as "FlowReservation Service (:8027)"

    Client->>GW: POST /edev/123/frq (Create FlowReservationRequest)
    GW->>FR: Forward POST /edev/123/frq (Contains 5-minute duration request)
    Note over FR: FR checks real-time EIM window limits
    FR-->>GW: 201 Created (Location: /edev/123/frp/15)
    GW-->>Client: 201 Created (Location: /edev/123/frp/15)

    Client->>GW: GET /edev/123/frp/15 (Fetch approved FlowReservationResponse)
    GW->>FR: Forward GET /edev/123/frp/15
    FR-->>GW: 200 OK (FlowReservationResponse with status code [Accepted])
    GW-->>Client: 200 OK (FlowReservationResponse)
```

---

## 3. Operation

During the Operation phase, the client monitors real-time active controls and modulates its active power generation or load consumption rapidly in response to dynamic dispatch signals.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /edev/123/der/1/cdc (Fetch Current DER Controls)
    GW->>DER: Forward GET /edev/123/der/1/cdc
    DER-->>GW: 200 OK (CurrentDERControls with active regulation signal)
    GW-->>Client: 200 OK (CurrentDERControls)

    Note over Client: Client parses dynamic active power setpoint
    Note over Client: Client adjusts power injection/consumption (kW)
```

---

## 4. Verify

In the Verification phase, the client submits high-frequency operational status updates (DERStatus) to verify exact active power tracking performance.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: POST /edev/123/der/1/ders (Update DER Status)
    GW->>DER: Forward POST /edev/123/der/1/ders (Sends current power status)
    DER-->>GW: 201 Created (Stores status telemetry)
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, fast telemetry data logged via the Mirror Usage Point is evaluated by the Billing service to compute performance-based regulation payments.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Client->>GW: POST /mup/123 (Submit fast telemetry logs)
    GW->>MUP: Forward POST /mup/123
    MUP-->>GW: 201 Created (Stores reading logs)
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch regulation telemetry readings)
    MUP-->>Bill: 200 OK (Telemetry list)
    Bill->>Bill: Reconcile regulation tracking performance (accuracy and response time)
    Bill->>Bill: Apply performance-based regulation tariff & credit CustomerAccount
```
