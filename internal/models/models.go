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
	// Enhanced intelligence fields
	IntelligenceType string `json:"intelligence_type" db:"intelligence_type"` // osint, humint, sigint, geoint, masint, techint, finint, cybint
	TargetScope      string `json:"target_scope" db:"target_scope"`           // individual, organization, infrastructure, campaign
	Classification   string `json:"classification" db:"classification"`       // unclassified, confidential, secret, top_secret
	ThreatLevel      string `json:"threat_level" db:"threat_level"`           // low, medium, high, critical
	GeographicScope  string `json:"geographic_scope" db:"geographic_scope"`   // local, regional, national, international
}

// MemoryEntry represents a piece of information stored in the memory system
type MemoryEntry struct {
	ID          string    `json:"id" db:"id"`
	SessionID   string    `json:"session_id" db:"session_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	ContentType string    `json:"content_type" db:"content_type"`
	Category    string    `json:"category" db:"category"`
	Priority    int       `json:"priority" db:"priority"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	Tags        []string  `json:"tags" db:"tags"`
	Source      string    `json:"source" db:"source"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	AccessedAt  time.Time `json:"accessed_at" db:"accessed_at"`
	AccessCount int       `json:"access_count" db:"access_count"`
	// Enhanced intelligence fields
	IntelligenceType string                 `json:"intelligence_type" db:"intelligence_type"` // osint, humint, sigint, geoint, masint, techint, finint, cybint
	Classification   string                 `json:"classification" db:"classification"`       // unclassified, confidential, secret, top_secret
	ThreatLevel      string                 `json:"threat_level" db:"threat_level"`           // low, medium, high, critical
	GeographicScope  string                 `json:"geographic_scope" db:"geographic_scope"`   // local, regional, national, international
	Attribution      string                 `json:"attribution" db:"attribution"`             // threat actor attribution
	IOCType          string                 `json:"ioc_type" db:"ioc_type"`                   // ip, domain, url, hash, email, file
	IOCValue         string                 `json:"ioc_value" db:"ioc_value"`                 // actual IOC value
	MITRETactic      string                 `json:"mitre_tactic" db:"mitre_tactic"`           // MITRE ATT&CK tactic
	MITRETechnique   string                 `json:"mitre_technique" db:"mitre_technique"`     // MITRE ATT&CK technique
	MITREProcedure   string                 `json:"mitre_procedure" db:"mitre_procedure"`     // MITRE ATT&CK procedure
	KillChainPhase   string                 `json:"kill_chain_phase" db:"kill_chain_phase"`   // reconnaissance, weaponization, delivery, exploitation, installation, c2, actions
	RiskScore        float64                `json:"risk_score" db:"risk_score"`               // calculated risk score
	ImpactScore      float64                `json:"impact_score" db:"impact_score"`           // calculated impact score
	LikelihoodScore  float64                `json:"likelihood_score" db:"likelihood_score"`   // calculated likelihood score
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`                   // additional metadata
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
	// Enhanced intelligence relationship fields
	CorrelationType string  `json:"correlation_type" db:"correlation_type"` // temporal, spatial, logical, statistical, causal, predictive
	Confidence      float64 `json:"confidence" db:"confidence"`             // confidence in the relationship
	Evidence        string  `json:"evidence" db:"evidence"`                 // evidence supporting the relationship
	Direction       string  `json:"direction" db:"direction"`               // unidirectional, bidirectional
	Weight          float64 `json:"weight" db:"weight"`                     // relationship weight for analysis
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
	ID                 string     `json:"id" db:"id"`
	SessionID          string     `json:"session_id" db:"session_id"`
	TaskName           string     `json:"task_name" db:"task_name"`
	Stage              string     `json:"stage" db:"stage"`
	Status             string     `json:"status" db:"status"`
	ProgressPercentage int        `json:"progress_percentage" db:"progress_percentage"`
	Notes              string     `json:"notes" db:"notes"`
	StartedAt          *time.Time `json:"started_at" db:"started_at"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// SearchRequest represents a search query with filters
type SearchRequest struct {
	Query         string   `json:"query"`
	SessionID     string   `json:"session_id,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	MinPriority   int      `json:"min_priority,omitempty"`
	MinConfidence float64  `json:"min_confidence,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	Offset        int      `json:"offset,omitempty"`
	SearchType    string   `json:"search_type,omitempty"` // semantic, exact, fuzzy, tag, category, relationship
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	MemoryEntry MemoryEntry `json:"memory_entry"`
	Relevance   float64     `json:"relevance"`
	Highlights  []string    `json:"highlights,omitempty"`
}

// MemorySummary represents a summary of relevant memories for context
type MemorySummary struct {
	SessionID        string         `json:"session_id"`
	RelevantMemories []MemoryEntry  `json:"relevant_memories"`
	RelatedTasks     []TaskProgress `json:"related_tasks"`
	KeyFindings      []MemoryEntry  `json:"key_findings"`
	Summary          string         `json:"summary"`
	GeneratedAt      time.Time      `json:"generated_at"`
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
	TableCounts        map[string]int           `json:"table_counts"`
	DatabaseSizeBytes  int64                    `json:"database_size_bytes"`
	TopAccessedEntries []map[string]interface{} `json:"top_accessed_entries"`
	RecentActivity     []map[string]interface{} `json:"recent_activity"`
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
	ID          string    `json:"id" db:"id"`
	SessionID   string    `json:"session_id" db:"session_id"`
	CWEID       string    `json:"cwe_id" db:"cwe_id"`
	CVEList     []string  `json:"cve_list" db:"cve_list"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	Source      string    `json:"source" db:"source"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RiskCorrelation represents correlation analysis between vulnerabilities
type RiskCorrelation struct {
	ID               string    `json:"id" db:"id"`
	SessionID        string    `json:"session_id" db:"session_id"`
	PrimaryVulnID    string    `json:"primary_vuln_id" db:"primary_vuln_id"`
	SecondaryVulnIDs []string  `json:"secondary_vuln_ids" db:"secondary_vuln_ids"`
	RiskMultiplier   float64   `json:"risk_multiplier" db:"risk_multiplier"`
	AttackChain      []string  `json:"attack_chain" db:"attack_chain"`
	BusinessImpact   string    `json:"business_impact" db:"business_impact"`
	Confidence       float64   `json:"confidence" db:"confidence"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ComplianceMapping represents mapping to security compliance standards
type ComplianceMapping struct {
	ID               string    `json:"id" db:"id"`
	SessionID        string    `json:"session_id" db:"session_id"`
	Standard         string    `json:"standard" db:"standard"`
	Requirement      string    `json:"requirement" db:"requirement"`
	VulnerabilityIDs []string  `json:"vulnerability_ids" db:"vulnerability_ids"`
	ComplianceScore  float64   `json:"compliance_score" db:"compliance_score"`
	GapAnalysis      []string  `json:"gap_analysis" db:"gap_analysis"`
	Recommendations  []string  `json:"recommendations" db:"recommendations"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// IntelligenceFinding represents an intelligence finding with enhanced metadata
type IntelligenceFinding struct {
	ID               string                 `json:"id" db:"id"`
	SessionID        string                 `json:"session_id" db:"session_id"`
	Title            string                 `json:"title" db:"title"`
	Description      string                 `json:"description" db:"description"`
	IntelligenceType string                 `json:"intelligence_type" db:"intelligence_type"`
	Classification   string                 `json:"classification" db:"classification"`
	ThreatLevel      string                 `json:"threat_level" db:"threat_level"`
	GeographicScope  string                 `json:"geographic_scope" db:"geographic_scope"`
	Attribution      string                 `json:"attribution" db:"attribution"`
	IOCType          string                 `json:"ioc_type" db:"ioc_type"`
	IOCValue         string                 `json:"ioc_value" db:"ioc_value"`
	MITRETactic      string                 `json:"mitre_tactic" db:"mitre_tactic"`
	MITRETechnique   string                 `json:"mitre_technique" db:"mitre_technique"`
	MITREProcedure   string                 `json:"mitre_procedure" db:"mitre_procedure"`
	KillChainPhase   string                 `json:"kill_chain_phase" db:"kill_chain_phase"`
	RiskScore        float64                `json:"risk_score" db:"risk_score"`
	ImpactScore      float64                `json:"impact_score" db:"impact_score"`
	LikelihoodScore  float64                `json:"likelihood_score" db:"likelihood_score"`
	Confidence       float64                `json:"confidence" db:"confidence"`
	Priority         int                    `json:"priority" db:"priority"`
	Tags             []string               `json:"tags" db:"tags"`
	Source           string                 `json:"source" db:"source"`
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// ThreatActor represents a threat actor or group
type ThreatActor struct {
	ID              string                 `json:"id" db:"id"`
	SessionID       string                 `json:"session_id" db:"session_id"`
	Name            string                 `json:"name" db:"name"`
	Aliases         []string               `json:"aliases" db:"aliases"`
	Description     string                 `json:"description" db:"description"`
	Motivation      string                 `json:"motivation" db:"motivation"`
	Capabilities    []string               `json:"capabilities" db:"capabilities"`
	Targets         []string               `json:"targets" db:"targets"`
	Tools           []string               `json:"tools" db:"tools"`
	Techniques      []string               `json:"techniques" db:"techniques"`
	Attribution     string                 `json:"attribution" db:"attribution"`
	Confidence      float64                `json:"confidence" db:"confidence"`
	ThreatLevel     string                 `json:"threat_level" db:"threat_level"`
	GeographicScope string                 `json:"geographic_scope" db:"geographic_scope"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// AttackCampaign represents an attack campaign or operation
type AttackCampaign struct {
	ID              string                 `json:"id" db:"id"`
	SessionID       string                 `json:"session_id" db:"session_id"`
	Name            string                 `json:"name" db:"name"`
	Description     string                 `json:"description" db:"description"`
	ThreatActors    []string               `json:"threat_actors" db:"threat_actors"`
	Targets         []string               `json:"targets" db:"targets"`
	Techniques      []string               `json:"techniques" db:"techniques"`
	Tools           []string               `json:"tools" db:"tools"`
	IOCs            []string               `json:"iocs" db:"iocs"`
	StartDate       time.Time              `json:"start_date" db:"start_date"`
	EndDate         *time.Time             `json:"end_date" db:"end_date"`
	Status          string                 `json:"status" db:"status"`
	ThreatLevel     string                 `json:"threat_level" db:"threat_level"`
	GeographicScope string                 `json:"geographic_scope" db:"geographic_scope"`
	Confidence      float64                `json:"confidence" db:"confidence"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// IndicatorOfCompromise represents an IOC with enhanced metadata
type IndicatorOfCompromise struct {
	ID          string                 `json:"id" db:"id"`
	SessionID   string                 `json:"session_id" db:"session_id"`
	Type        string                 `json:"type" db:"type"` // ip, domain, url, hash, email, file
	Value       string                 `json:"value" db:"value"`
	Description string                 `json:"description" db:"description"`
	ThreatLevel string                 `json:"threat_level" db:"threat_level"`
	Confidence  float64                `json:"confidence" db:"confidence"`
	FirstSeen   time.Time              `json:"first_seen" db:"first_seen"`
	LastSeen    time.Time              `json:"last_seen" db:"last_seen"`
	Source      string                 `json:"source" db:"source"`
	Attribution string                 `json:"attribution" db:"attribution"`
	Campaigns   []string               `json:"campaigns" db:"campaigns"`
	Techniques  []string               `json:"techniques" db:"techniques"`
	Tags        []string               `json:"tags" db:"tags"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// Pattern represents a security pattern or behavior
type Pattern struct {
	ID          string                 `json:"id" db:"id"`
	SessionID   string                 `json:"session_id" db:"session_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	PatternType string                 `json:"pattern_type" db:"pattern_type"` // behavioral, attack, defense, vulnerability, exploit, ioc, ttp
	Category    string                 `json:"category" db:"category"`
	Severity    string                 `json:"severity" db:"severity"`
	Confidence  float64                `json:"confidence" db:"confidence"`
	Frequency   int                    `json:"frequency" db:"frequency"`
	Examples    []string               `json:"examples" db:"examples"`
	Mitigations []string               `json:"mitigations" db:"mitigations"`
	Detections  []string               `json:"detections" db:"detections"`
	Tags        []string               `json:"tags" db:"tags"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// Correlation represents a correlation between findings
type Correlation struct {
	ID              string                 `json:"id" db:"id"`
	SessionID       string                 `json:"session_id" db:"session_id"`
	SourceFindingID string                 `json:"source_finding_id" db:"source_finding_id"`
	TargetFindingID string                 `json:"target_finding_id" db:"target_finding_id"`
	CorrelationType string                 `json:"correlation_type" db:"correlation_type"` // temporal, spatial, logical, statistical, causal, predictive
	Strength        float64                `json:"strength" db:"strength"`
	Confidence      float64                `json:"confidence" db:"confidence"`
	Evidence        string                 `json:"evidence" db:"evidence"`
	Description     string                 `json:"description" db:"description"`
	Weight          float64                `json:"weight" db:"weight"`
	Direction       string                 `json:"direction" db:"direction"` // unidirectional, bidirectional
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// Request/Response types for new features
type MapToCVERequest struct {
	SessionID string `json:"session_id"`
	CWEID     string `json:"cwe_id"`
}

// Intelligence-specific request/response types
type CreateIntelligenceFindingRequest struct {
	SessionID        string                 `json:"session_id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	IntelligenceType string                 `json:"intelligence_type"`
	Classification   string                 `json:"classification"`
	ThreatLevel      string                 `json:"threat_level"`
	GeographicScope  string                 `json:"geographic_scope"`
	Attribution      string                 `json:"attribution"`
	IOCType          string                 `json:"ioc_type"`
	IOCValue         string                 `json:"ioc_value"`
	MITRETactic      string                 `json:"mitre_tactic"`
	MITRETechnique   string                 `json:"mitre_technique"`
	MITREProcedure   string                 `json:"mitre_procedure"`
	KillChainPhase   string                 `json:"kill_chain_phase"`
	RiskScore        float64                `json:"risk_score"`
	ImpactScore      float64                `json:"impact_score"`
	LikelihoodScore  float64                `json:"likelihood_score"`
	Confidence       float64                `json:"confidence"`
	Priority         int                    `json:"priority"`
	Tags             []string               `json:"tags"`
	Source           string                 `json:"source"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type CreateThreatActorRequest struct {
	SessionID       string                 `json:"session_id"`
	Name            string                 `json:"name"`
	Aliases         []string               `json:"aliases"`
	Description     string                 `json:"description"`
	Motivation      string                 `json:"motivation"`
	Capabilities    []string               `json:"capabilities"`
	Targets         []string               `json:"targets"`
	Tools           []string               `json:"tools"`
	Techniques      []string               `json:"techniques"`
	Attribution     string                 `json:"attribution"`
	Confidence      float64                `json:"confidence"`
	ThreatLevel     string                 `json:"threat_level"`
	GeographicScope string                 `json:"geographic_scope"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type CreateAttackCampaignRequest struct {
	SessionID       string                 `json:"session_id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ThreatActors    []string               `json:"threat_actors"`
	Targets         []string               `json:"targets"`
	Techniques      []string               `json:"techniques"`
	Tools           []string               `json:"tools"`
	IOCs            []string               `json:"iocs"`
	StartDate       time.Time              `json:"start_date"`
	EndDate         *time.Time             `json:"end_date"`
	Status          string                 `json:"status"`
	ThreatLevel     string                 `json:"threat_level"`
	GeographicScope string                 `json:"geographic_scope"`
	Confidence      float64                `json:"confidence"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type CreateIOCRequest struct {
	SessionID   string                 `json:"session_id"`
	Type        string                 `json:"type"`
	Value       string                 `json:"value"`
	Description string                 `json:"description"`
	ThreatLevel string                 `json:"threat_level"`
	Confidence  float64                `json:"confidence"`
	FirstSeen   time.Time              `json:"first_seen"`
	LastSeen    time.Time              `json:"last_seen"`
	Source      string                 `json:"source"`
	Attribution string                 `json:"attribution"`
	Campaigns   []string               `json:"campaigns"`
	Techniques  []string               `json:"techniques"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreatePatternRequest struct {
	SessionID   string                 `json:"session_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	PatternType string                 `json:"pattern_type"`
	Category    string                 `json:"category"`
	Severity    string                 `json:"severity"`
	Confidence  float64                `json:"confidence"`
	Frequency   int                    `json:"frequency"`
	Examples    []string               `json:"examples"`
	Mitigations []string               `json:"mitigations"`
	Detections  []string               `json:"detections"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateCorrelationRequest struct {
	SessionID       string                 `json:"session_id"`
	SourceFindingID string                 `json:"source_finding_id"`
	TargetFindingID string                 `json:"target_finding_id"`
	CorrelationType string                 `json:"correlation_type"`
	Strength        float64                `json:"strength"`
	Confidence      float64                `json:"confidence"`
	Evidence        string                 `json:"evidence"`
	Description     string                 `json:"description"`
	Weight          float64                `json:"weight"`
	Direction       string                 `json:"direction"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type IntelligenceAnalysisRequest struct {
	SessionID         string   `json:"session_id"`
	AnalysisType      string   `json:"analysis_type"`      // threat_intelligence, attack_analysis, pattern_recognition, correlation_analysis
	TargetScope       string   `json:"target_scope"`       // individual, organization, infrastructure, campaign
	IntelligenceTypes []string `json:"intelligence_types"` // osint, humint, sigint, etc.
	TimeRange         struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"time_range"`
	Filters map[string]interface{} `json:"filters"`
}

type IntelligenceAnalysisResponse struct {
	AnalysisID      string                  `json:"analysis_id"`
	SessionID       string                  `json:"session_id"`
	AnalysisType    string                  `json:"analysis_type"`
	Findings        []IntelligenceFinding   `json:"findings"`
	ThreatActors    []ThreatActor           `json:"threat_actors"`
	AttackCampaigns []AttackCampaign        `json:"attack_campaigns"`
	IOCs            []IndicatorOfCompromise `json:"iocs"`
	Patterns        []Pattern               `json:"patterns"`
	Correlations    []Correlation           `json:"correlations"`
	Insights        []string                `json:"insights"`
	Recommendations []string                `json:"recommendations"`
	RiskAssessment  map[string]interface{}  `json:"risk_assessment"`
	Confidence      float64                 `json:"confidence"`
	CreatedAt       time.Time               `json:"created_at"`
}

type MITREAttackMappingRequest struct {
	SessionID            string   `json:"session_id"`
	TechniqueIDs         []string `json:"technique_ids"`
	TacticIDs            []string `json:"tactic_ids"`
	ProcedureIDs         []string `json:"procedure_ids"`
	IncludeSubtechniques bool     `json:"include_subtechniques"`
}

type MITREAttackMappingResponse struct {
	SessionID  string                   `json:"session_id"`
	Tactics    []map[string]interface{} `json:"tactics"`
	Techniques []map[string]interface{} `json:"techniques"`
	Procedures []map[string]interface{} `json:"procedures"`
	Mappings   []map[string]interface{} `json:"mappings"`
	CreatedAt  time.Time                `json:"created_at"`
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
