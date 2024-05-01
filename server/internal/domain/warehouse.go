package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Warehouse constraints.
const (
	warehouseTruckCapacityMinValue = 0
)

// Warehouse errors.
var (
	ErrWarehouseNotFound                     = errors.New("warehouse not found")                          // Returned when a warehouse is not found.
	ErrWarehouseTruckCapacityMinLimit        = errors.New("warehouse truck capacity below minimum limit") // Returned when a warehouse capacity is below the minimum limit.
	ErrWarehouseAssociatedWithWarehouseTruck = errors.New("warehouse associated with truck")              // Returned when a warehouse is associated with a truck.
	ErrWarehouseAssociatedWithRouteDeparture = errors.New("warehouse associated with route as departure") // Returned when a warehouse is associated with a route as a departure.
	ErrWarehouseAssociatedWithRouteArrival   = errors.New("warehouse associated with route as arrival")   // Returned when a warehouse is associated with a route as an arrival.
)

// WarehouseTruckCapacity defines the truck capacity of the warehouse.
type WarehouseTruckCapacity int

// Valid returns true if the truck capacity is valid, false otherwise.
func (tc WarehouseTruckCapacity) Valid() bool {
	return tc >= warehouseTruckCapacityMinValue
}

// EditableWarehouse defines the editable warehouse structure.
type EditableWarehouse struct {
	TruckCapacity WarehouseTruckCapacity
	GeoJSON       GeoJSON
}

// EditableWarehousePatch defines the patchable warehouse structure.
type EditableWarehousePatch struct {
	TruckCapacity *WarehouseTruckCapacity
	GeoJSON       GeoJSON
}

// Warehouse defines the warehouse structure.
type Warehouse struct {
	EditableWarehouse
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// WarehousePaginatedSort defines the field of the warehouse to sort.
type WarehousePaginatedSort string

const (
	WarehousePaginatedSortWayName          WarehousePaginatedSort = "wayName"
	WarehousePaginatedSortMunicipalityName WarehousePaginatedSort = "municipalityName"
	WarehousePaginatedSortCreatedAt        WarehousePaginatedSort = "createdAt"
	WarehousePaginatedSortModifiedAt       WarehousePaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s WarehousePaginatedSort) Field() WarehousePaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s WarehousePaginatedSort) Valid() bool {
	switch s {
	case WarehousePaginatedSortWayName,
		WarehousePaginatedSortMunicipalityName,
		WarehousePaginatedSortCreatedAt,
		WarehousePaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// WarehousesPaginatedFilter defines the warehouses filter structure.
type WarehousesPaginatedFilter struct {
	PaginatedRequest[WarehousePaginatedSort]
	LocationName *string
}
