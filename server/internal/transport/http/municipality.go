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
	errMunicipalityNotFound = "municipality not found"
)

// GetMunicipalityByReverseGeocoding handles the http request to get a municipality by reverse geocoding.
func (h *handler) GetMunicipalityByReverseGeocoding(w http.ResponseWriter, r *http.Request, params spec.GetMunicipalityByReverseGeocodingParams) {
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

	domainMunicipality, err := h.service.GetMunicipalityByGeometry(ctx, domainGeoJSONGeometryPoint)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMunicipalityNotFound):
			notFound(w, errMunicipalityNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	municipality := municipalityFromDomain(domainMunicipality)
	responseBody, err := json.Marshal(municipality)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// municipalityFromDomain returns a standardized municipality based on the domain model.
func municipalityFromDomain(municipality domain.Municipality) spec.Municipality {
	return spec.Municipality{
		Id:        municipality.ID,
		FeatureId: municipality.FeatureID,
		Name:      municipality.Name,
		District:  municipality.District,
		Nuts1:     municipality.NUTS1,
		Nuts2:     municipality.NUTS2,
		Nuts3:     municipality.NUTS3,
		Area:      municipality.Area,
		Perimeter: municipality.Perimeter,
	}
}
