package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// CreateLandfill handles the http request to create a landfill.
func (h *handler) CreateLandfill(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListLandfills handles the http request to list landfills.
func (h *handler) ListLandfills(w http.ResponseWriter, r *http.Request, params spec.ListLandfillsParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetLandfillByID handles the http request to get a landfill by ID.
func (h *handler) GetLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// PatchLandfillByID handles the http request to modify a landfill by ID.
func (h *handler) PatchLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteLandfillByID handles the http request to delete a landfill by ID.
func (h *handler) DeleteLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
