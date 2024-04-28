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
	errContainerNotFound                            = "container not found"
	errContainerAssociatedWithContainerReport       = "container associated with user report"
	errContainerAssociatedWithUserContainerBookmark = "container associated with user bookmark"
	errContainerAssociatedWithRouteContainer        = "container associated with route"
)

// CreateContainer handles the http request to create a container.
func (h *handler) CreateContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var containerPost spec.ContainerPost
	err = json.Unmarshal(requestBody, &containerPost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableContainer, err := containerPostToDomain(containerPost)
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

	domainContainer, err := h.service.CreateContainer(ctx, domainEditableContainer)
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

	container, err := containerFromDomain(domainContainer)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(container)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListContainers handles the http request to list containers.
func (h *handler) ListContainers(w http.ResponseWriter, r *http.Request, params spec.ListContainersParams) {
	ctx := r.Context()

	domainContainersFilter := listContainersParamsToDomain(params)
	domainPaginatedContainers, err := h.service.ListContainers(ctx, domainContainersFilter)
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

	containersPaginated, err := containersPaginatedFromDomain(domainPaginatedContainers)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(containersPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetContainerByID handles the http request to get a container by ID.
func (h *handler) GetContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	domainContainer, err := h.service.GetContainerByID(ctx, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrContainerNotFound):
			notFound(w, errContainerNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	container, err := containerFromDomain(domainContainer)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(container)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchContainerByID handles the http request to modify a container by ID.
func (h *handler) PatchContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var containerPatch spec.ContainerPatch
	err = json.Unmarshal(requestBody, &containerPatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableContainer, err := containerPatchToDomain(containerPatch)
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

	domainContainer, err := h.service.PatchContainer(ctx, containerID, domainEditableContainer)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrContainerNotFound):
			notFound(w, errContainerNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	container, err := containerFromDomain(domainContainer)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(container)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteContainerByID handles the http request to delete a container by ID.
func (h *handler) DeleteContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	domainContainer, err := h.service.DeleteContainerByID(ctx, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrContainerNotFound):
			notFound(w, errContainerNotFound)
		case errors.Is(err, domain.ErrContainerAssociatedWithContainerReport):
			conflict(w, errContainerAssociatedWithContainerReport)
		case errors.Is(err, domain.ErrContainerAssociatedWithUserContainerBookmark):
			conflict(w, errContainerAssociatedWithUserContainerBookmark)
		case errors.Is(err, domain.ErrContainerAssociatedWithRouteContainer):
			conflict(w, errContainerAssociatedWithRouteContainer)
		default:
			internalServerError(w)
		}

		return
	}

	container, err := containerFromDomain(domainContainer)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(container)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// containerCategoryToDomain returns a domain container category based on the standardized model.
func containerCategoryToDomain(category spec.ContainerCategory) domain.ContainerCategory {
	switch category {
	case spec.General:
		return domain.ContainerCategoryGeneral
	case spec.Paper:
		return domain.ContainerCategoryPaper
	case spec.Plastic:
		return domain.ContainerCategoryPlastic
	case spec.Metal:
		return domain.ContainerCategoryMetal
	case spec.Glass:
		return domain.ContainerCategoryGlass
	case spec.Organic:
		return domain.ContainerCategoryOrganic
	case spec.Hazardous:
		return domain.ContainerCategoryHazardous
	default:
		return domain.ContainerCategory(category)
	}
}

// containerCategoryFromDomain returns a standardized container category based on the domain model.
func containerCategoryFromDomain(category domain.ContainerCategory) spec.ContainerCategory {
	switch category {
	case domain.ContainerCategoryGeneral:
		return spec.General
	case domain.ContainerCategoryPaper:
		return spec.Paper
	case domain.ContainerCategoryPlastic:
		return spec.Plastic
	case domain.ContainerCategoryMetal:
		return spec.Metal
	case domain.ContainerCategoryGlass:
		return spec.Glass
	case domain.ContainerCategoryOrganic:
		return spec.Organic
	case domain.ContainerCategoryHazardous:
		return spec.Hazardous
	default:
		return spec.ContainerCategory(category)
	}
}

// containerPostToDomain returns a domain editable container based on the standardized container post.
func containerPostToDomain(containerPost spec.ContainerPost) (domain.EditableContainer, error) {
	geoJSON, err := geoJSONFeaturePointToDomain(&containerPost.GeoJson)
	if err != nil {
		return domain.EditableContainer{}, err
	}

	return domain.EditableContainer{
		Category: containerCategoryToDomain(containerPost.Category),
		GeoJSON:  geoJSON,
	}, nil
}

// containerPatchToDomain returns a domain patchable container based on the standardized container patch.
func containerPatchToDomain(containerPatch spec.ContainerPatch) (domain.EditableContainerPatch, error) {
	var category *domain.ContainerCategory
	if containerPatch.Category != nil {
		c := containerCategoryToDomain(*containerPatch.Category)
		category = &c
	}

	geoJSON, err := geoJSONFeaturePointToDomain(containerPatch.GeoJson)
	if err != nil {
		return domain.EditableContainerPatch{}, err
	}

	return domain.EditableContainerPatch{
		Category: category,
		GeoJSON:  geoJSON,
	}, nil
}

// listContainersParamsToDomain returns a domain containers paginated filter based on the standardized list containers parameters.
func listContainersParamsToDomain(params spec.ListContainersParams) domain.ContainersPaginatedFilter {
	domainSort := domain.ContainerPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListContainersParamsSortCategory:
			domainSort = domain.ContainerPaginatedSortCategory
		case spec.ListContainersParamsSortWayName:
			domainSort = domain.ContainerPaginatedSortWayName
		case spec.ListContainersParamsSortMunicipalityName:
			domainSort = domain.ContainerPaginatedSortMunicipalityName
		case spec.ListContainersParamsSortCreatedAt:
			domainSort = domain.ContainerPaginatedSortCreatedAt
		case spec.ListContainersParamsSortModifiedAt:
			domainSort = domain.ContainerPaginatedSortModifiedAt
		default:
			domainSort = domain.ContainerPaginatedSort(*params.Sort)
		}
	}

	var domainCategory *domain.ContainerCategory
	if params.Category != nil {
		category := containerCategoryToDomain(*params.Category)
		domainCategory = &category
	}

	return domain.ContainersPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			(*spec.LogicalOperatorQueryParam)(params.LogicalOperator),
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		Category:         domainCategory,
		WayName:          params.WayName,
		MunicipalityName: params.MunicipalityName,
	}
}

// containerFromDomain returns a standardized container based on the domain model.
func containerFromDomain(container domain.Container) (spec.Container, error) {
	geoJSON, err := geoJSONFeaturePointFromDomain(container.GeoJSON)
	if err != nil {
		return spec.Container{}, err
	}

	return spec.Container{
		Id:         container.ID,
		Category:   containerCategoryFromDomain(container.Category),
		GeoJson:    geoJSON,
		CreatedAt:  container.CreatedAt,
		ModifiedAt: container.ModifiedAt,
	}, nil
}

// containersFromDomain returns standardized containers based on the domain model.
func containersFromDomain(containers []domain.Container) ([]spec.Container, error) {
	specContainers := make([]spec.Container, len(containers))
	var err error

	for i, container := range containers {
		specContainers[i], err = containerFromDomain(container)
		if err != nil {
			return []spec.Container{}, err
		}
	}

	return specContainers, nil
}

// containersPaginatedFromDomain returns a standardized containers paginated response based on the domain model.
func containersPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.Container]) (spec.ContainersPaginated, error) {
	containers, err := containersFromDomain(paginatedResponse.Results)
	if err != nil {
		return spec.ContainersPaginated{}, err
	}

	return spec.ContainersPaginated{
		Total:      paginatedResponse.Total,
		Containers: containers,
	}, nil
}
