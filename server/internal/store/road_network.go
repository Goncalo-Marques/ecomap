package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// GetRoadByGeometry executes a query to return the road that is closest to the given geometry in the road network.
func (s *store) GetRoadByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Road, error) {
	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	// Only consider roads that are not motorways, trunk roads or primary roads (clazz > 20).
	row := tx.QueryRow(ctx, `
		SELECT rn.id, rn.osm_id, rn.osm_name, rn.osm_meta, rn.osm_source_id, rn.osm_target_id, rn.clazz, rn.flags, rn.source, rn.target, rn.km, rn.kmh, rn.cost, rn.reverse_cost, rn.x1, rn.y1, rn.x2, rn.y2
		FROM pgr_findCloseEdges(
			$$SELECT id, geom_way as geom FROM road_network WHERE clazz > 20$$,
			ST_GeomFromGeoJSON($1),
			0.5
		) AS ce
		INNER JOIN road_network AS rn ON ce.edge_id = rn.id
	`,
		string(geoJSON),
	)

	road, err := getRoadFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRoadNotFound)
		}

		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return road, nil
}

// GetRoadByWarehouseID executes a query to return the road associated with the warehouse with the specified identifier.
func (s *store) GetRoadByWarehouseID(ctx context.Context, tx pgx.Tx, warehouseID uuid.UUID) (domain.Road, error) {
	row := tx.QueryRow(ctx, `
		SELECT rn.id, rn.osm_id, rn.osm_name, rn.osm_meta, rn.osm_source_id, rn.osm_target_id, rn.clazz, rn.flags, rn.source, rn.target, rn.km, rn.kmh, rn.cost, rn.reverse_cost, rn.x1, rn.y1, rn.x2, rn.y2
		FROM warehouses AS w
		INNER JOIN road_network AS rn ON w.road_id = rn.id
		WHERE w.id = $1
	`,
		warehouseID,
	)

	road, err := getRoadFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRoadNotFound)
		}

		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return road, nil
}

// GetRoadByLandfillID executes a query to return the road associated with the landfill with the specified identifier.
func (s *store) GetRoadByLandfillID(ctx context.Context, tx pgx.Tx, landfillID uuid.UUID) (domain.Road, error) {
	row := tx.QueryRow(ctx, `
		SELECT rn.id, rn.osm_id, rn.osm_name, rn.osm_meta, rn.osm_source_id, rn.osm_target_id, rn.clazz, rn.flags, rn.source, rn.target, rn.km, rn.kmh, rn.cost, rn.reverse_cost, rn.x1, rn.y1, rn.x2, rn.y2
		FROM landfills AS l
		INNER JOIN road_network AS rn ON l.road_id = rn.id
		WHERE l.id = $1
	`,
		landfillID,
	)

	road, err := getRoadFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrRoadNotFound)
		}

		return domain.Road{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return road, nil
}

// GetContainerRoadsByRouteID executes a query to return the roads associated with the route containers.
func (s *store) GetContainerRoadsByRouteID(ctx context.Context, tx pgx.Tx, routeID uuid.UUID) ([]domain.Road, error) {
	rows, err := tx.Query(ctx, `
		SELECT rn.id, rn.osm_id, rn.osm_name, rn.osm_meta, rn.osm_source_id, rn.osm_target_id, rn.clazz, rn.flags, rn.source, rn.target, rn.km, rn.kmh, rn.cost, rn.reverse_cost, rn.x1, rn.y1, rn.x2, rn.y2
		FROM routes AS r
		INNER JOIN routes_containers AS rc ON r.id = rc.route_id
		INNER JOIN containers AS c ON rc.container_id = c.id
		INNER JOIN road_network AS rn ON c.road_id = rn.id
		WHERE r.id = $1
	`,
		routeID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	roads, err := getRoadsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return roads, nil
}

// GetRoadVerticesTSP executes a query to return the sequential vertex identifiers using the TSP algorithm and the A*
// cost matrix.
func (s *store) GetRoadVerticesTSP(ctx context.Context, tx pgx.Tx, vertexIDs []int, startVertexID, endVertexID int, directed bool) ([]int, error) {
	if len(vertexIDs) == 0 {
		return nil, nil
	}

	strVertexIDs := make([]string, len(vertexIDs))
	for i, id := range vertexIDs {
		strVertexIDs[i] = strconv.Itoa(id)
	}

	sqlMatrix := fmt.Sprintf(`
		$$SELECT * FROM pgr_aStarCostMatrix(
			'WITH bounding_box AS (
				SELECT ST_Buffer(ST_ConvexHull(ST_Collect(geom_vertex)), 0.3) FROM road_network_vertex WHERE id IN (%s)
			)
			SELECT id, source, target, cost, reverse_cost, x1, y1, x2, y2 
				FROM road_network
				WHERE ST_Contains((SELECT * FROM bounding_box), geom_way)',
			'{%s}'::bigint[],
			directed => %t
		)$$
	`,
		strings.Join(strVertexIDs, ", "),
		strings.Join(strVertexIDs, ", "),
		directed,
	)

	rows, err := tx.Query(ctx, `
		SELECT node 
		FROM pgr_TSP(
			`+sqlMatrix+`,
			start_id => $1,
			end_id => $2
		)
	`,
		startVertexID,
		endVertexID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	var seqVertexIDs []int
	for rows.Next() {
		var vertexID int
		err := rows.Scan(&vertexID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
		}

		seqVertexIDs = append(seqVertexIDs, vertexID)
	}

	return seqVertexIDs, nil
}

// GetRoadsGeometryAStar executes a query to return the geometry of sequential roads using the shortest path A*
// algorithm.
func (s *store) GetRoadsGeometryAStar(ctx context.Context, tx pgx.Tx, seqVertexIDs []int, directed bool) ([]domain.GeoJSONGeometryLineString, error) {
	if len(seqVertexIDs) == 0 {
		return nil, nil
	}

	batch := new(pgx.Batch)

	strSeqVertexIDs := make([]string, len(seqVertexIDs))
	for i, id := range seqVertexIDs {
		strSeqVertexIDs[i] = strconv.Itoa(id)
	}

	prevVertexID := seqVertexIDs[0]
	for i := 1; i < len(seqVertexIDs); i++ {
		batch.Queue(fmt.Sprintf(`
			SELECT ST_AsGeoJSON(
				CASE 
					WHEN a.node = rn.source THEN rn.geom_way
					ELSE ST_Reverse(rn.geom_way)
				END
			)::jsonb
			FROM pgr_aStar(
				'WITH bounding_box AS (
					SELECT ST_Buffer(ST_ConvexHull(ST_Collect(geom_vertex)), 0.3) FROM road_network_vertex WHERE id IN (%s)
				)
				SELECT id, source, target, cost, reverse_cost, x1, y1, x2, y2 
					FROM road_network
					WHERE ST_Contains((SELECT * FROM bounding_box), geom_way)',
				%d, %d,
				directed => %t
			) AS a
			INNER JOIN road_network AS rn ON a.edge = rn.id
		`,
			strings.Join(strSeqVertexIDs, ", "),
			prevVertexID,
			seqVertexIDs[i],
			directed,
		))

		prevVertexID = seqVertexIDs[i]
	}

	batchResult := tx.SendBatch(ctx, batch)
	defer batchResult.Close()

	var geoJSONLineStrings []domain.GeoJSONGeometryLineString
	for i := 0; i < len(batch.QueuedQueries); i++ {
		rows, err := batchResult.Query()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
		}
		defer rows.Close()

		for rows.Next() {
			var geometry domain.GeoJSONGeometryLineString
			err := rows.Scan(&geometry)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
			}

			geoJSONLineStrings = append(geoJSONLineStrings, geometry)
		}
	}

	return geoJSONLineStrings, nil
}

// getRoadFromRow returns the road by scanning the given row.
func getRoadFromRow(row pgx.Row) (domain.Road, error) {
	var road domain.Road
	err := row.Scan(
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
		return domain.Road{}, err
	}

	return road, nil
}

// getRoadsFromRows returns the roads by scanning the given rows.
func getRoadsFromRows(rows pgx.Rows) ([]domain.Road, error) {
	var roads []domain.Road
	for rows.Next() {
		road, err := getRoadFromRow(rows)
		if err != nil {
			return nil, err
		}

		roads = append(roads, road)
	}

	return roads, nil
}
