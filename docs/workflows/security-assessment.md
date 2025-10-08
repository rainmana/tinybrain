---
layout: default
title: Security Assessment
nav_order: 8
parent: Workflows
description: "Security code review workflow using TinyBrain"
---

# Security Assessment Workflow

This document describes the complete workflow for conducting security code reviews using TinyBrain, based on examples from the repository.

## Table of Contents
- [Overview](#overview)
- [Workflow Phases](#workflow-phases)
- [Session Creation](#session-creation)
- [Finding Storage](#finding-storage)
- [Relationship Mapping](#relationship-mapping)
- [Progress Tracking](#progress-tracking)
- [Context Management](#context-management)
- [Report Generation](#report-generation)

## Overview

Security code reviews using TinyBrain follow a structured workflow that ensures comprehensive vulnerability discovery, proper categorization, and relationship mapping between findings.

**Workflow Pattern**:
1. Create session with security_review task type
2. Store vulnerabilities with proper categorization
3. Create relationships between related findings
4. Track progress through review stages
5. Use context summaries for analysis
6. Generate comprehensive reports

## Workflow Phases

### Phase 1: Initialization
- Create security review session
- Define scope and metadata
- Initialize task progress tracking

### Phase 2: Discovery
- Store identified vulnerabilities
- Categorize by severity and type
- Tag for easy retrieval

### Phase 3: Analysis
- Create relationships between findings
- Store supporting evidence
- Document exploitation paths

### Phase 4: Validation
- Test and validate vulnerabilities
- Update confidence scores
- Store proof-of-concept code

### Phase 5: Reporting
- Generate context summaries
- Export session data
- Create remediation recommendations

## Session Creation

### Create Security Review Session

Initialize a new security code review session:

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

**Key Fields**:
- `task_type`: Always `"security_review"` for code reviews
- `metadata`: Store target information, scope, frameworks, etc.

## Finding Storage

### Store SQL Injection Vulnerability

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

### Store XSS Vulnerability

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "XSS in Product Reviews",
    "content": "Cross-site scripting vulnerability in product review form. User input is not properly escaped before rendering. Can be exploited to steal session cookies.",
    "category": "vulnerability",
    "priority": 7,
    "confidence": 0.85,
    "tags": "[\"xss\", \"reviews\", \"session-hijacking\"]",
    "source": "code-review"
  }
}
```

### Store Evidence

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "Error Message Revealing Database Structure",
    "content": "SQL error messages expose database table structure and column names. Error: 'Column \"password_hash\" does not exist in table \"users\"'",
    "category": "evidence",
    "priority": 7,
    "confidence": 1.0,
    "tags": "[\"information-disclosure\", \"sql-error\", \"database-structure\"]",
    "source": "testing"
  }
}
```

## Relationship Mapping

### Link Related Vulnerabilities

Create relationships between vulnerabilities that can be chained:

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

### Link Vulnerability to Evidence

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "memory_sql_injection",
    "target_memory_id": "memory_error_message",
    "relationship_type": "supports",
    "strength": 0.9,
    "description": "Error message confirms SQL injection vulnerability"
  }
}
```

### Link Vulnerability to Recommendation

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "memory_recommendation_prepared_stmt",
    "target_memory_id": "memory_sql_injection",
    "relationship_type": "mitigates",
    "strength": 1.0,
    "description": "Prepared statements will eliminate SQL injection vulnerability"
  }
}
```

## Progress Tracking

### Initialize Task Progress

Create task progress entry for the review:

```json
{
  "name": "create_task_progress",
  "arguments": {
    "session_id": "session_123",
    "task_name": "Authentication Module Review",
    "stage": "Initial Discovery",
    "status": "in_progress",
    "progress_percentage": 25,
    "notes": "Completed initial vulnerability discovery. Found SQL injection and XSS vulnerabilities."
  }
}
```

### Update Progress

Update task as review progresses:

```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_123",
    "task_name": "Authentication Module Review",
    "stage": "vulnerability-analysis",
    "status": "in_progress",
    "progress_percentage": 60,
    "notes": "Completed SQL injection analysis, now reviewing session management"
  }
}
```

### Mark Task Complete

```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_123",
    "task_name": "Authentication Module Review",
    "stage": "completed",
    "status": "completed",
    "progress_percentage": 100,
    "notes": "Review complete. Identified 12 vulnerabilities, 5 high-priority."
  }
}
```

## Context Management

### Get Context Before Analysis

Before starting analysis of a new component, get relevant context:

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

This returns relevant memories about authentication, helping avoid duplication and building on previous findings.

### Create Context Snapshot

Save context at key milestones:

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

### Search for Similar Vulnerabilities

Check for similar issues to avoid duplication:

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "SQL injection",
    "session_id": "session_123",
    "search_type": "semantic",
    "limit": 10
  }
}
```

## Report Generation

### Get High-Priority Findings

Retrieve all critical vulnerabilities for reporting:

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "vulnerability",
    "session_id": "session_123",
    "min_priority": 8,
    "category": "vulnerability",
    "limit": 100
  }
}
```

### Export Session Data

Export complete session for reporting:

```json
{
  "name": "export_session_data",
  "arguments": {
    "session_id": "session_123",
    "include_relationships": true
  }
}
```

### Get Database Statistics

Include statistics in executive summary:

```json
{
  "name": "get_database_stats",
  "arguments": {}
}
```

## Best Practices

Based on repository examples and workflows:

1. **Session Creation**: Always create a session before starting the review
2. **Consistent Categorization**: Use appropriate categories (vulnerability, evidence, recommendation)
3. **Priority Assignment**: Use 8-10 for critical issues, 5-7 for medium, 1-4 for low
4. **Tag Consistently**: Use standard tags (sql-injection, xss, authentication, etc.)
5. **Create Relationships**: Link vulnerabilities to evidence and recommendations
6. **Track Progress**: Update task progress regularly throughout the review
7. **Use Context**: Get context summaries before analyzing new components
8. **Search First**: Check for similar findings to avoid duplication

## Example: Complete Review Workflow

### 1. Initialize
```json
{
  "name": "create_session",
  "arguments": {
    "name": "Payment Module Security Review",
    "task_type": "security_review"
  }
}
```

### 2. Store Findings
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "SQL Injection in Payment Query",
    "category": "vulnerability",
    "priority": 10,
    "confidence": 0.95
  }
}
```

### 3. Add Evidence
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "SQL Injection Test Results",
    "category": "evidence",
    "priority": 9,
    "confidence": 1.0
  }
}
```

### 4. Link Findings
```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "evidence_001",
    "target_memory_id": "vuln_001",
    "relationship_type": "supports"
  }
}
```

### 5. Add Recommendation
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Implement Prepared Statements",
    "category": "recommendation",
    "priority": 10,
    "confidence": 1.0
  }
}
```

### 6. Track Progress
```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_456",
    "status": "completed",
    "progress_percentage": 100
  }
}
```

## Next Steps

- See [Penetration Testing Workflow](penetration-testing.md) for pen-test patterns
- Review [Exploit Development Workflow](exploit-development.md) for exploit development
- Check [API Examples](../api/examples.md) for more JSON examples
