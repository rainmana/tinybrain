// +build ignore

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
		// MCP endpoint
		e.Router.POST("/mcp", func(re *core.RequestEvent) error {
			var mcpRequest MCPRequest
			if err := re.BindBody(&mcpRequest); err != nil {
				return re.BadRequestError("Invalid MCP request", err)
			}

			requestInfo, _ := re.RequestInfo()
			response, err := s.mcp.HandleRequest(requestInfo, mcpRequest)
			if err != nil {
				return re.InternalServerError("MCP processing failed", err)
			}

			return re.JSON(http.StatusOK, response)
		})

		// Security data endpoints
		e.Router.GET("/api/security/nvd", func(re *core.RequestEvent) error {
			// Query NVD data using PocketBase
			filter := "severity = 'CRITICAL'"
			records, err := s.app.RecordQuery("nvd_cves").AndWhere(dbx.NewExp(filter)).Limit(100).All()
			if err != nil {
				return re.InternalServerError("Failed to query NVD data", err)
			}

			return re.JSON(http.StatusOK, records)
		}, apis.RequireRecordAuth())

		// Memory search endpoint
		e.Router.GET("/api/memories/search", func(re *core.RequestEvent) error {
			query := re.Request.URL.Query().Get("q")
			if query == "" {
				return re.BadRequestError("Query parameter required", nil)
			}

			// Search memories using PocketBase
			filter := fmt.Sprintf("title ~ '%s' OR content ~ '%s'", query, query)
			records, err := s.app.RecordQuery("memories").AndWhere(dbx.NewExp(filter)).Limit(20).All()
			if err != nil {
				return re.InternalServerError("Failed to search memories", err)
			}

			return re.JSON(http.StatusOK, records)
		})

		return e.Next()
	})
}

// MCP Server interface (your existing implementation)
type MCPServer struct {
	// Your existing MCP server fields
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

