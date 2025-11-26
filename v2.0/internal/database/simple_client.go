package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

// SimplePocketBaseClient wraps PocketBase with minimal configuration
type SimplePocketBaseClient struct {
	app *pocketbase.PocketBase
}

// NewSimplePocketBaseClient creates a new simple PocketBase client
func NewSimplePocketBaseClient(dataDir string) (*SimplePocketBaseClient, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize PocketBase with config
	config := pocketbase.Config{
		DefaultDataDir: dataDir,
	}
	app := pocketbase.NewWithConfig(config)

	// Bootstrap the app to ensure DB connections are open
	if err := app.Bootstrap(); err != nil {
		return nil, fmt.Errorf("failed to bootstrap PocketBase app: %w", err)
	}

	return &SimplePocketBaseClient{
		app: app,
	}, nil
}

// GetApp returns the underlying PocketBase app instance
func (c *SimplePocketBaseClient) GetApp() *pocketbase.PocketBase {
	return c.app
}

// Bootstrap initializes the database with required collections
func (c *SimplePocketBaseClient) Bootstrap(ctx context.Context) error {
	log.Println("Bootstrapping PocketBase database with collections...")

	// Create collections
	collections := []*models.Collection{
		CreateSessionsCollection(),
		CreateMemoryEntriesCollection(),
		CreateRelationshipsCollection(),
		CreateContextSnapshotsCollection(),
		CreateTaskProgressCollection(),
	}

	for _, collection := range collections {
		// Check if collection already exists
		existing, err := c.app.Dao().FindCollectionByNameOrId(collection.Name)
		if err != nil {
			// Collection doesn't exist, create it
			if err := c.app.Dao().SaveCollection(collection); err != nil {
				return fmt.Errorf("failed to create collection %s: %w", collection.Name, err)
			}
			log.Printf("Created collection: %s", collection.Name)
		} else {
			log.Printf("Collection '%s' already exists", existing.Name)
		}
	}

	log.Println("Database bootstrap completed successfully")
	return nil
}

// Close gracefully shuts down the PocketBase client
func (c *SimplePocketBaseClient) Close() error {
	// PocketBase doesn't have a direct Close method on the app instance.
	// The underlying DB connections are managed by the app's Dao.
	// For a clean shutdown, typically app.Start() handles this on process exit.
	// If explicit DB closing is needed, it would be on app.Dao().DB().Close()
	// For now, we'll just log a message.
	log.Println("PocketBase client close called (note: PocketBase app manages its own lifecycle)")
	return nil
}
