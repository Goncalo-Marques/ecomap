package http

import (
	"time"

	oapitypes "github.com/oapi-codegen/runtime/types"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	paginationLimitDefaultValue  = 100
	paginationOffsetDefaultValue = 0
	orderDefaultValue            = spec.OrderQueryParamAsc
)

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{
		Time: time,
	}
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
func orderToDomain(order *spec.OrderQueryParam) domain.Order {
	if order == nil {
		return domain.Order(orderDefaultValue)
	}

	return domain.Order(*order)
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

// userFromDomain returns a standardized user based on the domain model.
func userFromDomain(user domain.User) spec.User {
	return spec.User{
		Id:           user.ID,
		Username:     string(user.Username),
		FirstName:    string(user.FirstName),
		LastName:     string(user.LastName),
		CreatedTime:  user.CreatedTime,
		ModifiedTime: user.ModifiedTime,
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

// employeeFromDomain returns a standardized employee based on the domain model.
func employeeFromDomain(employee domain.Employee) spec.Employee {
	return spec.Employee{
		Id:          employee.ID,
		Name:        string(employee.FirstName),
		DateOfBirth: dateFromTime(employee.DateOfBirth),
	}
}
