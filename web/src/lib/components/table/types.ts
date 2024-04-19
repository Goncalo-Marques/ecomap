import type { ComponentType } from "svelte";

/**
 * Base structure of the table row.
 */
export type Row = Record<string, unknown>;

/**
 * Base structure of a sortable field.
 */
export type SortableField = string;

/**
 * Table columns.
 * @template TRow Type of the table row.
 */
export type Columns<TRow extends Row> = Column<TRow>[];

/**
 * Table column.
 * @template TRow Type of the table row.
 */
export type Column<TRow extends Row> = {
	[Field in keyof TRow]:
		| ColumnAccessor<TRow, Field, TRow[Field]>
		| ColumnDisplay<TRow>;
}[keyof TRow];

/**
 * Table cell Svelte component.
 */
export interface TableCellSvelteComponent {
	/**
	 * The Svelte component to be displayed inside the table cell.
	 */
	component: ComponentType;

	/**
	 * The props to pass into the Svelte component.
	 */
	props: Record<string, unknown>;

	/**
	 * The slot content to be displayed inside the Svelte component.
	 * Needs to be set when the Svelte component accepts `<slot />`.
	 */
	slotContent?: string;
}

/**
 * Table accessor column.
 * @template TRow Type of the table row.
 * @template TField Type of the cell field.
 * @template TValue Type of the cell value.
 */
export interface ColumnAccessor<TRow, TField, TValue> extends ColumnCommon {
	/**
	 * The type of the column.
	 *
	 * Accessor columns have an underlying data model which means they can be sorted, filtered, grouped, etc.
	 */
	type: "accessor";

	/**
	 * The field in the row.
	 */
	field: TField;

	/**
	 * The content displayed inside a table body cell.
	 * @param value Value of the cell.
	 * @param row Row to which the cell belongs to.
	 * @returns Content to be displayed inside the table body cell.
	 */
	cell(value: TValue, row: TRow): string | TableCellSvelteComponent;

	/**
	 * Indicates if the column is sortable.
	 */
	enableSorting: boolean;
}

/**
 * Table display column.
 * @template TRow Type of the table row.
 */
export interface ColumnDisplay<TRow> extends ColumnCommon {
	/**
	 * The type of the column.
	 *
	 * Display columns do not have a data model which means they cannot be sorted, filtered, etc, but they
	 * can be used to display arbitrary content in the table, eg. a row actions button, checkbox, expander, etc.
	 */
	type: "display";

	/**
	 * The content displayed inside a table body cell.
	 * @param row Row to which the cell belongs to.
	 * @returns Content to be displayed inside the table body cell.
	 */
	cell(row: TRow): string | TableCellSvelteComponent;
}

/**
 * Table column common properties.
 */
export interface ColumnCommon {
	/**
	 * The content displayed inside a table header cell.
	 */
	header: string;

	/**
	 * The alignment of the content in the column.
	 * @default "left"
	 */
	align?: "left" | "center" | "right";

	/**
	 * The minimum size of the column.
	 */
	minSize?: number;

	/**
	 * The fixed size of the column.
	 */
	size?: number;

	/**
	 * The maximum size of the column.
	 */
	maxSize?: number;
}

/**
 * Callback fired when sorting state changes.
 * @template TSortableFields Type of the sortable fields.
 * @param sorting Updated sorting state.
 */
export type onSortingChangeFn<TSortableFields extends SortableField> = (
	sorting: Sorting<TSortableFields>,
) => void;

/**
 * The available directions of a column sorting state.
 */
export type SortingDirection = "asc" | "desc" | undefined;

/**
 * Sorting state.
 * @template TSortableFields Type of the sortable fields.
 */
export interface Sorting<TSortableFields extends SortableField> {
	/**
	 * The name of the row field with the sorting.
	 */
	field: TSortableFields;

	/**
	 * The direction of the sort of the field.
	 */
	direction: SortingDirection;
}

/**
 * Columns sorting state.
 * @template TRow Type of the table row.
 */
export type SortingColumns<TRow extends Row> = {
	[Field in keyof TRow]: SortingDirection;
};

/**
 * Table pagination.
 */
export interface Pagination {
	/**
	 * The selected page index.
	 */
	pageIndex: number;

	/**
	 * The page size of the table.
	 */
	pageSize: number;

	/**
	 * The total number of rows in all pages.
	 */
	total: number;

	/**
	 * The name for the rows being displayed.
	 * @example "users"
	 */
	name: string;

	/**
	 * Callback fired when the page changes.
	 * @param pageIndex New page index.
	 */
	onPageChange(pageIndex: number): void;
}
