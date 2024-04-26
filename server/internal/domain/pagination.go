package domain

import "fmt"

// Pagination constraints.
const (
	paginationLimitMinValue = 1
	paginationLimitMaxValue = 100

	paginationOffsetMinValue = 0
)

// Returned when a filter contains an invalid value.
type ErrFilterValueInvalid struct {
	FilterName string
}

func (e *ErrFilterValueInvalid) Error() string {
	return fmt.Sprintf("invalid filter value: %s", e.FilterName)
}

// PaginationLogicalOperator defines the pagination logical operator type.
type PaginationLogicalOperator string

const (
	PaginationLogicalOperatorAnd PaginationLogicalOperator = "And"
	PaginationLogicalOperatorOr  PaginationLogicalOperator = "Or"
)

// Valid returns true if the pagination logical operator is valid, false otherwise.
func (lo PaginationLogicalOperator) Valid() bool {
	switch lo {
	case PaginationLogicalOperatorAnd, PaginationLogicalOperatorOr:
		return true
	default:
		return false
	}
}

// PaginationSort defines the pagination sort interface.
type PaginationSort[T any] interface {
	// Field returns the name of the field to sort by.
	Field() T
	// Valid returns true if the field is valid, false otherwise.
	Valid() bool
}

// PaginationOrder defines the pagination order to sort by.
type PaginationOrder string

const (
	PaginationOrderAsc  PaginationOrder = "Asc"
	PaginationOrderDesc PaginationOrder = "Desc"
)

// Valid returns true if the pagination order is valid, false otherwise.
func (o PaginationOrder) Valid() bool {
	switch o {
	case PaginationOrderAsc, PaginationOrderDesc:
		return true
	default:
		return false
	}
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

// PaginatedRequest defines the paginated request structure.
type PaginatedRequest[SortField any] struct {
	LogicalOperator PaginationLogicalOperator // Logical operator used for the provided filter.
	Sort            PaginationSort[SortField] // Field to sort by.
	Order           PaginationOrder           // Order to sort by.
	Limit           PaginationLimit           // Amount of resources to get for the provided filter.
	Offset          PaginationOffset          // Amount of resources to skip for the provided filter.
}

// PaginatedResponse defines the paginated response structure.
type PaginatedResponse[T any] struct {
	Total   int
	Results []T
}
