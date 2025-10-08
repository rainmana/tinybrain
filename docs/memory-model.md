---
layout: default
title: Memory Model
nav_order: 5
description: "Understanding TinyBrain's memory entry model and semantics"
---

# Memory Model

This document describes TinyBrain's memory entry model, including field semantics, categories, relationship types, and cleanup operations.

## Table of Contents
- [Memory Entry Fields](#memory-entry-fields)
- [Categories](#categories)
- [Relationship Types](#relationship-types)
- [Priority & Confidence Semantics](#priority--confidence-semantics)
- [Access Tracking](#access-tracking)
- [Cleanup Operations](#cleanup-operations)

## Memory Entry Fields

Each memory entry in TinyBrain contains the following fields:

### Core Fields

**id** (string, required)
- Unique identifier for the memory entry
- Automatically generated on creation
- Used for relationships and queries

**session_id** (string, required)
- ID of the session this memory belongs to
- Links memories to specific security assessments
- Used for filtering and context summaries

**title** (string, required)
- Short summary or title of the memory
- Searchable via full-text search
- Should be descriptive and specific

**content** (string, required)
- Full content of the memory
- Can contain detailed descriptions, code snippets, evidence
- Indexed for full-text search

**content_type** (string, default: "text")
- Type of content stored
- Valid values: `text`, `code`, `json`, `yaml`, `markdown`, `binary_ref`
- Used for proper rendering and processing

**category** (string, required)
- Security-focused categorization
- See [Categories](#categories) section for valid values
- Critical for organizing and filtering memories

### Security Fields

**priority** (integer, 0-10, default: 0)
- Priority level for the finding
- 0 = low priority, 10 = critical
- See [Priority Semantics](#priority-semantics) for guidance

**confidence** (float, 0.0-1.0, default: 0.5)
- Confidence score for the finding
- 0.0 = unverified, 1.0 = confirmed
- See [Confidence Semantics](#confidence-semantics) for guidance

**tags** (JSON array, optional)
- Array of tag strings for additional categorization
- Example: `["sql-injection", "authentication", "critical"]`
- Used for tag-based search

**source** (string, optional)
- Where the information came from
- Examples: `"code-review"`, `"manual-testing"`, `"nmap-scan"`
- Helps track information provenance

### Timestamp Fields

**created_at** (datetime, automatic)
- When the memory was created
- Set automatically on creation
- Immutable after creation

**updated_at** (datetime, automatic)
- When the memory was last updated
- Set automatically on creation and updates
- Updated by `update_memory` tool

**accessed_at** (datetime, automatic)
- When the memory was last accessed
- Updated when memory is retrieved or appears in search results
- Used for access-based cleanup

**access_count** (integer, default: 0)
- How many times the memory has been accessed
- Incremented on retrievals and searches
- Used for relevance scoring and cleanup

## Categories

TinyBrain supports specialized security-focused categories for memory entries:

### finding
General security finding or observation that doesn't fit other categories.

**Usage**: Initial discoveries, general observations, unclassified findings

**Example**:
```json
{
  "title": "Weak Password Policy",
  "category": "finding",
  "priority": 5,
  "confidence": 0.8
}
```

### vulnerability
Confirmed or suspected security vulnerability.

**Usage**: Security weaknesses, flaws, misconfigurations

**Example**:
```json
{
  "title": "SQL Injection in Login Form",
  "category": "vulnerability",
  "priority": 9,
  "confidence": 0.95
}
```

### exploit
Proof-of-concept exploit or exploitation technique.

**Usage**: Working exploits, exploitation methods, attack vectors

**Example**:
```json
{
  "title": "SQL Injection Bypass Authentication",
  "category": "exploit",
  "priority": 9,
  "confidence": 0.9
}
```

### payload
Specific payload for exploitation.

**Usage**: Attack payloads, malicious inputs, exploit code

**Example**:
```json
{
  "title": "XSS Payload for Cookie Theft",
  "category": "payload",
  "priority": 7,
  "confidence": 0.85
}
```

### technique
Security testing or attack technique.

**Usage**: Testing methodologies, attack techniques, procedures

**Example**:
```json
{
  "title": "Kerberoasting Technique",
  "category": "technique",
  "priority": 8,
  "confidence": 0.9
}
```

### tool
Security tool or software used in assessment.

**Usage**: Tools, scripts, software references

**Example**:
```json
{
  "title": "Burp Suite Pro Configuration",
  "category": "tool",
  "priority": 5,
  "confidence": 1.0
}
```

### reference
External reference, documentation, or resource.

**Usage**: CVEs, CWEs, documentation links, research papers

**Example**:
```json
{
  "title": "CWE-89: SQL Injection",
  "category": "reference",
  "priority": 6,
  "confidence": 1.0
}
```

### context
Contextual information about the target or environment.

**Usage**: Environment details, architecture, configurations

**Example**:
```json
{
  "title": "Application Architecture Overview",
  "category": "context",
  "priority": 4,
  "confidence": 0.9
}
```

### hypothesis
Untested hypothesis or theory about potential vulnerabilities.

**Usage**: Suspected issues, theories, areas to investigate

**Example**:
```json
{
  "title": "Possible Race Condition in Payment Processing",
  "category": "hypothesis",
  "priority": 6,
  "confidence": 0.4
}
```

### evidence
Evidence supporting a finding or vulnerability.

**Usage**: Proof, logs, screenshots, test results

**Example**:
```json
{
  "title": "Error Message Revealing Database Structure",
  "category": "evidence",
  "priority": 7,
  "confidence": 1.0
}
```

### recommendation
Remediation or security recommendation.

**Usage**: Fixes, mitigations, security improvements

**Example**:
```json
{
  "title": "Implement Prepared Statements",
  "category": "recommendation",
  "priority": 9,
  "confidence": 1.0
}
```

### note
General note or comment.

**Usage**: Observations, comments, status updates

**Example**:
```json
{
  "title": "Testing Completed for Login Module",
  "category": "note",
  "priority": 2,
  "confidence": 1.0
}
```

## Relationship Types

TinyBrain supports various relationship types to map connections between memories:

### depends_on
Source depends on target for success or execution.

**Example**: Exploit depends on vulnerability

### causes
Source causes or leads to target.

**Example**: Vulnerability causes data breach

### mitigates
Source mitigates or reduces target.

**Example**: Recommendation mitigates vulnerability

### exploits
Source exploits target.

**Example**: Exploit exploits vulnerability

### references
Source references target for additional information.

**Example**: Vulnerability references CWE

### contradicts
Source contradicts target.

**Example**: Test result contradicts hypothesis

### supports
Source supports or confirms target.

**Example**: Evidence supports vulnerability

### related_to
General relationship between source and target.

**Example**: Two vulnerabilities are related

### parent_of
Source is parent of target (hierarchical).

**Example**: Main finding is parent of sub-findings

### child_of
Source is child of target (hierarchical).

**Example**: Sub-finding is child of main finding

## Priority & Confidence Semantics

### Priority Semantics

Priority levels guide risk assessment and remediation prioritization:

**Priority 10 (Critical)**
- Immediate exploitation possible
- Severe business impact
- Data breach or system compromise
- Examples: Remote code execution, SQL injection with data access

**Priority 8-9 (High)**
- Exploitation likely
- Significant business impact
- Potential for major security incident
- Examples: Authentication bypass, privilege escalation

**Priority 5-7 (Medium)**
- Exploitation possible with effort
- Moderate business impact
- Security weakness requiring attention
- Examples: XSS, information disclosure

**Priority 1-4 (Low)**
- Difficult to exploit or limited impact
- Minor security concern
- Best practice improvement
- Examples: Weak password policy, missing security headers

**Priority 0 (Informational)**
- No direct security impact
- Contextual information
- General observations

### Confidence Semantics

Confidence scores indicate the certainty of findings:

**Confidence 0.9-1.0 (Confirmed)**
- Fully validated and confirmed
- Exploitation demonstrated
- Evidence collected
- High certainty

**Confidence 0.7-0.8 (Likely)**
- Strong evidence
- Testing indicates presence
- Not fully confirmed
- High probability

**Confidence 0.5-0.6 (Possible)**
- Some evidence or indicators
- Requires further testing
- Moderate probability
- Unverified hypothesis

**Confidence 0.0-0.4 (Unverified)**
- Speculation or theory
- Minimal evidence
- Requires investigation
- Low certainty

## Access Tracking

TinyBrain tracks memory access patterns for optimization and cleanup:

### accessed_at Field
Updated whenever a memory is:
- Retrieved via `get_memory`
- Included in search results
- Referenced in relationships
- Part of context summaries

### access_count Field
Incremented when a memory is:
- Explicitly retrieved
- Returned in search results
- Used in context generation

### Usage
Access tracking enables:
- **Relevance scoring**: Frequently accessed memories rank higher
- **Context optimization**: Popular memories prioritized in summaries
- **Cleanup decisions**: Unused memories can be identified and cleaned

## Cleanup Operations

TinyBrain provides three cleanup operations for managing memory lifecycle:

### cleanup_old_memories

**Purpose**: Remove memories based on age

**Parameters**:
- `max_age_days` (integer, required) - Maximum age in days
- `dry_run` (boolean, optional) - Preview without deleting

**Usage**: Remove outdated findings from completed assessments

**Example**:
```json
{
  "name": "cleanup_old_memories",
  "arguments": {
    "max_age_days": 90,
    "dry_run": true
  }
}
```

### cleanup_low_priority_memories

**Purpose**: Remove memories based on priority threshold

**Parameters**:
- `max_priority` (integer, required) - Maximum priority to delete (0-10)
- `dry_run` (boolean, optional) - Preview without deleting

**Usage**: Remove low-priority notes and observations

**Example**:
```json
{
  "name": "cleanup_low_priority_memories",
  "arguments": {
    "max_priority": 3,
    "dry_run": false
  }
}
```

### cleanup_unused_memories

**Purpose**: Remove memories based on access patterns

**Parameters**:
- `max_unused_days` (integer, required) - Days since last access
- `dry_run` (boolean, optional) - Preview without deleting

**Usage**: Remove memories that haven't been accessed recently

**Example**:
```json
{
  "name": "cleanup_unused_memories",
  "arguments": {
    "max_unused_days": 30,
    "dry_run": true
  }
}
```

**Dry Run Example Response**:
```json
{
  "deleted_count": 15,
  "max_unused_days": 30,
  "dry_run": true,
  "message": "Cleaned up 15 unused memories"
}
```

When `dry_run` is true, the operation shows what would be deleted without actually removing the memories.

## Best Practices

1. **Use Descriptive Titles**: Make memory titles clear and searchable
2. **Set Appropriate Priorities**: Use 8-10 for critical vulnerabilities, 5-7 for medium, 1-4 for low
3. **Tag Consistently**: Use consistent tagging conventions across sessions
4. **Create Relationships**: Link related vulnerabilities and exploits
5. **Update Progress**: Track task progress for long-running security assessments
6. **Use Context Summaries**: Get relevant memories before starting new tasks
7. **Regular Searches**: Search for similar issues to avoid duplication

## Next Steps

- Explore [API Examples](api/examples.md) for practical memory usage patterns
- Review [Workflows](workflows/security-assessment.md) for real-world assessment examples
- See [Advanced Features](advanced-features.md) for context snapshots and task tracking
