package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*database.Database, *MemoryRepository) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(t, err)

	repo := NewMemoryRepository(db.GetDB(), logger)
	return db, repo
}

func TestCreateSession(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	session := &models.Session{
		ID:          "test-session-1",
		Name:        "Test Security Review",
		Description: "A test session for security review",
		TaskType:    "security_review",
		Status:      "active",
		Metadata: map[string]interface{}{
			"target": "test-application",
			"scope":  "web-application",
		},
	}

	err := repo.CreateSession(ctx, session)
	assert.NoError(t, err)

	// Verify session was created
	retrieved, err := repo.GetSession(ctx, session.ID)
	assert.NoError(t, err)
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.Name, retrieved.Name)
	assert.Equal(t, session.TaskType, retrieved.TaskType)
	assert.Equal(t, session.Status, retrieved.Status)
	assert.Equal(t, session.Metadata, retrieved.Metadata)
}

func TestCreateMemoryEntry(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session first
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create memory entry
	req := &models.CreateMemoryEntryRequest{
		SessionID:   session.ID,
		Title:       "SQL Injection Vulnerability",
		Content:     "Found SQL injection in login form parameter 'username'",
		ContentType: "text",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.9,
		Tags:        []string{"sql-injection", "authentication", "critical"},
		Source:      "manual-testing",
	}

	entry, err := repo.CreateMemoryEntry(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, req.Title, entry.Title)
	assert.Equal(t, req.Content, entry.Content)
	assert.Equal(t, req.Category, entry.Category)
	assert.Equal(t, req.Priority, entry.Priority)
	assert.Equal(t, req.Confidence, entry.Confidence)
	assert.Equal(t, req.Tags, entry.Tags)
	assert.Equal(t, req.Source, entry.Source)
	assert.Equal(t, 0, entry.AccessCount) // Initial access count

	// Verify entry was created
	retrieved, err := repo.GetMemoryEntry(ctx, entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, retrieved.ID)
	assert.Equal(t, entry.Title, retrieved.Title)
	assert.Equal(t, 1, retrieved.AccessCount) // Should be incremented after retrieval
}

func TestSearchMemoryEntries(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create multiple memory entries
	entries := []*models.CreateMemoryEntryRequest{
		{
			SessionID: session.ID,
			Title:     "SQL Injection in Login",
			Content:   "Found SQL injection vulnerability in login form",
			Category:  "vulnerability",
			Priority:  8,
			Tags:      []string{"sql-injection", "authentication"},
		},
		{
			SessionID: session.ID,
			Title:     "XSS in Search",
			Content:   "Cross-site scripting vulnerability in search functionality",
			Category:  "vulnerability",
			Priority:  6,
			Tags:      []string{"xss", "search"},
		},
		{
			SessionID: session.ID,
			Title:     "Authentication Bypass",
			Content:   "Found way to bypass authentication using session manipulation",
			Category:  "exploit",
			Priority:  9,
			Tags:      []string{"authentication", "session"},
		},
	}

	var createdEntries []*models.MemoryEntry
	for _, req := range entries {
		entry, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(t, err)
		createdEntries = append(createdEntries, entry)
	}

	// Test semantic search
	searchReq := &models.SearchRequest{
		Query:      "SQL injection",
		SessionID:  session.ID,
		SearchType: "semantic",
		Limit:      10,
	}

	results, err := repo.SearchMemoryEntries(ctx, searchReq)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "SQL Injection in Login", results[0].MemoryEntry.Title)
	assert.Greater(t, results[0].Relevance, 0.0)

	// Test category filter
	searchReq = &models.SearchRequest{
		Query:       "vulnerability",
		SessionID:   session.ID,
		Categories:  []string{"vulnerability"},
		SearchType:  "exact",
		Limit:       10,
	}

	results, err = repo.SearchMemoryEntries(ctx, searchReq)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Should find both vulnerability entries

	// Test priority filter
	searchReq = &models.SearchRequest{
		Query:        "authentication",
		SessionID:    session.ID,
		MinPriority:  7,
		SearchType:   "exact",
		Limit:        10,
	}

	results, err = repo.SearchMemoryEntries(ctx, searchReq)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Should find entries with priority >= 7
}

func TestCreateRelationship(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create two memory entries
	entry1Req := &models.CreateMemoryEntryRequest{
		SessionID: session.ID,
		Title:     "SQL Injection Vulnerability",
		Content:   "Found SQL injection in login form",
		Category:  "vulnerability",
		Priority:  8,
	}

	entry2Req := &models.CreateMemoryEntryRequest{
		SessionID: session.ID,
		Title:     "Authentication Bypass",
		Content:   "Found way to bypass authentication",
		Category:  "exploit",
		Priority:  9,
	}

	entry1, err := repo.CreateMemoryEntry(ctx, entry1Req)
	require.NoError(t, err)

	entry2, err := repo.CreateMemoryEntry(ctx, entry2Req)
	require.NoError(t, err)

	// Create relationship
	relReq := &models.CreateRelationshipRequest{
		SourceEntryID:    entry1.ID,
		TargetEntryID:    entry2.ID,
		RelationshipType: "exploits",
		Strength:         0.8,
		Description:      "SQL injection can be used to bypass authentication",
	}

	relationship, err := repo.CreateRelationship(ctx, relReq)
	assert.NoError(t, err)
	assert.NotNil(t, relationship)
	assert.Equal(t, entry1.ID, relationship.SourceEntryID)
	assert.Equal(t, entry2.ID, relationship.TargetEntryID)
	assert.Equal(t, "exploits", relationship.RelationshipType)
	assert.Equal(t, 0.8, relationship.Strength)
}

func TestGetRelatedEntries(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create multiple memory entries
	entries := []*models.CreateMemoryEntryRequest{
		{
			SessionID: session.ID,
			Title:     "SQL Injection Vulnerability",
			Content:   "Found SQL injection in login form",
			Category:  "vulnerability",
			Priority:  8,
		},
		{
			SessionID: session.ID,
			Title:     "Authentication Bypass",
			Content:   "Found way to bypass authentication",
			Category:  "exploit",
			Priority:  9,
		},
		{
			SessionID: session.ID,
			Title:     "Session Management Issue",
			Content:   "Weak session management allows hijacking",
			Category:  "vulnerability",
			Priority:  7,
		},
	}

	var createdEntries []*models.MemoryEntry
	for _, req := range entries {
		entry, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(t, err)
		createdEntries = append(createdEntries, entry)
	}

	// Create relationships
	relationships := []*models.CreateRelationshipRequest{
		{
			SourceEntryID:    createdEntries[0].ID, // SQL Injection
			TargetEntryID:    createdEntries[1].ID, // Authentication Bypass
			RelationshipType: "exploits",
			Strength:         0.8,
		},
		{
			SourceEntryID:    createdEntries[0].ID, // SQL Injection
			TargetEntryID:    createdEntries[2].ID, // Session Management
			RelationshipType: "related_to",
			Strength:         0.6,
		},
	}

	for _, relReq := range relationships {
		_, err := repo.CreateRelationship(ctx, relReq)
		require.NoError(t, err)
	}

	// Get related entries for SQL Injection
	related, err := repo.GetRelatedEntries(ctx, createdEntries[0].ID, "", 10)
	assert.NoError(t, err)
	assert.Len(t, related, 2) // Should find both related entries

	// Get related entries with specific relationship type
	related, err = repo.GetRelatedEntries(ctx, createdEntries[0].ID, "exploits", 10)
	assert.NoError(t, err)
	assert.Len(t, related, 1) // Should find only the exploit relationship
	assert.Equal(t, createdEntries[1].ID, related[0].ID)
}

func TestListSessions(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create multiple sessions
	sessions := []*models.Session{
		{
			ID:       "session-1",
			Name:     "Security Review 1",
			TaskType: "security_review",
			Status:   "active",
		},
		{
			ID:       "session-2",
			Name:     "Penetration Test 1",
			TaskType: "penetration_test",
			Status:   "active",
		},
		{
			ID:       "session-3",
			Name:     "Security Review 2",
			TaskType: "security_review",
			Status:   "completed",
		},
	}

	for _, session := range sessions {
		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)
	}

	// Test listing all sessions
	allSessions, err := repo.ListSessions(ctx, "", "", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, allSessions, 3)

	// Test filtering by task type
	securityReviews, err := repo.ListSessions(ctx, "security_review", "", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, securityReviews, 2)

	// Test filtering by status
	activeSessions, err := repo.ListSessions(ctx, "", "active", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, activeSessions, 2)

	// Test pagination
	paginatedSessions, err := repo.ListSessions(ctx, "", "", 2, 0)
	assert.NoError(t, err)
	assert.Len(t, paginatedSessions, 2)

	paginatedSessions, err = repo.ListSessions(ctx, "", "", 2, 2)
	assert.NoError(t, err)
	assert.Len(t, paginatedSessions, 1)
}

func TestAccessTracking(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create memory entry
	req := &models.CreateMemoryEntryRequest{
		SessionID: session.ID,
		Title:     "Test Entry",
		Content:   "Test content",
		Category:  "note",
	}

	entry, err := repo.CreateMemoryEntry(ctx, req)
	require.NoError(t, err)

	// Verify initial access count
	assert.Equal(t, 0, entry.AccessCount)

	// Retrieve entry multiple times
	for i := 0; i < 5; i++ {
		retrieved, err := repo.GetMemoryEntry(ctx, entry.ID)
		assert.NoError(t, err)
		assert.Equal(t, i+1, retrieved.AccessCount)
	}
}

func BenchmarkCreateMemoryEntry(b *testing.B) {
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(b, err)
	defer db.Close()

	repo := NewMemoryRepository(db.GetDB(), logger)

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "bench-session",
		Name:     "Benchmark Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err = repo.CreateSession(ctx, session)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &models.CreateMemoryEntryRequest{
			SessionID: session.ID,
			Title:     fmt.Sprintf("Benchmark Entry %d", i),
			Content:   "Benchmark content",
			Category:  "note",
			Priority:  5,
		}

		_, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(b, err)
	}
}

func BenchmarkSearchMemoryEntries(b *testing.B) {
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(b, err)
	defer db.Close()

	repo := NewMemoryRepository(db.GetDB(), logger)

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "bench-session",
		Name:     "Benchmark Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err = repo.CreateSession(ctx, session)
	require.NoError(b, err)

	// Create test entries
	for i := 0; i < 100; i++ {
		req := &models.CreateMemoryEntryRequest{
			SessionID: session.ID,
			Title:     fmt.Sprintf("Entry %d", i),
			Content:   fmt.Sprintf("Content for entry %d with some keywords", i),
			Category:  "note",
			Priority:  i % 10,
		}

		_, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(b, err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		searchReq := &models.SearchRequest{
			Query:      "keywords",
			SessionID:  session.ID,
			SearchType: "semantic",
			Limit:      20,
		}

		_, err := repo.SearchMemoryEntries(ctx, searchReq)
		require.NoError(b, err)
	}
}
