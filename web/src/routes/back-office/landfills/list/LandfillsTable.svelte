<script lang="ts">
	import type { ComponentProps } from "svelte";

	import type { Landfill } from "$domain/landfill";
	import Table from "$lib/components/table/Table.svelte";
	import TableDetailsAction from "$lib/components/table/TableDetailsAction.svelte";
	import type { Columns } from "$lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import { BackOfficeRoutes } from "../../../constants/routes";
	import landfillsStore from "./landfillsStore";

	const { loading, data, filters } = landfillsStore;

	const columns: Columns<Landfill> = [
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
			type: "display",
			header: "",
			align: "center",
			size: 120,
			cell(row) {
				const props: ComponentProps<TableDetailsAction> = {
					href: `${BackOfficeRoutes.LANDFILLS}/${row.id}`,
				};

				return {
					component: TableDetailsAction,
					props,
				};
			},
		},
	];

	/**
	 * Handles changes of the landfills table pages.
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
	rows={$data.landfills}
	pagination={{
		name: $t("landfills").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
