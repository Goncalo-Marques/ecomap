package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintUsersContainerBookmarksPkey            = "users_container_bookmarks_pkey"
	constraintUsersContainerBookmarksUserIDFkey      = "users_container_bookmarks_user_id_fkey"
	constraintUsersContainerBookmarksContainerIDFkey = "users_container_bookmarks_container_id_fkey"
)

// CreateUserContainerBookmark executes a query to create a user container bookmark with the specified identifiers.
func (s *store) CreateUserContainerBookmark(ctx context.Context, tx pgx.Tx, userID, containerID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO users_container_bookmarks (user_id, container_id)
		VALUES ($1, $2)
	`,
		userID,
		containerID,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintUsersContainerBookmarksPkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserContainerBookmarkAlreadyExists)
		case constraintUsersContainerBookmarksUserIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserNotFound)
		case constraintUsersContainerBookmarksContainerIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerNotFound)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	return nil
}

// ListUserContainerBookmarks executes a query to return the user container bookmarks for the specified filter.
func (s *store) ListUserContainerBookmarks(ctx context.Context, tx pgx.Tx, userID uuid.UUID, filter domain.UserContainerBookmarksPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	var filterFields []string
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	filterFields = append(filterFields, "u.id::text")
	argsWhere = append(argsWhere, userID)
	if filter.ContainerCategory != nil {
		filterFields = append(filterFields, "c.category::text")
		argsWhere = append(argsWhere, containerCategoryFromDomain(*filter.ContainerCategory))
	}
	if filter.LocationName != nil {
		filterLocationFields = []string{"rn.osm_name", "m.name"}
		argsWhere = append(argsWhere, *filter.LocationName)
	}

	sqlWhere := listSQLWhere(filterFields, filterLocationFields)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(ucb.user_id)
		FROM users_container_bookmarks AS ucb
		INNER JOIN users AS u ON ucb.user_id = u.id
		INNER JOIN containers AS c ON ucb.container_id = c.id
		LEFT JOIN road_network AS rn ON c.road_id = rn.id
		LEFT JOIN municipalities AS m ON c.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.UserContainerBookmarkPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "ucb.created_at"
	switch domainSortField {
	case domain.UserContainerBookmarkPaginatedSortContainerCategory:
		sortField = "c.category"
	case domain.UserContainerBookmarkPaginatedSortContainerWayName:
		sortField = "rn.osm_name"
	case domain.UserContainerBookmarkPaginatedSortContainerMunicipalityName:
		sortField = "m.name"
	case domain.UserContainerBookmarkPaginatedSortCreatedAt:
		sortField = "ucb.created_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT c.id, c.category, ST_AsGeoJSON(c.geom)::jsonb, rn.osm_name, m.name, c.created_at, c.modified_at
		FROM users_container_bookmarks AS ucb
		INNER JOIN users AS u ON ucb.user_id = u.id
		INNER JOIN containers AS c ON ucb.container_id = c.id
		LEFT JOIN road_network AS rn ON c.road_id = rn.id
		LEFT JOIN municipalities AS m ON c.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	containers, err := getContainersFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Container]{
		Total:   total,
		Results: containers,
	}, nil
}

// DeleteUserContainerBookmark executes a query to delete the user container bookmark with the specified identifiers.
func (s *store) DeleteUserContainerBookmark(ctx context.Context, tx pgx.Tx, userID, containerID uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM users_container_bookmarks
		WHERE user_id = $1 AND container_id = $2
	`,
		userID,
		containerID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrUserContainerBookmarkNotFound)
	}

	return nil
}
