<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import type { Warehouse } from "$domain/warehouse";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Spinner from "$lib/components/Spinner.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import Card from "../../../components/Card.svelte";
	import WarehouseForm from "../../components/WarehouseForm.svelte";

	/**
	 * Warehouse ID.
	 */
	const id: string = $page.params.id;

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

		goto(`${BackOfficeRoutes.WAREHOUSES}/${id}`);
	}

	const warehousePromise = fetchWarehouse();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await warehousePromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then warehouse}
		{@const locationName = getLocationName(
			warehouse.geoJson.properties.wayName,
			warehouse.geoJson.properties.municipalityName,
		)}
		<WarehouseForm
			{warehouse}
			isSubmitting={isSubmittingForm}
			back={`${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}`}
			title={locationName}
			onSave={updateWarehouse}
		/>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("warehouses.notFound.title")}</h2>
			<p>{$t("warehouses.notFound.description")}</p>
		</div>
	{/await}
</Card>
