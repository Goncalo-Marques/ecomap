<script
	lang="ts"
	generics="TRow extends Record<string, unknown>, TValue extends TRow[keyof TRow]"
>
	import { t } from "../../utils/i8n";

	import type { FilterOption } from "./types";

	import Button from "../Button.svelte";
	import Icon from "../Icon.svelte";
	import Popover from "../Popover.svelte";
	import Radio from "../Radio.svelte";

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
	<div slot="trigger" class="filter">
		<Icon name="filter_alt" size="small" />
		<span class="badge" style:display={selectedValue ? "" : "none"} />
	</div>

	<div class="content">
		<fieldset>
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

		<div class="actions">
			<Button variant="primary" size="small" onClick={clearFilter}>
				{$t("clear")}
			</Button>
		</div>
	</div>
</Popover>

<style>
	.filter {
		position: relative;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.badge {
		position: absolute;
		top: 2px;
		right: 2px;
		height: 0.375rem;
		width: 0.375rem;
		background-color: var(--green-700);
		border-radius: 50%;
	}

	.content {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		min-width: 8rem;
	}

	.actions {
		display: flex;

		& > button {
			flex: 1;
		}
	}

	fieldset {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
</style>
