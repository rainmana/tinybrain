package repository

import (
	"context"
	"testing"
	"time"

	"github.com/rainmana/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapToCVE(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create a test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for CVE mapping",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Test mapping CWE-89 (SQL Injection)
	cveMapping, err := repo.MapToCVE(ctx, "test-session", "CWE-89")
	require.NoError(t, err)
	assert.NotNil(t, cveMapping)
	assert.Equal(t, "test-session", cveMapping.SessionID)
	assert.Equal(t, "CWE-89", cveMapping.CWEID)
	assert.Equal(t, 0.95, cveMapping.Confidence)
	assert.Equal(t, "nvd", cveMapping.Source)
	assert.Len(t, cveMapping.CVEList, 3)
	assert.Contains(t, cveMapping.CVEList, "CVE-2023-1234")
	assert.Contains(t, cveMapping.CVEList, "CVE-2023-5678")
	assert.Contains(t, cveMapping.CVEList, "CVE-2024-9012")

	// Test mapping CWE-79 (XSS)
	cveMapping, err = repo.MapToCVE(ctx, "test-session", "CWE-79")
	require.NoError(t, err)
	assert.Equal(t, "CWE-79", cveMapping.CWEID)
	assert.Equal(t, 0.90, cveMapping.Confidence)
	assert.Len(t, cveMapping.CVEList, 3)

	// Test mapping unknown CWE
	cveMapping, err = repo.MapToCVE(ctx, "test-session", "CWE-999")
	require.NoError(t, err)
	assert.Equal(t, "CWE-999", cveMapping.CWEID)
	assert.Equal(t, 0.50, cveMapping.Confidence)
	assert.Len(t, cveMapping.CVEList, 1)
	assert.Contains(t, cveMapping.CVEList, "CVE-2023-0000")
}

func TestGetCVEMapping(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create a test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for CVE mapping",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// First create a mapping
	_, err = repo.MapToCVE(ctx, "test-session", "CWE-89")
	require.NoError(t, err)

	// Then retrieve it
	mapping, err := repo.GetCVEMapping(ctx, "test-session", "CWE-89")
	require.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "test-session", mapping.SessionID)
	assert.Equal(t, "CWE-89", mapping.CWEID)
	assert.Equal(t, 0.95, mapping.Confidence)

	// Test non-existent mapping
	mapping, err = repo.GetCVEMapping(ctx, "test-session", "CWE-999")
	require.NoError(t, err)
	assert.Nil(t, mapping)
}

func TestMapToCVE_DatabaseStorage(t *testing.T) {
	db, repo := setupTestDB(t)
	ctx := context.Background()

	// Create a test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for CVE mapping",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create mapping
	cveMapping, err := repo.MapToCVE(ctx, "test-session", "CWE-78")
	require.NoError(t, err)

	// Verify it's stored in database
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM cve_mappings WHERE session_id = ? AND cwe_id = ?",
		"test-session", "CWE-78").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify data integrity
	var storedMapping models.CVEMapping
	var cveListJSON string
	err = db.GetDB().QueryRow(`
		SELECT id, session_id, cwe_id, cve_list, confidence, source, created_at, updated_at
		FROM cve_mappings WHERE session_id = ? AND cwe_id = ?
	`, "test-session", "CWE-78").Scan(
		&storedMapping.ID,
		&storedMapping.SessionID,
		&storedMapping.CWEID,
		&cveListJSON,
		&storedMapping.Confidence,
		&storedMapping.Source,
		&storedMapping.CreatedAt,
		&storedMapping.UpdatedAt,
	)
	require.NoError(t, err)
	assert.Equal(t, cveMapping.ID, storedMapping.ID)
	assert.Equal(t, "test-session", storedMapping.SessionID)
	assert.Equal(t, "CWE-78", storedMapping.CWEID)
	assert.Equal(t, 0.92, storedMapping.Confidence)
	assert.Equal(t, "nvd", storedMapping.Source)
}

func TestMapToCVE_MultipleSessions(t *testing.T) {
	db, repo := setupTestDB(t)
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

	// Create mappings for different sessions
	_, err = repo.MapToCVE(ctx, "session-1", "CWE-89")
	require.NoError(t, err)

	_, err = repo.MapToCVE(ctx, "session-2", "CWE-89")
	require.NoError(t, err)

	// Verify both are stored separately
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM cve_mappings WHERE cwe_id = ?", "CWE-89").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Verify session isolation
	mapping1, err := repo.GetCVEMapping(ctx, "session-1", "CWE-89")
	require.NoError(t, err)
	assert.NotNil(t, mapping1)
	assert.Equal(t, "session-1", mapping1.SessionID)

	mapping2, err := repo.GetCVEMapping(ctx, "session-2", "CWE-89")
	require.NoError(t, err)
	assert.NotNil(t, mapping2)
	assert.Equal(t, "session-2", mapping2.SessionID)
}

func TestMapToCVE_AllCWETypes(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for CWE types",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	testCases := []struct {
		cweID      string
		confidence float64
		cveCount   int
	}{
		{"CWE-89", 0.95, 3},  // SQL Injection
		{"CWE-79", 0.90, 3},  // XSS
		{"CWE-78", 0.92, 3},  // Command Injection
		{"CWE-502", 0.88, 3}, // Deserialization
		{"CWE-22", 0.87, 3},  // Path Traversal
		{"CWE-999", 0.50, 1}, // Unknown
	}

	for _, tc := range testCases {
		t.Run(tc.cweID, func(t *testing.T) {
			mapping, err := repo.MapToCVE(ctx, "test-session", tc.cweID)
			require.NoError(t, err)
			assert.Equal(t, tc.cweID, mapping.CWEID)
			assert.Equal(t, tc.confidence, mapping.Confidence)
			assert.Len(t, mapping.CVEList, tc.cveCount)
		})
	}
}

func TestMapToCVE_TimestampHandling(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test session first
	session := &models.Session{
		ID:          "test-session",
		Name:        "Test Session",
		Description: "Test session for timestamp handling",
		TaskType:    "security_review",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	before := time.Now()
	mapping, err := repo.MapToCVE(ctx, "test-session", "CWE-89")
	after := time.Now()

	require.NoError(t, err)
	assert.True(t, mapping.CreatedAt.After(before) || mapping.CreatedAt.Equal(before))
	assert.True(t, mapping.CreatedAt.Before(after) || mapping.CreatedAt.Equal(after))
	assert.True(t, mapping.UpdatedAt.After(before) || mapping.UpdatedAt.Equal(before))
	assert.True(t, mapping.UpdatedAt.Before(after) || mapping.UpdatedAt.Equal(after))
	assert.True(t, mapping.LastUpdated.After(before) || mapping.LastUpdated.Equal(before))
	assert.True(t, mapping.LastUpdated.Before(after) || mapping.LastUpdated.Equal(after))
}
