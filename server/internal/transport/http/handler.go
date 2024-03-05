package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
)

// Service defines the service interface.
type Service interface{}

type handler struct {
	handler http.Handler
	service Service
}

func New(service Service) *handler {
	h := &handler{
		service: service,
	}

	h.handler = spec.HandlerWithOptions(h, spec.StdHTTPServerOptions{
		BaseRouter: http.DefaultServeMux,
	})

	return h
}
