---
layout: default
title: TinyBrain
description: Security-focused LLM memory storage with intelligence gathering, reverse engineering, and MITRE ATT&CK integration
---

<!-- Force rebuild - Dark theme test -->

# TinyBrain

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![MCP Protocol](https://img.shields.io/badge/MCP-Protocol-green.svg)](https://modelcontextprotocol.io/)
[![Security Focused](https://img.shields.io/badge/Security-Focused-red.svg)](https://github.com/rainmana/tinybrain)

> Security-focused LLM memory storage with intelligence gathering, reverse engineering, and MITRE ATT&CK integration.

TinyBrain is a Model Context Protocol (MCP) server designed for security professionals, penetration testers, and AI assistants working on offensive security tasks. It provides intelligent memory management, pattern recognition, and comprehensive intelligence gathering capabilities.

## 🚀 Quick Start

```bash
# Install from source (recommended)
go install github.com/rainmana/tinybrain/cmd/server@latest

# Start the server
tinybrain serve --dir ~/.tinybrain

# Access admin dashboard
open http://127.0.0.1:8090/_/

# Or use Docker
docker build -t tinybrain .
docker run -p 8090:8090 tinybrain
```

## ✨ Key Features

- 🧠 **Intelligence Gathering**: OSINT, HUMINT, SIGINT, GEOINT, MASINT, TECHINT, FININT, CYBINT
- 🔍 **Reverse Engineering**: Malware analysis, binary analysis, vulnerability research, protocol analysis
- 🎯 **MITRE ATT&CK**: Complete framework integration with tactics, techniques, and procedures
- 🛡️ **Security Patterns**: CWE, OWASP, and multi-language vulnerability patterns
- 📊 **Memory Management**: 30+ memory categories for comprehensive security data organization
- 🔗 **MCP Protocol**: Seamless integration with AI assistants and LLMs
- 🔍 **Pattern Recognition**: Advanced insight mapping and correlation analysis
- 🎯 **Threat Intelligence**: Threat actor profiling, attack campaign tracking, IOC management
- 🚀 **PocketBase Integration**: Single binary with admin dashboard, REST API, and real-time capabilities
- 📈 **Comprehensive Testing**: 17/17 tests passing with full functionality verification

## 📚 Documentation

### 📖 **[Complete Documentation](documentation/)** - Comprehensive documentation index

**Quick Access:**
- [Getting Started](getting-started/) - Installation and basic usage
- [Core Features](core-features/) - Memory management, sessions, and search
- [Advanced Features](advanced-features/) - Advanced memory management and pattern recognition
- [Intelligence & Reconnaissance](intelligence/) - OSINT, HUMINT, SIGINT, and more
- [Reverse Engineering](reverse-engineering/) - Malware analysis and vulnerability research
- [Security Patterns](security-patterns/) - CWE, OWASP, and multi-language patterns
- [AI Integration](ai-integration/) - AI assistant integration and development setup
- [Authorization](authorization/) - Access control and authorization patterns
- [Integration](integration/) - General integration capabilities and tools
- [API Reference](api-reference/) - Complete API documentation

## 🛠️ Installation

### From Source
```bash
go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Docker
```bash
docker pull rainmana/tinybrain:latest
```

### Pre-built Binaries
Download from [Releases](https://github.com/rainmana/tinybrain/releases)

## 💡 Usage

### Basic Example
```go
// Create a new intelligence session
session := &Session{
    Name: "OSINT Intelligence Gathering",
    TaskType: "intelligence_analysis",
    IntelligenceType: "osint",
    Classification: "unclassified",
    ThreatLevel: "medium",
}

// Store intelligence findings
finding := &IntelligenceFinding{
    Title: "Social Media Intelligence",
    IntelligenceType: "osint",
    ThreatLevel: "medium",
    MITRETactic: "TA0001",
    MITRETechnique: "T1591",
    KillChainPhase: "reconnaissance",
}
```

### MCP Integration
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "mcp_tinybrain-mcp-server_create_session",
    "arguments": {
      "name": "Security Assessment",
      "task_type": "penetration_test",
      "intelligence_type": "osint"
    }
  }
}
```

## 🏗️ Architecture

TinyBrain is built with:
- **Go** - High-performance backend
- **PocketBase** - Single binary with embedded SQLite, REST API, and real-time capabilities
- **MCP Protocol** - LLM integration standard with 40+ tools
- **MITRE ATT&CK** - Security framework integration

### PocketBase Integration Benefits

- **Single Binary Deployment**: No external dependencies, works anywhere Go runs
- **Embedded Database**: SQLite database embedded in the binary
- **Web Admin Interface**: Built-in dashboard at http://127.0.0.1:8090/_/ for data management
- **REST API**: Full REST API at http://127.0.0.1:8090/api/ for external integrations
- **Real-time Capabilities**: Server-sent events for live updates
- **Zero Configuration**: Works out of the box with sensible defaults
- **Data Persistence**: All data automatically persisted across restarts

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test suite
go test -v ./internal/repository
```

## 📊 Performance

- **Memory Operations**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Intelligence Analysis**: ~100 analyses/second
- **Database Size**: ~1MB per 10,000 memory entries

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](contributing/) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](https://github.com/rainmana/tinybrain/blob/main/LICENSE) file for details.

## 🙏 Acknowledgments

- [MITRE ATT&CK](https://attack.mitre.org/) for the security framework
- [Model Context Protocol](https://modelcontextprotocol.io/) for LLM integration
- [OWASP](https://owasp.org/) for security patterns
- [CWE](https://cwe.mitre.org/) for vulnerability classification

## 📞 Support

- 📖 [Documentation](https://rainmana.github.io/tinybrain/)
- 🐛 [Issue Tracker](https://github.com/rainmana/tinybrain/issues)
- 💬 [Discussions](https://github.com/rainmana/tinybrain/discussions)

---

**TinyBrain** - Empowering security professionals with intelligent memory management and comprehensive intelligence gathering capabilities.
