package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

func main() {
	// Get data directory from environment or use default
	dataDir := os.Getenv("TINYBRAIN_DATA_DIR")
	if dataDir == "" {
		dataDir = "./pb_data"
	}

	log.Printf("Starting TinyBrain v2.0 with data directory: %s", dataDir)

	// Initialize PocketBase directly
	config := pocketbase.Config{
		DefaultDataDir: dataDir,
	}
	app := pocketbase.NewWithConfig(config)

	// Bootstrap the app to ensure DB connections are open
	if err := app.Bootstrap(); err != nil {
		log.Fatalf("Failed to bootstrap PocketBase app: %v", err)
	}

	// Initialize database collections immediately after bootstrap
	log.Println("Initializing database collections...")
	collections := []*pbmodels.Collection{
		database.CreateSessionsCollection(),
		database.CreateMemoryEntriesCollection(),
		database.CreateRelationshipsCollection(),
		database.CreateContextSnapshotsCollection(),
		database.CreateTaskProgressCollection(),
	}

	for _, collection := range collections {
		// Check if collection already exists
		existing, err := app.Dao().FindCollectionByNameOrId(collection.Name)
		if err != nil {
			// Collection doesn't exist, create it
			if err := app.Dao().SaveCollection(collection); err != nil {
				log.Printf("Warning: Failed to create collection %s: %v", collection.Name, err)
			} else {
				log.Printf("Created collection: %s", collection.Name)
			}
		} else {
			log.Printf("Collection '%s' already exists", existing.Name)
		}
	}
	log.Println("Database collections initialization completed")

	// Initialize repositories and services
	sessionRepo := repository.NewSessionRepositoryV2(app)
	memoryRepo := repository.NewMemoryRepositoryV2(app)
	relationshipRepo := repository.NewRelationshipRepositoryV2(app)
	contextRepo := repository.NewContextRepositoryV2(app)
	taskRepo := repository.NewTaskRepositoryV2(app)

	sessionService := services.NewSessionServiceV2(sessionRepo)
	memoryService := services.NewMemoryServiceV2(memoryRepo)
	relationshipService := services.NewRelationshipServiceV2(relationshipRepo)
	contextService := services.NewContextServiceV2(contextRepo)
	taskService := services.NewTaskServiceV2(taskRepo)

	// Add PocketBase web server hook for custom routes
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		log.Println("PocketBase web server starting on :8090")

		// Custom health check endpoint
		e.Router.GET("/health", func(c echo.Context) error {
			return c.JSON(200, map[string]interface{}{
				"status":  "healthy",
				"service": "TinyBrain v2.0",
				"version": "2.0.0",
				"features": []string{
					"session_management",
					"memory_storage",
					"relationship_tracking",
					"context_snapshots",
					"task_progress",
					"pocketbase_database",
					"mcp_protocol",
				},
			})
		})

		// Custom hello world endpoint
		e.Router.GET("/hello", func(c echo.Context) error {
			return c.String(200, "Hello from TinyBrain v2.0!")
		})

		return nil
	})

	// Create MCP server
	mcpServer := server.NewMCPServer("TinyBrain v2.0", "2.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Register Session Management MCP Tools
	mcpServer.AddTool(
		mcp.NewTool("create_session",
			mcp.WithDescription("Create a new LLM interaction session for security assessments"),
			mcp.WithString("name", mcp.Required()),
			mcp.WithString("task_type", mcp.Required()),
			mcp.WithString("description", mcp.Description("Optional session description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, _ := req.RequireString("name")
			taskType, _ := req.RequireString("task_type")
			description := req.GetString("description", "")

			sessionReq := &models.SessionCreateRequest{
				Name:        name,
				TaskType:    taskType,
				Description: description,
				Metadata:    map[string]interface{}{},
			}
			session, err := sessionService.CreateSession(ctx, sessionReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(session)
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("get_session",
			mcp.WithDescription("Retrieve an LLM interaction session by ID"),
			mcp.WithString("id", mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, _ := req.RequireString("id")
			session, err := sessionService.GetSession(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(session)
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("list_sessions",
			mcp.WithDescription("List LLM interaction sessions with optional filtering"),
			mcp.WithString("task_type", mcp.Description("Filter by task type")),
			mcp.WithString("status", mcp.Description("Filter by status")),
			mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Max(100)),
			mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskType := req.GetString("task_type", "")
			status := req.GetString("status", "")
			limit := int(req.GetFloat("limit", 20))
			offset := int(req.GetFloat("offset", 0))

			listReq := &models.SessionListRequest{
				TaskType: taskType,
				Status:   status,
				Limit:    limit,
				Offset:   offset,
			}
			sessions, totalCount, err := sessionService.ListSessions(ctx, listReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(map[string]interface{}{
				"sessions":    sessions,
				"total_count": totalCount,
			})
		},
	)

	// Register Memory Management MCP Tools
	mcpServer.AddTool(
		mcp.NewTool("store_memory",
			mcp.WithDescription("Store a new piece of memory for a session (vulnerabilities, findings, etc.)"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("title", mcp.Required()),
			mcp.WithString("content", mcp.Required()),
			mcp.WithString("category", mcp.Required()),
			mcp.WithNumber("priority", mcp.Required(), mcp.Min(1), mcp.Max(10)),
			mcp.WithNumber("confidence", mcp.Required(), mcp.Min(0.0), mcp.Max(1.0)),
			mcp.WithString("tags", mcp.Description("Comma-separated tags")),
			mcp.WithString("source", mcp.Description("Source of the memory")),
			mcp.WithString("content_type", mcp.Description("Type of content (text, json, etc.)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			title, _ := req.RequireString("title")
			content, _ := req.RequireString("content")
			category, _ := req.RequireString("category")
			priority := int(req.GetFloat("priority", 5))
			confidence := float32(req.GetFloat("confidence", 0.5))
			tagsStr := req.GetString("tags", "")
			source := req.GetString("source", "")
			contentType := req.GetString("content_type", "text")

			// Parse tags
			var tags []string
			if tagsStr != "" {
				// Simple comma-separated parsing
				tagParts := strings.Split(tagsStr, ",")
				for _, tag := range tagParts {
					if trimmed := strings.TrimSpace(tag); trimmed != "" {
						tags = append(tags, trimmed)
					}
				}
			}

			memoryReq := &models.MemoryCreateRequest{
				SessionID:   sessionID,
				Title:       title,
				Content:     content,
				Category:    category,
				Priority:    priority,
				Confidence:  confidence,
				Tags:        tags,
				Source:      source,
				ContentType: contentType,
			}
			memory, err := memoryService.StoreMemory(ctx, memoryReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(memory)
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("search_memories",
			mcp.WithDescription("Search for memories within a session"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("query", mcp.Description("Search query")),
			mcp.WithString("category", mcp.Description("Filter by category")),
			mcp.WithString("tags", mcp.Description("Comma-separated tags to filter by")),
			mcp.WithString("source", mcp.Description("Filter by source")),
			mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Max(100)),
			mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			query := req.GetString("query", "")
			category := req.GetString("category", "")
			tagsStr := req.GetString("tags", "")
			source := req.GetString("source", "")
			limit := int(req.GetFloat("limit", 10))
			offset := int(req.GetFloat("offset", 0))

			// Parse tags
			var tags []string
			if tagsStr != "" {
				tagParts := strings.Split(tagsStr, ",")
				for _, tag := range tagParts {
					if trimmed := strings.TrimSpace(tag); trimmed != "" {
						tags = append(tags, trimmed)
					}
				}
			}

			searchReq := &models.MemorySearchRequest{
				SessionID: sessionID,
				Query:     query,
				Category:  category,
				Tags:      tags,
				Source:    source,
				Limit:     limit,
				Offset:    offset,
			}
			memories, total, err := memoryService.SearchMemories(ctx, searchReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(map[string]interface{}{
				"memories":    memories,
				"total_count": total,
			})
		},
	)

	// Register Relationship Management MCP Tools
	mcpServer.AddTool(
		mcp.NewTool("create_relationship",
			mcp.WithDescription("Create a relationship between two memories"),
			mcp.WithString("source_memory_id", mcp.Required()),
			mcp.WithString("target_memory_id", mcp.Required()),
			mcp.WithString("relationship_type", mcp.Required()),
			mcp.WithNumber("strength", mcp.DefaultNumber(0.5), mcp.Min(0.0), mcp.Max(1.0)),
			mcp.WithString("description", mcp.Description("Relationship description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sourceID, _ := req.RequireString("source_memory_id")
			targetID, _ := req.RequireString("target_memory_id")
			relType, _ := req.RequireString("relationship_type")
			strength := float32(req.GetFloat("strength", 0.5))
			description := req.GetString("description", "")

			relReq := &models.RelationshipCreateRequest{
				SourceID:    sourceID,
				TargetID:    targetID,
				Type:        models.RelationshipType(relType),
				Strength:    strength,
				Description: description,
			}
			relationship, err := relationshipService.CreateRelationship(ctx, relReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(relationship)
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("list_relationships",
			mcp.WithDescription("List relationships based on criteria"),
			mcp.WithString("source_id", mcp.Description("Filter by source memory ID")),
			mcp.WithString("target_id", mcp.Description("Filter by target memory ID")),
			mcp.WithString("type", mcp.Description("Filter by relationship type")),
			mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Max(100)),
			mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sourceID := req.GetString("source_id", "")
			targetID := req.GetString("target_id", "")
			relType := req.GetString("type", "")
			limit := int(req.GetFloat("limit", 20))
			offset := int(req.GetFloat("offset", 0))

			listReq := &models.RelationshipListRequest{
				SourceID: sourceID,
				TargetID: targetID,
				Type:     relType,
				Limit:    limit,
				Offset:   offset,
			}
			relationships, total, err := relationshipService.ListRelationships(ctx, listReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(map[string]interface{}{
				"relationships": relationships,
				"total_count":   total,
			})
		},
	)

	// Register Context Snapshot MCP Tools
	mcpServer.AddTool(
		mcp.NewTool("create_context_snapshot",
			mcp.WithDescription("Create a snapshot of the LLM's context"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("name", mcp.Required()),
			mcp.WithString("context_data", mcp.Required()),
			mcp.WithString("description", mcp.Description("Snapshot description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			name, _ := req.RequireString("name")
			contextDataStr, _ := req.RequireString("context_data")
			description := req.GetString("description", "")

			// Parse context data as JSON
			var contextData map[string]interface{}
			// For now, we'll store it as a simple map with the raw string
			contextData = map[string]interface{}{
				"raw": contextDataStr,
			}

			snapshotReq := &models.ContextSnapshotCreateRequest{
				SessionID:   sessionID,
				Name:        name,
				ContextData: contextData,
				Description: description,
			}
			snapshot, err := contextService.CreateContextSnapshot(ctx, snapshotReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(snapshot)
		},
	)

	// Register Task Progress MCP Tools
	mcpServer.AddTool(
		mcp.NewTool("create_task_progress",
			mcp.WithDescription("Create a new task progress entry for a session"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("task_name", mcp.Required()),
			mcp.WithString("stage", mcp.Required()),
			mcp.WithString("status", mcp.Required()),
			mcp.WithNumber("progress_percentage", mcp.Required(), mcp.Min(0.0), mcp.Max(100.0)),
			mcp.WithString("notes", mcp.Description("Progress notes")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			taskName, _ := req.RequireString("task_name")
			stage, _ := req.RequireString("stage")
			status, _ := req.RequireString("status")
			progressPercentage := float32(req.GetFloat("progress_percentage", 0.0))
			notes := req.GetString("notes", "")

			taskReq := &models.TaskProgressCreateRequest{
				SessionID:          sessionID,
				TaskName:           taskName,
				Stage:              stage,
				Status:             status,
				ProgressPercentage: progressPercentage,
				Notes:              notes,
			}
			task, err := taskService.CreateTaskProgress(ctx, taskReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(task)
		},
	)

	// Register MCP Resources
	mcpServer.AddResource(
		mcp.NewResource(
			"tinybrain://status",
			"TinyBrain Status",
			mcp.WithResourceDescription("Current TinyBrain v2.0 status and capabilities"),
			mcp.WithMIMEType("application/json"),
		),
		func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			statusJSON := `{"status": "healthy", "service": "TinyBrain v2.0", "version": "2.0.0", "message": "Ready for security assessments!", "features": ["session_management", "memory_storage", "relationship_tracking", "context_snapshots", "task_progress", "pocketbase_database", "mcp_protocol"], "uptime": "0s"}`
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     statusJSON,
				},
			}, nil
		},
	)

	// Start PocketBase server in background
	go func() {
		if err := app.Start(); err != nil {
			log.Fatalf("PocketBase server failed to start: %v", err)
		}
	}()

	// Start MCP STDIO server
	log.Println("Starting MCP STDIO server...")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("MCP STDIO server failed to start: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down TinyBrain v2.0...")
}
