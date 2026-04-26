#!/bin/bash

# Stop nginx API gateway gracefully
# Usage: ./stop-nginx.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
PIDS_DIR="$REPO_ROOT/.pids"
NGINX_PID_FILE="$PIDS_DIR/nginx.pid"

# Check if nginx PID file exists
if [ ! -f "$NGINX_PID_FILE" ]; then
    echo "⚠️  nginx PID file not found ($NGINX_PID_FILE)"
    
    # Try to find and kill nginx anyway
    if pgrep -f "nginx: master" > /dev/null; then
        echo "Found running nginx process, terminating..."
        pkill -f "nginx: master"
        sleep 1
        if pgrep -f "nginx: master" > /dev/null; then
            echo "Force killing nginx..."
            pkill -9 -f "nginx: master"
        fi
        echo "✅ nginx stopped"
    else
        echo "ℹ️  nginx is not running"
    fi
    exit 0
fi

NGINX_PID=$(cat "$NGINX_PID_FILE")

# Check if the process is still running
if ! kill -0 "$NGINX_PID" 2>/dev/null; then
    echo "⚠️  nginx PID $NGINX_PID is not running"
    rm -f "$NGINX_PID_FILE"
    exit 0
fi

# Graceful shutdown with SIGTERM
echo "Stopping nginx (PID: $NGINX_PID)..."
kill -TERM "$NGINX_PID"

# Wait for graceful shutdown (up to 10 seconds)
for i in {1..10}; do
    if ! kill -0 "$NGINX_PID" 2>/dev/null; then
        echo "✅ nginx stopped gracefully"
        rm -f "$NGINX_PID_FILE"
        exit 0
    fi
    sleep 1
done

# Force kill if still running
echo "⚠️  nginx did not stop gracefully, force killing..."
kill -9 "$NGINX_PID"
rm -f "$NGINX_PID_FILE"
echo "✅ nginx force stopped"
