package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Landfill errors.
var (
	ErrLandfillNotFound = errors.New("landfill not found") // Returned when a landfill is not found.
)

// EditableLandfill defines the editable landfill structure.
type EditableLandfill struct {
	GeoJSON GeoJSON
}

// EditableLandfillPatch defines the patchable landfill structure.
type EditableLandfillPatch struct {
	GeoJSON GeoJSON
}

// Landfill defines the landfill structure.
type Landfill struct {
	EditableLandfill
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// LandfillPaginatedSort defines the field of the landfill to sort.
type LandfillPaginatedSort string

const (
	LandfillPaginatedSortWayName          LandfillPaginatedSort = "wayName"
	LandfillPaginatedSortMunicipalityName LandfillPaginatedSort = "municipalityName"
	LandfillPaginatedSortCreatedAt        LandfillPaginatedSort = "createdAt"
	LandfillPaginatedSortModifiedAt       LandfillPaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s LandfillPaginatedSort) Field() LandfillPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s LandfillPaginatedSort) Valid() bool {
	switch s {
	case LandfillPaginatedSortWayName,
		LandfillPaginatedSortMunicipalityName,
		LandfillPaginatedSortCreatedAt,
		LandfillPaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// LandfillsPaginatedFilter defines the landfills filter structure.
type LandfillsPaginatedFilter struct {
	PaginatedRequest[LandfillPaginatedSort]
	LocationName *string
}
