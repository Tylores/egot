# WADL Extractor

A command-line tool to extract portions of a WADL (Web Application Description Language) specification file based on resource path prefixes.

## Overview

The WADL Extractor extracts all resources from a WADL file that match a given path prefix, creating a new WADL file containing only those resources. This is useful for:

- Extracting specific API modules (e.g., demand-response, end-device)
- Creating focused WADL specifications for service implementation
- Splitting large WADL files into smaller, more manageable pieces

## Usage

```bash
wadl-extract -wadl <input_wadl> -path <path_prefix> -output <output_wadl> [-v]
```

### Flags

- `-wadl string` - **Required**. Path to input WADL specification file
- `-path string` - **Required**. Resource path prefix to extract (e.g., `/dr`, `/edev`, `/msg`)
- `-output string` - **Required**. Output path for extracted WADL file
- `-v` - Optional. Enable verbose output showing all extracted resources

## Common Resource Paths

| Path | Resources | Description |
|------|-----------|-------------|
| `/bill` | ~10 | Customer billing information |
| `/brs` | 4 | Billing reading sets |
| `/dcap` | 1 | Device capability |
| `/derp` | 6 | DER programs |
| `/dr` | 5 | Demand response programs |
| `/edev` | 53 | End device resources |
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

## Examples

### Extract demand-response resources

```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /dr \
  -output wadl/demand-response.wadl \
  -v
```

Output:
```
✓ Loaded WADL from: wadl/sep_wadl.xml
✓ Found 5 resources with prefix '/dr'
  - /dr
  - /dr/{id1}
  - /dr/{id1}/actedc
  - /dr/{id1}/edc
  - /dr/{id1}/edc/{id2}
✓ Extracted WADL saved to: wadl/demand-response.wadl
  Resources: 5
  Path prefix: /dr
```

### Extract end-device resources

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /edev -output wadl/end-device.wadl
```

### Extract messaging resources

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /msg -output wadl/messaging.wadl -v
```

### Extract usage points

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /upt -output wadl/upt.wadl
```

## Building

From the project root:

```bash
go build -o bin/wadl-extract ./cmd/wadl-extract
```

## Integration with scaffold-gen

Extract a WADL slice, then generate a full service from it:

```bash
./bin/wadl-extract -wadl wadl/sep_wadl.xml -path /dr -output wadl/dr.wadl -v
go run ./cmd/scaffold-gen -wadl wadl/dr.wadl -output .
```

See [docs/tools/scaffold-gen.md](../../docs/tools/scaffold-gen.md) for the scaffold generator guide.

## How It Works

1. **Loads** the input WADL file and parses it as XML
2. **Filters** all resource elements by matching their `samplePath` attribute against the provided prefix
3. **Preserves** WADL header elements (doc, grammars) in the output
4. **Generates** a new valid WADL file with only the matching resources
5. **Saves** the extracted WADL with proper XML formatting

## Implementation

The tool uses two main components:

- **`internal/wadlext/extractor.go`** - Core extraction logic that handles WADL parsing and resource filtering
- **`cmd/wadl-extract/main.go`** - Command-line interface

## Exit Codes

- `0` — Success
- `1` — Error (invalid arguments, file not found, parsing error, no resources found)

## Error Reference

| Error | Cause | Fix |
|-------|-------|-----|
| `"failed to open WADL file"` | Input file doesn't exist | Verify `-wadl` path |
| `"No resources found"` | No resources match the prefix | Check `-path` against the table above |
| `"failed to write WADL file"` | Output directory missing or unwritable | Verify `-output` directory exists |

## Limitations

- Extraction is based on exact path prefix matching
- Resources must have a `samplePath` attribute to be included
- All matched resources are included; you cannot selectively include/exclude specific resources

## Related

- Original WADL file: `wadl/sep_wadl.xml`
- SEP 2.2 specification: https://standards.ieee.org/standard/2030_5-2018.html
