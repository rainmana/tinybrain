#!/bin/bash

# Test script for TinyBrain Security Knowledge Hub integration

echo "🧠 Testing TinyBrain Security Knowledge Hub Integration"
echo "=================================================="

# Function to send JSON-RPC request
send_request() {
    local method="$1"
    local params="$2"
    local id="$3"
    
    echo "📤 Sending request: $method"
    echo "📋 Params: $params"
    
    echo "{\"jsonrpc\": \"2.0\", \"id\": $id, \"method\": \"$method\", \"params\": $params}" | ./tinybrain
    echo ""
}

echo "1️⃣ Testing Security Data Summary..."
send_request "get_security_data_summary" "{}" 1

echo "2️⃣ Testing NVD Data Download (small test)..."
send_request "download_security_data" "{\"data_source\": \"nvd\"}" 2

echo "3️⃣ Testing ATT&CK Data Download..."
send_request "download_security_data" "{\"data_source\": \"attack\"}" 3

echo "4️⃣ Testing OWASP Data Download..."
send_request "download_security_data" "{\"data_source\": \"owasp\"}" 4

echo "5️⃣ Testing Security Data Summary After Downloads..."
send_request "get_security_data_summary" "{}" 5

echo "6️⃣ Testing NVD Query..."
send_request "query_nvd" "{\"query\": \"SQL injection\", \"limit\": 5}" 6

echo "7️⃣ Testing ATT&CK Query..."
send_request "query_attack" "{\"query\": \"process injection\", \"limit\": 5}" 7

echo "8️⃣ Testing OWASP Query..."
send_request "query_owasp" "{\"query\": \"authentication\", \"limit\": 5}" 8

echo "✅ Security Knowledge Hub Integration Test Complete!"
