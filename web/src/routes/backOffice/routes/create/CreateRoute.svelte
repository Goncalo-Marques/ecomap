<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import RouteForm from "../components/RouteForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import type { SelectedRouteContainersIds } from "../../../../domain/container";
	import type { SelectedRouteEmployees } from "../../../../domain/routeEmployee";
	import type { Truck } from "../../../../domain/truck";

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
	 * Creates a route with a given name, departure warehouse,
	 * arrival warehouse, truck and assigns containers and employees
	 * to that route.
	 * @param name Route name.
	 * @param departureWarehouseId Route departure warehouse ID.
	 * @param arrivalWarehouseId Route arrival warehouse ID.
	 * @param truck Route truck.
	 * @param containersIds Container IDs.
	 * @param routeEmployees Route employees.
	 */
	async function createRoute(
		name: string,
		departureWarehouseId: string,
		arrivalWarehouseId: string,
		truck: Truck,
		containersIds: SelectedRouteContainersIds,
		routeEmployees: SelectedRouteEmployees,
	) {
		isSubmittingForm = true;

		const routeRes = await ecomapHttpClient.POST("/routes", {
			body: {
				name,
				arrivalWarehouseId,
				departureWarehouseId,
				truckId: truck.id,
			},
		});

		isSubmittingForm = false;

		if (routeRes.error) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});
			return;
		}

		const promises = [];

		// Add promises that associate each added container with the created route.
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

		// Add promises that associate each added employee with the created route.
		for (const routeEmployee of routeEmployees.added) {
			promises.push(
				ecomapHttpClient.POST("/routes/{routeId}/employees/{employeeId}", {
					params: {
						path: {
							routeId: routeRes.data.id,
							employeeId: routeEmployee.id,
						},
					},
					body: {
						routeRole: routeEmployee.routeRole,
					},
				}),
			);
		}

		// Execute all promises.
		await Promise.allSettled(promises);

		toast.show({
			type: "success",
			title: $t("routes.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.ROUTES}/${routeRes.data.id}`);
	}
</script>

<Card class="m-10 flex flex-col gap-10">
	<RouteForm
		back=""
		isSubmitting={isSubmittingForm}
		title={$t("routes.create.title")}
		onSave={createRoute}
	/>
</Card>
