package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// CreateContainer handles the http request to create a container.
func (h *handler) CreateContainer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListContainers handles the http request to list containers.
func (h *handler) ListContainers(w http.ResponseWriter, r *http.Request, params spec.ListContainersParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetContainerByID handles the http request to get a container by ID.
func (h *handler) GetContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteContainerByID handles the http request to delete a container by ID.
func (h *handler) DeleteContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
