package services

import (
	"context"
	"fmt"

	"tinybrain-v2/internal/models"
	"tinybrain-v2/internal/repository"
)

// SessionServiceV2 provides business logic for session management
type SessionServiceV2 struct {
	sessionRepo *repository.SessionRepositoryV2
}

// NewSessionServiceV2 creates a new SessionService
func NewSessionServiceV2(sessionRepo *repository.SessionRepositoryV2) *SessionServiceV2 {
	return &SessionServiceV2{sessionRepo: sessionRepo}
}

// CreateSession creates a new session
func (s *SessionServiceV2) CreateSession(ctx context.Context, req *models.SessionCreateRequest) (*models.Session, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("session name cannot be empty")
	}
	if req.TaskType == "" {
		return nil, fmt.Errorf("task type cannot be empty")
	}

	return s.sessionRepo.CreateSession(ctx, req)
}

// GetSession retrieves a session by ID
func (s *SessionServiceV2) GetSession(ctx context.Context, id string) (*models.Session, error) {
	if id == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	return s.sessionRepo.GetSession(ctx, id)
}

// ListSessions lists sessions based on criteria
func (s *SessionServiceV2) ListSessions(ctx context.Context, req *models.SessionListRequest) ([]*models.Session, int, error) {
	return s.sessionRepo.ListSessions(ctx, req)
}

// UpdateSession updates an existing session
func (s *SessionServiceV2) UpdateSession(ctx context.Context, id string, req *models.SessionUpdateRequest) (*models.Session, error) {
	if id == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	return s.sessionRepo.UpdateSession(ctx, id, req)
}

// DeleteSession deletes a session by ID
func (s *SessionServiceV2) DeleteSession(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	return s.sessionRepo.DeleteSession(ctx, id)
}
