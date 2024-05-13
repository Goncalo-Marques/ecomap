package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedCreateLandfill     = "service: failed to create landfill"
	descriptionFailedListLandfills      = "service: failed to list landfills"
	descriptionFailedGetLandfillByID    = "service: failed to get landfill by id"
	descriptionFailedPatchLandfill      = "service: failed to patch landfill"
	descriptionFailedDeleteLandfillByID = "service: failed to delete landfill by id"
)

// CreateLandfill creates a new landfill with the specified data.
func (s *service) CreateLandfill(ctx context.Context, editableLandfill domain.EditableLandfill) (domain.Landfill, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateLandfill"),
	}

	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableLandfill.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	var landfill domain.Landfill

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrRoadNotFound) {
				return err
			}
		} else {
			roadID = &road.ID
		}

		var municipalityID *int
		municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrMunicipalityNotFound) {
				return err
			}
		} else {
			municipalityID = &municipality.ID
		}

		id, err := s.store.CreateLandfill(ctx, tx, editableLandfill, roadID, municipalityID)
		if err != nil {
			return err
		}

		landfill, err = s.store.GetLandfillByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Landfill{}, logAndWrapError(ctx, err, descriptionFailedCreateLandfill, logAttrs...)
	}

	return landfill, nil
}

// ListLandfills returns the landfills with the specified filter.
func (s *service) ListLandfills(ctx context.Context, filter domain.LandfillsPaginatedFilter) (domain.PaginatedResponse[domain.Landfill], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListLandfills"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Landfill]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Landfill]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Landfill]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Landfill]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedLandfills domain.PaginatedResponse[domain.Landfill]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedLandfills, err = s.store.ListLandfills(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Landfill]{}, logAndWrapError(ctx, err, descriptionFailedListLandfills, logAttrs...)
	}

	return paginatedLandfills, nil
}

// GetLandfillByID returns the landfill with the specified identifier.
func (s *service) GetLandfillByID(ctx context.Context, id uuid.UUID) (domain.Landfill, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetLandfillByID"),
		slog.String(logging.LandfillID, id.String()),
	}

	var landfill domain.Landfill
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		landfill, err = s.store.GetLandfillByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			return domain.Landfill{}, logInfoAndWrapError(ctx, err, descriptionFailedGetLandfillByID, logAttrs...)
		default:
			return domain.Landfill{}, logAndWrapError(ctx, err, descriptionFailedGetLandfillByID, logAttrs...)
		}
	}

	return landfill, nil
}

// PatchLandfill modifies the landfill with the specified identifier.
func (s *service) PatchLandfill(ctx context.Context, id uuid.UUID, editableLandfill domain.EditableLandfillPatch) (domain.Landfill, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchLandfill"),
		slog.String(logging.LandfillID, id.String()),
	}

	var geometry domain.GeoJSONGeometryPoint
	if editableLandfill.GeoJSON != nil {
		if feature, ok := editableLandfill.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}
	}

	var landfill domain.Landfill

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		var municipalityID *int

		if editableLandfill.GeoJSON != nil {
			road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrRoadNotFound) {
					return err
				}
			} else {
				roadID = &road.ID
			}

			municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrMunicipalityNotFound) {
					return err
				}
			} else {
				municipalityID = &municipality.ID
			}
		}

		err := s.store.PatchLandfill(ctx, tx, id, editableLandfill, roadID, municipalityID)
		if err != nil {
			return err
		}

		landfill, err = s.store.GetLandfillByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			return domain.Landfill{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchLandfill, logAttrs...)
		default:
			return domain.Landfill{}, logAndWrapError(ctx, err, descriptionFailedPatchLandfill, logAttrs...)
		}
	}

	return landfill, nil
}

// DeleteLandfillByID deletes the landfill with the specified identifier.
func (s *service) DeleteLandfillByID(ctx context.Context, id uuid.UUID) (domain.Landfill, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteLandfillByID"),
		slog.String(logging.LandfillID, id.String()),
	}

	var landfill domain.Landfill
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		landfill, err = s.store.GetLandfillByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteLandfillByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			return domain.Landfill{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteLandfillByID, logAttrs...)
		default:
			return domain.Landfill{}, logAndWrapError(ctx, err, descriptionFailedDeleteLandfillByID, logAttrs...)
		}
	}

	return landfill, nil
}
