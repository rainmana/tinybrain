# TinyBrain Cursor Setup Guide

## 🚀 **Quick Setup (5 minutes)**

### **Step 1: Copy the Rules File**
```bash
# Copy the main rules file to your project root
cp .cursorrules /path/to/your/project/

# Or copy to your global Cursor config
cp .cursorrules ~/.cursor/rules/
```

### **Step 2: Customize for Your Needs (Optional)**
```bash
# Copy the user template
cp .cursorrules.user-template .cursorrules.user

# Edit the user file to add your custom MCP servers and preferences
nano .cursorrules.user
```

### **Step 3: Configure MCP Servers in Cursor**
Add these to your Cursor MCP configuration:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": []
    },
    "clear-thought": {
      "command": "npx",
      "args": ["@clear-thought/mcp-server"]
    },
    "stochastic-thinking": {
      "command": "npx", 
      "args": ["@stochastic-thinking/mcp-server"]
    }
  }
}
```

## 🎯 **What This Gives You**

### **Enforced TinyBrain Usage:**
- ✅ **Mandatory session creation** for all security tasks
- ✅ **Required memory storage** with proper categorization
- ✅ **Relationship tracking** between findings
- ✅ **Progress monitoring** through assessment phases
- ✅ **Quality assurance** checklists

### **Prevents Common LLM Issues:**
- 🚫 **No more lost findings** - everything stored in TinyBrain
- 🚫 **No more inconsistent categorization** - enforced standards
- 🚫 **No more forgotten relationships** - mandatory linking
- 🚫 **No more incomplete reports** - comprehensive validation
- 🚫 **No more LLM loops** - structured reasoning tools

### **Customizable for Your Workflow:**
- 🔧 **Add your MCP servers** (clear-thought, stochastic-thinking, etc.)
- 🔧 **Customize security templates** for your domain
- 🔧 **Set your quality standards** and validation requirements
- 🔧 **Configure your assessment phases** and methodology

## 🧠 **For Claude 3.7+ Users**

### **Recommended MCP Servers:**
```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain"
    },
    "clear-thought": {
      "command": "npx",
      "args": ["@clear-thought/mcp-server"]
    },
    "stochastic-thinking": {
      "command": "npx",
      "args": ["@stochastic-thinking/mcp-server"]
    }
  }
}
```

### **Why These Help:**
- **clear-thought**: Prevents reasoning loops and provides structured thinking
- **stochastic-thinking**: Handles uncertainty and probabilistic analysis
- **tinybrain**: Core memory management and security assessment tracking

## 🎯 **Usage Examples**

### **Starting a Security Assessment:**
```
I'm starting a security code review of Express.js. Please help me set up a comprehensive TinyBrain session for this assessment.
```

### **The LLM will automatically:**
1. Create a TinyBrain session
2. Set up task progress tracking
3. Get security templates
4. Create initial context snapshot
5. Proceed with structured analysis

### **Developing Exploits:**
```
Based on our TinyBrain findings, please develop proof-of-concept exploits for the most critical vulnerabilities.
```

### **The LLM will automatically:**
1. Search TinyBrain for high-priority findings
2. Develop exploit code
3. Store exploits with relationships to vulnerabilities
4. Generate comprehensive reports
5. Export session data

## 🔧 **Customization Options**

### **Add Your MCP Servers:**
Edit `.cursorrules.user` and add:
```
# Your custom MCP servers
- "mcp_my-security-tool": For specialized security testing
- "mcp_my-analysis-tool": For custom analysis workflows
```

### **Customize Security Templates:**
Add your own vulnerability patterns:
```
{
  "name": "my_custom_vulnerability",
  "title": "Custom Vulnerability Pattern",
  "content": "Description of your custom vulnerability pattern...",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9
}
```

### **Set Your Quality Standards:**
```
QUALITY_CHECKS = {
    "min_findings_per_assessment": 10,
    "min_high_priority_findings": 3,
    "require_exploit_code": true,
    "require_remediation_steps": true
}
```

## 🚨 **Troubleshooting**

### **If TinyBrain isn't being used:**
1. Check that `.cursorrules` is in your project root
2. Verify MCP server configuration in Cursor
3. Restart Cursor after configuration changes
4. Check that TinyBrain is running: `tinybrain` in terminal

### **If you get MCP errors:**
1. Verify all MCP servers are installed and running
2. Check MCP server configuration syntax
3. Test individual MCP servers separately
4. Check Cursor logs for detailed error messages

### **If rules aren't being enforced:**
1. Ensure `.cursorrules` is in the correct location
2. Check file permissions
3. Restart Cursor
4. Verify the rules file syntax is correct

## 🎉 **Success Indicators**

After setup, you should see:
- ✅ **Automatic TinyBrain session creation** for security tasks
- ✅ **Structured memory storage** with proper categorization
- ✅ **Relationship tracking** between findings
- ✅ **Progress monitoring** through assessment phases
- ✅ **Quality validation** before task completion
- ✅ **Comprehensive reports** generated from TinyBrain data

## 🚀 **Next Steps**

1. **Test the setup** with a simple security assessment
2. **Customize the rules** to match your workflow
3. **Add your preferred MCP servers** for enhanced functionality
4. **Create custom security templates** for your domain
5. **Set your quality standards** and validation requirements

**You're now ready to conduct professional security assessments with TinyBrain!** 🧠🔒🚀
