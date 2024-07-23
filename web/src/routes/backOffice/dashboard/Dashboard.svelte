<script lang="ts">
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import { t } from "../../../lib/utils/i8n";
	import { getBatchPaginatedResponse } from "../../../lib/utils/request";
	import ContainersAdded from "./partials/ContainersAdded.svelte";
	import ContainersByCategory from "./partials/ContainersByCategory.svelte";
	import ContainersByMunicipality from "./partials/ContainersByMunicipality.svelte";
	import TruckAmount from "./partials/TruckAmount.svelte";
	import WarehouseAmount from "./partials/WarehouseAmount.svelte";
	import ActiveEmployees from "./partials/ActiveEmployees.svelte";

	/**
	 * Retrieves the containers to be displayed in the charts.
	 * @returns Containers.
	 */
	async function getContainers() {
		const containers = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/containers", {
					params: {
						query: {
							offset,
							limit,
							sort: "createdAt",
							order: "asc",
						},
					},
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.containers };
			},
		);

		return containers;
	}

	let containersPromise = getContainers();
</script>

<main class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("dashboard")}</h1>
	<div
		class="grid grid-cols-3 gap-4 [grid-template-areas:'activeEmployees_warehouseAmount_truckAmount''containersAdded_containersAdded_containersByCategory''containersByMunicipality_containersByMunicipality_containersByMunicipality']"
	>
		<ActiveEmployees />
		<WarehouseAmount />
		<TruckAmount />
		<ContainersAdded {containersPromise} />
		<ContainersByCategory {containersPromise} />
		<ContainersByMunicipality {containersPromise} />
	</div>
</main>
