package repository

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/models"
)

// SessionRepositoryV2 handles session data operations with PocketBase
type SessionRepositoryV2 struct {
	app *pocketbase.PocketBase
}

// NewSessionRepositoryV2 creates a new session repository
func NewSessionRepositoryV2(app *pocketbase.PocketBase) *SessionRepositoryV2 {
	return &SessionRepositoryV2{app: app}
}

// CreateSession creates a new session in PocketBase
func (r *SessionRepositoryV2) CreateSession(ctx context.Context, req *models.SessionCreateRequest) (*models.Session, error) {
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

// GetSession retrieves a session by ID from PocketBase
func (r *SessionRepositoryV2) GetSession(ctx context.Context, id string) (*models.Session, error) {
	record, err := r.app.Dao().FindRecordById("sessions", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("session with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get session record: %w", err)
	}

	return r.recordToSession(record), nil
}

// ListSessions lists sessions based on criteria
func (r *SessionRepositoryV2) ListSessions(ctx context.Context, req *models.SessionListRequest) ([]*models.Session, int, error) {
	// Build filter
	filter := ""
	if req.TaskType != "" {
		filter = fmt.Sprintf("task_type = '%s'", req.TaskType)
	}
	if req.Status != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("status = '%s'", req.Status)
	}
	// Note: Query field not available in SessionListRequest, skipping for now

	// Set defaults
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	records, err := r.app.Dao().FindRecordsByFilter(
		"sessions",
		filter,
		"",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list session records: %w", err)
	}

	sessions := make([]*models.Session, len(records))
	for i, record := range records {
		sessions[i] = r.recordToSession(record)
	}

	// For total count, we'll do a separate query
	// This is simplified - in production you'd want proper counting
	totalCount := len(records)

	return sessions, totalCount, nil
}

// UpdateSession updates an existing session in PocketBase
func (r *SessionRepositoryV2) UpdateSession(ctx context.Context, id string, req *models.SessionUpdateRequest) (*models.Session, error) {
	record, err := r.app.Dao().FindRecordById("sessions", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("session with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find session record for update: %w", err)
	}

	if req.Name != "" {
		record.Set("name", req.Name)
	}
	if req.Status != "" {
		record.Set("status", req.Status)
	}
	if req.Description != "" {
		record.Set("description", req.Description)
	}
	if req.Metadata != nil {
		record.Set("metadata", req.Metadata)
	}

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to update session record: %w", err)
	}

	return r.recordToSession(record), nil
}

// DeleteSession deletes a session by ID from PocketBase
func (r *SessionRepositoryV2) DeleteSession(ctx context.Context, id string) error {
	record, err := r.app.Dao().FindRecordById("sessions", id)
	if err != nil {
		if err.Error() == "not found" {
			return fmt.Errorf("session with ID %s not found", id)
		}
		return fmt.Errorf("failed to find session record for deletion: %w", err)
	}

	if err := r.app.Dao().DeleteRecord(record); err != nil {
		return fmt.Errorf("failed to delete session record: %w", err)
	}

	return nil
}

// recordToSession converts a PocketBase record to a Session model
func (r *SessionRepositoryV2) recordToSession(record *pbmodels.Record) *models.Session {
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
