# 🧠 TinyBrain Security Templates
# Pre-configured templates for storing security findings and patterns

## 🚨 AUTHORIZATION & COMPLIANCE STATEMENT

**LEGAL AUTHORIZATION**: These templates are designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS** under explicit written authorization from authorized stakeholders.

## 🎯 TinyBrain Memory Templates

### **Vulnerability Finding Template:**
```json
{
  "title": "[VULNERABILITY_TYPE] in [LOCATION]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [VULNERABILITY_TYPE] vulnerability identified in [LOCATION]. [DETAILED_DESCRIPTION]. EXPLOITATION CONFIRMED: [EXPLOITATION_DETAILS].",
  "category": "vulnerability",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[vulnerability-type]", "[severity]", "[location]", "authorized-testing"],
  "source": "authorized-security-assessment",
  "exploitation_status": "confirmed",
  "business_impact": "[critical|high|medium|low]",
  "cwe_id": "[CWE-XXX]",
  "owasp_category": "[A0X:2021]",
  "cvss_score": "[X.X]"
}
```

### **Exploit Code Template:**
```json
{
  "title": "[EXPLOIT_TYPE] Exploit for [VULNERABILITY]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Proof-of-concept exploit for [VULNERABILITY_TYPE] in [LOCATION]. [EXPLOITATION_DESCRIPTION]. VALIDATION CONFIRMED: [VALIDATION_DETAILS].",
  "category": "exploit",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[exploit-type]", "[vulnerability-type]", "proof-of-concept", "authorized-testing"],
  "source": "authorized-security-assessment",
  "exploit_code": "[EXPLOIT_CODE_HERE]",
  "target_vulnerability": "[VULNERABILITY_ID]",
  "exploitation_method": "[METHOD_DESCRIPTION]"
}
```

### **Security Pattern Template:**
```json
{
  "title": "[PATTERN_TYPE] Security Pattern - [PATTERN_NAME]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [PATTERN_TYPE] pattern identified in [LOCATION]. [PATTERN_DESCRIPTION]. SECURITY IMPACT: [IMPACT_ANALYSIS].",
  "category": "technique",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[pattern-type]", "[security-category]", "code-pattern", "authorized-testing"],
  "source": "authorized-security-assessment",
  "pattern_code": "[PATTERN_CODE_HERE]",
  "security_implications": "[IMPLICATIONS]",
  "remediation_guidance": "[REMEDIATION_STEPS]"
}
```

### **Code Smell Template:**
```json
{
  "title": "[SMELL_TYPE] Code Smell in [LOCATION]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [SMELL_TYPE] code smell identified in [LOCATION]. [SMELL_DESCRIPTION]. SECURITY RISK: [RISK_ANALYSIS].",
  "category": "finding",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[smell-type]", "code-quality", "security-risk", "authorized-testing"],
  "source": "authorized-security-assessment",
  "smell_code": "[SMELL_CODE_HERE]",
  "security_risk": "[RISK_LEVEL]",
  "improvement_suggestions": "[SUGGESTIONS]"
}
```

### **Remediation Recommendation Template:**
```json
{
  "title": "Remediation for [VULNERABILITY_TYPE] in [LOCATION]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Remediation recommendations for [VULNERABILITY_TYPE] in [LOCATION]. [REMEDIATION_DESCRIPTION]. IMPLEMENTATION GUIDANCE: [IMPLEMENTATION_STEPS].",
  "category": "recommendation",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["remediation", "[vulnerability-type]", "security-fix", "authorized-testing"],
  "source": "authorized-security-assessment",
  "target_vulnerability": "[VULNERABILITY_ID]",
  "remediation_code": "[FIXED_CODE_HERE]",
  "implementation_steps": "[STEP_BY_STEP_GUIDANCE]"
}
```

## 🔧 TinyBrain Integration Commands

### **Storing Vulnerability Findings:**
```bash
# Store SQL injection vulnerability
cline "Store SQL injection vulnerability in TinyBrain: Priority 9, affects user authentication endpoint, exploit confirmed with authentication bypass"

# Store XSS vulnerability  
cline "Store stored XSS vulnerability in TinyBrain: Priority 8, affects admin panel, exploit confirmed with session hijacking"

# Store RCE vulnerability
cline "Store remote code execution vulnerability in TinyBrain: Priority 10, affects main application server, exploit confirmed with system compromise"
```

### **Storing Exploit Code:**
```bash
# Store SQL injection exploit
cline "Store SQL injection exploit code in TinyBrain with relationship to vulnerability ID 12345, includes authentication bypass payload"

# Store XSS exploit
cline "Store XSS exploit code in TinyBrain with relationship to vulnerability ID 12346, includes session hijacking payload"

# Store command injection exploit
cline "Store command injection exploit code in TinyBrain with relationship to vulnerability ID 12347, includes system command execution"
```

### **Storing Security Patterns:**
```bash
# Store vulnerability pattern
cline "Store SQL injection pattern in TinyBrain: String concatenation vulnerability pattern, affects multiple endpoints"

# Store code smell
cline "Store hardcoded credentials code smell in TinyBrain: Default admin credentials found in configuration files"

# Store security anti-pattern
cline "Store weak password validation anti-pattern in TinyBrain: Insufficient password complexity requirements"
```

### **Creating Relationships:**
```bash
# Link exploit to vulnerability
cline "Create relationship in TinyBrain between exploit ID 67890 and vulnerability ID 12345 with type 'exploits'"

# Link pattern to vulnerability
cline "Create relationship in TinyBrain between pattern ID 54321 and vulnerability ID 12345 with type 'causes'"

# Link recommendation to vulnerability
cline "Create relationship in TinyBrain between recommendation ID 98765 and vulnerability ID 12345 with type 'mitigates'"
```

## 🎯 Pre-configured Security Templates

### **OWASP Top 10 Templates:**

#### **A01: Broken Access Control**
```json
{
  "title": "Broken Access Control - IDOR Vulnerability",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Insecure Direct Object Reference (IDOR) vulnerability identified. Direct access to user data without proper authorization checks. EXPLOITATION CONFIRMED: Unauthorized access to other users' sensitive information.",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9,
  "tags": ["broken-access-control", "idor", "authorization", "owasp-a01", "authorized-testing"],
  "source": "authorized-security-assessment",
  "owasp_category": "A01:2021",
  "cwe_id": "CWE-639",
  "exploitation_status": "confirmed"
}
```

#### **A02: Cryptographic Failures**
```json
{
  "title": "Cryptographic Failure - Weak Password Hashing",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Weak password hashing identified using MD5 algorithm. Vulnerable to rainbow table attacks and brute force. EXPLOITATION CONFIRMED: Password hashes can be cracked using common tools.",
  "category": "vulnerability",
  "priority": 9,
  "confidence": 0.95,
  "tags": ["cryptographic-failure", "weak-hashing", "md5", "owasp-a02", "authorized-testing"],
  "source": "authorized-security-assessment",
  "owasp_category": "A02:2021",
  "cwe_id": "CWE-327",
  "exploitation_status": "confirmed"
}
```

#### **A03: Injection**
```json
{
  "title": "SQL Injection - Authentication Bypass",
  "content": "AUTHORIZED SECURITY ASSESSMENT: SQL injection vulnerability in login form allows authentication bypass. User input directly concatenated into SQL query. EXPLOITATION CONFIRMED: Successful authentication bypass and database access.",
  "category": "vulnerability",
  "priority": 10,
  "confidence": 0.98,
  "tags": ["sql-injection", "authentication-bypass", "owasp-a03", "authorized-testing"],
  "source": "authorized-security-assessment",
  "owasp_category": "A03:2021",
  "cwe_id": "CWE-89",
  "exploitation_status": "confirmed"
}
```

### **CWE Pattern Templates:**

#### **CWE-79: Cross-site Scripting (XSS)**
```json
{
  "title": "Cross-site Scripting (XSS) - Stored XSS",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Stored Cross-Site Scripting (XSS) vulnerability in user comments. Malicious scripts stored in database and executed when viewed. EXPLOITATION CONFIRMED: Session hijacking and administrative access achieved.",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9,
  "tags": ["xss", "stored-xss", "session-hijacking", "cwe-79", "authorized-testing"],
  "source": "authorized-security-assessment",
  "cwe_id": "CWE-79",
  "exploitation_status": "confirmed"
}
```

#### **CWE-22: Path Traversal**
```json
{
  "title": "Path Traversal - File Access Vulnerability",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Path traversal vulnerability in file download functionality. Allows access to arbitrary files on the server. EXPLOITATION CONFIRMED: Successful access to system files including /etc/passwd.",
  "category": "vulnerability",
  "priority": 9,
  "confidence": 0.95,
  "tags": ["path-traversal", "file-access", "cwe-22", "authorized-testing"],
  "source": "authorized-security-assessment",
  "cwe_id": "CWE-22",
  "exploitation_status": "confirmed"
}
```

## 🚀 Usage Examples

### **Complete Security Assessment Workflow:**

```bash
# 1. Initialize session
cline "Create TinyBrain session for web application security assessment of target.com"

# 2. Store vulnerability findings
cline "Store SQL injection vulnerability in TinyBrain: Priority 10, affects login endpoint, exploit confirmed with authentication bypass"

# 3. Store exploit code
cline "Store SQL injection exploit code in TinyBrain with relationship to vulnerability, includes authentication bypass payload"

# 4. Store security patterns
cline "Store SQL injection pattern in TinyBrain: String concatenation vulnerability pattern affecting multiple endpoints"

# 5. Store recommendations
cline "Store remediation recommendation in TinyBrain: Use parameterized queries to prevent SQL injection"

# 6. Create relationships
cline "Create relationship in TinyBrain between exploit and vulnerability with type 'exploits'"
cline "Create relationship in TinyBrain between recommendation and vulnerability with type 'mitigates'"

# 7. Generate report
cline "Generate comprehensive security report from TinyBrain data for client presentation"
```

### **Pattern Recognition Workflow:**

```bash
# 1. Search for similar patterns
cline "Search TinyBrain for similar SQL injection patterns in other projects"

# 2. Store new pattern
cline "Store new SQL injection pattern in TinyBrain: Dynamic query construction in search functionality"

# 3. Link to existing patterns
cline "Create relationship in TinyBrain between new pattern and existing SQL injection patterns with type 'related_to'"

# 4. Update recommendations
cline "Update remediation recommendations in TinyBrain based on pattern analysis"
```

## 🎯 Integration with Security Dataset

### **Loading Security Patterns:**
```bash
# Load OWASP Top 10 patterns
cline "Load OWASP Top 10 vulnerability patterns into TinyBrain session for reference"

# Load CWE patterns
cline "Load CWE vulnerability patterns into TinyBrain session for pattern matching"

# Load exploitation techniques
cline "Load exploitation techniques into TinyBrain session for proof-of-concept development"
```

### **Pattern Matching:**
```bash
# Match code against patterns
cline "Match target code against stored vulnerability patterns in TinyBrain"

# Find similar vulnerabilities
cline "Find similar vulnerabilities in TinyBrain based on code patterns and characteristics"

# Generate exploit suggestions
cline "Generate exploit suggestions based on stored exploitation techniques in TinyBrain"
```

## 🚨 Remember

These templates are designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS**. All templates include proper authorization language and are intended for legitimate security testing activities with proper authorization and compliance.

**Use responsibly and ethically!** 🛡️
