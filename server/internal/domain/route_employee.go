package domain

// RouteEmployeeRole defines the role of the route employee.
type RouteEmployeeRole string

const (
	RouteEmployeeRoleDriver    RouteEmployeeRole = "driver"
	RouteEmployeeRoleCollector RouteEmployeeRole = "collector"
)

// Valid returns true if the employee role is valid, false otherwise.
func (r RouteEmployeeRole) Valid() bool {
	switch r {
	case RouteEmployeeRoleDriver,
		RouteEmployeeRoleCollector:
		return true
	default:
		return false
	}
}

// RouteEmployee defines the route employee structure.
type RouteEmployee struct {
	Employee
	RouteRole RouteEmployeeRole
}

// RouteEmployeePaginatedSort defines the field of the route employee to sort.
type RouteEmployeePaginatedSort string

const (
	RouteEmployeePaginatedSortRouteRole RouteEmployeePaginatedSort = "routeRole"
	RouteEmployeePaginatedSortCreatedAt RouteEmployeePaginatedSort = "createdAt"
)

// Field returns the name of the field to sort by.
func (s RouteEmployeePaginatedSort) Field() RouteEmployeePaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s RouteEmployeePaginatedSort) Valid() bool {
	switch s {
	case RouteEmployeePaginatedSortRouteRole,
		RouteEmployeePaginatedSortCreatedAt:
		return true
	default:
		return false
	}
}

// RouteEmployeesPaginatedFilter defines the route employees filter structure.
type RouteEmployeesPaginatedFilter struct {
	PaginatedRequest[RouteEmployeePaginatedSort]
	RouteRole *RouteEmployeeRole
}
