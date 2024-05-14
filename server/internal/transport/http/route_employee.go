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
	errRouteEmployeeAlreadyExists = "route employee association already exists"
	errRouteEmployeeNotFound      = "route employee association does not exist"
)

// ListRouteEmployees handles the http request to list route employees.
func (h *handler) ListRouteEmployees(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, params spec.ListRouteEmployeesParams) {
	ctx := r.Context()

	domainRouteEmployeesFilter := listRouteEmployeesParamsToDomain(params)
	domainPaginatedRouteEmployees, err := h.service.ListRouteEmployees(ctx, routeID, domainRouteEmployeesFilter)
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

	routeEmployeesPaginated, err := routeEmployeesPaginatedFromDomain(domainPaginatedRouteEmployees)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(routeEmployeesPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// CreateRouteEmployee handles the http request to create a route employee association.
func (h *handler) CreateRouteEmployee(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, employeeID spec.EmployeeIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var routeEmployeePost spec.RouteEmployeePost
	err = json.Unmarshal(requestBody, &routeEmployeePost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableRouteEmployee := routeEmployeePostToDomain(routeEmployeePost)
	err = h.service.CreateRouteEmployee(ctx, routeID, employeeID, domainEditableRouteEmployee)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrRouteNotFound):
			notFound(w, errRouteNotFound)
		case errors.Is(err, domain.ErrEmployeeNotFound):
			notFound(w, errEmployeeNotFound)
		case errors.Is(err, domain.ErrRouteTruckPersonCapacityMaxLimit):
			conflict(w, errRouteTruckPersonCapacityMaxLimit)
		case errors.Is(err, domain.ErrRouteEmployeeAlreadyExists):
			conflict(w, errRouteEmployeeAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// DeleteRouteEmployee handles the http request to delete a route employee association.
func (h *handler) DeleteRouteEmployee(w http.ResponseWriter, r *http.Request, routeID spec.RouteIdPathParam, employeeID spec.EmployeeIdPathParam) {
	ctx := r.Context()

	err := h.service.DeleteRouteEmployee(ctx, routeID, employeeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteEmployeeNotFound):
			conflict(w, errRouteEmployeeNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// routeEmployeeRoleToDomain returns a domain route employee role based on the standardized model.
func routeEmployeeRoleToDomain(role spec.RouteEmployeeRole) domain.RouteEmployeeRole {
	switch role {
	case spec.Driver:
		return domain.RouteEmployeeRoleDriver
	case spec.Collector:
		return domain.RouteEmployeeRoleCollector
	default:
		return domain.RouteEmployeeRole(role)
	}
}

// routeEmployeeRoleFromDomain returns a standardized route employee role based on the domain model.
func routeEmployeeRoleFromDomain(role domain.RouteEmployeeRole) spec.RouteEmployeeRole {
	switch role {
	case domain.RouteEmployeeRoleDriver:
		return spec.Driver
	case domain.RouteEmployeeRoleCollector:
		return spec.Collector
	default:
		return spec.RouteEmployeeRole(role)
	}
}

// routeEmployeePostToDomain returns a domain editable route employee based on the standardized route employee post.
func routeEmployeePostToDomain(routeEmployeePost spec.RouteEmployeePost) domain.EditableRouteEmployee {
	return domain.EditableRouteEmployee{
		RouteRole: routeEmployeeRoleToDomain(routeEmployeePost.RouteRole),
	}
}

// listRouteEmployeesParamsToDomain returns a domain route employees paginated filter based on the standardized list
// route Employees parameters.
func listRouteEmployeesParamsToDomain(params spec.ListRouteEmployeesParams) domain.RouteEmployeesPaginatedFilter {
	domainSort := domain.RouteEmployeePaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListRouteEmployeesParamsSortRouteRole:
			domainSort = domain.RouteEmployeePaginatedSortRouteRole
		case spec.ListRouteEmployeesParamsSortCreatedAt:
			domainSort = domain.RouteEmployeePaginatedSortCreatedAt
		default:
			domainSort = domain.RouteEmployeePaginatedSort(*params.Sort)
		}
	}

	var domainRouteEmployeeRole *domain.RouteEmployeeRole
	if params.RouteRole != nil {
		role := routeEmployeeRoleToDomain(*params.RouteRole)
		domainRouteEmployeeRole = &role
	}

	return domain.RouteEmployeesPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		RouteRole: domainRouteEmployeeRole,
	}
}

// routeEmployeeFromDomain returns a standardized route employee based on the domain model.
func routeEmployeeFromDomain(routeEmployee domain.RouteEmployee) (spec.RouteEmployee, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(routeEmployee.GeoJSON)
	if err != nil {
		return spec.RouteEmployee{}, err
	}

	return spec.RouteEmployee{
		Id:            routeEmployee.ID,
		Username:      string(routeEmployee.Username),
		FirstName:     string(routeEmployee.FirstName),
		LastName:      string(routeEmployee.LastName),
		Role:          employeeRoleFromDomain(routeEmployee.Role),
		DateOfBirth:   dateFromTime(routeEmployee.DateOfBirth),
		PhoneNumber:   string(routeEmployee.PhoneNumber),
		GeoJson:       geoJSON,
		ScheduleStart: timeStringFromTime(routeEmployee.ScheduleStart),
		ScheduleEnd:   timeStringFromTime(routeEmployee.ScheduleEnd),
		CreatedAt:     routeEmployee.CreatedAt,
		ModifiedAt:    routeEmployee.ModifiedAt,
		RouteRole:     routeEmployeeRoleFromDomain(routeEmployee.RouteRole),
	}, nil
}

// routeEmployeesFromDomain returns standardized employees based on the domain model.
func routeEmployeesFromDomain(routeEmployees []domain.RouteEmployee) ([]spec.RouteEmployee, error) {
	specRouteEmployees := make([]spec.RouteEmployee, len(routeEmployees))
	var err error

	for i, employee := range routeEmployees {
		specRouteEmployees[i], err = routeEmployeeFromDomain(employee)
		if err != nil {
			return []spec.RouteEmployee{}, err
		}
	}

	return specRouteEmployees, nil
}

// routeEmployeesPaginatedFromDomain returns a standardized route employees paginated response based on the domain model.
func routeEmployeesPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.RouteEmployee]) (spec.RouteEmployeesPaginated, error) {
	routeEmployees, err := routeEmployeesFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.RouteEmployeesPaginated{}, err
	}

	return spec.RouteEmployeesPaginated{
		Total:     paginatedResponse.Total,
		Employees: routeEmployees,
	}, nil
}
