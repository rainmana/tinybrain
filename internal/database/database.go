package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
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
    relationship_type TEXT NOT NULL CHECK (relationship_type IN ('depends_on', 'causes', 'mitigates', 'exploits', 'refs', 'contradicts', 'supports', 'related_to', 'parent_of', 'child_of')),
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

-- Create notifications table for real-time alerts
CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    memory_id TEXT,
    session_id TEXT,
    priority INTEGER NOT NULL DEFAULT 5,
    metadata TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (memory_id) REFERENCES memory_entries(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- CVE mappings table - maps CWE to CVE entries
CREATE TABLE IF NOT EXISTS cve_mappings (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    cwe_id TEXT NOT NULL,
    cve_list TEXT NOT NULL, -- JSON array of CVE IDs
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confidence REAL NOT NULL DEFAULT 0.0,
    source TEXT NOT NULL DEFAULT 'nvd',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Risk correlations table - stores vulnerability correlation analysis
CREATE TABLE IF NOT EXISTS risk_correlations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    primary_vuln_id TEXT NOT NULL,
    secondary_vuln_ids TEXT NOT NULL, -- JSON array of vulnerability IDs
    risk_multiplier REAL NOT NULL DEFAULT 1.0,
    attack_chain TEXT NOT NULL, -- JSON array of attack steps
    business_impact TEXT NOT NULL,
    confidence REAL NOT NULL DEFAULT 0.0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (primary_vuln_id) REFERENCES memory_entries(id) ON DELETE CASCADE
);

-- Compliance mappings table - maps vulnerabilities to compliance standards
CREATE TABLE IF NOT EXISTS compliance_mappings (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    standard TEXT NOT NULL,
    requirement TEXT NOT NULL,
    vulnerability_ids TEXT NOT NULL, -- JSON array of vulnerability IDs
    compliance_score REAL NOT NULL DEFAULT 0.0,
    gap_analysis TEXT NOT NULL, -- JSON array of gaps
    recommendations TEXT NOT NULL, -- JSON array of recommendations
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

CREATE INDEX IF NOT EXISTS idx_notifications_session_id ON notifications(session_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Indexes for new tables
CREATE INDEX IF NOT EXISTS idx_cve_mappings_session_id ON cve_mappings(session_id);
CREATE INDEX IF NOT EXISTS idx_cve_mappings_cwe_id ON cve_mappings(cwe_id);
CREATE INDEX IF NOT EXISTS idx_cve_mappings_confidence ON cve_mappings(confidence);

CREATE INDEX IF NOT EXISTS idx_risk_correlations_session_id ON risk_correlations(session_id);
CREATE INDEX IF NOT EXISTS idx_risk_correlations_primary_vuln_id ON risk_correlations(primary_vuln_id);
CREATE INDEX IF NOT EXISTS idx_risk_correlations_confidence ON risk_correlations(confidence);

CREATE INDEX IF NOT EXISTS idx_compliance_mappings_session_id ON compliance_mappings(session_id);
CREATE INDEX IF NOT EXISTS idx_compliance_mappings_standard ON compliance_mappings(standard);
CREATE INDEX IF NOT EXISTS idx_compliance_mappings_compliance_score ON compliance_mappings(compliance_score);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read);

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

-- Security Knowledge Hub Tables

-- NVD CVEs table
CREATE TABLE IF NOT EXISTS nvd_cves (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    cvss_v2_score REAL,
    cvss_v2_vector TEXT,
    cvss_v3_score REAL,
    cvss_v3_vector TEXT,
    severity TEXT,
    published_date DATETIME,
    last_modified_date DATETIME,
    cwe_ids TEXT, -- JSON array
    affected_products TEXT, -- JSON array
    refs TEXT, -- JSON array
    raw_data TEXT, -- JSON raw data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ATT&CK Techniques table
CREATE TABLE IF NOT EXISTS attack_techniques (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    tactic TEXT,
    tactics TEXT, -- JSON array
    platforms TEXT, -- JSON array
    kill_chain_phases TEXT, -- JSON array
    data_sources TEXT, -- JSON array
    detection TEXT,
    mitigation TEXT,
    refs TEXT, -- JSON array
    sub_techniques TEXT, -- JSON array
    parent_technique TEXT,
    raw_data TEXT, -- JSON raw data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ATT&CK Tactics table
CREATE TABLE IF NOT EXISTS attack_tactics (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    kill_chain_phases TEXT, -- JSON array
    techniques TEXT, -- JSON array
    refs TEXT, -- JSON array
    raw_data TEXT, -- JSON raw data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- OWASP Procedures table
CREATE TABLE IF NOT EXISTS owasp_procedures (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    subcategory TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    objective TEXT,
    how_to_test TEXT,
    tools TEXT, -- JSON array
    refs TEXT, -- JSON array
    severity TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Security Data Updates table
CREATE TABLE IF NOT EXISTS security_data_updates (
    id TEXT PRIMARY KEY,
    data_source TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    total_records INTEGER,
    error_message TEXT,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for security tables
CREATE INDEX IF NOT EXISTS idx_nvd_cves_severity ON nvd_cves(severity);
CREATE INDEX IF NOT EXISTS idx_nvd_cves_published ON nvd_cves(published_date);
CREATE INDEX IF NOT EXISTS idx_nvd_cves_cvss3 ON nvd_cves(cvss_v3_score);

CREATE INDEX IF NOT EXISTS idx_attack_techniques_tactic ON attack_techniques(tactic);
CREATE INDEX IF NOT EXISTS idx_attack_techniques_platforms ON attack_techniques(platforms);

CREATE INDEX IF NOT EXISTS idx_owasp_procedures_category ON owasp_procedures(category);
CREATE INDEX IF NOT EXISTS idx_owasp_procedures_severity ON owasp_procedures(severity);

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
