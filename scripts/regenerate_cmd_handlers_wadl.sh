#!/bin/bash
# regenerate_cmd_handlers_wadl.sh
# 
# This script regenerates cmd/<Service>/main.go files by:
# 1. Reading WADL specifications to extract resource paths and HTTP methods
# 2. Extracting handler methods from internal/<Service>/handler/handler.go
# 3. Mapping WADL methods to handlers using deterministic naming convention
# 4. Auto-generating cmd/<Service>/main.go with proper http.Handle() registration
#
# Usage:
#   ./scripts/regenerate_cmd_handlers_wadl.sh           # Regenerate all services
#   ./scripts/regenerate_cmd_handlers_wadl.sh Bill      # Regenerate one service

set -e

cd "$(dirname "$(dirname "$(readlink -f "$0")")")"

# Services with their port assignments and WADL file mappings
declare -A SERVICE_WADL=(
  [BRS]=brs
  [Bill]=bill
  [DCAP]=dcap
  [DERP]=derp
  [DR]=dr
  [EDevice]=edev
  [File]=file
  [MUP]=mup
  [Messaging]=msg
  [Notify]=ntfy
  [PPY]=ppy
  [SDevice]=sdev
  [TariffProfile]=tariff-profile
  [TimeOfUse]=time-of-use
  [UPT]=upt
  [rsps]=rsps
)

declare -A SERVICE_PORTS=(
  [BRS]=8010
  [Bill]=8011
  [DCAP]=8012
  [DERP]=8013
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
  [rsps]=8009
)

# Function to find WADL file for a service
find_wadl_file() {
  local service=$1
  local wadl_base=${SERVICE_WADL[$service]:-$service}
  
  # Try exact match first
  if [ -f "wadl/${wadl_base}.wadl" ]; then
    echo "wadl/${wadl_base}.wadl"
    return 0
  fi
  
  # Try lowercase
  if [ -f "wadl/$(echo $wadl_base | tr '[:upper:]' '[:lower:]').wadl" ]; then
    echo "wadl/$(echo $wadl_base | tr '[:upper:]' '[:lower:]').wadl"
    return 0
  fi
  
  return 1
}

# Function to extract route from handler method name
# Pattern: METHOD + ResourceName -> /(path)
# This is reverse-engineered from WADL or uses fallback naming
extract_route() {
  local method_name=$1
  local wadl_file=$2
  local http_method=""
  
  # Determine HTTP method from method name prefix
  if [[ $method_name =~ ^GET ]]; then
    http_method="GET"
  elif [[ $method_name =~ ^HEAD ]]; then
    http_method="HEAD"
  elif [[ $method_name =~ ^POST ]]; then
    http_method="POST"
  elif [[ $method_name =~ ^PUT ]]; then
    http_method="PUT"
  elif [[ $method_name =~ ^DELETE ]]; then
    http_method="DELETE"
  else
    echo "ERROR: Unknown method prefix in $method_name"
    return 1
  fi
  
  # Extract resource name from method (remove HTTP method prefix)
  local resource_name="${method_name:${#http_method}}"
  
  # Look up this method in WADL to find its path
  # WADL has: <method id="GETCustomerAccount" name="GET" />
  # associated with: <resource id="CustomerAccount" wx:samplePath="/bill/{id1}" />
  
  # First try: search WADL for exact method ID
  local path=$(grep -o "method id=\"$method_name\"" "$wadl_file" 2>/dev/null | head -1)
  
  if [ -n "$path" ]; then
    # Find the resource containing this method
    # Go backwards from method ID to find parent resource's samplePath
    local line_num=$(grep -n "method id=\"$method_name\"" "$wadl_file" 2>/dev/null | head -1 | cut -d: -f1)
    
    if [ -n "$line_num" ]; then
      # Search backwards for the opening <resource> tag with samplePath
      path=$(head -n $line_num "$wadl_file" | \
             tac | \
             grep -m1 'samplePath=' | \
             grep -o 'samplePath="[^"]*"' | \
             cut -d'"' -f2)
      
      if [ -n "$path" ]; then
        echo "$http_method|$path"
        return 0
      fi
    fi
  fi
  
  # Fallback: Not found in WADL
  echo "ERROR: Method $method_name not found in WADL"
  return 1
}

# Services to process
SERVICES=()
if [ $# -gt 0 ]; then
  SERVICES=("$@")
else
  # Only process services that need updating (not rsps which is already correct)
  for service in "${!SERVICE_WADL[@]}"; do
    if [ "$service" != "rsps" ]; then
      SERVICES+=("$service")
    fi
  done
fi

echo "Regenerating cmd/<Service>/main.go from WADL specifications..."
echo ""

for service in "${SERVICES[@]}"; do
  port=${SERVICE_PORTS[$service]:-8000}
  handler_file="internal/$service/handler/handler.go"
  cmd_file="cmd/$service/main.go"
  wadl_file=$(find_wadl_file "$service")
  
  if [ ! -f "$handler_file" ]; then
    echo "⚠ Skipping $service - $handler_file not found"
    continue
  fi
  
  if [ ! -f "$wadl_file" ]; then
    echo "⚠ Skipping $service - WADL file not found (tried: wadl/$(echo ${SERVICE_WADL[$service]} | tr '[:upper:]' '[:lower:]').wadl)"
    continue
  fi
  
  # Extract public handler methods
  mapfile -t methods < <(grep "^func (h \*Handler)" "$handler_file" | \
                         grep -v "^func (h \*Handler) [a-z]" | \
                         awk '{print $4}' | cut -d'(' -f1 | sort -u)
  
  if [ ${#methods[@]} -eq 0 ]; then
    echo "⚠ No handler methods found in $service"
    continue
  fi
  
  # Generate cmd/<Service>/main.go
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

const MAX_ENTITIES memory.Entity = 100

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.$(echo $service | sed 's/^./\U&/'),
TLSConfig: cfg,
}

repo := memory.NewRepository(MAX_ENTITIES)
repo.InitRepository("./ssl")

h := handler.NewHandler(repo)
EOFMAIN
  
  # Add handler registrations with WADL-extracted routes
  local route_count=0
  for method in "${methods[@]}"; do
    # Extract route from handler method name using WADL
    route_info=$(extract_route "$method" "$wadl_file" 2>/dev/null || echo "ERROR")
    
    if [[ $route_info == "ERROR" ]]; then
      # Fallback: couldn't extract from WADL, skip this method
      echo "⚠ Warning: Could not map handler $method to WADL resource in $service"
      continue
    fi
    
    # Split route_info into HTTP method and path
    http_method=$(echo "$route_info" | cut -d'|' -f1)
    path=$(echo "$route_info" | cut -d'|' -f2)
    
    # Generate http.Handle() call
    echo -e "\thttp.Handle(\"$http_method $path\", http.HandlerFunc(h.$method))" >> "$cmd_file"
    ((route_count++))
  done
  
  # Add server startup
  cat >> "$cmd_file" << 'EOFEND'

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
EOFEND
  
  echo "✓ $service: Generated $cmd_file with $route_count routes from WADL"
done

echo ""
echo "✓ Done! All cmd/<Service>/main.go files regenerated with WADL-driven routes"
echo ""
echo "Next: go build ./cmd/... to build all services"
