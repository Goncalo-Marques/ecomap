<script lang="ts">
	import type {
		Container,
		ContainerSortableFields,
	} from "../../../../domain/container";
	import Table from "../../../../lib/components/table/Table.svelte";
	import type {
		Columns,
		SortingDirection,
	} from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { t } from "../../../../lib/utils/i8n";
	import { categoryOptions } from "../constants/category";
	import containersStore from "./containersStore";

	const { loading, data, filters } = containersStore;

	const categoryFilters = categoryOptions.map(category => {
		return {
			value: category,
			label: $t(`containers.category.${category}`),
		};
	});

	const columns: Columns<Container> = [
		{
			type: "accessor",
			field: "category",
			header: $t("containers.category"),
			enableSorting: false,
			enableFiltering: true,
			filterOptions: categoryFilters,
			cell(category) {
				return $t(`containers.category.${category}`);
			},
			onFilterChange(value) {
				containersStore.filters.update(filters => {
					return {
						...filters,
						pageIndex: 0,
						category: value,
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
					municipalityName,
					wayName = $t("containers.location.unknownWay"),
				} = geoJson.properties;

				return `${wayName}, ${municipalityName}`;
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

	/**
	 * Handles changes of the containers table sorting state.
	 * @param field New sorting field.
	 * @param order New sorting order.
	 */
	function handleSortingChange(
		field: ContainerSortableFields,
		order: SortingDirection,
	) {
		filters.update(store => {
			return {
				...store,
				sort: field,
				order,
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
	onSortingChange={handleSortingChange}
/>
