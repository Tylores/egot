#!/bin/bash
# regenerate_cmd_handlers.sh
# 
# This script regenerates cmd/<Service>/main.go files by:
# 1. Analyzing the existing internal/<Service>/handler/handler.go files
# 2. Extracting all public handler methods (GET*, POST*, PUT*, DELETE*, HEAD*)
# 3. Auto-generating cmd/<Service>/main.go with proper handler registration
#
# Usage:
#   ./scripts/regenerate_cmd_handlers.sh           # Regenerate all services
#   ./scripts/regenerate_cmd_handlers.sh Bill      # Regenerate one service

set -e

cd "$(dirname "$(dirname "$(readlink -f "$0")")")"

# Services with their port assignments
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

# Services to process
SERVICES=()
if [ $# -gt 0 ]; then
  SERVICES=("$@")
else
  SERVICES=("${!SERVICE_PORTS[@]}")
fi

echo "Regenerating cmd/<Service>/main.go from handler definitions..."
echo ""

for service in "${SERVICES[@]}"; do
  port=${SERVICE_PORTS[$service]:-8000}
  handler_file="internal/$service/handler/handler.go"
  cmd_file="cmd/$service/main.go"
  
  if [ ! -f "$handler_file" ]; then
    echo "⚠ Skipping $service - $handler_file not found"
    continue
  fi
  
  # Extract public handler methods (exclude private methods starting with lowercase)
  methods=$(grep "^func (h \*Handler)" "$handler_file" | \
            grep -v "^func (h \*Handler) [a-z]" | \
            awk '{print $4}' | cut -d'(' -f1 | sort -u)
  
  if [ -z "$methods" ]; then
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
)

func main() {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := http.Server{
		Addr:      "egot.internal.com:$port",
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()
	h := handler.NewHandler(repo)
	
	// Register all handler methods from internal/$service/handler/handler.go
EOFMAIN
  
  # Add handler registrations
  while IFS= read -r method; do
    if [ -n "$method" ]; then
      echo -e "\thttp.HandleFunc(\"/\", h.$method)" >> "$cmd_file"
    fi
  done <<< "$methods"
  
  # Add server startup
  cat >> "$cmd_file" << 'EOFEND'

	err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
EOFEND
  
  method_count=$(echo "$methods" | wc -l)
  echo "✓ $service: Generated $cmd_file with $method_count handlers"
done

echo ""
echo "✓ Done! All cmd/<Service>/main.go files regenerated with handlers"
echo ""
echo "Next: go build ./cmd/... to build all services"
