package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedCreateUser        = "service: failed to create user"
	descriptionFailedGetUserByUsername = "service: failed to get user by username"
	descriptionFailedGetUserSignIn     = "service: failed to get user sign-in"
)

// CreateUser creates a new user with the specified data.
func (s *service) CreateUser(ctx context.Context, editableUser domain.EditableUserWithPassword) (domain.User, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateUser"),
		slog.String(logging.UserUsername, string(editableUser.Username)),
		slog.String(logging.UserFirstName, string(editableUser.FirstName)),
		slog.String(logging.UserLastName, string(editableUser.LastName)),
	}

	editableUser.Username = domain.Username(replaceSpacesWithHyphen(string(editableUser.Username)))
	editableUser.FirstName = domain.Name(replaceSpacesWithHyphen(string(editableUser.FirstName)))
	editableUser.LastName = domain.Name(replaceSpacesWithHyphen(string(editableUser.LastName)))

	if !editableUser.Username.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldInvalid{FieldName: fieldUsername}, descriptionInvalidField, logAttrs...)
	}
	if !s.authnService.ValidPassword([]byte(editableUser.Password)) {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldInvalid{FieldName: fieldPassword}, descriptionInvalidField, logAttrs...)
	}
	if !editableUser.FirstName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldInvalid{FieldName: fieldFirstName}, descriptionInvalidField, logAttrs...)
	}
	if !editableUser.LastName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldInvalid{FieldName: fieldLastName}, descriptionInvalidField, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(editableUser.Password))
	if err != nil {
		return domain.User{}, logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	editableUser.Password = domain.Password(hashedPassword)

	var user domain.User

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		user, err = s.store.CreateUser(ctx, tx, editableUser)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			return domain.User{}, logInfoAndWrapError(ctx, err, descriptionFailedCreateUser, logAttrs...)
		default:
			return domain.User{}, logAndWrapError(ctx, err, descriptionFailedCreateUser, logAttrs...)
		}
	}

	return user, nil
}

// SignInUser returns a JSON Web Token for the specified username and password.
func (s *service) SignInUser(ctx context.Context, username domain.Username, password string) (string, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "SignInUser"),
		slog.String(logging.UserUsername, string(username)),
	}

	var signIn domain.SignIn
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		signIn, err = s.store.GetUserSignIn(ctx, tx, username)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedGetUserSignIn, logAttrs...)
		default:
			return "", logAndWrapError(ctx, err, descriptionFailedGetUserSignIn, logAttrs...)
		}
	}

	valid, err := s.authnService.CheckPasswordHash([]byte(password), []byte(signIn.Password))
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	if !valid {
		return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	var user domain.User

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		user, err = s.store.GetUserByUsername(ctx, tx, username)
		return err
	})
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedGetUserByUsername, logAttrs...)
	}

	role := authn.SubjectRoleUser
	token, err := s.authnService.NewJWT(user.ID.String(), []authn.SubjectRole{role})
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCreateJWT, logAttrs...)
	}

	return token, nil
}
