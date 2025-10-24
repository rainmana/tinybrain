package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/rainmana/tinybrain/internal/repository"
)

// SecurityRetrievalService handles intelligent retrieval and summarization of security data
type SecurityRetrievalService struct {
	securityRepo repository.SecurityRepositoryInterface
	logger       *log.Logger
}

// NewSecurityRetrievalService creates a new security retrieval service
func NewSecurityRetrievalService(securityRepo repository.SecurityRepositoryInterface, logger *log.Logger) *SecurityRetrievalService {
	return &SecurityRetrievalService{
		securityRepo: securityRepo,
		logger:       logger,
	}
}

// QuerySecurityData performs intelligent querying across security datasets
func (s *SecurityRetrievalService) QuerySecurityData(ctx context.Context, req models.SecurityQueryRequest) (*models.SecurityQueryResponse, error) {
	s.logger.Info("Querying security data", "data_source", req.DataSource, "query", req.Query)

	var results []interface{}
	var totalCount int
	var err error

	switch req.DataSource {
	case "nvd":
		results, totalCount, err = s.queryNVDData(ctx, req)
	case "attack":
		results, totalCount, err = s.queryATTACKData(ctx, req)
	case "owasp":
		results, totalCount, err = s.queryOWASPData(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported data source: %s", req.DataSource)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query %s data: %v", req.DataSource, err)
	}

	// Apply smart summarization for context efficiency
	summarizedResults := s.summarizeResults(results, req.DataSource, req.Limit)

	response := &models.SecurityQueryResponse{
		Results:    summarizedResults,
		TotalCount: totalCount,
		DataSource: req.DataSource,
		Query:      req.Query,
		Filters:    req.Filters,
		Limit:      req.Limit,
		Offset:     req.Offset,
		HasMore:    req.Offset+len(summarizedResults) < totalCount,
	}

	s.logger.Info("Security data query completed", "results", len(summarizedResults), "total", totalCount)
	return response, nil
}

// queryNVDData queries NVD data with intelligent filtering
func (s *SecurityRetrievalService) queryNVDData(ctx context.Context, req models.SecurityQueryRequest) ([]interface{}, int, error) {
	// Convert generic request to NVD-specific request
	nvdReq := models.NVDSearchRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	// Extract filters
	if filters, ok := req.Filters["cwe_id"].(string); ok {
		nvdReq.CWEID = &filters
	}
	if filters, ok := req.Filters["component"].(string); ok {
		nvdReq.Component = &filters
	}
	if filters, ok := req.Filters["severity"].(string); ok {
		nvdReq.Severity = &filters
	}
	if filters, ok := req.Filters["min_cvss"].(float64); ok {
		nvdReq.MinCVSS = &filters
	}

	// If no specific filters, use query for text search
	if nvdReq.CWEID == nil && nvdReq.Component == nil && req.Query != "" {
		// Try to extract CWE ID from query
		if strings.Contains(strings.ToUpper(req.Query), "CWE-") {
			cweID := extractCWEID(req.Query)
			nvdReq.CWEID = &cweID
		} else {
			// Use query as component search
			nvdReq.Component = &req.Query
		}
	}

	cves, totalCount, err := s.securityRepo.QueryNVD(ctx, nvdReq)
	if err != nil {
		return nil, 0, err
	}

	// Convert to interface{} slice
	results := make([]interface{}, len(cves))
	for i, cve := range cves {
		results[i] = cve
	}

	return results, totalCount, nil
}

// queryATTACKData queries ATT&CK data with intelligent filtering
func (s *SecurityRetrievalService) queryATTACKData(ctx context.Context, req models.SecurityQueryRequest) ([]interface{}, int, error) {
	// Convert generic request to ATT&CK-specific request
	attackReq := models.ATTACKSearchRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	// Extract filters
	if filters, ok := req.Filters["technique_id"].(string); ok {
		attackReq.TechniqueID = &filters
	}
	if filters, ok := req.Filters["tactic"].(string); ok {
		attackReq.Tactic = &filters
	}
	if filters, ok := req.Filters["platform"].(string); ok {
		attackReq.Platform = &filters
	}

	// Use query for text search if no specific filters
	if attackReq.TechniqueID == nil && attackReq.Tactic == nil && req.Query != "" {
		// Try to extract technique ID from query
		if strings.Contains(strings.ToUpper(req.Query), "T") && len(req.Query) > 3 {
			techniqueID := extractTechniqueID(req.Query)
			if techniqueID != "" {
				attackReq.TechniqueID = &techniqueID
			} else {
				attackReq.Query = &req.Query
			}
		} else {
			attackReq.Query = &req.Query
		}
	}

	techniques, totalCount, err := s.securityRepo.QueryATTACK(ctx, attackReq)
	if err != nil {
		return nil, 0, err
	}

	// Convert to interface{} slice
	results := make([]interface{}, len(techniques))
	for i, technique := range techniques {
		results[i] = technique
	}

	return results, totalCount, nil
}

// queryOWASPData queries OWASP data (placeholder for future implementation)
func (s *SecurityRetrievalService) queryOWASPData(ctx context.Context, req models.SecurityQueryRequest) ([]interface{}, int, error) {
	// TODO: Implement OWASP data querying
	return []interface{}{}, 0, fmt.Errorf("OWASP data querying not yet implemented")
}

// summarizeResults creates context-efficient summaries of security data
func (s *SecurityRetrievalService) summarizeResults(results []interface{}, dataSource string, limit int) []interface{} {
	if len(results) == 0 {
		return results
	}

	// Limit results for context efficiency
	maxResults := limit
	if maxResults <= 0 {
		maxResults = 10
	}
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	// Create summaries based on data source
	switch dataSource {
	case "nvd":
		return s.summarizeNVDResults(results)
	case "attack":
		return s.summarizeATTACKResults(results)
	case "owasp":
		return s.summarizeOWASPResults(results)
	default:
		return results
	}
}

// summarizeNVDResults creates concise summaries of CVE data
func (s *SecurityRetrievalService) summarizeNVDResults(results []interface{}) []interface{} {
	summaries := make([]interface{}, len(results))

	for i, result := range results {
		if cve, ok := result.(models.NVDCVE); ok {
			summary := map[string]interface{}{
				"id":                cve.ID,
				"description":       truncateString(cve.Description, 200),
				"severity":          cve.Severity,
				"cvss_v3":           cve.CVSSV3Score,
				"published":         cve.PublishedDate,
				"cwe_ids":           cve.CWEIDs,
				"affected_products": truncateSlice(cve.AffectedProducts, 5),
				"summary":           s.generateCVESummary(cve),
			}
			summaries[i] = summary
		}
	}

	return summaries
}

// summarizeATTACKResults creates concise summaries of ATT&CK techniques
func (s *SecurityRetrievalService) summarizeATTACKResults(results []interface{}) []interface{} {
	summaries := make([]interface{}, len(results))

	for i, result := range results {
		if technique, ok := result.(models.ATTACKTechnique); ok {
			summary := map[string]interface{}{
				"id":          technique.ID,
				"name":        technique.Name,
				"description": truncateString(technique.Description, 200),
				"tactic":      technique.Tactic,
				"platforms":   technique.Platforms,
				"detection":   technique.Detection,
				"mitigation":  technique.Mitigation,
				"summary":     s.generateTechniqueSummary(technique),
			}
			summaries[i] = summary
		}
	}

	return summaries
}

// summarizeOWASPResults creates concise summaries of OWASP procedures
func (s *SecurityRetrievalService) summarizeOWASPResults(results []interface{}) []interface{} {
	// TODO: Implement OWASP result summarization
	return results
}

// generateCVESummary creates a concise summary of a CVE
func (s *SecurityRetrievalService) generateCVESummary(cve models.NVDCVE) string {
	severity := "Unknown"
	if cve.Severity != nil {
		severity = *cve.Severity
	}

	cvss := "N/A"
	if cve.CVSSV3Score != nil {
		cvss = fmt.Sprintf("%.1f", *cve.CVSSV3Score)
	}

	products := "Unknown"
	if len(cve.AffectedProducts) > 0 {
		products = strings.Join(cve.AffectedProducts[:min(3, len(cve.AffectedProducts))], ", ")
		if len(cve.AffectedProducts) > 3 {
			products += "..."
		}
	}

	return fmt.Sprintf("CVE %s: %s severity (CVSS %s) affecting %s",
		cve.ID, severity, cvss, products)
}

// generateTechniqueSummary creates a concise summary of an ATT&CK technique
func (s *SecurityRetrievalService) generateTechniqueSummary(technique models.ATTACKTechnique) string {
	platforms := "Unknown"
	if len(technique.Platforms) > 0 {
		platforms = strings.Join(technique.Platforms, ", ")
	}

	return fmt.Sprintf("ATT&CK %s (%s): %s technique for %s platforms",
		technique.ID, technique.Name, technique.Tactic, platforms)
}

// Helper functions
func extractCWEID(query string) string {
	// Simple CWE ID extraction
	parts := strings.Split(strings.ToUpper(query), "CWE-")
	if len(parts) > 1 {
		cwePart := strings.Split(parts[1], " ")[0]
		return "CWE-" + cwePart
	}
	return ""
}

func extractTechniqueID(query string) string {
	// Simple technique ID extraction
	query = strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(query, "T") && len(query) >= 4 {
		// Extract T followed by numbers and optional sub-technique
		for i, char := range query[1:] {
			if char < '0' || char > '9' {
				if char == '.' && i > 0 {
					continue // Allow decimal point for sub-techniques
				}
				return query[:i+1]
			}
		}
		return query
	}
	return ""
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func truncateSlice(slice []string, maxLen int) []string {
	if len(slice) <= maxLen {
		return slice
	}
	return slice[:maxLen]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
