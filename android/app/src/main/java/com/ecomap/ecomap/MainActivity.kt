package com.ecomap.ecomap

import android.Manifest
import android.content.pm.PackageManager
import android.location.Location
import android.os.Bundle
import android.util.Log
import android.widget.Button
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
import com.google.android.libraries.places.api.Places
import com.google.android.libraries.places.api.net.PlacesClient

class MainActivity : AppCompatActivity(), OnMapReadyCallback {
    private var map: GoogleMap? = null

    private lateinit var mapLatLngBound: LatLngBounds

    // The entry point to the Places API.
    private lateinit var placesClient: PlacesClient

    // The entry point to the Fused Location Provider.
    private lateinit var fusedLocationProviderClient: FusedLocationProviderClient

    // Set Lisbon as the default location.
    // TODO: Check if coordinates are correct.
    private val defaultLocation = LatLng(38.72158286729078, -9.13913643581083)
    private var locationPermissionGranted = false
    private var lastKnownLocation: Location? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_main)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // Construct a place client.
        // For more information, see: https://developers.google.com/maps/documentation/places/android-sdk/overview.
        Places.initialize(applicationContext, BuildConfig.MAPS_API_KEY)
        placesClient = Places.createClient(this)

        // Construct the main entry point for Android location services.
        fusedLocationProviderClient = LocationServices.getFusedLocationProviderClient(this)

        // Google Map configurations.
        mapLatLngBound = LatLngBounds(
            LatLng(
                resources.getInteger(R.integer.camera_bound_southwest_latitude).toDouble(),
                resources.getInteger(R.integer.camera_bound_southwest_longitude).toDouble()
            ),
            LatLng(
                resources.getInteger(R.integer.camera_bound_northeast_latitude).toDouble(),
                resources.getInteger(R.integer.camera_bound_northeast_longitude).toDouble()
            )
        )

        val button: Button = findViewById(R.id.button)
        button.setOnClickListener { getDeviceLocation() }

        val googleMapOptions = GoogleMapOptions()
        googleMapOptions
            .mapType(GoogleMap.MAP_TYPE_NORMAL)
            .latLngBoundsForCameraTarget(mapLatLngBound)

        // Add support map fragment to the map container.
        val mapFragment = SupportMapFragment.newInstance(googleMapOptions)
        supportFragmentManager
            .beginTransaction()
            .add(R.id.fragment_container_view_map, mapFragment)
            .commit()

        // Register the map callback.
        mapFragment.getMapAsync(this)
    }

    override fun onMapReady(googleMap: GoogleMap) {
        this.map = googleMap

        // Prompt the user for permission.
        getLocationPermission()

        // Turn on the My Location layer and the related control on the map.
        updateLocationUI()

        // Get the current location of the device and set the position of the map.
        getDeviceLocation()

        googleMap.moveCamera(
            CameraUpdateFactory.newLatLngBounds(mapLatLngBound, 0)
        )

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
        when {
            ContextCompat.checkSelfPermission(
                this.applicationContext,
                Manifest.permission.ACCESS_FINE_LOCATION
            )
                    == PackageManager.PERMISSION_GRANTED -> {
                locationPermissionGranted = true
            }

            ContextCompat.checkSelfPermission(
                this.applicationContext,
                Manifest.permission.ACCESS_COARSE_LOCATION
            )
                    == PackageManager.PERMISSION_GRANTED -> {
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
                    locationPermissionGranted = true
                }
            }

            else -> super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        }

        updateLocationUI()
    }

    /**
     * Updates the map's UI settings based on whether the user has granted location permission.
     */
    private fun updateLocationUI() {
        if (map == null) {
            return
        }

        try {
            if (locationPermissionGranted) {
                map?.isMyLocationEnabled = true
//                map?.uiSettings?.isMyLocationButtonEnabled = true
            } else {
                map?.isMyLocationEnabled = false
//                map?.uiSettings?.isMyLocationButtonEnabled = false
                lastKnownLocation = null
//                getLocationPermission()
            }
        } catch (e: SecurityException) {
            Log.e("Exception: %s", e.message, e)
        }
    }

    /**
     * Use the fused location provider to find the device's last-known location, then use that location to position the map.
     */
    private fun getDeviceLocation() {
        /*
         * Get the best and most recent location of the device, which may be null in rare
         * cases when a location is not available.
         */
        try {
            if (locationPermissionGranted) {
                val locationResult = fusedLocationProviderClient.lastLocation
                locationResult.addOnCompleteListener(this) { task ->
                    if (task.isSuccessful) {
                        // Set the map's camera position to the current location of the device.
                        lastKnownLocation = task.result
                        if (lastKnownLocation != null) {
                            map?.animateCamera(
                                CameraUpdateFactory.newLatLngZoom(
                                    LatLng(
                                        lastKnownLocation!!.latitude,
                                        lastKnownLocation!!.longitude
                                    ), DEFAULT_ZOOM.toFloat()
                                )
                            )
                        }
                    } else {
                        Log.d(TAG, "Current location is null. Using defaults.")
                        Log.e(TAG, "Exception: %s", task.exception)
                        map?.animateCamera(
                            CameraUpdateFactory
                                .newLatLngZoom(defaultLocation, DEFAULT_ZOOM.toFloat())
                        )
                        map?.uiSettings?.isMyLocationButtonEnabled = false
                    }
                }
            }
        } catch (e: SecurityException) {
            Log.e("Exception: %s", e.message, e)
        }
    }

    companion object {
        private val TAG = MainActivity::class.java.simpleName
        private const val DEFAULT_ZOOM = 15
        private const val PERMISSIONS_REQUEST_ACCESS_LOCATION = 1
    }
}
