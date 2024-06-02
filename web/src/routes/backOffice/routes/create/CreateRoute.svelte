<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import RouteForm from "../components/RouteForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import type { RouteEmployeeRole } from "../../../../domain/routeEmployee";

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
		containersIds: {
			added: string[];
			deleted: string[];
		},
		operators: {
			added: { routeRole: RouteEmployeeRole; id: string }[];
			deleted: { routeRole: RouteEmployeeRole; id: string }[];
		},
	) {
		const routeRes = await ecomapHttpClient.POST("/routes", {
			body: {
				name,
				arrivalWarehouseId,
				departureWarehouseId,
				truckId,
			},
		});

		if (routeRes.error) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});
			return;
		}

		const promises = [];

		for (const containerId of containersIds.added) {
			promises.push(
				ecomapHttpClient.POST("/routes/{routeId}/containers/{containerId}", {
					params: {
						path: {
							routeId: routeRes.data.id,
							containerId,
						},
					},
				}),
			);
		}

		for (const operator of operators.added) {
			promises.push(
				ecomapHttpClient.POST("/routes/{routeId}/employees/{employeeId}", {
					params: {
						path: {
							routeId: routeRes.data.id,
							employeeId: operator.id,
						},
					},
					body: {
						routeRole: operator.routeRole,
					},
				}),
			);
		}

		await Promise.allSettled(promises);

		toast.show({
			type: "success",
			title: $t("routes.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.ROUTES}/${routeRes.data.id}`);
	}
</script>

<Card class="page-layout">
	<RouteForm back="" title={$t("routes.create.title")} onSave={createRoute} />
</Card>
