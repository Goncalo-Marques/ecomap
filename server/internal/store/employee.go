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
	errGetEmployeeByID = "store: failed to get employee by id"
)

// GetEmployeeByID executes a db query to return the employee by id.
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
			return domain.Employee{}, fmt.Errorf("%s: %w", errGetEmployeeByID, domain.ErrEmployeeNotFound)
		}

		return domain.Employee{}, fmt.Errorf("%s: %w", errGetEmployeeByID, err)
	}

	return employee, nil
}
