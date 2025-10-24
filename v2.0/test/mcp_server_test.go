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

// MCPServerTestSuite tests the MCP server functionality
type MCPServerTestSuite struct {
	suite.Suite
	app                 *pocketbase.PocketBase
	dataDir             string
	sessionService      *services.SessionServiceV2
	memoryService       *services.MemoryServiceV2
	relationshipService *services.RelationshipServiceV2
	contextService      *services.ContextServiceV2
	taskService         *services.TaskServiceV2
}

// SetupSuite initializes the test suite
func (suite *MCPServerTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_mcp_server_pb_data"
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
func (suite *MCPServerTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// initializeCollections creates all required database collections
func (suite *MCPServerTestSuite) initializeCollections() {
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

// TestSessionMCPTools tests the session-related MCP tools
func (suite *MCPServerTestSuite) TestSessionMCPTools() {
	ctx := context.Background()

	// Test create_session tool
	sessionReq := &models.SessionCreateRequest{
		Name:        "MCP Test Session",
		TaskType:    "penetration_test",
		Description: "Testing MCP session creation",
		Metadata: map[string]interface{}{
			"test":     true,
			"priority": "high",
		},
	}

	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.NoError(err, "Failed to create session via MCP tool")
	suite.NotNil(session, "Session should not be nil")
	suite.Equal("MCP Test Session", session.Name)
	suite.Equal("penetration_test", session.TaskType)
	suite.Equal("active", session.Status)

	// Test get_session tool
	retrievedSession, err := suite.sessionService.GetSession(ctx, session.ID)
	suite.NoError(err, "Failed to retrieve session via MCP tool")
	suite.Equal(session.ID, retrievedSession.ID)
	suite.Equal(session.Name, retrievedSession.Name)

	// Test list_sessions tool
	listReq := &models.SessionListRequest{
		TaskType: "penetration_test",
		Limit:    10,
		Offset:   0,
	}
	sessions, totalCount, err := suite.sessionService.ListSessions(ctx, listReq)
	suite.NoError(err, "Failed to list sessions via MCP tool")
	suite.GreaterOrEqual(totalCount, 1, "Should have at least one session")
	suite.NotEmpty(sessions, "Session list should not be empty")

	// Verify the session is in the list
	found := false
	for _, s := range sessions {
		if s.ID == session.ID {
			suite.Equal(session.Name, s.Name)
			found = true
			break
		}
	}
	suite.True(found, "Created session should be in the list")
}

// TestMemoryMCPTools tests the memory-related MCP tools
func (suite *MCPServerTestSuite) TestMemoryMCPTools() {
	ctx := context.Background()

	// First create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Memory MCP Test Session",
		TaskType: "vulnerability_assessment",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for memory test")

	// Test store_memory tool
	memoryReq := &models.MemoryCreateRequest{
		SessionID:   session.ID,
		Title:       "XSS Vulnerability Found",
		Content:     "Discovered reflected XSS vulnerability in search parameter",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.9,
		Tags:        []string{"xss", "reflected", "search", "critical"},
		Source:      "automated_scanning",
		ContentType: "text",
	}

	memory, err := suite.memoryService.StoreMemory(ctx, memoryReq)
	suite.NoError(err, "Failed to store memory via MCP tool")
	suite.NotNil(memory, "Memory should not be nil")
	suite.Equal("XSS Vulnerability Found", memory.Title)
	suite.Equal("vulnerability", memory.Category)
	suite.Equal(8, memory.Priority)
	suite.Equal(float32(0.9), memory.Confidence)
	suite.Equal(session.ID, memory.SessionID)

	// Test search_memories tool
	searchReq := &models.MemorySearchRequest{
		SessionID: session.ID,
		Query:     "XSS",
		Category:  "vulnerability",
		Tags:      []string{"xss"},
		Limit:     10,
		Offset:    0,
	}
	memories, total, err := suite.memoryService.SearchMemories(ctx, searchReq)
	suite.NoError(err, "Failed to search memories via MCP tool")
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

// TestRelationshipMCPTools tests the relationship-related MCP tools
func (suite *MCPServerTestSuite) TestRelationshipMCPTools() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Relationship MCP Test Session",
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

	// Test create_relationship tool
	relationshipReq := &models.RelationshipCreateRequest{
		SourceID:    memory1.ID,
		TargetID:    memory2.ID,
		Type:        models.RelationshipType("related_to"),
		Strength:    0.8,
		Description: "Both vulnerabilities affect authentication mechanisms",
	}

	relationship, err := suite.relationshipService.CreateRelationship(ctx, relationshipReq)
	suite.NoError(err, "Failed to create relationship via MCP tool")
	suite.NotNil(relationship, "Relationship should not be nil")
	suite.Equal(memory1.ID, relationship.SourceID)
	suite.Equal(memory2.ID, relationship.TargetID)
	suite.Equal(models.RelationshipType("related_to"), relationship.Type)
	suite.Equal(float32(0.8), relationship.Strength)

	// Test list_relationships tool
	listReq := &models.RelationshipListRequest{
		SourceID: memory1.ID,
		Limit:    10,
		Offset:   0,
	}
	relationships, total, err := suite.relationshipService.ListRelationships(ctx, listReq)
	suite.NoError(err, "Failed to list relationships via MCP tool")
	suite.GreaterOrEqual(total, 1, "Should find at least one relationship")
	suite.NotEmpty(relationships, "Relationship list should not be empty")

	// Verify the relationship is in the list
	found := false
	for _, r := range relationships {
		if r.ID == relationship.ID {
			suite.Equal(relationship.SourceID, r.SourceID)
			suite.Equal(relationship.TargetID, r.TargetID)
			found = true
			break
		}
	}
	suite.True(found, "Created relationship should be in the list")
}

// TestContextSnapshotMCPTools tests the context snapshot MCP tools
func (suite *MCPServerTestSuite) TestContextSnapshotMCPTools() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Context MCP Test Session",
		TaskType: "threat_modeling",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for context test")

	// Test create_context_snapshot tool
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
		"risk_level": "high",
	}

	snapshotReq := &models.ContextSnapshotCreateRequest{
		SessionID:   session.ID,
		Name:        "Authentication Analysis Context",
		ContextData: contextData,
		Description: "Current context for authentication system analysis",
	}

	snapshot, err := suite.contextService.CreateContextSnapshot(ctx, snapshotReq)
	suite.NoError(err, "Failed to create context snapshot via MCP tool")
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
		suite.Equal("high", parsedData["risk_level"])
	}
}

// TestTaskProgressMCPTools tests the task progress MCP tools
func (suite *MCPServerTestSuite) TestTaskProgressMCPTools() {
	ctx := context.Background()

	// Create a session
	sessionReq := &models.SessionCreateRequest{
		Name:     "Task Progress MCP Test Session",
		TaskType: "penetration_test",
	}
	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session for task progress test")

	// Test create_task_progress tool
	taskReq := &models.TaskProgressCreateRequest{
		SessionID:          session.ID,
		TaskName:           "Network Reconnaissance",
		Stage:              "reconnaissance",
		Status:             "in_progress",
		ProgressPercentage: 25.0,
		Notes:              "Completed port scanning, starting service enumeration",
	}

	task, err := suite.taskService.CreateTaskProgress(ctx, taskReq)
	suite.NoError(err, "Failed to create task progress via MCP tool")
	suite.NotNil(task, "Task progress should not be nil")
	suite.Equal("Network Reconnaissance", task.TaskName)
	suite.Equal("reconnaissance", task.Stage)
	suite.Equal("in_progress", task.Status)
	suite.Equal(float32(25.0), task.ProgressPercentage)
	suite.Equal(session.ID, task.SessionID)

	// Test updating task progress
	updateReq := &models.TaskProgressUpdateRequest{
		Stage:              "vulnerability_scanning",
		Status:             "in_progress",
		ProgressPercentage: 50.0,
		Notes:              "Completed reconnaissance, starting vulnerability scanning",
	}

	updatedTask, err := suite.taskService.UpdateTaskProgress(ctx, task.ID, updateReq)
	suite.NoError(err, "Failed to update task progress via MCP tool")
	suite.Equal("vulnerability_scanning", updatedTask.Stage)
	suite.Equal(float32(50.0), updatedTask.ProgressPercentage)
	suite.Equal("Completed reconnaissance, starting vulnerability scanning", updatedTask.Notes)
}

// TestMCPErrorHandling tests error handling in MCP tools
func (suite *MCPServerTestSuite) TestMCPErrorHandling() {
	ctx := context.Background()

	// Test invalid session creation
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

	// Test creating relationship with invalid memory IDs
	invalidRelationshipReq := &models.RelationshipCreateRequest{
		SourceID: "non-existent-source-id",
		TargetID: "non-existent-target-id",
		Type:     models.RelationshipType("related_to"),
		Strength: 0.5,
	}
	_, err = suite.relationshipService.CreateRelationship(ctx, invalidRelationshipReq)
	suite.Error(err, "Should return error for invalid memory IDs")
}

// TestMCPDataConsistency tests data consistency across MCP operations
func (suite *MCPServerTestSuite) TestMCPDataConsistency() {
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
		Name:        "Data Consistency Test Session",
		TaskType:    "security_review",
		Description: "Testing data consistency across MCP operations",
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

	// Create a memory for this session
	memoryReq := &models.MemoryCreateRequest{
		SessionID:   session.ID,
		Title:       "Consistency Test Memory",
		Content:     "Testing memory consistency",
		Category:    "test",
		Priority:    5,
		Confidence:  0.5,
		Tags:        []string{"consistency", "test"},
		Source:      "mcp_test",
		ContentType: "text",
	}

	memory, err := suite.memoryService.StoreMemory(ctx, memoryReq)
	suite.NoError(err, "Failed to store memory")
	suite.Equal(session.ID, memory.SessionID, "Memory should be linked to the correct session")

	// Verify the memory can be retrieved
	retrievedMemory, err := suite.memoryService.GetMemory(ctx, memory.ID)
	suite.NoError(err, "Failed to retrieve memory")
	suite.Equal(memory.ID, retrievedMemory.ID, "Memory IDs should match")
	suite.Equal(memory.Title, retrievedMemory.Title, "Memory titles should match")
	suite.Equal(memory.SessionID, retrievedMemory.SessionID, "Memory session IDs should match")
}

// Run the test suite
func TestMCPServerTestSuite(t *testing.T) {
	suite.Run(t, new(MCPServerTestSuite))
}

// TestMCPToolParameterValidation tests parameter validation for MCP tools
func TestMCPToolParameterValidation(t *testing.T) {
	// Create a temporary data directory
	dataDir := "./test_mcp_validation_pb_data"
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

	// Test session parameter validation
	t.Run("SessionParameterValidation", func(t *testing.T) {
		// Test empty name
		req := &models.SessionCreateRequest{
			Name:     "",
			TaskType: "penetration_test",
		}
		_, err := sessionService.CreateSession(ctx, req)
		assert.Error(t, err, "Should return error for empty name")

		// Test empty task type
		req = &models.SessionCreateRequest{
			Name:     "Valid Name",
			TaskType: "",
		}
		_, err = sessionService.CreateSession(ctx, req)
		assert.Error(t, err, "Should return error for empty task type")
	})

	// Test memory parameter validation
	t.Run("MemoryParameterValidation", func(t *testing.T) {
		// First create a valid session
		sessionReq := &models.SessionCreateRequest{
			Name:     "Validation Test Session",
			TaskType: "penetration_test",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		assert.NoError(t, err, "Failed to create session for validation test")

		// Test empty title
		memoryReq := &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "",
			Content:    "Valid content",
			Category:   "test",
			Priority:   5,
			Confidence: 0.5,
		}
		_, err = memoryService.StoreMemory(ctx, memoryReq)
		assert.Error(t, err, "Should return error for empty title")

		// Test empty content
		memoryReq = &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Valid title",
			Content:    "",
			Category:   "test",
			Priority:   5,
			Confidence: 0.5,
		}
		_, err = memoryService.StoreMemory(ctx, memoryReq)
		assert.Error(t, err, "Should return error for empty content")

		// Test empty category
		memoryReq = &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Valid title",
			Content:    "Valid content",
			Category:   "",
			Priority:   5,
			Confidence: 0.5,
		}
		_, err = memoryService.StoreMemory(ctx, memoryReq)
		assert.Error(t, err, "Should return error for empty category")

		// Test invalid priority
		memoryReq = &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Valid title",
			Content:    "Valid content",
			Category:   "test",
			Priority:   0, // Invalid priority
			Confidence: 0.5,
		}
		_, err = memoryService.StoreMemory(ctx, memoryReq)
		assert.Error(t, err, "Should return error for invalid priority")

		// Test invalid confidence
		memoryReq = &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Valid title",
			Content:    "Valid content",
			Category:   "test",
			Priority:   5,
			Confidence: 1.5, // Invalid confidence
		}
		_, err = memoryService.StoreMemory(ctx, memoryReq)
		assert.Error(t, err, "Should return error for invalid confidence")
	})
}
