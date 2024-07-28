<script lang="ts">
	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * The HTML element that's rendered.
	 * @default "main"
	 */
	export let element: string = "main";
</script>

<svelte:element
	this={element}
	class={`rounded border border-gray-300 bg-white p-8 ${className}`}
>
	<slot />
</svelte:element>
