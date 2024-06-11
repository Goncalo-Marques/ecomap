package com.ecomap.ecomap.data

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

/**
 * User data store.
 * @param context Application context.
 */
class UserStore(private val context: Context) {
    companion object {
        private val Context.dataStore by preferencesDataStore(name = "user")
        private val USER_TOKEN = stringPreferencesKey("user_token")
    }

    /**
     * Retrieves user JWT token.
     * @return Coroutine with JWT token or `null` when JWT token is not found.
     */
    fun getToken(): String? {
        return runBlocking {
            val preferences = context.dataStore.data.first()
            preferences[USER_TOKEN]
        }
    }

    /**
     * Stores JWT token.
     * @param token JWT token.
     */
    fun storeToken(token: String) {
        runBlocking {
            context.dataStore.edit { preferences ->
                preferences[USER_TOKEN] = token
            }
        }
    }

    /**
     * Removes JWT token.
     */
    fun removeToken() {
        runBlocking {
            context.dataStore.edit { preferences ->
                preferences.remove(USER_TOKEN)
            }
        }
    }
}
