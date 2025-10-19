# TinyBrain

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![MCP Protocol](https://img.shields.io/badge/MCP-Protocol-green.svg)](https://modelcontextprotocol.io/)
[![Security Focused](https://img.shields.io/badge/Security-Focused-red.svg)](https://github.com/rainmana/tinybrain)

> Security-focused LLM memory storage with intelligence gathering, reverse engineering, and MITRE ATT&CK integration.

TinyBrain is a Model Context Protocol (MCP) server designed for security professionals, penetration testers, and AI assistants working on offensive security tasks. It provides intelligent memory management, pattern recognition, and comprehensive intelligence gathering capabilities.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain

# Build and run
make build
make run

# Or use Docker
docker build -t tinybrain .
docker run -p 8080:8080 tinybrain
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

## 📚 Documentation

📖 **[Full Documentation Wiki](https://github.com/rainmana/tinybrain/wiki)**

- [Getting Started](https://github.com/rainmana/tinybrain/wiki/Getting-Started)
- [Intelligence & Reconnaissance](https://github.com/rainmana/tinybrain/wiki/Intelligence-&-Reconnaissance)
- [Reverse Engineering](https://github.com/rainmana/tinybrain/wiki/Reverse-Engineering)
- [Security Patterns](https://github.com/rainmana/tinybrain/wiki/Security-Patterns)
- [API Reference](https://github.com/rainmana/tinybrain/wiki/API-Reference)

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
- **SQLite** - Fast, reliable local storage
- **FTS5** - Full-text search capabilities
- **MCP Protocol** - LLM integration standard
- **MITRE ATT&CK** - Security framework integration

## 🔧 Configuration

```json
{
  "database": {
    "path": "./tinybrain.db",
    "max_connections": 10
  },
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "security": {
    "classification_levels": ["unclassified", "confidential", "secret"],
    "threat_levels": ["low", "medium", "high", "critical"]
  }
}
```

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

We welcome contributions! Please see our [Contributing Guide](https://github.com/rainmana/tinybrain/wiki/Contributing) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [MITRE ATT&CK](https://attack.mitre.org/) for the security framework
- [Model Context Protocol](https://modelcontextprotocol.io/) for LLM integration
- [OWASP](https://owasp.org/) for security patterns
- [CWE](https://cwe.mitre.org/) for vulnerability classification

## 📞 Support

- 📖 [Documentation Wiki](https://github.com/rainmana/tinybrain/wiki)
- 🐛 [Issue Tracker](https://github.com/rainmana/tinybrain/issues)
- 💬 [Discussions](https://github.com/rainmana/tinybrain/discussions)

---

**TinyBrain** - Empowering security professionals with intelligent memory management and comprehensive intelligence gathering capabilities.
