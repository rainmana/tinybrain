package models

import "time"

// RelationshipType represents the type of relationship between memories
type RelationshipType string

const (
	RelationshipTypeDependsOn   RelationshipType = "depends_on"
	RelationshipTypeCauses      RelationshipType = "causes"
	RelationshipTypeMitigates   RelationshipType = "mitigates"
	RelationshipTypeExploits    RelationshipType = "exploits"
	RelationshipTypeReferences  RelationshipType = "references"
	RelationshipTypeContradicts RelationshipType = "contradicts"
	RelationshipTypeSupports    RelationshipType = "supports"
	RelationshipTypeRelatedTo   RelationshipType = "related_to"
	RelationshipTypeParentOf    RelationshipType = "parent_of"
	RelationshipTypeChildOf     RelationshipType = "child_of"
)

// Relationship represents a relationship between two memories
type Relationship struct {
	ID          string           `json:"id"`
	SourceID    string           `json:"source_id"`
	TargetID    string           `json:"target_id"`
	Type        RelationshipType `json:"type"`
	Strength    float32          `json:"strength"` // 0.0-1.0 scale
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// RelationshipCreateRequest defines the structure for creating a new relationship
type RelationshipCreateRequest struct {
	SourceID    string           `json:"source_id"`
	TargetID    string           `json:"target_id"`
	Type        RelationshipType `json:"type"`
	Strength    float32          `json:"strength"`
	Description string           `json:"description,omitempty"`
}

// RelationshipUpdateRequest defines the structure for updating an existing relationship
type RelationshipUpdateRequest struct {
	Type        *RelationshipType `json:"type,omitempty"`
	Strength    *float32          `json:"strength,omitempty"`
	Description *string           `json:"description,omitempty"`
}

// RelationshipListRequest defines the structure for listing relationships
type RelationshipListRequest struct {
	SourceID string `json:"source_id,omitempty"`
	TargetID string `json:"target_id,omitempty"`
	Type     string `json:"type,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}
