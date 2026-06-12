package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// MCPRequest represents a generic MCP request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents a generic MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolHandler is a function that handles tool calls
type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// Server represents an MCP server
type Server struct {
	name        string
	version     string
	description string
	tools       map[string]Tool
	handlers    map[string]ToolHandler
	logger      *log.Logger
	mu          sync.RWMutex
}

// NewServer creates a new MCP server
func NewServer(name, version, description string, logger *log.Logger) *Server {
	return &Server{
		name:        name,
		version:     version,
		description: description,
		tools:       make(map[string]Tool),
		handlers:    make(map[string]ToolHandler),
		logger:      logger,
	}
}

// AddTool adds a tool to the server
func (s *Server) AddTool(name, description string, inputSchema map[string]interface{}, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools[name] = Tool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
	}
	s.handlers[name] = handler
}

// AddToolAlias registers an existing tool under an additional name,
// sharing the original tool's schema and handler.
func (s *Server) AddToolAlias(alias, existing, description string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tool, ok := s.tools[existing]
	if !ok {
		s.logger.Error("Cannot alias unknown tool", "alias", alias, "existing", existing)
		return
	}
	if description == "" {
		description = tool.Description
	}
	s.tools[alias] = Tool{
		Name:        alias,
		Description: description,
		InputSchema: tool.InputSchema,
	}
	s.handlers[alias] = s.handlers[existing]
}

// ServeStdio starts the server using stdio transport
func (s *Server) ServeStdio() error {
	s.logger.Info("Starting MCP server", "name", s.name, "version", s.version)

	scanner := bufio.NewScanner(os.Stdin)
	// Allow large payloads (memory contents, context snapshots) beyond the
	// 64KB bufio default
	scanner.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.logger.Error("Failed to parse request", "error", err, "line", line)
			continue
		}

		response := s.HandleRequest(context.Background(), &req)
		if response == nil {
			// Notifications expect no response
			continue
		}
		responseData, err := json.Marshal(response)
		if err != nil {
			s.logger.Error("Failed to marshal response", "error", err)
			continue
		}

		fmt.Println(string(responseData))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// HandleRequest handles an incoming MCP request. It returns nil for
// notifications, which must not receive a response.
func (s *Server) HandleRequest(ctx context.Context, req *MCPRequest) *MCPResponse {
	s.logger.Debug("Handling request", "method", req.Method, "id", req.ID)

	// MCP notifications (e.g. notifications/initialized, notifications/cancelled)
	// must be accepted silently
	if strings.HasPrefix(req.Method, "notifications/") {
		return nil
	}

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "ping":
		return &MCPResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]interface{}{}}
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(req *MCPRequest) *MCPResponse {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.name,
			"version": s.version,
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList handles the tools/list request
func (s *Server) handleToolsList(req *MCPRequest) *MCPResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tools []Tool
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsCall handles the tools/call request
func (s *Server) handleToolsCall(ctx context.Context, req *MCPRequest) *MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	name, ok := params["name"].(string)
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid tool name",
			},
		}
	}

	// The published docs show tool names prefixed by the MCP client's server
	// label (e.g. mcp_tinybrain-mcp-server_create_session); accept that form
	name = strings.TrimPrefix(name, "mcp_tinybrain-mcp-server_")

	s.mu.RLock()
	handler, exists := s.handlers[name]
	s.mu.RUnlock()

	if !exists {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	start := time.Now()
	result, err := handler(ctx, arguments)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("Tool execution failed", "tool", name, "error", err, "duration", duration)
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
		}
	}

	s.logger.Debug("Tool executed successfully", "tool", name, "duration", duration)

	// Convert result to JSON string for proper formatting
	resultJSON, err := json.Marshal(result)
	if err != nil {
		s.logger.Error("Failed to marshal result to JSON", "error", err)
		resultJSON = []byte(fmt.Sprintf("%v", result))
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": string(resultJSON),
				},
			},
		},
	}
}
