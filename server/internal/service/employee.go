package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedGetEmployeeSignIn     = "service: failed to get employee sign-in"
	descriptionFailedGetEmployeeByUsername = "service: failed to get employee by username"
	descriptionFailedGetEmployeeByID       = "service: failed to get employee by id"
)

// SignInEmployee returns a JSON Web Token for the specified username and password.
func (s *service) SignInEmployee(ctx context.Context, username string, password string) (string, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "SignInEmployee"),
		slog.String(logging.EmployeeUsername, username),
	}

	var signIn domain.SignIn
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		signIn, err = s.store.GetEmployeeSignIn(ctx, tx, username)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedGetEmployeeSignIn, logAttrs...)
		default:
			return "", logAndWrapError(ctx, err, descriptionFailedGetEmployeeSignIn, logAttrs...)
		}
	}

	valid, err := s.authnService.CheckPasswordHash(password, signIn.Password)
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	if !valid {
		return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	var employee domain.Employee

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		employee, err = s.store.GetEmployeeByUsername(ctx, tx, username)
		return err
	})
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedGetEmployeeByUsername, logAttrs...)
	}

	var role authn.SubjectRole
	switch employee.Role {
	case domain.EmployeeRoleWasteOperator:
		role = authn.SubjectRoleWasteOperator
	case domain.EmployeeRoleManager:
		role = authn.SubjectRoleManager
	}

	token, err := s.authnService.NewJWT(employee.ID.String(), role)
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCreateJWT, logAttrs...)
	}

	return token, nil
}

// GetEmployeeByID returns the employee with the specified identifier.
func (s *service) GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetEmployeeByID"),
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
			return domain.Employee{}, logInfoAndWrapError(ctx, err, descriptionFailedGetEmployeeByID, logAttrs...)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedGetEmployeeByID, logAttrs...)
		}
	}

	return employee, nil
}
