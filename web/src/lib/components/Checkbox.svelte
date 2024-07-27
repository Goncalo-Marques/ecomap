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

<label
	class="group relative flex items-center gap-2 data-[disabled=true]:text-gray-400"
	data-size={size}
	data-disabled={disabled}
>
	<input
		{disabled}
		{checked}
		type="checkbox"
		class={[
			"material-symbols-rounded appearance-none rounded border border-gray-300 leading-normal",
			"checked:relative checked:border-none checked:bg-green-700 checked:text-white before:checked:absolute before:checked:left-1/2 before:checked:top-1/2 before:checked:-translate-x-1/2 before:checked:-translate-y-1/2 before:checked:content-['check']",
			"group-data-[size=small]:size-3 group-data-[size=small]:text-[0.5rem]",
			"group-data-[size=medium]:size-4 group-data-[size=medium]:text-xs",
			"group-data-[size=large]:size-6 group-data-[size=large]:text-base",
			"group-data-[disabled=false]:cursor-pointer group-data-[disabled=true]:cursor-not-allowed group-data-[disabled=true]:bg-gray-300",
		].join(" ")}
		on:change={onChange}
	/>

	{#if label}
		{label}
	{/if}
</label>
