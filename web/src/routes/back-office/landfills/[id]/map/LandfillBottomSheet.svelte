<script lang="ts">
	import type { Landfill } from "$domain/landfill";
	import BottomSheet from "$lib/components/BottomSheet.svelte";
	import Field from "$lib/components/Field.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	/**
	 * The landfill whose information is displayed.
	 */
	export let landfill: Landfill;
</script>

<BottomSheet
	title={getLocationName(
		landfill.geoJson.properties.wayName,
		landfill.geoJson.properties.municipalityName,
	)}
>
	<Field
		label={$t("createdAt")}
		value={formatDate(landfill.createdAt, DateFormats.shortDateTime)}
	/>
	<Field
		label={$t("modifiedAt")}
		value={formatDate(landfill.modifiedAt, DateFormats.shortDateTime)}
	/>
</BottomSheet>
