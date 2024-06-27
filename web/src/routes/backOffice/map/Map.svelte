<script lang="ts">
	import { onMount } from "svelte";
	import OlMap from "ol/Map";
	import MapComponent from "../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../lib/components/map/mapUtils";
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import { Point, SimpleGeometry } from "ol/geom";
	import { Feature, MapBrowserEvent } from "ol";
	import { convertToMapProjection } from "../../../lib/utils/map";
	import { t } from "../../../lib/utils/i8n";
	import {
		CONTAINER_ICON_SRC,
		DEFAULT_MAX_ZOOM,
		LANDFILL_ICON_SRC,
		SELECTED_CONTAINER_ICON_SRC,
		SELECTED_LANDFILL_ICON_SRC,
		SELECTED_TRUCK_ICON_SRC,
		SELECTED_WAREHOUSE_ICON_SRC,
		TRUCK_ICON_SRC,
		WAREHOUSE_ICON_SRC,
	} from "../../../lib/constants/map";
	import type { Container } from "../../../domain/container";
	import type { Truck } from "../../../domain/truck";
	import type { Warehouse } from "../../../domain/warehouse";
	import { getCssVariable } from "../../../lib/utils/cssVars";
	import type { Coordinate } from "ol/coordinate";
	import { type Extent, boundingExtent } from "ol/extent";
	import MapBottomSheet from "./MapBottomSheet.svelte";
	import { getBatchPaginatedResponse } from "../../../lib/utils/request";
	import Button from "../../../lib/components/Button.svelte";
	import type { Landfill } from "../../../domain/landfill";
	import Spinner from "../../../lib/components/Spinner.svelte";

	/**
	 * The Open Layers map.
	 */
	let map: OlMap;

	/**
	 * The containers selected in the map.
	 */
	let selectedContainers: Container[] = [];

	/**
	 * The truck selected in the map.
	 */
	let selectedTruck: Truck | null = null;

	/**
	 * The warehouse selected in the map.
	 */
	let selectedWarehouse: Warehouse | null = null;

	/**
	 * The landfill selected in the map.
	 */
	let selectedLandfill: Landfill | null = null;

	/**
	 * The selected feature in the map.
	 */
	let selectedFeature: Feature | null = null;

	/**
	 * Indicates whether the map is loading features.
	 */
	let loading: boolean = false;

	/**
	 * Retrieves truck features to display in the map.
	 */
	async function getTruckFeatures(): Promise<Feature<Point>[]> {
		const trucks = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/trucks", {
				params: { query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.trucks };
		});

		const truckFeatures: Feature<Point>[] = [];

		for (const truck of trucks) {
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

			truckFeatures.push(feature);
		}

		return truckFeatures;
	}

	/**
	 * Retrieves warehouse features to display in the map.
	 */
	async function getWarehouseFeatures(): Promise<Feature[]> {
		const warehouses = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/warehouses", {
					params: { query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.warehouses };
			},
		);

		const warehouseFeatures: Feature<Point>[] = [];

		for (const warehouse of warehouses) {
			const { id, truckCapacity, createdAt, modifiedAt, geoJson } = warehouse;
			const transformedCoordinate = convertToMapProjection(
				geoJson.geometry.coordinates,
			);
			const point = new Point(transformedCoordinate);
			const feature = new Feature(point);
			feature.setProperties({
				type: "warehouse",
				id,
				truckCapacity,
				createdAt,
				modifiedAt,
				geoJson,
			});

			warehouseFeatures.push(feature);
		}

		return warehouseFeatures;
	}

	/**
	 * Retrieves container features to display in the map.
	 */
	async function getContainerFeatures(): Promise<Feature<Point>[]> {
		const containers = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/containers", {
					params: { query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.containers };
			},
		);

		const containerMap = new Map<string, Feature<Point>>();

		for (const container of containers) {
			const { id, category, createdAt, modifiedAt, geoJson } = container;

			const key = `${geoJson.geometry.coordinates[0]},${geoJson.geometry.coordinates[1]}`;

			if (!containerMap.has(key)) {
				const transformedCoordinate = convertToMapProjection(
					geoJson.geometry.coordinates,
				);
				const point = new Point(transformedCoordinate);
				const feature = new Feature(point);

				feature.setProperties({ type: "container", geoJson, items: [] });

				containerMap.set(key, feature);
			}

			const containerFeature = containerMap.get(key)!;
			containerFeature.get("items").push({
				id,
				category,
				createdAt,
				modifiedAt,
				geoJson,
			});
		}

		return Array.from(containerMap.values());
	}

	/**
	 * Retrieves landfills features to display in the map.
	 */
	async function getLandfillFeatures(): Promise<Feature[]> {
		const landfills = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/landfills", {
				params: { query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.landfills };
		});

		const landfillFeatures: Feature<Point>[] = [];

		for (const landfill of landfills) {
			const { id, createdAt, modifiedAt, geoJson } = landfill;
			const transformedCoordinate = convertToMapProjection(
				geoJson.geometry.coordinates,
			);
			const point = new Point(transformedCoordinate);
			const feature = new Feature(point);
			feature.setProperties({
				type: "landfill",
				id,
				createdAt,
				modifiedAt,
				geoJson,
			});

			landfillFeatures.push(feature);
		}

		return landfillFeatures;
	}

	/**
	 * Retrieves the containers from a container feature.
	 * @param feature Container feature.
	 * @returns Containers.
	 */
	function getContainersFromFeature(feature: Feature): Container[] {
		const containers: Container[] = feature.get("items");

		return containers;
	}

	/**
	 * Retrieves a truck from a truck feature.
	 * @param feature Truck feature.
	 * @returns Truck.
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
	 * Retrieves a warehouse from a warehouse feature.
	 * @param feature Warehouse feature.
	 * @returns Warehouse.
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
	 * Retrieves a landfill from a landfill feature.
	 * @param feature Landfill feature.
	 * @returns Landfill.
	 */
	function getLandfillFromFeature(feature: Feature): Landfill {
		const { id, createdAt, modifiedAt, geoJson } = feature.getProperties();

		return {
			id,
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
		loading = true;

		const [
			containerFeaturesRes,
			truckFeaturesRes,
			warehouseFeaturesRes,
			landfillFeaturesRes,
		] = await Promise.allSettled([
			getContainerFeatures(),
			getTruckFeatures(),
			getWarehouseFeatures(),
			getLandfillFeatures(),
		]);

		if (containerFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(containerFeaturesRes.value, {
				layerName: $t("containers"),
				layerColor: getCssVariable("--green-700"),
				iconSrc: CONTAINER_ICON_SRC,
				selectedIconSrc: SELECTED_CONTAINER_ICON_SRC,
				clusterBorderColor: getCssVariable("--green-700"),
			});
		}

		if (truckFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(truckFeaturesRes.value, {
				layerName: $t("trucks"),
				layerColor: getCssVariable("--cyan-900"),
				iconSrc: TRUCK_ICON_SRC,
				selectedIconSrc: SELECTED_TRUCK_ICON_SRC,
				clusterBorderColor: getCssVariable("--cyan-900"),
			});
		}

		if (warehouseFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(warehouseFeaturesRes.value, {
				layerName: $t("warehouses"),
				layerColor: getCssVariable("--indigo-400"),
				iconSrc: WAREHOUSE_ICON_SRC,
				selectedIconSrc: SELECTED_WAREHOUSE_ICON_SRC,
				clusterBorderColor: getCssVariable("--indigo-400"),
			});
		}

		if (landfillFeaturesRes.status === "fulfilled") {
			mapHelper.addClusterLayer(landfillFeaturesRes.value, {
				layerName: $t("landfills"),
				layerColor: getCssVariable("--yellow-900"),
				iconSrc: LANDFILL_ICON_SRC,
				selectedIconSrc: SELECTED_LANDFILL_ICON_SRC,
				clusterBorderColor: getCssVariable("--yellow-900"),
			});
		}

		loading = false;
	}

	/**
	 * Retrieves the extent from a list of features.
	 * @param features Features to retrieve the extent from.
	 * @returns Extent.
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
				selectedContainers = getContainersFromFeature(feature);
				break;

			case "truck":
				selectedTruck = getTruckFromFeature(feature);
				break;

			case "warehouse":
				selectedWarehouse = getWarehouseFromFeature(feature);
				break;

			case "landfill":
				selectedLandfill = getLandfillFromFeature(feature);
				break;
		}
	}

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {
		// Reset all previously selected features.
		selectedContainers = [];
		selectedTruck = null;
		selectedWarehouse = null;
		selectedLandfill = null;
		selectedFeature?.set("selected", false);

		// Get features that were clicked.
		const clickedFeatures = map.getFeaturesAtPixel(e.pixel);

		// Check if the click was not performed on any feature.
		if (!clickedFeatures.length) {
			selectedFeature = null;
			return;
		}

		const mapView = map.getView();
		let extent: Extent;
		let maxZoom: number | undefined = undefined;

		const selectedFeatures: Feature[] = clickedFeatures[0].get("features");
		if (selectedFeatures.length > 1) {
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

			// Set maximum zoom to the feature.
			maxZoom = DEFAULT_MAX_ZOOM;
		}

		// Zoom map.
		mapView.fit(extent, {
			duration: 800,
			padding: [50, 50, 50, 50],
			maxZoom,
		});
	}

	onMount(() => {
		const mapHelper = new MapHelper(map);
		loadFeatures(mapHelper);

		map.on("click", handleMapClick);
	});
</script>

<main>
	<MapComponent bind:map maxZoom={22} showLayers={!selectedFeature} />

	{#if loading}
		<div class="landfill-loading">
			<Spinner />
		</div>
	{/if}

	{#if selectedFeature}
		<div class="close-selected-feature">
			<Button
				variant="tertiary"
				startIcon="close"
				size="large"
				onClick={() => {
					selectedFeature?.set("selected", false);
					selectedFeature = null;
				}}
			/>
		</div>
	{/if}

	{#if selectedFeature && (selectedContainers.length || selectedTruck || selectedWarehouse || selectedLandfill)}
		<MapBottomSheet
			wayName={selectedFeature.get("geoJson").properties.wayName}
			municipalityName={selectedFeature.get("geoJson").properties
				.municipalityName}
			warehouse={selectedWarehouse}
			containers={selectedContainers}
			truck={selectedTruck}
			landfill={selectedLandfill}
		/>
	{/if}
</main>

<style>
	main {
		position: relative;
		width: 100%;
	}

	.landfill-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		z-index: 100;
	}

	.close-selected-feature {
		position: absolute;
		top: 2.5rem;
		left: 2.5rem;

		& > button {
			box-shadow: var(--shadow-md);
		}
	}
</style>
