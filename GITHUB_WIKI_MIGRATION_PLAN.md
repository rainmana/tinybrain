# GitHub Wiki Migration Plan for TinyBrain

## Overview
This document outlines the migration of TinyBrain documentation from individual markdown files to a structured GitHub wiki, improving organization and discoverability.

## Current Documentation Analysis

### Core Documentation (Keep in Main README)
- `README.md` - Main project overview and quick start
- `LICENSE` - Project license
- `Dockerfile` - Container setup
- `Makefile` - Build commands

### Wiki Migration Categories

#### 1. **Getting Started** (Wiki Home + Quick Start)
- `examples/basic_usage.md` → Wiki: "Getting Started"
- `CURSOR_SETUP.md` → Wiki: "Development Setup"
- `config.example.json` → Wiki: "Configuration"

#### 2. **Core Features** (Wiki Section)
- `ADVANCED_FEATURES.md` → Wiki: "Advanced Features"
- `ENHANCED_MEMORY_CATEGORIES.md` → Wiki: "Memory Categories"
- `TINYBRAIN_SECURITY_TEMPLATES.md` → Wiki: "Security Templates"

#### 3. **Intelligence & Reconnaissance** (Wiki Section)
- `INTELLIGENCE_RECON_FRAMEWORK.md` → Wiki: "Intelligence Framework"
- `INTELLIGENCE_SECURITY_TEMPLATES.md` → Wiki: "Intelligence Templates"
- `MITRE_ATTACK_INTEGRATION.md` → Wiki: "MITRE ATT&CK Integration"
- `TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md` → Wiki: "Intelligence Enhancements"

#### 4. **Reverse Engineering** (Wiki Section)
- `REVERSE_ENGINEERING_FRAMEWORK.md` → Wiki: "Reverse Engineering"
- `INSIGHT_MAPPING_FRAMEWORK.md` → Wiki: "Insight Mapping"

#### 5. **Security Patterns** (Wiki Section)
- `CWE_SECURITY_PATTERNS.md` → Wiki: "CWE Security Patterns"
- `CWE_TINYBRAIN_INTEGRATION.md` → Wiki: "CWE Integration"
- `MULTI_LANGUAGE_SECURITY_PATTERNS.md` → Wiki: "Multi-Language Patterns"
- `ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md` → Wiki: "Language Library Patterns"
- `ENHANCED_AUTHORIZATION_TEMPLATES.md` → Wiki: "Authorization Templates"
- `SECURITY_CODE_REVIEW_DATASET.md` → Wiki: "Security Datasets"

#### 6. **Integration & Development** (Wiki Section)
- `AI_ASSISTANT_INTEGRATION.md` → Wiki: "AI Assistant Integration"
- `INTEGRATION_TEST_RESULTS.md` → Wiki: "Test Results"
- `FRAMEWORK_LIBRARY_PATTERNS.json` → Wiki: "Framework Patterns"

## GitHub Wiki Structure

### Home Page
```
# TinyBrain - Security-Focused LLM Memory Storage

TinyBrain is a Model Context Protocol (MCP) server designed for security professionals, penetration testers, and AI assistants working on offensive security tasks.

## Quick Start
[Link to Getting Started page]

## Key Features
- **Intelligence Gathering**: OSINT, HUMINT, SIGINT, and more
- **Reverse Engineering**: Malware analysis, binary analysis, vulnerability research
- **MITRE ATT&CK Integration**: Complete framework support
- **Security Patterns**: CWE, OWASP, and multi-language patterns
- **Memory Management**: 30+ memory categories for security data

## Documentation
- [Getting Started](Getting-Started)
- [Core Features](Core-Features)
- [Intelligence & Reconnaissance](Intelligence-&-Reconnaissance)
- [Reverse Engineering](Reverse-Engineering)
- [Security Patterns](Security-Patterns)
- [Integration & Development](Integration-&-Development)
- [API Reference](API-Reference)
- [Contributing](Contributing)
```

### Page Structure

#### 1. Getting Started
- Installation
- Basic Usage
- Configuration
- Development Setup

#### 2. Core Features
- Memory Management
- Session Management
- Search Capabilities
- Advanced Features
- Security Templates

#### 3. Intelligence & Reconnaissance
- Intelligence Framework
- OSINT/HUMINT/SIGINT
- MITRE ATT&CK Integration
- Intelligence Templates
- Threat Actor Analysis
- Attack Campaign Tracking

#### 4. Reverse Engineering
- Malware Analysis
- Binary Analysis
- Vulnerability Research
- Protocol Analysis
- Insight Mapping

#### 5. Security Patterns
- CWE Security Patterns
- Multi-Language Patterns
- OWASP Integration
- Authorization Templates
- Security Datasets

#### 6. Integration & Development
- AI Assistant Integration
- MCP Protocol
- API Reference
- Testing
- Contributing

## Migration Steps

### Step 1: Enable GitHub Wiki
1. Go to your GitHub repository
2. Click on "Wiki" tab
3. Click "Create the first page"
4. Set up the home page structure

### Step 2: Create Wiki Pages
1. Create each page using the structure above
2. Copy content from existing markdown files
3. Update internal links to use wiki format
4. Add navigation between pages

### Step 3: Update Main README
1. Keep only essential information in main README
2. Add links to wiki pages
3. Include quick start and key features
4. Add badges and status information

### Step 4: Clean Up Repository
1. Move detailed documentation to wiki
2. Keep only essential files in main repo
3. Update any references to moved files
4. Add wiki links to README

## Recommended Main README Content

```markdown
# TinyBrain

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![MCP Protocol](https://img.shields.io/badge/MCP-Protocol-green.svg)](https://modelcontextprotocol.io/)

> Security-focused LLM memory storage with intelligence gathering, reverse engineering, and MITRE ATT&CK integration.

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/tinybrain.git
cd tinybrain

# Build and run
make build
make run

# Or use Docker
docker build -t tinybrain .
docker run -p 8080:8080 tinybrain
```

## Key Features

- 🧠 **Intelligence Gathering**: OSINT, HUMINT, SIGINT, and more
- 🔍 **Reverse Engineering**: Malware analysis, binary analysis, vulnerability research  
- 🎯 **MITRE ATT&CK**: Complete framework integration
- 🛡️ **Security Patterns**: CWE, OWASP, and multi-language patterns
- 📊 **Memory Management**: 30+ memory categories for security data
- 🔗 **MCP Protocol**: Seamless AI assistant integration

## Documentation

📚 **[Full Documentation Wiki](https://github.com/yourusername/tinybrain/wiki)**

- [Getting Started](https://github.com/yourusername/tinybrain/wiki/Getting-Started)
- [Intelligence & Reconnaissance](https://github.com/yourusername/tinybrain/wiki/Intelligence-&-Reconnaissance)
- [Reverse Engineering](https://github.com/yourusername/tinybrain/wiki/Reverse-Engineering)
- [API Reference](https://github.com/yourusername/tinybrain/wiki/API-Reference)

## Installation

```bash
go install github.com/yourusername/tinybrain/cmd/server@latest
```

## Usage

```go
// Create a new session
session := &Session{
    Name: "Security Assessment",
    TaskType: "penetration_test",
    IntelligenceType: "osint",
}

// Store intelligence findings
finding := &IntelligenceFinding{
    Title: "OSINT Finding: Social Media Intelligence",
    IntelligenceType: "osint",
    ThreatLevel: "medium",
    MITRETactic: "TA0001",
}
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/yourusername/tinybrain/wiki/Contributing) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.
```

## Benefits of Wiki Migration

1. **Better Organization**: Structured navigation and categorization
2. **Improved Discoverability**: Easy to find specific information
3. **Cleaner Repository**: Main repo focuses on code, not documentation
4. **Collaborative Editing**: Multiple contributors can edit wiki pages
5. **Version Control**: GitHub tracks wiki changes
6. **Search**: Built-in wiki search functionality
7. **Mobile Friendly**: Better mobile reading experience

## Files to Keep in Main Repository

- `README.md` (simplified)
- `LICENSE`
- `Dockerfile`
- `Makefile`
- `config.example.json`
- `go.mod`, `go.sum`
- All source code files
- Test files
- Build scripts

## Files to Move to Wiki

- All `*.md` files except `README.md`
- `*.json` files that are documentation
- Example files (move to wiki or examples directory)

## Next Steps

1. **Enable GitHub Wiki** for your repository
2. **Create the home page** with the structure above
3. **Migrate content** page by page
4. **Update main README** to be concise and link to wiki
5. **Test all links** and navigation
6. **Clean up repository** by removing moved files
7. **Update any references** in code or other files

This migration will significantly improve the organization and usability of your documentation while keeping your main repository clean and focused on the code.
