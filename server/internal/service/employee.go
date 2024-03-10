package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

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

	// TODO: create read transaction
	employee, err := s.store.GetEmployeeByID(ctx, id)
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
