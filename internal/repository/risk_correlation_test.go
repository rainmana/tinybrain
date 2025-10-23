package repository

import (
	"context"
	"testing"
	"time"

	"github.com/rainmana/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeRiskCorrelation(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for risk correlation",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerabilities
	vuln1 := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "SQL Injection Vulnerability",
		Content:     "SQL injection vulnerability in login form",
		Category:    "vulnerability",
		Priority:    9,
		Confidence:  0.9,
		Tags:        []string{"sql-injection", "authentication", "high"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	vuln2 := &models.MemoryEntry{
		ID:          "vuln-2",
		SessionID:   "test-session",
		Title:       "XSS Vulnerability",
		Content:     "Cross-site scripting vulnerability in search form",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.8,
		Tags:        []string{"xss", "search", "medium"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	// Store vulnerabilities
	req1 := &models.CreateMemoryEntryRequest{
		SessionID:   vuln1.SessionID,
		Title:       vuln1.Title,
		Content:     vuln1.Content,
		ContentType: vuln1.ContentType,
		Category:    vuln1.Category,
		Priority:    vuln1.Priority,
		Confidence:  vuln1.Confidence,
		Tags:        vuln1.Tags,
		Source:      vuln1.Source,
	}
	_, err = repo.CreateMemoryEntry(ctx, req1)
	require.NoError(t, err)

	req2 := &models.CreateMemoryEntryRequest{
		SessionID:   vuln2.SessionID,
		Title:       vuln2.Title,
		Content:     vuln2.Content,
		ContentType: vuln2.ContentType,
		Category:    vuln2.Category,
		Priority:    vuln2.Priority,
		Confidence:  vuln2.Confidence,
		Tags:        vuln2.Tags,
		Source:      vuln2.Source,
	}
	_, err = repo.CreateMemoryEntry(ctx, req2)
	require.NoError(t, err)

	// Analyze risk correlation
	correlations, err := repo.AnalyzeRiskCorrelation(ctx, "test-session")
	require.NoError(t, err)
	assert.Len(t, correlations, 1)

	correlation := correlations[0]
	assert.Equal(t, "test-session", correlation.SessionID)
	// Check that we have the expected number of vulnerabilities
	assert.Len(t, correlation.SecondaryVulnIDs, 1)
	assert.Greater(t, correlation.RiskMultiplier, 1.0)
	assert.Len(t, correlation.AttackChain, 1)
	assert.Contains(t, correlation.AttackChain[0], "SQL Injection")
	assert.Contains(t, correlation.AttackChain[0], "XSS")
	assert.Greater(t, correlation.Confidence, 0.0)
}

func TestCalculateRiskCorrelation_SQLInjectionXSS(t *testing.T) {
	_, repo := setupTestDB(t)

	primary := &models.MemoryEntry{
		ID:       "vuln-1",
		Tags:     []string{"sql-injection", "authentication"},
		Priority: 9,
	}

	secondary := &models.MemoryEntry{
		ID:       "vuln-2",
		Tags:     []string{"xss", "search"},
		Priority: 8,
	}

	correlation := repo.calculateRiskCorrelation(primary, secondary)
	require.NotNil(t, correlation)

	assert.Equal(t, "vuln-1", correlation.PrimaryVulnID)
	assert.Equal(t, []string{"vuln-2"}, correlation.SecondaryVulnIDs)
	assert.Greater(t, correlation.RiskMultiplier, 1.0)
	assert.Len(t, correlation.AttackChain, 1)
	assert.Contains(t, correlation.AttackChain[0], "SQL Injection → Data Extraction → XSS → Session Hijacking")
}

func TestCalculateRiskCorrelation_SSRFFileUpload(t *testing.T) {
	_, repo := setupTestDB(t)

	primary := &models.MemoryEntry{
		ID:       "vuln-1",
		Tags:     []string{"ssrf", "internal"},
		Priority: 9,
	}

	secondary := &models.MemoryEntry{
		ID:       "vuln-2",
		Tags:     []string{"file-upload", "malware"},
		Priority: 8,
	}

	correlation := repo.calculateRiskCorrelation(primary, secondary)
	require.NotNil(t, correlation)

	assert.Len(t, correlation.AttackChain, 1)
	assert.Contains(t, correlation.AttackChain[0], "SSRF → Internal Network Access → File Upload → Code Execution")
}

func TestCalculateRiskCorrelation_NoCorrelation(t *testing.T) {
	_, repo := setupTestDB(t)

	primary := &models.MemoryEntry{
		ID:       "vuln-1",
		Tags:     []string{"sql-injection"},
		Priority: 9,
	}

	secondary := &models.MemoryEntry{
		ID:       "vuln-2",
		Tags:     []string{"weak-crypto"},
		Priority: 6,
	}

	correlation := repo.calculateRiskCorrelation(primary, secondary)
	assert.Nil(t, correlation) // No correlation should be found
}

func TestCalculateRiskMultiplier(t *testing.T) {
	_, repo := setupTestDB(t)

	testCases := []struct {
		name        string
		primary     *models.MemoryEntry
		secondary   *models.MemoryEntry
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "Critical + Critical",
			primary:     &models.MemoryEntry{Priority: 10, Tags: []string{"rce"}},
			secondary:   &models.MemoryEntry{Priority: 10, Tags: []string{"rce"}},
			expectedMin: 2.0,
			expectedMax: 5.0,
		},
		{
			name:        "High + High",
			primary:     &models.MemoryEntry{Priority: 8, Tags: []string{"sql-injection"}},
			secondary:   &models.MemoryEntry{Priority: 8, Tags: []string{"xss"}},
			expectedMin: 1.5,
			expectedMax: 3.0,
		},
		{
			name:        "Medium + Medium",
			primary:     &models.MemoryEntry{Priority: 7, Tags: []string{"weak-crypto"}},
			secondary:   &models.MemoryEntry{Priority: 7, Tags: []string{"weak-crypto"}},
			expectedMin: 1.0,
			expectedMax: 2.0,
		},
		{
			name:        "SQL Injection + XSS",
			primary:     &models.MemoryEntry{Priority: 9, Tags: []string{"sql-injection"}},
			secondary:   &models.MemoryEntry{Priority: 8, Tags: []string{"xss"}},
			expectedMin: 2.0,
			expectedMax: 4.0,
		},
		{
			name:        "SSRF + File Upload",
			primary:     &models.MemoryEntry{Priority: 9, Tags: []string{"ssrf"}},
			secondary:   &models.MemoryEntry{Priority: 8, Tags: []string{"file-upload"}},
			expectedMin: 2.0,
			expectedMax: 4.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			multiplier := repo.calculateRiskMultiplier(tc.primary, tc.secondary)
			assert.GreaterOrEqual(t, multiplier, tc.expectedMin)
			assert.LessOrEqual(t, multiplier, tc.expectedMax)
		})
	}
}

func TestDetermineBusinessImpact(t *testing.T) {
	_, repo := setupTestDB(t)

	primary := &models.MemoryEntry{ID: "vuln-1"}
	secondary := &models.MemoryEntry{ID: "vuln-2"}

	testCases := []struct {
		riskMultiplier float64
		expectedLevel  string
	}{
		{3.5, "CRITICAL"},
		{2.5, "HIGH"},
		{1.8, "MEDIUM"},
		{1.2, "LOW"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedLevel, func(t *testing.T) {
			impact := repo.determineBusinessImpact(primary, secondary, tc.riskMultiplier)
			assert.Contains(t, impact, tc.expectedLevel)
		})
	}
}

func TestCalculateCorrelationConfidence(t *testing.T) {
	_, repo := setupTestDB(t)

	testCases := []struct {
		name        string
		primary     *models.MemoryEntry
		secondary   *models.MemoryEntry
		expectedMin float64
		expectedMax float64
	}{
		{
			name: "High priority with tag overlap",
			primary: &models.MemoryEntry{
				Priority: 10,
				Tags:     []string{"sql-injection", "authentication", "critical"},
			},
			secondary: &models.MemoryEntry{
				Priority: 9,
				Tags:     []string{"sql-injection", "authentication", "high"},
			},
			expectedMin: 0.7,
			expectedMax: 0.95,
		},
		{
			name: "Medium priority with some overlap",
			primary: &models.MemoryEntry{
				Priority: 7,
				Tags:     []string{"weak-crypto", "authentication"},
			},
			secondary: &models.MemoryEntry{
				Priority: 6,
				Tags:     []string{"weak-crypto", "password"},
			},
			expectedMin: 0.5,
			expectedMax: 0.8,
		},
		{
			name: "Low priority with no overlap",
			primary: &models.MemoryEntry{
				Priority: 5,
				Tags:     []string{"information-disclosure"},
			},
			secondary: &models.MemoryEntry{
				Priority: 4,
				Tags:     []string{"weak-crypto"},
			},
			expectedMin: 0.5,
			expectedMax: 0.6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			confidence := repo.calculateCorrelationConfidence(tc.primary, tc.secondary)
			assert.GreaterOrEqual(t, confidence, tc.expectedMin)
			assert.LessOrEqual(t, confidence, tc.expectedMax)
		})
	}
}

func TestStoreRiskCorrelation(t *testing.T) {
	db, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for risk correlation storage",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test memory entries first
	vuln1 := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Test Vulnerability 1",
		Content:     "Test vulnerability 1 for risk correlation",
		Category:    "vulnerability",
		Priority:    9,
		Confidence:  0.9,
		Tags:        []string{"sql-injection"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	vuln2 := &models.MemoryEntry{
		ID:          "vuln-2",
		SessionID:   "test-session",
		Title:       "Test Vulnerability 2",
		Content:     "Test vulnerability 2 for risk correlation",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.8,
		Tags:        []string{"xss"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	req1 := &models.CreateMemoryEntryRequest{
		SessionID:  vuln1.SessionID,
		Title:      vuln1.Title,
		Content:    vuln1.Content,
		Category:   vuln1.Category,
		Priority:   vuln1.Priority,
		Confidence: vuln1.Confidence,
		Tags:       vuln1.Tags,
		Source:     vuln1.Source,
	}
	vuln1ID, err := repo.CreateMemoryEntry(ctx, req1)
	require.NoError(t, err)

	req2 := &models.CreateMemoryEntryRequest{
		SessionID:  vuln2.SessionID,
		Title:      vuln2.Title,
		Content:    vuln2.Content,
		Category:   vuln2.Category,
		Priority:   vuln2.Priority,
		Confidence: vuln2.Confidence,
		Tags:       vuln2.Tags,
		Source:     vuln2.Source,
	}
	vuln2ID, err := repo.CreateMemoryEntry(ctx, req2)
	require.NoError(t, err)

	correlation := &models.RiskCorrelation{
		ID:               "test-correlation",
		SessionID:        "test-session",
		PrimaryVulnID:    vuln1ID.ID,
		SecondaryVulnIDs: []string{vuln2ID.ID},
		RiskMultiplier:   2.5,
		AttackChain:      []string{"SQL Injection → XSS → Session Hijacking"},
		BusinessImpact:   "HIGH: Significant data exposure",
		Confidence:       0.85,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = repo.storeRiskCorrelation(ctx, correlation)
	require.NoError(t, err)

	// Verify storage
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM risk_correlations WHERE id = ?", "test-correlation").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAnalyzeRiskCorrelation_EmptySession(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	correlations, err := repo.AnalyzeRiskCorrelation(ctx, "empty-session")
	require.NoError(t, err)
	assert.Len(t, correlations, 0)
}

func TestAnalyzeRiskCorrelation_SingleVulnerability(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for single vulnerability",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create single vulnerability
	vuln := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Single Vulnerability",
		Content:     "Single vulnerability for testing",
		Category:    "vulnerability",
		Priority:    9,
		Confidence:  0.9,
		Tags:        []string{"sql-injection"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	req := &models.CreateMemoryEntryRequest{
		SessionID:   vuln.SessionID,
		Title:       vuln.Title,
		Content:     vuln.Content,
		ContentType: vuln.ContentType,
		Category:    vuln.Category,
		Priority:    vuln.Priority,
		Confidence:  vuln.Confidence,
		Tags:        vuln.Tags,
		Source:      vuln.Source,
	}
	_, err = repo.CreateMemoryEntry(ctx, req)
	require.NoError(t, err)

	// Analyze risk correlation
	correlations, err := repo.AnalyzeRiskCorrelation(ctx, "test-session")
	require.NoError(t, err)
	assert.Len(t, correlations, 0) // No correlations with single vulnerability
}
