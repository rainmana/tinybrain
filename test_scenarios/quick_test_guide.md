# TinyBrain Quick Test Guide

## 🚀 **5-Minute Test Setup**

### **Step 1: Open New Cursor Window**
- Create a new Cursor window
- Make sure TinyBrain is configured as an MCP server
- Ignore the "too many tools" warning (it's just informational)

### **Step 2: Start with This Prompt**
```
I'm starting a security code review of the "express" npm package. Please help me set up a comprehensive TinyBrain session for this assessment.

First, create a session with:
- Name: "Security Review: Express.js"
- Task Type: "security_review" 
- Description: "Comprehensive security review focusing on authentication, input validation, and data handling vulnerabilities"

Then set up task progress tracking for these stages:
1. Reconnaissance & Setup
2. Authentication Analysis  
3. Input Validation Review
4. Data Handling Assessment
5. Final Report Generation

Finally, get the available security templates so we can use them for common vulnerability patterns.

Let's start with the session creation and then move into the actual code analysis.
```

### **Step 3: Vulnerability Discovery**
```
Now let's analyze the Express.js codebase for security vulnerabilities. Please:

1. Review the main entry points and identify potential security issues
2. For each vulnerability found, store it in TinyBrain with:
   - Appropriate category (vulnerability, exploit, technique, tool, evidence, etc.)
   - Priority level (1-10) based on severity
   - Confidence level (0.0-1.0) based on certainty
   - Detailed description and impact assessment
   - Tags for easy categorization

3. Create relationships between related vulnerabilities
4. Update task progress as we complete each analysis phase

Focus on these vulnerability types:
- Authentication bypasses
- Input validation flaws
- Injection vulnerabilities
- Cross-Site Scripting (XSS)
- Security misconfigurations
- Sensitive data exposure

Store detailed findings that could be used in a real security report.
```

### **Step 4: Exploit Development**
```
Based on our vulnerability analysis, let's develop proof-of-concept exploits for the most critical issues. Please:

1. Search TinyBrain for our highest priority vulnerabilities (priority >= 8)
2. For the top 3 most critical findings, create working exploit code
3. Store each exploit in TinyBrain with relationships to the original vulnerabilities
4. Create a comprehensive security report summarizing all findings
5. Export the session data for documentation

The exploits should be realistic and demonstrate actual security impact.
```

## 🎯 **What to Expect**

### **TinyBrain Usage:**
- **Sessions**: 1 main session
- **Memory Entries**: 10-15 (vulnerabilities, exploits, techniques)
- **Relationships**: 5-8 (linking related findings)
- **Task Progress**: 5 stages with updates
- **Notifications**: 2-3 high-priority alerts

### **Success Indicators:**
- ✅ Session created successfully
- ✅ Memory entries stored with proper categorization
- ✅ Relationships created between findings
- ✅ Task progress updated
- ✅ Exploit code generated
- ✅ Session data exported

## 🔧 **Troubleshooting**

### **If TinyBrain isn't responding:**
1. Check that TinyBrain is running: `tinybrain` in terminal
2. Verify MCP configuration in Cursor
3. Try a simple health check: `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"health_check","arguments":{}}}`

### **If you get errors:**
1. Check the TinyBrain logs in terminal
2. Verify database permissions
3. Try restarting TinyBrain

### **If tools aren't working:**
1. Check that all 40 tools are available
2. Verify JSON formatting in requests
3. Check parameter requirements for each tool

## 📊 **Validation Checklist**

- [ ] Session created with proper metadata
- [ ] Task progress tracking set up
- [ ] Security templates retrieved
- [ ] Vulnerabilities found and stored
- [ ] Relationships created between findings
- [ ] Exploit code generated
- [ ] High-priority notifications triggered
- [ ] Session data exported successfully
- [ ] All 40 MCP tools functional

## 🎉 **Success!**

If you complete this test successfully, you'll have:
- A working TinyBrain session with real security findings
- Proof-of-concept exploit code
- Comprehensive relationships between vulnerabilities
- A complete security assessment workflow
- Validation that all TinyBrain features work correctly

This demonstrates TinyBrain's full capabilities in a real-world security assessment scenario!
