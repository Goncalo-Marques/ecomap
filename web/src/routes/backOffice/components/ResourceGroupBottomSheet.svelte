<script lang="ts">
	import type { ResourceGroupLocation } from "../../../domain/map";
	import BottomSheet from "../../../lib/components/BottomSheet.svelte";
	import Field from "../../../lib/components/Field.svelte";
	import { t } from "../../../lib/utils/i8n";
	import { getLocationName } from "../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../constants/routes";
	import FeatureCard from "./FeatureCard.svelte";

	/**
	 * The grouped resources.
	 */
	export let group: ResourceGroupLocation;

	const { containers, trucks, warehouses, wayName, municipalityName } = group;
</script>

<BottomSheet title={getLocationName(wayName, municipalityName)}>
	{#each containers as container}
		<FeatureCard
			icon="delete"
			title={$t("container")}
			resourceLink={`${BackOfficeRoutes.CONTAINERS}/${container.id}`}
		>
			<Field
				label={$t("containers.category")}
				value={$t(`containers.category.${container.category}`)}
			/>
		</FeatureCard>
	{/each}

	{#each trucks as truck}
		<FeatureCard
			icon="local_shipping"
			title={$t("truck")}
			resourceLink={`${BackOfficeRoutes.TRUCKS}/${truck.id}`}
		>
			<Field label={$t("licensePlate")} value={truck.licensePlate} />
			<Field label={$t("personCapacity")} value={truck.personCapacity} />
		</FeatureCard>
	{/each}

	{#each warehouses as warehouse}
		<FeatureCard
			icon="warehouse"
			title={$t("warehouse")}
			resourceLink={`${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}`}
		>
			<Field label={$t("truckCapacity")} value={warehouse.truckCapacity} />
		</FeatureCard>
	{/each}
</BottomSheet>
