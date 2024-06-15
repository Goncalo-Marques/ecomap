package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

// CreateTemporaryTableRoadNetworkWithBuffer creates a temporary table that is dropped when the transaction is
// committed. The name of the table must be unique per storage session. This new table represents a buffer of the road
// network, taking into account the convex hull of the given vertices.
func (s *store) CreateTemporaryTableRoadNetworkWithBuffer(ctx context.Context, tx pgx.Tx, tableName string, verticesGeometry []domain.GeoJSONGeometryPoint) error {
	_, err := tx.Exec(ctx, fmt.Sprintf(`
		CREATE TEMP TABLE %s
		ON COMMIT DROP 
		AS
			WITH vertices AS (
				SELECT json_array_elements($1) AS geom
			),
			buffer AS (
				SELECT ST_Buffer(ST_ConvexHull(ST_Collect(ST_GeomFromGeoJSON(geom))), 0.3)
				FROM vertices
			)
			SELECT * FROM road_network
			WHERE ST_Contains((SELECT * FROM buffer), geom_way)
	`, tableName),
		verticesGeometry,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	return nil
}

// CreateVerticesCloseToRoadNetwork creates new vertices by dividing the existing road network by the nearest edge to
// each of the given vertices.
func (s *store) CreateVerticesCloseToRoadNetwork(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, verticesGeometry []domain.GeoJSONGeometryPoint) ([]int, error) {
	if len(verticesGeometry) == 0 {
		return nil, nil
	}
	batch := new(pgx.Batch)

	for i := 0; i < len(verticesGeometry); i++ {
		geoJSON, err := json.Marshal(verticesGeometry[i])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}

		batch.Queue(fmt.Sprintf(`
			WITH vertex AS (
				SELECT ST_GeomFromGeoJSON($1) AS geom
			),
			closest_way AS (
				SELECT *
				FROM pgr_findCloseEdges(
					$$SELECT id, geom_way as geom FROM %s WHERE clazz > 20$$,
					(SELECT geom FROM vertex),
					0.5
				) AS ce
				INNER JOIN %s AS rn ON ce.edge_id = rn.id
			),
			split_fraction AS (
				SELECT ST_LineLocatePoint(
					(SELECT geom_way FROM closest_way),
					(SELECT geom FROM vertex)
				) AS fraction
			),
			split_vertex AS (
				SELECT ST_LineInterpolatePoint(
					(SELECT geom_way FROM closest_way),
					(SELECT fraction FROM split_fraction)
				) AS geom
			),
			split_way AS (
				SELECT ST_Split(
					ST_Snap(
						(SELECT geom_way FROM closest_way),
						(SELECT geom FROM split_vertex),
						1
					),
					(SELECT geom FROM split_vertex)
				) AS geom
			),
			split_way_geom AS (
				SELECT
					ST_GeometryN((SELECT geom FROM split_way), 1) AS geom1,
					ST_GeometryN((SELECT geom FROM split_way), 2) AS geom2
			),
			new_way1 AS (
				SELECT
					(SELECT cost FROM closest_way) * ST_Length(geom1) / ST_Length((SELECT geom_way FROM closest_way)) AS cost,
					(SELECT reverse_cost FROM closest_way) * ST_Length(geom1) / ST_Length((SELECT geom_way FROM closest_way)) AS reverse_cost,
					ST_X(ST_StartPoint(geom1)) AS x1,
					ST_Y(ST_StartPoint(geom1)) AS y1,
					ST_X(ST_EndPoint(geom1)) AS x2,
					ST_Y(ST_EndPoint(geom1)) AS y2
				FROM split_way_geom
			),
			new_way2 AS (
				SELECT
					(SELECT cost FROM closest_way) * ST_Length(geom2) / ST_Length((SELECT geom_way FROM closest_way)) AS cost,
					(SELECT reverse_cost FROM closest_way) * ST_Length(geom2) / ST_Length((SELECT geom_way FROM closest_way)) AS reverse_cost,
					ST_X(ST_StartPoint(geom2)) AS x1,
					ST_Y(ST_StartPoint(geom2)) AS y1,
					ST_X(ST_EndPoint(geom2)) AS x2,
					ST_Y(ST_EndPoint(geom2)) AS y2
				FROM split_way_geom
			),
			delete_old_way AS (
				DELETE FROM %s
				WHERE id = (SELECT id FROM closest_way)
					AND (SELECT geom2 IS NOT NULL FROM split_way_geom) -- Prevent insertion if the path cannot be split.
			),
			insert1 AS (
				SELECT 
					(SELECT id FROM closest_way) AS id, -- Keep the original ID.
					(SELECT source FROM closest_way) AS source,
					%d AS target, -- New random ID.
					(SELECT cost FROM new_way1),
					(SELECT reverse_cost FROM new_way1),
					(SELECT x1 FROM new_way1),
					(SELECT y1 FROM new_way1),
					(SELECT x2 FROM new_way1),
					(SELECT y2 FROM new_way1),
					(SELECT geom1 FROM split_way_geom) AS geom_way,
					osm_id, osm_name, osm_meta, osm_source_id, osm_target_id, clazz, flags, km, kmh 
				FROM closest_way 
			),
			insert2 AS (
				SELECT 
					%d AS id, -- New random ID.
					(SELECT target FROM insert1) AS source,
					(SELECT target FROM closest_way) AS target,
					(SELECT cost FROM new_way2),
					(SELECT reverse_cost FROM new_way2),
					(SELECT x1 FROM new_way2),
					(SELECT y1 FROM new_way2),
					(SELECT x2 FROM new_way2),
					(SELECT y2 FROM new_way2),
					(SELECT geom2 FROM split_way_geom) AS geom_way,
					osm_id, osm_name, osm_meta, osm_source_id, osm_target_id, clazz, flags, km, kmh 
				FROM closest_way
			),
			insert_new_ways AS (
				INSERT INTO %s (id, source, target, cost, reverse_cost, x1, y1, x2, y2, geom_way, osm_id, osm_name, osm_meta, osm_source_id, osm_target_id, clazz, flags, km, kmh)
				SELECT * FROM insert1
					WHERE (SELECT geom2 IS NOT NULL FROM split_way_geom) -- Prevent insertion if the path cannot be split.
				UNION ALL
				SELECT * FROM insert2
					WHERE (SELECT geom2 IS NOT NULL FROM split_way_geom) -- Prevent insertion if the path cannot be split.
			)
			SELECT CASE 
				WHEN (SELECT geom2 IS NULL FROM split_way_geom) AND (SELECT fraction FROM split_fraction) <= 0.5 
					THEN (SELECT source FROM closest_way)
				WHEN (SELECT geom2 IS NULL FROM split_way_geom) AND (SELECT fraction FROM split_fraction) > 0.5 
					THEN (SELECT target FROM closest_way)
				ELSE (SELECT target FROM insert1)
			END
		`,
			roadNetworkTableName,
			roadNetworkTableName,
			roadNetworkTableName,
			-i-1, // New IDs are represented by a negative integer to quickly avoid conflicts with already existing IDs.
			-i-1,
			roadNetworkTableName,
		),
			string(geoJSON),
		)
	}

	batchResult := tx.SendBatch(ctx, batch)
	defer batchResult.Close()

	vertexIDs := make([]int, len(verticesGeometry))

	for i := 0; i < len(batch.QueuedQueries); i++ {
		var vertexID int

		err := batchResult.QueryRow().Scan(&vertexID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
		}

		vertexIDs[i] = vertexID
	}

	return vertexIDs, nil
}

// GetRoadVerticesTSP executes a query to return the sequential vertex identifiers using the TSP algorithm and the A*
// cost matrix.
func (s *store) GetRoadVerticesTSP(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, vertexIDs []int, startVertexID, endVertexID int, directed bool) ([]int, error) {
	if len(vertexIDs) == 0 {
		return nil, nil
	}

	strVertexIDs := make([]string, len(vertexIDs))
	for i, id := range vertexIDs {
		strVertexIDs[i] = strconv.Itoa(id)
	}

	sqlMatrix := fmt.Sprintf(`
		$$SELECT * FROM pgr_aStarCostMatrix(
			'SELECT id, source, target, cost, reverse_cost, x1, y1, x2, y2 FROM %s',
			'{%s}'::bigint[],
			directed => %t
		)$$
	`,
		roadNetworkTableName,
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
func (s *store) GetRoadsGeometryAStar(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, seqVertexIDs []int, directed bool) ([]domain.GeoJSONGeometryLineString, error) {
	if len(seqVertexIDs) == 0 {
		return nil, nil
	}

	batch := new(pgx.Batch)

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
				'SELECT id, source, target, cost, reverse_cost, x1, y1, x2, y2 FROM %s',
				%d, %d,
				directed => %t
			) AS a
			INNER JOIN %s AS rn ON a.edge = rn.id
		`,
			roadNetworkTableName,
			prevVertexID,
			seqVertexIDs[i],
			directed,
			roadNetworkTableName,
		))

		prevVertexID = seqVertexIDs[i]
	}

	batchShortestPathsResult := tx.SendBatch(ctx, batch)
	defer batchShortestPathsResult.Close()

	var geoJSONLineStrings []domain.GeoJSONGeometryLineString
	for i := 0; i < len(batch.QueuedQueries); i++ {
		rows, err := batchShortestPathsResult.Query()
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
