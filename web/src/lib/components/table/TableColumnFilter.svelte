<script
	lang="ts"
	generics="TRow extends Record<string, unknown>, TValue extends TRow[keyof TRow]"
>
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
	 * The selected radio button value.
	 */
	let selectedValue: TValue | undefined;

	/**
	 * The popover element that contains the filter options.
	 */
	let popover: HTMLElement;

	/**
	 * Clears the selected filter and closes the popover element.
	 */
	function clearFilter() {
		selectedValue = undefined;
		popover.hidePopover();
	}

	// Dispatch callback every time selectedValue changes.
	$: onFilterChange?.(selectedValue);
</script>

<Popover bind:popover>
	<div slot="trigger" class="filter">
		<Icon name="filter_alt" size="small" />
		<span class="badge" style:display={!selectedValue ? "none" : ""} />
	</div>

	<div class="content">
		<fieldset>
			{#each options as option}
				<Radio
					label={option.label}
					value={option.value}
					bind:group={selectedValue}
				/>
			{/each}
		</fieldset>

		<div class="actions">
			<Button variant="primary" size="small" onClick={clearFilter}>
				Clear
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
