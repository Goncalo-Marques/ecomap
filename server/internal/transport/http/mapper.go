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

	timeFormatDateOnly = "2006-01-02"
	timeFormatTimeOnly = "15:04:05"
)

var (
	errGeoJSONFeatureTypeUnexpected  = errors.New("unexpected geojson feature type")
	errGeoJSONGeometryTypeUnexpected = errors.New("unexpected geojson geometry type")
)

// dateStringFromTime returns a standardized date string based on the time model.
func dateStringFromTime(time time.Time) string {
	return time.UTC().Format(timeFormatDateOnly)
}

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{
		Time: time.UTC(),
	}
}

// timeStringFromTime returns a standardized time string based on the time model.
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

// geoJSONFeaturePointToDomain returns a domain GeoJSON based on the standardized GeoJSON feature point model.
func geoJSONFeaturePointToDomain(geoJSON *spec.GeoJSONFeaturePoint) (domain.GeoJSON, error) {
	if geoJSON == nil {
		return nil, nil
	}

	if len(geoJSON.Geometry.Coordinates) != 2 {
		return nil, &domain.ErrFieldValueInvalid{FieldName: domain.FieldGeoJSON}
	}

	return domain.GeoJSONFeature{
		Geometry: domain.GeoJSONGeometryPoint{
			Coordinates: [2]float64(geoJSON.Geometry.Coordinates),
		},
	}, nil
}

// geoJSONFeaturePointFromDomain returns standardized GeoJSON feature point based on the domain GeoJSON model.
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

// logicalOperatorToDomain returns a domain logical operator based on the standardized query parameter model.
func logicalOperatorToDomain(logicalOperator *spec.LogicalOperatorQueryParam) domain.PaginationLogicalOperator {
	if logicalOperator == nil {
		return domain.PaginationLogicalOperatorAnd
	}

	switch *logicalOperator {
	case spec.LogicalOperatorQueryParamAnd:
		return domain.PaginationLogicalOperatorAnd
	case spec.LogicalOperatorQueryParamOr:
		return domain.PaginationLogicalOperatorOr
	default:
		return domain.PaginationLogicalOperator(*logicalOperator)
	}
}

// orderToDomain returns a domain order based on the standardized query parameter model.
func orderToDomain(order *spec.OrderQueryParam) domain.PaginationOrder {
	if order == nil {
		return domain.PaginationOrderAsc
	}

	switch *order {
	case spec.OrderQueryParamAsc:
		return domain.PaginationOrderAsc
	case spec.OrderQueryParamDesc:
		return domain.PaginationOrderDesc
	default:
		return domain.PaginationOrder(*order)
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

// paginatedRequestToDomain returns a domain paginated request based on the standardized query parameter models.
func paginatedRequestToDomain[T any](logicalOperator *spec.LogicalOperatorQueryParam, sort domain.PaginationSort[T], order *spec.OrderQueryParam, limit *spec.LimitQueryParam, offset *spec.OffsetQueryParam) domain.PaginatedRequest[T] {
	return domain.PaginatedRequest[T]{
		LogicalOperator: logicalOperatorToDomain(logicalOperator),
		Sort:            sort,
		Order:           orderToDomain(order),
		Limit:           limitToDomain(limit),
		Offset:          offsetToDomain(offset),
	}
}
