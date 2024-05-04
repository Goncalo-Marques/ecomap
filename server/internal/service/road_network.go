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
	descriptionFailedGetRoadByGeometry = "service: failed to get road by geometry"
)

// GetRoadByGeometry returns the closest road to the given geometry point.
func (s *service) GetRoadByGeometry(ctx context.Context, geometry domain.GeoJSONGeometryPoint) (domain.Road, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetRoadByGeometry"),
	}

	var road domain.Road
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		road, err = s.store.GetRoadByGeometry(ctx, tx, geometry)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoadNotFound):
			return domain.Road{}, logInfoAndWrapError(ctx, err, descriptionFailedGetRoadByGeometry, logAttrs...)
		default:
			return domain.Road{}, logAndWrapError(ctx, err, descriptionFailedGetRoadByGeometry, logAttrs...)
		}
	}

	return road, nil
}
