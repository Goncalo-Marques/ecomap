<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import WarehouseForm from "../components/WarehouseForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

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
	 * Creates a warehouse with a given truck capacity and location.
	 * @param truckCapacity Warehouse truck capacity.
	 * @param location Warehouse location.
	 */
	async function createWarehouse(
		truckCapacity: number,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.POST("/warehouses", {
			body: {
				truckCapacity,
				geoJson: location,
			},
		});

		isSubmittingForm = false;

		if (res.error) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});
			return;
		}

		toast.show({
			type: "success",
			title: $t("warehouses.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.WAREHOUSES}/${res.data.id}`);
	}
</script>

<Card class="page-layout">
	<WarehouseForm
		back=""
		title={$t("warehouses.create.title")}
		isSubmitting={isSubmittingForm}
		onSave={createWarehouse}
	/>
</Card>
