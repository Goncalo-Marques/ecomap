package domain

import (
	"errors"
	"fmt"
)

// Common errors.
var (
	ErrCredentialsIncorrect = errors.New("incorrect credentials") // Returned when a username is not found or the password is incorrect.
)

// Returned when a field contains an invalid value.
type ErrFieldValueInvalid struct {
	FieldName string
}

func (e *ErrFieldValueInvalid) Error() string {
	return fmt.Sprintf("invalid field value: %s", e.FieldName)
}

// Returned when a filter contains an invalid value.
type ErrFilterValueInvalid struct {
	FilterName string
}

func (e *ErrFilterValueInvalid) Error() string {
	return fmt.Sprintf("invalid filter value: %s", e.FilterName)
}

// Field constraints.
const (
	usernameMinLength = 1
	usernameMaxLength = 50

	nameMinLength = 1
	nameMaxLength = 50

	paginationLimitMinValue = 1
	paginationLimitMaxValue = 100

	paginationOffsetMinValue = 0
)

// Username defines the username type.
type Username string

// Valid returns true if the username is valid, false otherwise.
func (u Username) Valid() bool {
	return len(u) >= usernameMinLength && len(u) <= usernameMaxLength
}

// Name defines the name type.
type Name string

// Valid returns true if the name is valid, false otherwise.
func (n Name) Valid() bool {
	return len(n) >= nameMinLength && len(n) <= nameMaxLength
}

// Password defines the password type.
type Password string

// SignIn defines the sign-in structure.
type SignIn struct {
	Username Username
	Password Password
}

// PaginationLimit defines the pagination limit type.
type PaginationLimit int

// Valid returns true if the pagination limit is valid, false otherwise.
func (l PaginationLimit) Valid() bool {
	return l >= paginationLimitMinValue && l <= paginationLimitMaxValue
}

// PaginationOffset defines the pagination offset type.
type PaginationOffset int

// Valid returns true if the pagination offset is valid, false otherwise.
func (o PaginationOffset) Valid() bool {
	return o >= paginationOffsetMinValue
}

// Order defines the order to sort by.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Valid returns true if the order is valid, false otherwise.
func (o Order) Valid() bool {
	switch o {
	case OrderAsc, OrderDesc:
		return true
	default:
		return false
	}
}

// Sort defines the sort interface.
type Sort[T any] interface {
	// Field returns the name of the field to sort by.
	Field() T
	// Valid returns true if the field is valid, false otherwise.
	Valid() bool
}

// PaginatedRequest defines the paginated request structure.
type PaginatedRequest[SortField any] struct {
	Limit  PaginationLimit  // Amount of resources to get for the provided filter.
	Offset PaginationOffset // Amount of resources to skip for the provided filter.
	Order  Order            // Order to sort by.
	Sort   Sort[SortField]  // Field to sort by.
}

// PaginatedResponse defines the paginated response structure.
type PaginatedResponse[T any] struct {
	Total   int
	Results []T
}
