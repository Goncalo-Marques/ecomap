package com.ecomap.ecomap.clients.ecomap.http

import android.util.Log
import com.android.volley.Request
import com.android.volley.Response.ErrorListener
import com.android.volley.Response.Listener
import com.android.volley.VolleyError
import com.android.volley.toolbox.JsonObjectRequest
import com.ecomap.ecomap.domain.Error
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

    // Authentication field names.
    private const val FIELD_NAME_TOKEN = "token"

    // Error field names.
    private const val FIELD_NAME_ERROR_CODE = "code"
    private const val FIELD_NAME_ERROR_MESSAGE = "message"

    // User field names.
    private const val FIELD_NAME_ID = "id"
    private const val FIELD_NAME_USERNAME = "username"
    private const val FIELD_NAME_PASSWORD = "password"
    private const val FIELD_NAME_FIRST_NAME = "firstName"
    private const val FIELD_NAME_LAST_NAME = "lastName"
    private const val FIELD_NAME_CREATED_AT = "createdAt"
    private const val FIELD_NAME_MODIFIED_AT = "modifiedAt"

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
        listener: Listener<String?>,
        errorListener: ErrorListener
    ): JsonObjectRequest {
        val requestPayload = JSONObject()
        requestPayload.put(FIELD_NAME_USERNAME, username)
        requestPayload.put(FIELD_NAME_PASSWORD, password)

        return JsonObjectRequest(
            Request.Method.POST, URL_USERS_SIGN_IN, requestPayload,
            { response ->
                var token: String? = null
                try {
                    token = response.optString(FIELD_NAME_TOKEN)
                } catch (e: JSONException) {
                    Log.e(LOG_TAG, e.message, e)
                }

                listener.onResponse(token)
            },
            errorListener
        )
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
            { response ->
                listener.onResponse(
                    User(
                        response.optString(FIELD_NAME_ID),
                        response.optString(FIELD_NAME_USERNAME),
                        response.optString(FIELD_NAME_FIRST_NAME),
                        response.optString(FIELD_NAME_LAST_NAME),
                        response.optString(FIELD_NAME_CREATED_AT),
                        response.optString(FIELD_NAME_MODIFIED_AT)
                    )
                )
            },
            errorListener
        )
    }

    /**
     * Returns a domain Error object based on the given VolleyError.
     * @param error Volley error.
     * @return Domain Error data object.
     */
    fun getError(error: VolleyError): Error {
        val body = String(error.networkResponse.data)
        val json = JSONObject(body)

        return Error(
            json.optString(FIELD_NAME_ERROR_CODE),
            json.optString(FIELD_NAME_ERROR_MESSAGE),
        )
    }
}
