package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// Base url const.
const (
	baseURLBackOffice = "/"
	baseURLAPI        = "/api"
	baseURLDocs       = "/docs"
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

	h.handler = spec.HandlerWithOptions(h, spec.StdHTTPServerOptions{
		BaseURL:    baseURLAPI,
		BaseRouter: http.DefaultServeMux,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			badRequest(w, err.Error())
		},
	})

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
