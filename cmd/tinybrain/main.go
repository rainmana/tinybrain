package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
	"strings"

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
	case "tools/call":
		return s.handleToolsCall(req)
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

// handleToolsCall handles the tools/call method which is the proper MCP protocol
func (s *TinyBrainPocketBaseServer) handleToolsCall(req MCPRequest) (MCPResponse, error) {
	// Parse params to get tool name and arguments
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

	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Tool name is required",
			},
		}, nil
	}

	// Extract arguments (may be nil or empty)
	arguments, _ := params["arguments"].(map[string]interface{})
	if arguments == nil {
		arguments = make(map[string]interface{})
	}

	// Route to appropriate handler based on tool name
	// Remove the "mcp_tinybrain-" prefix if present (some MCP clients add it)
	toolName = strings.TrimPrefix(toolName, "mcp_tinybrain-")
	toolName = strings.TrimPrefix(toolName, "tinybrain-")

	switch toolName {
	case "create_session":
		return s.handleCreateSessionWithArgs(req, arguments)
	case "store_memory":
		return s.handleStoreMemoryWithArgs(req, arguments)
	case "search_memories":
		return s.handleSearchMemoriesWithArgs(req, arguments)
	case "get_session":
		return s.handleGetSessionWithArgs(req, arguments)
	case "list_sessions":
		return s.handleListSessionsWithArgs(req, arguments)
	case "create_relationship":
		return s.handleCreateRelationshipWithArgs(req, arguments)
	case "get_related_memories":
		return s.handleGetRelatedEntriesWithArgs(req, arguments)
	case "create_context_snapshot":
		return s.handleCreateContextSnapshotWithArgs(req, arguments)
	case "get_context_snapshot":
		return s.handleGetContextSnapshotWithArgs(req, arguments)
	case "list_context_snapshots":
		return s.handleListContextSnapshotsWithArgs(req, arguments)
	case "create_task_progress":
		return s.handleCreateTaskProgressWithArgs(req, arguments)
	case "update_task_progress":
		return s.handleUpdateTaskProgressWithArgs(req, arguments)
	case "list_task_progress":
		return s.handleListTaskProgressWithArgs(req, arguments)
	case "get_memory_stats":
		return s.handleGetMemoryStatsWithArgs(req, arguments)
	case "get_system_diagnostics":
		return s.handleGetSystemDiagnosticsWithArgs(req, arguments)
	case "health_check":
		return s.handleHealthCheckWithArgs(req, arguments)
	case "download_security_data":
		return s.handleDownloadSecurityDataWithArgs(req, arguments)
	case "get_security_data_summary":
		return s.handleGetSecurityDataSummaryWithArgs(req, arguments)
	case "query_nvd":
		return s.handleQueryNVDWithArgs(req, arguments)
	case "query_attack":
		return s.handleQueryATTACKWithArgs(req, arguments)
	case "query_owasp":
		return s.handleQueryOWASPWithArgs(req, arguments)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", toolName),
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

// Wrapper functions that extract arguments from the tools/call format
func (s *TinyBrainPocketBaseServer) handleCreateSessionWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	name, _ := arguments["name"].(string)
	_, _ = arguments["description"].(string)
	taskType, _ := arguments["task_type"].(string)
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

// Legacy handler (kept for backward compatibility if needed)
func (s *TinyBrainPocketBaseServer) handleCreateSession(req MCPRequest) (MCPResponse, error) {
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
	return s.handleCreateSessionWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleStoreMemoryWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	title, _ := arguments["title"].(string)
	category, _ := arguments["category"].(string)

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

func (s *TinyBrainPocketBaseServer) handleStoreMemory(req MCPRequest) (MCPResponse, error) {
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
	return s.handleStoreMemoryWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleSearchMemoriesWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	query, _ := arguments["query"].(string)
	_ = arguments["limit"] // Acknowledge parameter exists but not used in mock

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

func (s *TinyBrainPocketBaseServer) handleSearchMemories(req MCPRequest) (MCPResponse, error) {
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
	return s.handleSearchMemoriesWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetSessionWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	sessionID, _ := arguments["session_id"].(string)
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

func (s *TinyBrainPocketBaseServer) handleGetSession(req MCPRequest) (MCPResponse, error) {
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
	return s.handleGetSessionWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleListSessionsWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	_ = arguments["limit"] // Acknowledge parameter exists but not used in mock

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

func (s *TinyBrainPocketBaseServer) handleListSessions(req MCPRequest) (MCPResponse, error) {
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
	return s.handleListSessionsWithArgs(req, params)
}

// Wrapper functions for remaining tools
func (s *TinyBrainPocketBaseServer) handleCreateRelationshipWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateRelationship(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleCreateRelationshipWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetRelatedEntriesWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetRelatedEntries(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleGetRelatedEntriesWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleCreateContextSnapshotWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateContextSnapshot(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleCreateContextSnapshotWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetContextSnapshotWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetContextSnapshot(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleGetContextSnapshotWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleListContextSnapshotsWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleListContextSnapshots(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleListContextSnapshotsWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleCreateTaskProgressWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleCreateTaskProgress(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleCreateTaskProgressWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleUpdateTaskProgressWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleUpdateTaskProgress(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleUpdateTaskProgressWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleListTaskProgressWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleListTaskProgress(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleListTaskProgressWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetMemoryStatsWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetMemoryStats(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleGetMemoryStatsWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetSystemDiagnosticsWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetSystemDiagnostics(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleGetSystemDiagnosticsWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleHealthCheckWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleHealthCheck(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleHealthCheckWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleDownloadSecurityDataWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleDownloadSecurityData(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleDownloadSecurityDataWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleGetSecurityDataSummaryWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleGetSecurityDataSummary(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleGetSecurityDataSummaryWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleQueryNVDWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryNVD(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleQueryNVDWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleQueryATTACKWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryATTACK(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleQueryATTACKWithArgs(req, params)
}

func (s *TinyBrainPocketBaseServer) handleQueryOWASPWithArgs(req MCPRequest, arguments map[string]interface{}) (MCPResponse, error) {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"status": "not implemented yet"},
	}, nil
}

func (s *TinyBrainPocketBaseServer) handleQueryOWASP(req MCPRequest) (MCPResponse, error) {
	params, _ := req.Params.(map[string]interface{})
	if params == nil {
		params = make(map[string]interface{})
	}
	return s.handleQueryOWASPWithArgs(req, params)
}

func printUsage() {
	fmt.Println()
	fmt.Println("🧠 TinyBrain MCP Server")
	fmt.Println("Security-focused LLM memory storage with intelligence gathering")
	fmt.Println()
	
	fmt.Println("USAGE:")
	fmt.Println("  tinybrain [command] [flags]")
	fmt.Println()
	
	fmt.Println("COMMANDS:")
	fmt.Println("  serve     Start the TinyBrain server (default)")
	fmt.Println("  --help    Show this help message")
	fmt.Println()
	
	fmt.Println("FLAGS:")
	fmt.Println("  --http=<address>  HTTP bind address (default: 127.0.0.1:8090)")
	fmt.Println("  --dir=<path>      Data directory (default: ~/.tinybrain)")
	fmt.Println()
	
	fmt.Println("EXAMPLES:")
	fmt.Println("  tinybrain                              # Start with defaults")
	fmt.Println("  tinybrain --http=127.0.0.1:9000        # Custom port")
	fmt.Println("  tinybrain serve --http=0.0.0.0:8090    # Bind to all interfaces")
	fmt.Println("  TINYBRAIN_HTTP=:9000 tinybrain         # Port via environment")
	fmt.Println()
	
	fmt.Println("ENVIRONMENT:")
	fmt.Println("  TINYBRAIN_HTTP      HTTP bind address")
	fmt.Println("  TINYBRAIN_DATA_DIR  Data directory")
	fmt.Println()
	
	fmt.Println("For more info: https://github.com/rainmana/tinybrain")
}


func main() {
	// Handle --help or -h
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		printUsage()
		return
	}
	
	// If no args provided, default to "serve"
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "serve")
	}
	
	// If first arg isn't "serve", assume they want to serve with those flags
	if os.Args[1] != "serve" {
		os.Args = append([]string{os.Args[0], "serve"}, os.Args[1:]...)
	}
	
	// Create the combined TinyBrain + PocketBase server
	app := pocketbase.New()
	logger := log.New(os.Stderr, "TinyBrain ", log.LstdFlags)
	
	server := &TinyBrainPocketBaseServer{
		app:    app,
		logger: logger,
	}
	
	// Setup before serving
	server.setupPocketBaseHooks()
	server.setupCustomRoutes()
	server.setupCollections()
	
	// Handle TINYBRAIN_HTTP environment variable
	httpAddr := os.Getenv("TINYBRAIN_HTTP")
	if httpAddr != "" {
		hasHTTPFlag := false
		for _, arg := range os.Args {
			if strings.HasPrefix(arg, "--http") {
				hasHTTPFlag = true
				break
			}
		}
		if !hasHTTPFlag {
			os.Args = append(os.Args, "--http="+httpAddr)
			logger.Printf("Using HTTP address from TINYBRAIN_HTTP: %s", httpAddr)
		}
	}
	
	// Setup data directory
	dataDir := filepath.Join(os.Getenv("HOME"), ".tinybrain")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.Fatalf("Failed to create data directory: %v", err)
	}
	
	logger.Printf("TinyBrain data directory: %s", dataDir)
	logger.Println("Starting TinyBrain MCP Server")
	logger.Println("Run 'tinybrain --help' for usage information")
	
	// Execute PocketBase
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
