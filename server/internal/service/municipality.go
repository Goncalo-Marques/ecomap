package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedGetMunicipalityByGeometry = "service: failed to get municipality by geometry"
)

// GetMunicipalityByGeometry returns the municipality that contains the given geometry point.
func (s *service) GetMunicipalityByGeometry(ctx context.Context, geometry domain.GeoJSONGeometryPoint) (domain.Municipality, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetMunicipalityByGeometry"),
	}

	var municipality domain.Municipality
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		municipality, err = s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMunicipalityNotFound):
			return domain.Municipality{}, logInfoAndWrapError(ctx, err, descriptionFailedGetMunicipalityByGeometry, logAttrs...)
		default:
			return domain.Municipality{}, logAndWrapError(ctx, err, descriptionFailedGetMunicipalityByGeometry, logAttrs...)
		}
	}

	return municipality, nil
}
