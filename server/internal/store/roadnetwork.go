package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetRoadIDByGeometry executes a query to return the identifier of the road that is closest to the given geometry in
// the road network.
func (s *store) GetRoadIDByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (int, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT rn.id
		FROM pgr_findCloseEdges(
			$$SELECT id, geom_way as geom FROM road_network$$,
			ST_GeomFromGeoJSON($1),
			0.5
		) AS ce
		INNER JOIN road_network AS rn ON ce.edge_id = rn.id
	`,
		string(geoJSON),
	)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRoadNotFound)
		}

		return 0, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return id, nil
}
