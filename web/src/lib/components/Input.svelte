<script lang="ts">
	import type {
		ChangeEventHandler,
		HTMLInputAttributes,
		MouseEventHandler,
	} from "svelte/elements";
	import Icon from "./Icon.svelte";

	/**
	 * The hint for form autofill feature.
	 * @default "off"
	 */
	export let autocomplete: HTMLInputAttributes["autocomplete"] = "off";

	/**
	 * The name of the icon placed at the end of the input.
	 * @default null
	 */
	export let endIcon: string | null = null;

	/**
	 * Indicates if the input contains an error.
	 * @default false
	 */
	export let error: boolean = false;

	/**
	 * Global attribute valid for all elements, including all the input types, it defines a unique identifier (ID) which must be unique in the whole document.
	 * @default null
	 */
	export let id: HTMLInputAttributes["id"] = null;

	/**
	 * Defines the maximum value that is acceptable and valid for the input.
	 * @default null
	 */
	export let max: HTMLInputAttributes["max"] = null;

	/**
	 * Defines the minimum value that is acceptable and valid for the input.
	 * @default null
	 */
	export let min: HTMLInputAttributes["min"] = null;

	/**
	 * The name of the form control. Submitted with the form as part of a name/value pair.
	 * @default null
	 */
	export let name: HTMLInputAttributes["name"] = null;

	/**
	 * Callback fired when input element is clicked.
	 * @default null
	 */
	export let onClick: MouseEventHandler<HTMLInputElement> | null = null;

	/**
	 * Callback fired when input value changes.
	 * @default null
	 */
	export let onInput: ChangeEventHandler<HTMLInputElement> | null = null;

	/**
	 * The text that appears in the form control when it has no value set.
	 * @default null
	 */
	export let placeholder: HTMLInputAttributes["placeholder"] = null;

	/**
	 * Indicates if input is not mutable.
	 * @default false
	 */
	export let readonly: boolean = false;

	/**
	 * The type of control to render.
	 * @default null
	 */
	export let type: HTMLInputAttributes["type"] = null;

	/**
	 * The value of the input.
	 * @default null
	 */
	export let value: number | string | null = null;
</script>

<div class="input-container">
	<input
		{autocomplete}
		class={error ? "error" : ""}
		{id}
		{name}
		{placeholder}
		{type}
		{value}
		{readonly}
		{max}
		{min}
		on:input={onInput}
		on:click={onClick}
		data-endIcon={endIcon}
	/>
	{#if endIcon}
		<div class="end-icon">
			<Icon name={endIcon} size="small" />
		</div>
	{/if}
</div>

<style>
	.input-container {
		position: relative;
	}

	input {
		width: 100%;
		padding: 0.375rem 0.5rem;
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
		color: var(--gray-900);

		&[data-endIcon] {
			padding-right: 2rem;
		}

		&:read-only {
			cursor: default;
		}

		&.error {
			border-color: var(--red-500);

			&:focus {
				outline-color: var(--red-500);
			}
		}

		&::placeholder {
			color: var(--gray-400);
		}
	}

	.end-icon {
		pointer-events: none;
		display: flex;
		align-items: center;
		color: var(--gray-900);
		position: absolute;
		top: 50%;
		right: 0.5rem;
		transform: translateY(-50%);
	}
</style>
