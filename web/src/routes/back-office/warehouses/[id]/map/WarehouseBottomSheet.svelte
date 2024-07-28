<script lang="ts">
	import type { Warehouse } from "$domain/warehouse";
	import BottomSheet from "$lib/components/BottomSheet.svelte";
	import Field from "$lib/components/Field.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	/**
	 * The warehouse whose information is displayed.
	 */
	export let warehouse: Warehouse;
</script>

<BottomSheet
	title={getLocationName(
		warehouse.geoJson.properties.wayName,
		warehouse.geoJson.properties.municipalityName,
	)}
>
	<Field label={$t("truckCapacity")} value={warehouse.truckCapacity} />
	<Field
		label={$t("createdAt")}
		value={formatDate(warehouse.createdAt, DateFormats.shortDateTime)}
	/>
	<Field
		label={$t("modifiedAt")}
		value={formatDate(warehouse.modifiedAt, DateFormats.shortDateTime)}
	/>
</BottomSheet>
