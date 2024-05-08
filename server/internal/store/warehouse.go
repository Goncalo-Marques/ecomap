package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// CreateWarehouse executes a query to create a warehouse with the specified data.
func (s *store) CreateWarehouse(ctx context.Context, tx pgx.Tx, editableWarehouse domain.EditableWarehouse, roadID, municipalityID *int) (uuid.UUID, error) {
	geoJSON, err := jsonMarshalGeoJSONGeometryPoint(editableWarehouse.GeoJSON)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO warehouses (truck_capacity, geom, road_id, municipality_id)
		VALUES ($1, ST_GeomFromGeoJSON($2), $3, $4) 
		RETURNING id
	`,
		editableWarehouse.TruckCapacity,
		geoJSON,
		roadID,
		municipalityID,
	)

	var id uuid.UUID

	err = row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}

// ListWarehouses executes a query to return the warehouses for the specified filter.
func (s *store) ListWarehouses(ctx context.Context, tx pgx.Tx, filter domain.WarehousesPaginatedFilter) (domain.PaginatedResponse[domain.Warehouse], error) {
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	if filter.LocationName != nil {
		filterLocationFields = []string{"rn.osm_name", "m.name"}
		argsWhere = append(argsWhere, *filter.LocationName)
	}

	sqlWhere := listSQLWhere(nil, filterLocationFields)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(w.id) 
		FROM warehouses AS w
		LEFT JOIN road_network AS rn ON w.road_id = rn.id
		LEFT JOIN municipalities AS m ON w.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Warehouse]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.WarehousePaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "w.created_at"
	switch domainSortField {
	case domain.WarehousePaginatedSortTruckCapacity:
		sortField = "w.truck_capacity"
	case domain.WarehousePaginatedSortWayName:
		sortField = "rn.osm_name"
	case domain.WarehousePaginatedSortMunicipalityName:
		sortField = "m.name"
	case domain.WarehousePaginatedSortCreatedAt:
		sortField = "w.created_at"
	case domain.WarehousePaginatedSortModifiedAt:
		sortField = "w.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT w.id, w.truck_capacity, ST_AsGeoJSON(w.geom)::jsonb, rn.osm_name, m.name, w.created_at, w.modified_at
		FROM warehouses AS w
		LEFT JOIN road_network AS rn ON w.road_id = rn.id
		LEFT JOIN municipalities AS m ON w.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Warehouse]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	warehouses, err := getWarehousesFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Warehouse]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Warehouse]{
		Total:   total,
		Results: warehouses,
	}, nil
}

// GetWarehouseByID executes a query to return the warehouse with the specified identifier.
func (s *store) GetWarehouseByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Warehouse, error) {
	row := tx.QueryRow(ctx, `
		SELECT w.id, w.truck_capacity, ST_AsGeoJSON(w.geom)::jsonb, rn.osm_name, m.name, w.created_at, w.modified_at 
		FROM warehouses AS w
		LEFT JOIN road_network AS rn ON w.road_id = rn.id
		LEFT JOIN municipalities AS m ON w.municipality_id = m.id
		WHERE w.id = $1 
	`,
		id,
	)

	warehouse, err := getWarehouseFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Warehouse{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrWarehouseNotFound)
		}

		return domain.Warehouse{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return warehouse, nil
}

// PatchWarehouse executes a query to patch an warehouse with the specified identifier and data.
func (s *store) PatchWarehouse(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableWarehouse domain.EditableWarehousePatch, roadID, municipalityID *int) error {
	var geoJSON []byte
	var err error

	if editableWarehouse.GeoJSON != nil {
		geoJSON, err = jsonMarshalGeoJSONGeometryPoint(editableWarehouse.GeoJSON)
		if err != nil {
			return fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE warehouses SET
			truck_capacity = coalesce($2, truck_capacity),
			geom = coalesce(ST_GeomFromGeoJSON($3), geom),
			road_id = CASE 
					WHEN $3 IS NOT NULL THEN $4 
					ELSE road_id
				END,
			municipality_id = CASE 
					WHEN $3 IS NOT NULL THEN $5
					ELSE municipality_id
				END
		WHERE id = $1
	`,
		id,
		editableWarehouse.TruckCapacity,
		geoJSON,
		roadID,
		municipalityID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseNotFound)
	}

	return nil
}

// DeleteWarehouseByID executes a query to delete the warehouse with the specified identifier.
func (s *store) DeleteWarehouseByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM warehouses
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintWarehousesTrucksWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseAssociatedWithWarehouseTruck)
		case constraintRoutesDepartureWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseAssociatedWithRouteDeparture)
		case constraintRoutesArrivalWarehouseIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseAssociatedWithRouteArrival)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrWarehouseNotFound)
	}

	return nil
}

// getWarehouseFromRow returns the warehouse by scanning the given row.
func getWarehouseFromRow(row pgx.Row) (domain.Warehouse, error) {
	var warehouse domain.Warehouse
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&warehouse.ID,
		&warehouse.TruckCapacity,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&warehouse.CreatedAt,
		&warehouse.ModifiedAt,
	)
	if err != nil {
		return domain.Warehouse{}, err
	}

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	warehouse.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return warehouse, nil
}

// getWarehousesFromRows returns the warehouses by scanning the given rows.
func getWarehousesFromRows(rows pgx.Rows) ([]domain.Warehouse, error) {
	var Warehouses []domain.Warehouse
	for rows.Next() {
		warehouse, err := getWarehouseFromRow(rows)
		if err != nil {
			return nil, err
		}

		Warehouses = append(Warehouses, warehouse)
	}

	return Warehouses, nil
}
