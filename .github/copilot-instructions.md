# Copilot Instructions

## Build, test, and lint commands

- Build everything: `go build ./...`
- Build repository binaries via Make: `make build-all`
- Build one binary: `make build-Bill` or `go build -o ./bin/Bill ./cmd/Bill`
- Run the full test suite: `go test ./...` or `make test`
- Run the fast test suite: `make test-unit`
- Run route-focused tests: `make test-routes`, `make test-route-registration`, `make test-http-methods`
- Run one package: `make test-package PKG=./cmd/DERP` or `make test-package PKG=./internal/Bill/handler`
- Run one named test: `go test ./cmd/DERP -run TestDerpHTTPMethods -v` or `go test ./internal/Bill/handler -run TestNewHandler -v`
- Lint: `make lint` (this only runs when `golangci-lint` is installed)
- Show available repo workflows: `make help`

## High-level architecture

- This repository implements an IEEE 2030.5 / SEP 2 service set in Go. The XML models live in `sep/`, and URI constants for the core service live in `sep/uri/`.
- Service binaries are under `cmd/`. `cmd/core` is the top-level SEP entrypoint, while many other services (`BRS`, `Bill`, `DCAP`, `DERP`, `DR`, `EDevice`, `File`, `MUP`, `Messaging`, `Notify`, `PPY`, `SDevice`, `TariffProfile`, `TimeOfUse`, `UPT`) are separate microservices with their own ports.
- For each service, `cmd/<Service>/main.go` wires mTLS, constructs the repository and handler, and registers routes directly with Go 1.22+ HTTP patterns like `http.Handle("GET /bill/{id1}", ...)`. The corresponding implementation lives under `internal/<Service>/handler` and `internal/<Service>/repository/memory`.
- `internal/routes/routes.go` is the cross-service routing map. It defines host:port constants and a `serviceMap` from SEP resource paths to the service responsible for them. If a change adds or moves a resource path, update this file too.
- `internal/tlsutil` owns the shared server TLS setup. Services are expected to use it so client certificate verification is consistent across binaries.
- WADL drives the generated-service workflow. `cmd/wadl-extract` slices `wadl/sep_wadl.xml` into service-specific WADLs, and `cmd/scaffold-gen` turns those WADLs into `cmd/<service>` plus `internal/<service>` scaffolds and updates `internal/routes/routes.go`.
- Route registration tests under `cmd/*/*_route_registration_test.go` verify that registered method/path pairs match generated expectations. The shared test helpers for this live in `test/routing/`.

## Key conventions

- Keep server TLS configuration centralized: use `tlsutil.NewServerConfig("./ssl")` for services instead of creating inline `tls.Config` values.
- Services assume local mTLS assets in `./ssl`. `make ssl-refresh` regenerates the CA, server certs, and default client cert, and many services call `repo.InitRepository("./ssl")` to load known client certificates into the in-memory repository.
- Handler methods follow a strict naming pattern: `<HTTP_METHOD><ResourceName>`, for example `GETCustomerAccountList` or `DELETECustomerAgreement`. Route registrations in `cmd/<Service>/main.go` are expected to mirror those names exactly.
- SEP path parameters consistently use numbered names like `{id1}`, `{id2}`, `{id3}`. When parsing request paths, handlers use `req.PathValue("id1")`, `req.PathValue("id2")`, etc. Match that numbering style when adding routes or WADL resources.
- Handler flow is consistent across services: validate the client certificate-derived LFDI first, validate numeric path parameters with `strconv.Atoi`, then set headers before calling `WriteHeader`, then encode SEP XML bodies. POST-to-list handlers generally set a `location` header; disallowed method/resource combinations return `405`.
- Prefer regenerating route registration code from existing handlers or WADL inputs instead of hand-editing large blocks of `http.Handle(...)` lines. Relevant scripts are `scripts/regenerate_cmd_handlers.sh` and the WADL-based generation workflow.
