#!/bin/bash

# Simple MCP Test Script for TinyBrain
# This script tests the MCP server with individual requests

echo "🧠 TinyBrain Simple MCP Test"
echo "============================"
echo ""

# Test 1: Initialize
echo "🔧 Test 1: Initialize MCP Connection"
echo "-----------------------------------"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/tinybrain
echo ""

# Test 2: List Tools
echo "🔧 Test 2: List Available Tools"
echo "-------------------------------"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./bin/tinybrain
echo ""

# Test 3: Create Session
echo "🔧 Test 3: Create Security Session"
echo "---------------------------------"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_session","arguments":{"name":"Test Security Session","description":"Testing TinyBrain MCP server","task_type":"security_review"}}}' | ./bin/tinybrain
echo ""

# Test 4: Health Check
echo "🔧 Test 4: Health Check"
echo "----------------------"
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"health_check","arguments":{}}}' | ./bin/tinybrain
echo ""

# Test 5: Database Stats
echo "🔧 Test 5: Database Statistics"
echo "-----------------------------"
echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_database_stats","arguments":{}}}' | ./bin/tinybrain
echo ""

echo "✅ Simple MCP Test Complete!"
echo "============================"
echo ""
echo "Note: Each request starts a new server instance."
echo "In production, the server would run continuously."
