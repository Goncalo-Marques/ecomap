package com.ecomap.ecomap

import android.Manifest
import android.content.pm.PackageManager
import android.content.Intent
import android.os.Bundle
import android.util.Log
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.GoogleMap
import com.google.android.gms.maps.GoogleMapOptions
import com.google.android.gms.maps.OnMapReadyCallback
import com.google.android.gms.maps.SupportMapFragment
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.LatLngBounds
import com.google.android.gms.maps.model.MarkerOptions
import com.google.android.material.floatingactionbutton.FloatingActionButton
import com.ecomap.ecomap.data.UserStore
import com.ecomap.ecomap.signin.SignInActivity
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class MainActivity : AppCompatActivity(), OnMapReadyCallback {
    /**
     * Defines the Google Map instance.
     * It is set when the map is ready.
     */
    private var map: GoogleMap? = null

    /**
     * The entry point to the Fused Location Provider.
     */
    private lateinit var fusedLocationProviderClient: FusedLocationProviderClient

    /**
     * Defines whether the location permission is granted.
     */
    private var locationPermissionGranted = false

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
        this.map = googleMap

        // Prompt the user for permission.
        getLocationPermission()

        // Turn on the My Location layer.
        updateLocationUI()

        // Get the current location of the device and set the position of the map.
        focusMyLocation()

        // TODO: Add the containers using the server.
        googleMap.addMarker(
            MarkerOptions()
                .position(LatLng(40.0, -9.0))
                .title("Marker")
        )
    }

    /**
     * Asks the user for the device location permission.
     * The result of the permission request is handled by the onRequestPermissionsResult callback.
     */
    private fun getLocationPermission() {
        when (PackageManager.PERMISSION_GRANTED) {
            ContextCompat.checkSelfPermission(
                this.applicationContext,
                Manifest.permission.ACCESS_FINE_LOCATION
            ) -> {
                locationPermissionGranted = true
            }

            ContextCompat.checkSelfPermission(
                this.applicationContext,
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
    private fun updateLocationUI(enableMyLocationLayer: Boolean = locationPermissionGranted) {
        if (map == null) {
            return
        }

        try {
            // Enable/disable the My Location layer based on the location permission.
            map?.isMyLocationEnabled = locationPermissionGranted

            // Disable the default My Location button because buttonMyLocation already performs the
            // same function.
            map?.uiSettings?.isMyLocationButtonEnabled = false
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
                        map?.animateCamera(
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

    companion object {
        private val LOG_TAG = MainActivity::class.java.simpleName

        private const val PERMISSIONS_REQUEST_ACCESS_LOCATION = 1

        private const val MAP_CAMERA_ZOOM_DEFAULT = 15.0

        private const val MAP_CAMERA_BOUND_SW_LAT = 38.0
        private const val MAP_CAMERA_BOUND_SW_LNG = -10.0
        private const val MAP_CAMERA_BOUND_NE_LAT = 41.0
        private const val MAP_CAMERA_BOUND_NE_LNG = -6.0
    }
}
