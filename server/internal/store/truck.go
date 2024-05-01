package store

import (
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

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
