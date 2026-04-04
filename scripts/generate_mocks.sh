#!/bin/bash

# Mock Generation Script for egot microservices
# This script helps generate mocks from interfaces using mockgen

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MOCKS_DIR="$PROJECT_ROOT/test/mocks"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if mockgen is installed
if ! command -v mockgen &> /dev/null; then
    echo -e "${YELLOW}mockgen not found. Installing...${NC}"
    go install github.com/golang/mock/mockgen@latest
fi

# Create mocks directory if it doesn't exist
mkdir -p "$MOCKS_DIR"

# Function to generate mock
generate_mock() {
    local source_file=$1
    local interface_name=$2
    local mock_filename="mock_${interface_name,,}.go"
    local destination="$MOCKS_DIR/$mock_filename"
    
    echo -e "${BLUE}Generating mock for $interface_name from $source_file${NC}"
    
    mockgen -source="$source_file" \
            -destination="$destination" \
            -package=mocks \
            "$interface_name"
    
    echo -e "${GREEN}✓ Mock generated: $destination${NC}"
}

# If no arguments, show usage
if [ $# -eq 0 ]; then
    echo -e "${BLUE}Mock Generation Helper${NC}"
    echo ""
    echo "Usage:"
    echo "  ./scripts/generate_mocks.sh <source_file> <InterfaceName>"
    echo ""
    echo "Examples:"
    echo "  ./scripts/generate_mocks.sh internal/Bill/repository/repository.go Repository"
    echo "  ./scripts/generate_mocks.sh internal/Messaging/handler/handler.go Handler"
    echo ""
    echo "To find interfaces in your project:"
    echo "  grep -r \"^type.*interface {\" internal/"
    echo ""
    exit 0
fi

if [ $# -lt 2 ]; then
    echo -e "${YELLOW}Error: Missing arguments${NC}"
    echo "Usage: ./scripts/generate_mocks.sh <source_file> <InterfaceName>"
    exit 1
fi

# Change to project root
cd "$PROJECT_ROOT"

# Generate the mock
generate_mock "$1" "$2"

echo -e "${GREEN}Done!${NC}"
echo ""
echo "Next steps:"
echo "  1. Review the generated mock in: $MOCKS_DIR"
echo "  2. Use it in your tests with: testhelpers.New(t)"
echo ""
