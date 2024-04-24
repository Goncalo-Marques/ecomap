<script lang="ts">
	import type { HTMLInputAttributes } from "svelte/elements";

	/**
	 * The hint for form autofill feature.
	 * @default "off"
	 */
	export let autocomplete: HTMLInputAttributes["autocomplete"] = "off";

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * Indicates if the input contains an error.
	 * @default false
	 */
	export let error: boolean = false;

	/**
	 * The text below the input that contains additional information about how the input will be used.
	 * @default null
	 */
	export let helperText: string | null = null;

	/**
	 * Global attribute valid for all elements, including all the input types, it defines a unique identifier (ID) which must be unique in the whole document.
	 * @default null
	 */
	export let id: HTMLInputAttributes["id"] = null;

	/**
	 * The name of the form control. Submitted with the form as part of a name/value pair.
	 * @default null
	 */
	export let name: HTMLInputAttributes["name"] = null;

	/**
	 * The text that appears in the form control when it has no value set.
	 * @default null
	 */
	export let placeholder: HTMLInputAttributes["placeholder"] = null;

	/**
	 * The type of control to render.
	 * @default null
	 */
	export let type: HTMLInputAttributes["type"] = null;

	/**
	 * The label for the input.
	 * @default null
	 */
	export let label: string | null = null;
</script>

<label class={className}>
	{#if label}
		{label}
	{/if}
	<input
		{autocomplete}
		class={error ? "error" : ""}
		{id}
		{name}
		{placeholder}
		{type}
	/>
	{#if helperText}
		<span class={`helper-text ${error ? "error" : ""}`}>{helperText}</span>
	{/if}
</label>

<style>
	label {
		display: block;
		color: var(--gray-500);
	}

	input {
		width: 100%;
		padding: 0.375rem 0.5rem;
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
		color: var(--gray-900);

		&::placeholder {
			color: var(--gray-400);
		}

		&.error {
			border-color: var(--red-500);

			&:focus {
				outline-color: var(--red-500);
			}
		}
	}

	.helper-text {
		font: var(--text-xs-regular);
		color: var(--gray-400);

		&.error {
			color: var(--red-500);
		}
	}
</style>
