package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// CreateWarehouse handles the http request to create a warehouse.
func (h *handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListWarehouses handles the http request to list warehouses.
func (h *handler) ListWarehouses(w http.ResponseWriter, r *http.Request, params spec.ListWarehousesParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetWarehouseByID handles the http request to get a warehouse by ID.
func (h *handler) GetWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// PatchWarehouseByID handles the http request to modify a warehouse by ID.
func (h *handler) PatchWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteWarehouseByID handles the http request to delete a warehouse by ID.
func (h *handler) DeleteWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
