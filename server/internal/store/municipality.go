package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetMunicipalityByGeometry executes a query to return the municipality that contains the given geometry.
func (s *store) GetMunicipalityByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Municipality, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return domain.Municipality{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT id, fid, name, district, nutsiii, nutsii, nutsi, area_ha, perimeter_km
		FROM municipalities
		WHERE ST_Within(ST_GeomFromGeoJSON($1), geom)
	`,
		string(geoJSON),
	)

	var municipality domain.Municipality

	err = row.Scan(
		&municipality.ID,
		&municipality.FeatureID,
		&municipality.Name,
		&municipality.District,
		&municipality.NUTS3,
		&municipality.NUTS2,
		&municipality.NUTS1,
		&municipality.Area,
		&municipality.Perimeter,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Municipality{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrMunicipalityNotFound)
		}

		return domain.Municipality{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return municipality, nil
}
