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
	descriptionFailedCreateTruck     = "service: failed to create truck"
	descriptionFailedListTrucks      = "service: failed to list trucks"
	descriptionFailedGetTruckByID    = "service: failed to get truck by id"
	descriptionFailedPatchTruck      = "service: failed to patch truck"
	descriptionFailedDeleteTruckByID = "service: failed to delete truck by id"
)

// CreateTruck creates a new truck with the specified data.
func (s *service) CreateTruck(ctx context.Context, editableTruck domain.EditableTruck) (domain.Truck, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateTruck"),
		slog.String(logging.TruckMake, string(editableTruck.Make)),
		slog.String(logging.TruckModel, string(editableTruck.Model)),
		slog.String(logging.TruckLicensePlate, string(editableTruck.LicensePlate)),
		slog.Int(logging.TruckPersonCapacity, int(editableTruck.PersonCapacity)),
	}

	if !editableTruck.Make.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldMake}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableTruck.Model.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldModel}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableTruck.LicensePlate.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldLicensePlate}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableTruck.PersonCapacity.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldPersonCapacity}, descriptionInvalidFieldValue, logAttrs...)
	}

	var truck domain.Truck

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		geometry := geometryPointFromGeoJSON(editableTruck.GeoJSON)

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

		id, err := s.store.CreateTruck(ctx, tx, editableTruck, roadID, municipalityID)
		if err != nil {
			return err
		}

		truck, err = s.store.GetTruckByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Truck{}, logAndWrapError(ctx, err, descriptionFailedCreateTruck, logAttrs...)
	}

	return truck, nil
}

// ListTrucks returns the trucks with the specified filter.
func (s *service) ListTrucks(ctx context.Context, filter domain.TrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListTrucks"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedTrucks domain.PaginatedResponse[domain.Truck]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedTrucks, err = s.store.ListTrucks(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Truck]{}, logAndWrapError(ctx, err, descriptionFailedListTrucks, logAttrs...)
	}

	return paginatedTrucks, nil
}

// GetTruckByID returns the truck with the specified identifier.
func (s *service) GetTruckByID(ctx context.Context, id uuid.UUID) (domain.Truck, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetTruckByID"),
		slog.String(logging.TruckID, id.String()),
	}

	var truck domain.Truck
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		truck, err = s.store.GetTruckByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound):
			return domain.Truck{}, logInfoAndWrapError(ctx, err, descriptionFailedGetTruckByID, logAttrs...)
		default:
			return domain.Truck{}, logAndWrapError(ctx, err, descriptionFailedGetTruckByID, logAttrs...)
		}
	}

	return truck, nil
}

// PatchTruck modifies the truck with the specified identifier.
func (s *service) PatchTruck(ctx context.Context, id uuid.UUID, editableTruck domain.EditableTruckPatch) (domain.Truck, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchTruck"),
		slog.String(logging.TruckID, id.String()),
	}

	if editableTruck.Make != nil && !editableTruck.Make.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldMake}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableTruck.Model != nil && !editableTruck.Model.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldModel}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableTruck.LicensePlate != nil && !editableTruck.LicensePlate.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldLicensePlate}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableTruck.PersonCapacity != nil && !editableTruck.PersonCapacity.Valid() {
		return domain.Truck{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldPersonCapacity}, descriptionInvalidFieldValue, logAttrs...)
	}

	var truck domain.Truck
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		geometry := geometryPointFromGeoJSON(editableTruck.GeoJSON)

		var roadID *int
		var municipalityID *int

		if editableTruck.GeoJSON != nil {
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

		err = s.store.PatchTruck(ctx, tx, id, editableTruck, roadID, municipalityID)
		if err != nil {
			return err
		}

		truck, err = s.store.GetTruckByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound):
			return domain.Truck{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchTruck, logAttrs...)
		default:
			return domain.Truck{}, logAndWrapError(ctx, err, descriptionFailedPatchTruck, logAttrs...)
		}
	}

	return truck, nil
}

// DeleteTruckByID deletes the truck with the specified identifier.
func (s *service) DeleteTruckByID(ctx context.Context, id uuid.UUID) (domain.Truck, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteTruckByID"),
		slog.String(logging.TruckID, id.String()),
	}

	var truck domain.Truck
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		truck, err = s.store.GetTruckByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteTruckByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound),
			errors.Is(err, domain.ErrTruckAssociatedWithWarehouseTruck),
			errors.Is(err, domain.ErrTruckAssociatedWithRoute):
			return domain.Truck{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteTruckByID, logAttrs...)
		default:
			return domain.Truck{}, logAndWrapError(ctx, err, descriptionFailedDeleteTruckByID, logAttrs...)
		}
	}

	return truck, nil
}
