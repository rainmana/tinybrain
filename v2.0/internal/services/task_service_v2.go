package services

import (
	"context"
	"fmt"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// TaskServiceV2 provides business logic for task progress tracking
type TaskServiceV2 struct {
	taskRepo *repository.TaskRepositoryV2
}

// NewTaskServiceV2 creates a new TaskService
func NewTaskServiceV2(taskRepo *repository.TaskRepositoryV2) *TaskServiceV2 {
	return &TaskServiceV2{taskRepo: taskRepo}
}

// CreateTaskProgress creates a new task progress entry
func (s *TaskServiceV2) CreateTaskProgress(ctx context.Context, req *models.TaskProgressCreateRequest) (*models.TaskProgress, error) {
	if req.SessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}
	if req.TaskName == "" {
		return nil, fmt.Errorf("task name cannot be empty")
	}
	if req.Stage == "" {
		return nil, fmt.Errorf("task stage cannot be empty")
	}
	if req.Status == "" {
		return nil, fmt.Errorf("task status cannot be empty")
	}
	if req.ProgressPercentage < 0.0 || req.ProgressPercentage > 100.0 {
		return nil, fmt.Errorf("progress percentage must be between 0.0 and 100.0")
	}

	return s.taskRepo.CreateTaskProgress(ctx, req)
}

// GetTaskProgress retrieves a task progress entry by ID
func (s *TaskServiceV2) GetTaskProgress(ctx context.Context, id string) (*models.TaskProgress, error) {
	if id == "" {
		return nil, fmt.Errorf("task progress ID cannot be empty")
	}

	return s.taskRepo.GetTaskProgress(ctx, id)
}

// ListTaskProgress lists task progress entries based on criteria
func (s *TaskServiceV2) ListTaskProgress(ctx context.Context, req *models.TaskProgressListRequest) ([]*models.TaskProgress, int, error) {
	return s.taskRepo.ListTaskProgress(ctx, req)
}

// UpdateTaskProgress updates an existing task progress entry
func (s *TaskServiceV2) UpdateTaskProgress(ctx context.Context, id string, req *models.TaskProgressUpdateRequest) (*models.TaskProgress, error) {
	if id == "" {
		return nil, fmt.Errorf("task progress ID cannot be empty")
	}

	return s.taskRepo.UpdateTaskProgress(ctx, id, req)
}

// DeleteTaskProgress deletes a task progress entry by ID
func (s *TaskServiceV2) DeleteTaskProgress(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("task progress ID cannot be empty")
	}

	return s.taskRepo.DeleteTaskProgress(ctx, id)
}
