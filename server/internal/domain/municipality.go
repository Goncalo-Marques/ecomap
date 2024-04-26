package domain

import "errors"

// Municipality errors.
var (
	ErrMunicipalityNotFound = errors.New("municipality not found") // Returned when a municipality is not found.
)

// Municipality defines the municipality structure.
type Municipality struct {
	ID        int
	FeatureID int
	Name      string
	District  string
	NUTS1     string
	NUTS2     string
	NUTS3     string
	Area      float64 // Area of the municipality in the metric unit hectare.
	Perimeter float64 // Perimeter of the municipality in kilometers.
}
