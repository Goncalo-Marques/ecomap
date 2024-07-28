<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import type { Landfill } from "$domain/landfill";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Spinner from "$lib/components/Spinner.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import Card from "../../../components/Card.svelte";
	import LandfillForm from "../../components/LandfillForm.svelte";

	/**
	 * Landfill ID.
	 */
	const id: string = $page.params.id;

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
	 * Fetches landfill data.
	 */
	async function fetchLandfill(): Promise<Landfill> {
		const res = await ecomapHttpClient.GET("/landfills/{landfillId}", {
			params: { path: { landfillId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve landfill details");
		}

		return res.data;
	}

	/**
	 * Updates a landfill with a given location.
	 * @param location Landfill location.
	 */
	async function updateLandfill(location: GeoJSONFeaturePoint) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.PATCH("/landfills/{landfillId}", {
			params: {
				path: {
					landfillId: id,
				},
			},
			body: {
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
			title: $t("landfills.update.success"),
			description: undefined,
		});

		goto(`${BackOfficeRoutes.LANDFILLS}/${id}`);
	}

	const landfillPromise = fetchLandfill();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await landfillPromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then landfill}
		{@const locationName = getLocationName(
			landfill.geoJson.properties.wayName,
			landfill.geoJson.properties.municipalityName,
		)}
		<LandfillForm
			{landfill}
			isSubmitting={isSubmittingForm}
			back={`${BackOfficeRoutes.LANDFILLS}/${landfill.id}`}
			title={locationName}
			onSave={updateLandfill}
		/>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("landfills.notFound.title")}</h2>
			<p>{$t("landfills.notFound.description")}</p>
		</div>
	{/await}
</Card>
