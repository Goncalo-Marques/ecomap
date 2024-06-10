<script lang="ts">
	import { type ComponentProps } from "svelte";
	import type { Employee } from "../../../../domain/employees";
	import Table from "../../../../lib/components/table/Table.svelte";
	import type {
		Columns,
		TableCellSvelteComponent,
	} from "../../../../lib/components/table/types";
	import { formatTime } from "../../../../lib/utils/date";
	import { t } from "../../../../lib/utils/i8n";
	import { getLocationName } from "../../../../lib/utils/location";
	import Checkbox from "../../../../lib/components/Checkbox.svelte";

	/**
	 * Indicates whether the operator data is loading.
	 */
	export let loading: boolean;

	/**
	 * The operators available.
	 */
	export let operators: Employee[];

	/**
	 * The operators that are selected.
	 * @default []
	 */
	export let selectedOperators: Employee[] = [];

	/**
	 * The operators that cannot be selected.
	 * @default []
	 */
	export let disabledOperators: Employee[] = [];

	/**
	 * The checkbox cell displayed for each operator.
	 *
	 * The cell function for the checkbox is extracted from the table columns to dynamically
	 * update the checkbox state as the selected and disabled operators change.
	 *
	 * @param selectedOps Selected operators.
	 * @param disabledOps Disabled operators.
	 * @returns Checkbox cell.
	 */
	function checkboxCell(selectedOps: Employee[], disabledOps: Employee[]) {
		return (operator: Employee): TableCellSvelteComponent => {
			const props: ComponentProps<Checkbox> = {
				size: "large",
				checked: selectedOps.some(selectedOp => selectedOp.id === operator.id),
				disabled: disabledOps.some(disabledOp => disabledOp.id === operator.id),
				onChange(e) {
					const newSelectedOperators = [...selectedOps];

					const checked = e.currentTarget.checked;
					if (checked) {
						// Add operator to selected operators list.
						newSelectedOperators.push(operator);
					} else {
						// Retrieve the index of the selected operator that needs to be removed.
						const selectedOperatorIndex = newSelectedOperators.findIndex(
							selectedOperator => selectedOperator.id === operator.id,
						);

						// If the operator was found, remove it from the list.
						if (selectedOperatorIndex !== -1) {
							newSelectedOperators.splice(selectedOperatorIndex, 1);
						}
					}

					// Update selected operators.
					selectedOperators = newSelectedOperators;
				},
			};

			return {
				component: Checkbox,
				props,
			};
		};
	}

	const columns: Columns<Employee> = [
		{
			type: "display",
			header: "",
			size: 40,
			cell(operator) {
				return checkboxCell(selectedOperators, disabledOperators)(operator);
			},
		},
		{
			type: "accessor",
			field: "firstName",
			enableSorting: false,
			enableFiltering: false,
			header: $t("employees.firstName"),
			cell(firstName) {
				return firstName;
			},
		},
		{
			type: "accessor",
			field: "lastName",
			enableSorting: false,
			enableFiltering: false,
			header: $t("employees.lastName"),
			cell(lastName) {
				return lastName;
			},
		},
		{
			type: "accessor",
			field: "scheduleStart",
			header: $t("employees.schedule"),
			enableSorting: false,
			enableFiltering: false,
			cell(scheduleStart, row) {
				return `${formatTime(scheduleStart)} - ${formatTime(row.scheduleEnd)}`;
			},
		},
		{
			type: "accessor",
			field: "geoJson",
			header: $t("location"),
			enableSorting: false,
			enableFiltering: false,
			cell(geoJson) {
				const { municipalityName, wayName } = geoJson.properties;

				return getLocationName(wayName, municipalityName);
			},
		},
	];

	// Reconstruct the checkbox cell each time the selected or disabled operators are changed.
	$: columns[0].cell = checkboxCell(selectedOperators, disabledOperators);
</script>

<Table {columns} {loading} rows={operators} />
