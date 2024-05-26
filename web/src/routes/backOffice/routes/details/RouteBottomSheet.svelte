<script lang="ts">
	import type { Route } from "../../../../domain/route";
	import BottomSheet from "../../../../lib/components/BottomSheet.svelte";
	import Field from "../../../../lib/components/Field.svelte";
	import { DateFormats } from "../../../../lib/constants/date";
	import { formatDate } from "../../../../lib/utils/date";
	import { t } from "../../../../lib/utils/i8n";
	import { getLocationName } from "../../../../lib/utils/location";

	/**
	 * The route whose information is displayed.
	 */
	export let route: Route;
</script>

<BottomSheet title={route.name}>
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
	<Field
		label={$t("createdAt")}
		value={formatDate(route.createdAt, DateFormats.shortDateTime)}
	/>
	<Field
		label={$t("modifiedAt")}
		value={formatDate(route.modifiedAt, DateFormats.shortDateTime)}
	/>
</BottomSheet>
