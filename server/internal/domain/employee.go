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

// EditableEmployee defines the editable employee structure.
type EditableEmployee struct {
	Name        string
	DateOfBirth time.Time
}

// Employee defines the employee structure.
type Employee struct {
	EditableEmployee
	ID uuid.UUID
}
