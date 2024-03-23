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
	let className: HTMLInputAttributes["class"] = "";
	export { className as class };

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

<label>
	{#if label}
		{label}
	{/if}
	<input {autocomplete} class={className} {id} {name} {placeholder} {type} />
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

		&::placeholder {
			color: var(--gray-400);
		}
	}
</style>
