package com.ecomap.ecomap

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.widget.Toast
import com.android.volley.VolleyError
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.signin.SignInActivity

/**
 * Class containing common methods and attributes.
 */
class Common {
    companion object {
        fun handleVolleyError(context: Context, activity: Activity, error: VolleyError) {
            if (activity.isFinishing || activity.isDestroyed) {
                return
            }

            // Check if the server found the token to be invalid, and if so, start the sign in
            // activity.
            if (error.networkResponse?.statusCode == 401) {
                Toast.makeText(
                    context.applicationContext,
                    context.getString(R.string.error_sign_in_expired),
                    Toast.LENGTH_LONG
                ).show()

                val intentSignInActivity = Intent(context, SignInActivity::class.java)

                // Flags the intent to mark the activity as the root in the history stack,
                // clearing out any other tasks.
                intentSignInActivity.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)

                context.startActivity(intentSignInActivity)
                activity.finish()
                return
            }

            // Handle other type of errors.
            val errorResponse = ApiClient.mapError(error)

            var errorMessage = errorResponse.code
            if (errorResponse.message.isNotEmpty()) {
                errorMessage = errorResponse.message
            }

            Toast.makeText(
                context.applicationContext,
                errorMessage,
                Toast.LENGTH_LONG
            ).show()
        }
    }
}
