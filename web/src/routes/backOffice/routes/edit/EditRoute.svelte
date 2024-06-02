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
	import type { RouteEmployeeRole } from "../../../../domain/routeEmployee";

	/**
	 * Route ID.
	 */
	export let id: string;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Fetches route data.
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
	 * arrival warehouse, truck and location.
	 * @param name Route name.
	 * @param departureWarehouseId Route departure warehouse ID.
	 * @param arrivalWarehouseId Route arrival warehouse ID.
	 * @param truckID Route truck ID.
	 */
	async function updateRoute(
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
		for (const containerId of containersIds.deleted) {
			promises.push(
				ecomapHttpClient.DELETE("/routes/{routeId}/containers/{containerId}", {
					params: {
						path: {
							routeId: routeRes.data.id,
							containerId,
						},
					},
				}),
			);
		}

		for (const operator of operators.deleted) {
			promises.push(
				ecomapHttpClient.DELETE("/routes/{routeId}/employees/{employeeId}", {
					params: {
						path: {
							routeId: routeRes.data.id,
							employeeId: operator.id,
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
			title: $t("routes.update.success"),
			description: undefined,
		});

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
			back={route.id}
			title={route.name}
			onSave={updateRoute}
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
