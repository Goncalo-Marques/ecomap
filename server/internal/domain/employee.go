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

// TODO: Create method to validate employee role.
// EmployeeRole defines the role of the employee.
type EmployeeRole string

const (
	EmployeeRoleWasteOperator EmployeeRole = "waste_operator"
	EmployeeRoleManager       EmployeeRole = "manager"
)

// EditableEmployee defines the editable employee structure.
type EditableEmployee struct {
	Username      Username
	FirstName     Name
	LastName      Name
	Role          EmployeeRole
	DateOfBirth   time.Time
	PhoneNumber   string
	GeoJSON       GeoJSON
	ScheduleStart time.Time
	ScheduleEnd   time.Time
}

// EditableEmployeeWithPassword defines the editable employee structure with a password.
type EditableEmployeeWithPassword struct {
	EditableEmployee
	Password
}

// Employee defines the employee structure.
type Employee struct {
	EditableEmployee
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}
