package services

import (
	"context"
	"fmt"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// RelationshipServiceV2 provides business logic for relationship management
type RelationshipServiceV2 struct {
	relationshipRepo *repository.RelationshipRepositoryV2
}

// NewRelationshipServiceV2 creates a new RelationshipService
func NewRelationshipServiceV2(relationshipRepo *repository.RelationshipRepositoryV2) *RelationshipServiceV2 {
	return &RelationshipServiceV2{relationshipRepo: relationshipRepo}
}

// CreateRelationship creates a new relationship between memories
func (s *RelationshipServiceV2) CreateRelationship(ctx context.Context, req *models.RelationshipCreateRequest) (*models.Relationship, error) {
	if req.SourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}
	if req.TargetID == "" {
		return nil, fmt.Errorf("target ID cannot be empty")
	}
	if req.SourceID == req.TargetID {
		return nil, fmt.Errorf("source and target cannot be the same")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("relationship type cannot be empty")
	}
	if req.Strength < 0.0 || req.Strength > 1.0 {
		return nil, fmt.Errorf("relationship strength must be between 0.0 and 1.0")
	}

	return s.relationshipRepo.CreateRelationship(ctx, req)
}

// GetRelationship retrieves a relationship by ID
func (s *RelationshipServiceV2) GetRelationship(ctx context.Context, id string) (*models.Relationship, error) {
	if id == "" {
		return nil, fmt.Errorf("relationship ID cannot be empty")
	}

	return s.relationshipRepo.GetRelationship(ctx, id)
}

// ListRelationships lists relationships based on criteria
func (s *RelationshipServiceV2) ListRelationships(ctx context.Context, req *models.RelationshipListRequest) ([]*models.Relationship, int, error) {
	return s.relationshipRepo.ListRelationships(ctx, req)
}

// UpdateRelationship updates an existing relationship
func (s *RelationshipServiceV2) UpdateRelationship(ctx context.Context, id string, req *models.RelationshipUpdateRequest) (*models.Relationship, error) {
	if id == "" {
		return nil, fmt.Errorf("relationship ID cannot be empty")
	}

	return s.relationshipRepo.UpdateRelationship(ctx, id, req)
}

// DeleteRelationship deletes a relationship by ID
func (s *RelationshipServiceV2) DeleteRelationship(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("relationship ID cannot be empty")
	}

	return s.relationshipRepo.DeleteRelationship(ctx, id)
}
