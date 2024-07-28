import { ROOT_FONT_SIZE } from "../../constants/fontSize";
import type {
	Column,
	Columns,
	Row,
	SortableField,
	SortingColumns,
	SortingDirection,
	TableCellSvelteComponent,
} from "./types";

/**
 * Retrieves the list of fields in a table's columns.
 * @template TRow Type of the table row.
 * @param columns Table columns.
 * @returns List of fields in a table's columns.
 */
export function getColumnsFields<TRow extends Row>(
	columns: Columns<TRow>,
): (keyof TRow)[] {
	const fields: (keyof TRow)[] = [];

	for (const column of columns) {
		if (column.type === "accessor") {
			fields.push(column.field);
		}
	}

	return fields;
}

/**
 * Retrieves a map with the sorting state for each table column.
 * @template TRow Type of the table row.
 * @template TSortableFields Type of the sortable fields.
 * @param columns Table columns.
 * @param sorting Current sorting state.
 * @returns Map with the sorting state for each table column.
 */
export function getColumnsSorting<
	TRow extends Row,
	TSortableFields extends SortableField,
>(
	columns: Columns<TRow>,
	sortingField: TSortableFields | null,
	sortingDirection: SortingDirection | null,
): SortingColumns<TRow> {
	const fields = getColumnsFields(columns);
	const columnsSorting = {} as SortingColumns<TRow>;

	for (const field of fields) {
		if (sortingField && sortingDirection && field === sortingField) {
			columnsSorting[field] = sortingDirection;
		} else {
			columnsSorting[field] = undefined;
		}
	}

	return columnsSorting;
}

/**
 * Toggles a column sorting direction into its next state.
 *
 * Sorting state transitions:
 * - `undefined` -> `asc`
 * - `asc` -> `desc`
 * - `desc` -> `asc`
 *
 * @param direction Current sorting direction.
 * @returns New sorting direction.
 */
export function toggleDirection(
	direction: SortingDirection | undefined,
): SortingDirection {
	let newDirection: SortingDirection;

	if (direction === "asc") {
		newDirection = "desc";
	} else {
		newDirection = "asc";
	}

	return newDirection;
}

/**
 * Retrieves a cell content given its column and row data.
 * @template TRow Type of the table row.
 * @param column Table column.
 * @param row Table row.
 * @returns Cell content.
 */
export function getCell<TRow extends Row>(
	column: Column<TRow>,
	row: TRow,
): string | TableCellSvelteComponent {
	if (column.type === "display") {
		return column.cell(row);
	}

	let value: TRow[keyof TRow] | null = null;

	for (const field of Object.keys(row) as (keyof TRow)[]) {
		if (field === column.field) {
			value = row[field];
			break;
		}
	}

	if (value === null) {
		throw new Error(
			`Cell value was not found for column '${String(column.field)}'`,
		);
	}

	return column.cell(value, row);
}

/**
 * Retrieves a cell `style` attribute.
 * @template TRow Type of the table row.
 * @param column Cell column.
 * @returns Cell style.
 */
export function getCellStyle<TRow extends Row>(column: Column<TRow>): string {
	let minWidth = "min-width: auto";
	if (column.minSize !== undefined) {
		minWidth = `min-width: ${column.minSize / ROOT_FONT_SIZE}rem`;
	}

	let maxWidth = "max-width: none";
	if (column.maxSize !== undefined) {
		maxWidth = `max-width: ${column.maxSize / ROOT_FONT_SIZE}rem`;
	}

	let width = "";
	if (column.size !== undefined) {
		width = `flex: 0 0 ${column.size / ROOT_FONT_SIZE}rem`;
	}

	return `${minWidth}; ${maxWidth}; ${width}`;
}

/**
 * Retrieves the list of the pages to display in the table pagination.
 * @param pages Pagination pages.
 * @param pageIndex Page index.
 * @param [size=5] Size of visible pages.
 * @returns List of the pages to display in the table pagination.
 */
export function getVisiblePages(
	pages: number[],
	pageIndex: number,
	size: number = 5,
): number[] {
	const half = Math.floor(size / 2);

	const remainderRightSide =
		half - Math.min(half, pages.length - 1 - pageIndex);
	const remainderLeftSide = Math.max(0, half - pageIndex);

	const startIndex = Math.max(0, pageIndex - half - remainderRightSide);
	const endIndex = Math.min(
		pages.length - 1,
		pageIndex + half + remainderLeftSide,
	);

	return pages.slice(startIndex, endIndex + 1);
}
