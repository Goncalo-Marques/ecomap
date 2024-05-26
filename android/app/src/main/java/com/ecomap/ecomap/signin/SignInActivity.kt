package com.ecomap.ecomap.signin

import android.content.Intent
import android.os.Bundle
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

class SignInActivity : AppCompatActivity() {
    private lateinit var textInputEditTextUsername: TextInputEditText
    private lateinit var textInputEditTextPassword: TextInputEditText
    private lateinit var progressBar: ProgressBar

    private lateinit var store: UserStore

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_sign_in)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // Retrieve user store.
        store = UserStore(applicationContext)

        // Get activity views.
        textInputEditTextUsername = findViewById(R.id.text_input_edit_text_username)
        textInputEditTextPassword = findViewById(R.id.text_input_edit_text_password)
        val buttonSignIn: Button = findViewById(R.id.button_sign_in)
        val buttonCreateAccount: Button = findViewById(R.id.button_create_account)
        progressBar = findViewById(R.id.progress_bar_sign_in)

        // Set up on click events for the sign in and create account button.
        buttonSignIn.setOnClickListener { signInUser() }
        buttonCreateAccount.setOnClickListener { openCreateAccountScreen() }

        // Hide progress bar when activity is created.
        progressBar.visibility = View.INVISIBLE
    }

    /**
     * Opens create account screen.
     */
    private fun openCreateAccountScreen() {
        val intentCreateAccountActivity = Intent(this, CreateAccountActivity::class.java)
        startActivity(intentCreateAccountActivity)
    }

    /**
     * Signs in a user performing the respective form validations.
     */
    private fun signInUser() {
        val username = textInputEditTextUsername.text.toString()
        val password = textInputEditTextPassword.text.toString()

        if (username.isBlank()) {
            textInputEditTextUsername.error = getString(R.string.sign_in_username_required_error)
        }
        if (password.isBlank()) {
            textInputEditTextPassword.error = getString(R.string.sign_in_password_required_error)
        }

        if (username.isBlank() || password.isBlank()) {
            return
        }

        // Display progress bar.
        progressBar.visibility = View.VISIBLE

        val request =
            ApiClient.signIn(
                username,
                password,
                { token ->
                    // Stores token in UserStore.
                    store.storeToken(token)

                    val intentMainActivity = Intent(this, MainActivity::class.java)
                    startActivity(intentMainActivity)

                    finishAffinity()
                },
                {
                    // Hide the progress bar when a network error occurs.
                    progressBar.visibility = View.INVISIBLE

                    Toast.makeText(
                        applicationContext,
                        getString(R.string.error_sign_in_invalid_credentials),
                        Toast.LENGTH_LONG
                    ).show()
                })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }
}
