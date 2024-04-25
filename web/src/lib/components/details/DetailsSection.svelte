<script lang="ts">
	/**
	 * The label of the section.
	 */
	export let label: string;
</script>

<section>
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
