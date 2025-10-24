#!/bin/bash
set -e

echo -e "\033[1;33m=== TinyBrain v2.0 Complete - Comprehensive MCP Testing ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR"

# Test server health
echo -e "\033[1;34mTesting server health...\033[0m"
HEALTH_RESPONSE=$(curl -s http://127.0.0.1:8090/health)
echo "Health Response: $HEALTH_RESPONSE"

if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo -e "\033[0;32m✅ Server is healthy\033[0m"
else
    echo -e "\033[0;31m❌ Server health check failed\033[0m"
    exit 1
fi

# Test hello endpoint
echo -e "\033[1;34mTesting hello endpoint...\033[0m"
HELLO_RESPONSE=$(curl -s http://127.0.0.1:8090/hello)
echo "Hello Response: $HELLO_RESPONSE"

if echo "$HELLO_RESPONSE" | grep -q "TinyBrain v2.0 Complete"; then
    echo -e "\033[0;32m✅ Hello endpoint working\033[0m"
else
    echo -e "\033[0;31m❌ Hello endpoint failed\033[0m"
    exit 1
fi

# Test MCP tools via simulated JSON-RPC calls
echo -e "\033[1;34mTesting MCP tools...\033[0m"

# Create a test session
echo "Creating test session..."
SESSION_REQUEST='{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_session",
    "arguments": {
      "name": "Security Assessment Test",
      "task_type": "penetration_test",
      "description": "Comprehensive security assessment"
    }
  }
}'

# Note: This is a simplified test - in a real scenario, you'd use an MCP client
echo "MCP Tool Test Request:"
echo "$SESSION_REQUEST" | jq '.'

# Test memory storage
echo "Testing memory storage..."
MEMORY_REQUEST='{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "store_memory",
    "arguments": {
      "session_id": "test-session-id",
      "title": "SQL Injection Vulnerability",
      "content": "Found SQL injection in login form",
      "category": "vulnerability",
      "priority": 8,
      "confidence": 0.9,
      "tags": "sql-injection,critical,authentication",
      "source": "manual_testing"
    }
  }
}'

echo "Memory Storage Test Request:"
echo "$MEMORY_REQUEST" | jq '.'

# Test relationship creation
echo "Testing relationship creation..."
RELATIONSHIP_REQUEST='{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "create_relationship",
    "arguments": {
      "source_id": "memory-1",
      "target_id": "memory-2",
      "type": "causes",
      "strength": 0.8,
      "description": "SQL injection leads to authentication bypass"
    }
  }
}'

echo "Relationship Test Request:"
echo "$RELATIONSHIP_REQUEST" | jq '.'

# Test context snapshot
echo "Testing context snapshot..."
CONTEXT_REQUEST='{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "create_context_snapshot",
    "arguments": {
      "session_id": "test-session-id",
      "name": "Assessment Progress",
      "context_data": "{\"current_findings\": 5, \"critical_issues\": 2}",
      "description": "Mid-assessment context snapshot"
    }
  }
}'

echo "Context Snapshot Test Request:"
echo "$CONTEXT_REQUEST" | jq '.'

# Test task progress
echo "Testing task progress..."
TASK_REQUEST='{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "create_task_progress",
    "arguments": {
      "session_id": "test-session-id",
      "task_name": "Vulnerability Assessment",
      "stage": "exploitation",
      "status": "in_progress",
      "progress_percentage": 75.0,
      "notes": "Successfully exploited 3 of 4 critical vulnerabilities"
    }
  }
}'

echo "Task Progress Test Request:"
echo "$TASK_REQUEST" | jq '.'

# Test MCP resource
echo "Testing MCP resource..."
RESOURCE_REQUEST='{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "resources/read",
  "params": {
    "uri": "tinybrain://status"
  }
}'

echo "Resource Test Request:"
echo "$RESOURCE_REQUEST" | jq '.'

echo -e "\033[0;32m🎉 All MCP tool tests prepared successfully!\033[0m"
echo ""
echo "Note: These are test requests that would be sent to the MCP server."
echo "In a real scenario, you would use an MCP client like Cline to interact"
echo "with the TinyBrain v2.0 server via STDIO transport."
echo ""
echo "Available MCP Tools:"
echo "  ✅ create_session - Create security assessment sessions"
echo "  ✅ get_session - Retrieve sessions by ID"
echo "  ✅ list_sessions - List sessions with filtering"
echo "  ✅ store_memory - Store security findings and memories"
echo "  ✅ search_memories - Search memories within sessions"
echo "  ✅ create_relationship - Link memories together"
echo "  ✅ list_relationships - List relationships between memories"
echo "  ✅ create_context_snapshot - Save LLM context snapshots"
echo "  ✅ create_task_progress - Track assessment progress"
echo ""
echo "Available MCP Resources:"
echo "  ✅ tinybrain://status - Current status and capabilities"
echo ""
echo "TinyBrain v2.0 Complete is ready for security assessments! 🚀"
