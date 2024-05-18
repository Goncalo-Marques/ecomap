package com.ecomap.ecomap.domain

/**
 * Represents the structure of GeoJSON geometry point.
 * @param coordinates GeoJSON geometry coordinates.
 */
data class GeoJSONGeometryPoint(
    var coordinates: DoubleArray = doubleArrayOf(0.0, 0.0), // [longitude, latitude]
)

/**
 * Represents the structure of GeoJSON feature properties.
 * @param wayName Way name of the feature.
 * @param municipalityName Name of the municipality where the feature is located.
 */
data class GeoJSONProperties(
    var wayName: String = "",
    var municipalityName: String = "",
)

/**
 * Represents the structure of a GeoJSON feature point.
 * @param geometry GeoJSON geometry.
 * @param properties GeoJSON properties.
 */
data class GeoJSONFeaturePoint(
    var geometry: GeoJSONGeometryPoint = GeoJSONGeometryPoint(),
    var properties: GeoJSONProperties = GeoJSONProperties(),
)
