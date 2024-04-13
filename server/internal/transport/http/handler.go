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

	errAuthorizationHeaderInvalid = "invalid authorization header"
	errJWTInvalid                 = "invalid jwt"
	errRolesInvalid               = "invalid subject roles"
	errAuthorizationInvalid       = "unauthorized subject"
	errParamInvalidFormat         = "invalid parameter format"
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

	GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error)
	SignInEmployee(ctx context.Context, username domain.Username, password domain.Password) (string, error)
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
