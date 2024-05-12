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
	import type { Container } from "../../../../domain/container";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import { CONTAINER_ICON_SRC } from "../../../../lib/constants/map";
	import ContainerBottomSheet from "../../components/ContainerBottomSheet.svelte";

	/**
	 * Container ID.
	 */
	export let id: string;

	/**
	 * Indicates if map is visible.
	 * The map is hidden when container details are being loaded or container details are not found.
	 */
	let isMapVisible = true;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * Adds a container to the map.
	 * @param coordinates Container coordinates.
	 */
	function addContainerToMap(coordinates: number[]) {
		const point = new Point(coordinates);
		const feature = new Feature(point);

		const mapHelper = new MapHelper(map);
		mapHelper.addPointLayer({ features: [feature] }, "container", "#fff", {
			"icon-src": CONTAINER_ICON_SRC,
		});

		const view = map.getView();
		view.fit(point);
	}

	/**
	 * Fetches container data and adds container to the map.
	 */
	async function fetchContainer(): Promise<Container> {
		const res = await ecomapHttpClient.GET("/containers/{containerId}", {
			params: { path: { containerId: id } },
		});

		if (res.error) {
			isMapVisible = false;
			throw new Error("Failed to retrieve container details");
		}

		const container = res.data;
		const containerCoordinates = fromLonLat(
			container.geoJson.geometry.coordinates,
		);
		addContainerToMap(containerCoordinates);

		isMapVisible = true;

		return container;
	}

	let containerPromise = fetchContainer();
</script>

<main class="map" data-mapVisible={isMapVisible}>
	<Map bind:map />

	{#await containerPromise}
		<div class="container-loading">
			<Spinner />
		</div>
	{:then container}
		<Link to={container.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<ContainerBottomSheet {container} />
	{:catch}
		<div class="container-not-found">
			<h2>{$t("containers.notFound.title")}</h2>
			<p>{$t("containers.notFound.description")}</p>
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

	.container-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}

	.container-not-found {
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
