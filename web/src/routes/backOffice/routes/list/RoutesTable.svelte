<script lang="ts">
	import type { ComponentProps } from "svelte";
	import Table from "../../../../lib/components/table/Table.svelte";
	import TableDetailsAction from "../../../../lib/components/table/TableDetailsAction.svelte";
	import type { Columns } from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { t } from "../../../../lib/utils/i8n";
	import routesStore from "./routesStore";
	import { getLocationName } from "../../../../lib/utils/location";
	import type { Route } from "../../../../domain/route";
	import { getTruckName } from "../utils/truck";

	const { loading, data, filters } = routesStore;

	const columns: Columns<Route> = [
		{
			type: "accessor",
			field: "name",
			header: $t("route"),
			enableSorting: false,
			enableFiltering: false,
			cell(route) {
				return route;
			},
		},
		{
			type: "accessor",
			field: "departureWarehouse",
			header: $t("departure"),
			enableSorting: false,
			enableFiltering: false,
			cell(departureWarehouse) {
				const { municipalityName, wayName } =
					departureWarehouse.geoJson.properties;

				return getLocationName(wayName, municipalityName);
			},
		},
		{
			type: "accessor",
			field: "arrivalWarehouse",
			header: $t("arrival"),
			enableSorting: false,
			enableFiltering: false,
			cell(arrivalWarehouse) {
				const { municipalityName, wayName } =
					arrivalWarehouse.geoJson.properties;

				return getLocationName(wayName, municipalityName);
			},
		},
		{
			type: "accessor",
			field: "truck",
			header: $t("truck"),
			enableSorting: false,
			enableFiltering: false,
			cell(truck) {
				return getTruckName(truck.make, truck.model, truck.licensePlate);
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
	 * Handles changes of the routes table pages.
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
	rows={$data.routes}
	pagination={{
		name: $t("routes").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
