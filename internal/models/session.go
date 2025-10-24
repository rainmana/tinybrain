package models

import "time"

// Session represents a security assessment session
type Session struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	TaskType    string                 `json:"task_type"` // security_review, penetration_test, etc.
	Status      string                 `json:"status"`    // active, paused, completed, archived
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"` // JSON object for additional data
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// SessionCreateRequest represents the request to create a new session
type SessionCreateRequest struct {
	Name        string                 `json:"name"`
	TaskType    string                 `json:"task_type"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SessionUpdateRequest represents the request to update a session
type SessionUpdateRequest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SessionListRequest represents the request to list sessions
type SessionListRequest struct {
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Status   string `json:"status,omitempty"`
	TaskType string `json:"task_type,omitempty"`
}
