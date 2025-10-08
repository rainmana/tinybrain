package repository

import (
	"context"
	"testing"

	"github.com/alec/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateContextSnapshot(t *testing.T) {
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

	// Create some memory entries for the summary
	entries := []*models.CreateMemoryEntryRequest{
		{
			SessionID: session.ID,
			Title:     "High Priority Finding",
			Content:   "This is a critical security issue",
			Category:  "vulnerability",
			Priority:  9,
			Confidence: 0.9,
		},
		{
			SessionID: session.ID,
			Title:     "Medium Priority Finding",
			Content:   "This is a moderate security issue",
			Category:  "finding",
			Priority:  5,
			Confidence: 0.7,
		},
	}

	for _, req := range entries {
		_, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(t, err)
	}

	// Create context snapshot
	contextData := map[string]interface{}{
		"current_stage": "assessment",
		"findings":      []string{"vuln1", "vuln2"},
		"next_steps":    []string{"validate", "remediate"},
	}

	snapshot, err := repo.CreateContextSnapshot(ctx, session.ID, "Test Snapshot", "Test description", contextData)
	require.NoError(t, err)

	assert.NotEmpty(t, snapshot.ID)
	assert.Equal(t, session.ID, snapshot.SessionID)
	assert.Equal(t, "Test Snapshot", snapshot.Name)
	assert.Equal(t, "Test description", snapshot.Description)
	assert.Equal(t, contextData, snapshot.ContextData)
	assert.Contains(t, snapshot.MemorySummary, "High Priority Finding")
	assert.Contains(t, snapshot.MemorySummary, "Priority: 9")
}

func TestGetContextSnapshot(t *testing.T) {
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

	// Create context snapshot
	contextData := map[string]interface{}{
		"test_key": "test_value",
	}

	snapshot, err := repo.CreateContextSnapshot(ctx, session.ID, "Test Snapshot", "Test description", contextData)
	require.NoError(t, err)

	// Retrieve the snapshot
	retrieved, err := repo.GetContextSnapshot(ctx, snapshot.ID)
	require.NoError(t, err)

	assert.Equal(t, snapshot.ID, retrieved.ID)
	assert.Equal(t, snapshot.SessionID, retrieved.SessionID)
	assert.Equal(t, snapshot.Name, retrieved.Name)
	assert.Equal(t, snapshot.Description, retrieved.Description)
	assert.Equal(t, snapshot.ContextData, retrieved.ContextData)
	assert.Equal(t, snapshot.MemorySummary, retrieved.MemorySummary)
}

func TestListContextSnapshots(t *testing.T) {
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

	// Create multiple snapshots
	for i := 0; i < 3; i++ {
		contextData := map[string]interface{}{
			"snapshot_number": i,
		}
		_, err := repo.CreateContextSnapshot(ctx, session.ID, 
			"Snapshot "+string(rune('A'+i)), 
			"Description "+string(rune('A'+i)), 
			contextData)
		require.NoError(t, err)
	}

	// List snapshots
	snapshots, err := repo.ListContextSnapshots(ctx, session.ID, 10, 0)
	require.NoError(t, err)

	assert.Len(t, snapshots, 3)
	// Should be ordered by created_at DESC
	assert.Equal(t, "Snapshot C", snapshots[0].Name)
	assert.Equal(t, "Snapshot B", snapshots[1].Name)
	assert.Equal(t, "Snapshot A", snapshots[2].Name)
}

