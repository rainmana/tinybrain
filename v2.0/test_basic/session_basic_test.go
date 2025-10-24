package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Basic session model
type BasicSession struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TaskType    string `json:"task_type"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type BasicSessionCreateRequest struct {
	Name        string `json:"name"`
	TaskType    string `json:"task_type"`
	Description string `json:"description"`
}

// Basic session repository
type BasicSessionRepository struct {
	app *pocketbase.PocketBase
}

func NewBasicSessionRepository(app *pocketbase.PocketBase) *BasicSessionRepository {
	return &BasicSessionRepository{app: app}
}

func (r *BasicSessionRepository) CreateSession(ctx context.Context, req *BasicSessionCreateRequest) (*BasicSession, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("sessions")
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("name", req.Name)
	record.Set("task_type", req.TaskType)
	record.Set("status", "active") // Default status
	record.Set("description", req.Description)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create session record: %w", err)
	}

	return r.recordToSession(record), nil
}

func (r *BasicSessionRepository) recordToSession(record *pbmodels.Record) *BasicSession {
	return &BasicSession{
		ID:          record.Id,
		Name:        record.GetString("name"),
		TaskType:    record.GetString("task_type"),
		Status:      record.GetString("status"),
		Description: record.GetString("description"),
	}
}

// Create basic sessions collection
func createBasicSessionsCollection() *pbmodels.Collection {
	collection := &pbmodels.Collection{
		Name:       "sessions",
		Type:       pbmodels.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add basic fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "name",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "task_type",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "status",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "description",
		Type: schema.FieldTypeText,
	})

	return collection
}

func TestBasicSessionManagement(t *testing.T) {
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

	// Create sessions collection
	collection := createBasicSessionsCollection()
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatalf("Failed to create sessions collection: %v", err)
	}

	// Initialize session repository
	sessionRepo := NewBasicSessionRepository(app)

	// Test 1: Create a session
	t.Run("CreateSession", func(t *testing.T) {
		req := &BasicSessionCreateRequest{
			Name:        "Security Assessment Test",
			TaskType:    "security_review",
			Description: "Testing TinyBrain v2.0 session creation",
		}

		session, err := sessionRepo.CreateSession(context.Background(), req)
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

	// Test 2: Create multiple sessions
	t.Run("CreateMultipleSessions", func(t *testing.T) {
		sessions := []*BasicSessionCreateRequest{
			{Name: "Penetration Test", TaskType: "penetration_test", Description: "Network penetration testing"},
			{Name: "Code Review", TaskType: "code_review", Description: "Security code review"},
			{Name: "Threat Modeling", TaskType: "threat_modeling", Description: "Application threat modeling"},
		}

		for i, req := range sessions {
			session, err := sessionRepo.CreateSession(context.Background(), req)
			if err != nil {
				t.Fatalf("Failed to create session %d: %v", i, err)
			}
			if session.ID == "" {
				t.Fatalf("Session %d ID should not be empty", i)
			}
			if session.Name != req.Name {
				t.Fatalf("Session %d: Expected name %s, got %s", i, req.Name, session.Name)
			}
		}

		t.Logf("✅ Created %d sessions", len(sessions))
	})

	t.Log("🎉 All basic session management tests passed!")
}
