<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import OlMap from "ol/Map";
	import Modal from "../clients/Modal.svelte";
	import Map from "./map/Map.svelte";
	import VectorLayer from "ol/layer/Vector";
	import { Feature, MapBrowserEvent } from "ol";
	import { Point } from "ol/geom";
	import VectorSource from "ol/source/Vector";
	import Button from "./Button.svelte";
	import type { Coordinate } from "ol/coordinate";
	import {
		convertToMapProjection,
		convertToResourceProjection,
	} from "../utils/map";
	import { t } from "../utils/i8n";
	import ecomapHttpClient from "../clients/ecomap/http";
	import { getLocationName } from "../utils/location";
	import {
		DEFAULT_ANIMATION_DURATION,
		DEFAULT_MAX_ZOOM,
		DEFAULT_PIN_ICON_SRC,
	} from "../constants/map";

	/**
	 * Indicates if the modal is open.
	 */
	export let open: boolean;

	/**
	 * Callback fired when open state of the modal changes.
	 * @param open New open state modal.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Callback fired when save action is triggered.
	 * @param coordinate Container coordinate.
	 * @param locationName Container location name.
	 */
	export let onSave: (coordinate: Coordinate, locationName: string) => void;

	/**
	 * Coordinate to display in the map.
	 * @default null
	 */
	export let coordinate: Coordinate | null = null;

	/**
	 * Callback fired when cancel action is triggered.
	 * @default null
	 */
	export let onCancel: (() => void) | null = null;

	/**
	 * Map layer style.
	 * @default "/images/pin.svg"
	 */
	export let iconSrc: string = DEFAULT_PIN_ICON_SRC;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * ID of the selected location.
	 */
	const featureId = "location";

	/**
	 * The map layer which displays the selected location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": iconSrc,
		},
	});

	/**
	 * Indicates if save action is disabled.
	 * @default true
	 */
	let isSaveActionDisabled = true;

	/**
	 * Removes the selected location from the map.
	 */
	function removeSelectedLocation() {
		isSaveActionDisabled = true;

		const source = layer.getSource();
		if (!source) {
			return;
		}

		source.clear();
		map.removeLayer(layer);
	}

	/**
	 * Adds a pin with the selected location to the map.
	 * @param coordinate Coordinate where pin will be placed.
	 * @param mode Map view mode.
	 */
	function addSelectedLocation(
		coordinate: Coordinate,
		mode: "fit" | "animate",
	) {
		const source = layer.getSource();
		if (!source) {
			return;
		}

		source.clear();
		map.removeLayer(layer);

		const point = new Point(coordinate);
		const feature = new Feature(point);
		feature.setId(featureId);

		source.addFeature(feature);
		map.addLayer(layer);

		const view = map.getView();
		if (mode === "fit") {
			view.fit(point);
		} else {
			view.animate({
				center: coordinate,
				duration: DEFAULT_ANIMATION_DURATION,
				zoom: DEFAULT_MAX_ZOOM,
			});
		}

		isSaveActionDisabled = false;
	}

	/**
	 * Retrieves the location of a resource given a map coordinate.
	 * @param coordinate Map coordinate.
	 */
	async function getLocation(coordinate: Coordinate) {
		const resourceCoordinate = convertToResourceProjection(coordinate);

		const wayNamePromise = ecomapHttpClient.GET("/ways/reverse-geocoding", {
			params: {
				query: {
					coordinates: resourceCoordinate,
				},
			},
		});

		const municipalityNamePromise = ecomapHttpClient.GET(
			"/municipalities/reverse-geocoding",
			{
				params: {
					query: {
						coordinates: resourceCoordinate,
					},
				},
			},
		);

		const [wayNameRes, municipalityNameRes] = await Promise.allSettled([
			wayNamePromise,
			municipalityNamePromise,
		]);

		let wayName: string | undefined;
		let municipalityName: string | undefined;

		if (wayNameRes.status === "fulfilled") {
			const wayRes = wayNameRes.value;

			if (!wayRes.error) {
				wayName = wayRes.data.osmName;
			}
		}

		if (municipalityNameRes.status === "fulfilled") {
			const municipalityRes = municipalityNameRes.value;

			if (!municipalityRes.error) {
				municipalityName = municipalityRes.data.name;
			}
		}

		return { wayName, municipalityName };
	}

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {
		addSelectedLocation(e.coordinate, "animate");
	}

	/**
	 * Handles the save action of the modal.
	 */
	async function handleSave() {
		const source = layer.getSource();
		if (!source) {
			return;
		}

		const feature = source.getFeatureById(featureId);
		if (!feature) {
			return;
		}

		const point = feature.getGeometry();
		if (!point) {
			return;
		}

		const coordinates = point.getCoordinates();
		const { wayName, municipalityName } = await getLocation(coordinates);

		const locationName = getLocationName(wayName, municipalityName);
		onSave(coordinates, locationName);

		onOpenChange(false);
	}

	/**
	 * Handles cancel action.
	 */
	function handleCancel() {
		onOpenChange(false);
		onCancel?.();
	}

	onMount(() => {
		map.on("click", handleMapClick);
	});

	onDestroy(() => {
		map.un("click", handleMapClick);
	});

	// Adds/removes selected location when modal is opened.
	$: if (open) {
		if (coordinate) {
			const mapCoordinate = convertToMapProjection(coordinate);
			addSelectedLocation(mapCoordinate, "fit");
		} else {
			removeSelectedLocation();
		}
	}
</script>

<Modal
	{open}
	{onOpenChange}
	onClickOutside={removeSelectedLocation}
	title={$t("selectLocation")}
>
	<Map bind:map mapId="select-location-map" --height="32rem" --width="60rem" />
	<svelte:fragment slot="actions">
		<Button variant="tertiary" onClick={handleCancel}>{$t("cancel")}</Button>
		<Button
			startIcon="check"
			disabled={isSaveActionDisabled}
			onClick={handleSave}
		>
			{$t("save")}
		</Button>
	</svelte:fragment>
</Modal>
