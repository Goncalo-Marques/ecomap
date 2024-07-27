<script lang="ts">
	import type { EventHandler } from "svelte/elements";

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
	 * Global attribute valid for all elements, including all the input types, it defines a unique identifier (ID) which must be unique in the whole document.
	 * @default null
	 */
	export let id: string | null = null;

	/**
	 * The name of the form control. Submitted with the form as part of a name/value pair.
	 * @default null
	 */
	export let name: string | null = null;

	/**
	 * Callback fired when select changes its value.
	 * @default null;
	 */
	export let onChange: EventHandler<Event, HTMLSelectElement> | null = null;

	/**
	 * The text that appears in the form control when it has no value set.
	 * @default null
	 */
	export let placeholder: string | null = null;

	/**
	 * Indicates that the user must specify a value for the input before the owning form can be submitted.
	 * @default false
	 */
	export let required: boolean = false;

	/**
	 * The value of the select.
	 * @default ""
	 */
	export let value: string = "";
</script>

<select
	{id}
	{name}
	{value}
	{placeholder}
	{required}
	on:change={onChange}
	data-error={error}
	class={`appearance-none rounded border border-gray-300 bg-expand-more bg-[calc(100%-0.5rem)_center] bg-no-repeat px-2 py-[0.375rem] has-[option:disabled:checked[hidden]]:text-gray-400 data-[error=true]:border-red-500 focus:data-[error=true]:outline-red-500 ${className}`}
>
	{#if placeholder}
		<option value="" disabled selected hidden>
			{placeholder}
		</option>
	{/if}
	<slot />
</select>
