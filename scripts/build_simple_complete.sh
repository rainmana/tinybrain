#!/bin/bash
set -e

echo -e "\033[1;33m=== Building TinyBrain v2.0 Complete Server ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR/.." # Go up to the v2.0 root directory

# Define output binary name
OUTPUT_NAME="tinybrain-v2-complete"

# Build the server
go build -o "$OUTPUT_NAME" ./cmd/server/simple_complete_v2.go

if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Build successful: $OUTPUT_NAME\033[0m"
    echo ""
    echo "To run the complete server:"
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
    echo "  - create_session: Create security assessment sessions"
    echo "  - get_session: Retrieve sessions by ID"
    echo "  - list_sessions: List sessions with filtering"
    echo "  - store_memory: Store security findings and memories"
    echo "  - search_memories: Search memories within sessions"
    echo "  - create_relationship: Link memories together"
    echo "  - list_relationships: List relationships between memories"
    echo "  - create_context_snapshot: Save LLM context snapshots"
    echo "  - create_task_progress: Track assessment progress"
    echo ""
    echo "MCP Resources available:"
    echo "  - tinybrain://status: Current status and capabilities"
    echo ""
    echo "Features:"
    echo "  ✅ Complete session management"
    echo "  ✅ Memory storage with categories and priorities"
    echo "  ✅ Relationship tracking between memories"
    echo "  ✅ Context snapshot functionality"
    echo "  ✅ Task progress tracking"
    echo "  ✅ PocketBase database with embedded SQLite"
    echo "  ✅ MCP protocol support via mcp-go"
    echo "  ✅ Admin dashboard for data management"
    echo "  ✅ REST API for web integration"
    echo "  ✅ Real-time capabilities"
    echo ""
    echo "Database collections:"
    echo "  ✅ sessions - Security assessment sessions"
    echo "  ✅ memory_entries - Security findings and memories"
    echo "  ✅ relationships - Links between memories"
    echo "  ✅ context_snapshots - LLM context snapshots"
    echo "  ✅ task_progress - Assessment progress tracking"
else
    echo -e "\033[0;31m❌ Build failed\033[0m"
    exit 1
fi
