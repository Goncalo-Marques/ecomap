<script lang="ts">
	import { navigate } from "svelte-routing";
	import WarehouseForm from "../components/WarehouseForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { t } from "../../../../lib/utils/i8n";
	import { getLocationName } from "../../../../lib/utils/location";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import type { Warehouse } from "../../../../domain/warehouse";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Warehouse ID.
	 */
	export let id: string;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Indicates if form is being submitted.
	 * @default false
	 */
	let isSubmittingForm: boolean = false;

	/**
	 * Fetches warehouse data.
	 */
	async function fetchWarehouse(): Promise<Warehouse> {
		const res = await ecomapHttpClient.GET("/warehouses/{warehouseId}", {
			params: { path: { warehouseId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve warehouse details");
		}

		return res.data;
	}

	/**
	 * Updates a warehouse with a given truck capacity and location.
	 * @param truckCapacity Warehouse truck capacity.
	 * @param location Warehouse location.
	 */
	async function updateWarehouse(
		truckCapacity: number,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.PATCH("/warehouses/{warehouseId}", {
			params: {
				path: {
					warehouseId: id,
				},
			},
			body: {
				truckCapacity,
				geoJson: location,
			},
		});

		isSubmittingForm = false;

		if (res.error) {
			switch (res.error.code) {
				case "conflict":
					toast.show({
						type: "error",
						title: $t("warehouses.update.conflict.title"),
						description: $t("warehouses.update.conflict.description"),
					});
					break;

				default:
					toast.show({
						type: "error",
						title: $t("error.unexpected.title"),
						description: $t("error.unexpected.description"),
					});
					break;
			}

			return;
		}

		toast.show({
			type: "success",
			title: $t("warehouses.update.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.WAREHOUSES}/${id}`);
	}

	const warehousePromise = fetchWarehouse();
</script>

<Card class="page-layout">
	{#await warehousePromise}
		<div class="warehouse-loading">
			<Spinner />
		</div>
	{:then warehouse}
		{@const locationName = getLocationName(
			warehouse.geoJson.properties.wayName,
			warehouse.geoJson.properties.municipalityName,
		)}
		<WarehouseForm
			{warehouse}
			isSubmitting={isSubmittingForm}
			back={warehouse.id}
			title={locationName}
			onSave={updateWarehouse}
		/>
	{:catch}
		<div class="warehouse-not-found">
			<h2>{$t("warehouses.notFound.title")}</h2>
			<p>{$t("warehouses.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.warehouse-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.warehouse-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
