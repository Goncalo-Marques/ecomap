package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/google/uuid"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/authz"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// Base URL const.
const (
	baseURLWebApp = "/"
	baseURLApi    = "/api"
	baseURLDocs   = "/api/docs/"
)

// Directories to serve.
const (
	dirWebApp    = "./dist/web"
	dirSwaggerUI = "./api/swagger"
	dirIndexHTML = "index.html"
)

// Request header const.
const (
	requestHeaderKeyAccept       = "Accept"
	requestHeaderValueAcceptHTML = "text/html"
)

// AuthorizationService defines the authorization service interface.
type AuthorizationService interface {
	Middleware(options authz.MiddlewareOptions) func(http.Handler) http.Handler
}

// Service defines the service interface.
type Service interface {
	CreateUser(ctx context.Context, editableUser domain.EditableUserWithPassword) (domain.User, error)
	ListUsers(ctx context.Context, filter domain.UsersPaginatedFilter) (domain.PaginatedResponse[domain.User], error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	PatchUser(ctx context.Context, id uuid.UUID, editableUser domain.EditableUserPatch) (domain.User, error)
	UpdateUserPassword(ctx context.Context, username domain.Username, oldPassword, newPassword domain.Password) error
	ResetUserPassword(ctx context.Context, username domain.Username, newPassword domain.Password) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	SignInUser(ctx context.Context, username domain.Username, password domain.Password) (string, error)

	CreateUserContainerBookmark(ctx context.Context, userID, containerID uuid.UUID) error
	ListUserContainerBookmarks(ctx context.Context, userID uuid.UUID, filter domain.UserContainerBookmarksPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	DeleteUserContainerBookmark(ctx context.Context, userID, containerID uuid.UUID) error

	CreateEmployee(ctx context.Context, editableEmployee domain.EditableEmployeeWithPassword) (domain.Employee, error)
	ListEmployees(ctx context.Context, filter domain.EmployeesPaginatedFilter) (domain.PaginatedResponse[domain.Employee], error)
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error)
	PatchEmployee(ctx context.Context, id uuid.UUID, editableEmployee domain.EditableEmployeePatch) (domain.Employee, error)
	UpdateEmployeePassword(ctx context.Context, username domain.Username, oldPassword, newPassword domain.Password) error
	ResetEmployeePassword(ctx context.Context, username domain.Username, newPassword domain.Password) error
	DeleteEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error)
	SignInEmployee(ctx context.Context, username domain.Username, password domain.Password) (string, error)

	CreateContainer(ctx context.Context, editableContainer domain.EditableContainer) (domain.Container, error)
	ListContainers(ctx context.Context, filter domain.ContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	GetContainerByID(ctx context.Context, id uuid.UUID) (domain.Container, error)
	PatchContainer(ctx context.Context, id uuid.UUID, editableContainer domain.EditableContainerPatch) (domain.Container, error)
	DeleteContainerByID(ctx context.Context, id uuid.UUID) (domain.Container, error)

	CreateTruck(ctx context.Context, editableTruck domain.EditableTruck) (domain.Truck, error)
	ListTrucks(ctx context.Context, filter domain.TrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error)
	GetTruckByID(ctx context.Context, id uuid.UUID) (domain.Truck, error)
	PatchTruck(ctx context.Context, id uuid.UUID, editableTruck domain.EditableTruckPatch) (domain.Truck, error)
	DeleteTruckByID(ctx context.Context, id uuid.UUID) (domain.Truck, error)

	CreateWarehouse(ctx context.Context, editableWarehouse domain.EditableWarehouse) (domain.Warehouse, error)
	ListWarehouses(ctx context.Context, filter domain.WarehousesPaginatedFilter) (domain.PaginatedResponse[domain.Warehouse], error)
	GetWarehouseByID(ctx context.Context, id uuid.UUID) (domain.Warehouse, error)
	PatchWarehouse(ctx context.Context, id uuid.UUID, editableWarehouse domain.EditableWarehousePatch) (domain.Warehouse, error)
	DeleteWarehouseByID(ctx context.Context, id uuid.UUID) (domain.Warehouse, error)

	CreateLandfill(ctx context.Context, editableLandfill domain.EditableLandfill) (domain.Landfill, error)
	ListLandfills(ctx context.Context, filter domain.LandfillsPaginatedFilter) (domain.PaginatedResponse[domain.Landfill], error)
	GetLandfillByID(ctx context.Context, id uuid.UUID) (domain.Landfill, error)
	PatchLandfill(ctx context.Context, id uuid.UUID, editableLandfill domain.EditableLandfillPatch) (domain.Landfill, error)
	DeleteLandfillByID(ctx context.Context, id uuid.UUID) (domain.Landfill, error)

	CreateRoute(ctx context.Context, editableRoute domain.EditableRoute) (domain.Route, error)
	ListRoutes(ctx context.Context, filter domain.RoutesPaginatedFilter) (domain.PaginatedResponse[domain.Route], error)
	GetRouteByID(ctx context.Context, id uuid.UUID) (domain.Route, error)
	PatchRoute(ctx context.Context, id uuid.UUID, editableRoute domain.EditableRoutePatch) (domain.Route, error)
	DeleteRouteByID(ctx context.Context, id uuid.UUID) (domain.Route, error)

	CreateRouteContainer(ctx context.Context, routeID, containerID uuid.UUID) error
	ListRouteContainers(ctx context.Context, routeID uuid.UUID, filter domain.RouteContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	DeleteRouteContainer(ctx context.Context, routeID, containerID uuid.UUID) error

	CreateRouteEmployee(ctx context.Context, routeID, employeeID uuid.UUID, editableRouteEmployee domain.EditableRouteEmployee) error
	ListRouteEmployees(ctx context.Context, routeID uuid.UUID, filter domain.RouteEmployeesPaginatedFilter) (domain.PaginatedResponse[domain.RouteEmployee], error)
	DeleteRouteEmployee(ctx context.Context, routeID, employeeID uuid.UUID) error

	GetRoadByGeometry(ctx context.Context, geometry domain.GeoJSONGeometryPoint) (domain.Road, error)

	GetMunicipalityByGeometry(ctx context.Context, geometry domain.GeoJSONGeometryPoint) (domain.Municipality, error)
}

// handler defines the http handler structure.
type handler struct {
	authzService AuthorizationService
	service      Service
	handler      http.Handler
}

// New returns a new http handler.
func New(authzService AuthorizationService, service Service) *handler {
	h := &handler{
		authzService: authzService,
		service:      service,
	}

	router := http.NewServeMux()

	// Handle web application.
	webAppFS := http.FileServer(http.Dir(dirWebApp))
	router.HandleFunc(baseURLWebApp, func(w http.ResponseWriter, r *http.Request) {
		// Handle single-page application routing.
		if r.URL.Path != baseURLWebApp && strings.Contains(r.Header.Get(requestHeaderKeyAccept), requestHeaderValueAcceptHTML) {
			http.ServeFile(w, r, path.Join(dirWebApp, dirIndexHTML))
			return
		}

		// Handle base file server.
		webAppFS.ServeHTTP(w, r)
	})

	// Handle API.
	authzMiddleware := authzService.Middleware(authz.MiddlewareOptions{
		UnauthorizedHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			switch {
			case errors.Is(err, authz.ErrAuthorizationHeaderInvalid):
				unauthorized(w, errAuthorizationHeaderInvalid)
			case errors.Is(err, authz.ErrJWTInvalid):
				unauthorized(w, errJWTInvalid)
			default:
				unauthorized(w, err.Error())
			}
		},
		ForbiddenHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			switch {
			case errors.Is(err, authz.ErrRolesInvalid):
				forbidden(w, errRolesInvalid)
			case errors.Is(err, authz.ErrAuthorizationInvalid):
				forbidden(w, errAuthorizationInvalid)
			default:
				forbidden(w, err.Error())
			}
		},
	})

	h.handler = spec.HandlerWithOptions(h, spec.StdHTTPServerOptions{
		BaseURL:     baseURLApi,
		BaseRouter:  router,
		Middlewares: []spec.MiddlewareFunc{authzMiddleware},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			var specInvalidParamFormatError *spec.InvalidParamFormatError

			switch {
			case errors.As(err, &specInvalidParamFormatError):
				badRequest(w, fmt.Sprintf("%s: %s", errParamInvalidFormat, specInvalidParamFormatError.ParamName))
			default:
				badRequest(w, err.Error())
			}
		},
	})

	// Handle swagger documentation.
	swaggerFS := http.FileServer(http.Dir(dirSwaggerUI))
	router.Handle(baseURLDocs, http.StripPrefix(baseURLDocs, swaggerFS))

	return h
}

// ServeHTTP responds to an http request.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

// setHeaderJSON sets the header with the content type json.
func setHeaderJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// writeResponseJSON writes the data to the response and sets the header with the provided status code and content type
// json.
func writeResponseJSON(w http.ResponseWriter, statusCode int, data []byte) {
	setHeaderJSON(w)
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}
