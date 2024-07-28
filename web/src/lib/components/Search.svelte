<script lang="ts">
	import type {
		ChangeEventHandler,
		HTMLInputAttributes,
	} from "svelte/elements";

	import Icon from "./Icon.svelte";

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
	 * The value of the search input.
	 * @default null
	 */
	export let value: string | null = null;
</script>

<label class={`group relative block ${className}`}>
	<div
		class="pointer-events-none absolute left-2 top-1/2 flex -translate-y-1/2 items-center justify-center text-gray-400 group-focus-within:hidden"
	>
		<Icon name="search" size="small" />
	</div>
	<input
		{autocomplete}
		{placeholder}
		type="search"
		{value}
		on:input={onInput}
		class="relative w-full rounded border border-gray-300 py-[0.375rem] pl-8 pr-2 text-gray-900 outline-offset-0 placeholder:text-gray-400 group-focus-within:pl-2 [&::-webkit-search-cancel-button]:relative [&::-webkit-search-cancel-button]:m-0 [&::-webkit-search-cancel-button]:size-5 [&::-webkit-search-cancel-button]:cursor-pointer [&::-webkit-search-cancel-button]:appearance-none [&::-webkit-search-cancel-button]:bg-close [&::-webkit-search-cancel-button]:bg-contain"
	/>
</label>
