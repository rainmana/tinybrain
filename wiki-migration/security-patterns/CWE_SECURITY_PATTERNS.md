# 🛡️ CWE Security Patterns Dataset
# Comprehensive Common Weakness Enumeration (CWE) vulnerability patterns

## 🚨 AUTHORIZATION & COMPLIANCE STATEMENT

**LEGAL AUTHORIZATION**: This dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS** under explicit written authorization from authorized stakeholders.

## 🎯 CWE Top 25 Most Dangerous Software Errors (2024)

### **CWE-787: Out-of-bounds Write**
**Description**: Writing past the end of allocated memory buffer
**Severity**: Critical
**Languages**: C, C++, Assembly

```c
// VULNERABLE: Buffer overflow
void vulnerable_function(char *input) {
    char buffer[100];
    strcpy(buffer, input); // No bounds checking
}

// SECURE: Bounds checking
void secure_function(char *input, size_t len) {
    char buffer[100];
    size_t copy_len = (len < sizeof(buffer)) ? len : sizeof(buffer) - 1;
    strncpy(buffer, input, copy_len);
    buffer[copy_len] = '\0';
}
```

### **CWE-79: Cross-site Scripting (XSS)**
**Description**: Improper neutralization of input during web page generation
**Severity**: High
**Languages**: JavaScript, HTML, PHP, Java, C#

```javascript
// VULNERABLE: Reflected XSS
app.get('/search', (req, res) => {
    const query = req.query.q;
    res.send(`<h1>Search results for: ${query}</h1>`); // No encoding
});

// SECURE: Proper encoding
app.get('/search', (req, res) => {
    const query = req.query.q;
    const encodedQuery = escapeHtml(query);
    res.send(`<h1>Search results for: ${encodedQuery}</h1>`);
});
```

### **CWE-20: Improper Input Validation**
**Description**: Product does not validate or incorrectly validates input
**Severity**: High
**Languages**: All

```python
# VULNERABLE: No input validation
def process_user_data(data):
    return database.save(data) # No validation

# SECURE: Input validation
def process_user_data_secure(data):
    if not validate_input(data):
        raise ValueError("Invalid input")
    return database.save(sanitize_input(data))
```

### **CWE-125: Out-of-bounds Read**
**Description**: Reading past the end of allocated memory buffer
**Severity**: High
**Languages**: C, C++, Assembly

```c
// VULNERABLE: Buffer overread
void vulnerable_read(char *buffer, int index) {
    return buffer[index]; // No bounds checking
}

// SECURE: Bounds checking
void secure_read(char *buffer, int index, int size) {
    if (index >= 0 && index < size) {
        return buffer[index];
    }
    return -1; // Error handling
}
```

### **CWE-78: OS Command Injection**
**Description**: Improper neutralization of special elements in OS command
**Severity**: Critical
**Languages**: All

```python
# VULNERABLE: Command injection
import os
def process_file(filename):
    os.system(f"cat {filename}") # Command injection

# SECURE: No shell, direct command
import subprocess
def process_file_secure(filename):
    subprocess.run(['cat', filename], check=True)
```

### **CWE-89: SQL Injection**
**Description**: Improper neutralization of special elements in SQL command
**Severity**: Critical
**Languages**: All with SQL

```python
# VULNERABLE: String concatenation
def get_user(user_id):
    query = f"SELECT * FROM users WHERE id = {user_id}"
    return database.execute(query)

# SECURE: Parameterized queries
def get_user_secure(user_id):
    query = "SELECT * FROM users WHERE id = ?"
    return database.execute(query, (user_id,))
```

### **CWE-416: Use After Free**
**Description**: Referencing memory after it has been freed
**Severity**: Critical
**Languages**: C, C++

```c
// VULNERABLE: Use after free
void use_after_free() {
    char *ptr = malloc(100);
    free(ptr);
    *ptr = 'A'; // Use after free
}

// SECURE: Set pointer to NULL
void secure_use_after_free() {
    char *ptr = malloc(100);
    free(ptr);
    ptr = NULL; // Prevent use after free
}
```

### **CWE-190: Integer Overflow or Wraparound**
**Description**: Integer operation causes wraparound
**Severity**: High
**Languages**: C, C++, Java, C#

```c
// VULNERABLE: Integer overflow
int vulnerable_add(int a, int b) {
    return a + b; // Can overflow
}

// SECURE: Overflow checking
int secure_add(int a, int b) {
    if (a > INT_MAX - b) {
        return -1; // Overflow detected
    }
    return a + b;
}
```

### **CWE-352: Cross-Site Request Forgery (CSRF)**
**Description**: Web application does not verify request source
**Severity**: High
**Languages**: Web applications

```html
<!-- VULNERABLE: No CSRF protection -->
<form action="/transfer" method="POST">
    <input type="hidden" name="amount" value="1000">
    <input type="submit" value="Transfer">
</form>

<!-- SECURE: CSRF token -->
<form action="/transfer" method="POST">
    <input type="hidden" name="csrf_token" value="random_token">
    <input type="hidden" name="amount" value="1000">
    <input type="submit" value="Transfer">
</form>
```

### **CWE-22: Path Traversal**
**Description**: Improper limitation of pathname to restricted directory
**Severity**: High
**Languages**: All

```python
# VULNERABLE: Path traversal
def read_file(filename):
    with open(f"uploads/{filename}", 'r') as f: # No validation
        return f.read()

# SECURE: Path validation
def read_file_secure(filename):
    safe_path = os.path.join("uploads", os.path.basename(filename))
    if not safe_path.startswith("uploads/"):
        raise ValueError("Invalid path")
    with open(safe_path, 'r') as f:
        return f.read()
```

### **CWE-494: Download of Code Without Integrity Check**
**Description**: Product downloads code without verifying integrity
**Severity**: High
**Languages**: All

```python
# VULNERABLE: No integrity check
def download_plugin(url):
    response = requests.get(url)
    with open("plugin.py", "wb") as f:
        f.write(response.content) # No verification

# SECURE: Integrity verification
def download_plugin_secure(url, expected_hash):
    response = requests.get(url)
    actual_hash = hashlib.sha256(response.content).hexdigest()
    if actual_hash != expected_hash:
        raise ValueError("Integrity check failed")
    with open("plugin.py", "wb") as f:
        f.write(response.content)
```

### **CWE-362: Concurrent Execution using Shared Resource with Improper Synchronization**
**Description**: Race condition in shared resource access
**Severity**: High
**Languages**: Multi-threaded applications

```java
// VULNERABLE: Race condition
public class BankAccount {
    private int balance = 100;
    
    public void withdraw(int amount) {
        if (balance >= amount) {
            balance -= amount; // Race condition
        }
    }
}

// SECURE: Synchronization
public class BankAccountSecure {
    private int balance = 100;
    private final Object lock = new Object();
    
    public void withdraw(int amount) {
        synchronized(lock) {
            if (balance >= amount) {
                balance -= amount;
            }
        }
    }
}
```

### **CWE-770: Allocation of Resources Without Limits or Throttling**
**Description**: No limit on resource allocation
**Severity**: Medium
**Languages**: All

```python
# VULNERABLE: No resource limits
def process_requests():
    while True:
        request = get_request()
        process_request(request) # No limits

# SECURE: Resource throttling
def process_requests_secure():
    semaphore = threading.Semaphore(10) # Limit to 10 concurrent
    while True:
        semaphore.acquire()
        request = get_request()
        threading.Thread(target=process_request, args=(request,)).start()
```

### **CWE-918: Server-Side Request Forgery (SSRF)**
**Description**: Web server makes requests to arbitrary URLs
**Severity**: High
**Languages**: Web applications

```python
# VULNERABLE: SSRF
import requests
def fetch_url(url):
    response = requests.get(url) # No validation
    return response.content

# SECURE: URL validation
def fetch_url_secure(url):
    if not is_allowed_url(url):
        raise ValueError("URL not allowed")
    response = requests.get(url)
    return response.content
```

### **CWE-311: Missing Encryption of Sensitive Data**
**Description**: Sensitive data not encrypted
**Severity**: High
**Languages**: All

```python
# VULNERABLE: No encryption
def store_password(password):
    database.save("password", password) # Plain text

# SECURE: Encryption
def store_password_secure(password):
    hashed = bcrypt.hashpw(password.encode(), bcrypt.gensalt())
    database.save("password", hashed)
```

### **CWE-74: Injection**
**Description**: Improper neutralization of special elements
**Severity**: High
**Languages**: All

```python
# VULNERABLE: LDAP injection
def search_user(username):
    query = f"(uid={username})" # LDAP injection
    return ldap.search(query)

# SECURE: Parameterized LDAP
def search_user_secure(username):
    query = "(uid={})".format(ldap.escape_filter_chars(username))
    return ldap.search(query)
```

### **CWE-434: Unrestricted Upload of File with Dangerous Type**
**Description**: File upload without type validation
**Severity**: High
**Languages**: Web applications

```php
<?php
// VULNERABLE: No file type validation
if (isset($_FILES['upload'])) {
    $uploadFile = 'uploads/' . $_FILES['upload']['name'];
    move_uploaded_file($_FILES['upload']['tmp_name'], $uploadFile);
}

// SECURE: File type validation
if (isset($_FILES['upload'])) {
    $allowedTypes = ['image/jpeg', 'image/png'];
    if (in_array($_FILES['upload']['type'], $allowedTypes)) {
        $uploadFile = 'uploads/' . $_FILES['upload']['name'];
        move_uploaded_file($_FILES['upload']['tmp_name'], $uploadFile);
    }
}
?>
```

### **CWE-807: Reliance on Untrusted Inputs in a Security Decision**
**Description**: Security decision based on untrusted input
**Severity**: High
**Languages**: All

```javascript
// VULNERABLE: Trust client-side input
function isAdmin(user) {
    return user.isAdmin; // Client can modify
}

// SECURE: Server-side validation
function isAdminSecure(userId) {
    const user = database.getUser(userId);
    return user.role === 'admin'; // Server-side check
}
```

### **CWE-250: Execution with Unnecessary Privileges**
**Description**: Application runs with excessive privileges
**Severity**: Medium
**Languages**: All

```bash
# VULNERABLE: Running as root
sudo ./myapp

# SECURE: Run with minimal privileges
sudo -u appuser ./myapp
```

### **CWE-863: Incorrect Authorization**
**Description**: Access control implementation is incorrect
**Severity**: High
**Languages**: All

```python
# VULNERABLE: Incorrect authorization
def delete_user(user_id, current_user):
    if current_user.id == user_id: # Can delete own account
        return database.delete_user(user_id)
    return False

# SECURE: Proper authorization
def delete_user_secure(user_id, current_user):
    if current_user.role == 'admin' or current_user.id == user_id:
        return database.delete_user(user_id)
    return False
```

### **CWE-639: Authorization Bypass Through User-Controlled Key**
**Description**: Authorization bypass using user-controlled key
**Severity**: High
**Languages**: All

```python
# VULNERABLE: User-controlled key
def get_user_data(user_id):
    return database.get(f"user_{user_id}") # IDOR

# SECURE: Authorization check
def get_user_data_secure(user_id, current_user):
    if current_user.id == user_id or current_user.role == 'admin':
        return database.get(f"user_{user_id}")
    return None
```

### **CWE-327: Use of a Broken or Risky Cryptographic Algorithm**
**Description**: Use of weak cryptographic algorithm
**Severity**: High
**Languages**: All

```python
# VULNERABLE: Weak encryption
import hashlib
def hash_password(password):
    return hashlib.md5(password.encode()).hexdigest() # MD5 is weak

# SECURE: Strong encryption
import bcrypt
def hash_password_secure(password):
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt())
```

### **CWE-306: Missing Authentication for Critical Function**
**Description**: Critical function lacks authentication
**Severity**: High
**Languages**: All

```python
# VULNERABLE: No authentication
def admin_panel():
    return render_admin_panel() # No auth check

# SECURE: Authentication required
def admin_panel_secure():
    if not is_authenticated() or not is_admin():
        return redirect('/login')
    return render_admin_panel()
```

### **CWE-862: Missing Authorization**
**Description**: Function lacks authorization check
**Severity**: High
**Languages**: All

```python
# VULNERABLE: No authorization
def delete_file(filename):
    os.remove(filename) # No permission check

# SECURE: Authorization check
def delete_file_secure(filename, user):
    if user.has_permission('delete', filename):
        os.remove(filename)
    else:
        raise PermissionError("Not authorized")
```

### **CWE-732: Incorrect Permission Assignment for Critical Resource**
**Description**: Critical resource has incorrect permissions
**Severity**: Medium
**Languages**: All

```bash
# VULNERABLE: World-writable config
chmod 666 /etc/myapp.conf

# SECURE: Restricted permissions
chmod 600 /etc/myapp.conf
```

## 🎯 Additional Critical CWEs

### **CWE-502: Deserialization of Untrusted Data**
**Description**: Deserializing untrusted data without validation
**Severity**: Critical
**Languages**: Java, Python, .NET, PHP

```java
// VULNERABLE: Unsafe deserialization
public Object deserializeObject(byte[] data) {
    ObjectInputStream ois = new ObjectInputStream(new ByteArrayInputStream(data));
    return ois.readObject(); // Dangerous!
}

// SECURE: Safe deserialization
public Object deserializeObjectSecure(byte[] data) {
    // Use safe deserialization libraries
    return JsonUtils.fromJson(new String(data), MyClass.class);
}
```

### **CWE-798: Use of Hard-coded Credentials**
**Description**: Hard-coded credentials in source code
**Severity**: High
**Languages**: All

```python
# VULNERABLE: Hard-coded credentials
DB_PASSWORD = "admin123"
API_KEY = "sk-1234567890abcdef"

# SECURE: Environment variables
import os
DB_PASSWORD = os.getenv('DB_PASSWORD')
API_KEY = os.getenv('API_KEY')
```

### **CWE-330: Use of Insufficiently Random Values**
**Description**: Use of predictable random values
**Severity**: Medium
**Languages**: All

```python
# VULNERABLE: Predictable random
import random
session_id = random.randint(1, 1000000) # Predictable

# SECURE: Cryptographically secure random
import secrets
session_id = secrets.token_hex(16)
```

## 🧠 TinyBrain Integration Templates

### **CWE Vulnerability Template:**
```json
{
  "title": "[CWE-XXX] [Vulnerability Name] in [Location]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [CWE-XXX] vulnerability identified in [LOCATION]. [DETAILED_DESCRIPTION]. EXPLOITATION CONFIRMED: [EXPLOITATION_DETAILS].",
  "category": "vulnerability",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[cwe-xxx]", "[vulnerability-type]", "[severity]", "authorized-testing"],
  "source": "authorized-security-assessment",
  "cwe_id": "CWE-XXX",
  "severity": "[critical|high|medium|low]",
  "exploitation_status": "confirmed"
}
```

### **CWE Pattern Template:**
```json
{
  "title": "[CWE-XXX] Pattern - [Pattern Name]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [CWE-XXX] pattern identified in [LOCATION]. [PATTERN_DESCRIPTION]. SECURITY IMPACT: [IMPACT_ANALYSIS].",
  "category": "technique",
  "priority": [1-10],
  "confidence": [0.0-1.0],
  "tags": ["[cwe-xxx]", "[pattern-type]", "code-pattern", "authorized-testing"],
  "source": "authorized-security-assessment",
  "cwe_id": "CWE-XXX",
  "pattern_code": "[PATTERN_CODE_HERE]",
  "remediation_guidance": "[REMEDIATION_STEPS]"
}
```

## 🚀 Usage Instructions

### **CWE-Based Assessment:**
```bash
# Load CWE patterns
cline "Load CWE Top 25 security patterns from TinyBrain dataset for comprehensive vulnerability assessment"

# Search by CWE ID
cline "Search TinyBrain for CWE-79 (XSS) patterns in target codebase"

# Store CWE finding
cline "Store CWE-89 SQL injection vulnerability in TinyBrain: Priority 10, affects login endpoint"
```

### **Standards Compliance:**
```bash
# Map to security standards
cline "Map CWE findings to OWASP Top 10 and NIST security standards in TinyBrain"

# Generate compliance report
cline "Generate CWE compliance report from TinyBrain data for security standards alignment"
```

## 🎯 Remember

This CWE dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS**. All CWE patterns are based on the official [CWE Top 25 Most Dangerous Software Errors](https://cwe.mitre.org/top25/) and are intended for legitimate security testing activities.

**Use responsibly and ethically!** 🛡️
