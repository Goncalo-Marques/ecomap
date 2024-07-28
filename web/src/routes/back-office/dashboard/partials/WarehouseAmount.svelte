<script lang="ts">
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Spinner from "$lib/components/Spinner.svelte";
	import { t } from "$lib/utils/i8n";

	import Card from "../../components/Card.svelte";

	/**
	 * Retrieves the number of warehouses.
	 * @returns Number of warehouses.
	 */
	async function getWarehouseAmount() {
		const res = await ecomapHttpClient.GET("/warehouses");

		if (res.error) {
			throw new Error("Failed to retrieve warehouse amount");
		}

		return res.data.total;
	}

	const warehouseAmountPromise = getWarehouseAmount();
</script>

<Card element="article" class="[grid-area:warehouseAmount]">
	<h2 class="text-gray-500">{$t("dashboard.amountOfWarehouses")}</h2>
	<p class="text-xl font-semibold">
		{#await warehouseAmountPromise}
			<Spinner />
		{:then warehouseAmount}
			{warehouseAmount}
		{:catch}
			{$t("error.unexpected.title")}
		{/await}
	</p>
</Card>
