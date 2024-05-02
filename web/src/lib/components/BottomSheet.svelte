<script lang="ts">
	/**
	 * The title of the bottom sheet.
	 */
	export let title: string;
</script>

<div class="layout">
	<div class="sheet">
		<h2>{title}</h2>
		<div class="details">
			<slot />
		</div>
	</div>
</div>

<style>
	.layout {
		position: absolute;
		left: 0;
		bottom: 0;
		width: 100%;
		display: flex;
		padding-inline: 2.5rem;
	}

	.sheet {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 2rem;
		border-top-left-radius: 0.25rem;
		border-top-right-radius: 0.25rem;
		background-color: var(--white);
		box-shadow: var(--shadow-lg);

		& h2 {
			font: var(--text-2xl-semibold);
		}
	}

	.details {
		--column-count: 4;
		--column-gap: 0.5rem;
		--column-min-width: 13rem;

		display: grid;
		grid-template-columns: repeat(
			auto-fill,
			minmax(
				max(
					var(--column-min-width),
					(100% - var(--column-gap)) / var(--column-count)
				),
				1fr
			)
		);
		gap: var(--column-gap);
	}
</style>
