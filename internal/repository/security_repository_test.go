package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupSecurityTestDB(t *testing.T) (*database.Database, *SecurityRepository) {
	// Create a temporary database for testing
	db, err := database.NewDatabase(":memory:", log.New(os.Stderr))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	repo := NewSecurityRepository(db, log.New(os.Stderr))
	return db, repo
}

func TestSecurityRepository_StoreNVDDataset(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Create test CVE data
	cves := []models.NVDCVE{
		{
			ID:               "CVE-2024-1234",
			Description:      "Test vulnerability 1",
			CVSSV3Score:      floatPtr(9.8),
			Severity:         stringPtr("CRITICAL"),
			CWEIDs:           []string{"CWE-89"},
			AffectedProducts: []string{"apache:http_server"},
			References:       []string{"https://example.com/ref1"},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               "CVE-2024-1235",
			Description:      "Test vulnerability 2",
			CVSSV3Score:      floatPtr(7.5),
			Severity:         stringPtr("HIGH"),
			CWEIDs:           []string{"CWE-79"},
			AffectedProducts: []string{"nginx"},
			References:       []string{"https://example.com/ref2"},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	// Store the dataset
	err := repo.StoreNVDDataset(ctx, cves)
	assert.NoError(t, err)

	// Verify data was stored
	// We can't easily query without implementing the query method, but we can check for errors
	assert.NoError(t, err)
}

func TestSecurityRepository_StoreATTACKDataset(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Create test ATT&CK data
	techniques := []models.ATTACKTechnique{
		{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Test technique",
			Tactic:      "defense-evasion",
			Tactics:     []string{"defense-evasion"},
			Platforms:   []string{"windows"},
			DataSources: []string{"process-monitoring"},
			References:  []string{"https://attack.mitre.org/techniques/T1055/001/"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tactics := []models.ATTACKTactic{
		{
			ID:              "TA0005",
			Name:            "Defense Evasion",
			Description:     "Test tactic",
			ExternalID:      stringPtr("defense-evasion"),
			KillChainPhases: []string{"defense-evasion"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	// Store the dataset
	err := repo.StoreATTACKDataset(ctx, techniques, tactics)
	assert.NoError(t, err)
}

func TestSecurityRepository_QueryNVD(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// First store some test data
	cves := []models.NVDCVE{
		{
			ID:               "CVE-2024-1234",
			Description:      "SQL injection vulnerability",
			CVSSV3Score:      floatPtr(9.8),
			Severity:         stringPtr("CRITICAL"),
			CWEIDs:           []string{"CWE-89"},
			AffectedProducts: []string{"apache:http_server"},
			References:       []string{"https://example.com/ref1"},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	err := repo.StoreNVDDataset(ctx, cves)
	assert.NoError(t, err)

	// Test querying by CWE ID
	cweID := "CWE-89"
	req := models.NVDSearchRequest{
		CWEID: &cweID,
		Limit: 10,
	}

	results, totalCount, err := repo.QueryNVD(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, "CVE-2024-1234", results[0].ID)
}

func TestSecurityRepository_QueryATTACK(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// First store some test data
	techniques := []models.ATTACKTechnique{
		{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Test technique for process injection",
			Tactic:      "defense-evasion",
			Tactics:     []string{"defense-evasion"},
			Platforms:   []string{"windows"},
			DataSources: []string{"process-monitoring"},
			References:  []string{"https://attack.mitre.org/techniques/T1055/001/"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tactics := []models.ATTACKTactic{
		{
			ID:              "TA0005",
			Name:            "Defense Evasion",
			Description:     "Test tactic",
			ExternalID:      stringPtr("defense-evasion"),
			KillChainPhases: []string{"defense-evasion"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	err := repo.StoreATTACKDataset(ctx, techniques, tactics)
	assert.NoError(t, err)

	// Test querying by tactic
	tactic := "defense-evasion"
	req := models.ATTACKSearchRequest{
		Tactic: &tactic,
		Limit:  10,
	}

	results, totalCount, err := repo.QueryATTACK(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, "T1055.001", results[0].ID)
}

func TestSecurityRepository_GetSecurityDataSummary(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Store some test data first
	cves := []models.NVDCVE{
		{
			ID:          "CVE-2024-1234",
			Description: "Test vulnerability",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	techniques := []models.ATTACKTechnique{
		{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Test technique",
			Tactic:      "defense-evasion",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	err := repo.StoreNVDDataset(ctx, cves)
	assert.NoError(t, err)

	err = repo.StoreATTACKDataset(ctx, techniques, []models.ATTACKTactic{})
	assert.NoError(t, err)

	// Get summary
	summaries, err := repo.GetSecurityDataSummary(ctx)
	assert.NoError(t, err)
	assert.Contains(t, summaries, "nvd")
	assert.Contains(t, summaries, "attack")

	nvdSummary := summaries["nvd"]
	assert.Equal(t, "nvd", nvdSummary.DataSource)
	assert.Equal(t, 1, nvdSummary.TotalRecords)

	attackSummary := summaries["attack"]
	assert.Equal(t, "attack", attackSummary.DataSource)
	assert.Equal(t, 1, attackSummary.TotalRecords)
}

func TestSecurityRepository_UpdateSecurityDataStatus(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Test updating status
	totalRecords := 100
	errorMessage := "Test error"

	err := repo.UpdateSecurityDataStatus(ctx, "nvd", "completed", &totalRecords, &errorMessage)
	assert.NoError(t, err)

	// Test updating with nil values
	err = repo.UpdateSecurityDataStatus(ctx, "attack", "failed", nil, nil)
	assert.NoError(t, err)
}

func TestSecurityRepository_QueryNVD_WithFilters(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Store test data with different severities
	cves := []models.NVDCVE{
		{
			ID:          "CVE-2024-1234",
			Description: "Critical vulnerability",
			CVSSV3Score: floatPtr(9.8),
			Severity:    stringPtr("CRITICAL"),
			CWEIDs:      []string{"CWE-89"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "CVE-2024-1235",
			Description: "High vulnerability",
			CVSSV3Score: floatPtr(7.5),
			Severity:    stringPtr("HIGH"),
			CWEIDs:      []string{"CWE-79"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	err := repo.StoreNVDDataset(ctx, cves)
	assert.NoError(t, err)

	// Test filtering by severity
	severity := "CRITICAL"
	req := models.NVDSearchRequest{
		Severity: &severity,
		Limit:    10,
	}

	results, totalCount, err := repo.QueryNVD(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, "CRITICAL", *results[0].Severity)

	// Test filtering by minimum CVSS score
	minCVSS := 8.0
	req = models.NVDSearchRequest{
		MinCVSS: &minCVSS,
		Limit:   10,
	}

	results, totalCount, err = repo.QueryNVD(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, 9.8, *results[0].CVSSV3Score)
}

func TestSecurityRepository_QueryATTACK_WithFilters(t *testing.T) {
	_, repo := setupSecurityTestDB(t)
	ctx := context.Background()

	// Store test data with different platforms
	techniques := []models.ATTACKTechnique{
		{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Windows process injection",
			Tactic:      "defense-evasion",
			Platforms:   []string{"windows"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "T1055.002",
			Name:        "Thread Execution Hijacking",
			Description: "Linux thread hijacking",
			Tactic:      "defense-evasion",
			Platforms:   []string{"linux"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	err := repo.StoreATTACKDataset(ctx, techniques, []models.ATTACKTactic{})
	assert.NoError(t, err)

	// Test filtering by platform
	platform := "windows"
	req := models.ATTACKSearchRequest{
		Platform: &platform,
		Limit:    10,
	}

	results, totalCount, err := repo.QueryATTACK(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, "T1055.001", results[0].ID)
	assert.Contains(t, results[0].Platforms, "windows")

	// Test filtering by technique ID
	techniqueID := "T1055.002"
	req = models.ATTACKSearchRequest{
		TechniqueID: &techniqueID,
		Limit:       10,
	}

	results, totalCount, err = repo.QueryATTACK(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, results, 1)
	assert.Equal(t, "T1055.002", results[0].ID)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
