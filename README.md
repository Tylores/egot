# sep-golang
Smart Energy Profile using Go

## Setup

### Host
Add the server address to your hosts file for local testing:

```shell
127.0.0.1   egot.internal.com
```

### SSL

Use [step-cli](https://smallstep.com/docs/step-cli/reference/) to set yourself as a CA and generate mutual TLS certificates:

* https://smallstep.com/hello-mtls/doc/combined/go/go

Extend the default certificate duration (24h is too short for development):

```shell
step ca provisioner update you@smallstep.com \
   --x509-min-dur=24h \
   --x509-max-dur=8760h \
   --x509-default-dur=8760h
```

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
| SEP 2 resource paths | [wadl/RESOURCE_URI_PATHS.md](wadl/RESOURCE_URI_PATHS.md) |
