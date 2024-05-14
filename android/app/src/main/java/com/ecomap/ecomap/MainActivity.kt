package com.ecomap.ecomap

import android.content.Intent
import android.os.Bundle
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.signin.SignInActivity
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class MainActivity : AppCompatActivity() {
    private lateinit var store: UserStore

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_main)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        store = UserStore(this.applicationContext)

        // User token validation.
        val intent = Intent(this, SignInActivity::class.java)

        // Flags the intent to mark the activity as the root in the history stack,
        // clearing out any other tasks.
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)

        runBlocking {
            val token = store.getToken().first()
            if (token == null) {
                // Open SignIn Activity if token is not found.
                startActivity(intent)
            }
        }
    }
}
