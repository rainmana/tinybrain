package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Standalone session model
type Session struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	TaskType    string                 `json:"task_type"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type SessionCreateRequest struct {
	Name        string                 `json:"name"`
	TaskType    string                 `json:"task_type"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Standalone session repository
type StandaloneSessionRepository struct {
	app *pocketbase.PocketBase
}

func NewStandaloneSessionRepository(app *pocketbase.PocketBase) *StandaloneSessionRepository {
	return &StandaloneSessionRepository{app: app}
}

func (r *StandaloneSessionRepository) CreateSession(ctx context.Context, req *SessionCreateRequest) (*Session, error) {
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

func (r *StandaloneSessionRepository) GetSession(ctx context.Context, id string) (*Session, error) {
	record, err := r.app.Dao().FindRecordById("sessions", id)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "not found" {
			return nil, fmt.Errorf("session with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get session record: %w", err)
	}

	return r.recordToSession(record), nil
}

func (r *StandaloneSessionRepository) recordToSession(record *pbmodels.Record) *Session {
	// Handle metadata safely
	var metadata map[string]interface{}
	if raw := record.Get("metadata"); raw != nil {
		if m, ok := raw.(map[string]interface{}); ok {
			metadata = m
		}
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &Session{
		ID:          record.Id,
		Name:        record.GetString("name"),
		TaskType:    record.GetString("task_type"),
		Status:      record.GetString("status"),
		Description: record.GetString("description"),
		Metadata:    metadata,
		CreatedAt:   record.Created.String(),
		UpdatedAt:   record.Updated.String(),
	}
}

// Create sessions collection
func createSessionsCollection() *pbmodels.Collection {
	collection := &pbmodels.Collection{
		Name:       "sessions",
		Type:       pbmodels.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "name",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: types.Pointer(1),
			Max: types.Pointer(255),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "task_type",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: types.Pointer(1),
			Max: types.Pointer(100),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "status",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: types.Pointer(1),
			Max: types.Pointer(50),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "description",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: types.Pointer(2000),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "metadata",
		Type: schema.FieldTypeJson,
	})

	return collection
}

func TestSessionManagementStandalone(t *testing.T) {
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
	collection := createSessionsCollection()
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatalf("Failed to create sessions collection: %v", err)
	}
	log.Printf("Created collection: %s", collection.Name)

	// Initialize session repository
	sessionRepo := NewStandaloneSessionRepository(app)

	// Test 1: Create a session
	t.Run("CreateSession", func(t *testing.T) {
		req := &SessionCreateRequest{
			Name:        "Security Assessment Test",
			TaskType:    "security_review",
			Description: "Testing TinyBrain v2.0 session creation",
			Metadata:    map[string]interface{}{"priority": "high", "client": "test-client"},
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

	// Test 2: Get session by ID
	t.Run("GetSession", func(t *testing.T) {
		// First create a session
		req := &SessionCreateRequest{
			Name:     "Get Session Test",
			TaskType: "vulnerability_analysis",
		}

		createdSession, err := sessionRepo.CreateSession(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Now retrieve it
		retrievedSession, err := sessionRepo.GetSession(context.Background(), createdSession.ID)
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

	// Test 3: Create multiple sessions
	t.Run("CreateMultipleSessions", func(t *testing.T) {
		sessions := []*SessionCreateRequest{
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

	t.Log("🎉 All standalone session management tests passed!")
}
