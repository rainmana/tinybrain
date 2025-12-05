-- TinyBrain Row Level Security Policies
-- Migration: 002_row_level_security.sql
-- Created: 2024-12-04

-- =============================================================================
-- ENABLE RLS ON ALL TABLES
-- =============================================================================

ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.team_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.memories ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.relationships ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.context_snapshots ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.task_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.notifications ENABLE ROW LEVEL SECURITY;

-- Security knowledge hub tables are read-only for all authenticated users
ALTER TABLE public.nvd_cves ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.mitre_attack ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.owasp_tests ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.cwe_patterns ENABLE ROW LEVEL SECURITY;

-- =============================================================================
-- HELPER FUNCTIONS FOR RLS
-- =============================================================================

-- Check if user is member of a team
CREATE OR REPLACE FUNCTION is_team_member(team_id_param UUID)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM public.team_members
        WHERE team_id = team_id_param
        AND user_id = auth.uid()
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Check if user has specific role in team
CREATE OR REPLACE FUNCTION has_team_role(team_id_param UUID, required_role TEXT)
RETURNS BOOLEAN AS $$
DECLARE
    user_role TEXT;
BEGIN
    SELECT role INTO user_role
    FROM public.team_members
    WHERE team_id = team_id_param
    AND user_id = auth.uid();
    
    IF user_role IS NULL THEN
        RETURN FALSE;
    END IF;
    
    -- Role hierarchy: owner > admin > member > viewer
    CASE required_role
        WHEN 'viewer' THEN
            RETURN user_role IN ('viewer', 'member', 'admin', 'owner');
        WHEN 'member' THEN
            RETURN user_role IN ('member', 'admin', 'owner');
        WHEN 'admin' THEN
            RETURN user_role IN ('admin', 'owner');
        WHEN 'owner' THEN
            RETURN user_role = 'owner';
        ELSE
            RETURN FALSE;
    END CASE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- =============================================================================
-- USERS TABLE POLICIES
-- =============================================================================

-- Users can view their own profile
CREATE POLICY "Users can view own profile"
    ON public.users FOR SELECT
    USING (auth.uid() = id);

-- Users can update their own profile
CREATE POLICY "Users can update own profile"
    ON public.users FOR UPDATE
    USING (auth.uid() = id);

-- Users can view profiles of team members
CREATE POLICY "Users can view team member profiles"
    ON public.users FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM public.team_members tm1
            JOIN public.team_members tm2 ON tm1.team_id = tm2.team_id
            WHERE tm1.user_id = auth.uid()
            AND tm2.user_id = public.users.id
        )
    );

-- =============================================================================
-- TEAMS TABLE POLICIES
-- =============================================================================

-- Users can view their own teams
CREATE POLICY "Users can view own teams"
    ON public.teams FOR SELECT
    USING (
        owner_id = auth.uid()
        OR is_team_member(id)
    );

-- Users can create teams
CREATE POLICY "Users can create teams"
    ON public.teams FOR INSERT
    WITH CHECK (auth.uid() = owner_id);

-- Team owners can update their teams
CREATE POLICY "Team owners can update teams"
    ON public.teams FOR UPDATE
    USING (owner_id = auth.uid() OR has_team_role(id, 'admin'));

-- Team owners can delete their teams
CREATE POLICY "Team owners can delete teams"
    ON public.teams FOR DELETE
    USING (owner_id = auth.uid());

-- =============================================================================
-- TEAM MEMBERS TABLE POLICIES
-- =============================================================================

-- Team members can view team membership
CREATE POLICY "Team members can view membership"
    ON public.team_members FOR SELECT
    USING (is_team_member(team_id));

-- Team admins can add members
CREATE POLICY "Team admins can add members"
    ON public.team_members FOR INSERT
    WITH CHECK (has_team_role(team_id, 'admin'));

-- Team admins can update member roles (except owner)
CREATE POLICY "Team admins can update member roles"
    ON public.team_members FOR UPDATE
    USING (
        has_team_role(team_id, 'admin')
        AND role != 'owner'
    );

-- Team admins can remove members (except owner)
CREATE POLICY "Team admins can remove members"
    ON public.team_members FOR DELETE
    USING (
        has_team_role(team_id, 'admin')
        AND role != 'owner'
    );

-- Users can leave teams (except owner)
CREATE POLICY "Users can leave teams"
    ON public.team_members FOR DELETE
    USING (
        user_id = auth.uid()
        AND role != 'owner'
    );

-- =============================================================================
-- SESSIONS TABLE POLICIES
-- =============================================================================

-- Users can view their own sessions
CREATE POLICY "Users can view own sessions"
    ON public.sessions FOR SELECT
    USING (user_id = auth.uid());

-- Team members can view team sessions
CREATE POLICY "Team members can view team sessions"
    ON public.sessions FOR SELECT
    USING (
        team_id IS NOT NULL
        AND is_team_member(team_id)
    );

-- Users can create sessions
CREATE POLICY "Users can create sessions"
    ON public.sessions FOR INSERT
    WITH CHECK (
        auth.uid() = user_id
        AND (team_id IS NULL OR is_team_member(team_id))
    );

-- Users can update their own sessions
CREATE POLICY "Users can update own sessions"
    ON public.sessions FOR UPDATE
    USING (user_id = auth.uid());

-- Team members with write access can update team sessions
CREATE POLICY "Team members can update team sessions"
    ON public.sessions FOR UPDATE
    USING (
        team_id IS NOT NULL
        AND has_team_role(team_id, 'member')
    );

-- Users can delete their own sessions
CREATE POLICY "Users can delete own sessions"
    ON public.sessions FOR DELETE
    USING (user_id = auth.uid());

-- =============================================================================
-- MEMORIES TABLE POLICIES
-- =============================================================================

-- Users can view their own memories
CREATE POLICY "Users can view own memories"
    ON public.memories FOR SELECT
    USING (user_id = auth.uid());

-- Team members can view team memories
CREATE POLICY "Team members can view team memories"
    ON public.memories FOR SELECT
    USING (
        team_id IS NOT NULL
        AND is_team_member(team_id)
    );

-- Users can create memories
CREATE POLICY "Users can create memories"
    ON public.memories FOR INSERT
    WITH CHECK (
        auth.uid() = user_id
        AND (team_id IS NULL OR is_team_member(team_id))
    );

-- Users can update their own memories
CREATE POLICY "Users can update own memories"
    ON public.memories FOR UPDATE
    USING (user_id = auth.uid());

-- Team members with write access can update team memories
CREATE POLICY "Team members can update team memories"
    ON public.memories FOR UPDATE
    USING (
        team_id IS NOT NULL
        AND has_team_role(team_id, 'member')
    );

-- Users can delete their own memories
CREATE POLICY "Users can delete own memories"
    ON public.memories FOR DELETE
    USING (user_id = auth.uid());

-- Team admins can delete team memories
CREATE POLICY "Team admins can delete team memories"
    ON public.memories FOR DELETE
    USING (
        team_id IS NOT NULL
        AND has_team_role(team_id, 'admin')
    );

-- =============================================================================
-- RELATIONSHIPS TABLE POLICIES
-- =============================================================================

-- Users can view relationships for their accessible memories
CREATE POLICY "Users can view accessible relationships"
    ON public.relationships FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM public.memories
            WHERE id = source_memory_id
            AND (
                user_id = auth.uid()
                OR (team_id IS NOT NULL AND is_team_member(team_id))
            )
        )
    );

-- Users can create relationships between accessible memories
CREATE POLICY "Users can create relationships"
    ON public.relationships FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM public.memories
            WHERE id = source_memory_id
            AND (
                user_id = auth.uid()
                OR (team_id IS NOT NULL AND has_team_role(team_id, 'member'))
            )
        )
        AND EXISTS (
            SELECT 1 FROM public.memories
            WHERE id = target_memory_id
            AND (
                user_id = auth.uid()
                OR (team_id IS NOT NULL AND has_team_role(team_id, 'member'))
            )
        )
    );

-- Users can update relationships they created
CREATE POLICY "Users can update relationships"
    ON public.relationships FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM public.memories
            WHERE id = source_memory_id
            AND (
                user_id = auth.uid()
                OR (team_id IS NOT NULL AND has_team_role(team_id, 'member'))
            )
        )
    );

-- Users can delete relationships they created
CREATE POLICY "Users can delete relationships"
    ON public.relationships FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM public.memories
            WHERE id = source_memory_id
            AND (
                user_id = auth.uid()
                OR (team_id IS NOT NULL AND has_team_role(team_id, 'admin'))
            )
        )
    );

-- =============================================================================
-- CONTEXT SNAPSHOTS TABLE POLICIES
-- =============================================================================

-- Users can view their own context snapshots
CREATE POLICY "Users can view own context snapshots"
    ON public.context_snapshots FOR SELECT
    USING (user_id = auth.uid());

-- Team members can view team session snapshots
CREATE POLICY "Team members can view team snapshots"
    ON public.context_snapshots FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM public.sessions
            WHERE id = context_snapshots.session_id
            AND team_id IS NOT NULL
            AND is_team_member(team_id)
        )
    );

-- Users can create context snapshots
CREATE POLICY "Users can create context snapshots"
    ON public.context_snapshots FOR INSERT
    WITH CHECK (auth.uid() = user_id);

-- Users can delete their own context snapshots
CREATE POLICY "Users can delete own context snapshots"
    ON public.context_snapshots FOR DELETE
    USING (user_id = auth.uid());

-- =============================================================================
-- TASK PROGRESS TABLE POLICIES
-- =============================================================================

-- Users can view their own task progress
CREATE POLICY "Users can view own task progress"
    ON public.task_progress FOR SELECT
    USING (user_id = auth.uid());

-- Team members can view team session task progress
CREATE POLICY "Team members can view team task progress"
    ON public.task_progress FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM public.sessions
            WHERE id = task_progress.session_id
            AND team_id IS NOT NULL
            AND is_team_member(team_id)
        )
    );

-- Users can create task progress
CREATE POLICY "Users can create task progress"
    ON public.task_progress FOR INSERT
    WITH CHECK (auth.uid() = user_id);

-- Users can update their own task progress
CREATE POLICY "Users can update own task progress"
    ON public.task_progress FOR UPDATE
    USING (user_id = auth.uid());

-- Users can delete their own task progress
CREATE POLICY "Users can delete own task progress"
    ON public.task_progress FOR DELETE
    USING (user_id = auth.uid());

-- =============================================================================
-- NOTIFICATIONS TABLE POLICIES
-- =============================================================================

-- Users can view their own notifications
CREATE POLICY "Users can view own notifications"
    ON public.notifications FOR SELECT
    USING (user_id = auth.uid());

-- System can create notifications (service role)
-- This is handled by service role key, no policy needed for INSERT

-- Users can update their own notifications (mark as read)
CREATE POLICY "Users can update own notifications"
    ON public.notifications FOR UPDATE
    USING (user_id = auth.uid());

-- Users can delete their own notifications
CREATE POLICY "Users can delete own notifications"
    ON public.notifications FOR DELETE
    USING (user_id = auth.uid());

-- =============================================================================
-- SECURITY KNOWLEDGE HUB POLICIES (READ-ONLY)
-- =============================================================================

-- All authenticated users can read NVD CVE data
CREATE POLICY "Authenticated users can read NVD data"
    ON public.nvd_cves FOR SELECT
    TO authenticated
    USING (true);

-- All authenticated users can read MITRE ATT&CK data
CREATE POLICY "Authenticated users can read MITRE data"
    ON public.mitre_attack FOR SELECT
    TO authenticated
    USING (true);

-- All authenticated users can read OWASP data
CREATE POLICY "Authenticated users can read OWASP data"
    ON public.owasp_tests FOR SELECT
    TO authenticated
    USING (true);

-- All authenticated users can read CWE data
CREATE POLICY "Authenticated users can read CWE data"
    ON public.cwe_patterns FOR SELECT
    TO authenticated
    USING (true);

-- Only service role can write to security knowledge hub tables
-- (These policies are implicitly enforced by not creating INSERT/UPDATE/DELETE policies)

-- =============================================================================
-- GRANT PERMISSIONS
-- =============================================================================

-- Grant usage on schema
GRANT USAGE ON SCHEMA public TO authenticated;
GRANT USAGE ON SCHEMA public TO anon;

-- Grant select on all tables to authenticated users (filtered by RLS)
GRANT SELECT ON ALL TABLES IN SCHEMA public TO authenticated;

-- Grant insert/update/delete on specific tables to authenticated users (filtered by RLS)
GRANT INSERT, UPDATE, DELETE ON public.users TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.teams TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.team_members TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.sessions TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.memories TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.relationships TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.context_snapshots TO authenticated;
GRANT INSERT, UPDATE, DELETE ON public.task_progress TO authenticated;
GRANT UPDATE, DELETE ON public.notifications TO authenticated;

-- Grant select on security knowledge hub to all authenticated users
GRANT SELECT ON public.nvd_cves TO authenticated;
GRANT SELECT ON public.mitre_attack TO authenticated;
GRANT SELECT ON public.owasp_tests TO authenticated;
GRANT SELECT ON public.cwe_patterns TO authenticated;

-- =============================================================================
-- COMMENTS
-- =============================================================================

COMMENT ON POLICY "Users can view own profile" ON public.users IS 
    'Users can view their own profile information';

COMMENT ON POLICY "Team members can view team sessions" ON public.sessions IS 
    'Team members can view all sessions within their team';

COMMENT ON POLICY "Users can view own memories" ON public.memories IS 
    'Users can view all their personal memories';

COMMENT ON FUNCTION is_team_member(UUID) IS 
    'Helper function to check if current user is member of specified team';

COMMENT ON FUNCTION has_team_role(UUID, TEXT) IS 
    'Helper function to check if current user has required role in specified team';
