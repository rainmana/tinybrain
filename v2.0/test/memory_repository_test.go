package test

import (
	"context"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/stretchr/testify/suite"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// MemoryRepositoryTestSuite tests the MemoryRepositoryV2
type MemoryRepositoryTestSuite struct {
	suite.Suite
	app         *pocketbase.PocketBase
	repo        *repository.MemoryRepositoryV2
	sessionRepo *repository.SessionRepositoryV2
	dataDir     string
	sessionID   string
}

// SetupSuite initializes the test suite
func (suite *MemoryRepositoryTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_pb_data_memory"
	os.RemoveAll(suite.dataDir) // Clean up any existing test data

	// Initialize PocketBase with test config
	config := pocketbase.Config{
		DefaultDataDir: suite.dataDir,
	}
	suite.app = pocketbase.NewWithConfig(config)

	// Bootstrap the app
	err := suite.app.Bootstrap()
	suite.Require().NoError(err, "Failed to bootstrap PocketBase for testing")

	// Create the repositories
	suite.repo = repository.NewMemoryRepositoryV2(suite.app)
	suite.sessionRepo = repository.NewSessionRepositoryV2(suite.app)

	// Create a test session
	ctx := context.Background()
	sessionReq := &models.SessionCreateRequest{
		Name:     "Test Session for Memory",
		TaskType: "penetration_test",
	}

	session, err := suite.sessionRepo.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create test session")
	suite.sessionID = session.ID
}

// TearDownSuite cleans up after the test suite
func (suite *MemoryRepositoryTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// TestStoreMemory tests storing a new memory
func (suite *MemoryRepositoryTestSuite) TestStoreMemory() {
	ctx := context.Background()

	req := &models.MemoryCreateRequest{
		SessionID:   suite.sessionID,
		Title:       "SQL Injection Vulnerability",
		Content:     "Found SQL injection in login form at /login endpoint",
		Category:    "vulnerability",
		Priority:    8,
		Confidence:  0.9,
		Tags:        []string{"sql-injection", "critical", "authentication"},
		Source:      "manual_testing",
		ContentType: "text",
	}

	memory, err := suite.repo.StoreMemory(ctx, req)

	suite.NoError(err, "Failed to store memory")
	suite.NotNil(memory, "Memory should not be nil")
	suite.Equal(req.SessionID, memory.SessionID, "Session ID should match")
	suite.Equal(req.Title, memory.Title, "Title should match")
	suite.Equal(req.Content, memory.Content, "Content should match")
	suite.Equal(req.Category, memory.Category, "Category should match")
	suite.Equal(req.Priority, memory.Priority, "Priority should match")
	suite.Equal(req.Confidence, memory.Confidence, "Confidence should match")
	suite.Equal(req.Source, memory.Source, "Source should match")
	suite.Equal(req.ContentType, memory.ContentType, "Content type should match")
	suite.NotEmpty(memory.ID, "Memory ID should not be empty")
	suite.NotZero(memory.CreatedAt, "CreatedAt should not be zero")
	suite.NotZero(memory.UpdatedAt, "UpdatedAt should not be zero")

	// Verify tags
	suite.Equal(len(req.Tags), len(memory.Tags), "Number of tags should match")
	for i, tag := range req.Tags {
		suite.Equal(tag, memory.Tags[i], "Tag should match")
	}
}

// TestGetMemory tests retrieving a memory by ID
func (suite *MemoryRepositoryTestSuite) TestGetMemory() {
	ctx := context.Background()

	// First create a memory
	req := &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Test Memory for Get",
		Content:    "Test content for retrieval",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	createdMemory, err := suite.repo.StoreMemory(ctx, req)
	suite.Require().NoError(err, "Failed to create memory for get test")

	// Now retrieve it
	retrievedMemory, err := suite.repo.GetMemory(ctx, createdMemory.ID)

	suite.NoError(err, "Failed to get memory")
	suite.NotNil(retrievedMemory, "Retrieved memory should not be nil")
	suite.Equal(createdMemory.ID, retrievedMemory.ID, "Memory IDs should match")
	suite.Equal(createdMemory.Title, retrievedMemory.Title, "Memory titles should match")
	suite.Equal(createdMemory.Content, retrievedMemory.Content, "Memory content should match")
}

// TestGetMemoryNotFound tests retrieving a non-existent memory
func (suite *MemoryRepositoryTestSuite) TestGetMemoryNotFound() {
	ctx := context.Background()

	_, err := suite.repo.GetMemory(ctx, "non-existent-id")

	suite.Error(err, "Should return error for non-existent memory")
	suite.Contains(err.Error(), "not found", "Error should indicate memory not found")
}

// TestSearchMemories tests searching memories with various filters
func (suite *MemoryRepositoryTestSuite) TestSearchMemories() {
	ctx := context.Background()

	// Create multiple memories with different categories and tags
	memories := []*models.MemoryCreateRequest{
		{
			SessionID:  suite.sessionID,
			Title:      "SQL Injection",
			Content:    "SQL injection vulnerability found",
			Category:   "vulnerability",
			Priority:   8,
			Confidence: 0.9,
			Tags:       []string{"sql-injection", "critical"},
			Source:     "scanner",
		},
		{
			SessionID:  suite.sessionID,
			Title:      "XSS Vulnerability",
			Content:    "Cross-site scripting vulnerability found",
			Category:   "vulnerability",
			Priority:   6,
			Confidence: 0.8,
			Tags:       []string{"xss", "medium"},
			Source:     "manual",
		},
		{
			SessionID:  suite.sessionID,
			Title:      "Security Configuration",
			Content:    "Security headers properly configured",
			Category:   "finding",
			Priority:   3,
			Confidence: 0.7,
			Tags:       []string{"configuration", "positive"},
			Source:     "review",
		},
	}

	for _, req := range memories {
		_, err := suite.repo.StoreMemory(ctx, req)
		suite.Require().NoError(err, "Failed to create memory for search test")
	}

	// Test searching all memories
	searchReq := &models.MemorySearchRequest{
		SessionID: suite.sessionID,
		Limit:     10,
		Offset:    0,
	}

	memoryList, totalCount, err := suite.repo.SearchMemories(ctx, searchReq)

	suite.NoError(err, "Failed to search memories")
	suite.GreaterOrEqual(totalCount, len(memories), "Total count should be at least the number of created memories")
	suite.NotEmpty(memoryList, "Memory list should not be empty")

	// Test filtering by category
	categoryReq := &models.MemorySearchRequest{
		SessionID: suite.sessionID,
		Category:  "vulnerability",
		Limit:     10,
		Offset:    0,
	}

	vulnerabilityList, _, err := suite.repo.SearchMemories(ctx, categoryReq)

	suite.NoError(err, "Failed to search memories by category")
	suite.NotEmpty(vulnerabilityList, "Vulnerability list should not be empty")

	// Verify all returned memories have the correct category
	for _, memory := range vulnerabilityList {
		suite.Equal("vulnerability", memory.Category, "All memories should have vulnerability category")
	}

	// Test filtering by tags
	tagReq := &models.MemorySearchRequest{
		SessionID: suite.sessionID,
		Tags:      []string{"critical"},
		Limit:     10,
		Offset:    0,
	}

	criticalList, _, err := suite.repo.SearchMemories(ctx, tagReq)

	suite.NoError(err, "Failed to search memories by tags")
	suite.NotEmpty(criticalList, "Critical memories list should not be empty")

	// Test filtering by source
	sourceReq := &models.MemorySearchRequest{
		SessionID: suite.sessionID,
		Source:    "scanner",
		Limit:     10,
		Offset:    0,
	}

	scannerList, _, err := suite.repo.SearchMemories(ctx, sourceReq)

	suite.NoError(err, "Failed to search memories by source")
	suite.NotEmpty(scannerList, "Scanner memories list should not be empty")

	// Verify all returned memories have the correct source
	for _, memory := range scannerList {
		suite.Equal("scanner", memory.Source, "All memories should have scanner source")
	}
}

// TestUpdateMemory tests updating an existing memory
func (suite *MemoryRepositoryTestSuite) TestUpdateMemory() {
	ctx := context.Background()

	// Create a memory
	req := &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Original Memory",
		Content:    "Original content",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	createdMemory, err := suite.repo.StoreMemory(ctx, req)
	suite.Require().NoError(err, "Failed to create memory for update test")

	// Update the memory
	updateReq := &models.MemoryUpdateRequest{
		Title:      stringPtr("Updated Memory"),
		Content:    stringPtr("Updated content with more details"),
		Category:   stringPtr("vulnerability"),
		Priority:   intPtr(8),
		Confidence: float32Ptr(0.9),
		Tags:       []string{"updated", "critical"},
		Source:     stringPtr("manual_review"),
	}

	updatedMemory, err := suite.repo.UpdateMemory(ctx, createdMemory.ID, updateReq)

	suite.NoError(err, "Failed to update memory")
	suite.NotNil(updatedMemory, "Updated memory should not be nil")
	suite.Equal(createdMemory.ID, updatedMemory.ID, "Memory ID should remain the same")
	suite.Equal("Updated Memory", updatedMemory.Title, "Title should be updated")
	suite.Equal("Updated content with more details", updatedMemory.Content, "Content should be updated")
	suite.Equal("vulnerability", updatedMemory.Category, "Category should be updated")
	suite.Equal(8, updatedMemory.Priority, "Priority should be updated")
	suite.Equal(float32(0.9), updatedMemory.Confidence, "Confidence should be updated")
	suite.Equal("manual_review", updatedMemory.Source, "Source should be updated")
	suite.True(updatedMemory.UpdatedAt.After(createdMemory.UpdatedAt), "UpdatedAt should be newer")
}

// TestDeleteMemory tests deleting a memory
func (suite *MemoryRepositoryTestSuite) TestDeleteMemory() {
	ctx := context.Background()

	// Create a memory
	req := &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Memory to Delete",
		Content:    "This memory will be deleted",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	createdMemory, err := suite.repo.StoreMemory(ctx, req)
	suite.Require().NoError(err, "Failed to create memory for delete test")

	// Delete the memory
	err = suite.repo.DeleteMemory(ctx, createdMemory.ID)

	suite.NoError(err, "Failed to delete memory")

	// Verify the memory is deleted
	_, err = suite.repo.GetMemory(ctx, createdMemory.ID)
	suite.Error(err, "Should return error when trying to get deleted memory")
	suite.Contains(err.Error(), "not found", "Error should indicate memory not found")
}

// TestMemoryValidation tests memory validation
func (suite *MemoryRepositoryTestSuite) TestMemoryValidation() {
	ctx := context.Background()

	// Test empty session ID
	req := &models.MemoryCreateRequest{
		SessionID:  "",
		Title:      "Valid Title",
		Content:    "Valid content",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	_, err := suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for empty session ID")

	// Test empty title
	req = &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "",
		Content:    "Valid content",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	_, err = suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for empty title")

	// Test empty content
	req = &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Valid Title",
		Content:    "",
		Category:   "finding",
		Priority:   5,
		Confidence: 0.7,
	}

	_, err = suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for empty content")

	// Test empty category
	req = &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Valid Title",
		Content:    "Valid content",
		Category:   "",
		Priority:   5,
		Confidence: 0.7,
	}

	_, err = suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for empty category")

	// Test invalid priority
	req = &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Valid Title",
		Content:    "Valid content",
		Category:   "finding",
		Priority:   15, // Invalid: should be 1-10
		Confidence: 0.7,
	}

	_, err = suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for invalid priority")

	// Test invalid confidence
	req = &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Valid Title",
		Content:    "Valid content",
		Category:   "finding",
		Priority:   5,
		Confidence: 1.5, // Invalid: should be 0.0-1.0
	}

	_, err = suite.repo.StoreMemory(ctx, req)
	suite.Error(err, "Should return error for invalid confidence")
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}

// Run the test suite
func TestMemoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryRepositoryTestSuite))
}
