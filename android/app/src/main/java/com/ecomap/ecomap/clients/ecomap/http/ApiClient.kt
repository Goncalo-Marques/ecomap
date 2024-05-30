package com.ecomap.ecomap.clients.ecomap.http

import android.util.Log
import com.android.volley.DefaultRetryPolicy
import com.android.volley.Request
import com.android.volley.Response.ErrorListener
import com.android.volley.Response.Listener
import com.android.volley.VolleyError
import com.android.volley.toolbox.JsonObjectRequest
import com.android.volley.toolbox.StringRequest
import com.ecomap.ecomap.domain.Container
import com.ecomap.ecomap.domain.ContainerCategory
import com.ecomap.ecomap.domain.ContainersPaginated
import com.ecomap.ecomap.domain.Error
import com.ecomap.ecomap.domain.GeoJSONFeaturePoint
import com.ecomap.ecomap.domain.GeoJSONProperties
import com.ecomap.ecomap.domain.User
import org.json.JSONException
import org.json.JSONObject

/**
 * EcoMap HTTP API client.
 */
object ApiClient {
    private val LOG_TAG = ApiClient::class.java.simpleName

    // URLs.
    private const val URL_BASE = "https://server-7fzc7ivuwa-ew.a.run.app/api"
    private const val URL_USERS = "$URL_BASE/users"
    private const val URL_USERS_SIGN_IN = "$URL_BASE/users/signin"
    private const val URL_CONTAINERS = "$URL_BASE/containers"
    private const val URL_BOOKMARK_CONTAINERS = "/bookmarks/containers"

    // Authentication field names.
    private const val FIELD_NAME_TOKEN = "token"
    private const val HEADER_KEY_AUTHORIZATION = "Authorization"
    private const val HEADER_VALUE_BEARER_PREFIX = "Bearer "

    // Error field names.
    private const val FIELD_NAME_ERROR_CODE = "code"
    private const val FIELD_NAME_ERROR_MESSAGE = "message"

    // Pagination field names.
    private const val FIELD_NAME_PAGINATION_LIMIT = "limit"
    private const val FIELD_NAME_PAGINATION_OFFSET = "offset"
    private const val FIELD_NAME_PAGINATION_TOTAL = "total"

    // User field names.
    private const val FIELD_NAME_ID = "id"
    private const val FIELD_NAME_USERNAME = "username"
    private const val FIELD_NAME_PASSWORD = "password"
    private const val FIELD_NAME_FIRST_NAME = "firstName"
    private const val FIELD_NAME_LAST_NAME = "lastName"
    private const val FIELD_NAME_CREATED_AT = "createdAt"
    private const val FIELD_NAME_MODIFIED_AT = "modifiedAt"

    // User field names.
    private const val FIELD_CONTAINER_ID = "id"
    private const val FIELD_CONTAINER_CATEGORY = "category"
    private const val FIELD_CONTAINER_GENERAL = "general"
    private const val FIELD_CONTAINER_PAPER = "paper"
    private const val FIELD_CONTAINER_PLASTIC = "plastic"
    private const val FIELD_CONTAINER_METAL = "metal"
    private const val FIELD_CONTAINER_GLASS = "glass"
    private const val FIELD_CONTAINER_ORGANIC = "organic"
    private const val FIELD_CONTAINER_HAZARDOUS = "hazardous"
    private const val FIELD_CONTAINER_GEO_JSON = "geoJson"
    private const val FIELD_CONTAINER_GEO_JSON_GEOMETRY = "geometry"
    private const val FIELD_CONTAINER_GEO_JSON_GEOMETRY_COORDINATES = "coordinates"
    private const val FIELD_CONTAINER_GEO_JSON_PROPERTIES = "properties"
    private const val FIELD_CONTAINER_GEO_JSON_PROPERTIES_WAY_NAME = "wayName"
    private const val FIELD_CONTAINER_GEO_JSON_PROPERTIES_MUNICIPALITY_NAME = "municipalityName"
    private const val FIELD_CONTAINER_CREATED_AT = "createdAt"
    private const val FIELD_CONTAINER_MODIFIED_AT = "modifiedAt"
    private const val FIELD_CONTAINERS = "containers"

    /**
     * Signs in a user with a given username and password.
     * @param username User username.
     * @param password User password.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun signIn(
        username: String,
        password: String,
        listener: Listener<String>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        val requestPayload = JSONObject()
        requestPayload.put(FIELD_NAME_USERNAME, username)
        requestPayload.put(FIELD_NAME_PASSWORD, password)

        return JsonObjectRequest(
            Request.Method.POST, URL_USERS_SIGN_IN, requestPayload,
            { response -> listener.onResponse(response.optString(FIELD_NAME_TOKEN)) },
            errorListener
        ).apply {
            // Avoid retrying a request that failed due to invalid credentials, as this is an
            // expected error.
            retryPolicy = DefaultRetryPolicy(0, 0, 0f)
        }
    }

    /**
     * Creates a user account.
     * @param firstName User first name.
     * @param lastName User last name.
     * @param username User username.
     * @param password User password.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun createAccount(
        firstName: String,
        lastName: String,
        username: String,
        password: String,
        listener: Listener<User>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        val requestPayload = JSONObject()

        requestPayload.put(FIELD_NAME_FIRST_NAME, firstName)
        requestPayload.put(FIELD_NAME_LAST_NAME, lastName)
        requestPayload.put(FIELD_NAME_USERNAME, username)
        requestPayload.put(FIELD_NAME_PASSWORD, password)

        return JsonObjectRequest(
            Request.Method.POST, URL_USERS, requestPayload,
            { response -> listener.onResponse(mapUser(response)) },
            errorListener
        )
    }

    /**
     * Returns the details of a user account.
     * @param userID User identifier.
     * @param token JWT authorization token.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun getAccount(
        userID: String,
        token: String,
        listener: Listener<User>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        val url = "$URL_USERS/$userID"

        return object : JsonObjectRequest(
            Method.GET, url, null,
            { response -> listener.onResponse(mapUser(response)) },
            errorListener
        ) {
            override fun getHeaders(): MutableMap<String, String> {
                return getHeaders(token)
            }
        }
    }

    /**
     * Returns the containers with the specified filter.
     * @param containerCategory Container category to filter by.
     * @param limit Amount of resources to get for the provided filter.
     * @param offset Amount of resources to skip for the provided filter.
     * @param token JWT authorization token.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun listContainers(
        containerCategory: ContainerCategory? = null,
        limit: Int,
        offset: Int,
        token: String,
        listener: Listener<ContainersPaginated>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        var url = URL_CONTAINERS +
                "?$FIELD_NAME_PAGINATION_LIMIT=$limit" +
                "&$FIELD_NAME_PAGINATION_OFFSET=$offset"
        if (containerCategory != null) {
            url += "&$FIELD_CONTAINER_CATEGORY=${mapDomainContainerCategory(containerCategory)}"
        }

        return object : JsonObjectRequest(
            Method.GET, url, null,
            { response ->
                val containers = ArrayList<Container>(limit)

                val jsonContainers = response.optJSONArray(FIELD_CONTAINERS)
                if (jsonContainers != null) {
                    for (i in 0 until jsonContainers.length()) {
                        containers.add(mapContainer(jsonContainers.optJSONObject(i)))
                    }
                }

                listener.onResponse(
                    ContainersPaginated(
                        response.optInt(FIELD_NAME_PAGINATION_TOTAL),
                        containers
                    )
                )
            },
            errorListener,
        ) {
            override fun getHeaders(): MutableMap<String, String> {
                return getHeaders(token)
            }
        }
    }

    /**
     * Returns the user container bookmarks with the specified filter. The bookmarks are sorted by
     * descending order of the date they were created.
     * @param userID User identifier.
     * @param limit Amount of resources to get for the provided filter.
     * @param offset Amount of resources to skip for the provided filter.
     * @param token JWT authorization token.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun listUserContainerBookmarks(
        userID: String,
        limit: Int,
        offset: Int,
        token: String,
        listener: Listener<ContainersPaginated>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        val url = "$URL_USERS/$userID$URL_BOOKMARK_CONTAINERS" +
                "?$FIELD_NAME_PAGINATION_LIMIT=$limit" +
                "&$FIELD_NAME_PAGINATION_OFFSET=$offset" +
                "&sort=createdAt&order=desc"

        return object : JsonObjectRequest(
            Method.GET, url, null,
            { response ->
                val containers = ArrayList<Container>(limit)

                val jsonContainers = response.optJSONArray(FIELD_CONTAINERS)
                if (jsonContainers != null) {
                    for (i in 0 until jsonContainers.length()) {
                        containers.add(mapContainer(jsonContainers.optJSONObject(i)))
                    }
                }

                listener.onResponse(
                    ContainersPaginated(
                        response.optInt(FIELD_NAME_PAGINATION_TOTAL),
                        containers
                    )
                )
            },
            errorListener,
        ) {
            override fun getHeaders(): MutableMap<String, String> {
                return getHeaders(token)
            }
        }
    }

    /**
     * Creates a user container bookmark for the specified identifiers.
     * @param userID User identifier.
     * @param containerID Container identifier.
     * @param token JWT authorization token.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun createUserContainerBookmark(
        userID: String,
        containerID: String,
        token: String,
        listener: Listener<Unit>,
        errorListener: ErrorListener
    ): StringRequest {
        val url = "$URL_USERS/$userID$URL_BOOKMARK_CONTAINERS/$containerID"

        return object : StringRequest(
            Method.POST, url,
            { listener.onResponse(Unit) },
            errorListener,
        ) {
            override fun getHeaders(): MutableMap<String, String> {
                return getHeaders(token)
            }
        }
    }

    /**
     * Removes a user container bookmark for the specified identifiers.
     * @param userID User identifier.
     * @param containerID Container identifier.
     * @param token JWT authorization token.
     * @param listener Volley response listener.
     * @param errorListener Volley response error listener.
     * @return Volley request.
     */
    fun removeUserContainerBookmark(
        userID: String,
        containerID: String,
        token: String,
        listener: Listener<Unit>,
        errorListener: ErrorListener
    ): StringRequest {
        val url = "$URL_USERS/$userID$URL_BOOKMARK_CONTAINERS/$containerID"

        return object : StringRequest(
            Method.DELETE, url,
            { listener.onResponse(Unit) },
            errorListener,
        ) {
            override fun getHeaders(): MutableMap<String, String> {
                return getHeaders(token)
            }
        }
    }

    /**
     * Returns the headers commonly used in the API client.
     * @param token Authorization bearer token.
     * @return Headers map.
     */
    private fun getHeaders(token: String): HashMap<String, String> {
        val headers = HashMap<String, String>()
        headers[HEADER_KEY_AUTHORIZATION] = HEADER_VALUE_BEARER_PREFIX + token
        return headers
    }

    /**
     * Returns a domain Error object based on the given VolleyError.
     * @param error Volley error to map.
     * @return Domain Error data class.
     */
    fun mapError(error: VolleyError): Error {
        if (error.networkResponse == null) {
            return Error("", "")
        }
        if (error.networkResponse.data.isEmpty()) {
            return Error("", "")
        }

        val body = String(error.networkResponse.data)
        val json = JSONObject(body)

        return Error(
            json.optString(FIELD_NAME_ERROR_CODE),
            json.optString(FIELD_NAME_ERROR_MESSAGE),
        )
    }

    /**
     * Returns a domain User object based on the given JSONObject.
     * @param json JSON object to map.
     * @return Domain User data class.
     */
    private fun mapUser(json: JSONObject): User {
        return User(
            json.optString(FIELD_NAME_ID),
            json.optString(FIELD_NAME_USERNAME),
            json.optString(FIELD_NAME_FIRST_NAME),
            json.optString(FIELD_NAME_LAST_NAME),
            json.optString(FIELD_NAME_CREATED_AT),
            json.optString(FIELD_NAME_MODIFIED_AT)
        )
    }

    /**
     * Returns a domain Container object based on the given JSONObject.
     * @param json JSON object to map.
     * @return Domain Container data class.
     */
    private fun mapContainer(json: JSONObject): Container {
        val geoJSON = GeoJSONFeaturePoint()

        val geoJSONObject = json.optJSONObject(FIELD_CONTAINER_GEO_JSON)
        if (geoJSONObject != null) {
            val geoJSONGeometryObject =
                geoJSONObject.optJSONObject(FIELD_CONTAINER_GEO_JSON_GEOMETRY)
            if (geoJSONGeometryObject != null) {
                val geoJSONGeometryCoordinates =
                    geoJSONGeometryObject.optJSONArray(FIELD_CONTAINER_GEO_JSON_GEOMETRY_COORDINATES)
                if (geoJSONGeometryCoordinates != null && geoJSONGeometryCoordinates.length() == 2) {
                    try {
                        geoJSON.geometry.coordinates[0] = geoJSONGeometryCoordinates.getDouble(0)
                        geoJSON.geometry.coordinates[1] = geoJSONGeometryCoordinates.getDouble(1)
                    } catch (e: JSONException) {
                        Log.e(LOG_TAG, e.message, e)
                    }
                }
            }

            val geoJSONPropertiesObject =
                geoJSONObject.optJSONObject(FIELD_CONTAINER_GEO_JSON_PROPERTIES)
            if (geoJSONPropertiesObject != null) {
                geoJSON.properties = GeoJSONProperties(
                    geoJSONPropertiesObject.optString(
                        FIELD_CONTAINER_GEO_JSON_PROPERTIES_WAY_NAME
                    ),
                    geoJSONPropertiesObject.optString(
                        FIELD_CONTAINER_GEO_JSON_PROPERTIES_MUNICIPALITY_NAME
                    )
                )
            }
        }

        return Container(
            json.optString(FIELD_CONTAINER_ID),
            mapDomainContainerCategory(json.optString(FIELD_CONTAINER_CATEGORY)),
            geoJSON,
            json.optString(FIELD_CONTAINER_CREATED_AT),
            json.optString(FIELD_CONTAINER_MODIFIED_AT)
        )
    }

    /**
     * Returns a domain ContainerCategory object based on the given string value. If the given value
     * is unexpected, it defaults to the general category.
     * @param value Value to map.
     * @return Domain ContainerCategory data class.
     */
    private fun mapDomainContainerCategory(value: String): ContainerCategory {
        return when (value) {
            FIELD_CONTAINER_GENERAL -> ContainerCategory.GENERAL
            FIELD_CONTAINER_PAPER -> ContainerCategory.PAPER
            FIELD_CONTAINER_PLASTIC -> ContainerCategory.PLASTIC
            FIELD_CONTAINER_METAL -> ContainerCategory.METAL
            FIELD_CONTAINER_GLASS -> ContainerCategory.GLASS
            FIELD_CONTAINER_ORGANIC -> ContainerCategory.ORGANIC
            FIELD_CONTAINER_HAZARDOUS -> ContainerCategory.HAZARDOUS
            else -> ContainerCategory.GENERAL
        }
    }

    /**
     * Returns a string value based on the given domain ContainerCategory object.
     * @param value Value to map.
     * @return Domain ContainerCategory data class.
     */
    private fun mapDomainContainerCategory(value: ContainerCategory): String {
        return when (value) {
            ContainerCategory.GENERAL -> FIELD_CONTAINER_GENERAL
            ContainerCategory.PAPER -> FIELD_CONTAINER_PAPER
            ContainerCategory.PLASTIC -> FIELD_CONTAINER_PLASTIC
            ContainerCategory.METAL -> FIELD_CONTAINER_METAL
            ContainerCategory.GLASS -> FIELD_CONTAINER_GLASS
            ContainerCategory.ORGANIC -> FIELD_CONTAINER_ORGANIC
            ContainerCategory.HAZARDOUS -> FIELD_CONTAINER_HAZARDOUS
        }
    }
}
