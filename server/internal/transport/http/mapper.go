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

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{
		Time: time.UTC(),
	}
}

// timeFromTime returns a standardized time based on the time model.
func timeFromTime(time time.Time) string {
	return time.UTC().Format(timeFormatTimeOnly)
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

// userPostToDomainEditableUserWithPassword returns a domain editable user with password based on the standardized user
// post.
func userPostToDomainEditableUserWithPassword(userPost spec.UserPost) domain.EditableUserWithPassword {
	return domain.EditableUserWithPassword{
		EditableUser: domain.EditableUser{
			Username:  domain.Username(userPost.Username),
			FirstName: domain.Name(userPost.FirstName),
			LastName:  domain.Name(userPost.LastName),
		},
		Password: domain.Password(userPost.Password),
	}
}

// userPatchToDomainEditableUserPatch returns a domain patchable user based on the standardized user patch.
func userPatchToDomainEditableUserPatch(userPatch spec.UserPatch) domain.EditableUserPatch {
	return domain.EditableUserPatch{
		Username:  (*domain.Username)(userPatch.Username),
		FirstName: (*domain.Name)(userPatch.FirstName),
		LastName:  (*domain.Name)(userPatch.LastName),
	}
}

// listUsersParamsToDomainUsersPaginatedFilter returns a domain users paginated filter based on the standardized list users parameters.
func listUsersParamsToDomainUsersPaginatedFilter(params spec.ListUsersParams) domain.UsersPaginatedFilter {
	var domainSort domain.PaginationSort[domain.UserPaginatedSort]
	if params.Sort != nil {
		domainSort = domain.UserPaginatedSort(*params.Sort)
	}

	return domain.UsersPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			params.Limit,
			params.Offset,
			(*spec.OrderQueryParam)(params.Order),
			domainSort,
		),
		Username:  (*domain.Username)(params.Username),
		FirstName: (*domain.Name)(params.FirstName),
		LastName:  (*domain.Name)(params.LastName),
	}
}

// geoJSONFeatureFromDomain returns standardized GeoJSON feature based on the domain model.
func geoJSONFeatureFromDomain(geoJSONFeature domain.GeoJSONFeature) (spec.GeoJSONFeature, error) {
	specGeoJSONFeature := spec.GeoJSONFeature{
		Type: spec.Feature,
		Properties: spec.GeoJSONFeatureProperties{
			WayOsmId: geoJSONFeature.Properties.WayOSM(),
		},
	}
	var err error

	switch geometry := geoJSONFeature.Geometry.(type) {
	case domain.GeoJSONGeometryPoint:
		err = specGeoJSONFeature.Geometry.FromGeoJSONGeometryPoint(spec.GeoJSONGeometryPoint{
			Type:        spec.Point,
			Coordinates: geometry.Coordinates[:],
		})
		if err != nil {
			return spec.GeoJSONFeature{}, err
		}

	case domain.GeoJSONGeometryLineString:
		coordinates := make([][]float64, len(geometry.Coordinates))
		for i, c := range geometry.Coordinates {
			coordinates[i] = c[:]
		}

		err = specGeoJSONFeature.Geometry.FromGeoJSONGeometryLineString(spec.GeoJSONGeometryLineString{
			Type:        spec.LineString,
			Coordinates: coordinates,
		})
		if err != nil {
			return spec.GeoJSONFeature{}, err
		}

	default:
		return spec.GeoJSONFeature{}, errors.New("unexpected geojson geometry type")
	}

	return specGeoJSONFeature, nil
}

// geoJSONFromDomain returns standardized GeoJSON based on the domain model.
func geoJSONFromDomain(geoJSON domain.GeoJSON) (spec.GeoJSON, error) {
	var specGeoJSON spec.GeoJSON
	var err error

	switch feature := geoJSON.(type) {
	case domain.GeoJSONFeature:
		specFeature, err := geoJSONFeatureFromDomain(feature)
		if err != nil {
			return spec.GeoJSON{}, err
		}

		err = specGeoJSON.FromGeoJSONFeature(specFeature)
		if err != nil {
			return spec.GeoJSON{}, err
		}

	case domain.GeoJSONFeatureCollection:
		specFeatureCollection := spec.GeoJSONFeatureCollection{
			Type:     spec.FeatureCollection,
			Features: make([]spec.GeoJSONFeature, len(feature.Features)),
		}

		for i, f := range feature.Features {
			specFeatureCollection.Features[i], err = geoJSONFeatureFromDomain(f)
			if err != nil {
				return spec.GeoJSON{}, err
			}
		}

		err = specGeoJSON.FromGeoJSONFeatureCollection(specFeatureCollection)
		if err != nil {
			return spec.GeoJSON{}, err
		}

	default:
		return spec.GeoJSON{}, errors.New("unexpected geojson feature type")
	}

	return specGeoJSON, nil
}

// userFromDomain returns a standardized user based on the domain model.
func userFromDomain(user domain.User) spec.User {
	return spec.User{
		Id:         user.ID,
		Username:   string(user.Username),
		FirstName:  string(user.FirstName),
		LastName:   string(user.LastName),
		CreatedAt:  user.CreatedAt,
		ModifiedAt: user.ModifiedAt,
	}
}

// usersFromDomain returns standardized users based on the domain model.
func usersFromDomain(users []domain.User) []spec.User {
	specUsers := make([]spec.User, len(users))
	for i, user := range users {
		specUsers[i] = userFromDomain(user)
	}

	return specUsers
}

// usersPaginatedFromDomain returns a standardized users paginated response based on the domain model.
func usersPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.User]) spec.UsersPaginated {
	return spec.UsersPaginated{
		Total: paginatedResponse.Total,
		Users: usersFromDomain(paginatedResponse.Results),
	}
}

// employeeFromDomain returns a standardized employee based on the domain model.
func employeeFromDomain(employee domain.Employee) (spec.Employee, error) {
	geoJSON, err := geoJSONFromDomain(employee.GeoJSON)
	if err != nil {
		return spec.Employee{}, err
	}

	return spec.Employee{
		Id:            employee.ID,
		Username:      string(employee.Username),
		FirstName:     string(employee.FirstName),
		LastName:      string(employee.LastName),
		Role:          spec.EmployeeRole(employee.Role),
		DateOfBirth:   dateFromTime(employee.DateOfBirth),
		PhoneNumber:   employee.PhoneNumber,
		GeoJSON:       geoJSON,
		ScheduleStart: timeFromTime(employee.ScheduleStart),
		ScheduleEnd:   timeFromTime(employee.ScheduleEnd),
		CreatedAt:     employee.CreatedAt,
		ModifiedAt:    employee.ModifiedAt,
	}, nil
}
