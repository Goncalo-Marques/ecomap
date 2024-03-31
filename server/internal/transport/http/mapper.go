package http

import (
	"time"

	oapitypes "github.com/oapi-codegen/runtime/types"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// jwtFromJWTToken returns a standardized JWT based on the JWT token.
func jwtFromJWTToken(token string) spec.JWT {
	return spec.JWT{
		Token: token,
	}
}

// employeeFromDomain returns a standardized employee based on the domain model.
func employeeFromDomain(employee domain.Employee) spec.Employee {
	return spec.Employee{
		Id:          employee.ID,
		Name:        employee.FirstName,
		DateOfBirth: dateFromTime(employee.DateOfBirth),
	}
}

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{Time: time}
}
