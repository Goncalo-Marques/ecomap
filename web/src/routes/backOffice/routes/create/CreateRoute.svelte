<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import RouteForm from "../components/RouteForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Creates a route with a given name, departure warehouse,
	 * arrival warehouse, truck and location.
	 * @param name Route name.
	 * @param departureWarehouseId Route departure warehouse ID.
	 * @param arrivalWarehouseId Route arrival warehouse ID.
	 * @param truckID Route truck ID.
	 */
	async function createRoute(
		name: string,
		departureWarehouseId: string,
		arrivalWarehouseId: string,
		truckId: string,
	) {
		const res = await ecomapHttpClient.POST("/routes", {
			body: {
				name,
				arrivalWarehouseId,
				departureWarehouseId,
				truckId,
			},
		});

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
			title: $t("routes.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.ROUTES}/${res.data.id}`);
	}
</script>

<Card class="page-layout">
	<RouteForm back="" title={$t("routes.create.title")} onSave={createRoute} />
</Card>
