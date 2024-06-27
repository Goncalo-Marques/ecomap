<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import OlMap from "ol/Map";
	import { Feature, MapBrowserEvent } from "ol";
	import { Point, SimpleGeometry } from "ol/geom";
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import type { Container } from "../../../../domain/container";
	import Modal from "../../../../lib/components/Modal.svelte";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import { boundingExtent, type Extent } from "ol/extent";
	import type { Coordinate } from "ol/coordinate";
	import {
		CONTAINER_ICON_SRC,
		SELECTED_CONTAINER_ICON_SRC,
	} from "../../../../lib/constants/map";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { convertToMapProjection } from "../../../../lib/utils/map";
	import MapComponent from "../../../../lib/components/map/Map.svelte";
	import { getCssVariable } from "../../../../lib/utils/cssVars";

	/**
	 * The route ID.
	 */
	export let routeId: string | undefined;

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
	 * @param addedContainers Container added.
	 * @param deletedContainers Container removed.
	 */
	export let onSave: (
		addedContainers: Container[],
		deletedContainers: Container[],
	) => void;

	/**
	 * Added containers.
	 */
	export let addedContainersMap = new Map<Container["id"], Container>();

	/**
	 * Deleted containers.
	 */
	export let deletedContainersMap = new Map<Container["id"], Container>();

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
	 * The original containers.
	 */
	let originalContainersMap = new Map<Container["id"], Container>();

	/**
	 * The selected features in the map.
	 */
	let selectedFeatures: Feature[] = [];

	/**
	 * Indicates if save action is disabled.
	 * @default true
	 */
	let isSaveActionDisabled = true;

	/**
	 * Retrieves the ID of a feature.
	 * @param coordinate Feature coordinate.
	 */
	function getFeatureId(coordinate: Coordinate) {
		return `${coordinate[0]},${coordinate[1]}`;
	}

	/**
	 * Retrieves the containers from the route ID.
	 * @param id Route ID.
	 * @returns Route containers.
	 */
	async function getRouteContainers(id: string) {
		const containers = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/routes/{routeId}/containers", {
					params: { path: { routeId: id }, query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.containers };
			},
		);

		return containers;
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

			const featureId = getFeatureId(geoJson.geometry.coordinates);

			if (!containerMap.has(featureId)) {
				const transformedCoordinate = convertToMapProjection(
					geoJson.geometry.coordinates,
				);
				const point = new Point(transformedCoordinate);
				const feature = new Feature(point);

				feature.setId(featureId);
				feature.setProperties({ geoJson, items: [] });

				containerMap.set(featureId, feature);
			}

			const containerFeature = containerMap.get(featureId)!;
			containerFeature.get("items").push({
				id,
				category,
				createdAt,
				modifiedAt,
				geoJson,
			});

			if (originalContainersMap.has(container.id)) {
				containerFeature.set("selected", true);
				selectedFeatures.push(containerFeature);
				isSaveActionDisabled = false;
			}
		}

		return Array.from(containerMap.values());
	}

	/**
	 * Loads all features into the map.
	 * @param mapHelper Map helper to add features to the map.
	 */
	async function loadFeatures(mapHelper: MapHelper) {
		const containerFeatures = await getContainerFeatures();
		mapHelper.addClusterLayer(containerFeatures, {
			iconSrc: CONTAINER_ICON_SRC,
			selectedIconSrc: SELECTED_CONTAINER_ICON_SRC,
			clusterBorderColor: getCssVariable("--green-700"),
		});
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

		// Add feature to the list of selected features.
		selectedFeatures.push(feature);

		const containers: Container[] = feature.get("items");
		for (const container of containers) {
			if (!originalContainersMap.has(container.id)) {
				addedContainersMap.set(container.id, container);
			}

			if (deletedContainersMap.has(container.id)) {
				deletedContainersMap.delete(container.id);
			}
		}
	}

	/**
	 * Unselects a single feature.
	 * @param feature Feature.
	 */
	function unselectFeature(feature: Feature) {
		// Mark feature as unselected.
		feature.set("selected", false);

		// Unselect feature from selected features list.
		selectedFeatures = selectedFeatures.filter(
			selectedFeature => selectedFeature.getId() !== feature.getId(),
		);

		const containers: Container[] = feature.get("items");
		for (const container of containers) {
			if (originalContainersMap.has(container.id)) {
				deletedContainersMap.set(container.id, container);
			}

			if (addedContainersMap.has(container.id)) {
				addedContainersMap.delete(container.id);
			}
		}
	}

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {
		// Get features that were clicked.
		const clickedFeatures = map.getFeaturesAtPixel(e.pixel);

		// Check if the click was not performed on any feature.
		if (!clickedFeatures.length) {
			return;
		}

		const features: Feature[] = clickedFeatures[0].get("features");
		if (features.length > 1) {
			const extent = getExtentFromSelectedFeatures(features);

			// Zoom map.
			map.getView().fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
		} else {
			const feature = features[0];
			if (feature.get("selected")) {
				unselectFeature(feature);
			} else {
				selectFeature(feature);
			}
		}

		isSaveActionDisabled = !selectedFeatures.length;
	}

	/**
	 * Handles the save action of the modal.
	 */
	async function handleSave() {
		onSave(
			Array.from(addedContainersMap.values()),
			Array.from(deletedContainersMap.values()),
		);
		onOpenChange(false);
	}

	/**
	 * Handles cancel action.
	 */
	function handleCancel() {
		onOpenChange(false);
		onCancel?.();
	}

	onMount(async () => {
		// Fill original containers map when a route ID is defined.
		if (routeId) {
			const containers = await getRouteContainers(routeId);
			for (const container of containers) {
				originalContainersMap.set(container.id, container);
			}
		}

		const mapHelper = new MapHelper(map);
		loadFeatures(mapHelper);

		map.on("click", handleMapClick);
	});

	onDestroy(() => {
		map.un("click", handleMapClick);
	});
</script>

<Modal {open} {onOpenChange} title={$t("routes.selectContainers")}>
	<MapComponent
		bind:map
		mapId="select-containers-map"
		maxZoom={22}
		--height="50vh"
		--width="60vw"
		--max-height="32rem"
		--max-width="60rem"
	/>
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
