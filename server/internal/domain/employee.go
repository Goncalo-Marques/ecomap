package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Employee errors.
var (
	ErrEmployeeNotFound = errors.New("employee not found") // Returned when an employee is not found.
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
	FirstName     Name
	LastName      Name
	Role          EmployeeRole
	DateOfBirth   time.Time
	PhoneNumber   string
	Geom          string // Defined in the GeoJSON format.
	ScheduleStart time.Time
	ScheduleEnd   time.Time
}

// Employee defines the employee structure.
type Employee struct {
	EditableEmployee
	ID           uuid.UUID
	CreatedTime  time.Time
	ModifiedTime time.Time
}
