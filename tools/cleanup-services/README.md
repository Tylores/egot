# Microservice Cleanup Tool

This tool helps identify and remove stale microservices that were created during development and testing.

## Overview

The egot project contains 15 core production microservices plus several infrastructure services. This tool:
- Identifies which services are production vs. stale
- Shows what would be deleted (dry-run by default)
- Safely removes stale services with confirmation
- Guides you through manual cleanup of routes.go

## Production Services (Protected)

These services are never removed:
- **BRS** - Bill Report Server (8010)
- **Bill** - Billing Service (8011)
- **DCAP** - Device Capability (8012)
- **DERP** - Demand Response Program (8013)
- **DR** - Demand Response (8014)
- **EDevice** - End Device (8015)
- **File** - File Service (8016)
- **MUP** - Meter Usage and Pricing (8017)
- **Messaging** - Messaging Service (8018)
- **Notify** - Notification Service (8019)
- **PPY** - Pricing and Pricing Yields (8020)
- **SDevice** - Smart Device (8021)
- **TariffProfile** - Tariff Profile Service (8022)
- **TimeOfUse** - Time of Use (8023)
- **UPT** - User Price Tracking (8024)

Plus infrastructure services:
- **Core** - Core SEP service (8000)
- **FlowReservation** - Flow Reservation (8001)
- **der** - Legacy DER service (8025)
- **rsps** - RSPS service (8003)
- **flow** - Flow service (8001)

## Usage

### List all services
```bash
go run ./tools/cleanup-services --list
```

Output shows:
- ✓ Production services (protected from deletion)
- ✗ Stale services (candidates for removal)

### Dry-run (show what would be deleted)
```bash
go run ./tools/cleanup-services --dry-run
```

This is the default behavior. Shows:
- Which directories would be removed
- Manual steps required for routes.go cleanup

### Interactive removal
```bash
go run ./tools/cleanup-services --dry-run=false --interactive=true
```

Prompts you to confirm before deleting.

### Remove all stale services automatically
```bash
go run ./tools/cleanup-services --dry-run=false --remove-all-stale
```

⚠️ No confirmation prompts - use with care!

## What Gets Removed

For each stale service:
1. `cmd/{ServiceName}/` - Main entry point and main.go
2. `internal/{ServiceName}/` - All handler, repository, and server code

What requires manual cleanup:
- Service constant from `internal/routes/routes.go` const block
- Service routes from `internal/routes/routes.go` serviceMap

## Example Workflow

1. **Identify stale services:**
   ```bash
   go run ./tools/cleanup-services --list
   ```

2. **Preview what will be deleted:**
   ```bash
   go run ./tools/cleanup-services
   ```

3. **Perform cleanup (with confirmation):**
   ```bash
   go run ./tools/cleanup-services --dry-run=false
   ```

4. **Manually update routes.go** - Follow the instructions printed by the tool

5. **Verify build:**
   ```bash
   go build ./...
   make build-all
   ```

6. **Commit changes:**
   ```bash
   git add -A
   git commit -m "cleanup: remove stale microservices"
   ```

## How Services Are Classified

Services are classified as:
- **Production**: Listed in `PRODUCTION_SERVICES` map in the tool source
- **Stale**: Any service in `cmd/` or `internal/` not in the production list

## Adding Production Services

If you create a new production microservice, add it to the `PRODUCTION_SERVICES` map in `tools/cleanup-services/main.go`:

```go
var PRODUCTION_SERVICES = map[string]bool{
    "MyNewService": true,
    // ... existing services ...
}
```

## Safety Notes

- Tool requires confirmation before deletion (can be disabled with `--remove-all-stale`)
- Always run with `--dry-run` first to see what will be deleted
- Make sure to commit any work to git before running
- The tool cannot remove git history - deleted files can be recovered from git
- Manual routes.go cleanup is required

## Troubleshooting

**Tool says "Not in egot repository root"**
- Make sure you're in the /home/tylor/dev/egot directory
- Look for go.mod file

**Tool finds services you want to keep**
- Add them to `PRODUCTION_SERVICES` map in main.go
- Rebuild and re-run

**routes.go changes cause compilation errors**
- Carefully remove only the stale service entries
- Each service has a constant and one or more serviceMap entries
- Use grep to find all occurrences: `grep -n "ServiceName" internal/routes/routes.go`

## Related

- Makefile: `make cleanup-services` or `make cleanup-services-dry-run`
- Routes definition: `internal/routes/routes.go`
- Scaffold generator: `cmd/scaffold-gen`
