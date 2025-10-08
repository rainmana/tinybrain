---
layout: default
title: CWE Integration
nav_order: 12
parent: Datasets
description: "CWE dataset integration guide for TinyBrain"
---

# CWE Integration Guide

This document describes how to integrate and use the CWE (Common Weakness Enumeration) dataset with TinyBrain for security assessments.

## Table of Contents
- [Overview](#overview)
- [Quick Start Commands](#quick-start-commands)
- [Search Strategies](#search-strategies)
- [Storing CWE Findings](#storing-cwe-findings)
- [Standards Compliance Mapping](#standards-compliance-mapping)
- [Workflow Examples](#workflow-examples)

## Overview

The CWE LLM Dataset (`CWE_LLM_DATASET.json`) provides structured CWE data optimized for AI assistant consumption and TinyBrain integration.

**Key Features**:
- Structured JSON format for easy LLM parsing
- Searchable by severity, language, and category
- Context-optimized to avoid overloading
- TinyBrain-ready templates and commands
- Correlation capabilities for pattern matching

## Quick Start Commands

### Load CWE Dataset into TinyBrain

**Load Complete CWE Top 25 Dataset**:
```bash
cline "Load CWE Top 25 security patterns from CWE_LLM_DATASET.json into TinyBrain session"
```

**Load by Severity**:
```bash
# Critical severity CWE patterns
cline "Load critical CWE vulnerabilities (CWE-787, CWE-78, CWE-89, CWE-416, CWE-502) into TinyBrain"
```

**Load by Language**:
```bash
# Python-specific patterns
cline "Load Python-specific CWE patterns into TinyBrain for Python codebase assessment"

# JavaScript-specific patterns
cline "Load JavaScript CWE patterns into TinyBrain for web application assessment"
```

## Search Strategies

### By Severity (Priority-based)

**Critical Vulnerabilities (Priority 10)**:
```bash
cline "Search TinyBrain for CWE vulnerabilities with priority 10: CWE-787, CWE-78, CWE-89, CWE-416, CWE-502"
```

**High-Severity Vulnerabilities (Priority 8-9)**:
```bash
cline "Search TinyBrain for high-priority CWE patterns: CWE-79, CWE-20, CWE-125, CWE-190, CWE-352"
```

**Medium-Severity Vulnerabilities (Priority 6-7)**:
```bash
cline "Search TinyBrain for medium-priority CWE patterns: CWE-770, CWE-250, CWE-732, CWE-330"
```

### By Language (Target-specific)

**Python Codebase**:
```bash
cline "Search TinyBrain for Python CWE patterns: CWE-20, CWE-78, CWE-89, CWE-22, CWE-434, CWE-918"
```

**JavaScript/Node.js Codebase**:
```bash
cline "Search TinyBrain for JavaScript CWE patterns: CWE-79, CWE-352, CWE-918, CWE-434, CWE-807"
```

**Java Codebase**:
```bash
cline "Search TinyBrain for Java CWE patterns: CWE-79, CWE-20, CWE-78, CWE-89, CWE-362, CWE-502"
```

**C/C++ Codebase**:
```bash
cline "Search TinyBrain for C/C++ CWE patterns: CWE-787, CWE-125, CWE-416, CWE-190"
```

### By Category (Attack-type specific)

**Injection Attacks**:
```bash
cline "Search TinyBrain for injection CWE patterns: CWE-78 (Command), CWE-89 (SQL), CWE-74 (LDAP)"
```

**Memory Corruption**:
```bash
cline "Search TinyBrain for memory corruption CWE patterns: CWE-787, CWE-125, CWE-416, CWE-190"
```

**Web Security**:
```bash
cline "Search TinyBrain for web security CWE patterns: CWE-79 (XSS), CWE-352 (CSRF), CWE-918 (SSRF)"
```

**Authentication/Authorization**:
```bash
cline "Search TinyBrain for auth CWE patterns: CWE-306, CWE-862, CWE-863, CWE-639"
```

## Storing CWE Findings

### CWE Vulnerability Template

Store a CWE vulnerability finding:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "[CWE-89] SQL Injection in Login Form",
    "content": "AUTHORIZED SECURITY ASSESSMENT: CWE-89 SQL injection vulnerability identified in user login endpoint. User input concatenated directly into SQL query without parameterization. EXPLOITATION CONFIRMED: Authentication bypass successful with payload: ' OR 1=1--",
    "category": "vulnerability",
    "priority": 10,
    "confidence": 0.95,
    "tags": "[\"cwe-89\", \"sql-injection\", \"critical\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

### CWE Pattern Template

Store a CWE pattern finding:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "[CWE-79] XSS Pattern - Reflected XSS",
    "content": "AUTHORIZED SECURITY ASSESSMENT: CWE-79 reflected XSS pattern identified in search functionality. User input reflected in response without encoding. SECURITY IMPACT: Session hijacking possible via cookie theft.",
    "category": "technique",
    "priority": 9,
    "confidence": 0.9,
    "tags": "[\"cwe-79\", \"xss\", \"reflected\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

### CWE Exploit Example

Store a CWE exploitation finding:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "[CWE-78] Command Injection Exploit",
    "content": "AUTHORIZED SECURITY ASSESSMENT: CWE-78 command injection exploit confirmed in file processing module. Priority 10, affects file upload functionality. EXPLOITATION: Remote code execution achieved via payload: ; cat /etc/passwd",
    "category": "exploit",
    "priority": 10,
    "confidence": 1.0,
    "tags": "[\"cwe-78\", \"command-injection\", \"rce\", \"authorized-testing\"]",
    "source": "exploitation"
  }
}
```

## Standards Compliance Mapping

### Map CWE to OWASP Top 10

Link CWE findings to OWASP categories:

```bash
# CWE-89 (SQL Injection) → OWASP A03:2021 (Injection)
cline "Map CWE-89 (SQL injection) to OWASP A03:2021 (Injection) in TinyBrain"

# CWE-79 (XSS) → OWASP A03:2021 (Injection)
cline "Map CWE-79 (XSS) to OWASP A03:2021 (Injection) in TinyBrain"

# CWE-352 (CSRF) → OWASP A01:2021 (Broken Access Control)
cline "Map CWE-352 (CSRF) to OWASP A01:2021 (Broken Access Control) in TinyBrain"
```

### Map CWE to NIST Standards

```bash
# Map to NIST SP 800-115
cline "Map CWE findings to NIST SP 800-115 security testing controls in TinyBrain"
```

### Generate Compliance Reports

```bash
# CWE compliance report
cline "Generate CWE compliance report from TinyBrain data for security standards alignment"

# OWASP Top 10 coverage report
cline "Generate OWASP Top 10 coverage report from CWE findings in TinyBrain"
```

## Workflow Examples

### Complete CWE Assessment Workflow

**1. Initialize CWE Assessment**:
```bash
cline "Create TinyBrain session for CWE-based security assessment of target application"
```

**2. Load Relevant CWE Patterns**:
```bash
cline "Load critical and high-severity CWE patterns into TinyBrain session"
```

**3. Search for Vulnerabilities**:
```bash
cline "Search TinyBrain for CWE patterns matching target codebase characteristics"
```

**4. Store Findings**:
```bash
cline "Store CWE-89 SQL injection vulnerability in TinyBrain: Priority 10, affects user authentication"
```

**5. Create Relationships**:
```bash
cline "Create relationship in TinyBrain between CWE-89 finding and CWE-20 input validation pattern"
```

**6. Generate Report**:
```bash
cline "Generate CWE compliance report from TinyBrain data for security standards alignment"
```

### Language-Specific Assessment

**Python Security Assessment**:
```bash
# Load Python patterns
cline "Load Python-specific CWE patterns (CWE-20, CWE-78, CWE-89, CWE-22, CWE-434) into TinyBrain"

# Search for vulnerabilities
cline "Search TinyBrain for Python CWE vulnerabilities in target Python codebase"

# Store findings
cline "Store CWE-78 command injection finding in TinyBrain: Priority 10, affects file processing module"
```

**JavaScript Security Assessment**:
```bash
# Load JavaScript patterns
cline "Load JavaScript-specific CWE patterns (CWE-79, CWE-352, CWE-918, CWE-434) into TinyBrain"

# Search for vulnerabilities
cline "Search TinyBrain for JavaScript CWE vulnerabilities in target web application"

# Store findings
cline "Store CWE-79 XSS finding in TinyBrain: Priority 9, affects search functionality"
```

## Context Optimization

### Avoid Context Overload

**❌ DON'T**: Load entire dataset at once
```bash
cline "Load all 28 CWE patterns into TinyBrain"  # Too much context
```

**✅ DO**: Load specific patterns
```bash
cline "Load critical CWE patterns (CWE-787, CWE-78, CWE-89, CWE-416, CWE-502) into TinyBrain"
```

**✅ DO**: Load by assessment phase
```bash
cline "Load injection CWE patterns for current assessment phase"
```

### Efficient Pattern Matching

**Search for Specific Vulnerability Type**:
```bash
cline "Search TinyBrain for SQL injection patterns (CWE-89) in database interaction code"
```

**Correlate with Existing Findings**:
```bash
cline "Find CWE patterns similar to stored vulnerability ID 12345 in TinyBrain"
```

**Map to Security Standards**:
```bash
cline "Map CWE findings to OWASP Top 10 categories in TinyBrain"
```

## Best Practices

1. **Load Selectively**: Only load CWE patterns relevant to the assessment
2. **Use Priorities**: Map CWE severity to TinyBrain priority levels (1-10)
3. **Tag Consistently**: Always include CWE ID in tags (e.g., "cwe-89")
4. **Create Relationships**: Link CWE findings to related patterns and standards
5. **Generate Reports**: Use TinyBrain's export features for compliance reporting
6. **Update Regularly**: Review new CWE additions and update patterns

## Next Steps

- Review [Security Templates](security-templates.md) for more template examples
- See [Workflows](../workflows/security-assessment.md) for comprehensive assessment workflows
- Check [API Examples](../api/examples.md) for TinyBrain tool usage
