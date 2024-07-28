<script lang="ts">
	import { Feature } from "ol";
	import { Point } from "ol/geom";
	import OlMap from "ol/Map";
	import { fromLonLat } from "ol/proj";

	import { page } from "$app/stores";
	import type { Warehouse } from "$domain/warehouse";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import Map from "$lib/components/map/Map.svelte";
	import { MapHelper } from "$lib/components/map/mapUtils";
	import Spinner from "$lib/components/Spinner.svelte";
	import { WAREHOUSE_ICON_SRC } from "$lib/constants/map";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { t } from "$lib/utils/i8n";

	import WarehouseBottomSheet from "./WarehouseBottomSheet.svelte";

	/**
	 * Warehouse ID.
	 */
	const id: string = $page.params.id;

	/**
	 * Indicates if map is visible.
	 * The map is hidden when warehouse details are being loaded or warehouse details are not found.
	 */
	let isMapVisible = true;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * Adds a warehouse to the map.
	 * @param coordinates Warehouse coordinates.
	 */
	function addWarehouseToMap(coordinates: number[]) {
		const point = new Point(coordinates);
		const feature = new Feature(point);

		const mapHelper = new MapHelper(map);
		mapHelper.addPointLayer([feature], { iconSrc: WAREHOUSE_ICON_SRC });

		const view = map.getView();
		view.fit(point);
	}

	/**
	 * Fetches warehouse data and adds warehouse to the map.
	 */
	async function fetchWarehouse(): Promise<Warehouse> {
		const res = await ecomapHttpClient.GET("/warehouses/{warehouseId}", {
			params: { path: { warehouseId: id } },
		});

		if (res.error) {
			isMapVisible = false;
			throw new Error("Failed to retrieve warehouse details");
		}

		const warehouse = res.data;
		const warehouseCoordinates = fromLonLat(
			warehouse.geoJson.geometry.coordinates,
		);
		addWarehouseToMap(warehouseCoordinates);

		isMapVisible = true;

		return warehouse;
	}

	let warehousePromise = fetchWarehouse();
</script>

<main class="group relative h-auto w-full" data-mapVisible={isMapVisible}>
	<Map class="group-data-[mapVisible=false]:hidden" bind:map />

	{#await warehousePromise}
		<Spinner
			class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
		/>
	{:then warehouse}
		<a href={`${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}`} class="contents">
			<div class="absolute left-10 top-10">
				<Button
					class="shadow-md"
					startIcon="arrow_back"
					size="large"
					variant="tertiary"
				/>
			</div>
		</a>
		<WarehouseBottomSheet {warehouse} />
	{:catch}
		<div
			class="absolute left-1/2 top-1/2 flex -translate-x-1/2 -translate-y-1/2 flex-col items-center justify-center"
		>
			<h2 class="text-2xl font-semibold">{$t("warehouses.notFound.title")}</h2>
			<p>{$t("warehouses.notFound.description")}</p>
		</div>
	{/await}
</main>
