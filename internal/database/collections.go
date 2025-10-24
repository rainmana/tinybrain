package database

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

// CreateSessionsCollection creates the sessions collection with proper schema
func CreateSessionsCollection() *models.Collection {
	collection := &models.Collection{
		Name:       "sessions",
		Type:       models.CollectionTypeBase,
		System:     false,
		CreateRule: nil, // Allow creation
		UpdateRule: nil, // Allow updates
		DeleteRule: nil, // Allow deletion
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "name",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(255),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "task_type",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(100),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "status",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(50),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "description",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(2000),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "metadata",
		Type: schema.FieldTypeJson,
	})

	return collection
}

// CreateMemoryEntriesCollection creates the memory_entries collection
func CreateMemoryEntriesCollection() *models.Collection {
	collection := &models.Collection{
		Name:       "memory_entries",
		Type:       models.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "session_id",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "title",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(500),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "content",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(65535), // TEXT field max length
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "category",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(100),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "priority",
		Type:     schema.FieldTypeNumber,
		Required: true,
		Options: &schema.NumberOptions{
			Min: types.Pointer(1.0),
			Max: types.Pointer(10.0),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "confidence",
		Type:     schema.FieldTypeNumber,
		Required: true,
		Options: &schema.NumberOptions{
			Min: types.Pointer(0.0),
			Max: types.Pointer(1.0),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "tags",
		Type: schema.FieldTypeJson,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "source",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(200),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "content_type",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(50),
		},
	})

	return collection
}

// CreateRelationshipsCollection creates the relationships collection
func CreateRelationshipsCollection() *models.Collection {
	collection := &models.Collection{
		Name:       "relationships",
		Type:       models.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "source_id",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "target_id",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "type",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(50),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "strength",
		Type:     schema.FieldTypeNumber,
		Required: true,
		Options: &schema.NumberOptions{
			Min: types.Pointer(0.0),
			Max: types.Pointer(1.0),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "description",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(1000),
		},
	})

	return collection
}

// CreateContextSnapshotsCollection creates the context_snapshots collection
func CreateContextSnapshotsCollection() *models.Collection {
	collection := &models.Collection{
		Name:       "context_snapshots",
		Type:       models.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "session_id",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "name",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(200),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "description",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(1000),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "context_data",
		Type:     schema.FieldTypeJson,
		Required: true,
	})

	return collection
}

// CreateTaskProgressCollection creates the task_progress collection
func CreateTaskProgressCollection() *models.Collection {
	collection := &models.Collection{
		Name:       "task_progress",
		Type:       models.CollectionTypeBase,
		System:     false,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
	}

	// Add fields
	collection.Schema.AddField(&schema.SchemaField{
		Name:     "session_id",
		Type:     schema.FieldTypeText,
		Required: true,
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "task_name",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(200),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "stage",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(100),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "status",
		Type:     schema.FieldTypeText,
		Required: true,
		Options: &schema.TextOptions{
			Min: intPtr(1),
			Max: intPtr(50),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "progress_percentage",
		Type:     schema.FieldTypeNumber,
		Required: true,
		Options: &schema.NumberOptions{
			Min: types.Pointer(0.0),
			Max: types.Pointer(100.0),
		},
	})

	collection.Schema.AddField(&schema.SchemaField{
		Name: "notes",
		Type: schema.FieldTypeText,
		Options: &schema.TextOptions{
			Max: intPtr(2000),
		},
	})

	return collection
}
