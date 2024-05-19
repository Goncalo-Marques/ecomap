package com.ecomap.ecomap.domain

import android.content.Context
import com.ecomap.ecomap.R

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
) {
    /**
     * Returns the location name based on the way and municipality name.
     * @return Location name.
     */
    fun getLocationName(context: Context): String {
        var locationName = context.getString(R.string.way_unknown)

        if (wayName.isNotBlank()) {
            locationName = wayName
        }
        if (municipalityName.isNotBlank()) {
            locationName += ", $municipalityName"
        }

        return locationName
    }
}

/**
 * Represents the structure of a GeoJSON feature point.
 * @param geometry GeoJSON geometry.
 * @param properties GeoJSON properties.
 */
data class GeoJSONFeaturePoint(
    var geometry: GeoJSONGeometryPoint = GeoJSONGeometryPoint(),
    var properties: GeoJSONProperties = GeoJSONProperties(),
)
