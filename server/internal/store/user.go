package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetUserSignIn executes a query to return the sign-in of the user with the specified username.
func (s *store) GetUserSignIn(ctx context.Context, tx pgx.Tx, username string) (domain.SignIn, error) {
	row := tx.QueryRow(ctx, `
		SELECT username, password 
		FROM users 
		WHERE username = $1 
	`, username)

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
func (s *store) GetUserByUsername(ctx context.Context, tx pgx.Tx, username string) (domain.User, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, first_name, last_name, created_time, modified_time 
		FROM users 
		WHERE username = $1 
	`, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	users, err := getUsersFromRows(rows)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	if len(users) == 0 {
		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, domain.ErrUserNotFound)
	}

	return users[0], nil
}

// getUsersFromRows returns the users by scanning the given rows.
func getUsersFromRows(rows pgx.Rows) ([]domain.User, error) {
	var users []domain.User
	var err error

	for rows.Next() {
		var user domain.User

		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedTime,
			&user.ModifiedTime,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
