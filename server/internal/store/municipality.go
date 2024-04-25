package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetMunicipalityIDByGeometry executes a query to return the identifier of the municipality that is closest to the
// given geometry.
func (s *store) GetMunicipalityIDByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (int, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT id
		FROM municipalities
		WHERE ST_Within(ST_GeomFromGeoJSON($1), geom)
	`,
		string(geoJSON),
	)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrMunicipalityNotFound)
		}

		return 0, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}
