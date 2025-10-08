package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/mcp"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/rainmana/tinybrain/internal/repository"
	"github.com/charmbracelet/log"
)

// TinyBrainServer represents the main MCP server for security-focused memory storage
type TinyBrainServer struct {
	db         *database.Database
	repo       *repository.MemoryRepository
	logger     *log.Logger
	dbPath     string
}

func main() {
	// Initialize logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "TinyBrain 🧠 ",
		Level:           log.InfoLevel,
	})

	// Get database path from environment or use default
	dbPath := os.Getenv("TINYBRAIN_DB_PATH")
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatal("Failed to get user home directory", "error", err)
		}
		dbPath = filepath.Join(homeDir, ".tinybrain", "memory.db")
	}

	// Initialize database
	db, err := database.NewDatabase(dbPath, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", "error", err)
	}
	defer db.Close()

	// Initialize repository
	repo := repository.NewMemoryRepository(db.GetDB(), logger)

	// Create server instance
	tinyBrain := &TinyBrainServer{
		db:     db,
		repo:   repo,
		logger: logger,
		dbPath: dbPath,
	}

	// Create MCP server
	mcpServer := mcp.NewServer("TinyBrain Memory Storage", "1.0.0", 
		"Security-focused LLM memory storage MCP server", logger)

	// Register tools
	tinyBrain.registerTools(mcpServer)

	logger.Info("Starting TinyBrain MCP Server", "db_path", dbPath)

	// Start server
	if err := mcpServer.ServeStdio(); err != nil {
		logger.Fatal("Server error", "error", err)
	}
}

// registerTools registers all MCP tools for memory operations
func (t *TinyBrainServer) registerTools(s *mcp.Server) {
	// Session management tools
	s.AddTool("create_session", 
		"Create a new security-focused session for tracking LLM interactions",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the session",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the session",
				},
				"task_type": map[string]interface{}{
					"type":        "string",
					"description": "Type of security task: security_review, penetration_test, exploit_dev, vulnerability_analysis, threat_modeling, incident_response, general",
				},
				"metadata": map[string]interface{}{
					"type":        "string",
					"description": "JSON metadata for the session",
				},
			},
			"required": []string{"name", "task_type"},
		},
		t.handleCreateSession,
	)

	s.AddTool("get_session",
		"Retrieve a session by ID",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to retrieve",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleGetSession,
	)

	s.AddTool("list_sessions",
		"List all sessions with optional filtering",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_type": map[string]interface{}{
					"type":        "string",
					"description": "Filter by task type",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Filter by status: active, paused, completed, archived",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of sessions to return (default: 50)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of sessions to skip (default: 0)",
				},
			},
		},
		t.handleListSessions,
	)

	// Memory entry tools
	s.AddTool("store_memory",
		"Store a new piece of information in memory with security-focused categorization",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session this memory belongs to",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title/summary of the memory",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "Content of the memory",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"description": "Category: finding, vulnerability, exploit, payload, technique, tool, reference, context, hypothesis, evidence, recommendation, note",
				},
				"content_type": map[string]interface{}{
					"type":        "string",
					"description": "Content type: text, code, json, yaml, markdown, binary_ref (default: text)",
				},
				"priority": map[string]interface{}{
					"type":        "number",
					"description": "Priority level 0-10 (default: 5)",
				},
				"confidence": map[string]interface{}{
					"type":        "number",
					"description": "Confidence level 0.0-1.0 (default: 0.5)",
				},
				"tags": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of tags",
				},
				"source": map[string]interface{}{
					"type":        "string",
					"description": "Source of this information",
				},
			},
			"required": []string{"session_id", "title", "content", "category"},
		},
		t.handleStoreMemory,
	)

	s.AddTool("get_memory",
		"Retrieve a specific memory entry by ID",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"memory_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the memory entry to retrieve",
				},
			},
			"required": []string{"memory_id"},
		},
		t.handleGetMemory,
	)

	s.AddTool("search_memories",
		"Search for memories using various search strategies optimized for security tasks",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query",
				},
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "Limit search to specific session",
				},
				"search_type": map[string]interface{}{
					"type":        "string",
					"description": "Search type: semantic, exact, fuzzy, tag, category, relationship (default: semantic)",
				},
				"categories": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of categories to filter by",
				},
				"tags": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of tags to filter by",
				},
				"min_priority": map[string]interface{}{
					"type":        "number",
					"description": "Minimum priority level (0-10)",
				},
				"min_confidence": map[string]interface{}{
					"type":        "number",
					"description": "Minimum confidence level (0.0-1.0)",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results (default: 20)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of results to skip (default: 0)",
				},
			},
			"required": []string{"query"},
		},
		t.handleSearchMemories,
	)

	s.AddTool("get_related_memories",
		"Get memories related to a specific memory entry",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"memory_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the memory entry to find related memories for",
				},
				"relationship_type": map[string]interface{}{
					"type":        "string",
					"description": "Type of relationship: depends_on, causes, mitigates, exploits, references, contradicts, supports, related_to, parent_of, child_of",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of related memories (default: 10)",
				},
			},
			"required": []string{"memory_id"},
		},
		t.handleGetRelatedMemories,
	)

	// Relationship tools
	s.AddTool("create_relationship",
		"Create a relationship between two memory entries",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"source_memory_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the source memory entry",
				},
				"target_memory_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the target memory entry",
				},
				"relationship_type": map[string]interface{}{
					"type":        "string",
					"description": "Type of relationship",
				},
				"strength": map[string]interface{}{
					"type":        "number",
					"description": "Strength of relationship 0.0-1.0 (default: 0.5)",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the relationship",
				},
			},
			"required": []string{"source_memory_id", "target_memory_id", "relationship_type"},
		},
		t.handleCreateRelationship,
	)

	// Context management tools
	s.AddTool("get_context_summary",
		"Get a summary of relevant memories for current context",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to get context for",
				},
				"current_task": map[string]interface{}{
					"type":        "string",
					"description": "Description of current task for context relevance",
				},
				"max_memories": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of memories to include (default: 20)",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleGetContextSummary,
	)

	// Task progress tools
	s.AddTool("update_task_progress",
		"Update progress on a multi-stage security task",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"task_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the task",
				},
				"stage": map[string]interface{}{
					"type":        "string",
					"description": "Current stage of the task",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Status: pending, in_progress, completed, failed, blocked",
				},
				"progress_percentage": map[string]interface{}{
					"type":        "number",
					"description": "Progress percentage 0-100",
				},
				"notes": map[string]interface{}{
					"type":        "string",
					"description": "Notes about the current progress",
				},
			},
			"required": []string{"session_id", "task_name", "stage", "status"},
		},
		t.handleUpdateTaskProgress,
	)

	// Utility tools
	s.AddTool("get_database_stats",
		"Get database statistics and health information",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		t.handleGetDatabaseStats,
	)

	// Context snapshot tools
	s.AddTool("create_context_snapshot", 
		"Create a snapshot of the current context for a session",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the context snapshot",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the context snapshot",
				},
				"context_data": map[string]interface{}{
					"type":        "string",
					"description": "JSON string containing context data",
				},
			},
			"required": []string{"session_id", "name"},
		},
		t.handleCreateContextSnapshot,
	)

	s.AddTool("get_context_snapshot", 
		"Retrieve a context snapshot by ID",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"snapshot_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the context snapshot to retrieve",
				},
			},
			"required": []string{"snapshot_id"},
		},
		t.handleGetContextSnapshot,
	)

	s.AddTool("list_context_snapshots", 
		"List context snapshots for a session",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of snapshots to return (default: 20)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of snapshots to skip (default: 0)",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleListContextSnapshots,
	)

	// Task progress tools
	s.AddTool("create_task_progress", 
		"Create a new task progress entry for tracking multi-stage security tasks",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"task_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the task",
				},
				"stage": map[string]interface{}{
					"type":        "string",
					"description": "Current stage of the task",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Status: pending, in_progress, completed, failed, blocked",
				},
				"progress_percentage": map[string]interface{}{
					"type":        "number",
					"description": "Progress percentage 0-100",
				},
				"notes": map[string]interface{}{
					"type":        "string",
					"description": "Notes about the current progress",
				},
			},
			"required": []string{"session_id", "task_name", "stage", "status"},
		},
		t.handleCreateTaskProgress,
	)

	s.AddTool("get_task_progress", 
		"Retrieve a task progress entry by ID",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the task progress entry",
				},
			},
			"required": []string{"task_id"},
		},
		t.handleGetTaskProgress,
	)

	s.AddTool("list_task_progress", 
		"List task progress entries for a session",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Filter by status: pending, in_progress, completed, failed, blocked",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of tasks to return (default: 20)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of tasks to skip (default: 0)",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleListTaskProgress,
	)

	s.AddTool("find_similar_memories",
		"Find memories similar to the given content for deduplication",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "Content to find similar memories for",
				},
				"threshold": map[string]interface{}{
					"type":        "number",
					"description": "Similarity threshold (0.0-1.0, default: 0.7)",
				},
			},
			"required": []string{"session_id", "content"},
		},
		t.handleFindSimilarMemories,
	)

	s.AddTool("check_duplicates",
		"Check if a memory entry is a duplicate of existing entries",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title of the memory entry",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "Content of the memory entry",
				},
			},
			"required": []string{"session_id", "title", "content"},
		},
		t.handleCheckDuplicates,
	)

	s.AddTool("export_session_data",
		"Export all data for a session in JSON format for backup or migration",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to export",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleExportSessionData,
	)

	s.AddTool("import_session_data",
		"Import session data from JSON format for restoration or migration",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"import_data": map[string]interface{}{
					"type":        "string",
					"description": "JSON string containing the session data to import",
				},
			},
			"required": []string{"import_data"},
		},
		t.handleImportSessionData,
	)

	s.AddTool("get_security_templates",
		"Get predefined templates for common security patterns",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		t.handleGetSecurityTemplates,
	)

	s.AddTool("create_memory_from_template",
		"Create a memory entry from a predefined security template",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"template_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the template to use",
				},
				"replacements": map[string]interface{}{
					"type":        "string",
					"description": "JSON string containing placeholder replacements",
				},
			},
			"required": []string{"session_id", "template_name"},
		},
		t.handleCreateMemoryFromTemplate,
	)

	s.AddTool("batch_create_memories",
		"Create multiple memory entries in a single transaction for bulk operations",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session",
				},
				"memory_requests": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of memory creation requests",
				},
			},
			"required": []string{"session_id", "memory_requests"},
		},
		t.handleBatchCreateMemories,
	)

	s.AddTool("batch_update_memories",
		"Update multiple memory entries in a single transaction for bulk operations",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"memory_updates": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of memory update requests",
				},
			},
			"required": []string{"memory_updates"},
		},
		t.handleBatchUpdateMemories,
	)

	s.AddTool("batch_delete_memories",
		"Delete multiple memory entries in a single transaction for bulk operations",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"memory_ids": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of memory IDs to delete",
				},
			},
			"required": []string{"memory_ids"},
		},
		t.handleBatchDeleteMemories,
	)

	s.AddTool("cleanup_old_memories",
		"Remove memories older than specified age with optional dry run",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"max_age_days": map[string]interface{}{
					"type":        "number",
					"description": "Maximum age in days for memories to keep",
				},
				"dry_run": map[string]interface{}{
					"type":        "boolean",
					"description": "If true, only show what would be deleted without actually deleting",
				},
			},
			"required": []string{"max_age_days"},
		},
		t.handleCleanupOldMemories,
	)

	s.AddTool("cleanup_low_priority_memories",
		"Remove memories with low priority and confidence scores",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"max_priority": map[string]interface{}{
					"type":        "number",
					"description": "Maximum priority level to consider for deletion (0-10)",
				},
				"max_confidence": map[string]interface{}{
					"type":        "number",
					"description": "Maximum confidence level to consider for deletion (0.0-1.0)",
				},
				"dry_run": map[string]interface{}{
					"type":        "boolean",
					"description": "If true, only show what would be deleted without actually deleting",
				},
			},
			"required": []string{"max_priority", "max_confidence"},
		},
		t.handleCleanupLowPriorityMemories,
	)

	s.AddTool("cleanup_unused_memories",
		"Remove memories that haven't been accessed recently",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"max_unused_days": map[string]interface{}{
					"type":        "number",
					"description": "Maximum days since last access to consider for deletion",
				},
				"dry_run": map[string]interface{}{
					"type":        "boolean",
					"description": "If true, only show what would be deleted without actually deleting",
				},
			},
			"required": []string{"max_unused_days"},
		},
		t.handleCleanupUnusedMemories,
	)

	s.AddTool("get_memory_stats",
		"Get comprehensive statistics about memory usage and aging",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		t.handleGetMemoryStats,
	)

	s.AddTool("get_detailed_memory_info",
		"Get comprehensive debugging information about a specific memory entry",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"memory_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the memory entry to get detailed info for",
				},
			},
			"required": []string{"memory_id"},
		},
		t.handleGetDetailedMemoryInfo,
	)

	s.AddTool("get_system_diagnostics",
		"Get comprehensive system diagnostics and debugging information",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		t.handleGetSystemDiagnostics,
	)

	s.AddTool("semantic_search",
		"Perform semantic search using embeddings for finding conceptually similar memories",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query for semantic matching",
				},
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to search within",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results to return (default: 20)",
				},
			},
			"required": []string{"query", "session_id"},
		},
		t.handleSemanticSearch,
	)

	s.AddTool("generate_embedding",
		"Generate an embedding vector for text (placeholder for future AI integration)",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"text": map[string]interface{}{
					"type":        "string",
					"description": "Text to generate embedding for",
				},
			},
			"required": []string{"text"},
		},
		t.handleGenerateEmbedding,
	)

	s.AddTool("calculate_similarity",
		"Calculate semantic similarity between two embeddings",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"embedding1": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of first embedding vector",
				},
				"embedding2": map[string]interface{}{
					"type":        "string",
					"description": "JSON array of second embedding vector",
				},
			},
			"required": []string{"embedding1", "embedding2"},
		},
		t.handleCalculateSimilarity,
	)

	s.AddTool("get_notifications",
		"Get notifications and alerts for a session",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to get notifications for",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of notifications to return (default: 20)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of notifications to skip (default: 0)",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleGetNotifications,
	)

	s.AddTool("mark_notification_read",
		"Mark a notification as read",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"notification_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the notification to mark as read",
				},
			},
			"required": []string{"notification_id"},
		},
		t.handleMarkNotificationRead,
	)

	s.AddTool("check_high_priority_memories",
		"Check for high-priority memories and create notifications",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to check for high-priority memories",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleCheckHighPriorityMemories,
	)

	s.AddTool("check_duplicate_memories",
		"Check for duplicate memories and create notifications",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the session to check for duplicate memories",
				},
			},
			"required": []string{"session_id"},
		},
		t.handleCheckDuplicateMemories,
	)

	s.AddTool("health_check",
		"Perform a health check on the database and server",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		t.handleHealthCheck,
	)
}


// Tool handlers

func (t *TinyBrainServer) handleCreateSession(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is required")
	}

	taskType, ok := params["task_type"].(string)
	if !ok {
		return nil, fmt.Errorf("task_type is required")
	}

	description, _ := params["description"].(string)
	metadataStr, _ := params["metadata"].(string)

	var metadata map[string]interface{}
	if metadataStr != "" {
		if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata JSON: %v", err)
		}
	}

	session := &models.Session{
		ID:          fmt.Sprintf("session_%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		TaskType:    taskType,
		Status:      "active",
		Metadata:    metadata,
	}

	if err := t.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

func (t *TinyBrainServer) handleGetSession(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	session, err := t.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return session, nil
}

func (t *TinyBrainServer) handleListSessions(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	taskType, _ := params["task_type"].(string)
	status, _ := params["status"].(string)
	
	limit := 50
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	offset := 0
	if offsetVal, ok := params["offset"].(float64); ok {
		offset = int(offsetVal)
	}

	sessions, err := t.repo.ListSessions(ctx, taskType, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %v", err)
	}

	return sessions, nil
}

func (t *TinyBrainServer) handleStoreMemory(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	title, ok := params["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content is required")
	}

	category, ok := params["category"].(string)
	if !ok {
		return nil, fmt.Errorf("category is required")
	}

	contentType, _ := params["content_type"].(string)
	priority := 5
	if priorityVal, ok := params["priority"].(float64); ok {
		priority = int(priorityVal)
	}

	confidence := 0.5
	if confidenceVal, ok := params["confidence"].(float64); ok {
		confidence = confidenceVal
	}

	var tags []string
	if tagsStr, ok := params["tags"].(string); ok && tagsStr != "" {
		if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
			return nil, fmt.Errorf("invalid tags JSON: %v", err)
		}
	}

	source, _ := params["source"].(string)

	memoryReq := &models.CreateMemoryEntryRequest{
		SessionID:   sessionID,
		Title:       title,
		Content:     content,
		ContentType: contentType,
		Category:    category,
		Priority:    priority,
		Confidence:  confidence,
		Tags:        tags,
		Source:      source,
	}

	entry, err := t.repo.CreateMemoryEntry(ctx, memoryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to store memory: %v", err)
	}

	return entry, nil
}

func (t *TinyBrainServer) handleGetMemory(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	memoryID, ok := params["memory_id"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_id is required")
	}

	entry, err := t.repo.GetMemoryEntry(ctx, memoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %v", err)
	}

	return entry, nil
}

func (t *TinyBrainServer) handleSearchMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	query, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query is required")
	}

	sessionID, _ := params["session_id"].(string)
	searchType, _ := params["search_type"].(string)
	if searchType == "" {
		searchType = "semantic"
	}

	var categories []string
	if categoriesStr, ok := params["categories"].(string); ok && categoriesStr != "" {
		if err := json.Unmarshal([]byte(categoriesStr), &categories); err != nil {
			return nil, fmt.Errorf("invalid categories JSON: %v", err)
		}
	}

	var tags []string
	if tagsStr, ok := params["tags"].(string); ok && tagsStr != "" {
		if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
			return nil, fmt.Errorf("invalid tags JSON: %v", err)
		}
	}

	minPriority := 0
	if priorityVal, ok := params["min_priority"].(float64); ok {
		minPriority = int(priorityVal)
	}

	minConfidence := 0.0
	if confidenceVal, ok := params["min_confidence"].(float64); ok {
		minConfidence = confidenceVal
	}

	limit := 20
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	offset := 0
	if offsetVal, ok := params["offset"].(float64); ok {
		offset = int(offsetVal)
	}

	searchReq := &models.SearchRequest{
		Query:        query,
		SessionID:    sessionID,
		Categories:   categories,
		Tags:         tags,
		MinPriority:  minPriority,
		MinConfidence: minConfidence,
		Limit:        limit,
		Offset:       offset,
		SearchType:   searchType,
	}

	results, err := t.repo.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %v", err)
	}

	return results, nil
}

func (t *TinyBrainServer) handleGetRelatedMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	memoryID, ok := params["memory_id"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_id is required")
	}

	relationshipType, _ := params["relationship_type"].(string)
	
	limit := 10
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	entries, err := t.repo.GetRelatedEntries(ctx, memoryID, relationshipType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get related memories: %v", err)
	}

	return entries, nil
}

func (t *TinyBrainServer) handleCreateRelationship(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sourceID, ok := params["source_memory_id"].(string)
	if !ok {
		return nil, fmt.Errorf("source_memory_id is required")
	}

	targetID, ok := params["target_memory_id"].(string)
	if !ok {
		return nil, fmt.Errorf("target_memory_id is required")
	}

	relationshipType, ok := params["relationship_type"].(string)
	if !ok {
		return nil, fmt.Errorf("relationship_type is required")
	}

	strength := 0.5
	if strengthVal, ok := params["strength"].(float64); ok {
		strength = strengthVal
	}

	description, _ := params["description"].(string)

	relReq := &models.CreateRelationshipRequest{
		SourceEntryID:    sourceID,
		TargetEntryID:    targetID,
		RelationshipType: relationshipType,
		Strength:         strength,
		Description:      description,
	}

	relationship, err := t.repo.CreateRelationship(ctx, relReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %v", err)
	}

	return relationship, nil
}

func (t *TinyBrainServer) handleGetContextSummary(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	currentTask, _ := params["current_task"].(string)
	
	maxMemories := 20
	if maxVal, ok := params["max_memories"].(float64); ok {
		maxMemories = int(maxVal)
	}

	// Search for relevant memories
	searchReq := &models.SearchRequest{
		Query:      currentTask,
		SessionID:  sessionID,
		Limit:      maxMemories,
		SearchType: "semantic",
	}

	results, err := t.repo.SearchMemoryEntries(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get context summary: %v", err)
	}

	// Create summary
	var relevantMemories []models.MemoryEntry
	for _, result := range results {
		relevantMemories = append(relevantMemories, result.MemoryEntry)
	}

	summary := &models.MemorySummary{
		SessionID:         sessionID,
		RelevantMemories:  relevantMemories,
		Summary:           fmt.Sprintf("Found %d relevant memories for current context", len(relevantMemories)),
		GeneratedAt:       time.Now(),
	}

	return summary, nil
}

func (t *TinyBrainServer) handleUpdateTaskProgress(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	taskName, ok := params["task_name"].(string)
	if !ok {
		return nil, fmt.Errorf("task_name is required")
	}

	stage, ok := params["stage"].(string)
	if !ok {
		return nil, fmt.Errorf("stage is required")
	}

	status, ok := params["status"].(string)
	if !ok {
		return nil, fmt.Errorf("status is required")
	}

	notes, _ := params["notes"].(string)
	
	progressPercentage := 0
	if progressVal, ok := params["progress_percentage"].(float64); ok {
		progressPercentage = int(progressVal)
	}

	// First, get the task ID by finding the task with the given name and session
	tasks, err := t.repo.ListTaskProgress(ctx, sessionID, "", 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %v", err)
	}

	var taskID string
	for _, task := range tasks {
		if task.TaskName == taskName {
			taskID = task.ID
			break
		}
	}

	if taskID == "" {
		return nil, fmt.Errorf("task not found: %s", taskName)
	}

	progress, err := t.repo.UpdateTaskProgress(ctx, taskID, stage, status, notes, progressPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to update task progress: %v", err)
	}

	return progress, nil
}

func (t *TinyBrainServer) handleFindSimilarMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content is required")
	}

	threshold := 0.7
	if thresholdVal, ok := params["threshold"].(float64); ok {
		threshold = thresholdVal
	}

	similarMemories, err := t.repo.FindSimilarMemories(ctx, sessionID, content, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar memories: %v", err)
	}

	return similarMemories, nil
}

func (t *TinyBrainServer) handleCheckDuplicates(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	title, ok := params["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content is required")
	}

	duplicates, err := t.repo.CheckForDuplicates(ctx, sessionID, title, content)
	if err != nil {
		return nil, fmt.Errorf("failed to check for duplicates: %v", err)
	}

	return duplicates, nil
}

func (t *TinyBrainServer) handleExportSessionData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	exportData, err := t.repo.ExportSessionData(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to export session data: %v", err)
	}

	return exportData, nil
}

func (t *TinyBrainServer) handleImportSessionData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	importDataStr, ok := params["import_data"].(string)
	if !ok {
		return nil, fmt.Errorf("import_data is required")
	}

	var importData map[string]interface{}
	if err := json.Unmarshal([]byte(importDataStr), &importData); err != nil {
		return nil, fmt.Errorf("invalid import_data JSON: %v", err)
	}

	sessionID, err := t.repo.ImportSessionData(ctx, importData)
	if err != nil {
		return nil, fmt.Errorf("failed to import session data: %v", err)
	}

	return map[string]interface{}{
		"imported_session_id": sessionID,
		"message":             "Session data imported successfully",
	}, nil
}

func (t *TinyBrainServer) handleGetSecurityTemplates(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	templates := t.repo.GetSecurityTemplates()
	return templates, nil
}

func (t *TinyBrainServer) handleCreateMemoryFromTemplate(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	templateName, ok := params["template_name"].(string)
	if !ok {
		return nil, fmt.Errorf("template_name is required")
	}

	replacements := make(map[string]string)
	if replacementsStr, ok := params["replacements"].(string); ok && replacementsStr != "" {
		if err := json.Unmarshal([]byte(replacementsStr), &replacements); err != nil {
			return nil, fmt.Errorf("invalid replacements JSON: %v", err)
		}
	}

	memory, err := t.repo.CreateMemoryFromTemplate(ctx, sessionID, templateName, replacements)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory from template: %v", err)
	}

	return memory, nil
}

func (t *TinyBrainServer) handleBatchCreateMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	memoryRequestsStr, ok := params["memory_requests"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_requests is required")
	}

	var memoryRequests []*models.CreateMemoryEntryRequest
	if err := json.Unmarshal([]byte(memoryRequestsStr), &memoryRequests); err != nil {
		return nil, fmt.Errorf("invalid memory_requests JSON: %v", err)
	}

	memories, err := t.repo.BatchCreateMemoryEntries(ctx, sessionID, memoryRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create memories: %v", err)
	}

	return map[string]interface{}{
		"created_memories": memories,
		"count":            len(memories),
	}, nil
}

func (t *TinyBrainServer) handleBatchUpdateMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	memoryUpdatesStr, ok := params["memory_updates"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_updates is required")
	}

	var memoryUpdates []*models.UpdateMemoryEntryRequest
	if err := json.Unmarshal([]byte(memoryUpdatesStr), &memoryUpdates); err != nil {
		return nil, fmt.Errorf("invalid memory_updates JSON: %v", err)
	}

	memories, err := t.repo.BatchUpdateMemoryEntries(ctx, memoryUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to batch update memories: %v", err)
	}

	return map[string]interface{}{
		"updated_memories": memories,
		"count":            len(memories),
	}, nil
}

func (t *TinyBrainServer) handleBatchDeleteMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	memoryIDsStr, ok := params["memory_ids"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_ids is required")
	}

	var memoryIDs []string
	if err := json.Unmarshal([]byte(memoryIDsStr), &memoryIDs); err != nil {
		return nil, fmt.Errorf("invalid memory_ids JSON: %v", err)
	}

	err := t.repo.BatchDeleteMemoryEntries(ctx, memoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch delete memories: %v", err)
	}

	return map[string]interface{}{
		"deleted_count": len(memoryIDs),
		"message":       "Memory entries deleted successfully",
	}, nil
}

func (t *TinyBrainServer) handleCleanupOldMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	maxAgeDays, ok := params["max_age_days"].(float64)
	if !ok {
		return nil, fmt.Errorf("max_age_days is required")
	}

	dryRun := false
	if dryRunVal, ok := params["dry_run"].(bool); ok {
		dryRun = dryRunVal
	}

	deletedCount, err := t.repo.CleanupOldMemories(ctx, int(maxAgeDays), dryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup old memories: %v", err)
	}

	return map[string]interface{}{
		"deleted_count": deletedCount,
		"max_age_days":  int(maxAgeDays),
		"dry_run":       dryRun,
		"message":       fmt.Sprintf("Cleaned up %d old memories", deletedCount),
	}, nil
}

func (t *TinyBrainServer) handleCleanupLowPriorityMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	maxPriority, ok := params["max_priority"].(float64)
	if !ok {
		return nil, fmt.Errorf("max_priority is required")
	}

	maxConfidence, ok := params["max_confidence"].(float64)
	if !ok {
		return nil, fmt.Errorf("max_confidence is required")
	}

	dryRun := false
	if dryRunVal, ok := params["dry_run"].(bool); ok {
		dryRun = dryRunVal
	}

	deletedCount, err := t.repo.CleanupLowPriorityMemories(ctx, int(maxPriority), maxConfidence, dryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup low priority memories: %v", err)
	}

	return map[string]interface{}{
		"deleted_count":  deletedCount,
		"max_priority":   int(maxPriority),
		"max_confidence": maxConfidence,
		"dry_run":        dryRun,
		"message":        fmt.Sprintf("Cleaned up %d low priority memories", deletedCount),
	}, nil
}

func (t *TinyBrainServer) handleCleanupUnusedMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	maxUnusedDays, ok := params["max_unused_days"].(float64)
	if !ok {
		return nil, fmt.Errorf("max_unused_days is required")
	}

	dryRun := false
	if dryRunVal, ok := params["dry_run"].(bool); ok {
		dryRun = dryRunVal
	}

	deletedCount, err := t.repo.CleanupUnusedMemories(ctx, int(maxUnusedDays), dryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup unused memories: %v", err)
	}

	return map[string]interface{}{
		"deleted_count":    deletedCount,
		"max_unused_days":  int(maxUnusedDays),
		"dry_run":          dryRun,
		"message":          fmt.Sprintf("Cleaned up %d unused memories", deletedCount),
	}, nil
}

func (t *TinyBrainServer) handleGetMemoryStats(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	stats, err := t.repo.GetMemoryStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %v", err)
	}

	return stats, nil
}

func (t *TinyBrainServer) handleGetDetailedMemoryInfo(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	memoryID, ok := params["memory_id"].(string)
	if !ok {
		return nil, fmt.Errorf("memory_id is required")
	}

	detailedInfo, err := t.repo.GetDetailedMemoryInfo(ctx, memoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed memory info: %v", err)
	}

	return detailedInfo, nil
}

func (t *TinyBrainServer) handleGetSystemDiagnostics(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	diagnostics, err := t.repo.GetSystemDiagnostics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system diagnostics: %v", err)
	}

	return diagnostics, nil
}

func (t *TinyBrainServer) handleSemanticSearch(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	query, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query is required")
	}

	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	limit := 20
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	memories, err := t.repo.SemanticSearch(ctx, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic search: %v", err)
	}

	return map[string]interface{}{
		"memories": memories,
		"count":    len(memories),
		"query":    query,
		"session_id": sessionID,
	}, nil
}

func (t *TinyBrainServer) handleGenerateEmbedding(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text is required")
	}

	embedding, err := t.repo.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %v", err)
	}

	return map[string]interface{}{
		"embedding": embedding,
		"dimension": len(embedding),
		"text":      text,
	}, nil
}

func (t *TinyBrainServer) handleCalculateSimilarity(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	embedding1Str, ok := params["embedding1"].(string)
	if !ok {
		return nil, fmt.Errorf("embedding1 is required")
	}

	embedding2Str, ok := params["embedding2"].(string)
	if !ok {
		return nil, fmt.Errorf("embedding2 is required")
	}

	var embedding1, embedding2 []float64
	if err := json.Unmarshal([]byte(embedding1Str), &embedding1); err != nil {
		return nil, fmt.Errorf("invalid embedding1 JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(embedding2Str), &embedding2); err != nil {
		return nil, fmt.Errorf("invalid embedding2 JSON: %v", err)
	}

	similarity, err := t.repo.CalculateSemanticSimilarity(embedding1, embedding2)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate similarity: %v", err)
	}

	return map[string]interface{}{
		"similarity": similarity,
		"embedding1_dimension": len(embedding1),
		"embedding2_dimension": len(embedding2),
	}, nil
}

func (t *TinyBrainServer) handleGetNotifications(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	limit := 20
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	offset := 0
	if offsetVal, ok := params["offset"].(float64); ok {
		offset = int(offsetVal)
	}

	notifications, err := t.repo.GetNotifications(ctx, sessionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %v", err)
	}

	return map[string]interface{}{
		"notifications": notifications,
		"count":         len(notifications),
		"session_id":    sessionID,
		"limit":         limit,
		"offset":        offset,
	}, nil
}

func (t *TinyBrainServer) handleMarkNotificationRead(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	notificationID, ok := params["notification_id"].(string)
	if !ok {
		return nil, fmt.Errorf("notification_id is required")
	}

	err := t.repo.MarkNotificationRead(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %v", err)
	}

	return map[string]interface{}{
		"notification_id": notificationID,
		"message":         "Notification marked as read",
	}, nil
}

func (t *TinyBrainServer) handleCheckHighPriorityMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	err := t.repo.CheckForHighPriorityMemories(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check high priority memories: %v", err)
	}

	return map[string]interface{}{
		"session_id": sessionID,
		"message":    "High priority memory check completed",
	}, nil
}

func (t *TinyBrainServer) handleCheckDuplicateMemories(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	err := t.repo.CheckForDuplicateMemories(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate memories: %v", err)
	}

	return map[string]interface{}{
		"session_id": sessionID,
		"message":    "Duplicate memory check completed",
	}, nil
}

func (t *TinyBrainServer) handleGetDatabaseStats(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	stats, err := t.db.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get database stats: %v", err)
	}

	return stats, nil
}

func (t *TinyBrainServer) handleHealthCheck(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	if err := t.db.HealthCheck(); err != nil {
		return nil, fmt.Errorf("health check failed: %v", err)
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"db_path":   t.dbPath,
	}

	return health, nil
}

func (t *TinyBrainServer) handleCreateContextSnapshot(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is required")
	}

	description, _ := params["description"].(string)
	contextDataStr, _ := params["context_data"].(string)

	var contextData map[string]interface{}
	if contextDataStr != "" {
		if err := json.Unmarshal([]byte(contextDataStr), &contextData); err != nil {
			return nil, fmt.Errorf("invalid context_data JSON: %v", err)
		}
	}

	snapshot, err := t.repo.CreateContextSnapshot(ctx, sessionID, name, description, contextData)
	if err != nil {
		return nil, fmt.Errorf("failed to create context snapshot: %v", err)
	}

	return snapshot, nil
}

func (t *TinyBrainServer) handleGetContextSnapshot(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	snapshotID, ok := params["snapshot_id"].(string)
	if !ok {
		return nil, fmt.Errorf("snapshot_id is required")
	}

	snapshot, err := t.repo.GetContextSnapshot(ctx, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get context snapshot: %v", err)
	}

	return snapshot, nil
}

func (t *TinyBrainServer) handleListContextSnapshots(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	limit := 20
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	offset := 0
	if offsetVal, ok := params["offset"].(float64); ok {
		offset = int(offsetVal)
	}

	snapshots, err := t.repo.ListContextSnapshots(ctx, sessionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list context snapshots: %v", err)
	}

	return snapshots, nil
}

func (t *TinyBrainServer) handleCreateTaskProgress(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	taskName, ok := params["task_name"].(string)
	if !ok {
		return nil, fmt.Errorf("task_name is required")
	}

	stage, ok := params["stage"].(string)
	if !ok {
		return nil, fmt.Errorf("stage is required")
	}

	status, ok := params["status"].(string)
	if !ok {
		return nil, fmt.Errorf("status is required")
	}

	notes, _ := params["notes"].(string)
	
	progressPercentage := 0
	if progressVal, ok := params["progress_percentage"].(float64); ok {
		progressPercentage = int(progressVal)
	}

	progress, err := t.repo.CreateTaskProgress(ctx, sessionID, taskName, stage, status, notes, progressPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to create task progress: %v", err)
	}

	return progress, nil
}

func (t *TinyBrainServer) handleGetTaskProgress(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	taskID, ok := params["task_id"].(string)
	if !ok {
		return nil, fmt.Errorf("task_id is required")
	}

	progress, err := t.repo.GetTaskProgress(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task progress: %v", err)
	}

	return progress, nil
}

func (t *TinyBrainServer) handleListTaskProgress(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	sessionID, ok := params["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id is required")
	}

	status, _ := params["status"].(string)
	
	limit := 20
	if limitVal, ok := params["limit"].(float64); ok {
		limit = int(limitVal)
	}

	offset := 0
	if offsetVal, ok := params["offset"].(float64); ok {
		offset = int(offsetVal)
	}

	tasks, err := t.repo.ListTaskProgress(ctx, sessionID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list task progress: %v", err)
	}

	return tasks, nil
}

