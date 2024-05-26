package com.ecomap.ecomap.map

import android.content.Context
import com.ecomap.ecomap.domain.Container
import com.ecomap.ecomap.domain.GeoJSONFeaturePoint
import com.google.android.gms.maps.model.LatLng
import com.google.maps.android.clustering.ClusterItem

/**
 * A marker that represents one or more containers at the same position on the map.
 */
class ContainerMarker(
    internal val context: Context,
    internal val containers: ArrayList<Container>
) : ClusterItem {
    internal var geoJSON = GeoJSONFeaturePoint()

    init {
        if (containers.isNotEmpty()) {
            geoJSON = containers[0].geoJSON
        }
    }

    override fun getPosition(): LatLng {
        val coordinates = geoJSON.geometry.coordinates
        return LatLng(coordinates[1], coordinates[0])
    }

    override fun getTitle(): String {
        return geoJSON.properties.getLocationName(context)
    }

    override fun getSnippet(): String {
        val categoriesString = ArrayList<String>(containers.size)
        for (container in containers) {
            categoriesString.add(container.category.getStringResource(context))
        }

        // Convert to set to avoid duplicate categories.
        return categoriesString.toSet().joinToString(", ")
    }

    override fun getZIndex(): Float {
        return 0f
    }
}
