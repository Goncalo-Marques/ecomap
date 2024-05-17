package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// CreateLandfill executes a query to create a landfill with the specified data.
func (s *store) CreateLandfill(ctx context.Context, tx pgx.Tx, editableLandfill domain.EditableLandfill, roadID, municipalityID *int) (uuid.UUID, error) {
	geoJSON, err := jsonMarshalGeoJSONGeometryPoint(editableLandfill.GeoJSON)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO landfills (geom, road_id, municipality_id)
		VALUES (ST_GeomFromGeoJSON($1), $2, $3) 
		RETURNING id
	`,
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

// ListLandfills executes a query to return the landfills for the specified filter.
func (s *store) ListLandfills(ctx context.Context, tx pgx.Tx, filter domain.LandfillsPaginatedFilter) (domain.PaginatedResponse[domain.Landfill], error) {
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
		SELECT count(l.id) 
		FROM landfills AS l
		LEFT JOIN road_network AS rn ON l.road_id = rn.id
		LEFT JOIN municipalities AS m ON l.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Landfill]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.LandfillPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "l.created_at"
	switch domainSortField {
	case domain.LandfillPaginatedSortWayName:
		sortField = "rn.osm_name"
	case domain.LandfillPaginatedSortMunicipalityName:
		sortField = "m.name"
	case domain.LandfillPaginatedSortCreatedAt:
		sortField = "l.created_at"
	case domain.LandfillPaginatedSortModifiedAt:
		sortField = "l.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT l.id, ST_AsGeoJSON(l.geom)::jsonb, rn.osm_name, m.name, l.created_at, l.modified_at
		FROM landfills AS l
		LEFT JOIN road_network AS rn ON l.road_id = rn.id
		LEFT JOIN municipalities AS m ON l.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Landfill]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	landfills, err := getLandfillsFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Landfill]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Landfill]{
		Total:   total,
		Results: landfills,
	}, nil
}

// GetLandfillByID executes a query to return the landfill with the specified identifier.
func (s *store) GetLandfillByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Landfill, error) {
	row := tx.QueryRow(ctx, `
		SELECT l.id, ST_AsGeoJSON(l.geom)::jsonb, rn.osm_name, m.name, l.created_at, l.modified_at 
		FROM landfills AS l
		LEFT JOIN road_network AS rn ON l.road_id = rn.id
		LEFT JOIN municipalities AS m ON l.municipality_id = m.id
		WHERE l.id = $1 
	`,
		id,
	)

	landfill, err := getLandfillFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Landfill{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrLandfillNotFound)
		}

		return domain.Landfill{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return landfill, nil
}

// GetLandfillClosestGeometry executes a query to return the landfill that is closest to the given geometry.
func (s *store) GetLandfillClosestGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Landfill, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return domain.Landfill{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT l.id, ST_AsGeoJSON(l.geom)::jsonb, rn.osm_name, m.name, l.created_at, l.modified_at 
		FROM landfills AS l
		LEFT JOIN road_network AS rn ON l.road_id = rn.id
		LEFT JOIN municipalities AS m ON l.municipality_id = m.id
		ORDER BY ST_Distance(l.geom, ST_GeomFromGeoJSON($1))
		LIMIT 1
	`,
		string(geoJSON),
	)

	landfill, err := getLandfillFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Landfill{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrLandfillNotFound)
		}

		return domain.Landfill{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return landfill, nil
}

// PatchLandfill executes a query to patch an landfill with the specified identifier and data.
func (s *store) PatchLandfill(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableLandfill domain.EditableLandfillPatch, roadID, municipalityID *int) error {
	var geoJSON []byte
	var err error

	if editableLandfill.GeoJSON != nil {
		geoJSON, err = jsonMarshalGeoJSONGeometryPoint(editableLandfill.GeoJSON)
		if err != nil {
			return fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE landfills SET
			geom = coalesce(ST_GeomFromGeoJSON($2), geom),
			road_id = CASE 
					WHEN $2 IS NOT NULL THEN $3
					ELSE road_id
				END,
			municipality_id = CASE 
					WHEN $2 IS NOT NULL THEN $4
					ELSE municipality_id
				END
		WHERE id = $1
	`,
		id,
		geoJSON,
		roadID,
		municipalityID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrLandfillNotFound)
	}

	return nil
}

// DeleteLandfillByID executes a query to delete the landfill with the specified identifier.
func (s *store) DeleteLandfillByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM landfills
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrLandfillNotFound)
	}

	return nil
}

// getLandfillFromRow returns the landfill by scanning the given row.
func getLandfillFromRow(row pgx.Row) (domain.Landfill, error) {
	var landfill domain.Landfill
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&landfill.ID,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&landfill.CreatedAt,
		&landfill.ModifiedAt,
	)
	if err != nil {
		return domain.Landfill{}, err
	}

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	landfill.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return landfill, nil
}

// getLandfillsFromRows returns the landfills by scanning the given rows.
func getLandfillsFromRows(rows pgx.Rows) ([]domain.Landfill, error) {
	var Landfills []domain.Landfill
	for rows.Next() {
		landfill, err := getLandfillFromRow(rows)
		if err != nil {
			return nil, err
		}

		Landfills = append(Landfills, landfill)
	}

	return Landfills, nil
}
