package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errEmployeeNotFound      = "employee not found"
	errEmployeeAlreadyExists = "username already exists"
)

// CreateEmployee handles the http request to create an employee.
func (h *handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var employeePost spec.EmployeePost
	err = json.Unmarshal(requestBody, &employeePost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableEmployee, err := employeePostToDomain(employeePost)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainEmployee, err := h.service.CreateEmployee(ctx, domainEditableEmployee)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrEmployeeAlreadyExists):
			conflict(w, errEmployeeAlreadyExists)
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

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListEmployees handles the http request to list employees.
func (h *handler) ListEmployees(w http.ResponseWriter, r *http.Request, params spec.ListEmployeesParams) {
	ctx := r.Context()

	domainEmployeesFilter := listEmployeesParamsToDomain(params)
	domainPaginatedEmployees, err := h.service.ListEmployees(ctx, domainEmployeesFilter)
	if err != nil {
		var domainErrFilterValueInvalid *domain.ErrFilterValueInvalid

		switch {
		case errors.As(err, &domainErrFilterValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFilterValueInvalid, domainErrFilterValueInvalid.FilterName))
		default:
			internalServerError(w)
		}

		return
	}

	employeesPaginated, err := employeesPaginatedFromDomain(domainPaginatedEmployees)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(employeesPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
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
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var employeePatch spec.EmployeePatch
	err = json.Unmarshal(requestBody, &employeePatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableEmployee, err := employeePatchToDomain(employeePatch)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainEmployee, err := h.service.PatchEmployee(ctx, employeeID, domainEditableEmployee)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrEmployeeNotFound):
			notFound(w, errEmployeeNotFound)
		case errors.Is(err, domain.ErrEmployeeAlreadyExists):
			conflict(w, errEmployeeAlreadyExists)
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

// DeleteEmployeeByID handles the http request to delete an employee by ID.
func (h *handler) DeleteEmployeeByID(w http.ResponseWriter, r *http.Request, employeeID spec.EmployeeIdPathParam) {
	ctx := r.Context()

	domainEmployee, err := h.service.DeleteEmployeeByID(ctx, employeeID)
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

// UpdateEmployeePassword handles the http request to update an employee password.
func (h *handler) UpdateEmployeePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var passwordChange spec.PasswordChange
	err = json.Unmarshal(requestBody, &passwordChange)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	err = h.service.UpdateEmployeePassword(ctx, domain.Username(passwordChange.Username), domain.Password(passwordChange.OldPassword), domain.Password(passwordChange.NewPassword))
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrCredentialsIncorrect):
			unauthorized(w, errCredentialsIncorrect)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// ResetEmployeePassword handles the http request to reset an employee password.
func (h *handler) ResetEmployeePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var passwordReset spec.PasswordReset
	err = json.Unmarshal(requestBody, &passwordReset)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	err = h.service.ResetEmployeePassword(ctx, domain.Username(passwordReset.Username), domain.Password(passwordReset.NewPassword))
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrEmployeeNotFound):
			notFound(w, errEmployeeNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
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

// employeePostToDomain returns a domain editable employee with password based on the standardized employee post.
func employeePostToDomain(employeePost spec.EmployeePost) (domain.EditableEmployeeWithPassword, error) {
	if len(employeePost.GeoJson.Geometry.Coordinates) != 2 {
		return domain.EditableEmployeeWithPassword{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldGeoJSON}
	}

	scheduleStart, err := timeStringToTime(employeePost.ScheduleStart)
	if err != nil {
		return domain.EditableEmployeeWithPassword{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleStart}
	}

	scheduleEnd, err := timeStringToTime(employeePost.ScheduleEnd)
	if err != nil {
		return domain.EditableEmployeeWithPassword{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleEnd}
	}

	var role domain.EmployeeRole
	switch employeePost.Role {
	case spec.WasteOperator:
		role = domain.EmployeeRoleWasteOperator
	case spec.Manager:
		role = domain.EmployeeRoleManager
	}

	return domain.EditableEmployeeWithPassword{
		EditableEmployee: domain.EditableEmployee{
			Username:    domain.Username(employeePost.Username),
			FirstName:   domain.Name(employeePost.FirstName),
			LastName:    domain.Name(employeePost.LastName),
			Role:        role,
			DateOfBirth: employeePost.DateOfBirth.Time,
			PhoneNumber: domain.PhoneNumber(employeePost.PhoneNumber),
			GeoJSON: domain.GeoJSONFeature{
				Geometry: domain.GeoJSONGeometryPoint{
					Coordinates: [2]float64(employeePost.GeoJson.Geometry.Coordinates),
				},
			},
			ScheduleStart: scheduleStart,
			ScheduleEnd:   scheduleEnd,
		},
		Password: domain.Password(employeePost.Password),
	}, nil
}

// employeePatchToDomain returns a domain patchable employee based on the standardized employee patch.
func employeePatchToDomain(employeePatch spec.EmployeePatch) (domain.EditableEmployeePatch, error) {
	var dateOfBirth *time.Time
	if employeePatch.DateOfBirth != nil {
		dateOfBirth = &employeePatch.DateOfBirth.Time
	}

	var geoJSON domain.GeoJSON
	if employeePatch.GeoJson != nil {
		if len(employeePatch.GeoJson.Geometry.Coordinates) != 2 {
			return domain.EditableEmployeePatch{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldGeoJSON}
		}

		geoJSON = domain.GeoJSONFeature{
			Geometry: domain.GeoJSONGeometryPoint{
				Coordinates: [2]float64(employeePatch.GeoJson.Geometry.Coordinates),
			},
		}
	}

	var scheduleStart *time.Time
	if employeePatch.ScheduleStart != nil {
		schedule, err := timeStringToTime(*employeePatch.ScheduleStart)
		if err != nil {
			return domain.EditableEmployeePatch{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleStart}
		}

		scheduleStart = &schedule
	}

	var scheduleEnd *time.Time
	if employeePatch.ScheduleEnd != nil {
		schedule, err := timeStringToTime(*employeePatch.ScheduleEnd)
		if err != nil {
			return domain.EditableEmployeePatch{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleEnd}
		}

		scheduleEnd = &schedule
	}

	return domain.EditableEmployeePatch{
		Username:      (*domain.Username)(employeePatch.Username),
		FirstName:     (*domain.Name)(employeePatch.FirstName),
		LastName:      (*domain.Name)(employeePatch.LastName),
		DateOfBirth:   dateOfBirth,
		PhoneNumber:   (*domain.PhoneNumber)(employeePatch.PhoneNumber),
		GeoJSON:       geoJSON,
		ScheduleStart: scheduleStart,
		ScheduleEnd:   scheduleEnd,
	}, nil
}

// listEmployeesParamsToDomain returns a domain employees paginated filter based on the standardized list employees parameters.
func listEmployeesParamsToDomain(params spec.ListEmployeesParams) domain.EmployeesPaginatedFilter {
	var domainSort domain.PaginationSort[domain.EmployeePaginatedSort]
	if params.Sort != nil {
		domainSort = domain.EmployeePaginatedSort(*params.Sort)
	}

	var domainRole *domain.EmployeeRole
	if params.Role != nil {
		role := domain.EmployeeRole(*params.Role)

		switch *params.Role {
		case spec.WasteOperator:
			role = domain.EmployeeRoleWasteOperator
		case spec.Manager:
			role = domain.EmployeeRoleManager
		}

		domainRole = &role
	}

	var domainDateOfBirth *string
	if params.DateOfBirth != nil {
		timeString := dateStringFromTime(params.DateOfBirth.Time)
		domainDateOfBirth = &timeString
	}

	return domain.EmployeesPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			(*spec.LogicalOperatorQueryParam)(params.LogicalOperator),
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		Username:         (*domain.Username)(params.Username),
		FirstName:        (*domain.Name)(params.FirstName),
		LastName:         (*domain.Name)(params.LastName),
		Role:             domainRole,
		DateOfBirth:      domainDateOfBirth,
		PhoneNumber:      (*domain.PhoneNumber)(params.PhoneNumber),
		ScheduleStart:    params.ScheduleStart,
		ScheduleEnd:      params.ScheduleEnd,
		WayName:          params.WayName,
		MunicipalityName: params.MunicipalityName,
	}
}

// employeeFromDomain returns a standardized employee based on the domain model.
func employeeFromDomain(employee domain.Employee) (spec.Employee, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(employee.GeoJSON)
	if err != nil {
		return spec.Employee{}, err
	}

	return spec.Employee{
		Id:            employee.ID,
		Username:      string(employee.Username),
		FirstName:     string(employee.FirstName),
		LastName:      string(employee.LastName),
		Role:          spec.EmployeeRole(employee.Role),
		DateOfBirth:   dateFromTime(employee.DateOfBirth),
		PhoneNumber:   string(employee.PhoneNumber),
		GeoJson:       geoJSON,
		ScheduleStart: timeStringFromTime(employee.ScheduleStart),
		ScheduleEnd:   timeStringFromTime(employee.ScheduleEnd),
		CreatedAt:     employee.CreatedAt,
		ModifiedAt:    employee.ModifiedAt,
	}, nil
}

// employeesFromDomain returns standardized employees based on the domain model.
func employeesFromDomain(employees []domain.Employee) ([]spec.Employee, error) {
	specEmployees := make([]spec.Employee, len(employees))
	var err error

	for i, employee := range employees {
		specEmployees[i], err = employeeFromDomain(employee)
		if err != nil {
			return []spec.Employee{}, err
		}
	}

	return specEmployees, nil
}

// employeesPaginatedFromDomain returns a standardized employees paginated response based on the domain model.
func employeesPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Employee]) (spec.EmployeesPaginated, error) {
	employees, err := employeesFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.EmployeesPaginated{}, err
	}

	return spec.EmployeesPaginated{
		Total:     paginatedResponse.Total,
		Employees: employees,
	}, nil
}
