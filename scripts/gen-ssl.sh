#!/usr/bin/env bash
# gen-ssl.sh — Generate mTLS certificates for egot development
#
# Usage:
#   gen-ssl.sh refresh          Generate CA, server, and default client certs
#   gen-ssl.sh clients <N>      Generate N numbered client certs (requires existing CA)
#
# Output directory: ssl/ (relative to repo root)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SSL_DIR="$(dirname "$SCRIPT_DIR")/ssl"

# Certificate validity
CERT_DAYS=365
CA_DAYS=3650

# Server identity
SERVER_CN="egot.internal.com"
SERVER_SAN="DNS:egot.internal.com"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

usage() {
    echo "Usage:"
    echo "  $0 refresh          Generate CA + server + default client certs"
    echo "  $0 clients <N>      Generate N numbered client certs (client-0001 ... client-NNNN)"
    exit 1
}

require_ca() {
    if [ ! -f "$SSL_DIR/ca.key" ] || [ ! -f "$SSL_DIR/ca.crt" ]; then
        echo -e "${RED}Error: CA not found. Run '$0 refresh' first.${NC}"
        exit 1
    fi
}

gen_ca() {
    echo -e "${BLUE}Generating CA key and certificate...${NC}"
    openssl ecparam -genkey -name prime256v1 -noout -out "$SSL_DIR/ca.key"
    openssl req -new -x509 \
        -key "$SSL_DIR/ca.key" \
        -out "$SSL_DIR/ca.crt" \
        -days "$CA_DAYS" \
        -subj "/O=egot/CN=egot Root CA" \
        -addext "basicConstraints=critical,CA:TRUE" \
        -addext "keyUsage=critical,keyCertSign,cRLSign"
    echo -e "${GREEN}  ✓ ssl/ca.key + ssl/ca.crt${NC}"
}

gen_server() {
    echo -e "${BLUE}Generating server certificate for ${SERVER_CN}...${NC}"

    openssl ecparam -genkey -name prime256v1 -noout -out "$SSL_DIR/server.key"

    openssl req -new \
        -key "$SSL_DIR/server.key" \
        -out "$SSL_DIR/server.csr" \
        -subj "/CN=${SERVER_CN}"

    openssl x509 -req \
        -in "$SSL_DIR/server.csr" \
        -CA "$SSL_DIR/ca.crt" \
        -CAkey "$SSL_DIR/ca.key" \
        -CAcreateserial \
        -out "$SSL_DIR/server.crt" \
        -days "$CERT_DAYS" \
        -extfile <(printf "subjectAltName=%s\nkeyUsage=critical,digitalSignature,keyEncipherment\nextendedKeyUsage=serverAuth" "$SERVER_SAN")

    rm -f "$SSL_DIR/server.csr"

    # flowreservation and operator reference srv.crt / srv.key
    cp "$SSL_DIR/server.crt" "$SSL_DIR/srv.crt"
    cp "$SSL_DIR/server.key" "$SSL_DIR/srv.key"

    echo -e "${GREEN}  ✓ ssl/server.crt + ssl/server.key (and srv.crt / srv.key)${NC}"
}

gen_client() {
    local cn="$1"
    local cert_out="$2"
    local key_out="$3"

    openssl ecparam -genkey -name prime256v1 -noout -out "$key_out"

    openssl req -new \
        -key "$key_out" \
        -out "${cert_out%.crt}.csr" \
        -subj "/CN=${cn}"

    openssl x509 -req \
        -in "${cert_out%.crt}.csr" \
        -CA "$SSL_DIR/ca.crt" \
        -CAkey "$SSL_DIR/ca.key" \
        -CAcreateserial \
        -out "$cert_out" \
        -days "$CERT_DAYS" \
        -extfile <(printf "keyUsage=critical,digitalSignature\nextendedKeyUsage=clientAuth")

    rm -f "${cert_out%.crt}.csr"
}

cmd_refresh() {
    mkdir -p "$SSL_DIR"
    gen_ca
    gen_server

    echo -e "${BLUE}Generating default client certificate (user-test)...${NC}"
    gen_client "user-test" "$SSL_DIR/client.crt" "$SSL_DIR/client.key"
    echo -e "${GREEN}  ✓ ssl/client.crt + ssl/client.key${NC}"

    echo ""
    echo -e "${GREEN}✅ SSL certificates refreshed. Validity: ${CERT_DAYS} days.${NC}"
    echo -e "   CA valid for: ${CA_DAYS} days"
    openssl x509 -in "$SSL_DIR/server.crt" -noout -enddate
}

cmd_clients() {
    local n="${1:-}"
    if [ -z "$n" ] || ! [[ "$n" =~ ^[0-9]+$ ]] || [ "$n" -lt 1 ]; then
        echo -e "${RED}Error: clients requires a positive integer. e.g. $0 clients 5${NC}"
        exit 1
    fi

    require_ca

    echo -e "${BLUE}Generating ${n} numbered client certificate(s)...${NC}"
    for i in $(seq -f "%04g" 1 "$n"); do
        local cn="user-${i}"
        local cert="$SSL_DIR/client-${i}.crt"
        local key="$SSL_DIR/client-${i}.key"
        gen_client "$cn" "$cert" "$key"
        echo -e "${GREEN}  ✓ ssl/client-${i}.crt + ssl/client-${i}.key (CN=${cn})${NC}"
    done

    echo ""
    echo -e "${GREEN}✅ Generated ${n} client certificate(s).${NC}"
}

# ── Main ──────────────────────────────────────────────────────────────────────

case "${1:-}" in
    refresh)  cmd_refresh ;;
    clients)  cmd_clients "${2:-}" ;;
    *)        usage ;;
esac
