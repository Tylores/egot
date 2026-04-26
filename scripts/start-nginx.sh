#!/bin/bash

# Start nginx API gateway
# Usage: ./start-nginx.sh [--config path/to/nginx.yaml]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
NGINX_DIR="$REPO_ROOT/nginx"
LOGS_DIR="$REPO_ROOT/logs"
PIDS_DIR="$REPO_ROOT/.pids"

# Configuration
CONFIG_FILE="${1:-$NGINX_DIR/nginx.yaml}"
NGINX_CONF="$NGINX_DIR/nginx.conf"
NGINX_PID_FILE="$PIDS_DIR/nginx.pid"

# Create directories
mkdir -p "$LOGS_DIR" "$PIDS_DIR"

# Validate config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Configuration file not found: $CONFIG_FILE"
    exit 1
fi

# Generate nginx config if generate-nginx-config.sh exists
if [ -f "$SCRIPT_DIR/generate-nginx-config.sh" ]; then
    echo "📝 Generating nginx.conf from configuration..."
    bash "$SCRIPT_DIR/generate-nginx-config.sh" "$CONFIG_FILE" "$NGINX_CONF"
    if [ $? -ne 0 ]; then
        echo "❌ Failed to generate nginx.conf"
        exit 1
    fi
fi

# Validate nginx config
if ! command -v nginx &> /dev/null; then
    echo "❌ nginx not found. Please install nginx."
    exit 1
fi

if ! nginx -t -c "$(cd "$NGINX_DIR" && pwd)/nginx.conf" 2>&1 | grep -q "successful"; then
    echo "❌ nginx configuration test failed"
    nginx -t -c "$(cd "$NGINX_DIR" && pwd)/nginx.conf"
    exit 1
fi

# Check if already running
if [ -f "$NGINX_PID_FILE" ] && kill -0 "$(cat "$NGINX_PID_FILE")" 2>/dev/null; then
    echo "⚠️  nginx is already running (PID: $(cat "$NGINX_PID_FILE"))"
    exit 0
fi

# Start nginx
echo "▶ Starting nginx API gateway..."
nginx -c "$(cd "$NGINX_DIR" && pwd)/nginx.conf"

# Get nginx PID (wait a moment for it to start)
sleep 0.5
if pgrep -f "nginx: master" > /dev/null; then
    NGINX_MAIN_PID=$(pgrep -f "nginx: master" | head -1)
    echo "$NGINX_MAIN_PID" > "$NGINX_PID_FILE"
    echo "✅ nginx started (PID: $NGINX_MAIN_PID)"
    echo "📍 Listening on $(grep 'listen' "$CONFIG_FILE" | grep -oP '(?<=listen: ")[^"]*' || echo 'localhost'):$(grep 'port:' "$CONFIG_FILE" | grep -oP '(?<=port: )[0-9]+')"
else
    echo "❌ Failed to start nginx"
    exit 1
fi
