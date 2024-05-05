package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// ListRouteEmployees handles the http request to list route employees.
func (h *handler) ListRouteEmployees(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, params spec.ListRouteEmployeesParams) {
	w.WriteHeader(http.StatusNotFound)
}

// CreateRouteEmployee handles the http request to create a route employee association.
func (h *handler) CreateRouteEmployee(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, employeeID spec.EmployeeIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteRouteEmployee handles the http request to delete a route employee association.
func (h *handler) DeleteRouteEmployee(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, employeeID spec.EmployeeIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
