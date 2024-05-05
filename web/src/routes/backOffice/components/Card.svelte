<script lang="ts">
	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };
</script>

<main class={className}>
	<slot />
</main>

<style>
	main {
		padding: 2rem;
		background-color: var(--white);
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
	}
</style>
