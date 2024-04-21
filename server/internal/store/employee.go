package store

import (
	"context"
	"encoding/json"
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
func (s *store) CreateEmployee(ctx context.Context, tx pgx.Tx, editableEmployee domain.EditableEmployeeWithPassword) (domain.Employee, error) {
	var feature domain.GeoJSONFeature
	if f, ok := editableEmployee.GeoJSON.(domain.GeoJSONFeature); ok {
		feature = f
	}

	var geometry domain.GeoJSONGeometryPoint
	if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
		geometry = g
	}

	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	wayOSM, err := s.GetWayOSMByGeoJSON(ctx, tx, geometry)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedGetWayOSM, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO employees (username, password, first_name, last_name, role, date_of_birth, phone_number, geom, way_osm_id, schedule_start, schedule_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, ST_GeomFromGeoJSON($8), $9, $10, $11, $12, $13) 
		RETURNING id, username, first_name, last_name, role, date_of_birth, phone_number, ST_AsGeoJSON(geom)::jsonb, way_osm_id, schedule_start, schedule_end, created_at, modified_at 
	`,
		editableEmployee.Username,
		editableEmployee.Password,
		editableEmployee.FirstName,
		editableEmployee.LastName,
		editableEmployee.Role,
		editableEmployee.DateOfBirth,
		editableEmployee.PhoneNumber,
		geoJSON,
		wayOSM,
		editableEmployee.ScheduleStart,
		editableEmployee.ScheduleEnd,
	)

	employee, err := getEmployeeFromRow(row)
	if err != nil {
		if getConstraintName(err) == constraintEmployeesUsernameKey {
			return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeAlreadyExists)
		}

		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return employee, nil
}

// GetEmployeeByID executes a query to return the employee with the specified identifier.
func (s *store) GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, username, first_name, last_name, role, date_of_birth, phone_number, ST_AsGeoJSON(geom)::jsonb, way_osm_id, schedule_start, schedule_end, created_at, modified_at 
		FROM employees 
		WHERE id = $1 
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
		SELECT id, username, first_name, last_name, role, date_of_birth, phone_number, ST_AsGeoJSON(geom)::jsonb, way_osm_id, schedule_start, schedule_end, created_at, modified_at 
		FROM employees 
		WHERE username = $1 
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

// getEmployeeFromRow returns the employee by scanning the given row.
func getEmployeeFromRow(row pgx.Row) (domain.Employee, error) {
	var employee domain.Employee
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayOSM *int

	err := row.Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.Role,
		&employee.DateOfBirth,
		&employee.PhoneNumber,
		&geoJSONPoint,
		&wayOSM,
		&employee.ScheduleStart,
		&employee.ScheduleEnd,
		&employee.CreatedAt,
		&employee.ModifiedAt,
	)
	if err != nil {
		return domain.Employee{}, err
	}

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayOSM != nil {
		geoJSONProperties.SetWayOSM(*wayOSM)
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
