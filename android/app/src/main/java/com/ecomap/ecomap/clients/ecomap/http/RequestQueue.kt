package com.ecomap.ecomap.clients.ecomap.http

import android.content.Context
import com.android.volley.Request
import com.android.volley.RequestQueue
import com.android.volley.toolbox.Volley

/**
 * EcoMap API Volley request queue singleton.
 */
class ApiRequestQueue(context: Context) {
    companion object {
        @Volatile
        private var INSTANCE: ApiRequestQueue? = null
        fun getInstance(context: Context) =
            INSTANCE ?: synchronized(this) {
                INSTANCE ?: ApiRequestQueue(context).also {
                    INSTANCE = it
                }
            }
    }

    /**
     * Request queue.
     */
    private val requestQueue: RequestQueue by lazy {
        Volley.newRequestQueue(context.applicationContext)
    }

    /**
     * Adds a request to the request queue.
     * @param request Request to be added to the queue.
     */
    fun <T> add(request: Request<T>) {
        requestQueue.add(request)
    }
}
