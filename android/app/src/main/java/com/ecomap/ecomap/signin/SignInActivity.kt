package com.ecomap.ecomap.signin

import android.content.Intent
import android.os.Bundle
import android.widget.Button
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

class SignInActivity : AppCompatActivity() {
    private lateinit var textInputEditTextUsername: TextInputEditText
    private lateinit var textInputEditTextPassword: TextInputEditText
    private lateinit var buttonSignIn: Button
    private lateinit var buttonCreateAccount: Button

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
        store = UserStore(this.applicationContext)

        // Get activity views.
        textInputEditTextUsername = findViewById(R.id.text_input_edit_text_username)
        textInputEditTextPassword = findViewById(R.id.text_input_edit_text_password)
        buttonSignIn = findViewById(R.id.button_sign_in)
        buttonCreateAccount = findViewById(R.id.button_create_account)

        // Set up on click events for the sign in and create account button.
        buttonSignIn.setOnClickListener { signInUser() }
        buttonCreateAccount.setOnClickListener { openCreateAccountScreen() }
    }

    /**
     * Opens create account screen.
     */
    private fun openCreateAccountScreen() {
        val intent = Intent(this, CreateAccountActivity::class.java)
        startActivity(intent)
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

        val request =
            ApiClient.signIn(
                username,
                password,
                { token ->
                    if (token.isEmpty()) {
                        Toast.makeText(
                            this.applicationContext,
                            getString(R.string.error_sign_in),
                            Toast.LENGTH_LONG
                        )
                            .show()
                        return@signIn
                    }

                    val intent = Intent(this, MainActivity::class.java)

                    // Flags the intent to mark the activity as the root in the history stack,
                    // clearing out any other tasks.
                    intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)

                    // Stores token in UserStore.
                    runBlocking { store.storeToken(token) }

                    startActivity(intent)
                },
                { _ ->
                    Toast.makeText(
                        this.applicationContext,
                        getString(R.string.error_sign_in),
                        Toast.LENGTH_LONG
                    )
                        .show()
                })

        ApiRequestQueue.getInstance(this.applicationContext).add(request)
    }
}
