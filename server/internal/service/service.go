package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

// Store defines the store interface.
type Store interface {
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error)
}

// handler defines the http handler structure.
type service struct {
	store Store
}

// New returns a new http handler.
func New(store Store) *service {
	return &service{
		store: store,
	}
}

// logInfoAndWrapError logs the error at the info level and returns the error wrapped with the provided description.
func logInfoAndWrapError(ctx context.Context, err error, description string, logAttrs ...any) error {
	logAttrs = append(logAttrs, logging.Error(err))
	logging.Logger.InfoContext(ctx, description, logAttrs...)
	return fmt.Errorf("%s: %w", description, err)
}

// logAndWrapError logs the error and returns the error wrapped with the provided description.
func logAndWrapError(ctx context.Context, err error, description string, logAttrs ...any) error {
	logAttrs = append(logAttrs, logging.Error(err))
	logging.Logger.ErrorContext(ctx, description, logAttrs...)
	return fmt.Errorf("%s: %w", description, err)
}
