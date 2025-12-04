-- TinyBrain Initial Database Schema for Supabase/PostgreSQL
-- Migration: 001_initial_schema.sql
-- Created: 2024-12-04

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search
CREATE EXTENSION IF NOT EXISTS "vector"; -- For embeddings (if using pgvector)

-- =============================================================================
-- CORE TABLES
-- =============================================================================

-- Users table (extends Supabase auth.users)
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    username TEXT UNIQUE,
    full_name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Teams/Organizations
CREATE TABLE IF NOT EXISTS public.teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settings JSONB DEFAULT '{}'::jsonb
);

-- Team membership
CREATE TABLE IF NOT EXISTS public.team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL REFERENCES public.teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

-- Sessions (security assessment sessions)
CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES public.teams(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    description TEXT,
    task_type TEXT NOT NULL CHECK (task_type IN (
        'security_review', 'penetration_test', 'exploit_dev',
        'vulnerability_analysis', 'threat_modeling', 'incident_response',
        'intelligence_analysis', 'general'
    )),
    intelligence_type TEXT CHECK (intelligence_type IN (
        'osint', 'humint', 'sigint', 'geoint', 'masint',
        'techint', 'finint', 'cybint'
    )),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'archived')),
    classification TEXT DEFAULT 'unclassified' CHECK (classification IN (
        'unclassified', 'confidential', 'secret', 'top_secret'
    )),
    threat_level TEXT CHECK (threat_level IN ('low', 'medium', 'high', 'critical')),
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Memory entries
CREATE TABLE IF NOT EXISTS public.memories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES public.sessions(id) ON DELETE CASCADE,
    team_id UUID REFERENCES public.teams(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'text' CHECK (content_type IN (
        'text', 'code', 'binary', 'json', 'xml', 'markdown', 'other'
    )),
    category TEXT NOT NULL CHECK (category IN (
        'finding', 'vulnerability', 'exploit', 'payload', 'technique',
        'tool', 'reference', 'context', 'hypothesis', 'evidence',
        'recommendation', 'note', 'ioc', 'ttp', 'threat_actor',
        'attack_pattern', 'malware', 'correlation'
    )),
    priority INTEGER NOT NULL DEFAULT 5 CHECK (priority BETWEEN 1 AND 10),
    confidence DECIMAL(3,2) DEFAULT 0.5 CHECK (confidence BETWEEN 0 AND 1),
    source TEXT,
    tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    mitre_tactic TEXT,
    mitre_technique TEXT,
    kill_chain_phase TEXT,
    access_count INTEGER NOT NULL DEFAULT 0,
    last_accessed TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb,
    embedding vector(1536) -- For semantic search with OpenAI embeddings
);

-- Memory relationships
CREATE TABLE IF NOT EXISTS public.relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_memory_id UUID NOT NULL REFERENCES public.memories(id) ON DELETE CASCADE,
    target_memory_id UUID NOT NULL REFERENCES public.memories(id) ON DELETE CASCADE,
    relationship_type TEXT NOT NULL CHECK (relationship_type IN (
        'depends_on', 'causes', 'mitigates', 'exploits', 'references',
        'contradicts', 'supports', 'related_to', 'parent_of', 'child_of'
    )),
    strength DECIMAL(3,2) DEFAULT 0.5 CHECK (strength BETWEEN 0 AND 1),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb,
    UNIQUE(source_memory_id, target_memory_id, relationship_type)
);

-- Context snapshots
CREATE TABLE IF NOT EXISTS public.context_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES public.sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    snapshot_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Task progress tracking
CREATE TABLE IF NOT EXISTS public.task_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES public.sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    task_name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('not_started', 'in_progress', 'completed', 'blocked')),
    progress INTEGER DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
    notes TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Notifications
CREATE TABLE IF NOT EXISTS public.notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES public.sessions(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN (
        'memory_created', 'memory_updated', 'high_priority_alert',
        'duplicate_detected', 'cleanup_performed', 'system_alert'
    )),
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    priority INTEGER DEFAULT 5 CHECK (priority BETWEEN 1 AND 10),
    read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

-- =============================================================================
-- SECURITY KNOWLEDGE HUB TABLES
-- =============================================================================

-- NVD CVE data
CREATE TABLE IF NOT EXISTS public.nvd_cves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cve_id TEXT UNIQUE NOT NULL,
    description TEXT,
    published_date TIMESTAMPTZ,
    last_modified_date TIMESTAMPTZ,
    cvss_v3_score DECIMAL(3,1),
    cvss_v3_vector TEXT,
    severity TEXT CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    cwe_ids TEXT[],
    references JSONB,
    configurations JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- MITRE ATT&CK data
CREATE TABLE IF NOT EXISTS public.mitre_attack (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attack_id TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('tactic', 'technique', 'subtechnique', 'mitigation', 'group', 'software')),
    name TEXT NOT NULL,
    description TEXT,
    platforms TEXT[],
    tactics TEXT[],
    parent_id TEXT,
    detection TEXT,
    references JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- OWASP testing procedures
CREATE TABLE IF NOT EXISTS public.owasp_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    test_id TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    testing_guide TEXT,
    owasp_top_10 TEXT[],
    references JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CWE patterns
CREATE TABLE IF NOT EXISTS public.cwe_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cwe_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    extended_description TEXT,
    likelihood TEXT,
    severity TEXT,
    languages TEXT[],
    consequences JSONB,
    mitigations JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- INDEXES
-- =============================================================================

-- Users indexes
CREATE INDEX idx_users_username ON public.users(username);
CREATE INDEX idx_users_created_at ON public.users(created_at);

-- Teams indexes
CREATE INDEX idx_teams_owner_id ON public.teams(owner_id);
CREATE INDEX idx_team_members_user_id ON public.team_members(user_id);
CREATE INDEX idx_team_members_team_id ON public.team_members(team_id);

-- Sessions indexes
CREATE INDEX idx_sessions_user_id ON public.sessions(user_id);
CREATE INDEX idx_sessions_team_id ON public.sessions(team_id);
CREATE INDEX idx_sessions_status ON public.sessions(status);
CREATE INDEX idx_sessions_task_type ON public.sessions(task_type);
CREATE INDEX idx_sessions_created_at ON public.sessions(created_at);

-- Memories indexes
CREATE INDEX idx_memories_user_id ON public.memories(user_id);
CREATE INDEX idx_memories_session_id ON public.memories(session_id);
CREATE INDEX idx_memories_team_id ON public.memories(team_id);
CREATE INDEX idx_memories_category ON public.memories(category);
CREATE INDEX idx_memories_priority ON public.memories(priority DESC);
CREATE INDEX idx_memories_created_at ON public.memories(created_at DESC);
CREATE INDEX idx_memories_tags ON public.memories USING GIN(tags);
CREATE INDEX idx_memories_content_fts ON public.memories USING GIN(to_tsvector('english', content));
CREATE INDEX idx_memories_title_fts ON public.memories USING GIN(to_tsvector('english', title));
CREATE INDEX idx_memories_mitre_technique ON public.memories(mitre_technique);

-- Relationships indexes
CREATE INDEX idx_relationships_source ON public.relationships(source_memory_id);
CREATE INDEX idx_relationships_target ON public.relationships(target_memory_id);
CREATE INDEX idx_relationships_type ON public.relationships(relationship_type);

-- Context snapshots indexes
CREATE INDEX idx_context_snapshots_session_id ON public.context_snapshots(session_id);
CREATE INDEX idx_context_snapshots_user_id ON public.context_snapshots(user_id);

-- Task progress indexes
CREATE INDEX idx_task_progress_session_id ON public.task_progress(session_id);
CREATE INDEX idx_task_progress_status ON public.task_progress(status);

-- Notifications indexes
CREATE INDEX idx_notifications_user_id ON public.notifications(user_id);
CREATE INDEX idx_notifications_read ON public.notifications(read, created_at DESC);
CREATE INDEX idx_notifications_session_id ON public.notifications(session_id);

-- Security knowledge hub indexes
CREATE INDEX idx_nvd_cves_cve_id ON public.nvd_cves(cve_id);
CREATE INDEX idx_nvd_cves_severity ON public.nvd_cves(severity);
CREATE INDEX idx_nvd_cves_published ON public.nvd_cves(published_date DESC);
CREATE INDEX idx_mitre_attack_attack_id ON public.mitre_attack(attack_id);
CREATE INDEX idx_mitre_attack_type ON public.mitre_attack(type);
CREATE INDEX idx_owasp_tests_test_id ON public.owasp_tests(test_id);
CREATE INDEX idx_owasp_tests_category ON public.owasp_tests(category);
CREATE INDEX idx_cwe_patterns_cwe_id ON public.cwe_patterns(cwe_id);

-- =============================================================================
-- TRIGGERS
-- =============================================================================

-- Update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sessions_updated_at BEFORE UPDATE ON public.sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_memories_updated_at BEFORE UPDATE ON public.memories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_progress_updated_at BEFORE UPDATE ON public.task_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Increment access count on memory access
CREATE OR REPLACE FUNCTION increment_memory_access()
RETURNS TRIGGER AS $$
BEGIN
    NEW.access_count = OLD.access_count + 1;
    NEW.last_accessed = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCTIONS
-- =============================================================================

-- Full-text search function
CREATE OR REPLACE FUNCTION search_memories_fts(
    search_query TEXT,
    user_id_param UUID,
    limit_param INTEGER DEFAULT 20
)
RETURNS SETOF public.memories AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM public.memories
    WHERE user_id = user_id_param
        AND (
            to_tsvector('english', title) @@ plainto_tsquery('english', search_query)
            OR to_tsvector('english', content) @@ plainto_tsquery('english', search_query)
        )
    ORDER BY
        ts_rank(to_tsvector('english', title || ' ' || content), plainto_tsquery('english', search_query)) DESC,
        created_at DESC
    LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- Get related memories function
CREATE OR REPLACE FUNCTION get_related_memories(
    memory_id_param UUID,
    limit_param INTEGER DEFAULT 10
)
RETURNS TABLE (
    memory_id UUID,
    title TEXT,
    relationship_type TEXT,
    strength DECIMAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        m.id,
        m.title,
        r.relationship_type,
        r.strength
    FROM public.memories m
    JOIN public.relationships r ON (
        (r.source_memory_id = memory_id_param AND r.target_memory_id = m.id)
        OR (r.target_memory_id = memory_id_param AND r.source_memory_id = m.id)
    )
    ORDER BY r.strength DESC, m.priority DESC
    LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMMENTS
-- =============================================================================

COMMENT ON TABLE public.users IS 'User accounts extending Supabase auth.users';
COMMENT ON TABLE public.teams IS 'Teams/organizations for collaborative work';
COMMENT ON TABLE public.sessions IS 'Security assessment sessions';
COMMENT ON TABLE public.memories IS 'Memory entries storing security findings and intelligence';
COMMENT ON TABLE public.relationships IS 'Relationships between memory entries';
COMMENT ON TABLE public.context_snapshots IS 'Snapshots of session context at specific points';
COMMENT ON TABLE public.task_progress IS 'Progress tracking for security tasks';
COMMENT ON TABLE public.notifications IS 'User notifications and alerts';
COMMENT ON TABLE public.nvd_cves IS 'National Vulnerability Database CVE entries';
COMMENT ON TABLE public.mitre_attack IS 'MITRE ATT&CK framework data';
COMMENT ON TABLE public.owasp_tests IS 'OWASP testing procedures';
COMMENT ON TABLE public.cwe_patterns IS 'Common Weakness Enumeration patterns';
