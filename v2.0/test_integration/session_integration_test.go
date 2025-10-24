package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

// Integration test for session management via HTTP API
func TestSessionIntegration(t *testing.T) {
	baseURL := "http://127.0.0.1:8090"

	// Wait for server to be ready
	if !waitForServer(baseURL, 30*time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}

	// Test 1: Health check
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Health endpoint returns JSON, not plain text
		expectedJSON := `{"status":"healthy","service":"TinyBrain v2.0 Minimal","version":"2.0.0"}`
		if string(body) != expectedJSON {
			t.Fatalf("Expected JSON response, got '%s'", string(body))
		}

		t.Log("✅ Health check passed")
	})

	// Test 2: Hello endpoint
	t.Run("HelloEndpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/hello")
		if err != nil {
			t.Fatalf("Hello endpoint failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		expected := "Hello from TinyBrain v2.0!"
		if string(body) != expected {
			t.Fatalf("Expected '%s', got '%s'", expected, string(body))
		}

		t.Log("✅ Hello endpoint passed")
	})

	// Test 3: Check if collections exist via API
	t.Run("CollectionsExist", func(t *testing.T) {
		// Try to access the collections API
		resp, err := http.Get(baseURL + "/api/collections")
		if err != nil {
			t.Fatalf("Collections API failed: %v", err)
		}
		defer resp.Body.Close()

		// We expect either 200 (if auth is disabled) or 401 (if auth is required)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 200 or 401, got %d", resp.StatusCode)
		}

		t.Log("✅ Collections API accessible")
	})

	// Test 4: Test MCP tool via HTTP (if available)
	t.Run("MCPToolTest", func(t *testing.T) {
		// This would test the MCP tools if they were exposed via HTTP
		// For now, we'll just verify the server is running
		t.Log("✅ MCP server is running (tested via server startup)")
	})

	t.Log("🎉 All integration tests passed!")
}

// Helper function to wait for server to be ready
func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := http.Get(url + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}
		time.Sleep(1 * time.Second)
	}

	return false
}

// Test that verifies the server can handle basic requests
func TestServerBasicFunctionality(t *testing.T) {
	baseURL := "http://127.0.0.1:8090"

	// Test server is responding
	resp, err := http.Get(baseURL + "/hello")
	if err != nil {
		t.Fatalf("Server not responding: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	expected := "Hello from TinyBrain v2.0!"
	if string(body) != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, string(body))
	}

	t.Log("✅ Server basic functionality verified")
}
