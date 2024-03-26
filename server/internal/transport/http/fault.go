package http

import (
	"encoding/json"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// Common fault messages.
const (
	errRequestBodyInvalid   = "invalid request body"
	errIncorrectCredentials = "incorrect credentials"
)

// Common fault descriptions.
const (
	descriptionFailedToMarshalResponseBody = "http: failed to marshal response body"
)

// Fault code const.
const (
	faultCodeBadRequest          = "bad_request"
	faultCodeUnauthorized        = "unauthorized"
	faultCodeNotFound            = "not_found"
	faultCodeInternalServerError = "internal_server_error"
)

// badRequest writes an error response and sets the header with the bad request status code.
func badRequest(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusBadRequest, faultCodeBadRequest, message)
}

// unauthorized writes an error response and sets the header with the unauthorized status code.
func unauthorized(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusUnauthorized, faultCodeUnauthorized, message)
}

// notFound writes an error response and sets the header with the not found status code.
func notFound(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusNotFound, faultCodeNotFound, message)
}

// internalServerError sets the header with the internal server error status code.
func internalServerError(w http.ResponseWriter) {
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
