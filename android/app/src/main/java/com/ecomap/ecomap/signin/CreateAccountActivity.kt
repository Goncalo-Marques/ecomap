package com.ecomap.ecomap.signin

import android.content.Intent
import android.os.Bundle
import android.view.MenuItem
import android.view.View
import android.widget.Button
import android.widget.ProgressBar
import android.widget.Toast
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.ecomap.ecomap.MainActivity
import com.ecomap.ecomap.R
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.clients.ecomap.http.ApiRequestQueue
import com.ecomap.ecomap.data.UserStore
import com.google.android.material.textfield.TextInputEditText
import kotlinx.coroutines.runBlocking

class CreateAccountActivity : AppCompatActivity() {
    private lateinit var textInputEditTextFirstName: TextInputEditText
    private lateinit var textInputEditTextLastName: TextInputEditText
    private lateinit var textInputEditTextUsername: TextInputEditText
    private lateinit var textInputEditTextPassword: TextInputEditText
    private lateinit var buttonCreateAccount: Button
    private lateinit var progressBar: ProgressBar

    private lateinit var store: UserStore

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_create_account)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // Display action bar.
        setSupportActionBar(findViewById(R.id.toolbar_create_account))

        // Enable back button on action bar.
        supportActionBar?.setDisplayHomeAsUpEnabled(true)

        // Retrieve user store.
        store = UserStore(applicationContext)

        // Get activity views.
        textInputEditTextFirstName = findViewById(R.id.text_input_edit_text_first_name)
        textInputEditTextLastName = findViewById(R.id.text_input_edit_text_last_name)
        textInputEditTextUsername = findViewById(R.id.text_input_edit_text_username)
        textInputEditTextPassword = findViewById(R.id.text_input_edit_text_password)
        buttonCreateAccount = findViewById(R.id.button_create_account)
        progressBar = findViewById(R.id.progress_bar_create_account)

        // Set up on click event for the create account button.
        buttonCreateAccount.setOnClickListener { createUser() }

        // Hide progress bar when activity is created.
        progressBar.visibility = View.INVISIBLE
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
     * Creates a user performing the respective form validations.
     */
    private fun createUser() {
        // Get values.
        val firstName = textInputEditTextFirstName.text.toString()
        val lastName = textInputEditTextLastName.text.toString()
        val username = textInputEditTextUsername.text.toString()
        val password = textInputEditTextPassword.text.toString()

        // Validate values.
        if (firstName.isBlank()) {
            textInputEditTextFirstName.error = getString(R.string.first_name_required_error)
        }
        if (lastName.isBlank()) {
            textInputEditTextLastName.error = getString(R.string.last_name_required_error)
        }
        if (username.isBlank()) {
            textInputEditTextUsername.error = getString(R.string.sign_in_username_required_error)
        }
        if (password.isBlank()) {
            textInputEditTextPassword.error = getString(R.string.sign_in_password_required_error)
        }

        // Prevent server request if fields any field is blank.
        if (firstName.isBlank() || lastName.isBlank() || username.isBlank() || password.isBlank()) {
            return
        }

        // Display progress bar.
        progressBar.visibility = View.VISIBLE

        // Create the request to create the user.
        val request =
            ApiClient.createAccount(firstName, lastName, username, password,
                { signInUser(username, password) },
                { error ->
                    // Hide the progress bar when a network error occurs.
                    progressBar.visibility = View.INVISIBLE

                    val errorResponse = ApiClient.mapError(error)

                    var errorMessage = errorResponse.code
                    if (errorResponse.message.isNotEmpty()) {
                        errorMessage = errorResponse.message
                    }

                    Toast.makeText(
                        applicationContext,
                        errorMessage,
                        Toast.LENGTH_LONG
                    ).show()
                })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }

    /**
     * Signs in a user.
     * @param username User username.
     * @param password User password.
     */
    private fun signInUser(username: String, password: String) {
        val request = ApiClient.signIn(username, password,
            { token ->
                if (token.isEmpty()) {
                    Toast.makeText(
                        applicationContext,
                        getString(R.string.error_create_account),
                        Toast.LENGTH_LONG
                    ).show()
                    return@signIn
                }

                val intentMainActivity = Intent(this, MainActivity::class.java)

                // Flags the intent to mark the activity as the root in the history stack,
                // clearing out any other tasks.
                intentMainActivity.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)

                // Stores token in UserStore.
                runBlocking { store.storeToken(token) }

                startActivity(intentMainActivity)
                finish()
            },
            { _ ->
                // Hide the progress bar when a network error occurs.
                progressBar.visibility = View.INVISIBLE

                Toast.makeText(
                    applicationContext,
                    getString(R.string.error_create_account),
                    Toast.LENGTH_LONG
                ).show()
            })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }
}
