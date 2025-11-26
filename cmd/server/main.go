package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// TinyBrainPocketBaseServer combines MCP and PocketBase in a single binary
type TinyBrainPocketBaseServer struct {
	app    *pocketbase.PocketBase
	logger *log.Logger
}

func NewTinyBrainPocketBaseServer() *TinyBrainPocketBaseServer {
	app := pocketbase.New()

	// Set up logging
	logger := log.New(os.Stderr, "TinyBrain ", log.LstdFlags)

	server := &TinyBrainPocketBaseServer{
		app:    app,
		logger: logger,
	}

	// Set up PocketBase hooks and custom routes
	server.setupPocketBaseHooks()
	server.setupCustomRoutes()
	server.setupCollections()

	return server
}

func (s *TinyBrainPocketBaseServer) setupCollections() {
	s.logger.Println("Setting up TinyBrain collections...")

	// For now, just log that we would set up collections
	// This is a safe approach that doesn't break existing functionality
	s.logger.Println("Collections setup will be implemented in the next phase")
	s.logger.Println("Current version uses mock responses for all MCP tools")
	s.logger.Println("This ensures the app remains working and testable")
	s.logger.Println("Real database operations will be added gradually")
}

func (s *TinyBrainPocketBaseServer) setupPocketBaseHooks() {
	// For now, just log that we would set up hooks
	s.logger.Println("PocketBase hooks will be set up after collections are created")
}

func (s *TinyBrainPocketBaseServer) setupCustomRoutes() {
	s.app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// MCP endpoint - maintain compatibility with existing MCP tools
		e.Router.POST("/mcp", func(re *core.RequestEvent) error {
			// Handle MCP JSON-RPC requests
			var mcpRequest MCPRequest
			if err := re.BindBody(&mcpRequest); err != nil {
				return re.BadRequestError("Invalid MCP request", err)
			}

			// Get request info, handle error appropriately
			requestInfo, err := re.RequestInfo()
			if err != nil {
				// Log the error but continue with nil requestInfo
				// handleMCPRequest doesn't currently use requestInfo, so this is safe
				s.logger.Printf("Warning: Failed to get request info: %v", err)
				requestInfo = nil
			}

			// Process through MCP handler
			response, err := s.handleMCPRequest(requestInfo, mcpRequest)
			if err != nil {
				return re.InternalServerError("MCP processing failed", err)
			}

			return re.JSON(http.StatusOK, response)
		})

		// Enhanced security data endpoints using PocketBase
		e.Router.GET("/api/security/nvd", func(re *core.RequestEvent) error {
			// For now, return mock data until collections are set up
			return re.JSON(http.StatusOK, map[string]interface{}{
				"message": "NVD data endpoint - collections not yet set up",
			})
		})

		// Real-time memory search endpoint
		e.Router.GET("/api/memories/search", func(re *core.RequestEvent) error {
			query := re.Request.URL.Query().Get("q")
			if query == "" {
				return re.BadRequestError("Query parameter required", nil)
			}

			// For now, return mock data until collections are set up
			return re.JSON(http.StatusOK, map[string]interface{}{
				"message": "Memory search endpoint - collections not yet set up",
				"query":   query,
			})
		})

		return e.Next()
	})
}

// MCP Request/Response structures
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

func (s *TinyBrainPocketBaseServer) handleMCPRequest(requestInfo *core.RequestInfo, req MCPRequest) (MCPResponse, error) {
	s.logger.Printf("Handling MCP request: %s", req.Method)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "create_session":
		return s.handleCreateSession(req)
	case "store_memory":
		return s.handleStoreMemory(req)
	case "search_memories":
		return s.handleSearchMemories(req)
	case "get_session":
		return s.handleGetSession(req)
	case "list_sessions":
		return s.handleListSessions(req)
	case "create_relationship":
		return s.handleCreateRelationship(req)
	case "get_related_entries":
		return s.handleGetRelatedEntries(req)
	case "create_context_snapshot":
		return s.handleCreateContextSnapshot(req)
	case "get_context_snapshot":
		return s.handleGetContextSnapshot(req)
	case "list_context_snapshots":
		return s.handleListContextSnapshots(req)
	case "create_task_progress":
		return s.handleCreateTaskProgress(req)
	case "update_task_progress":
		return s.handleUpdateTaskProgress(req)
	case "list_task_progress":
		return s.handleListTaskProgress(req)
	case "get_memory_stats":
		return s.handleGetMemoryStats(req)
	case "get_system_diagnostics":
		return s.handleGetSystemDiagnostics(req)
	case "health_check":
		return s.handleHealthCheck(req)
	case "download_security_data":
		return s.handleDownloadSecurityData(req)
	case "get_security_data_summary":
		return s.handleGetSecurityDataSummary(req)
	case "query_nvd":
		return s.handleQueryNVD(req)
	case "query_attack":
		return s.handleQueryATTACK(req)
	case "query_owasp":
		return s.handleQueryOWASP(req)
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

// MCP Tool Handlers - All return mock responses for now
func (s *TinyBrainPocketBaseServer) handleInitialize(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "TinyBrain Memory Storage",
				"version": "1.0.0",
			},
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleToolsList(req MCPRequest) (MCPResponse, error) {
	tools := []map[string]interface{}{
		{
			"name":        "create_session",
			"description": "Create a new security assessment session",
		},
		{
			"name":        "store_memory",
			"description": "Store a new piece of information in memory",
		},
		{
			"name":        "search_memories",
			"description": "Search for memories using various strategies",
		},
		{
			"name":        "get_session",
			"description": "Get session details by ID",
		},
		{
			"name":        "list_sessions",
			"description": "List all sessions with optional filtering",
		},
		{
			"name":        "create_relationship",
			"description": "Create a relationship between two memory entries",
		},
		{
			"name":        "get_related_entries",
			"description": "Get memory entries related to a specific entry",
		},
		{
			"name":        "create_context_snapshot",
			"description": "Create a snapshot of the current context",
		},
		{
			"name":        "get_context_snapshot",
			"description": "Get a context snapshot by ID",
		},
		{
			"name":        "list_context_snapshots",
			"description": "List context snapshots for a session",
		},
		{
			"name":        "create_task_progress",
			"description": "Create a new task progress entry",
		},
		{
			"name":        "update_task_progress",
			"description": "Update progress on a task",
		},
		{
			"name":        "list_task_progress",
			"description": "List task progress entries for a session",
		},
		{
			"name":        "get_memory_stats",
			"description": "Get comprehensive statistics about memory usage",
		},
		{
			"name":        "get_system_diagnostics",
			"description": "Get system diagnostics and debugging information",
		},
		{
			"name":        "health_check",
			"description": "Perform a health check on the database and server",
		},
		{
			"name":        "download_security_data",
			"description": "Download security datasets from external sources (NVD, ATT&CK, OWASP)",
		},
		{
			"name":        "get_security_data_summary",
			"description": "Get summary of security data in the knowledge hub",
		},
		{
			"name":        "query_nvd",
			"description": "Query NVD CVE data from the security knowledge hub",
		},
		{
			"name":        "query_attack",
			"description": "Query MITRE ATT&CK data from the security knowledge hub",
		},
		{
			"name":        "query_owasp",
			"description": "Query OWASP testing procedures from the security knowledge hub",
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateSession(req MCPRequest) (MCPResponse, error) {
	// Parse parameters
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	name, _ := params["name"].(string)
	_, _ = params["description"].(string)
	taskType, _ := params["task_type"].(string)
	if taskType == "" {
		taskType = "general"
	}

	// Return mock response for now
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"session_id": "mock-session-id",
			"name":       name,
			"status":     "active",
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleStoreMemory(req MCPRequest) (MCPResponse, error) {
	// Parse parameters
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	title, _ := params["title"].(string)
	category, _ := params["category"].(string)

	// Return mock response for now
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"memory_id": "mock-memory-id",
			"title":     title,
			"category":  category,
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleSearchMemories(req MCPRequest) (MCPResponse, error) {
	// Parse parameters
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	query, _ := params["query"].(string)
	// Note: limit parameter will be used when real database operations are implemented
	_ = params["limit"] // Acknowledge parameter exists but not used in mock

	// Return mock response for now
	results := make([]map[string]interface{}, 0)
	if query != "" {
		results = append(results, map[string]interface{}{
			"id":         "mock-memory-1",
			"title":      "Mock Memory for: " + query,
			"content":    "This is a mock memory result",
			"category":   "note",
			"priority":   5.0,
			"confidence": 0.8,
		})
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"memories": results,
			"count":    len(results),
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetSession(req MCPRequest) (MCPResponse, error) {
	// Parse parameters
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	sessionID, _ := params["session_id"].(string)
	if sessionID == "" {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "session_id is required",
			},
		}, nil
	}

	// Return mock response for now
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"session_id":  sessionID,
			"name":        "Mock Session",
			"description": "This is a mock session",
			"task_type":   "security_review",
			"status":      "active",
		},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleListSessions(req MCPRequest) (MCPResponse, error) {
	// Parse parameters
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	// Note: limit parameter will be used when real database operations are implemented
	_ = params["limit"] // Acknowledge parameter exists but not used in mock

	// Return mock response for now
	results := []map[string]interface{}{
		{
			"session_id":  "mock-session-1",
			"name":        "Mock Security Review",
			"description": "A mock security review session",
			"task_type":   "security_review",
			"status":      "active",
		},
		{
			"session_id":  "mock-session-2",
			"name":        "Mock Penetration Test",
			"description": "A mock penetration test session",
			"task_type":   "penetration_test",
			"status":      "active",
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"sessions": results,
			"count":    len(results),
		},
	}, nil
}

// Placeholder handlers for other MCP tools
func (s *TinyBrainPocketBaseServer) handleCreateRelationship(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetRelatedEntries(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateContextSnapshot(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetContextSnapshot(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleListContextSnapshots(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateTaskProgress(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleUpdateTaskProgress(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleListTaskProgress(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetMemoryStats(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetSystemDiagnostics(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleHealthCheck(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleDownloadSecurityData(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetSecurityDataSummary(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryNVD(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryATTACK(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryOWASP(req MCPRequest) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) Start() error {
	// Set up data directory
	dataDir := filepath.Join(os.Getenv("HOME"), ".tinybrain")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Configure PocketBase
	// Note: DataDir is set through command line flags, not directly

	s.logger.Printf("Starting TinyBrain with PocketBase backend, data_dir: %s", dataDir)

	// Start PocketBase (includes database, REST API, real-time, admin UI)
	return s.app.Start()
}

func main() {
	// Create the combined TinyBrain + PocketBase server
	server := NewTinyBrainPocketBaseServer()

	// Start the server (single binary with both MCP and PocketBase)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
