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
	errContainerNotFound = "container not found"
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
	w.WriteHeader(http.StatusNotFound)
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
	w.WriteHeader(http.StatusNotFound)
}

// DeleteContainerByID handles the http request to delete a container by ID.
func (h *handler) DeleteContainerByID(w http.ResponseWriter, r *http.Request, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
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
	if len(containerPost.GeoJson.Geometry.Coordinates) != 2 {
		return domain.EditableContainer{}, &domain.ErrFieldValueInvalid{FieldName: domain.FieldGeoJSON}
	}

	return domain.EditableContainer{
		Category: containerCategoryToDomain(containerPost.Category),
		GeoJSON: domain.GeoJSONFeature{
			Geometry: domain.GeoJSONGeometryPoint{
				Coordinates: [2]float64(containerPost.GeoJson.Geometry.Coordinates),
			},
		},
	}, nil
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
