<script lang="ts">
	import type { Truck } from "../../../../domain/truck";
	import BottomSheet from "../../../../lib/components/BottomSheet.svelte";
	import Field from "../../../../lib/components/Field.svelte";
	import { DateFormats } from "../../../../lib/constants/date";
	import { formatDate } from "../../../../lib/utils/date";
	import { t } from "../../../../lib/utils/i8n";
	import { getLocationName } from "../../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../../constants/routes";

	/**
	 * The truck whose information is displayed.
	 */
	export let truck: Truck;
</script>

<BottomSheet
	title={getLocationName(
		truck.geoJson.properties.wayName,
		truck.geoJson.properties.municipalityName,
	)}
	resourceLink={`${BackOfficeRoutes.TRUCKS}/${truck.id}`}
>
	<Field label={$t("make")} value={truck.make} />
	<Field label={$t("model")} value={truck.model} />
	<Field label={$t("licensePlate")} value={truck.licensePlate} />
	<Field label={$t("personCapacity")} value={truck.personCapacity} />
	<Field
		label={$t("createdAt")}
		value={formatDate(truck.createdAt, DateFormats.shortDateTime)}
	/>
	<Field
		label={$t("modifiedAt")}
		value={formatDate(truck.modifiedAt, DateFormats.shortDateTime)}
	/>
</BottomSheet>
