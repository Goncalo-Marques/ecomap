package store

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
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

// ListUsers executes a query to return the users for the specified filter.
func (s *store) ListUsers(ctx context.Context, tx pgx.Tx, filter domain.UsersFilter) (domain.PaginatedResponse[domain.User], error) {
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

	var sqlWhere string
	if len(filterFields) > 0 {
		for i, field := range filterFields {
			filterFields[i] = field + " ILIKE '%%' || $%d || '%%'"
		}

		sqlWhere = " WHERE " + strings.Join(filterFields, " AND ")
	}

	// Format the where sql parameters.
	if len(argsWhere) > 0 {
		sqlParamIndices := make([]any, len(argsWhere))
		for i := range argsWhere {
			sqlParamIndices[i] = i + 1
		}

		sqlWhere = fmt.Sprintf(sqlWhere, sqlParamIndices...)
	}

	// Get users count.
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

	var sqlSort string
	argsSort := make([]any, 0, 2)

	// Append the field to sort, if provided.
	var sortField domain.UserSort
	if filter.Sort != nil {
		sortField = filter.Sort.Field()
	}

	sqlSort = " ORDER BY "
	switch sortField {
	case domain.UserSortUsername:
		sqlSort += "username"
	case domain.UserSortFirstName:
		sqlSort += "first_name"
	case domain.UserSortLastName:
		sqlSort += "last_name"
	case domain.UserSortCreatedTime:
		sqlSort += "created_time"
	case domain.UserSortModifiedTime:
		sqlSort += "modified_time"
	default:
		sqlSort += "created_time"
	}

	order := " ASC"
	if filter.Order == domain.OrderDesc {
		order = " DESC"
	}
	sqlSort += order

	// Append the limit and offset.
	sqlSort += " LIMIT $%d OFFSET $%d"
	argsSort = append(argsSort, filter.Limit, filter.Offset)

	// Format the sort sql parameters.
	if len(argsSort) > 0 {
		sqlParamIndices := make([]any, len(argsSort))
		for i := range argsSort {
			sqlParamIndices[i] = len(argsWhere) + i + 1
		}

		sqlSort = fmt.Sprintf(sqlSort, sqlParamIndices...)
	}

	// Get users.
	rows, err := tx.Query(ctx, `
		SELECT id, username, first_name, last_name, created_time, modified_time 
		FROM users 
	`+sqlWhere+sqlSort,
		slices.Concat(argsWhere, argsSort)...,
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
		SELECT id, username, first_name, last_name, created_time, modified_time 
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
func (s *store) PatchUser(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableUser domain.EditableUserPatch) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		UPDATE users SET
			username = coalesce($2, username),
			first_name = coalesce($3, first_name),
			last_name = coalesce($4, last_name)
		WHERE id = $1
		RETURNING id, username, first_name, last_name, created_time, modified_time
	`,
		id,
		editableUser.Username,
		editableUser.FirstName,
		editableUser.LastName,
	)

	user, err := getUserFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserNotFound)
		}
		if getConstraintName(err) == constraintUsersUsernameKey {
			return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrUserAlreadyExists)
		}

		return domain.User{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return user, nil
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
func (s *store) DeleteUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error) {
	row := tx.QueryRow(ctx, `
		DELETE FROM users
		WHERE id = $1
		RETURNING id, username, first_name, last_name, created_time, modified_time
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
