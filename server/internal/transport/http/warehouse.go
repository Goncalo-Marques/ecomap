package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errWarehouseNotFound                            = "warehouse not found"
	errWarehouseAlreadyHasMoreTrucksThanNewCapacity = "warehouse already has more trucks than the new capacity"
	errWarehouseAssociatedWithWarehouseTruck        = "warehouse associated with truck"
	errWarehouseAssociatedWithRouteDeparture        = "warehouse associated with route as departure"
	errWarehouseAssociatedWithRouteArrival          = "warehouse associated with route as arrival"
)

// CreateWarehouse handles the http request to create a warehouse.
func (h *handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var warehousePost spec.WarehousePost
	err = json.Unmarshal(requestBody, &warehousePost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableWarehouse, err := warehousePostToDomain(warehousePost)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainWarehouse, err := h.service.CreateWarehouse(ctx, domainEditableWarehouse)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	warehouse, err := warehouseFromDomain(domainWarehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(warehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListWarehouses handles the http request to list warehouses.
func (h *handler) ListWarehouses(w http.ResponseWriter, r *http.Request, params spec.ListWarehousesParams) {
	ctx := r.Context()

	domainWarehousesFilter := listWarehousesParamsToDomain(params)
	domainPaginatedWarehouses, err := h.service.ListWarehouses(ctx, domainWarehousesFilter)
	if err != nil {
		var domainErrFilterValueInvalid *domain.ErrFilterValueInvalid

		switch {
		case errors.As(err, &domainErrFilterValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFilterValueInvalid, domainErrFilterValueInvalid.FilterName))
		default:
			internalServerError(w)
		}

		return
	}

	warehousesPaginated, err := warehousesPaginatedFromDomain(domainPaginatedWarehouses)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(warehousesPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetWarehouseByID handles the http request to get a warehouse by ID.
func (h *handler) GetWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	ctx := r.Context()

	domainWarehouse, err := h.service.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseNotFound):
			notFound(w, errWarehouseNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	warehouse, err := warehouseFromDomain(domainWarehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(warehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchWarehouseByID handles the http request to modify a warehouse by ID.
func (h *handler) PatchWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var warehousePatch spec.WarehousePatch
	err = json.Unmarshal(requestBody, &warehousePatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableWarehouse, err := warehousePatchToDomain(warehousePatch)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainWarehouse, err := h.service.PatchWarehouse(ctx, warehouseID, domainEditableWarehouse)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrWarehouseNotFound):
			notFound(w, errWarehouseNotFound)
		case errors.Is(err, domain.ErrWarehouseTruckCapacityMinLimit):
			conflict(w, errWarehouseAlreadyHasMoreTrucksThanNewCapacity)
		default:
			internalServerError(w)
		}

		return
	}

	warehouse, err := warehouseFromDomain(domainWarehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(warehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteWarehouseByID handles the http request to delete a warehouse by ID.
func (h *handler) DeleteWarehouseByID(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam) {
	ctx := r.Context()

	domainWarehouse, err := h.service.DeleteWarehouseByID(ctx, warehouseID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseNotFound):
			notFound(w, errWarehouseNotFound)
		case errors.Is(err, domain.ErrWarehouseAssociatedWithWarehouseTruck):
			conflict(w, errWarehouseAssociatedWithWarehouseTruck)
		case errors.Is(err, domain.ErrWarehouseAssociatedWithRouteDeparture):
			conflict(w, errWarehouseAssociatedWithRouteDeparture)
		case errors.Is(err, domain.ErrWarehouseAssociatedWithRouteArrival):
			conflict(w, errWarehouseAssociatedWithRouteArrival)
		default:
			internalServerError(w)
		}

		return
	}

	warehouse, err := warehouseFromDomain(domainWarehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(warehouse)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// warehousePostToDomain returns a domain editable warehouse based on the standardized warehouse post.
func warehousePostToDomain(warehousePost spec.WarehousePost) (domain.EditableWarehouse, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(&warehousePost.GeoJson)
	if err != nil {
		return domain.EditableWarehouse{}, err
	}

	return domain.EditableWarehouse{
		TruckCapacity: domain.WarehouseTruckCapacity(warehousePost.TruckCapacity),
		GeoJSON:       geoJSON,
	}, nil
}

// warehousePatchToDomain returns a domain patchable warehouse based on the standardized warehouse patch.
func warehousePatchToDomain(warehousePatch spec.WarehousePatch) (domain.EditableWarehousePatch, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(warehousePatch.GeoJson)
	if err != nil {
		return domain.EditableWarehousePatch{}, err
	}

	return domain.EditableWarehousePatch{
		TruckCapacity: (*domain.WarehouseTruckCapacity)(warehousePatch.TruckCapacity),
		GeoJSON:       geoJSON,
	}, nil
}

// listWarehousesParamsToDomain returns a domain warehouses paginated filter based on the standardized list warehouses parameters.
func listWarehousesParamsToDomain(params spec.ListWarehousesParams) domain.WarehousesPaginatedFilter {
	domainSort := domain.WarehousePaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListWarehousesParamsSortTruckCapacity:
			domainSort = domain.WarehousePaginatedSortTruckCapacity
		case spec.ListWarehousesParamsSortWayName:
			domainSort = domain.WarehousePaginatedSortWayName
		case spec.ListWarehousesParamsSortMunicipalityName:
			domainSort = domain.WarehousePaginatedSortMunicipalityName
		case spec.ListWarehousesParamsSortCreatedAt:
			domainSort = domain.WarehousePaginatedSortCreatedAt
		case spec.ListWarehousesParamsSortModifiedAt:
			domainSort = domain.WarehousePaginatedSortModifiedAt
		default:
			domainSort = domain.WarehousePaginatedSort(*params.Sort)
		}
	}

	return domain.WarehousesPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		LocationName: params.LocationName,
	}
}

// warehouseFromDomain returns a standardized warehouse based on the domain model.
func warehouseFromDomain(warehouse domain.Warehouse) (spec.Warehouse, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(warehouse.GeoJSON)
	if err != nil {
		return spec.Warehouse{}, err
	}

	return spec.Warehouse{
		Id:            warehouse.ID,
		TruckCapacity: int(warehouse.TruckCapacity),
		GeoJson:       geoJSON,
		CreatedAt:     warehouse.CreatedAt,
		ModifiedAt:    warehouse.ModifiedAt,
	}, nil
}

// warehousesFromDomain returns standardized warehouses based on the domain model.
func warehousesFromDomain(warehouses []domain.Warehouse) ([]spec.Warehouse, error) {
	specWarehouses := make([]spec.Warehouse, len(warehouses))
	var err error

	for i, warehouse := range warehouses {
		specWarehouses[i], err = warehouseFromDomain(warehouse)
		if err != nil {
			return []spec.Warehouse{}, err
		}
	}

	return specWarehouses, nil
}

// warehousesPaginatedFromDomain returns a standardized warehouses paginated response based on the domain model.
func warehousesPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Warehouse]) (spec.WarehousesPaginated, error) {
	warehouses, err := warehousesFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.WarehousesPaginated{}, err
	}

	return spec.WarehousesPaginated{
		Total:      paginatedResponse.Total,
		Warehouses: warehouses,
	}, nil
}
