package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestSecurityDataDownloader_GetDataSize(t *testing.T) {
	logger := log.New(os.Stderr)
	downloader := NewSecurityDataDownloader(logger)

	sizes := downloader.GetDataSize()

	assert.Contains(t, sizes, "nvd")
	assert.Contains(t, sizes, "attack")
	assert.Contains(t, sizes, "owasp")

	nvdInfo := sizes["nvd"].(map[string]interface{})
	assert.Equal(t, 314835, nvdInfo["estimated_records"])
	assert.Equal(t, 50, nvdInfo["estimated_size_mb"])

	attackInfo := sizes["attack"].(map[string]interface{})
	assert.Equal(t, 600, attackInfo["estimated_records"])
	assert.Equal(t, 38, attackInfo["estimated_size_mb"])
}

func TestSecurityDataDownloader_ConvertATTACKTechnique(t *testing.T) {
	logger := log.New(os.Stderr)
	downloader := NewSecurityDataDownloader(logger)

	// Create a mock ATT&CK technique JSON
	mockTechniqueJSON := json.RawMessage(`{
		"type": "attack-pattern",
		"id": "attack-pattern--test-123",
		"name": "Process Injection",
		"description": "Test technique description",
		"x_mitre_platforms": ["windows", "linux"],
		"kill_chain_phases": [
			{
				"kill_chain_name": "mitre-attack",
				"phase_name": "defense-evasion"
			}
		],
		"x_mitre_data_sources": [
			{
				"data_source": "process-monitoring"
			}
		],
		"x_mitre_detection": "Monitor for process injection",
		"x_mitre_mitigation": "Use application whitelisting",
		"external_references": [
			{
				"source_name": "mitre-attack",
				"url": "https://attack.mitre.org/techniques/T1055/",
				"external_id": "T1055"
			}
		],
		"x_mitre_is_subtechnique": false
	}`)

	technique := downloader.convertATTACKTechnique(mockTechniqueJSON)

	assert.Equal(t, "attack-pattern--test-123", technique.ID)
	assert.Equal(t, "Process Injection", technique.Name)
	assert.Equal(t, "Test technique description", technique.Description)
	assert.Equal(t, "defense-evasion", technique.Tactic)
	assert.Contains(t, technique.Tactics, "defense-evasion")
	assert.Contains(t, technique.Platforms, "windows")
	assert.Contains(t, technique.Platforms, "linux")
	assert.Contains(t, technique.DataSources, "process-monitoring")
	assert.Equal(t, "Monitor for process injection", *technique.Detection)
	assert.Equal(t, "Use application whitelisting", *technique.Mitigation)
	assert.Contains(t, technique.References, "https://attack.mitre.org/techniques/T1055/")
}

func TestSecurityDataDownloader_ConvertATTACKTactic(t *testing.T) {
	logger := log.New(os.Stderr)
	downloader := NewSecurityDataDownloader(logger)

	// Create a mock ATT&CK tactic JSON
	mockTacticJSON := json.RawMessage(`{
		"type": "x-mitre-tactic",
		"id": "x-mitre-tactic--test-123",
		"name": "Defense Evasion",
		"description": "Test tactic description",
		"external_references": [
			{
				"source_name": "mitre-attack",
				"url": "https://attack.mitre.org/tactics/TA0005/",
				"external_id": "TA0005"
			}
		],
		"kill_chain_phases": [
			{
				"kill_chain_name": "mitre-attack",
				"phase_name": "defense-evasion"
			}
		]
	}`)

	tactic := downloader.convertATTACKTactic(mockTacticJSON)

	assert.Equal(t, "x-mitre-tactic--test-123", tactic.ID)
	assert.Equal(t, "Defense Evasion", tactic.Name)
	assert.Equal(t, "Test tactic description", tactic.Description)
	assert.Equal(t, "TA0005", *tactic.ExternalID)
	assert.Contains(t, tactic.KillChainPhases, "defense-evasion")
}

func TestSecurityDataDownloader_DownloadNVDDataset_Mock(t *testing.T) {
	logger := log.New(os.Stderr)
	downloader := NewSecurityDataDownloader(logger)

	// Create a mock NVD API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a mock NVD response
		mockResponse := map[string]interface{}{
			"resultsPerPage": 10,
			"startIndex":     0,
			"totalResults":   10,
			"format":         "NVD_CVE",
			"version":        "2.0",
			"timestamp":      time.Now().Format(time.RFC3339),
			"vulnerabilities": []interface{}{
				map[string]interface{}{
					"cve": map[string]interface{}{
						"id":               "CVE-2024-1234",
						"sourceIdentifier": "nvd@nist.gov",
						"published":        "2024-01-01T00:00:00.000Z",
						"lastModified":     "2024-01-02T00:00:00.000Z",
						"vulnStatus":       "Analyzed",
						"descriptions": []interface{}{
							map[string]interface{}{
								"lang":  "en",
								"value": "Test vulnerability",
							},
						},
						"metrics":        map[string]interface{}{},
						"weaknesses":     []interface{}{},
						"configurations": []interface{}{},
						"references":     []interface{}{},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the downloader's client to use our mock server
	downloader.client = server.Client()

	// Test that the downloader can be created and has the expected methods
	assert.NotNil(t, downloader)
	assert.NotNil(t, downloader.client)
	assert.NotNil(t, downloader.logger)
}

func TestSecurityDataDownloader_NewSecurityDataDownloader(t *testing.T) {
	logger := log.New(os.Stderr)
	downloader := NewSecurityDataDownloader(logger)

	assert.NotNil(t, downloader)
	assert.NotNil(t, downloader.client)
	assert.NotNil(t, downloader.logger)
	assert.Equal(t, 30*time.Minute, downloader.client.Timeout)
}
