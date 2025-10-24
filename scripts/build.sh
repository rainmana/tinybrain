#!/bin/bash

# TinyBrain v2.0 Build Script
set -e

echo -e "\033[1;33m=== Building TinyBrain v2.0 ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR/.."

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo -e "\033[0;31m❌ go.mod not found. Are you in the right directory?\033[0m"
    exit 1
fi

# Download dependencies
echo -e "\033[0;34m📦 Downloading dependencies...\033[0m"
go mod tidy

# Run tests
echo -e "\033[0;34m🧪 Running tests...\033[0m"
go test ./test/... -v

# Build the server
echo -e "\033[0;34m🔨 Building server...\033[0m"
go build -o tinybrain-v2 ./cmd/server/main.go

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Build successful: tinybrain-v2\033[0m"
    echo ""
    echo "To run the server:"
    echo "  ./tinybrain-v2"
    echo ""
    echo "Environment variables:"
    echo "  TINYBRAIN_DATA_DIR - Data directory (default: ./data)"
    echo ""
    echo "The server will:"
    echo "  - Initialize PocketBase database"
    echo "  - Create required collections"
    echo "  - Start MCP server on stdio"
else
    echo -e "\033[0;31m❌ Build failed\033[0m"
    exit 1
fi
