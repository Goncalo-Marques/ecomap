<script lang="ts">
	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * The DOM element that's rendered.
	 * @default "main"
	 */
	export let element: string = "main";
</script>

<svelte:element this={element} class={`card ${className}`}>
	<slot />
</svelte:element>

<style>
	.card {
		padding: 2rem;
		background-color: var(--white);
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
	}
</style>
