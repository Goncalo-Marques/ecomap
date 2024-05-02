package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// CreateTruck handles the http request to create a truck.
func (h *handler) CreateTruck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListTrucks handles the http request to list trucks.
func (h *handler) ListTrucks(w http.ResponseWriter, r *http.Request, params spec.ListTrucksParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetTruckByID handles the http request to get a truck by ID.
func (h *handler) GetTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// PatchTruckByID handles the http request to modify a truck by ID.
func (h *handler) PatchTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteTruckByID handles the http request to delete a truck by ID.
func (h *handler) DeleteTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
