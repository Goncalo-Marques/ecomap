package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintWarehousesTrucksTruckIDFkey = "warehouses_trucks_truck_id_fkey"
)

// ListWarehouseTrucks executes a query to return the trucks associated with the warehouse with the specified identifier.
func (s *store) ListWarehouseTrucks(ctx context.Context, tx pgx.Tx, id uuid.UUID) ([]domain.Truck, error) {
	rows, err := tx.Query(ctx, `
		SELECT t.id, t.make, t.model, t.license_plate, t.person_capacity, ST_AsGeoJSON(t.geom)::jsonb, rn.osm_name, m.name, t.created_at, t.modified_at 
		FROM warehouses_trucks AS wt
		INNER JOIN trucks AS t ON wt.truck_id = t.id
		LEFT JOIN road_network AS rn ON t.road_id = rn.id
		LEFT JOIN municipalities AS m ON t.municipality_id = m.id
		WHERE wt.warehouse_id = $1 
	`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	trucks, err := getTrucksFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return trucks, nil
}
