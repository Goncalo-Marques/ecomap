package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// CreateRoute handles the http request to create a route.
func (h *handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListRoutes handles the http request to list routes.
func (h *handler) ListRoutes(w http.ResponseWriter, r *http.Request, params spec.ListRoutesParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetRouteByID handles the http request to get a route by ID.
func (h *handler) GetRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// PatchRouteByID handles the http request to modify a route by ID.
func (h *handler) PatchRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteRouteByID handles the http request to delete a route by ID.
func (h *handler) DeleteRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
