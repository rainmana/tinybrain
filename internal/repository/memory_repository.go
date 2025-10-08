package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rainmana/tinybrain/internal/models"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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
		ID:                  fmt.Sprintf("task_%d", time.Now().UnixNano()),
		SessionID:           sessionID,
		TaskName:            taskName,
		Stage:               stage,
		Status:              status,
		ProgressPercentage:  progressPercentage,
		Notes:               notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
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
