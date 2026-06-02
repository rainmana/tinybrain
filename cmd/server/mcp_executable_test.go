package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMCPExecutableAuthorizedAssessmentFlow(t *testing.T) {
	binaryPath := buildTinyBrainExecutable(t)
	dbPath := filepath.Join(t.TempDir(), "memory.db")

	client := startMCPExecutable(t, binaryPath, dbPath)
	initialize := client.request(t, "initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "tinybrain-user-flow-test",
			"version": "1.0.0",
		},
	})
	initResult := decodeResultMap(t, initialize)
	require.Equal(t, "2024-11-05", initResult["protocolVersion"])
	require.Equal(t, "TinyBrain Memory Storage", initResult["serverInfo"].(map[string]interface{})["name"])

	client.notify(t, "notifications/initialized", map[string]interface{}{})

	tools := toolNamesFromExecutableResponse(t, client.request(t, "tools/list", map[string]interface{}{}))
	for _, name := range []string{
		"create_session",
		"store_memory",
		"search_memories",
		"update_memory",
		"get_memory",
		"create_relationship",
		"get_related_memories",
		"create_context_snapshot",
		"create_task_progress",
		"update_task_progress",
		"export_session_data",
		"delete_memory",
	} {
		require.Contains(t, tools, name)
	}

	health := client.callToolMap(t, "health_check", map[string]interface{}{})
	require.Equal(t, "healthy", health["status"])
	require.Equal(t, dbPath, health["db_path"])

	session := client.callToolMap(t, "create_session", map[string]interface{}{
		"name":        "Authorized Web Application Assessment",
		"description": "Enterprise red-team validation in an owned lab environment",
		"task_type":   "penetration_test",
		"metadata":    `{"engagement_id":"ENG-2026-001","scope":"owned lab"}`,
	})
	sessionID := requireString(t, session, "id")
	require.Equal(t, "ENG-2026-001", session["metadata"].(map[string]interface{})["engagement_id"])

	evidence := client.callToolMap(t, "store_memory", map[string]interface{}{
		"session_id":   sessionID,
		"title":        "Login Request Evidence",
		"content":      "Captured the baseline login request and response pair for authorized testing.",
		"category":     "evidence",
		"priority":     6,
		"confidence":   0.95,
		"tags":         `["burp","auth","evidence"]`,
		"source":       "user-flow-test",
		"content_type": "text",
	})
	evidenceID := requireString(t, evidence, "id")
	require.Contains(t, stringSlice(t, evidence["tags"]), "burp")

	vulnerability := client.callToolMap(t, "store_memory", map[string]interface{}{
		"session_id":      sessionID,
		"title":           "SQL Injection in Login",
		"content":         "The username field accepts a tautology payload and bypasses authentication in the lab target.",
		"category":        "vulnerability",
		"priority":        9,
		"confidence":      0.91,
		"tags":            []interface{}{"sql-injection", "authentication", "critical-path"},
		"source":          "manual-validation",
		"mitre_tactic":    "TA0001",
		"mitre_technique": "T1190",
	})
	vulnerabilityID := requireString(t, vulnerability, "id")
	require.Contains(t, stringSlice(t, vulnerability["tags"]), "mitre_technique:T1190")

	searchResults := client.callToolSlice(t, "search_memories", map[string]interface{}{
		"query":          "login authentication",
		"session_id":     sessionID,
		"categories":     `["vulnerability"]`,
		"tags":           []interface{}{"sql-injection"},
		"min_priority":   8,
		"min_confidence": 0.8,
		"limit":          5,
	})
	require.Len(t, searchResults, 1)
	require.Equal(t, "SQL Injection in Login", searchResults[0].(map[string]interface{})["memory_entry"].(map[string]interface{})["title"])

	updated := client.callToolMap(t, "update_memory", map[string]interface{}{
		"memory_id":  vulnerabilityID,
		"priority":   10,
		"confidence": 0.97,
		"content":    "The username field accepts a tautology payload and bypasses authentication in the lab target; retest confirmed.",
		"tags":       `["sql-injection","authentication","validated"]`,
	})
	require.Equal(t, float64(10), updated["priority"])
	require.Contains(t, stringSlice(t, updated["tags"]), "validated")

	reloaded := client.callToolMap(t, "get_memory", map[string]interface{}{"memory_id": vulnerabilityID})
	require.Equal(t, "SQL Injection in Login", reloaded["title"])
	require.Equal(t, float64(10), reloaded["priority"])

	relationship := client.callToolMap(t, "create_relationship", map[string]interface{}{
		"source_memory_id":  evidenceID,
		"target_memory_id":  vulnerabilityID,
		"relationship_type": "references",
		"strength":          0.9,
		"description":       "Captured evidence supports the authentication bypass finding.",
	})
	require.Equal(t, "references", relationship["relationship_type"])

	related := client.callToolSlice(t, "get_related_memories", map[string]interface{}{
		"memory_id":         evidenceID,
		"relationship_type": "references",
		"limit":             5,
	})
	require.Len(t, related, 1)
	require.Equal(t, vulnerabilityID, related[0].(map[string]interface{})["id"])

	snapshot := client.callToolMap(t, "create_context_snapshot", map[string]interface{}{
		"session_id":   sessionID,
		"name":         "Post-validation checkpoint",
		"description":  "State after confirming the authentication finding",
		"context_data": `{"phase":"validation","next_steps":["write remediation","collect owner signoff"]}`,
	})
	require.NotEmpty(t, snapshot["id"])
	require.Contains(t, snapshot["memory_summary"], "SQL Injection")

	task := client.callToolMap(t, "create_task_progress", map[string]interface{}{
		"session_id":          sessionID,
		"task_name":           "Validate login authentication controls",
		"stage":               "validation",
		"status":              "in_progress",
		"progress_percentage": 60,
		"notes":               "Reproduction complete; remediation notes pending.",
	})
	require.Equal(t, "in_progress", task["status"])

	task = client.callToolMap(t, "update_task_progress", map[string]interface{}{
		"session_id":          sessionID,
		"task_name":           "Validate login authentication controls",
		"stage":               "reporting",
		"status":              "completed",
		"progress_percentage": 100,
		"notes":               "Finding validated and ready for report.",
	})
	require.Equal(t, "completed", task["status"])

	tasks := client.callToolSlice(t, "list_task_progress", map[string]interface{}{
		"session_id": sessionID,
		"status":     "completed",
	})
	require.Len(t, tasks, 1)
	require.Equal(t, "Validate login authentication controls", tasks[0].(map[string]interface{})["task_name"])

	contextSummary := client.callToolMap(t, "get_context_summary", map[string]interface{}{
		"session_id":   sessionID,
		"current_task": "login authentication tautology payload",
		"max_memories": 10,
	})
	require.Equal(t, sessionID, contextSummary["session_id"])
	require.NotEmpty(t, contextSummary["relevant_memories"])

	exported := client.callToolMap(t, "export_session_data", map[string]interface{}{"session_id": sessionID})
	require.Equal(t, sessionID, exported["session"].(map[string]interface{})["id"])
	require.Len(t, exported["memory_entries"], 2)
	require.Len(t, exported["relationships"], 1)
	require.Len(t, exported["snapshots"], 1)
	require.Len(t, exported["tasks"], 1)

	stats := client.callToolMap(t, "get_database_stats", map[string]interface{}{})
	require.Equal(t, float64(1), stats["sessions_count"])
	require.Equal(t, float64(2), stats["memory_entries_count"])
	require.Equal(t, float64(1), stats["relationships_count"])

	client.stop(t)

	restarted := startMCPExecutable(t, binaryPath, dbPath)
	restarted.request(t, "initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "tinybrain-user-flow-test-restart",
			"version": "1.0.0",
		},
	})
	restarted.notify(t, "notifications/initialized", map[string]interface{}{})

	persistedSession := restarted.callToolMap(t, "get_session", map[string]interface{}{"session_id": sessionID})
	require.Equal(t, "Authorized Web Application Assessment", persistedSession["name"])

	persistedSearch := restarted.callToolSlice(t, "search_memories", map[string]interface{}{
		"query":      "tautology payload",
		"session_id": sessionID,
		"category":   "vulnerability",
		"tags":       `["validated"]`,
		"limit":      5,
	})
	require.Len(t, persistedSearch, 1)
	require.Equal(t, vulnerabilityID, persistedSearch[0].(map[string]interface{})["memory_entry"].(map[string]interface{})["id"])

	deleted := restarted.callToolMap(t, "delete_memory", map[string]interface{}{"memory_id": evidenceID})
	require.True(t, deleted["deleted"].(bool))

	finalStats := restarted.callToolMap(t, "get_database_stats", map[string]interface{}{})
	require.Equal(t, float64(1), finalStats["memory_entries_count"])
}

type mcpExecutableClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr *bytes.Buffer
	nextID int
}

type executableRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	} `json:"error,omitempty"`
}

func buildTinyBrainExecutable(t *testing.T) string {
	t.Helper()

	binaryName := "tinybrain-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(t.TempDir(), binaryName)
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "failed to build test executable: %s", string(output))
	return binaryPath
}

func startMCPExecutable(t *testing.T, binaryPath, dbPath string) *mcpExecutableClient {
	t.Helper()

	cmd := exec.Command(binaryPath)
	cmd.Dir = filepath.Dir(binaryPath)
	cmd.Env = append(os.Environ(), "TINYBRAIN_DB_PATH="+dbPath)

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdoutPipe, err := cmd.StdoutPipe()
	require.NoError(t, err)

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	require.NoError(t, cmd.Start())

	client := &mcpExecutableClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
		stderr: stderr,
		nextID: 1,
	}
	t.Cleanup(func() {
		client.stop(t)
	})

	return client
}

func (c *mcpExecutableClient) request(t *testing.T, method string, params map[string]interface{}) executableRPCResponse {
	t.Helper()

	id := c.nextID
	c.nextID++
	c.write(t, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	})

	resp := c.readResponse(t)
	require.Equalf(t, id, resp.ID, "unexpected response id for method %s", method)
	if resp.Error != nil {
		require.Failf(t, "unexpected JSON-RPC error", "method %s returned %d %s: %v", method, resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}
	return resp
}

func (c *mcpExecutableClient) notify(t *testing.T, method string, params map[string]interface{}) {
	t.Helper()

	c.write(t, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	})
}

func (c *mcpExecutableClient) callToolMap(t *testing.T, name string, arguments map[string]interface{}) map[string]interface{} {
	t.Helper()

	result := c.callTool(t, name, arguments)
	resultMap, ok := result.(map[string]interface{})
	require.Truef(t, ok, "expected map result for tool %s, got %T", name, result)
	return resultMap
}

func (c *mcpExecutableClient) callToolSlice(t *testing.T, name string, arguments map[string]interface{}) []interface{} {
	t.Helper()

	result := c.callTool(t, name, arguments)
	resultSlice, ok := result.([]interface{})
	require.Truef(t, ok, "expected slice result for tool %s, got %T", name, result)
	return resultSlice
}

func (c *mcpExecutableClient) callTool(t *testing.T, name string, arguments map[string]interface{}) interface{} {
	t.Helper()

	resp := c.request(t, "tools/call", map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	})

	var wrapper struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	require.NoError(t, json.Unmarshal(resp.Result, &wrapper))
	require.Len(t, wrapper.Content, 1)
	require.Equal(t, "text", wrapper.Content[0].Type)

	var result interface{}
	require.NoError(t, json.Unmarshal([]byte(wrapper.Content[0].Text), &result))
	return result
}

func (c *mcpExecutableClient) write(t *testing.T, payload map[string]interface{}) {
	t.Helper()

	data, err := json.Marshal(payload)
	require.NoError(t, err)
	_, err = fmt.Fprintln(c.stdin, string(data))
	require.NoError(t, err)
}

func (c *mcpExecutableClient) readResponse(t *testing.T) executableRPCResponse {
	t.Helper()

	type readResult struct {
		line string
		err  error
	}
	done := make(chan readResult, 1)
	go func() {
		line, err := c.stdout.ReadString('\n')
		done <- readResult{line: line, err: err}
	}()

	select {
	case result := <-done:
		require.NoErrorf(t, result.err, "server stderr:\n%s", c.stderr.String())
		var resp executableRPCResponse
		require.NoErrorf(t, json.Unmarshal([]byte(result.line), &resp), "raw response: %s\nserver stderr:\n%s", result.line, c.stderr.String())
		return resp
	case <-time.After(5 * time.Second):
		require.Failf(t, "timed out waiting for MCP response", "server stderr:\n%s", c.stderr.String())
	}

	return executableRPCResponse{}
}

func (c *mcpExecutableClient) stop(t *testing.T) {
	t.Helper()

	if c.cmd == nil || c.cmd.Process == nil {
		return
	}

	_ = c.stdin.Close()
	waitDone := make(chan error, 1)
	go func() {
		waitDone <- c.cmd.Wait()
	}()

	select {
	case err := <-waitDone:
		if err != nil && !strings.Contains(err.Error(), "file already closed") {
			require.NoErrorf(t, err, "server stderr:\n%s", c.stderr.String())
		}
	case <-time.After(2 * time.Second):
		_ = c.cmd.Process.Kill()
		<-waitDone
	}

	c.cmd.Process = nil
}

func decodeResultMap(t *testing.T, resp executableRPCResponse) map[string]interface{} {
	t.Helper()

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(resp.Result, &result))
	return result
}

func toolNamesFromExecutableResponse(t *testing.T, resp executableRPCResponse) []string {
	t.Helper()

	var result struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	require.NoError(t, json.Unmarshal(resp.Result, &result))

	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
	}
	return names
}
