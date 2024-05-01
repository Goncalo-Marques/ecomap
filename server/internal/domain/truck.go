package domain

import (
	"time"

	"github.com/google/uuid"
)

// Truck constraints.
const (
	truckMakeMaxLength          = 50
	truckModelMaxLength         = 50
	truckLicensePlateMaxLength  = 30
	truckPersonCapacityMinValue = 1
)

// TruckMake defines the make of the truck.
type TruckMake string

// Valid returns true if the make is valid, false otherwise.
func (m TruckMake) Valid() bool {
	return len(m) <= truckMakeMaxLength
}

// TruckModel defines the model of the truck.
type TruckModel string

// Valid returns true if the model is valid, false otherwise.
func (m TruckModel) Valid() bool {
	return len(m) <= truckModelMaxLength
}

// TruckLicensePlate defines the license plate of the truck.
type TruckLicensePlate string

// Valid returns true if the license plate is valid, false otherwise.
func (lp TruckLicensePlate) Valid() bool {
	return len(lp) <= truckLicensePlateMaxLength
}

// TruckPersonCapacity defines the person capacity of the truck.
type TruckPersonCapacity int

// Valid returns true if the person capacity is valid, false otherwise.
func (pc TruckPersonCapacity) Valid() bool {
	return pc >= truckPersonCapacityMinValue
}

// EditableTruck defines the editable truck structure.
type EditableTruck struct {
	Make           TruckMake
	Model          TruckModel
	LicensePlate   TruckLicensePlate
	PersonCapacity TruckPersonCapacity
	GeoJSON        GeoJSON
}

// EditableTruckPatch defines the patchable truck structure.
type EditableTruckPatch struct {
	Make           *TruckMake
	Model          *TruckModel
	LicensePlate   *TruckLicensePlate
	PersonCapacity *TruckPersonCapacity
	GeoJSON        GeoJSON
}

// Truck defines the truck structure.
type Truck struct {
	EditableTruck
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}
