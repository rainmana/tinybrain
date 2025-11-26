# TinyBrain Code Review Test Scenario

## 🎯 **Test Objective**
Demonstrate TinyBrain's capabilities through a complete security code review workflow, from initial assessment to exploit validation.

## 📋 **Test Setup**

### **Step 1: Initialize TinyBrain Session**
```
Create a new security assessment session for code review:
- Session Name: "Express.js Web App Security Review"
- Task Type: "security_review"
- Description: "Comprehensive security review of Express.js web application with focus on authentication, input validation, and data handling"
```

### **Step 2: Select Target Package**
Choose a real-world package to review. Good candidates:
- **express**: Core Express.js framework
- **passport**: Authentication middleware
- **helmet**: Security middleware
- **express-session**: Session management
- **body-parser**: Request parsing
- **multer**: File upload handling

## 🔍 **Code Review Workflow**

### **Phase 1: Initial Assessment & Setup**
```
Prompt: "I'm starting a security code review of [PACKAGE_NAME]. Please help me:

1. Create a TinyBrain session for this security review
2. Set up initial task progress tracking with these stages:
   - Reconnaissance & Setup
   - Authentication Analysis
   - Input Validation Review
   - Data Handling Assessment
   - Session Management Review
   - File Upload Security
   - Final Report Generation

3. Create context snapshots at key milestones
4. Use security templates for common vulnerability patterns
"

Expected TinyBrain Usage:
- create_session
- create_task_progress (multiple stages)
- create_context_snapshot
- get_security_templates
```

### **Phase 2: Vulnerability Discovery**
```
Prompt: "Now let's analyze the [PACKAGE_NAME] codebase for security vulnerabilities. Please:

1. Review the main entry points and identify potential security issues
2. For each vulnerability found, store it in TinyBrain with:
   - Appropriate category (vulnerability, exploit, technique, etc.)
   - Priority level (1-10) based on severity
   - Confidence level (0.0-1.0) based on certainty
   - Detailed description and impact assessment
   - Affected code sections and line numbers

3. Create relationships between related vulnerabilities
4. Update task progress as we complete each analysis phase
5. Use semantic search to find similar patterns across the codebase

Focus on these vulnerability types:
- Authentication bypasses
- Input validation flaws
- Injection vulnerabilities (SQL, NoSQL, Command)
- Cross-Site Scripting (XSS)
- Cross-Site Request Forgery (CSRF)
- Insecure direct object references
- Security misconfigurations
- Sensitive data exposure
"

Expected TinyBrain Usage:
- store_memory (multiple vulnerability entries)
- create_relationship (linking related issues)
- update_task_progress
- semantic_search
- find_similar_memories
```

### **Phase 3: Exploit Development & Validation**
```
Prompt: "For the high-priority vulnerabilities we found, let's develop proof-of-concept exploits. Please:

1. Identify the top 3 highest priority vulnerabilities
2. For each vulnerability, create exploit code that demonstrates the issue
3. Store the exploit code in TinyBrain with:
   - Category: 'exploit' or 'payload'
   - Relationship to the original vulnerability
   - Detailed explanation of how the exploit works
   - Potential impact and remediation steps

4. Create a comprehensive security report summarizing all findings
5. Export the session data for documentation

Focus on creating realistic, working exploits that could be used for:
- Demonstrating the vulnerability to developers
- Penetration testing scenarios
- Security training materials
"

Expected TinyBrain Usage:
- search_memories (filter by priority/confidence)
- store_memory (exploit entries)
- create_relationship (exploit -> vulnerability)
- get_context_summary
- export_session_data
```

### **Phase 4: Advanced Features Testing**
```
Prompt: "Let's test TinyBrain's advanced features:

1. Check for duplicate findings and similar vulnerabilities
2. Create memory templates for common vulnerability patterns
3. Test the notification system for high-priority findings
4. Perform cleanup operations on low-priority test entries
5. Generate comprehensive system diagnostics
6. Test semantic search for finding conceptually similar issues

This will help validate that all TinyBrain features work correctly in a real-world scenario.
"

Expected TinyBrain Usage:
- check_duplicates
- check_duplicate_memories
- get_notifications
- check_high_priority_memories
- cleanup_low_priority_memories
- get_system_diagnostics
- semantic_search
```

## 📊 **Success Metrics**

### **TinyBrain Usage Targets:**
- **Sessions**: 1 main session + context snapshots
- **Memory Entries**: 15-25 (vulnerabilities, exploits, techniques, evidence)
- **Relationships**: 8-12 (linking related findings)
- **Task Progress**: 7 stages with updates
- **Notifications**: 3-5 high-priority alerts
- **Templates**: 2-3 custom vulnerability templates

### **Quality Metrics:**
- All 40 MCP tools tested and working
- Realistic vulnerability findings with proper categorization
- Working exploit code that demonstrates actual issues
- Comprehensive relationships between findings
- Complete session export with all data intact

## 🎯 **Sample Prompts for Each Phase**

### **Phase 1 Prompt:**
```
I'm conducting a security code review of the [PACKAGE_NAME] npm package. Please help me set up a comprehensive TinyBrain session to track this assessment. I want to:

1. Create a session called "Security Review: [PACKAGE_NAME]"
2. Set up task progress tracking for a complete security assessment
3. Initialize context snapshots for key milestones
4. Get the available security templates for common vulnerability patterns

Let's start with the session creation and then move into the actual code analysis.
```

### **Phase 2 Prompt:**
```
Now let's dive into the actual security analysis. I've selected [PACKAGE_NAME] for review. Please:

1. Analyze the codebase for security vulnerabilities
2. For each finding, store it in TinyBrain with proper categorization, priority, and confidence
3. Create relationships between related vulnerabilities
4. Update our task progress as we complete each analysis phase
5. Use semantic search to find similar patterns

Focus on authentication, input validation, and data handling security issues. Store detailed findings that could be used in a real security report.
```

### **Phase 3 Prompt:**
```
Based on our vulnerability analysis, let's develop proof-of-concept exploits for the most critical issues. Please:

1. Identify the top 3 highest priority vulnerabilities
2. Create working exploit code for each
3. Store the exploits in TinyBrain with relationships to the original vulnerabilities
4. Generate a comprehensive security report
5. Export all our findings for documentation

The exploits should be realistic and demonstrate actual security impact.
```

## 🚀 **Expected Outcomes**

After completing this test scenario, you should have:

1. **A complete security assessment** with real findings
2. **Working exploit code** that demonstrates vulnerabilities
3. **Comprehensive TinyBrain session** with all features tested
4. **Exportable security report** ready for documentation
5. **Validation that all 40 MCP tools** work correctly in real-world scenarios

This test will demonstrate TinyBrain's full capabilities while producing valuable security research output!
