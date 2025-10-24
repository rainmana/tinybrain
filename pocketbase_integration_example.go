package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// TinyBrainPocketBaseServer combines MCP and PocketBase in a single binary
type TinyBrainPocketBaseServer struct {
	app *pocketbase.PocketBase
	mcp *MCPServer // Your existing MCP server
}

func NewTinyBrainPocketBaseServer() *TinyBrainPocketBaseServer {
	app := pocketbase.New()

	// Initialize your existing MCP server
	mcp := NewMCPServer()

	server := &TinyBrainPocketBaseServer{
		app: app,
		mcp: mcp,
	}

	// Set up PocketBase hooks and custom routes
	server.setupPocketBaseHooks()
	server.setupCustomRoutes()

	return server
}

func (s *TinyBrainPocketBaseServer) setupPocketBaseHooks() {
	// Real-time memory updates
	s.app.OnRecordAfterCreate("memories").BindFunc(func(e *core.RecordEvent) error {
		// Notify all connected clients about new memory
		s.app.RealtimeServer().NotifyRecord(e.Record)
		return e.Next()
	})

	// Security data access control
	s.app.OnRecordViewRequest("nvd_cves").BindFunc(func(e *core.RecordRequestEvent) error {
		requestInfo, _ := e.RequestInfo()
		if requestInfo.Auth == nil {
			return e.ForbiddenError("Authentication required for security data", nil)
		}
		return e.Next()
	})

	// Auto-populate memory metadata
	s.app.OnRecordBeforeCreate("memories").BindFunc(func(e *core.RecordEvent) error {
		// Auto-set session and metadata
		if e.RequestInfo.Auth != nil {
			e.Record.Set("user", e.RequestInfo.Auth.Id)
		}
		e.Record.Set("created_at", core.NowDateTime())
		return e.Next()
	})
}

func (s *TinyBrainPocketBaseServer) setupCustomRoutes() {
	s.app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// MCP endpoint - keep your existing MCP functionality
		e.Router.POST("/mcp", func(re *core.RequestEvent) error {
			// Handle MCP JSON-RPC requests
			var mcpRequest MCPRequest
			if err := re.BindBody(&mcpRequest); err != nil {
				return re.BadRequestError("Invalid MCP request", err)
			}

			// Process through your existing MCP server
			response, err := s.mcp.HandleRequest(context.Background(), mcpRequest)
			if err != nil {
				return re.InternalServerError("MCP processing failed", err)
			}

			return re.JSON(http.StatusOK, response)
		})

		// Enhanced security data endpoints using PocketBase
		e.Router.GET("/api/security/nvd", func(re *core.RequestEvent) error {
			requestInfo, _ := re.RequestInfo()
			if requestInfo.Auth == nil {
				return re.UnauthorizedError("Authentication required", nil)
			}

			// Use PocketBase's built-in filtering
			severity := re.Request.URL.Query().Get("severity")
			filter := ""
			if severity != "" {
				filter = `severity = "` + severity + `"`
			}

			// Query using PocketBase's built-in API
			records, err := s.app.RecordQuery("nvd_cves").
				AndWhere(dbx.NewExp(filter)).
				Limit(100).
				All()
			if err != nil {
				return re.InternalServerError("Failed to query NVD data", err)
			}

			return re.JSON(http.StatusOK, records)
		}, apis.RequireAuth())

		// Real-time memory search endpoint
		e.Router.GET("/api/memories/search", func(re *core.RequestEvent) error {
			query := re.Request.URL.Query().Get("q")
			if query == "" {
				return re.BadRequestError("Query parameter required", nil)
			}

			// Use PocketBase's full-text search
			records, err := s.app.RecordQuery("memories").
				AndWhere(dbx.NewExp("content LIKE {:query}", dbx.Params{"query": "%" + query + "%"})).
				Limit(50).
				All()
			if err != nil {
				return re.InternalServerError("Search failed", err)
			}

			return re.JSON(http.StatusOK, records)
		}, apis.RequireAuth())

		// Security knowledge hub endpoints
		e.Router.GET("/api/security/attack", func(re *core.RequestEvent) error {
			// Query ATT&CK techniques using PocketBase
			records, err := s.app.RecordQuery("attack_techniques").
				AndWhere(dbx.NewExp("tactic = {:tactic}", dbx.Params{"tactic": re.Request.URL.Query().Get("tactic")})).
				Limit(100).
				All()
			if err != nil {
				return re.InternalServerError("Failed to query ATT&CK data", err)
			}

			return re.JSON(http.StatusOK, records)
		}, apis.RequireAuth())

		return e.Next()
	})
}

func (s *TinyBrainPocketBaseServer) Start() error {
	// Start PocketBase (includes database, REST API, real-time, admin UI)
	return s.app.Start()
}

// Your existing MCP server (simplified example)
type MCPServer struct {
	// Your existing MCP implementation
}

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewMCPServer() *MCPServer {
	return &MCPServer{}
}

func (m *MCPServer) HandleRequest(ctx context.Context, req MCPRequest) (MCPResponse, error) {
	// Your existing MCP tool handling logic
	switch req.Method {
	case "create_session":
		return m.handleCreateSession(req)
	case "store_memory":
		return m.handleStoreMemory(req)
	case "search_memories":
		return m.handleSearchMemories(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}, nil
	}
}

func (m *MCPServer) handleCreateSession(req MCPRequest) (MCPResponse, error) {
	// Your existing session creation logic
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"session_id": "new-session-id",
			"status":     "created",
		},
	}, nil
}

func (m *MCPServer) handleStoreMemory(req MCPRequest) (MCPResponse, error) {
	// Your existing memory storage logic
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"memory_id": "new-memory-id",
			"status":    "stored",
		},
	}, nil
}

func (m *MCPServer) handleSearchMemories(req MCPRequest) (MCPResponse, error) {
	// Your existing memory search logic
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"memories": []interface{}{},
			"count":    0,
		},
	}, nil
}

func main() {
	// Create the combined TinyBrain + PocketBase server
	server := NewTinyBrainPocketBaseServer()

	// Start the server (single binary with both MCP and PocketBase)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
