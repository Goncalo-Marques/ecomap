<script lang="ts" generics="TValue">
	import type { ChangeEventHandler } from "svelte/elements";

	/**
	 * The label of the radio button.
	 * @default null
	 */
	export let label: string | null = null;

	/**
	 * Indicates that the radio button is the default selected one in the group.
	 * @default false
	 */
	export let checked: boolean = false;

	/**
	 * The ID of the radio button.
	 * @default null
	 */
	export let id: string | null = null;

	/**
	 * The name of the radio button.
	 * @default null
	 */
	export let name: string | null = null;

	/**
	 * The value of the radio button.
	 * @default null
	 */
	export let value: TValue | null = null;

	/**
	 * Callback fired when radio button value changes.
	 * @default null
	 */
	export let onChange: ChangeEventHandler<HTMLInputElement> | null = null;
</script>

<label class="flex gap-2 text-base font-normal">
	<input type="radio" {id} {name} {value} {checked} on:change={onChange} />
	{label}
</label>
