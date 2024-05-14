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
	import type { Truck } from "../../../../domain/truck";

	/**
	 * Truck ID.
	 */
	export let id: string;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Fetches truck data.
	 */
	async function fetchTruck(): Promise<Truck> {
		const res = await ecomapHttpClient.GET("/trucks/{truckId}", {
			params: { path: { truckId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve truck details");
		}

		return res.data;
	}

	/**
	 * Deletes the truck displayed on the page.
	 */
	async function deleteTruck() {
		const res = await ecomapHttpClient.DELETE("/trucks/{truckId}", {
			params: {
				path: {
					truckId: id,
				},
			},
		});

		if (res.error) {
			switch (res.error.code) {
				case "conflict":
					toast.show({
						type: "error",
						title: $t("trucks.delete.conflict.title"),
						description: $t("trucks.delete.conflict.description"),
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
			title: $t("trucks.delete.success"),
			description: undefined,
		});

		navigate(BackOfficeRoutes.TRUCKS);
	}

	const truckPromise = fetchTruck();
</script>

<Card class="page-layout">
	{#await truckPromise}
		<div class="truck-loading">
			<Spinner />
		</div>
	{:then truck}
		{@const locationName = getLocationName(
			truck.geoJson.properties.wayName,
			truck.geoJson.properties.municipalityName,
		)}
		<DetailsHeader to="" title={`${truck.make} ${truck.model}`}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteTruck}
			/>
			<Link to={`${truck.id}/map`} style="display:contents">
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</Link>
			<Link to={`${truck.id}/edit`} style="display:contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</Link>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field label={$t("make")} value={truck.make} />
					<Field label={$t("model")} value={truck.model} />
					<Field label={$t("licensePlate")} value={truck.licensePlate} />
					<Field label={$t("personCapacity")} value={truck.personCapacity} />
					<Field label={$t("location")} value={locationName} />
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(truck.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(truck.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
		</DetailsContent>
	{:catch}
		<div class="truck-not-found">
			<h2>{$t("trucks.notFound.title")}</h2>
			<p>{$t("trucks.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.truck-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.truck-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
