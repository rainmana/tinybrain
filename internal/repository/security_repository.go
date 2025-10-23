package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/database"
	"github.com/rainmana/tinybrain/internal/models"
)

// SecurityRepositoryInterface defines the interface for security data operations
type SecurityRepositoryInterface interface {
	StoreNVDDataset(ctx context.Context, cves []models.NVDCVE) error
	StoreATTACKDataset(ctx context.Context, techniques []models.ATTACKTechnique, tactics []models.ATTACKTactic) error
	QueryNVD(ctx context.Context, req models.NVDSearchRequest) ([]models.NVDCVE, int, error)
	QueryATTACK(ctx context.Context, req models.ATTACKSearchRequest) ([]models.ATTACKTechnique, int, error)
	GetSecurityDataSummary(ctx context.Context) (map[string]models.SecurityDataSummary, error)
	UpdateSecurityDataStatus(ctx context.Context, dataSource string, status string, totalRecords *int, errorMessage *string) error
}

// SecurityRepository handles security data operations
type SecurityRepository struct {
	db     *database.Database
	logger *log.Logger
}

// NewSecurityRepository creates a new security repository
func NewSecurityRepository(db *database.Database, logger *log.Logger) *SecurityRepository {
	return &SecurityRepository{
		db:     db,
		logger: logger,
	}
}

// StoreNVDDataset stores NVD CVE data in the database
func (r *SecurityRepository) StoreNVDDataset(ctx context.Context, cves []models.NVDCVE) error {
	r.logger.Info("Storing NVD dataset", "count", len(cves))

	tx, err := r.db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Clear existing data
	if _, err := tx.ExecContext(ctx, "DELETE FROM nvd_cves"); err != nil {
		return fmt.Errorf("failed to clear existing NVD data: %v", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO nvd_cves (
			id, description, cvss_v2_score, cvss_v2_vector, cvss_v3_score, cvss_v3_vector,
			severity, published_date, last_modified_date, cwe_ids, affected_products,
			refs, raw_data, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare NVD insert statement: %v", err)
	}
	defer stmt.Close()

	for _, cve := range cves {
		cweIDsJSON, _ := json.Marshal(cve.CWEIDs)
		productsJSON, _ := json.Marshal(cve.AffectedProducts)
		referencesJSON, _ := json.Marshal(cve.References)

		_, err := stmt.ExecContext(ctx,
			cve.ID, cve.Description, cve.CVSSV2Score, cve.CVSSV2Vector,
			cve.CVSSV3Score, cve.CVSSV3Vector, cve.Severity,
			cve.PublishedDate, cve.LastModifiedDate,
			string(cweIDsJSON), string(productsJSON), string(referencesJSON),
			cve.RawData, cve.CreatedAt, cve.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert CVE %s: %v", cve.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit NVD transaction: %v", err)
	}

	r.logger.Info("NVD dataset stored successfully", "count", len(cves))
	return nil
}

// StoreATTACKDataset stores MITRE ATT&CK data in the database
func (r *SecurityRepository) StoreATTACKDataset(ctx context.Context, techniques []models.ATTACKTechnique, tactics []models.ATTACKTactic) error {
	r.logger.Info("Storing ATT&CK dataset", "techniques", len(techniques), "tactics", len(tactics))

	tx, err := r.db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Clear existing data
	if _, err := tx.ExecContext(ctx, "DELETE FROM attack_techniques"); err != nil {
		return fmt.Errorf("failed to clear existing ATT&CK techniques: %v", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM attack_tactics"); err != nil {
		return fmt.Errorf("failed to clear existing ATT&CK tactics: %v", err)
	}

	// Store tactics first
	tacticStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO attack_tactics (
			id, name, description, external_id, kill_chain_phases, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare tactic insert statement: %v", err)
	}
	defer tacticStmt.Close()

	for _, tactic := range tactics {
		phasesJSON, _ := json.Marshal(tactic.KillChainPhases)
		_, err := tacticStmt.ExecContext(ctx,
			tactic.ID, tactic.Name, tactic.Description, tactic.ExternalID,
			string(phasesJSON), tactic.CreatedAt, tactic.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert tactic %s: %v", tactic.ID, err)
		}
	}

	// Store techniques
	techniqueStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO attack_techniques (
			id, name, description, tactic, tactics, platforms, kill_chain_phases,
			data_sources, detection, mitigation, refs, sub_techniques,
			parent_technique, raw_data, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare technique insert statement: %v", err)
	}
	defer techniqueStmt.Close()

	for _, technique := range techniques {
		tacticsJSON, _ := json.Marshal(technique.Tactics)
		platformsJSON, _ := json.Marshal(technique.Platforms)
		phasesJSON, _ := json.Marshal(technique.KillChainPhases)
		sourcesJSON, _ := json.Marshal(technique.DataSources)
		referencesJSON, _ := json.Marshal(technique.References)
		subTechniquesJSON, _ := json.Marshal(technique.SubTechniques)

		_, err := techniqueStmt.ExecContext(ctx,
			technique.ID, technique.Name, technique.Description, technique.Tactic,
			string(tacticsJSON), string(platformsJSON), string(phasesJSON),
			string(sourcesJSON), technique.Detection, technique.Mitigation,
			string(referencesJSON), string(subTechniquesJSON), technique.ParentTechnique,
			technique.RawData, technique.CreatedAt, technique.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert technique %s: %v", technique.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit ATT&CK transaction: %v", err)
	}

	r.logger.Info("ATT&CK dataset stored successfully", "techniques", len(techniques), "tactics", len(tactics))
	return nil
}

// QueryNVD searches NVD data based on criteria
func (r *SecurityRepository) QueryNVD(ctx context.Context, req models.NVDSearchRequest) ([]models.NVDCVE, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause
	if req.CWEID != nil {
		conditions = append(conditions, "cwe_ids LIKE ?")
		args = append(args, "%"+*req.CWEID+"%")
	}

	if req.Component != nil {
		conditions = append(conditions, "affected_products LIKE ?")
		args = append(args, "%"+*req.Component+"%")
	}

	if req.Severity != nil {
		conditions = append(conditions, "severity = ?")
		args = append(args, *req.Severity)
	}

	if req.MinCVSS != nil {
		conditions = append(conditions, "cvss_v3_score >= ?")
		args = append(args, *req.MinCVSS)
	}

	if req.MaxCVSS != nil {
		conditions = append(conditions, "cvss_v3_score <= ?")
		args = append(args, *req.MaxCVSS)
	}

	if req.PublishedAfter != nil {
		conditions = append(conditions, "published_date >= ?")
		args = append(args, *req.PublishedAfter)
	}

	if req.PublishedBefore != nil {
		conditions = append(conditions, "published_date <= ?")
		args = append(args, *req.PublishedBefore)
	}

	// Build query
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total results
	countQuery := "SELECT COUNT(*) FROM nvd_cves " + whereClause
	var totalCount int
	err := r.db.GetDB().QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count NVD results: %v", err)
	}

	// Get results with pagination
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, description, cvss_v2_score, cvss_v2_vector, cvss_v3_score, cvss_v3_vector,
		       severity, published_date, last_modified_date, cwe_ids, affected_products,
		       refs, raw_data, created_at, updated_at
		FROM nvd_cves %s
		ORDER BY cvss_v3_score DESC, published_date DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)

	rows, err := r.db.GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query NVD data: %v", err)
	}
	defer rows.Close()

	var cves []models.NVDCVE
	for rows.Next() {
		var cve models.NVDCVE
		var cweIDsJSON, productsJSON, referencesJSON string

		err := rows.Scan(
			&cve.ID, &cve.Description, &cve.CVSSV2Score, &cve.CVSSV2Vector,
			&cve.CVSSV3Score, &cve.CVSSV3Vector, &cve.Severity,
			&cve.PublishedDate, &cve.LastModifiedDate,
			&cweIDsJSON, &productsJSON, &referencesJSON,
			&cve.RawData, &cve.CreatedAt, &cve.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan NVD row: %v", err)
		}

		// Parse JSON fields
		json.Unmarshal([]byte(cweIDsJSON), &cve.CWEIDs)
		json.Unmarshal([]byte(productsJSON), &cve.AffectedProducts)
		json.Unmarshal([]byte(referencesJSON), &cve.References)

		cves = append(cves, cve)
	}

	return cves, totalCount, nil
}

// QueryATTACK searches ATT&CK data based on criteria
func (r *SecurityRepository) QueryATTACK(ctx context.Context, req models.ATTACKSearchRequest) ([]models.ATTACKTechnique, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause
	if req.TechniqueID != nil {
		conditions = append(conditions, "id = ?")
		args = append(args, *req.TechniqueID)
	}

	if req.Tactic != nil {
		conditions = append(conditions, "tactic = ?")
		args = append(args, *req.Tactic)
	}

	if req.Platform != nil {
		conditions = append(conditions, "platforms LIKE ?")
		args = append(args, "%"+*req.Platform+"%")
	}

	if req.Query != nil {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		args = append(args, "%"+*req.Query+"%", "%"+*req.Query+"%")
	}

	// Build query
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total results
	countQuery := "SELECT COUNT(*) FROM attack_techniques " + whereClause
	var totalCount int
	err := r.db.GetDB().QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ATT&CK results: %v", err)
	}

	// Get results with pagination
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, tactic, tactics, platforms, kill_chain_phases,
		       data_sources, detection, mitigation, refs, sub_techniques,
		       parent_technique, raw_data, created_at, updated_at
		FROM attack_techniques %s
		ORDER BY name
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)

	rows, err := r.db.GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query ATT&CK data: %v", err)
	}
	defer rows.Close()

	var techniques []models.ATTACKTechnique
	for rows.Next() {
		var technique models.ATTACKTechnique
		var tacticsJSON, platformsJSON, phasesJSON, sourcesJSON, referencesJSON, subTechniquesJSON string

		err := rows.Scan(
			&technique.ID, &technique.Name, &technique.Description, &technique.Tactic,
			&tacticsJSON, &platformsJSON, &phasesJSON, &sourcesJSON,
			&technique.Detection, &technique.Mitigation, &referencesJSON, &subTechniquesJSON,
			&technique.ParentTechnique, &technique.RawData, &technique.CreatedAt, &technique.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan ATT&CK row: %v", err)
		}

		// Parse JSON fields
		json.Unmarshal([]byte(tacticsJSON), &technique.Tactics)
		json.Unmarshal([]byte(platformsJSON), &technique.Platforms)
		json.Unmarshal([]byte(phasesJSON), &technique.KillChainPhases)
		json.Unmarshal([]byte(sourcesJSON), &technique.DataSources)
		json.Unmarshal([]byte(referencesJSON), &technique.References)
		json.Unmarshal([]byte(subTechniquesJSON), &technique.SubTechniques)

		techniques = append(techniques, technique)
	}

	return techniques, totalCount, nil
}

// GetSecurityDataSummary returns a summary of security data
func (r *SecurityRepository) GetSecurityDataSummary(ctx context.Context) (map[string]models.SecurityDataSummary, error) {
	summaries := make(map[string]models.SecurityDataSummary)

	// NVD Summary
	var nvdCount int
	var nvdLastUpdate sql.NullTime
	err := r.db.GetDB().QueryRowContext(ctx, "SELECT COUNT(*), MAX(updated_at) FROM nvd_cves").Scan(&nvdCount, &nvdLastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to get NVD summary: %v", err)
	}

	summaries["nvd"] = models.SecurityDataSummary{
		DataSource:   "nvd",
		TotalRecords: nvdCount,
		LastUpdate:   &nvdLastUpdate.Time,
		Summary:      fmt.Sprintf("NVD database contains %d CVE entries", nvdCount),
	}

	// ATT&CK Summary
	var attackCount int
	var attackLastUpdate sql.NullTime
	err = r.db.GetDB().QueryRowContext(ctx, "SELECT COUNT(*), MAX(updated_at) FROM attack_techniques").Scan(&attackCount, &attackLastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to get ATT&CK summary: %v", err)
	}

	summaries["attack"] = models.SecurityDataSummary{
		DataSource:   "attack",
		TotalRecords: attackCount,
		LastUpdate:   &attackLastUpdate.Time,
		Summary:      fmt.Sprintf("MITRE ATT&CK database contains %d techniques", attackCount),
	}

	return summaries, nil
}

// UpdateSecurityDataStatus updates the status of security data updates
func (r *SecurityRepository) UpdateSecurityDataStatus(ctx context.Context, dataSource string, status string, totalRecords *int, errorMessage *string) error {
	now := time.Now()

	query := `
		INSERT INTO security_data_updates (id, data_source, last_update, total_records, update_status, error_message, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			last_update = excluded.last_update,
			total_records = excluded.total_records,
			update_status = excluded.update_status,
			error_message = excluded.error_message,
			updated_at = excluded.updated_at
	`

	_, err := r.db.GetDB().ExecContext(ctx, query,
		dataSource, dataSource, now, totalRecords, status, errorMessage, now, now,
	)

	return err
}
