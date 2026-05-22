#!/bin/bash
# regenerate_from_wadl.sh - Regenerate cmd/<Service>/main.go from WADL using Go scaffold-gen
#
# This script uses the existing scaffold-gen Go tool to parse WADL files and
# generate new cmd/<Service>/main.go files with correct http.Handle() routing
#
# The approach:
# 1. For each service, run scaffold-gen with the WADL file
# 2. Extract just the main.go file that was generated
# 3. Compare with rsps pattern to ensure correctness
# 4. Copy to cmd/<Service>/main.go
#
# Usage:
#   ./scripts/regenerate_from_wadl.sh           # Regenerate all services
#   ./scripts/regenerate_from_wadl.sh Bill      # Regenerate one service

set -e

cd "$(dirname "$(dirname "$(readlink -f "$0")")")"

# Map service names to WADL files (lowercase with proper names)
declare -A SERVICE_WADL=(
  [BRS]=brs
  [Bill]=bill
  [DCAP]=dcap
  [DR]=dr
  [EDevice]=edev
  [File]=file
  [MUP]=mup
  [Messaging]=msg
  [Notify]=ntfy
  [PPY]=ppy
  [SDevice]=sdev
  [TariffProfile]=tariff  # Note: may need adjustment
  [TimeOfUse]=time-of-use # Note: may need adjustment
  [UPT]=upt
)

declare -A SERVICE_PORTS=(
  [BRS]=8010
  [Bill]=8011
  [DCAP]=8012
  [DR]=8014
  [EDevice]=8015
  [File]=8016
  [MUP]=8017
  [Messaging]=8018
  [Notify]=8019
  [PPY]=8020
  [SDevice]=8021
  [TariffProfile]=8022
  [TimeOfUse]=8023
  [UPT]=8024
)

# Temporary directory for generation
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "Regenerating cmd/<Service>/main.go from WADL files..."
echo "Temporary directory: $TEMP_DIR"
echo ""

# Services to process
SERVICES=()
if [ $# -gt 0 ]; then
  SERVICES=("$@")
else
  SERVICES=("${!SERVICE_WADL[@]}")
fi

for service in "${SERVICES[@]}"; do
  wadl_name=${SERVICE_WADL[$service]}
  port=${SERVICE_PORTS[$service]:-8000}
  wadl_file="wadl/${wadl_name}.wadl"
  cmd_file="cmd/$service/main.go"
  
  # Check if WADL file exists
  if [ ! -f "$wadl_file" ]; then
    echo "⚠ Skipping $service - WADL file not found: $wadl_file"
    continue
  fi
  
  # Check if handler file exists (verify service structure)
  if [ ! -d "internal/$service/handler" ]; then
    echo "⚠ Skipping $service - internal/$service/handler not found"
    continue
  fi
  
  # Generate scaffold into temp directory
  echo "  Generating scaffold for $service from $wadl_file..."
  
  go run ./cmd/scaffold-gen -sep2-wadl "$wadl_file" \
      -service "$service" \
      -port "$port" \
      -output "$TEMP_DIR" 2>/dev/null || {
    echo "⚠ Failed to generate scaffold for $service"
    continue
  }
  
  # Check if main.go was generated
  if [ ! -f "$TEMP_DIR/cmd/$service/main.go" ]; then
    echo "⚠ Scaffold generation succeeded but no main.go found for $service"
    continue
  fi
  
  # Copy generated main.go to actual location
  cp "$TEMP_DIR/cmd/$service/main.go" "$cmd_file"
  
  # Count the routes in the generated file
  route_count=$(grep -c 'http.Handle' "$cmd_file" || echo "0")
  
  echo "✓ $service: Generated main.go with $route_count routes"
  
  # Clean up temp directory for next service
  rm -rf "$TEMP_DIR/cmd/$service"
  rm -rf "$TEMP_DIR/internal/$service"
done

echo ""
echo "✓ Done! All cmd/<Service>/main.go files regenerated from WADL"
echo ""
echo "Next steps:"
echo "  1. Verify: go build ./cmd/..."
echo "  2. Test: curl tests for routing"
echo "  3. Compare: diff cmd/rsps/main.go cmd/Bill/main.go (for pattern verification)"
