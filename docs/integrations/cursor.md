---
layout: default
title: Cursor Setup
nav_order: 16
parent: Integrations
description: "Cursor AI assistant setup guide for TinyBrain"
---

# Cursor Setup Guide

This guide provides step-by-step instructions for setting up Cursor with TinyBrain for security assessments.

## Table of Contents
- [Quick Setup](#quick-setup)
- [Configuration Steps](#configuration-steps)
- [What This Gives You](#what-this-gives-you)
- [Advanced Configuration](#advanced-configuration)
- [Troubleshooting](#troubleshooting)

## Quick Setup

Get started with Cursor and TinyBrain in 5 minutes.

### Step 1: Copy the Rules File

Copy the main rules file to your project:

```bash
# Copy to your project root
cp .cursorrules /path/to/your/project/

# Or copy to your global Cursor config
cp .cursorrules ~/.cursor/rules/
```

### Step 2: Customize for Your Needs (Optional)

Create a user-specific configuration:

```bash
# Copy the user template
cp .cursorrules.user-template .cursorrules.user

# Edit the user file to add your custom MCP servers and preferences
nano .cursorrules.user
```

### Step 3: Configure MCP Servers in Cursor

Add TinyBrain and optional MCP servers to your Cursor configuration:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": []
    },
    "clear-thought": {
      "command": "npx",
      "args": ["@clear-thought/mcp-server"]
    },
    "stochastic-thinking": {
      "command": "npx",
      "args": ["@stochastic-thinking/mcp-server"]
    }
  }
}
```

## Configuration Steps

### 1. Install TinyBrain

First, ensure TinyBrain is installed:

```bash
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
make install
```

Verify installation:

```bash
tinybrain --version
```

### 2. Configure Cursor MCP Settings

Open Cursor settings and add the MCP server configuration.

**Required Configuration**:
```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": []
    }
  }
}
```

**Optional: Add Custom Database Path**:
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

### 3. Copy Rules Files

Copy the Cursor rules to your project:

```bash
# Main rules file (required)
cp .cursorrules /path/to/your/project/

# User customization (optional)
cp .cursorrules.user-template .cursorrules.user
```

### 4. Test Configuration

Test that Cursor can communicate with TinyBrain:

```
Ask Cursor: "Create a test TinyBrain session named 'Configuration Test'"
```

If successful, you'll see a session created response.

## What This Gives You

### Enforced TinyBrain Usage

- ✅ **Mandatory session creation** for all security tasks
- ✅ **Required memory storage** with proper categorization
- ✅ **Relationship tracking** between findings
- ✅ **Progress monitoring** through assessment phases
- ✅ **Quality assurance** checklists

### Prevents Common LLM Issues

- 🚫 **No more lost findings** - everything stored in TinyBrain
- 🚫 **No more inconsistent categorization** - enforced standards
- 🚫 **No more forgotten relationships** - mandatory linking
- 🚫 **No more incomplete reports** - comprehensive validation
- 🚫 **No more LLM loops** - structured reasoning tools

### Customizable for Your Workflow

- 🔧 **Add your MCP servers** (clear-thought, stochastic-thinking, etc.)
- 🔧 **Customize security templates** for your domain
- 🔧 **Set your quality standards** and validation requirements
- 🔧 **Configure your assessment phases** and methodology

## Advanced Configuration

### For Claude 3.7+ Users

**Recommended MCP Servers**:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain"
    },
    "clear-thought": {
      "command": "npx",
      "args": ["@clear-thought/mcp-server"]
    },
    "stochastic-thinking": {
      "command": "npx",
      "args": ["@stochastic-thinking/mcp-server"]
    },
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem"]
    }
  }
}
```

### Custom Database Locations

**Project-Specific Database**:
```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "env": {
        "TINYBRAIN_DB_PATH": "./project-memory.db"
      }
    }
  }
}
```

**Team Shared Database** (network filesystem):
```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "env": {
        "TINYBRAIN_DB_PATH": "/shared/team/security-assessments.db"
      }
    }
  }
}
```

### Debug Logging

Enable debug logging for troubleshooting:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "env": {
        "TINYBRAIN_LOG_LEVEL": "debug"
      }
    }
  }
}
```

## Troubleshooting

### TinyBrain Not Found

**Problem**: Cursor cannot find the `tinybrain` command

**Solution**:
1. Verify TinyBrain is in your PATH:
   ```bash
   which tinybrain
   ```
2. If not found, use absolute path in config:
   ```json
   {
     "command": "/full/path/to/tinybrain"
   }
   ```

### Connection Errors

**Problem**: Cursor cannot connect to TinyBrain

**Solution**:
1. Restart Cursor after configuration changes
2. Check TinyBrain is running:
   ```bash
   ps aux | grep tinybrain
   ```
3. Review Cursor MCP logs for error messages

### Database Permission Issues

**Problem**: Cannot write to database

**Solution**:
1. Check database directory permissions:
   ```bash
   ls -la ~/.tinybrain/
   ```
2. Ensure directory is writable:
   ```bash
   chmod 755 ~/.tinybrain
   ```

### Rules Not Applied

**Problem**: Cursor doesn't follow TinyBrain rules

**Solution**:
1. Verify `.cursorrules` file is in project root
2. Check file permissions:
   ```bash
   ls -la .cursorrules
   ```
3. Restart Cursor to reload rules

### Health Check

Verify TinyBrain is working correctly:

```
Ask Cursor: "Run a health check on TinyBrain"
```

Expected response should include database status and connection info.

## Usage Examples

### Start Security Assessment

```
"I'm starting a security code review of the authentication module. Please:
1. Create a TinyBrain session
2. Initialize task progress tracking
3. Begin analyzing the code for vulnerabilities"
```

### Store Findings

```
"I found a SQL injection vulnerability in the login endpoint. Please store this in TinyBrain with:
- Category: vulnerability
- Priority: 9
- Confidence: 0.95
- Tags: sql-injection, authentication, critical"
```

### Generate Report

```
"Generate a comprehensive security assessment report from the TinyBrain session data, including:
1. Executive summary
2. Vulnerability list with priorities
3. Recommended remediation steps"
```

## Best Practices

1. **Always Create Sessions**: Start every assessment with a session
2. **Use Consistent Tags**: Apply standard security tags
3. **Track Progress**: Update task progress regularly
4. **Create Relationships**: Link related findings
5. **Regular Snapshots**: Save context at key milestones
6. **Generate Reports**: Export findings for documentation

## Next Steps

- Review [AI Assistant Integration](ai-assistants.md) for multi-assistant workflows
- See [Workflows](../workflows/security-assessment.md) for assessment patterns
- Check [API Examples](../api/examples.md) for TinyBrain tool usage
