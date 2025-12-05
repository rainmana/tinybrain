# Dual Endpoint Architecture Diagram

## Visual Overview: Web App + MCP Endpoint

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              END USERS                                       │
│                                                                              │
│  ┌───────────────────────┐              ┌──────────────────────────────┐   │
│  │  Human Users          │              │  AI Assistant Users          │   │
│  │  (Browser/Mobile)     │              │  (Claude, Cursor, Copilot)   │   │
│  │                       │              │                              │   │
│  │  Use Cases:           │              │  Use Cases:                  │   │
│  │  • Browse memories    │              │  • Store findings            │   │
│  │  • Create sessions    │              │  • Search memories           │   │
│  │  • Team collaboration │              │  • Query security data       │   │
│  │  • Generate API keys  │              │  • Relationship management   │   │
│  └───────────┬───────────┘              └──────────────┬───────────────┘   │
│              │                                         │                    │
└──────────────┼─────────────────────────────────────────┼────────────────────┘
               │                                         │
               │ HTTPS (Bearer JWT)                      │ HTTPS (Bearer API Key)
               │                                         │
               ▼                                         ▼
┌──────────────────────────────┐         ┌──────────────────────────────────┐
│   app.tinybrain.io           │         │   mcp.tinybrain.io               │
│   ┌──────────────────────┐   │         │   ┌──────────────────────────┐   │
│   │  Cloudflare Pages    │   │         │   │  Cloudflare Worker       │   │
│   │                      │   │         │   │  (Proxy + Rate Limit)    │   │
│   │  • Static Frontend   │   │         │   │                          │   │
│   │  • Next.js App       │   │         │   │  • Validate API key      │   │
│   │  • React Components  │   │         │   │  • Rate limiting (1K/hr) │   │
│   │  • Supabase Auth     │   │         │   │  • DDoS protection       │   │
│   └──────────────────────┘   │         │   │  • Forward to Railway    │   │
└──────────────┬───────────────┘         │   └──────────────────────────┘   │
               │                         └──────────────┬───────────────────┘
               │                                        │
               │ REST API calls                         │ MCP JSON-RPC
               │ (Authorization: Bearer <JWT>)          │ (Authorization: Bearer <API_KEY>)
               │                                        │
               └────────────────┬───────────────────────┘
                                │
                                ▼
                ┌───────────────────────────────────────────────┐
                │          Railway.app Backend                  │
                │                                               │
                │  ┌─────────────────────────────────────────┐ │
                │  │  Authentication Middleware              │ │
                │  │                                         │ │
                │  │  JWT Validator ─────► Extract user_id   │ │
                │  │  API Key Validator ─► Lookup user_id    │ │
                │  └──────────────┬──────────────────────────┘ │
                │                 │                             │
                │                 ▼                             │
                │  ┌─────────────────────────────────────────┐ │
                │  │  Dual Interface Layer                   │ │
                │  │                                         │ │
                │  │  ┌──────────────┐  ┌────────────────┐  │ │
                │  │  │  REST API    │  │  MCP Adapter   │  │ │
                │  │  │  /api/*      │  │  /v1/mcp       │  │ │
                │  │  │              │  │                │  │ │
                │  │  │ • Sessions   │  │ • store_memory │  │ │
                │  │  │ • Memories   │  │ • search_mem   │  │ │
                │  │  │ • Users      │  │ • create_sess  │  │ │
                │  │  │ • Teams      │  │ • 40+ tools    │  │ │
                │  │  └──────────────┘  └────────────────┘  │ │
                │  │                                         │ │
                │  │  Both interfaces use same:              │ │
                │  │  • Business logic                       │ │
                │  │  • Database layer                       │ │
                │  │  • User context (from auth)             │ │
                │  └─────────────────────────────────────────┘ │
                │                                               │
                │  ┌─────────────────────────────────────────┐ │
                │  │  Business Logic Layer                   │ │
                │  │                                         │ │
                │  │  • Memory management                    │ │
                │  │  • Session handling                     │ │
                │  │  • Search & indexing                    │ │
                │  │  • Relationship mapping                 │ │
                │  │  • MITRE ATT&CK integration             │ │
                │  └─────────────────────────────────────────┘ │
                └───────────────────┬───────────────────────────┘
                                    │
                                    │ PostgreSQL Protocol (TLS)
                                    │ SET LOCAL app.user_id = '<user_id>'
                                    │
                                    ▼
                ┌───────────────────────────────────────────────┐
                │              Supabase                         │
                │                                               │
                │  ┌─────────────────────────────────────────┐ │
                │  │  PostgreSQL + Row Level Security        │ │
                │  │                                         │ │
                │  │  All queries automatically filtered by: │ │
                │  │  • user_id (personal data)              │ │
                │  │  • team_id (shared data)                │ │
                │  └─────────────────────────────────────────┘ │
                │                                               │
                │  ┌─────────────────────────────────────────┐ │
                │  │  Tables:                                │ │
                │  │  • users                                │ │
                │  │  • api_keys (for MCP auth)              │ │
                │  │  • teams                                │ │
                │  │  • team_members                         │ │
                │  │  • sessions                             │ │
                │  │  • memories                             │ │
                │  │  • relationships                        │ │
                │  │  • context_snapshots                    │ │
                │  │  • task_progress                        │ │
                │  │  • notifications                        │ │
                │  └─────────────────────────────────────────┘ │
                └───────────────────────────────────────────────┘
```

## Authentication Flow Comparison

### Web App Flow (JWT)

```
1. User visits app.tinybrain.io
2. User signs in (email/password or OAuth)
3. Supabase Auth returns JWT token
4. Frontend stores token in localStorage
5. Every API call includes: Authorization: Bearer <JWT>
6. Railway validates JWT → extracts user_id
7. Database queries filtered by user_id (RLS)
```

### MCP Client Flow (API Key)

```
1. User generates API key in web app
2. User configures MCP client with API key
3. AI assistant connects to mcp.tinybrain.io
4. Every MCP call includes: Authorization: Bearer <API_KEY>
5. Cloudflare Worker validates rate limit
6. Railway validates API key → looks up user_id
7. Database queries filtered by user_id (RLS)
```

## Data Flow Examples

### Example 1: Web User Creates Memory

```
Browser
  │
  │ POST /api/memories
  │ Authorization: Bearer eyJhbGc...  (JWT)
  │ Body: { title: "Finding", content: "..." }
  │
  ▼
Railway Backend
  │
  │ 1. Validate JWT → user_id = "alice"
  │ 2. Insert memory with user_id = "alice"
  │
  ▼
Supabase
  │
  │ INSERT INTO memories (user_id, title, content)
  │ VALUES ('alice', 'Finding', '...')
  │
  │ RLS Policy: ✓ Allowed (user owns this memory)
  │
  └─► Success
```

### Example 2: MCP Client Searches Memories

```
Claude Desktop (MCP Client)
  │
  │ POST /v1/mcp
  │ Authorization: Bearer <API_KEY>  (User's API key)
  │ Body: { method: "search_memories", params: { query: "SQL" } }
  │
  ▼
Cloudflare Worker
  │
  │ Rate limit check: ✓ 450/1000 requests used
  │
  ▼
Railway Backend
  │
  │ 1. Validate API key
  │ 2. Lookup: API key belongs to user_id = "alice"
  │ 3. Process: search_memories(query: "SQL", user: "alice")
  │
  ▼
Supabase
  │
  │ SELECT * FROM memories 
  │ WHERE to_tsvector(content) @@ plainto_tsquery('SQL')
  │   AND user_id = 'alice'  <-- RLS adds this automatically
  │
  │ RLS Policy: ✓ Only returns alice's memories
  │
  └─► Returns: [{ memory_id: "...", title: "SQL Injection Finding" }]
```

### Example 3: Team Shared Memory

```
Alice (Web) creates memory with team_id = "security-team"
  │
  ▼
Supabase stores: 
  user_id = "alice"
  team_id = "security-team"
  │
  ▼
Bob (MCP client) searches
  │
  │ API key belongs to user_id = "bob"
  │ Bob is member of "security-team"
  │
  ▼
RLS Policy evaluates:
  │
  │ WHERE user_id = 'bob'           <-- Bob's own memories
  │    OR (team_id = 'security-team' 
  │        AND 'bob' IN team_members)  <-- Bob's team memories
  │
  └─► Bob sees Alice's team memory ✓
```

## Key Design Decisions

### 1. Why Two Domains?

**Separation of Concerns**:
- `app.tinybrain.io` - Human-focused web interface
- `mcp.tinybrain.io` - Machine-focused protocol endpoint

**Benefits**:
- Independent rate limiting
- Separate analytics/monitoring
- Different caching strategies
- Clearer user mental model

### 2. Why Single Backend Service?

**Code Reuse**:
- Same business logic
- Same database layer
- Same authentication layer
- Easier to maintain

**Alternative** (not chosen):
- Separate services would duplicate logic
- Harder to keep in sync
- More operational complexity

### 3. Why API Keys for MCP?

**JWT Limitations**:
- JWTs typically short-lived (15 min)
- MCP clients run long sessions
- Refresh token flow complex for CLI tools

**API Key Benefits**:
- Long-lived (no expiration or user-controlled)
- Simpler for MCP client config
- Easy to revoke
- Can be scoped (future: read-only keys)

### 4. Why Cloudflare Proxy for MCP?

**Protection**:
- DDoS protection (Cloudflare network)
- Rate limiting at edge (before hitting Railway)
- SSL/TLS termination
- Geographic distribution

**Cost**:
- Railway charges for bandwidth
- Cloudflare absorbs attack traffic
- Edge rate limiting prevents backend overload

## Scaling Considerations

### Current (Phase 1)
- Single Railway instance
- Single Supabase database
- Handles ~100 concurrent users

### Future (Scaling)

**Railway**:
- Horizontal scaling (multiple instances)
- Load balancer (built-in)
- Auto-scaling based on CPU/memory

**Supabase**:
- Connection pooling (PgBouncer)
- Read replicas (for queries)
- Vertical scaling (larger instance)

**Cloudflare**:
- Already global (no changes needed)
- KV replication (multi-region)

## Cost Breakdown

| Component | Free Tier | Recommended | Cost |
|-----------|-----------|-------------|------|
| Supabase | 500MB DB | Pro ($25/mo) | $25 |
| Railway | $5 credit | Usage-based | $30-50 |
| Cloudflare Pages | Unlimited | Free | $0 |
| Cloudflare Workers | 100K req/day | Paid ($5/mo) | $5 |
| **Total** | Limited | Production | **$60-80/mo** |

_Note: Prices as of December 2024. Check current pricing at service provider websites._

## Summary

The dual-endpoint architecture provides:

✅ **Single Backend**: One Railway service, two interfaces  
✅ **Shared Database**: Supabase with RLS ensures isolation  
✅ **Dual Authentication**: JWT (web) + API Keys (MCP)  
✅ **Separate Domains**: Clear separation of concerns  
✅ **Cloudflare Protection**: DDoS, rate limiting, edge caching  
✅ **Multi-User**: Full isolation via RLS policies  
✅ **Team Collaboration**: Shared memories via team_id  
✅ **Scalable**: Horizontal scaling for both interfaces  

Both web users and MCP users access the same data, with proper authentication, authorization, and rate limiting.
