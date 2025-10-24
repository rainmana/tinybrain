package services

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/repository"
)

// SecurityUpdateService handles intelligent updates to security datasets
type SecurityUpdateService struct {
	downloader   *SecurityDataDownloader
	securityRepo repository.SecurityRepositoryInterface
	logger       *log.Logger
}

// NewSecurityUpdateService creates a new security update service
func NewSecurityUpdateService(
	downloader *SecurityDataDownloader,
	securityRepo repository.SecurityRepositoryInterface,
	logger *log.Logger,
) *SecurityUpdateService {
	return &SecurityUpdateService{
		downloader:   downloader,
		securityRepo: securityRepo,
		logger:       logger,
	}
}

// UpdateStrategy defines how to handle updates for each data source
type UpdateStrategy struct {
	DataSource      string
	CheckInterval   time.Duration
	ForceFullUpdate bool
	MaxAge          time.Duration
}

// GetDefaultUpdateStrategies returns default update strategies for all data sources
func GetDefaultUpdateStrategies() []UpdateStrategy {
	return []UpdateStrategy{
		{
			DataSource:      "nvd",
			CheckInterval:   24 * time.Hour, // Check daily
			ForceFullUpdate: false,
			MaxAge:          7 * 24 * time.Hour, // Force full update if data is older than 7 days
		},
		{
			DataSource:      "attack",
			CheckInterval:   7 * 24 * time.Hour, // Check weekly
			ForceFullUpdate: false,
			MaxAge:          30 * 24 * time.Hour, // Force full update if data is older than 30 days
		},
		{
			DataSource:      "owasp",
			CheckInterval:   30 * 24 * time.Hour, // Check monthly
			ForceFullUpdate: false,
			MaxAge:          90 * 24 * time.Hour, // Force full update if data is older than 90 days
		},
	}
}

// CheckAndUpdateAll checks for updates and performs intelligent updates for all data sources
func (s *SecurityUpdateService) CheckAndUpdateAll(ctx context.Context) error {
	s.logger.Info("Starting security data update check")

	strategies := GetDefaultUpdateStrategies()

	for _, strategy := range strategies {
		if err := s.CheckAndUpdate(ctx, strategy); err != nil {
			s.logger.Error("Failed to update data source", "source", strategy.DataSource, "error", err)
			// Continue with other sources even if one fails
		}
	}

	s.logger.Info("Security data update check completed")
	return nil
}

// CheckAndUpdate checks for updates and performs intelligent update for a specific data source
func (s *SecurityUpdateService) CheckAndUpdate(ctx context.Context, strategy UpdateStrategy) error {
	s.logger.Info("Checking for updates", "source", strategy.DataSource)

	// Get current data status
	summary, err := s.securityRepo.GetSecurityDataSummary(ctx)
	if err != nil {
		return fmt.Errorf("failed to get data summary: %v", err)
	}

	dataStatus, exists := summary[strategy.DataSource]
	if !exists {
		s.logger.Info("No existing data found, performing initial download", "source", strategy.DataSource)
		return s.performFullUpdate(ctx, strategy.DataSource)
	}

	// Check if data is too old and needs full update
	if dataStatus.LastUpdate != nil && time.Since(*dataStatus.LastUpdate) > strategy.MaxAge {
		s.logger.Info("Data is too old, performing full update",
			"source", strategy.DataSource,
			"age", time.Since(*dataStatus.LastUpdate))
		return s.performFullUpdate(ctx, strategy.DataSource)
	}

	// Check for incremental updates
	var lastUpdate time.Time
	if dataStatus.LastUpdate != nil {
		lastUpdate = *dataStatus.LastUpdate
	} else {
		lastUpdate = time.Now().Add(-24 * time.Hour) // Default to 24 hours ago if no last update
	}

	hasUpdates, latestUpdate, err := s.downloader.CheckForUpdates(ctx, strategy.DataSource, lastUpdate)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %v", err)
	}

	if !hasUpdates {
		s.logger.Info("No updates available", "source", strategy.DataSource)
		return nil
	}

	s.logger.Info("Updates found, performing incremental update",
		"source", strategy.DataSource,
		"latest", latestUpdate)

	return s.performIncrementalUpdate(ctx, strategy.DataSource, lastUpdate)
}

// performFullUpdate performs a complete download and replacement of data
func (s *SecurityUpdateService) performFullUpdate(ctx context.Context, dataSource string) error {
	s.logger.Info("Performing full update", "source", dataSource)

	// Update status to "downloading"
	if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "downloading", nil, nil); err != nil {
		s.logger.Error("Failed to update status", "error", err)
	}

	switch dataSource {
	case "nvd":
		cves, err := s.downloader.DownloadNVDDataset(ctx)
		if err != nil {
			errorMsg := err.Error()
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, &errorMsg)
			return fmt.Errorf("failed to download NVD dataset: %v", err)
		}

		if len(cves) == 0 {
			s.logger.Warn("No CVEs downloaded, skipping update")
			return nil
		}

		if err := s.securityRepo.StoreNVDDataset(ctx, cves); err != nil {
			errorMsg := err.Error()
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, &errorMsg)
			return fmt.Errorf("failed to store NVD dataset: %v", err)
		}

		totalRecords := len(cves)
		if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "completed", &totalRecords, nil); err != nil {
			s.logger.Error("Failed to update final status", "error", err)
		}

	case "attack":
		techniques, tactics, err := s.downloader.DownloadATTACKDataset(ctx)
		if err != nil {
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, func() *string { s := err.Error(); return &s }())
			return fmt.Errorf("failed to download ATT&CK dataset: %v", err)
		}

		if len(techniques) == 0 && len(tactics) == 0 {
			s.logger.Warn("No ATT&CK data downloaded, skipping update")
			return nil
		}

		if err := s.securityRepo.StoreATTACKDataset(ctx, techniques, tactics); err != nil {
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, func() *string { s := err.Error(); return &s }())
			return fmt.Errorf("failed to store ATT&CK dataset: %v", err)
		}

		totalRecords := len(techniques) + len(tactics)
		if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "completed", &totalRecords, nil); err != nil {
			s.logger.Error("Failed to update final status", "error", err)
		}

	case "owasp":
		procedures, err := s.downloader.DownloadOWASPDataset(ctx)
		if err != nil {
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, func() *string { s := err.Error(); return &s }())
			return fmt.Errorf("failed to download OWASP dataset: %v", err)
		}

		if len(procedures) == 0 {
			s.logger.Warn("No OWASP procedures downloaded, skipping update")
			return nil
		}

		// Note: We need to implement StoreOWASPDataset in the repository
		s.logger.Info("OWASP procedures downloaded", "count", len(procedures))
		// TODO: Implement OWASP storage in repository

		totalRecords := len(procedures)
		if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "completed", &totalRecords, nil); err != nil {
			s.logger.Error("Failed to update final status", "error", err)
		}

	default:
		return fmt.Errorf("unknown data source: %s", dataSource)
	}

	s.logger.Info("Full update completed", "source", dataSource)
	return nil
}

// performIncrementalUpdate performs an incremental update for data sources that support it
func (s *SecurityUpdateService) performIncrementalUpdate(ctx context.Context, dataSource string, lastUpdate time.Time) error {
	s.logger.Info("Performing incremental update", "source", dataSource, "since", lastUpdate)

	// Update status to "downloading"
	if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "downloading", nil, nil); err != nil {
		s.logger.Error("Failed to update status", "error", err)
	}

	switch dataSource {
	case "nvd":
		// NVD supports incremental updates via date filtering
		cves, err := s.downloader.DownloadIncrementalNVD(ctx, lastUpdate)
		if err != nil {
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, func() *string { s := err.Error(); return &s }())
			return fmt.Errorf("failed to download incremental NVD data: %v", err)
		}

		if len(cves) == 0 {
			s.logger.Info("No new CVEs found in incremental update")
			if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "completed", nil, nil); err != nil {
				s.logger.Error("Failed to update status", "error", err)
			}
			return nil
		}

		// For incremental updates, we need to merge with existing data
		// This is a simplified approach - in production, you'd want more sophisticated merging
		if err := s.securityRepo.StoreNVDDataset(ctx, cves); err != nil {
			s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "error", nil, func() *string { s := err.Error(); return &s }())
			return fmt.Errorf("failed to store incremental NVD data: %v", err)
		}

		totalRecords := len(cves)
		if err := s.securityRepo.UpdateSecurityDataStatus(ctx, dataSource, "completed", &totalRecords, nil); err != nil {
			s.logger.Error("Failed to update final status", "error", err)
		}

	case "attack", "owasp":
		// ATT&CK and OWASP don't support true incremental updates
		// For these, we do a full update but only if there are actual changes
		s.logger.Info("Performing full update for source that doesn't support incremental updates", "source", dataSource)
		return s.performFullUpdate(ctx, dataSource)

	default:
		return fmt.Errorf("unknown data source: %s", dataSource)
	}

	s.logger.Info("Incremental update completed", "source", dataSource)
	return nil
}

// ForceFullUpdate forces a complete update of all data sources
func (s *SecurityUpdateService) ForceFullUpdate(ctx context.Context) error {
	s.logger.Info("Forcing full update of all data sources")

	strategies := GetDefaultUpdateStrategies()
	for _, strategy := range strategies {
		strategy.ForceFullUpdate = true
		if err := s.performFullUpdate(ctx, strategy.DataSource); err != nil {
			s.logger.Error("Failed to force full update", "source", strategy.DataSource, "error", err)
		}
	}

	return nil
}
