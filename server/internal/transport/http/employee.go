package http

import (
	"encoding/json"
	"errors"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	errEmployeeNotFound = "employee not found"
)

// GetEmployeeByID returns the employee by id.
func (h *handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdParam) {
	ctx := r.Context()

	domainEmployee, err := h.service.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			notFound(w, errEmployeeNotFound)
		default:
			internalServerError(ctx, w, err)
		}

		return
	}

	employee := fromDomainEmployee(domainEmployee)

	responseBody, err := json.Marshal(employee)
	if err != nil {
		internalServerError(ctx, w, err)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}
