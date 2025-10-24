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

func TestSessionManagementSimple(t *testing.T) {
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

	// Test 2: Get session by ID
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

	// Test 3: List sessions
	t.Run("ListSessions", func(t *testing.T) {
		// List all sessions
		listReq := &models.SessionListRequest{
			Limit:  10,
			Offset: 0,
		}

		sessions, totalCount, err := sessionRepo.ListSessions(ctx, listReq)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, totalCount, 2) // At least 2 sessions from previous tests
		assert.GreaterOrEqual(t, len(sessions), 2)

		t.Logf("✅ Listed %d sessions (total: %d)", len(sessions), totalCount)
	})

	// Test 4: Update session
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

	// Test 5: Delete session
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

	t.Log("🎉 All session management tests passed!")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}


