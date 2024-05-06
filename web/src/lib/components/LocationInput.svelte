<script lang="ts">
	import Input from "./Input.svelte";

	/**
	 * Indicates if the input contains an error.
	 * @default false
	 */
	export let error: boolean = false;

	/**
	 * The name of the form control. Submitted with the form as part of a name/value pair.
	 * @default null
	 */
	export let name: string | null = null;

	/**
	 * Callback fired when input element is clicked.
	 * @default null
	 */
	export let onClick: (() => void) | null = null;

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
	 * The value of the input.
	 * @default null
	 */
	export let value: string | null = null;

	/**
	 * Handle the keyboard event of the input by preventing typing a value inside the input.
	 * @param e Keyboard event.
	 */
	function handleKeyDown(e: KeyboardEvent) {
		switch (e.key) {
			// Ignore Tab key since it's used to interact with the HTML form element.
			case "Tab":
				break;

			case "Enter":
				onClick?.();
				e.preventDefault();
				break;

			default:
				e.preventDefault();
				break;
		}
	}
</script>

<Input
	{required}
	{name}
	{value}
	{placeholder}
	{error}
	{onClick}
	style="cursor: pointer; caret-color: transparent;"
	endIcon="location_on"
	onKeyDown={handleKeyDown}
/>
