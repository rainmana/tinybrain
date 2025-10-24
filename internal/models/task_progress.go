package models

import "time"

// TaskProgress represents the progress of a specific task within a session
type TaskProgress struct {
	ID                 string    `json:"id"`
	SessionID          string    `json:"session_id"`
	TaskName           string    `json:"task_name"`
	Stage              string    `json:"stage"`               // e.g., "data_collection", "analysis", "report_generation"
	Status             string    `json:"status"`              // e.g., "pending", "in_progress", "completed", "failed"
	ProgressPercentage float32   `json:"progress_percentage"` // 0.0 to 100.0
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TaskProgressCreateRequest defines the structure for creating a new task progress entry
type TaskProgressCreateRequest struct {
	SessionID          string  `json:"session_id"`
	TaskName           string  `json:"task_name"`
	Stage              string  `json:"stage"`
	Status             string  `json:"status"`
	ProgressPercentage float32 `json:"progress_percentage"`
	Notes              string  `json:"notes,omitempty"`
}

// TaskProgressUpdateRequest defines the structure for updating an existing task progress entry
type TaskProgressUpdateRequest struct {
	Stage              *string  `json:"stage,omitempty"`
	Status             *string  `json:"status,omitempty"`
	ProgressPercentage *float32 `json:"progress_percentage,omitempty"`
	Notes              *string  `json:"notes,omitempty"`
}

// TaskProgressListRequest defines the structure for listing task progress entries
type TaskProgressListRequest struct {
	SessionID string `json:"session_id,omitempty"`
	TaskName  string `json:"task_name,omitempty"`
	Status    string `json:"status,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}
