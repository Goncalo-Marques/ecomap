package com.ecomap.ecomap.map

import android.content.Context
import com.ecomap.ecomap.domain.GeoJSONFeaturePoint
import com.google.android.gms.maps.model.LatLng
import com.google.maps.android.clustering.ClusterItem

/**
 * A marker that represents one or more containers at the same position on the map.
 */
class ContainerMarker(
    internal val context: Context,
    internal val id: String,
    internal val geoJSON: GeoJSONFeaturePoint,
    val categories: ArrayList<String>
) : ClusterItem {
    override fun getPosition(): LatLng {
        val coordinates = geoJSON.geometry.coordinates
        return LatLng(coordinates[1], coordinates[0])
    }

    override fun getTitle(): String {
        return geoJSON.properties.getLocationName(context)
    }

    override fun getSnippet(): String {
        return categories.toSet().joinToString(", ")
    }

    override fun getZIndex(): Float {
        return 0f
    }
}
