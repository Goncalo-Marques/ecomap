<script lang="ts">
	import { navigate } from "svelte-routing";
	import RouteForm from "../components/RouteForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { t } from "../../../../lib/utils/i8n";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { Route } from "../../../../domain/route";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import type { SelectedRouteContainersIds } from "../../../../domain/container";
	import type { SelectedRouteEmployees } from "../../../../domain/routeEmployee";
	import type { Truck } from "../../../../domain/truck";

	/**
	 * Route ID.
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
	 * Fetches route data.
	 * @returns Route data.
	 */
	async function fetchRoute(): Promise<Route> {
		const res = await ecomapHttpClient.GET("/routes/{routeId}", {
			params: { path: { routeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve route details");
		}

		return res.data;
	}

	/**
	 * Updates a route with a given name, departure warehouse,
	 * arrival warehouse, truck and assigns containers and employees
	 * to that route.
	 * @param currentRoute Current route.
	 * @param name Route name.
	 * @param departureWarehouseId Route departure warehouse ID.
	 * @param arrivalWarehouseId Route arrival warehouse ID.
	 * @param truck Route truck.
	 * @param containersIds Container IDs.
	 * @param routeEmployees Route employees.
	 */
	async function updateRoute(
		currentRoute: Route,
		name: string,
		departureWarehouseId: string,
		arrivalWarehouseId: string,
		truck: Truck,
		containersIds: SelectedRouteContainersIds,
		routeEmployees: SelectedRouteEmployees,
	) {
		async function performRouteUpdate() {
			const routeRes = await ecomapHttpClient.PATCH("/routes/{routeId}", {
				params: {
					path: {
						routeId: id,
					},
				},
				body: {
					name,
					departureWarehouseId,
					arrivalWarehouseId,
					truckId: truck.id,
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

			toast.show({
				type: "success",
				title: $t("routes.update.success"),
				description: undefined,
			});
		}

		async function performRouteAssociations() {
			const promises = [];

			// Add promises that remove each removed container with the updated route.
			for (const containerId of containersIds.deleted) {
				promises.push(
					ecomapHttpClient.DELETE(
						"/routes/{routeId}/containers/{containerId}",
						{
							params: {
								path: {
									routeId: id,
									containerId,
								},
							},
						},
					),
				);
			}

			// Add promises that associate each added container with the updated route.
			for (const containerId of containersIds.added) {
				promises.push(
					ecomapHttpClient.POST("/routes/{routeId}/containers/{containerId}", {
						params: {
							path: {
								routeId: id,
								containerId,
							},
						},
					}),
				);
			}

			// Add promises that remove the association of each removed container with the updated route.
			for (const routeEmployee of routeEmployees.deleted) {
				promises.push(
					ecomapHttpClient.DELETE("/routes/{routeId}/employees/{employeeId}", {
						params: {
							path: {
								routeId: id,
								employeeId: routeEmployee.id,
							},
						},
					}),
				);
			}

			// Add promises that add the association of each added container with the updated route.
			for (const routeEmployee of routeEmployees.added) {
				promises.push(
					ecomapHttpClient.POST("/routes/{routeId}/employees/{employeeId}", {
						params: {
							path: {
								routeId: id,
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
		}

		isSubmittingForm = true;

		// Perform requests based on truck person capacity.
		// If the selected truck has a higher capacity than the current truck associated with the route,
		// the route must be updated first to avoid conflicts with the number of employees associated.
		// Otherwise, perform the associations first.
		if (truck.personCapacity > currentRoute.truck.personCapacity) {
			await performRouteUpdate();
			await performRouteAssociations();
		} else {
			await performRouteAssociations();
			await performRouteUpdate();
		}

		isSubmittingForm = false;

		navigate(`${BackOfficeRoutes.ROUTES}/${id}`);
	}

	const routePromise = fetchRoute();
</script>

<Card class="page-layout">
	{#await routePromise}
		<div class="route-loading">
			<Spinner />
		</div>
	{:then route}
		<RouteForm
			{route}
			isSubmitting={isSubmittingForm}
			back={route.id}
			title={route.name}
			onSave={(
				name,
				departureWarehouseId,
				arrivalWarehouseId,
				truck,
				containersIds,
				routeEmployees,
			) => {
				updateRoute(
					route,
					name,
					departureWarehouseId,
					arrivalWarehouseId,
					truck,
					containersIds,
					routeEmployees,
				);
			}}
		/>
	{:catch}
		<div class="route-not-found">
			<h2>{$t("routes.notFound.title")}</h2>
			<p>{$t("routes.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.route-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.route-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
