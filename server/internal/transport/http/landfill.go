package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errLandfillNotFound = "landfill not found"
)

// CreateLandfill handles the http request to create a landfill.
func (h *handler) CreateLandfill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var landfillPost spec.LandfillPost
	err = json.Unmarshal(requestBody, &landfillPost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableLandfill, err := landfillPostToDomain(landfillPost)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainLandfill, err := h.service.CreateLandfill(ctx, domainEditableLandfill)
	if err != nil {
		internalServerError(w)
		return
	}

	landfill, err := landfillFromDomain(domainLandfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(landfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListLandfills handles the http request to list landfills.
func (h *handler) ListLandfills(w http.ResponseWriter, r *http.Request, params spec.ListLandfillsParams) {
	ctx := r.Context()

	domainLandfillsFilter := listLandfillsParamsToDomain(params)
	domainPaginatedLandfills, err := h.service.ListLandfills(ctx, domainLandfillsFilter)
	if err != nil {
		var domainErrFilterValueInvalid *domain.ErrFilterValueInvalid

		switch {
		case errors.As(err, &domainErrFilterValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFilterValueInvalid, domainErrFilterValueInvalid.FilterName))
		default:
			internalServerError(w)
		}

		return
	}

	landfillsPaginated, err := landfillsPaginatedFromDomain(domainPaginatedLandfills)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(landfillsPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetLandfillByID handles the http request to get a landfill by ID.
func (h *handler) GetLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	ctx := r.Context()

	domainLandfill, err := h.service.GetLandfillByID(ctx, landfillID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			notFound(w, errLandfillNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	landfill, err := landfillFromDomain(domainLandfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(landfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchLandfillByID handles the http request to modify a landfill by ID.
func (h *handler) PatchLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var landfillPatch spec.LandfillPatch
	err = json.Unmarshal(requestBody, &landfillPatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableLandfill, err := landfillPatchToDomain(landfillPatch)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		default:
			internalServerError(w)
		}

		return
	}

	domainLandfill, err := h.service.PatchLandfill(ctx, landfillID, domainEditableLandfill)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			notFound(w, errLandfillNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	landfill, err := landfillFromDomain(domainLandfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(landfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteLandfillByID handles the http request to delete a landfill by ID.
func (h *handler) DeleteLandfillByID(w http.ResponseWriter, r *http.Request, landfillID spec.LandfillIdPathParam) {
	ctx := r.Context()

	domainLandfill, err := h.service.DeleteLandfillByID(ctx, landfillID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLandfillNotFound):
			notFound(w, errLandfillNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	landfill, err := landfillFromDomain(domainLandfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(landfill)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// landfillPostToDomain returns a domain editable landfill based on the standardized landfill post.
func landfillPostToDomain(landfillPost spec.LandfillPost) (domain.EditableLandfill, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(&landfillPost.GeoJson)
	if err != nil {
		return domain.EditableLandfill{}, err
	}

	return domain.EditableLandfill{
		GeoJSON: geoJSON,
	}, nil
}

// landfillPatchToDomain returns a domain patchable landfill based on the standardized landfill patch.
func landfillPatchToDomain(landfillPatch spec.LandfillPatch) (domain.EditableLandfillPatch, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(landfillPatch.GeoJson)
	if err != nil {
		return domain.EditableLandfillPatch{}, err
	}

	return domain.EditableLandfillPatch{
		GeoJSON: geoJSON,
	}, nil
}

// listLandfillsParamsToDomain returns a domain landfills paginated filter based on the standardized list landfills parameters.
func listLandfillsParamsToDomain(params spec.ListLandfillsParams) domain.LandfillsPaginatedFilter {
	domainSort := domain.LandfillPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListLandfillsParamsSortWayName:
			domainSort = domain.LandfillPaginatedSortWayName
		case spec.ListLandfillsParamsSortMunicipalityName:
			domainSort = domain.LandfillPaginatedSortMunicipalityName
		case spec.ListLandfillsParamsSortCreatedAt:
			domainSort = domain.LandfillPaginatedSortCreatedAt
		case spec.ListLandfillsParamsSortModifiedAt:
			domainSort = domain.LandfillPaginatedSortModifiedAt
		default:
			domainSort = domain.LandfillPaginatedSort(*params.Sort)
		}
	}

	return domain.LandfillsPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		LocationName: params.LocationName,
	}
}

// landfillFromDomain returns a standardized landfill based on the domain model.
func landfillFromDomain(landfill domain.Landfill) (spec.Landfill, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(landfill.GeoJSON)
	if err != nil {
		return spec.Landfill{}, err
	}

	return spec.Landfill{
		Id:         landfill.ID,
		GeoJson:    geoJSON,
		CreatedAt:  landfill.CreatedAt,
		ModifiedAt: landfill.ModifiedAt,
	}, nil
}

// landfillsFromDomain returns standardized landfills based on the domain model.
func landfillsFromDomain(landfills []domain.Landfill) ([]spec.Landfill, error) {
	specLandfills := make([]spec.Landfill, len(landfills))
	var err error

	for i, landfill := range landfills {
		specLandfills[i], err = landfillFromDomain(landfill)
		if err != nil {
			return []spec.Landfill{}, err
		}
	}

	return specLandfills, nil
}

// landfillsPaginatedFromDomain returns a standardized landfills paginated response based on the domain model.
func landfillsPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Landfill]) (spec.LandfillsPaginated, error) {
	landfills, err := landfillsFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.LandfillsPaginated{}, err
	}

	return spec.LandfillsPaginated{
		Total:     paginatedResponse.Total,
		Landfills: landfills,
	}, nil
}
