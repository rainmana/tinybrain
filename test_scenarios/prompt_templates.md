# TinyBrain Code Review Prompt Templates

## 🚀 **Quick Start Prompts**

### **1. Session Setup Prompt**
```
I'm starting a security code review of [PACKAGE_NAME]. Please help me set up a comprehensive TinyBrain session for this assessment.

First, create a session with:
- Name: "Security Review: [PACKAGE_NAME]"
- Task Type: "security_review" 
- Description: "Comprehensive security review focusing on authentication, input validation, and data handling vulnerabilities"

Then set up task progress tracking for these stages:
1. Reconnaissance & Setup
2. Authentication Analysis  
3. Input Validation Review
4. Data Handling Assessment
5. Session Management Review
6. File Upload Security
7. Final Report Generation

Finally, get the available security templates so we can use them for common vulnerability patterns.

Let's start with the session creation and then move into the actual code analysis.
```

### **2. Vulnerability Discovery Prompt**
```
Now let's analyze the [PACKAGE_NAME] codebase for security vulnerabilities. Please:

1. Review the main entry points and identify potential security issues
2. For each vulnerability found, store it in TinyBrain with:
   - Appropriate category (vulnerability, exploit, technique, tool, evidence, etc.)
   - Priority level (1-10) based on severity and exploitability
   - Confidence level (0.0-1.0) based on certainty of the finding
   - Detailed description, impact assessment, and affected code sections
   - Tags for easy categorization

3. Create relationships between related vulnerabilities (e.g., "exploits", "depends_on", "causes")
4. Update task progress as we complete each analysis phase
5. Use semantic search to find similar patterns across the codebase

Focus on these vulnerability types:
- Authentication bypasses and session management flaws
- Input validation and injection vulnerabilities
- Cross-Site Scripting (XSS) and Cross-Site Request Forgery (CSRF)
- Insecure direct object references
- Security misconfigurations
- Sensitive data exposure
- File upload vulnerabilities

Store detailed findings that could be used in a real security report.
```

### **3. Exploit Development Prompt**
```
Based on our vulnerability analysis, let's develop proof-of-concept exploits for the most critical issues. Please:

1. Search TinyBrain for our highest priority vulnerabilities (priority >= 8)
2. For the top 3 most critical findings, create working exploit code that demonstrates the issue
3. Store each exploit in TinyBrain with:
   - Category: "exploit" or "payload"
   - Relationship to the original vulnerability
   - Detailed explanation of how the exploit works
   - Potential impact and remediation steps
   - Working code examples

4. Create a comprehensive security report summarizing all findings
5. Export the session data for documentation

The exploits should be realistic and demonstrate actual security impact. Include both the vulnerability details and the working exploit code.
```

### **4. Advanced Features Test Prompt**
```
Let's test TinyBrain's advanced features to ensure everything works correctly:

1. Check for duplicate findings and similar vulnerabilities
2. Create custom memory templates for common vulnerability patterns we've seen
3. Test the notification system for high-priority findings
4. Perform cleanup operations on any low-priority test entries
5. Generate comprehensive system diagnostics
6. Test semantic search for finding conceptually similar issues
7. Create context snapshots at key milestones
8. Get memory statistics and usage analytics

This will help validate that all TinyBrain features work correctly in a real-world scenario.
```

## 🎯 **Package-Specific Prompts**

### **For Express.js Core:**
```
I'm reviewing the Express.js framework for security vulnerabilities. Focus on:
- Route handling and middleware security
- Request/response processing vulnerabilities
- Error handling and information disclosure
- HTTP header manipulation
- Path traversal in static file serving
```

### **For Passport.js:**
```
I'm reviewing Passport.js authentication middleware. Focus on:
- Authentication strategy vulnerabilities
- Session management and persistence
- OAuth and social login security
- Password handling and storage
- Authentication bypass techniques
```

### **For Helmet.js:**
```
I'm reviewing Helmet.js security middleware. Focus on:
- HTTP security header implementation
- Content Security Policy (CSP) handling
- XSS protection mechanisms
- Clickjacking prevention
- Security header bypass techniques
```

### **For Express-Session:**
```
I'm reviewing Express-Session for security vulnerabilities. Focus on:
- Session storage and serialization
- Session hijacking and fixation
- Cookie security and configuration
- Session timeout and invalidation
- Cross-site scripting in session handling
```

## 🔧 **TinyBrain Integration Prompts**

### **Memory Storage Template:**
```
Store this finding in TinyBrain:
- Title: "[VULNERABILITY_NAME]"
- Content: "[DETAILED_DESCRIPTION]"
- Category: "[vulnerability/exploit/technique/tool/evidence]"
- Priority: [1-10]
- Confidence: [0.0-1.0]
- Tags: ["[tag1]", "[tag2]", "[tag3]"]
- Source: "[CODE_SECTION_OR_FILE]"
```

### **Relationship Creation Template:**
```
Create a relationship in TinyBrain:
- Source: "[MEMORY_ID_1]"
- Target: "[MEMORY_ID_2]"
- Type: "[exploits/depends_on/causes/mitigates/references]"
- Description: "[RELATIONSHIP_DESCRIPTION]"
- Strength: [0.0-1.0]
```

### **Search and Analysis Template:**
```
Search TinyBrain for:
- Query: "[SEARCH_TERMS]"
- Session: "[SESSION_ID]"
- Categories: ["[category1]", "[category2]"]
- Min Priority: [PRIORITY_THRESHOLD]
- Min Confidence: [CONFIDENCE_THRESHOLD]
```

## 📊 **Success Validation Prompts**

### **Comprehensive Report Generation:**
```
Generate a comprehensive security report from our TinyBrain session:

1. Get all memory entries from our session
2. Organize findings by category and priority
3. Include relationships between vulnerabilities and exploits
4. Provide executive summary with risk assessment
5. Include detailed technical findings with code examples
6. Export the complete session data

This should be a professional security report ready for stakeholders.
```

### **System Validation:**
```
Validate that TinyBrain is working correctly by:

1. Checking system health and diagnostics
2. Verifying all 40 MCP tools are functional
3. Testing memory statistics and usage analytics
4. Confirming notification system is working
5. Validating export/import functionality
6. Checking cleanup operations

Report any issues or confirm everything is working perfectly.
```

## 🎯 **Expected Outcomes**

After using these prompts, you should have:

- **Complete security assessment** with 15-25 memory entries
- **Working exploit code** for 3-5 critical vulnerabilities  
- **Comprehensive relationships** between findings
- **Professional security report** ready for documentation
- **Full TinyBrain validation** with all 40 tools tested
- **Exportable session data** for sharing and backup

This will demonstrate TinyBrain's full capabilities while producing valuable security research!
