package services

import (
	"context"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecurityRepository is a mock implementation of SecurityRepositoryInterface
type MockSecurityRepository struct {
	mock.Mock
}

func (m *MockSecurityRepository) StoreNVDDataset(ctx context.Context, cves []models.NVDCVE) error {
	args := m.Called(ctx, cves)
	return args.Error(0)
}

func (m *MockSecurityRepository) StoreATTACKDataset(ctx context.Context, techniques []models.ATTACKTechnique, tactics []models.ATTACKTactic) error {
	args := m.Called(ctx, techniques, tactics)
	return args.Error(0)
}

func (m *MockSecurityRepository) QueryNVD(ctx context.Context, req models.NVDSearchRequest) ([]models.NVDCVE, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.NVDCVE), args.Int(1), args.Error(2)
}

func (m *MockSecurityRepository) QueryATTACK(ctx context.Context, req models.ATTACKSearchRequest) ([]models.ATTACKTechnique, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.ATTACKTechnique), args.Int(1), args.Error(2)
}

func (m *MockSecurityRepository) GetSecurityDataSummary(ctx context.Context) (map[string]models.SecurityDataSummary, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]models.SecurityDataSummary), args.Error(1)
}

func (m *MockSecurityRepository) UpdateSecurityDataStatus(ctx context.Context, dataSource string, status string, totalRecords *int, errorMessage *string) error {
	args := m.Called(ctx, dataSource, status, totalRecords, errorMessage)
	return args.Error(0)
}

func TestSecurityRetrievalService_QuerySecurityData_NVD(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	ctx := context.Background()

	// Create test CVE data
	testCVEs := []models.NVDCVE{
		{
			ID:               "CVE-2024-1234",
			Description:      "SQL injection vulnerability in Apache HTTP Server",
			CVSSV3Score:      floatPtr(9.8),
			Severity:         stringPtr("CRITICAL"),
			CWEIDs:           []string{"CWE-89"},
			AffectedProducts: []string{"apache:http_server"},
			References:       []string{"https://example.com/ref1"},
		},
	}

	// Set up mock expectations
	mockRepo.On("QueryNVD", ctx, mock.AnythingOfType("models.NVDSearchRequest")).Return(testCVEs, 1, nil)

	// Create request
	req := models.SecurityQueryRequest{
		Query:      "SQL injection",
		DataSource: "nvd",
		Filters: map[string]interface{}{
			"cwe_id": "CWE-89",
		},
		Limit: 10,
	}

	// Execute query
	response, err := service.QuerySecurityData(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "nvd", response.DataSource)
	assert.Equal(t, 1, response.TotalCount)
	assert.Len(t, response.Results, 1)
	assert.True(t, response.HasMore)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestSecurityRetrievalService_QuerySecurityData_ATTACK(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	ctx := context.Background()

	// Create test technique data
	testTechniques := []models.ATTACKTechnique{
		{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Adversaries may inject code into processes in order to evade process-based defenses",
			Tactic:      "defense-evasion",
			Tactics:     []string{"defense-evasion"},
			Platforms:   []string{"windows"},
			DataSources: []string{"process-monitoring"},
			References:  []string{"https://attack.mitre.org/techniques/T1055/001/"},
		},
	}

	// Set up mock expectations
	mockRepo.On("QueryATTACK", ctx, mock.AnythingOfType("models.ATTACKSearchRequest")).Return(testTechniques, 1, nil)

	// Create request
	req := models.SecurityQueryRequest{
		Query:      "process injection",
		DataSource: "attack",
		Filters: map[string]interface{}{
			"tactic": "defense-evasion",
		},
		Limit: 10,
	}

	// Execute query
	response, err := service.QuerySecurityData(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "attack", response.DataSource)
	assert.Equal(t, 1, response.TotalCount)
	assert.Len(t, response.Results, 1)
	assert.True(t, response.HasMore)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestSecurityRetrievalService_QuerySecurityData_OWASP(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	ctx := context.Background()

	// Create request for OWASP (not yet implemented)
	req := models.SecurityQueryRequest{
		Query:      "authentication testing",
		DataSource: "owasp",
		Filters: map[string]interface{}{
			"category": "authentication",
		},
		Limit: 10,
	}

	// Execute query
	_, err := service.QuerySecurityData(ctx, req)

	// Should return error since OWASP is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}

func TestSecurityRetrievalService_QuerySecurityData_UnsupportedDataSource(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	ctx := context.Background()

	// Create request with unsupported data source
	req := models.SecurityQueryRequest{
		Query:      "test query",
		DataSource: "unsupported",
		Limit:      10,
	}

	// Execute query
	response, err := service.QuerySecurityData(ctx, req)

	// Should return error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported data source")
	assert.Nil(t, response)
}

func TestSecurityRetrievalService_SummarizeNVDResults(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	// Create test CVE data
	cves := []interface{}{
		models.NVDCVE{
			ID:               "CVE-2024-1234",
			Description:      "SQL injection vulnerability in Apache HTTP Server that allows remote attackers to execute arbitrary SQL commands via the username parameter in the login form. This vulnerability affects Apache HTTP Server versions 2.4.0 through 2.4.41 and can be exploited by sending specially crafted requests to the authentication endpoint.",
			CVSSV3Score:      floatPtr(9.8),
			Severity:         stringPtr("CRITICAL"),
			CWEIDs:           []string{"CWE-89"},
			AffectedProducts: []string{"apache:http_server:2.4.0", "apache:http_server:2.4.1", "apache:http_server:2.4.2"},
			References:       []string{"https://example.com/ref1", "https://example.com/ref2"},
		},
		models.NVDCVE{
			ID:               "CVE-2024-1235",
			Description:      "Cross-site scripting vulnerability in Nginx web server",
			CVSSV3Score:      floatPtr(7.5),
			Severity:         stringPtr("HIGH"),
			CWEIDs:           []string{"CWE-79"},
			AffectedProducts: []string{"nginx:1.18.0"},
			References:       []string{"https://example.com/ref3"},
		},
	}

	// Test summarization
	summaries := service.summarizeNVDResults(cves)

	// Assertions
	assert.Len(t, summaries, 2)

	// Check first summary
	summary1 := summaries[0].(map[string]interface{})
	assert.Equal(t, "CVE-2024-1234", summary1["id"])
	assert.Equal(t, "CRITICAL", summary1["severity"])
	assert.Equal(t, 9.8, summary1["cvss_v3"])
	assert.Contains(t, summary1["description"], "SQL injection vulnerability")
	assert.Len(t, summary1["affected_products"], 3) // Should be truncated
	assert.NotEmpty(t, summary1["summary"])

	// Check second summary
	summary2 := summaries[1].(map[string]interface{})
	assert.Equal(t, "CVE-2024-1235", summary2["id"])
	assert.Equal(t, "HIGH", summary2["severity"])
	assert.Equal(t, 7.5, summary2["cvss_v3"])
}

func TestSecurityRetrievalService_SummarizeATTACKResults(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	// Create test technique data
	techniques := []interface{}{
		models.ATTACKTechnique{
			ID:          "T1055.001",
			Name:        "Process Injection",
			Description: "Adversaries may inject code into processes in order to evade process-based defenses as well as possibly elevate privileges. Process injection is a method of executing arbitrary code in the address space of a separate live process. Running code in the context of another process may allow access to the process's memory, system/network resources, and possibly elevated privileges.",
			Tactic:      "defense-evasion",
			Tactics:     []string{"defense-evasion", "persistence"},
			Platforms:   []string{"windows", "linux"},
			DataSources: []string{"process-monitoring", "file-monitoring"},
			References:  []string{"https://attack.mitre.org/techniques/T1055/001/"},
		},
	}

	// Test summarization
	summaries := service.summarizeATTACKResults(techniques)

	// Assertions
	assert.Len(t, summaries, 1)

	summary := summaries[0].(map[string]interface{})
	assert.Equal(t, "T1055.001", summary["id"])
	assert.Equal(t, "Process Injection", summary["name"])
	assert.Equal(t, "defense-evasion", summary["tactic"])
	assert.Contains(t, summary["description"], "Adversaries may inject code")
	assert.Contains(t, summary["platforms"], "windows")
	assert.Contains(t, summary["platforms"], "linux")
	assert.NotEmpty(t, summary["summary"])
}

func TestSecurityRetrievalService_GenerateCVESummary(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	// Test CVE with all fields
	cve := models.NVDCVE{
		ID:               "CVE-2024-1234",
		Description:      "SQL injection vulnerability",
		CVSSV3Score:      floatPtr(9.8),
		Severity:         stringPtr("CRITICAL"),
		AffectedProducts: []string{"apache:http_server:2.4.41", "apache:http_server:2.4.40"},
	}

	summary := service.generateCVESummary(cve)
	assert.Contains(t, summary, "CVE-2024-1234")
	assert.Contains(t, summary, "CRITICAL")
	assert.Contains(t, summary, "9.8")
	assert.Contains(t, summary, "apache:http_server")

	// Test CVE with minimal fields
	cve2 := models.NVDCVE{
		ID:          "CVE-2024-1235",
		Description: "Test vulnerability",
	}

	summary2 := service.generateCVESummary(cve2)
	assert.Contains(t, summary2, "CVE-2024-1235")
	assert.Contains(t, summary2, "Unknown")
}

func TestSecurityRetrievalService_GenerateTechniqueSummary(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockSecurityRepository)
	logger := log.New(os.Stderr)
	service := NewSecurityRetrievalService(mockRepo, logger)

	// Test technique with all fields
	technique := models.ATTACKTechnique{
		ID:          "T1055.001",
		Name:        "Process Injection",
		Description: "Test technique",
		Tactic:      "defense-evasion",
		Platforms:   []string{"windows", "linux"},
	}

	summary := service.generateTechniqueSummary(technique)
	assert.Contains(t, summary, "T1055.001")
	assert.Contains(t, summary, "Process Injection")
	assert.Contains(t, summary, "defense-evasion")
	assert.Contains(t, summary, "windows")
	assert.Contains(t, summary, "linux")

	// Test technique with minimal fields
	technique2 := models.ATTACKTechnique{
		ID:          "T1001",
		Name:        "Data Obfuscation",
		Description: "Test technique",
		Tactic:      "defense-evasion",
	}

	summary2 := service.generateTechniqueSummary(technique2)
	assert.Contains(t, summary2, "T1001")
	assert.Contains(t, summary2, "Data Obfuscation")
	assert.Contains(t, summary2, "Unknown")
}

func TestSecurityRetrievalService_ExtractCWEID(t *testing.T) {
	// Test CWE ID extraction
	cweID := extractCWEID("CWE-89 SQL injection")
	assert.Equal(t, "CWE-89", cweID)

	cweID = extractCWEID("CWE-79 Cross-site scripting")
	assert.Equal(t, "CWE-79", cweID)

	// Test with no CWE ID
	cweID = extractCWEID("No CWE ID here")
	assert.Equal(t, "", cweID)
}

func TestSecurityRetrievalService_ExtractTechniqueID(t *testing.T) {
	// Test technique ID extraction
	techniqueID := extractTechniqueID("T1055.001")
	assert.Equal(t, "T1055.001", techniqueID)

	techniqueID = extractTechniqueID("T1001")
	assert.Equal(t, "T1001", techniqueID)

	// Test with no technique ID
	techniqueID = extractTechniqueID("No technique ID here")
	assert.Equal(t, "", techniqueID)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
