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
	import type { Landfill } from "../../../../domain/landfill";
	import { LANDFILL_ICON_SRC } from "../../../../lib/constants/map";
	import LandfillBottomSheet from "./LandfillBottomSheet.svelte";

	/**
	 * Landfill ID.
	 */
	export let id: string;

	/**
	 * Indicates if map is visible.
	 * The map is hidden when landfill details are being loaded or landfill details are not found.
	 */
	let isMapVisible = true;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * Adds a landfill to the map.
	 * @param coordinates Landfill coordinates.
	 */
	function addLandfillToMap(coordinates: number[]) {
		const point = new Point(coordinates);
		const feature = new Feature(point);

		const mapHelper = new MapHelper(map);
		mapHelper.addPointLayer(
			{
				features: [feature],
			},
			"landfill",
			"#fff",
			{ "icon-src": LANDFILL_ICON_SRC },
		);

		const view = map.getView();
		view.fit(point);
	}

	/**
	 * Fetches landfill data and adds landfill to the map.
	 */
	async function fetchLandfill(): Promise<Landfill> {
		const res = await ecomapHttpClient.GET("/landfills/{landfillId}", {
			params: { path: { landfillId: id } },
		});

		if (res.error) {
			isMapVisible = false;
			throw new Error("Failed to retrieve landfill details");
		}

		const landfill = res.data;
		const landfillCoordinates = fromLonLat(
			landfill.geoJson.geometry.coordinates,
		);
		addLandfillToMap(landfillCoordinates);

		isMapVisible = true;

		return landfill;
	}

	let landfillPromise = fetchLandfill();
</script>

<main class="map" data-mapVisible={isMapVisible}>
	<Map bind:map />

	{#await landfillPromise}
		<div class="landfill-loading">
			<Spinner />
		</div>
	{:then landfill}
		<Link to={landfill.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<LandfillBottomSheet {landfill} />
	{:catch}
		<div class="landfill-not-found">
			<h2>{$t("landfills.notFound.title")}</h2>
			<p>{$t("landfills.notFound.description")}</p>
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

	.landfill-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}

	.landfill-not-found {
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
