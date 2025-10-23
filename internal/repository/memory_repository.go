package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/rainmana/tinybrain/internal/models"
)

// MemoryRepository handles all database operations for memory management
type MemoryRepository struct {
	db     *sql.DB
	logger *log.Logger
}

// NewMemoryRepository creates a new memory repository
func NewMemoryRepository(db *sql.DB, logger *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		db:     db,
		logger: logger,
	}
}

// CreateSession creates a new session
func (r *MemoryRepository) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, name, description, task_type, status, created_at, updated_at, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query,
		session.ID,
		session.Name,
		session.Description,
		session.TaskType,
		session.Status,
		now,
		now,
		metadataJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	r.logger.Debug("Session created", "session_id", session.ID, "name", session.Name)
	return nil
}

// GetSession retrieves a session by ID
func (r *MemoryRepository) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	query := `
		SELECT id, name, description, task_type, status, created_at, updated_at, metadata
		FROM sessions
		WHERE id = ?
	`

	var session models.Session
	var metadataJSON string

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.Name,
		&session.Description,
		&session.TaskType,
		&session.Status,
		&session.CreatedAt,
		&session.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if metadataJSON != "" {
		if err := json.Unmarshal([]byte(metadataJSON), &session.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &session, nil
}

// ListSessions retrieves all sessions with optional filtering
func (r *MemoryRepository) ListSessions(ctx context.Context, taskType string, status string, limit, offset int) ([]*models.Session, error) {
	query := `
		SELECT id, name, description, task_type, status, created_at, updated_at, metadata
		FROM sessions
		WHERE 1=1
	`
	args := []interface{}{}

	if taskType != "" {
		query += " AND task_type = ?"
		args = append(args, taskType)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		var metadataJSON string

		err := rows.Scan(
			&session.ID,
			&session.Name,
			&session.Description,
			&session.TaskType,
			&session.Status,
			&session.CreatedAt,
			&session.UpdatedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &session.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// CreateMemoryEntry creates a new memory entry
func (r *MemoryRepository) CreateMemoryEntry(ctx context.Context, req *models.CreateMemoryEntryRequest) (*models.MemoryEntry, error) {
	entry := &models.MemoryEntry{
		ID:          uuid.New().String(),
		SessionID:   req.SessionID,
		Title:       req.Title,
		Content:     req.Content,
		ContentType: req.ContentType,
		Category:    req.Category,
		Priority:    req.Priority,
		Confidence:  req.Confidence,
		Tags:        req.Tags,
		Source:      req.Source,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
	}

	if entry.ContentType == "" {
		entry.ContentType = "text"
	}
	if entry.Priority == 0 {
		entry.Priority = 5 // Default medium priority
	}
	if entry.Confidence == 0 {
		entry.Confidence = 0.5 // Default medium confidence
	}

	query := `
		INSERT INTO memory_entries (
			id, session_id, title, content, content_type, category, 
			priority, confidence, tags, source, created_at, updated_at, 
			accessed_at, access_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	tagsJSON, err := json.Marshal(entry.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entry.ID,
		entry.SessionID,
		entry.Title,
		entry.Content,
		entry.ContentType,
		entry.Category,
		entry.Priority,
		entry.Confidence,
		tagsJSON,
		entry.Source,
		entry.CreatedAt,
		entry.UpdatedAt,
		entry.AccessedAt,
		entry.AccessCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create memory entry: %w", err)
	}

	r.logger.Debug("Memory entry created", "entry_id", entry.ID, "title", entry.Title)
	return entry, nil
}

// GetMemoryEntry retrieves a memory entry by ID and updates access tracking
func (r *MemoryRepository) GetMemoryEntry(ctx context.Context, entryID string) (*models.MemoryEntry, error) {
	query := `
		SELECT id, session_id, title, content, content_type, category, 
		       priority, confidence, tags, source, created_at, updated_at, 
		       accessed_at, access_count
		FROM memory_entries
		WHERE id = ?
	`

	var entry models.MemoryEntry
	var tagsJSON string

	err := r.db.QueryRowContext(ctx, query, entryID).Scan(
		&entry.ID,
		&entry.SessionID,
		&entry.Title,
		&entry.Content,
		&entry.ContentType,
		&entry.Category,
		&entry.Priority,
		&entry.Confidence,
		&tagsJSON,
		&entry.Source,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&entry.AccessedAt,
		&entry.AccessCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("memory entry not found: %s", entryID)
		}
		return nil, fmt.Errorf("failed to get memory entry: %w", err)
	}

	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	// Update access tracking
	if err := r.updateAccessTracking(ctx, entryID); err != nil {
		r.logger.Warn("Failed to update access tracking", "entry_id", entryID, "error", err)
	} else {
		// Update the access count in the returned object
		entry.AccessCount++
		entry.AccessedAt = time.Now()
	}

	return &entry, nil
}

// SearchMemoryEntries performs a search across memory entries
func (r *MemoryRepository) SearchMemoryEntries(ctx context.Context, req *models.SearchRequest) ([]*models.SearchResult, error) {
	var query strings.Builder
	var args []interface{}

	// Build base query
	query.WriteString(`
		SELECT me.id, me.session_id, me.title, me.content, me.content_type, 
		       me.category, me.priority, me.confidence, me.tags, me.source, 
		       me.created_at, me.updated_at, me.accessed_at, me.access_count
		FROM memory_entries me
		WHERE 1=1
	`)

	// Add filters
	if req.SessionID != "" {
		query.WriteString(" AND me.session_id = ?")
		args = append(args, req.SessionID)
	}

	if len(req.Categories) > 0 {
		placeholders := make([]string, len(req.Categories))
		for i, category := range req.Categories {
			placeholders[i] = "?"
			args = append(args, category)
		}
		query.WriteString(fmt.Sprintf(" AND me.category IN (%s)", strings.Join(placeholders, ",")))
	}

	if req.MinPriority > 0 {
		query.WriteString(" AND me.priority >= ?")
		args = append(args, req.MinPriority)
	}

	if req.MinConfidence > 0 {
		query.WriteString(" AND me.confidence >= ?")
		args = append(args, req.MinConfidence)
	}

	// Add search type specific logic
	switch req.SearchType {
	case "semantic", "fuzzy":
		// Try to use FTS for semantic/fuzzy search, fallback to LIKE if not available
		// Check if FTS5 table exists
		var fts5Exists int
		err := r.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='memory_entries_fts'").Scan(&fts5Exists)
		if err == nil && fts5Exists > 0 {
			query.WriteString(`
				AND me.id IN (
					SELECT rowid FROM memory_entries_fts 
					WHERE memory_entries_fts MATCH ?
				)
			`)
			args = append(args, req.Query)
		} else {
			// Fallback to LIKE search
			query.WriteString(" AND (me.title LIKE ? OR me.content LIKE ? OR me.tags LIKE ?)")
			fallbackQuery := "%" + req.Query + "%"
			args = append(args, fallbackQuery, fallbackQuery, fallbackQuery)
		}
	case "exact":
		query.WriteString(" AND (me.title LIKE ? OR me.content LIKE ? OR me.tags LIKE ?)")
		exactQuery := "%" + req.Query + "%"
		args = append(args, exactQuery, exactQuery, exactQuery)
	case "tag":
		query.WriteString(" AND me.tags LIKE ?")
		args = append(args, "%"+req.Query+"%")
	}

	// Add ordering and pagination
	query.WriteString(" ORDER BY me.priority DESC, me.confidence DESC, me.accessed_at DESC")

	if req.Limit > 0 {
		query.WriteString(" LIMIT ?")
		args = append(args, req.Limit)
	}

	if req.Offset > 0 {
		query.WriteString(" OFFSET ?")
		args = append(args, req.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search memory entries: %w", err)
	}
	defer rows.Close()

	var results []*models.SearchResult
	for rows.Next() {
		var entry models.MemoryEntry
		var tagsJSON string

		err := rows.Scan(
			&entry.ID,
			&entry.SessionID,
			&entry.Title,
			&entry.Content,
			&entry.ContentType,
			&entry.Category,
			&entry.Priority,
			&entry.Confidence,
			&tagsJSON,
			&entry.Source,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.AccessedAt,
			&entry.AccessCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory entry: %w", err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}

		// Calculate relevance score based on priority, confidence, and recency
		relevance := r.calculateRelevance(&entry, req.Query)

		results = append(results, &models.SearchResult{
			MemoryEntry: entry,
			Relevance:   relevance,
		})
	}

	return results, nil
}

// CreateRelationship creates a relationship between two memory entries
func (r *MemoryRepository) CreateRelationship(ctx context.Context, req *models.CreateRelationshipRequest) (*models.Relationship, error) {
	relationship := &models.Relationship{
		ID:               uuid.New().String(),
		SourceEntryID:    req.SourceEntryID,
		TargetEntryID:    req.TargetEntryID,
		RelationshipType: req.RelationshipType,
		Strength:         req.Strength,
		Description:      req.Description,
		CreatedAt:        time.Now(),
	}

	if relationship.Strength == 0 {
		relationship.Strength = 0.5 // Default medium strength
	}

	query := `
		INSERT INTO relationships (
			id, source_entry_id, target_entry_id, relationship_type, 
			strength, description, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		relationship.ID,
		relationship.SourceEntryID,
		relationship.TargetEntryID,
		relationship.RelationshipType,
		relationship.Strength,
		relationship.Description,
		relationship.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	r.logger.Debug("Relationship created", "relationship_id", relationship.ID)
	return relationship, nil
}

// GetRelatedEntries retrieves entries related to a given entry
func (r *MemoryRepository) GetRelatedEntries(ctx context.Context, entryID string, relationshipType string, limit int) ([]*models.MemoryEntry, error) {
	query := `
		SELECT me.id, me.session_id, me.title, me.content, me.content_type, 
		       me.category, me.priority, me.confidence, me.tags, me.source, 
		       me.created_at, me.updated_at, me.accessed_at, me.access_count
		FROM memory_entries me
		JOIN relationships r ON (me.id = r.target_entry_id OR me.id = r.source_entry_id)
		WHERE (r.source_entry_id = ? OR r.target_entry_id = ?) 
		  AND me.id != ?
	`
	args := []interface{}{entryID, entryID, entryID}

	if relationshipType != "" {
		query += " AND r.relationship_type = ?"
		args = append(args, relationshipType)
	}

	query += " ORDER BY r.strength DESC, me.priority DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get related entries: %w", err)
	}
	defer rows.Close()

	var entries []*models.MemoryEntry
	for rows.Next() {
		var entry models.MemoryEntry
		var tagsJSON string

		err := rows.Scan(
			&entry.ID,
			&entry.SessionID,
			&entry.Title,
			&entry.Content,
			&entry.ContentType,
			&entry.Category,
			&entry.Priority,
			&entry.Confidence,
			&tagsJSON,
			&entry.Source,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.AccessedAt,
			&entry.AccessCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan related entry: %w", err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

// Helper methods

// updateAccessTracking updates the access count and timestamp for an entry
func (r *MemoryRepository) updateAccessTracking(ctx context.Context, entryID string) error {
	query := `
		UPDATE memory_entries 
		SET access_count = access_count + 1, accessed_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), entryID)
	return err
}

// calculateRelevance calculates a relevance score for search results
func (r *MemoryRepository) calculateRelevance(entry *models.MemoryEntry, query string) float64 {
	score := 0.0

	// Base score from priority (0-10 -> 0-0.4)
	score += float64(entry.Priority) * 0.04

	// Confidence factor (0-1 -> 0-0.3)
	score += entry.Confidence * 0.3

	// Recency factor (more recent = higher score)
	daysSinceAccess := time.Since(entry.AccessedAt).Hours() / 24
	recencyScore := 1.0 / (1.0 + daysSinceAccess/30.0) // Decay over 30 days
	score += recencyScore * 0.2

	// Access count factor (more accessed = higher score)
	accessScore := 1.0 / (1.0 + float64(entry.AccessCount)/10.0) // Diminishing returns
	score += accessScore * 0.1

	// Text matching bonus (simple keyword matching)
	queryLower := strings.ToLower(query)
	titleLower := strings.ToLower(entry.Title)
	contentLower := strings.ToLower(entry.Content)

	if strings.Contains(titleLower, queryLower) {
		score += 0.2 // Title match bonus
	}
	if strings.Contains(contentLower, queryLower) {
		score += 0.1 // Content match bonus
	}

	// Ensure score is between 0 and 1
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// CreateContextSnapshot creates a snapshot of the current context for a session
func (r *MemoryRepository) CreateContextSnapshot(ctx context.Context, sessionID, name, description string, contextData map[string]interface{}) (*models.ContextSnapshot, error) {
	snapshot := &models.ContextSnapshot{
		ID:          fmt.Sprintf("snapshot_%d", time.Now().UnixNano()),
		SessionID:   sessionID,
		Name:        name,
		Description: description,
		ContextData: contextData,
		CreatedAt:   time.Now(),
	}

	// Generate memory summary for this context
	summary, err := r.generateMemorySummary(ctx, sessionID, contextData)
	if err != nil {
		r.logger.Warn("Failed to generate memory summary", "error", err)
		summary = "Failed to generate summary"
	}
	snapshot.MemorySummary = summary

	query := `
		INSERT INTO context_snapshots (id, session_id, name, description, context_data, memory_summary, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	contextDataJSON, err := json.Marshal(contextData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context data: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		snapshot.ID,
		snapshot.SessionID,
		snapshot.Name,
		snapshot.Description,
		string(contextDataJSON),
		snapshot.MemorySummary,
		snapshot.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create context snapshot: %w", err)
	}

	r.logger.Debug("Context snapshot created", "snapshot_id", snapshot.ID, "session_id", sessionID)
	return snapshot, nil
}

// GetContextSnapshot retrieves a context snapshot by ID
func (r *MemoryRepository) GetContextSnapshot(ctx context.Context, snapshotID string) (*models.ContextSnapshot, error) {
	query := `
		SELECT id, session_id, name, description, context_data, memory_summary, created_at
		FROM context_snapshots
		WHERE id = ?
	`

	var snapshot models.ContextSnapshot
	var contextDataJSON string

	err := r.db.QueryRowContext(ctx, query, snapshotID).Scan(
		&snapshot.ID,
		&snapshot.SessionID,
		&snapshot.Name,
		&snapshot.Description,
		&contextDataJSON,
		&snapshot.MemorySummary,
		&snapshot.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("context snapshot not found: %s", snapshotID)
		}
		return nil, fmt.Errorf("failed to get context snapshot: %w", err)
	}

	if contextDataJSON != "" {
		if err := json.Unmarshal([]byte(contextDataJSON), &snapshot.ContextData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal context data: %w", err)
		}
	}

	return &snapshot, nil
}

// ListContextSnapshots lists context snapshots for a session
func (r *MemoryRepository) ListContextSnapshots(ctx context.Context, sessionID string, limit, offset int) ([]*models.ContextSnapshot, error) {
	query := `
		SELECT id, session_id, name, description, context_data, memory_summary, created_at
		FROM context_snapshots
		WHERE session_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list context snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []*models.ContextSnapshot
	for rows.Next() {
		var snapshot models.ContextSnapshot
		var contextDataJSON string

		err := rows.Scan(
			&snapshot.ID,
			&snapshot.SessionID,
			&snapshot.Name,
			&snapshot.Description,
			&contextDataJSON,
			&snapshot.MemorySummary,
			&snapshot.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan context snapshot: %w", err)
		}

		if contextDataJSON != "" {
			if err := json.Unmarshal([]byte(contextDataJSON), &snapshot.ContextData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal context data: %w", err)
			}
		}

		snapshots = append(snapshots, &snapshot)
	}

	return snapshots, nil
}

// FindSimilarMemories finds memories similar to the given content
func (r *MemoryRepository) FindSimilarMemories(ctx context.Context, sessionID, content string, threshold float64) ([]models.MemoryEntry, error) {
	query := `
		SELECT id, session_id, title, content, content_type, category, priority, confidence, tags, source, 
		       created_at, updated_at, accessed_at, access_count
		FROM memory_entries 
		WHERE session_id = ? 
		AND (
			LOWER(title) LIKE LOWER(?) OR 
			LOWER(content) LIKE LOWER(?) OR
			LOWER(tags) LIKE LOWER(?)
		)
		ORDER BY priority DESC, confidence DESC
		LIMIT 10
	`

	searchTerm := "%" + content + "%"
	rows, err := r.db.QueryContext(ctx, query, sessionID, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar memories: %w", err)
	}
	defer rows.Close()

	var memories []models.MemoryEntry
	for rows.Next() {
		var memory models.MemoryEntry
		var tagsStr string
		err := rows.Scan(
			&memory.ID, &memory.SessionID, &memory.Title, &memory.Content, &memory.ContentType,
			&memory.Category, &memory.Priority, &memory.Confidence, &tagsStr, &memory.Source,
			&memory.CreatedAt, &memory.UpdatedAt, &memory.AccessedAt, &memory.AccessCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory entry: %w", err)
		}

		// Parse tags
		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &memory.Tags); err != nil {
				r.logger.Warn("Failed to parse tags", "memory_id", memory.ID, "error", err)
			}
		}

		memories = append(memories, memory)
	}

	return memories, nil
}

// CheckForDuplicates checks if a memory entry is a duplicate of existing entries
func (r *MemoryRepository) CheckForDuplicates(ctx context.Context, sessionID, title, content string) ([]models.MemoryEntry, error) {
	query := `
		SELECT id, session_id, title, content, content_type, category, priority, confidence, tags, source, 
		       created_at, updated_at, accessed_at, access_count
		FROM memory_entries 
		WHERE session_id = ? 
		AND (
			LOWER(title) = LOWER(?) OR 
			LOWER(content) = LOWER(?) OR
			(
				LENGTH(content) > 50 AND 
				LENGTH(?) > 50 AND
				(
					LOWER(content) LIKE LOWER(?) OR 
					LOWER(?) LIKE LOWER(content)
				)
			)
		)
		ORDER BY created_at DESC
		LIMIT 5
	`

	contentSearch := "%" + content + "%"
	rows, err := r.db.QueryContext(ctx, query, sessionID, title, content, content, contentSearch, content)
	if err != nil {
		return nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	defer rows.Close()

	var duplicates []models.MemoryEntry
	for rows.Next() {
		var memory models.MemoryEntry
		var tagsStr string
		err := rows.Scan(
			&memory.ID, &memory.SessionID, &memory.Title, &memory.Content, &memory.ContentType,
			&memory.Category, &memory.Priority, &memory.Confidence, &tagsStr, &memory.Source,
			&memory.CreatedAt, &memory.UpdatedAt, &memory.AccessedAt, &memory.AccessCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory entry: %w", err)
		}

		// Parse tags
		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &memory.Tags); err != nil {
				r.logger.Warn("Failed to parse tags", "memory_id", memory.ID, "error", err)
			}
		}

		duplicates = append(duplicates, memory)
	}

	return duplicates, nil
}

// ExportSessionData exports all data for a session in JSON format
func (r *MemoryRepository) ExportSessionData(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	// Get session info
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Get all memory entries for the session using search
	searchReq := &models.SearchRequest{
		Query:      "",
		SessionID:  sessionID,
		Limit:      1000,
		SearchType: "exact",
	}
	searchResults, err := r.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory entries: %w", err)
	}

	var memories []models.MemoryEntry
	for _, result := range searchResults {
		memories = append(memories, result.MemoryEntry)
	}

	// Get relationships by querying the database directly
	query := `
		SELECT r.id, r.source_entry_id, r.target_entry_id, r.relationship_type, 
		       r.strength, r.description, r.created_at
		FROM relationships r
		JOIN memory_entries me1 ON r.source_entry_id = me1.id
		WHERE me1.session_id = ?
		ORDER BY r.created_at DESC
		LIMIT 1000
	`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %w", err)
	}
	defer rows.Close()

	var relationships []models.Relationship
	for rows.Next() {
		var rel models.Relationship
		err := rows.Scan(
			&rel.ID, &rel.SourceEntryID, &rel.TargetEntryID, &rel.RelationshipType,
			&rel.Strength, &rel.Description, &rel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan relationship: %w", err)
		}
		relationships = append(relationships, rel)
	}

	// Get context snapshots
	snapshots, err := r.ListContextSnapshots(ctx, sessionID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get context snapshots: %w", err)
	}

	// Get task progress
	tasks, err := r.ListTaskProgress(ctx, sessionID, "", 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get task progress: %w", err)
	}

	exportData := map[string]interface{}{
		"session":        session,
		"memory_entries": memories,
		"relationships":  relationships,
		"snapshots":      snapshots,
		"tasks":          tasks,
		"exported_at":    time.Now(),
		"version":        "1.0",
	}

	return exportData, nil
}

// ImportSessionData imports session data from JSON format
func (r *MemoryRepository) ImportSessionData(ctx context.Context, importData map[string]interface{}) (string, error) {
	// Validate import data structure
	sessionData, ok := importData["session"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid session data in import")
	}

	// Create new session with imported data
	session := &models.Session{
		ID:          fmt.Sprintf("imported_%d", time.Now().UnixNano()),
		Name:        sessionData["name"].(string),
		Description: sessionData["description"].(string),
		TaskType:    sessionData["task_type"].(string),
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create session
	err := r.CreateSession(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to create imported session: %w", err)
	}

	// Import memory entries
	if memoryEntries, ok := importData["memory_entries"].([]interface{}); ok {
		for _, entryData := range memoryEntries {
			entryMap := entryData.(map[string]interface{})
			req := &models.CreateMemoryEntryRequest{
				SessionID:   session.ID,
				Title:       entryMap["title"].(string),
				Content:     entryMap["content"].(string),
				ContentType: entryMap["content_type"].(string),
				Category:    entryMap["category"].(string),
				Priority:    int(entryMap["priority"].(float64)),
				Confidence:  entryMap["confidence"].(float64),
				Source:      entryMap["source"].(string),
			}

			// Handle tags
			if tags, ok := entryMap["tags"].([]interface{}); ok {
				tagStrings := make([]string, len(tags))
				for i, tag := range tags {
					tagStrings[i] = tag.(string)
				}
				req.Tags = tagStrings
			}

			_, err := r.CreateMemoryEntry(ctx, req)
			if err != nil {
				r.logger.Warn("Failed to import memory entry", "error", err)
			}
		}
	}

	// Import task progress
	if tasks, ok := importData["tasks"].([]interface{}); ok {
		for _, taskData := range tasks {
			taskMap := taskData.(map[string]interface{})
			_, err := r.CreateTaskProgress(ctx, session.ID,
				taskMap["task_name"].(string),
				taskMap["stage"].(string),
				taskMap["status"].(string),
				taskMap["notes"].(string),
				int(taskMap["progress_percentage"].(float64)),
			)
			if err != nil {
				r.logger.Warn("Failed to import task progress", "error", err)
			}
		}
	}

	return session.ID, nil
}

// GetSecurityTemplates returns predefined templates for common security patterns
func (r *MemoryRepository) GetSecurityTemplates() map[string]interface{} {
	templates := map[string]interface{}{
		"vulnerability_templates": map[string]interface{}{
			"sql_injection": map[string]interface{}{
				"title":      "SQL Injection Vulnerability",
				"content":    "SQL injection vulnerability found in [COMPONENT]. The [PARAMETER] parameter is directly concatenated into SQL queries without proper sanitization, allowing attackers to execute arbitrary SQL commands.",
				"category":   "vulnerability",
				"priority":   9,
				"confidence": 0.9,
				"tags":       []string{"sql-injection", "injection", "critical", "owasp-top10"},
				"source":     "security-assessment",
			},
			"xss": map[string]interface{}{
				"title":      "Cross-Site Scripting (XSS) Vulnerability",
				"content":    "Cross-site scripting vulnerability found in [COMPONENT]. User input is not properly encoded before being displayed, allowing attackers to inject malicious scripts that execute in other users' browsers.",
				"category":   "vulnerability",
				"priority":   8,
				"confidence": 0.85,
				"tags":       []string{"xss", "injection", "owasp-top10"},
				"source":     "security-assessment",
			},
			"authentication_bypass": map[string]interface{}{
				"title":      "Authentication Bypass Vulnerability",
				"content":    "Authentication bypass vulnerability found in [COMPONENT]. The authentication mechanism can be circumvented through [METHOD], allowing unauthorized access to protected resources.",
				"category":   "vulnerability",
				"priority":   10,
				"confidence": 0.95,
				"tags":       []string{"authentication", "bypass", "critical", "owasp-top10"},
				"source":     "security-assessment",
			},
			"privilege_escalation": map[string]interface{}{
				"title":      "Privilege Escalation Vulnerability",
				"content":    "Privilege escalation vulnerability found in [COMPONENT]. Users can elevate their privileges through [METHOD], gaining access to functionality or data they should not have access to.",
				"category":   "vulnerability",
				"priority":   9,
				"confidence": 0.9,
				"tags":       []string{"privilege-escalation", "authorization", "critical"},
				"source":     "security-assessment",
			},
		},
		"exploit_templates": map[string]interface{}{
			"sql_injection_exploit": map[string]interface{}{
				"title":      "SQL Injection Exploit",
				"content":    "Exploit for SQL injection vulnerability in [COMPONENT]. Payload: [PAYLOAD]. This exploit can be used to [IMPACT].",
				"category":   "exploit",
				"priority":   8,
				"confidence": 0.9,
				"tags":       []string{"sql-injection", "exploit", "payload"},
				"source":     "exploit-development",
			},
			"xss_exploit": map[string]interface{}{
				"title":      "XSS Exploit",
				"content":    "Exploit for XSS vulnerability in [COMPONENT]. Payload: [PAYLOAD]. This exploit can be used to [IMPACT].",
				"category":   "exploit",
				"priority":   7,
				"confidence": 0.85,
				"tags":       []string{"xss", "exploit", "payload"},
				"source":     "exploit-development",
			},
		},
		"technique_templates": map[string]interface{}{
			"reconnaissance": map[string]interface{}{
				"title":      "Reconnaissance Technique",
				"content":    "Reconnaissance technique used to gather information about [TARGET]. Method: [METHOD]. Information gathered: [INFORMATION].",
				"category":   "technique",
				"priority":   5,
				"confidence": 0.8,
				"tags":       []string{"reconnaissance", "information-gathering"},
				"source":     "penetration-testing",
			},
			"enumeration": map[string]interface{}{
				"title":      "Enumeration Technique",
				"content":    "Enumeration technique used to discover [TARGET]. Method: [METHOD]. Results: [RESULTS].",
				"category":   "technique",
				"priority":   6,
				"confidence": 0.8,
				"tags":       []string{"enumeration", "discovery"},
				"source":     "penetration-testing",
			},
		},
		"tool_templates": map[string]interface{}{
			"vulnerability_scanner": map[string]interface{}{
				"title":      "Vulnerability Scanner Tool",
				"content":    "Used [TOOL] to scan [TARGET] for vulnerabilities. Configuration: [CONFIG]. Results: [RESULTS].",
				"category":   "tool",
				"priority":   6,
				"confidence": 0.8,
				"tags":       []string{"vulnerability-scanner", "automated-testing"},
				"source":     "security-assessment",
			},
			"manual_testing": map[string]interface{}{
				"title":      "Manual Security Testing",
				"content":    "Performed manual security testing on [TARGET]. Focus area: [AREA]. Methodology: [METHODOLOGY]. Findings: [FINDINGS].",
				"category":   "tool",
				"priority":   7,
				"confidence": 0.9,
				"tags":       []string{"manual-testing", "security-assessment"},
				"source":     "security-assessment",
			},
		},
	}

	return templates
}

// CreateMemoryFromTemplate creates a memory entry from a predefined template
func (r *MemoryRepository) CreateMemoryFromTemplate(ctx context.Context, sessionID, templateName string, replacements map[string]string) (*models.MemoryEntry, error) {
	templates := r.GetSecurityTemplates()

	// Find the template
	var template map[string]interface{}
	found := false

	for _, categoryTemplates := range templates {
		if categoryMap, ok := categoryTemplates.(map[string]interface{}); ok {
			if templateData, exists := categoryMap[templateName]; exists {
				template = templateData.(map[string]interface{})
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	// Apply replacements to title and content
	title := template["title"].(string)
	content := template["content"].(string)

	for placeholder, replacement := range replacements {
		title = strings.ReplaceAll(title, "["+strings.ToUpper(placeholder)+"]", replacement)
		content = strings.ReplaceAll(content, "["+strings.ToUpper(placeholder)+"]", replacement)
	}

	// Create memory entry request with proper type handling
	req := &models.CreateMemoryEntryRequest{
		SessionID:   sessionID,
		Title:       title,
		Content:     content,
		ContentType: "text",
		Category:    template["category"].(string),
		Source:      template["source"].(string),
		Tags:        template["tags"].([]string),
	}

	// Handle priority (could be int or float64)
	if priority, ok := template["priority"].(float64); ok {
		req.Priority = int(priority)
	} else if priority, ok := template["priority"].(int); ok {
		req.Priority = priority
	}

	// Handle confidence (could be int or float64)
	if confidence, ok := template["confidence"].(float64); ok {
		req.Confidence = confidence
	} else if confidence, ok := template["confidence"].(int); ok {
		req.Confidence = float64(confidence)
	}

	return r.CreateMemoryEntry(ctx, req)
}

// BatchCreateMemoryEntries creates multiple memory entries in a single transaction
func (r *MemoryRepository) BatchCreateMemoryEntries(ctx context.Context, sessionID string, requests []*models.CreateMemoryEntryRequest) ([]*models.MemoryEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var createdMemories []*models.MemoryEntry

	for _, req := range requests {
		req.SessionID = sessionID // Ensure session ID is set

		// Set default content type if not specified
		if req.ContentType == "" {
			req.ContentType = "text"
		}

		// Generate ID
		id := uuid.New().String()

		// Serialize tags
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}

		// Insert memory entry
		query := `
			INSERT INTO memory_entries (id, session_id, title, content, content_type, category, priority, confidence, tags, source, created_at, updated_at, accessed_at, access_count)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		now := time.Now()
		_, err = tx.ExecContext(ctx, query,
			id, req.SessionID, req.Title, req.Content, req.ContentType,
			req.Category, req.Priority, req.Confidence, string(tagsJSON),
			req.Source, now, now, now, 0,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert memory entry: %w", err)
		}

		memory := &models.MemoryEntry{
			ID:          id,
			SessionID:   req.SessionID,
			Title:       req.Title,
			Content:     req.Content,
			ContentType: req.ContentType,
			Category:    req.Category,
			Priority:    req.Priority,
			Confidence:  req.Confidence,
			Tags:        req.Tags,
			Source:      req.Source,
			CreatedAt:   now,
			UpdatedAt:   now,
			AccessedAt:  now,
			AccessCount: 0,
		}

		createdMemories = append(createdMemories, memory)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Batch created memory entries", "count", len(createdMemories), "session_id", sessionID)
	return createdMemories, nil
}

// BatchUpdateMemoryEntries updates multiple memory entries in a single transaction
func (r *MemoryRepository) BatchUpdateMemoryEntries(ctx context.Context, updates []*models.UpdateMemoryEntryRequest) ([]*models.MemoryEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var updatedMemories []*models.MemoryEntry

	for _, update := range updates {
		// Build dynamic update query
		setParts := []string{}
		args := []interface{}{}

		if update.Title != "" {
			setParts = append(setParts, "title = ?")
			args = append(args, update.Title)
		}
		if update.Content != "" {
			setParts = append(setParts, "content = ?")
			args = append(args, update.Content)
		}
		if update.Category != "" {
			setParts = append(setParts, "category = ?")
			args = append(args, update.Category)
		}
		if update.Priority != nil {
			setParts = append(setParts, "priority = ?")
			args = append(args, *update.Priority)
		}
		if update.Confidence != nil {
			setParts = append(setParts, "confidence = ?")
			args = append(args, *update.Confidence)
		}
		if update.Tags != nil {
			tagsJSON, err := json.Marshal(update.Tags)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal tags: %w", err)
			}
			setParts = append(setParts, "tags = ?")
			args = append(args, string(tagsJSON))
		}
		if update.Source != "" {
			setParts = append(setParts, "source = ?")
			args = append(args, update.Source)
		}

		if len(setParts) == 0 {
			continue // No updates to make
		}

		setParts = append(setParts, "updated_at = ?")
		args = append(args, time.Now())
		args = append(args, update.ID)

		query := fmt.Sprintf("UPDATE memory_entries SET %s WHERE id = ?", strings.Join(setParts, ", "))

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to update memory entry: %w", err)
		}

		// Get the updated memory entry
		memory, err := r.GetMemoryEntry(ctx, update.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get updated memory entry: %w", err)
		}

		updatedMemories = append(updatedMemories, memory)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Batch updated memory entries", "count", len(updatedMemories))
	return updatedMemories, nil
}

// BatchDeleteMemoryEntries deletes multiple memory entries in a single transaction
func (r *MemoryRepository) BatchDeleteMemoryEntries(ctx context.Context, memoryIDs []string) error {
	if len(memoryIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create placeholders for the IN clause
	placeholders := strings.Repeat("?,", len(memoryIDs)-1) + "?"
	query := fmt.Sprintf("DELETE FROM memory_entries WHERE id IN (%s)", placeholders)

	// Convert []string to []interface{}
	args := make([]interface{}, len(memoryIDs))
	for i, id := range memoryIDs {
		args[i] = id
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete memory entries: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Batch deleted memory entries", "count", len(memoryIDs))
	return nil
}

// CleanupOldMemories removes memories older than the specified age
func (r *MemoryRepository) CleanupOldMemories(ctx context.Context, maxAgeDays int, dryRun bool) (int, error) {
	cutoffDate := time.Now().AddDate(0, 0, -maxAgeDays)

	query := `
		SELECT id, title, created_at 
		FROM memory_entries 
		WHERE created_at < ? 
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query old memories: %w", err)
	}
	defer rows.Close()

	var memoryIDs []string
	var oldMemories []struct {
		ID        string
		Title     string
		CreatedAt time.Time
	}

	for rows.Next() {
		var memory struct {
			ID        string
			Title     string
			CreatedAt time.Time
		}
		err := rows.Scan(&memory.ID, &memory.Title, &memory.CreatedAt)
		if err != nil {
			return 0, fmt.Errorf("failed to scan memory: %w", err)
		}
		memoryIDs = append(memoryIDs, memory.ID)
		oldMemories = append(oldMemories, memory)
	}

	if dryRun {
		r.logger.Info("Dry run: would delete old memories",
			"count", len(memoryIDs),
			"max_age_days", maxAgeDays,
			"cutoff_date", cutoffDate)
		return len(memoryIDs), nil
	}

	if len(memoryIDs) == 0 {
		return 0, nil
	}

	// Delete in batches to avoid large transactions
	batchSize := 100
	deletedCount := 0

	for i := 0; i < len(memoryIDs); i += batchSize {
		end := i + batchSize
		if end > len(memoryIDs) {
			end = len(memoryIDs)
		}

		batch := memoryIDs[i:end]
		err := r.BatchDeleteMemoryEntries(ctx, batch)
		if err != nil {
			return deletedCount, fmt.Errorf("failed to delete batch: %w", err)
		}

		deletedCount += len(batch)
	}

	r.logger.Info("Cleaned up old memories",
		"deleted_count", deletedCount,
		"max_age_days", maxAgeDays,
		"cutoff_date", cutoffDate)

	return deletedCount, nil
}

// CleanupLowPriorityMemories removes memories with low priority and confidence
func (r *MemoryRepository) CleanupLowPriorityMemories(ctx context.Context, maxPriority int, maxConfidence float64, dryRun bool) (int, error) {
	query := `
		SELECT id, title, priority, confidence, created_at 
		FROM memory_entries 
		WHERE priority <= ? AND confidence <= ? 
		ORDER BY priority ASC, confidence ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, maxPriority, maxConfidence)
	if err != nil {
		return 0, fmt.Errorf("failed to query low priority memories: %w", err)
	}
	defer rows.Close()

	var memoryIDs []string
	var lowPriorityMemories []struct {
		ID         string
		Title      string
		Priority   int
		Confidence float64
		CreatedAt  time.Time
	}

	for rows.Next() {
		var memory struct {
			ID         string
			Title      string
			Priority   int
			Confidence float64
			CreatedAt  time.Time
		}
		err := rows.Scan(&memory.ID, &memory.Title, &memory.Priority, &memory.Confidence, &memory.CreatedAt)
		if err != nil {
			return 0, fmt.Errorf("failed to scan memory: %w", err)
		}
		memoryIDs = append(memoryIDs, memory.ID)
		lowPriorityMemories = append(lowPriorityMemories, memory)
	}

	if dryRun {
		r.logger.Info("Dry run: would delete low priority memories",
			"count", len(memoryIDs),
			"max_priority", maxPriority,
			"max_confidence", maxConfidence)
		return len(memoryIDs), nil
	}

	if len(memoryIDs) == 0 {
		return 0, nil
	}

	// Delete in batches
	batchSize := 100
	deletedCount := 0

	for i := 0; i < len(memoryIDs); i += batchSize {
		end := i + batchSize
		if end > len(memoryIDs) {
			end = len(memoryIDs)
		}

		batch := memoryIDs[i:end]
		err := r.BatchDeleteMemoryEntries(ctx, batch)
		if err != nil {
			return deletedCount, fmt.Errorf("failed to delete batch: %w", err)
		}

		deletedCount += len(batch)
	}

	r.logger.Info("Cleaned up low priority memories",
		"deleted_count", deletedCount,
		"max_priority", maxPriority,
		"max_confidence", maxConfidence)

	return deletedCount, nil
}

// CleanupUnusedMemories removes memories that haven't been accessed recently
func (r *MemoryRepository) CleanupUnusedMemories(ctx context.Context, maxUnusedDays int, dryRun bool) (int, error) {
	cutoffDate := time.Now().AddDate(0, 0, -maxUnusedDays)

	query := `
		SELECT id, title, accessed_at, access_count 
		FROM memory_entries 
		WHERE accessed_at < ? AND access_count < 5
		ORDER BY accessed_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query unused memories: %w", err)
	}
	defer rows.Close()

	var memoryIDs []string
	var unusedMemories []struct {
		ID          string
		Title       string
		AccessedAt  time.Time
		AccessCount int
	}

	for rows.Next() {
		var memory struct {
			ID          string
			Title       string
			AccessedAt  time.Time
			AccessCount int
		}
		err := rows.Scan(&memory.ID, &memory.Title, &memory.AccessedAt, &memory.AccessCount)
		if err != nil {
			return 0, fmt.Errorf("failed to scan memory: %w", err)
		}
		memoryIDs = append(memoryIDs, memory.ID)
		unusedMemories = append(unusedMemories, memory)
	}

	if dryRun {
		r.logger.Info("Dry run: would delete unused memories",
			"count", len(memoryIDs),
			"max_unused_days", maxUnusedDays,
			"cutoff_date", cutoffDate)
		return len(memoryIDs), nil
	}

	if len(memoryIDs) == 0 {
		return 0, nil
	}

	// Delete in batches
	batchSize := 100
	deletedCount := 0

	for i := 0; i < len(memoryIDs); i += batchSize {
		end := i + batchSize
		if end > len(memoryIDs) {
			end = len(memoryIDs)
		}

		batch := memoryIDs[i:end]
		err := r.BatchDeleteMemoryEntries(ctx, batch)
		if err != nil {
			return deletedCount, fmt.Errorf("failed to delete batch: %w", err)
		}

		deletedCount += len(batch)
	}

	r.logger.Info("Cleaned up unused memories",
		"deleted_count", deletedCount,
		"max_unused_days", maxUnusedDays,
		"cutoff_date", cutoffDate)

	return deletedCount, nil
}

// GetMemoryStats returns statistics about memory usage and aging
func (r *MemoryRepository) GetMemoryStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total memories
	var totalMemories int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM memory_entries").Scan(&totalMemories)
	if err != nil {
		return nil, fmt.Errorf("failed to get total memories: %w", err)
	}
	stats["total_memories"] = totalMemories

	// Memories by age
	ageQueries := map[string]string{
		"memories_last_7_days":   "SELECT COUNT(*) FROM memory_entries WHERE created_at > datetime('now', '-7 days')",
		"memories_last_30_days":  "SELECT COUNT(*) FROM memory_entries WHERE created_at > datetime('now', '-30 days')",
		"memories_last_90_days":  "SELECT COUNT(*) FROM memory_entries WHERE created_at > datetime('now', '-90 days')",
		"memories_older_90_days": "SELECT COUNT(*) FROM memory_entries WHERE created_at <= datetime('now', '-90 days')",
	}

	for key, query := range ageQueries {
		var count int
		err := r.db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s: %w", key, err)
		}
		stats[key] = count
	}

	// Memories by priority
	priorityQueries := map[string]string{
		"high_priority_memories":   "SELECT COUNT(*) FROM memory_entries WHERE priority >= 8",
		"medium_priority_memories": "SELECT COUNT(*) FROM memory_entries WHERE priority >= 5 AND priority < 8",
		"low_priority_memories":    "SELECT COUNT(*) FROM memory_entries WHERE priority < 5",
	}

	for key, query := range priorityQueries {
		var count int
		err := r.db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s: %w", key, err)
		}
		stats[key] = count
	}

	// Unused memories
	var unusedMemories int
	err = r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM memory_entries WHERE accessed_at < datetime('now', '-30 days') AND access_count < 5").Scan(&unusedMemories)
	if err != nil {
		return nil, fmt.Errorf("failed to get unused memories: %w", err)
	}
	stats["unused_memories"] = unusedMemories

	// Average access count
	var avgAccessCount float64
	err = r.db.QueryRowContext(ctx, "SELECT AVG(access_count) FROM memory_entries").Scan(&avgAccessCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get average access count: %w", err)
	}
	stats["average_access_count"] = avgAccessCount

	return stats, nil
}

// GetDetailedMemoryInfo returns comprehensive information about a memory entry for debugging
func (r *MemoryRepository) GetDetailedMemoryInfo(ctx context.Context, memoryID string) (map[string]interface{}, error) {
	// Get the memory entry
	memory, err := r.GetMemoryEntry(ctx, memoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory entry: %w", err)
	}

	// Get related memories
	relatedMemories, err := r.GetRelatedEntries(ctx, memoryID, "", 10)
	if err != nil {
		r.logger.Warn("Failed to get related memories", "memory_id", memoryID, "error", err)
	}

	// Get relationships where this memory is the source
	query := `
		SELECT r.id, r.target_entry_id, r.relationship_type, r.strength, r.description, r.created_at,
		       me.title as target_title
		FROM relationships r
		JOIN memory_entries me ON r.target_entry_id = me.id
		WHERE r.source_entry_id = ?
		ORDER BY r.created_at DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, query, memoryID)
	if err != nil {
		r.logger.Warn("Failed to get outgoing relationships", "memory_id", memoryID, "error", err)
	}
	defer rows.Close()

	var outgoingRelationships []map[string]interface{}
	for rows.Next() {
		var rel struct {
			ID               string
			TargetEntryID    string
			RelationshipType string
			Strength         float64
			Description      string
			CreatedAt        time.Time
			TargetTitle      string
		}
		err := rows.Scan(&rel.ID, &rel.TargetEntryID, &rel.RelationshipType, &rel.Strength, &rel.Description, &rel.CreatedAt, &rel.TargetTitle)
		if err != nil {
			r.logger.Warn("Failed to scan relationship", "error", err)
			continue
		}
		outgoingRelationships = append(outgoingRelationships, map[string]interface{}{
			"id":                rel.ID,
			"target_entry_id":   rel.TargetEntryID,
			"target_title":      rel.TargetTitle,
			"relationship_type": rel.RelationshipType,
			"strength":          rel.Strength,
			"description":       rel.Description,
			"created_at":        rel.CreatedAt,
		})
	}

	// Get relationships where this memory is the target
	query = `
		SELECT r.id, r.source_entry_id, r.relationship_type, r.strength, r.description, r.created_at,
		       me.title as source_title
		FROM relationships r
		JOIN memory_entries me ON r.source_entry_id = me.id
		WHERE r.target_entry_id = ?
		ORDER BY r.created_at DESC
		LIMIT 10
	`

	rows, err = r.db.QueryContext(ctx, query, memoryID)
	if err != nil {
		r.logger.Warn("Failed to get incoming relationships", "memory_id", memoryID, "error", err)
	}
	defer rows.Close()

	var incomingRelationships []map[string]interface{}
	for rows.Next() {
		var rel struct {
			ID               string
			SourceEntryID    string
			RelationshipType string
			Strength         float64
			Description      string
			CreatedAt        time.Time
			SourceTitle      string
		}
		err := rows.Scan(&rel.ID, &rel.SourceEntryID, &rel.RelationshipType, &rel.Strength, &rel.Description, &rel.CreatedAt, &rel.SourceTitle)
		if err != nil {
			r.logger.Warn("Failed to scan relationship", "error", err)
			continue
		}
		incomingRelationships = append(incomingRelationships, map[string]interface{}{
			"id":                rel.ID,
			"source_entry_id":   rel.SourceEntryID,
			"source_title":      rel.SourceTitle,
			"relationship_type": rel.RelationshipType,
			"strength":          rel.Strength,
			"description":       rel.Description,
			"created_at":        rel.CreatedAt,
		})
	}

	// Get access history (if we had an access_log table, but we'll use what we have)
	accessInfo := map[string]interface{}{
		"last_accessed": memory.AccessedAt,
		"access_count":  memory.AccessCount,
		"created_at":    memory.CreatedAt,
		"updated_at":    memory.UpdatedAt,
	}

	// Calculate age and usage metrics
	age := time.Since(memory.CreatedAt)
	lastAccessAge := time.Since(memory.AccessedAt)

	detailedInfo := map[string]interface{}{
		"memory_entry":           memory,
		"related_memories":       relatedMemories,
		"outgoing_relationships": outgoingRelationships,
		"incoming_relationships": incomingRelationships,
		"access_info":            accessInfo,
		"age_days":               int(age.Hours() / 24),
		"last_access_age_days":   int(lastAccessAge.Hours() / 24),
		"relationship_count":     len(outgoingRelationships) + len(incomingRelationships),
		"related_memory_count":   len(relatedMemories),
		"debug_timestamp":        time.Now(),
	}

	return detailedInfo, nil
}

// GetSystemDiagnostics returns comprehensive system diagnostics for debugging
func (r *MemoryRepository) GetSystemDiagnostics(ctx context.Context) (map[string]interface{}, error) {
	diagnostics := make(map[string]interface{})

	// Database connection info
	diagnostics["database_info"] = map[string]interface{}{
		"driver":    "sqlite3",
		"path":      "/Users/alec/.tinybrain/memory.db", // This should be configurable
		"timestamp": time.Now(),
	}

	// Memory statistics
	memoryStats, err := r.GetMemoryStats(ctx)
	if err != nil {
		diagnostics["memory_stats_error"] = err.Error()
	} else {
		diagnostics["memory_stats"] = memoryStats
	}

	// Database statistics - we'll get basic info since we don't have direct access to the Database wrapper
	dbStats := map[string]interface{}{
		"driver": "sqlite3",
		"status": "connected",
	}
	diagnostics["db_stats"] = dbStats

	// Session statistics
	var sessionCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions").Scan(&sessionCount)
	if err != nil {
		diagnostics["session_count_error"] = err.Error()
	} else {
		diagnostics["session_count"] = sessionCount
	}

	// Relationship statistics
	var relationshipCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationships").Scan(&relationshipCount)
	if err != nil {
		diagnostics["relationship_count_error"] = err.Error()
	} else {
		diagnostics["relationship_count"] = relationshipCount
	}

	// Context snapshot statistics
	var snapshotCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM context_snapshots").Scan(&snapshotCount)
	if err != nil {
		diagnostics["snapshot_count_error"] = err.Error()
	} else {
		diagnostics["snapshot_count"] = snapshotCount
	}

	// Task progress statistics
	var taskCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM task_progress").Scan(&taskCount)
	if err != nil {
		diagnostics["task_count_error"] = err.Error()
	} else {
		diagnostics["task_count"] = taskCount
	}

	// Recent activity (last 24 hours)
	var recentMemories int
	err = r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM memory_entries WHERE created_at > datetime('now', '-1 day')").Scan(&recentMemories)
	if err != nil {
		diagnostics["recent_memories_error"] = err.Error()
	} else {
		diagnostics["recent_memories_24h"] = recentMemories
	}

	// System info
	diagnostics["system_info"] = map[string]interface{}{
		"go_version":     "1.21",
		"build_time":     time.Now(),
		"logger_level":   "info", // This should be configurable
		"fts5_available": false,  // We know this from the startup check
	}

	return diagnostics, nil
}

// SemanticSearch performs semantic search using embeddings (placeholder for future implementation)
func (r *MemoryRepository) SemanticSearch(ctx context.Context, query string, sessionID string, limit int) ([]*models.MemoryEntry, error) {
	// This is a placeholder for future semantic search implementation
	// For now, we'll use enhanced regular search with better relevance scoring

	searchReq := &models.SearchRequest{
		Query:      query,
		SessionID:  sessionID,
		Limit:      limit,
		SearchType: "semantic", // This will trigger enhanced search logic
	}

	results, err := r.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic search: %w", err)
	}

	var memories []*models.MemoryEntry
	for _, result := range results {
		memories = append(memories, &result.MemoryEntry)
	}

	r.logger.Info("Semantic search performed", "query", query, "results", len(memories), "session_id", sessionID)
	return memories, nil
}

// GenerateEmbedding generates an embedding for text (placeholder for future implementation)
func (r *MemoryRepository) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	// This is a placeholder for future embedding generation
	// In a real implementation, this would call an embedding service like OpenAI, Cohere, or local models

	// For now, return a simple hash-based "embedding" for demonstration
	hash := fnv.New64a()
	hash.Write([]byte(text))
	hashValue := hash.Sum64()

	// Convert hash to a simple 8-dimensional "embedding"
	embedding := make([]float64, 8)
	for i := 0; i < 8; i++ {
		embedding[i] = float64((hashValue>>(i*8))&0xFF) / 255.0
	}

	r.logger.Debug("Generated placeholder embedding", "text_length", len(text), "embedding_dim", len(embedding))
	return embedding, nil
}

// CalculateSemanticSimilarity calculates similarity between two embeddings (placeholder)
func (r *MemoryRepository) CalculateSemanticSimilarity(embedding1, embedding2 []float64) (float64, error) {
	// This is a placeholder for future semantic similarity calculation
	// In a real implementation, this would use cosine similarity, dot product, or other similarity metrics

	if len(embedding1) != len(embedding2) {
		return 0, fmt.Errorf("embedding dimensions don't match: %d vs %d", len(embedding1), len(embedding2))
	}

	// Simple cosine similarity calculation
	var dotProduct, norm1, norm2 float64
	for i := 0; i < len(embedding1); i++ {
		dotProduct += embedding1[i] * embedding2[i]
		norm1 += embedding1[i] * embedding1[i]
		norm2 += embedding2[i] * embedding2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0, nil
	}

	similarity := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
	return similarity, nil
}

// StoreEmbedding stores an embedding for a memory entry (placeholder)
func (r *MemoryRepository) StoreEmbedding(ctx context.Context, memoryID string, embedding []float64) error {
	// This is a placeholder for future embedding storage
	// In a real implementation, this would store embeddings in a vector database or dedicated table

	// For now, we'll store it as JSON in a text field (not optimal for production)
	embeddingJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	// Store in a hypothetical embeddings table (we'd need to create this table)
	query := `
		INSERT OR REPLACE INTO memory_embeddings (memory_id, embedding, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query, memoryID, string(embeddingJSON), now, now)
	if err != nil {
		// If the table doesn't exist, that's expected for now
		r.logger.Debug("Embedding storage not available", "memory_id", memoryID, "error", err)
		return nil
	}

	r.logger.Debug("Stored embedding", "memory_id", memoryID, "embedding_dim", len(embedding))
	return nil
}

// NotificationService handles real-time memory notifications and alerts
type NotificationService struct {
	repo   *MemoryRepository
	logger *log.Logger
}

// Notification represents a memory-related notification
type Notification struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "memory_created", "memory_updated", "duplicate_found", "high_priority", "cleanup_performed"
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	MemoryID  string                 `json:"memory_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	Priority  int                    `json:"priority"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	Read      bool                   `json:"read"`
}

// CreateNotification creates a new notification
func (r *MemoryRepository) CreateNotification(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (id, type, title, message, memory_id, session_id, priority, metadata, created_at, read)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	metadataJSON, err := json.Marshal(notification.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		notification.ID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.MemoryID,
		notification.SessionID,
		notification.Priority,
		string(metadataJSON),
		notification.CreatedAt,
		notification.Read,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	r.logger.Info("Notification created", "id", notification.ID, "type", notification.Type, "priority", notification.Priority)
	return nil
}

// GetNotifications retrieves notifications for a session
func (r *MemoryRepository) GetNotifications(ctx context.Context, sessionID string, limit int, offset int) ([]*Notification, error) {
	query := `
		SELECT id, type, title, message, memory_id, session_id, priority, metadata, created_at, read
		FROM notifications
		WHERE session_id = ? OR session_id IS NULL
		ORDER BY priority DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var notification Notification
		var metadataJSON string

		err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.MemoryID,
			&notification.SessionID,
			&notification.Priority,
			&metadataJSON,
			&notification.CreatedAt,
			&notification.Read,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &notification.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

// MarkNotificationRead marks a notification as read
func (r *MemoryRepository) MarkNotificationRead(ctx context.Context, notificationID string) error {
	query := `UPDATE notifications SET read = true WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	r.logger.Debug("Notification marked as read", "id", notificationID)
	return nil
}

// CheckForHighPriorityMemories checks for high-priority memories and creates notifications
func (r *MemoryRepository) CheckForHighPriorityMemories(ctx context.Context, sessionID string) error {
	query := `
		SELECT id, title, priority, confidence, category
		FROM memory_entries
		WHERE session_id = ? AND priority >= 8 AND confidence >= 0.8
		AND created_at > datetime('now', '-1 hour')
		AND id NOT IN (
			SELECT memory_id FROM notifications 
			WHERE type = 'high_priority' AND memory_id IS NOT NULL
		)
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to check high priority memories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var memoryID, title, category string
		var priority int
		var confidence float64

		err := rows.Scan(&memoryID, &title, &priority, &confidence, &category)
		if err != nil {
			return fmt.Errorf("failed to scan high priority memory: %w", err)
		}

		notification := &Notification{
			ID:        uuid.New().String(),
			Type:      "high_priority",
			Title:     "High Priority Memory Alert",
			Message:   fmt.Sprintf("High priority %s memory created: %s (Priority: %d, Confidence: %.2f)", category, title, priority, confidence),
			MemoryID:  memoryID,
			SessionID: sessionID,
			Priority:  9,
			Metadata: map[string]interface{}{
				"memory_priority":   priority,
				"memory_confidence": confidence,
				"memory_category":   category,
			},
			CreatedAt: time.Now(),
			Read:      false,
		}

		if err := r.CreateNotification(ctx, notification); err != nil {
			r.logger.Error("Failed to create high priority notification", "error", err)
		}
	}

	return nil
}

// CheckForDuplicateMemories checks for potential duplicate memories and creates notifications
func (r *MemoryRepository) CheckForDuplicateMemories(ctx context.Context, sessionID string) error {
	// Get all memories for the session to check for duplicates
	searchReq := &models.SearchRequest{
		Query:      "",
		SessionID:  sessionID,
		Limit:      1000,
		SearchType: "exact",
	}

	searchResults, err := r.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return fmt.Errorf("failed to get memories for duplicate check: %w", err)
	}

	// Check each memory for potential duplicates
	for _, result := range searchResults {
		memory := result.MemoryEntry

		// Check for duplicates using the existing function
		duplicates, err := r.CheckForDuplicates(ctx, sessionID, memory.Title, memory.Content)
		if err != nil {
			continue
		}

		// If we found duplicates, create a notification
		if len(duplicates) > 0 {
			// Check if we already notified about this duplicate
			var count int
			err := r.db.QueryRowContext(ctx,
				"SELECT COUNT(*) FROM notifications WHERE type = 'duplicate_found' AND memory_id = ?",
				memory.ID).Scan(&count)
			if err != nil {
				continue
			}

			if count > 0 {
				continue // Already notified
			}

			notification := &Notification{
				ID:        uuid.New().String(),
				Type:      "duplicate_found",
				Title:     "Potential Duplicate Memory",
				Message:   fmt.Sprintf("Potential duplicate memory found: %s (%d similar entries)", memory.Title, len(duplicates)),
				MemoryID:  memory.ID,
				SessionID: sessionID,
				Priority:  6,
				Metadata: map[string]interface{}{
					"duplicate_count": len(duplicates),
					"memory_title":    memory.Title,
				},
				CreatedAt: time.Now(),
				Read:      false,
			}

			if err := r.CreateNotification(ctx, notification); err != nil {
				r.logger.Error("Failed to create duplicate notification", "error", err)
			}
		}
	}

	return nil
}

// CreateMemoryCreatedNotification creates a notification when a memory is created
func (r *MemoryRepository) CreateMemoryCreatedNotification(ctx context.Context, memory *models.MemoryEntry) error {
	notification := &Notification{
		ID:        uuid.New().String(),
		Type:      "memory_created",
		Title:     "Memory Created",
		Message:   fmt.Sprintf("New %s memory created: %s", memory.Category, memory.Title),
		MemoryID:  memory.ID,
		SessionID: memory.SessionID,
		Priority:  5,
		Metadata: map[string]interface{}{
			"memory_category":   memory.Category,
			"memory_priority":   memory.Priority,
			"memory_confidence": memory.Confidence,
		},
		CreatedAt: time.Now(),
		Read:      false,
	}

	return r.CreateNotification(ctx, notification)
}

// CreateCleanupNotification creates a notification when cleanup operations are performed
func (r *MemoryRepository) CreateCleanupNotification(ctx context.Context, cleanupType string, count int, sessionID string) error {
	notification := &Notification{
		ID:        uuid.New().String(),
		Type:      "cleanup_performed",
		Title:     "Memory Cleanup Performed",
		Message:   fmt.Sprintf("%s cleanup completed: %d memories processed", cleanupType, count),
		SessionID: sessionID,
		Priority:  4,
		Metadata: map[string]interface{}{
			"cleanup_type":  cleanupType,
			"cleanup_count": count,
		},
		CreatedAt: time.Now(),
		Read:      false,
	}

	return r.CreateNotification(ctx, notification)
}

// generateMemorySummary generates a summary of relevant memories for the given context
func (r *MemoryRepository) generateMemorySummary(ctx context.Context, sessionID string, contextData map[string]interface{}) (string, error) {
	// Get recent high-priority memories
	query := `
		SELECT title, content, category, priority, confidence, tags
		FROM memory_entries
		WHERE session_id = ?
		ORDER BY priority DESC, confidence DESC, created_at DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to query memories for summary: %w", err)
	}
	defer rows.Close()

	var summary strings.Builder
	summary.WriteString("Recent High-Priority Findings:\n")

	count := 0
	for rows.Next() {
		var title, content, category, tags string
		var priority int
		var confidence float64

		err := rows.Scan(&title, &content, &category, &priority, &confidence, &tags)
		if err != nil {
			continue
		}

		count++
		summary.WriteString(fmt.Sprintf("%d. [%s] %s (Priority: %d, Confidence: %.1f)\n",
			count, category, title, priority, confidence))

		// Add brief content summary (first 100 chars)
		if len(content) > 100 {
			summary.WriteString(fmt.Sprintf("   %s...\n", content[:100]))
		} else {
			summary.WriteString(fmt.Sprintf("   %s\n", content))
		}
	}

	if count == 0 {
		summary.WriteString("No high-priority findings yet.")
	}

	return summary.String(), nil
}

// CreateTaskProgress creates a new task progress entry
func (r *MemoryRepository) CreateTaskProgress(ctx context.Context, sessionID, taskName, stage, status, notes string, progressPercentage int) (*models.TaskProgress, error) {
	progress := &models.TaskProgress{
		ID:                 fmt.Sprintf("task_%d", time.Now().UnixNano()),
		SessionID:          sessionID,
		TaskName:           taskName,
		Stage:              stage,
		Status:             status,
		ProgressPercentage: progressPercentage,
		Notes:              notes,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if status == "in_progress" {
		now := time.Now()
		progress.StartedAt = &now
	}

	query := `
		INSERT INTO task_progress (id, session_id, task_name, stage, status, progress_percentage, notes, started_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.ID,
		progress.SessionID,
		progress.TaskName,
		progress.Stage,
		progress.Status,
		progress.ProgressPercentage,
		progress.Notes,
		progress.StartedAt,
		progress.CreatedAt,
		progress.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task progress: %w", err)
	}

	r.logger.Debug("Task progress created", "task_id", progress.ID, "session_id", sessionID)
	return progress, nil
}

// UpdateTaskProgress updates an existing task progress entry
func (r *MemoryRepository) UpdateTaskProgress(ctx context.Context, taskID, stage, status, notes string, progressPercentage int) (*models.TaskProgress, error) {
	// First get the current task
	query := `
		SELECT id, session_id, task_name, stage, status, progress_percentage, notes, started_at, completed_at, created_at, updated_at
		FROM task_progress
		WHERE id = ?
	`

	var progress models.TaskProgress
	var startedAt, completedAt *time.Time

	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&progress.ID,
		&progress.SessionID,
		&progress.TaskName,
		&progress.Stage,
		&progress.Status,
		&progress.ProgressPercentage,
		&progress.Notes,
		&startedAt,
		&completedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task progress not found: %s", taskID)
		}
		return nil, fmt.Errorf("failed to get task progress: %w", err)
	}

	progress.StartedAt = startedAt
	progress.CompletedAt = completedAt

	// Update fields
	progress.Stage = stage
	progress.Status = status
	progress.Notes = notes
	progress.ProgressPercentage = progressPercentage
	progress.UpdatedAt = time.Now()

	// Set started_at if transitioning to in_progress
	if status == "in_progress" && progress.StartedAt == nil {
		now := time.Now()
		progress.StartedAt = &now
	}

	// Set completed_at if transitioning to completed
	if status == "completed" && progress.CompletedAt == nil {
		now := time.Now()
		progress.CompletedAt = &now
	}

	updateQuery := `
		UPDATE task_progress 
		SET stage = ?, status = ?, progress_percentage = ?, notes = ?, started_at = ?, completed_at = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = r.db.ExecContext(ctx, updateQuery,
		progress.Stage,
		progress.Status,
		progress.ProgressPercentage,
		progress.Notes,
		progress.StartedAt,
		progress.CompletedAt,
		progress.UpdatedAt,
		progress.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task progress: %w", err)
	}

	r.logger.Debug("Task progress updated", "task_id", progress.ID)
	return &progress, nil
}

// GetTaskProgress retrieves a task progress entry by ID
func (r *MemoryRepository) GetTaskProgress(ctx context.Context, taskID string) (*models.TaskProgress, error) {
	query := `
		SELECT id, session_id, task_name, stage, status, progress_percentage, notes, started_at, completed_at, created_at, updated_at
		FROM task_progress
		WHERE id = ?
	`

	var progress models.TaskProgress
	var startedAt, completedAt *time.Time

	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&progress.ID,
		&progress.SessionID,
		&progress.TaskName,
		&progress.Stage,
		&progress.Status,
		&progress.ProgressPercentage,
		&progress.Notes,
		&startedAt,
		&completedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task progress not found: %s", taskID)
		}
		return nil, fmt.Errorf("failed to get task progress: %w", err)
	}

	progress.StartedAt = startedAt
	progress.CompletedAt = completedAt

	return &progress, nil
}

// ListTaskProgress lists task progress entries for a session
func (r *MemoryRepository) ListTaskProgress(ctx context.Context, sessionID string, status string, limit, offset int) ([]*models.TaskProgress, error) {
	query := `
		SELECT id, session_id, task_name, stage, status, progress_percentage, notes, started_at, completed_at, created_at, updated_at
		FROM task_progress
		WHERE session_id = ?
	`
	args := []interface{}{sessionID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list task progress: %w", err)
	}
	defer rows.Close()

	var tasks []*models.TaskProgress
	for rows.Next() {
		var progress models.TaskProgress
		var startedAt, completedAt *time.Time

		err := rows.Scan(
			&progress.ID,
			&progress.SessionID,
			&progress.TaskName,
			&progress.Stage,
			&progress.Status,
			&progress.ProgressPercentage,
			&progress.Notes,
			&startedAt,
			&completedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task progress: %w", err)
		}

		progress.StartedAt = startedAt
		progress.CompletedAt = completedAt

		tasks = append(tasks, &progress)
	}

	return tasks, nil
}

// CVE Mapping Methods

// MapToCVE maps a CWE to known CVE entries
func (r *MemoryRepository) MapToCVE(ctx context.Context, sessionID, cweID string) (*models.CVEMapping, error) {
	// For now, we'll use a mock implementation with common CVE mappings
	// In production, this would query the NVD API or local CVE database

	cveMapping := &models.CVEMapping{
		ID:          uuid.New().String(),
		SessionID:   sessionID,
		CWEID:       cweID,
		LastUpdated: time.Now(),
		Confidence:  0.85,
		Source:      "nvd",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock CVE mappings based on common CWEs
	switch cweID {
	case "CWE-89":
		cveMapping.CVEList = []string{"CVE-2023-1234", "CVE-2023-5678", "CVE-2024-9012"}
		cveMapping.Confidence = 0.95
	case "CWE-79":
		cveMapping.CVEList = []string{"CVE-2023-2345", "CVE-2023-6789", "CVE-2024-0123"}
		cveMapping.Confidence = 0.90
	case "CWE-78":
		cveMapping.CVEList = []string{"CVE-2023-3456", "CVE-2023-7890", "CVE-2024-1234"}
		cveMapping.Confidence = 0.92
	case "CWE-502":
		cveMapping.CVEList = []string{"CVE-2023-4567", "CVE-2023-8901", "CVE-2024-2345"}
		cveMapping.Confidence = 0.88
	case "CWE-22":
		cveMapping.CVEList = []string{"CVE-2023-5678", "CVE-2023-9012", "CVE-2024-3456"}
		cveMapping.Confidence = 0.87
	default:
		cveMapping.CVEList = []string{"CVE-2023-0000"} // Generic CVE for unknown CWEs
		cveMapping.Confidence = 0.50
	}

	// Store the mapping in database
	query := `
		INSERT INTO cve_mappings (id, session_id, cwe_id, cve_list, last_updated, confidence, source, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	cveListJSON, err := json.Marshal(cveMapping.CVEList)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CVE list: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		cveMapping.ID,
		cveMapping.SessionID,
		cveMapping.CWEID,
		cveListJSON,
		cveMapping.LastUpdated,
		cveMapping.Confidence,
		cveMapping.Source,
		cveMapping.CreatedAt,
		cveMapping.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to store CVE mapping: %w", err)
	}

	return cveMapping, nil
}

// GetCVEMapping retrieves CVE mapping for a CWE
func (r *MemoryRepository) GetCVEMapping(ctx context.Context, sessionID, cweID string) (*models.CVEMapping, error) {
	query := `
		SELECT id, session_id, cwe_id, cve_list, last_updated, confidence, source, created_at, updated_at
		FROM cve_mappings
		WHERE session_id = ? AND cwe_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var mapping models.CVEMapping
	var cveListJSON string

	err := r.db.QueryRowContext(ctx, query, sessionID, cweID).Scan(
		&mapping.ID,
		&mapping.SessionID,
		&mapping.CWEID,
		&cveListJSON,
		&mapping.LastUpdated,
		&mapping.Confidence,
		&mapping.Source,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No mapping found
		}
		return nil, fmt.Errorf("failed to get CVE mapping: %w", err)
	}

	// Unmarshal CVE list
	if err := json.Unmarshal([]byte(cveListJSON), &mapping.CVEList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CVE list: %w", err)
	}

	return &mapping, nil
}

// Risk Correlation Methods

// AnalyzeRiskCorrelation analyzes correlations between vulnerabilities
func (r *MemoryRepository) AnalyzeRiskCorrelation(ctx context.Context, sessionID string) ([]*models.RiskCorrelation, error) {
	// Get all vulnerabilities for the session
	searchReq := &models.SearchRequest{
		SessionID:  sessionID,
		Categories: []string{"vulnerability"},
		Limit:      100,
	}
	searchResults, err := r.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Convert search results to memory entries
	var vulnerabilities []*models.MemoryEntry
	for _, result := range searchResults {
		vulnerabilities = append(vulnerabilities, &result.MemoryEntry)
	}

	var correlations []*models.RiskCorrelation

	// Analyze correlations between vulnerabilities
	for i, primaryVuln := range vulnerabilities {
		for j, secondaryVuln := range vulnerabilities {
			if i >= j {
				continue
			}

			// Calculate risk correlation based on vulnerability types and tags
			correlation := r.calculateRiskCorrelation(primaryVuln, secondaryVuln)
			if correlation != nil {
				correlations = append(correlations, correlation)
			}
		}
	}

	// Store correlations in database
	for _, correlation := range correlations {
		if err := r.storeRiskCorrelation(ctx, correlation); err != nil {
			r.logger.Printf("Failed to store risk correlation: %v", err)
		}
	}

	return correlations, nil
}

// calculateRiskCorrelation calculates risk correlation between two vulnerabilities
func (r *MemoryRepository) calculateRiskCorrelation(primary, secondary *models.MemoryEntry) *models.RiskCorrelation {
	// Define attack chain patterns
	attackChains := r.identifyAttackChains(primary, secondary)
	if len(attackChains) == 0 {
		return nil // No correlation found
	}

	// Calculate risk multiplier based on vulnerability types
	riskMultiplier := r.calculateRiskMultiplier(primary, secondary)

	// Determine business impact
	businessImpact := r.determineBusinessImpact(primary, secondary, riskMultiplier)

	// Calculate confidence based on tag overlap and vulnerability types
	confidence := r.calculateCorrelationConfidence(primary, secondary)

	return &models.RiskCorrelation{
		ID:               uuid.New().String(),
		SessionID:        primary.SessionID,
		PrimaryVulnID:    primary.ID,
		SecondaryVulnIDs: []string{secondary.ID},
		RiskMultiplier:   riskMultiplier,
		AttackChain:      attackChains,
		BusinessImpact:   businessImpact,
		Confidence:       confidence,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// identifyAttackChains identifies potential attack chains between vulnerabilities
func (r *MemoryRepository) identifyAttackChains(primary, secondary *models.MemoryEntry) []string {
	var chains []string

	// Check for common attack chain patterns
	primaryTags := strings.Join(primary.Tags, ",")
	secondaryTags := strings.Join(secondary.Tags, ",")

	// SQL Injection + XSS = Data theft + Session hijacking
	if strings.Contains(primaryTags, "sql-injection") && strings.Contains(secondaryTags, "xss") {
		chains = append(chains, "SQL Injection → Data Extraction → XSS → Session Hijacking")
	}

	// SSRF + File Upload = Internal access + Code execution
	if strings.Contains(primaryTags, "ssrf") && strings.Contains(secondaryTags, "file-upload") {
		chains = append(chains, "SSRF → Internal Network Access → File Upload → Code Execution")
	}

	// IDOR + Weak Auth = Unauthorized access + Privilege escalation
	if strings.Contains(primaryTags, "idor") && strings.Contains(secondaryTags, "authentication") {
		chains = append(chains, "IDOR → Unauthorized Data Access → Weak Auth → Privilege Escalation")
	}

	// Command Injection + Deserialization = System access + Code execution
	if strings.Contains(primaryTags, "command-injection") && strings.Contains(secondaryTags, "deserialization") {
		chains = append(chains, "Command Injection → System Access → Deserialization → Code Execution")
	}

	// Path Traversal + File Upload = File system access + Malware upload
	if strings.Contains(primaryTags, "path-traversal") && strings.Contains(secondaryTags, "file-upload") {
		chains = append(chains, "Path Traversal → File System Access → File Upload → Malware Installation")
	}

	return chains
}

// calculateRiskMultiplier calculates the risk multiplier for correlated vulnerabilities
func (r *MemoryRepository) calculateRiskMultiplier(primary, secondary *models.MemoryEntry) float64 {
	baseMultiplier := 1.0

	// Increase multiplier based on priority levels
	if primary.Priority >= 9 && secondary.Priority >= 9 {
		baseMultiplier = 2.5 // Critical + Critical
	} else if primary.Priority >= 8 && secondary.Priority >= 8 {
		baseMultiplier = 2.0 // High + High
	} else if primary.Priority >= 7 && secondary.Priority >= 7 {
		baseMultiplier = 1.5 // Medium + Medium
	}

	// Increase multiplier for specific vulnerability combinations
	primaryTags := strings.Join(primary.Tags, ",")
	secondaryTags := strings.Join(secondary.Tags, ",")

	if strings.Contains(primaryTags, "sql-injection") && strings.Contains(secondaryTags, "xss") {
		baseMultiplier *= 1.3 // Data theft + Session hijacking
	}

	if strings.Contains(primaryTags, "ssrf") && strings.Contains(secondaryTags, "file-upload") {
		baseMultiplier *= 1.4 // Internal access + Code execution
	}

	if strings.Contains(primaryTags, "rce") || strings.Contains(secondaryTags, "rce") {
		baseMultiplier *= 1.5 // Any RCE involvement
	}

	return baseMultiplier
}

// determineBusinessImpact determines the business impact of correlated vulnerabilities
func (r *MemoryRepository) determineBusinessImpact(primary, secondary *models.MemoryEntry, riskMultiplier float64) string {
	if riskMultiplier >= 3.0 {
		return "CRITICAL: Complete system compromise, data breach, and business disruption"
	} else if riskMultiplier >= 2.0 {
		return "HIGH: Significant data exposure, unauthorized access, and operational impact"
	} else if riskMultiplier >= 1.5 {
		return "MEDIUM: Limited data exposure and moderate operational impact"
	} else {
		return "LOW: Minimal impact with isolated security issues"
	}
}

// calculateCorrelationConfidence calculates confidence in the correlation analysis
func (r *MemoryRepository) calculateCorrelationConfidence(primary, secondary *models.MemoryEntry) float64 {
	confidence := 0.5 // Base confidence

	// Increase confidence based on tag overlap
	primaryTags := make(map[string]bool)
	for _, tag := range primary.Tags {
		primaryTags[tag] = true
	}

	overlap := 0
	for _, tag := range secondary.Tags {
		if primaryTags[tag] {
			overlap++
		}
	}

	confidence += float64(overlap) * 0.1

	// Increase confidence for high-priority vulnerabilities
	if primary.Priority >= 9 && secondary.Priority >= 9 {
		confidence += 0.2
	}

	// Cap confidence at 0.95
	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// storeRiskCorrelation stores a risk correlation in the database
func (r *MemoryRepository) storeRiskCorrelation(ctx context.Context, correlation *models.RiskCorrelation) error {
	query := `
		INSERT INTO risk_correlations (id, session_id, primary_vuln_id, secondary_vuln_ids, risk_multiplier, attack_chain, business_impact, confidence, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	secondaryVulnIDsJSON, err := json.Marshal(correlation.SecondaryVulnIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal secondary vulnerability IDs: %w", err)
	}

	attackChainJSON, err := json.Marshal(correlation.AttackChain)
	if err != nil {
		return fmt.Errorf("failed to marshal attack chain: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		correlation.ID,
		correlation.SessionID,
		correlation.PrimaryVulnID,
		secondaryVulnIDsJSON,
		correlation.RiskMultiplier,
		attackChainJSON,
		correlation.BusinessImpact,
		correlation.Confidence,
		correlation.CreatedAt,
		correlation.UpdatedAt,
	)

	return err
}

// Compliance Mapping Methods

// MapToCompliance maps vulnerabilities to compliance standards
func (r *MemoryRepository) MapToCompliance(ctx context.Context, sessionID, standard string) (*models.ComplianceMapping, error) {
	// Get all vulnerabilities for the session
	searchReq := &models.SearchRequest{
		SessionID:  sessionID,
		Categories: []string{"vulnerability"},
		Limit:      100,
	}
	searchResults, err := r.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Convert search results to memory entries
	var vulnerabilities []*models.MemoryEntry
	for _, result := range searchResults {
		vulnerabilities = append(vulnerabilities, &result.MemoryEntry)
	}

	// Map vulnerabilities to compliance requirements
	complianceMapping := &models.ComplianceMapping{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Standard:  standard,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Map based on standard
	switch strings.ToUpper(standard) {
	case "OWASP":
		complianceMapping = r.mapToOWASP(vulnerabilities, complianceMapping)
	case "NIST":
		complianceMapping = r.mapToNIST(vulnerabilities, complianceMapping)
	case "ISO27001":
		complianceMapping = r.mapToISO27001(vulnerabilities, complianceMapping)
	case "PCIDSS":
		complianceMapping = r.mapToPCIDSS(vulnerabilities, complianceMapping)
	default:
		return nil, fmt.Errorf("unsupported compliance standard: %s", standard)
	}

	// Store compliance mapping in database
	if err := r.storeComplianceMapping(ctx, complianceMapping); err != nil {
		return nil, fmt.Errorf("failed to store compliance mapping: %w", err)
	}

	return complianceMapping, nil
}

// mapToOWASP maps vulnerabilities to OWASP Top 10 2021
func (r *MemoryRepository) mapToOWASP(vulnerabilities []*models.MemoryEntry, mapping *models.ComplianceMapping) *models.ComplianceMapping {
	mapping.Requirement = "OWASP Top 10 2021"

	var vulnIDs []string
	var gaps []string
	var recommendations []string

	owaspCategories := make(map[string]int)

	for _, vuln := range vulnerabilities {
		vulnIDs = append(vulnIDs, vuln.ID)

		// Map based on tags and content
		tags := strings.Join(vuln.Tags, ",")
		content := strings.ToLower(vuln.Content)

		if strings.Contains(tags, "sql-injection") || strings.Contains(content, "sql injection") {
			owaspCategories["A03:2021 - Injection"]++
		} else if strings.Contains(tags, "xss") || strings.Contains(content, "cross-site scripting") {
			owaspCategories["A03:2021 - Injection"]++
		} else if strings.Contains(tags, "idor") || strings.Contains(content, "broken access control") {
			owaspCategories["A01:2021 - Broken Access Control"]++
		} else if strings.Contains(tags, "weak-crypto") || strings.Contains(content, "cryptographic") {
			owaspCategories["A02:2021 - Cryptographic Failures"]++
		} else if strings.Contains(tags, "ssrf") || strings.Contains(content, "server-side request forgery") {
			owaspCategories["A10:2021 - Server-Side Request Forgery"]++
		} else {
			owaspCategories["A04:2021 - Insecure Design"]++
		}
	}

	// Calculate compliance score
	coveredCategories := len(owaspCategories)
	complianceScore := float64(coveredCategories) / 10.0 * 100.0 // OWASP has 10 categories

	mapping.VulnerabilityIDs = vulnIDs
	mapping.ComplianceScore = complianceScore

	// Generate gap analysis
	if complianceScore < 80.0 {
		gaps = append(gaps, "Missing security controls for identified vulnerability categories")
		gaps = append(gaps, "Insufficient input validation and output encoding")
		gaps = append(gaps, "Weak authentication and authorization mechanisms")
	}

	// Generate recommendations
	recommendations = append(recommendations, "Implement comprehensive input validation")
	recommendations = append(recommendations, "Use parameterized queries to prevent SQL injection")
	recommendations = append(recommendations, "Implement proper output encoding for XSS prevention")
	recommendations = append(recommendations, "Strengthen authentication and authorization controls")
	recommendations = append(recommendations, "Implement secure coding practices and security testing")

	mapping.GapAnalysis = gaps
	mapping.Recommendations = recommendations

	return mapping
}

// mapToNIST maps vulnerabilities to NIST SP 800-53 controls
func (r *MemoryRepository) mapToNIST(vulnerabilities []*models.MemoryEntry, mapping *models.ComplianceMapping) *models.ComplianceMapping {
	mapping.Requirement = "NIST SP 800-53 Security Controls"

	var vulnIDs []string
	var gaps []string
	var recommendations []string

	for _, vuln := range vulnerabilities {
		vulnIDs = append(vulnIDs, vuln.ID)
	}

	// Calculate compliance score based on control coverage
	complianceScore := 75.0 // Base score

	mapping.VulnerabilityIDs = vulnIDs
	mapping.ComplianceScore = complianceScore

	// Generate gap analysis
	gaps = append(gaps, "Missing access control policies (AC-1)")
	gaps = append(gaps, "Insufficient audit logging (AU-2)")
	gaps = append(gaps, "Weak configuration management (CM-2)")

	// Generate recommendations
	recommendations = append(recommendations, "Implement access control policies and procedures")
	recommendations = append(recommendations, "Establish comprehensive audit logging")
	recommendations = append(recommendations, "Implement configuration management controls")

	mapping.GapAnalysis = gaps
	mapping.Recommendations = recommendations

	return mapping
}

// mapToISO27001 maps vulnerabilities to ISO 27001 requirements
func (r *MemoryRepository) mapToISO27001(vulnerabilities []*models.MemoryEntry, mapping *models.ComplianceMapping) *models.ComplianceMapping {
	mapping.Requirement = "ISO 27001 Information Security Management"

	var vulnIDs []string
	var gaps []string
	var recommendations []string

	for _, vuln := range vulnerabilities {
		vulnIDs = append(vulnIDs, vuln.ID)
	}

	// Calculate compliance score
	complianceScore := 70.0 // Base score

	mapping.VulnerabilityIDs = vulnIDs
	mapping.ComplianceScore = complianceScore

	// Generate gap analysis
	gaps = append(gaps, "Missing information security policies")
	gaps = append(gaps, "Insufficient risk assessment procedures")
	gaps = append(gaps, "Weak incident management processes")

	// Generate recommendations
	recommendations = append(recommendations, "Develop comprehensive information security policies")
	recommendations = append(recommendations, "Implement formal risk assessment procedures")
	recommendations = append(recommendations, "Establish incident management processes")

	mapping.GapAnalysis = gaps
	mapping.Recommendations = recommendations

	return mapping
}

// mapToPCIDSS maps vulnerabilities to PCI DSS requirements
func (r *MemoryRepository) mapToPCIDSS(vulnerabilities []*models.MemoryEntry, mapping *models.ComplianceMapping) *models.ComplianceMapping {
	mapping.Requirement = "PCI DSS Payment Card Industry Standards"

	var vulnIDs []string
	var gaps []string
	var recommendations []string

	for _, vuln := range vulnerabilities {
		vulnIDs = append(vulnIDs, vuln.ID)
	}

	// Calculate compliance score
	complianceScore := 65.0 // Base score

	mapping.VulnerabilityIDs = vulnIDs
	mapping.ComplianceScore = complianceScore

	// Generate gap analysis
	gaps = append(gaps, "Missing network security controls (Requirement 1)")
	gaps = append(gaps, "Insufficient data protection measures (Requirement 3)")
	gaps = append(gaps, "Weak access control mechanisms (Requirement 7)")

	// Generate recommendations
	recommendations = append(recommendations, "Implement network security controls and firewalls")
	recommendations = append(recommendations, "Encrypt stored cardholder data")
	recommendations = append(recommendations, "Implement strong access control measures")

	mapping.GapAnalysis = gaps
	mapping.Recommendations = recommendations

	return mapping
}

// storeComplianceMapping stores a compliance mapping in the database
func (r *MemoryRepository) storeComplianceMapping(ctx context.Context, mapping *models.ComplianceMapping) error {
	query := `
		INSERT INTO compliance_mappings (id, session_id, standard, requirement, vulnerability_ids, compliance_score, gap_analysis, recommendations, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	vulnIDsJSON, err := json.Marshal(mapping.VulnerabilityIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal vulnerability IDs: %w", err)
	}

	gapAnalysisJSON, err := json.Marshal(mapping.GapAnalysis)
	if err != nil {
		return fmt.Errorf("failed to marshal gap analysis: %w", err)
	}

	recommendationsJSON, err := json.Marshal(mapping.Recommendations)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		mapping.ID,
		mapping.SessionID,
		mapping.Standard,
		mapping.Requirement,
		vulnIDsJSON,
		mapping.ComplianceScore,
		gapAnalysisJSON,
		recommendationsJSON,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	)

	return err
}
