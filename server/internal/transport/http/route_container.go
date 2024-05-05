package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// ListRouteContainers handles the http request to list route containers.
func (h *handler) ListRouteContainers(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, params spec.ListRouteContainersParams) {
	w.WriteHeader(http.StatusNotFound)
}

// CreateRouteContainer handles the http request to create a route container association.
func (h *handler) CreateRouteContainer(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteRouteContainer handles the http request to delete a route container association.
func (h *handler) DeleteRouteContainer(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
