package http

import (
	"time"

	oapitypes "github.com/oapi-codegen/runtime/types"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// fromDomainEmployee returns a standardized employee based on the domain model.
func fromDomainEmployee(employee domain.Employee) spec.Employee {
	return spec.Employee{
		Id:          employee.ID,
		Name:        employee.Name,
		DateOfBirth: dateFromTime(employee.DateOfBirth),
	}
}

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{Time: time}
}
