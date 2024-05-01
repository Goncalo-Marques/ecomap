package domain

import "errors"

// Road network errors.
var (
	ErrRoadNotFound = errors.New("road not found") // Returned when a road is not found.
)

// Road defines the road structure.
type Road struct {
	ID          int
	OsmID       *int
	OsmName     *string
	OsmMeta     *string
	OsmSourceID *int
	OsmTargetID *int
	Clazz       *int
	Flags       *int
	Source      *int
	Target      *int
	KM          *float64
	KMH         *int
	Cost        *float64
	ReverseCost *float64
	X1          *float64
	Y1          *float64
	X2          *float64
	Y2          *float64
}
