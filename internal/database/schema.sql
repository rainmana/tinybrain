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
    task_type TEXT NOT NULL CHECK (task_type IN ('security_review', 'penetration_test', 'exploit_dev', 'vulnerability_analysis', 'threat_modeling', 'incident_response', 'intelligence_analysis', 'reverse_engineering', 'general')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'archived')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT, -- JSON metadata for session-specific data
    -- Enhanced intelligence fields
    intelligence_type TEXT CHECK (intelligence_type IN ('osint', 'humint', 'sigint', 'geoint', 'masint', 'techint', 'finint', 'cybint', 'mixed')),
    target_scope TEXT CHECK (target_scope IN ('individual', 'organization', 'infrastructure', 'campaign')),
    classification TEXT CHECK (classification IN ('unclassified', 'confidential', 'secret', 'top_secret')),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    geographic_scope TEXT CHECK (geographic_scope IN ('local', 'regional', 'national', 'international'))
);

-- Memory entries table - stores individual pieces of information
CREATE TABLE IF NOT EXISTS memory_entries (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'text' CHECK (content_type IN ('text', 'code', 'json', 'yaml', 'markdown', 'binary_ref', 'ioc', 'ttp', 'campaign', 'threat_actor', 'pattern', 'correlation')),
    category TEXT NOT NULL CHECK (category IN ('finding', 'vulnerability', 'exploit', 'payload', 'technique', 'tool', 'reference', 'context', 'hypothesis', 'evidence', 'recommendation', 'note', 'intelligence', 'ioc', 'ttp', 'campaign', 'threat_actor', 'pattern', 'correlation', 'reverse_engineering')),
    priority INTEGER DEFAULT 0 CHECK (priority >= 0 AND priority <= 10), -- 0=low, 10=critical
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    tags TEXT, -- JSON array of tags
    source TEXT, -- Where this information came from
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    access_count INTEGER DEFAULT 0,
    -- Enhanced intelligence fields
    intelligence_type TEXT CHECK (intelligence_type IN ('osint', 'humint', 'sigint', 'geoint', 'masint', 'techint', 'finint', 'cybint', 'mixed')),
    classification TEXT CHECK (classification IN ('unclassified', 'confidential', 'secret', 'top_secret')),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    geographic_scope TEXT CHECK (geographic_scope IN ('local', 'regional', 'national', 'international')),
    attribution TEXT, -- threat actor attribution
    ioc_type TEXT CHECK (ioc_type IN ('ip', 'domain', 'url', 'hash', 'email', 'file', 'registry', 'mutex', 'service')),
    ioc_value TEXT, -- actual IOC value
    mitre_tactic TEXT, -- MITRE ATT&CK tactic
    mitre_technique TEXT, -- MITRE ATT&CK technique
    mitre_procedure TEXT, -- MITRE ATT&CK procedure
    kill_chain_phase TEXT CHECK (kill_chain_phase IN ('reconnaissance', 'weaponization', 'delivery', 'exploitation', 'installation', 'c2', 'actions')),
    risk_score REAL DEFAULT 0.0 CHECK (risk_score >= 0.0 AND risk_score <= 10.0),
    impact_score REAL DEFAULT 0.0 CHECK (impact_score >= 0.0 AND impact_score <= 10.0),
    likelihood_score REAL DEFAULT 0.0 CHECK (likelihood_score >= 0.0 AND likelihood_score <= 10.0),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Relationships table - links related memory entries
CREATE TABLE IF NOT EXISTS relationships (
    id TEXT PRIMARY KEY,
    source_entry_id TEXT NOT NULL,
    target_entry_id TEXT NOT NULL,
    relationship_type TEXT NOT NULL CHECK (relationship_type IN ('depends_on', 'causes', 'mitigates', 'exploits', 'references', 'contradicts', 'supports', 'related_to', 'parent_of', 'child_of', 'correlates_with', 'attributed_to', 'part_of_campaign', 'uses_technique', 'targets', 'indicates')),
    strength REAL DEFAULT 0.5 CHECK (strength >= 0.0 AND strength <= 1.0),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- Enhanced intelligence relationship fields
    correlation_type TEXT CHECK (correlation_type IN ('temporal', 'spatial', 'logical', 'statistical', 'causal', 'predictive')),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    evidence TEXT, -- evidence supporting the relationship
    direction TEXT CHECK (direction IN ('unidirectional', 'bidirectional')),
    weight REAL DEFAULT 1.0 CHECK (weight >= 0.0), -- relationship weight for analysis
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

-- Intelligence-specific tables

-- Intelligence findings table - specialized findings with enhanced metadata
CREATE TABLE IF NOT EXISTS intelligence_findings (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    intelligence_type TEXT CHECK (intelligence_type IN ('osint', 'humint', 'sigint', 'geoint', 'masint', 'techint', 'finint', 'cybint', 'mixed')),
    classification TEXT CHECK (classification IN ('unclassified', 'confidential', 'secret', 'top_secret')),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    geographic_scope TEXT CHECK (geographic_scope IN ('local', 'regional', 'national', 'international')),
    attribution TEXT, -- threat actor attribution
    ioc_type TEXT CHECK (ioc_type IN ('ip', 'domain', 'url', 'hash', 'email', 'file', 'registry', 'mutex', 'service')),
    ioc_value TEXT, -- actual IOC value
    mitre_tactic TEXT, -- MITRE ATT&CK tactic
    mitre_technique TEXT, -- MITRE ATT&CK technique
    mitre_procedure TEXT, -- MITRE ATT&CK procedure
    kill_chain_phase TEXT CHECK (kill_chain_phase IN ('reconnaissance', 'weaponization', 'delivery', 'exploitation', 'installation', 'c2', 'actions')),
    risk_score REAL DEFAULT 0.0 CHECK (risk_score >= 0.0 AND risk_score <= 10.0),
    impact_score REAL DEFAULT 0.0 CHECK (impact_score >= 0.0 AND impact_score <= 10.0),
    likelihood_score REAL DEFAULT 0.0 CHECK (likelihood_score >= 0.0 AND likelihood_score <= 10.0),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    priority INTEGER DEFAULT 0 CHECK (priority >= 0 AND priority <= 10),
    tags TEXT, -- JSON array of tags
    source TEXT,
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Threat actors table - tracks threat actors and groups
CREATE TABLE IF NOT EXISTS threat_actors (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    aliases TEXT, -- JSON array of aliases
    description TEXT,
    motivation TEXT,
    capabilities TEXT, -- JSON array of capabilities
    targets TEXT, -- JSON array of targets
    tools TEXT, -- JSON array of tools
    techniques TEXT, -- JSON array of techniques
    attribution TEXT,
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    geographic_scope TEXT CHECK (geographic_scope IN ('local', 'regional', 'national', 'international')),
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Attack campaigns table - tracks attack campaigns and operations
CREATE TABLE IF NOT EXISTS attack_campaigns (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    threat_actors TEXT, -- JSON array of threat actor IDs
    targets TEXT, -- JSON array of targets
    techniques TEXT, -- JSON array of techniques
    tools TEXT, -- JSON array of tools
    iocs TEXT, -- JSON array of IOC IDs
    start_date DATETIME,
    end_date DATETIME,
    status TEXT CHECK (status IN ('active', 'inactive', 'completed', 'suspended')),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    geographic_scope TEXT CHECK (geographic_scope IN ('local', 'regional', 'national', 'international')),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Indicators of Compromise table - specialized IOC storage
CREATE TABLE IF NOT EXISTS indicators_of_compromise (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('ip', 'domain', 'url', 'hash', 'email', 'file', 'registry', 'mutex', 'service')),
    value TEXT NOT NULL,
    description TEXT,
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    source TEXT,
    attribution TEXT,
    campaigns TEXT, -- JSON array of campaign IDs
    techniques TEXT, -- JSON array of techniques
    tags TEXT, -- JSON array of tags
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Patterns table - security patterns and behaviors
CREATE TABLE IF NOT EXISTS patterns (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    pattern_type TEXT CHECK (pattern_type IN ('behavioral', 'attack', 'defense', 'vulnerability', 'exploit', 'ioc', 'ttp')),
    category TEXT,
    severity TEXT CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    frequency INTEGER DEFAULT 1,
    examples TEXT, -- JSON array of examples
    mitigations TEXT, -- JSON array of mitigations
    detections TEXT, -- JSON array of detection methods
    tags TEXT, -- JSON array of tags
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Correlations table - correlations between findings
CREATE TABLE IF NOT EXISTS correlations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    source_finding_id TEXT NOT NULL,
    target_finding_id TEXT NOT NULL,
    correlation_type TEXT CHECK (correlation_type IN ('temporal', 'spatial', 'logical', 'statistical', 'causal', 'predictive')),
    strength REAL DEFAULT 0.5 CHECK (strength >= 0.0 AND strength <= 1.0),
    confidence REAL DEFAULT 0.5 CHECK (confidence >= 0.0 AND confidence <= 1.0),
    evidence TEXT,
    description TEXT,
    weight REAL DEFAULT 1.0 CHECK (weight >= 0.0),
    direction TEXT CHECK (direction IN ('unidirectional', 'bidirectional')),
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (source_finding_id) REFERENCES intelligence_findings(id) ON DELETE CASCADE,
    FOREIGN KEY (target_finding_id) REFERENCES intelligence_findings(id) ON DELETE CASCADE
);

-- Create indexes for new tables
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_session_id ON intelligence_findings(session_id);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_intelligence_type ON intelligence_findings(intelligence_type);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_threat_level ON intelligence_findings(threat_level);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_mitre_tactic ON intelligence_findings(mitre_tactic);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_mitre_technique ON intelligence_findings(mitre_technique);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_ioc_type ON intelligence_findings(ioc_type);
CREATE INDEX IF NOT EXISTS idx_intelligence_findings_ioc_value ON intelligence_findings(ioc_value);

CREATE INDEX IF NOT EXISTS idx_threat_actors_session_id ON threat_actors(session_id);
CREATE INDEX IF NOT EXISTS idx_threat_actors_name ON threat_actors(name);
CREATE INDEX IF NOT EXISTS idx_threat_actors_threat_level ON threat_actors(threat_level);

CREATE INDEX IF NOT EXISTS idx_attack_campaigns_session_id ON attack_campaigns(session_id);
CREATE INDEX IF NOT EXISTS idx_attack_campaigns_status ON attack_campaigns(status);
CREATE INDEX IF NOT EXISTS idx_attack_campaigns_threat_level ON attack_campaigns(threat_level);

CREATE INDEX IF NOT EXISTS idx_indicators_of_compromise_session_id ON indicators_of_compromise(session_id);
CREATE INDEX IF NOT EXISTS idx_indicators_of_compromise_type ON indicators_of_compromise(type);
CREATE INDEX IF NOT EXISTS idx_indicators_of_compromise_value ON indicators_of_compromise(value);
CREATE INDEX IF NOT EXISTS idx_indicators_of_compromise_threat_level ON indicators_of_compromise(threat_level);

CREATE INDEX IF NOT EXISTS idx_patterns_session_id ON patterns(session_id);
CREATE INDEX IF NOT EXISTS idx_patterns_pattern_type ON patterns(pattern_type);
CREATE INDEX IF NOT EXISTS idx_patterns_severity ON patterns(severity);

CREATE INDEX IF NOT EXISTS idx_correlations_session_id ON correlations(session_id);
CREATE INDEX IF NOT EXISTS idx_correlations_source_finding_id ON correlations(source_finding_id);
CREATE INDEX IF NOT EXISTS idx_correlations_target_finding_id ON correlations(target_finding_id);
CREATE INDEX IF NOT EXISTS idx_correlations_correlation_type ON correlations(correlation_type);

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
