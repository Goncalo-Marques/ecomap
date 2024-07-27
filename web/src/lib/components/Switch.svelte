<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";

	/**
	 * Indicates if switch is checked.
	 */
	export let checked: boolean = false;

	/**
	 * Callback fired when clicking the component.
	 *
	 * @default null
	 */
	export let onClick: MouseEventHandler<HTMLInputElement> | null = null;
</script>

<label class="relative ml-auto inline-block h-5 w-10">
	<input
		class="peer opacity-0"
		type="checkbox"
		bind:checked
		on:click={onClick}
	/>
	<span
		class={[
			"absolute inset-0 cursor-pointer rounded-xl outline outline-1 outline-gray-400",
			"before:absolute before:bottom-1 before:size-3 before:translate-x-1 before:rounded-full before:bg-gray-400 before:transition before:content-['']",
			"peer-checked:bg-green-100 peer-checked:outline-0 peer-checked:before:translate-x-[1.375rem] peer-checked:before:bg-green-800",
		].join(" ")}
	/>
</label>
