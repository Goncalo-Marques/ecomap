<script lang="ts">
	import type { Container } from "../../../domain/container";
	import BottomSheet from "../../../lib/components/BottomSheet.svelte";
	import Field from "../../../lib/components/Field.svelte";
	import { DateFormats } from "../../../lib/constants/date";
	import { formatDate } from "../../../lib/utils/date";
	import { t } from "../../../lib/utils/i8n";
	import { getLocationName } from "../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../constants/routes";

	/**
	 * The container whose information is displayed.
	 */
	export let container: Container;
</script>

<BottomSheet
	title={getLocationName(
		container.geoJson.properties.wayName,
		container.geoJson.properties.municipalityName,
	)}
	resourceLink={`${BackOfficeRoutes.CONTAINERS}/${container.id}`}
>
	<Field
		label={$t("containers.category")}
		value={$t(`containers.category.${container.category}`)}
	/>
	<Field
		label={$t("createdAt")}
		value={formatDate(container.createdAt, DateFormats.shortDateTime)}
	/>
	<Field
		label={$t("modifiedAt")}
		value={formatDate(container.modifiedAt, DateFormats.shortDateTime)}
	/>
</BottomSheet>
