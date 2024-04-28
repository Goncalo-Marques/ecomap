package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Employee errors.
var (
	ErrEmployeeAlreadyExists = errors.New("username already exists") // Returned when an employee already exists with the same username.
	ErrEmployeeNotFound      = errors.New("employee not found")      // Returned when an employee is not found.
)

// EmployeeRole defines the role of the employee.
type EmployeeRole string

const (
	EmployeeRoleWasteOperator EmployeeRole = "wasteOperator"
	EmployeeRoleManager       EmployeeRole = "manager"
)

// Valid returns true if the role is valid, false otherwise.
func (r EmployeeRole) Valid() bool {
	switch r {
	case EmployeeRoleWasteOperator,
		EmployeeRoleManager:
		return true
	default:
		return false
	}
}

// EditableEmployee defines the editable employee structure.
type EditableEmployee struct {
	Username      Username
	FirstName     Name
	LastName      Name
	Role          EmployeeRole
	DateOfBirth   time.Time
	PhoneNumber   PhoneNumber
	GeoJSON       GeoJSON
	ScheduleStart time.Time
	ScheduleEnd   time.Time
}

// EditableEmployeeWithPassword defines the editable employee structure with a password.
type EditableEmployeeWithPassword struct {
	EditableEmployee
	Password
}

// EditableEmployeePatch defines the patchable employee structure.
type EditableEmployeePatch struct {
	Username      *Username
	FirstName     *Name
	LastName      *Name
	DateOfBirth   *time.Time
	PhoneNumber   *PhoneNumber
	GeoJSON       GeoJSON
	ScheduleStart *time.Time
	ScheduleEnd   *time.Time
}

// Employee defines the employee structure.
type Employee struct {
	EditableEmployee
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// EmployeePaginatedSort defines the field of the employee to sort.
type EmployeePaginatedSort string

const (
	EmployeePaginatedSortUsername         EmployeePaginatedSort = "username"
	EmployeePaginatedSortFirstName        EmployeePaginatedSort = "firstName"
	EmployeePaginatedSortLastName         EmployeePaginatedSort = "lastName"
	EmployeePaginatedSortRole             EmployeePaginatedSort = "role"
	EmployeePaginatedSortDateOfBirth      EmployeePaginatedSort = "dateOfBirth"
	EmployeePaginatedSortScheduleStart    EmployeePaginatedSort = "scheduleStart"
	EmployeePaginatedSortScheduleEnd      EmployeePaginatedSort = "scheduleEnd"
	EmployeePaginatedSortWayName          EmployeePaginatedSort = "wayName"
	EmployeePaginatedSortMunicipalityName EmployeePaginatedSort = "municipalityName"
	EmployeePaginatedSortCreatedAt        EmployeePaginatedSort = "createdAt"
	EmployeePaginatedSortModifiedAt       EmployeePaginatedSort = "modifiedAt"
)

// Field returns the name of the field to sort by.
func (s EmployeePaginatedSort) Field() EmployeePaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s EmployeePaginatedSort) Valid() bool {
	switch s {
	case EmployeePaginatedSortUsername,
		EmployeePaginatedSortFirstName,
		EmployeePaginatedSortLastName,
		EmployeePaginatedSortRole,
		EmployeePaginatedSortDateOfBirth,
		EmployeePaginatedSortScheduleStart,
		EmployeePaginatedSortScheduleEnd,
		EmployeePaginatedSortWayName,
		EmployeePaginatedSortMunicipalityName,
		EmployeePaginatedSortCreatedAt,
		EmployeePaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// EmployeesPaginatedFilter defines the employees filter structure.
type EmployeesPaginatedFilter struct {
	PaginatedRequest[EmployeePaginatedSort]
	Username         *Username
	FirstName        *Name
	LastName         *Name
	Role             *EmployeeRole
	DateOfBirth      *string
	PhoneNumber      *PhoneNumber
	ScheduleStart    *string
	ScheduleEnd      *string
	WayName          *string
	MunicipalityName *string
}
