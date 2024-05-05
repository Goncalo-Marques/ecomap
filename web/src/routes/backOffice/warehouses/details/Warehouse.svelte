<script lang="ts">
	import { Link, navigate } from "svelte-routing";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Button from "../../../../lib/components/Button.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import Field from "../../../../lib/components/Field.svelte";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import { formatDate } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";
	import { getLocationName } from "../../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import type { Warehouse } from "../../../../domain/warehouse";

	/**
	 * Warehouse ID.
	 */
	export let id: string;

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

		navigate(BackOfficeRoutes.WAREHOUSES);
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
		<DetailsHeader to="" title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteWarehouse}
			/>
			<Link to={`${warehouse.id}/map`} style="display:contents">
				<Button variant="secondary" startIcon="map">
					{$t("sidebar.map")}
				</Button>
			</Link>
			<Link to={`${warehouse.id}/edit`} style="display:contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</Link>
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
