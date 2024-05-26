package com.ecomap.ecomap.domain

import android.content.Context
import com.ecomap.ecomap.R

/**
 * Represents all the available container categories.
 */
enum class ContainerCategory {
    GENERAL,
    PAPER,
    PLASTIC,
    METAL,
    GLASS,
    ORGANIC,
    HAZARDOUS;

    /**
     * Returns the respective drawable icon resource.
     * @return Resource ID.
     */
    fun getIconResource(): Int {
        return when (this) {
            GENERAL -> R.drawable.general
            PAPER -> R.drawable.paper
            PLASTIC -> R.drawable.plastic
            METAL -> R.drawable.metal
            GLASS -> R.drawable.glass
            ORGANIC -> R.drawable.organic
            HAZARDOUS -> R.drawable.hazardous
        }
    }

    /**
     * Returns the respective human readable string resource.
     * @param context Activity context.
     * @return String resource.
     */
    fun getStringResource(context: Context): String {
        return when (this) {
            GENERAL -> context.getString(R.string.container_category_general)
            PAPER -> context.getString(R.string.container_category_paper)
            PLASTIC -> context.getString(R.string.container_category_plastic)
            METAL -> context.getString(R.string.container_category_metal)
            GLASS -> context.getString(R.string.container_category_glass)
            ORGANIC -> context.getString(R.string.container_category_organic)
            HAZARDOUS -> context.getString(R.string.container_category_hazardous)
        }
    }
}

/**
 * Represents the structure of a container.
 * @param id Container id.
 * @param category Container category.
 * @param geoJSON Container location in the GeoJSON format.
 * @param createdAt Timestamp of when the container was created.
 * @param modifiedAt Timestamp of when the container was modified.
 */
data class Container(
    val id: String,
    val category: ContainerCategory,
    val geoJSON: GeoJSONFeaturePoint,
    val createdAt: String,
    val modifiedAt: String,
)

/**
 * Represents the structure of paginated containers.
 * @param total The total amount of resources available for the provided filter.
 * @param containers The list of containers for the provided filter.
 */
data class ContainersPaginated(
    val total: Int,
    val containers: ArrayList<Container>
)
