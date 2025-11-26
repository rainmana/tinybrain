package test

import (
	"context"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/stretchr/testify/suite"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
	"tinybrain-v2/internal/services"
)

// ServiceIntegrationTestSuite tests the complete service layer integration
type ServiceIntegrationTestSuite struct {
	suite.Suite
	app                 *pocketbase.PocketBase
	sessionService      *services.SessionServiceV2
	memoryService       *services.MemoryServiceV2
	relationshipService *services.RelationshipServiceV2
	contextService      *services.ContextServiceV2
	taskService         *services.TaskServiceV2
	dataDir             string
	sessionID           string
}

// SetupSuite initializes the test suite
func (suite *ServiceIntegrationTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_pb_data_service"
	os.RemoveAll(suite.dataDir) // Clean up any existing test data

	// Initialize PocketBase with test config
	config := pocketbase.Config{
		DefaultDataDir: suite.dataDir,
	}
	suite.app = pocketbase.NewWithConfig(config)

	// Bootstrap the app
	err := suite.app.Bootstrap()
	suite.Require().NoError(err, "Failed to bootstrap PocketBase for testing")

	// Create repositories
	sessionRepo := repository.NewSessionRepositoryV2(suite.app)
	memoryRepo := repository.NewMemoryRepositoryV2(suite.app)
	relationshipRepo := repository.NewRelationshipRepositoryV2(suite.app)
	contextRepo := repository.NewContextRepositoryV2(suite.app)
	taskRepo := repository.NewTaskRepositoryV2(suite.app)

	// Create services
	suite.sessionService = services.NewSessionServiceV2(sessionRepo)
	suite.memoryService = services.NewMemoryServiceV2(memoryRepo)
	suite.relationshipService = services.NewRelationshipServiceV2(relationshipRepo)
	suite.contextService = services.NewContextServiceV2(contextRepo)
	suite.taskService = services.NewTaskServiceV2(taskRepo)
}

// TearDownSuite cleans up after the test suite
func (suite *ServiceIntegrationTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// TestCompleteWorkflow tests a complete security assessment workflow
func (suite *ServiceIntegrationTestSuite) TestCompleteWorkflow() {
	ctx := context.Background()

	// Step 1: Create a security assessment session
	sessionReq := &models.SessionCreateRequest{
		Name:        "Comprehensive Security Assessment",
		TaskType:    "penetration_test",
		Description: "Full security assessment of web application",
		Metadata: map[string]interface{}{
			"client":   "test-client",
			"priority": "high",
			"scope":    "web-application",
		},
	}

	session, err := suite.sessionService.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create session")
	suite.sessionID = session.ID

	// Step 2: Store security findings as memories
	findings := []*models.MemoryCreateRequest{
		{
			SessionID:   suite.sessionID,
			Title:       "SQL Injection Vulnerability",
			Content:     "Found SQL injection in login form at /login endpoint. Parameter 'username' is vulnerable to SQL injection attacks.",
			Category:    "vulnerability",
			Priority:    8,
			Confidence:  0.9,
			Tags:        []string{"sql-injection", "critical", "authentication"},
			Source:      "manual_testing",
			ContentType: "text",
		},
		{
			SessionID:   suite.sessionID,
			Title:       "XSS Vulnerability",
			Content:     "Cross-site scripting vulnerability found in user profile page. User input is not properly sanitized.",
			Category:    "vulnerability",
			Priority:    6,
			Confidence:  0.8,
			Tags:        []string{"xss", "medium", "user-input"},
			Source:      "automated_scanner",
			ContentType: "text",
		},
		{
			SessionID:   suite.sessionID,
			Title:       "Security Headers Missing",
			Content:     "Application is missing important security headers like X-Frame-Options, X-XSS-Protection, etc.",
			Category:    "configuration",
			Priority:    4,
			Confidence:  0.7,
			Tags:        []string{"headers", "configuration", "low"},
			Source:      "security_headers_check",
			ContentType: "text",
		},
	}

	var memoryIDs []string
	for _, finding := range findings {
		memory, err := suite.memoryService.StoreMemory(ctx, finding)
		suite.Require().NoError(err, "Failed to store memory")
		memoryIDs = append(memoryIDs, memory.ID)
	}

	// Step 3: Create relationships between findings
	relationships := []*models.RelationshipCreateRequest{
		{
			SourceID:    memoryIDs[0], // SQL Injection
			TargetID:    memoryIDs[1], // XSS
			Type:        models.RelationshipTypeRelatedTo,
			Strength:    0.7,
			Description: "Both vulnerabilities are related to input validation issues",
		},
		{
			SourceID:    memoryIDs[0], // SQL Injection
			TargetID:    memoryIDs[2], // Security Headers
			Type:        models.RelationshipTypeCauses,
			Strength:    0.6,
			Description: "SQL injection vulnerability is exacerbated by missing security headers",
		},
	}

	for _, rel := range relationships {
		relationship, err := suite.relationshipService.CreateRelationship(ctx, rel)
		suite.Require().NoError(err, "Failed to create relationship")
		suite.NotNil(relationship, "Relationship should not be nil")
	}

	// Step 4: Create context snapshots at different stages
	contextSnapshots := []*models.ContextSnapshotCreateRequest{
		{
			SessionID: suite.sessionID,
			Name:      "Initial Assessment Context",
			ContextData: map[string]interface{}{
				"assessment_phase":   "reconnaissance",
				"targets_identified": 3,
				"tools_used":         []string{"nmap", "nikto", "burp_suite"},
			},
			Description: "Context at the beginning of the assessment",
		},
		{
			SessionID: suite.sessionID,
			Name:      "Mid-Assessment Context",
			ContextData: map[string]interface{}{
				"assessment_phase":      "vulnerability_discovery",
				"vulnerabilities_found": 3,
				"critical_issues":       1,
				"tools_used":            []string{"burp_suite", "sqlmap", "custom_scripts"},
			},
			Description: "Context during vulnerability discovery phase",
		},
	}

	for _, snapshot := range contextSnapshots {
		contextSnapshot, err := suite.contextService.CreateContextSnapshot(ctx, snapshot)
		suite.Require().NoError(err, "Failed to create context snapshot")
		suite.NotNil(contextSnapshot, "Context snapshot should not be nil")
	}

	// Step 5: Track task progress
	taskProgress := []*models.TaskProgressCreateRequest{
		{
			SessionID:          suite.sessionID,
			TaskName:           "Reconnaissance",
			Stage:              "reconnaissance",
			Status:             "completed",
			ProgressPercentage: 100.0,
			Notes:              "Completed network scanning and service enumeration",
		},
		{
			SessionID:          suite.sessionID,
			TaskName:           "Vulnerability Assessment",
			Stage:              "vulnerability_discovery",
			Status:             "in_progress",
			ProgressPercentage: 75.0,
			Notes:              "Found 3 vulnerabilities, still testing for more",
		},
		{
			SessionID:          suite.sessionID,
			TaskName:           "Exploitation",
			Stage:              "exploitation",
			Status:             "pending",
			ProgressPercentage: 0.0,
			Notes:              "Will begin exploitation phase after completing vulnerability assessment",
		},
	}

	for _, task := range taskProgress {
		taskProgress, err := suite.taskService.CreateTaskProgress(ctx, task)
		suite.Require().NoError(err, "Failed to create task progress")
		suite.NotNil(taskProgress, "Task progress should not be nil")
	}

	// Step 6: Search and retrieve data
	// Search memories by category
	memorySearchReq := &models.MemorySearchRequest{
		SessionID: suite.sessionID,
		Category:  "vulnerability",
		Limit:     10,
		Offset:    0,
	}

	vulnerabilities, total, err := suite.memoryService.SearchMemories(ctx, memorySearchReq)
	suite.NoError(err, "Failed to search memories")
	suite.Equal(2, total, "Should find 2 vulnerabilities")
	suite.Len(vulnerabilities, 2, "Should return 2 vulnerabilities")

	// List relationships
	relationshipListReq := &models.RelationshipListRequest{
		Limit:  10,
		Offset: 0,
	}

	relationships, total, err := suite.relationshipService.ListRelationships(ctx, relationshipListReq)
	suite.NoError(err, "Failed to list relationships")
	suite.Equal(2, total, "Should find 2 relationships")
	suite.Len(relationships, 2, "Should return 2 relationships")

	// List context snapshots
	contextListReq := &models.ContextSnapshotListRequest{
		SessionID: suite.sessionID,
		Limit:     10,
		Offset:    0,
	}

	snapshots, total, err := suite.contextService.ListContextSnapshots(ctx, contextListReq)
	suite.NoError(err, "Failed to list context snapshots")
	suite.Equal(2, total, "Should find 2 context snapshots")
	suite.Len(snapshots, 2, "Should return 2 context snapshots")

	// List task progress
	taskListReq := &models.TaskProgressListRequest{
		SessionID: suite.sessionID,
		Limit:     10,
		Offset:    0,
	}

	tasks, total, err := suite.taskService.ListTaskProgress(ctx, taskListReq)
	suite.NoError(err, "Failed to list task progress")
	suite.Equal(3, total, "Should find 3 task progress entries")
	suite.Len(tasks, 3, "Should return 3 task progress entries")

	// Step 7: Update session status
	updateReq := &models.SessionUpdateRequest{
		Status: "in_progress",
		Metadata: map[string]interface{}{
			"progress":     "75%",
			"last_updated": "2025-10-24",
		},
	}

	updatedSession, err := suite.sessionService.UpdateSession(ctx, suite.sessionID, updateReq)
	suite.NoError(err, "Failed to update session")
	suite.Equal("in_progress", updatedSession.Status, "Session status should be updated")

	// Step 8: Verify data integrity
	// Verify session exists and has correct data
	retrievedSession, err := suite.sessionService.GetSession(ctx, suite.sessionID)
	suite.NoError(err, "Failed to retrieve session")
	suite.Equal("Comprehensive Security Assessment", retrievedSession.Name, "Session name should match")
	suite.Equal("penetration_test", retrievedSession.TaskType, "Task type should match")
	suite.Equal("in_progress", retrievedSession.Status, "Status should be updated")

	// Verify memories are linked to session
	for _, memoryID := range memoryIDs {
		memory, err := suite.memoryService.GetMemory(ctx, memoryID)
		suite.NoError(err, "Failed to retrieve memory")
		suite.Equal(suite.sessionID, memory.SessionID, "Memory should be linked to session")
	}

	// Verify relationships exist
	for _, relationship := range relationships {
		suite.Contains(memoryIDs, relationship.SourceID, "Source ID should be in memory IDs")
		suite.Contains(memoryIDs, relationship.TargetID, "Target ID should be in memory IDs")
	}
}

// TestServiceValidation tests service layer validation
func (suite *ServiceIntegrationTestSuite) TestServiceValidation() {
	ctx := context.Background()

	// Test session validation
	invalidSessionReq := &models.SessionCreateRequest{
		Name:     "", // Empty name should fail
		TaskType: "penetration_test",
	}

	_, err := suite.sessionService.CreateSession(ctx, invalidSessionReq)
	suite.Error(err, "Should return error for empty session name")

	// Test memory validation
	invalidMemoryReq := &models.MemoryCreateRequest{
		SessionID:  "invalid-session-id",
		Title:      "Test Memory",
		Content:    "Test content",
		Category:   "vulnerability",
		Priority:   15, // Invalid priority
		Confidence: 0.8,
	}

	_, err = suite.memoryService.StoreMemory(ctx, invalidMemoryReq)
	suite.Error(err, "Should return error for invalid memory priority")

	// Test relationship validation
	invalidRelationshipReq := &models.RelationshipCreateRequest{
		SourceID: "invalid-source",
		TargetID: "invalid-target",
		Type:     models.RelationshipTypeCauses,
		Strength: 1.5, // Invalid strength
	}

	_, err = suite.relationshipService.CreateRelationship(ctx, invalidRelationshipReq)
	suite.Error(err, "Should return error for invalid relationship strength")

	// Test context snapshot validation
	invalidContextReq := &models.ContextSnapshotCreateRequest{
		SessionID:   "invalid-session-id",
		Name:        "", // Empty name should fail
		ContextData: map[string]interface{}{"test": "data"},
	}

	_, err = suite.contextService.CreateContextSnapshot(ctx, invalidContextReq)
	suite.Error(err, "Should return error for empty context snapshot name")

	// Test task progress validation
	invalidTaskReq := &models.TaskProgressCreateRequest{
		SessionID:          "invalid-session-id",
		TaskName:           "Test Task",
		Stage:              "test",
		Status:             "pending",
		ProgressPercentage: 150.0, // Invalid progress percentage
	}

	_, err = suite.taskService.CreateTaskProgress(ctx, invalidTaskReq)
	suite.Error(err, "Should return error for invalid progress percentage")
}

// TestServiceErrorHandling tests error handling in services
func (suite *ServiceIntegrationTestSuite) TestServiceErrorHandling() {
	ctx := context.Background()

	// Test getting non-existent session
	_, err := suite.sessionService.GetSession(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent session")

	// Test getting non-existent memory
	_, err = suite.memoryService.GetMemory(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent memory")

	// Test getting non-existent relationship
	_, err = suite.relationshipService.GetRelationship(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent relationship")

	// Test getting non-existent context snapshot
	_, err = suite.contextService.GetContextSnapshot(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent context snapshot")

	// Test getting non-existent task progress
	_, err = suite.taskService.GetTaskProgress(ctx, "non-existent-id")
	suite.Error(err, "Should return error for non-existent task progress")
}

// Run the test suite
func TestServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationTestSuite))
}
