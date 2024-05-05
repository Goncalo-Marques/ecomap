<script lang="ts">
	import { navigate } from "svelte-routing";
	import type { ContainerCategory } from "../../../../domain/container";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import ContainerForm from "../components/ContainerForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Creates a container with a given category and location.
	 * @param category Container category.
	 * @param location Container location.
	 */
	async function createContainer(
		category: ContainerCategory,
		location: GeoJSONFeaturePoint,
	) {
		const res = await ecomapHttpClient.POST("/containers", {
			body: {
				category,
				geoJson: location,
			},
		});

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

		navigate(`${BackOfficeRoutes.CONTAINERS}/${res.data.id}`);
	}
</script>

<Card class="page-layout">
	<ContainerForm
		back=""
		title={$t("containers.create.title")}
		onSave={createContainer}
	/>
</Card>
