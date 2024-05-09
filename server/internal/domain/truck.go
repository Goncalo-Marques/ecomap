package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Truck constraints.
const (
	truckMakeMaxLength          = 50
	truckModelMaxLength         = 50
	truckLicensePlateMaxLength  = 30
	truckPersonCapacityMinValue = 1
)

// Truck errors.
var (
	ErrTruckNotFound                     = errors.New("truck not found")                 // Returned when a truck is not found.
	ErrTruckAssociatedWithWarehouseTruck = errors.New("truck associated with warehouse") // Returned when a truck is associated with a warehouse.
	ErrTruckAssociatedWithRoute          = errors.New("truck associated with route")     // Returned when a truck is associated with a route.
)

// TruckMake defines the make of the truck.
type TruckMake string

// Valid returns true if the make is valid, false otherwise.
func (m TruckMake) Valid() bool {
	return len(m) <= truckMakeMaxLength
}

// TruckModel defines the model of the truck.
type TruckModel string

// Valid returns true if the model is valid, false otherwise.
func (m TruckModel) Valid() bool {
	return len(m) <= truckModelMaxLength
}

// TruckLicensePlate defines the license plate of the truck.
type TruckLicensePlate string

// Valid returns true if the license plate is valid, false otherwise.
func (lp TruckLicensePlate) Valid() bool {
	return len(lp) <= truckLicensePlateMaxLength
}

// TruckPersonCapacity defines the person capacity of the truck.
type TruckPersonCapacity int

// Valid returns true if the person capacity is valid, false otherwise.
func (pc TruckPersonCapacity) Valid() bool {
	return pc >= truckPersonCapacityMinValue
}

// EditableTruck defines the editable truck structure.
type EditableTruck struct {
	Make           TruckMake
	Model          TruckModel
	LicensePlate   TruckLicensePlate
	PersonCapacity TruckPersonCapacity
	GeoJSON        GeoJSON
}

// EditableTruckPatch defines the patchable truck structure.
type EditableTruckPatch struct {
	Make           *TruckMake
	Model          *TruckModel
	LicensePlate   *TruckLicensePlate
	PersonCapacity *TruckPersonCapacity
	GeoJSON        GeoJSON
}

// Truck defines the truck structure.
type Truck struct {
	EditableTruck
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// TruckPaginatedSort defines the field of the truck to sort.
type TruckPaginatedSort string

const (
	TruckPaginatedSortMake             TruckPaginatedSort = "make"
	TruckPaginatedSortModel            TruckPaginatedSort = "model"
	TruckPaginatedSortLicensePlate     TruckPaginatedSort = "licensePlate"
	TruckPaginatedSortPersonCapacity   TruckPaginatedSort = "personCapacity"
	TruckPaginatedSortWayName          TruckPaginatedSort = "wayName"
	TruckPaginatedSortMunicipalityName TruckPaginatedSort = "municipalityName"
	TruckPaginatedSortCreatedAt        TruckPaginatedSort = "createdAt"
	TruckPaginatedSortModifiedAt       TruckPaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s TruckPaginatedSort) Field() TruckPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s TruckPaginatedSort) Valid() bool {
	switch s {
	case TruckPaginatedSortMake,
		TruckPaginatedSortModel,
		TruckPaginatedSortLicensePlate,
		TruckPaginatedSortPersonCapacity,
		TruckPaginatedSortWayName,
		TruckPaginatedSortMunicipalityName,
		TruckPaginatedSortCreatedAt,
		TruckPaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// TrucksPaginatedFilter defines the trucks filter structure.
type TrucksPaginatedFilter struct {
	PaginatedRequest[TruckPaginatedSort]
	Make         *TruckMake
	Model        *TruckModel
	LicensePlate *TruckLicensePlate
	LocationName *string
}
