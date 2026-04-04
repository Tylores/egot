# Microservice Scaffold Generator

A Go-based tool for rapidly generating complete microservices from WADL specifications.

## Quick Links

- **[SCAFFOLD_GENERATOR.md](./SCAFFOLD_GENERATOR.md)** - Start here! Overview and usage guide
- **[SCAFFOLD_SPEC.md](./SCAFFOLD_SPEC.md)** - Complete WADL specification and schema
- **[cmd/scaffold-gen/README.md](./cmd/scaffold-gen/README.md)** - Quick reference and examples
- **[wadl/device-manager.wadl](./wadl/device-manager.wadl)** - Example WADL file

## What It Does

The scaffold generator creates complete microservices from WADL (Web Application Description Language) XML specifications. Each generated service includes:

- HTTP handlers with TLS certificate validation
- Thread-safe in-memory repository
- Mutual TLS server setup
- Automatic route registration
- All code generated and ready to build

## Get Started in 3 Steps

### 1. Create WADL Specification

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application name="my-service" port="8003" max_entities="50">
  <resources>
    <resource path="/items">
      <methods>
        <method name="GetItems" http_method="GET" response_type="ItemList"/>
      </methods>
    </resource>
  </resources>
</application>
```

### 2. Generate Scaffold

```bash
go run ./cmd/scaffold-gen -wadl wadl/my-service.wadl -output .
```

### 3. Implement & Deploy

```bash
# Edit generated files to add business logic
vim internal/my-service/handler/handler.go
vim internal/my-service/repository/memory/memory.go

# Build and run
go build ./cmd/my-service
./my-service
```

## Generated Structure

```
cmd/{service}/
  main.go                      # Entry point with TLS & routes

internal/{service}/
  handler/handler.go          # HTTP handlers
  repository/
    memory/memory.go          # Storage
    error.go                  # Errors
  server/server.go            # Utilities
```

## Features

✓ XML WADL parsing  
✓ Complete code generation  
✓ Mutual TLS support  
✓ Thread-safe storage  
✓ Automatic route registration  
✓ Path parameter extraction  
✓ Request/response type specification  
✓ Entity pool management  

## Architecture

Each generated service follows the established pattern from core and flowreservation:

1. **Handler** - Validates client certificate, processes requests
2. **Repository** - Thread-safe in-memory storage
3. **Server** - TLS configuration and utilities
4. **Main** - Entry point with route setup

## WADL Format

Simple XML specification:

```xml
<application name="service-name" port="8002" max_entities="50">
  <resources>
    <resource path="/resource/{id}">
      <methods>
        <method name="GetResource" http_method="GET" response_type="Type">
          <path_params>
            <param name="id" type="int"/>
          </path_params>
        </method>
      </methods>
    </resource>
  </resources>
</application>
```

## Documentation

- **Overview** → SCAFFOLD_GENERATOR.md
- **Full Spec** → SCAFFOLD_SPEC.md  
- **Quick Start** → cmd/scaffold-gen/README.md
- **Example** → wadl/device-manager.wadl

## For More Information

See SCAFFOLD_GENERATOR.md for complete overview and usage guide.

---

**Status**: ✓ Production Ready  
**Build**: `go build ./cmd/scaffold-gen`  
**Usage**: `go run ./cmd/scaffold-gen -wadl <wadl> -output <dir>`
