package repository

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/models"
)

// TaskRepositoryV2 handles task progress data operations with PocketBase
type TaskRepositoryV2 struct {
	app *pocketbase.PocketBase
}

// NewTaskRepositoryV2 creates a new task repository
func NewTaskRepositoryV2(app *pocketbase.PocketBase) *TaskRepositoryV2 {
	return &TaskRepositoryV2{app: app}
}

// CreateTaskProgress creates a new task progress entry in PocketBase
func (r *TaskRepositoryV2) CreateTaskProgress(ctx context.Context, req *models.TaskProgressCreateRequest) (*models.TaskProgress, error) {
	collection, err := r.app.Dao().FindCollectionByNameOrId("task_progress")
	if err != nil {
		return nil, fmt.Errorf("failed to find task_progress collection: %w", err)
	}

	record := pbmodels.NewRecord(collection)
	record.Set("session_id", req.SessionID)
	record.Set("task_name", req.TaskName)
	record.Set("stage", req.Stage)
	record.Set("status", req.Status)
	record.Set("progress_percentage", req.ProgressPercentage)
	record.Set("notes", req.Notes)

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create task progress record: %w", err)
	}

	return r.recordToTaskProgress(record), nil
}

// GetTaskProgress retrieves a task progress entry by ID from PocketBase
func (r *TaskRepositoryV2) GetTaskProgress(ctx context.Context, id string) (*models.TaskProgress, error) {
	record, err := r.app.Dao().FindRecordById("task_progress", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("task progress with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get task progress record: %w", err)
	}

	return r.recordToTaskProgress(record), nil
}

// ListTaskProgress lists task progress entries based on criteria
func (r *TaskRepositoryV2) ListTaskProgress(ctx context.Context, req *models.TaskProgressListRequest) ([]*models.TaskProgress, int, error) {
	filter := ""
	if req.SessionID != "" {
		filter = fmt.Sprintf("session_id = '%s'", req.SessionID)
	}
	if req.TaskName != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("task_name = '%s'", req.TaskName)
	}
	if req.Status != "" {
		if filter != "" {
			filter += " && "
		}
		filter += fmt.Sprintf("status = '%s'", req.Status)
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
		"task_progress",
		filter,
		"",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list task progress records: %w", err)
	}

	progressEntries := make([]*models.TaskProgress, len(records))
	for i, record := range records {
		progressEntries[i] = r.recordToTaskProgress(record)
	}

	totalCount := len(records) // Simplified
	return progressEntries, totalCount, nil
}

// UpdateTaskProgress updates an existing task progress entry in PocketBase
func (r *TaskRepositoryV2) UpdateTaskProgress(ctx context.Context, id string, req *models.TaskProgressUpdateRequest) (*models.TaskProgress, error) {
	record, err := r.app.Dao().FindRecordById("task_progress", id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, fmt.Errorf("task progress with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find task progress record for update: %w", err)
	}

	if req.Stage != nil {
		record.Set("stage", *req.Stage)
	}
	if req.Status != nil {
		record.Set("status", *req.Status)
	}
	if req.ProgressPercentage != nil {
		record.Set("progress_percentage", *req.ProgressPercentage)
	}
	if req.Notes != nil {
		record.Set("notes", *req.Notes)
	}

	if err := r.app.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("failed to update task progress record: %w", err)
	}

	return r.recordToTaskProgress(record), nil
}

// DeleteTaskProgress deletes a task progress entry by ID from PocketBase
func (r *TaskRepositoryV2) DeleteTaskProgress(ctx context.Context, id string) error {
	record, err := r.app.Dao().FindRecordById("task_progress", id)
	if err != nil {
		if err.Error() == "not found" {
			return fmt.Errorf("task progress with ID %s not found", id)
		}
		return fmt.Errorf("failed to find task progress record for deletion: %w", err)
	}

	if err := r.app.Dao().DeleteRecord(record); err != nil {
		return fmt.Errorf("failed to delete task progress record: %w", err)
	}

	return nil
}

// recordToTaskProgress converts a PocketBase record to a TaskProgress model
func (r *TaskRepositoryV2) recordToTaskProgress(record *pbmodels.Record) *models.TaskProgress {
	return &models.TaskProgress{
		ID:                 record.Id,
		SessionID:          record.GetString("session_id"),
		TaskName:           record.GetString("task_name"),
		Stage:              record.GetString("stage"),
		Status:             record.GetString("status"),
		ProgressPercentage: float32(record.GetFloat("progress_percentage")),
		Notes:              record.GetString("notes"),
		CreatedAt:          record.Created.Time(),
		UpdatedAt:          record.Updated.Time(),
	}
}
