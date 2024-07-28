<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Route } from "$domain/route";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import Field from "$lib/components/Field.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import Card from "../../components/Card.svelte";
	import RouteOperatorsTable from "./details/routeEmployees/RouteEmployeesTable.svelte";

	/**
	 * Route ID.
	 */
	const id: string = $page.params.id;

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
	 * Deletes the route displayed on the page.
	 */
	async function deleteRoute() {
		const res = await ecomapHttpClient.DELETE("/routes/{routeId}", {
			params: {
				path: {
					routeId: id,
				},
			},
		});

		if (res.error) {
			switch (res.error.code) {
				case "conflict":
					toast.show({
						type: "error",
						title: $t("routes.delete.conflict.title"),
						description: $t("routes.delete.conflict.description"),
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
			title: $t("routes.delete.success"),
			description: undefined,
		});

		goto(BackOfficeRoutes.ROUTES);
	}

	const routePromise = fetchRoute();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await routePromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then route}
		<DetailsHeader href={BackOfficeRoutes.ROUTES} title={route.name}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteRoute}
			/>
			<a href={`${BackOfficeRoutes.ROUTES}/${route.id}/map`} class="contents">
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</a>
			<a href={`${BackOfficeRoutes.ROUTES}/${route.id}/edit`} class="contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</a>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field
						label={$t("departure")}
						value={getLocationName(
							route.departureWarehouse.geoJson.properties.wayName,
							route.departureWarehouse.geoJson.properties.municipalityName,
						)}
					/>
					<Field
						label={$t("arrival")}
						value={getLocationName(
							route.arrivalWarehouse.geoJson.properties.wayName,
							route.arrivalWarehouse.geoJson.properties.municipalityName,
						)}
					/>
					<Field
						label={$t("truck")}
						value={`${route.truck.make} ${route.truck.model} (${route.truck.licensePlate})`}
					/>
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(route.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(route.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
			<DetailsSection class="h-full" label={$t("operators")}>
				<RouteOperatorsTable routeId={id} />
			</DetailsSection>
		</DetailsContent>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("routes.notFound.title")}</h2>
			<p>{$t("routes.notFound.description")}</p>
		</div>
	{/await}
</Card>
