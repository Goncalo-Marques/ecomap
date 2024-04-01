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
	descriptionFailedGetUserSignIn     = "service: failed to get user sign-in"
	descriptionFailedGetUserByUsername = "service: failed to get user by username"
)

// SignInUser returns a JSON Web Token for the specified username and password.
func (s *service) SignInUser(ctx context.Context, username string, password string) (string, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "SignInUser"),
		slog.String(logging.UserUsername, username),
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

	valid, err := s.authnService.CheckPasswordHash(password, signIn.Password)
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
