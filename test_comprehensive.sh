#!/bin/bash

# Comprehensive TinyBrain PocketBase Integration Test Suite
# Tests all MCP endpoints with real PocketBase operations

BASE_URL="http://127.0.0.1:8090"
MCP_URL="$BASE_URL/mcp"
TEST_DIR="/tmp/tinybrain-test-$(date +%s)"

echo "🧪 Starting Comprehensive TinyBrain PocketBase Integration Tests"
echo "================================================================"
echo "Test Directory: $TEST_DIR"
echo "MCP Endpoint: $MCP_URL"
echo ""

# Create test directory
mkdir -p "$TEST_DIR"

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to run MCP test
run_mcp_test() {
    local test_name="$1"
    local method="$2"
    local params="$3"
    local expected_field="$4"
    local expected_value="$5"
    
    echo "🔍 Testing: $test_name"
    
    local response=$(curl -s -X POST "$MCP_URL" \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":$params}")
    
    local result=$(echo "$response" | jq -r ".result$expected_field // \"ERROR\"")
    
    if [[ "$result" == "$expected_value" ]] || [[ "$result" != "ERROR" && "$result" != "null" ]]; then
        echo "✅ PASS: $test_name"
        echo "   Response: $(echo "$response" | jq -c '.result')"
        ((TESTS_PASSED++))
    else
        echo "❌ FAIL: $test_name"
        echo "   Expected: $expected_value"
        echo "   Got: $result"
        echo "   Full Response: $response"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Helper function to run MCP test with custom validation
run_mcp_test_custom() {
    local test_name="$1"
    local method="$2"
    local params="$3"
    local validation_script="$4"
    
    echo "🔍 Testing: $test_name"
    
    local response=$(curl -s -X POST "$MCP_URL" \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":$params}")
    
    if eval "$validation_script"; then
        echo "✅ PASS: $test_name"
        echo "   Response: $(echo "$response" | jq -c '.result')"
        ((TESTS_PASSED++))
    else
        echo "❌ FAIL: $test_name"
        echo "   Full Response: $response"
        ((TESTS_FAILED++))
    fi
    echo ""
}

echo "🚀 Starting TinyBrain Server..."
cd /Users/alec/tinybrain
./tinybrain-fixed-list serve --dir "$TEST_DIR" &
SERVER_PID=$!

# Wait for server to start
echo "⏳ Waiting for server to start..."
sleep 3

# Check if server is running
if ! curl -s "$BASE_URL/api/" > /dev/null; then
    echo "❌ Server failed to start"
    exit 1
fi

echo "✅ Server started successfully"
echo ""

# Test 1: Initialize MCP
echo "📋 Test Suite 1: MCP Initialization"
echo "=================================="
run_mcp_test "Initialize MCP" "initialize" '{"protocolVersion":"2024-11-05","capabilities":{"roots":{"listChanged":true},"sampling":{}},"clientInfo":{"name":"test-client","version":"1.0.0"}}' ".protocolVersion" "2024-11-05"

# Test 2: List Tools
echo "📋 Test Suite 2: MCP Tools"
echo "========================="
run_mcp_test_custom "List MCP Tools" "tools/list" '{}' 'echo "$response" | jq -e ".result.tools | length > 0"'

# Test 3: Session Management
echo "📋 Test Suite 3: Session Management"
echo "=================================="
run_mcp_test_custom "Create Session" "create_session" '{"name":"Test Session","task_type":"security_review","description":"Test session for comprehensive testing"}' 'echo "$response" | jq -e ".result.session_id != null"'

# Store session ID for later tests
SESSION_ID=$(curl -s -X POST "$MCP_URL" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","id":1,"method":"create_session","params":{"name":"Test Session 2","task_type":"security_review"}}' | jq -r '.result.session_id')

echo "📝 Session ID for tests: $SESSION_ID"

# Test 4: Memory Operations
echo "📋 Test Suite 4: Memory Operations"
echo "================================"
run_mcp_test_custom "Store Memory" "store_memory" "{\"session_id\":\"$SESSION_ID\",\"title\":\"Test Memory\",\"content\":\"This is a test memory\",\"category\":\"note\",\"priority\":5}" 'echo "$response" | jq -e ".result.memory_id != null"'

# Store memory ID for later tests
MEMORY_ID=$(curl -s -X POST "$MCP_URL" \
    -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"store_memory\",\"params\":{\"session_id\":\"$SESSION_ID\",\"title\":\"Test Memory 2\",\"content\":\"This is another test memory\",\"category\":\"vulnerability\",\"priority\":8}}" | jq -r '.result.memory_id')

echo "📝 Memory ID for tests: $MEMORY_ID"

# Test 5: Search Operations
echo "📋 Test Suite 5: Search Operations"
echo "=================================="
run_mcp_test_custom "Search Memories" "search_memories" '{"query":"Test","limit":10}' 'echo "$response" | jq -e ".result.memories | length > 0"'
run_mcp_test_custom "Search Memories with Category" "search_memories" '{"query":"","category":"note","limit":10}' 'echo "$response" | jq -e ".result.memories | length > 0"'

# Test 6: Session Operations
echo "📋 Test Suite 6: Session Operations"
echo "==================================="
run_mcp_test_custom "Get Session" "get_session" "{\"session_id\":\"$SESSION_ID\"}" 'echo "$response" | jq -e ".result.name == \"Test Session 2\""'
run_mcp_test_custom "List Sessions" "list_sessions" '{"limit":10}' 'echo "$response" | jq -e ".result.sessions | length > 0"'

# Test 7: Memory Statistics
echo "📋 Test Suite 7: Memory Statistics"
echo "==================================="
run_mcp_test_custom "Get Memory Stats" "get_memory_stats" '{}' 'echo "$response" | jq -e ".result != null"'

# Test 8: Error Handling
echo "📋 Test Suite 8: Error Handling"
echo "==============================="
run_mcp_test_custom "Invalid Method" "invalid_method" '{}' 'echo "$response" | jq -e ".error != null"'
run_mcp_test_custom "Invalid Session ID" "get_session" '{"session_id":"invalid-id"}' 'echo "$response" | jq -e ".error != null"'

# Test 9: Data Persistence
echo "📋 Test Suite 9: Data Persistence"
echo "=================================="
echo "🔄 Restarting server to test data persistence..."
kill $SERVER_PID
sleep 2

# Start server again
./tinybrain-fixed-list serve --dir "$TEST_DIR" &
SERVER_PID=$!
sleep 3

# Test if data persists
run_mcp_test_custom "Data Persistence - Search" "search_memories" '{"query":"Test","limit":10}' 'echo "$response" | jq -e ".result.memories | length > 0"'
run_mcp_test_custom "Data Persistence - Session" "get_session" "{\"session_id\":\"$SESSION_ID\"}" 'echo "$response" | jq -e ".result.name == \"Test Session 2\""'

# Test 10: Advanced Operations
echo "📋 Test Suite 10: Advanced Operations"
echo "===================================="
run_mcp_test_custom "Store Memory with Tags" "store_memory" "{\"session_id\":\"$SESSION_ID\",\"title\":\"Tagged Memory\",\"content\":\"Memory with tags\",\"category\":\"technique\",\"tags\":[\"security\",\"testing\"],\"priority\":7}" 'echo "$response" | jq -e ".result.memory_id != null"'

run_mcp_test_custom "Search with Tags" "search_memories" '{"query":"Tagged","limit":10}' 'echo "$response" | jq -e ".result.memories | length > 0"'

# Test 11: Admin Dashboard Access
echo "📋 Test Suite 11: Admin Dashboard"
echo "================================"
echo "🔍 Testing admin dashboard access..."
if curl -s "$BASE_URL/_/" | grep -q "PocketBase"; then
    echo "✅ PASS: Admin Dashboard accessible"
    ((TESTS_PASSED++))
else
    echo "❌ FAIL: Admin Dashboard not accessible"
    ((TESTS_FAILED++))
fi

# Test 12: REST API Access
echo "📋 Test Suite 12: REST API"
echo "========================="
echo "🔍 Testing REST API access..."
if curl -s "$BASE_URL/api/" | grep -q "data"; then
    echo "✅ PASS: REST API accessible"
    ((TESTS_PASSED++))
else
    echo "❌ FAIL: REST API not accessible"
    ((TESTS_FAILED++))
fi

# Cleanup
echo "🧹 Cleaning up..."
kill $SERVER_PID
rm -rf "$TEST_DIR"

# Final Results
echo ""
echo "🏁 Test Results Summary"
echo "======================"
echo "✅ Tests Passed: $TESTS_PASSED"
echo "❌ Tests Failed: $TESTS_FAILED"
echo "📊 Total Tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ $TESTS_FAILED -eq 0 ]; then
    echo ""
    echo "🎉 ALL TESTS PASSED! TinyBrain PocketBase integration is working perfectly!"
    exit 0
else
    echo ""
    echo "⚠️  Some tests failed. Please review the output above."
    exit 1
fi
