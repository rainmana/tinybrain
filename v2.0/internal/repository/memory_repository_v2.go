package repository

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/models"
)

// MemoryRepositoryV2 handles memory data operations with PocketBase
type MemoryRepositoryV2 struct {
	app *pocketbase.PocketBase
}

// NewMemoryRepositoryV2 creates a new memory repository
func NewMemoryRepositoryV2(app *pocketbase.PocketBase) *MemoryRepositoryV2 {
	return &MemoryRepositoryV2{app: app}
}

// StoreMemory creates a new memory entry in PocketBase
func (r *MemoryRepositoryV2) StoreMemory(ctx context.Context, req *models.MemoryCreateRequest) (*models.Memory, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("memory_entries")
	if err != nil {
		return nil, fmt.Errorf("failed to find memory_entries collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("session_id", req.SessionID)
	record.Set("title", req.Title)
	record.Set("content", req.Content)
	record.Set("category", req.Category)
	record.Set("priority", req.Priority)
	record.Set("confidence", req.Confidence)
	record.Set("tags", req.Tags)
	record.Set("source", req.Source)
	record.Set("content_type", req.ContentType)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create memory record: %w", err)
	}

	return r.recordToMemory(record), nil
}

// GetMemory retrieves a memory entry by ID from PocketBase
func (r *MemoryRepositoryV2) GetMemory(ctx context.Context, id string) (*models.Memory, error) {
	record, err := r.app.Dao().FindRecordById("memory_entries", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("memory with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get memory record: %w", err)
	}

	return r.recordToMemory(record), nil
}

// SearchMemories searches for memories based on criteria
func (r *MemoryRepositoryV2) SearchMemories(ctx context.Context, req *models.MemorySearchRequest) ([]*models.Memory, int, error) {
	// Build filter
	filter := fmt.Sprintf("session_id = '%s'", req.SessionID)

	if req.Query != "" {
		filter += fmt.Sprintf(" && (title ~ '%s' || content ~ '%s')", req.Query, req.Query)
	}
	if req.Category != "" {
		filter += fmt.Sprintf(" && category = '%s'", req.Category)
	}
	if req.Source != "" {
		filter += fmt.Sprintf(" && source = '%s'", req.Source)
	}

	// Set defaults
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	records, err := r.app.Dao().FindRecordsByFilter(
		"memory_entries",
		filter,
		"",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search memory records: %w", err)
	}

	memories := make([]*models.Memory, len(records))
	for i, record := range records {
		memories[i] = r.recordToMemory(record)
	}

	// For total count, we'll do a separate query
	// This is simplified - in production you'd want proper counting
	totalCount := len(records)

	return memories, totalCount, nil
}

// UpdateMemory updates an existing memory entry in PocketBase
func (r *MemoryRepositoryV2) UpdateMemory(ctx context.Context, id string, req *models.MemoryUpdateRequest) (*models.Memory, error) {
	record, err := r.app.Dao().FindRecordById("memory_entries", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("memory with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find memory record for update: %w", err)
	}

	if req.Title != nil {
		record.Set("title", *req.Title)
	}
	if req.Content != nil {
		record.Set("content", *req.Content)
	}
	if req.Category != nil {
		record.Set("category", *req.Category)
	}
	if req.Priority != nil {
		record.Set("priority", *req.Priority)
	}
	if req.Confidence != nil {
		record.Set("confidence", *req.Confidence)
	}
	if req.Tags != nil {
		record.Set("tags", req.Tags)
	}
	if req.Source != nil {
		record.Set("source", *req.Source)
	}
	if req.ContentType != nil {
		record.Set("content_type", *req.ContentType)
	}

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to update memory record: %w", err)
	}

	return r.recordToMemory(record), nil
}

// DeleteMemory deletes a memory entry by ID from PocketBase
func (r *MemoryRepositoryV2) DeleteMemory(ctx context.Context, id string) error {
	record, err := r.app.Dao().FindRecordById("memory_entries", id)
	if err != nil {
		if err.Error() == "not found" {
			return fmt.Errorf("memory with ID %s not found", id)
		}
		return fmt.Errorf("failed to find memory record for deletion: %w", err)
	}

	if err := r.app.Dao().DeleteRecord(record); err != nil {
		return fmt.Errorf("failed to delete memory record: %w", err)
	}

	return nil
}

// recordToMemory converts a PocketBase record to a Memory model
func (r *MemoryRepositoryV2) recordToMemory(record *pbmodels.Record) *models.Memory {
	// Handle tags safely
	var tags []string
	if raw := record.Get("tags"); raw != nil {
		if tagSlice, ok := raw.([]interface{}); ok {
			for _, v := range tagSlice {
				if s, ok := v.(string); ok {
					tags = append(tags, s)
				}
			}
		}
	}
	if tags == nil {
		tags = make([]string, 0)
	}

	return &models.Memory{
		ID:          record.Id,
		SessionID:   record.GetString("session_id"),
		Title:       record.GetString("title"),
		Content:     record.GetString("content"),
		Category:    record.GetString("category"),
		Priority:    int(record.GetFloat("priority")),
		Confidence:  float32(record.GetFloat("confidence")),
		Tags:        tags,
		Source:      record.GetString("source"),
		ContentType: record.GetString("content_type"),
		CreatedAt:   record.Created.Time(),
		UpdatedAt:   record.Updated.Time(),
	}
}
