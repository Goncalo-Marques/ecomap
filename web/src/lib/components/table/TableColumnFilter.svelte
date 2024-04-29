<script lang="ts" generics="TValue">
	import Button from "../Button.svelte";
	import Icon from "../Icon.svelte";
	import Popover from "../Popover.svelte";

	export let options: { value: TValue; label: string }[] = [];

	export let onValueChange: ((value: TValue) => void) | null = null;

	let popover: HTMLElement;

	let selectedValue: TValue | null = null;

	$: if (selectedValue) {
		onValueChange?.(selectedValue);
	}
</script>

<Popover bind:popover>
	<div slot="trigger" class="filter">
		<Icon name="filter_alt" size="small" />
		<span class="badge" style:display={!selectedValue ? "none" : ""} />
	</div>

	<div class="content">
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

		<div class="actions">
			<Button
				type="reset"
				variant="primary"
				size="small"
				onClick={() => {
					selectedValue = null;
					popover.hidePopover();
				}}
			>
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
		justify-content: flex-end;
	}

	fieldset {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	label {
		font: var(--text-base-regular);
	}
</style>
