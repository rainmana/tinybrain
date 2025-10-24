package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

func main() {
	// Get data directory from environment or use default
	dataDir := os.Getenv("TINYBRAIN_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	// Initialize PocketBase client
	pbClient, err := database.NewPocketBaseClient(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase client: %v", err)
	}
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	if err := pbClient.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap database: %v", err)
	}

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(pbClient.GetApp())
	memoryRepo := repository.NewMemoryRepository(pbClient.GetApp())
	relationshipRepo := repository.NewRelationshipRepository(pbClient.GetApp())
	contextRepo := repository.NewContextRepository(pbClient.GetApp())
	taskRepo := repository.NewTaskRepository(pbClient.GetApp())

	// Initialize services
	sessionService := services.NewSessionService(sessionRepo)
	memoryService := services.NewMemoryService(memoryRepo)
	relationshipService := services.NewRelationshipService(relationshipRepo)
	contextService := services.NewContextService(contextRepo)
	taskService := services.NewTaskService(taskRepo)

	// Create MCP server
	mcpServer := server.NewMCPServer("TinyBrain v2.0", "2.0.0",
		server.WithAllCapabilities(),
		server.WithRecovery(),
		server.WithHooks(&server.Hooks{
			OnSessionStart: func(sessionID string) {
				log.Printf("MCP session started: %s", sessionID)
			},
			OnSessionEnd: func(sessionID string) {
				log.Printf("MCP session ended: %s", sessionID)
			},
		}),
	)

	// Register tools
	registerSessionTools(mcpServer, sessionService)
	registerMemoryTools(mcpServer, memoryService)
	registerRelationshipTools(mcpServer, relationshipService)
	registerContextTools(mcpServer, contextService)
	registerTaskTools(mcpServer, taskService)

	// Setup graceful shutdown
	setupGracefulShutdown(mcpServer, pbClient)

	// Start PocketBase server in background
	go func() {
		if err := pbClient.Start(); err != nil {
			log.Fatalf("Failed to start PocketBase: %v", err)
		}
	}()

	// Wait a moment for PocketBase to start
	time.Sleep(2 * time.Second)

	// Start MCP server
	log.Println("Starting TinyBrain v2.0 MCP Server...")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

// registerSessionTools registers session management tools
func registerSessionTools(s *server.MCPServer, service *services.SessionService) {
	s.AddTool(
		mcp.NewTool("create_session",
			mcp.WithDescription("Create a new security assessment session"),
			mcp.WithString("name", mcp.Required()),
			mcp.WithString("task_type", mcp.Required()),
			mcp.WithString("description", mcp.Optional()),
			mcp.WithString("metadata", mcp.Optional()),
		),
		service.CreateSession,
	)

	s.AddTool(
		mcp.NewTool("get_session",
			mcp.WithDescription("Get session details by ID"),
			mcp.WithString("id", mcp.Required()),
		),
		service.GetSession,
	)

	s.AddTool(
		mcp.NewTool("list_sessions",
			mcp.WithDescription("List all sessions with optional filtering"),
			mcp.WithInteger("limit", mcp.Optional()),
			mcp.WithInteger("offset", mcp.Optional()),
			mcp.WithString("status", mcp.Optional()),
			mcp.WithString("task_type", mcp.Optional()),
		),
		service.ListSessions,
	)

	s.AddTool(
		mcp.NewTool("update_session",
			mcp.WithDescription("Update session details"),
			mcp.WithString("id", mcp.Required()),
			mcp.WithString("name", mcp.Optional()),
			mcp.WithString("status", mcp.Optional()),
			mcp.WithString("description", mcp.Optional()),
			mcp.WithString("metadata", mcp.Optional()),
		),
		service.UpdateSession,
	)
}

// registerMemoryTools registers memory management tools
func registerMemoryTools(s *server.MCPServer, service *services.MemoryService) {
	s.AddTool(
		mcp.NewTool("store_memory",
			mcp.WithDescription("Store a new memory entry"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("title", mcp.Required()),
			mcp.WithString("content", mcp.Required()),
			mcp.WithString("category", mcp.Required()),
			mcp.WithInteger("priority", mcp.Optional()),
			mcp.WithNumber("confidence", mcp.Optional()),
			mcp.WithArray("tags", mcp.Optional()),
			mcp.WithString("source", mcp.Optional()),
			mcp.WithString("content_type", mcp.Optional()),
		),
		service.StoreMemory,
	)

	s.AddTool(
		mcp.NewTool("get_memory",
			mcp.WithDescription("Get memory entry by ID"),
			mcp.WithString("id", mcp.Required()),
		),
		service.GetMemory,
	)

	s.AddTool(
		mcp.NewTool("search_memories",
			mcp.WithDescription("Search memory entries with various filters"),
			mcp.WithString("session_id", mcp.Optional()),
			mcp.WithString("query", mcp.Optional()),
			mcp.WithString("category", mcp.Optional()),
			mcp.WithArray("tags", mcp.Optional()),
			mcp.WithInteger("min_priority", mcp.Optional()),
			mcp.WithNumber("min_confidence", mcp.Optional()),
			mcp.WithString("search_type", mcp.Optional()),
			mcp.WithInteger("limit", mcp.Optional()),
			mcp.WithInteger("offset", mcp.Optional()),
		),
		service.SearchMemories,
	)

	s.AddTool(
		mcp.NewTool("update_memory",
			mcp.WithDescription("Update memory entry"),
			mcp.WithString("id", mcp.Required()),
			mcp.WithString("title", mcp.Optional()),
			mcp.WithString("content", mcp.Optional()),
			mcp.WithString("category", mcp.Optional()),
			mcp.WithInteger("priority", mcp.Optional()),
			mcp.WithNumber("confidence", mcp.Optional()),
			mcp.WithArray("tags", mcp.Optional()),
			mcp.WithString("source", mcp.Optional()),
			mcp.WithString("content_type", mcp.Optional()),
		),
		service.UpdateMemory,
	)

	s.AddTool(
		mcp.NewTool("delete_memory",
			mcp.WithDescription("Delete memory entry"),
			mcp.WithString("id", mcp.Required()),
		),
		service.DeleteMemory,
	)
}

// registerRelationshipTools registers relationship management tools
func registerRelationshipTools(s *server.MCPServer, service *services.RelationshipService) {
	s.AddTool(
		mcp.NewTool("create_relationship",
			mcp.WithDescription("Create a relationship between memory entries"),
			mcp.WithString("source_id", mcp.Required()),
			mcp.WithString("target_id", mcp.Required()),
			mcp.WithString("type", mcp.Required()),
			mcp.WithNumber("strength", mcp.Optional()),
			mcp.WithString("description", mcp.Optional()),
		),
		service.CreateRelationship,
	)

	s.AddTool(
		mcp.NewTool("get_relationships",
			mcp.WithDescription("Get relationships for a memory entry"),
			mcp.WithString("memory_id", mcp.Required()),
			mcp.WithString("type", mcp.Optional()),
			mcp.WithInteger("limit", mcp.Optional()),
			mcp.WithInteger("offset", mcp.Optional()),
		),
		service.GetRelationships,
	)

	s.AddTool(
		mcp.NewTool("delete_relationship",
			mcp.WithDescription("Delete a relationship"),
			mcp.WithString("id", mcp.Required()),
		),
		service.DeleteRelationship,
	)
}

// registerContextTools registers context management tools
func registerContextTools(s *server.MCPServer, service *services.ContextService) {
	s.AddTool(
		mcp.NewTool("create_context_snapshot",
			mcp.WithDescription("Create a context snapshot"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("name", mcp.Required()),
			mcp.WithObject("context_data", mcp.Required()),
			mcp.WithString("description", mcp.Optional()),
		),
		service.CreateContextSnapshot,
	)

	s.AddTool(
		mcp.NewTool("get_context_snapshot",
			mcp.WithDescription("Get context snapshot by ID"),
			mcp.WithString("id", mcp.Required()),
		),
		service.GetContextSnapshot,
	)

	s.AddTool(
		mcp.NewTool("get_context_summary",
			mcp.WithDescription("Get context summary for a session"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithInteger("max_memories", mcp.Optional()),
		),
		service.GetContextSummary,
	)
}

// registerTaskTools registers task progress tools
func registerTaskTools(s *server.MCPServer, service *services.TaskService) {
	s.AddTool(
		mcp.NewTool("create_task_progress",
			mcp.WithDescription("Create a new task progress entry"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("task_name", mcp.Required()),
			mcp.WithString("stage", mcp.Required()),
			mcp.WithString("status", mcp.Required()),
			mcp.WithInteger("progress_percentage", mcp.Optional()),
			mcp.WithString("notes", mcp.Optional()),
		),
		service.CreateTaskProgress,
	)

	s.AddTool(
		mcp.NewTool("update_task_progress",
			mcp.WithDescription("Update task progress"),
			mcp.WithString("id", mcp.Required()),
			mcp.WithString("stage", mcp.Optional()),
			mcp.WithString("status", mcp.Optional()),
			mcp.WithInteger("progress_percentage", mcp.Optional()),
			mcp.WithString("notes", mcp.Optional()),
		),
		service.UpdateTaskProgress,
	)

	s.AddTool(
		mcp.NewTool("get_task_progress",
			mcp.WithDescription("Get task progress by ID"),
			mcp.WithString("id", mcp.Required()),
		),
		service.GetTaskProgress,
	)

	s.AddTool(
		mcp.NewTool("list_task_progress",
			mcp.WithDescription("List task progress for a session"),
			mcp.WithString("session_id", mcp.Required()),
			mcp.WithString("status", mcp.Optional()),
			mcp.WithInteger("limit", mcp.Optional()),
			mcp.WithInteger("offset", mcp.Optional()),
		),
		service.ListTaskProgress,
	)
}

// setupGracefulShutdown sets up graceful shutdown handling
func setupGracefulShutdown(mcpServer *server.MCPServer, pbClient *database.PocketBaseClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := mcpServer.Shutdown(ctx); err != nil {
			log.Printf("MCP server shutdown error: %v", err)
		}

		if err := pbClient.Close(); err != nil {
			log.Printf("PocketBase shutdown error: %v", err)
		}

		log.Println("Graceful shutdown completed")
		os.Exit(0)
	}()
}
