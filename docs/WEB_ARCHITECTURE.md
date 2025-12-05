# TinyBrain Web Architecture

## Overview

This document describes the architecture for the web-based version of TinyBrain, which migrates from a local Go MCP server with PocketBase to a distributed cloud architecture using Supabase, Railway.app, and Cloudflare Pages.

## Current Architecture (v1.2.1)

```
┌─────────────────────────────────────────────────┐
│                  LLM Client                     │
│              (Claude, GPT, etc.)                │
└────────────────┬────────────────────────────────┘
                 │ MCP Protocol
                 │ (JSON-RPC over stdio/http)
                 ▼
┌─────────────────────────────────────────────────┐
│           TinyBrain MCP Server (Go)             │
│  ┌─────────────────────────────────────────┐   │
│  │         PocketBase Backend              │   │
│  │  ┌───────────────────────────────────┐  │   │
│  │  │      SQLite Database              │  │   │
│  │  │   (local file: pb_data/data.db)   │  │   │
│  │  └───────────────────────────────────┘  │   │
│  │  ┌───────────────────────────────────┐  │   │
│  │  │      Admin Dashboard              │  │   │
│  │  │   (http://127.0.0.1:8090/_/)      │  │   │
│  │  └───────────────────────────────────┘  │   │
│  └─────────────────────────────────────────┘   │
│                                                  │
│  Features:                                       │
│  • 40+ MCP Tools                                 │
│  • Memory Management                             │
│  • Intelligence Gathering                        │
│  • MITRE ATT&CK Integration                      │
│  • Security Knowledge Hub                        │
└─────────────────────────────────────────────────┘
```

**Limitations:**
- Single-user local deployment
- No collaborative features
- Limited scalability
- Requires local installation
- No web interface for non-technical users

## Target Architecture (Web-Based)

```
┌───────────────────────────────────────────────────────────────────┐
│                         End Users                                  │
│                  (Browser, Mobile, API Clients)                   │
└────────────────┬──────────────────────────────────────────────────┘
                 │ HTTPS
                 ▼
┌───────────────────────────────────────────────────────────────────┐
│                    Cloudflare Pages                                │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Static Frontend Assets                       │    │
│  │  • React/Next.js Web Dashboard                            │    │
│  │  • Memory Browser & Search UI                             │    │
│  │  • Session Management Interface                           │    │
│  │  • Real-time Data Visualization                           │    │
│  │  • Security Intelligence Dashboard                        │    │
│  └──────────────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │            Cloudflare Workers                             │    │
│  │  • API Request Routing                                    │    │
│  │  • Edge Caching & CDN                                     │    │
│  │  • Rate Limiting & Security                               │    │
│  │  • WebSocket Proxy for Real-time                          │    │
│  └──────────────────────────────────────────────────────────┘    │
└────────────────┬──────────────────────────────────────────────────┘
                 │ HTTPS/WSS
                 ▼
┌───────────────────────────────────────────────────────────────────┐
│                       Railway.app                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              TinyBrain API Server (Go)                    │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │           REST/GraphQL API Layer                │      │    │
│  │  │  • Authentication & Authorization               │      │    │
│  │  │  • Session Management API                       │      │    │
│  │  │  • Memory Operations API                        │      │    │
│  │  │  • Search & Query API                           │      │    │
│  │  │  • Real-time Subscription API (WebSocket)       │      │    │
│  │  │  • Intelligence Gathering API                   │      │    │
│  │  │  • Security Knowledge Hub API                   │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │           MCP Protocol Adapter                  │      │    │
│  │  │  • Backward compatibility with MCP clients      │      │    │
│  │  │  • JSON-RPC over HTTP endpoint                  │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │           Business Logic Layer                  │      │    │
│  │  │  • Memory categorization                        │      │    │
│  │  │  • Semantic search                              │      │    │
│  │  │  • Relationship management                      │      │    │
│  │  │  • Context analysis                             │      │    │
│  │  │  • MITRE ATT&CK mapping                         │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                    │
│  Infrastructure:                                                   │
│  • Auto-scaling containers                                         │
│  • Health monitoring                                               │
│  • Logging & metrics                                               │
│  • Zero-downtime deployment                                        │
└────────────────┬──────────────────────────────────────────────────┘
                 │ PostgreSQL Protocol (TLS)
                 ▼
┌───────────────────────────────────────────────────────────────────┐
│                         Supabase                                   │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              PostgreSQL Database                          │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │           Core Collections                      │      │    │
│  │  │  • memories (memory entries)                    │      │    │
│  │  │  • sessions (assessment sessions)               │      │    │
│  │  │  • relationships (memory links)                 │      │    │
│  │  │  • context_snapshots (context states)           │      │    │
│  │  │  • task_progress (progress tracking)            │      │    │
│  │  │  • notifications (alerts & events)              │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │       Security Knowledge Collections            │      │    │
│  │  │  • nvd_cves (vulnerability database)            │      │    │
│  │  │  • mitre_attack (tactics & techniques)          │      │    │
│  │  │  • owasp_tests (security testing procedures)    │      │    │
│  │  │  • cwe_patterns (weakness patterns)             │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  │  ┌────────────────────────────────────────────────┐      │    │
│  │  │           User & Access Management              │      │    │
│  │  │  • users (user accounts)                        │      │    │
│  │  │  • teams (team/org structures)                  │      │    │
│  │  │  • permissions (access control)                 │      │    │
│  │  └────────────────────────────────────────────────┘      │    │
│  └──────────────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Row Level Security (RLS)                     │    │
│  │  • User-level data isolation                              │    │
│  │  • Team-based access control                              │    │
│  │  • Read/write permission policies                         │    │
│  └──────────────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Real-time Subscriptions                      │    │
│  │  • Memory change notifications                             │    │
│  │  • Session activity streams                               │    │
│  │  • Alert broadcasts                                        │    │
│  └──────────────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Authentication                                │    │
│  │  • Email/Password                                          │    │
│  │  • OAuth (Google, GitHub)                                 │    │
│  │  • API Keys for programmatic access                       │    │
│  │  • JWT token management                                   │    │
│  └──────────────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Storage                                       │    │
│  │  • File attachments                                        │    │
│  │  • Export archives                                         │    │
│  │  • Backup snapshots                                        │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                    │
│  Features:                                                         │
│  • Automatic backups                                               │
│  • Point-in-time recovery                                          │
│  • Connection pooling                                              │
│  • Built-in monitoring                                             │
└───────────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

### 1. Cloudflare Pages (Frontend & Edge)

**Purpose:** Serve static frontend assets and provide edge computing capabilities

**Responsibilities:**
- Host the React/Next.js web application
- Provide global CDN for fast content delivery
- Handle SSL/TLS termination
- Execute Cloudflare Workers for edge logic
- Route API requests to Railway backend
- Implement edge caching strategies
- Rate limiting and DDoS protection

**Technologies:**
- Next.js/React for UI framework
- Cloudflare Pages for static hosting
- Cloudflare Workers for serverless edge functions
- TypeScript for type-safe frontend code

### 2. Railway.app (Backend API)

**Purpose:** Host the TinyBrain API server with auto-scaling and monitoring

**Responsibilities:**
- REST/GraphQL API endpoints
- MCP protocol adapter for backward compatibility
- Business logic implementation
- Authentication and authorization
- WebSocket server for real-time features
- Background job processing
- Integration with external APIs (MITRE, NVD, OWASP)
- Logging and metrics collection

**Technologies:**
- Go 1.24+ for high-performance backend
- Gorilla/Mux or Gin for HTTP routing
- WebSocket library for real-time communication
- Supabase Go client
- Docker containers for deployment

### 3. Supabase (Database & Backend Services)

**Purpose:** Managed PostgreSQL database with authentication and real-time capabilities

**Responsibilities:**
- Primary data storage
- User authentication and management
- Row-level security policies
- Real-time subscriptions
- File storage for attachments
- Automated backups
- Database functions and triggers
- Connection pooling

**Technologies:**
- PostgreSQL 15+ with extensions (pg_vector for embeddings)
- Supabase Auth for authentication
- Supabase Storage for files
- Supabase Realtime for subscriptions

## Data Flow

### 1. User Authentication Flow

```
User → Cloudflare Pages → Supabase Auth → JWT Token
                              ↓
                        Store in browser
                              ↓
                    Include in API requests
```

### 2. Memory Storage Flow

```
User → Frontend → Cloudflare Worker → Railway API
                                          ↓
                                   Validate JWT
                                          ↓
                                   Business Logic
                                          ↓
                                   Supabase Insert
                                          ↓
                                   Real-time Event
                                          ↓
                                   Notify Subscribers
```

### 3. Search Query Flow

```
User → Frontend → Railway API → Supabase Query
                                      ↓
                                Full-text search
                                      ↓
                                Semantic search
                                      ↓
                                Return results
```

### 4. Real-time Updates Flow

```
Supabase Change → Supabase Realtime → Railway WebSocket → Frontend
                                                              ↓
                                                        Update UI
```

## Security Architecture

### Authentication Layers

1. **User Authentication** (Supabase Auth)
   - Email/password with verification
   - OAuth providers (Google, GitHub)
   - Magic links
   - MFA support

2. **API Authentication** (JWT)
   - Short-lived access tokens
   - Refresh token rotation
   - Token revocation

3. **API Keys** (for programmatic access)
   - Scoped permissions
   - Rate limiting per key
   - Usage tracking

### Authorization Strategy

1. **Row Level Security (RLS)** in Supabase
   ```sql
   -- Users can only access their own memories
   CREATE POLICY "Users can view own memories"
     ON memories FOR SELECT
     USING (auth.uid() = user_id);
   
   -- Team members can view shared memories
   CREATE POLICY "Team members can view shared memories"
     ON memories FOR SELECT
     USING (
       EXISTS (
         SELECT 1 FROM team_members
         WHERE team_id = memories.team_id
         AND user_id = auth.uid()
       )
     );
   ```

2. **API-Level Authorization**
   - Role-based access control (RBAC)
   - Resource ownership validation
   - Team/organization boundaries

### Data Protection

1. **Encryption**
   - TLS 1.3 for all connections
   - At-rest encryption in Supabase
   - Optional field-level encryption for sensitive data

2. **Secrets Management**
   - Environment variables for credentials
   - Railway secrets for API keys
   - No hardcoded secrets in code

## Deployment Strategy

### Environment Structure

1. **Development**
   - Local Supabase instance (optional)
   - Railway preview deployments
   - Cloudflare Pages preview

2. **Staging**
   - Shared Supabase staging project
   - Railway staging environment
   - Cloudflare Pages branch deployments

3. **Production**
   - Production Supabase project
   - Railway production environment
   - Cloudflare Pages production deployment

### CI/CD Pipeline

```
Git Push → GitHub Actions → Run Tests
                                ↓
                          Build Backend
                                ↓
                     Deploy to Railway (auto)
                                ↓
                          Build Frontend
                                ↓
                 Deploy to Cloudflare Pages (auto)
                                ↓
                        Run E2E Tests
                                ↓
                          Monitor Health
```

## Scalability Considerations

### Database Optimization

1. **Indexing Strategy**
   - B-tree indexes for lookups
   - GIN indexes for JSON/array fields
   - Full-text search indexes
   - Vector indexes for semantic search

2. **Query Optimization**
   - Prepared statements
   - Connection pooling
   - Query result caching
   - Pagination for large datasets

3. **Partitioning**
   - Time-based partitioning for large tables
   - Archive old data to cold storage

### API Server Scaling

1. **Horizontal Scaling**
   - Railway auto-scaling based on CPU/memory
   - Stateless API design
   - Load balancing across instances

2. **Caching**
   - Redis for session caching
   - CDN for static responses
   - Application-level caching

3. **Rate Limiting**
   - Per-user rate limits
   - API key rate limits
   - DDoS protection at edge (Cloudflare)

## Migration Path

### Phase 1: Database Migration
- Export current SQLite data
- Transform to PostgreSQL schema
- Import into Supabase
- Validate data integrity

### Phase 2: Backend Adaptation
- Add REST API endpoints
- Implement authentication
- Maintain MCP compatibility
- Deploy to Railway

### Phase 3: Frontend Development
- Build web dashboard
- Implement authentication UI
- Connect to API
- Deploy to Cloudflare Pages

### Phase 4: Cutover
- Run both systems in parallel
- Migrate users gradually
- Monitor and validate
- Sunset local version

## Monitoring & Observability

### Metrics to Track

1. **Application Metrics**
   - API request rate
   - Response times
   - Error rates
   - Active users

2. **Database Metrics**
   - Connection pool usage
   - Query performance
   - Database size
   - Replication lag

3. **Infrastructure Metrics**
   - CPU/memory usage
   - Network throughput
   - Disk I/O
   - Container health

### Logging Strategy

1. **Structured Logging**
   - JSON format
   - Correlation IDs
   - User context
   - Request/response tracing

2. **Log Aggregation**
   - Centralized logging (Railway logs)
   - Log retention policies
   - Alert on errors
   - Performance analysis

## Cost Considerations

### Supabase
- Free tier: 500MB database, 1GB storage, 2GB bandwidth
- Pro tier: $25/month for production
- Additional costs for high storage/bandwidth

### Railway
- Free tier: $5 credit/month
- Usage-based pricing: ~$20-50/month for small-medium apps
- Scales with traffic

### Cloudflare Pages
- Free tier: Unlimited requests
- Very cost-effective for static hosting

**Estimated Total:** $25-75/month for small to medium deployment

## Multi-User MCP Server

For detailed information on supporting multiple users with the MCP protocol endpoint, see:

**[Multi-User MCP Architecture Guide](./MULTI_USER_MCP_ARCHITECTURE.md)**

This covers:
- Two-domain architecture (`app.tinybrain.io` vs `mcp.tinybrain.io`)
- API key-based authentication for MCP clients
- Per-user data isolation with RLS
- Cloudflare proxy for the MCP endpoint
- Rate limiting and security considerations

## Future Enhancements

1. **AI Integration**
   - OpenAI/Anthropic embeddings
   - Semantic similarity search
   - Intelligent categorization

2. **Collaboration Features**
   - Team workspaces
   - Shared sessions
   - Real-time collaboration

3. **Mobile Apps**
   - React Native mobile app
   - Push notifications
   - Offline support

4. **Advanced Analytics**
   - Memory usage patterns
   - Search analytics
   - Security insights

5. **Integration Marketplace**
   - Third-party integrations
   - Webhook support
   - API marketplace
