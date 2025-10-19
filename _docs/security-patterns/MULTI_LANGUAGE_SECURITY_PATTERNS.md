# 🌐 Multi-Language Security Patterns
# Comprehensive vulnerability patterns across programming languages

## 🚨 AUTHORIZATION & COMPLIANCE STATEMENT

**LEGAL AUTHORIZATION**: This dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS** under explicit written authorization from authorized stakeholders.

## 🎯 Language-Specific Vulnerability Patterns

### **Java Security Patterns**

#### **SQL Injection in Java:**
```java
// VULNERABLE: String concatenation
public User getUserById(String userId) {
    String query = "SELECT * FROM users WHERE id = " + userId;
    return database.execute(query);
}

// VULNERABLE: PreparedStatement misuse
public User getUserByName(String name) {
    String query = "SELECT * FROM users WHERE name = '" + name + "'";
    PreparedStatement stmt = connection.prepareStatement(query);
    return stmt.executeQuery();
}

// SECURE: Proper PreparedStatement
public User getUserByIdSecure(String userId) {
    String query = "SELECT * FROM users WHERE id = ?";
    PreparedStatement stmt = connection.prepareStatement(query);
    stmt.setString(1, userId);
    return stmt.executeQuery();
}
```

#### **Java Deserialization Vulnerabilities:**
```java
// VULNERABLE: Unsafe deserialization
public Object deserializeObject(byte[] data) {
    ObjectInputStream ois = new ObjectInputStream(new ByteArrayInputStream(data));
    return ois.readObject(); // Dangerous!
}

// VULNERABLE: Jackson deserialization
@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS)
public class User {
    private String name;
    // Vulnerable to deserialization attacks
}
```

#### **Java Authentication Issues:**
```java
// VULNERABLE: Weak password hashing
public String hashPassword(String password) {
    return DigestUtils.md5Hex(password); // MD5 is weak
}

// VULNERABLE: Hardcoded credentials
private static final String ADMIN_USER = "admin";
private static final String ADMIN_PASS = "password123";
```

### **C#/.NET Security Patterns**

#### **SQL Injection in C#:**
```csharp
// VULNERABLE: String concatenation
public User GetUserById(string userId)
{
    string query = "SELECT * FROM users WHERE id = " + userId;
    return database.Execute(query);
}

// VULNERABLE: String.Format
public User GetUserByName(string name)
{
    string query = string.Format("SELECT * FROM users WHERE name = '{0}'", name);
    return database.Execute(query);
}

// SECURE: Parameterized queries
public User GetUserByIdSecure(string userId)
{
    string query = "SELECT * FROM users WHERE id = @userId";
    using (SqlCommand cmd = new SqlCommand(query, connection))
    {
        cmd.Parameters.AddWithValue("@userId", userId);
        return cmd.ExecuteReader();
    }
}
```

#### **C# Deserialization Vulnerabilities:**
```csharp
// VULNERABLE: BinaryFormatter deserialization
public object DeserializeObject(byte[] data)
{
    BinaryFormatter formatter = new BinaryFormatter();
    using (MemoryStream stream = new MemoryStream(data))
    {
        return formatter.Deserialize(stream); // Dangerous!
    }
}

// VULNERABLE: JSON.NET TypeNameHandling
public class User
{
    [JsonProperty(TypeNameHandling = TypeNameHandling.All)]
    public object Data { get; set; } // Vulnerable
}
```

#### **C# Authentication Issues:**
```csharp
// VULNERABLE: Weak password hashing
public string HashPassword(string password)
{
    using (MD5 md5 = MD5.Create())
    {
        byte[] hash = md5.ComputeHash(Encoding.UTF8.GetBytes(password));
        return Convert.ToBase64String(hash);
    }
}

// VULNERABLE: Hardcoded connection strings
private static readonly string ConnectionString = 
    "Server=localhost;Database=MyDB;User=admin;Password=admin123;";
```

### **PHP Security Patterns**

#### **SQL Injection in PHP:**
```php
<?php
// VULNERABLE: String concatenation
function getUserById($userId) {
    $query = "SELECT * FROM users WHERE id = " . $userId;
    return mysql_query($query);
}

// VULNERABLE: Direct variable interpolation
function getUserByName($name) {
    $query = "SELECT * FROM users WHERE name = '$name'";
    return mysql_query($query);
}

// SECURE: Prepared statements
function getUserByIdSecure($userId) {
    $stmt = $pdo->prepare("SELECT * FROM users WHERE id = ?");
    $stmt->execute([$userId]);
    return $stmt->fetch();
}
?>
```

#### **PHP File Upload Vulnerabilities:**
```php
<?php
// VULNERABLE: No file type validation
if (isset($_FILES['upload'])) {
    $uploadDir = 'uploads/';
    $uploadFile = $uploadDir . basename($_FILES['upload']['name']);
    move_uploaded_file($_FILES['upload']['tmp_name'], $uploadFile);
}

// VULNERABLE: Path traversal
function readFile($filename) {
    return file_get_contents("uploads/" . $filename);
}
?>
```

#### **PHP Authentication Issues:**
```php
<?php
// VULNERABLE: Weak session management
session_start();
$_SESSION['user_id'] = $userId; // No session regeneration

// VULNERABLE: Weak password hashing
function hashPassword($password) {
    return md5($password); // MD5 is weak
}

// VULNERABLE: SQL injection in authentication
function authenticate($username, $password) {
    $query = "SELECT * FROM users WHERE username = '$username' AND password = '$password'";
    return mysql_query($query);
}
?>
```

### **Ruby Security Patterns**

#### **SQL Injection in Ruby:**
```ruby
# VULNERABLE: String interpolation
def get_user_by_id(user_id)
  query = "SELECT * FROM users WHERE id = #{user_id}"
  ActiveRecord::Base.connection.execute(query)
end

# VULNERABLE: String concatenation
def get_user_by_name(name)
  query = "SELECT * FROM users WHERE name = '" + name + "'"
  ActiveRecord::Base.connection.execute(query)
end

# SECURE: Parameterized queries
def get_user_by_id_secure(user_id)
  User.where("id = ?", user_id)
end
```

#### **Ruby Mass Assignment Vulnerabilities:**
```ruby
# VULNERABLE: Mass assignment
class UsersController < ApplicationController
  def update
    @user = User.find(params[:id])
    @user.update(params[:user]) # Dangerous!
  end
end

# SECURE: Strong parameters
class UsersController < ApplicationController
  def update
    @user = User.find(params[:id])
    @user.update(user_params)
  end

  private
  def user_params
    params.require(:user).permit(:name, :email)
  end
end
```

#### **Ruby Authentication Issues:**
```ruby
# VULNERABLE: Weak password hashing
def hash_password(password)
  Digest::MD5.hexdigest(password) # MD5 is weak
end

# VULNERABLE: Hardcoded secrets
SECRET_KEY = "mysecretkey123"
```

### **Go Security Patterns**

#### **SQL Injection in Go:**
```go
// VULNERABLE: String concatenation
func GetUserByID(userID string) (*User, error) {
    query := "SELECT * FROM users WHERE id = " + userID
    rows, err := db.Query(query)
    // ...
}

// VULNERABLE: fmt.Sprintf
func GetUserByName(name string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
    rows, err := db.Query(query)
    // ...
}

// SECURE: Parameterized queries
func GetUserByIDSecure(userID string) (*User, error) {
    query := "SELECT * FROM users WHERE id = ?"
    rows, err := db.Query(query, userID)
    // ...
}
```

#### **Go Command Injection:**
```go
// VULNERABLE: Command injection
func ProcessFile(filename string) error {
    cmd := exec.Command("cat", filename)
    return cmd.Run()
}

// VULNERABLE: Shell injection
func ProcessFileShell(filename string) error {
    cmd := exec.Command("sh", "-c", "cat "+filename)
    return cmd.Run()
}

// SECURE: No shell, direct command
func ProcessFileSecure(filename string) error {
    cmd := exec.Command("cat", filename)
    return cmd.Run()
}
```

#### **Go Authentication Issues:**
```go
// VULNERABLE: Weak password hashing
func HashPassword(password string) string {
    h := md5.New()
    h.Write([]byte(password))
    return hex.EncodeToString(h.Sum(nil))
}

// VULNERABLE: Hardcoded secrets
const SecretKey = "mysecretkey123"
```

### **C/C++ Security Patterns**

#### **Buffer Overflow in C:**
```c
// VULNERABLE: Buffer overflow
void copy_string(char *input) {
    char buffer[100];
    strcpy(buffer, input); // No bounds checking
}

// VULNERABLE: Format string vulnerability
void print_user(char *username) {
    printf(username); // Format string vulnerability
}

// SECURE: Bounds checking
void copy_string_secure(char *input, size_t input_len) {
    char buffer[100];
    size_t copy_len = (input_len < sizeof(buffer)) ? input_len : sizeof(buffer) - 1;
    strncpy(buffer, input, copy_len);
    buffer[copy_len] = '\0';
}
```

#### **C Memory Management Issues:**
```c
// VULNERABLE: Memory leak
char* get_data() {
    char *data = malloc(1000);
    // No free() call
    return data;
}

// VULNERABLE: Use after free
void use_after_free() {
    char *ptr = malloc(100);
    free(ptr);
    *ptr = 'A'; // Use after free
}

// VULNERABLE: Double free
void double_free() {
    char *ptr = malloc(100);
    free(ptr);
    free(ptr); // Double free
}
```

### **TypeScript Security Patterns**

#### **TypeScript XSS Vulnerabilities:**
```typescript
// VULNERABLE: InnerHTML without sanitization
function displayUserInput(input: string) {
    document.getElementById('output').innerHTML = input; // XSS
}

// VULNERABLE: Template literal injection
function createQuery(userInput: string) {
    return `SELECT * FROM users WHERE name = '${userInput}'`; // SQL injection
}

// SECURE: Proper sanitization
function displayUserInputSecure(input: string) {
    const sanitized = DOMPurify.sanitize(input);
    document.getElementById('output').innerHTML = sanitized;
}
```

#### **TypeScript Authentication Issues:**
```typescript
// VULNERABLE: Client-side authentication
function authenticate(username: string, password: string) {
    if (username === 'admin' && password === 'admin') {
        localStorage.setItem('authenticated', 'true'); // Client-side only
        return true;
    }
    return false;
}

// VULNERABLE: Weak session management
interface Session {
    userId: string;
    isAdmin: boolean;
}

function createSession(userId: string): Session {
    return {
        userId: userId,
        isAdmin: userId === 'admin' // Client-side role check
    };
}
```

### **Rust Security Patterns**

#### **Rust Unsafe Code Issues:**
```rust
// VULNERABLE: Unsafe pointer dereference
unsafe fn dereference_pointer(ptr: *const i32) -> i32 {
    *ptr // Can cause segmentation fault
}

// VULNERABLE: Buffer overflow potential
fn copy_string(input: &str) {
    let mut buffer = [0u8; 100];
    let bytes = input.as_bytes();
    for i in 0..bytes.len() {
        buffer[i] = bytes[i]; // No bounds checking
    }
}

// SECURE: Safe Rust patterns
fn copy_string_secure(input: &str) {
    let mut buffer = [0u8; 100];
    let bytes = input.as_bytes();
    let copy_len = bytes.len().min(buffer.len());
    buffer[..copy_len].copy_from_slice(&bytes[..copy_len]);
}
```

## 🎯 Language-Specific TinyBrain Templates

### **Java Template:**
```json
{
  "title": "Java SQL Injection - PreparedStatement Misuse",
  "content": "AUTHORIZED SECURITY ASSESSMENT: Java SQL injection vulnerability identified. PreparedStatement used incorrectly with string concatenation. EXPLOITATION CONFIRMED: Authentication bypass achieved.",
  "category": "vulnerability",
  "priority": 9,
  "confidence": 0.95,
  "tags": ["java", "sql-injection", "preparedstatement", "authentication-bypass", "authorized-testing"],
  "source": "authorized-security-assessment",
  "language": "java",
  "vulnerability_type": "sql-injection",
  "exploitation_status": "confirmed"
}
```

### **C# Template:**
```json
{
  "title": "C# SQL Injection - String.Format Vulnerability",
  "content": "AUTHORIZED SECURITY ASSESSMENT: C# SQL injection vulnerability using String.Format. User input directly interpolated into SQL query. EXPLOITATION CONFIRMED: Database access achieved.",
  "category": "vulnerability",
  "priority": 9,
  "confidence": 0.95,
  "tags": ["csharp", "sql-injection", "string-format", "database-access", "authorized-testing"],
  "source": "authorized-security-assessment",
  "language": "csharp",
  "vulnerability_type": "sql-injection",
  "exploitation_status": "confirmed"
}
```

## 🚀 Usage Instructions

### **Language-Specific Assessment:**
```bash
# Java security assessment
cline "Load Java security patterns from TinyBrain dataset and analyze target Java codebase"

# C# security assessment  
cline "Load C# security patterns from TinyBrain dataset and analyze target .NET codebase"

# Multi-language assessment
cline "Load multi-language security patterns from TinyBrain dataset for comprehensive code review"
```

### **Pattern Matching by Language:**
```bash
# Match Java patterns
cline "Search TinyBrain for Java-specific vulnerability patterns in target codebase"

# Match C# patterns
cline "Search TinyBrain for C#-specific vulnerability patterns in target codebase"

# Cross-language analysis
cline "Compare vulnerability patterns across different languages in TinyBrain dataset"
```

## 🎯 Remember

This multi-language dataset is designed for **AUTHORIZED SECURITY ASSESSMENTS** conducted by **CERTIFIED CYBERSECURITY PROFESSIONALS**. All patterns include proper authorization language and are intended for legitimate security testing activities.

**Use responsibly and ethically!** 🛡️
