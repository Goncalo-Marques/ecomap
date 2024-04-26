package domain

import (
	"encoding/json"
)

const (
	geoJSONFeaturePropertyWayName          = "wayName"
	geoJSONFeaturePropertyMunicipalityName = "municipalityName"
)

// GeoJSONGeometry defines the GeoJSON feature interface.
type GeoJSONGeometry interface {
	// GeometryType returns the geometry type name.
	GeometryType() string
}

// GeoJSON defines the GeoJSON interface.
type GeoJSON interface {
	// Type returns the GeoJSON feature type.
	Type() string
}

// GeoJSONGeometryJSON defines the GeoJSON geometry JSON helper structure.
type GeoJSONGeometryJSON struct {
	Type        string `json:"type"`
	Coordinates any    `json:"coordinates"`
}

// GeoJSONGeometryPoint defines the GeoJSON geometry point structure.
type GeoJSONGeometryPoint struct {
	Coordinates [2]float64
}

func (g GeoJSONGeometryPoint) GeometryType() string {
	return "Point"
}

func (g GeoJSONGeometryPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(GeoJSONGeometryJSON{
		Type:        g.GeometryType(),
		Coordinates: g.Coordinates,
	})
}

// GeoJSONGeometryLineString defines the GeoJSON geometry line string structure.
type GeoJSONGeometryLineString struct {
	Coordinates [][2]float64
}

func (g GeoJSONGeometryLineString) GeometryType() string {
	return "LineString"
}

func (g GeoJSONGeometryLineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(GeoJSONGeometryJSON{
		Type:        g.GeometryType(),
		Coordinates: g.Coordinates,
	})
}

// GeoJSONFeatureProperties defines the GeoJSON feature properties.
type GeoJSONFeatureProperties map[string]any

// WayName returns the name of the closest way to the geometry.
func (p GeoJSONFeatureProperties) WayName() *string {
	if p == nil {
		return nil
	}

	value, ok := p[geoJSONFeaturePropertyWayName]
	if !ok {
		return nil
	}

	valueString, ok := value.(string)
	if !ok {
		return nil
	}

	return &valueString
}

// SetWayName sets the name of the closest way to the geometry.
func (p GeoJSONFeatureProperties) SetWayName(name string) {
	if p == nil {
		return
	}

	p[geoJSONFeaturePropertyWayName] = name
}

// MunicipalityName returns the name of the municipality containing the geometry.
func (p GeoJSONFeatureProperties) MunicipalityName() *string {
	if p == nil {
		return nil
	}

	value, ok := p[geoJSONFeaturePropertyMunicipalityName]
	if !ok {
		return nil
	}

	valueString, ok := value.(string)
	if !ok {
		return nil
	}

	return &valueString
}

// SetMunicipalityName sets the name of the municipality containing the geometry.
func (p GeoJSONFeatureProperties) SetMunicipalityName(name string) {
	if p == nil {
		return
	}

	p[geoJSONFeaturePropertyMunicipalityName] = name
}

// GeoJSONFeature defines the GeoJSON feature structure.
type GeoJSONFeature struct {
	Geometry   GeoJSONGeometry
	Properties GeoJSONFeatureProperties
}

func (f GeoJSONFeature) Type() string {
	return "Feature"
}

// GeoJSONFeatureCollection defines the GeoJSON feature collection structure.
type GeoJSONFeatureCollection struct {
	Features []GeoJSONFeature
}

func (f GeoJSONFeatureCollection) Type() string {
	return "FeatureCollection"
}
