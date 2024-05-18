package com.ecomap.ecomap.domain

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
    HAZARDOUS
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
