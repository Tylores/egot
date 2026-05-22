# EGoT PhD Research Roadmap

This roadmap outlines the task-oriented steps required to complete the PhD dissertation and the EGoT (Energy Grid of Things) implementation, now focused on GMLC Common Grid Services.

## Phase 1: GMLC Grid Service Implementation
*Focus: Aligning the microservices fleet with IEEE 2030.5 resources to serve GMLC requirements.*

- [x] **Energy Scheduling & 5-minute Regulation (EIM)**
    - [x] Task 1.1: Update `FlowReservation` service to support 5-minute interval scheduling for Energy Imbalance Market (EIM) emulation.
    - [x] Task 1.2: Implement logic in `DemandResponse` to allow devices to reserve load capacity prior to flow reservation requests.
- [x] **Blackstart & Reserve Services**
    - [x] Task 1.3: Enhance `EDevice` microservice to calculate and expose **Load Shed Availability** telemetry for real-time reserve monitoring.
    - [x] Task 1.4: Implement targeted `DemandResponse` dispatch logic for coordinating generation-load balance during cold-start (Blackstart) scenarios.
- [x] **Frequency & Voltage Management (Implementation Only)**
    - [x] Task 1.5: Implement **Frequency-Watt** and **Volt-Var** curve support in the `DER` service for local transient stability.
    - [x] Task 1.6: Develop a **Feeder-Aware Dispatcher** in the `Operator` service that utilizes the grid topology to issue localized `DERControl` events.
- [x] **Operator Service Orchestration**
    - [x] Task 1.7: Integrate the `GreedyScheduler` to handle both 24-hour day-ahead and 5-minute EIM windows.
    - [x] **Gemini Support**: Generate Go boilerplate for inter-service coordination between `FlowReservation` and `DERP`.

## Phase 2: GMLC Service Validation (Simulation)
*Focus: Evaluating grid service performance via steady-state simulation (excluding transient frequency/voltage).*

- [x] **Scenario 2.1: Day-Ahead & 5-min EIM Regulation**
    - [x] Task 2.1: Define a simulation window using the IEEE 13-node model featuring hourly scheduling and 5-minute imbalance corrections.
    - [x] Task 2.2: Validate the coordination between `FlowReservation` and `DERControl` dispatch.
- [x] **Scenario 2.2: Blackstart Coordination**
    - [x] Task 2.3: Execute a "Cold-Start" scenario where `DemandResponse` manages load pick-up for devices within specific feeder segments.
    - [x] Task 2.4: Measure the impact of managed load pick-up on generation ramp-rate requirements.
- [x] **Scenario 2.3: Reserve Capacity Validation**
    - [x] Task 2.5: Trigger a simulated contingency and compare "Load Shed Availability" telemetry from `EDevice` against actual DR performance.
    - [x] Task 2.6: Analyze "Grid Service Reliability" metrics using the `SettlementEngine`.
- [x] **Data Export & Visualization**
    - [x] Task 2.7: Generate publication-ready figures for EIM and Blackstart scenarios using `generate_dissertation_plots.py`.

## Phase 3: Dissertation Composition (LaTeX)
*Focus: Translating technical success into academic contribution.*

- [ ] **Chapter 4: Results and Discussion**
    - [ ] Task 3.1: Document the impact of 5-minute EIM on grid stability and settlement.
    - [ ] Task 3.2: Discuss the effectiveness of DR-coordinated Blackstart and Reserve services.
    - [ ] Task 3.3: Analyze the "Grid Service Reliability" metrics from the Settlement Engine.
    - [ ] **Gemini Support**: Drafting section summaries and formatting complex LaTeX tables from CSV data.
- [ ] **Chapter 2: Literature Review Refinement**
    - [ ] Task 3.4: Expand the IEEE 2030.5 vs. OpenADR comparison for GMLC services.
    - [ ] Task 3.5: Strengthen the "Energy Service Interface" (ESI) theoretical framework.
    - [ ] **Gemini Support**: Summarizing recent papers (2024-2025) on DER interoperability.

## Phase 4: Review & Defense Preparation
*Focus: Final presentation and quality assurance.*

- [ ] **Internal Review**
    - [ ] Task 4.1: Submit individual chapters to the advisor/committee for feedback.
    - [ ] Task 4.2: Address technical comments and refine implementation details.
- [ ] **Defense Presentation**
    - [ ] Task 4.3: Create slides highlighting the "GMLC Grid Service Orchestration" in IEEE 2030.5.
    - [ ] Task 4.4: Prepare a live or recorded demo of the EGoT fleet in action.
    - [ ] **Gemini Support**: Drafting speaker notes and creating architecture diagrams (Mermaid).
