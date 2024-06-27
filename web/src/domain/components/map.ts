/**
 * Represents the options for creating a map.
 */
export interface CreateMapOptions {
	/**
	 * The initial longitude where map is centered.
	 */
	lon: number;

	/**
	 * The initial latitude where map is centered.
	 */
	lat: number;

	/**
	 * The initial map zoom.
	 */
	zoom: number;

	/**
	 * The minimum map zoom.
	 */
	minZoom?: number;

	/**
	 * The maximum map zoon.
	 */
	maxZoom?: number;

	/**
	 * The projection system used for the map.
	 */
	projection?: string;
}

/**
 * Represents the options of the point layer of `MapHelper`.
 */
export interface MapHelperPointLayerOptions {
	/**
	 * The name of the layer.
	 */
	layerName?: string;

	/**
	 * The color of the layer.
	 */
	layerColor?: string;

	/**
	 * The source of the icon to be displayed for each feature of the layer.
	 */
	iconSrc?: string;
}

/**
 * Represents the options of the cluster layer of `MapHelper`.
 */
export interface MapHelperClusterLayerOptions {
	/**
	 * The name of the layer.
	 */
	layerName?: string;

	/**
	 * The color of the layer.
	 */
	layerColor?: string;

	/**
	 * The source of the icon to be displayed for each feature of the layer.
	 */
	iconSrc?: string;

	/**
	 * The source of the selected icon to be displayed for each selected feature of the layer.
	 */
	selectedIconSrc?: string;

	/**
	 * The color of the cluster border.
	 */
	clusterBorderColor?: string;
}
