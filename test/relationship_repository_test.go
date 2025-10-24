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

// RelationshipRepositoryTestSuite tests the RelationshipRepositoryV2
type RelationshipRepositoryTestSuite struct {
	suite.Suite
	app         *pocketbase.PocketBase
	repo        *repository.RelationshipRepositoryV2
	sessionRepo *repository.SessionRepositoryV2
	memoryRepo  *repository.MemoryRepositoryV2
	dataDir     string
	sessionID   string
	memory1ID   string
	memory2ID   string
}

// SetupSuite initializes the test suite
func (suite *RelationshipRepositoryTestSuite) SetupSuite() {
	// Create a temporary data directory for testing
	suite.dataDir = "./test_pb_data_relationship"
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
	suite.repo = repository.NewRelationshipRepositoryV2(suite.app)
	suite.sessionRepo = repository.NewSessionRepositoryV2(suite.app)
	suite.memoryRepo = repository.NewMemoryRepositoryV2(suite.app)

	// Create a test session
	ctx := context.Background()
	sessionReq := &models.SessionCreateRequest{
		Name:     "Test Session for Relationships",
		TaskType: "penetration_test",
	}

	session, err := suite.sessionRepo.CreateSession(ctx, sessionReq)
	suite.Require().NoError(err, "Failed to create test session")
	suite.sessionID = session.ID

	// Create test memories
	memory1Req := &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "SQL Injection Vulnerability",
		Content:    "Found SQL injection in login form",
		Category:   "vulnerability",
		Priority:   8,
		Confidence: 0.9,
	}

	memory1, err := suite.memoryRepo.StoreMemory(ctx, memory1Req)
	suite.Require().NoError(err, "Failed to create test memory 1")
	suite.memory1ID = memory1.ID

	memory2Req := &models.MemoryCreateRequest{
		SessionID:  suite.sessionID,
		Title:      "Authentication Bypass",
		Content:    "Authentication can be bypassed",
		Category:   "vulnerability",
		Priority:   9,
		Confidence: 0.8,
	}

	memory2, err := suite.memoryRepo.StoreMemory(ctx, memory2Req)
	suite.Require().NoError(err, "Failed to create test memory 2")
	suite.memory2ID = memory2.ID
}

// TearDownSuite cleans up after the test suite
func (suite *RelationshipRepositoryTestSuite) TearDownSuite() {
	// Clean up test data directory
	os.RemoveAll(suite.dataDir)
}

// TestCreateRelationship tests creating a new relationship
func (suite *RelationshipRepositoryTestSuite) TestCreateRelationship() {
	ctx := context.Background()

	req := &models.RelationshipCreateRequest{
		SourceID:    suite.memory1ID,
		TargetID:    suite.memory2ID,
		Type:        models.RelationshipTypeCauses,
		Strength:    0.8,
		Description: "SQL injection leads to authentication bypass",
	}

	relationship, err := suite.repo.CreateRelationship(ctx, req)

	suite.NoError(err, "Failed to create relationship")
	suite.NotNil(relationship, "Relationship should not be nil")
	suite.Equal(req.SourceID, relationship.SourceID, "Source ID should match")
	suite.Equal(req.TargetID, relationship.TargetID, "Target ID should match")
	suite.Equal(req.Type, relationship.Type, "Type should match")
	suite.Equal(req.Strength, relationship.Strength, "Strength should match")
	suite.Equal(req.Description, relationship.Description, "Description should match")
	suite.NotEmpty(relationship.ID, "Relationship ID should not be empty")
	suite.NotZero(relationship.CreatedAt, "CreatedAt should not be zero")
	suite.NotZero(relationship.UpdatedAt, "UpdatedAt should not be zero")
}

// TestGetRelationship tests retrieving a relationship by ID
func (suite *RelationshipRepositoryTestSuite) TestGetRelationship() {
	ctx := context.Background()

	// First create a relationship
	req := &models.RelationshipCreateRequest{
		SourceID:    suite.memory1ID,
		TargetID:    suite.memory2ID,
		Type:        models.RelationshipTypeDependsOn,
		Strength:    0.7,
		Description: "Memory 1 depends on Memory 2",
	}

	createdRelationship, err := suite.repo.CreateRelationship(ctx, req)
	suite.Require().NoError(err, "Failed to create relationship for get test")

	// Now retrieve it
	retrievedRelationship, err := suite.repo.GetRelationship(ctx, createdRelationship.ID)

	suite.NoError(err, "Failed to get relationship")
	suite.NotNil(retrievedRelationship, "Retrieved relationship should not be nil")
	suite.Equal(createdRelationship.ID, retrievedRelationship.ID, "Relationship IDs should match")
	suite.Equal(createdRelationship.SourceID, retrievedRelationship.SourceID, "Source IDs should match")
	suite.Equal(createdRelationship.TargetID, retrievedRelationship.TargetID, "Target IDs should match")
	suite.Equal(createdRelationship.Type, retrievedRelationship.Type, "Types should match")
}

// TestGetRelationshipNotFound tests retrieving a non-existent relationship
func (suite *RelationshipRepositoryTestSuite) TestGetRelationshipNotFound() {
	ctx := context.Background()

	_, err := suite.repo.GetRelationship(ctx, "non-existent-id")

	suite.Error(err, "Should return error for non-existent relationship")
	suite.Contains(err.Error(), "not found", "Error should indicate relationship not found")
}

// TestListRelationships tests listing relationships with various filters
func (suite *RelationshipRepositoryTestSuite) TestListRelationships() {
	ctx := context.Background()

	// Create multiple relationships
	relationships := []*models.RelationshipCreateRequest{
		{
			SourceID:    suite.memory1ID,
			TargetID:    suite.memory2ID,
			Type:        models.RelationshipTypeCauses,
			Strength:    0.8,
			Description: "Memory 1 causes Memory 2",
		},
		{
			SourceID:    suite.memory2ID,
			TargetID:    suite.memory1ID,
			Type:        models.RelationshipTypeMitigates,
			Strength:    0.6,
			Description: "Memory 2 mitigates Memory 1",
		},
	}

	for _, req := range relationships {
		_, err := suite.repo.CreateRelationship(ctx, req)
		suite.Require().NoError(err, "Failed to create relationship for list test")
	}

	// Test listing all relationships
	listReq := &models.RelationshipListRequest{
		Limit:  10,
		Offset: 0,
	}

	relationshipList, totalCount, err := suite.repo.ListRelationships(ctx, listReq)

	suite.NoError(err, "Failed to list relationships")
	suite.GreaterOrEqual(totalCount, len(relationships), "Total count should be at least the number of created relationships")
	suite.NotEmpty(relationshipList, "Relationship list should not be empty")

	// Test filtering by source ID
	sourceReq := &models.RelationshipListRequest{
		SourceID: suite.memory1ID,
		Limit:    10,
		Offset:   0,
	}

	sourceList, _, err := suite.repo.ListRelationships(ctx, sourceReq)

	suite.NoError(err, "Failed to list relationships by source")
	suite.NotEmpty(sourceList, "Source relationship list should not be empty")

	// Verify all returned relationships have the correct source ID
	for _, relationship := range sourceList {
		suite.Equal(suite.memory1ID, relationship.SourceID, "All relationships should have the correct source ID")
	}

	// Test filtering by target ID
	targetReq := &models.RelationshipListRequest{
		TargetID: suite.memory2ID,
		Limit:    10,
		Offset:   0,
	}

	targetList, _, err := suite.repo.ListRelationships(ctx, targetReq)

	suite.NoError(err, "Failed to list relationships by target")
	suite.NotEmpty(targetList, "Target relationship list should not be empty")

	// Verify all returned relationships have the correct target ID
	for _, relationship := range targetList {
		suite.Equal(suite.memory2ID, relationship.TargetID, "All relationships should have the correct target ID")
	}

	// Test filtering by type
	typeReq := &models.RelationshipListRequest{
		Type:   string(models.RelationshipTypeCauses),
		Limit:  10,
		Offset: 0,
	}

	typeList, _, err := suite.repo.ListRelationships(ctx, typeReq)

	suite.NoError(err, "Failed to list relationships by type")
	suite.NotEmpty(typeList, "Type relationship list should not be empty")

	// Verify all returned relationships have the correct type
	for _, relationship := range typeList {
		suite.Equal(models.RelationshipTypeCauses, relationship.Type, "All relationships should have the correct type")
	}
}

// TestUpdateRelationship tests updating an existing relationship
func (suite *RelationshipRepositoryTestSuite) TestUpdateRelationship() {
	ctx := context.Background()

	// Create a relationship
	req := &models.RelationshipCreateRequest{
		SourceID:    suite.memory1ID,
		TargetID:    suite.memory2ID,
		Type:        models.RelationshipTypeRelatedTo,
		Strength:    0.5,
		Description: "Original relationship description",
	}

	createdRelationship, err := suite.repo.CreateRelationship(ctx, req)
	suite.Require().NoError(err, "Failed to create relationship for update test")

	// Update the relationship
	updateReq := &models.RelationshipUpdateRequest{
		Type:        &models.RelationshipTypeExploits,
		Strength:    float32Ptr(0.9),
		Description: stringPtr("Updated relationship description"),
	}

	updatedRelationship, err := suite.repo.UpdateRelationship(ctx, createdRelationship.ID, updateReq)

	suite.NoError(err, "Failed to update relationship")
	suite.NotNil(updatedRelationship, "Updated relationship should not be nil")
	suite.Equal(createdRelationship.ID, updatedRelationship.ID, "Relationship ID should remain the same")
	suite.Equal(models.RelationshipTypeExploits, updatedRelationship.Type, "Type should be updated")
	suite.Equal(float32(0.9), updatedRelationship.Strength, "Strength should be updated")
	suite.Equal("Updated relationship description", updatedRelationship.Description, "Description should be updated")
	suite.True(updatedRelationship.UpdatedAt.After(createdRelationship.UpdatedAt), "UpdatedAt should be newer")
}

// TestDeleteRelationship tests deleting a relationship
func (suite *RelationshipRepositoryTestSuite) TestDeleteRelationship() {
	ctx := context.Background()

	// Create a relationship
	req := &models.RelationshipCreateRequest{
		SourceID:    suite.memory1ID,
		TargetID:    suite.memory2ID,
		Type:        models.RelationshipTypeReferences,
		Strength:    0.6,
		Description: "Relationship to be deleted",
	}

	createdRelationship, err := suite.repo.CreateRelationship(ctx, req)
	suite.Require().NoError(err, "Failed to create relationship for delete test")

	// Delete the relationship
	err = suite.repo.DeleteRelationship(ctx, createdRelationship.ID)

	suite.NoError(err, "Failed to delete relationship")

	// Verify the relationship is deleted
	_, err = suite.repo.GetRelationship(ctx, createdRelationship.ID)
	suite.Error(err, "Should return error when trying to get deleted relationship")
	suite.Contains(err.Error(), "not found", "Error should indicate relationship not found")
}

// TestRelationshipValidation tests relationship validation
func (suite *RelationshipRepositoryTestSuite) TestRelationshipValidation() {
	ctx := context.Background()

	// Test empty source ID
	req := &models.RelationshipCreateRequest{
		SourceID: "",
		TargetID: suite.memory2ID,
		Type:     models.RelationshipTypeCauses,
		Strength: 0.8,
	}

	_, err := suite.repo.CreateRelationship(ctx, req)
	suite.Error(err, "Should return error for empty source ID")

	// Test empty target ID
	req = &models.RelationshipCreateRequest{
		SourceID: suite.memory1ID,
		TargetID: "",
		Type:     models.RelationshipTypeCauses,
		Strength: 0.8,
	}

	_, err = suite.repo.CreateRelationship(ctx, req)
	suite.Error(err, "Should return error for empty target ID")

	// Test same source and target
	req = &models.RelationshipCreateRequest{
		SourceID: suite.memory1ID,
		TargetID: suite.memory1ID,
		Type:     models.RelationshipTypeCauses,
		Strength: 0.8,
	}

	_, err = suite.repo.CreateRelationship(ctx, req)
	suite.Error(err, "Should return error for same source and target")

	// Test empty type
	req = &models.RelationshipCreateRequest{
		SourceID: suite.memory1ID,
		TargetID: suite.memory2ID,
		Type:     "",
		Strength: 0.8,
	}

	_, err = suite.repo.CreateRelationship(ctx, req)
	suite.Error(err, "Should return error for empty type")

	// Test invalid strength
	req = &models.RelationshipCreateRequest{
		SourceID: suite.memory1ID,
		TargetID: suite.memory2ID,
		Type:     models.RelationshipTypeCauses,
		Strength: 1.5, // Invalid: should be 0.0-1.0
	}

	_, err = suite.repo.CreateRelationship(ctx, req)
	suite.Error(err, "Should return error for invalid strength")
}

// TestRelationshipTypes tests all relationship types
func (suite *RelationshipRepositoryTestSuite) TestRelationshipTypes() {
	ctx := context.Background()

	relationshipTypes := []models.RelationshipType{
		models.RelationshipTypeDependsOn,
		models.RelationshipTypeCauses,
		models.RelationshipTypeMitigates,
		models.RelationshipTypeExploits,
		models.RelationshipTypeReferences,
		models.RelationshipTypeContradicts,
		models.RelationshipTypeSupports,
		models.RelationshipTypeRelatedTo,
		models.RelationshipTypeParentOf,
		models.RelationshipTypeChildOf,
	}

	for i, relType := range relationshipTypes {
		req := &models.RelationshipCreateRequest{
			SourceID:    suite.memory1ID,
			TargetID:    suite.memory2ID,
			Type:        relType,
			Strength:    0.5 + float32(i)*0.05, // Varying strength
			Description: "Test relationship of type " + string(relType),
		}

		relationship, err := suite.repo.CreateRelationship(ctx, req)

		suite.NoError(err, "Failed to create relationship with type %s", relType)
		suite.NotNil(relationship, "Relationship should not be nil")
		suite.Equal(relType, relationship.Type, "Relationship type should match")
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func float32Ptr(f float32) *float32 {
	return &f
}

// Run the test suite
func TestRelationshipRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RelationshipRepositoryTestSuite))
}
