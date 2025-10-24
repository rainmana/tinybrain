package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tinybrain-v2/internal/database"
	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

func TestTinyBrainV2Integration(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()

	// Initialize PocketBase client
	pbClient, err := database.NewPocketBaseClient(tempDir)
	require.NoError(t, err)
	defer pbClient.Close()

	// Bootstrap database
	ctx := context.Background()
	err = pbClient.Bootstrap(ctx)
	require.NoError(t, err)

	// Get PocketBase app
	app := pbClient.GetApp()

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(app)
	memoryRepo := repository.NewMemoryRepository(app)
	relationshipRepo := repository.NewRelationshipRepository(app)
	contextRepo := repository.NewContextRepository(app)
	taskRepo := repository.NewTaskRepository(app)

	// Initialize services
	sessionService := services.NewSessionService(sessionRepo)
	memoryService := services.NewMemoryService(memoryRepo)
	relationshipService := services.NewRelationshipService(relationshipRepo)
	contextService := services.NewContextService(contextRepo)
	taskService := services.NewTaskService(taskRepo)

	// Test 1: Create a session
	t.Run("CreateSession", func(t *testing.T) {
		req := &models.SessionCreateRequest{
			Name:        "Security Assessment Test",
			TaskType:    "security_review",
			Description: "Testing TinyBrain v2.0 session creation",
			Metadata:    map[string]interface{}{"priority": "high"},
		}

		session, err := sessionService.CreateSession(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, req.Name, session.Name)
		assert.Equal(t, req.TaskType, session.TaskType)
		assert.Equal(t, "active", session.Status)

		t.Logf("✅ Created session: %s", session.ID)
	})

	// Test 2: Store a memory
	t.Run("StoreMemory", func(t *testing.T) {
		// First create a session
		sessionReq := &models.SessionCreateRequest{
			Name:     "Memory Test Session",
			TaskType: "vulnerability_analysis",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err)

		// Store a memory
		memoryReq := &models.MemoryCreateRequest{
			SessionID:   session.ID,
			Title:       "SQL Injection Vulnerability",
			Content:     "Found SQL injection in login form parameter 'username'",
			Category:    "vulnerability",
			Priority:    8,
			Confidence:  0.9,
			Tags:        []string{"sql-injection", "authentication", "critical"},
			Source:      "manual_testing",
			ContentType: "text",
		}

		memory, err := memoryService.StoreMemory(ctx, memoryReq)
		require.NoError(t, err)
		assert.NotEmpty(t, memory.ID)
		assert.Equal(t, memoryReq.Title, memory.Title)
		assert.Equal(t, memoryReq.Content, memory.Content)
		assert.Equal(t, memoryReq.Category, memory.Category)
		assert.Equal(t, memoryReq.Priority, memory.Priority)
		assert.Equal(t, memoryReq.Confidence, memory.Confidence)

		t.Logf("✅ Stored memory: %s", memory.ID)
	})

	// Test 3: Create a relationship
	t.Run("CreateRelationship", func(t *testing.T) {
		// Create two memories first
		sessionReq := &models.SessionCreateRequest{
			Name:     "Relationship Test Session",
			TaskType: "threat_modeling",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err)

		// Create first memory
		memory1Req := &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Authentication Bypass",
			Content:    "Weak authentication mechanism allows bypass",
			Category:   "vulnerability",
			Priority:   9,
			Confidence: 0.95,
		}
		memory1, err := memoryService.StoreMemory(ctx, memory1Req)
		require.NoError(t, err)

		// Create second memory
		memory2Req := &models.MemoryCreateRequest{
			SessionID:  session.ID,
			Title:      "Privilege Escalation",
			Content:    "User can escalate privileges after authentication bypass",
			Category:   "vulnerability",
			Priority:   8,
			Confidence: 0.85,
		}
		memory2, err := memoryService.StoreMemory(ctx, memory2Req)
		require.NoError(t, err)

		// Create relationship
		relReq := &models.RelationshipCreateRequest{
			SourceID:    memory1.ID,
			TargetID:    memory2.ID,
			Type:        "causes",
			Strength:    0.9,
			Description: "Authentication bypass leads to privilege escalation",
		}

		relationship, err := relationshipService.CreateRelationship(ctx, relReq)
		require.NoError(t, err)
		assert.NotEmpty(t, relationship.ID)
		assert.Equal(t, relReq.SourceID, relationship.SourceID)
		assert.Equal(t, relReq.TargetID, relationship.TargetID)
		assert.Equal(t, relReq.Type, relationship.Type)
		assert.Equal(t, relReq.Strength, relationship.Strength)

		t.Logf("✅ Created relationship: %s", relationship.ID)
	})

	// Test 4: Create context snapshot
	t.Run("CreateContextSnapshot", func(t *testing.T) {
		// Create a session first
		sessionReq := &models.SessionCreateRequest{
			Name:     "Context Test Session",
			TaskType: "incident_response",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err)

		// Create context snapshot
		contextReq := &models.ContextSnapshotCreateRequest{
			SessionID:   session.ID,
			Name:        "Initial Assessment Context",
			Description: "Context at the start of incident response",
			ContextData: map[string]interface{}{
				"threat_level":     "high",
				"affected_systems": []string{"web-server", "database"},
				"timeline":         "2024-01-15T10:30:00Z",
			},
		}

		snapshot, err := contextService.CreateContextSnapshot(ctx, contextReq)
		require.NoError(t, err)
		assert.NotEmpty(t, snapshot.ID)
		assert.Equal(t, contextReq.SessionID, snapshot.SessionID)
		assert.Equal(t, contextReq.Name, snapshot.Name)
		assert.Equal(t, contextReq.Description, snapshot.Description)
		assert.Equal(t, contextReq.ContextData, snapshot.ContextData)

		t.Logf("✅ Created context snapshot: %s", snapshot.ID)
	})

	// Test 5: Create task progress
	t.Run("CreateTaskProgress", func(t *testing.T) {
		// Create a session first
		sessionReq := &models.SessionCreateRequest{
			Name:     "Task Progress Test Session",
			TaskType: "penetration_test",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err)

		// Create task progress
		taskReq := &models.TaskProgressCreateRequest{
			SessionID:          session.ID,
			TaskName:           "Network Reconnaissance",
			Stage:              "reconnaissance",
			Status:             "in_progress",
			ProgressPercentage: 25,
			Notes:              "Completed port scanning, starting service enumeration",
		}

		task, err := taskService.CreateTaskProgress(ctx, taskReq)
		require.NoError(t, err)
		assert.NotEmpty(t, task.ID)
		assert.Equal(t, taskReq.SessionID, task.SessionID)
		assert.Equal(t, taskReq.TaskName, task.TaskName)
		assert.Equal(t, taskReq.Stage, task.Stage)
		assert.Equal(t, taskReq.Status, task.Status)
		assert.Equal(t, taskReq.ProgressPercentage, task.ProgressPercentage)
		assert.Equal(t, taskReq.Notes, task.Notes)

		t.Logf("✅ Created task progress: %s", task.ID)
	})

	// Test 6: Search memories
	t.Run("SearchMemories", func(t *testing.T) {
		// Create a session and some memories
		sessionReq := &models.SessionCreateRequest{
			Name:     "Search Test Session",
			TaskType: "code_review",
		}
		session, err := sessionService.CreateSession(ctx, sessionReq)
		require.NoError(t, err)

		// Store multiple memories
		memories := []*models.MemoryCreateRequest{
			{
				SessionID:  session.ID,
				Title:      "XSS Vulnerability",
				Content:    "Cross-site scripting found in user input field",
				Category:   "vulnerability",
				Priority:   7,
				Confidence: 0.8,
				Tags:       []string{"xss", "input-validation"},
			},
			{
				SessionID:  session.ID,
				Title:      "CSRF Protection",
				Content:    "Missing CSRF tokens in forms",
				Category:   "vulnerability",
				Priority:   6,
				Confidence: 0.7,
				Tags:       []string{"csrf", "authentication"},
			},
		}

		for _, memReq := range memories {
			_, err := memoryService.StoreMemory(ctx, memReq)
			require.NoError(t, err)
		}

		// Search for vulnerabilities
		searchReq := &models.MemorySearchRequest{
			SessionID:   session.ID,
			Category:    "vulnerability",
			MinPriority: 6,
			Limit:       10,
		}

		results, totalCount, err := memoryService.SearchMemories(ctx, searchReq)
		require.NoError(t, err)
		assert.Equal(t, 2, totalCount)
		assert.Len(t, results, 2)

		// Verify all results are vulnerabilities with priority >= 6
		for _, result := range results {
			assert.Equal(t, "vulnerability", result.Category)
			assert.GreaterOrEqual(t, result.Priority, 6)
		}

		t.Logf("✅ Found %d memories matching search criteria", totalCount)
	})

	t.Log("🎉 All TinyBrain v2.0 integration tests passed!")
}
