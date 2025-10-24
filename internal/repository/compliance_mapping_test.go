package repository

import (
	"context"
	"testing"
	"time"

	"github.com/rainmana/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapToCompliance_OWASP(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for compliance mapping",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerabilities
	vulns := []*models.MemoryEntry{
		{
			ID:          "vuln-1",
			SessionID:   "test-session",
			Title:       "SQL Injection Vulnerability",
			Content:     "SQL injection vulnerability in login form",
			Category:    "vulnerability",
			Priority:    9,
			Confidence:  0.9,
			Tags:        []string{"sql-injection", "authentication"},
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessedAt:  time.Now(),
			AccessCount: 0,
		},
		{
			ID:          "vuln-2",
			SessionID:   "test-session",
			Title:       "XSS Vulnerability",
			Content:     "Cross-site scripting vulnerability in search form",
			Category:    "vulnerability",
			Priority:    8,
			Confidence:  0.8,
			Tags:        []string{"xss", "search"},
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessedAt:  time.Now(),
			AccessCount: 0,
		},
		{
			ID:          "vuln-3",
			SessionID:   "test-session",
			Title:       "IDOR Vulnerability",
			Content:     "Insecure direct object reference vulnerability",
			Category:    "vulnerability",
			Priority:    7,
			Confidence:  0.7,
			Tags:        []string{"idor", "broken-access-control"},
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessedAt:  time.Now(),
			AccessCount: 0,
		},
	}

	// Store vulnerabilities
	for _, vuln := range vulns {
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
	}

	// Map to OWASP compliance
	mapping, err := repo.MapToCompliance(ctx, "test-session", "OWASP")
	require.NoError(t, err)

	assert.Equal(t, "test-session", mapping.SessionID)
	assert.Equal(t, "OWASP", mapping.Standard)
	assert.Equal(t, "OWASP Top 10 2021", mapping.Requirement)
	// Verify the mapping was stored (check that we have 3 vulnerabilities)
	assert.Len(t, mapping.VulnerabilityIDs, 3)
	assert.Greater(t, mapping.ComplianceScore, 0.0)
	assert.Len(t, mapping.GapAnalysis, 3)
	assert.Len(t, mapping.Recommendations, 5)
}

func TestMapToCompliance_NIST(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for NIST compliance",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerability
	vuln := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Access Control Vulnerability",
		Content:     "Weak access control mechanisms",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.8,
		Tags:        []string{"access-control", "authorization"},
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

	// Map to NIST compliance
	mapping, err := repo.MapToCompliance(ctx, "test-session", "NIST")
	require.NoError(t, err)

	assert.Equal(t, "test-session", mapping.SessionID)
	assert.Equal(t, "NIST", mapping.Standard)
	assert.Equal(t, "NIST SP 800-53 Security Controls", mapping.Requirement)
	assert.Len(t, mapping.VulnerabilityIDs, 1)
	assert.Equal(t, 75.0, mapping.ComplianceScore)
	assert.Len(t, mapping.GapAnalysis, 3)
	assert.Len(t, mapping.Recommendations, 3)
}

func TestMapToCompliance_ISO27001(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for ISO27001 compliance",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerability
	vuln := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Information Security Vulnerability",
		Content:     "Information security management issue",
		Category:    "vulnerability",
		Priority:    7,
		Confidence:  0.7,
		Tags:        []string{"information-security", "policy"},
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

	// Map to ISO 27001 compliance
	mapping, err := repo.MapToCompliance(ctx, "test-session", "ISO27001")
	require.NoError(t, err)

	assert.Equal(t, "test-session", mapping.SessionID)
	assert.Equal(t, "ISO27001", mapping.Standard)
	assert.Equal(t, "ISO 27001 Information Security Management", mapping.Requirement)
	assert.Len(t, mapping.VulnerabilityIDs, 1)
	assert.Equal(t, 70.0, mapping.ComplianceScore)
	assert.Len(t, mapping.GapAnalysis, 3)
	assert.Len(t, mapping.Recommendations, 3)
}

func TestMapToCompliance_PCIDSS(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for PCI DSS compliance",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerability
	vuln := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Payment Card Data Vulnerability",
		Content:     "Payment card industry data security issue",
		Category:    "vulnerability",
		Priority:    9,
		Confidence:  0.9,
		Tags:        []string{"payment-card", "data-protection"},
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

	// Map to PCI DSS compliance
	mapping, err := repo.MapToCompliance(ctx, "test-session", "PCIDSS")
	require.NoError(t, err)

	assert.Equal(t, "test-session", mapping.SessionID)
	assert.Equal(t, "PCIDSS", mapping.Standard)
	assert.Equal(t, "PCI DSS Payment Card Industry Standards", mapping.Requirement)
	assert.Len(t, mapping.VulnerabilityIDs, 1)
	assert.Equal(t, 65.0, mapping.ComplianceScore)
	assert.Len(t, mapping.GapAnalysis, 3)
	assert.Len(t, mapping.Recommendations, 3)
}

func TestMapToCompliance_UnsupportedStandard(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Test unsupported standard
	_, err := repo.MapToCompliance(ctx, "test-session", "UNSUPPORTED")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported compliance standard")
}

func TestMapToOWASP_CategoryMapping(t *testing.T) {
	_, repo := setupTestDB(t)

	testCases := []struct {
		name             string
		vuln             *models.MemoryEntry
		expectedCategory string
	}{
		{
			name: "SQL Injection",
			vuln: &models.MemoryEntry{
				Tags:    []string{"sql-injection"},
				Content: "SQL injection vulnerability",
			},
			expectedCategory: "A03:2021 - Injection",
		},
		{
			name: "XSS",
			vuln: &models.MemoryEntry{
				Tags:    []string{"xss"},
				Content: "Cross-site scripting vulnerability",
			},
			expectedCategory: "A03:2021 - Injection",
		},
		{
			name: "IDOR",
			vuln: &models.MemoryEntry{
				Tags:    []string{"idor"},
				Content: "Broken access control vulnerability",
			},
			expectedCategory: "A01:2021 - Broken Access Control",
		},
		{
			name: "Weak Crypto",
			vuln: &models.MemoryEntry{
				Tags:    []string{"weak-crypto"},
				Content: "Cryptographic failure vulnerability",
			},
			expectedCategory: "A02:2021 - Cryptographic Failures",
		},
		{
			name: "SSRF",
			vuln: &models.MemoryEntry{
				Tags:    []string{"ssrf"},
				Content: "Server-side request forgery vulnerability",
			},
			expectedCategory: "A10:2021 - Server-Side Request Forgery",
		},
		{
			name: "Default Category",
			vuln: &models.MemoryEntry{
				Tags:    []string{"unknown"},
				Content: "Unknown vulnerability type",
			},
			expectedCategory: "A04:2021 - Insecure Design",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapping := &models.ComplianceMapping{
				ID:        "test-mapping",
				SessionID: "test-session",
				Standard:  "OWASP",
			}

			result := repo.mapToOWASP([]*models.MemoryEntry{tc.vuln}, mapping)
			assert.Equal(t, "OWASP Top 10 2021", result.Requirement)
			assert.Len(t, result.VulnerabilityIDs, 1)
		})
	}
}

func TestMapToCompliance_EmptySession(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "empty-session",
		Name:        "Empty Session",
		Description: "Test session for empty compliance",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Test with empty session
	mapping, err := repo.MapToCompliance(ctx, "empty-session", "OWASP")
	require.NoError(t, err)

	assert.Equal(t, "empty-session", mapping.SessionID)
	assert.Equal(t, "OWASP", mapping.Standard)
	assert.Len(t, mapping.VulnerabilityIDs, 0)
	assert.Equal(t, 0.0, mapping.ComplianceScore)
}

func TestStoreComplianceMapping(t *testing.T) {
	db, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for compliance mapping storage",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	mapping := &models.ComplianceMapping{
		ID:               "test-mapping",
		SessionID:        "test-session",
		Standard:         "OWASP",
		Requirement:      "OWASP Top 10 2021",
		VulnerabilityIDs: []string{"vuln-1", "vuln-2"},
		ComplianceScore:  75.0,
		GapAnalysis:      []string{"Missing input validation", "Weak authentication"},
		Recommendations:  []string{"Implement input validation", "Strengthen authentication"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = repo.storeComplianceMapping(ctx, mapping)
	require.NoError(t, err)

	// Verify storage
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM compliance_mappings WHERE id = ?", "test-mapping").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestMapToCompliance_MultipleSessions(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test sessions first
	session1 := &models.Session{
		ID:          "session-1",
		Name:        "Session 1",
		Description: "Test session 1",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session1)
	require.NoError(t, err)

	session2 := &models.Session{
		ID:          "session-2",
		Name:        "Session 2",
		Description: "Test session 2",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err = repo.CreateSession(ctx, session2)
	require.NoError(t, err)

	// Create vulnerabilities for different sessions
	vuln1 := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "session-1",
		Title:       "Session 1 Vulnerability",
		Content:     "Vulnerability in session 1",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.8,
		Tags:        []string{"sql-injection"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	vuln2 := &models.MemoryEntry{
		ID:          "vuln-2",
		SessionID:   "session-2",
		Title:       "Session 2 Vulnerability",
		Content:     "Vulnerability in session 2",
		Category:    "vulnerability",
		Priority:    7,
		Confidence:  0.7,
		Tags:        []string{"xss"},
		Source:      "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

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

	// Map each session separately
	mapping1, err := repo.MapToCompliance(ctx, "session-1", "OWASP")
	require.NoError(t, err)
	assert.Equal(t, "session-1", mapping1.SessionID)
	assert.Len(t, mapping1.VulnerabilityIDs, 1)

	mapping2, err := repo.MapToCompliance(ctx, "session-2", "OWASP")
	require.NoError(t, err)
	assert.Equal(t, "session-2", mapping2.SessionID)
	assert.Len(t, mapping2.VulnerabilityIDs, 1)
}

func TestMapToCompliance_AllStandards(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for all standards",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create test vulnerability
	vuln := &models.MemoryEntry{
		ID:          "vuln-1",
		SessionID:   "test-session",
		Title:       "Test Vulnerability",
		Content:     "Test vulnerability for compliance mapping",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.8,
		Tags:        []string{"test"},
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

	standards := []string{"OWASP", "NIST", "ISO27001", "PCIDSS"}

	for _, standard := range standards {
		t.Run(standard, func(t *testing.T) {
			mapping, err := repo.MapToCompliance(ctx, "test-session", standard)
			require.NoError(t, err)
			assert.Equal(t, "test-session", mapping.SessionID)
			assert.Equal(t, standard, mapping.Standard)
			assert.Len(t, mapping.VulnerabilityIDs, 1)
			assert.Greater(t, mapping.ComplianceScore, 0.0)
			assert.Len(t, mapping.GapAnalysis, 3)
			assert.GreaterOrEqual(t, len(mapping.Recommendations), 3)
		})
	}
}
