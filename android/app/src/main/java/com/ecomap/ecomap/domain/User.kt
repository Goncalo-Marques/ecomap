package com.ecomap.ecomap.domain

/**
 * Represents the structure of a user.
 * @param id User ID.
 * @param username User username.
 * @param firstName User first name.
 * @param lastName User last name.
 * @param createdAt Date the user was created.
 * @param modifiedAt Date the user was last modified.
 */
data class User(
    val id: String,
    val username: String,
    val firstName: String,
    val lastName: String,
    val createdAt: String,
    val modifiedAt: String,
)
