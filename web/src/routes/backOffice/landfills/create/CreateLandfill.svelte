<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import LandfillForm from "../components/LandfillForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Creates a landfill with a given location.
	 * @param location Landfill location.
	 */
	async function createLandfill(location: GeoJSONFeaturePoint) {
		const res = await ecomapHttpClient.POST("/landfills", {
			body: {
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
			title: $t("landfills.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.LANDFILLS}/${res.data.id}`);
	}
</script>

<Card class="page-layout">
	<LandfillForm
		back=""
		title={$t("landfills.create.title")}
		onSave={createLandfill}
	/>
</Card>
