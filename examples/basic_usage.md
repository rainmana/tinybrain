# TinyBrain Basic Usage Examples

This document provides practical examples of using TinyBrain for security-focused tasks.

## 1. Security Code Review Session

### Create a Security Review Session

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

### Store Security Findings

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

### Create Exploit Chain Relationship

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

## 2. Penetration Testing Session

### Create Penetration Test Session

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

### Store Reconnaissance Data

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

### Store Exploit Techniques

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

### Store Payload Information

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

## 3. Exploit Development Session

### Create Exploit Development Session

```json
{
  "name": "create_session",
  "arguments": {
    "name": "Buffer Overflow Exploit Development",
    "description": "Developing exploit for buffer overflow vulnerability in custom application",
    "task_type": "exploit_dev",
    "metadata": "{\"target\": \"custom-app-v1.2\", \"vulnerability\": \"buffer-overflow\", \"platform\": \"linux-x86\"}"
  }
}
```

### Store Vulnerability Analysis

```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_789",
    "title": "Buffer Overflow in strcpy()",
    "content": "Buffer overflow occurs in handle_request() function when copying user input to fixed-size buffer using strcpy(). Buffer size: 256 bytes. No stack canary or ASLR protection.",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.9,
    "tags": "[\"buffer-overflow\", \"strcpy\", \"no-protection\", \"code-execution\"]",
    "source": "reverse-engineering"
  }
}
```

### Store Exploit Development Progress

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

## 4. Search and Context Examples

### Search for High-Priority Vulnerabilities

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

### Get Context for Current Task

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

### Find Related Exploits

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

## 5. Task Progress Tracking

### Update Task Progress

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

## 6. Advanced Search Examples

### Search by Category and Tags

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

### Fuzzy Search for Similar Content

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

## 7. Database Management

### Check Database Health

```json
{
  "name": "health_check",
  "arguments": {}
}
```

### Get Database Statistics

```json
{
  "name": "get_database_stats",
  "arguments": {}
}
```

## Best Practices

1. **Use Descriptive Titles**: Make memory titles clear and searchable
2. **Set Appropriate Priorities**: Use 8-10 for critical vulnerabilities, 5-7 for medium, 1-4 for low
3. **Tag Consistently**: Use consistent tagging conventions across sessions
4. **Create Relationships**: Link related vulnerabilities and exploits
5. **Update Progress**: Track task progress for long-running security assessments
6. **Use Context Summaries**: Get relevant memories before starting new tasks
7. **Regular Searches**: Search for similar issues to avoid duplication

## Common Workflows

### 1. Vulnerability Assessment Workflow
1. Create security review session
2. Store each vulnerability finding
3. Create relationships between related vulnerabilities
4. Search for similar issues
5. Generate context summary for report writing

### 2. Penetration Testing Workflow
1. Create penetration test session
2. Store reconnaissance findings
3. Store exploit techniques and payloads
4. Track progress through different phases
5. Link exploits to vulnerabilities

### 3. Exploit Development Workflow
1. Create exploit development session
2. Store vulnerability analysis
3. Store exploit development progress
4. Create relationships between techniques
5. Reference previous similar exploits
