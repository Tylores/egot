# sep-golang
Smart Energy Profile using Go

## Setup

### Host
Add the server address to your hosts file for local testing:

```shell
127.0.0.1   egot.internal.com
```

### SSL

Generate mTLS certificates for local development using the included script (requires `openssl`):

```shell
make ssl-refresh
```

This creates a self-signed CA and issues a server certificate for `egot.internal.com` plus a default client certificate, all valid for 1 year. The CA itself is valid for 10 years.

To generate additional numbered client certificates (e.g. 5 clients):

```shell
make ssl-clients N=5
# Generates ssl/client-0001.crt ... ssl/client-0005.crt
```

> **Note:** The original setup used [step-cli](https://smallstep.com/docs/step-cli/reference/) with a Smallstep CA. The `make ssl-refresh` approach uses plain `openssl` and requires no external CA service. See [docs/development/microservices.md](docs/development/microservices.md) for details on the generated files.

### Go

```shell
go mod tidy
go run ./cmd/crawler
```

Import the SEP 2.0 models package:

```go
import "github.com/Tylores/egot/sep"
```

## Documentation

| Topic | Location |
|-------|----------|
| Microservices (services, ports, building) | [docs/development/microservices.md](docs/development/microservices.md) |
| HTTP handlers (patterns, WADL compliance) | [docs/development/handlers.md](docs/development/handlers.md) |
| Testing | [docs/testing/README.md](docs/testing/README.md) |
| Testing quick reference | [docs/testing/quick-reference.md](docs/testing/quick-reference.md) |
| Route testing | [docs/testing/routing.md](docs/testing/routing.md) |
| CI/CD integration | [docs/testing/ci-cd.md](docs/testing/ci-cd.md) |
| Scaffold generator | [docs/tools/scaffold-gen.md](docs/tools/scaffold-gen.md) |
| WADL extractor | [cmd/wadl-extract/README.md](cmd/wadl-extract/README.md) |
| WADL extractor | [cmd/wadl-extract/README.md](cmd/wadl-extract/README.md) |
| SEP 2 resource paths | [wadl/RESOURCE_URI_PATHS.md](wadl/RESOURCE_URI_PATHS.md) |

## Grid Services & Simulation

### Operator Service
The `operator` service handles high-level grid coordination:
- **Scheduling**: Uses a Greedy Algorithm to match `FlowReservationRequests` with grid needs.
- **Settlement**: Compares scheduled DER performance against reported telemetry.

### OpenDSS Simulation
The platform includes a closed-loop simulation environment using OpenDSS:
- `scripts/egot_sim.py`: Advanced simulation script supporting DER-to-node mapping.
- `model/IEEE13Nodeckt.dss`: Standard distribution feeder model for validation.

To run a simulation:
```shell
python3 scripts/egot_sim.py
```
