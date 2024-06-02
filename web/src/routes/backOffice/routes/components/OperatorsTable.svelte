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

	export let loading: boolean;

	export let operators: Employee[];

	export let selectedOperators: Employee[] = [];

	export let disabledOperators: Employee[];

	function checkboxColumn(selectedOps: Employee[], disabledOps: Employee[]) {
		return (employee: Employee): TableCellSvelteComponent => {
			const props: ComponentProps<Checkbox> = {
				size: "large",
				checked: selectedOps.some(operator => operator.id === employee.id),
				disabled: disabledOps.some(operator => operator.id === employee.id),
				onChange(e) {
					const newSelectedOperators = [...selectedOps];

					const checked = e.currentTarget.checked;
					if (checked) {
						newSelectedOperators.push(employee);
					} else {
						const selectedOperatorIndex = newSelectedOperators.findIndex(
							operator => operator.id === employee.id,
						);
						if (selectedOperatorIndex !== -1) {
							newSelectedOperators.splice(selectedOperatorIndex, 1);
						}
					}

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
				return checkboxColumn(selectedOperators, disabledOperators)(operator);
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

	$: columns[0].cell = checkboxColumn(selectedOperators, disabledOperators);
</script>

<Table {columns} {loading} rows={operators} />
