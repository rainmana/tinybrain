package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNVDCVE_JSONMarshaling(t *testing.T) {
	cve := NVDCVE{
		ID:               "CVE-2024-1234",
		Description:      "Test vulnerability",
		CVSSV3Score:      floatPtr(9.8),
		Severity:         stringPtr("CRITICAL"),
		CWEIDs:           []string{"CWE-89", "CWE-79"},
		AffectedProducts: []string{"product1", "product2"},
		References:       []string{"https://example.com/ref1", "https://example.com/ref2"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(cve)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var unmarshaled NVDCVE
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, cve.ID, unmarshaled.ID)
	assert.Equal(t, cve.Description, unmarshaled.Description)
	assert.Equal(t, cve.CWEIDs, unmarshaled.CWEIDs)
	assert.Equal(t, cve.AffectedProducts, unmarshaled.AffectedProducts)
	assert.Equal(t, cve.References, unmarshaled.References)
}

func TestATTACKTechnique_JSONMarshaling(t *testing.T) {
	technique := ATTACKTechnique{
		ID:              "T1055.001",
		Name:            "Process Injection",
		Description:     "Test technique",
		Tactic:          "defense-evasion",
		Tactics:         []string{"defense-evasion", "persistence"},
		Platforms:       []string{"windows", "linux"},
		KillChainPhases: []string{"execution"},
		DataSources:     []string{"process-monitoring"},
		References:      []string{"https://attack.mitre.org/techniques/T1055/001/"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(technique)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var unmarshaled ATTACKTechnique
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, technique.ID, unmarshaled.ID)
	assert.Equal(t, technique.Name, unmarshaled.Name)
	assert.Equal(t, technique.Tactics, unmarshaled.Tactics)
	assert.Equal(t, technique.Platforms, unmarshaled.Platforms)
	assert.Equal(t, technique.References, unmarshaled.References)
}

func TestATTACKTactic_JSONMarshaling(t *testing.T) {
	tactic := ATTACKTactic{
		ID:              "tactic-123",
		Name:            "Defense Evasion",
		Description:     "Test tactic",
		ExternalID:      stringPtr("defense-evasion"),
		KillChainPhases: []string{"execution"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(tactic)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var unmarshaled ATTACKTactic
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, tactic.ID, unmarshaled.ID)
	assert.Equal(t, tactic.Name, unmarshaled.Name)
	assert.Equal(t, tactic.KillChainPhases, unmarshaled.KillChainPhases)
}

func TestOWASPProcedure_JSONMarshaling(t *testing.T) {
	procedure := OWASPProcedure{
		ID:          "proc-123",
		Category:    "authentication",
		Title:       "Test Procedure",
		Description: "Test description",
		Tools:       []string{"tool1", "tool2"},
		References:  []string{"https://owasp.org/ref1"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(procedure)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var unmarshaled OWASPProcedure
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, procedure.ID, unmarshaled.ID)
	assert.Equal(t, procedure.Category, unmarshaled.Category)
	assert.Equal(t, procedure.Tools, unmarshaled.Tools)
	assert.Equal(t, procedure.References, unmarshaled.References)
}

func TestSecurityQueryRequest_Validation(t *testing.T) {
	req := SecurityQueryRequest{
		Query:      "test query",
		DataSource: "nvd",
		Limit:      10,
		Offset:     0,
		SortBy:     "published_date",
		SortOrder:  "desc",
	}

	assert.Equal(t, "test query", req.Query)
	assert.Equal(t, "nvd", req.DataSource)
	assert.Equal(t, 10, req.Limit)
	assert.Equal(t, 0, req.Offset)
}

func TestNVDSearchRequest_Validation(t *testing.T) {
	cweID := "CWE-89"
	component := "apache"
	severity := "HIGH"
	minCVSS := 7.0

	req := NVDSearchRequest{
		CWEID:     &cweID,
		Component: &component,
		Severity:  &severity,
		MinCVSS:   &minCVSS,
		Limit:     5,
		Offset:    0,
	}

	assert.Equal(t, "CWE-89", *req.CWEID)
	assert.Equal(t, "apache", *req.Component)
	assert.Equal(t, "HIGH", *req.Severity)
	assert.Equal(t, 7.0, *req.MinCVSS)
}

func TestATTACKSearchRequest_Validation(t *testing.T) {
	techniqueID := "T1055.001"
	tactic := "defense-evasion"
	platform := "windows"
	query := "injection"

	req := ATTACKSearchRequest{
		TechniqueID: &techniqueID,
		Tactic:      &tactic,
		Platform:    &platform,
		Query:       &query,
		Limit:       10,
		Offset:      0,
	}

	assert.Equal(t, "T1055.001", *req.TechniqueID)
	assert.Equal(t, "defense-evasion", *req.Tactic)
	assert.Equal(t, "windows", *req.Platform)
	assert.Equal(t, "injection", *req.Query)
}

func TestOWASPSearchRequest_Validation(t *testing.T) {
	category := "authentication"
	vulnType := "sql-injection"
	testingPhase := "dynamic"
	severity := "HIGH"
	query := "login"

	req := OWASPSearchRequest{
		Category:          &category,
		VulnerabilityType: &vulnType,
		TestingPhase:      &testingPhase,
		Severity:          &severity,
		Query:             &query,
		Limit:             5,
		Offset:            0,
	}

	assert.Equal(t, "authentication", *req.Category)
	assert.Equal(t, "sql-injection", *req.VulnerabilityType)
	assert.Equal(t, "dynamic", *req.TestingPhase)
	assert.Equal(t, "HIGH", *req.Severity)
	assert.Equal(t, "login", *req.Query)
}

func TestSecurityDataSummary_Structure(t *testing.T) {
	now := time.Now()
	summary := SecurityDataSummary{
		DataSource:   "nvd",
		TotalRecords: 1000,
		LastUpdate:   &now,
		TopCategories: map[string]int{
			"CRITICAL": 100,
			"HIGH":     200,
			"MEDIUM":   500,
			"LOW":      200,
		},
		RecentEntries: []interface{}{
			map[string]string{"id": "CVE-2024-1234", "severity": "CRITICAL"},
			map[string]string{"id": "CVE-2024-1235", "severity": "HIGH"},
		},
		Summary: "NVD database contains 1000 CVE entries",
	}

	assert.Equal(t, "nvd", summary.DataSource)
	assert.Equal(t, 1000, summary.TotalRecords)
	assert.NotNil(t, summary.LastUpdate)
	assert.Equal(t, 4, len(summary.TopCategories))
	assert.Equal(t, 2, len(summary.RecentEntries))
	assert.NotEmpty(t, summary.Summary)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
