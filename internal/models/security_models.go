package models

import (
	"encoding/json"
	"time"
)

// NVDCVE represents a CVE entry from the National Vulnerability Database
type NVDCVE struct {
	ID               string     `json:"id" db:"id"`
	Description      string     `json:"description" db:"description"`
	CVSSV2Score      *float64   `json:"cvss_v2_score" db:"cvss_v2_score"`
	CVSSV2Vector     *string    `json:"cvss_v2_vector" db:"cvss_v2_vector"`
	CVSSV3Score      *float64   `json:"cvss_v3_score" db:"cvss_v3_score"`
	CVSSV3Vector     *string    `json:"cvss_v3_vector" db:"cvss_v3_vector"`
	Severity         *string    `json:"severity" db:"severity"`
	PublishedDate    *time.Time `json:"published_date" db:"published_date"`
	LastModifiedDate *time.Time `json:"last_modified_date" db:"last_modified_date"`
	CWEIDs           []string   `json:"cwe_ids" db:"cwe_ids"`
	AffectedProducts []string   `json:"affected_products" db:"affected_products"`
	References       []string   `json:"references" db:"refs"`
	RawData          string     `json:"raw_data" db:"raw_data"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// ATTACKTechnique represents a MITRE ATT&CK technique
type ATTACKTechnique struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	Tactic          string    `json:"tactic" db:"tactic"`
	Tactics         []string  `json:"tactics" db:"tactics"`
	Platforms       []string  `json:"platforms" db:"platforms"`
	KillChainPhases []string  `json:"kill_chain_phases" db:"kill_chain_phases"`
	DataSources     []string  `json:"data_sources" db:"data_sources"`
	Detection       *string   `json:"detection" db:"detection"`
	Mitigation      *string   `json:"mitigation" db:"mitigation"`
	References      []string  `json:"references" db:"refs"`
	SubTechniques   []string  `json:"sub_techniques" db:"sub_techniques"`
	ParentTechnique *string   `json:"parent_technique" db:"parent_technique"`
	RawData         string    `json:"raw_data" db:"raw_data"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ATTACKTactic represents a MITRE ATT&CK tactic
type ATTACKTactic struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	ExternalID      *string   `json:"external_id" db:"external_id"`
	KillChainPhases []string  `json:"kill_chain_phases" db:"kill_chain_phases"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// OWASPProcedure represents an OWASP testing procedure
type OWASPProcedure struct {
	ID          string    `json:"id" db:"id"`
	Category    string    `json:"category" db:"category"`
	Subcategory *string   `json:"subcategory" db:"subcategory"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Objective   *string   `json:"objective" db:"objective"`
	HowToTest   *string   `json:"how_to_test" db:"how_to_test"`
	Tools       []string  `json:"tools" db:"tools"`
	References  []string  `json:"references" db:"refs"`
	Severity    *string   `json:"severity" db:"severity"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SecurityDataUpdate represents the status of security data updates
type SecurityDataUpdate struct {
	ID           string     `json:"id" db:"id"`
	DataSource   string     `json:"data_source" db:"data_source"`
	LastUpdate   *time.Time `json:"last_update" db:"last_update"`
	TotalRecords *int       `json:"total_records" db:"total_records"`
	UpdateStatus string     `json:"update_status" db:"update_status"`
	ErrorMessage *string    `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// SecurityQueryRequest represents a request to query security data
type SecurityQueryRequest struct {
	Query      string                 `json:"query"`
	DataSource string                 `json:"data_source"` // "nvd", "attack", "owasp"
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
	SortBy     string                 `json:"sort_by,omitempty"`
	SortOrder  string                 `json:"sort_order,omitempty"` // "asc", "desc"
}

// SecurityQueryResponse represents the response from a security data query
type SecurityQueryResponse struct {
	Results    []interface{}          `json:"results"`
	TotalCount int                    `json:"total_count"`
	DataSource string                 `json:"data_source"`
	Query      string                 `json:"query"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Limit      int                    `json:"limit"`
	Offset     int                    `json:"offset"`
	HasMore    bool                   `json:"has_more"`
}

// NVDSearchRequest represents a request to search NVD data
type NVDSearchRequest struct {
	CWEID           *string    `json:"cwe_id,omitempty"`
	Component       *string    `json:"component,omitempty"`
	Severity        *string    `json:"severity,omitempty"`
	MinCVSS         *float64   `json:"min_cvss,omitempty"`
	MaxCVSS         *float64   `json:"max_cvss,omitempty"`
	PublishedAfter  *time.Time `json:"published_after,omitempty"`
	PublishedBefore *time.Time `json:"published_before,omitempty"`
	Limit           int        `json:"limit,omitempty"`
	Offset          int        `json:"offset,omitempty"`
}

// ATTACKSearchRequest represents a request to search ATT&CK data
type ATTACKSearchRequest struct {
	TechniqueID *string `json:"technique_id,omitempty"`
	Tactic      *string `json:"tactic,omitempty"`
	Platform    *string `json:"platform,omitempty"`
	Query       *string `json:"query,omitempty"`
	Limit       int     `json:"limit,omitempty"`
	Offset      int     `json:"offset,omitempty"`
}

// OWASPSearchRequest represents a request to search OWASP data
type OWASPSearchRequest struct {
	Category          *string `json:"category,omitempty"`
	VulnerabilityType *string `json:"vulnerability_type,omitempty"`
	TestingPhase      *string `json:"testing_phase,omitempty"`
	Severity          *string `json:"severity,omitempty"`
	Query             *string `json:"query,omitempty"`
	Limit             int     `json:"limit,omitempty"`
	Offset            int     `json:"offset,omitempty"`
}

// SecurityDataSummary represents a summary of security data for context
type SecurityDataSummary struct {
	DataSource    string         `json:"data_source"`
	TotalRecords  int            `json:"total_records"`
	LastUpdate    *time.Time     `json:"last_update"`
	TopCategories map[string]int `json:"top_categories,omitempty"`
	RecentEntries []interface{}  `json:"recent_entries,omitempty"`
	Summary       string         `json:"summary"`
}

// Custom JSON marshaling for slices to handle database storage
func (n *NVDCVE) MarshalJSON() ([]byte, error) {
	type Alias NVDCVE
	return json.Marshal(&struct {
		*Alias
		CWEIDs           []string `json:"cwe_ids"`
		AffectedProducts []string `json:"affected_products"`
		References       []string `json:"references"`
	}{
		Alias:            (*Alias)(n),
		CWEIDs:           n.CWEIDs,
		AffectedProducts: n.AffectedProducts,
		References:       n.References,
	})
}

func (n *NVDCVE) UnmarshalJSON(data []byte) error {
	type Alias NVDCVE
	aux := &struct {
		*Alias
		CWEIDs           []string `json:"cwe_ids"`
		AffectedProducts []string `json:"affected_products"`
		References       []string `json:"references"`
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n.CWEIDs = aux.CWEIDs
	n.AffectedProducts = aux.AffectedProducts
	n.References = aux.References
	return nil
}

func (a *ATTACKTechnique) MarshalJSON() ([]byte, error) {
	type Alias ATTACKTechnique
	return json.Marshal(&struct {
		*Alias
		Tactics         []string `json:"tactics"`
		Platforms       []string `json:"platforms"`
		KillChainPhases []string `json:"kill_chain_phases"`
		DataSources     []string `json:"data_sources"`
		References      []string `json:"references"`
		SubTechniques   []string `json:"sub_techniques"`
	}{
		Alias:           (*Alias)(a),
		Tactics:         a.Tactics,
		Platforms:       a.Platforms,
		KillChainPhases: a.KillChainPhases,
		DataSources:     a.DataSources,
		References:      a.References,
		SubTechniques:   a.SubTechniques,
	})
}

func (a *ATTACKTechnique) UnmarshalJSON(data []byte) error {
	type Alias ATTACKTechnique
	aux := &struct {
		*Alias
		Tactics         []string `json:"tactics"`
		Platforms       []string `json:"platforms"`
		KillChainPhases []string `json:"kill_chain_phases"`
		DataSources     []string `json:"data_sources"`
		References      []string `json:"references"`
		SubTechniques   []string `json:"sub_techniques"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.Tactics = aux.Tactics
	a.Platforms = aux.Platforms
	a.KillChainPhases = aux.KillChainPhases
	a.DataSources = aux.DataSources
	a.References = aux.References
	a.SubTechniques = aux.SubTechniques
	return nil
}

func (o *OWASPProcedure) MarshalJSON() ([]byte, error) {
	type Alias OWASPProcedure
	return json.Marshal(&struct {
		*Alias
		Tools      []string `json:"tools"`
		References []string `json:"references"`
	}{
		Alias:      (*Alias)(o),
		Tools:      o.Tools,
		References: o.References,
	})
}

func (o *OWASPProcedure) UnmarshalJSON(data []byte) error {
	type Alias OWASPProcedure
	aux := &struct {
		*Alias
		Tools      []string `json:"tools"`
		References []string `json:"references"`
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	o.Tools = aux.Tools
	o.References = aux.References
	return nil
}
