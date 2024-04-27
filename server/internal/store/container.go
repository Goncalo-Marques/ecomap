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

// CreateContainer executes a query to create a container with the specified data.
func (s *store) CreateContainer(ctx context.Context, tx pgx.Tx, editableContainer domain.EditableContainer, roadID, municipalityID *int) (uuid.UUID, error) {
	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableContainer.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	geoJSON, err := json.Marshal(geometry)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO containers (category, geom, road_id, municipality_id)
		VALUES ($1, ST_GeomFromGeoJSON($2), $3, $4) 
		RETURNING id
	`,
		containerCategoryFromDomain(editableContainer.Category),
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

// GetContainerByID executes a query to return the container with the specified identifier.
func (s *store) GetContainerByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Container, error) {
	row := tx.QueryRow(ctx, `
		SELECT c.id, c.category, ST_AsGeoJSON(c.geom)::jsonb, rn.osm_name, m.name, c.created_at, c.modified_at 
		FROM containers AS c
		LEFT JOIN road_network AS rn ON c.road_id = rn.id
		LEFT JOIN municipalities AS m ON c.municipality_id = m.id
		WHERE c.id = $1 
	`,
		id,
	)

	container, err := getContainerFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Container{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, domain.ErrContainerNotFound)
		}

		return domain.Container{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	return container, nil
}

// containerCategoryFromDomain returns a store container category based on the domain model.
func containerCategoryFromDomain(category domain.ContainerCategory) string {
	switch category {
	case domain.ContainerCategoryGeneral:
		return "general"
	case domain.ContainerCategoryPaper:
		return "paper"
	case domain.ContainerCategoryPlastic:
		return "plastic"
	case domain.ContainerCategoryMetal:
		return "metal"
	case domain.ContainerCategoryGlass:
		return "glass"
	case domain.ContainerCategoryOrganic:
		return "organic"
	case domain.ContainerCategoryHazardous:
		return "hazardous"
	default:
		return string(category)
	}
}

// containerCategoryToDomain returns a domain container category based on the store model.
func containerCategoryToDomain(category string) domain.ContainerCategory {
	switch category {
	case "general":
		return domain.ContainerCategoryGeneral
	case "paper":
		return domain.ContainerCategoryPaper
	case "plastic":
		return domain.ContainerCategoryPlastic
	case "metal":
		return domain.ContainerCategoryMetal
	case "glass":
		return domain.ContainerCategoryGlass
	case "organic":
		return domain.ContainerCategoryOrganic
	case "hazardous":
		return domain.ContainerCategoryHazardous
	default:
		return domain.ContainerCategory(category)
	}
}

// getContainerFromRow returns the container by scanning the given row.
func getContainerFromRow(row pgx.Row) (domain.Container, error) {
	var container domain.Container
	var category string
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&container.ID,
		&category,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&container.CreatedAt,
		&container.ModifiedAt,
	)
	if err != nil {
		return domain.Container{}, err
	}

	container.Category = containerCategoryToDomain(category)

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	container.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return container, nil
}
