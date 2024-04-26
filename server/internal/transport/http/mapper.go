package http

import (
	"errors"
	"time"

	oapitypes "github.com/oapi-codegen/runtime/types"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	paginationLimitDefaultValue  = 100
	paginationOffsetDefaultValue = 0
	orderDefaultValue            = spec.OrderQueryParamAsc

	timeFormatTimeOnly = "15:04:05"
)

var (
	errGeoJSONFeatureTypeUnexpected  = errors.New("unexpected geojson feature type")
	errGeoJSONGeometryTypeUnexpected = errors.New("unexpected geojson geometry type")
)

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{
		Time: time.UTC(),
	}
}

// timeStringFromTime returns a standardized time based on the time model.
func timeStringFromTime(time time.Time) string {
	return time.UTC().Format(timeFormatTimeOnly)
}

// timeStringToTime returns a time structure based on the standardized time model.
func timeStringToTime(timeString string) (time.Time, error) {
	return time.Parse(timeFormatTimeOnly, timeString)
}

// jwtFromJWTToken returns a standardized JWT based on the JWT token.
func jwtFromJWTToken(token string) spec.JWT {
	return spec.JWT{
		Token: token,
	}
}

// limitToDomain returns a domain pagination limit based on the standardized query parameter model.
func limitToDomain(limit *spec.LimitQueryParam) domain.PaginationLimit {
	if limit == nil {
		return domain.PaginationLimit(paginationLimitDefaultValue)
	}

	return domain.PaginationLimit(*limit)
}

// offsetToDomain returns a domain pagination offset based on the standardized query parameter model.
func offsetToDomain(offset *spec.OffsetQueryParam) domain.PaginationOffset {
	if offset == nil {
		return domain.PaginationOffset(paginationOffsetDefaultValue)
	}

	return domain.PaginationOffset(*offset)
}

// orderToDomain returns a domain order based on the standardized query parameter model.
func orderToDomain(order *spec.OrderQueryParam) domain.PaginationOrder {
	if order == nil {
		return domain.PaginationOrder(orderDefaultValue)
	}

	return domain.PaginationOrder(*order)
}

// paginatedRequestToDomain returns a domain paginated request based on the standardized query parameter models.
func paginatedRequestToDomain[T any](limit *spec.LimitQueryParam, offset *spec.OffsetQueryParam, order *spec.OrderQueryParam, sort domain.PaginationSort[T]) domain.PaginatedRequest[T] {
	return domain.PaginatedRequest[T]{
		Limit:  limitToDomain(limit),
		Offset: offsetToDomain(offset),
		Order:  orderToDomain(order),
		Sort:   sort,
	}
}

// geoJSONFeaturePointFromDomain returns standardized GeoJSON feature point based on the domain model.
func geoJSONFeaturePointFromDomain(geoJSON domain.GeoJSON) (spec.GeoJSONFeaturePoint, error) {
	geoJSONFeature, ok := geoJSON.(domain.GeoJSONFeature)
	if !ok {
		return spec.GeoJSONFeaturePoint{}, errGeoJSONFeatureTypeUnexpected
	}

	geoJSONGeometry, ok := geoJSONFeature.Geometry.(domain.GeoJSONGeometryPoint)
	if !ok {
		return spec.GeoJSONFeaturePoint{}, errGeoJSONGeometryTypeUnexpected
	}

	return spec.GeoJSONFeaturePoint{
		Type: spec.Feature,
		Geometry: spec.GeoJSONGeometryPoint{
			Type:        spec.Point,
			Coordinates: geoJSONGeometry.Coordinates[:],
		},
		Properties: spec.GeoJSONFeatureProperties{
			WayName:          geoJSONFeature.Properties.WayName(),
			MunicipalityName: geoJSONFeature.Properties.MunicipalityName(),
		},
	}, nil
}
