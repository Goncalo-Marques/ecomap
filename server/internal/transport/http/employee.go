package http

import (
	"encoding/json"
	"errors"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errEmployeeNotFound = "employee not found"
)

// GetEmployeeByID handles the http request that gets the employee by id.
func (h *handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdParam) {
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

	employee := fromDomainEmployee(domainEmployee)

	responseBody, err := json.Marshal(employee)
	if err != nil {
		logging.Logger.ErrorContext(ctx, errFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}
