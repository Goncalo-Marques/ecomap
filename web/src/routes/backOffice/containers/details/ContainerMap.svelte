<script lang="ts">
	import OlMap from "ol/Map";
	import { Feature } from "ol";
	import { Point } from "ol/geom";
	import { fromLonLat } from "ol/proj";
	import { Link } from "svelte-routing";
	import Button from "../../../../lib/components/Button.svelte";
	import BottomSheet from "../../../../lib/components/BottomSheet.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import Field from "../../../../lib/components/Field.svelte";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import type { Container } from "../../../../domain/container";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import { formatDate } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";
	import { getLocationName } from "../../../../lib/utils/location";

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
		mapHelper.addPointLayer({ features: [feature] }, "container");

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
		{@const { wayName, municipalityName } = container.geoJson.properties}
		<Link to={container.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<BottomSheet title={getLocationName(wayName, municipalityName)}>
			<Field
				label={$t("containers.category")}
				value={$t(`containers.category.${container.category}`)}
			/>
			<Field
				label={$t("createdAt")}
				value={formatDate(container.createdAt, DateFormats.shortDateTime)}
			/>
			<Field
				label={$t("modifiedAt")}
				value={formatDate(container.modifiedAt, DateFormats.shortDateTime)}
			/>
		</BottomSheet>
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
