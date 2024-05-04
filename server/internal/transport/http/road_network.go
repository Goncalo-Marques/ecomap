package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errRoadNotFound = "way not found"
)

// GetWayByReverseGeocoding handles the http request to get a way by reverse geocoding.
func (h *handler) GetWayByReverseGeocoding(w http.ResponseWriter, r *http.Request, params spec.GetWayByReverseGeocodingParams) {
	ctx := r.Context()

	domainGeoJSONGeometryPoint, err := coordinatesToDomain(params.Coordinates)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errParamInvalidFormat, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainRoad, err := h.service.GetRoadByGeometry(ctx, domainGeoJSONGeometryPoint)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoadNotFound):
			notFound(w, errRoadNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	way := wayFromDomainRoad(domainRoad)
	responseBody, err := json.Marshal(way)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// wayFromDomainRoad returns a standardized way based on the road domain model.
func wayFromDomainRoad(road domain.Road) spec.Way {
	return spec.Way{
		Id:          road.ID,
		OsmId:       road.OsmID,
		OsmName:     road.OsmName,
		OsmMeta:     road.OsmMeta,
		OsmSourceId: road.OsmSourceID,
		OsmTargetId: road.OsmTargetID,
		Clazz:       road.Clazz,
		Flags:       road.Flags,
		Source:      road.Source,
		Target:      road.Target,
		Km:          road.KM,
		Kmh:         road.KMH,
		Cost:        road.Cost,
		ReverseCost: road.ReverseCost,
		X1:          road.X1,
		Y1:          road.Y1,
		X2:          road.X2,
		Y2:          road.Y2,
	}
}
