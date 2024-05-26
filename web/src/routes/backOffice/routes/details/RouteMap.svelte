<script lang="ts">
	import OlMap from "ol/Map";
	import { Feature } from "ol";
	import { LineString, MultiLineString, Point } from "ol/geom";
	import { fromLonLat } from "ol/proj";
	import { Link } from "svelte-routing";
	import Button from "../../../../lib/components/Button.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import type { Route } from "../../../../domain/route";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import { CONTAINER_ICON_SRC } from "../../../../lib/constants/map";
	import RouteBottomSheet from "./RouteBottomSheet.svelte";
	import type { FeatureObject, SimpleGeometryObject } from "ol/format/Feature";
	import VectorSource from "ol/source/Vector";
	import VectorLayer from "ol/layer/Vector";
	import type { Coordinate } from "ol/coordinate";
	import { Stroke, Style } from "ol/style";

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
	 * @param coordinates Route coordinates.
	 */
	function addRouteToMap(coordinates: Coordinate[][]) {
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
	 * Fetches route data and adds route to the map.
	 */
	async function fetchRoute(): Promise<Route> {
		const [routeRes, routeWaysRes] = await Promise.allSettled([
			ecomapHttpClient.GET("/routes/{routeId}", {
				params: { path: { routeId: id } },
			}),
			ecomapHttpClient.GET("/routes/{routeId}/ways", {
				params: { path: { routeId: id } },
			}),
		]);

		if (
			routeRes.status === "rejected" ||
			routeWaysRes.status === "rejected" ||
			(routeRes.status === "fulfilled" && routeRes.value.error) ||
			(routeWaysRes.status === "fulfilled" && routeWaysRes.value.error)
		) {
			isMapVisible = false;
			throw new Error("Failed to retrieve route details");
		}

		const route = routeRes.value.data!;
		const routeWays = routeWaysRes.value.data!;

		const coordinates: Coordinate[][] = [];
		for (const feature of routeWays.features) {
			const transformedCoordinates = feature.geometry.coordinates.map(
				coordinate => fromLonLat(coordinate),
			);
			coordinates.push(transformedCoordinates);
		}

		addRouteToMap(coordinates);

		isMapVisible = true;

		return route;
	}

	let routePromise = fetchRoute();
</script>

<main class="map" data-mapVisible={isMapVisible}>
	<Map bind:map />

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
