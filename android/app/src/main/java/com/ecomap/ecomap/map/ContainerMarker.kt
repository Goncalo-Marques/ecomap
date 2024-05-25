package com.ecomap.ecomap.map

import com.google.android.gms.maps.model.LatLng
import com.google.maps.android.clustering.ClusterItem

/**
 * A marker that represents one or more containers at the same position on the map.
 */
class ContainerMarker(
    private val position: LatLng,
    private val locationName: String,
    val categories: ArrayList<String>
) : ClusterItem {
    override fun getPosition(): LatLng {
        return position
    }

    override fun getTitle(): String {
        return locationName
    }

    override fun getSnippet(): String {
        return categories.toSet().joinToString(", ")
    }

    override fun getZIndex(): Float {
        return 0f
    }
}
