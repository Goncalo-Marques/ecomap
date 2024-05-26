package com.ecomap.ecomap.domain

/**
 * Represents the structure of an error message.
 * @param code Error code.
 * @param message Error message.
 */
data class Error(
    val code: String,
    val message: String,
)
