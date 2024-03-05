package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/swagger/ecomap"
)

func (h *handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request, employeeIDParam spec.EmployeeIdParam) {
}
