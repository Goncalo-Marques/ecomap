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
	import { convertToMapProjection } from "../utils/map";
	import { t } from "../utils/i8n";

	/**
	 * Indicates if the modal is open.
	 */
	export let open: boolean;

	/**
	 * Callback fired when open state of the modal changes.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Callback fired when save action is triggered.
	 */
	export let onSave: (coordinate: Coordinate) => void;

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
			"icon-src": "/images/pin.svg",
		},
	});

	/**
	 * Indicates if save action is disabled.
	 * @default true
	 */
	let disabled = true;

	/**
	 * Removes the selected location from the map.
	 */
	function removeSelectedLocation() {
		disabled = true;

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
	 */
	function addSelectedLocation(coordinate: Coordinate) {
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

		disabled = false;
	}

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {
		addSelectedLocation(e.coordinate);
	}

	/**
	 * Handles the save action of the modal.
	 */
	function handleSave() {
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
		onSave(coordinates);
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
			addSelectedLocation(mapCoordinate);
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
		<Button startIcon="check" {disabled} onClick={handleSave}>
			{$t("save")}
		</Button>
	</svelte:fragment>
</Modal>
