<script lang="ts" generics="TValue">
	import Icon from "../Icon.svelte";
	import Popover from "../Popover.svelte";

	export let options: { value: TValue; label: string }[] = [];

	export let onValueChange: ((value: TValue) => void) | null = null;

	let selectedValue: TValue | null = null;

	$: if (selectedValue) {
		onValueChange?.(selectedValue);
	}
</script>

<Popover>
	<Icon slot="trigger" name="filter_alt" size="x-small" />

	<fieldset>
		{#each options as option}
			<label>
				<input
					type="radio"
					name="option"
					value={option.value}
					bind:group={selectedValue}
				/>
				{option.label}
			</label>
		{/each}
	</fieldset>
</Popover>

<style>
	fieldset {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	label {
		font: var(--text-base-regular);
	}
</style>
