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
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.ecomap.ecomap.Common
import com.ecomap.ecomap.R
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.clients.ecomap.http.ApiRequestQueue
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.domain.ContainersPaginated
import com.ecomap.ecomap.signin.SignInActivity
import com.google.android.gms.maps.model.LatLng

class UserAccountActivity : AppCompatActivity() {
    private lateinit var textViewFirstName: TextView
    private lateinit var textViewLastName: TextView
    private lateinit var textViewUsername: TextView
    private lateinit var textViewContainerBookmarksEmpty: TextView
    private lateinit var recyclerViewContainerBookmarks: RecyclerView
    private lateinit var progressBar: ProgressBar

    private lateinit var store: UserStore
    private lateinit var token: String
    private lateinit var userID: String

    private lateinit var recyclerViewContainerBookmarksDataSet: ArrayList<ContainerBookmarkRecyclerViewData>

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
        textViewContainerBookmarksEmpty = findViewById(R.id.text_view_container_bookmarks_empty)
        recyclerViewContainerBookmarks = findViewById(R.id.recycler_container_bookmarks)
        progressBar = findViewById(R.id.progress_bar_user_account)

        // Set up on click events for the buttons.
        buttonSignOut.setOnClickListener { signOutUser() }

        // Set up container bookmarks recycler view.
        recyclerViewContainerBookmarksDataSet = arrayListOf()

        recyclerViewContainerBookmarks.layoutManager = LinearLayoutManager(this)
        val recyclerViewContainerBookmarksAdapter =
            ContainerBookmarksRecyclerViewAdapter(this, recyclerViewContainerBookmarksDataSet)
        recyclerViewContainerBookmarks.adapter = recyclerViewContainerBookmarksAdapter

        updateUserContainerBookmarksVisibility()
        recyclerViewContainerBookmarksAdapter.onButtonContainerBookmarkClicked = { itemPosition ->
            removeUserContainerBookmark(itemPosition)
            updateUserContainerBookmarksVisibility()
        }

        // Get user store and token.
        store = UserStore(applicationContext)
        token = store.getToken().toString()
        userID = Common.getSubjectFromJWT(token)

        // Show progress bar while data is still loading.
        progressBar.visibility = View.VISIBLE

        // Update UI with the user personal information and bookmarks.
        updateUserPersonalInformationUI()
        updateUserContainerBookmarksUI()
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
     * Removes the user container bookmark associated with the data set at the specified position.
     * It also notifies the recycler view adapter to update the affected item.
     */
    private fun removeUserContainerBookmark(position: Int) {
        // Remove the user container bookmark.
        for (container in recyclerViewContainerBookmarksDataSet[position].containers) {
            val request =
                ApiClient.removeUserContainerBookmark(userID, container.id, token,
                    {}, {})
            ApiRequestQueue.getInstance(this).add(request)
        }

        // Remove the item from the data set.
        recyclerViewContainerBookmarksDataSet.removeAt(position)
        recyclerViewContainerBookmarks.adapter?.notifyItemRemoved(position)
        recyclerViewContainerBookmarks.adapter?.notifyItemRangeChanged(
            position,
            recyclerViewContainerBookmarksDataSet.size
        )
    }

    /**
     * Makes the container bookmark list visible if the data set is not empty, otherwise makes it
     * invisible.
     */
    private fun updateUserContainerBookmarksVisibility() {
        if (recyclerViewContainerBookmarksDataSet.isEmpty()) {
            textViewContainerBookmarksEmpty.visibility = View.VISIBLE
            recyclerViewContainerBookmarks.visibility = View.GONE
        } else {
            textViewContainerBookmarksEmpty.visibility = View.GONE
            recyclerViewContainerBookmarks.visibility = View.VISIBLE
        }
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

    /**
     * Gets the current list of containers that the user has bookmarked and sets them in the UI.
     */
    private fun updateUserContainerBookmarksUI() {
        // Map containing the containers, merging those that are in the same position to be
        // contained in the same item.
        val mappedContainers = mutableMapOf<LatLng, ContainerBookmarkRecyclerViewData>()

        // Helper function to handle a successful response.
        val handleSuccess = fun(paginatedContainers: ContainersPaginated) {
            if (isFinishing || isDestroyed) {
                return
            }

            for (container in paginatedContainers.containers) {
                val containerCoordinates = container.geoJSON.geometry.coordinates
                val containerPosition = LatLng(containerCoordinates[1], containerCoordinates[0])

                // Add the container if it is not currently in the data set, otherwise append the
                // container category to the existing item.
                val existingContainer = mappedContainers[containerPosition]
                if (existingContainer == null) {
                    val containerBookmarkData =
                        ContainerBookmarkRecyclerViewData(arrayListOf(container))

                    recyclerViewContainerBookmarksDataSet.add(containerBookmarkData)
                    recyclerViewContainerBookmarks.adapter?.notifyItemInserted(
                        recyclerViewContainerBookmarksDataSet.size - 1
                    )

                    mappedContainers[containerPosition] = containerBookmarkData
                } else {
                    existingContainer.containers.add(container)

                    // Find the position of the changed item.
                    for ((index, data) in recyclerViewContainerBookmarksDataSet.withIndex()) {
                        for (c in data.containers) {
                            if (c.id == container.id) {
                                recyclerViewContainerBookmarks.adapter?.notifyItemChanged(index)
                                break
                            }
                        }
                    }
                }
            }

            updateUserContainerBookmarksVisibility()
        }

        // Execute the request to get all existing user container bookmarks and add them to the list.
        val request = ApiClient.listUserContainerBookmarks(
            userID,
            REQUEST_LIST_CONTAINER_LIMIT,
            0,
            token,
            { paginatedContainers ->
                val remainingRequest =
                    paginatedContainers.total / REQUEST_LIST_CONTAINER_LIMIT
                for (i in 1..remainingRequest) {
                    ApiRequestQueue.getInstance(applicationContext).add(
                        ApiClient.listUserContainerBookmarks(
                            userID,
                            REQUEST_LIST_CONTAINER_LIMIT,
                            REQUEST_LIST_CONTAINER_LIMIT * i,
                            token,
                            { handleSuccess(it) },
                            { Common.handleVolleyError(this, this, it) }
                        )
                    )
                }

                handleSuccess(paginatedContainers)
            },
            { Common.handleVolleyError(this, this, it) })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }

    companion object {
        private const val REQUEST_LIST_CONTAINER_LIMIT = 100
    }
}
