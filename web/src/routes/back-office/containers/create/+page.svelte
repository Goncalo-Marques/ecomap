<script lang="ts">
	import { goto } from "$app/navigation";
	import type { ContainerCategory } from "$domain/container";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../../constants/routes";
	import Card from "../../components/Card.svelte";
	import ContainerForm from "../components/ContainerForm.svelte";

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
	 * Creates a container with a given category and location.
	 * @param category Container category.
	 * @param location Container location.
	 */
	async function createContainer(
		category: ContainerCategory,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.POST("/containers", {
			body: {
				category,
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
			title: $t("containers.create.success"),
			description: undefined,
		});

		goto(`${BackOfficeRoutes.CONTAINERS}/${res.data.id}`);
	}
</script>

<Card class="m-10 flex flex-col gap-10">
	<ContainerForm
		back={BackOfficeRoutes.CONTAINERS}
		title={$t("containers.create.title")}
		isSubmitting={isSubmittingForm}
		onSave={createContainer}
	/>
</Card>
