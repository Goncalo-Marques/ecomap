<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import Field from "$lib/components/Field.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { getToastContext } from "$lib/contexts/toast";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import type { Warehouse } from "../../../../domain/warehouse";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import Card from "../../components/Card.svelte";

	/**
	 * Warehouse ID.
	 */
	const id: string = $page.params.id;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

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
	 * Deletes the warehouse displayed on the page.
	 */
	async function deleteWarehouse() {
		const res = await ecomapHttpClient.DELETE("/warehouses/{warehouseId}", {
			params: {
				path: {
					warehouseId: id,
				},
			},
		});

		if (res.error) {
			switch (res.error.code) {
				case "conflict":
					toast.show({
						type: "error",
						title: $t("warehouses.delete.conflict.title"),
						description: $t("warehouses.delete.conflict.description"),
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
			title: $t("warehouses.delete.success"),
			description: undefined,
		});

		goto(BackOfficeRoutes.WAREHOUSES);
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
		<DetailsHeader href={BackOfficeRoutes.WAREHOUSES} title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteWarehouse}
			/>
			<a
				href={`${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}/map`}
				class="contents"
			>
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</a>
			<a
				href={`${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}/edit`}
				class="contents"
			>
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</a>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field label={$t("location")} value={locationName} />
					<Field label={$t("truckCapacity")} value={warehouse.truckCapacity} />
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(warehouse.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(warehouse.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
		</DetailsContent>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("warehouses.notFound.title")}</h2>
			<p>{$t("warehouses.notFound.description")}</p>
		</div>
	{/await}
</Card>
