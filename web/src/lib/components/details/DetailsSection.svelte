<script lang="ts">
	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * The label of the section.
	 */
	export let label: string;
</script>

<section class={className}>
	<header>
		<h2>{label}</h2>
	</header>
	<slot />
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	h2 {
		font: var(--text-base-semibold);
	}
</style>
