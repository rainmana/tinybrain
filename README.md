# TinyBrain 🧠

**Security-Focused LLM Memory Storage MCP Server**

TinyBrain is a comprehensive memory storage system designed specifically for security-focused tasks like security code review, penetration testing, and exploit development. It provides an LLM with persistent, intelligent memory capabilities through the Model Context Protocol (MCP).

## 🎯 Key Features

### Security-Focused Design
- **Specialized Categories**: Vulnerability, exploit, payload, technique, tool, reference, context, hypothesis, evidence, recommendation
- **Priority & Confidence Tracking**: 0-10 priority levels and 0.0-1.0 confidence scores
- **Relationship Mapping**: Track dependencies, causes, mitigations, and exploit chains
- **Task Progress Tracking**: Multi-stage security task management

### Intelligent Memory Management
- **Context-Aware Storage**: Automatically categorizes and prioritizes information
- **Advanced Search**: Semantic, exact, fuzzy, tag-based, and relationship-based search
- **Access Tracking**: Monitors which memories are most relevant and frequently accessed
- **Context Summaries**: Provides relevant memory summaries for current tasks

### High Performance & Reliability
- **SQLite Backend**: Fast, reliable, local storage with full-text search
- **Optimized Queries**: Indexed searches and efficient relationship traversal
- **Transaction Safety**: ACID compliance for data integrity
- **Concurrent Access**: Thread-safe operations for multiple LLM interactions

### Developer Experience
- **Simple Installation**: `go install` or `go build`
- **Comprehensive Logging**: Detailed logging with structured output
- **Extensive Testing**: 90%+ test coverage with benchmarks
- **Docker Support**: Containerized deployment ready

## 🚀 Quick Start

### Installation

```bash
# Install from source
go install github.com/rainmana/tinybrain/cmd/server@latest

# Or build locally
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
make install
```

### Basic Usage

```bash
# Start the server (uses ~/.tinybrain/memory.db by default)
tinybrain

# Or with custom database path
TINYBRAIN_DB_PATH=/path/to/your/memory.db tinybrain
```

### MCP Client Configuration

Add to your MCP client configuration (e.g., Claude Desktop):

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

## 📚 API Reference

### Session Management

#### Create Session
```json
{
  "name": "create_session",
  "arguments": {
    "name": "Web App Security Review",
    "description": "Comprehensive security review of web application",
    "task_type": "security_review",
    "metadata": "{\"target\": \"myapp.com\", \"scope\": \"web-application\"}"
  }
}
```

**Task Types**: `security_review`, `penetration_test`, `exploit_dev`, `vulnerability_analysis`, `threat_modeling`, `incident_response`, `general`

#### List Sessions
```json
{
  "name": "list_sessions",
  "arguments": {
    "task_type": "security_review",
    "status": "active",
    "limit": 10
  }
}
```

### Memory Operations

#### Store Memory
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "SQL Injection in Login Form",
    "content": "Found SQL injection vulnerability in username parameter of login form. Payload: ' OR 1=1--",
    "category": "vulnerability",
    "content_type": "text",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"sql-injection\", \"authentication\", \"critical\"]",
    "source": "manual-testing"
  }
}
```

**Categories**: `finding`, `vulnerability`, `exploit`, `payload`, `technique`, `tool`, `reference`, `context`, `hypothesis`, `evidence`, `recommendation`, `note`

#### Search Memories
```json
{
  "name": "search_memories",
  "arguments": {
    "query": "SQL injection authentication",
    "session_id": "session_123",
    "search_type": "semantic",
    "categories": "[\"vulnerability\", \"exploit\"]",
    "min_priority": 7,
    "limit": 20
  }
}
```

**Search Types**: `semantic`, `exact`, `fuzzy`, `tag`, `category`, `relationship`

#### Get Related Memories
```json
{
  "name": "get_related_memories",
  "arguments": {
    "memory_id": "memory_456",
    "relationship_type": "exploits",
    "limit": 10
  }
}
```

### Relationship Management

#### Create Relationship
```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "memory_123",
    "target_memory_id": "memory_456",
    "relationship_type": "exploits",
    "strength": 0.8,
    "description": "SQL injection can be used to bypass authentication"
  }
}
```

**Relationship Types**: `depends_on`, `causes`, `mitigates`, `exploits`, `references`, `contradicts`, `supports`, `related_to`, `parent_of`, `child_of`

### Context Management

#### Get Context Summary
```json
{
  "name": "get_context_summary",
  "arguments": {
    "session_id": "session_123",
    "current_task": "Analyzing authentication vulnerabilities",
    "max_memories": 20
  }
}
```

## 🏗️ Architecture

### Database Schema

```sql
-- Core tables
sessions              -- LLM interaction sessions
memory_entries        -- Individual pieces of information
relationships         -- Links between memory entries
context_snapshots     -- Saved context states
search_history        -- Search query tracking
task_progress         -- Multi-stage task progress

-- Views and indexes
memory_entries_fts    -- Full-text search virtual table
memory_entries_with_session  -- Enhanced memory view
relationship_network  -- Relationship analysis view
```

### Key Design Principles

1. **Security-First**: All data structures and operations designed for security tasks
2. **Performance**: Optimized queries and indexes for fast retrieval
3. **Flexibility**: Extensible schema and relationship system
4. **Reliability**: ACID transactions and data integrity checks
5. **Usability**: Simple API with comprehensive documentation

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
make bench

# Run specific test
go test -v ./internal/database -run TestNewDatabase
```

## 🐳 Docker Deployment

```bash
# Build Docker image
make docker-build

# Run container
docker run --rm -it \
  -v ~/.tinybrain:/app/data \
  tinybrain:latest
```

## 📊 Performance

### Benchmarks

- **Memory Entry Creation**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Relationship Queries**: ~500 queries/second
- **Database Size**: ~1MB per 10,000 memory entries

### Optimization Features

- **Connection Pooling**: Optimized for SQLite single-writer model
- **Index Strategy**: Comprehensive indexing for all query patterns
- **Full-Text Search**: FTS5 virtual tables for semantic search
- **Caching**: Access pattern tracking for intelligent caching

## 🔧 Configuration

### Environment Variables

- `TINYBRAIN_DB_PATH`: Path to SQLite database (default: `~/.tinybrain/memory.db`)
- `TINYBRAIN_LOG_LEVEL`: Log level (debug, info, warn, error)

### Database Configuration

The SQLite database is configured with:
- WAL mode for better concurrency
- Foreign key constraints enabled
- Full-text search enabled
- Optimized pragma settings

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Development Setup

```bash
# Clone repository
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain

# Setup development environment
make dev-setup

# Run tests
make test

# Build
make build
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [mcp-go](https://github.com/mark3labs/mcp-go) for MCP server framework
- [go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite driver
- [charmbracelet/log](https://github.com/charmbracelet/log) for structured logging

## 📈 Roadmap

- [ ] HTTP transport support
- [ ] Memory compression and archiving
- [ ] Advanced analytics and insights
- [ ] Multi-user support with access controls
- [ ] Plugin system for custom memory types
- [ ] Integration with popular security tools
- [ ] Web dashboard for memory visualization

---

**TinyBrain** - Making LLM memory storage intelligent, fast, and security-focused. 🧠🔒
