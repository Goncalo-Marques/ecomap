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
	import containersStore from "./containersStore";

	const { loading, data, filters } = containersStore;

	const columns: Columns<Container> = [
		{
			type: "accessor",
			field: "category",
			header: $t("containers.category"),
			enableSorting: true,
			cell(category) {
				return $t(`containers.category.${category}`);
			},
		},
		{
			type: "accessor",
			field: "geoJson",
			header: $t("containers.location"),
			enableSorting: false,
			cell(geoJson) {
				const { municipalityName, wayName } = geoJson.properties;

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
	sortingField={$filters.sort}
	sortingOrder={$filters.order}
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
