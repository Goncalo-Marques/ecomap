package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// CreateTruck executes a query to create a truck with the specified data.
func (s *store) CreateTruck(ctx context.Context, tx pgx.Tx, editableTruck domain.EditableTruck, roadID, municipalityID *int) (uuid.UUID, error) {
	geoJSON, err := jsonMarshalGeoJSONGeometryPoint(editableTruck.GeoJSON)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO trucks (make, model, license_plate, person_capacity, geom, road_id, municipality_id)
		VALUES ($1, $2, $3, $4, ST_GeomFromGeoJSON($5), $6, $7) 
		RETURNING id
	`,
		editableTruck.Make,
		editableTruck.Model,
		editableTruck.LicensePlate,
		editableTruck.PersonCapacity,
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

// ListTrucks executes a query to return the trucks for the specified filter.
func (s *store) ListTrucks(ctx context.Context, tx pgx.Tx, filter domain.TrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error) {
	var filterFields []string
	var filterLocationFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	if filter.Make != nil {
		filterFields = append(filterFields, "t.make")
		argsWhere = append(argsWhere, *filter.Make)
	}
	if filter.Model != nil {
		filterFields = append(filterFields, "t.model")
		argsWhere = append(argsWhere, *filter.Model)
	}
	if filter.LicensePlate != nil {
		filterFields = append(filterFields, "t.license_plate")
		argsWhere = append(argsWhere, *filter.LicensePlate)
	}
	if filter.LocationName != nil {
		filterLocationFields = []string{"rn.osm_name", "m.name"}
		argsWhere = append(argsWhere, *filter.LocationName)
	}

	sqlWhere := listSQLWhere(filterFields, filterLocationFields)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(t.id) 
		FROM trucks AS t
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
	var domainSortField domain.TruckPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "t.created_at"
	switch domainSortField {
	case domain.TruckPaginatedSortMake:
		sortField = "t.make"
	case domain.TruckPaginatedSortModel:
		sortField = "t.model"
	case domain.TruckPaginatedSortLicensePlate:
		sortField = "t.license_plate"
	case domain.TruckPaginatedSortWayName:
		sortField = "rn.osm_name"
	case domain.TruckPaginatedSortMunicipalityName:
		sortField = "m.name"
	case domain.TruckPaginatedSortCreatedAt:
		sortField = "t.created_at"
	case domain.TruckPaginatedSortModifiedAt:
		sortField = "t.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT t.id, t.make, t.model, t.license_plate, t.person_capacity, ST_AsGeoJSON(t.geom)::jsonb, rn.osm_name, m.name, t.created_at, t.modified_at
		FROM trucks AS t
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

// GetTruckByID executes a query to return the truck with the specified identifier.
func (s *store) GetTruckByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Truck, error) {
	row := tx.QueryRow(ctx, `
		SELECT t.id, t.make, t.model, t.license_plate, t.person_capacity, ST_AsGeoJSON(t.geom)::jsonb, rn.osm_name, m.name, t.created_at, t.modified_at
		FROM trucks AS t
		LEFT JOIN road_network AS rn ON t.road_id = rn.id
		LEFT JOIN municipalities AS m ON t.municipality_id = m.id
		WHERE t.id = $1 
	`,
		id,
	)

	truck, err := getTruckFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Truck{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrTruckNotFound)
		}

		return domain.Truck{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return truck, nil
}

// PatchTruck executes a query to patch an truck with the specified identifier and data.
func (s *store) PatchTruck(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableTruck domain.EditableTruckPatch, roadID, municipalityID *int) error {
	var geoJSON []byte
	var err error

	if editableTruck.GeoJSON != nil {
		geoJSON, err = jsonMarshalGeoJSONGeometryPoint(editableTruck.GeoJSON)
		if err != nil {
			return fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE trucks SET
			make = coalesce($2, make),
			model = coalesce($3, model),
			license_plate = coalesce($4, license_plate),
			person_capacity = coalesce($5, person_capacity),
			geom = coalesce(ST_GeomFromGeoJSON($6), geom),
			road_id = CASE 
					WHEN $6 IS NOT NULL THEN $7 
					ELSE road_id
				END,
			municipality_id = CASE 
					WHEN $6 IS NOT NULL THEN $8
					ELSE municipality_id
				END
		WHERE id = $1
	`,
		id,
		editableTruck.Make,
		editableTruck.Model,
		editableTruck.LicensePlate,
		editableTruck.PersonCapacity,
		geoJSON,
		roadID,
		municipalityID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrTruckNotFound)
	}

	return nil
}

// DeleteTruckByID executes a query to delete the truck with the specified identifier.
func (s *store) DeleteTruckByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM trucks
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintWarehousesTrucksTruckIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrTruckAssociatedWithWarehouseTruck)
		case constraintRoutesTruckIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrTruckAssociatedWithRoute)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrTruckNotFound)
	}

	return nil
}

// getTruckFromRow returns the truck by scanning the given row.
func getTruckFromRow(row pgx.Row) (domain.Truck, error) {
	var truck domain.Truck
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&truck.ID,
		&truck.Make,
		&truck.Model,
		&truck.LicensePlate,
		&truck.PersonCapacity,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&truck.CreatedAt,
		&truck.ModifiedAt,
	)
	if err != nil {
		return domain.Truck{}, err
	}

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	truck.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return truck, nil
}

// getTrucksFromRows returns the trucks by scanning the given rows.
func getTrucksFromRows(rows pgx.Rows) ([]domain.Truck, error) {
	var Trucks []domain.Truck
	for rows.Next() {
		truck, err := getTruckFromRow(rows)
		if err != nil {
			return nil, err
		}

		Trucks = append(Trucks, truck)
	}

	return Trucks, nil
}
