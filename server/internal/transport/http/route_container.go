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
	errRouteContainerAlreadyExists = "route container association already exists"
	errRouteContainerNotFound      = "route container association does not exist"
)

// ListRouteContainers handles the http request to list route containers.
func (h *handler) ListRouteContainers(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, params spec.ListRouteContainersParams) {
	ctx := r.Context()

	domainRouteContainersFilter := listRouteContainersParamsToDomain(params)
	domainPaginatedContainers, err := h.service.ListRouteContainers(ctx, routeID, domainRouteContainersFilter)
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

	containersPaginated, err := containersPaginatedFromDomain(domainPaginatedContainers)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(containersPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// CreateRouteContainer handles the http request to create a route container association.
func (h *handler) CreateRouteContainer(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	err := h.service.CreateRouteContainer(ctx, routeID, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		case errors.Is(err, domain.ErrContainerNotFound):
			notFound(w, errContainerNotFound)
		case errors.Is(err, domain.ErrRouteContainerAlreadyExists):
			conflict(w, errRouteContainerAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// DeleteRouteContainer handles the http request to delete a route container association.
func (h *handler) DeleteRouteContainer(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	err := h.service.DeleteRouteContainer(ctx, routeID, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteContainerNotFound):
			conflict(w, errRouteContainerNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// listRouteContainersParamsToDomain returns a domain route containers paginated filter based on the standardized list
// route containers parameters.
func listRouteContainersParamsToDomain(params spec.ListRouteContainersParams) domain.RouteContainersPaginatedFilter {
	domainSort := domain.RouteContainerPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListRouteContainersParamsSortContainerCategory:
			domainSort = domain.RouteContainerPaginatedSortContainerCategory
		case spec.ListRouteContainersParamsSortContainerWayName:
			domainSort = domain.RouteContainerPaginatedSortContainerWayName
		case spec.ListRouteContainersParamsSortContainerMunicipalityName:
			domainSort = domain.RouteContainerPaginatedSortContainerMunicipalityName
		case spec.ListRouteContainersParamsSortCreatedAt:
			domainSort = domain.RouteContainerPaginatedSortCreatedAt
		default:
			domainSort = domain.RouteContainerPaginatedSort(*params.Sort)
		}
	}

	var domainContainerCategory *domain.ContainerCategory
	if params.ContainerCategory != nil {
		category := containerCategoryToDomain(*params.ContainerCategory)
		domainContainerCategory = &category
	}

	return domain.RouteContainersPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		ContainerCategory: domainContainerCategory,
		LocationName:      params.LocationName,
	}
}
