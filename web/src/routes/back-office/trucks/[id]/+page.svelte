<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Truck } from "$domain/truck";
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

	/**
	 * Truck ID.
	 */
	const id: string = $page.params.id;

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

		goto(BackOfficeRoutes.TRUCKS);
	}

	const truckPromise = fetchTruck();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await truckPromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then truck}
		{@const locationName = getLocationName(
			truck.geoJson.properties.wayName,
			truck.geoJson.properties.municipalityName,
		)}
		<DetailsHeader
			href={BackOfficeRoutes.TRUCKS}
			title={`${truck.make} ${truck.model}`}
		>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteTruck}
			/>
			<a href={`${BackOfficeRoutes.TRUCKS}/${truck.id}/map`} class="contents">
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</a>
			<a href={`${BackOfficeRoutes.TRUCKS}/${truck.id}/edit`} class="contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</a>
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
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("trucks.notFound.title")}</h2>
			<p>{$t("trucks.notFound.description")}</p>
		</div>
	{/await}
</Card>
