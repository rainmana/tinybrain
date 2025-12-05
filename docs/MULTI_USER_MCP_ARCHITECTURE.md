# Multi-User MCP Server Architecture

## Overview

This document addresses the architecture for supporting both a **multi-user web application** and a **multi-user MCP server endpoint** backed by the same Supabase database.

## Problem Statement

How do we:
1. Enable multiple users to use the web-based MCP server (multi-user support)
2. Provide a single MCP endpoint for programmatic access
3. Have separate hostnames for the web app vs. MCP endpoint
4. Share the same Supabase backend between both interfaces

## Proposed Architecture

### Two-Domain Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      End Users                                   │
│                                                                   │
│  ┌──────────────────────┐    ┌───────────────────────────┐     │
│  │   Web Users          │    │   MCP Clients             │     │
│  │   (Browser/Mobile)   │    │   (Claude, Cursor, etc.)  │     │
│  └──────────┬───────────┘    └───────────┬───────────────┘     │
└─────────────┼──────────────────────────────┼───────────────────┘
              │ HTTPS                        │ HTTPS/WSS
              │                              │
              ▼                              ▼
┌─────────────────────────┐    ┌──────────────────────────────┐
│   app.tinybrain.io      │    │   mcp.tinybrain.io           │
│   (Cloudflare Pages)    │    │   (Railway + CF Proxy)       │
│                         │    │                              │
│  • Web Dashboard        │    │  • MCP Protocol Endpoint     │
│  • User Authentication  │    │  • Per-User Authentication   │
│  • Session Management   │    │  • User-Scoped Data Access   │
│  • Memory Browser       │    │  • Rate Limiting per User    │
│  • Team Collaboration   │    │  • Connection Pooling        │
└─────────────┬───────────┘    └──────────────┬───────────────┘
              │                                │
              │ HTTPS                          │ HTTPS
              │                                │
              └────────────────┬───────────────┘
                               │
                               ▼
                ┌──────────────────────────────────┐
                │     Railway.app Backend          │
                │                                  │
                │  ┌────────────────────────────┐ │
                │  │  REST API                  │ │
                │  │  (for web app)             │ │
                │  └────────────────────────────┘ │
                │                                  │
                │  ┌────────────────────────────┐ │
                │  │  MCP Protocol Adapter      │ │
                │  │  (for MCP clients)         │ │
                │  │  • Multi-user support      │ │
                │  │  • API key authentication  │ │
                │  │  • Per-user data isolation │ │
                │  └────────────────────────────┘ │
                │                                  │
                │  ┌────────────────────────────┐ │
                │  │  Authentication Service    │ │
                │  │  • JWT validation          │ │
                │  │  • API key validation      │ │
                │  │  • User context injection  │ │
                │  └────────────────────────────┘ │
                └──────────────┬───────────────────┘
                               │ PostgreSQL
                               ▼
                ┌──────────────────────────────────┐
                │         Supabase                 │
                │                                  │
                │  • PostgreSQL Database           │
                │  • Row Level Security (RLS)      │
                │  • User Authentication           │
                │  • API Keys Management           │
                │  • Multi-tenant data isolation   │
                └──────────────────────────────────┘
```

## Key Components

### 1. Web Application Domain: `app.tinybrain.io`

**Purpose**: Human-facing web interface

**Features**:
- Web dashboard for managing memories and sessions
- User authentication via Supabase Auth (email/password, OAuth)
- Team collaboration features
- Real-time updates
- Visual data exploration

**Authentication Flow**:
1. User signs in via web form
2. Supabase Auth returns JWT token
3. Token stored in browser (cookie/localStorage)
4. Token sent with every API request
5. Backend validates JWT and applies RLS

### 2. MCP Server Domain: `mcp.tinybrain.io`

**Purpose**: Programmatic MCP protocol access for AI assistants

**Features**:
- Standard MCP JSON-RPC endpoint
- Per-user API key authentication
- User-scoped data access
- Rate limiting per user
- WebSocket support for real-time

**Authentication Flow**:
1. User generates API key in web app
2. API key stored in user's MCP client config
3. MCP client sends API key in Authorization header
4. Backend validates API key → maps to user ID
5. User context injected → RLS applied

## Multi-User MCP Server Implementation

### Challenge: Traditional MCP is Single-User

Traditional MCP servers (like current TinyBrain) run locally and are inherently single-user. To make it multi-user:

### Solution: API Key-Based User Identification

```go
// MCP endpoint with multi-user support
POST https://mcp.tinybrain.io/v1/mcp

Headers:
  Authorization: Bearer api_key_user123xyz
  Content-Type: application/json

Body:
  {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "store_memory",
    "params": {
      "title": "Security Finding",
      "content": "SQL injection in login form"
    }
  }

Response:
  {
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
      "memory_id": "uuid-123",
      "user_id": "user123",  // Derived from API key
      "title": "Security Finding",
      "created_at": "2024-12-05T00:00:00Z"
    }
  }
```

### Backend Flow

```go
// Pseudo-code for MCP endpoint handler
func HandleMCPRequest(w http.ResponseWriter, r *http.Request) {
    // 1. Extract API key from Authorization header
    apiKey := extractBearerToken(r.Header.Get("Authorization"))
    
    // 2. Validate API key and get user ID
    userID, err := validateAPIKey(apiKey)
    if err != nil {
        return JSONError(w, "Invalid API key", 401)
    }
    
    // 3. Inject user context (used by RLS)
    ctx := context.WithValue(r.Context(), "user_id", userID)
    
    // 4. Parse MCP JSON-RPC request
    var mcpRequest MCPRequest
    json.NewDecoder(r.Body).Decode(&mcpRequest)
    
    // 5. Process MCP method with user context
    result, err := processMCPMethod(ctx, mcpRequest.Method, mcpRequest.Params)
    if err != nil {
        return MCPError(w, err)
    }
    
    // 6. Return MCP response
    return MCPResponse(w, mcpRequest.ID, result)
}

// All database operations respect user_id from context
func StoreMemory(ctx context.Context, memory Memory) error {
    userID := ctx.Value("user_id").(string)
    memory.UserID = userID // Enforce user ownership
    
    // Supabase RLS automatically filters by user_id
    return db.Insert("memories", memory)
}
```

## API Key Management

### Generating API Keys (Web App)

Users generate API keys in the web dashboard:

```
Settings → API Keys → Generate New Key
```

**API Key Format**: `tbrain_live_xxxxxxxxxxxxxxxx` (prefix identifies environment)

_Note: Example format only. Actual keys are 32+ random bytes, base64-encoded (~43 characters)._

**Storage in Supabase**:
```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash TEXT NOT NULL UNIQUE, -- Hashed with bcrypt
    name TEXT NOT NULL, -- User-provided description
    last_used TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    revoked BOOLEAN DEFAULT FALSE
);
```

**Security**:
- API keys are hashed (bcrypt) before storage
- Only show full key once at creation time
- Support key expiration and revocation
- Log last usage timestamp
- Rate limit per key

### MCP Client Configuration

Users configure their MCP clients (Claude Desktop, Cursor, etc.) with their API key:

**Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "https://mcp.tinybrain.io/v1/mcp",
        "-H", "Authorization: Bearer YOUR_API_KEY_HERE",
        "-H", "Content-Type: application/json",
        "-d", "@-"
      ]
    }
  }
}
```

**Alternative: HTTP Transport** (when MCP supports it):
```json
{
  "mcpServers": {
    "tinybrain": {
      "url": "https://mcp.tinybrain.io/v1/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY_HERE"
      }
    }
  }
}
```

## Cloudflare Proxy for MCP Endpoint

Use Cloudflare Workers to proxy MCP endpoint with additional features:

### Why Proxy?

1. **DDoS Protection**: Cloudflare's network protects Railway backend
2. **Rate Limiting**: Per-user rate limits at the edge
3. **Caching**: Cache security knowledge hub responses
4. **Analytics**: Track MCP usage patterns
5. **Custom Domain**: `mcp.tinybrain.io` → Railway backend

### Cloudflare Worker for MCP

```typescript
// cloudflare/workers/mcp-proxy.ts

interface Env {
  BACKEND_URL: string;
  RATE_LIMIT_KV: KVNamespace;
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    // 1. Extract API key
    const apiKey = request.headers.get('Authorization')?.replace('Bearer ', '');
    if (!apiKey) {
      return jsonError('Missing API key', 401);
    }
    
    // 2. Rate limiting (per API key)
    const rateLimitKey = `ratelimit:${apiKey}`;
    const count = await env.RATE_LIMIT_KV.get(rateLimitKey);
    
    if (count && parseInt(count) > 1000) { // 1000 requests/hour
      return jsonError('Rate limit exceeded', 429);
    }
    
    // 3. Increment counter
    const newCount = count ? parseInt(count) + 1 : 1;
    await env.RATE_LIMIT_KV.put(rateLimitKey, String(newCount), {
      expirationTtl: 3600, // 1 hour
    });
    
    // 4. Forward to Railway backend
    const backendURL = new URL('/v1/mcp', env.BACKEND_URL);
    const backendRequest = new Request(backendURL.toString(), {
      method: request.method,
      headers: request.headers,
      body: request.body,
    });
    
    const response = await fetch(backendRequest);
    
    // 5. Add custom headers
    const newResponse = new Response(response.body, response);
    newResponse.headers.set('X-Proxy', 'Cloudflare Workers');
    newResponse.headers.set('X-Rate-Limit-Remaining', String(1000 - newCount));
    
    return newResponse;
  },
};
```

## Data Isolation with Row Level Security

### How RLS Ensures Multi-User Isolation

Every database operation automatically filters by the authenticated user:

```sql
-- Example RLS policy on memories table
CREATE POLICY "Users can only access their own memories"
  ON memories
  FOR ALL
  USING (user_id = auth.uid());

-- When user "alice" queries:
SELECT * FROM memories WHERE session_id = 'session-123';

-- PostgreSQL automatically rewrites to:
SELECT * FROM memories 
WHERE session_id = 'session-123' 
  AND user_id = 'alice-user-id';  -- Added by RLS
```

### Setting User Context

**For Web API** (JWT):
```go
// Extract user ID from JWT token
claims := parseJWT(token)
userID := claims["sub"]

// Set Supabase context
supabase.SetAuthContext(userID)
```

**For MCP API** (API Key):
```go
// Look up user ID from API key
userID := lookupUserIDFromAPIKey(apiKey)

// Set Supabase context
supabase.SetAuthContext(userID)
```

## Team/Shared Access

For team collaboration (where users share memories):

```sql
-- RLS policy for team access
CREATE POLICY "Team members can access team memories"
  ON memories
  FOR ALL
  USING (
    user_id = auth.uid() 
    OR (
      team_id IS NOT NULL 
      AND EXISTS (
        SELECT 1 FROM team_members
        WHERE team_id = memories.team_id
          AND user_id = auth.uid()
      )
    )
  );
```

**Impact**:
- Users see their own memories + team memories
- Works for both web app and MCP endpoint
- RLS handles all filtering automatically

## Deployment Architecture

### DNS Configuration

```
app.tinybrain.io    → Cloudflare Pages (A/CNAME)
mcp.tinybrain.io    → Cloudflare Worker → Railway (A/CNAME)
```

### Railway Configuration

**Single Railway Service** with multiple routes:

```yaml
# Railway service
tinybrain-backend:
  routes:
    - /api/*        # REST API for web app
    - /v1/mcp       # MCP protocol endpoint
  
  environment:
    - SUPABASE_URL
    - SUPABASE_SERVICE_KEY
    - JWT_SECRET
```

**Or Two Railway Services** (better isolation):

```yaml
# Web API service
tinybrain-web-api:
  routes:
    - /api/*
  
# MCP API service  
tinybrain-mcp-api:
  routes:
    - /v1/mcp
```

## User Journey

### Web App User

1. Visit `app.tinybrain.io`
2. Sign up with email/password or OAuth
3. Use web dashboard to create memories
4. Invite team members
5. Generate API key for MCP access

### MCP User

1. Get API key from web app (or via CLI tool)
2. Configure MCP client with:
   - Endpoint: `https://mcp.tinybrain.io/v1/mcp`
   - API Key: `tbrain_live_xxxxxxxxxxxxxxxx` (example format)
3. Use AI assistant (Claude, Cursor) normally
4. Memories stored with user isolation
5. Can view/manage memories in web app

### Both Interfaces

- Same Supabase backend
- Same user account
- Same memories and sessions
- RLS ensures data isolation
- Web app shows MCP activity
- MCP endpoint has web app's data

## Implementation Phases

### Phase 1: Single-User MCP (Current)
- ✅ Local MCP server
- ✅ No authentication
- ✅ Single user

### Phase 2: Multi-User Backend (In Progress)
- 🔄 Supabase integration
- 🔄 REST API with JWT auth
- 🔄 RLS policies

### Phase 3: Multi-User MCP Endpoint
- ⏳ API key generation in web app
- ⏳ API key validation middleware
- ⏳ MCP endpoint with user context
- ⏳ Rate limiting per API key

### Phase 4: Cloudflare Proxy
- ⏳ Worker for MCP endpoint
- ⏳ Edge rate limiting
- ⏳ Custom domain setup

### Phase 5: Web Dashboard
- ⏳ API key management UI
- ⏳ MCP usage analytics
- ⏳ Connection status display

## Security Considerations

### API Key Security

1. **Generation**: 
   - Cryptographically random (32+ raw bytes)
   - Base64-encoded for transmission (~43 characters)
   - Format: `tbrain_{env}_{base64_encoded_random}`
   - Example length: 60-70 characters total
2. **Storage**: Bcrypt hashed in database (never store plaintext)
3. **Transmission**: HTTPS only, Bearer token in Authorization header
4. **Rotation**: Support key rotation without downtime
5. **Scope**: Per-user, optionally per-team (future)
6. **Expiration**: Optional expiration dates
7. **Revocation**: Instant revocation support

### Rate Limiting

- **Per API Key**: 1000 requests/hour
- **Per IP**: 10000 requests/hour (DDoS protection)
- **Burst Protection**: Max 10 req/second per key
- **Edge Enforcement**: Cloudflare Workers

### Monitoring

- **Failed Auth Attempts**: Alert on repeated failures
- **Unusual Usage**: Alert on spike in API calls
- **Data Access**: Audit log for sensitive operations
- **Key Usage**: Track last used timestamp

## Cost Impact

### Additional Costs

- **Cloudflare Workers**: $5/month (includes KV for rate limiting)
- **API Key Storage**: Minimal (< 1KB per user)
- **Additional Traffic**: Included in Railway/Supabase plans

**Total Additional Cost**: ~$5/month

## FAQ

### Q: Can multiple users use the same API key?

**A**: Not recommended. Each user should have their own API key for proper data isolation and rate limiting. However, team API keys could be supported in the future.

### Q: How do I revoke an API key?

**A**: In the web app: Settings → API Keys → Revoke. Takes effect immediately (keys are validated on every request).

### Q: Can I use the MCP endpoint without the web app?

**A**: Yes, but you need to create an account and generate an API key first. This can be done via CLI tool or web app.

### Q: What happens if my API key is leaked?

**A**: Revoke it immediately in the web app. Generate a new one. The leaked key stops working instantly.

### Q: Does the MCP endpoint support all MCP features?

**A**: Yes, all 40+ MCP tools are available. The only difference is authentication (API key instead of local stdio).

### Q: Can I have multiple API keys?

**A**: Yes, you can generate multiple keys (e.g., one per device, one per AI assistant). Each key is tracked separately.

## Summary

The multi-user MCP architecture provides:

1. **Two Domains**:
   - `app.tinybrain.io` - Web dashboard
   - `mcp.tinybrain.io` - MCP protocol endpoint

2. **Single Backend**: Railway service with dual interfaces

3. **Shared Database**: Supabase with RLS for isolation

4. **Authentication**:
   - Web: JWT tokens (Supabase Auth)
   - MCP: API keys (generated in web app)

5. **User Isolation**: RLS policies ensure data separation

6. **Cloudflare Proxy**: Rate limiting, DDoS protection, custom domain

This architecture enables both human (web) and programmatic (MCP) access to the same multi-user TinyBrain instance, with proper authentication, authorization, and data isolation.
