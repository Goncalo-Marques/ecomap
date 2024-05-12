package domain

import "errors"

// Route container errors.
var (
	ErrRouteContainerAlreadyExists = errors.New("route container already exists") // Returned when a route container association already exists.
	ErrRouteContainerNotFound      = errors.New("route container not found")      // Returned when a route container association is not found.
)

// RouteContainerPaginatedSort defines the field of the route container to sort.
type RouteContainerPaginatedSort string

const (
	RouteContainerPaginatedSortContainerCategory         RouteContainerPaginatedSort = "containerCategory"
	RouteContainerPaginatedSortContainerWayName          RouteContainerPaginatedSort = "containerWayName"
	RouteContainerPaginatedSortContainerMunicipalityName RouteContainerPaginatedSort = "containerMunicipalityName"
	RouteContainerPaginatedSortCreatedAt                 RouteContainerPaginatedSort = "createdAt"
)

// Field returns the name of the field to sort by.
func (s RouteContainerPaginatedSort) Field() RouteContainerPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s RouteContainerPaginatedSort) Valid() bool {
	switch s {
	case RouteContainerPaginatedSortContainerCategory,
		RouteContainerPaginatedSortContainerWayName,
		RouteContainerPaginatedSortContainerMunicipalityName,
		RouteContainerPaginatedSortCreatedAt:
		return true
	default:
		return false
	}
}

// RouteContainersPaginatedFilter defines the route containers filter structure.
type RouteContainersPaginatedFilter struct {
	PaginatedRequest[RouteContainerPaginatedSort]
	ContainerCategory *ContainerCategory
	LocationName      *string
}
