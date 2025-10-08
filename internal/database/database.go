package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/charmbracelet/log"
)

// Database represents the SQLite database connection and operations
type Database struct {
	db *sql.DB
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
	dsn := fmt.Sprintf("file:%s?cache=shared&_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=30000", dbPath)
	db, err := sql.Open("sqlite3", dsn)
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

// initializeSchema creates all tables and indexes
func (d *Database) initializeSchema() error {
	// Embedded schema for testing and deployment
	schema := `
-- Security-focused LLM Memory Storage Database Schema
-- Designed for security code review, penetration testing, and exploit development

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Enable WAL mode for better concurrency
PRAGMA journal_mode = WAL;

-- Enable full-text search
PRAGMA case_sensitive_like = OFF;

-- Sessions table - tracks LLM interaction sessions
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    task_type TEXT NOT NULL CHECK (task_type IN ('security_review', 'penetration_test', 'exploit_dev', 'vulnerability_analysis', 'threat_modeling', 'incident_response', 'general')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'archived')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT -- JSON metadata for session-specific data
);

-- Memory entries table - stores individual pieces of information
CREATE TABLE IF NOT EXISTS memory_entries (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'text' CHECK (content_type IN ('text', 'code', 'json', 'yaml', 'markdown', 'binary_ref')),
    category TEXT NOT NULL CHECK (category IN ('finding', 'vulnerability', 'exploit', 'payload', 'technique', 'tool', 'reference', 'context', 'hypothesis', 'evidence', 'recommendation', 'note')),
    priority INTEGER DEFAULT 0 CHECK (priority >= 0 AND priority <= 10), -- 0=low, 10=critical
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    tags TEXT, -- JSON array of tags
    source TEXT, -- Where this information came from
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    access_count INTEGER DEFAULT 0,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Relationships table - links related memory entries
CREATE TABLE IF NOT EXISTS relationships (
    id TEXT PRIMARY KEY,
    source_entry_id TEXT NOT NULL,
    target_entry_id TEXT NOT NULL,
    relationship_type TEXT NOT NULL CHECK (relationship_type IN ('depends_on', 'causes', 'mitigates', 'exploits', 'references', 'contradicts', 'supports', 'related_to', 'parent_of', 'child_of')),
    strength REAL DEFAULT 0.5 CHECK (strength >= 0.0 AND strength <= 1.0),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_entry_id) REFERENCES memory_entries(id) ON DELETE CASCADE,
    FOREIGN KEY (target_entry_id) REFERENCES memory_entries(id) ON DELETE CASCADE,
    UNIQUE(source_entry_id, target_entry_id, relationship_type)
);

-- Context snapshots - stores context state at specific points
CREATE TABLE IF NOT EXISTS context_snapshots (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    context_data TEXT NOT NULL, -- JSON representation of current context
    memory_summary TEXT, -- Summary of relevant memories
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Search history - tracks what the LLM has searched for
CREATE TABLE IF NOT EXISTS search_history (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    query TEXT NOT NULL,
    search_type TEXT NOT NULL CHECK (search_type IN ('semantic', 'exact', 'fuzzy', 'tag', 'category', 'relationship')),
    results_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Task progress - tracks multi-stage task progress
CREATE TABLE IF NOT EXISTS task_progress (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    task_name TEXT NOT NULL,
    stage TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'failed', 'blocked')),
    progress_percentage INTEGER DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    notes TEXT,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_memory_entries_session_id ON memory_entries(session_id);
CREATE INDEX IF NOT EXISTS idx_memory_entries_category ON memory_entries(category);
CREATE INDEX IF NOT EXISTS idx_memory_entries_priority ON memory_entries(priority);
CREATE INDEX IF NOT EXISTS idx_memory_entries_created_at ON memory_entries(created_at);
CREATE INDEX IF NOT EXISTS idx_memory_entries_accessed_at ON memory_entries(accessed_at);
CREATE INDEX IF NOT EXISTS idx_memory_entries_access_count ON memory_entries(access_count);

CREATE INDEX IF NOT EXISTS idx_relationships_source ON relationships(source_entry_id);
CREATE INDEX IF NOT EXISTS idx_relationships_target ON relationships(target_entry_id);
CREATE INDEX IF NOT EXISTS idx_relationships_type ON relationships(relationship_type);

CREATE INDEX IF NOT EXISTS idx_context_snapshots_session_id ON context_snapshots(session_id);
CREATE INDEX IF NOT EXISTS idx_search_history_session_id ON search_history(session_id);
CREATE INDEX IF NOT EXISTS idx_task_progress_session_id ON task_progress(session_id);
CREATE INDEX IF NOT EXISTS idx_task_progress_status ON task_progress(status);

-- Note: FTS5 virtual table creation is handled separately to avoid errors
-- if FTS5 is not available in the SQLite build

-- Create view for easy access to memory entries with session info
CREATE VIEW IF NOT EXISTS memory_entries_with_session AS
SELECT 
    me.*,
    s.name as session_name,
    s.task_type,
    s.status as session_status
FROM memory_entries me
JOIN sessions s ON me.session_id = s.id;

-- Create view for relationship analysis
CREATE VIEW IF NOT EXISTS relationship_network AS
SELECT 
    r.*,
    s1.title as source_title,
    s1.category as source_category,
    t1.title as target_title,
    t1.category as target_category
FROM relationships r
JOIN memory_entries s1 ON r.source_entry_id = s1.id
JOIN memory_entries t1 ON r.target_entry_id = t1.id;
`

	// Execute schema
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Try to create FTS5 table if available
	d.createFTS5Table()

	d.logger.Debug("Database schema initialized")
	return nil
}

// createFTS5Table creates the FTS5 virtual table if available
func (d *Database) createFTS5Table() {
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
`

	if _, err := d.db.Exec(fts5Schema); err != nil {
		d.logger.Warn("FTS5 not available, using regular search", "error", err)
	} else {
		d.logger.Debug("FTS5 table created successfully")
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
		"sessions_count":         "SELECT COUNT(*) FROM sessions",
		"memory_entries_count":   "SELECT COUNT(*) FROM memory_entries",
		"relationships_count":    "SELECT COUNT(*) FROM relationships",
		"context_snapshots_count": "SELECT COUNT(*) FROM context_snapshots",
		"search_history_count":   "SELECT COUNT(*) FROM search_history",
		"task_progress_count":    "SELECT COUNT(*) FROM task_progress",
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
