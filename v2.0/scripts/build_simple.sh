#!/bin/bash
set -e

echo -e "\033[1;33m=== Building TinyBrain v2.0 Simple Server ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR/.." # Go up to the v2.0 root directory

# Define output binary name
OUTPUT_NAME="tinybrain-v2-simple"

# Build the simple server
go build -o "$OUTPUT_NAME" ./cmd/server/simple_v2.go

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Build successful: $OUTPUT_NAME\033[0m"
    echo ""
    echo "To run the simple server:"
    echo "  ./$OUTPUT_NAME serve --dev"
    echo ""
    echo "Server will provide:"
    echo "  - PocketBase web server: http://127.0.0.1:8090"
    echo "  - Admin dashboard: http://127.0.0.1:8090/_/"
    echo "  - Health check: http://127.0.0.1:8090/health"
    echo "  - Hello endpoint: http://127.0.0.1:8090/hello"
    echo "  - STDIO MCP server for LLM integration"
    echo ""
    echo "MCP Tools available:"
    echo "  - get_status: Get TinyBrain v2.0 status"
    echo ""
    echo "MCP Resources available:"
    echo "  - tinybrain://status: Current status information"
else
    echo -e "\033[0;31m❌ Build failed\033[0m"
    exit 1
fi
