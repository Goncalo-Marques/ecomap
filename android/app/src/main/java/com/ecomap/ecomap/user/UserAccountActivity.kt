package com.ecomap.ecomap.user

import android.content.Intent
import android.os.Bundle
import android.view.MenuItem
import android.view.View
import android.widget.Button
import android.widget.ProgressBar
import android.widget.TextView
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.ecomap.ecomap.Common
import com.ecomap.ecomap.R
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.clients.ecomap.http.ApiRequestQueue
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.signin.SignInActivity

class UserAccountActivity : AppCompatActivity() {
    private lateinit var textViewFirstName: TextView
    private lateinit var textViewLastName: TextView
    private lateinit var textViewUsername: TextView
    private lateinit var progressBar: ProgressBar

    private lateinit var store: UserStore
    private lateinit var token: String
    private lateinit var userID: String

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_user_account)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // Display action bar.
        setSupportActionBar(findViewById(R.id.toolbar_user_account))

        // Enable back button on action bar.
        supportActionBar?.setDisplayHomeAsUpEnabled(true)

        // Get activity views.
        val buttonSignOut: Button = findViewById(R.id.button_sign_out)
        textViewFirstName = findViewById(R.id.text_view_first_name_value)
        textViewLastName = findViewById(R.id.text_view_last_name_value)
        textViewUsername = findViewById(R.id.text_view_username_value)
        progressBar = findViewById(R.id.progress_bar_user_account)

        // Set up on click event for the sign out button.
        buttonSignOut.setOnClickListener { signOutUser() }

        // Show progress bar while data is still loading.
        progressBar.visibility = View.VISIBLE

        // Get user store and token.
        store = UserStore(applicationContext)
        token = store.getToken().toString()
        userID = Common.getSubjectFromJWT(token)

        // Update UI with the user personal information and bookmarks.
        updateUserPersonalInformationUI()
        // TODO: Load bookmarks.
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            android.R.id.home -> {
                // Finish the current activity when the home button on the action bar is clicked.
                finish()
                true
            }

            else -> super.onOptionsItemSelected(item)
        }
    }

    /**
     * Signs the user out by deleting the token from the store and starting the sign in activity.
     */
    private fun signOutUser() {
        // Removes token from UserStore.
        store.removeToken()

        val intentSignInActivity = Intent(this, SignInActivity::class.java)
        startActivity(intentSignInActivity)

        finishAffinity()
    }

    /**
     * Gets the user's personal information and sets it in the UI.
     */
    private fun updateUserPersonalInformationUI() {
        val request = ApiClient.getAccount(
            userID, token,
            { userAccount ->
                // Set user information.
                textViewFirstName.text = userAccount.firstName
                textViewLastName.text = userAccount.lastName
                textViewUsername.text = userAccount.username

                // Hide the progress bar when the user information is loaded.
                progressBar.visibility = View.INVISIBLE
            },
            {
                // Hide the progress bar when a network error occurs.
                progressBar.visibility = View.INVISIBLE

                Common.handleVolleyError(this, this, it)
            }
        )

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }
}
