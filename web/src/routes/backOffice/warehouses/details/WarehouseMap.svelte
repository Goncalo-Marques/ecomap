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
	import type { Warehouse } from "../../../../domain/warehouse";
	import { WAREHOUSE_ICON_SRC } from "../../../../lib/constants/map";
	import WarehouseBottomSheet from "../../components/WarehouseBottomSheet.svelte";

	/**
	 * Warehouse ID.
	 */
	export let id: string;

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
		mapHelper.addPointLayer(
			{
				features: [feature],
			},
			"warehouse",
			"#fff",
			{ "icon-src": WAREHOUSE_ICON_SRC },
		);

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

<main class="map" data-mapVisible={isMapVisible}>
	<Map bind:map />

	{#await warehousePromise}
		<div class="warehouse-loading">
			<Spinner />
		</div>
	{:then warehouse}
		<Link to={warehouse.id} style="display:contents">
			<div class="back">
				<Button startIcon="arrow_back" size="large" variant="tertiary" />
			</div>
		</Link>
		<WarehouseBottomSheet {warehouse} />
	{:catch}
		<div class="warehouse-not-found">
			<h2>{$t("warehouses.notFound.title")}</h2>
			<p>{$t("warehouses.notFound.description")}</p>
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

	.warehouse-loading {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}

	.warehouse-not-found {
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
