#!/bin/bash

# Build script for TinyBrain Intelligence Final Server
# This script builds ONLY the intelligence final server to avoid conflicts with other main packages

set -e

echo "Building TinyBrain Intelligence Final Server..."

# Build the intelligence final server
go build -o tinybrain-intelligence-final cmd/server/pocketbase_intelligence_final.go

echo "✅ Build successful: tinybrain-intelligence-final"
echo ""
echo "To run the server:"
echo "  ./tinybrain-intelligence-final"
echo ""
echo "Server will listen on: http://127.0.0.1:8090"
echo "Admin dashboard: http://127.0.0.1:8090/_/"
echo "REST API: http://127.0.0.1:8090/api/"
echo "MCP endpoint: http://127.0.0.1:8090/mcp"

