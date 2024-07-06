package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// ListWarehouseTrucks handles the http request to list warehouse trucks.
func (h *handler) ListWarehouseTrucks(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, params spec.ListWarehouseTrucksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// CreateWarehouseTruck handles the http request to create a warehouse truck association.
func (h *handler) CreateWarehouseTruck(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, truckID spec.TruckIdPathParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

// DeleteWarehouseTruck handles the http request to delete a warehouse truck association.
func (h *handler) DeleteWarehouseTruck(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, truckID spec.TruckIdPathParam) {
	w.WriteHeader(http.StatusNotImplemented)
}
