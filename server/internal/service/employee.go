package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errGetEmployeeByID = "service: failed to get employee by id"
)

// GetEmployeeByID returns the employee by id.
func (s *service) GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	logAttrs := []slog.Attr{
		slog.String(logging.EmployeeID, id.String()),
	}

	var employee domain.Employee
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		employee, err = s.store.GetEmployeeByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return domain.Employee{}, logInfoAndWrapError(ctx, err, errGetEmployeeByID, logAttrs)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, errGetEmployeeByID, logAttrs)
		}
	}

	return employee, nil
}
