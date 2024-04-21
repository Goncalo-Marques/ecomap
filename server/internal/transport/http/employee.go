package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errEmployeeNotFound = "employee not found"
)

// CreateEmployee handles the http request to create an employee.
func (h *handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ListEmployees handles the http request to list employees.
func (h *handler) ListEmployees(w http.ResponseWriter, r *http.Request, params spec.ListEmployeesParams) {
	w.WriteHeader(http.StatusNotFound)
}

// GetEmployeeByID handles the http request to get an employee by id.
func (h *handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdPathParam) {
	ctx := r.Context()

	domainEmployee, err := h.service.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			notFound(w, errEmployeeNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	employee, err := employeeFromDomain(domainEmployee)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(employee)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchEmployeeByID handles the http request to modify an employee by ID.
func (h *handler) PatchEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteEmployeeByID handles the http request to delete an employee by ID.
func (h *handler) DeleteEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// UpdateEmployeePassword handles the http request to update an employee password.
func (h *handler) UpdateEmployeePassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// ResetEmployeePassword handles the http request to reset an employee password.
func (h *handler) ResetEmployeePassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// SignInEmployee handles the http request to sign in an employee.
func (h *handler) SignInEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var signIn spec.SignIn
	err = json.Unmarshal(requestBody, &signIn)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	token, err := h.service.SignInEmployee(ctx, domain.Username(signIn.Username), domain.Password(signIn.Password))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCredentialsIncorrect):
			unauthorized(w, errCredentialsIncorrect)
		default:
			internalServerError(w)
		}

		return
	}

	jwt := jwtFromJWTToken(token)
	responseBody, err := json.Marshal(jwt)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}
