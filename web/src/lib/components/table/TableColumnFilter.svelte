<script
	lang="ts"
	generics="TRow extends Record<string, unknown>, TValue extends TRow[keyof TRow]"
>
	import { t } from "../../utils/i8n";
	import Button from "../Button.svelte";
	import Icon from "../Icon.svelte";
	import Popover from "../Popover.svelte";
	import Radio from "../Radio.svelte";
	import type { FilterOption } from "./types";

	/**
	 * The filter options of the column.
	 * @default []
	 */
	export let options: FilterOption<TValue>[] = [];

	/**
	 * Callback fired when a column filter changes.
	 * @default null
	 */
	export let onFilterChange: ((value: TValue | undefined) => void) | null =
		null;

	/**
	 * The initial value of the filter.
	 */
	export let initialValue: TValue | undefined;

	/**
	 * Internal value to manage the selected radio button.
	 */
	let selectedValue: TValue | undefined = initialValue;

	/**
	 * The popover element that contains the filter options.
	 */
	let popover: HTMLElement;

	/**
	 * Updates the selected value of the filter and fires the
	 * `onFilterChanges` callback.
	 * @param value Updated filter value.
	 */
	function updatedSelectedValue(value: TValue | undefined) {
		selectedValue = value;
		onFilterChange?.(value);
	}

	/**
	 * Clears the selected filter and closes the popover element.
	 */
	function clearFilter() {
		updatedSelectedValue(undefined);
		popover.hidePopover();
	}

	/**
	 * Handles change events of a radio button.
	 * @param e Change event.
	 */
	function handleRadioChange(e: Event) {
		const radioButton = e.currentTarget as HTMLInputElement;
		const value = radioButton.value as TValue;
		updatedSelectedValue(value);
		popover.hidePopover();
	}
</script>

<Popover bind:popover>
	<div slot="trigger" class="relative flex items-center justify-center">
		<Icon name="filter_alt" size="small" />
		<span
			class="absolute right-[2px] top-[2px] size-[0.375rem] rounded-full bg-green-700"
			style:display={selectedValue ? "" : "none"}
		/>
	</div>

	<div class="flex min-w-32 flex-col gap-4 p-4">
		<fieldset class="flex flex-col gap-1">
			{#each options as option}
				<Radio
					name="option"
					label={option.label}
					value={option.value}
					checked={option.value === selectedValue}
					onChange={handleRadioChange}
				/>
			{/each}
		</fieldset>

		<div class="flex">
			<Button
				class="flex-1"
				variant="primary"
				size="small"
				onClick={clearFilter}
			>
				{$t("clear")}
			</Button>
		</div>
	</div>
</Popover>
