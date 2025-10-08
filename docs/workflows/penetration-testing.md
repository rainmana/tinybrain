---
layout: default
title: Penetration Testing
nav_order: 9
parent: Workflows
description: "Penetration testing workflow using TinyBrain"
---

# Penetration Testing Workflow

This document describes the complete workflow for conducting penetration tests using TinyBrain, based on examples from the repository.

## Table of Contents
- [Overview](#overview)
- [Workflow Phases](#workflow-phases)
- [Session Setup](#session-setup)
- [Reconnaissance](#reconnaissance)
- [Vulnerability Discovery](#vulnerability-discovery)
- [Exploitation](#exploitation)
- [Post-Exploitation](#post-exploitation)
- [Reporting](#reporting)

## Overview

Penetration testing with TinyBrain follows the standard pen-test methodology with structured memory management for all findings, techniques, and evidence.

**Workflow Pattern**:
1. Create penetration_test session
2. Store reconnaissance findings
3. Document identified vulnerabilities
4. Store exploitation techniques and payloads
5. Track relationships between findings
6. Generate comprehensive reports

## Workflow Phases

### Phase 1: Reconnaissance
- Network scanning and enumeration
- Service identification
- Initial foothold discovery

### Phase 2: Vulnerability Discovery
- Identify exploitable weaknesses
- Assess attack surface
- Prioritize targets

### Phase 3: Exploitation
- Develop and execute exploits
- Document successful attacks
- Maintain access

### Phase 4: Post-Exploitation
- Privilege escalation
- Lateral movement
- Data exfiltration paths

### Phase 5: Reporting
- Compile findings
- Document attack paths
- Provide remediation guidance

## Session Setup

### Create Penetration Test Session

Initialize a new penetration testing session:

```json
{
  "name": "create_session",
  "arguments": {
    "name": "Internal Network Penetration Test",
    "description": "Internal network penetration test for corporate environment",
    "task_type": "penetration_test",
    "metadata": "{\"target\": \"192.168.1.0/24\", \"scope\": \"internal-network\", \"duration\": \"2-weeks\"}"
  }
}
```

**Key Fields**:
- `task_type`: Always `"penetration_test"` for pen-tests
- `metadata`: Store target ranges, scope boundaries, duration

## Reconnaissance

### Store Network Scan Results

Document discovered hosts and services:

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Active Directory Domain Controller",
    "content": "Discovered domain controller at 192.168.1.10 running Windows Server 2019. Ports 88, 389, 636, 3268, 3269 open. Kerberos and LDAP services active.",
    "category": "finding",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"ad\", \"domain-controller\", \"kerberos\", \"ldap\"]",
    "source": "nmap-scan"
  }
}
```

### Store Service Enumeration

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "SMB Service Enumeration",
    "content": "SMB signing disabled on file server at 192.168.1.20. Vulnerable to relay attacks. Anonymous enumeration enabled, revealing user lists.",
    "category": "finding",
    "priority": 7,
    "confidence": 0.85,
    "tags": "[\"smb\", \"relay\", \"anonymous-enum\", \"file-server\"]",
    "source": "enum4linux"
  }
}
```

### Store Web Application Discovery

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Internal Web Portal Discovery",
    "content": "Found internal web portal at https://portal.internal.corp on port 8443. Running Apache Tomcat 9.0.45. Login page vulnerable to username enumeration.",
    "category": "finding",
    "priority": 6,
    "confidence": 0.8,
    "tags": "[\"web-app\", \"tomcat\", \"username-enum\", \"portal\"]",
    "source": "web-discovery"
  }
}
```

## Vulnerability Discovery

### Store Identified Vulnerability

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Weak Kerberos Service Account Passwords",
    "content": "Multiple service accounts with weak passwords identified through Kerberoasting. Accounts: svc_backup, svc_sql, svc_web. Passwords crackable in under 2 hours.",
    "category": "vulnerability",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"kerberoasting\", \"weak-passwords\", \"service-accounts\"]",
    "source": "manual-analysis"
  }
}
```

### Store Privilege Escalation Path

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Unquoted Service Path Vulnerability",
    "content": "Windows service 'BackupService' running with unquoted path: C:\\Program Files\\Backup Service\\service.exe. Service runs as SYSTEM. Can be exploited for privilege escalation.",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.95,
    "tags": "[\"privilege-escalation\", \"unquoted-path\", \"windows\", \"system\"]",
    "source": "local-enumeration"
  }
}
```

## Exploitation

### Store Attack Technique

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Kerberoasting Attack Vector",
    "content": "Service accounts with weak passwords identified. Can perform Kerberoasting attack using GetUserSPNs.py and hashcat. Expected success rate: 60%",
    "category": "technique",
    "priority": 7,
    "confidence": 0.8,
    "tags": "[\"kerberoasting\", \"service-accounts\", \"password-attack\"]",
    "source": "manual-analysis"
  }
}
```

### Store Payload

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "PowerShell Reverse Shell Payload",
    "content": "Base64 encoded PowerShell reverse shell payload for Windows targets:\n\npowershell -e <base64-encoded-payload>\n\nListener: nc -lvp 4444",
    "category": "payload",
    "priority": 6,
    "confidence": 0.95,
    "tags": "[\"powershell\", \"reverse-shell\", \"windows\", \"base64\"]",
    "source": "custom-payload"
  }
}
```

### Store Successful Exploit

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Successful SMB Relay to Domain Controller",
    "content": "Successfully relayed SMB authentication to domain controller using Responder and ntlmrelayx. Gained SYSTEM access on DC. Dumped NTDS.dit for credential access.",
    "category": "exploit",
    "priority": 10,
    "confidence": 1.0,
    "tags": "[\"smb-relay\", \"domain-controller\", \"system\", \"credential-dump\"]",
    "source": "exploitation"
  }
}
```

## Post-Exploitation

### Store Credential Findings

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Domain Admin Credentials Obtained",
    "content": "Extracted domain admin credentials from NTDS.dit dump. Account: corp\\admin. Password hash cracked. Full domain compromise achieved.",
    "category": "evidence",
    "priority": 10,
    "confidence": 1.0,
    "tags": "[\"domain-admin\", \"credentials\", \"ntds\", \"full-compromise\"]",
    "source": "post-exploitation"
  }
}
```

### Store Lateral Movement Path

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_456",
    "title": "Lateral Movement via WMI",
    "content": "Used domain admin credentials to perform lateral movement via WMI. Compromised file server, database server, and backup server. Full network access achieved.",
    "category": "technique",
    "priority": 9,
    "confidence": 1.0,
    "tags": "[\"lateral-movement\", \"wmi\", \"domain-admin\", \"network-compromise\"]",
    "source": "post-exploitation"
  }
}
```

## Relationship Mapping

### Link Reconnaissance to Vulnerability

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "finding_smb_service",
    "target_memory_id": "vuln_smb_relay",
    "relationship_type": "causes",
    "strength": 0.9,
    "description": "SMB signing disabled enables relay attacks"
  }
}
```

### Link Vulnerability to Exploit

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "exploit_smb_relay",
    "target_memory_id": "vuln_smb_relay",
    "relationship_type": "exploits",
    "strength": 1.0,
    "description": "SMB relay exploit leverages disabled signing"
  }
}
```

### Link Exploit to Evidence

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "evidence_domain_admin",
    "target_memory_id": "exploit_smb_relay",
    "relationship_type": "supports",
    "strength": 1.0,
    "description": "Credentials obtained via successful SMB relay"
  }
}
```

### Map Attack Chain

```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "technique_lateral_movement",
    "target_memory_id": "evidence_domain_admin",
    "relationship_type": "depends_on",
    "strength": 1.0,
    "description": "Lateral movement requires domain admin credentials"
  }
}
```

## Progress Tracking

### Track Penetration Test Phases

```json
{
  "name": "create_task_progress",
  "arguments": {
    "session_id": "session_456",
    "task_name": "Network Penetration Test",
    "stage": "Reconnaissance",
    "status": "in_progress",
    "progress_percentage": 20,
    "notes": "Completed network scanning. Identified 45 hosts, 12 critical services."
  }
}
```

### Update for Exploitation Phase

```json
{
  "name": "update_task_progress",
  "arguments": {
    "session_id": "session_456",
    "task_name": "Network Penetration Test",
    "stage": "Exploitation",
    "status": "in_progress",
    "progress_percentage": 60,
    "notes": "Successful SMB relay. Domain admin access achieved. Beginning lateral movement."
  }
}
```

## Reporting

### Search for Critical Findings

```json
{
  "name": "search_memories",
  "arguments": {
    "query": "critical vulnerability exploit",
    "session_id": "session_456",
    "min_priority": 8,
    "limit": 50
  }
}
```

### Get Context Summary for Executive Report

```json
{
  "name": "get_context_summary",
  "arguments": {
    "session_id": "session_456",
    "current_task": "Preparing executive summary of penetration test findings",
    "max_memories": 20
  }
}
```

### Get Related Attack Chain

Trace the complete attack path:

```json
{
  "name": "get_related_memories",
  "arguments": {
    "memory_id": "exploit_smb_relay",
    "limit": 10
  }
}
```

This retrieves all related memories showing the complete attack chain from reconnaissance to domain compromise.

## Best Practices

Based on repository examples:

1. **Document Everything**: Store all findings, even unsuccessful attempts
2. **Use Proper Categories**: finding (recon), vulnerability (weaknesses), technique (methods), exploit (successful attacks), evidence (proof)
3. **Priority Assignment**: 10 for domain compromise, 8-9 for privilege escalation, 5-7 for information disclosure
4. **Create Attack Chains**: Link recon → vulnerability → exploit → evidence
5. **Track Progress**: Update task progress through pen-test phases
6. **Tag Consistently**: Use standard tags (smb-relay, kerberoasting, privilege-escalation)
7. **Store Payloads**: Document working payloads for future reference

## Example: Complete Pen-Test Workflow

### 1. Initialize
```json
{"name": "create_session", "arguments": {"task_type": "penetration_test"}}
```

### 2. Recon
```json
{"name": "store_memory", "arguments": {"category": "finding", "title": "Network Scan Results"}}
```

### 3. Vulnerability
```json
{"name": "store_memory", "arguments": {"category": "vulnerability", "title": "SMB Relay Vulnerability"}}
```

### 4. Exploitation
```json
{"name": "store_memory", "arguments": {"category": "exploit", "title": "Successful SMB Relay"}}
```

### 5. Evidence
```json
{"name": "store_memory", "arguments": {"category": "evidence", "title": "Domain Admin Access"}}
```

### 6. Relationships
```json
{"name": "create_relationship", "arguments": {"relationship_type": "exploits"}}
```

### 7. Progress
```json
{"name": "update_task_progress", "arguments": {"status": "completed"}}
```

## Next Steps

- See [Security Assessment Workflow](security-assessment.md) for code review patterns
- Review [Exploit Development Workflow](exploit-development.md) for exploit development
- Check [API Examples](../api/examples.md) for more JSON examples
