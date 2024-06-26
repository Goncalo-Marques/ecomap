<script lang="ts">
	import { navigate } from "svelte-routing";
	import LandfillForm from "../components/LandfillForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { t } from "../../../../lib/utils/i8n";
	import { getLocationName } from "../../../../lib/utils/location";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import type { Landfill } from "../../../../domain/landfill";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Landfill ID.
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

		navigate(`${BackOfficeRoutes.LANDFILLS}/${id}`);
	}

	const landfillPromise = fetchLandfill();
</script>

<Card class="page-layout">
	{#await landfillPromise}
		<div class="landfill-loading">
			<Spinner />
		</div>
	{:then landfill}
		{@const locationName = getLocationName(
			landfill.geoJson.properties.wayName,
			landfill.geoJson.properties.municipalityName,
		)}
		<LandfillForm
			{landfill}
			isSubmitting={isSubmittingForm}
			back={landfill.id}
			title={locationName}
			onSave={updateLandfill}
		/>
	{:catch}
		<div class="landfill-not-found">
			<h2>{$t("landfills.notFound.title")}</h2>
			<p>{$t("landfills.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.landfill-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.landfill-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
