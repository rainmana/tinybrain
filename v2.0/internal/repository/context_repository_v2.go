package repository

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/models"
)

// ContextRepositoryV2 handles context snapshot data operations with PocketBase
type ContextRepositoryV2 struct {
	app *pocketbase.PocketBase
}

// NewContextRepositoryV2 creates a new context repository
func NewContextRepositoryV2(app *pocketbase.PocketBase) *ContextRepositoryV2 {
	return &ContextRepositoryV2{app: app}
}

// CreateContextSnapshot creates a new context snapshot in PocketBase
func (r *ContextRepositoryV2) CreateContextSnapshot(ctx context.Context, req *models.ContextSnapshotCreateRequest) (*models.ContextSnapshot, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("context_snapshots")
	if err != nil {
		return nil, fmt.Errorf("failed to find context_snapshots collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("session_id", req.SessionID)
	record.Set("name", req.Name)
	record.Set("context_data", req.ContextData)
	record.Set("description", req.Description)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create context snapshot record: %w", err)
	}

	return r.recordToContextSnapshot(record), nil
}

// GetContextSnapshot retrieves a context snapshot by ID from PocketBase
func (r *ContextRepositoryV2) GetContextSnapshot(ctx context.Context, id string) (*models.ContextSnapshot, error) {
	record, err := r.app.Dao().FindRecordById("context_snapshots", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("context snapshot with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get context snapshot record: %w", err)
	}

	return r.recordToContextSnapshot(record), nil
}

// ListContextSnapshots lists context snapshots based on criteria
func (r *ContextRepositoryV2) ListContextSnapshots(ctx context.Context, req *models.ContextSnapshotListRequest) ([]*models.ContextSnapshot, int, error) {
	filter := ""
	if req.SessionID != "" {
		filter = fmt.Sprintf("session_id = '%s'", req.SessionID)
	}
	if req.Query != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("name ~ '%s' || description ~ '%s'", req.Query, req.Query)
	}

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
		"context_snapshots",
		filter,
		"",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list context snapshot records: %w", err)
	}

	snapshots := make([]*models.ContextSnapshot, len(records))
	for i, record := range records {
		snapshots[i] = r.recordToContextSnapshot(record)
	}

	totalCount := len(records) // Simplified
	return snapshots, totalCount, nil
}

// UpdateContextSnapshot updates an existing context snapshot in PocketBase
func (r *ContextRepositoryV2) UpdateContextSnapshot(ctx context.Context, id string, req *models.ContextSnapshotUpdateRequest) (*models.ContextSnapshot, error) {
	record, err := r.app.Dao().FindRecordById("context_snapshots", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("context snapshot with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find context snapshot record for update: %w", err)
	}

	if req.Name != nil {
		record.Set("name", *req.Name)
	}
	if req.ContextData != nil {
		record.Set("context_data", req.ContextData)
	}
	if req.Description != nil {
		record.Set("description", *req.Description)
	}

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to update context snapshot record: %w", err)
	}

	return r.recordToContextSnapshot(record), nil
}

// DeleteContextSnapshot deletes a context snapshot by ID from PocketBase
func (r *ContextRepositoryV2) DeleteContextSnapshot(ctx context.Context, id string) error {
	record, err := r.app.Dao().FindRecordById("context_snapshots", id)
	if err != nil {
		if err.Error() == "not found" {
			return fmt.Errorf("context snapshot with ID %s not found", id)
		}
		return fmt.Errorf("failed to find context snapshot record for deletion: %w", err)
	}

	if err := r.app.Dao().DeleteRecord(record); err != nil {
		return fmt.Errorf("failed to delete context snapshot record: %w", err)
	}

	return nil
}

// recordToContextSnapshot converts a PocketBase record to a ContextSnapshot model
func (r *ContextRepositoryV2) recordToContextSnapshot(record *pbmodels.Record) *models.ContextSnapshot {
	// Handle context data safely
	var contextData map[string]interface{}
	if raw := record.Get("context_data"); raw != nil {
		if m, ok := raw.(map[string]interface{}); ok {
			contextData = m
		}
	}
	if contextData == nil {
		contextData = make(map[string]interface{})
	}

	return &models.ContextSnapshot{
		ID:          record.Id,
		SessionID:   record.GetString("session_id"),
		Name:        record.GetString("name"),
		ContextData: contextData,
		Description: record.GetString("description"),
		CreatedAt:   record.Created.Time(),
		UpdatedAt:   record.Updated.Time(),
	}
}
