package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// TODO: Replicate User methods.

// GetEmployeeByID executes a query to return the employee with the specified identifier.
func (s *store) GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, first_name, last_name, role, date_of_birth, phone_number, ST_AsGeoJSON(geom), schedule_start, schedule_end, created_at, modified_at 
		FROM employees 
		WHERE id = $1 
	`,
		id,
	)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	employees, err := getEmployeesFromRows(rows)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	if len(employees) == 0 {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, domain.ErrEmployeeNotFound)
	}

	return employees[0], nil
}

// GetEmployeeByUsername executes a query to return the employee with the specified username.
func (s *store) GetEmployeeByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.Employee, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, first_name, last_name, role, date_of_birth, phone_number, ST_AsGeoJSON(geom), schedule_start, schedule_end, created_at, modified_at 
		FROM employees 
		WHERE username = $1 
	`,
		username,
	)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	employees, err := getEmployeesFromRows(rows)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	if len(employees) == 0 {
		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, domain.ErrEmployeeNotFound)
	}

	return employees[0], nil
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

// getEmployeesFromRows returns the employees by scanning the given rows.
func getEmployeesFromRows(rows pgx.Rows) ([]domain.Employee, error) {
	var employees []domain.Employee
	var err error

	for rows.Next() {
		var employee domain.Employee

		err = rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.Role,
			&employee.DateOfBirth,
			&employee.PhoneNumber,
			&employee.Geom,
			&employee.ScheduleStart,
			&employee.ScheduleEnd,
			&employee.CreatedTime,
			&employee.ModifiedTime,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}
