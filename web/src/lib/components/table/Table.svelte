<script
	lang="ts"
	generics="TRow extends Record<string, unknown>, TSortableFields extends string"
>
	import { t } from "../../utils/i8n";
	import Icon from "../Icon.svelte";
	import Spinner from "../Spinner.svelte";
	import TableColumnFilter from "./TableColumnFilter.svelte";
	import type {
		Columns,
		onSortingChangeFn,
		Pagination,
		SortingColumns,
		SortingDirection,
	} from "./types";
	import {
		getCell,
		getCellStyle,
		getColumnsSorting,
		getVisiblePages,
		toggleDirection,
	} from "./utils";

	/**
	 * The columns of the table.
	 */
	export let columns: Columns<TRow>;

	/**
	 * Callback fired when sorting state changes.
	 * @default null
	 */
	export let onSortingChange: onSortingChangeFn<TSortableFields> | null = null;

	/**
	 * The pagination configuration for the table.
	 * @default null
	 */
	export let pagination: Pagination | null = null;

	/**
	 * The rows to display in the table.
	 * @default []
	 */
	export let rows: TRow[] = [];

	/**
	 * Indicates if table is loading in data.
	 * @default false
	 */
	export let loading: boolean = false;

	/**
	 * The sorting field of the table.
	 * @default null
	 */
	export let sortingField: TSortableFields | null = null;

	/**
	 * The sorting order of the table.
	 * @default null
	 */
	export let sortingOrder: SortingDirection | null = null;

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * Map that contains the sorting state for each column.
	 * @example
	 * { id: "asc", name: undefined }
	 */
	let columnsSorting: SortingColumns<TRow>;

	/**
	 * Handles on click event for each table header cell sorting button.
	 * @param e Click event.
	 */
	function handleSortingClick(e: Event) {
		const sortingButton = e.currentTarget as HTMLButtonElement;
		const headerCell = sortingButton.parentElement;

		if (!headerCell) {
			return;
		}

		// Retrieves data attributes from the header cell element.
		const { field, sortable, direction } = headerCell.dataset;

		// Ignore sorting update when column is not a accessor type column.
		if (!field) {
			return;
		}

		// Ignore sorting update when column is not sortable.
		if (!sortable || sortable === "false") {
			return;
		}

		if (
			direction !== undefined &&
			direction !== "asc" &&
			direction !== "desc"
		) {
			throw new Error(
				`The sorting direction '${direction}' is invalid for the column field '${field}'`,
			);
		}

		const updatedSortingField = field as TSortableFields;
		const updatedSortingDirection = toggleDirection(direction);

		// Retrieve the updated columns sorting map with the updated sorting state.
		columnsSorting = getColumnsSorting(
			columns,
			updatedSortingField,
			updatedSortingDirection,
		);

		// Dispatch onSortingChange callback with the updated sorting state.
		onSortingChange?.(updatedSortingField, updatedSortingDirection);
	}

	/**
	 * Handles on click event for each page of the table.
	 * @param e Click event.
	 */
	function handlePageClick(e: Event) {
		const pageElement = e.currentTarget as HTMLButtonElement;
		const { active, index: pageIndex } = pageElement.dataset;

		// Prevent page update when page index is not defined or page index is not a number.
		if (!pageIndex || Number.isNaN(Number(pageIndex))) {
			return;
		}

		// Prevent page update when active is not defined or is already the selected page.
		if (!active || active === "true") {
			return;
		}

		pagination?.onPageChange(Number(pageIndex));
	}

	/**
	 * Handles on click event for the previous page button.
	 * @param e Click event.
	 */
	function handlePreviousPageClick(e: Event) {
		const previousPageElement = e.currentTarget as HTMLButtonElement;

		// Prevent page update when button is disabled.
		if (previousPageElement.disabled) {
			return;
		}

		pagination?.onPageChange(pagination.pageIndex - 1);
	}

	/**
	 * Handles on click event for the next page button.
	 * @param e Click event.
	 */
	function handleNextPageClick(e: Event) {
		const nextPageElement = e.currentTarget as HTMLButtonElement;

		// Prevent page update when button is disabled.
		if (nextPageElement.disabled) {
			return;
		}

		pagination?.onPageChange(pagination.pageIndex + 1);
	}

	// Re-constructs columnsSorting every time columns, sortingField or sortingOrder changes.
	$: columnsSorting = getColumnsSorting(columns, sortingField, sortingOrder);
</script>

<div class={`relative flex h-full flex-col ${className}`}>
	<table class="flex flex-1 flex-col">
		<thead class="flex-shrink-0 overflow-y-auto [scrollbar-gutter:stable]">
			<tr class="flex border-b border-gray-300">
				{#each columns as column}
					<th
						class="flex flex-shrink-0 flex-grow basis-0 items-center gap-2 overflow-hidden px-2 py-3 text-left font-semibold data-[align=left]:justify-start data-[align=right]:justify-end data-[align=center]:justify-center [&:not([data-align])]:justify-start"
						data-align={column.align}
						data-field={column.type === "accessor" ? column.field : null}
						data-sortable={column.type === "accessor"
							? column.enableSorting
							: null}
						data-direction={column.type === "accessor" && column.enableSorting
							? columnsSorting[column.field]
							: null}
						style={getCellStyle(column)}
					>
						<span class="truncate">{column.header}</span>
						{#if column.type === "accessor"}
							{#if column.enableSorting}
								{@const arrowDirection =
									columnsSorting[column.field] === "asc"
										? "upward"
										: "downward"}
								{@const sortingClass = columnsSorting[column.field]
									? "text-green-700"
									: ""}

								<button
									on:click={handleSortingClick}
									class={`flex items-center justify-center ${sortingClass}`}
								>
									<Icon name={`arrow_${arrowDirection}`} size="small" />
								</button>
							{/if}

							{#if column.enableFiltering}
								<TableColumnFilter
									options={column.filterOptions}
									initialValue={column.filterInitialValue}
									onFilterChange={column.onFilterChange}
								/>
							{/if}
						{/if}
					</th>
				{/each}
			</tr>
		</thead>
		<tbody
			class="flex-shrink flex-grow basis-0 overflow-y-auto overflow-x-hidden border-b border-gray-300 [scrollbar-gutter:stable]"
		>
			{#each rows as row}
				<tr class="flex border-b border-gray-300 last:border-b-0">
					{#each columns as column}
						{@const cell = getCell(column, row)}
						<td
							class="flex flex-shrink-0 flex-grow basis-0 items-center overflow-hidden px-2 py-3 text-left data-[align=left]:justify-start data-[align=right]:justify-end data-[align=center]:justify-center [&:not([data-align])]:justify-start"
							style={getCellStyle(column)}
							data-align={column.align}
							data-field={column.type === "accessor" ? column.field : null}
						>
							{#if typeof cell === "string"}
								<span class="truncate">{cell}</span>
							{:else if cell.slotContent}
								<svelte:component this={cell.component} {...cell.props}>
									{cell.slotContent}
								</svelte:component>
							{:else}
								<svelte:component this={cell.component} {...cell.props} />
							{/if}
						</td>
					{/each}
				</tr>
			{/each}
		</tbody>
	</table>
	{#if loading}
		<Spinner
			class="absolute left-1/2 top-1/2 z-10 -translate-x-1/2 -translate-y-1/2"
		/>
	{/if}
	{#if pagination}
		{@const start = pagination.pageSize * pagination.pageIndex + 1}
		{@const end = pagination.pageSize * (pagination.pageIndex + 1)}
		{@const pages =
			pagination.total > 0
				? Math.ceil(pagination.total / pagination.pageSize)
				: 1}
		{@const pagesArray = Array.from({ length: pages }, (_, idx) => idx)}
		{@const visiblePages = getVisiblePages(pagesArray, pagination.pageIndex)}

		<div class="flex items-center gap-2 pt-4">
			<span class="flex-1 text-sm">
				{start > pagination.total ? pagination.total : start}-{end >
				pagination.total
					? pagination.total
					: end}
				{$t("pagination.of")}
				{pagination.total}
				{pagination.name}
			</span>

			<div class="flex items-center justify-center gap-2">
				<button
					class="flex items-center justify-center disabled:opacity-60"
					disabled={pagination.pageIndex === 0}
					on:click={handlePreviousPageClick}
				>
					<Icon name="arrow_back" size="x-small" />
				</button>

				{#each visiblePages as pageIndex}
					<button
						class="size-6 rounded text-xs font-semibold data-[active=true]:bg-green-700 data-[active=true]:text-white"
						data-index={pageIndex}
						data-active={pageIndex === pagination.pageIndex}
						on:click={handlePageClick}
					>
						{pageIndex + 1}
					</button>
				{/each}

				<button
					class="flex items-center justify-center disabled:opacity-60"
					disabled={pagination.pageIndex === pages - 1}
					on:click={handleNextPageClick}
				>
					<Icon name="arrow_forward" size="x-small" />
				</button>
			</div>
		</div>
	{/if}
</div>
