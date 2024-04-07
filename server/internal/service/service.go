package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	fieldUsername  = "username"
	fieldPassword  = "password"
	fieldFirstName = "firstName"
	fieldLastName  = "lastName"

	filterLimit  = "limit"
	filterOffset = "offset"
	filterSort   = "sort"
	filterOrder  = "order"

	descriptionInvalidFieldValue       = "service: invalid field value"
	descriptionInvalidFilterValue      = "service: invalid filter value"
	descriptionFailedHashPassword      = "service: failed to hash password"
	descriptionFailedCheckPasswordHash = "service: failed to check password hash"
	descriptionFailedCreateJWT         = "service: failed to create jwt"
)

// AuthenticationService defines the authentication service interface.
type AuthenticationService interface {
	ValidPassword(password []byte) bool
	HashPassword(password []byte) ([]byte, error)
	CheckPasswordHash(password, hash []byte) (bool, error)

	NewJWT(subject string, subjectRoles []authn.SubjectRole) (string, error)
}

// Store defines the store interface.
type Store interface {
	CreateUser(ctx context.Context, tx pgx.Tx, editableUser domain.EditableUserWithPassword) (domain.User, error)
	ListUsers(ctx context.Context, tx pgx.Tx, filter domain.UsersFilter) (domain.PaginatedResponse[domain.User], error)
	GetUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error)
	GetUserByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.User, error)
	GetUserSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error)
	PatchUser(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableUser domain.EditableUserPatch) (domain.User, error)
	DeleteUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error)

	GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error)
	GetEmployeeByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.Employee, error)
	GetEmployeeSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error)

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

// readWriteTx returns a read and write transaction wrapper.
func (s *service) readWriteTx(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.store.NewTx(ctx, pgx.RepeatableRead, pgx.ReadWrite)
	if err != nil {
		return err
	}
	defer rollbackFunc(ctx, tx)()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
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

// replaceSpacesWithHyphen returns s with no extra spaces and separates it with a hyphen.
func replaceSpacesWithHyphen(s string) string {
	return strings.Join(strings.Fields(s), "-")
}
