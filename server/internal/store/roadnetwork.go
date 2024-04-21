package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetWayOSMByGeoJSON executes a query to return the identifier of the OpenStreetMap way that is closest to the given
// geometry in the road network.
func (s *store) GetWayOSMByGeoJSON(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (*int, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT b.osm_id
		FROM pgr_findCloseEdges(
			$$SELECT id, geom_way as geom FROM road_network$$,
			ST_GeomFromGeoJSON($1),
			0.5
		) AS ce
		INNER JOIN road_network AS rn 
			ON ce.edge_id = rn.id
	`,
		string(geoJSON),
	)

	var wayOSM *int
	err = row.Scan(&wayOSM)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return wayOSM, nil
}
