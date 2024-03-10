package http

import (
	"context"
	"encoding/json"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

// Fault code const.
const (
	faultCodeBadRequest          = "bad_request"
	faultCodeNotFound            = "not_found"
	faultCodeInternalServerError = "internal_server_error"
)

// badRequest writes an error response and sets the header with the bad request status code.
func badRequest(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusBadRequest, faultCodeBadRequest, message)
}

// notFound writes an error response and sets the header with the not found status code.
func notFound(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusNotFound, faultCodeNotFound, message)
}

// internalServerError logs the error and sets the header with the internal server error status code.
func internalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	logging.Logger.ErrorContext(ctx, "http: internal server error", logging.Error(err))
	w.WriteHeader(http.StatusInternalServerError)
}

// fault writes an error response and sets the header with the provided status code and content type json.
func fault(w http.ResponseWriter, statusCode int, code, message string) error {
	setHeaderJSON(w)
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	return enc.Encode(spec.Error{
		Code:    code,
		Message: message,
	})
}
