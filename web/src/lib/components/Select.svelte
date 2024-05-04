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
	on:change={onChange}
	class={`${className} ${error ? "error" : ""}`}
>
	{#if placeholder}
		<option value="" disabled selected hidden>
			{placeholder}
		</option>
	{/if}
	<slot />
</select>

<style>
	select {
		appearance: none;
		border: 1px solid var(--gray-300);
		padding: 0.375rem 0.5rem;
		border-radius: 0.25rem;
		background-image: url("data:image/svg+xml;utf8,<svg width='16' height='16' viewBox='0 0 16 16' xmlns='http://www.w3.org/2000/svg'><path d='M7.99999 9.96664C7.9111 9.96664 7.82777 9.95275 7.74999 9.92498C7.67221 9.8972 7.59999 9.84998 7.53333 9.78331L4.46666 6.71664C4.34444 6.59442 4.28333 6.43886 4.28333 6.24998C4.28333 6.06109 4.34444 5.90553 4.46666 5.78331C4.58888 5.66109 4.74444 5.59998 4.93333 5.59998C5.12221 5.59998 5.27777 5.66109 5.39999 5.78331L7.99999 8.38331L10.6 5.78331C10.7222 5.66109 10.8778 5.59998 11.0667 5.59998C11.2555 5.59998 11.4111 5.66109 11.5333 5.78331C11.6555 5.90553 11.7167 6.06109 11.7167 6.24998C11.7167 6.43886 11.6555 6.59442 11.5333 6.71664L8.46666 9.78331C8.39999 9.84998 8.32777 9.8972 8.24999 9.92498C8.17221 9.95275 8.08888 9.96664 7.99999 9.96664Z' fill='%23111827'/></svg>");
		background-position: calc(100% - 0.5rem) center;
		background-repeat: no-repeat;

		&.error {
			border-color: var(--red-500);

			&:focus {
				outline-color: var(--red-500);
			}
		}

		&:has(option:disabled:checked[hidden]) {
			color: var(--gray-400);
		}
	}
</style>
