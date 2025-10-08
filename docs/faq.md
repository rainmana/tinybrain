---
layout: default
title: FAQ
nav_order: 19
description: "Frequently asked questions about TinyBrain"
---

# Frequently Asked Questions

Common questions and answers about TinyBrain.

## General Questions

### What is TinyBrain?

TinyBrain is a security-focused LLM memory storage system that provides persistent, intelligent memory capabilities through the Model Context Protocol (MCP). It's designed specifically for security tasks like code review, penetration testing, and exploit development.

### Why use TinyBrain?

TinyBrain solves the problem of LLM context limitations by:
- Providing persistent memory across conversations
- Organizing security findings with proper categorization
- Tracking relationships between vulnerabilities and exploits
- Enabling context preservation for long-running assessments
- Supporting multi-stage security workflows

### Is TinyBrain open source?

Yes, TinyBrain is open source software released under the MIT License.

## Installation and Setup

### Where is the data stored?

By default, TinyBrain stores data in a SQLite database at:
```
~/.tinybrain/memory.db
```

You can customize this location using the `TINYBRAIN_DB_PATH` environment variable.

### How do I reset the database?

Reset the database using the Makefile target:

```bash
make db-reset
```

Or manually remove the database files:

```bash
rm -f ~/.tinybrain/memory.db*
```

### Can I use a different database location for each project?

Yes! Set the `TINYBRAIN_DB_PATH` environment variable when starting the server:

```bash
TINYBRAIN_DB_PATH=./project-memory.db tinybrain
```

Or configure it in your MCP client:

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

## Features and Capabilities

### How many MCP tools does TinyBrain provide?

TinyBrain provides 40 MCP tools organized into 8 categories:
- Core Memory Operations (8 tools)
- Session & Task Management (6 tools)
- Advanced Memory Features (8 tools)
- Security Templates & Batch Operations (6 tools)
- Memory Lifecycle & Cleanup (4 tools)
- AI-Enhanced Search (3 tools)
- Real-Time Notifications (4 tools)
- System Monitoring (1 tool)

### Is semantic search fully implemented?

Semantic search infrastructure is in place, but embedding generation currently uses a hash-based placeholder implementation as noted in the repository code comments. This provides a foundation for future AI model integration (OpenAI, Cohere, or local models).

### What is the test coverage?

The repository README states "90%+ test coverage." You can verify current coverage by running:

```bash
go test -cover ./...
```

### Can multiple users access the same database?

Currently, TinyBrain is designed for single-user access. Multi-user support with access controls is planned for a future release (see [Roadmap](roadmap.md)).

## Usage Questions

### How do I check the system health?

Use the `health_check` tool:

```json
{
  "name": "health_check",
  "arguments": {}
}
```

This returns database status and connection information.

### How do I get database statistics?

Use the `get_database_stats` tool:

```json
{
  "name": "get_database_stats",
  "arguments": {}
}
```

This returns counts by category, priority distribution, and access patterns.

### Can I export my session data?

Yes, use the `export_session_data` tool:

```json
{
  "name": "export_session_data",
  "arguments": {
    "session_id": "session_123",
    "include_relationships": true
  }
}
```

### How do I clean up old memories?

TinyBrain provides three cleanup operations:

**Age-based cleanup**:
```json
{
  "name": "cleanup_old_memories",
  "arguments": {
    "max_age_days": 90,
    "dry_run": true
  }
}
```

**Priority-based cleanup**:
```json
{
  "name": "cleanup_low_priority_memories",
  "arguments": {
    "max_priority": 3,
    "dry_run": true
  }
}
```

**Access-based cleanup**:
```json
{
  "name": "cleanup_unused_memories",
  "arguments": {
    "max_unused_days": 30,
    "dry_run": true
  }
}
```

Use `dry_run: true` to preview what would be deleted without actually removing anything.

## Performance Questions

### What are the performance benchmarks?

Current benchmark figures from repository testing:

- **Memory Entry Creation**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Relationship Queries**: ~500 queries/second
- **Database Size**: ~1MB per 10,000 memory entries

These are current benchmarks and may vary based on hardware and usage patterns.

### How is the database optimized?

TinyBrain uses several optimization techniques:
- WAL mode for better concurrency
- Comprehensive indexing for all query patterns
- FTS5 virtual tables for full-text search
- Access pattern tracking for intelligent caching
- Connection pooling optimized for SQLite

### Can I run TinyBrain in production?

Yes, TinyBrain is production-ready for single-user deployments. The SQLite backend is reliable and ACID-compliant. For multi-user scenarios, wait for the planned multi-user support release.

## Security Questions

### Does TinyBrain support authentication?

Currently, TinyBrain uses MCP protocol over stdio and does not include authentication. For network deployment, use appropriate network-level security (VPN, SSH tunnels, etc.). Authentication is planned for the HTTP transport feature.

### Is data encrypted?

Database files are not encrypted by default. You can use:
- File system encryption (e.g., encrypted volumes)
- Full disk encryption
- Database-level encryption (SQLCipher integration planned)

### Can I restrict access to certain sessions?

Access control is not currently implemented. This is planned for the multi-user support feature (see [Roadmap](roadmap.md)).

## Integration Questions

### Which AI assistants work with TinyBrain?

TinyBrain works with any MCP-compatible AI assistant. Specific configurations are available for:
- Cursor (`.cursorrules`)
- Cline (`.clinerules`)
- Roo (`.roo-mode`)
- Claude Desktop
- Any other MCP-compatible client

See the [Integrations](integrations/ai-assistants.md) section for details.

### Can I use TinyBrain with non-security tasks?

Yes! While TinyBrain is optimized for security assessments, it can be used for any task requiring persistent memory. The categorization system is flexible and can adapt to different domains.

### Does TinyBrain integrate with security tools?

Integration with popular security tools (Burp Suite, Metasploit, etc.) is planned but not yet implemented. See the [Roadmap](roadmap.md) for details.

## Troubleshooting

### TinyBrain won't start

Check:
1. Is TinyBrain installed? Run `which tinybrain`
2. Is the database directory writable? Check `~/.tinybrain` permissions
3. Are there any error messages in the logs?

### Database corruption

If the database becomes corrupted:
1. Stop the TinyBrain server
2. Backup the database: `cp ~/.tinybrain/memory.db ~/.tinybrain/memory.db.backup`
3. Try SQLite recovery: `sqlite3 ~/.tinybrain/memory.db ".recover" > recovered.sql`
4. If recovery fails, reset the database: `make db-reset`

### Search not returning results

Check:
1. Are you searching in the correct session?
2. Is the query matching the memory content?
3. Try different search types (semantic, exact, fuzzy)
4. Verify memories exist: `get_database_stats`

### Performance is slow

Optimization tips:
1. Clean up old/unused memories regularly
2. Avoid excessive relationships (1000s per memory)
3. Use appropriate search limits
4. Monitor database size and vacuum if needed
5. Consider database optimization: `sqlite3 ~/.tinybrain/memory.db "VACUUM;"`

## Development Questions

### How do I contribute?

See the [Contributing Guide](contributing.md) for detailed instructions on:
- Development setup
- Coding standards
- Testing requirements
- Pull request process

### Where can I report bugs?

Report bugs on the [GitHub Issues](https://github.com/rainmana/tinybrain/issues) page.

### Where can I request features?

Feature requests can be submitted as GitHub Issues or discussed in GitHub Discussions.

### How do I run the tests?

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run specific test
go test -v ./internal/database -run TestNewDatabase

# Run benchmarks
make bench
```

## Still Have Questions?

If your question isn't answered here:

1. **Check the Documentation**: Review the [complete documentation](index.md)
2. **Search Issues**: Look for similar questions in [GitHub Issues](https://github.com/rainmana/tinybrain/issues)
3. **Ask the Community**: Post in GitHub Discussions
4. **Open an Issue**: Create a new issue with your question

## Next Steps

- Review [Getting Started](getting-started.md) for installation instructions
- See [API Examples](api/examples.md) for usage patterns
- Check [Workflows](workflows/security-assessment.md) for real-world examples
