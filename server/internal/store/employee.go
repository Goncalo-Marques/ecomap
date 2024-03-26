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
	descriptionFailedGetEmployeeSignIn = "store: failed to get employee sign-in"
	descriptionFailedGetEmployeeByID   = "store: failed to get employee by id"
)

// GetEmployeeSignIn executes a query to return the sign-in of the employee with the specified username.
func (s *store) GetEmployeeSignIn(ctx context.Context, tx pgx.Tx, username string) (domain.SignIn, error) {
	row := tx.QueryRow(ctx, `
		SELECT username, password FROM public.employees 
		WHERE username = $1 
	`, username)

	var signIn domain.SignIn

	err := row.Scan(
		&signIn.Username,
		&signIn.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedGetEmployeeSignIn, domain.ErrEmployeeNotFound)
		}

		return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedGetEmployeeSignIn, err)
	}

	return signIn, nil
}

// GetEmployeeByID executes a query to return the employee with the specified identifier.
func (s *store) GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, name, date_of_birth
		FROM employee
		WHERE id = $1
	`, id)

	var employee domain.Employee

	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.DateOfBirth,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedGetEmployeeByID, domain.ErrEmployeeNotFound)
		}

		return domain.Employee{}, fmt.Errorf("%s: %w", descriptionFailedGetEmployeeByID, err)
	}

	return employee, nil
}
