# TinyBrain v2.0 - Clean PocketBase Implementation

## Vision
A clean, modular, PocketBase-powered MCP server for security-focused LLM memory storage.

## Core Principles
- **KISS**: Keep it simple, stupid
- **Test-Driven**: Unit tests from day 1
- **Modular**: Easy to add data sources later
- **Clean**: No technical debt from v1.0

## Architecture

### Core Features (MVP)
1. **Session Management** - Create, get, list security assessment sessions
2. **Memory Storage** - Store, retrieve, search security findings
3. **Relationships** - Link vulnerabilities to exploits, techniques, etc.
4. **Context Management** - Snapshots for long-running assessments
5. **Task Progress** - Track multi-stage security tasks

### Data Sources (Manual Import)
- **OWASP Top 10** - Manual import of testing procedures
- **MITRE ATT&CK** - Manual import of techniques/tactics
- **NVD CVEs** - Manual import of critical CVEs

### Technology Stack
- **Backend**: PocketBase (embedded)
- **MCP**: Go MCP library (with proper docs)
- **Testing**: Go testing + MCP debugger
- **Data**: Manual JSON imports (no API complexity)

## Development Phases

### Phase 1: Core Memory Storage
- [ ] PocketBase setup
- [ ] Session management
- [ ] Memory CRUD operations
- [ ] Unit tests

### Phase 2: MCP Integration
- [ ] MCP protocol implementation
- [ ] Tool registration
- [ ] Error handling
- [ ] Integration tests

### Phase 3: Advanced Features
- [ ] Relationships
- [ ] Context snapshots
- [ ] Task progress
- [ ] Search functionality

### Phase 4: Data Import
- [ ] OWASP data import
- [ ] ATT&CK data import
- [ ] NVD data import
- [ ] Data validation

## Success Criteria
- [ ] All core features work
- [ ] 90%+ test coverage
- [ ] MCP debugger integration
- [ ] Performance benchmarks
- [ ] Clean, maintainable code
