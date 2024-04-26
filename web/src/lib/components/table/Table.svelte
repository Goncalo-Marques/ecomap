<script
	lang="ts"
	generics="TRow extends Record<string, unknown>, TSortableFields extends string"
>
	import { t } from "../../utils/i8n";

	import Icon from "../Icon.svelte";
	import {
		getCell,
		getCellStyle,
		getColumnsSorting,
		toggleDirection,
	} from "./utils";
	import type {
		Columns,
		Pagination,
		SortingColumns,
		SortingDirection,
		onSortingChangeFn,
	} from "./types";
	import Spinner from "../Spinner.svelte";

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
	 * Map that contains the sorting state for each column.
	 * @example
	 * { id: "asc", name: undefined }
	 */
	let columnsSorting: SortingColumns<TRow>;

	/**
	 * Handles on click event for each table header cell.
	 * @param e Click event.
	 */
	function handleHeaderCellClick(e: Event) {
		const headerCell = e.currentTarget as HTMLTableCellElement;

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

<div class="table-container">
	<table>
		<thead>
			<tr>
				{#each columns as column}
					<th
						align={column.align}
						data-field={column.type === "accessor" ? column.field : null}
						data-sortable={column.type === "accessor"
							? column.enableSorting
							: null}
						data-direction={column.type === "accessor" && column.enableSorting
							? columnsSorting[column.field]
							: null}
						style={getCellStyle(column)}
						on:click={handleHeaderCellClick}
					>
						{column.header}
						{#if column.type === "accessor" && column.enableSorting}
							{@const arrowDirection =
								columnsSorting[column.field] === "asc" ? "upward" : "downward"}
							{@const sortingClass = columnsSorting[column.field]
								? "sorted"
								: undefined}

							<button class={sortingClass}>
								<Icon name={`arrow_${arrowDirection}`} size="small" />
							</button>
						{/if}
					</th>
				{/each}
			</tr>
		</thead>
		<tbody>
			{#each rows as row}
				<tr>
					{#each columns as column}
						{@const cell = getCell(column, row)}

						<td
							align={column.align}
							style={getCellStyle(column)}
							data-field={column.type === "accessor" ? column.field : null}
						>
							{#if typeof cell === "string"}
								{cell}
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
		<Spinner class="table-loading-spinner" />
	{/if}
	{#if pagination}
		{@const start = pagination.pageSize * pagination.pageIndex + 1}
		{@const end = pagination.pageSize * (pagination.pageIndex + 1)}
		{@const pages = Math.ceil(pagination.total / pagination.pageSize)}
		{@const pagesArray = Array.from({ length: pages }, (_, idx) => idx)}

		<div class="pagination">
			<span class="pagination-info">
				{start > pagination.total ? pagination.total : start}-{end >
				pagination.total
					? pagination.total
					: end}
				{$t("pagination.of")}
				{pagination.total}
				{pagination.name}
			</span>

			<div class="pagination-pages">
				<button
					class="pagination-page-previous"
					disabled={pagination.pageIndex === 0}
					on:click={handlePreviousPageClick}
				>
					<Icon name="arrow_back" size="small" />
				</button>

				{#each pagesArray as pageIndex}
					<button
						class="pagination-page"
						data-index={pageIndex}
						data-active={pageIndex === pagination.pageIndex}
						on:click={handlePageClick}
					>
						{pageIndex + 1}
					</button>
				{/each}

				<button
					class="pagination-page-next"
					disabled={pagination.pageIndex === pages - 1}
					on:click={handleNextPageClick}
				>
					<Icon name="arrow_forward" size="small" />
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	.table-container {
		position: relative;
		display: flex;
		flex-direction: column;
		height: 100%;
	}
	table {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
	thead,
	tbody {
		overflow-x: hidden;
		overflow-y: auto;
		scrollbar-gutter: stable;
	}
	thead {
		flex-shrink: 0;
	}
	tbody {
		flex: 1 1 0;
		border-bottom: 1px solid var(--gray-300);
	}
	tr {
		display: flex;
		border-bottom: 1px solid var(--gray-300);
	}
	tbody tr:last-child {
		border-bottom: none;
	}
	th,
	td {
		display: flex;
		align-items: center;
		justify-content: flex-start;
		flex: 1 0;
		padding: 0.75rem 0.5rem;
		text-align: left;

		&[align="center"] {
			justify-content: center;
		}
		&[align="right"] {
			justify-content: flex-end;
		}
	}
	th {
		gap: 0.5rem;
		font: var(--text-base-semibold);
	}
	th :global(.sorted) {
		color: var(--green-700);
	}
	td {
		overflow: hidden;
		text-overflow: ellipsis;
		text-wrap: nowrap;
	}

	.table-container :global(.table-loading-spinner) {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		z-index: 10;
	}

	.pagination {
		display: flex;
		align-items: center;
		padding-top: 1rem;
		gap: 0.5rem;
	}
	.pagination-info {
		font: var(--text-sm-regular);
		flex: 1;
	}
	.pagination-pages {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.5rem;
	}
	.pagination-page {
		width: 1.5rem;
		height: 1.5rem;
		border-radius: 0.25rem;
		font: var(--text-xs-semibold);

		&[data-active="true"] {
			background-color: var(--green-700);
			color: var(--white);
		}
	}
	.pagination-page-previous,
	.pagination-page-next {
		display: flex;
		justify-content: center;
		align-items: center;

		&:disabled {
			opacity: 0.6;
		}
	}
</style>
