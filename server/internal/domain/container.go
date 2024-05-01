package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Container errors.
var (
	ErrContainerNotFound                            = errors.New("container not found")                     // Returned when a container is not found.
	ErrContainerAssociatedWithContainerReport       = errors.New("container associated with user report")   // Returned when a container is associated with a user report.
	ErrContainerAssociatedWithUserContainerBookmark = errors.New("container associated with user bookmark") // Returned when a container is associated with a user bookmark.
	ErrContainerAssociatedWithRouteContainer        = errors.New("container associated with route")         // Returned when a container is associated with a route.
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
func (c ContainerCategory) Valid() bool {
	switch c {
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

// EditableContainerPatch defines the patchable container structure.
type EditableContainerPatch struct {
	Category *ContainerCategory
	GeoJSON  GeoJSON
}

// Container defines the container structure.
type Container struct {
	EditableContainer
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// ContainerPaginatedSort defines the field of the container to sort.
type ContainerPaginatedSort string

const (
	ContainerPaginatedSortCategory         ContainerPaginatedSort = "category"
	ContainerPaginatedSortWayName          ContainerPaginatedSort = "wayName"
	ContainerPaginatedSortMunicipalityName ContainerPaginatedSort = "municipalityName"
	ContainerPaginatedSortCreatedAt        ContainerPaginatedSort = "createdAt"
	ContainerPaginatedSortModifiedAt       ContainerPaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s ContainerPaginatedSort) Field() ContainerPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s ContainerPaginatedSort) Valid() bool {
	switch s {
	case ContainerPaginatedSortCategory,
		ContainerPaginatedSortWayName,
		ContainerPaginatedSortMunicipalityName,
		ContainerPaginatedSortCreatedAt,
		ContainerPaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// ContainersPaginatedFilter defines the containers filter structure.
type ContainersPaginatedFilter struct {
	PaginatedRequest[ContainerPaginatedSort]
	Category         *ContainerCategory
	WayName          *string
	MunicipalityName *string
}
