#!/bin/bash

# Test script for TinyBrain MCP Server
# This script sends JSON-RPC requests to test the MCP server functionality

echo "Testing TinyBrain MCP Server..."

# Function to send JSON-RPC request
send_request() {
    local method="$1"
    local params="$2"
    local id="$3"
    
    echo "Sending request: $method"
    echo "{\"jsonrpc\":\"2.0\",\"id\":$id,\"method\":\"$method\",\"params\":$params}" | ./bin/tinybrain
    echo ""
}

# Test 1: Initialize
echo "=== Test 1: Initialize ==="
send_request "initialize" "{}" "1"

# Test 2: List tools
echo "=== Test 2: List Tools ==="
send_request "list_tools" "{}" "2"

# Test 3: Create a session
echo "=== Test 3: Create Session ==="
send_request "call_tool" '{"name":"create_session","arguments":{"name":"Security Code Review","description":"Testing security vulnerabilities in web application","task_type":"security_review"}}' "3"

# Test 4: Store a memory entry
echo "=== Test 4: Store Memory ==="
send_request "call_tool" '{"name":"store_memory","arguments":{"session_id":"session_1","title":"SQL Injection in Login","content":"Found SQL injection vulnerability in login form. User input is not properly sanitized.","category":"vulnerability","priority":8,"tags":"[\"sql-injection\",\"authentication\",\"critical\"]"}}' "4"

# Test 5: Search memories
echo "=== Test 5: Search Memories ==="
send_request "call_tool" '{"name":"search_memories","arguments":{"session_id":"session_1","query":"SQL injection","search_type":"exact","limit":10}}' "5"

# Test 6: Get session info
echo "=== Test 6: Get Session ==="
send_request "call_tool" '{"name":"get_session","arguments":{"session_id":"session_1"}}' "6"

echo "Testing complete!"
