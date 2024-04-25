<script lang="ts">
	import type { ComponentProps } from "svelte";
	import type {
		Container,
		ContainerSortableFields,
	} from "../../../../domain/container";
	import Table from "../../../../lib/components/table/Table.svelte";
	import type {
		Columns,
		Sorting,
	} from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { MISSING_VALUE_SYMBOL } from "../../../../lib/constants/symbols";
	import { t } from "../../../../lib/utils/i8n";
	import containersStore from "./containersStore";
	import TableDetailsAction from "../../../../lib/components/table/TableDetailsAction.svelte";

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
			field: "geom",
			header: $t("containers.location"),
			enableSorting: false,
			cell() {
				// TODO: Display address of the container.
				return MISSING_VALUE_SYMBOL;
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

	/**
	 * Handles changes of the containers table sorting state.
	 * @param sorting New sorting state.
	 */
	function handleSortingChange(sorting: Sorting<ContainerSortableFields>) {
		filters.update(store => {
			return {
				...store,
				sorting,
			};
		});
	}
</script>

<Table
	{columns}
	loading={$loading}
	sorting={$filters.sorting}
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
