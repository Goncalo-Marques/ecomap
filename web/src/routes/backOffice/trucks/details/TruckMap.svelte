<script lang="ts">
	import OlMap from "ol/Map";
	import { Feature } from "ol";
	import { Point } from "ol/geom";
	import { fromLonLat } from "ol/proj";
	import { Link } from "svelte-routing";
	import Button from "../../../../lib/components/Button.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { Truck } from "../../../../domain/truck";
	import { TRUCK_ICON_SRC } from "../../../../lib/constants/map";
	import TruckBottomSheet from "../../components/TruckBottomSheet.svelte";

	/**
	 * Truck ID.
	 */
	export let id: string;

	/**
	 * Indicates if map is visible.
	 * The map is hidden when truck details are being loaded or truck details are not found.
	 */
	let isMapVisible = true;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * Adds a truck to the map.
	 * @param coordinates Truck coordinates.
	 */
	function addTruckToMap(coordinates: number[]) {
		const point = new Point(coordinates);
		const feature = new Feature(point);

		const mapHelper = new MapHelper(map);
		mapHelper.addPointLayer(
			{
				features: [feature],
			},
			"truck",
			"#fff",
			{ "icon-src": TRUCK_ICON_SRC },
		);

		const view = map.getView();
		view.fit(point);
	}

	/**
	 * Fetches truck data and adds truck to the map.
	 */
	async function fetchTruck(): Promise<Truck> {
		const res = await ecomapHttpClient.GET("/trucks/{truckId}", {
			params: { path: { truckId: id } },
		});

		if (res.error) {
			isMapVisible = false;
			throw new Error("Failed to retrieve truck details");
		}

		const truck = res.data;
		const truckCoordinates = fromLonLat(truck.geoJson.geometry.coordinates);
		addTruckToMap(truckCoordinates);

		isMapVisible = true;

		return truck;
	}

	let truckPromise = fetchTruck();
</script>

<main class="map" data-mapVisible={isMapVisible}>
	<Map bind:map />

	{#await truckPromise}
		<div class="truck-loading">
			<Spinner />
		</div>
	{:then truck}
		<Link to={truck.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<TruckBottomSheet {truck} />
	{:catch}
		<div class="truck-not-found">
			<h2>{$t("trucks.notFound.title")}</h2>
			<p>{$t("trucks.notFound.description")}</p>
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

	.truck-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}

	.truck-not-found {
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
