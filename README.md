# TinyBrain 🧠

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![MCP Protocol](https://img.shields.io/badge/MCP-Protocol-green.svg)](https://modelcontextprotocol.io/)
[![Security Focused](https://img.shields.io/badge/Security-Focused-red.svg)](https://github.com/rainmana/tinybrain)
[![Version](https://img.shields.io/badge/version-v1.2.1-blue.svg)](https://github.com/rainmana/tinybrain/releases)

**Security-Focused LLM Memory Storage with Intelligence Gathering, Reverse Engineering, and MITRE ATT&CK Integration**

TinyBrain is a comprehensive memory storage system designed specifically for security professionals, penetration testers, and AI assistants working on offensive security tasks. It provides intelligent memory management, pattern recognition, and comprehensive intelligence gathering capabilities through the Model Context Protocol (MCP).

📖 **[View Complete Documentation](https://rainmana.github.io/tinybrain/)** | 🐛 [Report Issues](https://github.com/rainmana/tinybrain/issues) | 💬 [Discussions](https://github.com/rainmana/tinybrain/discussions)

## ✨ Key Features

### 🧠 Intelligence Gathering
- **OSINT**: Open Source Intelligence collection and analysis
- **HUMINT**: Human Intelligence gathering and social engineering assessment
- **SIGINT**: Signals Intelligence and communications analysis
- **GEOINT**: Geospatial Intelligence and location-based analysis
- **MASINT**: Measurement and Signature Intelligence
- **TECHINT**: Technical Intelligence and technology assessment
- **FININT**: Financial Intelligence and cryptocurrency tracking
- **CYBINT**: Cyber Intelligence and threat analysis

### 🔍 Reverse Engineering
- **Malware Analysis**: Static and dynamic malware analysis capabilities
- **Binary Analysis**: PE, ELF, Mach-O file format analysis
- **Vulnerability Research**: Fuzzing, exploit development, and vulnerability analysis
- **Protocol Analysis**: Network and application protocol reverse engineering
- **Code Analysis**: Source code and assembly analysis tools

### 🎯 MITRE ATT&CK Integration
- **Complete Framework**: All 14 Enterprise tactics and 200+ techniques
- **TTP Mapping**: Map findings to specific tactics, techniques, and procedures
- **Attack Chain Analysis**: Complete attack chain mapping and analysis
- **Threat Hunting**: Hunt for specific TTPs and attack patterns
- **Campaign Tracking**: Track attack campaigns and threat actor activities
- **Real-Time Data**: Live MITRE ATT&CK dataset with 823+ techniques and 14 tactics
- **Intelligent Querying**: Semantic search across attack techniques and procedures

### 🛡️ Security Patterns & Standards
- **CWE Integration**: Common Weakness Enumeration patterns and classifications
- **OWASP Compliance**: OWASP Top 10 2021 and testing guide integration
- **Multi-Language Support**: Security patterns for 10+ programming languages
- **Authorization Templates**: RBAC, ABAC, and DAC access control patterns
- **Standards Compliance**: NIST, ISO 27001, PTES, and industry standards
- **NVD Integration**: National Vulnerability Database with 314,835+ CVEs
- **OWASP Testing Guide**: Complete web application security testing procedures

### 🔬 Security Knowledge Hub
- **NVD Integration**: National Vulnerability Database with 314,835+ CVEs
- **MITRE ATT&CK**: Complete framework with 823+ techniques and 14 tactics
- **OWASP Testing Guide**: Comprehensive web application security testing procedures
- **Intelligent Retrieval**: Context-aware querying with summarization
- **Real-Time Updates**: Incremental data updates from official sources
- **Rate Limiting**: Respectful API usage with proper rate limiting

### 📊 Memory Management
- **30+ Memory Categories**: Comprehensive categorization for intelligence, reconnaissance, and analysis data
- **Intelligence Objects**: Threat actors, attack campaigns, IOCs, TTPs, patterns, and correlations
- **Context-Aware Storage**: Automatically categorizes and prioritizes information
- **Advanced Search**: Semantic, exact, fuzzy, tag-based, and relationship-based search
- **Access Tracking**: Monitors which memories are most relevant and frequently accessed
- **Context Summaries**: Provides relevant memory summaries for current tasks

### High Performance & Reliability
- **PocketBase Backend**: Single binary with embedded SQLite, REST API, and real-time capabilities (v1.2.1+)
- **Admin Dashboard**: Web-based interface for data management and visualization at http://127.0.0.1:8090/_/
- **MCP Endpoint**: Custom MCP protocol endpoint at http://127.0.0.1:8090/mcp
- **Real-time Updates**: Server-sent events for live memory updates
- **Optimized Queries**: Indexed searches and efficient relationship traversal
- **Transaction Safety**: ACID compliance for data integrity
- **Concurrent Access**: Thread-safe operations for multiple LLM interactions
- **Zero Configuration**: Works out of the box with minimal setup
- **Mock Implementation**: Current version (v1.2.1) provides mock responses; real database operations coming in future releases

### AI-Enhanced Search & Intelligence
- **Semantic Search**: AI-powered memory search using embeddings for conceptual similarity
- **Embedding Generation**: Generate embeddings for text (placeholder for AI integration)
- **Similarity Calculation**: Calculate semantic similarity between embeddings
- **Future-Ready**: Complete foundation for OpenAI, Cohere, or local model integration

### Real-Time Notifications & Alerts
- **Memory Notifications**: Real-time alerts for memory events and system activities
- **High Priority Alerts**: Notifications for high-priority memories (priority ≥8, confidence ≥0.8)
- **Duplicate Detection**: Alerts for potential duplicate memories with similarity scores
- **Cleanup Notifications**: Notifications for automated cleanup operations
- **Notification Management**: Mark notifications as read, filter by session, priority-based sorting

### Developer Experience
- **Simple Installation**: `go install github.com/rainmana/tinybrain/cmd/server@latest`
- **Comprehensive Logging**: Detailed logging with structured output
- **Extensive Testing**: Full test coverage for all MCP tool handlers
- **Docker Support**: Containerized deployment ready
- **40 MCP Tools**: Complete API for all memory management operations
- **PocketBase Integration**: Single binary with admin dashboard and REST API

## 🛠️ Complete MCP Tool Set (40 Tools)

TinyBrain provides a comprehensive set of 40 MCP tools for complete LLM memory management:

### Core Memory Operations (8 tools)
- `store_memory` - Store new memory entries
- `get_memory` - Retrieve memory by ID
- `search_memories` - Advanced search with multiple strategies
- `update_memory` - Update existing memory entries
- `delete_memory` - Delete memory entries
- `find_similar_memories` - Find similar memories by content
- `check_duplicates` - Check for duplicate memories
- `get_memory_stats` - Get comprehensive memory statistics

### Session & Task Management (6 tools)
- `create_session` - Create new security assessment sessions
- `get_session` - Retrieve session information
- `list_sessions` - List all sessions with filtering
- `create_task_progress` - Create task progress entries
- `update_task_progress` - Update task progress
- `list_task_progress` - List task progress entries

### Advanced Memory Features (8 tools)
- `create_relationship` - Create memory relationships
- `get_related_memories` - Get related memories
- `create_context_snapshot` - Create context snapshots
- `get_context_snapshot` - Retrieve context snapshots
- `list_context_snapshots` - List context snapshots
- `get_context_summary` - Get memory summaries for context
- `export_session_data` - Export session data
- `import_session_data` - Import session data

### Security Templates & Batch Operations (6 tools)
- `get_security_templates` - Get predefined security templates
- `create_memory_from_template` - Create memories from templates
- `batch_create_memories` - Bulk create memory entries
- `batch_update_memories` - Bulk update memory entries
- `batch_delete_memories` - Bulk delete memory entries
- `get_detailed_memory_info` - Get detailed memory debugging info

### Memory Lifecycle & Cleanup (4 tools)
- `cleanup_old_memories` - Age-based memory cleanup
- `cleanup_low_priority_memories` - Priority-based cleanup
- `cleanup_unused_memories` - Access-based cleanup
- `get_system_diagnostics` - Comprehensive system diagnostics

### AI-Enhanced Search (3 tools)
- `semantic_search` - AI-powered semantic search
- `generate_embedding` - Generate embeddings for text
- `calculate_similarity` - Calculate embedding similarity

### Real-Time Notifications (4 tools)
- `get_notifications` - Get notifications and alerts
- `mark_notification_read` - Mark notifications as read
- `check_high_priority_memories` - Check for high-priority alerts
- `check_duplicate_memories` - Check for duplicate alerts

### System Monitoring (1 tool)
- `health_check` - Perform system health checks

## 🛡️ Security Standards & Source Attribution

### **Standards Compliance**
TinyBrain's security patterns and vulnerability datasets are aligned with industry-standard security frameworks:

- **[OWASP Top 10 2021](https://owasp.org/Top10/)** - Web Application Security Risks
- **[CWE (Common Weakness Enumeration)](https://cwe.mitre.org/)** - Software Weakness Classification
- **[NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)** - Technical Guide to Information Security Testing
- **[ISO 27001](https://www.iso.org/isoiec-27001-information-security.html)** - Information Security Management Systems
- **[PTES (Penetration Testing Execution Standard)](http://www.pentest-standard.org/)** - Penetration Testing Methodology

### **Source Attribution**
Our security patterns and vulnerability datasets are based on authoritative sources:

- **[OWASP Code Review Guide](https://github.com/OWASP/www-project-code-review-guide)** - Comprehensive secure code review methodology
- **[OWASP Secure Coding Dojo](https://owasp.org/SecureCodingDojo/codereview101/)** - Interactive security code review training
- **[OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)** - Web application security testing methodology
- **[SANS Top 25 CWE](https://cwe.mitre.org/top25/)** - Most dangerous software errors
- **[NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)** - Cybersecurity risk management

### **Multi-Language Coverage**
Our security patterns cover 10 major programming languages with language-specific vulnerability patterns:

- **JavaScript/Node.js** - Web application security patterns
- **Python** - Backend and automation security patterns  
- **Java** - Enterprise application security patterns
- **C#/.NET** - Microsoft ecosystem security patterns
- **PHP** - Web application security patterns
- **Ruby** - Web framework security patterns
- **Go** - System and API security patterns
- **C/C++** - System-level security patterns
- **TypeScript** - Type-safe web application patterns
- **Rust** - Memory-safe system programming patterns

## 🚀 Quick Start

### Installation

#### Option 1: Pre-built Binaries (Recommended)

Download the latest release for your platform from [Releases](https://github.com/rainmana/tinybrain/releases):

- **macOS (Apple Silicon)**: `tinybrain_*_Darwin_arm64.tar.gz`
- **macOS (Intel)**: `tinybrain_*_Darwin_x86_64.tar.gz`
- **Linux (x86_64)**: `tinybrain_*_Linux_x86_64.tar.gz`
- **Linux (ARM64)**: `tinybrain_*_Linux_arm64.tar.gz`
- **Windows**: `tinybrain_*_Windows_x86_64.zip`

Extract and run:
```bash
# Extract (Linux/macOS)
tar -xzf tinybrain_*_*.tar.gz
cd tinybrain_*

# Make executable
chmod +x server

# Run
./server serve
```

#### Option 2: Install from Source with go install

```bash
# Install latest version
go install github.com/rainmana/tinybrain/cmd/server@latest

# Install specific version
go install github.com/rainmana/tinybrain/cmd/server@v1.2.2

# The binary will be installed as 'server' in your $GOPATH/bin or $GOBIN
# Make sure $GOPATH/bin or $GOBIN is in your PATH
```

#### Option 3: Clone and Build Locally

```bash
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
go build -o server ./cmd/server
```

#### Option 4: Docker

```bash
docker pull rainmana/tinybrain:latest
docker run -p 8090:8090 rainmana/tinybrain
```

### Basic Usage

```bash
# Start the server (uses ./pb_data by default for PocketBase)
server serve

# Or specify a custom data directory
server serve --dir /path/to/your/data

# Customize the port (default is 127.0.0.1:8090)
server serve --http=127.0.0.1:9000

# Or use environment variable
TINYBRAIN_HTTP=127.0.0.1:9000 server serve

# Combine options
server serve --dir ~/.tinybrain --http=0.0.0.0:8090

# Access admin dashboard
open http://127.0.0.1:8090/_/

# The MCP endpoint is available at:
# http://127.0.0.1:8090/mcp
```

**Important**: With PocketBase, the data directory structure is different from the previous SQLite-only version. PocketBase stores its data in a `pb_data` subdirectory within the specified directory (or `./pb_data` by default).

### PocketBase Integration Features

TinyBrain uses PocketBase as its backend (v1.2.1+), providing:

- **Single Binary**: Everything in one executable with zero configuration
- **Admin Dashboard**: Web interface at http://127.0.0.1:8090/_/ for data management and visualization
- **REST API**: Full REST API at http://127.0.0.1:8090/api/ for external integrations
- **Real-time Updates**: Server-sent events for live memory updates
- **Data Persistence**: All data persists across server restarts in `pb_data` directory
- **MCP Endpoint**: Custom MCP protocol endpoint at http://127.0.0.1:8090/mcp
- **Comprehensive Testing**: Full test coverage with all MCP tools verified
- **Mock Implementation**: Current version provides mock responses for all MCP tools, with real database operations to be implemented in future releases

### Intelligence Gathering Example

```go
// Create an OSINT intelligence session
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

### Troubleshooting Installation

If you encounter issues with `go install`, try these solutions:

```bash
# If you get "main redeclared" errors, ensure you're using v1.2.1 or later
# Older versions (v1.2.0) had duplicate files that caused build errors
go install github.com/rainmana/tinybrain/cmd/server@v1.2.1

# If you get authentication errors, use direct clone method
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
go build -o server ./cmd/server

# If repository is private, ensure you have access
git config --global url."git@github.com:".insteadOf "https://github.com/"

# For Go module proxy issues, use direct mode
GOPROXY=direct go install github.com/rainmana/tinybrain/cmd/server@latest

# For checksum database errors (temporary issue with new releases)
GOSUMDB=off go install github.com/rainmana/tinybrain/cmd/server@latest
```

**Common Issues**:
- **"main redeclared" error**: You're using an old version. Use `@v1.2.1` or `@latest`
- **Binary not found**: Ensure `$GOPATH/bin` or `$GOBIN` is in your PATH
- **Checksum errors**: Wait a few minutes after a new release, or temporarily use `GOSUMDB=off`

### MCP Client Configuration

Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "server",
      "args": ["serve"],
      "env": {
        "POCKETBASE_DATA_DIR": "~/.tinybrain"
      }
    }
  }
}
```

**Note**: The binary name is `server` (from `cmd/server`), not `tinybrain`. PocketBase will create a `pb_data` subdirectory in the specified data directory.

## 📚 Documentation

For complete documentation, API reference, and detailed guides, visit our comprehensive documentation site:

### 📖 **[Complete Documentation](https://rainmana.github.io/tinybrain/)**

The documentation includes:
- **[Getting Started](https://rainmana.github.io/tinybrain/getting-started/)** - Installation and basic usage
- **[Core Features](https://rainmana.github.io/tinybrain/core-features/)** - Memory management, sessions, and search
- **[Intelligence & Reconnaissance](https://rainmana.github.io/tinybrain/intelligence/)** - OSINT, HUMINT, SIGINT, and more
- **[Reverse Engineering](https://rainmana.github.io/tinybrain/reverse-engineering/)** - Malware analysis and vulnerability research
- **[Security Patterns](https://rainmana.github.io/tinybrain/security-patterns/)** - CWE, OWASP, and multi-language patterns
- **[Integration](https://rainmana.github.io/tinybrain/integration/)** - AI assistant integration and development setup
- **[API Reference](https://rainmana.github.io/tinybrain/api-reference/)** - Complete MCP tools and REST API documentation
- **[Contributing](https://rainmana.github.io/tinybrain/contributing/)** - Guidelines for contributors

### Quick API Reference

#### Session Management

**Task Types**: `security_review`, `penetration_test`, `exploit_dev`, `vulnerability_analysis`, `threat_modeling`, `incident_response`, `intelligence_analysis`

**Intelligence Types**: `osint`, `humint`, `sigint`, `geoint`, `masint`, `techint`, `finint`, `cybint`

#### List Sessions
```json
{
  "name": "list_sessions",
  "arguments": {
    "task_type": "security_review",
    "status": "active",
    "limit": 10
  }
}
```

### Memory Operations

#### Store Memory
```json
{
  "name": "store_memory",
  "arguments": {
    "session_id": "session_123",
    "title": "SQL Injection in Login Form",
    "content": "Found SQL injection vulnerability in username parameter of login form. Payload: ' OR 1=1--",
    "category": "vulnerability",
    "content_type": "text",
    "priority": 8,
    "confidence": 0.9,
    "tags": "[\"sql-injection\", \"authentication\", \"critical\"]",
    "source": "manual-testing"
  }
}
```

**Categories**: `finding`, `vulnerability`, `exploit`, `payload`, `technique`, `tool`, `reference`, `context`, `hypothesis`, `evidence`, `recommendation`, `note`

#### Search Memories
```json
{
  "name": "search_memories",
  "arguments": {
    "query": "SQL injection authentication",
    "session_id": "session_123",
    "search_type": "semantic",
    "categories": "[\"vulnerability\", \"exploit\"]",
    "min_priority": 7,
    "limit": 20
  }
}
```

**Search Types**: `semantic`, `exact`, `fuzzy`, `tag`, `category`, `relationship`

#### Get Related Memories
```json
{
  "name": "get_related_memories",
  "arguments": {
    "memory_id": "memory_456",
    "relationship_type": "exploits",
    "limit": 10
  }
}
```

### Relationship Management

#### Create Relationship
```json
{
  "name": "create_relationship",
  "arguments": {
    "source_memory_id": "memory_123",
    "target_memory_id": "memory_456",
    "relationship_type": "exploits",
    "strength": 0.8,
    "description": "SQL injection can be used to bypass authentication"
  }
}
```

**Relationship Types**: `depends_on`, `causes`, `mitigates`, `exploits`, `references`, `contradicts`, `supports`, `related_to`, `parent_of`, `child_of`

### Context Management

#### Get Context Summary
```json
{
  "name": "get_context_summary",
  "arguments": {
    "session_id": "session_123",
    "current_task": "Analyzing authentication vulnerabilities",
    "max_memories": 20
  }
}
```

## 🏗️ Architecture

TinyBrain is built with:
- **Go 1.24+** - High-performance backend
- **PocketBase v0.30.4** - Single binary with embedded SQLite, REST API, and real-time capabilities
- **MCP Protocol** - LLM integration standard with 40+ tools
- **MITRE ATT&CK** - Security framework integration
- **Jekyll** - Documentation site with Minimal theme

**Current Version**: v1.2.1 (PocketBase backend with mock MCP tool responses)

### PocketBase Integration Benefits

- **Single Binary Deployment**: No external dependencies, works anywhere Go runs
- **Embedded Database**: SQLite database embedded in the binary
- **Web Admin Interface**: Built-in dashboard for data management and visualization
- **REST API**: Full REST API for external integrations and automation
- **Real-time Capabilities**: Server-sent events for live updates
- **Zero Configuration**: Works out of the box with sensible defaults
- **Data Persistence**: All data automatically persisted across restarts

### Key Design Principles

1. **Security-First**: All data structures and operations designed for security tasks
2. **Intelligence-Focused**: Comprehensive intelligence gathering and analysis capabilities
3. **Performance**: Optimized queries and indexes for fast retrieval
4. **Flexibility**: Extensible schema and relationship system
5. **Reliability**: ACID transactions and data integrity checks
6. **Usability**: Simple API with comprehensive documentation

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
make bench

# Run specific test
go test -v ./internal/database -run TestNewDatabase
```

## 🐳 Docker Deployment

```bash
# Build Docker image
make docker-build

# Run container
docker run --rm -it \
  -v ~/.tinybrain:/app/data \
  tinybrain:latest
```

## 📊 Performance

### Benchmarks

- **Memory Entry Creation**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Relationship Queries**: ~500 queries/second
- **Database Size**: ~1MB per 10,000 memory entries

### Optimization Features

- **Connection Pooling**: Optimized for SQLite single-writer model
- **Index Strategy**: Comprehensive indexing for all query patterns
- **Full-Text Search**: FTS5 virtual tables for semantic search
- **Caching**: Access pattern tracking for intelligent caching

## 🔧 Configuration

### Environment Variables

- `POCKETBASE_DATA_DIR`: Path to PocketBase data directory (default: `./pb_data`)
- `TINYBRAIN_HTTP`: HTTP server bind address (default: `127.0.0.1:8090`, e.g., `127.0.0.1:9000` or `0.0.0.0:8090`)
- `TINYBRAIN_LOG_LEVEL`: Log level (debug, info, warn, error)

### PocketBase Configuration

PocketBase provides:
- **Embedded SQLite**: Database stored in `pb_data/data.db` within the data directory
- **Admin Dashboard**: Accessible at http://127.0.0.1:8090/_/ after first run
- **REST API**: Full REST API at http://127.0.0.1:8090/api/
- **Collections**: Automatically managed by PocketBase
- **Zero Configuration**: Works out of the box with sensible defaults

### Data Directory Structure

```
~/.tinybrain/
├── pb_data/
│   ├── data.db          # PocketBase SQLite database
│   ├── logs.db          # PocketBase logs
│   └── storage/         # File storage (if used)
└── ...
```

## 🛡️ Security Datasets & Templates

### **Comprehensive Security Patterns**
- **[Security Code Review Dataset](SECURITY_CODE_REVIEW_DATASET.md)** - OWASP Top 10 2021 patterns, CWE vulnerabilities, and exploitation techniques
- **[Multi-Language Security Patterns](MULTI_LANGUAGE_SECURITY_PATTERNS.md)** - Language-specific vulnerability patterns for 10 programming languages
- **[CWE Security Patterns](CWE_SECURITY_PATTERNS.md)** - CWE Top 25 Most Dangerous Software Errors with comprehensive vulnerability patterns
- **[CWE LLM Dataset](CWE_LLM_DATASET.json)** - LLM-optimized CWE dataset in structured JSON format for efficient consumption
- **[CWE TinyBrain Integration](CWE_TINYBRAIN_INTEGRATION.md)** - Integration guide for CWE dataset with TinyBrain memory system
- **[TinyBrain Security Templates](TINYBRAIN_SECURITY_TEMPLATES.md)** - Pre-configured memory templates for consistent security assessment storage

### **AI Assistant Configurations**
- **[Cursor Rules](.cursorrules)** - Security assessment rules for Cursor AI assistant
- **[Cline Rules](.clinerules)** - Code review and exploitation framework for Cline
- **[Roo Mode](.roo-mode)** - Penetration testing configuration for Roo AI assistant
- **[User Configuration Template](.cursorrules.user-template)** - Customizable user configuration template

### **Standards-Based Approach**
All security patterns are derived from authoritative sources and aligned with industry standards:

- **OWASP Top 10 2021** - Based on [OWASP Top 10](https://owasp.org/Top10/) web application security risks
- **CWE Patterns** - Derived from [Common Weakness Enumeration](https://cwe.mitre.org/) software weakness classification
- **Code Review Standards** - Aligned with [OWASP Code Review Guide](https://github.com/OWASP/www-project-code-review-guide)
- **Training Integration** - Compatible with [OWASP Secure Coding Dojo](https://owasp.org/SecureCodingDojo/codereview101/)
- **Testing Methodology** - Follows [NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final) security testing guidelines

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Development Setup

```bash
# Clone repository
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain

# Setup development environment
make dev-setup

# Run tests
make test

# Build
make build
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

### **Technical Dependencies**
- [mcp-go](https://github.com/mark3labs/mcp-go) for MCP server framework
- [go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite driver
- [charmbracelet/log](https://github.com/charmbracelet/log) for structured logging

### **Security Standards & Sources**
- [OWASP](https://owasp.org/) for security standards and vulnerability classifications
- [OWASP Code Review Guide](https://github.com/OWASP/www-project-code-review-guide) for secure code review methodology
- [OWASP Secure Coding Dojo](https://owasp.org/SecureCodingDojo/codereview101/) for interactive security training
- [CWE (Common Weakness Enumeration)](https://cwe.mitre.org/) for software weakness classification
- [NIST](https://www.nist.gov/) for cybersecurity frameworks and testing guidelines
- [SANS](https://www.sans.org/) for security research and training materials

## 📈 Roadmap

- [x] Intelligence gathering frameworks (OSINT, HUMINT, SIGINT, etc.)
- [x] MITRE ATT&CK integration
- [x] Reverse engineering capabilities
- [x] Enhanced memory categories
- [x] Comprehensive documentation site
- [ ] HTTP transport support
- [ ] Memory compression and archiving
- [ ] Advanced analytics and insights
- [ ] Multi-user support with access controls
- [ ] Plugin system for custom memory types
- [ ] Integration with popular security tools
- [ ] Web dashboard for memory visualization

---

**TinyBrain** - Making LLM memory storage intelligent, fast, and security-focused. 🧠🔒
