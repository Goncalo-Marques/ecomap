<script lang="ts">
	import type { ComponentProps } from "svelte";
	import type { Container } from "../../../../domain/container";
	import Table from "../../../../lib/components/table/Table.svelte";
	import TableDetailsAction from "../../../../lib/components/table/TableDetailsAction.svelte";
	import type { Columns } from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { t } from "../../../../lib/utils/i8n";
	import { categoryOptions } from "../constants/category";
	import containersStore from "./containersStore";
	import { getContainerLocation } from "../utils/location";

	const { loading, data, filters } = containersStore;

	const columns: Columns<Container> = [
		{
			type: "accessor",
			field: "category",
			header: $t("containers.category"),
			enableSorting: false,
			enableFiltering: true,
			filterOptions: categoryOptions.map(category => {
				return {
					value: category,
					label: $t(`containers.category.${category}`),
				};
			}),
			filterInitialValue: $filters.category,
			cell(category) {
				return $t(`containers.category.${category}`);
			},
			onFilterChange(category) {
				filters.update(filters => {
					return {
						...filters,
						pageIndex: 0,
						category,
					};
				});
			},
		},
		{
			type: "accessor",
			field: "geoJson",
			header: $t("containers.location"),
			enableSorting: false,
			enableFiltering: false,
			cell(geoJson) {
				const {
					municipalityName = $t("containers.location.unknownMunicipality"),
					wayName = $t("containers.location.unknownWay"),
				} = geoJson.properties;

				return getContainerLocation(wayName, municipalityName);
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
	 * Handles changes of the containers table pages.
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
	rows={$data.containers}
	pagination={{
		name: $t("containers.title").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
