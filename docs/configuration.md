---
layout: default
title: Configuration
nav_order: 3
description: "Environment variables and database configuration for TinyBrain"
---

# Configuration

This guide covers all configuration options for TinyBrain, including environment variables, database settings, and Makefile helper targets.

## Table of Contents
- [Environment Variables](#environment-variables)
- [Database Configuration](#database-configuration)
- [Makefile Helper Targets](#makefile-helper-targets)

## Environment Variables

TinyBrain supports the following environment variables for configuration:

### TINYBRAIN_DB_PATH

**Description**: Path to the SQLite database file

**Default**: `~/.tinybrain/memory.db`

**Example**:
```bash
# Use custom database path
TINYBRAIN_DB_PATH=/var/lib/tinybrain/memory.db tinybrain

# Use project-local database
TINYBRAIN_DB_PATH=./project-memory.db tinybrain
```

### TINYBRAIN_LOG_LEVEL

**Description**: Log level for server output

**Values**: `debug`, `info`, `warn`, `error`

**Default**: `info`

**Example**:
```bash
# Enable debug logging
TINYBRAIN_LOG_LEVEL=debug tinybrain

# Reduce logging to errors only
TINYBRAIN_LOG_LEVEL=error tinybrain
```

## Database Configuration

TinyBrain uses SQLite as its storage backend with optimized settings for performance and reliability.

### SQLite Settings

The database is automatically configured with:

- **WAL Mode**: Write-Ahead Logging for better concurrency
  - Allows concurrent reads while writing
  - Improves performance for multiple client connections
  
- **Foreign Key Constraints**: Enabled for referential integrity
  - Ensures data consistency across relationships
  - Prevents orphaned memory entries
  
- **Full-Text Search**: FTS5 virtual tables enabled
  - Powers semantic search capabilities
  - Optimized for natural language queries
  
- **Optimized Pragma Settings**:
  - Connection pooling optimized for SQLite single-writer model
  - Comprehensive indexing for all query patterns
  - Access pattern tracking for intelligent caching

### Database Schema

The database includes the following tables:

**Core Tables:**
- `sessions` - LLM interaction sessions
- `memory_entries` - Individual pieces of information
- `relationships` - Links between memory entries
- `context_snapshots` - Saved context states
- `search_history` - Search query tracking
- `task_progress` - Multi-stage task progress
- `notifications` - Real-time alerts and notifications
- `cve_mappings` - CVE vulnerability mappings
- `risk_correlations` - Risk correlation data
- `compliance_mappings` - Compliance framework mappings

**Views:**
- `memory_entries_fts` - Full-text search virtual table
- `memory_entries_with_session` - Enhanced memory view with session data
- `relationship_network` - Relationship analysis view

### Database Maintenance

**Initialize Database Directory:**
```bash
make db-init
```

This creates the `~/.tinybrain` directory with appropriate permissions.

**Reset Database:**
```bash
make db-reset
```

This removes the database file and all its data. **Warning**: This operation is irreversible.

**Backup Database:**
```bash
# Simple file copy (ensure server is stopped)
cp ~/.tinybrain/memory.db ~/.tinybrain/memory.db.backup

# Or use SQLite backup command
sqlite3 ~/.tinybrain/memory.db ".backup ~/.tinybrain/memory.db.backup"
```

## Makefile Helper Targets

The TinyBrain repository includes a comprehensive Makefile with helper targets for development and operations.

### Development Targets

**dev-setup** - Set up development environment
```bash
make dev-setup
```
Downloads dependencies and tidies Go modules.

**build** - Build the binary
```bash
make build
```
Builds the `tinybrain` binary to `bin/tinybrain`.

**install** - Install to GOPATH/bin
```bash
make install
```
Installs the binary to your Go binary directory.

**run** - Run the server
```bash
make run
```
Runs the server with default settings.

**run-dev** - Run with development database
```bash
make run-dev
```
Runs the server using `./dev.db` as the database path.

### Database Targets

**db-init** - Initialize database directory
```bash
make db-init
```
Creates the `~/.tinybrain` directory.

**db-reset** - Reset database
```bash
make db-reset
```
Removes all database files from `~/.tinybrain`.

### Docker Targets

**docker-build** - Build Docker image
```bash
make docker-build
```
Builds a Docker image tagged with the current version.

**docker-run** - Run Docker container
```bash
make docker-run
```
Runs the TinyBrain server in a Docker container.

### Testing and Quality Targets

**test** - Run tests with coverage
```bash
make test
```
Runs all tests and generates coverage reports.

**lint** - Run linter
```bash
make lint
```
Runs `go vet` to check for code issues.

**security** - Run security checks
```bash
make security
```
Runs security vulnerability scanning on dependencies.

### Documentation Target

**docs** - Generate documentation
```bash
make docs
```
Generates API documentation from Go doc comments.

### Other Targets

**clean** - Clean build artifacts
```bash
make clean
```
Removes compiled binaries, test coverage files, and development databases.

**help** - Show help
```bash
make help
```
Displays all available Makefile targets with descriptions.

## Configuration Best Practices

1. **Use Absolute Paths**: When specifying `TINYBRAIN_DB_PATH`, use absolute paths to avoid confusion
2. **Backup Regularly**: Implement regular database backups for production use
3. **Monitor Database Size**: Track database growth using `get_database_stats` tool
4. **Log Appropriately**: Use `info` level for production, `debug` for troubleshooting
5. **Secure the Database**: Ensure proper file permissions on the database file and directory

## Next Steps

- Learn about the [Architecture](architecture.md) to understand the system design
- Explore [API Examples](api/examples.md) for practical usage patterns
- Review [Advanced Features](advanced-features.md) for context snapshots and task tracking
