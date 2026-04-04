# WADL Extractor - Usage Examples

## Overview

The WADL Extractor (`wadl-extract`) is a command-line tool that extracts portions of a WADL (Web Application Description Language) specification file based on resource path prefixes.

## Building

From the project root:

```bash
go build -o bin/wadl-extract ./cmd/wadl-extract
```

The binary will be available at `bin/wadl-extract`.

## Quick Examples

### Extract Demand Response resources
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

### Extract End Device resources
```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /edev \
  -output wadl/end-device.wadl
```

### Extract Messaging resources
```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /msg \
  -output wadl/messaging.wadl \
  -v
```

### Extract Usage Point resources
```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /upt \
  -output wadl/usage-point.wadl
```

### Extract Billing resources
```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /bill \
  -output wadl/billing.wadl
```

### Extract DER Program resources
```bash
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /derp \
  -output wadl/der-program.wadl
```

## Flags

- `-wadl string` (**Required**) - Path to input WADL specification file
- `-path string` (**Required**) - Resource path prefix to extract (e.g., `/dr`, `/edev`, `/msg`)
- `-output string` (**Required**) - Output path for extracted WADL file
- `-v` (**Optional**) - Enable verbose output (shows all extracted resource paths)

## Common Resource Paths

Based on `RESOURCE_URI_PATHS.md`, here are the primary resource paths available for extraction:

| Path | Count | Description |
|------|-------|-------------|
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

## How It Works

1. **Load** - Reads and parses the input WADL file
2. **Filter** - Identifies all resources with `samplePath` attributes matching the given prefix
3. **Extract** - Preserves the raw XML including namespaced attributes
4. **Save** - Generates a new, valid WADL file with only the matched resources

## Output Format

The extracted WADL file maintains:
- XML declaration
- All namespace declarations
- WADL header elements (doc, grammars)
- Only the filtered resource elements
- Proper XML formatting and indentation

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| "failed to open WADL file" | Input file doesn't exist | Verify `-wadl` path is correct |
| "No resources found" | No resources match the prefix | Check `-path` parameter against `RESOURCE_URI_PATHS.md` |
| "failed to write WADL file" | Output directory doesn't exist or permission denied | Verify `-output` directory exists and is writable |

## Integration with Other Tools

The extracted WADL files can be used with other tools in the project:

```bash
# Extract and then use with scaffold-gen
./bin/wadl-extract \
  -wadl wadl/sep_wadl.xml \
  -path /dr \
  -output wadl/dr.wadl

go run ./cmd/scaffold-gen \
  -wadl wadl/dr.wadl \
  -output internal/demand-response
```

## Notes

- Extraction is based on exact prefix matching of the `wx:samplePath` attribute
- All resources matching the prefix are extracted (selective exclusion is not supported)
- The tool preserves all method definitions, parameters, and representations
- Namespaced attributes (e.g., `wx:samplePath`, `wx:mode`) are preserved in the output
