#!/bin/bash
set -e

echo -e "\033[1;33m=== Testing TinyBrain v2.0 Complete Implementation ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR/.." # Go up to the v2.0 root directory

# Run unit tests
echo -e "\033[1;34mRunning unit tests...\033[0m"
go test ./test/... -v

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Unit tests passed\033[0m"
else
    echo -e "\033[0;31m❌ Unit tests failed\033[0m"
    exit 1
fi

# Run integration tests
echo -e "\033[1;34mRunning integration tests...\033[0m"
go test ./test/integration_test.go -v

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Integration tests passed\033[0m"
else
    echo -e "\033[0;31m❌ Integration tests failed\033[0m"
    exit 1
fi

# Test the complete server build
echo -e "\033[1;34mTesting complete server build...\033[0m"
./scripts/build_complete.sh

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Complete server build successful\033[0m"
else
    echo -e "\033[0;31m❌ Complete server build failed\033[0m"
    exit 1
fi

# Test server startup (briefly)
echo -e "\033[1;34mTesting server startup...\033[0m"
timeout 10s ./tinybrain-v2-complete serve --dev &
SERVER_PID=$!

# Wait a moment for server to start
sleep 3

# Test health endpoint
echo -e "\033[1;34mTesting health endpoint...\033[0m"
if curl -s http://127.0.0.1:8090/health | grep -q "healthy"; then
    echo -e "\033[0;32m✅ Health endpoint working\033[0m"
else
    echo -e "\033[0;31m❌ Health endpoint failed\033[0m"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

# Test MCP endpoint
echo -e "\033[1;34mTesting MCP endpoint...\033[0m"
if curl -s http://127.0.0.1:8090/mcp | grep -q "MCP endpoint"; then
    echo -e "\033[0;32m✅ MCP endpoint working\033[0m"
else
    echo -e "\033[0;31m❌ MCP endpoint failed\033[0m"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

# Clean up
kill $SERVER_PID 2>/dev/null || true

echo -e "\033[0;32m🎉 All TinyBrain v2.0 tests passed!\033[0m"
echo ""
echo "TinyBrain v2.0 is ready with:"
echo "  ✅ PocketBase database with all collections"
echo "  ✅ MCP-Go server with all security tools"
echo "  ✅ Complete memory storage system"
echo "  ✅ Session management"
echo "  ✅ Relationship tracking"
echo "  ✅ Context snapshots"
echo "  ✅ Task progress tracking"
echo "  ✅ Web server with admin UI"
echo "  ✅ Health monitoring"
echo ""
echo "To start the server:"
echo "  ./tinybrain-v2-complete serve --dev"
echo ""
echo "To use with LLM:"
echo "  Add to your MCP configuration:"
echo "  {"
echo "    \"name\": \"tinybrain-v2\","
echo "    \"command\": \"$(pwd)/tinybrain-v2-complete\","
echo "    \"args\": [\"serve\", \"--dev\"]"
echo "  }"
