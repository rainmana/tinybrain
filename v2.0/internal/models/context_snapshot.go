package models

import "time"

// ContextSnapshot represents a snapshot of the LLM's context at a specific point in time
type ContextSnapshot struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	Name        string                 `json:"name"`
	ContextData map[string]interface{} `json:"context_data"` // The actual context data (e.g., JSON representation of LLM state)
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ContextSnapshotCreateRequest defines the structure for creating a new context snapshot
type ContextSnapshotCreateRequest struct {
	SessionID   string                 `json:"session_id"`
	Name        string                 `json:"name"`
	ContextData map[string]interface{} `json:"context_data"`
	Description string                 `json:"description,omitempty"`
}

// ContextSnapshotUpdateRequest defines the structure for updating an existing context snapshot
type ContextSnapshotUpdateRequest struct {
	Name        *string                `json:"name,omitempty"`
	ContextData map[string]interface{} `json:"context_data,omitempty"`
	Description *string                `json:"description,omitempty"`
}

// ContextSnapshotListRequest defines the structure for listing context snapshots
type ContextSnapshotListRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Query     string `json:"query,omitempty"` // General search query
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}
