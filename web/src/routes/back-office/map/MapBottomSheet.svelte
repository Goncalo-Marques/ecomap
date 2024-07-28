<script lang="ts">
	import type { Container } from "$domain/container";
	import type { Landfill } from "$domain/landfill";
	import type { Truck } from "$domain/truck";
	import type { Warehouse } from "$domain/warehouse";
	import BottomSheet from "$lib/components/BottomSheet.svelte";
	import Field from "$lib/components/Field.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import FeatureCard from "../components/FeatureCard.svelte";

	/**
	 * The containers in the location.
	 */
	export let containers: Container[];

	/**
	 * The truck in the location.
	 */
	export let truck: Truck | null;

	/**
	 * The warehouse in the location.
	 */
	export let warehouse: Warehouse | null;

	/**
	 * The landfill in the location.
	 */
	export let landfill: Landfill | null;

	/**
	 * The way name of the location.
	 */
	export let wayName: string | undefined;

	/**
	 * The municipality name of the location.
	 */
	export let municipalityName: string | undefined;

	/**
	 * Retrieves the respective resource href to display in the bottom sheet.
	 * @param truck Selected truck.
	 * @param warehouse Selected warehouse.
	 * @param landfill Selected landfill.
	 * @param containers Selected containers.
	 * @returns Resource href or `undefined` when bottom sheet displays multiple containers.
	 */
	function getResourceHref(
		truck: Truck | null,
		warehouse: Warehouse | null,
		landfill: Landfill | null,
		containers: Container[],
	) {
		if (truck) {
			return `${BackOfficeRoutes.TRUCKS}/${truck.id}`;
		}

		if (warehouse) {
			return `${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}`;
		}

		if (landfill) {
			return `${BackOfficeRoutes.LANDFILLS}/${landfill.id}`;
		}

		if (containers.length === 1) {
			return `${BackOfficeRoutes.CONTAINERS}/${containers[0].id}`;
		}

		return undefined;
	}
</script>

<BottomSheet
	title={getLocationName(wayName, municipalityName)}
	resourceHref={getResourceHref(truck, warehouse, landfill, containers)}
>
	{#if truck}
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
	{:else if warehouse}
		<Field label={$t("truckCapacity")} value={warehouse.truckCapacity} />
		<Field
			label={$t("createdAt")}
			value={formatDate(warehouse.createdAt, DateFormats.shortDateTime)}
		/>
		<Field
			label={$t("modifiedAt")}
			value={formatDate(warehouse.modifiedAt, DateFormats.shortDateTime)}
		/>
	{:else if landfill}
		<Field
			label={$t("createdAt")}
			value={formatDate(landfill.createdAt, DateFormats.shortDateTime)}
		/>
		<Field
			label={$t("modifiedAt")}
			value={formatDate(landfill.modifiedAt, DateFormats.shortDateTime)}
		/>
	{:else if containers.length === 1}
		{@const container = containers[0]}
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
	{:else}
		{#each containers as container}
			<FeatureCard
				icon="delete"
				title={$t("container")}
				resourceHref={`${BackOfficeRoutes.CONTAINERS}/${container.id}`}
			>
				<Field
					label={$t("containers.category")}
					value={$t(`containers.category.${container.category}`)}
				/>
			</FeatureCard>
		{/each}
	{/if}
</BottomSheet>
