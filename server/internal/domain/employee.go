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
	EmployeeRoleWasteOperator EmployeeRole = "waste_operator"
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
