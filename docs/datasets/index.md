---
layout: default
title: Datasets
nav_order: 11
description: "Security datasets and templates available in TinyBrain"
has_children: true
---

# Security Datasets & Templates

TinyBrain includes comprehensive security datasets and templates aligned with industry standards for security assessments.

## Overview

The repository provides multiple dataset files containing security patterns, vulnerability examples, and assessment templates. All content is aligned with OWASP, CWE, NIST, and ISO security standards.

## Available Datasets

### Security Code Review Dataset
**File**: `SECURITY_CODE_REVIEW_DATASET.md`

OWASP Top 10 2021 patterns, CWE vulnerabilities, and exploitation techniques for security code reviews.

**Contents**:
- OWASP Top 10 2021 vulnerability patterns
- Code smell detection patterns
- Exploitation technique examples
- Language-specific vulnerability patterns

### Multi-Language Security Patterns
**File**: `MULTI_LANGUAGE_SECURITY_PATTERNS.md`

Language-specific vulnerability patterns covering 10 major programming languages.

**Covered Languages**:
- JavaScript/Node.js
- Python
- Java
- C#/.NET
- PHP
- Ruby
- Go
- C/C++
- TypeScript
- Rust

### CWE Security Patterns
**File**: `CWE_SECURITY_PATTERNS.md`

CWE Top 25 Most Dangerous Software Errors with comprehensive vulnerability patterns.

**Contents**:
- CWE Top 25 vulnerability patterns
- Severity classifications
- Language-specific examples
- Exploitation techniques

### CWE LLM Dataset
**File**: `CWE_LLM_DATASET.json`

LLM-optimized CWE dataset in structured JSON format for efficient consumption.

**Features**:
- Structured JSON format for LLM parsing
- Searchable by severity, language, and category
- Context-optimized for AI assistants
- Ready for TinyBrain integration

### CWE TinyBrain Integration
**File**: `CWE_TINYBRAIN_INTEGRATION.md`

Integration guide for using the CWE dataset with TinyBrain.

**See**: [CWE Integration Guide](cwe-integration.md)

### TinyBrain Security Templates
**File**: `TINYBRAIN_SECURITY_TEMPLATES.md`

Pre-configured memory templates for consistent security assessment storage.

**See**: [Security Templates Guide](security-templates.md)

### Enhanced Language Library Patterns
**File**: `ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md`

Framework and library-specific security patterns for popular web frameworks and libraries.

**Contents**:
- Framework-specific vulnerability patterns
- Library security best practices
- Common misconfiguration examples

## Standards Compliance

All datasets are aligned with industry-standard security frameworks:

### OWASP Standards
- **[OWASP Top 10 2021](https://owasp.org/Top10/)** - Web Application Security Risks
- **[OWASP Code Review Guide](https://github.com/OWASP/www-project-code-review-guide)** - Secure code review methodology
- **[OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)** - Security testing methodology

### CWE Standards
- **[CWE (Common Weakness Enumeration)](https://cwe.mitre.org/)** - Software Weakness Classification
- **[SANS Top 25 CWE](https://cwe.mitre.org/top25/)** - Most dangerous software errors

### NIST Standards
- **[NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)** - Technical Guide to Information Security Testing
- **[NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)** - Cybersecurity risk management

### ISO Standards
- **[ISO 27001](https://www.iso.org/isoiec-27001-information-security.html)** - Information Security Management Systems

### Penetration Testing Standards
- **[PTES (Penetration Testing Execution Standard)](http://www.pentest-standard.org/)** - Penetration Testing Methodology

## Source Attribution

Security patterns and vulnerability datasets are based on authoritative sources:

- **[OWASP Secure Coding Dojo](https://owasp.org/SecureCodingDojo/codereview101/)** - Interactive security code review training
- **OWASP Code Review Guide** - Comprehensive secure code review methodology
- **OWASP Testing Guide** - Web application security testing
- **CWE/SANS Top 25** - Most dangerous software errors
- **NIST publications** - Government security standards

## Multi-Language Coverage

The datasets provide comprehensive coverage across 10 major programming languages:

### Web Languages
- **JavaScript/Node.js** - Web application security patterns
- **TypeScript** - Type-safe web application patterns
- **PHP** - Web application security patterns
- **Ruby** - Web framework security patterns

### Enterprise Languages
- **Java** - Enterprise application security patterns
- **C#/.NET** - Microsoft ecosystem security patterns

### System Languages
- **C/C++** - System-level security patterns
- **Go** - System and API security patterns
- **Rust** - Memory-safe system programming patterns

### Scripting Languages
- **Python** - Backend and automation security patterns

## Usage with TinyBrain

### Loading Datasets

See the [CWE Integration Guide](cwe-integration.md) for commands to load datasets into TinyBrain sessions.

### Using Templates

See the [Security Templates Guide](security-templates.md) for template usage and example commands.

### Integration Workflows

Datasets can be integrated into TinyBrain workflows:

1. **Code Review**: Use OWASP patterns for systematic code review
2. **Vulnerability Assessment**: Reference CWE patterns for weakness identification
3. **Penetration Testing**: Apply exploitation techniques from datasets
4. **Compliance**: Map findings to frameworks (OWASP, CWE, NIST)

## Best Practices

1. **Use Standard Categories**: Follow dataset categorization for consistency
2. **Reference Standards**: Link findings to CWE IDs and OWASP categories
3. **Apply Templates**: Use security templates for uniform documentation
4. **Cross-Reference**: Link related patterns across datasets
5. **Update Regularly**: Review new dataset additions and updates

## Next Steps

- Explore [CWE Integration](cwe-integration.md) for dataset loading commands
- Review [Security Templates](security-templates.md) for template usage
- See [Workflows](../workflows/security-assessment.md) for practical examples
