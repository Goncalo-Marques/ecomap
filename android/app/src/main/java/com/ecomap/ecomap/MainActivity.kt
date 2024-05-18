package com.ecomap.ecomap

import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Bundle
import android.util.Log
import android.widget.Toast
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.android.volley.VolleyError
import com.ecomap.ecomap.clients.ecomap.http.ApiClient
import com.ecomap.ecomap.clients.ecomap.http.ApiRequestQueue
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.domain.ContainerCategory
import com.ecomap.ecomap.domain.ContainersPaginated
import com.ecomap.ecomap.signin.SignInActivity
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.GoogleMap
import com.google.android.gms.maps.GoogleMapOptions
import com.google.android.gms.maps.OnMapReadyCallback
import com.google.android.gms.maps.SupportMapFragment
import com.google.android.gms.maps.model.BitmapDescriptorFactory
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.LatLngBounds
import com.google.android.gms.maps.model.MarkerOptions
import com.google.android.material.floatingactionbutton.FloatingActionButton
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
     * Defines the authentication token.
     */
    private lateinit var token: String

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

        // Flags the intent to mark the activity as the root in the history stack,
        // clearing out any other tasks.
        intentSignInActivity.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)

        // Check whether the user store contains the login token.
        // If not, start the SignIn Activity.
        val store = UserStore(applicationContext)
        runBlocking {
            val storeToken = store.getToken().first()
            if (storeToken == null) {
                startActivity(intentSignInActivity)
            }

            token = storeToken.toString()
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

        // Add support map fragment to the map container.
        val mapFragment = SupportMapFragment.newInstance(googleMapOptions)
        supportFragmentManager
            .beginTransaction()
            .add(R.id.fragment_container_view_map, mapFragment)
            .commit()

        // Register the map callback.
        mapFragment.getMapAsync(this)

        // Get activity views.
        val buttonMyLocation: FloatingActionButton = findViewById(R.id.button_my_location)

        // Set button functions.
        buttonMyLocation.setOnClickListener { focusMyLocation() }
    }

    /**
     * Function called when the Google Map is ready.
     */
    override fun onMapReady(googleMap: GoogleMap) {
        map = googleMap

        // Prompt the user for permission.
        getLocationPermission()

        // Turn on the My Location layer.
        updateLocationUI()

        // Get the current location of the device and set the position of the map.
        focusMyLocation()

        // Adds the containers in the map.
        updateContainersUI()
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
                                ), MAP_CAMERA_ZOOM_DEFAULT.toFloat()
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
     * Updates the map UI by adding the containers as markers using the provided filter.
     */
    private fun updateContainersUI(containerCategoryFilter: ContainerCategory? = null) {
        if (map == null) {
            return
        }

        // Clear the current markers.
        map?.clear()

        // Helper function to handle a successful response.
        val handleSuccess = fun(paginatedContainers: ContainersPaginated) {
            val markerIcon = BitmapDescriptorFactory.fromResource(R.drawable.marker_icon)
            for (container in paginatedContainers.containers) {
                val containerCoordinates = container.geoJSON.geometry.coordinates
                map?.addMarker(
                    MarkerOptions()
                        .position(LatLng(containerCoordinates[1], containerCoordinates[0]))
                        .icon(markerIcon)
                )
            }
        }

        // Helper function to handle a failed response.
        val handleError = fun(error: VolleyError) {
            val errorResponse = ApiClient.mapError(error)

            var errorMessage = errorResponse.code
            if (errorResponse.message.isNotEmpty()) {
                errorMessage = errorResponse.message
            }

            Toast.makeText(
                this.applicationContext,
                errorMessage,
                Toast.LENGTH_LONG
            ).show()
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
                    ApiRequestQueue.getInstance(this.applicationContext).add(
                        ApiClient.listContainers(
                            containerCategoryFilter,
                            REQUEST_LIST_CONTAINER_LIMIT,
                            REQUEST_LIST_CONTAINER_LIMIT * i,
                            token,
                            { handleSuccess(it) },
                            { handleError(it) }
                        )
                    )
                }

                handleSuccess(paginatedContainers)
            },
            { error -> handleError(error) })

        ApiRequestQueue.getInstance(this.applicationContext).add(request)
    }

    companion object {
        private val LOG_TAG = MainActivity::class.java.simpleName

        private const val PERMISSIONS_REQUEST_ACCESS_LOCATION = 1

        private const val MAP_CAMERA_ZOOM_DEFAULT = 15.0

        private const val MAP_CAMERA_BOUND_SW_LAT = 38.0
        private const val MAP_CAMERA_BOUND_SW_LNG = -10.0
        private const val MAP_CAMERA_BOUND_NE_LAT = 41.0
        private const val MAP_CAMERA_BOUND_NE_LNG = -6.0

        private const val REQUEST_LIST_CONTAINER_LIMIT = 100
    }
}
