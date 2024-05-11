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
	constraintRoutesTruckIDFkey              = "routes_truck_id_fkey"
	constraintRoutesDepartureWarehouseIDFkey = "routes_departure_warehouse_id_fkey"
	constraintRoutesArrivalWarehouseIDFkey   = "routes_arrival_warehouse_id_fkey"
)

// CreateRoute executes a query to create a route with the specified data.
func (s *store) CreateRoute(ctx context.Context, tx pgx.Tx, editableRoute domain.EditableRoute) (uuid.UUID, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO routes (name, truck_id, departure_warehouse_id, arrival_warehouse_id)
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`,
		editableRoute.Name,
		editableRoute.TruckID,
		editableRoute.DepartureWarehouseID,
		editableRoute.ArrivalWarehouseID,
	)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintRoutesTruckIDFkey:
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrTruckNotFound)
		case constraintRoutesDepartureWarehouseIDFkey:
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRouteDepartureWarehouseNotFound)
		case constraintRoutesArrivalWarehouseIDFkey:
			return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRouteArrivalWarehouseNotFound)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}

// ListRoutes executes a query to return the routes for the specified filter.
func (s *store) ListRoutes(ctx context.Context, tx pgx.Tx, filter domain.RoutesPaginatedFilter) (domain.PaginatedResponse[domain.Route], error) {
	var filterFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	if filter.Name != nil {
		filterFields = append(filterFields, "r.name")
		argsWhere = append(argsWhere, *filter.Name)
	}
	if filter.TruckID != nil {
		filterFields = append(filterFields, "r.truck_id")
		argsWhere = append(argsWhere, *filter.TruckID)
	}
	if filter.DepartureWarehouseID != nil {
		filterFields = append(filterFields, "r.departure_warehouse_id")
		argsWhere = append(argsWhere, *filter.DepartureWarehouseID)
	}
	if filter.ArrivalWarehouseID != nil {
		filterFields = append(filterFields, "r.arrival_warehouse_id")
		argsWhere = append(argsWhere, *filter.ArrivalWarehouseID)
	}

	sqlWhere := listSQLWhere(filterFields, nil)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(r.id) 
		FROM routes AS r
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Route]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.RoutePaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "r.created_at"
	switch domainSortField {
	case domain.RoutePaginatedSortName:
		sortField = "r.name"
	case domain.RoutePaginatedSortTruckID:
		sortField = "r.truck_id"
	case domain.RoutePaginatedSortDepartureWarehouseID:
		sortField = "r.departure_warehouse_id"
	case domain.RoutePaginatedSortArrivalWarehouseID:
		sortField = "r.arrival_warehouse_id"
	case domain.RoutePaginatedSortCreatedAt:
		sortField = "r.created_at"
	case domain.RoutePaginatedSortModifiedAt:
		sortField = "r.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT r.id, r.name, r.truck_id, r.departure_warehouse_id, r.arrival_warehouse_id, r.created_at, r.modified_at
		FROM routes AS r
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Route]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	routes, err := getRoutesFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Route]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Route]{
		Total:   total,
		Results: routes,
	}, nil
}

// GetRouteByID executes a query to return the route with the specified identifier.
func (s *store) GetRouteByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Route, error) {
	row := tx.QueryRow(ctx, `
		SELECT r.id, r.name, r.truck_id, r.departure_warehouse_id, r.arrival_warehouse_id, r.created_at, r.modified_at
		FROM routes AS r
		WHERE r.id = $1 
	`,
		id,
	)

	route, err := getRouteFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Route{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRouteNotFound)
		}

		return domain.Route{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return route, nil
}

// PatchRoute executes a query to patch an route with the specified identifier and data.
func (s *store) PatchRoute(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableRoute domain.EditableRoutePatch) error {
	commandTag, err := tx.Exec(ctx, `
		UPDATE routes SET
			name = coalesce($2, name),
			truck_id = coalesce($3, truck_id),
			departure_warehouse_id = coalesce($4, departure_warehouse_id),
			arrival_warehouse_id = coalesce($5, arrival_warehouse_id)
		WHERE id = $1
	`,
		id,
		editableRoute.Name,
		editableRoute.TruckID,
		editableRoute.DepartureWarehouseID,
		editableRoute.ArrivalWarehouseID,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintRoutesTruckIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrTruckNotFound)
		case constraintRoutesDepartureWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRouteDepartureWarehouseNotFound)
		case constraintRoutesArrivalWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRouteArrivalWarehouseNotFound)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteNotFound)
	}

	return nil
}

// DeleteRouteByID executes a query to delete the route with the specified identifier.
func (s *store) DeleteRouteByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM routes
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintRoutesContainersRouteIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteAssociatedWithRouteContainer)
		case constraintRoutesEmployeesRouteIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteAssociatedWithRouteEmployee)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrRouteNotFound)
	}

	return nil
}

// getRouteFromRow returns the route by scanning the given row.
func getRouteFromRow(row pgx.Row) (domain.Route, error) {
	var route domain.Route

	err := row.Scan(
		&route.ID,
		&route.Name,
		&route.TruckID,
		&route.DepartureWarehouseID,
		&route.ArrivalWarehouseID,
		&route.CreatedAt,
		&route.ModifiedAt,
	)
	if err != nil {
		return domain.Route{}, err
	}

	return route, nil
}

// getRoutesFromRows returns the routes by scanning the given rows.
func getRoutesFromRows(rows pgx.Rows) ([]domain.Route, error) {
	var routes []domain.Route
	for rows.Next() {
		route, err := getRouteFromRow(rows)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}
