package domain

import "errors"

// Warehouse truck errors.
var (
	ErrWarehouseTruckAlreadyExists                = errors.New("warehouse truck already exists")                     // Returned when a warehouse truck association already exists.
	ErrWarehouseTruckNotFound                     = errors.New("warehouse truck not found")                          // Returned when a warehouse truck association is not found.
	ErrWarehouseTruckAssociatedWithRouteDeparture = errors.New("warehouse truck associated with route as departure") // Returned when a warehouse truck is associated with a route as a departure.
	ErrWarehouseTruckAssociatedWithRouteArrival   = errors.New("warehouse truck associated with route as arrival")   // Returned when a warehouse truck is associated with a route as an arrival.
)

// WarehouseTruckPaginatedSort defines the field of the warehouse truck to sort.
type WarehouseTruckPaginatedSort string

const (
	WarehouseTruckPaginatedSortTruckMake             WarehouseTruckPaginatedSort = "truckMake"
	WarehouseTruckPaginatedSortTruckModel            WarehouseTruckPaginatedSort = "truckModel"
	WarehouseTruckPaginatedSortTruckLicensePlate     WarehouseTruckPaginatedSort = "truckLicensePlate"
	WarehouseTruckPaginatedSortTruckPersonCapacity   WarehouseTruckPaginatedSort = "truckPersonCapacity"
	WarehouseTruckPaginatedSortTruckWayName          WarehouseTruckPaginatedSort = "truckWayName"
	WarehouseTruckPaginatedSortTruckMunicipalityName WarehouseTruckPaginatedSort = "truckMunicipalityName"
	WarehouseTruckPaginatedSortCreatedAt             WarehouseTruckPaginatedSort = "createdAt"
)

// Field returns the name of the field to sort by.
func (s WarehouseTruckPaginatedSort) Field() WarehouseTruckPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s WarehouseTruckPaginatedSort) Valid() bool {
	switch s {
	case WarehouseTruckPaginatedSortTruckMake,
		WarehouseTruckPaginatedSortTruckModel,
		WarehouseTruckPaginatedSortTruckLicensePlate,
		WarehouseTruckPaginatedSortTruckPersonCapacity,
		WarehouseTruckPaginatedSortTruckWayName,
		WarehouseTruckPaginatedSortTruckMunicipalityName,
		WarehouseTruckPaginatedSortCreatedAt:
		return true
	default:
		return false
	}
}

// WarehouseTrucksPaginatedFilter defines the warehouse trucks filter structure.
type WarehouseTrucksPaginatedFilter struct {
	PaginatedRequest[WarehouseTruckPaginatedSort]
	TruckMake         *TruckMake
	TruckModel        *TruckModel
	TruckLicensePlate *TruckLicensePlate
	LocationName      *string
}
