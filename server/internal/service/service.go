package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

// AuthenticationService defines the authentication service interface.
type AuthenticationService interface {
	ValidPassword(password string) bool
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) (bool, error)

	NewJWT(subject string) (string, error)
	ParseJWT(tokenString string) (string, error)
}

// Store defines the store interface.
type Store interface {
	GetEmployeeSignIn(ctx context.Context, tx pgx.Tx, username string) (domain.SignIn, error)
	GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error)

	NewTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode) (pgx.Tx, error)
}

// service defines the service structure.
type service struct {
	authnService AuthenticationService
	store        Store
}

// New returns a new http handler.
func New(authnService AuthenticationService, store Store) *service {
	return &service{
		authnService: authnService,
		store:        store,
	}
}

// rollbackFunc returns a function to rollback a transaction.
func rollbackFunc(ctx context.Context, tx pgx.Tx) func() {
	return func() {
		err := tx.Rollback(ctx)
		if err != nil {
			logging.Logger.ErrorContext(ctx, "service: failed to rollback transaction", logging.Error(err))
		}
	}
}

// readOnlyTx returns a read only transaction wrapper.
func (s *service) readOnlyTx(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.store.NewTx(ctx, pgx.ReadCommitted, pgx.ReadOnly)
	if err != nil {
		return err
	}
	defer rollbackFunc(ctx, tx)()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// TODO: Avoid lint issue (remove this comments in the future)
// readWriteTx returns a read and write transaction wrapper.
// func (s *service) readWriteTx(ctx context.Context, f func(pgx.Tx) error) error {
// 	tx, err := s.store.NewTx(ctx, pgx.RepeatableRead, pgx.ReadWrite)
// 	if err != nil {
// 		return err
// 	}
// 	defer rollbackFunc(ctx, tx)()

// 	if err := f(tx); err != nil {
// 		return err
// 	}

// 	return tx.Commit(ctx)
// }

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
