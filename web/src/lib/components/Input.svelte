<script lang="ts">
	import type {
		ChangeEventHandler,
		HTMLInputAttributes,
		KeyboardEventHandler,
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
	 * Defines the maximum length that the user can enter into the input.
	 * @default null
	 */
	export let maxLength: HTMLInputAttributes["maxlength"] = null;

	/**
	 * Defines the minimum value that is acceptable and valid for the input.
	 * @default null
	 */
	export let min: HTMLInputAttributes["min"] = null;

	/**
	 * Defines the minimum length that the user can enter into the input.
	 * @default null
	 */
	export let minLength: HTMLInputAttributes["minlength"] = null;

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
	 * Callback fired when a key is pressed.
	 * @default null
	 */
	export let onKeyDown: KeyboardEventHandler<HTMLInputElement> | null = null;

	/**
	 * Specifies a regular expression the input value should match.
	 * @default null
	 */
	export let pattern: HTMLInputAttributes["pattern"] = null;

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
	 * Indicates that the user must specify a value for the input before the owning form can be submitted.
	 * @default false
	 */
	export let required: boolean = false;

	/**
	 * CSS styling declarations to be applied to the element.
	 * @default null
	 */
	export let style: string | null = null;

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

<div class="relative">
	<input
		{autocomplete}
		class="w-full rounded border border-gray-300 px-2 py-[0.375rem] text-gray-900 placeholder:text-gray-400 read-only:cursor-default data-[error=true]:border-red-500 data-[endIcon]:pr-8 focus:data-[error=true]:outline-red-500"
		{id}
		{name}
		{pattern}
		{placeholder}
		{type}
		{value}
		{readonly}
		{required}
		{max}
		{min}
		{style}
		data-error={error}
		maxlength={maxLength}
		minlength={minLength}
		on:keydown={onKeyDown}
		on:input={onInput}
		on:click={onClick}
		data-endIcon={endIcon}
	/>
	{#if endIcon}
		<div
			class="pointer-events-none absolute right-2 top-1/2 flex -translate-y-1/2 transform items-center text-gray-900"
		>
			<Icon name={endIcon} size="small" />
		</div>
	{/if}
</div>
