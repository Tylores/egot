# WADL Extractor

A command-line tool to extract and cluster portions of a WADL (Web Application Description Language) specification file based on resource path prefixes.

## Overview

The WADL Extractor extracts resources from a WADL file by path prefix. For large path trees like `/edev` (53 routes across 18 logical sub-trees), it supports a **config-file-driven grouping** workflow that lets you define named service groups with validation and port assignment. This is useful for:

- Extracting specific API modules (e.g., demand-response, end-device)
- Creating focused WADL specifications for service implementation
- Grouping clusters into logical named services via a version-controlled YAML config
- Splitting large path trees into per-cluster WADLs for individual microservices

## Usage

```bash
# Single or multi-prefix extract
wadl-extract -wadl <input_wadl> -path <prefix>[,<prefix>...] -output <output_wadl> [-v]

# Show suggested clusters (dry-run, no files written)
wadl-extract -wadl <input_wadl> -path <prefix> -suggest [-depth <n>]

# Write one WADL per cluster into a directory (auto-split)
wadl-extract -wadl <input_wadl> -path <prefix> -split <dir> [-depth <n>] [-v]

# Generate a starter grouping config
wadl-extract -wadl <input_wadl> -path <prefix> -init-config <config.yaml> [-depth <n>]

# Run a grouping config (the recommended workflow for large paths)
wadl-extract -config <config.yaml> [-v]
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `-wadl string` | extract/suggest/split/init | — | Path to input WADL specification file |
| `-path string` | extract/suggest/split/init | — | Comma-separated path prefix(es) |
| `-output string` | extract mode | — | Output path for the extracted WADL file |
| `-suggest` | — | false | Print cluster table (no files written; single prefix only) |
| `-split string` | — | — | Directory to write one WADL per cluster (single prefix only) |
| `-depth int` | — | `2` | Clustering depth for `-suggest`, `-split`, `-init-config` |
| `-init-config string` | — | — | Write a starter YAML grouping config to this path |
| `-config string` | — | — | Run a YAML grouping config; generates one WADL per group |
| `-v` | — | false | Verbose output (show each resource path) |

## Grouping Workflow (Recommended for Large Paths)

Use this workflow when a path like `/edev` has too many clusters for a single service and you want repeatable, version-controlled groupings.

### Step 1 — Explore clusters

```bash
wadl-extract -wadl wadl/sep_wadl.xml -path /edev -suggest
```

```
Clusters for /edev at depth 2 (18 clusters, 53 total paths):
  [core        ]   2 paths  → /edev, ...
  [der         ]  12 paths  → /edev/{id1}/der, ...
  [ns          ]  12 paths  → /edev/{id1}/ns, ...
  [cfg         ]   3 paths  → /edev/{id1}/cfg, ...
  ...
```

### Step 2 — Generate starter config

```bash
wadl-extract -wadl wadl/sep_wadl.xml -path /edev -init-config edev-split.yaml
```

This writes `edev-split.yaml` with all 18 clusters listed as comments. Every cluster must be assigned before the config will validate.

### Step 3 — Edit the config

```yaml
# edev-split.yaml
wadl: wadl/sep_wadl.xml
base: /edev
depth: 2
output: wadl/edev/

groups:
  edev-core:
    port: 8015
    clusters: [core, rg, dstat, aggp, fs, ps]
  edev-der:
    port: 8020
    clusters: [der]
  edev-network:
    port: 8021
    clusters: [ns, cfg, di, adev, prxy]
  edev-scheduling:
    port: 8022
    clusters: [frq, frp, fsa, lel, lsl, sub]
```

### Step 4 — Generate WADLs

```bash
wadl-extract -config edev-split.yaml
```

```
Generating 4 group WADLs from wadl/sep_wadl.xml → wadl/edev/
  ✓ [edev-core           ] port=8015    7 paths → wadl/edev/edev-core.wadl
  ✓ [edev-der            ] port=8020   12 paths → wadl/edev/edev-der.wadl
  ✓ [edev-network        ] port=8021   22 paths → wadl/edev/edev-network.wadl
  ✓ [edev-scheduling     ] port=8022   12 paths → wadl/edev/edev-scheduling.wadl
```

Each output WADL has the configured `port` injected and is immediately usable with `scaffold-gen`.

### Validation

The tool validates the config before writing any files:

```
config validation failed:
  - cluster key "derr" in group "edev-der" not found (available: adev, aggp, cfg, core, der, ...)
  - cluster key "der" assigned to both "edev-der" and "edev-network"
  - unassigned clusters: aggp, dstat — every cluster must be in a group
```

Rules enforced:
- Every cluster key referenced in the config must exist at the given `depth`
- No cluster key may appear in more than one group
- Every cluster discovered by `-suggest` must be assigned to exactly one group

## Common Resource Paths

| Path | Resources | Description |
|------|-----------|-------------|
| `/bill` | ~10 | Customer billing information |
| `/brs` | 4 | Billing reading sets |
| `/dcap` | 1 | Device capability |
| `/derp` | 6 | DER programs |
| `/dr` | 5 | Demand response programs |
| `/edev` | 53 | End device resources (18 sub-clusters at depth 2) |
| `/file` | 2 | File management |
| `/msg` | 5 | Messaging programs |
| `/mup` | 2 | Mirror usage points |
| `/ntfy` | 2 | Notifications |
| `/ppy` | 9 | Prepayment |
| `/rsps` | 3 | Response sets |
| `/sdev` | 1 | Self device |
| `/tm` | 1 | Time |
| `/tp` | 7 | Tariff profiles |
| `/upt` | 6 | Usage points |

## Other Examples

### Extract a single prefix

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /dr -output wadl/dr.wadl -v
```

### Combine specific clusters manually

```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path "/edev/{id1}/der,/edev/{id1}/ns" \
  -output wadl/edev-der-ns.wadl
```

### Auto-split into one WADL per cluster

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /edev -split wadl/edev/
```

## Building

```bash
go build -o bin/wadl-extract ./cmd/wadl-extract
```

## Integration with scaffold-gen

```bash
# Extract one service
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /dr -output wadl/dr.wadl
go run ./cmd/scaffold-gen -wadl wadl/dr.wadl -output .

# Scaffold all groups from a config
./bin/wadl-extract -config edev-split.yaml
for wadl in wadl/edev/*.wadl; do
  go run ./cmd/scaffold-gen -wadl "$wadl" -output .
done
```

See [docs/tools/scaffold-gen.md](../../docs/tools/scaffold-gen.md) for the scaffold generator guide.

## How It Works

### Extract mode
1. Filters all resources whose `samplePath` starts with any provided prefix
2. Preserves WADL header elements (doc, grammars) in the output
3. Injects `port` and `max_entities` attributes into `<application>`

### Config mode (`-config`)
1. Loads and validates the YAML config (fail-fast before any files are written)
2. Runs `-suggest` internally to discover cluster keys at the configured `depth`
3. For each group, collects exact resources from the named clusters (no prefix bleed-through)
4. Writes one `<group-name>.wadl` per group with the configured `port` injected

### Clustering depth
- depth 1 → groups by first segment after base (e.g. `{id1}`) — too coarse for `/edev`
- depth 2 → groups by second segment (`der`, `ns`, `cfg`, ...) — natural service boundaries
- depth 3 → finer splits within `ns`, `der`, etc.

## Implementation

- **`internal/wadlext/extractor.go`** — `Extract()`, `ExtractMany()`, `SuggestClusters()`, `ExtractCluster()`, `ExtractClusters()`
- **`internal/wadlext/config.go`** — `SplitConfig`, `GroupConfig`, `LoadSplitConfig()`, `ValidateConfig()`, `GenerateInitConfig()`
- **`cmd/wadl-extract/main.go`** — CLI flags and mode dispatch

## Exit Codes

- `0` — Success
- `1` — Error (invalid arguments, file not found, parse error, validation failure, no resources found)

## Related

- Original WADL file: `wadl/sep_wadl.xml`
- SEP 2.2 specification: https://standards.ieee.org/standard/2030_5-2018.html

