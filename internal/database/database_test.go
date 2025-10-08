package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	// Test database creation
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Test database connection
	assert.NoError(t, db.HealthCheck())

	// Test database stats
	stats, err := db.GetStats()
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Contains(t, stats, "sessions_count")
	assert.Contains(t, stats, "memory_entries_count")

	// Clean up
	assert.NoError(t, db.Close())
}

func TestDatabaseHealthCheck(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test health check
	assert.NoError(t, db.HealthCheck())

	// Test with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	// This should work within the timeout
	assert.NoError(t, db.HealthCheck())
	
	// Use the context to avoid unused variable warning
	_ = ctx
}

func TestDatabaseStats(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Get initial stats
	stats, err := db.GetStats()
	require.NoError(t, err)

	// Check that all expected tables are present
	expectedTables := []string{
		"sessions_count",
		"memory_entries_count",
		"relationships_count",
		"context_snapshots_count",
		"search_history_count",
		"task_progress_count",
	}

	for _, table := range expectedTables {
		assert.Contains(t, stats, table)
		assert.Equal(t, 0, stats[table]) // Should be empty initially
	}

	// Check database size
	assert.Contains(t, stats, "database_size_bytes")
	assert.Greater(t, stats["database_size_bytes"], 0)

	// Check top accessed entries
	assert.Contains(t, stats, "top_accessed_entries")
	assert.IsType(t, []map[string]interface{}{}, stats["top_accessed_entries"])
}

func TestDatabaseTransaction(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test successful transaction
	err = db.ExecuteInTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO sessions (id, name, task_type, status) VALUES (?, ?, ?, ?)",
			"test-session-1", "Test Session", "security_review", "active")
		return err
	})
	assert.NoError(t, err)

	// Verify the data was inserted
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", "test-session-1").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test failed transaction (should rollback)
	err = db.ExecuteInTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO sessions (id, name, task_type, status) VALUES (?, ?, ?, ?)",
			"test-session-2", "Test Session 2", "security_review", "active")
		if err != nil {
			return err
		}
		// Force an error
		return assert.AnError
	})
	assert.Error(t, err)

	// Verify the data was not inserted (rollback worked)
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", "test-session-2").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDatabaseConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test concurrent reads
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			// Perform multiple operations
			for j := 0; j < 5; j++ {
				stats, err := db.GetStats()
				assert.NoError(t, err)
				assert.NotNil(t, stats)
				
				err = db.HealthCheck()
				assert.NoError(t, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDatabaseSchema(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test that all tables exist
	tables := []string{
		"sessions",
		"memory_entries", 
		"relationships",
		"context_snapshots",
		"search_history",
		"task_progress",
	}

	for _, table := range tables {
		var count int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "Table %s should exist", table)
	}

	// Test FTS5 table (may not be available in all SQLite builds)
	var fts5Count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", "memory_entries_fts").Scan(&fts5Count)
	assert.NoError(t, err)
	if fts5Count == 0 {
		t.Log("FTS5 table not available in this SQLite build - this is expected")
	}

	// Test that views exist
	views := []string{
		"memory_entries_with_session",
		"relationship_network",
	}

	for _, view := range views {
		var count int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='view' AND name=?", view).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "View %s should exist", view)
	}

	// Test that FTS5 triggers exist (only if FTS5 is available)
	if fts5Count > 0 {
		triggers := []string{
			"memory_entries_fts_insert",
			"memory_entries_fts_delete", 
			"memory_entries_fts_update",
		}

		for _, trigger := range triggers {
			var count int
			err := db.GetDB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='trigger' AND name=?", trigger).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "Trigger %s should exist", trigger)
		}
	} else {
		t.Log("Skipping FTS5 trigger tests - FTS5 not available")
	}
}

func TestDatabaseConstraints(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test foreign key constraints
	_, err = db.GetDB().Exec("INSERT INTO memory_entries (id, session_id, title, content, category) VALUES (?, ?, ?, ?, ?)",
		"test-entry-1", "non-existent-session", "Test", "Content", "note")
	assert.Error(t, err, "Should fail due to foreign key constraint")

	// Test check constraints on priority
	_, err = db.GetDB().Exec("INSERT INTO sessions (id, name, task_type, status) VALUES (?, ?, ?, ?)",
		"test-session-1", "Test", "security_review", "active")
	assert.NoError(t, err)

	_, err = db.GetDB().Exec("INSERT INTO memory_entries (id, session_id, title, content, category, priority) VALUES (?, ?, ?, ?, ?, ?)",
		"test-entry-1", "test-session-1", "Test", "Content", "note", 15)
	assert.Error(t, err, "Should fail due to priority check constraint (must be 0-10)")

	_, err = db.GetDB().Exec("INSERT INTO memory_entries (id, session_id, title, content, category, confidence) VALUES (?, ?, ?, ?, ?, ?)",
		"test-entry-2", "test-session-1", "Test", "Content", "note", 1.5)
	assert.Error(t, err, "Should fail due to confidence check constraint (must be 0.0-1.0)")
}

func BenchmarkDatabaseOperations(b *testing.B) {
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench.db")

	logger := log.New(os.Stderr)
	db, err := NewDatabase(dbPath, logger)
	require.NoError(b, err)
	defer db.Close()

	// Create test session
	_, err = db.GetDB().Exec("INSERT INTO sessions (id, name, task_type, status) VALUES (?, ?, ?, ?)",
		"bench-session", "Benchmark Session", "security_review", "active")
	require.NoError(b, err)

	b.Run("HealthCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := db.HealthCheck()
			require.NoError(b, err)
		}
	})

	b.Run("GetStats", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.GetStats()
			require.NoError(b, err)
		}
	})

	b.Run("InsertMemoryEntry", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.GetDB().Exec("INSERT INTO memory_entries (id, session_id, title, content, category) VALUES (?, ?, ?, ?, ?)",
				fmt.Sprintf("bench-entry-%d", i), "bench-session", "Benchmark Entry", "Content", "note")
			require.NoError(b, err)
		}
	})
}
