# 🌐 Enhanced Language & Library Security Patterns
# Comprehensive vulnerability patterns across frameworks and libraries

## 🚨 AUTHORIZATION & COMPLIANCE STATEMENT

**LEGAL AUTHORIZATION**: This dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS** under explicit written authorization from authorized stakeholders.

## ⚠️ **IMPORTANT DISCLAIMER**

**Security patterns evolve constantly!** The examples provided are **illustrative** and represent **common vulnerability patterns** as of 2024. Security landscapes change rapidly with:
- New framework versions and security updates
- Emerging attack vectors and exploitation techniques
- Updated security best practices and recommendations
- New vulnerability discoveries and CVE publications

**Always use the latest security information and consult current documentation for your specific framework versions.**

## 🎯 Python Framework & Library Patterns

### **Django Framework**

#### **CWE-89: SQL Injection in Django ORM**
```python
# VULNERABLE: Raw SQL without parameterization
from django.db import connection

def get_user_vulnerable(user_id):
    with connection.cursor() as cursor:
        cursor.execute(f"SELECT * FROM auth_user WHERE id = {user_id}")
        return cursor.fetchone()

# VULNERABLE: Unsafe extra() usage
def search_users_vulnerable(query):
    return User.objects.extra(where=[f"username LIKE '%{query}%'"])

# SECURE: Django ORM with parameterized queries
def get_user_secure(user_id):
    return User.objects.get(id=user_id)

def search_users_secure(query):
    return User.objects.filter(username__icontains=query)
```

#### **CWE-79: XSS in Django Templates**
```python
# VULNERABLE: Unsafe template rendering
def render_user_content_vulnerable(request):
    user_input = request.GET.get('content', '')
    return render(request, 'template.html', {'content': user_input})

# SECURE: Auto-escaping enabled (Django default)
def render_user_content_secure(request):
    user_input = request.GET.get('content', '')
    return render(request, 'template.html', {'content': user_input})

# SECURE: Manual escaping for trusted content
from django.utils.html import escape
def render_trusted_content_secure(request):
    trusted_content = request.GET.get('content', '')
    return render(request, 'template.html', {'content': escape(trusted_content)})
```

#### **CWE-434: File Upload in Django**
```python
# VULNERABLE: No file type validation
def upload_file_vulnerable(request):
    if request.method == 'POST':
        uploaded_file = request.FILES['file']
        with open(f'uploads/{uploaded_file.name}', 'wb') as destination:
            for chunk in uploaded_file.chunks():
                destination.write(chunk)

# SECURE: File type validation and secure handling
import os
from django.core.files.storage import FileSystemStorage
from django.core.files.base import ContentFile

def upload_file_secure(request):
    if request.method == 'POST':
        uploaded_file = request.FILES['file']
        # Validate file type
        allowed_types = ['image/jpeg', 'image/png', 'application/pdf']
        if uploaded_file.content_type not in allowed_types:
            raise ValueError("Invalid file type")
        
        # Secure file handling
        fs = FileSystemStorage(location='secure_uploads/')
        filename = fs.save(uploaded_file.name, uploaded_file)
        return filename
```

### **Flask Framework**

#### **CWE-89: SQL Injection in Flask-SQLAlchemy**
```python
# VULNERABLE: String concatenation in raw SQL
from flask import Flask, request
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
db = SQLAlchemy(app)

@app.route('/user/<user_id>')
def get_user_vulnerable(user_id):
    result = db.engine.execute(f"SELECT * FROM users WHERE id = {user_id}")
    return result.fetchone()

# SECURE: SQLAlchemy ORM
@app.route('/user/<int:user_id>')
def get_user_secure(user_id):
    user = User.query.get(user_id)
    return user
```

#### **CWE-352: CSRF in Flask**
```python
# VULNERABLE: No CSRF protection
from flask import Flask, request, render_template

app = Flask(__name__)

@app.route('/transfer', methods=['POST'])
def transfer_money_vulnerable():
    amount = request.form['amount']
    # No CSRF protection
    return f"Transferred ${amount}"

# SECURE: Flask-WTF CSRF protection
from flask_wtf.csrf import CSRFProtect

app = Flask(__name__)
csrf = CSRFProtect(app)

@app.route('/transfer', methods=['POST'])
def transfer_money_secure():
    amount = request.form['amount']
    # CSRF token automatically validated
    return f"Transferred ${amount}"
```

### **Python Libraries**

#### **CWE-78: Command Injection with subprocess**
```python
# VULNERABLE: Shell injection
import subprocess
import os

def process_file_vulnerable(filename):
    # Dangerous - uses shell
    result = subprocess.run(f"cat {filename}", shell=True, capture_output=True)
    return result.stdout

# SECURE: No shell, direct command
def process_file_secure(filename):
    # Safe - no shell, direct command
    result = subprocess.run(['cat', filename], capture_output=True)
    return result.stdout
```

#### **CWE-502: Deserialization with pickle**
```python
# VULNERABLE: Unsafe pickle deserialization
import pickle

def deserialize_data_vulnerable(data):
    return pickle.loads(data)  # Dangerous!

# SECURE: Safe deserialization with validation
import json

def deserialize_data_secure(data):
    # Use JSON instead of pickle for untrusted data
    return json.loads(data)
```

## 🎯 JavaScript/Node.js Framework & Library Patterns

### **Express.js Framework**

#### **CWE-89: SQL Injection in Express**
```javascript
// VULNERABLE: String concatenation
const express = require('express');
const mysql = require('mysql');
const app = express();

app.get('/user/:id', (req, res) => {
    const query = `SELECT * FROM users WHERE id = ${req.params.id}`;
    connection.query(query, (err, results) => {
        res.json(results);
    });
});

// SECURE: Parameterized queries
app.get('/user/:id', (req, res) => {
    const query = 'SELECT * FROM users WHERE id = ?';
    connection.query(query, [req.params.id], (err, results) => {
        res.json(results);
    });
});
```

#### **CWE-79: XSS in Express**
```javascript
// VULNERABLE: No output encoding
app.get('/search', (req, res) => {
    const query = req.query.q;
    res.send(`<h1>Search results for: ${query}</h1>`);
});

// SECURE: Output encoding
const escapeHtml = require('escape-html');

app.get('/search', (req, res) => {
    const query = req.query.q;
    res.send(`<h1>Search results for: ${escapeHtml(query)}</h1>`);
});
```

### **React Framework**

#### **CWE-79: XSS in React**
```jsx
// VULNERABLE: dangerouslySetInnerHTML without sanitization
function UserProfile({ userContent }) {
    return (
        <div dangerouslySetInnerHTML={{ __html: userContent }} />
    );
}

// SECURE: React's built-in XSS protection
function UserProfile({ userContent }) {
    return <div>{userContent}</div>; // React auto-escapes
}

// SECURE: Manual sanitization for trusted HTML
import DOMPurify from 'dompurify';

function UserProfile({ userContent }) {
    const sanitizedContent = DOMPurify.sanitize(userContent);
    return <div dangerouslySetInnerHTML={{ __html: sanitizedContent }} />;
}
```

### **Node.js Libraries**

#### **CWE-78: Command Injection with child_process**
```javascript
// VULNERABLE: Shell injection
const { exec } = require('child_process');

function processFile(filename) {
    exec(`cat ${filename}`, (error, stdout, stderr) => {
        console.log(stdout);
    });
}

// SECURE: No shell, direct command
const { spawn } = require('child_process');

function processFileSecure(filename) {
    const child = spawn('cat', [filename]);
    child.stdout.on('data', (data) => {
        console.log(data.toString());
    });
}
```

## 🎯 Java Framework & Library Patterns

### **Spring Framework**

#### **CWE-89: SQL Injection in Spring JDBC**
```java
// VULNERABLE: String concatenation
@Repository
public class UserRepository {
    @Autowired
    private JdbcTemplate jdbcTemplate;
    
    public User findUserVulnerable(String userId) {
        String sql = "SELECT * FROM users WHERE id = " + userId;
        return jdbcTemplate.queryForObject(sql, User.class);
    }
}

// SECURE: Parameterized queries
@Repository
public class UserRepository {
    @Autowired
    private JdbcTemplate jdbcTemplate;
    
    public User findUserSecure(String userId) {
        String sql = "SELECT * FROM users WHERE id = ?";
        return jdbcTemplate.queryForObject(sql, new Object[]{userId}, User.class);
    }
}
```

#### **CWE-502: Deserialization in Spring**
```java
// VULNERABLE: Unsafe deserialization
public Object deserializeVulnerable(byte[] data) {
    ObjectInputStream ois = new ObjectInputStream(new ByteArrayInputStream(data));
    return ois.readObject(); // Dangerous!
}

// SECURE: Safe deserialization with validation
@JsonTypeInfo(use = JsonTypeInfo.Id.NONE) // Disable polymorphic deserialization
public class SafeDeserialization {
    // Use Jackson with safe configuration
}
```

### **Java Libraries**

#### **CWE-327: Weak Hashing with MessageDigest**
```java
// VULNERABLE: MD5 hashing
import java.security.MessageDigest;

public String hashPasswordVulnerable(String password) {
    MessageDigest md = MessageDigest.getInstance("MD5");
    byte[] hash = md.digest(password.getBytes());
    return bytesToHex(hash);
}

// SECURE: BCrypt hashing
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

public String hashPasswordSecure(String password) {
    BCryptPasswordEncoder encoder = new BCryptPasswordEncoder();
    return encoder.encode(password);
}
```

## 🎯 C#/.NET Framework & Library Patterns

### **ASP.NET Core**

#### **CWE-89: SQL Injection in Entity Framework**
```csharp
// VULNERABLE: Raw SQL without parameterization
public User GetUserVulnerable(string userId)
{
    string sql = $"SELECT * FROM Users WHERE Id = {userId}";
    return context.Users.FromSqlRaw(sql).FirstOrDefault();
}

// SECURE: Entity Framework with parameterized queries
public User GetUserSecure(string userId)
{
    return context.Users.Where(u => u.Id == userId).FirstOrDefault();
}
```

#### **CWE-79: XSS in Razor Views**
```csharp
// VULNERABLE: No encoding
@Html.Raw(Model.UserContent)

// SECURE: Auto-encoding (Razor default)
@Model.UserContent

// SECURE: Manual encoding for trusted content
@Html.Encode(Model.TrustedContent)
```

### **.NET Libraries**

#### **CWE-502: Deserialization with BinaryFormatter**
```csharp
// VULNERABLE: BinaryFormatter deserialization
public object DeserializeVulnerable(byte[] data)
{
    BinaryFormatter formatter = new BinaryFormatter();
    using (MemoryStream stream = new MemoryStream(data))
    {
        return formatter.Deserialize(stream); // Dangerous!
    }
}

// SECURE: JSON deserialization
public T DeserializeSecure<T>(string json)
{
    return JsonSerializer.Deserialize<T>(json);
}
```

## 🎯 PHP Framework & Library Patterns

### **Laravel Framework**

#### **CWE-89: SQL Injection in Laravel Eloquent**
```php
<?php
// VULNERABLE: Raw SQL without parameterization
public function getUserVulnerable($userId)
{
    $sql = "SELECT * FROM users WHERE id = " . $userId;
    return DB::select($sql);
}

// SECURE: Eloquent ORM
public function getUserSecure($userId)
{
    return User::find($userId);
}
?>
```

#### **CWE-352: CSRF in Laravel**
```php
<?php
// VULNERABLE: No CSRF protection
Route::post('/transfer', function(Request $request) {
    $amount = $request->input('amount');
    // No CSRF protection
    return "Transferred $amount";
});

// SECURE: Laravel CSRF protection
Route::post('/transfer', function(Request $request) {
    $amount = $request->input('amount');
    // CSRF token automatically validated
    return "Transferred $amount";
})->middleware('csrf');
?>
```

## 🎯 TinyBrain Integration Templates

### **Framework-Specific Vulnerability Template:**
```json
{
  "title": "[FRAMEWORK] [CWE-XXX] [Vulnerability Name] in [Location]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [FRAMEWORK] [CWE-XXX] vulnerability identified in [LOCATION]. [DETAILED_DESCRIPTION]. EXPLOITATION CONFIRMED: [EXPLOITATION_DETAILS].",
  "category": "vulnerability",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["[framework]", "[cwe-xxx]", "[vulnerability-type]", "[severity]", "authorized-testing"],
  "source": "authorized-security-assessment",
  "framework": "[FRAMEWORK_NAME]",
  "library": "[LIBRARY_NAME]",
  "cwe_id": "CWE-XXX",
  "exploitation_status": "confirmed"
}
```

### **Library-Specific Pattern Template:**
```json
{
  "title": "[LIBRARY] [CWE-XXX] Pattern - [Pattern Name]",
  "content": "AUTHORIZED SECURITY ASSESSMENT: [LIBRARY] [CWE-XXX] pattern identified in [LOCATION]. [PATTERN_DESCRIPTION]. SECURITY IMPACT: [IMPACT_ANALYSIS].",
  "category": "technique",
  "priority": "[1-10]",
  "confidence": "[0.0-1.0]",
  "tags": ["[library]", "[cwe-xxx]", "[pattern-type]", "code-pattern", "authorized-testing"],
  "source": "authorized-security-assessment",
  "framework": "[FRAMEWORK_NAME]",
  "library": "[LIBRARY_NAME]",
  "cwe_id": "CWE-XXX",
  "pattern_code": "[PATTERN_CODE_HERE]",
  "remediation_guidance": "[REMEDIATION_STEPS]"
}
```

## 🚀 Usage Instructions

### **Framework-Specific Assessment:**
```bash
# Django security assessment
cline "Load Django-specific CWE patterns from TinyBrain dataset and analyze target Django application"

# Express.js security assessment
cline "Load Express.js-specific CWE patterns from TinyBrain dataset and analyze target Node.js application"

# Spring security assessment
cline "Load Spring-specific CWE patterns from TinyBrain dataset and analyze target Java application"
```

### **Library-Specific Pattern Matching:**
```bash
# Python library patterns
cline "Search TinyBrain for Python library security patterns: subprocess, pickle, requests"

# JavaScript library patterns
cline "Search TinyBrain for JavaScript library security patterns: child_process, express, react"

# Java library patterns
cline "Search TinyBrain for Java library security patterns: MessageDigest, ObjectInputStream, JdbcTemplate"
```

## 🎯 Remember

**Security patterns evolve constantly!** These examples are **illustrative** and represent **common vulnerability patterns**. Always:

- **Consult latest documentation** for your specific framework versions
- **Stay updated** with security advisories and CVE publications
- **Use current best practices** and security recommendations
- **Validate patterns** against your specific environment and requirements

**Use responsibly and ethically!** 🛡️
