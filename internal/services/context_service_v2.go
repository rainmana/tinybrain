package services

import (
	"context"
	"fmt"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// ContextServiceV2 provides business logic for context snapshot management
type ContextServiceV2 struct {
	contextRepo *repository.ContextRepositoryV2
}

// NewContextServiceV2 creates a new ContextService
func NewContextServiceV2(contextRepo *repository.ContextRepositoryV2) *ContextServiceV2 {
	return &ContextServiceV2{contextRepo: contextRepo}
}

// CreateContextSnapshot creates a new context snapshot
func (s *ContextServiceV2) CreateContextSnapshot(ctx context.Context, req *models.ContextSnapshotCreateRequest) (*models.ContextSnapshot, error) {
	if req.SessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("context snapshot name cannot be empty")
	}
	if req.ContextData == nil {
		return nil, fmt.Errorf("context data cannot be empty")
	}

	return s.contextRepo.CreateContextSnapshot(ctx, req)
}

// GetContextSnapshot retrieves a context snapshot by ID
func (s *ContextServiceV2) GetContextSnapshot(ctx context.Context, id string) (*models.ContextSnapshot, error) {
	if id == "" {
		return nil, fmt.Errorf("context snapshot ID cannot be empty")
	}

	return s.contextRepo.GetContextSnapshot(ctx, id)
}

// ListContextSnapshots lists context snapshots based on criteria
func (s *ContextServiceV2) ListContextSnapshots(ctx context.Context, req *models.ContextSnapshotListRequest) ([]*models.ContextSnapshot, int, error) {
	return s.contextRepo.ListContextSnapshots(ctx, req)
}

// UpdateContextSnapshot updates an existing context snapshot
func (s *ContextServiceV2) UpdateContextSnapshot(ctx context.Context, id string, req *models.ContextSnapshotUpdateRequest) (*models.ContextSnapshot, error) {
	if id == "" {
		return nil, fmt.Errorf("context snapshot ID cannot be empty")
	}

	return s.contextRepo.UpdateContextSnapshot(ctx, id, req)
}

// DeleteContextSnapshot deletes a context snapshot by ID
func (s *ContextServiceV2) DeleteContextSnapshot(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("context snapshot ID cannot be empty")
	}

	return s.contextRepo.DeleteContextSnapshot(ctx, id)
}
