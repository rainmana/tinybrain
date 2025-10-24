package models

import "time"

// Memory represents a stored memory entry for security assessments
type Memory struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	Priority    int       `json:"priority"`   // 1-10 scale
	Confidence  float32   `json:"confidence"` // 0.0-1.0 scale
	Tags        []string  `json:"tags"`
	Source      string    `json:"source"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemoryCreateRequest defines the structure for creating a new memory entry
type MemoryCreateRequest struct {
	SessionID   string   `json:"session_id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Category    string   `json:"category"`
	Priority    int      `json:"priority"`
	Confidence  float32  `json:"confidence"`
	Tags        []string `json:"tags"`
	Source      string   `json:"source"`
	ContentType string   `json:"content_type"`
}

// MemoryUpdateRequest defines the structure for updating an existing memory entry
type MemoryUpdateRequest struct {
	Title       *string  `json:"title,omitempty"`
	Content     *string  `json:"content,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Confidence  *float32 `json:"confidence,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Source      *string  `json:"source,omitempty"`
	ContentType *string  `json:"content_type,omitempty"`
}

// MemorySearchRequest defines the structure for searching memory entries
type MemorySearchRequest struct {
	SessionID string   `json:"session_id"`
	Query     string   `json:"query,omitempty"`
	Category  string   `json:"category,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Source    string   `json:"source,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}
