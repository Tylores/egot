# ESI Lifecycle Sequence Diagrams: Voltage Management

This document covers the client/server interactions and sequence diagrams for the **Voltage Management** grid service (Volt-Var control / reactive power adjustment based on local voltage deviations) across all 5 ESI lifecycle phases, aligned with IEEE 2030.5 CSIP test procedures.

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

    Client->>GW: GET /derp (Discover available DER Programs)
    GW->>DER: Forward GET /derp
    Note over DER: [CORE-010] / [CORE-011] Query FSA assigned DER Programs
    DER-->>GW: 200 OK (DERProgramList with pollRate=900)
    GW-->>Client: 200 OK (DERProgramList)
```

---

## 2. Scheduling

In the Scheduling phase, the client monitors and parses future planned Volt-Var curves (`CurveType=11`) using either Polling or Subscriptions.

### Option A: Polling Interaction
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Note over Client: [CORE-003] Client polls at interval defined by pollRate (e.g. 900s)
    Client->>GW: GET /derp/1 (Fetch target DERProgram details)
    GW->>DER: Forward GET /derp/1
    DER-->>GW: 200 OK (DERProgram detailing DefaultDERControl & Curve links)
    GW-->>Client: 200 OK (DERProgram)

    Client->>GW: GET /derp/1/dderc (Fetch Default DER Control)
    GW->>DER: Forward GET /derp/1/dderc
    Note over DER: [CORE-012] Provide DefaultDERControl linking active curves
    DER-->>GW: 200 OK (DefaultDERControl)
    GW-->>Client: 200 OK (DefaultDERControl)

    Client->>GW: GET /derp/1/dc/990 (Fetch Volt-Var Curve)
    GW->>DER: Forward GET /derp/1/dc/990
    Note over DER: [CORE-012] Provide DERCurve: CurveType=11 (Volt-Var) & points
    DER-->>GW: 200 OK (DERCurve containing 10 points mapping Volt to % Max reactive power)
    GW-->>Client: 200 OK (DERCurve)

    Client->>GW: GET /derp/1/derc (Query scheduled controls list)
    GW->>DER: Forward GET /derp/1/derc
    DER-->>GW: 200 OK (DERControlList showing future active curve changes)
    GW-->>Client: 200 OK (DERControlList)
```

### Option B: Subscription/Notification
```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: POST /derp/1/derc/sub (Subscribe to DERControlList changes)
    GW->>DER: Forward POST /derp/1/derc/sub (Includes Client Notification URI)
    DER-->>GW: 201 Created (Location: /derp/1/derc/sub/6)
    GW-->>Client: 201 Created (Location: /derp/1/derc/sub/6)

    Note over DER: [CORE-018] / [CORE-019] Event curve is scheduled by Grid Operator
    DER->>GW: POST /client/notification (Notify client of resource change)
    GW->>Client: Forward POST /client/notification (Contains updated DERControlList)
    Client-->>GW: 204 No Content
    GW-->>DER: 204 No Content
```

---

## 3. Operation

During the Operation phase, the client tracks active controls, applies event randomization, manages priorities, and executes dynamic Volt-Var reactive power adjustments.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant DER as "DER Service (:8026)"

    Client->>GW: GET /derp/1/actderc (Fetch Active DER Controls)
    GW->>DER: Forward GET /derp/1/actderc
    DER-->>GW: 200 OK (ActiveDERControl detailing active voltage control curves)
    GW-->>Client: 200 OK (ActiveDERControl)

    Note over Client: [CORE-021] Client applies randomizeStart / randomizeDuration values
    Note over Client: [BASIC-021] Priority Check: SP (Service Point) level supersedes SY (System) level
    Note over Client: Client senses local grid voltage (e.g. 242V on a 240V nominal grid)
    Note over Client: Client interpolates reactive power adjustment based on DERCurve (CurveType=11)
    Note over Client: Client adjusts reactive power output (kVAr absorption/injection) dynamically
```

---

## 4. Verify

In the Verification phase, the client posts status confirmations, reports execution states, updates statuses, and logs alarms to verify curve enforcement.

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

    Note over Client: Event starts executing (Voltage management is active)
    Client->>GW: POST /rsps/123/rsp (Post Status = 2 [Started])
    GW->>RSPS: Forward POST /rsps/123/rsp
    RSPS-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Client: [BASIC-028] Inverter Status Update
    Client->>GW: PUT /der/1/ders (Periodically report device/inverter status)
    GW->>DER: Forward PUT /der/1/ders
    DER-->>GW: 204 No Content
    GW-->>Client: 204 No Content

    Note over Client: [BASIC-027] Alarm Logging (If grid voltage limits are exceeded)
    Client->>GW: POST /edev/123/lel (POST LogEvent for out-of-bounds voltage alert)
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

In the Settlement phase, telemetry data is evaluated by the Billing service to credit the customer for active/reactive voltage support.

> [!NOTE]
> Unlike standard active-power billing, Voltage Management verification requires tracking reactive energy exchange. The client packages:
> - **Reactive Energy (VARh)**: Captured using the standard IEEE 2030.5 `MirrorMeterReading` with the appropriate unit multiplier and type (Reactive Power).
> - **Voltage Telemetry (V)**: Transmitted as part of the `MirrorMeterReading` readings or device status parameters to allow the settlement engine to correlate Volt-Var curve execution with local voltage levels.

```mermaid
sequenceDiagram
    autonumber
    participant Client as "End Device / DER Client"
    participant GW as "Nginx API Gateway"
    participant MUP as "MUP Service (:8017)"
    participant Bill as "Bill Service (:8011)"

    Note over Client: [BASIC-029] Client periodically posts cumulative Wh, VARh, and V telemetry
    Client->>GW: POST /mup/123 (Submit MirrorUsagePoint payload containing VARh and V readings)
    GW->>MUP: Forward POST /mup/123 (XML payload validation)
    MUP-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over Bill: Periodic billing process runs
    Bill->>MUP: GET /mup/123 (Fetch delivered reactive support and local voltage readings)
    MUP-->>Bill: 200 OK (UsagePoint readings list containing VARh/V)
    
    Bill->>Bill: Reconcile actual delivered active/reactive energy against Volt-Var curve expectations at observed voltages
    Bill->>Bill: Fetch customer agreement active billing periods (/bill/123/ca/1/actbp)
    Bill->>Bill: Reconcile & credit CustomerAccount
```
