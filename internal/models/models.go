package models

import (
	"encoding/json"
	"time"
)

// Session represents a security-focused LLM interaction session
type Session struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	TaskType    string                 `json:"task_type" db:"task_type"`
	Status      string                 `json:"status" db:"status"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
}

// MemoryEntry represents a piece of information stored in the memory system
type MemoryEntry struct {
	ID           string                 `json:"id" db:"id"`
	SessionID    string                 `json:"session_id" db:"session_id"`
	Title        string                 `json:"title" db:"title"`
	Content      string                 `json:"content" db:"content"`
	ContentType  string                 `json:"content_type" db:"content_type"`
	Category     string                 `json:"category" db:"category"`
	Priority     int                    `json:"priority" db:"priority"`
	Confidence   float64                `json:"confidence" db:"confidence"`
	Tags         []string               `json:"tags" db:"tags"`
	Source       string                 `json:"source" db:"source"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	AccessedAt   time.Time              `json:"accessed_at" db:"accessed_at"`
	AccessCount  int                    `json:"access_count" db:"access_count"`
}

// Relationship represents a connection between two memory entries
type Relationship struct {
	ID               string    `json:"id" db:"id"`
	SourceEntryID    string    `json:"source_entry_id" db:"source_entry_id"`
	TargetEntryID    string    `json:"target_entry_id" db:"target_entry_id"`
	RelationshipType string    `json:"relationship_type" db:"relationship_type"`
	Strength         float64   `json:"strength" db:"strength"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// ContextSnapshot represents a saved state of the LLM's context
type ContextSnapshot struct {
	ID            string                 `json:"id" db:"id"`
	SessionID     string                 `json:"session_id" db:"session_id"`
	Name          string                 `json:"name" db:"name"`
	Description   string                 `json:"description" db:"description"`
	ContextData   map[string]interface{} `json:"context_data" db:"context_data"`
	MemorySummary string                 `json:"memory_summary" db:"memory_summary"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// SearchHistory represents a record of search queries
type SearchHistory struct {
	ID           string    `json:"id" db:"id"`
	SessionID    string    `json:"session_id" db:"session_id"`
	Query        string    `json:"query" db:"query"`
	SearchType   string    `json:"search_type" db:"search_type"`
	ResultsCount int       `json:"results_count" db:"results_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// TaskProgress represents progress on a multi-stage task
type TaskProgress struct {
	ID                  string     `json:"id" db:"id"`
	SessionID           string     `json:"session_id" db:"session_id"`
	TaskName            string     `json:"task_name" db:"task_name"`
	Stage               string     `json:"stage" db:"stage"`
	Status              string     `json:"status" db:"status"`
	ProgressPercentage  int        `json:"progress_percentage" db:"progress_percentage"`
	Notes               string     `json:"notes" db:"notes"`
	StartedAt           *time.Time `json:"started_at" db:"started_at"`
	CompletedAt         *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// SearchRequest represents a search query with filters
type SearchRequest struct {
	Query        string   `json:"query"`
	SessionID    string   `json:"session_id,omitempty"`
	Categories   []string `json:"categories,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	MinPriority  int      `json:"min_priority,omitempty"`
	MinConfidence float64 `json:"min_confidence,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Offset       int      `json:"offset,omitempty"`
	SearchType   string   `json:"search_type,omitempty"` // semantic, exact, fuzzy, tag, category, relationship
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	MemoryEntry MemoryEntry `json:"memory_entry"`
	Relevance   float64     `json:"relevance"`
	Highlights  []string    `json:"highlights,omitempty"`
}

// MemorySummary represents a summary of relevant memories for context
type MemorySummary struct {
	SessionID     string         `json:"session_id"`
	RelevantMemories []MemoryEntry `json:"relevant_memories"`
	RelatedTasks  []TaskProgress `json:"related_tasks"`
	KeyFindings   []MemoryEntry  `json:"key_findings"`
	Summary       string         `json:"summary"`
	GeneratedAt   time.Time      `json:"generated_at"`
}

// CreateMemoryEntryRequest represents a request to create a new memory entry
type CreateMemoryEntryRequest struct {
	SessionID   string   `json:"session_id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type,omitempty"`
	Category    string   `json:"category"`
	Priority    int      `json:"priority,omitempty"`
	Confidence  float64  `json:"confidence,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Source      string   `json:"source,omitempty"`
}

// UpdateMemoryEntryRequest represents a request to update a memory entry
type UpdateMemoryEntryRequest struct {
	ID          string   `json:"id"`
	Title       string   `json:"title,omitempty"`
	Content     string   `json:"content,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Category    string   `json:"category,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Confidence  *float64 `json:"confidence,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Source      string   `json:"source,omitempty"`
}

// CreateRelationshipRequest represents a request to create a relationship
type CreateRelationshipRequest struct {
	SourceEntryID    string  `json:"source_entry_id"`
	TargetEntryID    string  `json:"target_entry_id"`
	RelationshipType string  `json:"relationship_type"`
	Strength         float64 `json:"strength,omitempty"`
	Description      string  `json:"description,omitempty"`
}

// DatabaseStats represents database statistics
type DatabaseStats struct {
	TableCounts        map[string]int                    `json:"table_counts"`
	DatabaseSizeBytes  int64                             `json:"database_size_bytes"`
	TopAccessedEntries []map[string]interface{}          `json:"top_accessed_entries"`
	RecentActivity     []map[string]interface{}          `json:"recent_activity"`
}

// Helper methods for JSON serialization of complex fields

// MarshalJSON custom marshaling for Session.Metadata
func (s *Session) MarshalJSON() ([]byte, error) {
	type Alias Session
	metadataBytes, err := json.Marshal(s.Metadata)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		*Alias
		Metadata json.RawMessage `json:"metadata"`
	}{
		Alias:    (*Alias)(s),
		Metadata: json.RawMessage(metadataBytes),
	})
}

// UnmarshalJSON custom unmarshaling for Session.Metadata
func (s *Session) UnmarshalJSON(data []byte) error {
	type Alias Session
	aux := &struct {
		*Alias
		Metadata json.RawMessage `json:"metadata"`
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Metadata != nil {
		if err := json.Unmarshal(aux.Metadata, &s.Metadata); err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON custom marshaling for MemoryEntry.Tags
func (m *MemoryEntry) MarshalJSON() ([]byte, error) {
	type Alias MemoryEntry
	tagsBytes, err := json.Marshal(m.Tags)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		*Alias
		Tags json.RawMessage `json:"tags"`
	}{
		Alias: (*Alias)(m),
		Tags:  json.RawMessage(tagsBytes),
	})
}

// UnmarshalJSON custom unmarshaling for MemoryEntry.Tags
func (m *MemoryEntry) UnmarshalJSON(data []byte) error {
	type Alias MemoryEntry
	aux := &struct {
		*Alias
		Tags json.RawMessage `json:"tags"`
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Tags != nil {
		if err := json.Unmarshal(aux.Tags, &m.Tags); err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON custom marshaling for ContextSnapshot.ContextData
func (c *ContextSnapshot) MarshalJSON() ([]byte, error) {
	type Alias ContextSnapshot
	contextDataBytes, err := json.Marshal(c.ContextData)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		*Alias
		ContextData json.RawMessage `json:"context_data"`
	}{
		Alias:       (*Alias)(c),
		ContextData: json.RawMessage(contextDataBytes),
	})
}

// UnmarshalJSON custom unmarshaling for ContextSnapshot.ContextData
func (c *ContextSnapshot) UnmarshalJSON(data []byte) error {
	type Alias ContextSnapshot
	aux := &struct {
		*Alias
		ContextData json.RawMessage `json:"context_data"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.ContextData != nil {
		if err := json.Unmarshal(aux.ContextData, &c.ContextData); err != nil {
			return err
		}
	}
	return nil
}

// CVEMapping represents a mapping between CWE and CVE entries
type CVEMapping struct {
	ID           string    `json:"id" db:"id"`
	SessionID    string    `json:"session_id" db:"session_id"`
	CWEID        string    `json:"cwe_id" db:"cwe_id"`
	CVEList      []string  `json:"cve_list" db:"cve_list"`
	LastUpdated  time.Time `json:"last_updated" db:"last_updated"`
	Confidence   float64   `json:"confidence" db:"confidence"`
	Source       string    `json:"source" db:"source"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RiskCorrelation represents correlation analysis between vulnerabilities
type RiskCorrelation struct {
	ID                string    `json:"id" db:"id"`
	SessionID         string    `json:"session_id" db:"session_id"`
	PrimaryVulnID     string    `json:"primary_vuln_id" db:"primary_vuln_id"`
	SecondaryVulnIDs  []string  `json:"secondary_vuln_ids" db:"secondary_vuln_ids"`
	RiskMultiplier    float64   `json:"risk_multiplier" db:"risk_multiplier"`
	AttackChain       []string  `json:"attack_chain" db:"attack_chain"`
	BusinessImpact    string    `json:"business_impact" db:"business_impact"`
	Confidence        float64   `json:"confidence" db:"confidence"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// ComplianceMapping represents mapping to security compliance standards
type ComplianceMapping struct {
	ID                string    `json:"id" db:"id"`
	SessionID         string    `json:"session_id" db:"session_id"`
	Standard          string    `json:"standard" db:"standard"`
	Requirement       string    `json:"requirement" db:"requirement"`
	VulnerabilityIDs  []string  `json:"vulnerability_ids" db:"vulnerability_ids"`
	ComplianceScore   float64   `json:"compliance_score" db:"compliance_score"`
	GapAnalysis       []string  `json:"gap_analysis" db:"gap_analysis"`
	Recommendations   []string  `json:"recommendations" db:"recommendations"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response types for new features
type MapToCVERequest struct {
	SessionID string `json:"session_id"`
	CWEID     string `json:"cwe_id"`
}

type MapToCVEResponse struct {
	CVEMapping *CVEMapping `json:"cve_mapping"`
	Error      string      `json:"error,omitempty"`
}

type AnalyzeRiskCorrelationRequest struct {
	SessionID string `json:"session_id"`
}

type AnalyzeRiskCorrelationResponse struct {
	Correlations []*RiskCorrelation `json:"correlations"`
	Error        string             `json:"error,omitempty"`
}

type MapToComplianceRequest struct {
	SessionID string `json:"session_id"`
	Standard  string `json:"standard"`
}

type MapToComplianceResponse struct {
	ComplianceMapping *ComplianceMapping `json:"compliance_mapping"`
	Error             string             `json:"error,omitempty"`
}
