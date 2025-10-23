package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/models"
	"golang.org/x/time/rate"
)

// SecurityDataDownloader handles downloading and updating security datasets
type SecurityDataDownloader struct {
	logger      *log.Logger
	client      *http.Client
	rateLimiter *rate.Limiter
}

// NewSecurityDataDownloader creates a new security data downloader
func NewSecurityDataDownloader(logger *log.Logger) *SecurityDataDownloader {
	return &SecurityDataDownloader{
		logger: logger,
		client: &http.Client{
			Timeout: 30 * time.Minute, // Long timeout for large downloads
		},
		// Rate limiting: 1 request per second for NVD API (respects their rate limits)
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

// rateLimitedGet performs a rate-limited HTTP GET request
func (s *SecurityDataDownloader) rateLimitedGet(ctx context.Context, url string) (*http.Response, error) {
	// Wait for rate limiter
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add proper headers to be respectful
	req.Header.Set("User-Agent", "TinyBrain-SecurityHub/1.0 (Security Research Tool)")
	req.Header.Set("Accept", "application/json")

	return s.client.Do(req)
}

// NVDResponse represents the response from NVD API
type NVDResponse struct {
	ResultsPerPage  int    `json:"resultsPerPage"`
	StartIndex      int    `json:"startIndex"`
	TotalResults    int    `json:"totalResults"`
	Format          string `json:"format"`
	Version         string `json:"version"`
	Timestamp       string `json:"timestamp"`
	Vulnerabilities []struct {
		CVE struct {
			ID               string `json:"id"`
			SourceIdentifier string `json:"sourceIdentifier"`
			Published        string `json:"published"`
			LastModified     string `json:"lastModified"`
			VulnStatus       string `json:"vulnStatus"`
			Descriptions     []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"descriptions"`
			Metrics struct {
				CVSSMetricV2 []struct {
					Source   string `json:"source"`
					Type     string `json:"type"`
					CVSSData struct {
						Version               string  `json:"version"`
						VectorString          string  `json:"vectorString"`
						BaseScore             float64 `json:"baseScore"`
						AccessVector          string  `json:"accessVector"`
						AccessComplexity      string  `json:"accessComplexity"`
						Authentication        string  `json:"authentication"`
						ConfidentialityImpact string  `json:"confidentialityImpact"`
						IntegrityImpact       string  `json:"integrityImpact"`
						AvailabilityImpact    string  `json:"availabilityImpact"`
					} `json:"cvssData"`
					BaseSeverity        string  `json:"baseSeverity"`
					ExploitabilityScore float64 `json:"exploitabilityScore"`
					ImpactScore         float64 `json:"impactScore"`
				} `json:"cvssMetricV2"`
				CVSSMetricV3 []struct {
					Source   string `json:"source"`
					Type     string `json:"type"`
					CVSSData struct {
						Version               string  `json:"version"`
						VectorString          string  `json:"vectorString"`
						BaseScore             float64 `json:"baseScore"`
						BaseSeverity          string  `json:"baseSeverity"`
						AttackVector          string  `json:"attackVector"`
						AttackComplexity      string  `json:"attackComplexity"`
						PrivilegesRequired    string  `json:"privilegesRequired"`
						UserInteraction       string  `json:"userInteraction"`
						Scope                 string  `json:"scope"`
						ConfidentialityImpact string  `json:"confidentialityImpact"`
						IntegrityImpact       string  `json:"integrityImpact"`
						AvailabilityImpact    string  `json:"availabilityImpact"`
					} `json:"cvssData"`
					ExploitabilityScore float64 `json:"exploitabilityScore"`
					ImpactScore         float64 `json:"impactScore"`
				} `json:"cvssMetricV3"`
			} `json:"metrics"`
			Weaknesses []struct {
				Source      string `json:"source"`
				Type        string `json:"type"`
				Description []struct {
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"description"`
			} `json:"weaknesses"`
			Configurations []struct {
				Nodes []struct {
					Operator string `json:"operator"`
					Negate   bool   `json:"negate"`
					CPEMatch []struct {
						Vulnerable      bool   `json:"vulnerable"`
						Criteria        string `json:"criteria"`
						MatchCriteriaID string `json:"matchCriteriaId"`
					} `json:"cpeMatch"`
				} `json:"nodes"`
			} `json:"configurations"`
			References []struct {
				URL    string `json:"url"`
				Source string `json:"source"`
			} `json:"references"`
		} `json:"cve"`
	} `json:"vulnerabilities"`
}

// DownloadNVDDataset downloads the complete NVD dataset
func (s *SecurityDataDownloader) DownloadNVDDataset(ctx context.Context) ([]models.NVDCVE, error) {
	s.logger.Info("Starting NVD dataset download")

	var allCVEs []models.NVDCVE
	startIndex := 0
	resultsPerPage := 2000 // NVD API max

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?startIndex=%d&resultsPerPage=%d",
			startIndex, resultsPerPage)

		s.logger.Debug("Downloading NVD batch", "startIndex", startIndex, "resultsPerPage", resultsPerPage)

		resp, err := s.rateLimitedGet(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to download NVD data: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read NVD response: %v", err)
		}

		var nvdResp NVDResponse
		if err := json.Unmarshal(body, &nvdResp); err != nil {
			return nil, fmt.Errorf("failed to parse NVD response: %v", err)
		}

		// Convert to our models
		for _, vuln := range nvdResp.Vulnerabilities {
			cve := s.convertNVDToModel(vuln)
			allCVEs = append(allCVEs, cve)
		}

		s.logger.Info("Downloaded NVD batch", "count", len(nvdResp.Vulnerabilities), "total", len(allCVEs))

		// Check if we've downloaded everything
		if startIndex+resultsPerPage >= nvdResp.TotalResults {
			break
		}

		startIndex += resultsPerPage

		// Rate limiting - be respectful to NVD API
		time.Sleep(1 * time.Second)
	}

	s.logger.Info("NVD dataset download completed", "total_cves", len(allCVEs))
	return allCVEs, nil
}

// convertNVDToModel converts NVD API response to our model
func (s *SecurityDataDownloader) convertNVDToModel(vuln interface{}) models.NVDCVE {
	// Convert to map for easier access
	vulnMap, ok := vuln.(map[string]interface{})
	if !ok {
		return models.NVDCVE{}
	}

	cveMap, ok := vulnMap["cve"].(map[string]interface{})
	if !ok {
		return models.NVDCVE{}
	}

	cve := models.NVDCVE{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Extract basic fields
	if id, ok := cveMap["id"].(string); ok {
		cve.ID = id
	}

	// Extract description
	if descriptions, ok := cveMap["descriptions"].([]interface{}); ok {
		for _, desc := range descriptions {
			if descMap, ok := desc.(map[string]interface{}); ok {
				if lang, ok := descMap["lang"].(string); ok && lang == "en" {
					if value, ok := descMap["value"].(string); ok {
						cve.Description = value
						break
					}
				}
			}
		}
	}

	// Extract CVSS scores
	if metrics, ok := cveMap["metrics"].(map[string]interface{}); ok {
		if cvssV3, ok := metrics["cvssMetricV3"].([]interface{}); ok && len(cvssV3) > 0 {
			if cvssMap, ok := cvssV3[0].(map[string]interface{}); ok {
				if cvssData, ok := cvssMap["cvssData"].(map[string]interface{}); ok {
					if score, ok := cvssData["baseScore"].(float64); ok {
						cve.CVSSV3Score = &score
					}
					if severity, ok := cvssData["baseSeverity"].(string); ok {
						severityUpper := strings.ToUpper(severity)
						cve.Severity = &severityUpper
					}
				}
			}
		}
	}

	// Extract CWE IDs
	if weaknesses, ok := cveMap["weaknesses"].([]interface{}); ok {
		for _, weakness := range weaknesses {
			if weaknessMap, ok := weakness.(map[string]interface{}); ok {
				if descriptions, ok := weaknessMap["description"].([]interface{}); ok {
					for _, desc := range descriptions {
						if descMap, ok := desc.(map[string]interface{}); ok {
							if value, ok := descMap["value"].(string); ok {
								if strings.Contains(strings.ToUpper(value), "CWE-") {
									parts := strings.Split(value, "CWE-")
									if len(parts) > 1 {
										cwePart := strings.Split(parts[1], " ")[0]
										cve.CWEIDs = append(cve.CWEIDs, "CWE-"+cwePart)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Extract affected products
	if configurations, ok := cveMap["configurations"].([]interface{}); ok {
		for _, config := range configurations {
			if configMap, ok := config.(map[string]interface{}); ok {
				if nodes, ok := configMap["nodes"].([]interface{}); ok {
					for _, node := range nodes {
						if nodeMap, ok := node.(map[string]interface{}); ok {
							if cpeMatch, ok := nodeMap["cpeMatch"].([]interface{}); ok {
								for _, cpe := range cpeMatch {
									if cpeMap, ok := cpe.(map[string]interface{}); ok {
										if vulnerable, ok := cpeMap["vulnerable"].(bool); ok && vulnerable {
											if criteria, ok := cpeMap["criteria"].(string); ok {
												cve.AffectedProducts = append(cve.AffectedProducts, criteria)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Extract references
	if references, ok := cveMap["references"].([]interface{}); ok {
		for _, ref := range references {
			if refMap, ok := ref.(map[string]interface{}); ok {
				if url, ok := refMap["url"].(string); ok {
					cve.References = append(cve.References, url)
				}
			}
		}
	}

	// Store raw data
	if rawData, err := json.Marshal(vuln); err == nil {
		cve.RawData = string(rawData)
	}

	return cve
}

// DownloadATTACKDataset downloads the MITRE ATT&CK dataset
func (s *SecurityDataDownloader) DownloadATTACKDataset(ctx context.Context) ([]models.ATTACKTechnique, []models.ATTACKTactic, error) {
	s.logger.Info("Starting MITRE ATT&CK dataset download")

	url := "https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json"

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download ATT&CK data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ATT&CK response: %v", err)
	}

	var bundle struct {
		Objects []json.RawMessage `json:"objects"`
	}

	if err := json.Unmarshal(body, &bundle); err != nil {
		return nil, nil, fmt.Errorf("failed to parse ATT&CK response: %v", err)
	}

	var techniques []models.ATTACKTechnique
	var tactics []models.ATTACKTactic

	for _, obj := range bundle.Objects {
		var baseObj struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(obj, &baseObj); err != nil {
			continue
		}

		switch baseObj.Type {
		case "attack-pattern":
			technique := s.convertATTACKTechnique(obj)
			techniques = append(techniques, technique)
		case "x-mitre-tactic":
			tactic := s.convertATTACKTactic(obj)
			tactics = append(tactics, tactic)
		}
	}

	s.logger.Info("ATT&CK dataset download completed", "techniques", len(techniques), "tactics", len(tactics))
	return techniques, tactics, nil
}

// DownloadOWASPDataset downloads OWASP security testing procedures and guidelines
func (s *SecurityDataDownloader) DownloadOWASPDataset(ctx context.Context) ([]models.OWASPProcedure, error) {
	s.logger.Info("Starting OWASP dataset download")

	// OWASP Testing Guide - Official GitHub repository
	url := "https://raw.githubusercontent.com/OWASP/wstg/master/document/4_Web_Application_Security_Testing_Guide/README.md"

	resp, err := s.rateLimitedGet(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to download OWASP data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OWASP response: %v", err)
	}

	// Parse OWASP Testing Guide content
	procedures := s.parseOWASPContent(string(body))

	s.logger.Info("OWASP dataset download completed", "procedures", len(procedures))
	return procedures, nil
}

// parseOWASPContent parses OWASP content into structured procedures
func (s *SecurityDataDownloader) parseOWASPContent(content string) []models.OWASPProcedure {
	var procedures []models.OWASPProcedure

	// This is a simplified parser - in a real implementation, you'd want more sophisticated parsing
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Look for procedure headers (simplified pattern)
		if strings.Contains(line, "##") && (strings.Contains(strings.ToLower(line), "test") ||
			strings.Contains(strings.ToLower(line), "procedure")) {

			procedure := models.OWASPProcedure{
				ID:          fmt.Sprintf("OWASP-%d", i),
				Title:       strings.TrimSpace(strings.TrimPrefix(line, "##")),
				Category:    "Web Application Security Testing",
				Description: "OWASP Testing Guide procedure",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Try to extract description from next few lines
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if lines[j] != "" && !strings.HasPrefix(lines[j], "#") {
					procedure.Description = strings.TrimSpace(lines[j])
					break
				}
			}

			procedures = append(procedures, procedure)
		}
	}

	return procedures
}

// CheckForUpdates checks if there are updates available for security datasets
func (s *SecurityDataDownloader) CheckForUpdates(ctx context.Context, dataSource string, lastUpdate time.Time) (bool, *time.Time, error) {
	switch dataSource {
	case "nvd":
		return s.checkNVDUpdates(ctx, lastUpdate)
	case "attack":
		return s.checkATTACKUpdates(ctx, lastUpdate)
	case "owasp":
		return s.checkOWASPUpdates(ctx, lastUpdate)
	default:
		return false, nil, fmt.Errorf("unknown data source: %s", dataSource)
	}
}

// checkNVDUpdates checks for new CVEs since last update
func (s *SecurityDataDownloader) checkNVDUpdates(ctx context.Context, lastUpdate time.Time) (bool, *time.Time, error) {
	// NVD API supports date filtering - check for CVEs modified since last update
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?lastModStartDate=%s&resultsPerPage=1",
		lastUpdate.Format("2006-01-02T15:04:05.000"))

	resp, err := s.rateLimitedGet(ctx, url)
	if err != nil {
		return false, nil, fmt.Errorf("failed to check NVD updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read NVD update check response: %v", err)
	}

	var nvdResp NVDResponse
	if err := json.Unmarshal(body, &nvdResp); err != nil {
		return false, nil, fmt.Errorf("failed to parse NVD update check response: %v", err)
	}

	// If there are any results, there are updates
	hasUpdates := nvdResp.TotalResults > 0

	// Get the latest modification time from the first result
	var latestUpdate *time.Time
	if hasUpdates && len(nvdResp.Vulnerabilities) > 0 {
		if vuln := nvdResp.Vulnerabilities[0]; vuln.CVE.LastModified != "" {
			if modified, err := time.Parse(time.RFC3339, vuln.CVE.LastModified); err == nil {
				latestUpdate = &modified
			}
		}
	}

	return hasUpdates, latestUpdate, nil
}

// checkATTACKUpdates checks for updates to MITRE ATT&CK dataset
func (s *SecurityDataDownloader) checkATTACKUpdates(ctx context.Context, lastUpdate time.Time) (bool, *time.Time, error) {
	// Check GitHub API for latest commit to the ATT&CK repository
	url := "https://api.github.com/repos/mitre/cti/commits?path=enterprise-attack/enterprise-attack.json&per_page=1"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "TinyBrain-SecurityHub/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("failed to check ATT&CK updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read ATT&CK update check response: %v", err)
	}

	var commits []struct {
		Commit struct {
			Committer struct {
				Date string `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	}

	if err := json.Unmarshal(body, &commits); err != nil {
		return false, nil, fmt.Errorf("failed to parse ATT&CK update check response: %v", err)
	}

	if len(commits) == 0 {
		return false, nil, nil
	}

	// Parse the latest commit date
	latestCommitDate, err := time.Parse(time.RFC3339, commits[0].Commit.Committer.Date)
	if err != nil {
		return false, nil, fmt.Errorf("failed to parse commit date: %v", err)
	}

	hasUpdates := latestCommitDate.After(lastUpdate)
	return hasUpdates, &latestCommitDate, nil
}

// checkOWASPUpdates checks for updates to OWASP Testing Guide
func (s *SecurityDataDownloader) checkOWASPUpdates(ctx context.Context, lastUpdate time.Time) (bool, *time.Time, error) {
	// Check GitHub API for latest commit to the OWASP WSTG repository
	url := "https://api.github.com/repos/OWASP/wstg/commits?path=document/4_Web_Application_Security_Testing_Guide/README.md&per_page=1"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "TinyBrain-SecurityHub/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("failed to check OWASP updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read OWASP update check response: %v", err)
	}

	var commits []struct {
		Commit struct {
			Committer struct {
				Date string `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	}

	if err := json.Unmarshal(body, &commits); err != nil {
		return false, nil, fmt.Errorf("failed to parse OWASP update check response: %v", err)
	}

	if len(commits) == 0 {
		return false, nil, nil
	}

	// Parse the latest commit date
	latestCommitDate, err := time.Parse(time.RFC3339, commits[0].Commit.Committer.Date)
	if err != nil {
		return false, nil, fmt.Errorf("failed to parse commit date: %v", err)
	}

	hasUpdates := latestCommitDate.After(lastUpdate)
	return hasUpdates, &latestCommitDate, nil
}

// DownloadIncrementalNVD downloads only new CVEs since last update
func (s *SecurityDataDownloader) DownloadIncrementalNVD(ctx context.Context, lastUpdate time.Time) ([]models.NVDCVE, error) {
	s.logger.Info("Starting incremental NVD download", "since", lastUpdate)

	// Use NVD API's date filtering to get only new/updated CVEs
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?lastModStartDate=%s&resultsPerPage=2000",
		lastUpdate.Format("2006-01-02T15:04:05.000"))

	resp, err := s.rateLimitedGet(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to download incremental NVD data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read incremental NVD response: %v", err)
	}

	var nvdResp NVDResponse
	if err := json.Unmarshal(body, &nvdResp); err != nil {
		return nil, fmt.Errorf("failed to parse incremental NVD response: %v", err)
	}

	var cves []models.NVDCVE
	for _, vuln := range nvdResp.Vulnerabilities {
		cve := s.convertNVDToModel(vuln)
		cves = append(cves, cve)
	}

	s.logger.Info("Incremental NVD download completed", "new_cves", len(cves))
	return cves, nil
}

// convertATTACKTechnique converts ATT&CK technique to our model
func (s *SecurityDataDownloader) convertATTACKTechnique(obj json.RawMessage) models.ATTACKTechnique {
	var technique struct {
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		XMitrePlatforms []string `json:"x_mitre_platforms"`
		KillChainPhases []struct {
			KillChainName string `json:"kill_chain_name"`
			PhaseName     string `json:"phase_name"`
		} `json:"kill_chain_phases"`
		XMitreDataSources []struct {
			DataSource string `json:"data_source"`
		} `json:"x_mitre_data_sources"`
		XMitreDetection    string `json:"x_mitre_detection"`
		XMitreMitigation   string `json:"x_mitre_mitigation"`
		ExternalReferences []struct {
			SourceName string `json:"source_name"`
			URL        string `json:"url"`
			ExternalID string `json:"external_id"`
		} `json:"external_references"`
		XMitreIsSubtechnique bool `json:"x_mitre_is_subtechnique"`
		XMitreSubtechniqueOf *struct {
			TargetRef string `json:"target_ref"`
		} `json:"x_mitre_subtechnique_of"`
	}

	json.Unmarshal(obj, &technique)

	result := models.ATTACKTechnique{
		ID:          technique.ID,
		Name:        technique.Name,
		Description: technique.Description,
		Platforms:   technique.XMitrePlatforms,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Extract tactics from kill chain phases
	var tactics []string
	for _, phase := range technique.KillChainPhases {
		if phase.KillChainName == "mitre-attack" {
			tactics = append(tactics, phase.PhaseName)
		}
	}
	result.Tactics = tactics
	if len(tactics) > 0 {
		result.Tactic = tactics[0] // Primary tactic
	}

	// Extract kill chain phases
	var phases []string
	for _, phase := range technique.KillChainPhases {
		phases = append(phases, phase.PhaseName)
	}
	result.KillChainPhases = phases

	// Extract data sources
	var sources []string
	for _, source := range technique.XMitreDataSources {
		sources = append(sources, source.DataSource)
	}
	result.DataSources = sources

	// Set detection and mitigation
	if technique.XMitreDetection != "" {
		result.Detection = &technique.XMitreDetection
	}
	if technique.XMitreMitigation != "" {
		result.Mitigation = &technique.XMitreMitigation
	}

	// Extract references
	var references []string
	for _, ref := range technique.ExternalReferences {
		if ref.URL != "" {
			references = append(references, ref.URL)
		}
	}
	result.References = references

	// Handle sub-techniques
	if technique.XMitreIsSubtechnique && technique.XMitreSubtechniqueOf != nil {
		result.ParentTechnique = &technique.XMitreSubtechniqueOf.TargetRef
	}

	// Store raw data
	result.RawData = string(obj)

	return result
}

// convertATTACKTactic converts ATT&CK tactic to our model
func (s *SecurityDataDownloader) convertATTACKTactic(obj json.RawMessage) models.ATTACKTactic {
	var tactic struct {
		ID                 string `json:"id"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		ExternalReferences []struct {
			SourceName string `json:"source_name"`
			URL        string `json:"url"`
			ExternalID string `json:"external_id"`
		} `json:"external_references"`
		KillChainPhases []struct {
			KillChainName string `json:"kill_chain_name"`
			PhaseName     string `json:"phase_name"`
		} `json:"kill_chain_phases"`
	}

	json.Unmarshal(obj, &tactic)

	result := models.ATTACKTactic{
		ID:          tactic.ID,
		Name:        tactic.Name,
		Description: tactic.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Extract external ID
	for _, ref := range tactic.ExternalReferences {
		if ref.SourceName == "mitre-attack" && ref.ExternalID != "" {
			result.ExternalID = &ref.ExternalID
			break
		}
	}

	// Extract kill chain phases
	var phases []string
	for _, phase := range tactic.KillChainPhases {
		phases = append(phases, phase.PhaseName)
	}
	result.KillChainPhases = phases

	return result
}

// GetDataSize estimates the size of security datasets
func (s *SecurityDataDownloader) GetDataSize() map[string]interface{} {
	return map[string]interface{}{
		"nvd": map[string]interface{}{
			"estimated_records": 314835,
			"estimated_size_mb": 50,
			"api_endpoint":      "https://services.nvd.nist.gov/rest/json/cves/2.0",
		},
		"attack": map[string]interface{}{
			"estimated_records": 600,
			"estimated_size_mb": 38,
			"download_url":      "https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json",
		},
		"owasp": map[string]interface{}{
			"estimated_records": 1000,
			"estimated_size_mb": 10,
			"status":            "to_be_implemented",
		},
	}
}
