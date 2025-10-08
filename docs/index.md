---
layout: default
title: Home
nav_order: 1
description: "TinyBrain - Security-Focused LLM Memory Storage MCP Server"
permalink: /
---

# TinyBrain 🧠

**Security-Focused LLM Memory Storage MCP Server**

TinyBrain is a comprehensive memory storage system designed specifically for security-focused tasks like security code review, penetration testing, and exploit development. It provides an LLM with persistent, intelligent memory capabilities through the Model Context Protocol (MCP).

## Key Features

### Security-Focused Design
- **Specialized Categories**: Vulnerability, exploit, payload, technique, tool, reference, context, hypothesis, evidence, recommendation
- **Priority & Confidence Tracking**: 0-10 priority levels and 0.0-1.0 confidence scores
- **Relationship Mapping**: Track dependencies, causes, mitigations, and exploit chains
- **Task Progress Tracking**: Multi-stage security task management
- **Standards Compliance**: Aligned with OWASP, CWE, NIST, and ISO security standards

### Intelligent Memory Management
- **Context-Aware Storage**: Automatically categorizes and prioritizes information
- **Advanced Search**: Semantic, exact, fuzzy, tag-based, and relationship-based search
- **Access Tracking**: Monitors which memories are most relevant and frequently accessed
- **Context Summaries**: Provides relevant memory summaries for current tasks

### High Performance & Reliability
- **SQLite Backend**: Fast, reliable, local storage with full-text search
- **Optimized Queries**: Indexed searches and efficient relationship traversal
- **Transaction Safety**: ACID compliance for data integrity
- **Concurrent Access**: Thread-safe operations for multiple LLM interactions

### AI-Enhanced Search & Intelligence
- **Semantic Search**: AI-powered memory search using embeddings for conceptual similarity
- **Embedding Generation**: Generate embeddings for text (placeholder for AI integration)
- **Similarity Calculation**: Calculate semantic similarity between embeddings
- **Future-Ready**: Complete foundation for OpenAI, Cohere, or local model integration

## Quick Start

Get started with TinyBrain in just a few minutes:

```bash
# Method 1: Clone and build locally (recommended)
git clone https://github.com/rainmana/tinybrain.git
cd tinybrain
make install

# Start the server
tinybrain
```

See the [Getting Started Guide](getting-started.md) for detailed installation and configuration instructions.

## Navigation

- **[Getting Started](getting-started.md)** - Installation, setup, and basic usage
- **[Configuration](configuration.md)** - Environment variables and database settings
- **[Architecture](architecture.md)** - System architecture and database schema
- **[Memory Model](memory-model.md)** - Understanding memory entries and relationships
- **[API Overview](api/overview.md)** - Complete MCP tool reference (40 tools)
- **[API Examples](api/examples.md)** - Practical JSON examples
- **[Workflows](workflows/security-assessment.md)** - Security assessment workflows
- **[Datasets](datasets/index.md)** - Security datasets and templates
- **[Advanced Features](advanced-features.md)** - Context snapshots and task tracking
- **[Integrations](integrations/ai-assistants.md)** - AI assistant integrations
- **[Contributing](contributing.md)** - Development and contribution guide
- **[Roadmap](roadmap.md)** - Future development plans
- **[FAQ](faq.md)** - Frequently asked questions

## 40 MCP Tools

TinyBrain provides a comprehensive set of 40 MCP tools for complete LLM memory management across:
- **Core Memory Operations** (8 tools)
- **Session & Task Management** (6 tools)
- **Advanced Memory Features** (8 tools)
- **Security Templates & Batch Operations** (6 tools)
- **Memory Lifecycle & Cleanup** (4 tools)
- **AI-Enhanced Search** (3 tools)
- **Real-Time Notifications** (4 tools)
- **System Monitoring** (1 tool)

See the [API Overview](api/overview.md) for the complete tool listing.

## Performance

Current benchmark figures from testing:
- **Memory Entry Creation**: ~1000 entries/second
- **Search Operations**: ~100 searches/second
- **Relationship Queries**: ~500 queries/second
- **Database Size**: ~1MB per 10,000 memory entries

## Standards & Compliance

TinyBrain's security patterns are aligned with industry-standard security frameworks:
- **[OWASP Top 10 2021](https://owasp.org/Top10/)** - Web Application Security Risks
- **[CWE (Common Weakness Enumeration)](https://cwe.mitre.org/)** - Software Weakness Classification
- **[NIST SP 800-115](https://csrc.nist.gov/publications/detail/sp/800-115/final)** - Technical Guide to Information Security Testing
- **[ISO 27001](https://www.iso.org/isoiec-27001-information-security.html)** - Information Security Management Systems
- **[PTES (Penetration Testing Execution Standard)](http://www.pentest-standard.org/)** - Penetration Testing Methodology

## License

TinyBrain is open source software. See the LICENSE file in the repository for details.
