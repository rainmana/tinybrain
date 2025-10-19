---
layout: default
title: API Reference
permalink: /api-reference/
---

# API Reference

TinyBrain provides comprehensive APIs for integration with security tools, AI assistants, and custom applications.

## MCP Tools

### Session Management

#### create_session
Create a new security assessment session.

**Parameters:**
- `name` (string): Session name
- `description` (string, optional): Session description
- `task_type` (string): Type of security task
- `intelligence_type` (string, optional): Intelligence gathering type
- `classification` (string, optional): Security classification level
- `threat_level` (string, optional): Threat level assessment

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_create_session",
  "arguments": {
    "name": "OSINT Intelligence Gathering",
    "description": "Gathering open source intelligence on target organization",
    "task_type": "intelligence_analysis",
    "intelligence_type": "osint",
    "classification": "unclassified",
    "threat_level": "medium"
  }
}
```

#### get_session
Retrieve session information.

**Parameters:**
- `session_id` (string): Session identifier

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_get_session",
  "arguments": {
    "session_id": "session-123"
  }
}
```

#### list_sessions
List all sessions with optional filtering.

**Parameters:**
- `task_type` (string, optional): Filter by task type
- `status` (string, optional): Filter by session status
- `limit` (integer, optional): Maximum number of sessions to return

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_list_sessions",
  "arguments": {
    "task_type": "intelligence_analysis",
    "status": "active",
    "limit": 10
  }
}
```

### Memory Management

#### create_memory
Store a new memory entry.

**Parameters:**
- `session_id` (string): Session identifier
- `title` (string): Memory entry title
- `content` (string): Memory entry content
- `category` (string): Memory category
- `priority` (integer, optional): Priority level (0-10)
- `confidence` (float, optional): Confidence level (0.0-1.0)
- `tags` (array, optional): List of tags
- `source` (string, optional): Information source
- `intelligence_type` (string, optional): Intelligence type
- `threat_level` (string, optional): Threat level
- `mitre_tactic` (string, optional): MITRE ATT&CK tactic
- `mitre_technique` (string, optional): MITRE ATT&CK technique

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_create_memory",
  "arguments": {
    "session_id": "session-123",
    "title": "Social Media Intelligence Finding",
    "content": "Analysis reveals suspicious activity on social media platforms",
    "category": "intelligence",
    "priority": 7,
    "confidence": 0.8,
    "tags": ["osint", "social-media", "suspicious"],
    "source": "Social Media Platforms",
    "intelligence_type": "osint",
    "threat_level": "medium",
    "mitre_tactic": "TA0001",
    "mitre_technique": "T1591"
  }
}
```

#### search_memory
Search memory entries.

**Parameters:**
- `query` (string): Search query
- `session_id` (string, optional): Filter by session
- `category` (string, optional): Filter by category
- `search_type` (string, optional): Search type (exact, fuzzy, semantic)
- `limit` (integer, optional): Maximum number of results

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_search_memory",
  "arguments": {
    "query": "malware analysis",
    "session_id": "session-123",
    "category": "intelligence",
    "search_type": "semantic",
    "limit": 20
  }
}
```

#### get_memory
Retrieve a specific memory entry.

**Parameters:**
- `memory_id` (string): Memory entry identifier

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_get_memory",
  "arguments": {
    "memory_id": "memory-456"
  }
}
```

### Relationship Management

#### create_relationship
Create a relationship between memory entries.

**Parameters:**
- `source_entry_id` (string): Source memory entry ID
- `target_entry_id` (string): Target memory entry ID
- `relationship_type` (string): Type of relationship
- `strength` (float, optional): Relationship strength (0.0-1.0)
- `description` (string, optional): Relationship description

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_create_relationship",
  "arguments": {
    "source_entry_id": "memory-456",
    "target_entry_id": "memory-789",
    "relationship_type": "causes",
    "strength": 0.9,
    "description": "SQL injection vulnerability causes data breach"
  }
}
```

#### get_related_entries
Get entries related to a specific memory entry.

**Parameters:**
- `memory_id` (string): Memory entry identifier
- `relationship_type` (string, optional): Filter by relationship type
- `limit` (integer, optional): Maximum number of results

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_get_related_entries",
  "arguments": {
    "memory_id": "memory-456",
    "relationship_type": "causes",
    "limit": 10
  }
}
```

### Context Management

#### create_context_snapshot
Save current LLM context state.

**Parameters:**
- `session_id` (string): Session identifier
- `name` (string): Snapshot name
- `description` (string, optional): Snapshot description
- `context_data` (string): JSON-encoded context data

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_create_context_snapshot",
  "arguments": {
    "session_id": "session-123",
    "name": "Mid-Assessment Snapshot",
    "description": "Context state during mid-assessment phase",
    "context_data": "{\"current_findings\": [], \"active_targets\": []}"
  }
}
```

#### get_context_snapshot
Retrieve a context snapshot.

**Parameters:**
- `snapshot_id` (string): Snapshot identifier

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_get_context_snapshot",
  "arguments": {
    "snapshot_id": "snapshot-789"
  }
}
```

### Task Progress

#### create_task_progress
Create a new task progress entry.

**Parameters:**
- `session_id` (string): Session identifier
- `task_name` (string): Task name
- `stage` (string): Current stage
- `status` (string): Task status
- `progress_percentage` (integer, optional): Progress percentage (0-100)
- `notes` (string, optional): Task notes

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_create_task_progress",
  "arguments": {
    "session_id": "session-123",
    "task_name": "OSINT Intelligence Gathering",
    "stage": "reconnaissance",
    "status": "in_progress",
    "progress_percentage": 45,
    "notes": "Completed social media analysis, starting dark web monitoring"
  }
}
```

#### update_task_progress
Update task progress.

**Parameters:**
- `task_id` (string): Task identifier
- `stage` (string, optional): New stage
- `status` (string, optional): New status
- `progress_percentage` (integer, optional): New progress percentage
- `notes` (string, optional): Updated notes

**Example:**
```json
{
  "name": "mcp_tinybrain-mcp-server_update_task_progress",
  "arguments": {
    "task_id": "task-123",
    "stage": "analysis",
    "status": "in_progress",
    "progress_percentage": 75,
    "notes": "Completed data collection, starting analysis phase"
  }
}
```

## REST API

### Authentication

#### POST /auth
Authenticate and receive access token.

**Request:**
```json
{
  "username": "user",
  "password": "password"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### Sessions

#### POST /sessions
Create a new session.

**Request:**
```json
{
  "name": "Security Assessment",
  "description": "Penetration testing session",
  "task_type": "penetration_test",
  "intelligence_type": "osint",
  "classification": "unclassified",
  "threat_level": "medium"
}
```

**Response:**
```json
{
  "id": "session-123",
  "name": "Security Assessment",
  "description": "Penetration testing session",
  "task_type": "penetration_test",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /sessions
List sessions with optional filtering.

**Query Parameters:**
- `task_type` (string): Filter by task type
- `status` (string): Filter by status
- `limit` (integer): Maximum results (default: 50)
- `offset` (integer): Results offset (default: 0)

**Response:**
```json
{
  "sessions": [
    {
      "id": "session-123",
      "name": "Security Assessment",
      "task_type": "penetration_test",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

#### GET /sessions/{id}
Get session details.

**Response:**
```json
{
  "id": "session-123",
  "name": "Security Assessment",
  "description": "Penetration testing session",
  "task_type": "penetration_test",
  "intelligence_type": "osint",
  "classification": "unclassified",
  "threat_level": "medium",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "metadata": {}
}
```

### Memory Entries

#### POST /memory
Create a memory entry.

**Request:**
```json
{
  "session_id": "session-123",
  "title": "SQL Injection Vulnerability",
  "content": "Found SQL injection in login form",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9,
  "tags": ["sql-injection", "critical", "web-app"],
  "source": "Manual Testing",
  "mitre_tactic": "TA0001",
  "mitre_technique": "T1190"
}
```

**Response:**
```json
{
  "id": "memory-456",
  "session_id": "session-123",
  "title": "SQL Injection Vulnerability",
  "content": "Found SQL injection in login form",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9,
  "tags": ["sql-injection", "critical", "web-app"],
  "source": "Manual Testing",
  "mitre_tactic": "TA0001",
  "mitre_technique": "T1190",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /memory/search
Search memory entries.

**Query Parameters:**
- `query` (string): Search query
- `session_id` (string): Filter by session
- `category` (string): Filter by category
- `search_type` (string): Search type (exact, fuzzy, semantic)
- `limit` (integer): Maximum results (default: 20)
- `offset` (integer): Results offset (default: 0)

**Response:**
```json
{
  "entries": [
    {
      "id": "memory-456",
      "title": "SQL Injection Vulnerability",
      "content": "Found SQL injection in login form",
      "category": "vulnerability",
      "priority": 8,
      "confidence": 0.9,
      "relevance_score": 0.95
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

#### GET /memory/{id}
Get memory entry details.

**Response:**
```json
{
  "id": "memory-456",
  "session_id": "session-123",
  "title": "SQL Injection Vulnerability",
  "content": "Found SQL injection in login form",
  "category": "vulnerability",
  "priority": 8,
  "confidence": 0.9,
  "tags": ["sql-injection", "critical", "web-app"],
  "source": "Manual Testing",
  "mitre_tactic": "TA0001",
  "mitre_technique": "T1190",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "accessed_at": "2024-01-01T00:00:00Z",
  "access_count": 1
}
```

## GraphQL API

### Schema

```graphql
type Session {
  id: ID!
  name: String!
  description: String
  taskType: String!
  intelligenceType: String
  classification: String
  threatLevel: String
  status: String!
  createdAt: String!
  updatedAt: String!
  memoryEntries: [MemoryEntry!]!
  taskProgress: [TaskProgress!]!
}

type MemoryEntry {
  id: ID!
  sessionId: ID!
  title: String!
  content: String!
  category: String!
  priority: Int
  confidence: Float
  tags: [String!]
  source: String
  intelligenceType: String
  threatLevel: String
  mitreTactic: String
  mitreTechnique: String
  createdAt: String!
  updatedAt: String!
  relatedEntries: [MemoryEntry!]!
}

type TaskProgress {
  id: ID!
  sessionId: ID!
  taskName: String!
  stage: String!
  status: String!
  progressPercentage: Int
  notes: String
  createdAt: String!
  updatedAt: String!
}

type Query {
  session(id: ID!): Session
  sessions(taskType: String, status: String, limit: Int, offset: Int): [Session!]!
  memoryEntry(id: ID!): MemoryEntry
  searchMemory(query: String!, sessionId: ID, category: String, searchType: String, limit: Int): [MemoryEntry!]!
  taskProgress(sessionId: ID!): [TaskProgress!]!
}

type Mutation {
  createSession(input: CreateSessionInput!): Session!
  updateSession(id: ID!, input: UpdateSessionInput!): Session!
  createMemoryEntry(input: CreateMemoryEntryInput!): MemoryEntry!
  updateMemoryEntry(id: ID!, input: UpdateMemoryEntryInput!): MemoryEntry!
  createRelationship(input: CreateRelationshipInput!): Relationship!
  createTaskProgress(input: CreateTaskProgressInput!): TaskProgress!
  updateTaskProgress(id: ID!, input: UpdateTaskProgressInput!): TaskProgress!
}
```

### Example Queries

#### Get Session with Memory Entries
```graphql
query GetSession($id: ID!) {
  session(id: $id) {
    id
    name
    taskType
    intelligenceType
    memoryEntries {
      id
      title
      category
      priority
      threatLevel
      mitreTactic
      mitreTechnique
    }
  }
}
```

#### Search Memory Entries
```graphql
query SearchMemory($query: String!, $sessionId: ID) {
  searchMemory(query: $query, sessionId: $sessionId) {
    id
    title
    content
    category
    priority
    confidence
    relevanceScore
  }
}
```

#### Create Memory Entry
```graphql
mutation CreateMemoryEntry($input: CreateMemoryEntryInput!) {
  createMemoryEntry(input: $input) {
    id
    title
    category
    priority
    confidence
    createdAt
  }
}
```

## Error Handling

### HTTP Status Codes
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `422 Unprocessable Entity`: Validation error
- `500 Internal Server Error`: Server error

### Error Response Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "priority",
      "issue": "Priority must be between 0 and 10"
    }
  }
}
```

## Rate Limiting

### Limits
- **MCP Tools**: 100 requests per minute per session
- **REST API**: 1000 requests per hour per API key
- **GraphQL**: 100 queries per minute per user

### Headers
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Authentication

### API Keys
```bash
curl -H "Authorization: Bearer your-api-key" \
  https://api.tinybrain.com/sessions
```

### JWT Tokens
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  https://api.tinybrain.com/sessions
```

## SDKs

### Go SDK
```go
import "github.com/rainmana/tinybrain-sdk-go"

client := tinybrain.NewClient("https://api.tinybrain.com", "your-api-key")

session, err := client.Sessions.Create(&tinybrain.CreateSessionRequest{
    Name: "Security Assessment",
    TaskType: "penetration_test",
})
```

### Python SDK
```python
from tinybrain import TinyBrainClient

client = TinyBrainClient("https://api.tinybrain.com", "your-api-key")

session = client.sessions.create(
    name="Security Assessment",
    task_type="penetration_test"
)
```

### JavaScript SDK
```javascript
import { TinyBrainClient } from '@tinybrain/sdk-js';

const client = new TinyBrainClient('https://api.tinybrain.com', 'your-api-key');

const session = await client.sessions.create({
    name: 'Security Assessment',
    taskType: 'penetration_test'
});
```
