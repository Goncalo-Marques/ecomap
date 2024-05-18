package com.ecomap.ecomap.clients.ecomap.http

import android.util.Log
import com.android.volley.Request
import com.android.volley.Response.ErrorListener
import com.android.volley.Response.Listener
import com.android.volley.toolbox.JsonObjectRequest
import com.ecomap.ecomap.domain.User
import org.json.JSONException
import org.json.JSONObject

const val BASE_API_URL = "https://server-7fzc7ivuwa-ew.a.run.app/api"

/**
 * EcoMap HTTP API client.
 */
object ApiClient {
    /**
     * Retrieves the URL of the API.
     * @param endpoint API endpoint.
     * @return API URL.
     */
    private fun getApiUrl(endpoint: String): String {
        return BASE_API_URL + endpoint
    }

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
        requestPayload.put("username", username)
        requestPayload.put("password", password)

        return JsonObjectRequest(
            Request.Method.POST, getApiUrl("/users/signin"), requestPayload,
            { response ->
                var token: String? = null
                try {
                    token = response.getString("token")
                } catch (e: JSONException) {
                    Log.e(javaClass.name, e.message, e)
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
        requestPayload.put("firstName", firstName)
        requestPayload.put("lastName", lastName)
        requestPayload.put("username", username)
        requestPayload.put("password", password)

        return JsonObjectRequest(
            Request.Method.POST, getApiUrl("/users"), requestPayload,
            { response ->
                listener.onResponse(
                    User(
                        response.getString("id"),
                        response.getString("username"),
                        response.getString("firstName"),
                        response.getString("lastName"),
                        response.getString("createdAt"),
                        response.getString("modifiedAt")
                    )
                )
            },
            errorListener
        )
    }
}
