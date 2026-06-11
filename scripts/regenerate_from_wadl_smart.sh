#!/bin/bash
# regenerate_from_wadl_smart.sh - Regenerate cmd/<Service>/main.go from WADL
#
# This version generates a proper main.go that matches the existing
# repository interface of each service
#
# Usage:
#   ./scripts/regenerate_from_wadl_smart.sh           # Regenerate all services
#   ./scripts/regenerate_from_wadl_smart.sh Bill      # Regenerate one service

set -e

cd "$(dirname "$(dirname "$(readlink -f "$0")")")"

# Map service names to WADL files
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
)

# Temporary directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "Regenerating cmd/<Service>/main.go from WADL with smart repository mapping..."
echo ""

SERVICES=()
if [ $# -gt 0 ]; then
  SERVICES=("$@")
else
  SERVICES=("${!SERVICE_WADL[@]}")
fi

for service in "${SERVICES[@]}"; do
  wadl_name=${SERVICE_WADL[$service]:-$service}
  port=${SERVICE_PORTS[$service]:-8000}
  wadl_file="wadl/${wadl_name}.wadl"
  cmd_file="cmd/$service/main.go"
  
  if [ ! -f "$wadl_file" ]; then
    echo "⚠ Skipping $service - WADL file not found: $wadl_file"
    continue
  fi
  
  if [ ! -d "internal/$service/handler" ]; then
    echo "⚠ Skipping $service - internal/$service/handler not found"
    continue
  fi
  
  echo "  Generating scaffold for $service from $wadl_file..."
  
  # Generate to temp directory
  go run ./cmd/scaffold-gen -sep2-wadl "$wadl_file" \
      -service "$service" \
      -port "$port" \
      -output "$TEMP_DIR" 2>/dev/null || {
    echo "⚠ Failed to generate scaffold for $service"
    continue
  }
  
  # Extract generated main.go
  temp_main="$TEMP_DIR/cmd/$service/main.go"
  if [ ! -f "$temp_main" ]; then
    echo "⚠ Scaffold generation succeeded but no main.go found for $service"
    continue
  fi
  
  # Extract routes from generated main.go
  # These are the important parts - we'll rebuild main.go to fit existing repo
  routes=$(grep 'http.Handle' "$temp_main" || echo "")
  route_count=$(echo "$routes" | grep -c 'http.Handle' || echo "0")
  
  if [ "$route_count" -eq 0 ]; then
    echo "⚠ No routes found in generated main.go for $service"
    continue
  fi
  
  # Create custom main.go for this service
  # This matches the existing repository interface
  cat > "$cmd_file" << EOFMAIN
package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/$service/handler"
"github.com/Tylores/egot/internal/$service/repository/memory"
"github.com/Tylores/egot/internal/routes"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.$(echo $service | sed 's/.*/\U&/'),
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)
EOFMAIN
  
  # Add extracted routes
  echo "$routes" >> "$cmd_file"
  
  # Add server startup
  cat >> "$cmd_file" << 'EOFEND'

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
EOFEND
  
  echo "✓ $service: Generated main.go with $route_count routes"
  
  # Clean up temp directory
  rm -rf "$TEMP_DIR/cmd/$service"
  rm -rf "$TEMP_DIR/internal/$service"
done

echo ""
echo "✓ Done! All cmd/<Service>/main.go files regenerated from WADL"
echo ""
echo "Next: go build ./cmd/..."
