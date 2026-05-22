# ESI Lifecycle Sequence Diagrams: Voltage Management

This document covers the client/server interactions and sequence diagrams for the **Voltage Management** grid service (Volt-Var control / reactive power adjustment based on local voltage deviations) across all 5 ESI lifecycle phases.

---

## 1. Registration

In the Registration phase, the client registers as an End Device, registers its specific DER settings to define its active and reactive power capabilities, and discovers available Volt-Var management programs.

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

    Client->>GW: GET /derp (Discover available DER Programs)
    GW->>DER: Forward GET /derp
    DER-->>GW: 200 OK (DERProgramList)
    GW-->>Client: 200 OK (DERProgramList)
```

---

## 2. Scheduling

In the Scheduling phase, the client fetches the Volt-Var controls, schedules, and active Volt-Var curves (`CurveType=11`) that define the target voltage setpoints and reactive power injection thresholds.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /derp/1 (Fetch target DERProgram details)
    GW->>DER: Forward GET /derp/1
    DER-->>GW: 200 OK (DERProgram detailing Volt-Var controls & curve links)
    GW-->>Client: 200 OK (DERProgram)

    Client->>GW: GET /derp/1/derc/789 (Fetch scheduled Volt-Var DERControl)
    GW->>DER: Forward GET /derp/1/derc/789
    DER-->>GW: 200 OK (DERControl linking Volt-Var curve)
    GW-->>Client: 200 OK (DERControl)

    Client->>GW: GET /derp/1/dc/990 (Fetch Volt-Var Curve)
    GW->>DER: Forward GET /derp/1/dc/990
    DER-->>GW: 200 OK (DERCurve: CurveType=11 [Volt-Var] + CurveData voltage-reactive points)
    GW-->>Client: 200 OK (DERCurve)
```

---

## 3. Operation

During the Operation phase, the client monitors active controls, senses local grid voltage, and automatically adjusts its reactive power output (kVAr injection or absorption) matching the Volt-Var curve profile.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /edev/123/der/1/cdc (Fetch Current DER Controls)
    GW->>DER: Forward GET /edev/123/der/1/cdc
    DER-->>GW: 200 OK (CurrentDERControls with active Volt-Var mode)
    GW-->>Client: 200 OK (CurrentDERControls)

    Note over Client: Client senses local grid voltage (e.g., 242V on a 240V nominal grid)
    Note over Client: Client calculates reactive power adjustment from Curve 990 (absorb VARs)
    Note over Client: Client adjusts reactive power output (kVAr absorption)
```

---

## 4. Verify

In the Verification phase, the client posts confirmation responses when a control is received, and uploads real-time operation status to verify curve enforcement.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant RSPS as "Rsps Service (:8041)"
    participant DER as "DER Service (:8026)"

    Client->>GW: POST /rsps/123/rsp (Post control execution status)
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created (Stores Response payload)
    GW-->>Client: 201 Created

    Client->>GW: POST /edev/123/der/1/ders (Update DER Status)
    GW->>DER: Forward POST /edev/123/der/1/ders
    DER-->>GW: 201 Created (Stores current reactive power and voltage status)
    GW-->>Client: 201 Created
```

---

## 5. Settlement

In the Settlement phase, the telemetry reported via Mirror Usage Points is processed by the Billing service to reward or credit the customer account for providing reactive voltage support.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Client->>GW: POST /mup/123 (Submit telemetry meter readings)
    GW->>MUP: Forward POST /mup/123
    MUP-->>GW: 201 Created (Stores VARh and Wh readings)
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch delivered reactive support readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list)
    Bill->>Bill: Reconcile voltage management performance against agreement
    Bill->>Bill: Apply reactive support tariff & credit CustomerAccount
```
