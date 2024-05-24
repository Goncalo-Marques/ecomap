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
	errRouteNotFound                    = "route not found"
	errRouteDepartureWarehouseNotFound  = "route departure warehouse not found"
	errRouteArrivalWarehouseNotFound    = "route arrival warehouse not found"
	errRouteTruckPersonCapacityMinLimit = "route already has more employees than the new truck has capacity"
	errRouteTruckPersonCapacityMaxLimit = "route cannot have more employees than the truck has capacity"
)

// CreateRoute handles the http request to create a route.
func (h *handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var routePost spec.RoutePost
	err = json.Unmarshal(requestBody, &routePost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableRoute := routePostToDomain(routePost)
	domainRoute, err := h.service.CreateRoute(ctx, domainEditableRoute)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrTruckNotFound):
			conflict(w, errTruckNotFound)
		case errors.Is(err, domain.ErrRouteDepartureWarehouseNotFound):
			conflict(w, errRouteDepartureWarehouseNotFound)
		case errors.Is(err, domain.ErrRouteArrivalWarehouseNotFound):
			conflict(w, errRouteArrivalWarehouseNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	route, err := routeFromDomain(domainRoute)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(route)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListRoutes handles the http request to list routes.
func (h *handler) ListRoutes(w http.ResponseWriter, r *http.Request, params spec.ListRoutesParams) {
	ctx := r.Context()

	domainRoutesFilter := listRoutesParamsToDomain(params)
	domainPaginatedRoutes, err := h.service.ListRoutes(ctx, domainRoutesFilter)
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

	routesPaginated, err := routesPaginatedFromDomain(domainPaginatedRoutes)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(routesPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetRouteByID handles the http request to get a route by ID.
func (h *handler) GetRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	ctx := r.Context()

	domainRoute, err := h.service.GetRouteByID(ctx, routeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	route, err := routeFromDomain(domainRoute)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(route)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchRouteByID handles the http request to modify a route by ID.
func (h *handler) PatchRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var routePatch spec.RoutePatch
	err = json.Unmarshal(requestBody, &routePatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableRoute := routePatchToDomain(routePatch)
	domainRoute, err := h.service.PatchRoute(ctx, routeID, domainEditableRoute)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		case errors.Is(err, domain.ErrTruckNotFound):
			conflict(w, errTruckNotFound)
		case errors.Is(err, domain.ErrRouteDepartureWarehouseNotFound):
			conflict(w, errRouteDepartureWarehouseNotFound)
		case errors.Is(err, domain.ErrRouteArrivalWarehouseNotFound):
			conflict(w, errRouteArrivalWarehouseNotFound)
		case errors.Is(err, domain.ErrRouteTruckPersonCapacityMinLimit):
			conflict(w, errRouteTruckPersonCapacityMinLimit)
		default:
			internalServerError(w)
		}

		return
	}

	route, err := routeFromDomain(domainRoute)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(route)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteRouteByID handles the http request to delete a route by ID.
func (h *handler) DeleteRouteByID(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	ctx := r.Context()

	domainRoute, err := h.service.DeleteRouteByID(ctx, routeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	route, err := routeFromDomain(domainRoute)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(route)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetRouteWays handles the http request to get ways of a route.
func (h *handler) GetRouteWays(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam) {
	ctx := r.Context()

	domainGeoJSON, err := h.service.GetRouteRoads(ctx, routeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	geoJSON, err := geoJSONFeatureCollectionLineStringFromDomain(domainGeoJSON)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(geoJSON)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// routePostToDomain returns a domain editable route based on the standardized route post.
func routePostToDomain(routePost spec.RoutePost) domain.EditableRoute {
	return domain.EditableRoute{
		Name:                 domain.RouteName(routePost.Name),
		TruckID:              routePost.TruckId,
		DepartureWarehouseID: routePost.DepartureWarehouseId,
		ArrivalWarehouseID:   routePost.ArrivalWarehouseId,
	}
}

// routePatchToDomain returns a domain patchable route based on the standardized route patch.
func routePatchToDomain(routePatch spec.RoutePatch) domain.EditableRoutePatch {
	return domain.EditableRoutePatch{
		Name:                 (*domain.RouteName)(routePatch.Name),
		TruckID:              routePatch.TruckId,
		DepartureWarehouseID: routePatch.DepartureWarehouseId,
		ArrivalWarehouseID:   routePatch.ArrivalWarehouseId,
	}
}

// listRoutesParamsToDomain returns a domain routes paginated filter based on the standardized list routes parameters.
func listRoutesParamsToDomain(params spec.ListRoutesParams) domain.RoutesPaginatedFilter {
	domainSort := domain.RoutePaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListRoutesParamsSortName:
			domainSort = domain.RoutePaginatedSortName
		case spec.ListRoutesParamsSortTruckId:
			domainSort = domain.RoutePaginatedSortTruckID
		case spec.ListRoutesParamsSortDepartureWarehouseId:
			domainSort = domain.RoutePaginatedSortDepartureWarehouseID
		case spec.ListRoutesParamsSortArrivalWarehouseId:
			domainSort = domain.RoutePaginatedSortArrivalWarehouseID
		case spec.ListRoutesParamsSortCreatedAt:
			domainSort = domain.RoutePaginatedSortCreatedAt
		case spec.ListRoutesParamsSortModifiedAt:
			domainSort = domain.RoutePaginatedSortModifiedAt
		default:
			domainSort = domain.RoutePaginatedSort(*params.Sort)
		}
	}

	return domain.RoutesPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		Name:                 (*domain.RouteName)(params.Name),
		TruckID:              params.TruckId,
		DepartureWarehouseID: params.DepartureWarehouseId,
		ArrivalWarehouseID:   params.ArrivalWarehouseId,
	}
}

// routeFromDomain returns a standardized route based on the domain model.
func routeFromDomain(route domain.Route) (spec.Route, error) {
	truck, err := truckFromDomain(route.Truck)
	if err != nil {
		return spec.Route{}, err
	}

	departureWarehouse, err := warehouseFromDomain(route.DepartureWarehouse)
	if err != nil {
		return spec.Route{}, err
	}

	arrivalWarehouse, err := warehouseFromDomain(route.ArrivalWarehouse)
	if err != nil {
		return spec.Route{}, err
	}

	return spec.Route{
		Id:                 route.ID,
		Name:               string(route.Name),
		Truck:              truck,
		DepartureWarehouse: departureWarehouse,
		ArrivalWarehouse:   arrivalWarehouse,
		CreatedAt:          route.CreatedAt,
		ModifiedAt:         route.ModifiedAt,
	}, nil
}

// routesFromDomain returns standardized routes based on the domain model.
func routesFromDomain(routes []domain.Route) ([]spec.Route, error) {
	specRoutes := make([]spec.Route, len(routes))
	var err error

	for i, route := range routes {
		specRoutes[i], err = routeFromDomain(route)
		if err != nil {
			return []spec.Route{}, err
		}
	}

	return specRoutes, nil
}

// routesPaginatedFromDomain returns a standardized routes paginated response based on the domain model.
func routesPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Route]) (spec.RoutesPaginated, error) {
	routes, err := routesFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.RoutesPaginated{}, err
	}

	return spec.RoutesPaginated{
		Total:  paginatedResponse.Total,
		Routes: routes,
	}, nil
}
