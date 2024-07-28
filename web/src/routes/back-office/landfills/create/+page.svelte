<script lang="ts">
	import { goto } from "$app/navigation";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";

	import Card from "../../components/Card.svelte";
	import LandfillForm from "../components/LandfillForm.svelte";

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
	 * Creates a landfill with a given location.
	 * @param location Landfill location.
	 */
	async function createLandfill(location: GeoJSONFeaturePoint) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.POST("/landfills", {
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
			title: $t("landfills.create.success"),
			description: undefined,
		});

		goto(`${BackOfficeRoutes.LANDFILLS}/${res.data.id}`);
	}
</script>

<Card class="m-10 flex flex-col gap-10">
	<LandfillForm
		back={BackOfficeRoutes.LANDFILLS}
		title={$t("landfills.create.title")}
		isSubmitting={isSubmittingForm}
		onSave={createLandfill}
	/>
</Card>
