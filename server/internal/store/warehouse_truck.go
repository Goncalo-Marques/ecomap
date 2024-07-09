package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintWarehousesTrucksPkey            = "warehouses_trucks_pkey"
	constraintWarehousesTrucksWarehouseIDFkey = "warehouses_trucks_warehouse_id_fkey"
	constraintWarehousesTrucksTruckIDFkey     = "warehouses_trucks_truck_id_fkey"
)

// CreateWarehouseTruck executes a query to create a warehouse truck association with the specified identifiers.
func (s *store) CreateWarehouseTruck(ctx context.Context, tx pgx.Tx, warehouseID, truckID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO warehouses_trucks (warehouse_id, truck_id)
		VALUES ($1, $2)
	`,
		warehouseID,
		truckID,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintWarehousesTrucksPkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseTruckAlreadyExists)
		case constraintWarehousesTrucksWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseNotFound)
		case constraintWarehousesTrucksTruckIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrTruckNotFound)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	return nil
}

// ListWarehouseTrucks executes a query to return the warehouse trucks for the specified filter.
func (s *store) ListWarehouseTrucks(ctx context.Context, tx pgx.Tx, warehouseID uuid.UUID, filter domain.WarehouseTrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error) {
	var filterFields []string
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	filterFields = append(filterFields, "rc.warehouse_id::text")
	argsWhere = append(argsWhere, warehouseID)
	if filter.TruckMake != nil {
		filterFields = append(filterFields, "t.make")
		argsWhere = append(argsWhere, *filter.TruckMake)
	}
	if filter.TruckModel != nil {
		filterFields = append(filterFields, "t.model")
		argsWhere = append(argsWhere, *filter.TruckModel)
	}
	if filter.TruckLicensePlate != nil {
		filterFields = append(filterFields, "t.license_plate")
		argsWhere = append(argsWhere, *filter.TruckLicensePlate)
	}
	if filter.LocationName != nil {
		filterLocationFields = []string{"rn.osm_name", "m.name"}
		argsWhere = append(argsWhere, *filter.LocationName)
	}

	sqlWhere := listSQLWhere(filterFields, filterLocationFields)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(wt.warehouse_id)
		FROM warehouses_trucks AS wt
		INNER JOIN trucks AS t ON wt.truck_id = t.id
		LEFT JOIN road_network AS rn ON t.road_id = rn.id
		LEFT JOIN municipalities AS m ON t.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Truck]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.WarehouseTruckPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "wt.created_at"
	switch domainSortField {
	case domain.WarehouseTruckPaginatedSortTruckMake:
		sortField = "t.make"
	case domain.WarehouseTruckPaginatedSortTruckModel:
		sortField = "t.model"
	case domain.WarehouseTruckPaginatedSortTruckLicensePlate:
		sortField = "t.license_plate"
	case domain.WarehouseTruckPaginatedSortTruckPersonCapacity:
		sortField = "t.person_capacity"
	case domain.WarehouseTruckPaginatedSortTruckWayName:
		sortField = "rn.osm_name"
	case domain.WarehouseTruckPaginatedSortTruckMunicipalityName:
		sortField = "m.name"
	case domain.WarehouseTruckPaginatedSortCreatedAt:
		sortField = "wt.created_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT t.id, t.make, t.model, t.license_plate, t.person_capacity, ST_AsGeoJSON(t.geom)::jsonb, rn.osm_name, m.name, t.created_at, t.modified_at
		FROM warehouses_trucks AS wt
		INNER JOIN trucks AS t ON wt.truck_id = t.id
		LEFT JOIN road_network AS rn ON t.road_id = rn.id
		LEFT JOIN municipalities AS m ON t.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Truck]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	trucks, err := getTrucksFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Truck]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Truck]{
		Total:   total,
		Results: trucks,
	}, nil
}

// ExistsWarehouseTruck executes a query to return whether the warehouse truck association exists for the specified
// identifiers.
func (s *store) ExistsWarehouseTruck(ctx context.Context, tx pgx.Tx, warehouseID, truckID uuid.UUID) (bool, error) {
	row := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM warehouses_trucks
			WHERE warehouse_id = $1 AND truck_id = $2
		)
	`,
		warehouseID,
		truckID,
	)

	var exists bool

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return exists, nil
}

// DeleteWarehouseTruck executes a query to delete the warehouse truck association with the specified identifiers.
func (s *store) DeleteWarehouseTruck(ctx context.Context, tx pgx.Tx, warehouseID, truckID uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM warehouses_trucks
		WHERE warehouse_id = $1 AND truck_id = $2
	`,
		warehouseID,
		truckID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseTruckNotFound)
	}

	return nil
}
