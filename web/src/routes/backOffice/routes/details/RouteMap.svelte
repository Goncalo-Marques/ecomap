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
		LANDFILL_ICON_SRC,
		SELECTED_CONTAINER_ICON_SRC,
		TRUCK_ICON_SRC,
		WAREHOUSE_ICON_SRC,
	} from "../../../../lib/constants/map";
	import RouteBottomSheet from "./RouteBottomSheet.svelte";
	import VectorSource from "ol/source/Vector";
	import VectorLayer from "ol/layer/Vector";
	import type { Coordinate } from "ol/coordinate";
	import { Icon, Stroke, Style } from "ol/style";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import { convertToMapProjection } from "../../../../lib/utils/map";
	import type { GeoJSONFeatureCollectionLineString } from "../../../../domain/geojson";
	import { getCssVariable } from "../../../../lib/utils/cssVars";
	import { RouteWithNoContainersError } from "../../../../lib/errors/route";

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
	 * The interval of segments added to the map within a route.
	 */
	const SEGMENT_INTERVAL = 5;

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

		const view = map.getView();

		const layerLines = new VectorLayer({
			source: new VectorSource({
				features: [feature],
			}),
			style(feature) {
				const geometry = feature.getGeometry();
				if (!(geometry instanceof MultiLineString)) {
					throw new Error("Feature is not of type MultiLineString");
				}

				const zoom = view.getZoom();
				if (!zoom) {
					return;
				}

				const isDirectionsVisible = zoom > 14;

				const styles = [
					new Style({
						stroke: new Stroke({
							color: getCssVariable("--blue-600"),
							width: isDirectionsVisible ? 24 : 12,
						}),
					}),
					new Style({
						stroke: new Stroke({
							color: getCssVariable("--blue-500"),
							width: 8,
						}),
					}),
				];

				if (isDirectionsVisible) {
					let segmentsSkipped = SEGMENT_INTERVAL;

					// Add directions style to the layer.
					for (const lineString of geometry.getLineStrings()) {
						lineString.forEachSegment((start, end) => {
							// Skip segment if the number of skipped segments is less than
							// the segment interval.
							if (segmentsSkipped < SEGMENT_INTERVAL) {
								segmentsSkipped++;
								return;
							}

							segmentsSkipped = 0;

							const dx = end[0] - start[0];
							const dy = end[1] - start[1];
							const rotation = Math.atan2(dy, dx);

							// Point in the middle of the segment.
							const middlePoint = new Point([
								dx * 0.5 + start[0],
								dy * 0.5 + start[1],
							]);

							styles.push(
								new Style({
									geometry: middlePoint,
									image: new Icon({
										src: "/images/arrow.svg",
										rotateWithView: true,
										displacement: [0, -8],
										rotation: -rotation,
									}),
								}),
							);
						});
					}
				}

				return styles;
			},
		});

		map.addLayer(layerLines);

		view.fit(multiLineString, { padding: [40, 40, 280, 40] });
	}

	/**
	 * Retrieves route container features to display in the map.
	 * @returns Route container features.
	 */
	async function getRouteContainerFeatures(): Promise<Feature<Point>[]> {
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

	/**
	 * Retrieves landfill features to display in the map.
	 * @returns Landfill features.
	 */
	async function getLandfillFeatures(): Promise<Feature<Point>[]> {
		const landfills = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/landfills", {
				params: { path: { routeId: id }, query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.landfills };
		});

		const landfillFeatures: Feature<Point>[] = [];
		for (const landfill of landfills) {
			const transformedCoordinate = convertToMapProjection(
				landfill.geoJson.geometry.coordinates,
			);
			const landfillFeature = new Feature(new Point(transformedCoordinate));
			landfillFeatures.push(landfillFeature);
		}

		return landfillFeatures;
	}

	/**
	 * Retrieves route details.
	 * @returns Route details.
	 */
	async function getRoute() {
		const res = await ecomapHttpClient.GET("/routes/{routeId}", {
			params: { path: { routeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve route details");
		}

		return res.data;
	}

	/**
	 * Retrieves route ways.
	 * @returns Route ways.
	 */
	async function getRouteWays() {
		const res = await ecomapHttpClient.GET("/routes/{routeId}/ways", {
			params: { path: { routeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve route ways");
		}

		return res.data;
	}

	/**
	 * Loads route data and displays it in the map.
	 * @returns Route data.
	 */
	async function loadRoute() {
		const [routeRes, routeWaysRes, containerFeaturesRes, landfillFeaturesRes] =
			await Promise.allSettled([
				getRoute(),
				getRouteWays(),
				getRouteContainerFeatures(),
				getLandfillFeatures(),
			]);

		if (routeRes.status === "rejected" || routeWaysRes.status === "rejected") {
			isMapVisible = false;
			throw new Error("Failed to retrieve route details");
		}

		// Check if there are no containers associated with the route.
		if (
			containerFeaturesRes.status === "fulfilled" &&
			!containerFeaturesRes.value.length
		) {
			isMapVisible = false;
			throw new RouteWithNoContainersError();
		}

		const mapHelper = new MapHelper(map);

		addRouteToMap(routeWaysRes.value);

		// Add route warehouses to map.
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

		// Add route truck to map.
		const truckFeature = new Feature(
			new Point(
				convertToMapProjection(
					routeRes.value.truck.geoJson.geometry.coordinates,
				),
			),
		);
		mapHelper.addPointLayer([truckFeature], {
			iconSrc: TRUCK_ICON_SRC,
		});

		// Add route containers to map.
		if (containerFeaturesRes.status === "fulfilled") {
			mapHelper.addPointLayer(containerFeaturesRes.value, {
				iconSrc: SELECTED_CONTAINER_ICON_SRC,
			});
		}

		// Add landfills to map.
		if (landfillFeaturesRes.status === "fulfilled") {
			mapHelper.addPointLayer(landfillFeaturesRes.value, {
				iconSrc: LANDFILL_ICON_SRC,
			});
		}

		isMapVisible = true;

		return routeRes.value;
	}

	let routePromise = loadRoute();
</script>

<main class="group relative h-auto w-full" data-mapVisible={isMapVisible}>
	<MapComponent class="group-data-[mapVisible=false]:hidden" bind:map />

	{#await routePromise}
		<Spinner
			class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
		/>
	{:then route}
		<Link to={route.id} class="contents">
			<div class="absolute left-10 top-10">
				<Button
					class="shadow-md"
					startIcon="arrow_back"
					size="large"
					variant="tertiary"
				/>
			</div>
		</Link>
		<RouteBottomSheet {route} />
	{:catch error}
		{@const errorType =
			error instanceof RouteWithNoContainersError
				? "noContainersAssociated"
				: "notFound"}
		<div
			class="absolute left-1/2 top-1/2 flex -translate-x-1/2 -translate-y-1/2 flex-col items-center justify-center"
		>
			<h2 class="text-2xl font-semibold">
				{$t(`routes.${errorType}.title`)}
			</h2>
			<p>{$t(`routes.${errorType}.description`)}</p>
		</div>
	{/await}
</main>
