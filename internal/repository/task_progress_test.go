package repository

import (
	"context"
	"testing"

	"github.com/alec/tinybrain/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTaskProgress(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session first
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create task progress
	progress, err := repo.CreateTaskProgress(ctx, session.ID, "Vulnerability Assessment", 
		"Initial Discovery", "in_progress", "Starting assessment", 25)
	require.NoError(t, err)

	assert.NotEmpty(t, progress.ID)
	assert.Equal(t, session.ID, progress.SessionID)
	assert.Equal(t, "Vulnerability Assessment", progress.TaskName)
	assert.Equal(t, "Initial Discovery", progress.Stage)
	assert.Equal(t, "in_progress", progress.Status)
	assert.Equal(t, "Starting assessment", progress.Notes)
	assert.Equal(t, 25, progress.ProgressPercentage)
	assert.NotNil(t, progress.StartedAt) // Should be set for in_progress status
}

func TestUpdateTaskProgress(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create task progress
	progress, err := repo.CreateTaskProgress(ctx, session.ID, "Vulnerability Assessment", 
		"Initial Discovery", "pending", "Starting assessment", 0)
	require.NoError(t, err)

	// Update task progress
	updated, err := repo.UpdateTaskProgress(ctx, progress.ID, "Validation", "in_progress", 
		"Validating findings", 50)
	require.NoError(t, err)

	assert.Equal(t, progress.ID, updated.ID)
	assert.Equal(t, "Validation", updated.Stage)
	assert.Equal(t, "in_progress", updated.Status)
	assert.Equal(t, "Validating findings", updated.Notes)
	assert.Equal(t, 50, updated.ProgressPercentage)
	assert.NotNil(t, updated.StartedAt) // Should be set when transitioning to in_progress

	// Update to completed
	completed, err := repo.UpdateTaskProgress(ctx, progress.ID, "Completed", "completed", 
		"Assessment complete", 100)
	require.NoError(t, err)

	assert.Equal(t, "completed", completed.Status)
	assert.Equal(t, 100, completed.ProgressPercentage)
	assert.NotNil(t, completed.CompletedAt) // Should be set when transitioning to completed
}

func TestGetTaskProgress(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create task progress
	progress, err := repo.CreateTaskProgress(ctx, session.ID, "Vulnerability Assessment", 
		"Initial Discovery", "in_progress", "Starting assessment", 25)
	require.NoError(t, err)

	// Retrieve the task progress
	retrieved, err := repo.GetTaskProgress(ctx, progress.ID)
	require.NoError(t, err)

	assert.Equal(t, progress.ID, retrieved.ID)
	assert.Equal(t, progress.SessionID, retrieved.SessionID)
	assert.Equal(t, progress.TaskName, retrieved.TaskName)
	assert.Equal(t, progress.Stage, retrieved.Stage)
	assert.Equal(t, progress.Status, retrieved.Status)
	assert.Equal(t, progress.Notes, retrieved.Notes)
	assert.Equal(t, progress.ProgressPercentage, retrieved.ProgressPercentage)
}

func TestListTaskProgress(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create multiple task progress entries
	tasks := []struct {
		name   string
		stage  string
		status string
		notes  string
	}{
		{"Task 1", "Stage 1", "completed", "Task 1 completed"},
		{"Task 2", "Stage 2", "in_progress", "Task 2 in progress"},
		{"Task 3", "Stage 3", "pending", "Task 3 pending"},
	}

	for _, task := range tasks {
		_, err := repo.CreateTaskProgress(ctx, session.ID, task.name, task.stage, 
			task.status, task.notes, 0)
		require.NoError(t, err)
	}

	// List all tasks
	allTasks, err := repo.ListTaskProgress(ctx, session.ID, "", 10, 0)
	require.NoError(t, err)
	assert.Len(t, allTasks, 3)

	// List only in_progress tasks
	inProgressTasks, err := repo.ListTaskProgress(ctx, session.ID, "in_progress", 10, 0)
	require.NoError(t, err)
	assert.Len(t, inProgressTasks, 1)
	assert.Equal(t, "Task 2", inProgressTasks[0].TaskName)

	// List only completed tasks
	completedTasks, err := repo.ListTaskProgress(ctx, session.ID, "completed", 10, 0)
	require.NoError(t, err)
	assert.Len(t, completedTasks, 1)
	assert.Equal(t, "Task 1", completedTasks[0].TaskName)
}

func TestTaskProgressStatusTransitions(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a session
	session := &models.Session{
		ID:       "test-session-1",
		Name:     "Test Session",
		TaskType: "security_review",
		Status:   "active",
	}
	err := repo.CreateSession(ctx, session)
	require.NoError(t, err)

	// Create task in pending status
	progress, err := repo.CreateTaskProgress(ctx, session.ID, "Test Task", 
		"Initial", "pending", "Created", 0)
	require.NoError(t, err)
	assert.Nil(t, progress.StartedAt)
	assert.Nil(t, progress.CompletedAt)

	// Transition to in_progress
	progress, err = repo.UpdateTaskProgress(ctx, progress.ID, "Working", "in_progress", 
		"Started working", 25)
	require.NoError(t, err)
	assert.NotNil(t, progress.StartedAt)
	assert.Nil(t, progress.CompletedAt)

	// Transition to completed
	progress, err = repo.UpdateTaskProgress(ctx, progress.ID, "Done", "completed", 
		"Finished", 100)
	require.NoError(t, err)
	assert.NotNil(t, progress.StartedAt)
	assert.NotNil(t, progress.CompletedAt)
}
