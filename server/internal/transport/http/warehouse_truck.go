package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errWarehouseTruckAlreadyExists                = "warehouse truck association already exists"
	errWarehouseTruckAssociatedWithRouteDeparture = "warehouse truck associated with route as departure"
	errWarehouseTruckAssociatedWithRouteArrival   = "warehouse truck associated with route as arrival"
	errWarehouseTruckNotFound                     = "warehouse truck association does not exist"
)

// ListWarehouseTrucks handles the http request to list warehouse trucks.
func (h *handler) ListWarehouseTrucks(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, params spec.ListWarehouseTrucksParams) {
	ctx := r.Context()

	domainWarehouseTrucksFilter := listWarehouseTrucksParamsToDomain(params)
	domainPaginatedTrucks, err := h.service.ListWarehouseTrucks(ctx, warehouseID, domainWarehouseTrucksFilter)
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

	trucksPaginated, err := trucksPaginatedFromDomain(domainPaginatedTrucks)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(trucksPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// CreateWarehouseTruck handles the http request to create a warehouse truck association.
func (h *handler) CreateWarehouseTruck(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, truckID spec.TruckIdPathParam) {
	ctx := r.Context()

	err := h.service.CreateWarehouseTruck(ctx, warehouseID, truckID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseNotFound):
			notFound(w, errWarehouseNotFound)
		case errors.Is(err, domain.ErrTruckNotFound):
			notFound(w, errTruckNotFound)
		case errors.Is(err, domain.ErrWarehouseTruckCapacityMaxLimit):
			conflict(w, errWarehouseTruckCapacityMaxLimit)
		case errors.Is(err, domain.ErrWarehouseTruckAlreadyExists):
			conflict(w, errWarehouseTruckAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// DeleteWarehouseTruck handles the http request to delete a warehouse truck association.
func (h *handler) DeleteWarehouseTruck(w http.ResponseWriter, r *http.Request, warehouseID spec.WarehouseIdPathParam, truckID spec.TruckIdPathParam) {
	ctx := r.Context()

	err := h.service.DeleteWarehouseTruck(ctx, warehouseID, truckID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseTruckAssociatedWithRouteDeparture):
			conflict(w, errWarehouseTruckAssociatedWithRouteDeparture)
		case errors.Is(err, domain.ErrWarehouseTruckAssociatedWithRouteArrival):
			conflict(w, errWarehouseTruckAssociatedWithRouteArrival)
		case errors.Is(err, domain.ErrWarehouseTruckNotFound):
			conflict(w, errWarehouseTruckNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// listWarehouseTrucksParamsToDomain returns a domain warehouse trucks paginated filter based on the standardized list
// warehouse trucks parameters.
func listWarehouseTrucksParamsToDomain(params spec.ListWarehouseTrucksParams) domain.WarehouseTrucksPaginatedFilter {
	domainSort := domain.WarehouseTruckPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListWarehouseTrucksParamsSortTruckMake:
			domainSort = domain.WarehouseTruckPaginatedSortTruckMake
		case spec.ListWarehouseTrucksParamsSortTruckModel:
			domainSort = domain.WarehouseTruckPaginatedSortTruckModel
		case spec.ListWarehouseTrucksParamsSortTruckLicensePlate:
			domainSort = domain.WarehouseTruckPaginatedSortTruckLicensePlate
		case spec.ListWarehouseTrucksParamsSortTruckPersonCapacity:
			domainSort = domain.WarehouseTruckPaginatedSortTruckPersonCapacity
		case spec.ListWarehouseTrucksParamsSortTruckWayName:
			domainSort = domain.WarehouseTruckPaginatedSortTruckWayName
		case spec.ListWarehouseTrucksParamsSortTruckMunicipalityName:
			domainSort = domain.WarehouseTruckPaginatedSortTruckMunicipalityName
		case spec.ListWarehouseTrucksParamsSortCreatedAt:
			domainSort = domain.WarehouseTruckPaginatedSortCreatedAt
		default:
			domainSort = domain.WarehouseTruckPaginatedSort(*params.Sort)
		}
	}

	return domain.WarehouseTrucksPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		TruckMake:         (*domain.TruckMake)(params.TruckMake),
		TruckModel:        (*domain.TruckModel)(params.TruckModel),
		TruckLicensePlate: (*domain.TruckLicensePlate)(params.TruckLicensePlate),
		LocationName:      params.LocationName,
	}
}
