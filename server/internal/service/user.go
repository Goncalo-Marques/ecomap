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
	descriptionFailedCreateUser         = "service: failed to create user"
	descriptionFailedListUsers          = "service: failed to list users"
	descriptionFailedGetUserByID        = "service: failed to get user by id"
	descriptionFailedGetUserByUsername  = "service: failed to get user by username"
	descriptionFailedGetUserSignIn      = "service: failed to get user sign-in"
	descriptionFailedPatchUser          = "service: failed to patch user"
	descriptionFailedUpdateUserPassword = "service: failed to update user password"
	descriptionFailedDeleteUserByID     = "service: failed to delete user by id"
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
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldUsername}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !s.authnService.ValidPassword([]byte(editableUser.Password)) {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldPassword}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableUser.FirstName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldFirstName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableUser.LastName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldLastName}, descriptionInvalidFieldValue, logAttrs...)
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

// ListUsers returns the users with the specified filter.
func (s *service) ListUsers(ctx context.Context, filter domain.UsersFilter) (domain.PaginatedResponse[domain.User], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListUsers"),
	}

	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.User]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: filterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.User]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: filterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}
	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.User]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: filterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.User]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: filterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedUsers domain.PaginatedResponse[domain.User]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedUsers, err = s.store.ListUsers(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, logAndWrapError(ctx, err, descriptionFailedListUsers, logAttrs...)
	}

	return paginatedUsers, nil
}

// GetUserByID returns the user with the specified identifier.
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetUserByID"),
		slog.String(logging.UserID, id.String()),
	}

	var user domain.User
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		user, err = s.store.GetUserByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return domain.User{}, logInfoAndWrapError(ctx, err, descriptionFailedGetUserByID, logAttrs...)
		default:
			return domain.User{}, logAndWrapError(ctx, err, descriptionFailedGetUserByID, logAttrs...)
		}
	}

	return user, nil
}

// PatchUser modifies the user with the specified identifier. Only the specified fields in the request body are updated.
func (s *service) PatchUser(ctx context.Context, id uuid.UUID, editableUser domain.EditableUserPatch) (domain.User, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchUser"),
		slog.String(logging.UserID, id.String()),
	}

	if editableUser.Username != nil {
		username := domain.Username(replaceSpacesWithHyphen(string(*editableUser.Username)))
		editableUser.Username = &username
	}
	if editableUser.FirstName != nil {
		firstName := domain.Name(replaceSpacesWithHyphen(string(*editableUser.FirstName)))
		editableUser.FirstName = &firstName
	}
	if editableUser.LastName != nil {
		lastName := domain.Name(replaceSpacesWithHyphen(string(*editableUser.LastName)))
		editableUser.LastName = &lastName
	}

	if editableUser.Username != nil && !editableUser.Username.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldUsername}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableUser.FirstName != nil && !editableUser.FirstName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldFirstName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableUser.LastName != nil && !editableUser.LastName.Valid() {
		return domain.User{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldLastName}, descriptionInvalidFieldValue, logAttrs...)
	}

	var user domain.User
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		user, err = s.store.PatchUser(ctx, tx, id, editableUser)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound),
			errors.Is(err, domain.ErrUserAlreadyExists):
			return domain.User{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchUser, logAttrs...)
		default:
			return domain.User{}, logAndWrapError(ctx, err, descriptionFailedPatchUser, logAttrs...)
		}
	}

	return user, nil
}

// UpdateUserPassword updates the password of the user with the specified username.
func (s *service) UpdateUserPassword(ctx context.Context, username domain.Username, oldPassword, newPassword domain.Password) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "UpdateUserPassword"),
		slog.String(logging.UserUsername, string(username)),
	}

	if !s.authnService.ValidPassword([]byte(newPassword)) {
		return logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldNewPassword}, descriptionInvalidFieldValue, logAttrs...)
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
			return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedGetUserSignIn, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedGetUserSignIn, logAttrs...)
		}
	}

	valid, err := s.authnService.CheckPasswordHash([]byte(oldPassword), []byte(signIn.Password))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	if !valid {
		return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(newPassword))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	newPassword = domain.Password(hashedPassword)

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		err = s.store.UpdateUserPassword(ctx, tx, username, newPassword)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedUpdateUserPassword, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedUpdateUserPassword, logAttrs...)
		}
	}

	return nil
}

// ResetUserPassword resets the password of the user with the specified username.
func (s *service) ResetUserPassword(ctx context.Context, username domain.Username, newPassword domain.Password) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ResetUserPassword"),
		slog.String(logging.UserUsername, string(username)),
	}

	if !s.authnService.ValidPassword([]byte(newPassword)) {
		return logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: fieldNewPassword}, descriptionInvalidFieldValue, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(newPassword))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	newPassword = domain.Password(hashedPassword)

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		err = s.store.UpdateUserPassword(ctx, tx, username, newPassword)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedUpdateUserPassword, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedUpdateUserPassword, logAttrs...)
		}
	}

	return nil
}

// DeleteUserByID deletes the user with the specified identifier.
func (s *service) DeleteUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteUserByID"),
		slog.String(logging.UserID, id.String()),
	}

	var user domain.User
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		user, err = s.store.DeleteUserByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return domain.User{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteUserByID, logAttrs...)
		default:
			return domain.User{}, logAndWrapError(ctx, err, descriptionFailedDeleteUserByID, logAttrs...)
		}
	}

	return user, nil
}

// SignInUser returns a JSON Web Token for the specified username and password.
func (s *service) SignInUser(ctx context.Context, username domain.Username, password domain.Password) (string, error) {
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
