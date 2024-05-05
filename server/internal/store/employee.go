package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintEmployeesUsernameKey = "employees_username_key"
)

// CreateEmployee executes a query to create an employee with the specified data.
func (s *store) CreateEmployee(ctx context.Context, tx pgx.Tx, editableEmployee domain.EditableEmployeeWithPassword, roadID, municipalityID *int) (uuid.UUID, error) {
	geoJSON, err := jsonMarshalGeoJSONGeometryPoint(editableEmployee.GeoJSON)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO employees (username, password, first_name, last_name, role, date_of_birth, phone_number, geom, road_id, municipality_id, schedule_start, schedule_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, ST_GeomFromGeoJSON($8), $9, $10, $11, $12) 
		RETURNING id
	`,
		editableEmployee.Username,
		editableEmployee.Password,
		editableEmployee.FirstName,
		editableEmployee.LastName,
		employeeRoleFromDomain(editableEmployee.Role),
		editableEmployee.DateOfBirth,
		editableEmployee.PhoneNumber,
		geoJSON,
		roadID,
		municipalityID,
		editableEmployee.ScheduleStart,
		editableEmployee.ScheduleEnd,
	)

	var id uuid.UUID

	err = row.Scan(&id)
	if err != nil {
		if constraintNameFromError(err) == constraintEmployeesUsernameKey {
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeAlreadyExists)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}

// ListEmployees executes a query to return the employees for the specified filter.
func (s *store) ListEmployees(ctx context.Context, tx pgx.Tx, filter domain.EmployeesPaginatedFilter) (domain.PaginatedResponse[domain.Employee], error) {
	var filterFields []string
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	if filter.Username != nil {
		filterFields = append(filterFields, "e.username")
		argsWhere = append(argsWhere, *filter.Username)
	}
	if filter.FirstName != nil {
		filterFields = append(filterFields, "e.first_name")
		argsWhere = append(argsWhere, *filter.FirstName)
	}
	if filter.LastName != nil {
		filterFields = append(filterFields, "e.last_name")
		argsWhere = append(argsWhere, *filter.LastName)
	}
	if filter.Role != nil {
		filterFields = append(filterFields, "e.role::text")
		argsWhere = append(argsWhere, employeeRoleFromDomain(*filter.Role))
	}
	if filter.DateOfBirth != nil {
		filterFields = append(filterFields, "e.date_of_birth::text")
		argsWhere = append(argsWhere, *filter.DateOfBirth)
	}
	if filter.PhoneNumber != nil {
		filterFields = append(filterFields, "e.phone_number")
		argsWhere = append(argsWhere, *filter.PhoneNumber)
	}
	if filter.ScheduleStart != nil {
		filterFields = append(filterFields, "e.schedule_start::text")
		argsWhere = append(argsWhere, *filter.ScheduleStart)
	}
	if filter.ScheduleEnd != nil {
		filterFields = append(filterFields, "e.schedule_end::text")
		argsWhere = append(argsWhere, *filter.ScheduleEnd)
	}
	if filter.LocationName != nil {
		filterLocationFields = []string{"rn.osm_name", "m.name"}
		argsWhere = append(argsWhere, *filter.LocationName)
	}

	sqlWhere := listSQLWhere(filterFields, filterLocationFields)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(e.id) 
		FROM employees AS e
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Employee]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.EmployeePaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "e.created_at"
	switch domainSortField {
	case domain.EmployeePaginatedSortUsername:
		sortField = "e.username"
	case domain.EmployeePaginatedSortFirstName:
		sortField = "e.first_name"
	case domain.EmployeePaginatedSortLastName:
		sortField = "e.last_name"
	case domain.EmployeePaginatedSortRole:
		sortField = "e.role"
	case domain.EmployeePaginatedSortDateOfBirth:
		sortField = "e.date_of_birth"
	case domain.EmployeePaginatedSortScheduleStart:
		sortField = "e.schedule_start"
	case domain.EmployeePaginatedSortScheduleEnd:
		sortField = "e.schedule_end"
	case domain.EmployeePaginatedSortWayName:
		sortField = "rn.osm_name"
	case domain.EmployeePaginatedSortMunicipalityName:
		sortField = "m.name"
	case domain.EmployeePaginatedSortCreatedAt:
		sortField = "e.created_at"
	case domain.EmployeePaginatedSortModifiedAt:
		sortField = "e.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT e.id, e.username, e.first_name, e.last_name, e.role, e.date_of_birth, e.phone_number, ST_AsGeoJSON(e.geom)::jsonb, rn.osm_name, m.name, e.schedule_start, e.schedule_end, e.created_at, e.modified_at 
		FROM employees AS e
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Employee]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	employees, err := getEmployeesFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Employee]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Employee]{
		Total:   total,
		Results: employees,
	}, nil
}

// GetEmployeeByID executes a query to return the employee with the specified identifier.
func (s *store) GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error) {
	row := tx.QueryRow(ctx, `
		SELECT e.id, e.username, e.first_name, e.last_name, e.role, e.date_of_birth, e.phone_number, ST_AsGeoJSON(e.geom)::jsonb, rn.osm_name, m.name, e.schedule_start, e.schedule_end, e.created_at, e.modified_at 
		FROM employees AS e
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
		WHERE e.id = $1 
	`,
		id,
	)

	employee, err := getEmployeeFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeNotFound)
		}

		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return employee, nil
}

// GetEmployeeByUsername executes a query to return the employee with the specified username.
func (s *store) GetEmployeeByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.Employee, error) {
	row := tx.QueryRow(ctx, `
		SELECT e.id, e.username, e.first_name, e.last_name, e.role, e.date_of_birth, e.phone_number, ST_AsGeoJSON(e.geom)::jsonb, rn.osm_name, m.name, e.schedule_start, e.schedule_end, e.created_at, e.modified_at 
		FROM employees AS e
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
		WHERE e.username = $1 
	`,
		username,
	)

	employee, err := getEmployeeFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeNotFound)
		}

		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return employee, nil
}

// GetEmployeeSignIn executes a query to return the sign-in of the employee with the specified username.
func (s *store) GetEmployeeSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error) {
	row := tx.QueryRow(ctx, `
		SELECT username, password 
		FROM employees 
		WHERE username = $1 
	`,
		username,
	)

	var signIn domain.SignIn
	err := row.Scan(
		&signIn.Username,
		&signIn.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeNotFound)
		}

		return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return signIn, nil
}

// PatchEmployee executes a query to patch an employee with the specified identifier and data.
func (s *store) PatchEmployee(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableEmployee domain.EditableEmployeePatch, roadID, municipalityID *int) error {
	var geoJSON []byte
	var err error

	if editableEmployee.GeoJSON != nil {
		geoJSON, err = jsonMarshalGeoJSONGeometryPoint(editableEmployee.GeoJSON)
		if err != nil {
			return fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE employees SET
			username = coalesce($2, username),
			first_name = coalesce($3, first_name),
			last_name = coalesce($4, last_name),
			date_of_birth = coalesce($5, date_of_birth),
			phone_number = coalesce($6, phone_number),
			geom = coalesce(ST_GeomFromGeoJSON($7), geom),
			road_id = CASE 
					WHEN $7 IS NOT NULL THEN $8 
					ELSE road_id
				END,
			municipality_id = CASE 
					WHEN $7 IS NOT NULL THEN $9 
					ELSE municipality_id
				END,
			schedule_start = coalesce($10, schedule_start),
			schedule_end = coalesce($11, schedule_end)
		WHERE id = $1
	`,
		id,
		editableEmployee.Username,
		editableEmployee.FirstName,
		editableEmployee.LastName,
		editableEmployee.DateOfBirth,
		editableEmployee.PhoneNumber,
		geoJSON,
		roadID,
		municipalityID,
		editableEmployee.ScheduleStart,
		editableEmployee.ScheduleEnd,
	)
	if err != nil {
		if constraintNameFromError(err) == constraintEmployeesUsernameKey {
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeAlreadyExists)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeNotFound)
	}

	return nil
}

// UpdateEmployeePassword executes a query to update the password of the employee with the specified username.
func (s *store) UpdateEmployeePassword(ctx context.Context, tx pgx.Tx, username domain.Username, password domain.Password) error {
	commandTag, err := tx.Exec(ctx, `
		UPDATE employees SET
			password = $2
		WHERE username = $1 
	`,
		username,
		password,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeNotFound)
	}

	return nil
}

// DeleteEmployeeByID executes a query to delete the employee with the specified identifier.
func (s *store) DeleteEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM employees
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintContainersReportsResolverIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeAssociatedWithContainerReportAsResolver)
		case constraintRoutesContainersResponsibleIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeAssociatedWithRouteContainerAsResponsible)
		case constraintRoutesEmployeesEmployeeIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeAssociatedWithRouteEmployee)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeNotFound)
	}

	return nil
}

// employeeRoleFromDomain returns a store employee role based on the domain model.
func employeeRoleFromDomain(role domain.EmployeeRole) string {
	switch role {
	case domain.EmployeeRoleWasteOperator:
		return "waste_operator"
	case domain.EmployeeRoleManager:
		return "manager"
	default:
		return string(role)
	}
}

// employeeRoleToDomain returns a domain employee role based on the store model.
func employeeRoleToDomain(role string) domain.EmployeeRole {
	switch role {
	case "waste_operator":
		return domain.EmployeeRoleWasteOperator
	case "manager":
		return domain.EmployeeRoleManager
	default:
		return domain.EmployeeRole(role)
	}
}

// getEmployeeFromRow returns the employee by scanning the given row.
func getEmployeeFromRow(row pgx.Row) (domain.Employee, error) {
	var employee domain.Employee
	var role string
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&role,
		&employee.DateOfBirth,
		&employee.PhoneNumber,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&employee.ScheduleStart,
		&employee.ScheduleEnd,
		&employee.CreatedAt,
		&employee.ModifiedAt,
	)
	if err != nil {
		return domain.Employee{}, err
	}

	employee.Role = employeeRoleToDomain(role)

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	employee.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return employee, nil
}

// getEmployeesFromRows returns the employees by scanning the given rows.
func getEmployeesFromRows(rows pgx.Rows) ([]domain.Employee, error) {
	var employees []domain.Employee
	for rows.Next() {
		employee, err := getEmployeeFromRow(rows)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}
