# ESI Lifecycle Sequence Diagrams: Frequency Response

This document covers the client/server interactions and sequence diagrams for the **Frequency Response** grid service (Frequency-Watt / active power adjustment based on local frequency deviations) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client (End Device/DER) registers itself with the server to establish an authenticated session, registers its DER capability, and discovers available DER programs that offer Frequency Response.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant EDev as "EDevice Service (:8015)"
    participant DER as "DER Service (:8026)"

    Note over Client, GW: mTLS Handshake initiated by Client
    Client->>GW: POST /edev (SFDI / registration payload)
    GW->>EDev: Forward POST /edev (with peer certificate metadata)
    EDev-->>GW: 201 Created (Location: /edev/123)
    GW-->>Client: 201 Created (Location: /edev/123)

    Client->>GW: POST /edev/123/rg (Pin/SFDI verification)
    GW->>EDev: Forward POST /edev/123/rg
    EDev-->>GW: 201 Created
    GW-->>Client: 201 Created

    Client->>GW: GET /edev/123/der (Discover associated DER instances)
    GW->>DER: Forward GET /edev/123/der
    DER-->>GW: 200 OK (DERList containing DER MRID)
    GW-->>Client: 200 OK (DERList)

    Client->>GW: GET /derp (Discover available DER Programs)
    GW->>DER: Forward GET /derp
    DER-->>GW: 200 OK (DERProgramList)
    GW-->>Client: 200 OK (DERProgramList)
```

---

## 2. Scheduling

In the Scheduling phase, the client fetches the details of the Frequency Response program, subscribes to control lists, and retrieves the specific Frequency-Watt control curves configured by the utility.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /derp/1 (Fetch target DERProgram details)
    GW->>DER: Forward GET /derp/1
    DER-->>GW: 200 OK (DERProgram with DefaultDERControl & Curve links)
    GW-->>Client: 200 OK (DERProgram)

    Client->>GW: GET /derp/1/dderc (Fetch Default DER Control)
    GW->>DER: Forward GET /derp/1/dderc
    DER-->>GW: 200 OK (DefaultDERControl linking active curves)
    GW-->>Client: 200 OK (DefaultDERControl)

    Client->>GW: GET /derp/1/dc/456 (Fetch target Frequency-Watt Curve)
    GW->>DER: Forward GET /derp/1/dc/456
    DER-->>GW: 200 OK (DERCurve: CurveType=0 [Frequency-Watt] + CurveData points)
    GW-->>Client: 200 OK (DERCurve)

    Client->>GW: GET /derp/1/derc (Query scheduled controls)
    GW->>DER: Forward GET /derp/1/derc
    DER-->>GW: 200 OK (DERControlList)
    GW-->>Client: 200 OK (DERControlList)
```

---

## 3. Operation

During the Operation phase, the client monitors real-time active controls, senses local grid frequency, and adjusts its active power output dynamically in accordance with the configured Frequency-Watt curve.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /edev/123/der/1/cdc (Fetch Current DER Controls)
    GW->>DER: Forward GET /edev/123/der/1/cdc
    DER-->>GW: 200 OK (CurrentDERControls detailing active frequency control)
    GW-->>Client: 200 OK (CurrentDERControls)

    Note over Client: Client senses local grid frequency (e.g., 59.95 Hz)
    Note over Client: Client interpolates active power change using Curve 456
    Note over Client: Client adjusts active power output (kW injection/shed)
```

---

## 4. Verify

In the Verification phase, the client submits status and availability telemetry to verify compliance with the requested control parameters.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: POST /edev/123/der/1/ders (Update DER Status)
    GW->>DER: Forward POST /edev/123/der/1/ders
    DER-->>GW: 201 Created (Stores status / connection telemetry)
    GW-->>Client: 201 Created

    Client->>GW: POST /edev/123/der/1/dera (Update DER Availability)
    GW->>DER: Forward POST /edev/123/der/1/dera
    DER-->>GW: 201 Created (Stores available capacity telemetry)
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, telemetry data posted via Mirror Usage Points is reconciled against customer agreements and active tariffs by the Billing service to calculate financial incentives or penalties.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Client->>GW: POST /mup/123 (Submit telemetry meter readings)
    GW->>MUP: Forward POST /mup/123
    MUP-->>GW: 201 Created (Stores readings)
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch delivered energy readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list)
    Bill->>Bill: Reconcile actual frequency response performance
    Bill->>Bill: Apply tariff rules & credit CustomerAccount
```
