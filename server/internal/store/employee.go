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
func (s *store) CreateEmployee(ctx context.Context, tx pgx.Tx, editableEmployee domain.EditableEmployeeWithPassword, roadID, municipalityID *int) (uuid.UUID, error) {
	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableEmployee.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	geoJSON, err := json.Marshal(geometry)
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
		editableEmployee.Role,
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
		if getConstraintName(err) == constraintEmployeesUsernameKey {
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrEmployeeAlreadyExists)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
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
		var geometry domain.GeoJSONGeometryPoint
		if feature, ok := editableEmployee.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}

		geoJSON, err = json.Marshal(geometry)
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
		if getConstraintName(err) == constraintEmployeesUsernameKey {
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeAlreadyExists)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrEmployeeNotFound)
	}

	return nil
}

// getEmployeeFromRow returns the employee by scanning the given row.
func getEmployeeFromRow(row pgx.Row) (domain.Employee, error) {
	var employee domain.Employee
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.Role,
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
