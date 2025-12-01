package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewTinyBrainPocketBaseServer tests server initialization
func TestNewTinyBrainPocketBaseServer(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	assert.NotNil(t, server)
	assert.NotNil(t, server.app)
	assert.NotNil(t, server.logger)
}

// TestSetupCollections tests collection setup
func TestSetupCollections(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	// Should not panic
	assert.NotPanics(t, func() {
		server.setupCollections()
	})
}

// TestSetupPocketBaseHooks tests hook setup
func TestSetupPocketBaseHooks(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	// Should not panic
	assert.NotPanics(t, func() {
		server.setupPocketBaseHooks()
	})
}

// TestHandleInitialize tests the initialize MCP handler
func TestHandleInitialize(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  nil,
	}
	
	resp, err := server.handleInitialize(req)
	
	require.NoError(t, err)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)
	assert.NotNil(t, resp.Result)
	
	result := resp.Result.(map[string]interface{})
	assert.Equal(t, "2024-11-05", result["protocolVersion"])
	assert.NotNil(t, result["capabilities"])
	assert.NotNil(t, result["serverInfo"])
	
	serverInfo := result["serverInfo"].(map[string]interface{})
	assert.Equal(t, "TinyBrain Memory Storage", serverInfo["name"])
	assert.Equal(t, "1.0.0", serverInfo["version"])
}

// TestHandleToolsList tests the tools/list MCP handler
func TestHandleToolsList(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  nil,
	}
	
	resp, err := server.handleToolsList(req)
	
	require.NoError(t, err)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)
	assert.NotNil(t, resp.Result)
	
	result := resp.Result.(map[string]interface{})
	tools := result["tools"].([]map[string]interface{})
	
	assert.Greater(t, len(tools), 0)
	
	// Check for expected tools
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		name := tool["name"].(string)
		toolNames[name] = true
		assert.NotEmpty(t, tool["description"])
	}
	
	assert.True(t, toolNames["create_session"])
	assert.True(t, toolNames["store_memory"])
	assert.True(t, toolNames["search_memories"])
	assert.True(t, toolNames["get_session"])
	assert.True(t, toolNames["list_sessions"])
}

// TestHandleCreateSession tests the create_session MCP handler
func TestHandleCreateSession(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	tests := []struct {
		name          string
		params        interface{}
		expectError   bool
		expectedName  string
		expectedType  string
	}{
		{
			name: "valid session creation",
			params: map[string]interface{}{
				"name":      "Test Session",
				"task_type": "security_review",
			},
			expectError:  false,
			expectedName: "Test Session",
			expectedType: "security_review",
		},
		{
			name: "session with default task type",
			params: map[string]interface{}{
				"name": "Test Session 2",
			},
			expectError:  false,
			expectedName: "Test Session 2",
			expectedType: "general",
		},
		{
			name: "invalid params type - not a map",
			params: "not a map", // Invalid type - not a map
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "create_session",
				Params:  tt.params,
			}
			
			resp, err := server.handleCreateSession(req)
			
			if tt.expectError {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, -32602, resp.Error.Code)
			} else {
				require.NoError(t, err)
				assert.Nil(t, resp.Error)
				assert.NotNil(t, resp.Result)
				
				result := resp.Result.(map[string]interface{})
				assert.Equal(t, tt.expectedName, result["name"])
				assert.Equal(t, "active", result["status"])
				assert.NotEmpty(t, result["session_id"])
			}
		})
	}
}

// TestHandleStoreMemory tests the store_memory MCP handler
func TestHandleStoreMemory(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	tests := []struct {
		name         string
		params       interface{}
		expectError  bool
		expectedTitle string
	}{
		{
			name: "valid memory storage",
			params: map[string]interface{}{
				"title":    "Test Memory",
				"category": "vulnerability",
			},
			expectError:   false,
			expectedTitle: "Test Memory",
		},
		{
			name: "invalid params type - not a map",
			params: "not a map", // Invalid type - not a map
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "store_memory",
				Params:  tt.params,
			}
			
			resp, err := server.handleStoreMemory(req)
			
			if tt.expectError {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, -32602, resp.Error.Code)
			} else {
				require.NoError(t, err)
				assert.Nil(t, resp.Error)
				assert.NotNil(t, resp.Result)
				
				result := resp.Result.(map[string]interface{})
				assert.Equal(t, tt.expectedTitle, result["title"])
				assert.NotEmpty(t, result["memory_id"])
			}
		})
	}
}

// TestHandleSearchMemories tests the search_memories MCP handler
func TestHandleSearchMemories(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	tests := []struct {
		name          string
		params        interface{}
		expectError   bool
		expectedCount int
	}{
		{
			name: "search with query",
			params: map[string]interface{}{
				"query": "test query",
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "search with empty query",
			params: map[string]interface{}{
				"query": "",
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name: "invalid params type - not a map",
			params: "not a map", // Invalid type - not a map
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "search_memories",
				Params:  tt.params,
			}
			
			resp, err := server.handleSearchMemories(req)
			
			if tt.expectError {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, -32602, resp.Error.Code)
			} else {
				require.NoError(t, err)
				assert.Nil(t, resp.Error)
				assert.NotNil(t, resp.Result)
				
				result := resp.Result.(map[string]interface{})
				memories := result["memories"].([]map[string]interface{})
				assert.Equal(t, tt.expectedCount, result["count"])
				assert.Equal(t, tt.expectedCount, len(memories))
			}
		})
	}
}

// TestHandleGetSession tests the get_session MCP handler
func TestHandleGetSession(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	tests := []struct {
		name         string
		params       interface{}
		expectError  bool
		errorCode    int
	}{
		{
			name: "valid session retrieval",
			params: map[string]interface{}{
				"session_id": "test-session-id",
			},
			expectError: false,
		},
		{
			name: "missing session_id",
			params: map[string]interface{}{
				"session_id": "",
			},
			expectError: true,
			errorCode:   -32602,
		},
		{
			name: "invalid params type - not a map",
			params: "not a map", // Invalid type - not a map
			expectError: true,
			errorCode:   -32602,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "get_session",
				Params:  tt.params,
			}
			
			resp, err := server.handleGetSession(req)
			
			if tt.expectError {
				assert.NotNil(t, resp.Error)
				if tt.errorCode != 0 {
					assert.Equal(t, tt.errorCode, resp.Error.Code)
				}
			} else {
				require.NoError(t, err)
				assert.Nil(t, resp.Error)
				assert.NotNil(t, resp.Result)
				
				result := resp.Result.(map[string]interface{})
				assert.Equal(t, "test-session-id", result["session_id"])
				assert.Equal(t, "active", result["status"])
			}
		})
	}
}

// TestHandleListSessions tests the list_sessions MCP handler
func TestHandleListSessions(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "list_sessions",
		Params:  map[string]interface{}{},
	}
	
	resp, err := server.handleListSessions(req)
	
	require.NoError(t, err)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
	
	result := resp.Result.(map[string]interface{})
	sessions := result["sessions"].([]map[string]interface{})
	count := result["count"].(int)
	
	assert.Greater(t, count, 0)
	assert.Equal(t, count, len(sessions))
}

// TestHandleMCPRequest tests the main MCP request router
func TestHandleMCPRequest(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	tests := []struct {
		name        string
		method      string
		params      interface{}
		expectError bool
		errorCode   int
	}{
		{"initialize", "initialize", nil, false, 0},
		{"tools/list", "tools/list", nil, false, 0},
		{"create_session", "create_session", map[string]interface{}{"name": "test"}, false, 0},
		{"store_memory", "store_memory", map[string]interface{}{"title": "test"}, false, 0},
		{"search_memories", "search_memories", map[string]interface{}{"query": "test"}, false, 0},
		{"get_session", "get_session", map[string]interface{}{"session_id": "test-id"}, false, 0},
		{"list_sessions", "list_sessions", map[string]interface{}{}, false, 0},
		{"unknown method", "unknown_method", nil, true, -32601},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  tt.method,
				Params:  tt.params,
			}
			
			requestInfo := &core.RequestInfo{}
			resp, err := server.handleMCPRequest(requestInfo, req)
			
			if tt.expectError {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, tt.errorCode, resp.Error.Code)
			} else {
				require.NoError(t, err)
				assert.Nil(t, resp.Error)
				assert.NotNil(t, resp.Result)
			}
		})
	}
}

// TestPlaceholderHandlers tests all placeholder MCP handlers
func TestPlaceholderHandlers(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	handlers := []struct {
		name   string
		method func(MCPRequest) (MCPResponse, error)
	}{
		{"create_relationship", server.handleCreateRelationship},
		{"get_related_entries", server.handleGetRelatedEntries},
		{"create_context_snapshot", server.handleCreateContextSnapshot},
		{"get_context_snapshot", server.handleGetContextSnapshot},
		{"list_context_snapshots", server.handleListContextSnapshots},
		{"create_task_progress", server.handleCreateTaskProgress},
		{"update_task_progress", server.handleUpdateTaskProgress},
		{"list_task_progress", server.handleListTaskProgress},
		{"get_memory_stats", server.handleGetMemoryStats},
		{"get_system_diagnostics", server.handleGetSystemDiagnostics},
		{"health_check", server.handleHealthCheck},
		{"download_security_data", server.handleDownloadSecurityData},
		{"get_security_data_summary", server.handleGetSecurityDataSummary},
		{"query_nvd", server.handleQueryNVD},
		{"query_attack", server.handleQueryATTACK},
		{"query_owasp", server.handleQueryOWASP},
	}
	
	for _, tt := range handlers {
		t.Run(tt.name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  tt.name,
				Params:  map[string]interface{}{},
			}
			
			resp, err := tt.method(req)
			
			require.NoError(t, err)
			assert.Nil(t, resp.Error)
			assert.NotNil(t, resp.Result)
			
			result := resp.Result.(map[string]interface{})
			assert.Equal(t, "not implemented yet", result["status"])
		})
	}
}

// TestMCPEndpoint tests the HTTP MCP endpoint handler
func TestMCPEndpoint(t *testing.T) {
	// Create a temporary data directory
	tempDir := filepath.Join(os.TempDir(), "tinybrain-test-"+t.Name())
	defer os.RemoveAll(tempDir)
	
	// Initialize PocketBase with test config
	config := pocketbase.Config{
		DefaultDataDir: tempDir,
	}
	app := pocketbase.NewWithConfig(config)
	
	// Bootstrap the app
	err := app.Bootstrap()
	require.NoError(t, err)
	
	// Create server
	server := NewTinyBrainPocketBaseServer()
	server.app = app
	
	// Test MCP endpoint with initialize request
	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  nil,
	}
	
	// We test the handler directly since PocketBase router testing is complex
	requestInfo := &core.RequestInfo{}
	resp, err := server.handleMCPRequest(requestInfo, mcpReq)
	
	require.NoError(t, err)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)
	
	// Verify response can be marshaled
	_, err = json.Marshal(resp)
	assert.NoError(t, err)
}


// TestMCPErrorHandling tests error handling in MCP requests
func TestMCPErrorHandling(t *testing.T) {
	server := NewTinyBrainPocketBaseServer()
	
	// Test invalid method
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "invalid_method",
		Params:  nil,
	}
	
	requestInfo := &core.RequestInfo{}
	resp, err := server.handleMCPRequest(requestInfo, req)
	
	require.NoError(t, err)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32601, resp.Error.Code)
	assert.Equal(t, "Method not found", resp.Error.Message)
}

// TestMCPRequestResponseStructures tests the MCP request/response structures
func TestMCPRequestResponseStructures(t *testing.T) {
	// Test MCPRequest JSON marshaling
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
		Params:  map[string]interface{}{"key": "value"},
	}
	
	reqJSON, err := json.Marshal(req)
	require.NoError(t, err)
	
	var unmarshaledReq MCPRequest
	err = json.Unmarshal(reqJSON, &unmarshaledReq)
	require.NoError(t, err)
	
	assert.Equal(t, req.JSONRPC, unmarshaledReq.JSONRPC)
	assert.Equal(t, req.ID, unmarshaledReq.ID)
	assert.Equal(t, req.Method, unmarshaledReq.Method)
	
	// Test MCPResponse JSON marshaling
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]interface{}{"status": "ok"},
	}
	
	respJSON, err := json.Marshal(resp)
	require.NoError(t, err)
	
	var unmarshaledResp MCPResponse
	err = json.Unmarshal(respJSON, &unmarshaledResp)
	require.NoError(t, err)
	
	assert.Equal(t, resp.JSONRPC, unmarshaledResp.JSONRPC)
	assert.Equal(t, resp.ID, unmarshaledResp.ID)
	
	// Test MCPError JSON marshaling
	mcpError := MCPError{
		Code:    -32602,
		Message: "Invalid params",
	}
	
	errorJSON, err := json.Marshal(mcpError)
	require.NoError(t, err)
	
	var unmarshaledError MCPError
	err = json.Unmarshal(errorJSON, &unmarshaledError)
	require.NoError(t, err)
	
	assert.Equal(t, mcpError.Code, unmarshaledError.Code)
	assert.Equal(t, mcpError.Message, unmarshaledError.Message)
}

