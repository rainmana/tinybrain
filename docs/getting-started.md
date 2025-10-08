---
layout: default
title: Getting Started
nav_order: 2
description: "Installation and setup guide for TinyBrain"
---

# Getting Started

This guide will help you install, configure, and start using TinyBrain for security-focused memory management.

## Table of Contents
- [Installation](#installation)
- [Basic Usage](#basic-usage)
- [MCP Client Configuration](#mcp-client-configuration)
- [Troubleshooting](#troubleshooting)

## Installation

TinyBrain supports three installation methods. Choose the one that works best for your environment.

### Method 1: Clone and Build Locally (Recommended)

This is the recommended installation method as it gives you full control over the build process.

```bash
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
make install
```

The binary will be installed to your `$GOPATH/bin` directory.

### Method 2: Install from Source

If the repository is public, you can install directly using Go:

```bash
go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Method 3: Build Binary Directly

Build the binary and manually install it:

```bash
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
go build -o tinybrain cmd/server/main.go
sudo mv tinybrain /usr/local/bin/
```

## Basic Usage

Once installed, you can start the TinyBrain MCP server:

```bash
# Start the server (uses ~/.tinybrain/memory.db by default)
tinybrain

# Or with custom database path
TINYBRAIN_DB_PATH=/path/to/your/memory.db tinybrain
```

The server will:
1. Initialize the SQLite database (if not exists)
2. Start the MCP server on standard input/output
3. Wait for MCP client connections

## MCP Client Configuration

TinyBrain works with any MCP-compatible client. Here's how to configure it with popular clients:

### Claude Desktop

Add to your Claude Desktop MCP configuration file:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": [],
      "env": {
        "TINYBRAIN_DB_PATH": "~/.tinybrain/memory.db"
      }
    }
  }
}
```

### Other MCP Clients

The configuration is similar for other MCP clients:
- **Command**: `tinybrain`
- **Args**: `[]` (empty array)
- **Environment**: Optional `TINYBRAIN_DB_PATH` variable

See the [Integrations](integrations/ai-assistants.md) section for detailed AI assistant configurations.

## Troubleshooting

### Authentication Errors with `go install`

If you encounter authentication errors when using `go install`:

```bash
# Use direct clone method instead
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
go build -o tinybrain cmd/server/main.go
```

### Private Repository Access

If the repository is private, ensure you have proper access configured:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### Go Module Proxy Issues

If you experience issues with the Go module proxy:

```bash
GOPROXY=direct go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Database Permissions

Ensure the database directory is writable:

```bash
mkdir -p ~/.tinybrain
chmod 755 ~/.tinybrain
```

### Check Database Health

Once running, you can verify the database is working correctly using the `health_check` tool:

```json
{
  "name": "health_check",
  "arguments": {}
}
```

This will return database status and connection information.

### Reset Database

If you need to reset your database:

```bash
make db-reset
```

Or manually:

```bash
rm -f ~/.tinybrain/memory.db*
```

## Next Steps

Now that TinyBrain is installed:

1. **Review Configuration** - See [Configuration Guide](configuration.md) for environment variables and database settings
2. **Understand the Architecture** - Read the [Architecture Guide](architecture.md) to understand the system design
3. **Explore the API** - Check out [API Examples](api/examples.md) for practical usage
4. **Follow Workflows** - See [Security Assessment Workflow](workflows/security-assessment.md) for real-world examples

## Verification

To verify your installation is working correctly:

1. Start the TinyBrain server
2. Connect with your MCP client
3. Try creating a test session:

```json
{
  "name": "create_session",
  "arguments": {
    "name": "Test Session",
    "description": "Testing TinyBrain installation",
    "task_type": "general"
  }
}
```

If this succeeds, your installation is complete and ready to use!
