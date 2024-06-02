<script lang="ts">
	import OlMap from "ol/Map";
	import { Feature } from "ol";
	import { MultiLineString, Point } from "ol/geom";
	import { Link } from "svelte-routing";
	import Button from "../../../../lib/components/Button.svelte";
	import MapComponent from "../../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import {
		SELECTED_CONTAINER_ICON_SRC,
		WAREHOUSE_ICON_SRC,
	} from "../../../../lib/constants/map";
	import RouteBottomSheet from "./RouteBottomSheet.svelte";
	import VectorSource from "ol/source/Vector";
	import VectorLayer from "ol/layer/Vector";
	import type { Coordinate } from "ol/coordinate";
	import { Stroke, Style } from "ol/style";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import { convertToMapProjection } from "../../../../lib/utils/map";
	import type { GeoJSONFeatureCollectionLineString } from "../../../../domain/geojson";

	/**
	 * Route ID.
	 */
	export let id: string;

	/**
	 * Indicates if map is visible.
	 * The map is hidden when route details are being loaded or route details are not found.
	 */
	let isMapVisible = true;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * Adds a route to the map.
	 * @param geoJson Route ways.
	 */
	function addRouteToMap(geoJson: GeoJSONFeatureCollectionLineString) {
		const coordinates: Coordinate[][] = [];
		for (const feature of geoJson.features) {
			const transformedCoordinates = feature.geometry.coordinates.map(
				convertToMapProjection,
			);
			coordinates.push(transformedCoordinates);
		}

		const multiLineString = new MultiLineString(coordinates);
		const feature = new Feature(multiLineString);

		const layerLines = new VectorLayer({
			source: new VectorSource({
				features: [feature],
			}),
			style() {
				return new Style({
					stroke: new Stroke({
						width: 1,
						color: [0, 0, 255, 1],
					}),
				});
			},
		});

		map.addLayer(layerLines);

		const view = map.getView();
		view.fit(multiLineString, { padding: [40, 40, 280, 40] });
	}

	/**
	 * Retrieves container features to display in the map.
	 */
	async function getContainerFeatures(): Promise<Feature<Point>[]> {
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

		const containerFeatures: Feature<Point>[] = [];
		for (const container of containers) {
			const transformedCoordinate = convertToMapProjection(
				container.geoJson.geometry.coordinates,
			);
			const containerFeature = new Feature(new Point(transformedCoordinate));
			containerFeatures.push(containerFeature);
		}

		return containerFeatures;
	}

	async function getRoute() {
		const res = await ecomapHttpClient.GET("/routes/{routeId}", {
			params: { path: { routeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve route details");
		}

		return res.data;
	}

	async function getRouteWays() {
		const res = await ecomapHttpClient.GET("/routes/{routeId}/ways", {
			params: { path: { routeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve route ways");
		}

		return res.data;
	}

	async function loadRoute() {
		const [routeRes, routeWaysRes, containerFeaturesRes] =
			await Promise.allSettled([
				getRoute(),
				getRouteWays(),
				getContainerFeatures(),
			]);

		if (routeRes.status === "rejected" || routeWaysRes.status === "rejected") {
			isMapVisible = false;
			throw new Error("Failed to retrieve route details");
		}

		const mapHelper = new MapHelper(map);

		addRouteToMap(routeWaysRes.value);

		const departureWarehouseFeature = new Feature(
			new Point(
				convertToMapProjection(
					routeRes.value.departureWarehouse.geoJson.geometry.coordinates,
				),
			),
		);
		const arrivalWarehouseFeature = new Feature(
			new Point(
				convertToMapProjection(
					routeRes.value.arrivalWarehouse.geoJson.geometry.coordinates,
				),
			),
		);
		mapHelper.addPointLayer(
			[departureWarehouseFeature, arrivalWarehouseFeature],
			{
				iconSrc: WAREHOUSE_ICON_SRC,
			},
		);

		if (containerFeaturesRes.status === "fulfilled") {
			mapHelper.addPointLayer(containerFeaturesRes.value, {
				iconSrc: SELECTED_CONTAINER_ICON_SRC,
			});
		}

		isMapVisible = true;

		return routeRes.value;
	}

	let routePromise = loadRoute();
</script>

<main class="map" data-mapVisible={isMapVisible}>
	<MapComponent bind:map />

	{#await routePromise}
		<div class="route-loading">
			<Spinner />
		</div>
	{:then route}
		<Link to={route.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<RouteBottomSheet {route} />
	{:catch}
		<div class="route-not-found">
			<h2>{$t("routes.notFound.title")}</h2>
			<p>{$t("routes.notFound.description")}</p>
		</div>
	{/await}
</main>

<style>
	main {
		position: relative;
		height: auto;
		width: 100%;

		&[data-mapVisible="false"] {
			& #map_id {
				display: none;
			}
		}
	}

	.back {
		position: absolute;
		top: 2.5rem;
		left: 2.5rem;

		& > button {
			box-shadow: var(--shadow-md);
		}
	}

	.route-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}

	.route-not-found {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
