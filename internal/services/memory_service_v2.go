package services

import (
	"context"
	"fmt"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// MemoryServiceV2 provides business logic for memory management
type MemoryServiceV2 struct {
	memoryRepo *repository.MemoryRepositoryV2
}

// NewMemoryServiceV2 creates a new MemoryService
func NewMemoryServiceV2(memoryRepo *repository.MemoryRepositoryV2) *MemoryServiceV2 {
	return &MemoryServiceV2{memoryRepo: memoryRepo}
}

// StoreMemory stores a new memory entry
func (s *MemoryServiceV2) StoreMemory(ctx context.Context, req *models.MemoryCreateRequest) (*models.Memory, error) {
	if req.SessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("memory title cannot be empty")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("memory content cannot be empty")
	}
	if req.Category == "" {
		return nil, fmt.Errorf("memory category cannot be empty")
	}
	if req.Priority < 1 || req.Priority > 10 {
		return nil, fmt.Errorf("memory priority must be between 1 and 10")
	}
	if req.Confidence < 0.0 || req.Confidence > 1.0 {
		return nil, fmt.Errorf("memory confidence must be between 0.0 and 1.0")
	}

	return s.memoryRepo.StoreMemory(ctx, req)
}

// GetMemory retrieves a memory entry by ID
func (s *MemoryServiceV2) GetMemory(ctx context.Context, id string) (*models.Memory, error) {
	if id == "" {
		return nil, fmt.Errorf("memory ID cannot be empty")
	}

	return s.memoryRepo.GetMemory(ctx, id)
}

// SearchMemories searches for memories based on criteria
func (s *MemoryServiceV2) SearchMemories(ctx context.Context, req *models.MemorySearchRequest) ([]*models.Memory, int, error) {
	if req.SessionID == "" {
		return nil, 0, fmt.Errorf("session ID cannot be empty")
	}

	return s.memoryRepo.SearchMemories(ctx, req)
}

// UpdateMemory updates an existing memory entry
func (s *MemoryServiceV2) UpdateMemory(ctx context.Context, id string, req *models.MemoryUpdateRequest) (*models.Memory, error) {
	if id == "" {
		return nil, fmt.Errorf("memory ID cannot be empty")
	}

	return s.memoryRepo.UpdateMemory(ctx, id, req)
}

// DeleteMemory deletes a memory entry by ID
func (s *MemoryServiceV2) DeleteMemory(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("memory ID cannot be empty")
	}

	return s.memoryRepo.DeleteMemory(ctx, id)
}
