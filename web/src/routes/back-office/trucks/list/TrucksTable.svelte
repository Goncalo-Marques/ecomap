<script lang="ts">
	import type { ComponentProps } from "svelte";

	import Table from "$lib/components/table/Table.svelte";
	import TableDetailsAction from "$lib/components/table/TableDetailsAction.svelte";
	import type { Columns, SortingDirection } from "$lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
	import { t } from "$lib/utils/i8n";

	import type { Truck, TruckSortableFields } from "../../../../domain/truck";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import trucksStore from "./trucksStore";

	const { loading, data, filters } = trucksStore;

	const columns: Columns<Truck> = [
		{
			type: "accessor",
			field: "licensePlate",
			header: $t("licensePlate"),
			enableSorting: true,
			enableFiltering: false,
			cell(licensePlate) {
				return licensePlate;
			},
		},
		{
			type: "accessor",
			field: "make",
			header: $t("make"),
			enableSorting: true,
			enableFiltering: false,
			cell(make) {
				return make;
			},
		},
		{
			type: "accessor",
			field: "model",
			header: $t("model"),
			enableSorting: true,
			enableFiltering: false,
			cell(model) {
				return model;
			},
		},
		{
			type: "accessor",
			field: "personCapacity",
			header: $t("personCapacity"),
			enableSorting: true,
			enableFiltering: false,
			cell(personCapacity) {
				return personCapacity.toString();
			},
		},
		{
			type: "display",
			header: "",
			align: "center",
			size: 120,
			cell(row) {
				const props: ComponentProps<TableDetailsAction> = {
					href: `${BackOfficeRoutes.TRUCKS}/${row.id}`,
				};

				return {
					component: TableDetailsAction,
					props,
				};
			},
		},
	];

	/**
	 * Handles changes of the trucks table pages.
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

	/**
	 * Handles changes of the trucks table sorting state.
	 * @param field Sorting field.
	 * @param direction Sorting direction.
	 */
	function handleSortingChange(
		field: TruckSortableFields,
		direction: SortingDirection,
	) {
		filters.update(store => {
			return {
				...store,
				sort: field,
				order: direction,
			};
		});
	}
</script>

<Table
	{columns}
	loading={$loading}
	rows={$data.trucks}
	sortingField={$filters.sort}
	sortingOrder={$filters.order}
	onSortingChange={handleSortingChange}
	pagination={{
		name: $t("trucks").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
