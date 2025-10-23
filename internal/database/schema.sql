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

-- Create full-text search virtual table for memory entries
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

-- Create view for easy access to memory entries with session info
CREATE VIEW IF NOT EXISTS memory_entries_with_session AS
SELECT 
    me.*,
    s.name as session_name,
    s.task_type,
    s.status as session_status
FROM memory_entries me
JOIN sessions s ON me.session_id = s.id;

-- Security Knowledge Hub Tables

-- NVD CVE data table
CREATE TABLE IF NOT EXISTS nvd_cves (
    id TEXT PRIMARY KEY, -- CVE ID (e.g., CVE-2024-1234)
    description TEXT NOT NULL,
    cvss_v2_score REAL,
    cvss_v2_vector TEXT,
    cvss_v3_score REAL,
    cvss_v3_vector TEXT,
    severity TEXT CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    published_date DATETIME,
    last_modified_date DATETIME,
    cwe_ids TEXT, -- JSON array of CWE IDs
    affected_products TEXT, -- JSON array of affected products
    references TEXT, -- JSON array of references
    raw_data TEXT, -- Full JSON data from NVD
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- MITRE ATT&CK techniques table
CREATE TABLE IF NOT EXISTS attack_techniques (
    id TEXT PRIMARY KEY, -- Technique ID (e.g., T1055.001)
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    tactic TEXT NOT NULL, -- Primary tactic
    tactics TEXT, -- JSON array of all tactics
    platforms TEXT, -- JSON array of platforms
    kill_chain_phases TEXT, -- JSON array of kill chain phases
    data_sources TEXT, -- JSON array of data sources
    detection TEXT, -- Detection guidance
    mitigation TEXT, -- Mitigation guidance
    references TEXT, -- JSON array of references
    sub_techniques TEXT, -- JSON array of sub-technique IDs
    parent_technique TEXT, -- Parent technique ID if this is a sub-technique
    raw_data TEXT, -- Full STIX object data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- MITRE ATT&CK tactics table
CREATE TABLE IF NOT EXISTS attack_tactics (
    id TEXT PRIMARY KEY, -- Tactic ID
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    external_id TEXT, -- External reference ID
    kill_chain_phases TEXT, -- JSON array of kill chain phases
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- OWASP testing procedures table
CREATE TABLE IF NOT EXISTS owasp_procedures (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL, -- Testing category
    subcategory TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    objective TEXT, -- Testing objective
    how_to_test TEXT, -- How to test procedure
    tools TEXT, -- JSON array of recommended tools
    references TEXT, -- JSON array of references
    severity TEXT CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Security data update tracking
CREATE TABLE IF NOT EXISTS security_data_updates (
    id TEXT PRIMARY KEY,
    data_source TEXT NOT NULL CHECK (data_source IN ('nvd', 'attack', 'owasp')),
    last_update DATETIME,
    total_records INTEGER,
    update_status TEXT CHECK (update_status IN ('pending', 'in_progress', 'completed', 'failed')),
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for security data tables
CREATE INDEX IF NOT EXISTS idx_nvd_cves_severity ON nvd_cves(severity);
CREATE INDEX IF NOT EXISTS idx_nvd_cves_published ON nvd_cves(published_date);
CREATE INDEX IF NOT EXISTS idx_nvd_cves_cvss3 ON nvd_cves(cvss_v3_score);
CREATE INDEX IF NOT EXISTS idx_attack_techniques_tactic ON attack_techniques(tactic);
CREATE INDEX IF NOT EXISTS idx_attack_techniques_platforms ON attack_techniques(platforms);
CREATE INDEX IF NOT EXISTS idx_owasp_procedures_category ON owasp_procedures(category);
CREATE INDEX IF NOT EXISTS idx_owasp_procedures_severity ON owasp_procedures(severity);

-- Full-text search for security data (if FTS5 is available)
CREATE VIRTUAL TABLE IF NOT EXISTS nvd_cves_fts USING fts5(
    id, description, affected_products,
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
