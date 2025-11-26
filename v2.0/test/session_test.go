package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

func TestSessionManagement(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase client
	pbClient, err := database.NewSimplePocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	require.NoError(t, err)

	// Get PocketBase app
	app := pbClient.GetApp()

	// Initialize session repository
	sessionRepo := repository.NewSessionRepositoryV2(app)

	// Test 1: Create a session
	t.Run("CreateSession", func(t *testing.T) {
		req := &models.SessionCreateRequest{
			Name:        "Security Assessment Test",
			TaskType:    "security_review",
			Description: "Testing TinyBrain v2.0 session creation",
			Metadata:    map[string]interface{}{"priority": "high", "client": "test-client"},
		}

		session, err := sessionRepo.CreateSession(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, req.Name, session.Name)
		assert.Equal(t, req.TaskType, session.TaskType)
		assert.Equal(t, "active", session.Status) // Default status
		assert.Equal(t, req.Description, session.Description)
		assert.Equal(t, req.Metadata["priority"], session.Metadata["priority"])
		assert.Equal(t, req.Metadata["client"], session.Metadata["client"])

		t.Logf("✅ Created session: %s", session.ID)
	})

	// Test 2: Create multiple sessions
	t.Run("CreateMultipleSessions", func(t *testing.T) {
		sessions := []*models.SessionCreateRequest{
			{Name: "Penetration Test", TaskType: "penetration_test", Description: "Network penetration testing"},
			{Name: "Code Review", TaskType: "code_review", Description: "Security code review"},
			{Name: "Threat Modeling", TaskType: "threat_modeling", Description: "Application threat modeling"},
		}

		createdSessions := make([]*models.Session, len(sessions))
		for i, req := range sessions {
			session, err := sessionRepo.CreateSession(ctx, req)
			require.NoError(t, err)
			assert.NotEmpty(t, session.ID)
			assert.Equal(t, req.Name, session.Name)
			assert.Equal(t, req.TaskType, session.TaskType)
			createdSessions[i] = session
		}

		t.Logf("✅ Created %d sessions", len(createdSessions))
	})

	// Test 3: List sessions
	t.Run("ListSessions", func(t *testing.T) {
		// List all sessions
		listReq := &models.SessionListRequest{
			Limit:  10,
			Offset: 0,
		}

		sessions, totalCount, err := sessionRepo.ListSessions(ctx, listReq)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, 4) // At least 4 sessions from previous tests
		assert.GreaterOrEqual(t, len(sessions), 4)

		// List sessions by task type
		listReq.TaskType = "security_review"
		sessions, totalCount, err = sessionRepo.ListSessions(ctx, listReq)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, 1)
		assert.GreaterOrEqual(t, len(sessions), 1)

		// Verify all returned sessions match the filter
		for _, session := range sessions {
			assert.Equal(t, "security_review", session.TaskType)
		}

		t.Logf("✅ Listed %d sessions (total: %d)", len(sessions), totalCount)
	})

	// Test 4: Get session by ID
	t.Run("GetSession", func(t *testing.T) {
		// First create a session
		req := &models.SessionCreateRequest{
			Name:     "Get Session Test",
			TaskType: "vulnerability_analysis",
		}

		createdSession, err := sessionRepo.CreateSession(ctx, req)
		require.NoError(t, err)

		// Now retrieve it
		retrievedSession, err := sessionRepo.GetSession(ctx, createdSession.ID)
		require.NoError(t, err)
		assert.Equal(t, createdSession.ID, retrievedSession.ID)
		assert.Equal(t, createdSession.Name, retrievedSession.Name)
		assert.Equal(t, createdSession.TaskType, retrievedSession.TaskType)
		assert.Equal(t, createdSession.Status, retrievedSession.Status)

		t.Logf("✅ Retrieved session: %s", retrievedSession.ID)
	})

	// Test 5: Update session
	t.Run("UpdateSession", func(t *testing.T) {
		// First create a session
		req := &models.SessionCreateRequest{
			Name:     "Update Test Session",
			TaskType: "incident_response",
		}

		createdSession, err := sessionRepo.CreateSession(ctx, req)
		require.NoError(t, err)

		// Update the session
		updateReq := &models.SessionUpdateRequest{
			Name:        stringPtr("Updated Session Name"),
			Status:      stringPtr("completed"),
			Description: stringPtr("This session has been updated"),
			Metadata:    map[string]interface{}{"updated": true, "version": "2.0"},
		}

		updatedSession, err := sessionRepo.UpdateSession(ctx, createdSession.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, createdSession.ID, updatedSession.ID)
		assert.Equal(t, "Updated Session Name", updatedSession.Name)
		assert.Equal(t, "completed", updatedSession.Status)
		assert.Equal(t, "This session has been updated", updatedSession.Description)
		assert.Equal(t, true, updatedSession.Metadata["updated"])
		assert.Equal(t, "2.0", updatedSession.Metadata["version"])

		t.Logf("✅ Updated session: %s", updatedSession.ID)
	})

	// Test 6: Delete session
	t.Run("DeleteSession", func(t *testing.T) {
		// First create a session
		req := &models.SessionCreateRequest{
			Name:     "Delete Test Session",
			TaskType: "compliance_audit",
		}

		createdSession, err := sessionRepo.CreateSession(ctx, req)
		require.NoError(t, err)

		// Delete the session
		err = sessionRepo.DeleteSession(ctx, createdSession.ID)
		require.NoError(t, err)

		// Verify it's deleted
		_, err = sessionRepo.GetSession(ctx, createdSession.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		t.Logf("✅ Deleted session: %s", createdSession.ID)
	})

	// Test 7: Search sessions
	t.Run("SearchSessions", func(t *testing.T) {
		// Create a session with specific content
		req := &models.SessionCreateRequest{
			Name:        "Searchable Session",
			TaskType:    "security_review",
			Description: "This session contains the word 'vulnerability' for testing search",
		}

		createdSession, err := sessionRepo.CreateSession(ctx, req)
		require.NoError(t, err)

		// Search for sessions containing "vulnerability"
		searchReq := &models.SessionListRequest{
			Query:  "vulnerability",
			Limit:  10,
			Offset: 0,
		}

		sessions, totalCount, err := sessionRepo.ListSessions(ctx, searchReq)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, 1)
		assert.GreaterOrEqual(t, len(sessions), 1)

		// Verify the search found our session
		found := false
		for _, session := range sessions {
			if session.ID == createdSession.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Search should have found the created session")

		t.Logf("✅ Search found %d sessions", len(sessions))
	})

	t.Log("🎉 All session management tests passed!")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}


