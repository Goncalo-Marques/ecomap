package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintUsersUsernameKey = "users_username_key"
)

// CreateUser executes a query to create a user with the specified data.
func (s *store) CreateUser(ctx context.Context, tx pgx.Tx, editableUser domain.EditableUserWithPassword) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO users (username, password, first_name, last_name)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, username, first_name, last_name, created_time, modified_time
	`,
		editableUser.Username,
		editableUser.Password,
		editableUser.FirstName,
		editableUser.LastName,
	)

	user, err := getUserFromRow(row)
	if err != nil {
		if getConstraintName(err) == constraintUsersUsernameKey {
			return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserAlreadyExists)
		}

		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return user, nil
}

// GetUserSignIn executes a query to return the sign-in of the user with the specified username.
func (s *store) GetUserSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error) {
	row := tx.QueryRow(ctx, `
		SELECT username, password 
		FROM users 
		WHERE username = $1 
	`,
		username,
	)

	var signIn domain.SignIn
	err := row.Scan(
		&signIn.Username,
		&signIn.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserNotFound)
		}

		return domain.SignIn{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return signIn, nil
}

// GetUserByUsername executes a query to return the user with the specified username.
func (s *store) GetUserByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, username, first_name, last_name, created_time, modified_time 
		FROM users 
		WHERE username = $1 
	`,
		username,
	)

	user, err := getUserFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return user, nil
}

// getUserFromRow returns the user by scanning the given row.
func getUserFromRow(row pgx.Row) (domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.CreatedTime,
		&user.ModifiedTime,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// TODO: Remove comments when necessary.
// getUsersFromRows returns the users by scanning the given rows.
// func getUsersFromRows(rows pgx.Rows) ([]domain.User, error) {
// 	var users []domain.User
// 	for rows.Next() {
// 		user, err := getUserFromRow(rows)
// 		if err != nil {
// 			return nil, err
// 		}

// 		users = append(users, user)
// 	}

// 	return users, nil
// }
