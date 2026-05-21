#!/bin/bash

# Generate nginx.conf from nginx.yaml configuration and routes.go service map
# Usage: ./generate-nginx-config.sh [input_yaml] [output_conf]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

# Configuration parameters
INPUT_YAML="${1:-$REPO_ROOT/nginx/nginx.yaml}"
OUTPUT_CONF="${2:-$REPO_ROOT/nginx/nginx.conf}"
TEMPLATE="$REPO_ROOT/nginx/nginx.conf.template"
ROUTES_FILE="$REPO_ROOT/internal/routes/routes.go"

# Validate input files
if [ ! -f "$INPUT_YAML" ]; then
    echo "❌ Configuration file not found: $INPUT_YAML"
    exit 1
fi

if [ ! -f "$TEMPLATE" ]; then
    echo "❌ Template file not found: $TEMPLATE"
    exit 1
fi

if [ ! -f "$ROUTES_FILE" ]; then
    echo "❌ Routes file not found: $ROUTES_FILE"
    exit 1
fi

# Parse YAML configuration using grep and sed
parse_yaml() {
    local key="$1"
    local file="$2"
    grep "^  $key:" "$file" | sed "s/^  $key: ['\"]//g" | sed "s/['\"]$//g"
}

parse_nested_yaml() {
    local parent="$1"
    local key="$2"
    local file="$3"
    grep "    $key:" "$file" | sed "s/^    $key: ['\"]//g" | sed "s/['\"]$//g"
}

# Extract configuration values
LISTEN_ADDR=$(parse_nested_yaml "server" "listen" "$INPUT_YAML" || echo "127.0.0.1")
PORT=$(parse_nested_yaml "server" "port" "$INPUT_YAML" || echo "8443")
SERVER_NAME=$(parse_nested_yaml "server" "server_name" "$INPUT_YAML" || echo "localhost")
TLS_CERT=$(parse_nested_yaml "tls" "cert" "$INPUT_YAML" || echo "./ssl/server.crt")
TLS_KEY=$(parse_nested_yaml "tls" "key" "$INPUT_YAML" || echo "./ssl/server.key")
ACCESS_LOG=$(parse_nested_yaml "logging" "access_log" "$INPUT_YAML" || echo "./logs/nginx_access.log")
ERROR_LOG=$(parse_nested_yaml "logging" "error_log" "$INPUT_YAML" || echo "./logs/nginx_error.log")
ERROR_LOG_LEVEL=$(parse_nested_yaml "logging" "error_log_level" "$INPUT_YAML" || echo "warn")
PROXY_CONNECT_TIMEOUT=$(parse_nested_yaml "proxy" "connect_timeout" "$INPUT_YAML" || echo "10")
PROXY_READ_TIMEOUT=$(parse_nested_yaml "proxy" "read_timeout" "$INPUT_YAML" || echo "30")
PROXY_SEND_TIMEOUT=$(parse_nested_yaml "proxy" "send_timeout" "$INPUT_YAML" || echo "30")
PROXY_BUFFER_SIZE=$(parse_nested_yaml "proxy" "buffer_size" "$INPUT_YAML" || echo "4k")

# Extract service definitions from routes.go
# Services are defined as const declarations like: ServiceName = "host:port"
SERVICES=$(grep '^\s*[A-Z][a-zA-Z]* = ".*\.internal\.com:[0-9]\+"' "$ROUTES_FILE" | sed 's/^\s*//g' | sed 's/ = //g' | sed 's/"//g')

# Generate upstream blocks
UPSTREAM_BLOCKS=""
while IFS= read -r service_def; do
    if [ -z "$service_def" ]; then
        continue
    fi
    
    SERVICE_NAME=$(echo "$service_def" | cut -d' ' -f1)
    SERVICE_ADDR=$(echo "$service_def" | cut -d' ' -f2)
    SERVICE_HOST=$(echo "$SERVICE_ADDR" | cut -d':' -f1)
    SERVICE_PORT=$(echo "$SERVICE_ADDR" | cut -d':' -f2)
    
    UPSTREAM_BLOCKS="${UPSTREAM_BLOCKS}    upstream ${SERVICE_NAME,,}_backend {
        server $SERVICE_HOST:$SERVICE_PORT;
    }

"
done <<< "$SERVICES"

# Generate location blocks from serviceMap in routes.go
# Extract all path mappings from the serviceMap variable
LOCATION_BLOCKS=""

# Parse serviceMap entries (format: "/path": ServiceName,)
while IFS= read -r line; do
    # Extract path and service name
    ROUTE_PATH=$(echo "$line" | sed -n 's/.*"\(\/[^"]*\)": \([^,}]*\),.*/\1/p' | head -1)
    SERVICE=$(echo "$line" | sed -n 's/.*"\(\/[^"]*\)": \([^,}]*\),.*/\2/p' | head -1)
    
    if [ -n "$ROUTE_PATH" ] && [ -n "$SERVICE" ]; then
        # Convert path with {id} to nginx location pattern
        NGINX_PATH=$(echo "$ROUTE_PATH" | sed 's/{id[0-9]*}/~/g')
        
        LOCATION_BLOCKS="${LOCATION_BLOCKS}        location ~ ^${NGINX_PATH}(/.*)?$ {
            proxy_pass https://${SERVICE,,}_backend;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
            proxy_ssl_verify off;
        }

"
    fi
done < <(grep '^\s*"/.*":' "$ROUTES_FILE")

# If no location blocks were generated, add a catch-all that returns 404
if [ -z "$LOCATION_BLOCKS" ]; then
    LOCATION_BLOCKS="        # Default catch-all for unmapped paths
        location / {
            return 404 'No service mapping for this path';
        }"
fi

# Create output file from template
{
    sed "s|{{ LISTEN_ADDR }}|$LISTEN_ADDR|g" "$TEMPLATE" | \
    sed "s|{{ PORT }}|$PORT|g" | \
    sed "s|{{ SERVER_NAME }}|$SERVER_NAME|g" | \
    sed "s|{{ TLS_CERT }}|$TLS_CERT|g" | \
    sed "s|{{ TLS_KEY }}|$TLS_KEY|g" | \
    sed "s|{{ ACCESS_LOG }}|$ACCESS_LOG|g" | \
    sed "s|{{ ERROR_LOG }}|$ERROR_LOG|g" | \
    sed "s|{{ ERROR_LOG_LEVEL }}|$ERROR_LOG_LEVEL|g" | \
    sed "s|{{ PROXY_CONNECT_TIMEOUT }}|${PROXY_CONNECT_TIMEOUT}s|g" | \
    sed "s|{{ PROXY_READ_TIMEOUT }}|${PROXY_READ_TIMEOUT}s|g" | \
    sed "s|{{ PROXY_SEND_TIMEOUT }}|${PROXY_SEND_TIMEOUT}s|g" | \
    sed "s|{{ PROXY_BUFFER_SIZE }}|$PROXY_BUFFER_SIZE|g" | \
    sed "/{{ UPSTREAM_BLOCKS }}/r /dev/stdin" <<< "$UPSTREAM_BLOCKS" | \
    sed "/{{ UPSTREAM_BLOCKS }}/d" | \
    sed "/{{ LOCATION_BLOCKS }}/r /dev/stdin" <<< "$LOCATION_BLOCKS" | \
    sed "/{{ LOCATION_BLOCKS }}/d"
} > "$OUTPUT_CONF"

echo "✅ Generated nginx.conf: $OUTPUT_CONF"
echo "📍 Configuration: $INPUT_YAML"
echo "📋 Services mapped: $(echo "$SERVICES" | wc -l)"
