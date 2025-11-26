package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
)

// Clean session repository for testing
type CleanSessionRepository struct {
	app *pocketbase.PocketBase
}

func NewCleanSessionRepository(app *pocketbase.PocketBase) *CleanSessionRepository {
	return &CleanSessionRepository{app: app}
}

func (r *CleanSessionRepository) CreateSession(ctx context.Context, req *models.SessionCreateRequest) (*models.Session, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("sessions")
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("name", req.Name)
	record.Set("task_type", req.TaskType)
	record.Set("status", "active") // Default status
	record.Set("description", req.Description)
	record.Set("metadata", req.Metadata)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create session record: %w", err)
	}

	return r.recordToSession(record), nil
}

func (r *CleanSessionRepository) GetSession(ctx context.Context, id string) (*models.Session, error) {
	record, err := r.app.Dao().FindRecordById("sessions", id)
	if err != nil {
		if daos.IsNotFoundError(err) {
			return nil, fmt.Errorf("session with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get session record: %w", err)
	}

	return r.recordToSession(record), nil
}

func (r *CleanSessionRepository) recordToSession(record *pbmodels.Record) *models.Session {
	// Handle metadata safely
	var metadata map[string]interface{}
	if raw := record.Get(types.FieldTypeJson, "metadata"); raw != nil {
		if m, ok := raw.(map[string]interface{}); ok {
			metadata = m
		}
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &models.Session{
		ID:          record.Id,
		Name:        record.GetString("name"),
		TaskType:    record.GetString("task_type"),
		Status:      record.GetString("status"),
		Description: record.GetString("description"),
		Metadata:    metadata,
		CreatedAt:   record.Created.Time(),
		UpdatedAt:   record.Updated.Time(),
	}
}

func TestSessionManagementClean(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase with config
	config := pocketbase.Config{
		DefaultDataDir: tempDir,
	}
	app := pocketbase.NewWithConfig(config)

	// Bootstrap the app
	if err := app.Bootstrap(); err != nil {
		t.Fatalf("Failed to bootstrap PocketBase app: %v", err)
	}

	// Initialize PocketBase client
	pbClient, err := database.NewSimplePocketBaseClient(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize PocketBase client: %v", err)
	}
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	if err != nil {
		t.Fatalf("Failed to bootstrap database: %v", err)
	}

	// Initialize session repository
	sessionRepo := NewCleanSessionRepository(app)

	// Test 1: Create a session
	t.Run("CreateSession", func(t *testing.T) {
		req := &models.SessionCreateRequest{
			Name:        "Security Assessment Test",
			TaskType:    "security_review",
			Description: "Testing TinyBrain v2.0 session creation",
			Metadata:    map[string]interface{}{"priority": "high", "client": "test-client"},
		}

		session, err := sessionRepo.CreateSession(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.ID == "" {
			t.Fatal("Session ID should not be empty")
		}
		if session.Name != req.Name {
			t.Fatalf("Expected name %s, got %s", req.Name, session.Name)
		}
		if session.TaskType != req.TaskType {
			t.Fatalf("Expected task_type %s, got %s", req.TaskType, session.TaskType)
		}
		if session.Status != "active" {
			t.Fatalf("Expected status 'active', got %s", session.Status)
		}

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
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Now retrieve it
		retrievedSession, err := sessionRepo.GetSession(ctx, createdSession.ID)
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		if retrievedSession.ID != createdSession.ID {
			t.Fatalf("Expected ID %s, got %s", createdSession.ID, retrievedSession.ID)
		}
		if retrievedSession.Name != createdSession.Name {
			t.Fatalf("Expected name %s, got %s", createdSession.Name, retrievedSession.Name)
		}

		t.Logf("✅ Retrieved session: %s", retrievedSession.ID)
	})

	t.Log("🎉 All session management tests passed!")
}


