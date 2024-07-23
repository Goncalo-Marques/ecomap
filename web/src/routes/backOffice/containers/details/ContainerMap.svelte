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
	import ContainerBottomSheet from "./ContainerBottomSheet.svelte";

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
		mapHelper.addPointLayer([feature], { iconSrc: CONTAINER_ICON_SRC });

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

<main class="group relative h-auto w-full" data-mapVisible={isMapVisible}>
	<Map class="group-data-[mapVisible=false]:hidden" bind:map />

	{#await containerPromise}
		<Spinner
			class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
		/>
	{:then container}
		<Link to={container.id} class="contents">
			<div class="absolute left-10 top-10">
				<Button
					class="shadow-md"
					startIcon="arrow_back"
					size="large"
					variant="tertiary"
				/>
			</div>
		</Link>
		<ContainerBottomSheet {container} />
	{:catch}
		<div
			class="absolute left-1/2 top-1/2 flex -translate-x-1/2 -translate-y-1/2 flex-col items-center justify-center"
		>
			<h2 class="text-2xl font-semibold">{$t("containers.notFound.title")}</h2>
			<p>{$t("containers.notFound.description")}</p>
		</div>
	{/await}
</main>
