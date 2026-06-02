package main

import (
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/mcp"
	"github.com/rainmana/tinybrain/internal/repository"
	"github.com/rainmana/tinybrain/internal/services"
	"github.com/stretchr/testify/require"
)

func TestMCPUserFlowDocumentedCoreAliases(t *testing.T) {
	server, db := newTestMCPServer(t)
	defer db.Close()

	listResp := mcpRoundTrip(t, server, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	})
	require.Nil(t, listResp.Error)
	require.Contains(t, toolNames(t, listResp), "store_memory")
	require.Contains(t, toolNames(t, listResp), "update_memory")
	require.Contains(t, toolNames(t, listResp), "delete_memory")

	session := callToolMap(t, server, "mcp_tinybrain-mcp-server_create_session", map[string]interface{}{
		"name":              "OSINT Intelligence Gathering",
		"description":       "Authorized enterprise assessment",
		"task_type":         "intelligence_analysis",
		"intelligence_type": "osint",
		"classification":    "unclassified",
		"threat_level":      "medium",
	})
	sessionID := requireString(t, session, "id")
	require.Equal(t, "intelligence_analysis", session["task_type"])
	require.Equal(t, "osint", session["metadata"].(map[string]interface{})["intelligence_type"])

	intelMemory := callToolMap(t, server, "mcp_tinybrain-mcp-server_create_memory", map[string]interface{}{
		"session_id":        sessionID,
		"title":             "Social Media Intelligence Finding",
		"content":           "Social media analysis reveals suspicious activity around the target brand.",
		"category":          "intelligence",
		"priority":          7,
		"confidence":        0.8,
		"tags":              []string{"osint", "social-media", "reconnaissance"},
		"source":            "manual-analysis",
		"intelligence_type": "osint",
		"threat_level":      "medium",
		"mitre_tactic":      "TA0043",
		"mitre_technique":   "T1591",
	})
	intelMemoryID := requireString(t, intelMemory, "id")
	require.Equal(t, "intelligence", intelMemory["category"])
	require.Contains(t, stringSlice(t, intelMemory["tags"]), "mitre_technique:T1591")

	vulnMemory := callToolMap(t, server, "store_memory", map[string]interface{}{
		"session_id": sessionID,
		"title":      "SQL Injection in Login",
		"content":    "The username parameter is injectable and can bypass authentication.",
		"category":   "vulnerability",
		"priority":   9,
		"confidence": 0.92,
		"tags":       []interface{}{"sql-injection", "authentication"},
		"source":     "manual-testing",
	})
	vulnMemoryID := requireString(t, vulnMemory, "id")

	searchResults := callToolSlice(t, server, "mcp_tinybrain-mcp-server_search_memory", map[string]interface{}{
		"query":      "social media",
		"session_id": sessionID,
		"category":   "intelligence",
		"tags":       []string{"osint"},
		"limit":      5,
	})
	require.Len(t, searchResults, 1)
	require.Equal(t, "Social Media Intelligence Finding", searchResults[0].(map[string]interface{})["memory_entry"].(map[string]interface{})["title"])

	updated := callToolMap(t, server, "update_memory", map[string]interface{}{
		"memory_id": intelMemoryID,
		"priority":  8,
		"tags":      []string{"osint", "social-media", "validated"},
	})
	require.Equal(t, float64(8), updated["priority"])
	require.Contains(t, stringSlice(t, updated["tags"]), "validated")

	relationship := callToolMap(t, server, "create_relationship", map[string]interface{}{
		"source_memory_id":  intelMemoryID,
		"target_memory_id":  vulnMemoryID,
		"relationship_type": "references",
		"strength":          0.8,
		"description":       "Intelligence finding informs the exploitability hypothesis.",
	})
	require.Equal(t, "references", relationship["relationship_type"])

	related := callToolSlice(t, server, "mcp_tinybrain-mcp-server_get_related_entries", map[string]interface{}{
		"memory_id":         intelMemoryID,
		"relationship_type": "references",
		"limit":             5,
	})
	require.Len(t, related, 1)
	require.Equal(t, vulnMemoryID, related[0].(map[string]interface{})["id"])

	snapshot := callToolMap(t, server, "create_context_snapshot", map[string]interface{}{
		"session_id": sessionID,
		"name":       "Initial Assessment State",
		"context_data": map[string]interface{}{
			"current_stage": "reconnaissance",
			"next_steps":    []string{"validate login finding", "preserve evidence"},
		},
	})
	require.NotEmpty(t, snapshot["id"])
	require.Contains(t, snapshot["memory_summary"], "SQL Injection")

	task := callToolMap(t, server, "create_task_progress", map[string]interface{}{
		"session_id":          sessionID,
		"task_name":           "Validate Authentication Finding",
		"stage":               "triage",
		"status":              "in_progress",
		"progress_percentage": 25,
		"notes":               "Evidence captured, validation pending.",
	})
	require.Equal(t, "in_progress", task["status"])

	task = callToolMap(t, server, "update_task_progress", map[string]interface{}{
		"session_id":          sessionID,
		"task_name":           "Validate Authentication Finding",
		"stage":               "validation",
		"status":              "completed",
		"progress_percentage": 100,
		"notes":               "Validated in authorized test environment.",
	})
	require.Equal(t, "completed", task["status"])

	deleted := callToolMap(t, server, "delete_memory", map[string]interface{}{"memory_id": vulnMemoryID})
	require.Equal(t, true, deleted["deleted"])
}

func TestRunHealthCheck(t *testing.T) {
	logger := log.New(io.Discard)
	require.NoError(t, runHealthCheck(filepath.Join(t.TempDir(), "health.db"), logger))
}

func TestMCPNotificationsAreSilent(t *testing.T) {
	server, db := newTestMCPServer(t)
	defer db.Close()

	resp := mcpRoundTrip(t, server, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
		"params":  map[string]interface{}{},
	})
	require.Nil(t, resp)
}

func newTestMCPServer(t *testing.T) (*mcp.Server, *database.Database) {
	t.Helper()

	logger := log.New(io.Discard)
	db, err := database.NewDatabase(filepath.Join(t.TempDir(), "tinybrain.db"), logger)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	securityRepo := repository.NewSecurityRepository(db, logger)
	downloader := services.NewSecurityDataDownloader(logger)
	tinyBrain := &TinyBrainServer{
		db:                 db,
		repo:               repository.NewMemoryRepository(db.GetDB(), logger),
		securityRepo:       securityRepo,
		securityDownloader: downloader,
		securityRetrieval:  services.NewSecurityRetrievalService(securityRepo, logger),
		securityUpdate:     services.NewSecurityUpdateService(downloader, securityRepo, logger),
		logger:             logger,
		dbPath:             filepath.Join(t.TempDir(), "tinybrain.db"),
	}

	server := mcp.NewServer("TinyBrain Memory Storage", "test", "test server", logger)
	tinyBrain.registerTools(server)
	return server, db
}

func mcpRoundTrip(t *testing.T, server *mcp.Server, request map[string]interface{}) *mcp.MCPResponse {
	t.Helper()

	data, err := json.Marshal(request)
	require.NoError(t, err)

	var req mcp.MCPRequest
	require.NoError(t, json.Unmarshal(data, &req))
	return server.HandleRequest(context.Background(), &req)
}

func callToolMap(t *testing.T, server *mcp.Server, name string, arguments map[string]interface{}) map[string]interface{} {
	t.Helper()

	result := callTool(t, server, name, arguments)
	resultMap, ok := result.(map[string]interface{})
	require.Truef(t, ok, "expected map result for %s, got %T", name, result)
	return resultMap
}

func callToolSlice(t *testing.T, server *mcp.Server, name string, arguments map[string]interface{}) []interface{} {
	t.Helper()

	result := callTool(t, server, name, arguments)
	resultSlice, ok := result.([]interface{})
	require.Truef(t, ok, "expected slice result for %s, got %T", name, result)
	return resultSlice
}

func callTool(t *testing.T, server *mcp.Server, name string, arguments map[string]interface{}) interface{} {
	t.Helper()

	resp := mcpRoundTrip(t, server, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      name,
			"arguments": arguments,
		},
	})
	require.Nilf(t, resp.Error, "tool %s returned error: %#v", name, resp.Error)

	var wrapper struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	data, err := json.Marshal(resp.Result)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &wrapper))
	require.NotEmpty(t, wrapper.Content)

	var result interface{}
	require.NoError(t, json.Unmarshal([]byte(wrapper.Content[0].Text), &result))
	return result
}

func toolNames(t *testing.T, resp *mcp.MCPResponse) []string {
	t.Helper()

	var result struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	data, err := json.Marshal(resp.Result)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &result))

	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
	}
	return names
}

func requireString(t *testing.T, values map[string]interface{}, key string) string {
	t.Helper()

	value, ok := values[key].(string)
	require.Truef(t, ok, "expected %s to be string, got %T", key, values[key])
	require.NotEmpty(t, value)
	return value
}

func stringSlice(t *testing.T, value interface{}) []string {
	t.Helper()

	raw, ok := value.([]interface{})
	require.Truef(t, ok, "expected string slice, got %T", value)
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		str, ok := item.(string)
		require.Truef(t, ok, "expected string item, got %T", item)
		result = append(result, str)
	}
	return result
}
