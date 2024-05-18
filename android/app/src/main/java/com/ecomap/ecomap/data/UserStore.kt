package com.ecomap.ecomap.data

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

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
    fun getToken(): Flow<String?> {
        return context.dataStore.data
            .map { preferences ->
                preferences[USER_TOKEN]
            }
    }

    /**
     * Stores JWT token.
     * @param token JWT token.
     */
    suspend fun storeToken(token: String) {
        context.dataStore.edit { preferences ->
            preferences[USER_TOKEN] = token
        }
    }
}
