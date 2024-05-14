package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintRoutesContainersPkey            = "routes_containers_pkey"
	constraintRoutesContainersRouteIDFkey     = "routes_containers_route_id_fkey"
	constraintRoutesContainersContainerIDFkey = "routes_containers_container_id_fkey"
	constraintRoutesContainersRouteIDFkey     = "routes_containers_route_id_fkey"
)

// CreateRouteContainer executes a query to create a route container association with the specified identifiers.
func (s *store) CreateRouteContainer(ctx context.Context, tx pgx.Tx, routeID, containerID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO routes_containers (route_id, container_id)
		VALUES ($1, $2)
	`,
		routeID,
		containerID,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintRoutesContainersPkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteContainerAlreadyExists)
		case constraintRoutesContainersRouteIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteNotFound)
		case constraintRoutesContainersContainerIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerNotFound)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	return nil
}

// ListRouteContainers executes a query to return the route containers for the specified filter.
func (s *store) ListRouteContainers(ctx context.Context, tx pgx.Tx, routeID uuid.UUID, filter domain.RouteContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	var filterFields []string
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	filterFields = append(filterFields, "rc.route_id::text")
	argsWhere = append(argsWhere, routeID)
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
		SELECT count(rc.route_id)
		FROM routes_containers AS rc
		INNER JOIN containers AS c ON rc.container_id = c.id
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
	var domainSortField domain.RouteContainerPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "rc.created_at"
	switch domainSortField {
	case domain.RouteContainerPaginatedSortContainerCategory:
		sortField = "c.category"
	case domain.RouteContainerPaginatedSortContainerWayName:
		sortField = "rn.osm_name"
	case domain.RouteContainerPaginatedSortContainerMunicipalityName:
		sortField = "m.name"
	case domain.RouteContainerPaginatedSortCreatedAt:
		sortField = "rc.created_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT c.id, c.category, ST_AsGeoJSON(c.geom)::jsonb, rn.osm_name, m.name, c.created_at, c.modified_at
		FROM routes_containers AS rc
		INNER JOIN containers AS c ON rc.container_id = c.id
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

// DeleteRouteContainer executes a query to delete the route container association with the specified identifiers.
func (s *store) DeleteRouteContainer(ctx context.Context, tx pgx.Tx, routeID, containerID uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM routes_containers
		WHERE route_id = $1 AND container_id = $2
	`,
		routeID,
		containerID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteContainerNotFound)
	}

	return nil
}
