---
layout: default
title: AI Assistants
nav_order: 15
parent: Integrations
description: "AI assistant integration guide for TinyBrain"
---

# AI Assistant Integration

This guide provides integration instructions for using Cursor, Cline, and Roo AI assistants with TinyBrain for advanced security assessments.

## Table of Contents
- [Overview](#overview)
- [TinyBrain Integration (Universal)](#tinybrain-integration-universal)
- [Cursor Integration](#cursor-integration)
- [Cline Integration](#cline-integration)
- [Roo Integration](#roo-integration)
- [Multi-Assistant Workflow](#multi-assistant-workflow)

## Overview

TinyBrain integrates with multiple AI assistants, each with specialized configurations optimized for different aspects of security testing:

- **Cursor**: Unified security assessment framework, code review, reporting
- **Cline**: Advanced code analysis, static analysis, exploit development
- **Roo**: Advanced penetration testing, network security, post-exploitation

## TinyBrain Integration (Universal)

### Mandatory for ALL AI Assistants

All AI assistants must follow these TinyBrain integration requirements:

- **ALWAYS** use TinyBrain MCP server for memory management
- **STORE** all findings with proper categorization (vulnerability, exploit, technique, tool, evidence, recommendation)
- **USE** priority levels (1-10) and confidence scores (0.0-1.0)
- **CREATE** relationships between related vulnerabilities and exploits
- **TRACK** progress through assessment stages
- **EXPORT** session data for documentation

## Cursor Integration

### Configuration File

**File**: `.cursorrules`

**Purpose**: Unified security assessment framework

**Focus**: Code review, vulnerability analysis, exploit development

### Key Features

- **Security Assessment Rules**: Comprehensive framework for security testing
- **MCP Server Management**: Required and optional MCP server configurations
- **Quality Assurance**: Validation checklists and success metrics
- **Anti-Pattern Prevention**: Prevents LLM loops and inconsistent behavior

### Usage Commands

**Start Security Assessment**:
```bash
cursor "I'm starting an AUTHORIZED security assessment of target.com. Please create a TinyBrain session and begin vulnerability analysis."
```

**Store Findings**:
```bash
cursor "Store this SQL injection vulnerability in TinyBrain: Priority 9, affects login endpoint, exploit confirmed."
```

**Generate Security Report**:
```bash
cursor "Generate comprehensive security report from TinyBrain data for client presentation."
```

### Setup Steps

See the [Cursor Setup Guide](cursor.md) for detailed configuration instructions.

## Cline Integration

### Configuration File

**File**: `.clinerules`

**Purpose**: Advanced code review and exploitation framework

**Focus**: Static analysis, dynamic testing, exploit development

### Key Features

- **Advanced Exploitation Templates**: SQL injection, XSS, RCE, command injection
- **Security Code Review Framework**: Comprehensive vulnerability categories
- **Tool Integration**: Semgrep, CodeQL, Burp Suite, Metasploit
- **Workflow Automation**: Automated testing and exploitation scripts

### Usage Commands

**Initialize Security Assessment**:
```bash
cline "Create TinyBrain session for web application security assessment of target.com"
```

**Run Security Analysis**:
```bash
cline "Analyze the codebase for SQL injection vulnerabilities and store findings in TinyBrain"
```

**Develop Exploits**:
```bash
cline "Develop proof-of-concept exploits for high-priority vulnerabilities found in TinyBrain"
```

**Generate Reports**:
```bash
cline "Generate technical security report from TinyBrain findings for client documentation"
```

## Roo Integration

### Configuration File

**File**: `.roo-mode`

**Purpose**: Advanced penetration testing and vulnerability research

**Focus**: Network testing, social engineering, post-exploitation

### Key Features

- **Advanced Penetration Testing**: Network, web, social engineering
- **Exploitation Frameworks**: Custom exploitation scripts and tools
- **Persistence Techniques**: Advanced backdoors and persistence mechanisms
- **APT Simulation**: Advanced persistent threat techniques

### Usage Commands

**Start Penetration Test**:
```bash
roo "Initialize TinyBrain session for comprehensive penetration testing of target.com"
```

**Conduct Reconnaissance**:
```bash
roo "Perform OSINT reconnaissance and store findings in TinyBrain with proper categorization"
```

**Execute Exploitation**:
```bash
roo "Exploit high-priority vulnerabilities from TinyBrain and store results with relationships"
```

**Generate Penetration Test Report**:
```bash
roo "Generate executive summary and technical report from TinyBrain data for client presentation"
```

## Multi-Assistant Workflow

Combine multiple AI assistants for comprehensive security assessments.

### Phase 1: Assessment Planning (Cursor)

Use Cursor for initial planning and setup:

```bash
cursor "I'm starting an AUTHORIZED security assessment of target.com. Please:
1. Create TinyBrain session with name 'Target.com Security Assessment'
2. Set up task progress tracking for all phases
3. Get security templates relevant to web applications
4. Create initial context snapshot with assessment scope"
```

### Phase 2: Code Analysis (Cline)

Use Cline for detailed code analysis:

```bash
cline "Based on the TinyBrain session 'Target.com Security Assessment', please:
1. Perform static code analysis for security vulnerabilities
2. Store all findings in TinyBrain with proper categorization
3. Develop proof-of-concept exploits for critical findings
4. Update task progress to 50% - completed vulnerability discovery"
```

### Phase 3: Penetration Testing (Roo)

Use Roo for comprehensive penetration testing:

```bash
roo "Based on TinyBrain findings from 'Target.com Security Assessment', please:
1. Conduct dynamic testing and exploitation validation
2. Perform network penetration testing
3. Test social engineering vectors
4. Store all exploitation results in TinyBrain with relationships"
```

### Phase 4: Reporting (Cursor)

Use Cursor for final reporting and documentation:

```bash
cursor "Based on the complete TinyBrain session 'Target.com Security Assessment', please:
1. Generate executive summary for client presentation
2. Create detailed technical report with all findings
3. Export session data for client documentation
4. Update task progress to 100% - assessment completed"
```

## Specialized Use Cases

### Web Application Security Assessment

**Cursor (Planning & Reporting)**:
- Assessment scope definition
- Risk assessment and prioritization
- Executive summary generation
- Client communication

**Cline (Code Analysis)**:
- Static code analysis
- Vulnerability identification
- Exploit development
- Code review documentation

**Roo (Dynamic Testing)**:
- Web application penetration testing
- Authentication bypass testing
- Session management testing
- Business logic testing

### Network Penetration Testing

**Cursor (Coordination)**:
- Network scope definition
- Testing methodology planning
- Progress tracking
- Report coordination

**Cline (Tool Integration)**:
- Nmap scanning automation
- Vulnerability scanning
- Service enumeration
- Custom script development

**Roo (Exploitation)**:
- Network service exploitation
- Privilege escalation
- Lateral movement
- Post-exploitation activities

## Best Practices

1. **Use Session Continuity**: Share the same TinyBrain session across assistants
2. **Track Progress**: Update task progress when switching between assistants
3. **Create Snapshots**: Save context snapshots at phase transitions
4. **Maintain Relationships**: Link findings across different assessment phases
5. **Use Consistent Tags**: Apply consistent tagging across all assistants
6. **Generate Reports**: Consolidate findings from all assistants in final reports

## Common Workflow Patterns

### Pattern 1: Code-First Assessment

1. **Cline**: Analyze code, store findings
2. **Roo**: Validate findings with dynamic testing
3. **Cursor**: Generate reports and recommendations

### Pattern 2: Recon-First Assessment

1. **Roo**: Perform reconnaissance, store targets
2. **Cline**: Analyze discovered services
3. **Cursor**: Compile findings and prioritize

### Pattern 3: Continuous Assessment

1. **Cursor**: Initialize and monitor assessment
2. **Cline**: Continuous code analysis
3. **Roo**: Periodic penetration testing
4. **Cursor**: Regular progress reports

## Next Steps

- See [Cursor Setup Guide](cursor.md) for detailed Cursor configuration
- Review [Workflows](../workflows/security-assessment.md) for assessment patterns
- Check [API Examples](../api/examples.md) for TinyBrain tool usage
