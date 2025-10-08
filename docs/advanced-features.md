---
layout: default
title: Advanced Features
nav_order: 14
description: "Advanced TinyBrain features for complex security assessments"
---

# Advanced Features

TinyBrain includes advanced features for managing complex, long-running security tasks with sophisticated context management and task tracking capabilities.

## Table of Contents
- [Overview](#overview)
- [Context Snapshots](#context-snapshots)
- [Task Progress Tracking](#task-progress-tracking)
- [Enhanced Memory Summarization](#enhanced-memory-summarization)
- [Use Cases](#use-cases)
- [Performance](#performance)
- [Future Enhancements](#future-enhancements)

## Overview

Advanced features enable:
- **Context preservation** across long-running assessments
- **Progress monitoring** through assessment stages
- **Memory summarization** for efficient context management
- **Milestone tracking** for complex security tasks

## Context Snapshots

Context snapshots allow you to capture the current state of a security assessment at any point in time.

### Features

- **Context Data Storage**: Store complex context information as JSON
- **Automatic Memory Summarization**: Generate summaries of high-priority findings
- **Timestamp Tracking**: Track when snapshots were created
- **Session Association**: Link snapshots to specific security sessions

### MCP Tools

- `create_context_snapshot` - Create a new context snapshot
- `get_context_snapshot` - Retrieve a specific snapshot
- `list_context_snapshots` - List all snapshots for a session

### Example Usage

**Create Context Snapshot**:
```json
{
  "name": "create_context_snapshot",
  "arguments": {
    "session_id": "session_123",
    "name": "Initial Assessment Complete",
    "description": "Context after discovering initial vulnerabilities",
    "context_data": "{\"current_findings\":[\"SQL injection\",\"XSS\"],\"next_steps\":[\"validate\",\"remediate\"]}"
  }
}
```

### Database Schema

The `context_snapshots` table includes:
- `id` - Unique snapshot identifier
- `session_id` - Associated session
- `name` - Snapshot name
- `description` - Snapshot description
- `context_data` - JSON context data
- `memory_summary` - Generated memory summary
- `created_at` - Creation timestamp

## Task Progress Tracking

Multi-stage task progress tracking for complex security assessments.

### Features

- **Task Lifecycle Management**: Track tasks from pending to completed
- **Stage Tracking**: Monitor progress through different assessment stages
- **Status Transitions**: Automatic timestamp management (started_at, completed_at)
- **Progress Percentage**: Track completion percentage (0-100%)
- **Detailed Notes**: Add context and progress notes
- **Status Filtering**: List tasks by status (pending, in_progress, completed, failed, blocked)

### MCP Tools

- `create_task_progress` - Create a new task progress entry
- `get_task_progress` - Retrieve a specific task progress
- `list_task_progress` - List tasks with optional status filtering
- `update_task_progress` - Update task progress status

### Example Usage

**Create Task Progress**:
```json
{
  "name": "create_task_progress",
  "arguments": {
    "session_id": "session_123",
    "task_name": "Vulnerability Assessment",
    "stage": "Initial Discovery",
    "status": "in_progress",
    "progress_percentage": 25,
    "notes": "Completed initial vulnerability discovery. Found SQL injection and XSS vulnerabilities."
  }
}
```

**Update Task Progress**:
```json
{
  "name": "update_task_progress",
  "arguments": {
    "task_id": "task_456",
    "stage": "Validation",
    "status": "in_progress",
    "progress_percentage": 60,
    "notes": "Validated SQL injection. Moving to impact assessment."
  }
}
```

### Database Schema

The `task_progress` table includes:
- `id` - Unique task identifier
- `session_id` - Associated session
- `task_name` - Task name
- `stage` - Current stage
- `status` - Task status (pending, in_progress, completed, failed, blocked)
- `progress_percentage` - Completion percentage (0-100)
- `notes` - Progress notes
- `started_at` - When task started (nullable)
- `completed_at` - When task completed (nullable)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Status Types

- **pending** - Task not yet started
- **in_progress** - Task actively being worked on
- **completed** - Task finished successfully
- **failed** - Task encountered errors or cannot complete
- **blocked** - Task waiting on external dependency

## Enhanced Memory Summarization

Automatic generation of memory summaries for context snapshots.

### Features

- **High-Priority Focus**: Prioritizes findings with high priority and confidence
- **Structured Format**: Organized summary with categories and scores
- **Content Preview**: Includes brief content summaries (first 100 characters)
- **Real-time Generation**: Summaries generated when snapshots are created

### Example Summary Output

```
Recent High-Priority Findings:
1. [vulnerability] SQL Injection Vulnerability in Login Form (Priority: 9, Confidence: 0.9)
   Found a critical SQL injection vulnerability in the login form. The user input is not properly sanit...
2. [vulnerability] XSS Vulnerability in Search Function (Priority: 7, Confidence: 0.8)
   Discovered a cross-site scripting (XSS) vulnerability in the search functionality. User input is ref...
```

## Use Cases

### 1. Long-Running Security Assessments

Track complex security assessments across multiple stages:

**Assessment Stages**:
- Initial Discovery → Validation → Impact Assessment → Remediation Planning → Implementation → Verification

**Example Workflow**:
1. Create task progress for each stage
2. Update percentage as work progresses
3. Create snapshots at key milestones
4. Track blockers and failures

### 2. Context Preservation

Maintain context across long-running tasks:

**Benefits**:
- Capture assessment state at key milestones
- Restore context when resuming work
- Share context with team members
- Reduce context window overhead

**Example**:
```json
{
  "name": "create_context_snapshot",
  "arguments": {
    "session_id": "session_789",
    "name": "Week 1 Complete - Initial Findings",
    "description": "End of week 1: Discovered 12 vulnerabilities, 5 critical"
  }
}
```

### 3. Progress Monitoring

Monitor progress on security tasks:

**Capabilities**:
- Track completion percentage
- Identify blocked or failed tasks
- Generate progress reports
- Estimate time to completion

**Example**:
```json
{
  "name": "list_task_progress",
  "arguments": {
    "session_id": "session_789",
    "status": "blocked"
  }
}
```

### 4. Memory Management

Optimize memory usage for long-running tasks:

**Strategies**:
- Generate summaries of important findings
- Focus on high-priority information
- Reduce context window usage
- Archive completed work

## Performance

Current performance characteristics from testing:

- **Context Snapshots**: < 100ms creation time
- **Memory Summarization**: < 50ms generation time
- **Task Progress**: < 50ms CRUD operations
- **Database Queries**: Optimized with proper indexing

All operations are optimized for minimal latency in typical security assessment workflows.

## Future Enhancements

The advanced features provide a foundation for additional capabilities currently listed in the roadmap:

### Planned Enhancements

- **Semantic Search**: Vector embeddings for similarity search
- **Memory Deduplication**: Detect and merge similar findings
- **Export/Import**: Backup and restore memory data
- **Memory Templates**: Pre-defined patterns for common security issues
- **Real-time Notifications**: Alert on critical findings or task completions
- **Batch Operations**: Bulk memory management operations

**Note**: These are placeholder items for future development. Current implementation provides the foundational infrastructure.

## Testing

All advanced features have been thoroughly tested:

- ✅ Context snapshot creation and retrieval
- ✅ Memory summarization generation
- ✅ Task progress lifecycle management
- ✅ Status transition handling
- ✅ Database integrity and constraints
- ✅ MCP protocol compliance
- ✅ Error handling and edge cases

*Note*: The repository README states "90%+ test coverage." Users can verify current coverage by running `go test -cover` in the repository.

## Integration

Advanced features integrate seamlessly with existing TinyBrain functionality:

### With Workflows

- **Security Assessment**: Track progress through code review stages
- **Penetration Testing**: Create snapshots at each testing phase
- **Exploit Development**: Monitor exploit development milestones

### With Datasets

- **CWE Integration**: Track CWE pattern discovery progress
- **OWASP Mapping**: Monitor OWASP Top 10 coverage percentage
- **Template Usage**: Apply templates and track application progress

### With AI Assistants

- **Context Management**: Provide relevant context summaries to AI assistants
- **Progress Updates**: Keep AI assistants informed of assessment progress
- **Milestone Notifications**: Alert on critical progress events

## Best Practices

1. **Create Snapshots Regularly**: At key milestones and decision points
2. **Update Task Progress**: Keep progress current for accurate tracking
3. **Use Meaningful Names**: Make snapshot and task names descriptive
4. **Include Detailed Notes**: Add context in notes fields for future reference
5. **Monitor Completion**: Regularly review task progress and blockers
6. **Generate Summaries**: Use context summaries to reduce context window usage

## Next Steps

- Review [API Examples](api/examples.md) for practical usage
- See [Workflows](workflows/security-assessment.md) for integration examples
- Check [Memory Model](memory-model.md) for field semantics
