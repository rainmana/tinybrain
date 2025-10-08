---
layout: default
title: Security Templates
nav_order: 13
parent: Datasets
description: "Pre-configured security templates for TinyBrain"
---

# Security Templates

This document describes the pre-configured security templates available in TinyBrain for consistent security assessment storage, based on `TINYBRAIN_SECURITY_TEMPLATES.md`.

## Table of Contents
- [Overview](#overview)
- [Template Types](#template-types)
- [Usage Commands](#usage-commands)
- [OWASP Top 10 Templates](#owasp-top-10-templates)
- [CWE Pattern Templates](#cwe-pattern-templates)
- [Workflow Examples](#workflow-examples)

## Overview

TinyBrain Security Templates provide pre-configured memory templates for storing security findings consistently across assessments.

**Template Benefits**:
- Consistent categorization
- Standard field structure
- OWASP and CWE alignment
- Exploitation status tracking
- Relationship-ready design

## Template Types

### Vulnerability Finding Template

```json
{
  "title": "[VULNERABILITY_TYPE] in [LOCATION]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [VULNERABILITY_TYPE] vulnerability identified in [LOCATION]. [DETAILED_DESCRIPTION]. EXPLOITATION CONFIRMED: [EXPLOITATION_DETAILS].",
  "category": "vulnerability",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["[vulnerability-type]", "[severity]", "[location]", "authorized-testing"],
  "source": "authorized-security-assessment"
}
```

### Exploit Code Template

```json
{
  "title": "[EXPLOIT_TYPE] Exploit for [VULNERABILITY]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Proof-of-concept exploit for [VULNERABILITY_TYPE] in [LOCATION]. [EXPLOITATION_DESCRIPTION]. VALIDATION CONFIRMED: [VALIDATION_DETAILS].",
  "category": "exploit",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["[exploit-type]", "[vulnerability-type]", "proof-of-concept", "authorized-testing"],
  "source": "authorized-security-assessment"
}
```

### Security Pattern Template

```json
{
  "title": "[PATTERN_TYPE] Security Pattern - [PATTERN_NAME]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [PATTERN_TYPE] pattern identified in [LOCATION]. [PATTERN_DESCRIPTION]. SECURITY IMPACT: [IMPACT_ANALYSIS].",
  "category": "technique",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["[pattern-type]", "[security-category]", "code-pattern", "authorized-testing"],
  "source": "authorized-security-assessment"
}
```

### Remediation Recommendation Template

```json
{
  "title": "Remediation for [VULNERABILITY_TYPE] in [LOCATION]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Remediation recommendations for [VULNERABILITY_TYPE] in [LOCATION]. [REMEDIATION_DESCRIPTION]. IMPLEMENTATION GUIDANCE: [IMPLEMENTATION_STEPS].",
  "category": "recommendation",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["remediation", "[vulnerability-type]", "security-fix", "authorized-testing"],
  "source": "authorized-security-assessment"
}
```

## Usage Commands

### Storing Vulnerability Findings

**SQL Injection Vulnerability**:
```bash
cline "Store SQL injection vulnerability in TinyBrain: Priority 9, affects user authentication endpoint, exploit confirmed with authentication bypass"
```

**XSS Vulnerability**:
```bash
cline "Store stored XSS vulnerability in TinyBrain: Priority 8, affects admin panel, exploit confirmed with session hijacking"
```

**RCE Vulnerability**:
```bash
cline "Store remote code execution vulnerability in TinyBrain: Priority 10, affects main application server, exploit confirmed with system compromise"
```

### Storing Exploit Code

**SQL Injection Exploit**:
```bash
cline "Store SQL injection exploit code in TinyBrain with relationship to vulnerability ID 12345, includes authentication bypass payload"
```

**XSS Exploit**:
```bash
cline "Store XSS exploit code in TinyBrain with relationship to vulnerability ID 12346, includes session hijacking payload"
```

**Command Injection Exploit**:
```bash
cline "Store command injection exploit code in TinyBrain with relationship to vulnerability ID 12347, includes system command execution"
```

### Storing Security Patterns

**Vulnerability Pattern**:
```bash
cline "Store SQL injection pattern in TinyBrain: String concatenation vulnerability pattern, affects multiple endpoints"
```

**Code Smell**:
```bash
cline "Store hardcoded credentials code smell in TinyBrain: Default admin credentials found in configuration files"
```

**Security Anti-Pattern**:
```bash
cline "Store weak password validation anti-pattern in TinyBrain: Insufficient password complexity requirements"
```

### Creating Relationships

**Link Exploit to Vulnerability**:
```bash
cline "Create relationship in TinyBrain between exploit ID 67890 and vulnerability ID 12345 with type 'exploits'"
```

**Link Pattern to Vulnerability**:
```bash
cline "Create relationship in TinyBrain between pattern ID 54321 and vulnerability ID 12345 with type 'causes'"
```

**Link Recommendation to Vulnerability**:
```bash
cline "Create relationship in TinyBrain between recommendation ID 98765 and vulnerability ID 12345 with type 'mitigates'"
```

## OWASP Top 10 Templates

### A01: Broken Access Control

**IDOR Vulnerability Template**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "Broken Access Control - IDOR Vulnerability",
    "content": "AUTHORIZED SECURITY ASSESSMENT: Insecure Direct Object Reference (IDOR) vulnerability identified. Direct access to user data without proper authorization checks. EXPLOITATION CONFIRMED: Unauthorized access to other users' sensitive information.",
    "category": "vulnerability",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"broken-access-control\", \"idor\", \"authorization\", \"owasp-a01\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

### A02: Cryptographic Failures

**Weak Password Hashing Template**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "Cryptographic Failure - Weak Password Hashing",
    "content": "AUTHORIZED SECURITY ASSESSMENT: Weak password hashing identified using MD5 algorithm. Vulnerable to rainbow table attacks and brute force. EXPLOITATION CONFIRMED: Password hashes can be cracked using common tools.",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.95,
    "tags": "[\"cryptographic-failure\", \"weak-hashing\", \"md5\", \"owasp-a02\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

### A03: Injection

**SQL Injection Template**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "SQL Injection - Authentication Bypass",
    "content": "AUTHORIZED SECURITY ASSESSMENT: SQL injection vulnerability in login form allows authentication bypass. User input directly concatenated into SQL query. EXPLOITATION CONFIRMED: Successful authentication bypass and database access.",
    "category": "vulnerability",
    "priority": 10,
    "confidence": 0.98,
    "tags": "[\"sql-injection\", \"authentication-bypass\", \"owasp-a03\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

## CWE Pattern Templates

### CWE-79: Cross-site Scripting (XSS)

**Stored XSS Template**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "Cross-site Scripting (XSS) - Stored XSS",
    "content": "AUTHORIZED SECURITY ASSESSMENT: Stored Cross-Site Scripting (XSS) vulnerability in user comments. Malicious scripts stored in database and executed when viewed. EXPLOITATION CONFIRMED: Session hijacking and administrative access achieved.",
    "category": "vulnerability",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"xss\", \"stored-xss\", \"session-hijacking\", \"cwe-79\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

### CWE-22: Path Traversal

**File Access Vulnerability Template**:
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "Path Traversal - File Access Vulnerability",
    "content": "AUTHORIZED SECURITY ASSESSMENT: Path traversal vulnerability in file download functionality. Allows access to arbitrary files on the server. EXPLOITATION CONFIRMED: Successful access to system files including /etc/passwd.",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.95,
    "tags": "[\"path-traversal\", \"file-access\", \"cwe-22\", \"authorized-testing\"]",
    "source": "authorized-security-assessment"
  }
}
```

## Workflow Examples

### Complete Security Assessment Workflow

**1. Initialize Session**:
```bash
cline "Create TinyBrain session for web application security assessment of target.com"
```

**2. Store Vulnerability Findings**:
```bash
cline "Store SQL injection vulnerability in TinyBrain: Priority 10, affects login endpoint, exploit confirmed with authentication bypass"
```

**3. Store Exploit Code**:
```bash
cline "Store SQL injection exploit code in TinyBrain with relationship to vulnerability, includes authentication bypass payload"
```

**4. Store Security Patterns**:
```bash
cline "Store SQL injection pattern in TinyBrain: String concatenation vulnerability pattern affecting multiple endpoints"
```

**5. Store Recommendations**:
```bash
cline "Store remediation recommendation in TinyBrain: Use parameterized queries to prevent SQL injection"
```

**6. Create Relationships**:
```bash
cline "Create relationship in TinyBrain between exploit and vulnerability with type 'exploits'"
cline "Create relationship in TinyBrain between recommendation and vulnerability with type 'mitigates'"
```

**7. Generate Report**:
```bash
cline "Generate security report from TinyBrain session data for client presentation"
```

### Pattern Recognition Workflow

**1. Load Patterns**:
```bash
cline "Load OWASP Top 10 security patterns from TinyBrain templates"
```

**2. Search Code**:
```bash
cline "Search codebase for SQL injection patterns matching TinyBrain template signatures"
```

**3. Store Findings**:
```bash
cline "Store matched SQL injection pattern instances in TinyBrain with template metadata"
```

**4. Validate**:
```bash
cline "Validate stored patterns and update confidence scores based on exploitation results"
```

**5. Report**:
```bash
cline "Generate pattern recognition report from TinyBrain findings"
```

## Best Practices

1. **Use Consistent Templates**: Apply standard templates for all similar findings
2. **Include OWASP/CWE IDs**: Always tag with relevant standard identifiers
3. **Track Exploitation Status**: Mark findings as confirmed, suspected, or theoretical
4. **Create Relationships**: Link vulnerabilities to exploits and recommendations
5. **Update Confidence**: Adjust confidence scores as validation progresses
6. **Generate Reports**: Use templates for consistent reporting across assessments

## Next Steps

- Review [CWE Integration](cwe-integration.md) for CWE dataset usage
- See [Workflows](../workflows/security-assessment.md) for comprehensive examples
- Check [API Examples](../api/examples.md) for TinyBrain tool usage
