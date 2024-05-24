package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Field constraints.
const (
	routeNameMaxLength = 50
)

// Route errors.
var (
	ErrRouteNotFound                    = errors.New("route not found")                                 // Returned when a route is not found.
	ErrRouteDepartureWarehouseNotFound  = errors.New("route departure warehouse not found")             // Returned when a route departure warehouse is not found.
	ErrRouteArrivalWarehouseNotFound    = errors.New("route arrival warehouse not found")               // Returned when a route arrival warehouse is not found.
	ErrRouteTruckPersonCapacityMinLimit = errors.New("route truck person capacity below minimum limit") // Returned when a route truck person capacity is below the minimum limit.
	ErrRouteTruckPersonCapacityMaxLimit = errors.New("route truck person capacity above maximum limit") // Returned when a route truck person capacity is above the maximum limit.
)

// RouteName defines the route name type.
type RouteName string

// Valid returns true if the name is valid, false otherwise.
func (n RouteName) Valid() bool {
	return len(n) <= routeNameMaxLength
}

// EditableRoute defines the editable route structure.
type EditableRoute struct {
	Name                 RouteName
	TruckID              uuid.UUID
	DepartureWarehouseID uuid.UUID
	ArrivalWarehouseID   uuid.UUID
}

// EditableRoutePatch defines the patchable route structure.
type EditableRoutePatch struct {
	Name                 *RouteName
	TruckID              *uuid.UUID
	DepartureWarehouseID *uuid.UUID
	ArrivalWarehouseID   *uuid.UUID
}

// Route defines the route structure.
type Route struct {
	ID                 uuid.UUID
	Name               RouteName
	Truck              Truck
	DepartureWarehouse Warehouse
	ArrivalWarehouse   Warehouse
	CreatedAt          time.Time
	ModifiedAt         time.Time
}

// RoutePaginatedSort defines the field of the route to sort.
type RoutePaginatedSort string

const (
	RoutePaginatedSortName                 RoutePaginatedSort = "name"
	RoutePaginatedSortTruckID              RoutePaginatedSort = "truckID"
	RoutePaginatedSortDepartureWarehouseID RoutePaginatedSort = "departureWarehouseID"
	RoutePaginatedSortArrivalWarehouseID   RoutePaginatedSort = "arrivalWarehouseID"
	RoutePaginatedSortCreatedAt            RoutePaginatedSort = "createdAt"
	RoutePaginatedSortModifiedAt           RoutePaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s RoutePaginatedSort) Field() RoutePaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s RoutePaginatedSort) Valid() bool {
	switch s {
	case RoutePaginatedSortName,
		RoutePaginatedSortTruckID,
		RoutePaginatedSortDepartureWarehouseID,
		RoutePaginatedSortArrivalWarehouseID,
		RoutePaginatedSortCreatedAt,
		RoutePaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// RoutesPaginatedFilter defines the routes filter structure.
type RoutesPaginatedFilter struct {
	PaginatedRequest[RoutePaginatedSort]
	Name                 *RouteName
	TruckID              *uuid.UUID
	DepartureWarehouseID *uuid.UUID
	ArrivalWarehouseID   *uuid.UUID
}
