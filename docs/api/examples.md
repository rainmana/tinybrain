---
layout: default
title: API Examples
nav_order: 7
parent: API
description: "Practical JSON examples for TinyBrain MCP tools"
---

# API Examples

This page provides curated JSON examples for common TinyBrain operations, sourced from the repository's usage examples.

## Table of Contents
- [Session Management](#session-management)
- [Memory Operations](#memory-operations)
- [Relationship Management](#relationship-management)
- [Search Operations](#search-operations)
- [Context & Task Tracking](#context--task-tracking)
- [Cleanup Operations](#cleanup-operations)
- [System Monitoring](#system-monitoring)

## Session Management

### Create Session

Create a new security assessment session:

```json
{
  "name": "create_session",
  "arguments": {
    "name": "Web Application Security Review",
    "description": "Comprehensive security review of e-commerce web application",
    "task_type": "security_review",
    "metadata": "{\"target\": \"shop.example.com\", \"scope\": \"web-application\", \"framework\": \"django\"}"
  }
}
```

## Memory Operations

### Store Memory - Vulnerability

Store a security vulnerability finding:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "SQL Injection in User Search",
    "content": "Found SQL injection vulnerability in user search functionality. The search parameter is directly concatenated into SQL query without sanitization. Payload: '; DROP TABLE users; --",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.95,
    "tags": "[\"sql-injection\", \"search\", \"critical\", \"data-loss\"]",
    "source": "code-review"
  }
}
```

### Store Memory - Exploit

Store an exploit or proof-of-concept:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_789",
    "title": "ROP Chain Construction",
    "content": "Successfully constructed ROP chain to bypass NX bit. Found gadgets at:\n- pop rdi; ret (0x4006a3)\n- pop rsi; ret (0x4006a1)\n- system() (0x4004e0)\n\nChain length: 3 gadgets",
    "category": "exploit",
    "priority": 8,
    "confidence": 0.85,
    "tags": "[\"rop-chain\", \"nx-bypass\", \"gadgets\", \"system-call\"]",
    "source": "exploit-development"
  }
}
```

### Get Memory

Retrieve a specific memory by ID:

```json
{
  "name": "get_memory",
  "arguments": {
    "memory_id": "memory_sql_injection"
  }
}
```

## Relationship Management

### Create Relationship

Link related memories (e.g., vulnerability to exploit):

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "memory_sql_injection",
    "target_memory_id": "memory_xss",
    "relationship_type": "related_to",
    "strength": 0.7,
    "description": "Both vulnerabilities can be combined for advanced attacks"
  }
}
```

### Get Related Memories

Find memories related through relationships:

```json
{
  "name": "get_related_memories",
  "arguments": {
    "memory_id": "memory_sql_injection",
    "relationship_type": "exploits",
    "limit": 5
  }
}
```

## Search Operations

### Search Memories - Semantic

Search for high-priority vulnerabilities using semantic search:

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "critical vulnerability",
    "session_id": "session_123",
    "min_priority": 8,
    "search_type": "semantic",
    "limit": 10
  }
}
```

### Search Memories - Category and Tags

Search with category and tag filters:

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "authentication",
    "categories": "[\"vulnerability\", \"exploit\"]",
    "tags": "[\"sql-injection\", \"bypass\"]",
    "search_type": "exact",
    "limit": 20
  }
}
```

### Search Memories - Fuzzy

Fuzzy search for similar content:

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "buffer overflow",
    "search_type": "fuzzy",
    "min_confidence": 0.7,
    "limit": 15
  }
}
```

## Context & Task Tracking

### Get Context Summary

Get relevant memory summary for current task:

```json
{
  "name": "get_context_summary",
  "arguments": {
    "session_id": "session_123",
    "current_task": "Analyzing authentication vulnerabilities in web application",
    "max_memories": 15
  }
}
```

### Create Context Snapshot

Save current context state:

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

### Create Task Progress

Initialize task progress tracking:

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

### Update Task Progress

Update progress on an existing task:

```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_123",
    "task_name": "Authentication Security Review",
    "stage": "vulnerability-analysis",
    "status": "in_progress",
    "progress_percentage": 60,
    "notes": "Completed SQL injection analysis, now reviewing session management"
  }
}
```

## Cleanup Operations

### Cleanup Unused Memories - Dry Run

Preview cleanup without actually deleting:

```json
{
  "name": "cleanup_unused_memories",
  "arguments": {
    "max_unused_days": 30,
    "dry_run": true
  }
}
```

**Example Response**:
```json
{
  "deleted_count": 15,
  "max_unused_days": 30,
  "dry_run": true,
  "message": "Cleaned up 15 unused memories"
}
```

### Cleanup Old Memories

Remove memories older than specified age:

```json
{
  "name": "cleanup_old_memories",
  "arguments": {
    "max_age_days": 90,
    "dry_run": false
  }
}
```

### Cleanup Low Priority Memories

Remove low-priority memories:

```json
{
  "name": "cleanup_low_priority_memories",
  "arguments": {
    "max_priority": 3,
    "dry_run": false
  }
}
```

## System Monitoring

### Health Check

Verify database connectivity and system status:

```json
{
  "name": "health_check",
  "arguments": {}
}
```

### Get Database Statistics

Retrieve database usage statistics:

```json
{
  "name": "get_database_stats",
  "arguments": {}
}
```

**Example Response**:
```json
{
  "total_memories": 1523,
  "total_sessions": 42,
  "total_relationships": 856,
  "categories": {
    "vulnerability": 342,
    "exploit": 178,
    "finding": 445,
    "recommendation": 234
  },
  "priority_distribution": {
    "critical": 89,
    "high": 234,
    "medium": 567,
    "low": 633
  }
}
```

## Common Patterns

### Security Code Review Workflow

1. **Create Session**:
```json
{
  "name": "create_session",
  "arguments": {
    "name": "Code Review - Payment Module",
    "description": "Security review of payment processing code",
    "task_type": "security_review"
  }
}
```

2. **Store Vulnerabilities**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "SQL Injection in Payment Query",
    "content": "Direct SQL concatenation without parameterization",
    "category": "vulnerability",
    "priority": 10,
    "confidence": 0.95
  }
}
```

3. **Create Relationships**:
```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "vuln_001",
    "target_memory_id": "exploit_001",
    "relationship_type": "exploits"
  }
}
```

4. **Track Progress**:
```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_456",
    "task_name": "Payment Module Review",
    "status": "completed",
    "progress_percentage": 100
  }
}
```

### Penetration Testing Workflow

1. **Create Session** (penetration_test type)
2. **Store reconnaissance data** (category: finding)
3. **Store identified vulnerabilities** (category: vulnerability)
4. **Store exploits** (category: exploit)
5. **Link vulnerabilities to exploits** (relationship_type: exploits)
6. **Track progress through stages** (recon → exploitation → reporting)

### Exploit Development Workflow

1. **Create Session** (exploit_dev type)
2. **Store vulnerability analysis** (category: vulnerability)
3. **Store exploit techniques** (category: technique)
4. **Store working exploits** (category: exploit)
5. **Store payloads** (category: payload)
6. **Create dependency relationships** (relationship_type: depends_on)

## Best Practices

Based on repository examples:

1. **Use Descriptive Titles**: Make memory titles clear and searchable
2. **Set Appropriate Priorities**: Use 8-10 for critical vulnerabilities, 5-7 for medium, 1-4 for low
3. **Tag Consistently**: Use consistent tagging conventions across sessions
4. **Create Relationships**: Link related vulnerabilities and exploits
5. **Update Progress**: Track task progress for long-running security assessments
6. **Use Context Summaries**: Get relevant memories before starting new tasks
7. **Regular Searches**: Search for similar issues to avoid duplication

## Next Steps

- Review [API Overview](overview.md) for complete tool reference
- See [Workflows](../workflows/security-assessment.md) for comprehensive workflow examples
- Check [Memory Model](../memory-model.md) for field semantics and constraints
