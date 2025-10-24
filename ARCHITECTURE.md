# TinyBrain v2.0 Architecture Plan

## Decision: **Layered Architecture** (Highest Score)

Based on weighted scoring analysis:
- **Simplicity**: 8/10 - Clear separation of concerns
- **Testability**: 9/10 - Each layer can be tested independently  
- **Performance**: 7/10 - Good for security tasks, some overhead
- **Maintainability**: 8/10 - Easy to modify and extend
- **Future-Proof**: 7/10 - Flexible for new data sources

## Architecture Layers

```
┌─────────────────────────────────────┐
│           MCP Protocol Layer         │  ← JSON-RPC 2.0, Tool Registration
├─────────────────────────────────────┤
│         Business Logic Layer        │  ← Session, Memory, Relationships
├─────────────────────────────────────┤
│          Data Access Layer          │  ← Repository Pattern
├─────────────────────────────────────┤
│         PocketBase Backend          │  ← Embedded Database
└─────────────────────────────────────┘
```

## Core Components

### 1. MCP Protocol Layer
**File**: `cmd/server/main.go`
- MCP server initialization
- Tool registration
- Request/response handling
- Error handling

### 2. Business Logic Layer
**Files**: `internal/services/`
- `session_service.go` - Session management
- `memory_service.go` - Memory operations
- `relationship_service.go` - Memory relationships
- `context_service.go` - Context snapshots
- `task_service.go` - Task progress

### 3. Data Access Layer
**Files**: `internal/repository/`
- `session_repository.go` - Session CRUD
- `memory_repository.go` - Memory CRUD
- `relationship_repository.go` - Relationship CRUD
- `context_repository.go` - Context CRUD
- `task_repository.go` - Task CRUD

### 4. PocketBase Backend
**Files**: `internal/database/`
- `pocketbase_client.go` - PocketBase client wrapper
- `collections.go` - Collection definitions
- `migrations.go` - Database migrations

## Data Models

### Core Collections
```go
// Sessions
type Session struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    TaskType    string    `json:"task_type"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Memory Entries
type MemoryEntry struct {
    ID          string    `json:"id"`
    SessionID   string    `json:"session_id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Category    string    `json:"category"`
    Priority    int       `json:"priority"`
    Confidence  float64   `json:"confidence"`
    Tags        []string  `json:"tags"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Relationships
type Relationship struct {
    ID          string    `json:"id"`
    SourceID    string    `json:"source_id"`
    TargetID    string    `json:"target_id"`
    Type        string    `json:"type"`
    Strength    float64   `json:"strength"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Data Source Collections (Manual Import)
```go
// OWASP Procedures
type OWASPProcedure struct {
    ID              string   `json:"id"`
    Category        string   `json:"category"`
    VulnerabilityType string `json:"vulnerability_type"`
    Severity        string   `json:"severity"`
    Description     string   `json:"description"`
    Steps           []string `json:"steps"`
    Tools           []string `json:"tools"`
    References      []string `json:"references"`
}

// ATT&CK Techniques
type ATTACKTechnique struct {
    ID          string   `json:"id"`
    TechniqueID string   `json:"technique_id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Tactic      string   `json:"tactic"`
    Platform    []string `json:"platform"`
    DataSources []string `json:"data_sources"`
}

// NVD CVEs
type NVDEntry struct {
    ID          string    `json:"id"`
    CVEID       string    `json:"cve_id"`
    Description string    `json:"description"`
    Severity    string    `json:"severity"`
    CVSS        float64   `json:"cvss"`
    Published   time.Time `json:"published"`
    Modified    time.Time `json:"modified"`
}
```

## MCP Tools (Core Features)

### Session Management
- `create_session` - Create new security assessment session
- `get_session` - Retrieve session details
- `list_sessions` - List all sessions
- `update_session` - Update session status

### Memory Operations
- `store_memory` - Store security finding
- `get_memory` - Retrieve memory entry
- `search_memories` - Search with filters
- `update_memory` - Update memory entry
- `delete_memory` - Remove memory entry

### Relationships
- `create_relationship` - Link memories
- `get_relationships` - Get memory relationships
- `delete_relationship` - Remove relationship

### Context Management
- `create_context_snapshot` - Save current context
- `get_context_snapshot` - Retrieve context
- `get_context_summary` - Get context summary

### Task Progress
- `create_task_progress` - Start new task
- `update_task_progress` - Update task status
- `get_task_progress` - Get task details
- `list_task_progress` - List all tasks

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)
1. ✅ PocketBase setup
2. ✅ Collection definitions
3. ✅ Basic repository pattern
4. ✅ Unit test framework

### Phase 2: Core Features (Week 2)
1. ✅ Session management
2. ✅ Memory storage
3. ✅ Basic MCP integration
4. ✅ Unit tests

### Phase 3: Advanced Features (Week 3)
1. ✅ Relationships
2. ✅ Context snapshots
3. ✅ Task progress
4. ✅ Search functionality

### Phase 4: Data Import (Week 4)
1. ✅ OWASP data import
2. ✅ ATT&CK data import
3. ✅ NVD data import
4. ✅ Integration tests

## Success Metrics

### Functional Requirements
- [ ] All MCP tools work correctly
- [ ] Data persists in PocketBase
- [ ] Search returns relevant results
- [ ] Relationships are maintained

### Non-Functional Requirements
- [ ] 90%+ test coverage
- [ ] < 100ms response time for queries
- [ ] Handles 10,000+ memory entries
- [ ] Clean, maintainable code

### Quality Gates
- [ ] All unit tests pass
- [ ] Integration tests pass
- [ ] MCP debugger works
- [ ] Performance benchmarks met
- [ ] Code review approved

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Database**: PocketBase (embedded)
- **MCP Library**: github.com/mark3labs/mcp-go
- **Testing**: Go testing + testify

### Development Tools
- **MCP Debugger**: For protocol testing
- **Context7**: For documentation
- **Clear Thought**: For reasoning
- **Stochastic Thinking**: For decisions

### Data Sources
- **OWASP**: Manual JSON import
- **ATT&CK**: Manual JSON import  
- **NVD**: Manual JSON import
- **No API calls**: KISS principle

## Risk Mitigation

### Technical Risks
- **PocketBase Learning Curve**: Use official docs, start simple
- **MCP Protocol Complexity**: Use proven library, follow examples
- **Data Import Issues**: Start with small datasets, validate format

### Project Risks
- **Scope Creep**: Stick to core features first
- **Over-Engineering**: KISS principle, add complexity gradually
- **Testing Gaps**: Write tests from day 1

## Next Steps

1. ✅ Create project structure
2. ✅ Set up PocketBase
3. ✅ Implement basic MCP server
4. ✅ Add first tool (create_session)
5. ✅ Write unit tests
6. ✅ Iterate and improve
