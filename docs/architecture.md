---
layout: default
title: Architecture
nav_order: 4
description: "TinyBrain system architecture and database schema"
---

# Architecture

This document describes TinyBrain's architecture, database schema, design principles, and server initialization flow.

## Table of Contents
- [Database Schema](#database-schema)
- [Key Design Principles](#key-design-principles)
- [Server Initialization Flow](#server-initialization-flow)
- [Performance Considerations](#performance-considerations)

## Database Schema

TinyBrain uses SQLite with an optimized schema designed for security-focused memory management.

### Core Tables

#### sessions
Tracks LLM interaction sessions for organizing security assessments.

**Fields:**
- `id` - Unique session identifier (TEXT PRIMARY KEY)
- `name` - Session name (TEXT NOT NULL)
- `description` - Session description (TEXT)
- `task_type` - Type of security task (TEXT NOT NULL)
  - Valid values: `security_review`, `penetration_test`, `exploit_dev`, `vulnerability_analysis`, `threat_modeling`, `incident_response`, `general`
- `status` - Current session status (TEXT NOT NULL)
  - Valid values: `active`, `paused`, `completed`, `archived`
- `created_at` - Creation timestamp (DATETIME)
- `updated_at` - Last update timestamp (DATETIME)
- `metadata` - JSON metadata for session-specific data (TEXT)

#### memory_entries
Stores individual pieces of security information.

**Fields:**
- `id` - Unique memory entry identifier (TEXT PRIMARY KEY)
- `session_id` - Associated session ID (TEXT NOT NULL, FOREIGN KEY)
- `title` - Memory title/summary (TEXT NOT NULL)
- `content` - Memory content (TEXT NOT NULL)
- `content_type` - Type of content (TEXT NOT NULL)
  - Valid values: `text`, `code`, `json`, `yaml`, `markdown`, `binary_ref`
- `category` - Security category (TEXT NOT NULL)
  - Valid values: `finding`, `vulnerability`, `exploit`, `payload`, `technique`, `tool`, `reference`, `context`, `hypothesis`, `evidence`, `recommendation`, `note`
- `priority` - Priority level 0-10 (INTEGER, 0=low, 10=critical)
- `confidence` - Confidence score 0.0-1.0 (REAL)
- `tags` - JSON array of tags (TEXT)
- `source` - Information source (TEXT)
- `created_at` - Creation timestamp (DATETIME)
- `updated_at` - Last update timestamp (DATETIME)
- `accessed_at` - Last access timestamp (DATETIME)
- `access_count` - Number of times accessed (INTEGER)

#### relationships
Links related memory entries for relationship mapping.

**Fields:**
- `id` - Unique relationship identifier (TEXT PRIMARY KEY)
- `source_entry_id` - Source memory ID (TEXT NOT NULL, FOREIGN KEY)
- `target_entry_id` - Target memory ID (TEXT NOT NULL, FOREIGN KEY)
- `relationship_type` - Type of relationship (TEXT NOT NULL)
  - Valid values: `depends_on`, `causes`, `mitigates`, `exploits`, `references`, `contradicts`, `supports`, `related_to`, `parent_of`, `child_of`
- `strength` - Relationship strength 0.0-1.0 (REAL)
- `description` - Relationship description (TEXT)
- `created_at` - Creation timestamp (DATETIME)

#### context_snapshots
Stores saved context states for security assessments.

**Fields:**
- `id` - Unique snapshot identifier (TEXT PRIMARY KEY)
- `session_id` - Associated session ID (TEXT NOT NULL, FOREIGN KEY)
- `name` - Snapshot name (TEXT NOT NULL)
- `description` - Snapshot description (TEXT)
- `context_data` - JSON context data (TEXT NOT NULL)
- `memory_summary` - Generated memory summary (TEXT)
- `created_at` - Creation timestamp (DATETIME)

#### search_history
Tracks LLM search queries for pattern analysis.

**Fields:**
- `id` - Unique search identifier (TEXT PRIMARY KEY)
- `session_id` - Associated session ID (TEXT NOT NULL, FOREIGN KEY)
- `query` - Search query text (TEXT NOT NULL)
- `search_type` - Type of search (TEXT NOT NULL)
  - Valid values: `semantic`, `exact`, `fuzzy`, `tag`, `category`, `relationship`
- `results_count` - Number of results returned (INTEGER)
- `created_at` - Creation timestamp (DATETIME)

#### task_progress
Tracks multi-stage security task progress.

**Fields:**
- `id` - Unique task identifier (TEXT PRIMARY KEY)
- `session_id` - Associated session ID (TEXT NOT NULL, FOREIGN KEY)
- `task_name` - Task name (TEXT NOT NULL)
- `stage` - Current stage (TEXT NOT NULL)
- `status` - Task status (TEXT NOT NULL)
  - Valid values: `pending`, `in_progress`, `completed`, `failed`, `blocked`
- `progress_percentage` - Completion percentage 0-100 (INTEGER)
- `notes` - Progress notes (TEXT)
- `started_at` - Task start timestamp (DATETIME, nullable)
- `completed_at` - Task completion timestamp (DATETIME, nullable)
- `created_at` - Creation timestamp (DATETIME)
- `updated_at` - Last update timestamp (DATETIME)

#### notifications
Real-time alerts and notifications.

**Fields:**
- `id` - Unique notification identifier (TEXT PRIMARY KEY)
- `session_id` - Associated session ID (TEXT, FOREIGN KEY, nullable)
- `notification_type` - Type of notification (TEXT NOT NULL)
- `priority` - Notification priority (INTEGER)
- `message` - Notification message (TEXT NOT NULL)
- `metadata` - JSON metadata (TEXT)
- `read` - Read status (INTEGER, 0 or 1)
- `created_at` - Creation timestamp (DATETIME)

#### cve_mappings
CVE vulnerability mappings.

**Fields:**
- `id` - Unique mapping identifier (TEXT PRIMARY KEY)
- `memory_id` - Associated memory ID (TEXT, FOREIGN KEY)
- `cve_id` - CVE identifier (TEXT NOT NULL)
- `created_at` - Creation timestamp (DATETIME)

#### risk_correlations
Risk correlation data.

**Fields:**
- `id` - Unique correlation identifier (TEXT PRIMARY KEY)
- `source_memory_id` - Source memory ID (TEXT, FOREIGN KEY)
- `target_memory_id` - Target memory ID (TEXT, FOREIGN KEY)
- `correlation_type` - Type of correlation (TEXT NOT NULL)
- `risk_score` - Risk score (REAL)
- `created_at` - Creation timestamp (DATETIME)

#### compliance_mappings
Compliance framework mappings.

**Fields:**
- `id` - Unique mapping identifier (TEXT PRIMARY KEY)
- `memory_id` - Associated memory ID (TEXT, FOREIGN KEY)
- `framework` - Compliance framework (TEXT NOT NULL)
- `control_id` - Control identifier (TEXT NOT NULL)
- `created_at` - Creation timestamp (DATETIME)

### Views

#### memory_entries_fts
Full-text search virtual table for semantic search capabilities.

Provides FTS5-based indexing of memory content for natural language queries.

#### memory_entries_with_session
Enhanced memory view with session data joined.

Combines memory entries with their associated session information for efficient querying.

#### relationship_network
Relationship analysis view for traversing memory relationships.

Provides optimized queries for exploring the relationship graph.

## Key Design Principles

TinyBrain's architecture is built on five core principles:

### 1. Security-First
All data structures and operations are designed specifically for security tasks:
- Specialized categories for vulnerabilities, exploits, and techniques
- Priority and confidence tracking for risk assessment
- Relationship mapping for exploit chains and dependencies

### 2. Performance
Optimized queries and indexes ensure fast retrieval:
- Comprehensive indexing for all query patterns
- Connection pooling optimized for SQLite single-writer model
- Full-text search with FTS5 virtual tables
- Access pattern tracking for intelligent caching

### 3. Flexibility
Extensible schema and relationship system:
- JSON metadata fields for custom data
- Multiple content types (text, code, JSON, YAML, etc.)
- Flexible relationship types
- Tag-based organization

### 4. Reliability
ACID transactions and data integrity checks:
- Foreign key constraints enabled
- CASCADE delete for referential integrity
- WAL mode for better concurrency
- Transaction safety for all operations

### 5. Usability
Simple API with comprehensive documentation:
- 40 MCP tools covering all operations
- Clear categorization and search strategies
- Intuitive relationship types
- Extensive examples and workflows

## Server Initialization Flow

TinyBrain follows a structured initialization process:

### 1. Logger Initialization
```go
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportCaller:    true,
    ReportTimestamp: true,
    TimeFormat:      time.Kitchen,
    Prefix:          "TinyBrain 🧠 ",
    Level:           log.InfoLevel,
})
```

### 2. Database Path Resolution
```go
dbPath := os.Getenv("TINYBRAIN_DB_PATH")
if dbPath == "" {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        logger.Fatal("Failed to get user home directory", "error", err)
    }
    dbPath = filepath.Join(homeDir, ".tinybrain", "memory.db")
}
```

### 3. Database Initialization
```go
db, err := database.NewDatabase(dbPath, logger)
if err != nil {
    logger.Fatal("Failed to initialize database", "error", err)
}
defer db.Close()
```

The database initialization:
- Creates the database directory if needed
- Opens SQLite with optimized settings (WAL mode, foreign keys, etc.)
- Configures connection pool
- Initializes schema (creates tables, views, indexes)

### 4. Repository Initialization
```go
repo := repository.NewMemoryRepository(db.GetDB(), logger)
```

### 5. MCP Server Creation
```go
mcpServer := mcp.NewServer("TinyBrain Memory Storage", "1.0.0",
    "Security-focused LLM memory storage MCP server", logger)
```

### 6. Tool Registration
```go
tinyBrain.registerTools(mcpServer)
```

Registers all 40 MCP tools for memory operations.

### 7. Server Start
```go
if err := mcpServer.ServeStdio(); err != nil {
    logger.Fatal("Server error", "error", err)
}
```

The server starts listening on stdin/stdout for MCP protocol messages.

## Performance Considerations

TinyBrain is optimized for security assessment workloads with the following performance characteristics:

### Current Benchmark Figures

These are current benchmark figures from repository testing:

- **Memory Entry Creation**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Relationship Queries**: ~500 queries/second
- **Database Size**: ~1MB per 10,000 memory entries

### Optimization Features

**Connection Pooling**
- Optimized for SQLite single-writer model
- Single connection with maximum reuse
- No connection overhead for concurrent operations

**Index Strategy**
- Comprehensive indexing for all query patterns
- Indexes on frequently queried fields (session_id, category, priority, etc.)
- Composite indexes for common query combinations

**Full-Text Search**
- FTS5 virtual tables for semantic search
- Optimized for natural language queries
- Automatic index maintenance

**Caching**
- Access pattern tracking for intelligent caching
- `access_count` and `accessed_at` fields track usage
- Frequently accessed memories prioritized in context summaries

### Database Configuration

The SQLite database uses optimized pragma settings:
- **WAL mode**: Better concurrency for reads during writes
- **Foreign keys**: Enabled for referential integrity
- **Synchronous NORMAL**: Balance between safety and performance
- **Busy timeout**: 30 seconds to handle concurrent access
- **Shared cache**: Memory efficiency for single connection

## Next Steps

- Learn about the [Memory Model](memory-model.md) for details on memory entries and relationships
- Explore [API Overview](api/overview.md) for the complete MCP tool reference
- See [Advanced Features](advanced-features.md) for context snapshots and task tracking
