package test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

// MCPIntegrationTestSuite tests the complete MCP integration
type MCPIntegrationTestSuite struct {
	suite.Suite
	app                 *pocketbase.PocketBase
	dataDir             string
	sessionService      *services.SessionServiceV2
	memoryService       *services.MemoryServiceV2
	relationshipService *services.RelationshipServiceV2
	contextService      *services.ContextServiceV2
	taskService         *services.TaskServiceV2
}

// SetupSuite initializes the test suite with proper database setup
func (suite *MCPIntegrationTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_mcp_pb_data"
	os.RemoveAll(suite.dataDir) // Clean up any existing test data

	// Initialize PocketBase with test config
	config := pocketbase.Config{
		DefaultDataDir: suite.dataDir,
	}
	suite.app = pocketbase.NewWithConfig(config)

	// Bootstrap the app
	err := suite.app.Bootstrap()
	suite.Require().NoError(err, "Failed to bootstrap PocketBase for testing")

	// Initialize database collections
	suite.initializeCollections()

	// Initialize repositories and services
	sessionRepo := repository.NewSessionRepositoryV2(suite.app)
	memoryRepo := repository.NewMemoryRepositoryV2(suite.app)
	relationshipRepo := repository.NewRelationshipRepositoryV2(suite.app)
	contextRepo := repository.NewContextRepositoryV2(suite.app)
	taskRepo := repository.NewTaskRepositoryV2(suite.app)

	suite.sessionService = services.NewSessionServiceV2(sessionRepo)
	suite.memoryService = services.NewMemoryServiceV2(memoryRepo)
	suite.relationshipService = services.NewRelationshipServiceV2(relationshipRepo)
	suite.contextService = services.NewContextServiceV2(contextRepo)
	suite.taskService = services.NewTaskServiceV2(taskRepo)
}

// TearDownSuite cleans up after the test suite
func (suite *MCPIntegrationTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// initializeCollections creates all required database collections
func (suite *MCPIntegrationTestSuite) initializeCollections() {
	collections := []*models.Collection{
		database.CreateSessionsCollection(),
		database.CreateMemoryEntriesCollection(),
		database.CreateRelationshipsCollection(),
		database.CreateContextSnapshotsCollection(),
		database.CreateTaskProgressCollection(),
	}

	for _, collection := range collections {
		// Check if collection already exists
		existing, err := suite.app.Dao().FindCollectionByNameOrId(collection.Name)
		if err != nil {
			// Collection doesn't exist, create it
			err := suite.app.Dao().SaveCollection(collection)
			suite.Require().NoError(err, "Failed to create collection %s", collection.Name)
		} else {
			suite.T().Logf("Collection '%s' already exists", existing.Name)
		}
	}
}

// TestSessionManagement tests the complete session management workflow
func (suite *MCPIntegrationTestSuite) TestSessionManagement() {
	ctx := context.Background()

	// Test creating a session
	sessionReq := &models.SessionCreateRequest{
		Name:        "Security Assessment Session",
		TaskType:    "penetration_test",
		Description: "Comprehensive security assessment",
		Metadata: map[string]interface{}{
			"priority": "high",
			"client":   "test-client",
		},
	}

	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.NoError(err, "Failed to create session")
	suite.NotNil(session, "Session should not be nil")
	suite.Equal("Security Assessment Session", session.Name)
	suite.Equal("penetration_test", session.TaskType)
	suite.Equal("active", session.Status)

	// Test retrieving the session
	retrievedSession, err := suite.sessionService.GetSession(ctx, session.ID)
	suite.NoError(err, "Failed to retrieve session")
	suite.Equal(session.ID, retrievedSession.ID)
	suite.Equal(session.Name, retrievedSession.Name)

	// Test listing sessions
	listReq := &models.SessionListRequest{
		Limit:  10,
		Offset: 0,
	}
	sessions, totalCount, err := suite.sessionService.ListSessions(ctx, listReq)
	suite.NoError(err, "Failed to list sessions")
	suite.GreaterOrEqual(totalCount, 1, "Should have at least one session")
	suite.NotEmpty(sessions, "Session list should not be empty")

	// Test updating the session
	updateReq := &models.SessionUpdateRequest{
		Name:        "Updated Security Assessment",
		Status:      "in_progress",
		Description: "Updated description",
	}
	updatedSession, err := suite.sessionService.UpdateSession(ctx, session.ID, updateReq)
	suite.NoError(err, "Failed to update session")
	suite.Equal("Updated Security Assessment", updatedSession.Name)
	suite.Equal("in_progress", updatedSession.Status)
}

// TestMemoryManagement tests the complete memory management workflow
func (suite *MCPIntegrationTestSuite) TestMemoryManagement() {
	ctx := context.Background()

	// First create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Memory Test Session",
		TaskType: "vulnerability_assessment",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for memory test")

	// Test storing a memory
	memoryReq := &models.MemoryCreateRequest{
		SessionID:   session.ID,
		Title:       "SQL Injection Vulnerability",
		Content:     "Found SQL injection vulnerability in login form",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.9,
		Tags:        []string{"sql-injection", "authentication", "critical"},
		Source:      "manual_testing",
		ContentType: "text",
	}

	memory, err := suite.memoryService.StoreMemory(ctx, memoryReq)
	suite.NoError(err, "Failed to store memory")
	suite.NotNil(memory, "Memory should not be nil")
	suite.Equal("SQL Injection Vulnerability", memory.Title)
	suite.Equal("vulnerability", memory.Category)
	suite.Equal(8, memory.Priority)
	suite.Equal(float32(0.9), memory.Confidence)

	// Test searching memories
	searchReq := &models.MemorySearchRequest{
		SessionID: session.ID,
		Query:     "SQL injection",
		Category:  "vulnerability",
		Limit:     10,
		Offset:    0,
	}
	memories, total, err := suite.memoryService.SearchMemories(ctx, searchReq)
	suite.NoError(err, "Failed to search memories")
	suite.GreaterOrEqual(total, 1, "Should find at least one memory")
	suite.NotEmpty(memories, "Memory list should not be empty")

	// Verify the found memory
	found := false
	for _, m := range memories {
		if m.ID == memory.ID {
			suite.Equal(memory.Title, m.Title)
			suite.Equal(memory.Content, m.Content)
			found = true
			break
		}
	}
	suite.True(found, "Should find the stored memory in search results")
}

// TestRelationshipManagement tests the relationship management workflow
func (suite *MCPIntegrationTestSuite) TestRelationshipManagement() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Relationship Test Session",
		TaskType: "security_review",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for relationship test")

	// Create two memories
	memory1Req := &models.MemoryCreateRequest{
		SessionID:  session.ID,
		Title:      "Authentication Bypass",
		Content:    "Found authentication bypass vulnerability",
		Category:   "vulnerability",
		Priority:   9,
		Confidence: 0.95,
		Tags:       []string{"authentication", "bypass", "critical"},
	}
	memory1, err := suite.memoryService.StoreMemory(ctx, memory1Req)
	suite.Require().NoError(err, "Failed to create first memory")

	memory2Req := &models.MemoryCreateRequest{
		SessionID:  session.ID,
		Title:      "Session Management Issue",
		Content:    "Found session management vulnerability",
		Category:   "vulnerability",
		Priority:   7,
		Confidence: 0.8,
		Tags:       []string{"session", "management", "medium"},
	}
	memory2, err := suite.memoryService.StoreMemory(ctx, memory2Req)
	suite.Require().NoError(err, "Failed to create second memory")

	// Create a relationship between the memories
	relationshipReq := &models.RelationshipCreateRequest{
		SourceID:    memory1.ID,
		TargetID:    memory2.ID,
		Type:        models.RelationshipType("related_to"),
		Strength:    0.8,
		Description: "Both vulnerabilities affect authentication mechanisms",
	}

	relationship, err := suite.relationshipService.CreateRelationship(ctx, relationshipReq)
	suite.NoError(err, "Failed to create relationship")
	suite.NotNil(relationship, "Relationship should not be nil")
	suite.Equal(memory1.ID, relationship.SourceID)
	suite.Equal(memory2.ID, relationship.TargetID)
	suite.Equal(models.RelationshipType("related_to"), relationship.Type)
	suite.Equal(float32(0.8), relationship.Strength)

	// Test listing relationships
	listReq := &models.RelationshipListRequest{
		SourceID: memory1.ID,
		Limit:    10,
		Offset:   0,
	}
	relationships, total, err := suite.relationshipService.ListRelationships(ctx, listReq)
	suite.NoError(err, "Failed to list relationships")
	suite.GreaterOrEqual(total, 1, "Should find at least one relationship")
	suite.NotEmpty(relationships, "Relationship list should not be empty")
}

// TestContextSnapshot tests context snapshot functionality
func (suite *MCPIntegrationTestSuite) TestContextSnapshot() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Context Test Session",
		TaskType: "threat_modeling",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for context test")

	// Create a context snapshot
	contextData := map[string]interface{}{
		"current_focus": "authentication_system",
		"threats": []string{
			"credential_stuffing",
			"session_hijacking",
			"brute_force",
		},
		"assets": []string{
			"login_endpoint",
			"session_tokens",
			"user_database",
		},
	}

	snapshotReq := &models.ContextSnapshotCreateRequest{
		SessionID:   session.ID,
		Name:        "Authentication Analysis Context",
		ContextData: contextData,
		Description: "Current context for authentication system analysis",
	}

	snapshot, err := suite.contextService.CreateContextSnapshot(ctx, snapshotReq)
	suite.NoError(err, "Failed to create context snapshot")
	suite.NotNil(snapshot, "Context snapshot should not be nil")
	suite.Equal("Authentication Analysis Context", snapshot.Name)
	suite.Equal(session.ID, snapshot.SessionID)

	// Verify context data is preserved
	suite.NotNil(snapshot.ContextData, "Context data should not be nil")
	if raw, ok := snapshot.ContextData["raw"]; ok {
		// If stored as raw string, parse it
		var parsedData map[string]interface{}
		err := json.Unmarshal([]byte(raw.(string)), &parsedData)
		suite.NoError(err, "Failed to parse context data")
		suite.Equal("authentication_system", parsedData["current_focus"])
	}
}

// TestTaskProgress tests task progress functionality
func (suite *MCPIntegrationTestSuite) TestTaskProgress() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Task Progress Test Session",
		TaskType: "penetration_test",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for task progress test")

	// Create task progress entries
	taskReq := &models.TaskProgressCreateRequest{
		SessionID:          session.ID,
		TaskName:           "Network Reconnaissance",
		Stage:              "reconnaissance",
		Status:             "in_progress",
		ProgressPercentage: 25.0,
		Notes:              "Completed port scanning, starting service enumeration",
	}

	task, err := suite.taskService.CreateTaskProgress(ctx, taskReq)
	suite.NoError(err, "Failed to create task progress")
	suite.NotNil(task, "Task progress should not be nil")
	suite.Equal("Network Reconnaissance", task.TaskName)
	suite.Equal("reconnaissance", task.Stage)
	suite.Equal("in_progress", task.Status)
	suite.Equal(float32(25.0), task.ProgressPercentage)

	// Update task progress
	updateReq := &models.TaskProgressUpdateRequest{
		Stage:              "vulnerability_scanning",
		Status:             "in_progress",
		ProgressPercentage: 50.0,
		Notes:              "Completed reconnaissance, starting vulnerability scanning",
	}

	updatedTask, err := suite.taskService.UpdateTaskProgress(ctx, task.ID, updateReq)
	suite.NoError(err, "Failed to update task progress")
	suite.Equal("vulnerability_scanning", updatedTask.Stage)
	suite.Equal(float32(50.0), updatedTask.ProgressPercentage)
}

// TestErrorHandling tests error handling scenarios
func (suite *MCPIntegrationTestSuite) TestErrorHandling() {
	ctx := context.Background()

	// Test creating session with invalid data
	invalidSessionReq := &models.SessionCreateRequest{
		Name:     "", // Empty name should fail
		TaskType: "penetration_test",
	}
	_, err := suite.sessionService.CreateSession(ctx, invalidSessionReq)
	suite.Error(err, "Should return error for empty session name")

	// Test getting non-existent session
	_, err = suite.sessionService.GetSession(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent session")

	// Test storing memory with invalid session ID
	invalidMemoryReq := &models.MemoryCreateRequest{
		SessionID:  "non-existent-session-id",
		Title:      "Test Memory",
		Content:    "Test content",
		Category:   "test",
		Priority:   5,
		Confidence: 0.5,
	}
	_, err = suite.memoryService.StoreMemory(ctx, invalidMemoryReq)
	suite.Error(err, "Should return error for invalid session ID")
}

// TestDataIntegrity tests data integrity and consistency
func (suite *MCPIntegrationTestSuite) TestDataIntegrity() {
	ctx := context.Background()

	// Create a session with complex metadata
	metadata := map[string]interface{}{
		"priority": "high",
		"client":   "test-client",
		"tags":     []string{"security", "assessment"},
		"nested": map[string]interface{}{
			"value": 123,
			"flag":  true,
		},
	}

	sessionReq := &models.SessionCreateRequest{
		Name:        "Data Integrity Test Session",
		TaskType:    "security_review",
		Description: "Testing data integrity",
		Metadata:    metadata,
	}

	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.NoError(err, "Failed to create session with complex metadata")

	// Verify metadata is preserved correctly
	suite.NotNil(session.Metadata, "Metadata should not be nil")
	suite.Equal("high", session.Metadata["priority"])
	suite.Equal("test-client", session.Metadata["client"])

	// Verify nested metadata
	nested, ok := session.Metadata["nested"].(map[string]interface{})
	suite.True(ok, "Nested metadata should be preserved")
	suite.Equal(123, nested["value"])
	suite.Equal(true, nested["flag"])

	// Test that timestamps are set correctly
	suite.NotZero(session.CreatedAt, "CreatedAt should not be zero")
	suite.NotZero(session.UpdatedAt, "UpdatedAt should not be zero")
	suite.True(session.UpdatedAt.After(session.CreatedAt) || session.UpdatedAt.Equal(session.CreatedAt), "UpdatedAt should be >= CreatedAt")
}

// Run the test suite
func TestMCPIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MCPIntegrationTestSuite))
}

// TestMCPToolSimulation simulates MCP tool calls
func TestMCPToolSimulation(t *testing.T) {
	// This test simulates the MCP tool calls that would be made by Cursor
	// It tests the same functionality but in a more direct way

	// Create a temporary data directory
	dataDir := "./test_mcp_simulation_pb_data"
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	// Initialize PocketBase
	config := pocketbase.Config{
		DefaultDataDir: dataDir,
	}
	app := pocketbase.NewWithConfig(config)
	err := app.Bootstrap()
	assert.NoError(t, err, "Failed to bootstrap PocketBase")

	// Initialize collections
	collections := []*models.Collection{
		database.CreateSessionsCollection(),
		database.CreateMemoryEntriesCollection(),
		database.CreateRelationshipsCollection(),
		database.CreateContextSnapshotsCollection(),
		database.CreateTaskProgressCollection(),
	}

	for _, collection := range collections {
		existing, err := app.Dao().FindCollectionByNameOrId(collection.Name)
		if err != nil {
			err := app.Dao().SaveCollection(collection)
			assert.NoError(t, err, "Failed to create collection %s", collection.Name)
		}
	}

	// Initialize services
	sessionRepo := repository.NewSessionRepositoryV2(app)
	memoryRepo := repository.NewMemoryRepositoryV2(app)
	sessionService := services.NewSessionServiceV2(sessionRepo)
	memoryService := services.NewMemoryServiceV2(memoryRepo)

	ctx := context.Background()

	// Simulate create_session MCP tool call
	sessionReq := &models.SessionCreateRequest{
		Name:        "MCP Simulation Test",
		TaskType:    "penetration_test",
		Description: "Testing MCP tool simulation",
		Metadata: map[string]interface{}{
			"test": true,
		},
	}

	session, err := sessionService.CreateSession(ctx, sessionReq)
	assert.NoError(t, err, "Failed to create session via MCP simulation")
	assert.NotNil(t, session, "Session should not be nil")
	assert.Equal(t, "MCP Simulation Test", session.Name)

	// Simulate store_memory MCP tool call
	memoryReq := &models.MemoryCreateRequest{
		SessionID:   session.ID,
		Title:       "Simulated Vulnerability",
		Content:     "This is a simulated vulnerability finding",
		Category:    "vulnerability",
		Priority:    7,
		Confidence:  0.85,
		Tags:        []string{"simulation", "test"},
		Source:      "mcp_simulation",
		ContentType: "text",
	}

	memory, err := memoryService.StoreMemory(ctx, memoryReq)
	assert.NoError(t, err, "Failed to store memory via MCP simulation")
	assert.NotNil(t, memory, "Memory should not be nil")
	assert.Equal(t, "Simulated Vulnerability", memory.Title)
	assert.Equal(t, session.ID, memory.SessionID)
}
