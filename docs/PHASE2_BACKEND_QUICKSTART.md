# Phase 2: Backend API Development Quick Start

This guide helps you get started with Phase 2 of the web implementation - adapting the backend to work with Supabase and providing REST/GraphQL APIs.

## Prerequisites

Before starting Phase 2, ensure you have:

- ✅ Completed Phase 1 (documentation and configurations)
- ✅ Supabase project created and configured
- ✅ Railway account ready
- ✅ Go 1.24+ installed
- ✅ Local development environment set up

## Overview

Phase 2 transforms the current PocketBase-based backend into a cloud-native API server that:

1. Connects to Supabase (PostgreSQL)
2. Provides REST/GraphQL endpoints
3. Maintains MCP protocol compatibility
4. Implements authentication and authorization
5. Supports real-time features via WebSocket

## Directory Structure

```
tinybrain/
├── cmd/tinybrain/           # Main application entry point
│   └── main.go             # Server initialization
├── internal/
│   ├── api/                # API handlers (NEW)
│   │   ├── rest/          # REST API endpoints
│   │   │   ├── sessions.go
│   │   │   ├── memories.go
│   │   │   ├── relationships.go
│   │   │   └── security.go
│   │   ├── mcp/           # MCP protocol adapter
│   │   │   └── adapter.go
│   │   └── middleware/    # HTTP middleware
│   │       ├── auth.go
│   │       ├── cors.go
│   │       └── logging.go
│   ├── database/           # Database layer (MODIFY)
│   │   ├── supabase.go    # Supabase client
│   │   └── queries.go     # Database queries
│   ├── models/            # Data models (MODIFY)
│   │   ├── session.go
│   │   ├── memory.go
│   │   └── user.go
│   ├── repository/        # Data access layer (MODIFY)
│   │   ├── session_repo.go
│   │   └── memory_repo.go
│   └── services/          # Business logic (MODIFY)
│       ├── auth_service.go
│       ├── session_service.go
│       └── memory_service.go
└── pkg/                   # Shared packages
    └── websocket/         # WebSocket server (NEW)
```

## Step 1: Install Dependencies

Add Supabase Go client and other required packages:

```bash
go get github.com/supabase-community/supabase-go
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get github.com/golang-jwt/jwt/v5
go get github.com/rs/cors
```

Update `go.mod`:

```go
require (
    github.com/supabase-community/supabase-go v0.0.1
    github.com/gorilla/mux v1.8.1
    github.com/gorilla/websocket v1.5.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/rs/cors v1.10.1
    // ... existing dependencies
)
```

## Step 2: Create Supabase Client

Create `internal/database/supabase.go`:

```go
package database

import (
    "fmt"
    "os"
    
    supabase "github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
    Client *supabase.Client
}

func NewSupabaseClient() (*SupabaseClient, error) {
    url := os.Getenv("SUPABASE_URL")
    key := os.Getenv("SUPABASE_SERVICE_KEY") // Use service key for server-side
    
    if url == "" || key == "" {
        return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_KEY must be set")
    }
    
    client, err := supabase.NewClient(url, key, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create Supabase client: %w", err)
    }
    
    return &SupabaseClient{Client: client}, nil
}
```

## Step 3: Update Models for PostgreSQL

Update `internal/models/memory.go`:

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Memory struct {
    ID            uuid.UUID         `json:"id" db:"id"`
    UserID        uuid.UUID         `json:"user_id" db:"user_id"`
    SessionID     *uuid.UUID        `json:"session_id,omitempty" db:"session_id"`
    TeamID        *uuid.UUID        `json:"team_id,omitempty" db:"team_id"`
    Title         string            `json:"title" db:"title"`
    Content       string            `json:"content" db:"content"`
    ContentType   string            `json:"content_type" db:"content_type"`
    Category      string            `json:"category" db:"category"`
    Priority      int               `json:"priority" db:"priority"`
    Confidence    float64           `json:"confidence" db:"confidence"`
    Source        string            `json:"source,omitempty" db:"source"`
    Tags          []string          `json:"tags" db:"tags"`
    MITRETactic   string            `json:"mitre_tactic,omitempty" db:"mitre_tactic"`
    MITRETechnique string           `json:"mitre_technique,omitempty" db:"mitre_technique"`
    KillChainPhase string           `json:"kill_chain_phase,omitempty" db:"kill_chain_phase"`
    AccessCount   int               `json:"access_count" db:"access_count"`
    LastAccessed  *time.Time        `json:"last_accessed,omitempty" db:"last_accessed"`
    CreatedAt     time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
    Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
}
```

## Step 4: Implement Repository Layer

Create `internal/repository/memory_repo.go`:

```go
package repository

import (
    "context"
    "fmt"
    
    "github.com/rainmana/tinybrain/internal/database"
    "github.com/rainmana/tinybrain/internal/models"
    "github.com/google/uuid"
)

type MemoryRepository struct {
    db *database.SupabaseClient
}

func NewMemoryRepository(db *database.SupabaseClient) *MemoryRepository {
    return &MemoryRepository{db: db}
}

func (r *MemoryRepository) Create(ctx context.Context, memory *models.Memory) error {
    // Set UUID if not provided
    if memory.ID == uuid.Nil {
        memory.ID = uuid.New()
    }
    
    // Insert into Supabase
    var result []models.Memory
    err := r.db.Client.DB.From("memories").Insert(memory).Execute(&result)
    if err != nil {
        return fmt.Errorf("failed to create memory: %w", err)
    }
    
    if len(result) > 0 {
        *memory = result[0]
    }
    
    return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Memory, error) {
    var result []models.Memory
    err := r.db.Client.DB.From("memories").
        Select("*").
        Eq("id", id.String()).
        Single().
        Execute(&result)
    
    if err != nil {
        return nil, fmt.Errorf("failed to get memory: %w", err)
    }
    
    if len(result) == 0 {
        return nil, fmt.Errorf("memory not found")
    }
    
    return &result[0], nil
}

func (r *MemoryRepository) List(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Memory, error) {
    var result []models.Memory
    query := r.db.Client.DB.From("memories").
        Select("*").
        Eq("user_id", userID.String()).
        Order("created_at", &supabase.OrderOpts{Ascending: false})
    
    if limit > 0 {
        query = query.Limit(limit, "")
    }
    
    err := query.Execute(&result)
    if err != nil {
        return nil, fmt.Errorf("failed to list memories: %w", err)
    }
    
    // Convert to pointer slice
    memories := make([]*models.Memory, len(result))
    for i := range result {
        memories[i] = &result[i]
    }
    
    return memories, nil
}

// Add more methods: Update, Delete, Search, etc.
```

## Step 5: Create REST API Handlers

Create `internal/api/rest/memories.go`:

```go
package rest

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/rainmana/tinybrain/internal/repository"
    "github.com/rainmana/tinybrain/internal/models"
    "github.com/google/uuid"
)

type MemoryHandler struct {
    repo *repository.MemoryRepository
}

func NewMemoryHandler(repo *repository.MemoryRepository) *MemoryHandler {
    return &MemoryHandler{repo: repo}
}

// GET /api/memories
func (h *MemoryHandler) List(w http.ResponseWriter, r *http.Request) {
    // Get user ID from context (set by auth middleware)
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    memories, err := h.repo.List(r.Context(), userID, 50)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "memories": memories,
        "count":    len(memories),
    })
}

// POST /api/memories
func (h *MemoryHandler) Create(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    var memory models.Memory
    if err := json.NewDecoder(r.Body).Decode(&memory); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    memory.UserID = userID
    
    if err := h.repo.Create(r.Context(), &memory); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(memory)
}

// GET /api/memories/:id
func (h *MemoryHandler) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    memory, err := h.repo.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(memory)
}

// Register routes
func (h *MemoryHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/memories", h.List).Methods("GET")
    r.HandleFunc("/api/memories", h.Create).Methods("POST")
    r.HandleFunc("/api/memories/{id}", h.Get).Methods("GET")
    // Add PUT, DELETE, etc.
}
```

## Step 6: Implement Authentication Middleware

Create `internal/api/middleware/auth.go`:

```go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Parse and validate JWT
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            // Return secret key
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Extract user ID from claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }
        
        userIDStr, ok := claims["sub"].(string)
        if !ok {
            http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
            return
        }
        
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
            http.Error(w, "Invalid user ID format", http.StatusUnauthorized)
            return
        }
        
        // Add user ID to context
        ctx := context.WithValue(r.Context(), "user_id", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Step 7: Update Main Server

Update `cmd/tinybrain/main.go`:

```go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    
    "github.com/rainmana/tinybrain/internal/database"
    "github.com/rainmana/tinybrain/internal/repository"
    "github.com/rainmana/tinybrain/internal/api/rest"
    "github.com/rainmana/tinybrain/internal/api/middleware"
)

func main() {
    // Initialize Supabase client
    db, err := database.NewSupabaseClient()
    if err != nil {
        log.Fatalf("Failed to initialize Supabase: %v", err)
    }
    
    // Initialize repositories
    memoryRepo := repository.NewMemoryRepository(db)
    
    // Initialize handlers
    memoryHandler := rest.NewMemoryHandler(memoryRepo)
    
    // Create router
    r := mux.NewRouter()
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    }).Methods("GET")
    
    // Register API routes with auth middleware
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)
    memoryHandler.RegisterRoutes(api)
    
    // CORS configuration
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    }).Handler(r)
    
    // Start server
    addr := os.Getenv("TINYBRAIN_HTTP")
    if addr == "" {
        addr = "127.0.0.1:8090"
    }
    
    log.Printf("Starting server on %s", addr)
    if err := http.ListenAndServe(addr, corsHandler); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

## Step 8: Local Testing

1. **Set up environment variables:**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your Supabase credentials
   ```

2. **Run migrations:**
   ```bash
   psql $DATABASE_URL -f supabase/migrations/001_initial_schema.sql
   psql $DATABASE_URL -f supabase/migrations/002_row_level_security.sql
   ```

3. **Run the server:**
   ```bash
   go run ./cmd/tinybrain serve
   ```

4. **Test endpoints:**
   ```bash
   # Health check
   curl http://localhost:8090/health
   
   # Create memory (requires auth token)
   curl -X POST http://localhost:8090/api/memories \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Test Memory",
       "content": "This is a test",
       "category": "note",
       "priority": 5
     }'
   ```

## Step 9: Deploy to Railway

1. **Push changes to GitHub:**
   ```bash
   git add .
   git commit -m "Implement Phase 2: Backend API with Supabase"
   git push
   ```

2. **Railway will automatically deploy** (if connected to GitHub)

3. **Set environment variables in Railway dashboard:**
   - SUPABASE_URL
   - SUPABASE_SERVICE_KEY
   - DATABASE_URL
   - JWT_SECRET
   - CORS_ALLOWED_ORIGINS

4. **Monitor deployment logs:**
   ```bash
   railway logs
   ```

## Next Steps

After completing these core steps:

1. **Implement remaining repositories** (sessions, relationships, etc.)
2. **Add more API endpoints** (search, filters, aggregations)
3. **Implement WebSocket** for real-time features
4. **Add MCP adapter** for backward compatibility
5. **Write comprehensive tests**
6. **Add API documentation** (Swagger/OpenAPI)

## Useful Resources

- [Supabase Go Client Docs](https://github.com/supabase-community/supabase-go)
- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
- [JWT Go Documentation](https://github.com/golang-jwt/jwt)
- [Railway Documentation](https://docs.railway.app)

## Troubleshooting

### Common Issues

**Issue: Supabase connection fails**
- Verify SUPABASE_URL and SUPABASE_SERVICE_KEY are correct
- Check database is accessible (test with psql)
- Verify RLS policies allow service role

**Issue: Authentication errors**
- Ensure JWT_SECRET matches between services
- Verify token format (Bearer prefix)
- Check token expiration

**Issue: CORS errors**
- Add frontend origin to CORS_ALLOWED_ORIGINS
- Verify CORS middleware is applied
- Check browser console for specific CORS error

## Summary

This quick start guide provides the foundation for Phase 2 backend development. Follow the steps to:

1. Set up Supabase client
2. Update models for PostgreSQL
3. Implement repository pattern
4. Create REST API endpoints
5. Add authentication middleware
6. Update main server
7. Test locally
8. Deploy to Railway

For detailed architecture and design decisions, refer to `docs/WEB_ARCHITECTURE.md`.
