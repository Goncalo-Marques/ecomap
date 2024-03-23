package http

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/google/uuid"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
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
	requestHeaderAcceptKey       = "Accept"
	requestHeaderAcceptHTMLValue = "text/html"
)

// Service defines the service interface.
type Service interface {
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error)
}

// handler defines the http handler structure.
type handler struct {
	handler http.Handler
	service Service
}

// New returns a new http handler.
func New(service Service) *handler {
	h := &handler{
		service: service,
	}

	router := http.NewServeMux()

	// Handle web application.
	webAppFS := http.FileServer(http.Dir(dirWebApp))
	router.HandleFunc(baseURLWebApp, func(w http.ResponseWriter, r *http.Request) {
		// Handle single-page application routing.
		if r.URL.Path != baseURLWebApp && strings.Contains(r.Header.Get(requestHeaderAcceptKey), requestHeaderAcceptHTMLValue) {
			http.ServeFile(w, r, path.Join(dirWebApp, dirIndexHTML))
			return
		}

		// Handle base file server.
		webAppFS.ServeHTTP(w, r)
	})

	// Handle API.
	h.handler = spec.HandlerWithOptions(h, spec.StdHTTPServerOptions{
		BaseURL:    baseURLApi,
		BaseRouter: router,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			badRequest(w, err.Error())
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
