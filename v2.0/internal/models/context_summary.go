package models

// ContextSummary represents a summary of context for a session
type ContextSummary struct {
	SessionID   string             `json:"session_id"`
	Snapshots   []*ContextSnapshot `json:"snapshots"`
	TotalCount  int                `json:"total_count"`
	MaxMemories int                `json:"max_memories"`
}
