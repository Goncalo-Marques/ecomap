package http

import (
	"encoding/json"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// Common fault messages.
const (
	errRequestBodyInvalid   = "invalid request body"
	errFieldValueInvalid    = "invalid field value"
	errFilterValueInvalid   = "invalid filter value"
	errCredentialsIncorrect = "incorrect credentials"
)

// Common fault descriptions.
const (
	descriptionFailedToMarshalResponseBody = "http: failed to marshal response body"
)

// badRequest writes an error response and sets the header with the bad request status code.
func badRequest(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusBadRequest, spec.ErrorCodeBadRequest, &message)
}

// unauthorized writes an error response and sets the header with the unauthorized status code.
func unauthorized(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusUnauthorized, spec.ErrorCodeUnauthorized, &message)
}

// forbidden writes an error response and sets the header with the forbidden status code.
func forbidden(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusForbidden, spec.ErrorCodeForbidden, &message)
}

// notFound writes an error response and sets the header with the not found status code.
func notFound(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusNotFound, spec.ErrorCodeNotFound, &message)
}

// conflict writes an error response and sets the header with the conflict status code.
func conflict(w http.ResponseWriter, message string) {
	_ = fault(w, http.StatusConflict, spec.ErrorCodeConflict, &message)
}

// internalServerError sets the header with the internal server error status code.
func internalServerError(w http.ResponseWriter) {
	_ = fault(w, http.StatusInternalServerError, spec.ErrorCodeInternalServerError, nil)
}

// fault writes an error response and sets the header with the provided status code and content type json.
func fault(w http.ResponseWriter, statusCode int, code spec.ErrorCode, message *string) error {
	setHeaderJSON(w)
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	return enc.Encode(spec.Error{
		Code:    code,
		Message: message,
	})
}
