<script lang="ts">
	import type { ChangeEventHandler } from "svelte/elements";

	/**
	 * The label displayed for the checkbox.
	 * @default null
	 */
	export let label: string | null = null;

	/**
	 * The size of the checkbox.
	 * @default "medium"
	 */
	export let size: "small" | "medium" | "large" = "medium";

	/**
	 * Indicates if the checkbox is disabled.
	 * @default false
	 */
	export let disabled: boolean = false;

	/**
	 * Indicates if the checkbox is checked.
	 * @default false
	 */
	export let checked: boolean = false;

	/**
	 * Callback fired when checkbox checked state changes.
	 * @default null
	 */
	export let onChange: ChangeEventHandler<HTMLInputElement> | null = null;
</script>

<label data-size={size} data-disabled={disabled}>
	<input
		{disabled}
		{checked}
		type="checkbox"
		class="material-symbols-rounded"
		on:change={onChange}
	/>

	{#if label}
		{label}
	{/if}
</label>

<style>
	label {
		cursor: pointer;
		position: relative;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	input {
		appearance: none;
		cursor: inherit;
		line-height: normal;
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
		accent-color: var(--green-700);

		&:checked {
			position: relative;
			border: none;
			color: var(--white);
			background-color: var(--green-700);

			&::before {
				content: "check";
				position: absolute;
				top: 50%;
				left: 50%;
				transform: translate(-50%, -50%);
			}
		}

		&:hover {
			accent-color: var(--green-700);
		}
	}

	[data-size="small"] input {
		height: 0.75rem;
		width: 0.75rem;
		font-size: 0.5rem;
	}

	[data-size="medium"] input {
		height: 1rem;
		width: 1rem;
		font-size: 0.75rem;
	}

	[data-size="large"] input {
		height: 1.5rem;
		width: 1.5rem;
		font-size: 1rem;
	}

	[data-disabled="true"] {
		cursor: not-allowed;
		color: var(--gray-400);

		& input {
			background-color: var(--gray-300);
		}
	}
</style>
