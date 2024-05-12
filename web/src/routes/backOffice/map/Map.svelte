<script lang="ts">
	import { onMount } from "svelte";
	import OlMap from "ol/Map";
	import Map from "../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../lib/components/map/mapUtils";
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import { Point, SimpleGeometry } from "ol/geom";
	import { Feature, MapBrowserEvent } from "ol";
	import { convertToMapProjection } from "../../../lib/utils/map";
	import { t } from "../../../lib/utils/i8n";
	import {
		CONTAINER_ICON_SRC,
		SELECTED_CONTAINER_ICON_SRC,
		SELECTED_TRUCK_ICON_SRC,
		SELECTED_WAREHOUSE_ICON_SRC,
		TRUCK_ICON_SRC,
		WAREHOUSE_ICON_SRC,
	} from "../../../lib/constants/map";
	import type { Container } from "../../../domain/container";
	import type { Truck } from "../../../domain/truck";
	import type { Warehouse } from "../../../domain/warehouse";
	import ContainerBottomSheet from "../components/ContainerBottomSheet.svelte";
	import TruckBottomSheet from "../components/TruckBottomSheet.svelte";
	import WarehouseBottomSheet from "../components/WarehouseBottomSheet.svelte";
	import { getCssVariable } from "../../../lib/utils/cssVars";
	import type { Coordinate } from "ol/coordinate";
	import { type Extent, boundingExtent } from "ol/extent";
	import ResourceGroupBottomSheet from "../components/ResourceGroupBottomSheet.svelte";
	import type { ResourceGroupLocation } from "../../../domain/map";

	/**
	 * The Open Layers map.
	 */
	let map: OlMap;

	/**
	 * The container selected in the map.
	 */
	let selectedContainer: Container | null = null;

	/**
	 * The truck selected in the map.
	 */
	let selectedTruck: Truck | null = null;

	/**
	 * The warehouse selected in the map.
	 */
	let selectedWarehouse: Warehouse | null = null;

	/**
	 * The feature group selected in the map.
	 */
	let selectedGroup: ResourceGroupLocation | null = null;

	/**
	 * The selected feature in the map.
	 */
	let selectedFeature: Feature | null = null;

	/**
	 * Retrieves truck features to display in the map.
	 */
	async function getTruckFeatures(): Promise<Feature<Point>[]> {
		const res = await ecomapHttpClient.GET("/trucks");

		const features: Feature<Point>[] = [];

		if (res.error) {
			return features;
		}

		for (const truck of res.data.trucks) {
			const {
				id,
				make,
				model,
				licensePlate,
				personCapacity,
				createdAt,
				modifiedAt,
				geoJson,
			} = truck;
			const transformedCoordinate = convertToMapProjection(
				geoJson.geometry.coordinates,
			);
			const point = new Point(transformedCoordinate);
			const feature = new Feature(point);
			feature.setId(id);
			feature.setProperties({
				type: "truck",
				id,
				make,
				model,
				licensePlate,
				personCapacity,
				createdAt,
				modifiedAt,
				geoJson,
			});

			features.push(feature);
		}

		return features;
	}

	/**
	 * Retrieves warehouse features to display in the map.
	 */
	async function getWarehouseFeatures(): Promise<Feature<Point>[]> {
		const res = await ecomapHttpClient.GET("/warehouses");

		const features: Feature<Point>[] = [];

		if (res.error) {
			return features;
		}

		for (const warehouse of res.data.warehouses) {
			const { id, truckCapacity, createdAt, modifiedAt, geoJson } = warehouse;
			const transformedCoordinate = convertToMapProjection(
				geoJson.geometry.coordinates,
			);
			const point = new Point(transformedCoordinate);
			const feature = new Feature(point);
			feature.setId(id);
			feature.setProperties({
				type: "warehouse",
				id,
				truckCapacity,
				createdAt,
				modifiedAt,
				geoJson,
			});

			features.push(feature);
		}

		return features;
	}

	/**
	 * Retrieves container features to display in the map.
	 */
	async function getContainerFeatures(): Promise<Feature<Point>[]> {
		const res = await ecomapHttpClient.GET("/containers");

		const features: Feature<Point>[] = [];

		if (res.error) {
			return features;
		}

		for (const container of res.data.containers) {
			const { id, category, createdAt, modifiedAt, geoJson } = container;
			const transformedCoordinate = convertToMapProjection(
				geoJson.geometry.coordinates,
			);
			const point = new Point(transformedCoordinate);
			const feature = new Feature(point);
			feature.setId(id);
			feature.setProperties({
				type: "container",
				id,
				category,
				createdAt,
				modifiedAt,
				geoJson,
			});

			features.push(feature);
		}

		return features;
	}

	/**
	 * Retrieves a container from a container feature.
	 * @param feature Container feature.
	 */
	function getContainerFromFeature(feature: Feature): Container {
		const { id, category, createdAt, modifiedAt, geoJson } =
			feature.getProperties();

		return {
			id,
			category,
			createdAt,
			modifiedAt,
			geoJson,
		};
	}

	/**
	 * Retrieves a truck from a container feature.
	 * @param feature Truck feature.
	 */
	function getTruckFromFeature(feature: Feature): Truck {
		const {
			id,
			make,
			model,
			licensePlate,
			personCapacity,
			createdAt,
			modifiedAt,
			geoJson,
		} = feature.getProperties();

		return {
			id,
			make,
			model,
			licensePlate,
			personCapacity,
			createdAt,
			modifiedAt,
			geoJson,
		};
	}

	/**
	 * Retrieves a warehouse from a container feature.
	 * @param feature Warehouse feature.
	 */
	function getWarehouseFromFeature(feature: Feature): Warehouse {
		const { id, truckCapacity, createdAt, modifiedAt, geoJson } =
			feature.getProperties();

		return {
			id,
			truckCapacity,
			createdAt,
			modifiedAt,
			geoJson,
		};
	}

	/**
	 * Loads all features into the map.
	 * @param mapHelper Map helper to add features to the map.
	 */
	async function loadFeatures(mapHelper: MapHelper) {
		const [containerFeaturesRes, truckFeaturesRes, warehouseFeaturesRes] =
			await Promise.allSettled([
				getContainerFeatures(),
				getTruckFeatures(),
				getWarehouseFeatures(),
			]);

		if (containerFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(containerFeaturesRes.value, {
				layerName: $t("containers"),
				layerColor: getCssVariable("--green-700"),
				iconSrc: CONTAINER_ICON_SRC,
				selectedIconSrc: SELECTED_CONTAINER_ICON_SRC,
			});
		}

		if (truckFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(truckFeaturesRes.value, {
				layerName: $t("trucks"),
				layerColor: getCssVariable("--cyan-900"),
				iconSrc: TRUCK_ICON_SRC,
				selectedIconSrc: SELECTED_TRUCK_ICON_SRC,
			});
		}

		if (warehouseFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(warehouseFeaturesRes.value, {
				layerName: $t("warehouses"),
				layerColor: getCssVariable("--amber-800"),
				iconSrc: WAREHOUSE_ICON_SRC,
				selectedIconSrc: SELECTED_WAREHOUSE_ICON_SRC,
			});
		}
	}

	/**
	 * Retrieves the extent from a list of features.
	 * @param features Features to retrieve the extent from.
	 */
	function getExtentFromSelectedFeatures(features: Feature[]): Extent {
		const coordinates: Coordinate[] = [];

		for (const feature of features) {
			const geom = feature.getGeometry();
			if (!(geom instanceof SimpleGeometry)) {
				continue;
			}

			const coord = geom.getCoordinates();
			if (!coord) {
				continue;
			}

			coordinates.push(coord);
		}

		return boundingExtent(coordinates);
	}

	/**
	 * Selects a single feature.
	 * @param feature Feature.
	 */
	function selectFeature(feature: Feature) {
		// Mark feature as selected.
		feature.set("selected", true);

		// Keep track of the selected feature to reset it on the next map click event.
		selectedFeature = feature;

		// Determine which type of feature it is.
		switch (feature.get("type")) {
			case "container":
				selectedContainer = getContainerFromFeature(feature);
				break;

			case "truck":
				selectedTruck = getTruckFromFeature(feature);
				break;

			case "warehouse":
				selectedWarehouse = getWarehouseFromFeature(feature);
				break;
		}
	}

	/**
	 * Selects a group of features.
	 * @param features Features.
	 */
	function selectGroupedFeatures(features: Feature[]) {
		selectedGroup = {
			wayName: features[0].get("geoJson").properties.wayName,
			municipalityName: features[0].get("geoJson").properties.municipalityName,
			containers: [],
			trucks: [],
			warehouses: [],
		};

		for (const feature of features) {
			switch (feature.get("type")) {
				case "container":
					selectedGroup.containers.push(getContainerFromFeature(feature));
					break;

				case "truck":
					selectedGroup.trucks.push(getTruckFromFeature(feature));
					break;

				case "warehouse":
					selectedGroup.warehouses.push(getWarehouseFromFeature(feature));
					break;
			}
		}
	}

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {
		// Reset all previously selected features.
		selectedContainer = null;
		selectedTruck = null;
		selectedWarehouse = null;
		selectedGroup = null;
		selectedFeature?.set("selected", false);

		// Get features that were clicked.
		const clickedFeatures = map.getFeaturesAtPixel(e.pixel);

		// Check if the click was not performed on any feature.
		if (!clickedFeatures.length) {
			return;
		}

		const mapView = map.getView();
		let extent: Extent;

		const selectedFeatures: Feature[] = clickedFeatures[0].get("features");
		if (selectedFeatures.length > 1) {
			// When the map zoom is at maximum and there are still grouped features, select them.
			if (mapView.getZoom() === mapView.getMaxZoom()) {
				selectGroupedFeatures(selectedFeatures);
				return;
			}

			// Otherwise, retrieve the map extent of the features.
			extent = getExtentFromSelectedFeatures(selectedFeatures);
		} else {
			const feature = selectedFeatures[0];

			const geometry = feature.getGeometry();
			if (!geometry) {
				return;
			}

			// Mark feature as a selected feature.
			selectFeature(feature);

			// Get extent of the feature.
			extent = geometry.getExtent();
		}

		// Zoom map.
		mapView.fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
	}

	onMount(() => {
		const mapHelper = new MapHelper(map);
		loadFeatures(mapHelper);

		map.on("click", handleMapClick);
	});
</script>

<main>
	<Map
		bind:map
		showLayers={!selectedWarehouse &&
			!selectedContainer &&
			!selectedTruck &&
			!selectedGroup}
	/>

	{#if selectedGroup}
		<ResourceGroupBottomSheet group={selectedGroup} />
	{:else if selectedContainer}
		<ContainerBottomSheet container={selectedContainer} />
	{:else if selectedTruck}
		<TruckBottomSheet truck={selectedTruck} />
	{:else if selectedWarehouse}
		<WarehouseBottomSheet warehouse={selectedWarehouse} />
	{/if}
</main>

<style>
	main {
		position: relative;
		width: 100%;
	}
</style>
