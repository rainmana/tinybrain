package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// SimpleIntegrationTestSuite tests the TinyBrain v2.0 server via HTTP endpoints
type SimpleIntegrationTestSuite struct {
	suite.Suite
	baseURL string
	client  *http.Client
}

// SetupSuite initializes the test suite
func (suite *SimpleIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://127.0.0.1:8090"
	suite.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	// Wait for server to be ready
	suite.waitForServer()
}

// waitForServer waits for the server to be ready
func (suite *SimpleIntegrationTestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	suite.Fail("Server not ready after 30 seconds")
}

// TestHealthEndpoint tests the health endpoint
func (suite *SimpleIntegrationTestSuite) TestHealthEndpoint() {
	resp, err := suite.client.Get(suite.baseURL + "/health")
	suite.Require().NoError(err, "Failed to get health endpoint")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "Health endpoint should return 200")

	var healthResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	suite.NoError(err, "Failed to decode health response")

	suite.Equal("healthy", healthResponse["status"], "Status should be healthy")
	suite.Equal("TinyBrain v2.0 Complete", healthResponse["service"], "Service should be TinyBrain v2.0 Complete")
	suite.Equal("2.0.0", healthResponse["version"], "Version should be 2.0.0")

	// Verify features
	features, ok := healthResponse["features"].([]interface{})
	suite.True(ok, "Features should be an array")
	suite.Contains(features, "session_management", "Should include session_management")
	suite.Contains(features, "memory_storage", "Should include memory_storage")
	suite.Contains(features, "relationship_tracking", "Should include relationship_tracking")
	suite.Contains(features, "context_snapshots", "Should include context_snapshots")
	suite.Contains(features, "task_progress", "Should include task_progress")
	suite.Contains(features, "pocketbase_database", "Should include pocketbase_database")
	suite.Contains(features, "mcp_protocol", "Should include mcp_protocol")
}

// TestHelloEndpoint tests the hello endpoint
func (suite *SimpleIntegrationTestSuite) TestHelloEndpoint() {
	resp, err := suite.client.Get(suite.baseURL + "/hello")
	suite.Require().NoError(err, "Failed to get hello endpoint")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "Hello endpoint should return 200")

	// Read response body
	body := make([]byte, 200)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		suite.NoError(err, "Failed to read response body")
	}

	responseText := string(body[:n])
	suite.Contains(responseText, "TinyBrain v2.0 Complete", "Hello response should contain TinyBrain v2.0 Complete")
}

// TestPocketBaseAdmin tests that PocketBase admin is accessible
func (suite *SimpleIntegrationTestSuite) TestPocketBaseAdmin() {
	resp, err := suite.client.Get(suite.baseURL + "/_/")
	suite.Require().NoError(err, "Failed to get PocketBase admin")
	defer resp.Body.Close()

	// PocketBase admin might return 400 for root path, but should not return 404
	suite.NotEqual(http.StatusNotFound, resp.StatusCode, "PocketBase admin should not return 404")
}

// TestPocketBaseAPI tests that PocketBase API is accessible
func (suite *SimpleIntegrationTestSuite) TestPocketBaseAPI() {
	resp, err := suite.client.Get(suite.baseURL + "/api/")
	suite.Require().NoError(err, "Failed to get PocketBase API")
	defer resp.Body.Close()

	// PocketBase API might return 404 for root path, but should not return 500
	suite.NotEqual(http.StatusInternalServerError, resp.StatusCode, "PocketBase API should not return 500")
}

// TestServerCapabilities tests server capabilities
func (suite *SimpleIntegrationTestSuite) TestServerCapabilities() {
	// Test that the server is running and responding
	resp, err := suite.client.Get(suite.baseURL + "/health")
	suite.Require().NoError(err, "Server should be responding")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "Server should return 200 OK")

	// Test that we can make multiple requests
	for i := 0; i < 5; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		suite.Require().NoError(err, "Server should handle multiple requests")
		resp.Body.Close()
		suite.Equal(http.StatusOK, resp.StatusCode, "Server should consistently return 200 OK")
	}
}

// TestResponseHeaders tests response headers
func (suite *SimpleIntegrationTestSuite) TestResponseHeaders() {
	resp, err := suite.client.Get(suite.baseURL + "/health")
	suite.Require().NoError(err, "Failed to get health endpoint")
	defer resp.Body.Close()

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	suite.Contains(contentType, "application/json", "Health endpoint should return JSON")

	// Check that server is responding quickly
	suite.Less(resp.Header.Get("X-Response-Time"), "1s", "Server should respond quickly")
}

// TestServerStability tests server stability
func (suite *SimpleIntegrationTestSuite) TestServerStability() {
	// Make multiple concurrent requests
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer func() { done <- true }()

			for j := 0; j < 5; j++ {
				resp, err := suite.client.Get(suite.baseURL + "/health")
				if err != nil {
					suite.Fail("Concurrent request failed: " + err.Error())
					return
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					suite.Fail(fmt.Sprintf("Concurrent request returned %d", resp.StatusCode))
					return
				}

				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		select {
		case <-done:
			// Good
		case <-time.After(30 * time.Second):
			suite.Fail("Concurrent requests timed out")
		}
	}
}

// TestMCPToolsAvailability tests that MCP tools are available
func (suite *SimpleIntegrationTestSuite) TestMCPToolsAvailability() {
	// This test verifies that the server is running and ready for MCP connections
	// In a real scenario, you would connect via STDIO and test the MCP tools

	resp, err := suite.client.Get(suite.baseURL + "/health")
	suite.Require().NoError(err, "Failed to get health endpoint")
	defer resp.Body.Close()

	var healthResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	suite.NoError(err, "Failed to decode health response")

	// Verify MCP protocol is listed as a feature
	features, ok := healthResponse["features"].([]interface{})
	suite.True(ok, "Features should be an array")
	suite.Contains(features, "mcp_protocol", "Should include mcp_protocol feature")

	// Verify all expected features are present
	expectedFeatures := []string{
		"session_management",
		"memory_storage",
		"relationship_tracking",
		"context_snapshots",
		"task_progress",
		"pocketbase_database",
		"mcp_protocol",
	}

	for _, feature := range expectedFeatures {
		suite.Contains(features, feature, "Should include feature: %s", feature)
	}
}

// TestServerPerformance tests basic server performance
func (suite *SimpleIntegrationTestSuite) TestServerPerformance() {
	start := time.Now()

	// Make 10 requests and measure average response time
	requestCount := 10
	for i := 0; i < requestCount; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		suite.Require().NoError(err, "Performance test request failed")
		resp.Body.Close()
		suite.Equal(http.StatusOK, resp.StatusCode, "Performance test should return 200")
	}

	elapsed := time.Since(start)
	averageResponseTime := elapsed / time.Duration(requestCount)

	// Server should respond within 1 second on average
	suite.Less(averageResponseTime, 1*time.Second, "Average response time should be less than 1 second")

	fmt.Printf("Average response time: %v\n", averageResponseTime)
}

// Run the test suite
func TestSimpleIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleIntegrationTestSuite))
}
