package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Container errors.
var (
	ErrContainerNotFound = errors.New("container not found") // Returned when a container is not found.
)

// ContainerCategory defines the category of the container.
type ContainerCategory string

const (
	ContainerCategoryGeneral   ContainerCategory = "general"
	ContainerCategoryPaper     ContainerCategory = "paper"
	ContainerCategoryPlastic   ContainerCategory = "plastic"
	ContainerCategoryMetal     ContainerCategory = "metal"
	ContainerCategoryGlass     ContainerCategory = "glass"
	ContainerCategoryOrganic   ContainerCategory = "organic"
	ContainerCategoryHazardous ContainerCategory = "hazardous"
)

// Valid returns true if the category is valid, false otherwise.
func (r ContainerCategory) Valid() bool {
	switch r {
	case ContainerCategoryGeneral,
		ContainerCategoryPaper,
		ContainerCategoryPlastic,
		ContainerCategoryMetal,
		ContainerCategoryGlass,
		ContainerCategoryOrganic,
		ContainerCategoryHazardous:
		return true
	default:
		return false
	}
}

// EditableContainer defines the editable container structure.
type EditableContainer struct {
	Category ContainerCategory
	GeoJSON  GeoJSON
}

// Container defines the container structure.
type Container struct {
	EditableContainer
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}
