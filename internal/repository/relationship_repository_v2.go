package repository

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/models"
)

// RelationshipRepositoryV2 handles relationship data operations with PocketBase
type RelationshipRepositoryV2 struct {
	app *pocketbase.PocketBase
}

// NewRelationshipRepositoryV2 creates a new relationship repository
func NewRelationshipRepositoryV2(app *pocketbase.PocketBase) *RelationshipRepositoryV2 {
	return &RelationshipRepositoryV2{app: app}
}

// CreateRelationship creates a new relationship in PocketBase
func (r *RelationshipRepositoryV2) CreateRelationship(ctx context.Context, req *models.RelationshipCreateRequest) (*models.Relationship, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("relationships")
	if err != nil {
		return nil, fmt.Errorf("failed to find relationships collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("source_id", req.SourceID)
	record.Set("target_id", req.TargetID)
	record.Set("type", string(req.Type))
	record.Set("strength", req.Strength)
	record.Set("description", req.Description)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create relationship record: %w", err)
	}

	return r.recordToRelationship(record), nil
}

// GetRelationship retrieves a relationship by ID from PocketBase
func (r *RelationshipRepositoryV2) GetRelationship(ctx context.Context, id string) (*models.Relationship, error) {
	record, err := r.app.Dao().FindRecordById("relationships", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("relationship with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get relationship record: %w", err)
	}

	return r.recordToRelationship(record), nil
}

// ListRelationships lists relationships based on criteria
func (r *RelationshipRepositoryV2) ListRelationships(ctx context.Context, req *models.RelationshipListRequest) ([]*models.Relationship, int, error) {
	// Build filter
	filter := ""
	if req.SourceID != "" {
		filter = fmt.Sprintf("source_id = '%s'", req.SourceID)
	}
	if req.TargetID != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("target_id = '%s'", req.TargetID)
	}
	if req.Type != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("type = '%s'", req.Type)
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
		"relationships",
		filter,
		"",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list relationship records: %w", err)
	}

	relationships := make([]*models.Relationship, len(records))
	for i, record := range records {
		relationships[i] = r.recordToRelationship(record)
	}

	// For total count, we'll do a separate query
	// This is simplified - in production you'd want proper counting
	totalCount := len(records)

	return relationships, totalCount, nil
}

// UpdateRelationship updates an existing relationship in PocketBase
func (r *RelationshipRepositoryV2) UpdateRelationship(ctx context.Context, id string, req *models.RelationshipUpdateRequest) (*models.Relationship, error) {
	record, err := r.app.Dao().FindRecordById("relationships", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("relationship with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find relationship record for update: %w", err)
	}

	if req.Type != nil {
		record.Set("type", string(*req.Type))
	}
	if req.Strength != nil {
		record.Set("strength", *req.Strength)
	}
	if req.Description != nil {
		record.Set("description", *req.Description)
	}

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to update relationship record: %w", err)
	}

	return r.recordToRelationship(record), nil
}

// DeleteRelationship deletes a relationship by ID from PocketBase
func (r *RelationshipRepositoryV2) DeleteRelationship(ctx context.Context, id string) error {
	record, err := r.app.Dao().FindRecordById("relationships", id)
	if err != nil {
		if err.Error() == "not found" {
			return fmt.Errorf("relationship with ID %s not found", id)
		}
		return fmt.Errorf("failed to find relationship record for deletion: %w", err)
	}

	if err := r.app.Dao().DeleteRecord(record); err != nil {
		return fmt.Errorf("failed to delete relationship record: %w", err)
	}

	return nil
}

// recordToRelationship converts a PocketBase record to a Relationship model
func (r *RelationshipRepositoryV2) recordToRelationship(record *pbmodels.Record) *models.Relationship {
	return &models.Relationship{
		ID:          record.Id,
		SourceID:    record.GetString("source_id"),
		TargetID:    record.GetString("target_id"),
		Type:        models.RelationshipType(record.GetString("type")),
		Strength:    float32(record.GetFloat("strength")),
		Description: record.GetString("description"),
		CreatedAt:   record.Created.Time(),
		UpdatedAt:   record.Updated.Time(),
	}
}
