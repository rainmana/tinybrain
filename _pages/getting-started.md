---
layout: default
title: Getting Started
permalink: /getting-started/
---

# Getting Started

## Installation {#installation}

### From Source (Recommended)
```bash
go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Docker
```bash
docker pull rainmana/tinybrain:latest
docker run -p 8090:8090 rainmana/tinybrain
```

### Pre-built Binaries
Download from [Releases](https://github.com/rainmana/tinybrain/releases)

## PocketBase Integration {#pocketbase-integration}

TinyBrain now uses PocketBase as its backend, providing:

- **Single Binary**: Everything in one executable with zero configuration
- **Admin Dashboard**: Web interface at http://127.0.0.1:8090/_/ for data management
- **REST API**: Full REST API at http://127.0.0.1:8090/api/ for external integrations
- **Real-time Updates**: Server-sent events for live memory updates
- **Data Persistence**: All data persists across server restarts
- **Comprehensive Testing**: 17/17 tests passing with full functionality verification

## Basic Usage {#basic-usage}

### 1. Start the Server
```bash
# Start with default data directory
tinybrain serve --dir ~/.tinybrain

# Or with custom data directory
tinybrain serve --dir /path/to/your/data
```

### 2. Access Admin Dashboard
```bash
# Open admin dashboard in browser
open http://127.0.0.1:8090/_/

# Or access REST API
curl http://127.0.0.1:8090/api/
```

### 3. MCP Integration
Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": ["serve", "--dir", "~/.tinybrain"],
      "env": {
        "TINYBRAIN_DB_PATH": "~/.tinybrain/memory.db"
      }
    }
  }
}
```

### 4. Create a Session (MCP)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "mcp_tinybrain-mcp-server_create_session",
    "arguments": {
      "name": "Security Assessment",
      "task_type": "penetration_test",
      "intelligence_type": "osint"
    }
  }
}
```

### 5. Store Memory (MCP)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "mcp_tinybrain-mcp-server_store_memory",
    "arguments": {
      "session_id": "session-id",
      "title": "OSINT Finding",
      "content": "Social media analysis reveals...",
      "category": "intelligence",
      "priority": 7,
      "confidence": 0.8,
      "tags": "[\"osint\", \"social-media\", \"reconnaissance\"]"
    }
  }
}
```

## Configuration {#configuration}

### PocketBase Configuration
TinyBrain uses PocketBase with zero configuration by default:

```bash
# Default configuration
tinybrain serve --dir ~/.tinybrain

# Custom data directory
tinybrain serve --dir /path/to/your/data

# Custom port (if needed)
tinybrain serve --dir ~/.tinybrain --port 8090
```

### Environment Variables
```bash
export TINYBRAIN_DB_PATH="~/.tinybrain/memory.db"
export TINYBRAIN_HOST="localhost"
export TINYBRAIN_PORT="8090"
```

### Admin Dashboard Configuration
Access the admin dashboard at http://127.0.0.1:8090/_/ to:
- View and manage collections
- Monitor system performance
- Configure user access (future feature)
- View real-time logs and metrics

### REST API Configuration
The REST API is automatically available at http://127.0.0.1:8090/api/ with:
- Full CRUD operations for all collections
- Real-time subscriptions
- Authentication support (future feature)
- Rate limiting and security controls

## Next Steps

- [Core Features](core-features/) - Learn about memory management and sessions
- [Intelligence & Reconnaissance](intelligence/) - Explore intelligence gathering capabilities
- [API Reference](api-reference/) - Complete API documentation
