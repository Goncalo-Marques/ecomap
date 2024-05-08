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
	errTruckNotFound                     = "truck not found"
	errTruckAssociatedWithWarehouseTruck = "truck associated with warehouse"
	errTruckAssociatedWithRoute          = "truck associated with route"
)

// CreateTruck handles the http request to create a truck.
func (h *handler) CreateTruck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var truckPost spec.TruckPost
	err = json.Unmarshal(requestBody, &truckPost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableTruck, err := truckPostToDomain(truckPost)
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

	domainTruck, err := h.service.CreateTruck(ctx, domainEditableTruck)
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

	truck, err := truckFromDomain(domainTruck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(truck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListTrucks handles the http request to list trucks.
func (h *handler) ListTrucks(w http.ResponseWriter, r *http.Request, params spec.ListTrucksParams) {
	ctx := r.Context()

	domainTrucksFilter := listTrucksParamsToDomain(params)
	domainPaginatedTrucks, err := h.service.ListTrucks(ctx, domainTrucksFilter)
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

	trucksPaginated, err := trucksPaginatedFromDomain(domainPaginatedTrucks)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(trucksPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetTruckByID handles the http request to get a truck by ID.
func (h *handler) GetTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	ctx := r.Context()

	domainTruck, err := h.service.GetTruckByID(ctx, truckID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound):
			notFound(w, errTruckNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	truck, err := truckFromDomain(domainTruck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(truck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchTruckByID handles the http request to modify a truck by ID.
func (h *handler) PatchTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var truckPatch spec.TruckPatch
	err = json.Unmarshal(requestBody, &truckPatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableTruck, err := truckPatchToDomain(truckPatch)
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

	domainTruck, err := h.service.PatchTruck(ctx, truckID, domainEditableTruck)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrTruckNotFound):
			notFound(w, errTruckNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	truck, err := truckFromDomain(domainTruck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(truck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteTruckByID handles the http request to delete a truck by ID.
func (h *handler) DeleteTruckByID(w http.ResponseWriter, r *http.Request, truckID spec.TruckIdPathParam) {
	ctx := r.Context()

	domainTruck, err := h.service.DeleteTruckByID(ctx, truckID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound):
			notFound(w, errTruckNotFound)
		case errors.Is(err, domain.ErrTruckAssociatedWithWarehouseTruck):
			conflict(w, errTruckAssociatedWithWarehouseTruck)
		case errors.Is(err, domain.ErrTruckAssociatedWithRoute):
			conflict(w, errTruckAssociatedWithRoute)
		default:
			internalServerError(w)
		}

		return
	}

	truck, err := truckFromDomain(domainTruck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(truck)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// truckPostToDomain returns a domain editable truck based on the standardized truck post.
func truckPostToDomain(truckPost spec.TruckPost) (domain.EditableTruck, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(&truckPost.GeoJson)
	if err != nil {
		return domain.EditableTruck{}, err
	}

	return domain.EditableTruck{
		Make:           domain.TruckMake(truckPost.Make),
		Model:          domain.TruckModel(truckPost.Model),
		LicensePlate:   domain.TruckLicensePlate(truckPost.LicensePlate),
		PersonCapacity: domain.TruckPersonCapacity(truckPost.PersonCapacity),
		GeoJSON:        geoJSON,
	}, nil
}

// truckPatchToDomain returns a domain patchable truck based on the standardized truck patch.
func truckPatchToDomain(truckPatch spec.TruckPatch) (domain.EditableTruckPatch, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(truckPatch.GeoJson)
	if err != nil {
		return domain.EditableTruckPatch{}, err
	}

	return domain.EditableTruckPatch{
		Make:           (*domain.TruckMake)(truckPatch.Make),
		Model:          (*domain.TruckModel)(truckPatch.Model),
		LicensePlate:   (*domain.TruckLicensePlate)(truckPatch.LicensePlate),
		PersonCapacity: (*domain.TruckPersonCapacity)(truckPatch.PersonCapacity),
		GeoJSON:        geoJSON,
	}, nil
}

// listTrucksParamsToDomain returns a domain trucks paginated filter based on the standardized list trucks parameters.
func listTrucksParamsToDomain(params spec.ListTrucksParams) domain.TrucksPaginatedFilter {
	domainSort := domain.TruckPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListTrucksParamsSortMake:
			domainSort = domain.TruckPaginatedSortMake
		case spec.ListTrucksParamsSortModel:
			domainSort = domain.TruckPaginatedSortModel
		case spec.ListTrucksParamsSortLicensePlate:
			domainSort = domain.TruckPaginatedSortLicensePlate
		case spec.ListTrucksParamsSortPersonCapacity:
			domainSort = domain.TruckPaginatedSortPersonCapacity
		case spec.ListTrucksParamsSortWayName:
			domainSort = domain.TruckPaginatedSortWayName
		case spec.ListTrucksParamsSortMunicipalityName:
			domainSort = domain.TruckPaginatedSortMunicipalityName
		case spec.ListTrucksParamsSortCreatedAt:
			domainSort = domain.TruckPaginatedSortCreatedAt
		case spec.ListTrucksParamsSortModifiedAt:
			domainSort = domain.TruckPaginatedSortModifiedAt
		default:
			domainSort = domain.TruckPaginatedSort(*params.Sort)
		}
	}

	return domain.TrucksPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		Make:         (*domain.TruckMake)(params.Make),
		Model:        (*domain.TruckModel)(params.Model),
		LicensePlate: (*domain.TruckLicensePlate)(params.LicensePlate),
		LocationName: params.LocationName,
	}
}

// truckFromDomain returns a standardized truck based on the domain model.
func truckFromDomain(truck domain.Truck) (spec.Truck, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(truck.GeoJSON)
	if err != nil {
		return spec.Truck{}, err
	}

	return spec.Truck{
		Id:             truck.ID,
		Make:           string(truck.Make),
		Model:          string(truck.Model),
		LicensePlate:   string(truck.LicensePlate),
		PersonCapacity: int(truck.PersonCapacity),
		GeoJson:        geoJSON,
		CreatedAt:      truck.CreatedAt,
		ModifiedAt:     truck.ModifiedAt,
	}, nil
}

// trucksFromDomain returns standardized trucks based on the domain model.
func trucksFromDomain(trucks []domain.Truck) ([]spec.Truck, error) {
	specTrucks := make([]spec.Truck, len(trucks))
	var err error

	for i, truck := range trucks {
		specTrucks[i], err = truckFromDomain(truck)
		if err != nil {
			return []spec.Truck{}, err
		}
	}

	return specTrucks, nil
}

// trucksPaginatedFromDomain returns a standardized trucks paginated response based on the domain model.
func trucksPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Truck]) (spec.TrucksPaginated, error) {
	trucks, err := trucksFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.TrucksPaginated{}, err
	}

	return spec.TrucksPaginated{
		Total:  paginatedResponse.Total,
		Trucks: trucks,
	}, nil
}
