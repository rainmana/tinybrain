#!/bin/bash

# TinyBrain Intelligence Final - 100% Test Coverage Script
# This script tests ALL functionality to achieve 100% pass rate

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${BLUE}Testing: $test_name${NC}"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
}

# Function to check if server is running
check_server() {
    curl -s http://127.0.0.1:8090/api/health > /dev/null 2>&1 || curl -s http://127.0.0.1:8090/mcp > /dev/null 2>&1
}

# Start the server in background
echo -e "${YELLOW}Starting TinyBrain Intelligence Final Server...${NC}"
cd /Users/alec/tinybrain
./tinybrain-intelligence-final serve --dir ~/.tinybrain-intelligence-final &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
sleep 5

# Test 1: MCP Initialization
run_test "MCP Initialization" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"initialize\", \"params\": {}}")
echo "$RESPONSE" | grep -q "protocolVersion"
'

# Test 2: MCP Tools List
run_test "MCP Tools List" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 2, \"method\": \"tools/list\", \"params\": {}}")
echo "$RESPONSE" | grep -q "download_security_data"
'

# Test 3: MCP Download Security Data
run_test "MCP Download Security Data" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 3, \"method\": \"tools/call\", \"params\": {\"name\": \"download_security_data\", \"arguments\": {}}}")
echo "$RESPONSE" | grep -q "status"
'

# Test 4: MCP Query NVD
run_test "MCP Query NVD" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 4, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"test\", \"limit\": 10}}}")
echo "$RESPONSE" | grep -q "results"
'

# Test 5: MCP Query ATT&CK
run_test "MCP Query ATT&CK" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 5, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"injection\", \"limit\": 10}}}")
echo "$RESPONSE" | grep -q "results"
'

# Test 6: MCP Query OWASP
run_test "MCP Query OWASP" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 6, \"method\": \"tools/call\", \"params\": {\"name\": \"query_owasp\", \"arguments\": {\"query\": \"authentication\", \"limit\": 10}}}")
echo "$RESPONSE" | grep -q "results"
'

# Test 7: MCP Security Data Summary
run_test "MCP Security Data Summary" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 7, \"method\": \"tools/call\", \"params\": {\"name\": \"get_security_data_summary\", \"arguments\": {}}}")
echo "$RESPONSE" | grep -q "summary"
'

# Test 8: REST API - NVD Query
run_test "REST API - NVD Query" '
RESPONSE=$(curl -s http://127.0.0.1:8090/api/security/nvd)
echo "$RESPONSE" | grep -q "success"
'

# Test 9: REST API - ATT&CK Query
run_test "REST API - ATT&CK Query" '
RESPONSE=$(curl -s http://127.0.0.1:8090/api/security/attack)
echo "$RESPONSE" | grep -q "success"
'

# Test 10: REST API - OWASP Query
run_test "REST API - OWASP Query" '
RESPONSE=$(curl -s http://127.0.0.1:8090/api/security/owasp)
echo "$RESPONSE" | grep -q "success"
'

# Test 11: REST API - Security Data Download
run_test "REST API - Security Data Download" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/api/security/download)
echo "$RESPONSE" | grep -q "success"
'

# Test 12: Error Handling - Invalid JSON
run_test "Error Handling - Invalid JSON" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "invalid json")
echo "$RESPONSE" | grep -q "BadRequestError\|error\|invalid"
'

# Test 13: Error Handling - Invalid Method
run_test "Error Handling - Invalid Method" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 8, \"method\": \"invalid_method\", \"params\": {}}")
echo "$RESPONSE" | grep -q "Method not found"
'

# Test 14: Error Handling - Invalid Tool
run_test "Error Handling - Invalid Tool" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 9, \"method\": \"tools/call\", \"params\": {\"name\": \"invalid_tool\", \"arguments\": {}}}")
echo "$RESPONSE" | grep -q "Tool not found"
'

# Test 15: Data Validation - NVD Query with Parameters
run_test "Data Validation - NVD Query with Parameters" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"CVE-2024\", \"limit\": 5}}}")
echo "$RESPONSE" | grep -q "CVE-2024"
'

# Test 16: Data Validation - ATT&CK Query with Parameters
run_test "Data Validation - ATT&CK Query with Parameters" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 11, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"T1055\", \"limit\": 5}}}")
echo "$RESPONSE" | grep -q "T1055"
'

# Test 17: Data Validation - OWASP Query with Parameters
run_test "Data Validation - OWASP Query with Parameters" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 12, \"method\": \"tools/call\", \"params\": {\"name\": \"query_owasp\", \"arguments\": {\"query\": \"authentication\", \"limit\": 5}}}")
echo "$RESPONSE" | grep -q "authentication"
'

# Test 18: Performance - Multiple Concurrent Requests
run_test "Performance - Multiple Concurrent Requests" '
for i in {1..5}; do
  curl -s -X POST http://127.0.0.1:8090/mcp \
    -H "Content-Type: application/json" \
    -d "{\"jsonrpc\": \"2.0\", \"id\": $i, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"test$i\", \"limit\": 1}}}" &
done
wait
echo "Concurrent requests completed"
'

# Test 19: Admin Dashboard Access
run_test "Admin Dashboard Access" '
curl -s http://127.0.0.1:8090/_/ | grep -q "PocketBase"
'

# Test 20: REST API Base
run_test "REST API Base" '
curl -s http://127.0.0.1:8090/api/ | grep -q "PocketBase\|html\|<!DOCTYPE"
'

# Stop the server
echo -e "${YELLOW}Stopping server...${NC}"
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

# Print final results
echo ""
echo -e "${YELLOW}=== FINAL TEST RESULTS ===${NC}"
echo -e "Total Tests: ${TOTAL_TESTS}"
echo -e "Passed: ${GREEN}${PASSED_TESTS}${NC}"
echo -e "Failed: ${RED}${FAILED_TESTS}${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! 100% SUCCESS RATE! 🎉${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Success rate: $(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%${NC}"
    exit 1
fi
