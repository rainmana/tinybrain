---
layout: default
title: API Overview
nav_order: 6
parent: API
description: "Complete MCP tool reference for TinyBrain"
---

# API Overview

TinyBrain provides a comprehensive set of 40 MCP (Model Context Protocol) tools for complete LLM memory management.

## Table of Contents
- [About MCP](#about-mcp)
- [Tool Categories](#tool-categories)
- [Core Memory Operations](#core-memory-operations)
- [Session & Task Management](#session--task-management)
- [Advanced Memory Features](#advanced-memory-features)
- [Security Templates & Batch Operations](#security-templates--batch-operations)
- [Memory Lifecycle & Cleanup](#memory-lifecycle--cleanup)
- [AI-Enhanced Search](#ai-enhanced-search)
- [Real-Time Notifications](#real-time-notifications)
- [System Monitoring](#system-monitoring)

## About MCP

The Model Context Protocol (MCP) is a standardized protocol for LLM-tool communication. It allows AI assistants like Claude to interact with external systems through well-defined tool interfaces.

TinyBrain implements the MCP specification to provide:
- **Standardized Communication**: JSON-based request/response format
- **Type Safety**: Structured parameters with validation
- **Error Handling**: Clear error messages and status codes
- **Stateless Operations**: Each tool call is independent

All TinyBrain tools follow the MCP pattern:
```json
{
  "name": "tool_name",
  "arguments": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

## Tool Categories

TinyBrain's 40 tools are organized into 8 functional categories:

1. **Core Memory Operations** (8 tools) - Basic memory CRUD operations
2. **Session & Task Management** (6 tools) - Session and task tracking
3. **Advanced Memory Features** (8 tools) - Relationships, snapshots, context
4. **Security Templates & Batch Operations** (6 tools) - Templates and bulk operations
5. **Memory Lifecycle & Cleanup** (4 tools) - Memory maintenance
6. **AI-Enhanced Search** (3 tools) - Semantic search and embeddings
7. **Real-Time Notifications** (4 tools) - Alerts and notifications
8. **System Monitoring** (1 tool) - Health checks

## Core Memory Operations

**8 tools for basic memory management**

### store_memory
Store new memory entries with security categorization.

**Purpose**: Create and store new information in the memory system

**Key Parameters**:
- `session_id` - Session to associate memory with
- `title` - Memory title/summary
- `content` - Full memory content
- `category` - Security category (vulnerability, exploit, etc.)
- `priority` - Priority level 0-10
- `confidence` - Confidence score 0.0-1.0
- `tags` - Array of tags for categorization

### get_memory
Retrieve a specific memory entry by ID.

**Purpose**: Fetch detailed information about a specific memory

**Key Parameters**:
- `memory_id` - Unique identifier of the memory to retrieve

### search_memories
Advanced search with multiple search strategies.

**Purpose**: Find memories using various search types and filters

**Key Parameters**:
- `query` - Search query text
- `session_id` - Filter by session (optional)
- `search_type` - Type of search (semantic, exact, fuzzy, tag, category, relationship)
- `categories` - Filter by categories
- `min_priority` - Minimum priority threshold
- `limit` - Maximum results to return

### update_memory
Update existing memory entries.

**Purpose**: Modify fields of an existing memory

**Key Parameters**:
- `memory_id` - ID of memory to update
- All memory fields are optional update targets

### delete_memory
Delete memory entries.

**Purpose**: Remove a memory from the system

**Key Parameters**:
- `memory_id` - ID of memory to delete

### find_similar_memories
Find similar memories by content.

**Purpose**: Discover related memories based on content similarity

**Key Parameters**:
- `content` - Content to find similar memories for
- `limit` - Maximum results to return

### check_duplicates
Check for duplicate memories.

**Purpose**: Identify potential duplicate entries

**Key Parameters**:
- `session_id` - Session to check (optional)
- `threshold` - Similarity threshold for duplicates

### get_memory_stats
Get comprehensive memory statistics.

**Purpose**: Retrieve usage statistics and metrics

**Returns**: Counts by category, priority distribution, access patterns

## Session & Task Management

**6 tools for organizing security assessments**

### create_session
Create new security assessment sessions.

**Purpose**: Initialize a new security assessment session

**Key Parameters**:
- `name` - Session name
- `description` - Session description
- `task_type` - Type of assessment (security_review, penetration_test, etc.)
- `metadata` - JSON metadata for session-specific data

### get_session
Retrieve session information.

**Purpose**: Get details about a specific session

**Key Parameters**:
- `session_id` - ID of session to retrieve

### list_sessions
List all sessions with filtering.

**Purpose**: Query sessions with filters

**Key Parameters**:
- `task_type` - Filter by task type
- `status` - Filter by status (active, paused, completed, archived)
- `limit` - Maximum results to return
- `offset` - Pagination offset

### create_task_progress
Create task progress entries.

**Purpose**: Initialize progress tracking for a task

**Key Parameters**:
- `session_id` - Associated session
- `task_name` - Name of the task
- `stage` - Current stage
- `status` - Task status (pending, in_progress, completed, failed, blocked)
- `progress_percentage` - Completion percentage 0-100

### update_task_progress
Update task progress.

**Purpose**: Update progress for an existing task

**Key Parameters**:
- `task_id` - ID of task to update
- `status` - New status
- `progress_percentage` - New completion percentage
- `notes` - Progress notes

### list_task_progress
List task progress entries.

**Purpose**: Query task progress with filters

**Key Parameters**:
- `session_id` - Filter by session
- `status` - Filter by status

## Advanced Memory Features

**8 tools for relationships, snapshots, and context**

### create_relationship
Create memory relationships.

**Purpose**: Link related memories (exploits, dependencies, mitigations)

**Key Parameters**:
- `source_memory_id` - Source memory ID
- `target_memory_id` - Target memory ID
- `relationship_type` - Type (depends_on, causes, mitigates, exploits, etc.)
- `strength` - Relationship strength 0.0-1.0
- `description` - Relationship description

### get_related_memories
Get related memories.

**Purpose**: Retrieve memories connected via relationships

**Key Parameters**:
- `memory_id` - Base memory ID
- `relationship_type` - Type of relationship to follow (optional)
- `limit` - Maximum results

### create_context_snapshot
Create context snapshots.

**Purpose**: Save current context state for later reference

**Key Parameters**:
- `session_id` - Associated session
- `name` - Snapshot name
- `description` - Snapshot description
- `context_data` - JSON context data

### get_context_snapshot
Retrieve context snapshots.

**Purpose**: Load a saved context snapshot

**Key Parameters**:
- `snapshot_id` - ID of snapshot to retrieve

### list_context_snapshots
List context snapshots.

**Purpose**: Query available snapshots for a session

**Key Parameters**:
- `session_id` - Filter by session

### get_context_summary
Get memory summaries for context.

**Purpose**: Generate relevant memory summary for current task

**Key Parameters**:
- `session_id` - Session to summarize
- `current_task` - Description of current task
- `max_memories` - Maximum memories to include

### export_session_data
Export session data.

**Purpose**: Export complete session data for backup or sharing

**Key Parameters**:
- `session_id` - Session to export
- `include_relationships` - Include relationship data (optional)

### import_session_data
Import session data.

**Purpose**: Import previously exported session data

**Key Parameters**:
- `session_data` - JSON session data to import

## Security Templates & Batch Operations

**6 tools for templates and bulk operations**

### get_security_templates
Get predefined security templates.

**Purpose**: Retrieve available memory templates

**Key Parameters**:
- `category` - Filter by category (optional)

### create_memory_from_template
Create memories from templates.

**Purpose**: Create memory using a predefined template

**Key Parameters**:
- `template_id` - Template to use
- `session_id` - Target session
- `values` - Template field values

### batch_create_memories
Bulk create memory entries.

**Purpose**: Create multiple memories in a single operation

**Key Parameters**:
- `memories` - Array of memory objects to create

### batch_update_memories
Bulk update memory entries.

**Purpose**: Update multiple memories in a single operation

**Key Parameters**:
- `updates` - Array of memory updates

### batch_delete_memories
Bulk delete memory entries.

**Purpose**: Delete multiple memories in a single operation

**Key Parameters**:
- `memory_ids` - Array of memory IDs to delete

### get_detailed_memory_info
Get detailed memory debugging info.

**Purpose**: Retrieve comprehensive debugging information for a memory

**Key Parameters**:
- `memory_id` - Memory to inspect

## Memory Lifecycle & Cleanup

**4 tools for memory maintenance**

### cleanup_old_memories
Age-based memory cleanup.

**Purpose**: Remove memories older than specified age

**Key Parameters**:
- `max_age_days` - Maximum age in days
- `dry_run` - Preview without deleting (optional)

### cleanup_low_priority_memories
Priority-based cleanup.

**Purpose**: Remove memories below priority threshold

**Key Parameters**:
- `max_priority` - Maximum priority to delete (0-10)
- `dry_run` - Preview without deleting (optional)

### cleanup_unused_memories
Access-based cleanup.

**Purpose**: Remove memories not accessed recently

**Key Parameters**:
- `max_unused_days` - Days since last access
- `dry_run` - Preview without deleting (optional)

### get_system_diagnostics
Comprehensive system diagnostics.

**Purpose**: Retrieve system health and diagnostic information

**Returns**: Database statistics, performance metrics, health status

## AI-Enhanced Search

**3 tools for semantic search capabilities**

**Note**: Embedding generation currently uses a hash-based placeholder implementation as indicated in the repository code comments. This provides a foundation for future AI model integration (OpenAI, Cohere, or local models).

### semantic_search
AI-powered semantic search.

**Purpose**: Search memories using semantic similarity

**Key Parameters**:
- `query` - Natural language query
- `session_id` - Filter by session (optional)
- `limit` - Maximum results

### generate_embedding
Generate embeddings for text.

**Purpose**: Create embedding vectors for similarity comparison

**Key Parameters**:
- `text` - Text to generate embedding for

**Note**: Currently uses hash-based placeholder for embedding generation.

### calculate_similarity
Calculate embedding similarity.

**Purpose**: Compute similarity score between two embeddings

**Key Parameters**:
- `embedding1` - First embedding vector
- `embedding2` - Second embedding vector

## Real-Time Notifications

**4 tools for alerts and notifications**

### get_notifications
Get notifications and alerts.

**Purpose**: Retrieve notifications for events and system activities

**Key Parameters**:
- `session_id` - Filter by session (optional)
- `read` - Filter by read status (optional)
- `limit` - Maximum results

### mark_notification_read
Mark notifications as read.

**Purpose**: Update notification read status

**Key Parameters**:
- `notification_id` - ID of notification to mark as read

### check_high_priority_memories
Check for high-priority alerts.

**Purpose**: Get notifications for high-priority memories

**Criteria**: Priority ≥8 and confidence ≥0.8

**Key Parameters**:
- `session_id` - Filter by session (optional)

### check_duplicate_memories
Check for duplicate alerts.

**Purpose**: Get notifications for potential duplicate memories

**Key Parameters**:
- `session_id` - Filter by session (optional)
- `threshold` - Similarity threshold

## System Monitoring

**1 tool for health checks**

### health_check
Perform system health checks.

**Purpose**: Verify database connectivity and system status

**Returns**: Health status, database info, connection details

**Example**:
```json
{
  "name": "health_check",
  "arguments": {}
}
```

## Next Steps

- See [API Examples](examples.md) for practical JSON examples of all major tools
- Review [Workflows](../workflows/security-assessment.md) for real-world usage patterns
- Check [Memory Model](../memory-model.md) for detailed field documentation
