package test

import (
	"context"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/stretchr/testify/suite"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// SessionRepositoryTestSuite tests the SessionRepositoryV2
type SessionRepositoryTestSuite struct {
	suite.Suite
	app     *pocketbase.PocketBase
	repo    *repository.SessionRepositoryV2
	dataDir string
}

// SetupSuite initializes the test suite
func (suite *SessionRepositoryTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_pb_data"
	os.RemoveAll(suite.dataDir) // Clean up any existing test data

	// Initialize PocketBase with test config
	config := pocketbase.Config{
		DefaultDataDir: suite.dataDir,
	}
	suite.app = pocketbase.NewWithConfig(config)

	// Bootstrap the app
	err := suite.app.Bootstrap()
	suite.Require().NoError(err, "Failed to bootstrap PocketBase for testing")

	// Create the repository
	suite.repo = repository.NewSessionRepositoryV2(suite.app)
}

// TearDownSuite cleans up after the test suite
func (suite *SessionRepositoryTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// TestCreateSession tests creating a new session
func (suite *SessionRepositoryTestSuite) TestCreateSession() {
	ctx := context.Background()

	req := &models.SessionCreateRequest{
		Name:        "Test Security Assessment",
		TaskType:    "penetration_test",
		Description: "A comprehensive security assessment",
		Metadata: map[string]interface{}{
			"priority": "high",
			"client":   "test-client",
		},
	}

	session, err := suite.repo.CreateSession(ctx, req)

	suite.NoError(err, "Failed to create session")
	suite.NotNil(session, "Session should not be nil")
	suite.Equal(req.Name, session.Name, "Session name should match")
	suite.Equal(req.TaskType, session.TaskType, "Task type should match")
	suite.Equal(req.Description, session.Description, "Description should match")
	suite.Equal("active", session.Status, "Default status should be active")
	suite.NotEmpty(session.ID, "Session ID should not be empty")
	suite.NotZero(session.CreatedAt, "CreatedAt should not be zero")
	suite.NotZero(session.UpdatedAt, "UpdatedAt should not be zero")
}

// TestGetSession tests retrieving a session by ID
func (suite *SessionRepositoryTestSuite) TestGetSession() {
	ctx := context.Background()

	// First create a session
	req := &models.SessionCreateRequest{
		Name:     "Test Session for Get",
		TaskType: "vulnerability_assessment",
	}

	createdSession, err := suite.repo.CreateSession(ctx, req)
	suite.Require().NoError(err, "Failed to create session for get test")

	// Now retrieve it
	retrievedSession, err := suite.repo.GetSession(ctx, createdSession.ID)

	suite.NoError(err, "Failed to get session")
	suite.NotNil(retrievedSession, "Retrieved session should not be nil")
	suite.Equal(createdSession.ID, retrievedSession.ID, "Session IDs should match")
	suite.Equal(createdSession.Name, retrievedSession.Name, "Session names should match")
	suite.Equal(createdSession.TaskType, retrievedSession.TaskType, "Task types should match")
}

// TestGetSessionNotFound tests retrieving a non-existent session
func (suite *SessionRepositoryTestSuite) TestGetSessionNotFound() {
	ctx := context.Background()

	_, err := suite.repo.GetSession(ctx, "non-existent-id")

	suite.Error(err, "Should return error for non-existent session")
	suite.Contains(err.Error(), "not found", "Error should indicate session not found")
}

// TestListSessions tests listing sessions with filtering
func (suite *SessionRepositoryTestSuite) TestListSessions() {
	ctx := context.Background()

	// Create multiple sessions
	sessions := []*models.SessionCreateRequest{
		{Name: "Session 1", TaskType: "penetration_test"},
		{Name: "Session 2", TaskType: "vulnerability_assessment"},
		{Name: "Session 3", TaskType: "penetration_test"},
	}

	for _, req := range sessions {
		_, err := suite.repo.CreateSession(ctx, req)
		suite.Require().NoError(err, "Failed to create session for list test")
	}

	// Test listing all sessions
	listReq := &models.SessionListRequest{
		Limit:  10,
		Offset: 0,
	}

	sessionList, totalCount, err := suite.repo.ListSessions(ctx, listReq)

	suite.NoError(err, "Failed to list sessions")
	suite.GreaterOrEqual(totalCount, len(sessions), "Total count should be at least the number of created sessions")
	suite.NotEmpty(sessionList, "Session list should not be empty")

	// Test filtering by task type
	filteredReq := &models.SessionListRequest{
		TaskType: "penetration_test",
		Limit:    10,
		Offset:   0,
	}

	filteredList, _, err := suite.repo.ListSessions(ctx, filteredReq)

	suite.NoError(err, "Failed to list filtered sessions")
	suite.NotEmpty(filteredList, "Filtered session list should not be empty")

	// Verify all returned sessions have the correct task type
	for _, session := range filteredList {
		suite.Equal("penetration_test", session.TaskType, "All sessions should have penetration_test task type")
	}
}

// TestUpdateSession tests updating an existing session
func (suite *SessionRepositoryTestSuite) TestUpdateSession() {
	ctx := context.Background()

	// Create a session
	req := &models.SessionCreateRequest{
		Name:        "Original Session",
		TaskType:    "security_review",
		Description: "Original description",
	}

	createdSession, err := suite.repo.CreateSession(ctx, req)
	suite.Require().NoError(err, "Failed to create session for update test")

	// Update the session
	updateReq := &models.SessionUpdateRequest{
		Name:        "Updated Session",
		Status:      "completed",
		Description: "Updated description",
		Metadata: map[string]interface{}{
			"updated": true,
		},
	}

	updatedSession, err := suite.repo.UpdateSession(ctx, createdSession.ID, updateReq)

	suite.NoError(err, "Failed to update session")
	suite.NotNil(updatedSession, "Updated session should not be nil")
	suite.Equal(createdSession.ID, updatedSession.ID, "Session ID should remain the same")
	suite.Equal("Updated Session", updatedSession.Name, "Name should be updated")
	suite.Equal("completed", updatedSession.Status, "Status should be updated")
	suite.Equal("Updated description", updatedSession.Description, "Description should be updated")
	suite.True(updatedSession.UpdatedAt.After(createdSession.UpdatedAt), "UpdatedAt should be newer")
}

// TestDeleteSession tests deleting a session
func (suite *SessionRepositoryTestSuite) TestDeleteSession() {
	ctx := context.Background()

	// Create a session
	req := &models.SessionCreateRequest{
		Name:     "Session to Delete",
		TaskType: "security_review",
	}

	createdSession, err := suite.repo.CreateSession(ctx, req)
	suite.Require().NoError(err, "Failed to create session for delete test")

	// Delete the session
	err = suite.repo.DeleteSession(ctx, createdSession.ID)

	suite.NoError(err, "Failed to delete session")

	// Verify the session is deleted
	_, err = suite.repo.GetSession(ctx, createdSession.ID)
	suite.Error(err, "Should return error when trying to get deleted session")
	suite.Contains(err.Error(), "not found", "Error should indicate session not found")
}

// TestDeleteSessionNotFound tests deleting a non-existent session
func (suite *SessionRepositoryTestSuite) TestDeleteSessionNotFound() {
	ctx := context.Background()

	err := suite.repo.DeleteSession(ctx, "non-existent-id")

	suite.Error(err, "Should return error for non-existent session")
	suite.Contains(err.Error(), "not found", "Error should indicate session not found")
}

// TestSessionValidation tests session validation
func (suite *SessionRepositoryTestSuite) TestSessionValidation() {
	ctx := context.Background()

	// Test empty name
	req := &models.SessionCreateRequest{
		Name:     "",
		TaskType: "penetration_test",
	}

	_, err := suite.repo.CreateSession(ctx, req)
	suite.Error(err, "Should return error for empty name")

	// Test empty task type
	req = &models.SessionCreateRequest{
		Name:     "Valid Name",
		TaskType: "",
	}

	_, err = suite.repo.CreateSession(ctx, req)
	suite.Error(err, "Should return error for empty task type")
}

// TestSessionMetadataHandling tests metadata handling
func (suite *SessionRepositoryTestSuite) TestSessionMetadataHandling() {
	ctx := context.Background()

	metadata := map[string]interface{}{
		"priority": "high",
		"client":   "test-client",
		"tags":     []string{"security", "assessment"},
		"nested": map[string]interface{}{
			"value": 123,
			"flag":  true,
		},
	}

	req := &models.SessionCreateRequest{
		Name:        "Session with Metadata",
		TaskType:    "penetration_test",
		Description: "Session with complex metadata",
		Metadata:    metadata,
	}

	session, err := suite.repo.CreateSession(ctx, req)

	suite.NoError(err, "Failed to create session with metadata")
	suite.NotNil(session.Metadata, "Metadata should not be nil")
	suite.Equal("high", session.Metadata["priority"], "Priority should be preserved")
	suite.Equal("test-client", session.Metadata["client"], "Client should be preserved")

	// Verify nested metadata
	nested, ok := session.Metadata["nested"].(map[string]interface{})
	suite.True(ok, "Nested metadata should be preserved")
	suite.Equal(123, nested["value"], "Nested value should be preserved")
	suite.Equal(true, nested["flag"], "Nested flag should be preserved")
}

// Run the test suite
func TestSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}
