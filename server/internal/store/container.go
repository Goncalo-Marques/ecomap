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

// ListContainers executes a query to return the containers for the specified filter.
func (s *store) ListContainers(ctx context.Context, tx pgx.Tx, filter domain.ContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	filterFields := make([]string, 0, 3)
	argsWhere := make([]any, 0, 3)

	// Append the optional fields to filter.
	if filter.Category != nil {
		filterFields = append(filterFields, "c.category::text")
		argsWhere = append(argsWhere, containerCategoryFromDomain(*filter.Category))
	}
	if filter.WayName != nil {
		filterFields = append(filterFields, "rn.osm_name")
		argsWhere = append(argsWhere, *filter.WayName)
	}
	if filter.MunicipalityName != nil {
		filterFields = append(filterFields, "m.name")
		argsWhere = append(argsWhere, *filter.MunicipalityName)
	}

	sqlWhere := listSQLWhere(filterFields, filter.LogicalOperator)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(c.id) 
		FROM containers AS c
		LEFT JOIN road_network AS rn ON c.road_id = rn.id
		LEFT JOIN municipalities AS m ON c.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.ContainerPaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "c.created_at"
	switch domainSortField {
	case domain.ContainerPaginatedSortCategory:
		sortField = "c.category"
	case domain.ContainerPaginatedSortWayName:
		sortField = "rn.osm_name"
	case domain.ContainerPaginatedSortMunicipalityName:
		sortField = "m.name"
	case domain.ContainerPaginatedSortCreatedAt:
		sortField = "c.created_at"
	case domain.ContainerPaginatedSortModifiedAt:
		sortField = "c.modified_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT c.id, c.category, ST_AsGeoJSON(c.geom)::jsonb, rn.osm_name, m.name, c.created_at, c.modified_at
		FROM containers AS c
		LEFT JOIN road_network AS rn ON c.road_id = rn.id
		LEFT JOIN municipalities AS m ON c.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	containers, err := getContainersFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.Container]{
		Total:   total,
		Results: containers,
	}, nil
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

// PatchContainer executes a query to patch an container with the specified identifier and data.
func (s *store) PatchContainer(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableContainer domain.EditableContainerPatch, roadID, municipalityID *int) error {
	var category *string
	if editableContainer.Category != nil {
		c := containerCategoryFromDomain(*editableContainer.Category)
		category = &c
	}

	var geoJSON []byte
	var err error

	if editableContainer.GeoJSON != nil {
		var geometry domain.GeoJSONGeometryPoint
		if feature, ok := editableContainer.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}

		geoJSON, err = json.Marshal(geometry)
		if err != nil {
			return fmt.Errorf("%s: %w", descriptionFailedMarshalGeoJSON, err)
		}
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE containers SET
			category = coalesce($2, category),
			geom = coalesce(ST_GeomFromGeoJSON($3), geom),
			road_id = CASE 
					WHEN $3 IS NOT NULL THEN $4 
					ELSE road_id
				END,
			municipality_id = CASE 
					WHEN $3 IS NOT NULL THEN $5
					ELSE municipality_id
				END
		WHERE id = $1
	`,
		id,
		category,
		geoJSON,
		roadID,
		municipalityID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerNotFound)
	}

	return nil
}

// DeleteContainerByID executes a query to delete the container with the specified identifier.
func (s *store) DeleteContainerByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, `
		DELETE FROM containers
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		switch constraintNameFromError(err) {
		case constraintContainersReportsContainerIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerAssociatedWithContainerReport)
		case constraintUsersContainerBookmarksContainerIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerAssociatedWithUserContainerBookmark)
		case constraintRoutesContainersContainerIDFkey:
			return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerAssociatedWithRouteContainer)
		}

		return fmt.Errorf("%s: %w", descriptionFailedExec, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", descriptionFailedExec, domain.ErrContainerNotFound)
	}

	return nil
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

// getContainersFromRows returns the containers by scanning the given rows.
func getContainersFromRows(rows pgx.Rows) ([]domain.Container, error) {
	var containers []domain.Container
	for rows.Next() {
		container, err := getContainerFromRow(rows)
		if err != nil {
			return nil, err
		}

		containers = append(containers, container)
	}

	return containers, nil
}
