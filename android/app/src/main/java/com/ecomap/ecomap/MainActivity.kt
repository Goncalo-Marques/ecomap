package com.ecomap.ecomap

import android.Manifest
import android.annotation.SuppressLint
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import android.widget.ImageButton
import android.widget.TextView
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.constraintlayout.widget.ConstraintLayout
import androidx.constraintlayout.widget.Group
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import androidx.core.view.isVisible
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.clients.ecomap.http.ApiRequestQueue
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.domain.Container
import com.ecomap.ecomap.domain.ContainerCategory
import com.ecomap.ecomap.domain.ContainersPaginated
import com.ecomap.ecomap.map.ContainerCategoriesRecyclerViewAdapter
import com.ecomap.ecomap.map.ContainerCategoryRecyclerViewData
import com.ecomap.ecomap.map.ContainerClusterRenderer
import com.ecomap.ecomap.map.ContainerMarker
import com.ecomap.ecomap.signin.SignInActivity
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.GoogleMap
import com.google.android.gms.maps.GoogleMapOptions
import com.google.android.gms.maps.OnMapReadyCallback
import com.google.android.gms.maps.SupportMapFragment
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.LatLngBounds
import com.google.android.material.chip.Chip
import com.google.android.material.chip.ChipGroup
import com.google.android.material.floatingactionbutton.FloatingActionButton
import com.google.maps.android.clustering.ClusterManager
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class MainActivity : AppCompatActivity(), OnMapReadyCallback {
    /**
     * Defines the Google Map instance.
     */
    private lateinit var map: GoogleMap

    /**
     * The entry point to the Fused Location Provider.
     */
    private lateinit var fusedLocationProviderClient: FusedLocationProviderClient

    /**
     * Defines whether the location permission is granted.
     */
    private var locationPermissionGranted = false

    /**
     * Defines the cluster manager of the container markers.
     */
    private lateinit var containerClusterManager: ClusterManager<ContainerMarker>

    /**
     * Defines the group of buttons view.
     */
    private lateinit var groupButtonsView: Group

    /**
     * Defines the container info window view.
     */
    private lateinit var containerInfoWindowView: ConstraintLayout

    /**
     * Defines the container info window title view.
     */
    private lateinit var containerInfoWindowTitleText: TextView

    /**
     * Defines the container info window snippet view.
     */
    private lateinit var containerInfoWindowSnippetText: TextView

    /**
     * Defines the container info window bookmark button.
     */
    private lateinit var containerInfoWindowBookmarkButton: ImageButton


    /**
     * Defines the container info window categories recycler view.
     */
    private lateinit var containerInfoWindowRecyclerCategories: RecyclerView

    /**
     * Defines the container info window directions button.
     */
    private lateinit var containerInfoWindowDirectionsButton: Button

    /**
     * Defines the containers the user has bookmarked.
     * It is contained in a map with the container id as key.
     */
    private var userContainerBookmarks: MutableMap<String, Container> = mutableMapOf()

    /**
     * Defines the authentication token.
     */
    private lateinit var token: String

    /**
     * Defines the identifier of the current user.
     */
    private lateinit var userID: String

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_main)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // User token validation.
        val intentSignInActivity = Intent(this, SignInActivity::class.java)

        // Check whether the user store contains the login token.
        // If not, start the SignIn Activity.
        val store = UserStore(applicationContext)
        runBlocking {
            val storeToken = store.getToken().first()
            if (storeToken == null) {
                startActivity(intentSignInActivity)
                finish()
            }

            token = storeToken.toString()
            userID = Common.getSubjectFromJWT(token)
        }

        // Construct the main entry point for the Android location services.
        fusedLocationProviderClient = LocationServices.getFusedLocationProviderClient(this)

        // Google Map configurations.
        val mapLatLngBounds = LatLngBounds(
            LatLng(MAP_CAMERA_BOUND_SW_LAT, MAP_CAMERA_BOUND_SW_LNG),
            LatLng(MAP_CAMERA_BOUND_NE_LAT, MAP_CAMERA_BOUND_NE_LNG)
        )

        val googleMapOptions = GoogleMapOptions()
        googleMapOptions
            .mapType(GoogleMap.MAP_TYPE_NORMAL)
            .latLngBoundsForCameraTarget(mapLatLngBounds)
            .mapToolbarEnabled(false)

        // Add support map fragment to the map container.
        val mapFragment = SupportMapFragment.newInstance(googleMapOptions)
        supportFragmentManager
            .beginTransaction()
            .add(R.id.fragment_container_view_map, mapFragment)
            .commit()

        // Register the map callback.
        mapFragment.getMapAsync(this)

        // Get activity views.
        val chipGroupContainerFilter: ChipGroup = findViewById(R.id.chip_group_container_filter)
        val buttonMyLocation: FloatingActionButton = findViewById(R.id.button_my_location)
        groupButtonsView = findViewById(R.id.group_buttons)
        containerInfoWindowView = findViewById(R.id.info_window)
        containerInfoWindowTitleText = findViewById(R.id.info_window_text_title)
        containerInfoWindowSnippetText = findViewById(R.id.info_window_text_snippet)
        containerInfoWindowBookmarkButton = findViewById(R.id.info_window_button_container_bookmark)
        containerInfoWindowRecyclerCategories = findViewById(R.id.info_window_recycler_categories)
        containerInfoWindowRecyclerCategories.layoutManager =
            GridLayoutManager(this, CONTAINER_INFO_WINDOW_RECYCLER_CATEGORIES_SPAN_COUNT)
        containerInfoWindowDirectionsButton = findViewById(R.id.info_window_button_directions)

        // Set button functions.
        populateChipGroupContainerFilter(chipGroupContainerFilter)
        buttonMyLocation.setOnClickListener { focusMyLocation() }

        // Start with the container info window closed.
        closeContainerInfoWindow()
    }

    /**
     * Populates the given chip group with all the available container categories.
     */
    private fun populateChipGroupContainerFilter(chipGroup: ChipGroup) {
        for (category in ContainerCategory.entries) {
            val chip = Chip(this)

            // Set the chip style.
            chip.chipIcon = ContextCompat.getDrawable(this, category.getIconResource())
            chip.text = category.getStringResource(this)
            chip.isCheckable = true

            // Set the chip function.
            chip.setOnClickListener {
                // Filter the current containers on the map based on the chip container category.
                // If the chip is not checked, show all available containers regardless of their
                // category.
                if (chip.isChecked) {
                    updateContainersUI(category)
                } else {
                    updateContainersUI()
                }
            }

            // Add the chip to the group.
            chipGroup.addView(chip)
        }
    }

    /**
     * Function called when the Google Map is ready.
     */
    @SuppressLint("PotentialBehaviorOverride")
    override fun onMapReady(googleMap: GoogleMap) {
        map = googleMap
        map.setPadding(MAP_PADDING_LEFT, MAP_PADDING_TOP, MAP_PADDING_RIGHT, MAP_PADDING_BOTTOM)

        // Initialize the container cluster manager.
        containerClusterManager = ClusterManager(this, map)
        containerClusterManager.renderer =
            ContainerClusterRenderer(this, map, containerClusterManager)
        map.setOnCameraIdleListener(containerClusterManager)
        map.setOnMarkerClickListener(containerClusterManager)

        // Set container info window functions.
        containerClusterManager.setOnClusterItemClickListener { container ->
            // Display the container information window and move to the container location.
            showContainerInfoWindow(container)
            map.animateCamera(
                CameraUpdateFactory.newLatLngZoom(
                    container.position,
                    map.cameraPosition.zoom
                )
            )

            // Returns true so that the default info window is not displayed.
            true
        }
        containerClusterManager.setOnClusterClickListener {
            closeContainerInfoWindow()

            // Returns false, so the default behavior is still used.
            false
        }
        map.setOnMapClickListener { closeContainerInfoWindow() }

        // Prompt the user for permission.
        getLocationPermission()

        // Turn on the My Location layer.
        updateLocationUI()

        // Adds the containers in the map.
        updateContainersUI()

        // Get user container bookmarks.
        getUserContainerBookmarks()

        // Get the current location of the device and set the position of the map.
        focusMyLocation()
    }

    /**
     * Asks the user for the device location permission.
     * The result of the permission request is handled by the onRequestPermissionsResult callback.
     */
    private fun getLocationPermission() {
        when (PackageManager.PERMISSION_GRANTED) {
            ContextCompat.checkSelfPermission(
                applicationContext,
                Manifest.permission.ACCESS_FINE_LOCATION
            ) -> {
                locationPermissionGranted = true
            }

            ContextCompat.checkSelfPermission(
                applicationContext,
                Manifest.permission.ACCESS_COARSE_LOCATION
            ) -> {
                locationPermissionGranted = true
            }

            else -> {
                ActivityCompat.requestPermissions(
                    this,
                    arrayOf(
                        Manifest.permission.ACCESS_COARSE_LOCATION,
                        Manifest.permission.ACCESS_FINE_LOCATION
                    ),
                    PERMISSIONS_REQUEST_ACCESS_LOCATION
                )
            }
        }
    }

    /**
     * Function called when the user responds to the permissions request.
     */
    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray
    ) {
        locationPermissionGranted = false

        when (requestCode) {
            PERMISSIONS_REQUEST_ACCESS_LOCATION -> {
                // If request is cancelled, the result arrays are empty.
                if (grantResults.isNotEmpty() &&
                    grantResults[0] == PackageManager.PERMISSION_GRANTED
                ) {
                    // Location was successfully granted.
                    locationPermissionGranted = true
                }
            }

            else -> super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        }

        updateLocationUI()
        focusMyLocation()
    }

    /**
     * Updates the map's UI settings based on whether the user has granted location permission.
     * If the location permission is granted, enable the Google Map My Location layer.
     */
    private fun updateLocationUI() {
        try {
            // Enable/disable the My Location layer based on the location permission.
            map.isMyLocationEnabled = locationPermissionGranted

            // Disable the default My Location button because buttonMyLocation already performs the
            // same function.
            map.uiSettings.isMyLocationButtonEnabled = false
        } catch (e: SecurityException) {
            Log.e(LOG_TAG, e.message, e)
        }
    }

    /**
     * Updates the map UI by adding the containers as markers using the provided filter.
     */
    private fun updateContainersUI(containerCategoryFilter: ContainerCategory? = null) {
        // Clear the current markers.
        containerClusterManager.clearItems()

        // Map containing the filtered containers, merging those that are in the same position to be
        // contained in the same marker.
        val filteredContainers = mutableMapOf<LatLng, ContainerMarker>()

        // Helper function to handle a successful response.
        val handleSuccess = fun(paginatedContainers: ContainersPaginated) {
            if (isFinishing || isDestroyed) {
                return
            }

            for (container in paginatedContainers.containers) {
                val containerCoordinates = container.geoJSON.geometry.coordinates
                val containerPosition = LatLng(containerCoordinates[1], containerCoordinates[0])

                // Add the marker if it is not currently in the Cluster Manager, otherwise append
                // the container category to the existing marker.
                val existingContainer = filteredContainers[containerPosition]
                if (existingContainer == null) {
                    val containerMarker = ContainerMarker(
                        this,
                        arrayListOf(container)
                    )

                    containerClusterManager.addItem(containerMarker)
                    filteredContainers[containerPosition] = containerMarker
                } else {
                    existingContainer.containers.add(container)
                }
            }

            // Force a re-cluster on the map.
            containerClusterManager.cluster()
        }

        // Execute the request to get all existing containers and mark them in the map.
        val request = ApiClient.listContainers(
            containerCategoryFilter,
            REQUEST_LIST_CONTAINER_LIMIT,
            0,
            token,
            { paginatedContainers ->
                val remainingRequest = paginatedContainers.total / REQUEST_LIST_CONTAINER_LIMIT
                for (i in 1..remainingRequest) {
                    ApiRequestQueue.getInstance(applicationContext).add(
                        ApiClient.listContainers(
                            containerCategoryFilter,
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
            { error -> Common.handleVolleyError(this, this, error) })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }

    /**
     * Update the Google Map camera to focus on the user last-known location.
     */
    private fun focusMyLocation() {
        if (!locationPermissionGranted) {
            return
        }

        try {
            val locationResult = fusedLocationProviderClient.lastLocation
            locationResult.addOnCompleteListener(this) { task ->
                if (task.isSuccessful) {
                    // Set the map's camera position to the current location of the device.
                    val lastKnownLocation = task.result
                    if (lastKnownLocation != null) {
                        map.animateCamera(
                            CameraUpdateFactory.newLatLngZoom(
                                LatLng(
                                    lastKnownLocation.latitude,
                                    lastKnownLocation.longitude
                                ), MAP_CAMERA_ZOOM_DEFAULT
                            )
                        )
                    }
                } else {
                    Log.e(LOG_TAG, task.exception?.message, task.exception)
                }
            }
        } catch (e: SecurityException) {
            Log.e(LOG_TAG, e.message, e)
        }
    }

    /**
     * Opens the info window for the given container marker.
     * It hides the main buttons and shows the info window with the given container marker data.
     */
    private fun showContainerInfoWindow(containerMarker: ContainerMarker) {
        groupButtonsView.visibility = View.GONE
        containerInfoWindowView.visibility = View.VISIBLE
        containerInfoWindowTitleText.text = containerMarker.geoJSON.properties.municipalityName
        containerInfoWindowSnippetText.text = containerMarker.geoJSON.properties.getWayName(this)

        containerInfoWindowBookmarkButton.setImageResource(R.drawable.bookmark)
        for (container in containerMarker.containers) {
            if (userContainerBookmarks.contains(container.id)) {
                // Container is bookmarked.
                containerInfoWindowBookmarkButton.setImageResource(R.drawable.bookmark_fill)
                break
            }
        }
        containerInfoWindowBookmarkButton.setOnClickListener {
            var bookmarked = false
            for (container in containerMarker.containers) {
                if (userContainerBookmarks.contains(container.id)) {
                    bookmarked = true
                    break
                }
            }

            if (bookmarked) {
                containerInfoWindowBookmarkButton.setImageResource(R.drawable.bookmark)

                for (container in containerMarker.containers) {
                    userContainerBookmarks.remove(container.id)

                    val request = ApiClient.removeUserContainerBookmark(userID, container.id, token,
                        {}, {})

                    ApiRequestQueue.getInstance(applicationContext).add(request)
                }
            } else {
                containerInfoWindowBookmarkButton.setImageResource(R.drawable.bookmark_fill)

                for (container in containerMarker.containers) {
                    userContainerBookmarks[container.id] = container

                    val request = ApiClient.createUserContainerBookmark(userID, container.id, token,
                        {}, {})

                    ApiRequestQueue.getInstance(applicationContext).add(request)
                }
            }
        }

        val containerCategoriesDataSet =
            ArrayList<ContainerCategoryRecyclerViewData>(containerMarker.containers.size)
        for (container in containerMarker.containers) {
            val data = ContainerCategoryRecyclerViewData(
                container.category.getIconResource(),
                container.category.getStringResource(this)
            )
            if (containerCategoriesDataSet.contains(data)) {
                // The category already exists in the current data set.
                continue
            }

            containerCategoriesDataSet.add(data)
        }
        containerInfoWindowRecyclerCategories.adapter =
            ContainerCategoriesRecyclerViewAdapter(containerCategoriesDataSet.toTypedArray())

        containerInfoWindowDirectionsButton.setOnClickListener {
            val intentMapDirections = Intent(Intent.ACTION_VIEW)
            intentMapDirections.data =
                Uri.parse("geo:0,0?q=${containerMarker.position.latitude},${containerMarker.position.longitude}(${containerMarker.snippet})")

            if (intentMapDirections.resolveActivity(packageManager) != null) {
                // Start activity only if there is an app that can resolve it.
                startActivity(intentMapDirections)
            }
        }
    }

    /**
     * Closes the container info window.
     * It makes the main buttons visible and hides the info window.
     */
    private fun closeContainerInfoWindow() {
        if (!containerInfoWindowView.isVisible) {
            return
        }

        groupButtonsView.visibility = View.VISIBLE
        containerInfoWindowView.visibility = View.GONE
    }

    /**
     * Gets the current list of containers that the user has bookmarked.
     */
    private fun getUserContainerBookmarks() {
        userContainerBookmarks.clear()

        // Helper function to handle a successful response.
        val handleSuccess = fun(paginatedContainers: ContainersPaginated) {
            if (isFinishing || isDestroyed) {
                return
            }

            for (container in paginatedContainers.containers) {
                userContainerBookmarks[container.id] = container
            }
        }

        // Execute the request to get all existing user container bookmarks and add them to the list.
        val request = ApiClient.listUserContainerBookmarks(
            userID,
            REQUEST_LIST_CONTAINER_LIMIT,
            0,
            token,
            { paginatedContainers ->
                val remainingRequest = paginatedContainers.total / REQUEST_LIST_CONTAINER_LIMIT
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
            { error -> Common.handleVolleyError(this, this, error) })

        ApiRequestQueue.getInstance(applicationContext).add(request)
    }

    companion object {
        private val LOG_TAG = MainActivity::class.java.simpleName

        private const val PERMISSIONS_REQUEST_ACCESS_LOCATION = 1

        private const val MAP_CAMERA_BOUND_SW_LAT = 38.0
        private const val MAP_CAMERA_BOUND_SW_LNG = -10.0
        private const val MAP_CAMERA_BOUND_NE_LAT = 41.0
        private const val MAP_CAMERA_BOUND_NE_LNG = -6.0

        private const val MAP_PADDING_LEFT = 16
        private const val MAP_PADDING_TOP = 144
        private const val MAP_PADDING_RIGHT = 16
        private const val MAP_PADDING_BOTTOM = 32

        private const val MAP_CAMERA_ZOOM_DEFAULT = 15.0F

        private const val REQUEST_LIST_CONTAINER_LIMIT = 100

        private const val CONTAINER_INFO_WINDOW_RECYCLER_CATEGORIES_SPAN_COUNT = 2
    }
}
