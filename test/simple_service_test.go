package test

import (
	"context"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

// TestSimpleServiceIntegration tests the service layer directly
func TestSimpleServiceIntegration(t *testing.T) {
	// Create a temporary data directory
	dataDir := "./test_simple_service_pb_data"
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	// Initialize PocketBase
	config := pocketbase.Config{
		DefaultDataDir: dataDir,
	}
	app := pocketbase.NewWithConfig(config)
	err := app.Bootstrap()
	require.NoError(t, err, "Failed to bootstrap PocketBase")

	// Initialize collections
	collections := []*pbmodels.Collection{
		database.CreateSessionsCollection(),
		database.CreateMemoryEntriesCollection(),
		database.CreateRelationshipsCollection(),
		database.CreateContextSnapshotsCollection(),
		database.CreateTaskProgressCollection(),
	}

	for _, collection := range collections {
		_, err := app.Dao().FindCollectionByNameOrId(collection.Name)
		if err != nil {
			err := app.Dao().SaveCollection(collection)
			require.NoError(t, err, "Failed to create collection %s", collection.Name)
		}
	}

	// Initialize services
	sessionRepo := repository.NewSessionRepositoryV2(app)
	memoryRepo := repository.NewMemoryRepositoryV2(app)
	sessionService := services.NewSessionServiceV2(sessionRepo)
	memoryService := services.NewMemoryServiceV2(memoryRepo)

	ctx := context.Background()

	// Test session creation
	t.Run("SessionCreation", func(t *testing.T) {
		sessionReq := &models.SessionCreateRequest{
			Name:        "Test Session",
			TaskType:    "penetration_test",
			Description: "Testing session creation",
			Metadata: map[string]interface{}{
				"test": true,
			},
		}

		session, err := sessionService.CreateSession(ctx, sessionReq)
		assert.NoError(t, err, "Failed to create session")
		assert.NotNil(t, session, "Session should not be nil")
		assert.Equal(t, "Test Session", session.Name)
		assert.Equal(t, "penetration_test", session.TaskType)
		assert.Equal(t, "active", session.Status)
		assert.NotEmpty(t, session.ID, "Session ID should not be empty")

		// Test retrieving the session
		retrievedSession, err := sessionService.GetSession(ctx, session.ID)
		assert.NoError(t, err, "Failed to retrieve session")
		assert.Equal(t, session.ID, retrievedSession.ID)
		assert.Equal(t, session.Name, retrievedSession.Name)
	})

	// Test memory storage
	t.Run("MemoryStorage", func(t *testing.T) {
		// First create a session
		sessionReq := &models.SessionCreateRequest{
			Name:     "Memory Test Session",
			TaskType: "vulnerability_assessment",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err, "Failed to create session for memory test")

		// Test storing a memory
		memoryReq := &models.MemoryCreateRequest{
			SessionID:   session.ID,
			Title:       "Test Vulnerability",
			Content:     "This is a test vulnerability finding",
			Category:    "vulnerability",
			Priority:    7,
			Confidence:  0.85,
			Tags:        []string{"test", "vulnerability"},
			Source:      "manual_testing",
			ContentType: "text",
		}

		memory, err := memoryService.StoreMemory(ctx, memoryReq)
		assert.NoError(t, err, "Failed to store memory")
		assert.NotNil(t, memory, "Memory should not be nil")
		assert.Equal(t, "Test Vulnerability", memory.Title)
		assert.Equal(t, "vulnerability", memory.Category)
		assert.Equal(t, 7, memory.Priority)
		assert.Equal(t, float32(0.85), memory.Confidence)
		assert.Equal(t, session.ID, memory.SessionID)

		// Test retrieving the memory
		retrievedMemory, err := memoryService.GetMemory(ctx, memory.ID)
		assert.NoError(t, err, "Failed to retrieve memory")
		assert.Equal(t, memory.ID, retrievedMemory.ID)
		assert.Equal(t, memory.Title, retrievedMemory.Title)
		assert.Equal(t, memory.SessionID, retrievedMemory.SessionID)
	})

	// Test error handling
	t.Run("ErrorHandling", func(t *testing.T) {
		// Test invalid session creation
		invalidSessionReq := &models.SessionCreateRequest{
			Name:     "", // Empty name should fail
			TaskType: "penetration_test",
		}
		_, err := sessionService.CreateSession(ctx, invalidSessionReq)
		assert.Error(t, err, "Should return error for empty session name")

		// Test getting non-existent session
		_, err = sessionService.GetSession(ctx, "non-existent-id")
		assert.Error(t, err, "Should return error for non-existent session")

		// Test storing memory with invalid session ID
		invalidMemoryReq := &models.MemoryCreateRequest{
			SessionID:  "non-existent-session-id",
			Title:      "Test Memory",
			Content:    "Test content",
			Category:   "test",
			Priority:   5,
			Confidence: 0.5,
		}
		_, err = memoryService.StoreMemory(ctx, invalidMemoryReq)
		assert.Error(t, err, "Should return error for invalid session ID")
	})
}
