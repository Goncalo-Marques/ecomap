package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetRoadByGeometry executes a query to return the road that is closest to the given geometry in the road network.
func (s *store) GetRoadByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Road, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		SELECT rn.id, rn.osm_id, rn.osm_name, rn.osm_meta, rn.osm_source_id, rn.osm_target_id, rn.clazz, rn.flags, rn.source, rn.target, rn.km, rn.kmh, rn.cost, rn.reverse_cost, rn.x1, rn.y1, rn.x2, rn.y2
		FROM pgr_findCloseEdges(
			$$SELECT id, geom_way as geom FROM road_network$$,
			ST_GeomFromGeoJSON($1),
			0.5
		) AS ce
		INNER JOIN road_network AS rn ON ce.edge_id = rn.id
	`,
		string(geoJSON),
	)

	var road domain.Road

	err = row.Scan(
		&road.ID,
		&road.OsmID,
		&road.OsmName,
		&road.OsmMeta,
		&road.OsmSourceID,
		&road.OsmTargetID,
		&road.Clazz,
		&road.Flags,
		&road.Source,
		&road.Target,
		&road.KM,
		&road.KMH,
		&road.Cost,
		&road.ReverseCost,
		&road.X1,
		&road.Y1,
		&road.X2,
		&road.Y2,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRoadNotFound)
		}

		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return road, nil
}
