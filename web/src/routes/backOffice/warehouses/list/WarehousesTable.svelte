<script lang="ts">
	import type { ComponentProps } from "svelte";
	import Table from "../../../../lib/components/table/Table.svelte";
	import TableDetailsAction from "../../../../lib/components/table/TableDetailsAction.svelte";
	import type { Columns } from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { t } from "../../../../lib/utils/i8n";
	import warehousesStore from "./warehousesStore";
	import { getLocationName } from "../../../../lib/utils/location";
	import type { Warehouse } from "../../../../domain/warehouse";

	const { loading, data, filters } = warehousesStore;

	const columns: Columns<Warehouse> = [
		{
			type: "accessor",
			field: "geoJson",
			header: $t("location"),
			enableSorting: false,
			enableFiltering: false,
			cell(geoJson) {
				const { municipalityName, wayName } = geoJson.properties;

				return getLocationName(wayName, municipalityName);
			},
		},
		{
			type: "accessor",
			field: "truckCapacity",
			header: $t("truckCapacity"),
			enableSorting: false,
			enableFiltering: false,
			cell(truckCapacity) {
				return truckCapacity.toString();
			},
		},
		{
			type: "display",
			header: "",
			align: "center",
			size: 120,
			cell(row) {
				const props: ComponentProps<TableDetailsAction> = {
					id: row.id,
				};

				return {
					component: TableDetailsAction,
					props,
				};
			},
		},
	];

	/**
	 * Handles changes of the warehouses table pages.
	 * @param pageIndex New page index.
	 */
	function handlePageChange(pageIndex: number) {
		filters.update(store => {
			return {
				...store,
				pageIndex,
			};
		});
	}
</script>

<Table
	{columns}
	loading={$loading}
	rows={$data.warehouses}
	pagination={{
		name: $t("warehouses").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
