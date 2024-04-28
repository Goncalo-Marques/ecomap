package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintUsersUsernameKey = "users_username_key"
)

// CreateUser executes a query to create a user with the specified data.
func (s *store) CreateUser(ctx context.Context, tx pgx.Tx, editableUser domain.EditableUserWithPassword) (uuid.UUID, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO users (username, password, first_name, last_name)
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`,
		editableUser.Username,
		editableUser.Password,
		editableUser.FirstName,
		editableUser.LastName,
	)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		if getConstraintName(err) == constraintUsersUsernameKey {
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserAlreadyExists)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}

// ListUsers executes a query to return the users for the specified filter.
func (s *store) ListUsers(ctx context.Context, tx pgx.Tx, filter domain.UsersPaginatedFilter) (domain.PaginatedResponse[domain.User], error) {
	filterFields := make([]string, 0, 3)
	argsWhere := make([]any, 0, 3)

	// Append the optional fields to filter.
	if filter.Username != nil {
		filterFields = append(filterFields, "username")
		argsWhere = append(argsWhere, *filter.Username)
	}
	if filter.FirstName != nil {
		filterFields = append(filterFields, "first_name")
		argsWhere = append(argsWhere, *filter.FirstName)
	}
	if filter.LastName != nil {
		filterFields = append(filterFields, "last_name")
		argsWhere = append(argsWhere, *filter.LastName)
	}

	sqlWhere := listSQLWhere(filterFields, filter.LogicalOperator)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(id) 
		FROM users 
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.UserPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "created_at"
	switch domainSortField {
	case domain.UserPaginatedSortUsername:
		sortField = "username"
	case domain.UserPaginatedSortFirstName:
		sortField = "first_name"
	case domain.UserPaginatedSortLastName:
		sortField = "last_name"
	case domain.UserPaginatedSortCreatedAt:
		sortField = "created_at"
	case domain.UserPaginatedSortModifiedAt:
		sortField = "modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT id, username, first_name, last_name, created_at, modified_at 
		FROM users 
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	users, err := getUsersFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.User]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.User]{
		Total:   total,
		Results: users,
	}, nil
}

// GetUserByID executes a query to return the user with the specified identifier.
func (s *store) GetUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, username, first_name, last_name, created_at, modified_at 
		FROM users 
		WHERE id = $1 
	`,
		id,
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

// GetUserByUsername executes a query to return the user with the specified username.
func (s *store) GetUserByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, username, first_name, last_name, created_at, modified_at 
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

// PatchUser executes a query to patch a user with the specified identifier and data.
func (s *store) PatchUser(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableUser domain.EditableUserPatch) error {
	commandTag, err := tx.Exec(ctx, `
		UPDATE users SET
			username = coalesce($2, username),
			first_name = coalesce($3, first_name),
			last_name = coalesce($4, last_name)
		WHERE id = $1
	`,
		id,
		editableUser.Username,
		editableUser.FirstName,
		editableUser.LastName,
	)
	if err != nil {
		if getConstraintName(err) == constraintUsersUsernameKey {
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserAlreadyExists)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserNotFound)
	}

	return nil
}

// UpdateUserPassword executes a query to update the password of the user with the specified username.
func (s *store) UpdateUserPassword(ctx context.Context, tx pgx.Tx, username domain.Username, password domain.Password) error {
	commandTag, err := tx.Exec(ctx, `
		UPDATE users SET
			password = $2
		WHERE username = $1 
	`,
		username,
		password,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserNotFound)
	}

	return nil
}

// DeleteUserByID executes a query to delete the user with the specified identifier.
func (s *store) DeleteUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserNotFound)
	}

	return nil
}

// getUserFromRow returns the user by scanning the given row.
func getUserFromRow(row pgx.Row) (domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.ModifiedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// getUsersFromRows returns the users by scanning the given rows.
func getUsersFromRows(rows pgx.Rows) ([]domain.User, error) {
	var users []domain.User
	for rows.Next() {
		user, err := getUserFromRow(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
