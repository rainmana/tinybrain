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

func TestBasicSetup(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase client
	pbClient, err := database.NewPocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	require.NoError(t, err)

	// Test that we can create a session
	sessionRepo := repository.NewSessionRepository(pbClient.GetApp())

	req := &models.SessionCreateRequest{
		Name:        "Test Session",
		TaskType:    "security_review",
		Description: "A test security review session",
	}

	session, err := sessionRepo.Create(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, "Test Session", session.Name)
	assert.Equal(t, "security_review", session.TaskType)
	assert.Equal(t, "active", session.Status)

	// Test that we can retrieve the session
	retrieved, err := sessionRepo.Get(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.Name, retrieved.Name)

	// Test that we can list sessions
	sessions, count, err := sessionRepo.List(ctx, &models.SessionListRequest{
		Limit: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, sessions, 1)
	assert.Equal(t, session.ID, sessions[0].ID)

	t.Logf("✅ Basic setup test passed - created session %s", session.ID)
}

func TestDatabasePersistence(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase client
	pbClient, err := database.NewPocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	require.NoError(t, err)

	// Create a session
	sessionRepo := repository.NewSessionRepository(pbClient.GetApp())

	req := &models.SessionCreateRequest{
		Name:        "Persistence Test",
		TaskType:    "penetration_test",
		Description: "Testing database persistence",
	}

	session, err := sessionRepo.Create(ctx, req)
	require.NoError(t, err)

	// Close and reopen the database
	pbClient.Close()

	// Reinitialize with same directory
	pbClient2, err := database.NewPocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient2.Close()

	// Bootstrap again
	err = pbClient2.Bootstrap(ctx)
	require.NoError(t, err)

	// Try to retrieve the session
	sessionRepo2 := repository.NewSessionRepository(pbClient2.GetApp())
	retrieved, err := sessionRepo2.Get(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.Name, retrieved.Name)

	t.Logf("✅ Database persistence test passed - session %s persisted", session.ID)
}

func TestMultipleSessions(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase client
	pbClient, err := database.NewPocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	require.NoError(t, err)

	// Create multiple sessions
	sessionRepo := repository.NewSessionRepository(pbClient.GetApp())

	sessions := []*models.SessionCreateRequest{
		{Name: "Session 1", TaskType: "security_review", Description: "First session"},
		{Name: "Session 2", TaskType: "penetration_test", Description: "Second session"},
		{Name: "Session 3", TaskType: "vulnerability_analysis", Description: "Third session"},
	}

	var createdSessions []*models.Session
	for _, req := range sessions {
		session, err := sessionRepo.Create(ctx, req)
		require.NoError(t, err)
		createdSessions = append(createdSessions, session)
	}

	// List all sessions
	allSessions, count, err := sessionRepo.List(ctx, &models.SessionListRequest{
		Limit: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, allSessions, 3)

	// Filter by task type
	securitySessions, count, err := sessionRepo.List(ctx, &models.SessionListRequest{
		TaskType: "security_review",
		Limit:    10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, securitySessions, 1)
	assert.Equal(t, "Session 1", securitySessions[0].Name)

	t.Logf("✅ Multiple sessions test passed - created %d sessions", len(createdSessions))
}
