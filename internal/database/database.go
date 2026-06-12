package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite" // pure-Go SQLite driver (no cgo), includes FTS5
)

// Database represents the SQLite database connection and operations
type Database struct {
	db     *sql.DB
	logger *log.Logger
}

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase(dbPath string, logger *log.Logger) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database with optimized settings for security tasks
	dsn := fmt.Sprintf("file:%s?cache=shared&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(30000)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for security tasks
	db.SetMaxOpenConns(1) // SQLite doesn't benefit from multiple connections
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // Keep connections alive

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{
		db:     db,
		logger: logger,
	}

	// Initialize schema
	if err := database.initializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.Info("Database initialized successfully", "path", dbPath)
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

//go:embed schema.sql
var schemaSQL string

// initializeSchema creates all tables, indexes, and views from the embedded
// schema.sql, which is the single source of truth for the database schema.
func (d *Database) initializeSchema() error {
	if _, err := d.db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Try to create FTS5 virtual tables if this SQLite build supports them
	d.createFTS5Table()

	d.logger.Debug("Database schema initialized")
	return nil
}

// createFTS5Table creates the FTS5 virtual table if available
func (d *Database) createFTS5Table() {
	// Check if FTS5 is available before trying to create virtual table
	var fts5Available bool
	err := d.db.QueryRow("SELECT 1 FROM pragma_compile_options WHERE compile_options LIKE '%FTS5%'").Scan(&fts5Available)
	if err != nil {
		// FTS5 not available, skip virtual table creation
		d.logger.Info("FTS5 not compiled in SQLite, using regular search")
		return
	}

	// Try to create FTS5 table
	fts5Schema := `
CREATE VIRTUAL TABLE IF NOT EXISTS memory_entries_fts USING fts5(
    title,
    content,
    tags,
    content='memory_entries',
    content_rowid='rowid'
);

-- Create triggers to keep FTS table in sync
CREATE TRIGGER IF NOT EXISTS memory_entries_fts_insert AFTER INSERT ON memory_entries BEGIN
    INSERT INTO memory_entries_fts(rowid, title, content, tags) 
    VALUES (new.rowid, new.title, new.content, new.tags);
END;

CREATE TRIGGER IF NOT EXISTS memory_entries_fts_delete AFTER DELETE ON memory_entries BEGIN
    INSERT INTO memory_entries_fts(memory_entries_fts, rowid, title, content, tags) 
    VALUES('delete', old.rowid, old.title, old.content, old.tags);
END;

CREATE TRIGGER IF NOT EXISTS memory_entries_fts_update AFTER UPDATE ON memory_entries BEGIN
    INSERT INTO memory_entries_fts(memory_entries_fts, rowid, title, content, tags) 
    VALUES('delete', old.rowid, old.title, old.content, old.tags);
    INSERT INTO memory_entries_fts(rowid, title, content, tags) 
    VALUES (new.rowid, new.title, new.content, new.tags);
END;

-- FTS5 virtual tables for security data (if available)
CREATE VIRTUAL TABLE IF NOT EXISTS nvd_cves_fts USING fts5(
    id, description, cwe_ids,
    content='nvd_cves',
    content_rowid='rowid'
);

CREATE VIRTUAL TABLE IF NOT EXISTS attack_techniques_fts USING fts5(
    id, name, description, tactic,
    content='attack_techniques',
    content_rowid='rowid'
);

CREATE VIRTUAL TABLE IF NOT EXISTS owasp_procedures_fts USING fts5(
    id, title, description, category,
    content='owasp_procedures',
    content_rowid='rowid'
);
`

	if _, err := d.db.Exec(fts5Schema); err != nil {
		d.logger.Info("FTS5 virtual table creation failed, using regular search", "error", err)
	} else {
		d.logger.Info("FTS5 full-text search initialized")
	}
}

// GetDB returns the underlying sql.DB for direct operations
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// BeginTransaction starts a new transaction
func (d *Database) BeginTransaction() (*sql.Tx, error) {
	return d.db.Begin()
}

// ExecuteInTransaction executes a function within a transaction
func (d *Database) ExecuteInTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.logger.Error("Failed to rollback transaction", "error", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// HealthCheck performs a basic health check on the database
func (d *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Check if we can query a simple table
	var count int
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions").Scan(&count); err != nil {
		return fmt.Errorf("database query test failed: %w", err)
	}

	return nil
}

// GetStats returns database statistics
func (d *Database) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get table counts - using hardcoded table names to avoid SQL injection
	tableQueries := map[string]string{
		"sessions_count":          "SELECT COUNT(*) FROM sessions",
		"memory_entries_count":    "SELECT COUNT(*) FROM memory_entries",
		"relationships_count":     "SELECT COUNT(*) FROM relationships",
		"context_snapshots_count": "SELECT COUNT(*) FROM context_snapshots",
		"search_history_count":    "SELECT COUNT(*) FROM search_history",
		"task_progress_count":     "SELECT COUNT(*) FROM task_progress",
	}

	for statName, query := range tableQueries {
		var count int
		if err := d.db.QueryRow(query).Scan(&count); err != nil {
			return nil, fmt.Errorf("failed to get count for %s: %w", statName, err)
		}
		stats[statName] = count
	}

	// Get database size
	var pageCount, pageSize int
	if err := d.db.QueryRow("PRAGMA page_count").Scan(&pageCount); err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}
	if err := d.db.QueryRow("PRAGMA page_size").Scan(&pageSize); err != nil {
		return nil, fmt.Errorf("failed to get page size: %w", err)
	}
	stats["database_size_bytes"] = pageCount * pageSize

	// Get most accessed entries
	var topAccessed []map[string]interface{}
	rows, err := d.db.Query(`
		SELECT title, access_count, category, priority 
		FROM memory_entries 
		ORDER BY access_count DESC 
		LIMIT 5
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get top accessed entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var accessCount int
		var category string
		var priority int
		if err := rows.Scan(&title, &accessCount, &category, &priority); err != nil {
			return nil, fmt.Errorf("failed to scan top accessed entry: %w", err)
		}
		topAccessed = append(topAccessed, map[string]interface{}{
			"title":        title,
			"access_count": accessCount,
			"category":     category,
			"priority":     priority,
		})
	}
	stats["top_accessed_entries"] = topAccessed

	return stats, nil
}
