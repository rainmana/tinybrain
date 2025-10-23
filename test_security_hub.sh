#!/bin/bash

# Test script for TinyBrain Security Knowledge Hub
# This script tests the new security data querying capabilities

echo "Testing TinyBrain Security Knowledge Hub..."

# Function to send JSON-RPC request
send_request() {
    local method="$1"
    local params="$2"
    local id="$3"
    
    echo "Sending request: $method"
    echo "{\"jsonrpc\":\"2.0\",\"id\":$id,\"method\":\"tools/call\",\"params\":{\"name\":\"$method\",\"arguments\":$params}}" | ./tinybrain
    echo ""
}

# Test 1: Get security data summary
echo "=== Test 1: Get Security Data Summary ==="
send_request "get_security_data_summary" "{}" "1"

# Test 2: Query NVD (placeholder)
echo "=== Test 2: Query NVD ==="
send_request "query_nvd" '{"cwe_id":"CWE-89","limit":5}' "2"

# Test 3: Query ATT&CK (placeholder)
echo "=== Test 3: Query ATT&CK ==="
send_request "query_attack" '{"tactic":"persistence","limit":5}' "3"

# Test 4: Query OWASP (placeholder)
echo "=== Test 4: Query OWASP ==="
send_request "query_owasp" '{"category":"authentication","limit":5}' "4"

# Test 5: Download security data (placeholder)
echo "=== Test 5: Download Security Data ==="
send_request "download_security_data" '{"data_source":"nvd","force_update":false}' "5"

echo "Security Knowledge Hub testing complete!"
