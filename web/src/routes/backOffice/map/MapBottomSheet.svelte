<script lang="ts">
	import type { Container } from "../../../domain/container";
	import type { Truck } from "../../../domain/truck";
	import type { Warehouse } from "../../../domain/warehouse";
	import BottomSheet from "../../../lib/components/BottomSheet.svelte";
	import Field from "../../../lib/components/Field.svelte";
	import { DateFormats } from "../../../lib/constants/date";
	import { formatDate } from "../../../lib/utils/date";
	import { t } from "../../../lib/utils/i8n";
	import { getLocationName } from "../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../constants/routes";
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
	 * The way name of the location.
	 */
	export let wayName: string | undefined;

	/**
	 * The municipality name of the location.
	 */
	export let municipalityName: string | undefined;

	/**
	 * Retrieves the respective resource link to display in the bottom sheet.
	 * @param truck Selected truck.
	 * @param warehouse Selected warehouse.
	 * @param containers Selected containers.
	 * @returns Resource link or `undefined` when bottom sheet displays multiple containers.
	 */
	function getResourceLink(
		truck: Truck | null,
		warehouse: Warehouse | null,
		containers: Container[],
	) {
		if (truck) {
			return `${BackOfficeRoutes.TRUCKS}/${truck.id}`;
		}

		if (warehouse) {
			return `${BackOfficeRoutes.WAREHOUSES}/${warehouse.id}`;
		}

		if (containers.length === 1) {
			return `${BackOfficeRoutes.CONTAINERS}/${containers[0].id}`;
		}

		return undefined;
	}
</script>

<BottomSheet
	title={getLocationName(wayName, municipalityName)}
	resourceLink={getResourceLink(truck, warehouse, containers)}
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
	{:else if containers.length > 1}
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
	{:else}
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
	{/if}
</BottomSheet>
