# ESI Lifecycle Sequence Diagrams: Reserve Capacity

This document covers the client/server interactions and sequence diagrams for the **Reserve Capacity** grid service (capacity reservation / spinning and non-spinning reserves) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client registers as an End Device, registers its DER settings to advertise its reserve capacity limits (e.g. maximum discharge power), and registers with the Flow Reservation system.

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

In the Scheduling phase, the client requests a reserve capacity reservation by submitting a FlowReservationRequest, and retrieves the approved FlowReservationResponse indicating whether the reservation was accepted, modified, or rejected.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant FR as "FlowReservation Service (:8027)"

    Client->>GW: POST /frq (Create FlowReservationRequest)
    GW->>FR: Forward POST /frq (Contains requested reserve power & interval)
    Note over FR: FR checks transformer limits & active DR events
    FR-->>GW: 201 Created (Location: /frp/5)
    GW-->>Client: 201 Created (Location: /frp/5)

    Client->>GW: GET /frp/5 (Fetch approved FlowReservationResponse)
    GW->>FR: Forward GET /frp/5
    FR-->>GW: 200 OK (FlowReservationResponse with status code [Accepted])
    GW-->>Client: 200 OK (FlowReservationResponse)
```

---

## 3. Operation

During the Operation phase, the client transitions to a standby state, holding the requested power capacity ready for grid injection and ensuring it does not consume/discharge beyond the reserved baseline.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    
    Note over Client: Client enters Standby Operational mode
    Note over Client: Client holds battery capacity/power headroom online
    Note over Client: Client limits auxiliary loads to maintain reservation power
```

---

## 4. Verify

In the Verification phase, the client posts regular DERAvailability updates to prove to the utility that the reserved capacity remains fully online and available.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: POST /der/1/dera (Update DER Availability)
    GW->>DER: Forward POST /der/1/dera (Contains active available capacity)
    DER-->>GW: 201 Created (Stores availability telemetry)
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, availability readings are collected, and the Billing service issues standby credits to the customer account based on the reservation duration and tariff profile.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Client->>GW: POST /mup/123 (Submit availability telemetry)
    GW->>MUP: Forward POST /mup/123
    MUP-->>GW: 201 Created (Stores telemetry)
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch capacity availability readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list)
    Bill->>Bill: Reconcile actual availability against reservation duration
    Bill->>Bill: Apply reservation tariff & credit CustomerAccount
```
