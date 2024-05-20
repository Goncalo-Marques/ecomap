package com.ecomap.ecomap.map

import android.content.Context
import com.ecomap.ecomap.R
import com.google.android.gms.maps.GoogleMap
import com.google.android.gms.maps.model.BitmapDescriptorFactory
import com.google.android.gms.maps.model.MarkerOptions
import com.google.maps.android.clustering.ClusterManager
import com.google.maps.android.clustering.view.DefaultClusterRenderer

/**
 * Custom cluster renderer to display the respective container marker icon.
 */
class ContainerClusterRenderer(
    context: Context,
    map: GoogleMap,
    clusterManager: ClusterManager<ContainerMarker>
) : DefaultClusterRenderer<ContainerMarker>(context, map, clusterManager) {
    private val bitmapMarker = BitmapDescriptorFactory.fromResource(R.drawable.marker_icon)
    private val color = context.getColor(R.color.green_800)

    override fun onBeforeClusterItemRendered(item: ContainerMarker, markerOptions: MarkerOptions) {
        super.onBeforeClusterItemRendered(item, markerOptions)
        markerOptions.icon(bitmapMarker)
    }

    override fun getColor(clusterSize: Int): Int {
        return color
    }
}
