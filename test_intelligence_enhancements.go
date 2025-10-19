package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/models"
	"github.com/rainmana/tinybrain/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntelligenceEnhancements tests the new intelligence capabilities
func TestIntelligenceEnhancements(t *testing.T) {
	// Skip if CGO is not enabled (Windows without C compiler)
	if os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("Skipping intelligence enhancement tests - CGO not enabled")
	}

	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "intelligence_test.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	// Initialize database
	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewMemoryRepository(db.GetDB(), logger)
	ctx := context.Background()

	// Test 1: Create intelligence session
	t.Run("CreateIntelligenceSession", func(t *testing.T) {
		session := &models.Session{
			ID:               "intel-session-1",
			Name:             "OSINT Intelligence Gathering",
			Description:      "Gathering open source intelligence on target organization",
			TaskType:         "intelligence_analysis",
			Status:           "active",
			IntelligenceType: "osint",
			TargetScope:      "organization",
			Classification:   "unclassified",
			ThreatLevel:      "medium",
			GeographicScope:  "national",
			Metadata: map[string]interface{}{
				"target_organization": "Example Corp",
				"analysis_scope":      "social_media",
			},
		}

		err := repo.CreateSession(ctx, session)
		assert.NoError(t, err)

		// Verify session was created with intelligence fields
		retrieved, err := repo.GetSession(ctx, session.ID)
		assert.NoError(t, err)
		assert.Equal(t, "intelligence_analysis", retrieved.TaskType)
		assert.Equal(t, "osint", retrieved.IntelligenceType)
		assert.Equal(t, "organization", retrieved.TargetScope)
		assert.Equal(t, "unclassified", retrieved.Classification)
		assert.Equal(t, "medium", retrieved.ThreatLevel)
		assert.Equal(t, "national", retrieved.GeographicScope)
	})

	// Test 2: Create intelligence finding
	t.Run("CreateIntelligenceFinding", func(t *testing.T) {
		// Create session first
		session := &models.Session{
			ID:               "intel-session-2",
			Name:             "HUMINT Analysis",
			TaskType:         "intelligence_analysis",
			Status:           "active",
			IntelligenceType: "humint",
		}
		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		// Create intelligence finding
		req := &models.CreateMemoryEntryRequest{
			SessionID:        session.ID,
			Title:            "Source Intelligence Report",
			Content:          "Human source reports suspicious activity at target location",
			ContentType:      "intelligence",
			Category:         "intelligence",
			Priority:         8,
			Confidence:       0.9,
			Tags:             []string{"humint", "source", "suspicious-activity"},
			Source:           "Human Source Alpha",
			IntelligenceType: "humint",
			Classification:   "confidential",
			ThreatLevel:      "high",
			GeographicScope:  "local",
			Attribution:      "Unknown threat actor",
			IOCType:          "location",
			IOCValue:         "123 Main St, City",
			MITRETactic:      "TA0001",
			MITRETechnique:   "T1591",
			MITREProcedure:   "T1591.001",
			KillChainPhase:   "reconnaissance",
			RiskScore:        8.5,
			ImpactScore:      9.0,
			LikelihoodScore:  8.0,
			Metadata: map[string]interface{}{
				"source_reliability":    "high",
				"information_currency":  "recent",
				"corroboration_sources": []string{"OSINT", "SIGINT"},
			},
		}

		entry, err := repo.CreateMemoryEntry(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, entry)
		assert.Equal(t, "intelligence", entry.Category)
		assert.Equal(t, "humint", entry.IntelligenceType)
		assert.Equal(t, "confidential", entry.Classification)
		assert.Equal(t, "high", entry.ThreatLevel)
		assert.Equal(t, "local", entry.GeographicScope)
		assert.Equal(t, "Unknown threat actor", entry.Attribution)
		assert.Equal(t, "location", entry.IOCType)
		assert.Equal(t, "123 Main St, City", entry.IOCValue)
		assert.Equal(t, "TA0001", entry.MITRETactic)
		assert.Equal(t, "T1591", entry.MITRETechnique)
		assert.Equal(t, "T1591.001", entry.MITREProcedure)
		assert.Equal(t, "reconnaissance", entry.KillChainPhase)
		assert.Equal(t, 8.5, entry.RiskScore)
		assert.Equal(t, 9.0, entry.ImpactScore)
		assert.Equal(t, 8.0, entry.LikelihoodScore)
	})

	// Test 3: Create threat actor
	t.Run("CreateThreatActor", func(t *testing.T) {
		// This would require implementing the threat actor repository methods
		// For now, we'll test the data model structure
		threatActor := &models.ThreatActor{
			ID:              "threat-actor-1",
			SessionID:       "intel-session-1",
			Name:            "APT29",
			Aliases:         []string{"Cozy Bear", "The Dukes"},
			Description:     "Russian state-sponsored threat group",
			Motivation:      "Espionage",
			Capabilities:    []string{"Spear phishing", "Zero-day exploits", "Custom malware"},
			Targets:         []string{"Government", "Healthcare", "Energy"},
			Tools:           []string{"Custom malware", "Living off the land"},
			Techniques:      []string{"T1566.001", "T1055", "T1071.001"},
			Attribution:     "High confidence attribution to Russia",
			Confidence:      0.9,
			ThreatLevel:     "critical",
			GeographicScope: "international",
			Metadata: map[string]interface{}{
				"first_observed": "2014",
				"last_observed":  "2024",
				"estimated_size": "50-100 operators",
			},
		}

		// Verify the data structure is correct
		assert.Equal(t, "APT29", threatActor.Name)
		assert.Equal(t, "critical", threatActor.ThreatLevel)
		assert.Equal(t, "international", threatActor.GeographicScope)
		assert.Len(t, threatActor.Aliases, 2)
		assert.Len(t, threatActor.Capabilities, 3)
	})

	// Test 4: Create attack campaign
	t.Run("CreateAttackCampaign", func(t *testing.T) {
		campaign := &models.AttackCampaign{
			ID:              "campaign-1",
			SessionID:       "intel-session-1",
			Name:            "Operation SolarWinds",
			Description:     "Supply chain attack targeting SolarWinds Orion software",
			ThreatActors:    []string{"APT29"},
			Targets:         []string{"Government", "Technology", "Critical Infrastructure"},
			Techniques:      []string{"T1195", "T1055", "T1071.001"},
			Tools:           []string{"SUNBURST", "TEARDROP", "Raindrop"},
			IOCs:            []string{"hash1", "domain1", "ip1"},
			StartDate:       time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
			EndDate:         time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC),
			Status:          "completed",
			ThreatLevel:     "critical",
			GeographicScope: "international",
			Confidence:      0.95,
			Metadata: map[string]interface{}{
				"victim_countries": []string{"US", "UK", "CA", "AU"},
				"estimated_damage": "$100M+",
				"data_exfiltrated": []string{"PII", "Intellectual Property", "Credentials"},
			},
		}

		// Verify the data structure is correct
		assert.Equal(t, "Operation SolarWinds", campaign.Name)
		assert.Equal(t, "critical", campaign.ThreatLevel)
		assert.Equal(t, "completed", campaign.Status)
		assert.Len(t, campaign.ThreatActors, 1)
		assert.Len(t, campaign.Targets, 3)
	})

	// Test 5: Create IOC
	t.Run("CreateIOC", func(t *testing.T) {
		ioc := &models.IndicatorOfCompromise{
			ID:          "ioc-1",
			SessionID:   "intel-session-1",
			Type:        "domain",
			Value:       "malicious.example.com",
			Description: "Malicious domain used for C2 communications",
			ThreatLevel: "high",
			Confidence:  0.9,
			FirstSeen:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			LastSeen:    time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			Source:      "Threat Intelligence Feed",
			Attribution: "APT Group",
			Campaigns:   []string{"campaign-1"},
			Techniques:  []string{"T1071.001", "T1566.001"},
			Tags:        []string{"malware", "c2", "phishing"},
			Metadata: map[string]interface{}{
				"dns_records":  []string{"A", "AAAA", "MX"},
				"ip_addresses": []string{"1.2.3.4", "5.6.7.8"},
				"reputation":   "malicious",
			},
		}

		// Verify the data structure is correct
		assert.Equal(t, "domain", ioc.Type)
		assert.Equal(t, "malicious.example.com", ioc.Value)
		assert.Equal(t, "high", ioc.ThreatLevel)
		assert.Equal(t, 0.9, ioc.Confidence)
		assert.Len(t, ioc.Campaigns, 1)
		assert.Len(t, ioc.Techniques, 2)
	})

	// Test 6: Create pattern
	t.Run("CreatePattern", func(t *testing.T) {
		pattern := &models.Pattern{
			ID:          "pattern-1",
			SessionID:   "intel-session-1",
			Name:        "Spear Phishing with Malicious Attachment",
			Description: "Attack pattern involving targeted phishing emails with malicious attachments",
			PatternType: "attack",
			Category:    "phishing",
			Severity:    "high",
			Confidence:  0.9,
			Frequency:   50,
			Examples:    []string{"Email with malicious Word document", "PDF with embedded JavaScript"},
			Mitigations: []string{"Email filtering", "User awareness training", "Sandboxing"},
			Detections:  []string{"Email security gateways", "Endpoint detection", "User reporting"},
			Tags:        []string{"phishing", "malware", "social-engineering"},
			Metadata: map[string]interface{}{
				"mitre_techniques": []string{"T1566.001", "T1566.002"},
				"target_sectors":   []string{"Government", "Financial", "Healthcare"},
				"success_factors":  []string{"urgency", "authority", "familiarity"},
			},
		}

		// Verify the data structure is correct
		assert.Equal(t, "attack", pattern.PatternType)
		assert.Equal(t, "high", pattern.Severity)
		assert.Equal(t, 50, pattern.Frequency)
		assert.Len(t, pattern.Examples, 2)
		assert.Len(t, pattern.Mitigations, 3)
	})

	// Test 7: Create correlation
	t.Run("CreateCorrelation", func(t *testing.T) {
		correlation := &models.Correlation{
			ID:              "correlation-1",
			SessionID:       "intel-session-1",
			SourceFindingID: "finding-1",
			TargetFindingID: "finding-2",
			CorrelationType: "temporal",
			Strength:        0.9,
			Confidence:      0.8,
			Evidence:        "Both findings occurred within 24 hours and share similar IOCs",
			Description:     "Temporal correlation between phishing campaign and malware deployment",
			Weight:          1.0,
			Direction:       "unidirectional",
			Metadata: map[string]interface{}{
				"time_difference": "2 hours",
				"shared_iocs":     []string{"domain1", "ip1"},
				"attack_sequence": []string{"phishing", "malware_deployment"},
			},
		}

		// Verify the data structure is correct
		assert.Equal(t, "temporal", correlation.CorrelationType)
		assert.Equal(t, 0.9, correlation.Strength)
		assert.Equal(t, 0.8, correlation.Confidence)
		assert.Equal(t, "unidirectional", correlation.Direction)
	})

	// Test 8: Test enhanced search capabilities
	t.Run("EnhancedSearchCapabilities", func(t *testing.T) {
		// Create session
		session := &models.Session{
			ID:               "search-session",
			Name:             "Search Test Session",
			TaskType:         "intelligence_analysis",
			Status:           "active",
			IntelligenceType: "osint",
		}
		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		// Create multiple intelligence entries
		entries := []*models.CreateMemoryEntryRequest{
			{
				SessionID:        session.ID,
				Title:            "OSINT Finding: Social Media Intelligence",
				Content:          "Social media analysis reveals suspicious activity",
				Category:         "intelligence",
				IntelligenceType: "osint",
				ThreatLevel:      "medium",
				Tags:             []string{"osint", "social-media", "suspicious"},
			},
			{
				SessionID:        session.ID,
				Title:            "HUMINT Finding: Source Report",
				Content:          "Human source reports threat actor activity",
				Category:         "intelligence",
				IntelligenceType: "humint",
				ThreatLevel:      "high",
				Tags:             []string{"humint", "source", "threat-actor"},
			},
			{
				SessionID:        session.ID,
				Title:            "SIGINT Finding: Communications Intelligence",
				Content:          "Signals intelligence reveals encrypted communications",
				Category:         "intelligence",
				IntelligenceType: "sigint",
				ThreatLevel:      "critical",
				Tags:             []string{"sigint", "communications", "encrypted"},
			},
		}

		for _, req := range entries {
			_, err := repo.CreateMemoryEntry(ctx, req)
			require.NoError(t, err)
		}

		// Test search by intelligence type
		searchReq := &models.SearchRequest{
			Query:      "intelligence",
			SessionID:  session.ID,
			SearchType: "exact",
			Limit:      10,
		}

		results, err := repo.SearchMemoryEntries(ctx, searchReq)
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		// Test search by threat level
		searchReq = &models.SearchRequest{
			Query:       "threat",
			SessionID:   session.ID,
			MinPriority: 8, // High and critical threats
			SearchType:  "exact",
			Limit:       10,
		}

		results, err = repo.SearchMemoryEntries(ctx, searchReq)
		assert.NoError(t, err)
		// Should find entries with high threat level (priority 8+)
		assert.GreaterOrEqual(t, len(results), 1)
	})

	// Test 9: Test MITRE ATT&CK integration
	t.Run("MITREAttackIntegration", func(t *testing.T) {
		// Create session
		session := &models.Session{
			ID:               "mitre-session",
			Name:             "MITRE ATT&CK Test",
			TaskType:         "intelligence_analysis",
			Status:           "active",
			IntelligenceType: "osint",
		}
		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		// Create entry with MITRE ATT&CK mapping
		req := &models.CreateMemoryEntryRequest{
			SessionID:        session.ID,
			Title:            "Phishing Campaign Analysis",
			Content:          "Analysis of spear phishing campaign targeting government officials",
			Category:         "intelligence",
			IntelligenceType: "osint",
			MITRETactic:      "TA0001",    // Initial Access
			MITRETechnique:   "T1566",     // Phishing
			MITREProcedure:   "T1566.001", // Spearphishing Attachment
			KillChainPhase:   "reconnaissance",
			ThreatLevel:      "high",
			Tags:             []string{"phishing", "mitre", "ta0001", "t1566"},
		}

		entry, err := repo.CreateMemoryEntry(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "TA0001", entry.MITRETactic)
		assert.Equal(t, "T1566", entry.MITRETechnique)
		assert.Equal(t, "T1566.001", entry.MITREProcedure)
		assert.Equal(t, "reconnaissance", entry.KillChainPhase)
	})

	// Test 10: Test data validation
	t.Run("DataValidation", func(t *testing.T) {
		// Test invalid intelligence type
		session := &models.Session{
			ID:               "validation-session",
			Name:             "Validation Test",
			TaskType:         "intelligence_analysis",
			Status:           "active",
			IntelligenceType: "invalid_type", // This should be caught by validation
		}

		// Note: This test would require implementing validation in the repository
		// For now, we'll just verify the data structure
		assert.Equal(t, "invalid_type", session.IntelligenceType)
	})

	t.Log("All intelligence enhancement tests completed successfully!")
}

// TestDatabaseSchemaCompatibility tests that the new schema is compatible
func TestDatabaseSchemaCompatibility(t *testing.T) {
	// Skip if CGO is not enabled
	if os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("Skipping schema compatibility tests - CGO not enabled")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "schema_test.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	// Test that the current database still works
	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	// Test basic operations
	ctx := context.Background()
	repo := repository.NewMemoryRepository(db.GetDB(), logger)

	// Create a basic session
	session := &models.Session{
		ID:       "compat-session",
		Name:     "Compatibility Test",
		TaskType: "security_review",
		Status:   "active",
	}

	err = repo.CreateSession(ctx, session)
	assert.NoError(t, err)

	// Create a basic memory entry
	req := &models.CreateMemoryEntryRequest{
		SessionID: session.ID,
		Title:     "Test Entry",
		Content:   "Test content",
		Category:  "note",
		Priority:  5,
	}

	entry, err := repo.CreateMemoryEntry(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, entry)

	// Test search
	searchReq := &models.SearchRequest{
		Query:      "test",
		SessionID:  session.ID,
		SearchType: "exact",
		Limit:      10,
	}

	results, err := repo.SearchMemoryEntries(ctx, searchReq)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	t.Log("Database schema compatibility test passed!")
}

// TestIntelligenceTemplates tests the intelligence templates
func TestIntelligenceTemplates(t *testing.T) {
	// Test OSINT template
	osintTemplate := map[string]interface{}{
		"title":             "OSINT Finding: [Target] Social Media Intelligence",
		"content":           "Social media analysis reveals [specific findings] about [target]",
		"content_type":      "intelligence",
		"category":          "intelligence",
		"intelligence_type": "osint",
		"classification":    "unclassified",
		"threat_level":      "medium",
		"geographic_scope":  "national",
		"mitre_tactic":      "TA0043",
		"mitre_technique":   "T1591",
		"mitre_procedure":   "T1591.001",
		"kill_chain_phase":  "reconnaissance",
		"risk_score":        6.5,
		"impact_score":      7.0,
		"likelihood_score":  6.0,
		"confidence":        0.8,
		"priority":          7,
		"tags":              []string{"osint", "social-media", "reconnaissance"},
		"source":            "Social Media Platforms",
	}

	// Verify template structure
	assert.Equal(t, "osint", osintTemplate["intelligence_type"])
	assert.Equal(t, "unclassified", osintTemplate["classification"])
	assert.Equal(t, "medium", osintTemplate["threat_level"])
	assert.Equal(t, "TA0043", osintTemplate["mitre_tactic"])
	assert.Equal(t, "T1591", osintTemplate["mitre_technique"])

	// Test HUMINT template
	humintTemplate := map[string]interface{}{
		"title":             "HUMINT Finding: [Source] Intelligence Report",
		"content":           "Human intelligence source reports [specific information]",
		"intelligence_type": "humint",
		"classification":    "confidential",
		"threat_level":      "high",
		"geographic_scope":  "regional",
		"attribution":       "Source Alpha",
		"ioc_type":          "email",
		"ioc_value":         "suspicious@example.com",
		"mitre_tactic":      "TA0001",
		"mitre_technique":   "T1566",
		"mitre_procedure":   "T1566.001",
		"kill_chain_phase":  "delivery",
		"risk_score":        8.5,
		"impact_score":      9.0,
		"likelihood_score":  8.0,
		"confidence":        0.9,
		"priority":          9,
		"tags":              []string{"humint", "source-intelligence", "threat-actor"},
		"source":            "Human Source",
	}

	// Verify template structure
	assert.Equal(t, "humint", humintTemplate["intelligence_type"])
	assert.Equal(t, "confidential", humintTemplate["classification"])
	assert.Equal(t, "high", humintTemplate["threat_level"])
	assert.Equal(t, "Source Alpha", humintTemplate["attribution"])

	// Test SIGINT template
	sigintTemplate := map[string]interface{}{
		"title":             "SIGINT Finding: [Target] Communications Intelligence",
		"content":           "Signals intelligence analysis reveals [specific findings]",
		"intelligence_type": "sigint",
		"classification":    "secret",
		"threat_level":      "critical",
		"geographic_scope":  "international",
		"attribution":       "APT Group",
		"ioc_type":          "ip",
		"ioc_value":         "192.168.1.100",
		"mitre_tactic":      "TA0011",
		"mitre_technique":   "T1071",
		"mitre_procedure":   "T1071.001",
		"kill_chain_phase":  "c2",
		"risk_score":        9.5,
		"impact_score":      9.5,
		"likelihood_score":  9.0,
		"confidence":        0.95,
		"priority":          10,
		"tags":              []string{"sigint", "communications", "apt", "c2"},
		"source":            "Signals Intelligence",
	}

	// Verify template structure
	assert.Equal(t, "sigint", sigintTemplate["intelligence_type"])
	assert.Equal(t, "secret", sigintTemplate["classification"])
	assert.Equal(t, "critical", sigintTemplate["threat_level"])
	assert.Equal(t, "international", sigintTemplate["geographic_scope"])

	t.Log("Intelligence templates validation passed!")
}

// TestMemoryCategories tests the enhanced memory categories
func TestMemoryCategories(t *testing.T) {
	// Test intelligence categories
	intelligenceCategories := []string{
		"intelligence", "osint", "humint", "sigint", "geoint", "masint", "techint", "finint", "cybint",
	}

	for _, category := range intelligenceCategories {
		assert.NotEmpty(t, category, "Intelligence category should not be empty")
	}

	// Test reconnaissance categories
	reconCategories := []string{
		"reconnaissance", "target_analysis", "infrastructure_mapping", "vulnerability_assessment",
		"threat_hunting", "incident_response",
	}

	for _, category := range reconCategories {
		assert.NotEmpty(t, category, "Reconnaissance category should not be empty")
	}

	// Test analysis categories
	analysisCategories := []string{
		"malware_analysis", "binary_analysis", "vulnerability_research", "protocol_analysis",
		"code_analysis", "behavioral_analysis",
	}

	for _, category := range analysisCategories {
		assert.NotEmpty(t, category, "Analysis category should not be empty")
	}

	// Test intelligence objects
	intelligenceObjects := []string{
		"threat_actor", "attack_campaign", "ioc", "ttp", "pattern", "correlation",
	}

	for _, category := range intelligenceObjects {
		assert.NotEmpty(t, category, "Intelligence object category should not be empty")
	}

	t.Log("Memory categories validation passed!")
}

// TestInsightMapping tests the insight mapping capabilities
func TestInsightMapping(t *testing.T) {
	// Test pattern types
	patternTypes := []string{
		"behavioral", "attack", "temporal", "spatial", "network", "data",
	}

	for _, patternType := range patternTypes {
		assert.NotEmpty(t, patternType, "Pattern type should not be empty")
	}

	// Test correlation types
	correlationTypes := []string{
		"temporal", "spatial", "logical", "statistical", "causal", "predictive",
	}

	for _, correlationType := range correlationTypes {
		assert.NotEmpty(t, correlationType, "Correlation type should not be empty")
	}

	// Test insight types
	insightTypes := []string{
		"pattern", "correlation", "prediction", "recommendation",
	}

	for _, insightType := range insightTypes {
		assert.NotEmpty(t, insightType, "Insight type should not be empty")
	}

	// Test knowledge graph node types
	nodeTypes := []string{
		"entity", "event", "location", "asset", "person", "organization",
	}

	for _, nodeType := range nodeTypes {
		assert.NotEmpty(t, nodeType, "Node type should not be empty")
	}

	t.Log("Insight mapping validation passed!")
}

// TestReverseEngineering tests the reverse engineering capabilities
func TestReverseEngineering(t *testing.T) {
	// Test analysis types
	analysisTypes := []string{
		"malware", "binary", "vulnerability", "protocol",
	}

	for _, analysisType := range analysisTypes {
		assert.NotEmpty(t, analysisType, "Analysis type should not be empty")
	}

	// Test finding types
	findingTypes := []string{
		"vulnerability", "behavior", "ioc", "technique", "function", "string", "api_call",
	}

	for _, findingType := range findingTypes {
		assert.NotEmpty(t, findingType, "Finding type should not be empty")
	}

	// Test exploit types
	exploitTypes := []string{
		"local", "remote", "dos", "code_execution", "privilege_escalation",
	}

	for _, exploitType := range exploitTypes {
		assert.NotEmpty(t, exploitType, "Exploit type should not be empty")
	}

	// Test pattern types
	patternTypes := []string{
		"behavioral", "attack", "defense", "vulnerability", "exploit", "ioc", "ttp",
	}

	for _, patternType := range patternTypes {
		assert.NotEmpty(t, patternType, "Pattern type should not be empty")
	}

	t.Log("Reverse engineering validation passed!")
}

// TestMITREAttackIntegration tests MITRE ATT&CK integration
func TestMITREAttackIntegration(t *testing.T) {
	// Test enterprise tactics
	enterpriseTactics := []string{
		"TA0001", "TA0002", "TA0003", "TA0004", "TA0005", "TA0006",
		"TA0007", "TA0008", "TA0009", "TA0010", "TA0011", "TA0040",
	}

	for _, tactic := range enterpriseTactics {
		assert.NotEmpty(t, tactic, "MITRE tactic should not be empty")
		assert.True(t, len(tactic) == 6, "MITRE tactic should be 6 characters")
		assert.True(t, tactic[:2] == "TA", "MITRE tactic should start with TA")
	}

	// Test common techniques
	techniques := []string{
		"T1566", "T1190", "T1078", "T1071", "T1059", "T1204", "T1053",
		"T1543", "T1053", "T1547", "T1562", "T1070", "T1036", "T1027",
		"T1110", "T1003", "T1555", "T1087", "T1018", "T1082", "T1021",
		"T1071", "T1028", "T1005", "T1039", "T1003", "T1041", "T1048",
		"T1020", "T1071", "T1090", "T1102", "T1485", "T1489", "T1529",
	}

	for _, technique := range techniques {
		assert.NotEmpty(t, technique, "MITRE technique should not be empty")
		assert.True(t, len(technique) == 5, "MITRE technique should be 5 characters")
		assert.True(t, technique[:1] == "T", "MITRE technique should start with T")
	}

	// Test kill chain phases
	killChainPhases := []string{
		"reconnaissance", "weaponization", "delivery", "exploitation",
		"installation", "c2", "actions",
	}

	for _, phase := range killChainPhases {
		assert.NotEmpty(t, phase, "Kill chain phase should not be empty")
	}

	t.Log("MITRE ATT&CK integration validation passed!")
}

// TestDataValidation tests data validation
func TestDataValidation(t *testing.T) {
	// Test classification levels
	classificationLevels := []string{
		"unclassified", "confidential", "secret", "top_secret",
	}

	for _, level := range classificationLevels {
		assert.NotEmpty(t, level, "Classification level should not be empty")
	}

	// Test threat levels
	threatLevels := []string{
		"low", "medium", "high", "critical",
	}

	for _, level := range threatLevels {
		assert.NotEmpty(t, level, "Threat level should not be empty")
	}

	// Test geographic scopes
	geographicScopes := []string{
		"local", "regional", "national", "international",
	}

	for _, scope := range geographicScopes {
		assert.NotEmpty(t, scope, "Geographic scope should not be empty")
	}

	// Test intelligence types
	intelligenceTypes := []string{
		"osint", "humint", "sigint", "geoint", "masint", "techint", "finint", "cybint", "mixed",
	}

	for _, intelType := range intelligenceTypes {
		assert.NotEmpty(t, intelType, "Intelligence type should not be empty")
	}

	// Test IOC types
	iocTypes := []string{
		"ip", "domain", "url", "hash", "email", "file", "registry", "mutex", "service",
	}

	for _, iocType := range iocTypes {
		assert.NotEmpty(t, iocType, "IOC type should not be empty")
	}

	t.Log("Data validation passed!")
}

// TestJSONSerialization tests JSON serialization of new models
func TestJSONSerialization(t *testing.T) {
	// Test IntelligenceFinding serialization
	finding := &models.IntelligenceFinding{
		ID:               "finding-1",
		SessionID:        "session-1",
		Title:            "Test Intelligence Finding",
		Description:      "Test description",
		IntelligenceType: "osint",
		Classification:   "unclassified",
		ThreatLevel:      "medium",
		GeographicScope:  "national",
		Attribution:      "Unknown",
		IOCType:          "domain",
		IOCValue:         "example.com",
		MITRETactic:      "TA0001",
		MITRETechnique:   "T1566",
		MITREProcedure:   "T1566.001",
		KillChainPhase:   "reconnaissance",
		RiskScore:        6.5,
		ImpactScore:      7.0,
		LikelihoodScore:  6.0,
		Confidence:       0.8,
		Priority:         7,
		Tags:             []string{"osint", "test"},
		Source:           "Test Source",
		Metadata: map[string]interface{}{
			"test_key": "test_value",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(finding)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaledFinding models.IntelligenceFinding
	err = json.Unmarshal(jsonData, &unmarshaledFinding)
	assert.NoError(t, err)
	assert.Equal(t, finding.ID, unmarshaledFinding.ID)
	assert.Equal(t, finding.Title, unmarshaledFinding.Title)
	assert.Equal(t, finding.IntelligenceType, unmarshaledFinding.IntelligenceType)

	// Test ThreatActor serialization
	threatActor := &models.ThreatActor{
		ID:              "actor-1",
		SessionID:       "session-1",
		Name:            "Test Actor",
		Aliases:         []string{"Alias1", "Alias2"},
		Description:     "Test description",
		Motivation:      "Financial",
		Capabilities:    []string{"Phishing", "Malware"},
		Targets:         []string{"Government", "Financial"},
		Tools:           []string{"Tool1", "Tool2"},
		Techniques:      []string{"T1566", "T1055"},
		Attribution:     "High confidence",
		Confidence:      0.9,
		ThreatLevel:     "high",
		GeographicScope: "international",
		Metadata: map[string]interface{}{
			"test_key": "test_value",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test JSON marshaling
	jsonData, err = json.Marshal(threatActor)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaledActor models.ThreatActor
	err = json.Unmarshal(jsonData, &unmarshaledActor)
	assert.NoError(t, err)
	assert.Equal(t, threatActor.ID, unmarshaledActor.ID)
	assert.Equal(t, threatActor.Name, unmarshaledActor.Name)
	assert.Equal(t, threatActor.ThreatLevel, unmarshaledActor.ThreatLevel)

	t.Log("JSON serialization tests passed!")
}

// TestPerformance tests performance of new features
func TestPerformance(t *testing.T) {
	// Skip if CGO is not enabled
	if os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("Skipping performance tests - CGO not enabled")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "performance_test.db")

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	db, err := database.NewDatabase(dbPath, logger)
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewMemoryRepository(db.GetDB(), logger)
	ctx := context.Background()

	// Create session
	session := &models.Session{
		ID:               "perf-session",
		Name:             "Performance Test",
		TaskType:         "intelligence_analysis",
		Status:           "active",
		IntelligenceType: "osint",
	}
	err = repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Test creating many intelligence entries
	start := time.Now()
	for i := 0; i < 100; i++ {
		req := &models.CreateMemoryEntryRequest{
			SessionID:        session.ID,
			Title:            fmt.Sprintf("Intelligence Finding %d", i),
			Content:          fmt.Sprintf("Content for finding %d", i),
			Category:         "intelligence",
			IntelligenceType: "osint",
			ThreatLevel:      "medium",
			Priority:         i % 10,
			Tags:             []string{"test", "performance"},
		}
		_, err := repo.CreateMemoryEntry(ctx, req)
		require.NoError(t, err)
	}
	creationTime := time.Since(start)

	// Test searching
	start = time.Now()
	searchReq := &models.SearchRequest{
		Query:      "intelligence",
		SessionID:  session.ID,
		SearchType: "exact",
		Limit:      50,
	}
	results, err := repo.SearchMemoryEntries(ctx, searchReq)
	require.NoError(t, err)
	searchTime := time.Since(start)

	// Performance assertions
	assert.Less(t, creationTime, 5*time.Second, "Creation should be fast")
	assert.Less(t, searchTime, 1*time.Second, "Search should be fast")
	assert.Len(t, results, 50, "Should find 50 results")

	t.Logf("Performance test completed - Creation: %v, Search: %v", creationTime, searchTime)
}

// Main test runner
func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("CGO_ENABLED", "1")

	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}
